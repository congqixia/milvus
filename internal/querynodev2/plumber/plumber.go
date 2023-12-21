package plumber

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/golang/protobuf/proto"
	"github.com/milvus-io/milvus-proto/go-api/v2/schemapb"
	"github.com/milvus-io/milvus/internal/proto/internalpb"
	"github.com/milvus-io/milvus/internal/proto/querypb"
	"github.com/milvus-io/milvus/pkg/common"
	"github.com/milvus-io/milvus/pkg/log"
	"github.com/milvus-io/milvus/pkg/util/conc"
	"github.com/milvus-io/milvus/pkg/util/merr"
	"github.com/milvus-io/milvus/pkg/util/typeutil"
	"github.com/samber/lo"
)

var (
	Enabled          = true
	PresetTopKs      = []int64{1, 10, 100, 1000, 10000}
	cachedSearchData = typeutil.NewConcurrentMap[string, []byte]()
	collSearchMeta   = typeutil.NewConcurrentMap[int64, *SearchMeta]()
	sf               = &conc.Singleflight[[]byte]{}
)

type SearchMeta struct {
	PKField *schemapb.FieldSchema

	MockData *typeutil.ConcurrentMap[int64, *schemapb.FieldData]
}

func Register(collID int64, coll *schemapb.CollectionSchema) {
	if coll == nil {
		log.Warn("[Plumber] collection schema is nil")
		return
	}

	pkField := lo.FindOrElse(coll.GetFields(), nil, func(field *schemapb.FieldSchema) bool {
		return field.IsPrimaryKey
	})

	if pkField == nil {
		log.Warn("[Plumber] cannot find pk field")
		return
	}

	mockData := typeutil.NewConcurrentMap[int64, *schemapb.FieldData]()

	for _, field := range coll.GetFields() {
		if field.GetFieldID() < common.StartOfUserFieldID {
			// ignore system fields
			continue
		}

		mockData.Insert(field.GetFieldID(), generateMockData(field))
	}

	collSearchMeta.Insert(collID, &SearchMeta{
		MockData: mockData,
		PKField:  pkField,
	})
}

func generateMockData(field *schemapb.FieldSchema) *schemapb.FieldData {
	fieldData := &schemapb.FieldData{
		FieldId:   field.GetFieldID(),
		Type:      field.GetDataType(),
		FieldName: field.GetName(),
		IsDynamic: field.GetIsDynamic(),
	}
	switch field.GetDataType() {
	case schemapb.DataType_Int64:
		data := lo.RepeatBy(10000, func(idx int) int64 { return int64(idx) })
		field := &schemapb.FieldData_Scalars{
			Scalars: &schemapb.ScalarField{
				Data: &schemapb.ScalarField_LongData{
					LongData: &schemapb.LongArray{
						Data: data,
					},
				},
			},
		}
		fieldData.Field = field
	default:
		return nil
	}

	return fieldData
}

func PlumberSearch(ctx context.Context, req *querypb.SearchRequest) (*internalpb.SearchResults, error) {
	collectionID := req.GetReq().GetCollectionID()
	topk := req.Req.GetTopk()
	nq := req.Req.GetNq()
	metricType := req.GetReq().GetMetricType()

	outputFieldIDs := req.GetReq().GetOutputFieldsId()

	sort.Slice(outputFieldIDs, func(i, j int) bool { return outputFieldIDs[i] < outputFieldIDs[j] })

	tag := fmt.Sprintf("%d-%d-%d-%v", collectionID, nq, topk, outputFieldIDs)

	bs, ok := cachedSearchData.Get(tag)
	if !ok {
		var err error
		bs, err = getSearchResult(tag, collectionID, topk, nq, outputFieldIDs)
		if err != nil {
			return nil, err
		}
	}

	searchResults := &internalpb.SearchResults{
		Status:     merr.Success(),
		NumQueries: nq,
		TopK:       topk,
		MetricType: metricType,
		SlicedBlob: bs,
	}
	return searchResults, nil
}

func getSearchResult(tag string, collectionID, topk, nq int64, outputFieldIDs []int64) ([]byte, error) {
	bs, err, _ := sf.Do(tag, func() ([]byte, error) {
		searchMeta, ok := collSearchMeta.Get(collectionID)
		if !ok {
			return nil, errors.New("unexpected collection")
		}
		resultData := &schemapb.SearchResultData{
			NumQueries: nq,
			TopK:       topk,
			Scores:     lo.RepeatBy(int(topk*nq), func(_ int) float32 { return 0 }),
			Topks:      lo.RepeatBy(int(nq), func(_ int) int64 { return topk }),
		}

		data, ok := searchMeta.MockData.Get(searchMeta.PKField.FieldID)
		if !ok {
			return nil, errors.New("plumber not supported field")
		}
		resultData.Ids = &schemapb.IDs{
			IdField: &schemapb.IDs_IntId{
				IntId: &schemapb.LongArray{
					Data: data.GetScalars().GetLongData().GetData()[:int(topk*nq)],
				},
			},
		}

		for _, field := range outputFieldIDs {

			data, ok := searchMeta.MockData.Get(field)
			if !ok {
				return nil, errors.New("plumber not supported field")
			}
			resultData.FieldsData = append(resultData.FieldsData, sliceFieldData(data, int(topk)))
		}
		slicedBlob, err := proto.Marshal(resultData)
		if err != nil {
			return nil, err
		}
		return slicedBlob, nil
	})

	if err == nil {
		cachedSearchData.Insert(tag, bs)
	}
	return bs, err
}

func sliceFieldData(fd *schemapb.FieldData, length int) *schemapb.FieldData {
	fieldData := &schemapb.FieldData{
		FieldId:   fd.GetFieldId(),
		Type:      fd.GetType(),
		FieldName: fd.GetFieldName(),
		IsDynamic: fd.GetIsDynamic(),
	}
	switch fd.GetType() {
	case schemapb.DataType_Int64:
		field := &schemapb.FieldData_Scalars{
			Scalars: &schemapb.ScalarField{
				Data: &schemapb.ScalarField_LongData{
					LongData: &schemapb.LongArray{
						Data: fd.GetScalars().GetLongData().GetData()[:length],
					},
				},
			},
		}
		fieldData.Field = field
		return fieldData
	default:
		return nil
	}
}
