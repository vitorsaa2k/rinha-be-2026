package utils

import (
	"container/heap"
	"slices"
)

type DatasetHeapQuantized []CalculatedDistanceQuantized

func (h DatasetHeapQuantized) Len() int {
	return len(h)
}

func (h DatasetHeapQuantized) Less(i, j int) bool {
	return h[i].Distance > h[j].Distance
}

func (h DatasetHeapQuantized) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]

}
func (h *DatasetHeapQuantized) Push(x any) {
	*h = append(*h, x.(CalculatedDistanceQuantized))
}

func (h *DatasetHeapQuantized) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type BoundedCollectionQuantized struct {
	heap     *DatasetHeapQuantized
	capacity int
}

func NewBoundedCollectionQuantized(capacity int) *BoundedCollectionQuantized {
	h := &DatasetHeapQuantized{}
	heap.Init(h)
	return &BoundedCollectionQuantized{heap: h, capacity: capacity}
}

func (bc *BoundedCollectionQuantized) Add(val CalculatedDistanceQuantized) {
	if bc.heap.Len() < bc.capacity {
		heap.Push(bc.heap, val)
		return
	}

	if val.Distance < (*bc.heap)[0].Distance {
		(*bc.heap)[0] = val
		heap.Fix(bc.heap, 0)
	}
}

func (bc *BoundedCollectionQuantized) Reset() {
	*bc.heap = (*bc.heap)[:0]
}

func (bc *BoundedCollectionQuantized) Sorted() []CalculatedDistanceQuantized {
	result := append([]CalculatedDistanceQuantized(nil), (*bc.heap)...)

	slices.SortFunc(result, func(a, b CalculatedDistanceQuantized) int {
		switch {
		case a.Distance > b.Distance:
			return -1
		case a.Distance < b.Distance:
			return 1
		default:
			return 0
		}
	})

	return result
}

type DatasetHeap []CalculatedDistance

func (h DatasetHeap) Len() int {
	return len(h)
}

func (h DatasetHeap) Less(i, j int) bool {
	return h[i].Distance > h[j].Distance
}

func (h DatasetHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]

}
func (h *DatasetHeap) Push(x any) {
	*h = append(*h, x.(CalculatedDistance))
}

func (h *DatasetHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type BoundedCollection struct {
	heap     *DatasetHeap
	capacity int
}

func NewBoundedCollection(capacity int) *BoundedCollection {
	h := &DatasetHeap{}
	heap.Init(h)
	return &BoundedCollection{heap: h, capacity: capacity}
}

func (bc *BoundedCollection) Add(val CalculatedDistance) {
	if bc.heap.Len() < bc.capacity {
		heap.Push(bc.heap, val)
		return
	}

	if val.Distance < (*bc.heap)[0].Distance {
		(*bc.heap)[0] = val
		heap.Fix(bc.heap, 0)
	}
}

func (bc *BoundedCollection) Reset() {
	*bc.heap = (*bc.heap)[:0]
}

func (bc *BoundedCollection) Sorted() []CalculatedDistance {
	result := append([]CalculatedDistance(nil), (*bc.heap)...)

	slices.SortFunc(result, func(a, b CalculatedDistance) int {
		switch {
		case a.Distance > b.Distance:
			return -1
		case a.Distance < b.Distance:
			return 1
		default:
			return 0
		}
	})

	return result
}
