package ivf

import (
	"container/heap"
	"fmt"
	"sync"
)

const HEAP_SIZE = 25

type ClosestCentroids struct {
	Cluster  Cluster
	Distance float64
}

type CentroidMatch struct {
	PartitionIdx int32
	LocalCID     uint32
	Distance     uint64
}

// CentroidMatchHeap is a max-heap that keeps the K closest centroids.
type CentroidMatchHeap struct {
	items []CentroidMatch
}

func (h *CentroidMatchHeap) Reset()             { h.items = h.items[:0] }
func (h *CentroidMatchHeap) Len() int           { return len(h.items) }
func (h *CentroidMatchHeap) Less(i, j int) bool { return h.items[i].Distance > h.items[j].Distance }
func (h *CentroidMatchHeap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *CentroidMatchHeap) Push(x any)         { h.items = append(h.items, x.(CentroidMatch)) }
func (h *CentroidMatchHeap) Pop() any {
	old := h.items
	n := len(old)
	x := old[n-1]
	h.items = old[:n-1]
	return x
}

func (h *CentroidMatchHeap) Add(val CentroidMatch) {
	if len(h.items) < cap(h.items) {
		h.items = append(h.items, val)
		if len(h.items) == cap(h.items) {
			heap.Init(h)
		}
		return
	}
	if val.Distance < h.items[0].Distance {
		h.items[0] = val
		heap.Fix(h, 0)
	}
}

var centroidHeapPool = sync.Pool{
	New: func() any {
		return &CentroidMatchHeap{items: make([]CentroidMatch, 0, HEAP_SIZE)}
	},
}

var boundedPool = sync.Pool{
	New: func() any {
		return NewBoundedCollection(HEAP_SIZE)
	},
}

func SearchIVF(normalizedVector []float32, k uint8) []ClosestCentroids {
	boundedClosest := boundedPool.Get().(*BoundedCollection)
	defer boundedPool.Put(boundedClosest)
	boundedClosest.Reset()
	for _, cluster := range GlobalClusters {
		distance := Distance(cluster.Centroid, normalizedVector)
		boundedClosest.Add(ClosestCentroids{Cluster: cluster, Distance: distance})
	}
	return *boundedClosest.heap
}

func SearchIVFQuantized(query *[14]int16, k uint8) []CentroidMatch {
	h := centroidHeapPool.Get().(*CentroidMatchHeap)
	defer centroidHeapPool.Put(h)
	h.Reset()

	for pi, p := range GlobalPartitions {
		for c := uint32(0); c < p.K; c++ {
			base := int(c) * dims
			dist := QuantizedDistanceArr(query, p.Centroids[base:base+dims])
			h.Add(CentroidMatch{
				PartitionIdx: int32(pi),
				LocalCID:     c,
				Distance:     dist,
			})
		}
	}

	n := int(k)
	if h.Len() < n {
		n = h.Len()
	}
	out := make([]CentroidMatch, n)
	for i := n - 1; i >= 0; i-- {
		out[i] = heap.Pop(h).(CentroidMatch)
	}
	return out
}

func QuantizeQuery(normalized [14]float32) [14]int16 {
	var out [14]int16
	for i, v := range normalized {
		out[i] = QuantizeFloat32(v)
	}
	return out
}

func LoadIndexFromFile(path string) (*IVFIndex, error) {
	fmt.Println("Loading IVF index from:", path)
	idx, err := LoadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("load index: %w", err)
	}
	Clusters = idx.ToClusters()
	fmt.Printf("Loaded index: K=%d N=%d clusters=%d\n", idx.K, idx.N, len(Clusters))
	return idx, nil
}

func LoadQuantizedIndexFromFile(path string) (*IVFQuantizedIndex, error) {
	fmt.Println("Loading quantized IVF index from:", path)
	idx, err := LoadQuantizedFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("load quantized index: %w", err)
	}
	fmt.Printf("Loaded quantized index: K=%d N=%d\n", idx.K, idx.N)
	return idx, nil
}
