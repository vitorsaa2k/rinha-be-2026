package utils

import (
	"gin-test/internal/ivf"
	"sync"
)

const THRESHOLD = 0.6
const HEAP_SIZE = 10

type SearchResultStruct struct {
	IsPossibleFraud bool
	Score           float64
}

type CalculatedDistance struct {
	Distance float64
	Label    string
}

var boundedPool = sync.Pool{
	New: func() any {
		return NewBoundedCollection(HEAP_SIZE)
	},
}

func SearchInVector(vec []float64, totalDimensions int8, closestCentroids []ivf.ClosestCentroids) (SearchResultStruct, error) {
	boundedClosest := boundedPool.Get().(*BoundedCollection)
	defer boundedPool.Put(boundedClosest)
	boundedClosest.Reset()
	for _, c := range closestCentroids {
		for _, v := range c.Cluster.Lists {
			var totalSum float64
			for j := 0; j < int(totalDimensions); j++ {
				difference := float64(v.Vector[j]) - vec[j]
				totalSum += difference * difference
			}
			boundedClosest.Add(CalculatedDistance{Distance: totalSum, Label: v.Label})
		}
	}
	totalFraudsNeighbours := 0.0
	for _, value := range *boundedClosest.heap {
		if value.Label == "fraud" {
			totalFraudsNeighbours++
		}
	}
	//fmt.Println("Total fraud neighbours(out of 1000):", totalFraudsNeighbours)
	score := 0.0
	score = totalFraudsNeighbours / HEAP_SIZE

	//fmt.Println("Score:", score)
	if score < THRESHOLD {
		return SearchResultStruct{IsPossibleFraud: false, Score: score}, nil
	} else {
		return SearchResultStruct{IsPossibleFraud: true, Score: score}, nil
	}
}
