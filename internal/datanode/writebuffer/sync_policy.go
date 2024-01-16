package writebuffer

import (
	"container/heap"
	"math/rand"
	"time"

	"github.com/samber/lo"
	"go.uber.org/atomic"

	"github.com/milvus-io/milvus-proto/go-api/v2/commonpb"
	"github.com/milvus-io/milvus/internal/datanode/metacache"
	"github.com/milvus-io/milvus/internal/proto/datapb"
	"github.com/milvus-io/milvus/pkg/util/paramtable"
	"github.com/milvus-io/milvus/pkg/util/tsoutil"
	"github.com/milvus-io/milvus/pkg/util/typeutil"
)

type SyncPolicy interface {
	SelectSegments(buffers []*segmentBuffer, ts typeutil.Timestamp) []int64
	Reason() string
}

type SelectSegmentFunc func(buffer []*segmentBuffer, ts typeutil.Timestamp) []int64

type SelectSegmentFnPolicy struct {
	fn     SelectSegmentFunc
	reason string
}

func (f SelectSegmentFnPolicy) SelectSegments(buffers []*segmentBuffer, ts typeutil.Timestamp) []int64 {
	return f.fn(buffers, ts)
}

func (f SelectSegmentFnPolicy) Reason() string { return f.reason }

func wrapSelectSegmentFuncPolicy(fn SelectSegmentFunc, reason string) SelectSegmentFnPolicy {
	return SelectSegmentFnPolicy{
		fn:     fn,
		reason: reason,
	}
}

func GetFullBufferPolicy() SyncPolicy {
	return wrapSelectSegmentFuncPolicy(
		func(buffers []*segmentBuffer, _ typeutil.Timestamp) []int64 {
			return lo.FilterMap(buffers, func(buf *segmentBuffer, _ int) (int64, bool) {
				return buf.segmentID, buf.IsFull()
			})
		}, "buffer full")
}

func GetCompactedSegmentsPolicy(meta metacache.MetaCache) SyncPolicy {
	return wrapSelectSegmentFuncPolicy(func(buffers []*segmentBuffer, _ typeutil.Timestamp) []int64 {
		segmentIDs := lo.Map(buffers, func(buffer *segmentBuffer, _ int) int64 { return buffer.segmentID })
		return meta.GetSegmentIDsBy(metacache.WithSegmentIDs(segmentIDs...), metacache.WithCompacted())
	}, "segment compacted")
}

func GetSyncStaleBufferPolicy(staleDuration time.Duration) SyncPolicy {
	return wrapSelectSegmentFuncPolicy(func(buffers []*segmentBuffer, ts typeutil.Timestamp) []int64 {
		current := tsoutil.PhysicalTime(ts)
		return lo.FilterMap(buffers, func(buf *segmentBuffer, _ int) (int64, bool) {
			minTs := buf.MinTimestamp()
			start := tsoutil.PhysicalTime(minTs)
			jitter := time.Duration(rand.Float64() * 0.1 * float64(staleDuration))
			return buf.segmentID, current.Sub(start) > staleDuration+jitter
		})
	}, "buffer stale")
}

func GetSealedSegmentsPolicy(meta metacache.MetaCache) SyncPolicy {
	return wrapSelectSegmentFuncPolicy(func(_ []*segmentBuffer, _ typeutil.Timestamp) []int64 {
		ids := meta.GetSegmentIDsBy(metacache.WithSegmentState(commonpb.SegmentState_Sealed))
		meta.UpdateSegments(metacache.UpdateState(commonpb.SegmentState_Flushing),
			metacache.WithSegmentIDs(ids...), metacache.WithSegmentState(commonpb.SegmentState_Sealed))
		return ids
	}, "segment flushing")
}

func GetFlushTsPolicy(flushTimestamp *atomic.Uint64, meta metacache.MetaCache) SyncPolicy {
	return wrapSelectSegmentFuncPolicy(func(buffers []*segmentBuffer, ts typeutil.Timestamp) []int64 {
		flushTs := flushTimestamp.Load()
		if flushTs != nonFlushTS && ts >= flushTs {
			// flush segment start pos < flushTs && checkpoint > flushTs
			ids := lo.FilterMap(buffers, func(buf *segmentBuffer, _ int) (int64, bool) {
				seg, ok := meta.GetSegmentByID(buf.segmentID)
				if !ok {
					return buf.segmentID, false
				}
				inRange := seg.State() == commonpb.SegmentState_Flushed ||
					seg.Level() == datapb.SegmentLevel_L0
				return buf.segmentID, inRange && buf.MinTimestamp() < flushTs
			})

			// flush all buffer
			return ids
		}
		return nil
	}, "flush ts")
}

func GetMemoryHighPolicy(memoryHigh *atomic.Bool) SyncPolicy {
	return wrapSelectSegmentFuncPolicy(func(buffers []*segmentBuffer, ts typeutil.Timestamp) []int64 {
		if !memoryHigh.Load() {
			return nil
		}

		h := &SegMemSizeHeap{}
		heap.Init(h)

		maxNum := paramtable.Get().DataNodeCfg.MemoryForceSyncSegmentNum.GetAsInt()
		minNum := paramtable.Get().DataNodeCfg.MemoryForceSyncSegmentMinNum.GetAsInt()
		minSize := paramtable.Get().DataNodeCfg.MemoryForceSyncMinSize.GetAsInt64()
		for _, buf := range buffers {
			if h.Len() >= minNum && buf.MemorySize() < minSize {
				continue
			}
			heap.Push(h, buf)
			if h.Len() > maxNum {
				heap.Pop(h)
			}
		}

		return lo.Map(*h, func(buf *segmentBuffer, _ int) int64 { return buf.segmentID })
	}, "memory high")
}

// SegMemSizeHeap implement min-heap for sorting.
type SegMemSizeHeap []*segmentBuffer

func (h SegMemSizeHeap) Len() int { return len(h) }
func (h SegMemSizeHeap) Less(i, j int) bool {
	return h[i].MemorySize() < h[j].MemorySize()
}
func (h SegMemSizeHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *SegMemSizeHeap) Push(x any) {
	*h = append(*h, x.(*segmentBuffer))
}

func (h *SegMemSizeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
