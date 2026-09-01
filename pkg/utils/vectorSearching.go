package utils

import (
	"container/heap"
	"gin-test/internal/ivf"
	"sync"
)

const (
	THRESHOLD = 0.6
	HEAP_SIZE = 5
	dims      = 14
)

type SearchResultStruct struct {
	IsPossibleFraud bool
	Score           float64
}

type DistanceLabel struct {
	Distance uint64
	Label    uint8
}

type CalculatedDistance struct {
	Distance float32
	Label    string
}

// KNNHeap is a max-heap that keeps the K nearest neighbors.
type KNNHeap struct {
	items []DistanceLabel
}

func (h *KNNHeap) Reset()             { h.items = h.items[:0] }
func (h *KNNHeap) Len() int           { return len(h.items) }
func (h *KNNHeap) Less(i, j int) bool { return h.items[i].Distance > h.items[j].Distance }
func (h *KNNHeap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *KNNHeap) Push(x any)         { h.items = append(h.items, x.(DistanceLabel)) }
func (h *KNNHeap) Pop() any {
	old := h.items
	n := len(old)
	x := old[n-1]
	h.items = old[:n-1]
	return x
}

func (h *KNNHeap) Add(val DistanceLabel) {
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

var knnHeapPool = sync.Pool{
	New: func() any {
		return &KNNHeap{items: make([]DistanceLabel, 0, HEAP_SIZE)}
	},
}

var boundedPoolQuantized = sync.Pool{
	New: func() any {
		return NewBoundedCollection(HEAP_SIZE)
	},
}

func SearchInVectorQuantized(vec *[14]int16, matches []ivf.CentroidMatch) SearchResultStruct {
	h := knnHeapPool.Get().(*KNNHeap)
	defer knnHeapPool.Put(h)
	h.Reset()

	for _, m := range matches {
		p := ivf.GlobalPartitions[m.PartitionIdx]
		start := p.Offsets[m.LocalCID]
		end := p.Offsets[m.LocalCID+1]

		for i := start; i < end; i++ {
			base := int(i) * dims
			var sum uint64
			v := p.Vectors[base : base+dims]
			for j := range dims {
				diff := int64(v[j]) - int64(vec[j])
				sum += uint64(diff * diff)
			}
			h.Add(DistanceLabel{Distance: sum, Label: p.Labels[i]})
		}
	}

	fraudCount := 0
	for _, v := range h.items {
		fraudCount += int(v.Label)
	}
	score := float64(fraudCount) / float64(len(h.items))

	return SearchResultStruct{
		IsPossibleFraud: score >= THRESHOLD,
		Score:           score,
	}
}

func SearchInVector(vec []float32, totalDimensions int8, closestCentroids []ivf.ClosestCentroids) (SearchResultStruct, error) {
	boundedClosest := boundedPoolQuantized.Get().(*BoundedCollection)
	defer boundedPoolQuantized.Put(boundedClosest)
	boundedClosest.Reset()
	for _, c := range closestCentroids {
		for _, v := range c.Cluster.Lists {
			var totalSum float32
			for j := 0; j < int(totalDimensions); j++ {
				difference := float32(v.Vector[j]) - vec[j]
				totalSum += difference * difference
			}
			boundedClosest.Add(CalculatedDistance{Distance: totalSum, Label: v.Label})
		}
	}
	totalFraudCount := 0.0
	for _, value := range *boundedClosest.heap {
		if value.Label == "fraud" {
			totalFraudCount++
		}
	}
	score := 0.0
	score = totalFraudCount / HEAP_SIZE
	if score < THRESHOLD {
		return SearchResultStruct{IsPossibleFraud: false, Score: score}, nil
	} else {
		return SearchResultStruct{IsPossibleFraud: true, Score: score}, nil
	}
}
