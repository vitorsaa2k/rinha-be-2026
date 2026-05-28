package utils

import (
	"errors"
	"fmt"
	"gin-test/models"
	"math"
)

const TRESHOLD = 0.6

type SearchResultStruct struct {
	IsPossibleFraud bool
	Score           float64
}

type CalculatedDistance struct {
	Distance float64
	Label    string
}

func SearchInVector(vec []float64, totalDimensions int8, dataset []models.DatasetStruct) (SearchResultStruct, error) {
	boundedClosest := NewBoundedCollection(1000)
	var lowest float64
	for i, v := range dataset {
		var totalSum float64
		for j := 0; j < int(totalDimensions); j++ {
			difference := float64(v.Vector[j]) - vec[j]
			totalSum = totalSum + math.Pow(difference, 2)
		}
		finalValue := math.Sqrt(totalSum)
		boundedClosest.Add(CalculatedDistance{Distance: finalValue, Label: v.Label})
		if i == 0 {
			lowest = finalValue
		}
		if finalValue < lowest {
			lowest = finalValue
		}
	}
	totalFraudsNeighbours := 0
	for _, value := range *boundedClosest.heap {
		if value.Label == "fraud" {
			totalFraudsNeighbours++
		}
	}
	fmt.Println("Total fraud neighbours(out of 1000):", totalFraudsNeighbours)
	score := float64(totalFraudsNeighbours) / float64(boundedClosest.heap.Len())
	fmt.Println("Score:", score)
	if float64(score) < TRESHOLD {
		return SearchResultStruct{IsPossibleFraud: false, Score: score}, nil
	} else if float64(score) > TRESHOLD {
		return SearchResultStruct{IsPossibleFraud: true, Score: score}, nil
	}
	return SearchResultStruct{IsPossibleFraud: true, Score: score}, errors.New("Error when searching in vector")
}
