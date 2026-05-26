package store

import (
	"gin-test/models"
	"gin-test/pkg/utils"
)

var Dataset = []utils.DatasetStruct{
	{Vector: []float64{0.0100, 0.0833, 0.05}, Label: "legit"},
	{Vector: []float64{0.5796, 0.9167, 1.00}, Label: "fraud"},
	{Vector: []float64{0.0035, 0.1667, 0.05}, Label: "legit"},
	{Vector: []float64{0.9708, 1.0000, 1.00}, Label: "fraud"},
	{Vector: []float64{0.4082, 1.0000, 1.00}, Label: "fraud"},
	{Vector: []float64{0.0092, 0.0833, 0.05}, Label: "legit"},
}
var Normalizer = models.NormalizerStruct{Max_amount: 10000, Max_hour: 23, Max_avg: 5000}
