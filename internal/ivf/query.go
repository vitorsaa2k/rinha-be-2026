package ivf

import (
	"fmt"
	"runtime/debug"
	"sync"
)

const HEAP_SIZE = 25

type ClosestCentroids struct {
	Cluster  Cluster
	Distance float64
}

type QuantizedClosestCentroids struct {
	Cluster  QuantizedCluster
	Distance uint64
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
	for _, c := range *boundedClosest.heap {
		fmt.Println(c.Cluster.Centroid)
	}
	return *boundedClosest.heap
}

func SearchIVFQuantized(query []int16, k uint8) []QuantizedClosestCentroids {
	boundedClosest := quantizedBoundedPool.Get().(*QuantizedBoundedCollection)
	defer quantizedBoundedPool.Put(boundedClosest)
	boundedClosest.Reset()
	for _, cluster := range QuantizedClusters {
		distance := QuantizedDistance(cluster.Centroid, query)
		boundedClosest.Add(QuantizedClosestCentroids{Cluster: cluster, Distance: distance})
	}
	return *boundedClosest.heap
}

func QuantizeQuery(normalized []float32) []int16 {
	out := make([]int16, dims)
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
	QuantizedClusters = idx.ToQuantizedClusters()
	fmt.Printf("Loaded quantized index: K=%d N=%d clusters=%d\n", idx.K, idx.N, len(QuantizedClusters))
	idx = nil
	debug.FreeOSMemory()
	return idx, nil
}
