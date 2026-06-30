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

type CalculatedDistanceQuantized struct {
	Distance uint64
	Label    uint8
}

type CalculatedDistance struct {
	Distance float32
	Label    string
}

var boundedPool = sync.Pool{
	New: func() any {
		return NewBoundedCollectionQuantized(HEAP_SIZE)
	},
}

var boundedPoolQuantized = sync.Pool{
	New: func() any {
		return NewBoundedCollection(HEAP_SIZE)
	},
}

func SearchInVectorQuantized(vec []int16, totalDimensions int8, closestCentroids []ivf.QuantizedClosestCentroids) (SearchResultStruct, error) {
	boundedClosest := boundedPool.Get().(*BoundedCollectionQuantized)
	defer boundedPool.Put(boundedClosest)
	boundedClosest.Reset()
	for _, c := range closestCentroids {
		for _, v := range c.Cluster.Lists {
			var totalSum uint64
			for j := 0; j < int(totalDimensions); j++ {
				diff := int64(v.Vector[j]) - int64(vec[j])
				totalSum += uint64(diff * diff)
			}
			boundedClosest.Add(CalculatedDistanceQuantized{Distance: totalSum, Label: v.Label})
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
	//fmt.Println("Total fraud neighbours(out of 1000):", totalFraudsNeighbours)
	score := 0.0
	score = totalFraudCount / HEAP_SIZE
	//fmt.Println("Score:", score)
	if score < THRESHOLD {
		return SearchResultStruct{IsPossibleFraud: false, Score: score}, nil
	} else {
		return SearchResultStruct{IsPossibleFraud: true, Score: score}, nil
	}
}
