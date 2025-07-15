package segcore_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/milvus-io/milvus-proto/go-api/v2/schemapb"
	"github.com/milvus-io/milvus/internal/mocks/util/mock_segcore"
	"github.com/milvus-io/milvus/internal/util/segcore"
	"github.com/milvus-io/milvus/pkg/v2/util/paramtable"
	"github.com/milvus-io/milvus/pkg/v2/util/tsoutil"
)

func TestCollection(t *testing.T) {
	paramtable.Init()
	schema := mock_segcore.GenTestCollectionSchema("test", schemapb.DataType_Int64, false)
	indexMeta := mock_segcore.GenTestIndexMeta(1, schema)
	ccollection, err := segcore.CreateCCollection(&segcore.CreateCCollectionRequest{
		CollectionID: 1,
		Schema:       schema,
		IndexMeta:    indexMeta,
	})
	assert.NoError(t, err)
	assert.NotNil(t, ccollection)
	assert.NotNil(t, ccollection.Schema())
	assert.NotNil(t, ccollection.IndexMeta())
	assert.Equal(t, int64(1), ccollection.ID())
	defer ccollection.Release()
}

func TestCollectionUpdateSchema(t *testing.T) {
	paramtable.Init()
	schema := mock_segcore.GenTestCollectionSchema("test", schemapb.DataType_Int64, false)
	indexMeta := mock_segcore.GenTestIndexMeta(1, schema)
	ccollection, err := segcore.CreateCCollection(&segcore.CreateCCollectionRequest{
		CollectionID: 1,
		Schema:       schema,
		IndexMeta:    indexMeta,
	})
	assert.NoError(t, err)
	assert.NotNil(t, ccollection)
	assert.NotNil(t, ccollection.Schema())
	assert.NotNil(t, ccollection.IndexMeta())
	assert.Equal(t, int64(1), ccollection.ID())
	defer ccollection.Release()

	testUpdateSchema(ccollection, schema)
	// ccollection.UpdateSchema(schema, tsoutil.ComposeTSByTime(time.Now(), 0))
	printSchema(ccollection)
}

func testUpdateSchema(ccol *segcore.CCollection, schema *schemapb.CollectionSchema) {
	ccol.UpdateSchema(schema, tsoutil.ComposeTSByTime(time.Now(), 0))
}

func printSchema(ccol *segcore.CCollection) {
	ccol.PrintSchema()
}
