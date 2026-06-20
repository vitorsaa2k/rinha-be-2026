package ivf

import (
	"fmt"
	"sync"
)

const HEAP_SIZE = 10

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
	return *boundedClosest.heap
}

// LoadIndexFromFile loads a serialized index from path and populates
// the global Clusters variable. Returns the loaded index for inspection.
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
