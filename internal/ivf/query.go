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

func SearchIVF(normalizedVector []float32, k uint8) []ClosestCentroids {
	boundedClosest := boundedPool.Get().(*BoundedCollection)
	defer boundedPool.Put(boundedClosest)
	boundedClosest.Reset()
	for _, cluster := range Clusters {
		distance := Distance(cluster.Centroid, normalizedVector)
		boundedClosest.Add(ClosestCentroids{Cluster: cluster, Distance: distance})
	}
	/* for _, c := range *boundedClosest.heap {
		fmt.Println(c)
	} */
	return []ClosestCentroids{(*boundedClosest.heap)[HEAP_SIZE-1]}
}
