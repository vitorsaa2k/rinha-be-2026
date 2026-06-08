package utils

import (
	"errors"
	"gin-test/models"
	"math"
	"sync"
)

const THRESHOLD = 0.6

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
		return NewBoundedCollection(100)
	},
}

func SearchInVector(vec []float64, totalDimensions int8, dataset []models.DatasetStruct) (SearchResultStruct, error) {
	boundedClosest := boundedPool.Get().(*BoundedCollection)
	defer boundedPool.Put(boundedClosest)
	boundedClosest.Reset()
	for _, v := range dataset {
		var totalSum float64
		for j := 0; j < int(totalDimensions); j++ {
			difference := float64(v.Vector[j]) - vec[j]
			totalSum += difference * difference
		}
		finalValue := math.Sqrt(totalSum)
		boundedClosest.Add(CalculatedDistance{Distance: finalValue, Label: v.Label})
	}
	totalFraudsNeighbours := 0
	for _, value := range *boundedClosest.heap {
		if value.Label == "fraud" {
			totalFraudsNeighbours++
		}
	}
	//fmt.Println("Total fraud neighbours(out of 1000):", totalFraudsNeighbours)
	score := 0
	if boundedClosest.heap.Len() < 1 {
		score = 1
	} else {
		score = totalFraudsNeighbours / boundedClosest.heap.Len()
	}
	//fmt.Println("Score:", score)
	if float64(score) < THRESHOLD {
		return SearchResultStruct{IsPossibleFraud: false, Score: float64(score)}, nil
	} else if float64(score) > THRESHOLD {
		return SearchResultStruct{IsPossibleFraud: true, Score: float64(score)}, nil
	}
	return SearchResultStruct{IsPossibleFraud: true, Score: float64(score)}, errors.New("Error when searching in vector")
}
