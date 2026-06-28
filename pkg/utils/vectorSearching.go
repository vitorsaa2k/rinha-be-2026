package utils

import (
	"gin-test/internal/ivf"
	"sync"
)

const THRESHOLD = 0.6
const HEAP_SIZE = 5

type SearchResultStruct struct {
	IsPossibleFraud bool
	Score           float64
}

type CalculatedDistance struct {
	Distance uint64
	Label    uint8
}

var boundedPool = sync.Pool{
	New: func() any {
		return NewBoundedCollection(HEAP_SIZE)
	},
}

func SearchInVector(vec []int16, totalDimensions int8, closestCentroids []ivf.QuantizedClosestCentroids) (SearchResultStruct, error) {
	boundedClosest := boundedPool.Get().(*BoundedCollection)
	defer boundedPool.Put(boundedClosest)
	boundedClosest.Reset()
	for _, c := range closestCentroids {
		for _, v := range c.Cluster.Lists {
			var totalSum uint64
			for j := 0; j < int(totalDimensions); j++ {
				diff := int64(v.Vector[j]) - int64(vec[j])
				totalSum += uint64(diff * diff)
			}
			boundedClosest.Add(CalculatedDistance{Distance: totalSum, Label: v.Label})
		}
	}
	totalFraudCount := 0.0
	for _, value := range *boundedClosest.heap {
		if value.Label == 1 {
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
