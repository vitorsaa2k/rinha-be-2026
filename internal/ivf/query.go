package ivf

import (
	"sync"
)

const HEAP_SIZE = 5

type ClosestCentroids struct {
	Cluster  Cluster
	Distance float64
}

var boundedPool = sync.Pool{
	New: func() any {
		return NewBoundedCollection(HEAP_SIZE)
	},
}

func SearchIVF(normalizedVector []float64, k uint8) []ClosestCentroids {
	boundedClosest := boundedPool.Get().(*BoundedCollection)
	defer boundedPool.Put(boundedClosest)
	boundedClosest.Reset()
	parsedVector := make([]float32, len(normalizedVector))
	for i, v := range normalizedVector {
		parsedVector[i] = float32(v)
	}
	for _, cluster := range Clusters {
		distance := Distance(cluster.Centroid, parsedVector)
		boundedClosest.Add(ClosestCentroids{Cluster: cluster, Distance: distance})
	}
	return *boundedClosest.heap
}
