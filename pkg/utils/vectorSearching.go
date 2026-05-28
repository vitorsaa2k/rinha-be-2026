package utils

import (
	"errors"
	"fmt"
	"gin-test/models"
	"math"
)

type SearchResultStruct struct {
	IsPossibleFraud bool
	Number          float64
}

func SearchInVector(vec []float64, totalDimensions int8, dataset []models.DatasetStruct) (SearchResultStruct, error) {
	var lowest float64
	var label string
	for i, v := range dataset {
		var totalSum float64
		for j := 0; j < int(totalDimensions); j++ {
			difference := float64(v.Vector[j]) - vec[j]
			totalSum = totalSum + math.Pow(difference, 2)
		}
		finalValue := math.Sqrt(totalSum)
		if i == 0 {
			lowest = finalValue
		}
		if finalValue < lowest {
			lowest = finalValue
			label = v.Label
		}
		fmt.Println(v.Label, finalValue)
	}
	if label == "legit" {
		return SearchResultStruct{IsPossibleFraud: false, Number: lowest}, nil
	} else if label == "fraud" {
		return SearchResultStruct{IsPossibleFraud: true, Number: lowest}, nil
	}
	return SearchResultStruct{IsPossibleFraud: true, Number: lowest}, errors.New("Error when searching in vector")
}
