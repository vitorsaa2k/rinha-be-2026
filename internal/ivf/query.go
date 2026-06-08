package ivf

import (
	"gin-test/models"
	"math"
)

type ClosestCentroids struct {
	cluster  Cluster
	distance float64
}

func SearchIVF(normalizedVector []float64, k uint8) []models.DatasetStruct {
	parsedVector := make([]float32, len(normalizedVector))
	for i, v := range normalizedVector {
		parsedVector[i] = float32(v)
	}
	closestCentroids := make([]ClosestCentroids, k)
	// Initialize closestCentroids with high values
	for i := range closestCentroids {
		closestCentroids[i].distance = math.MaxFloat64
	}
	highestDistanceIndex := 0
	for _, cluster := range Clusters {
		distance := Distance(cluster.Centroid, parsedVector)
		for i, c := range closestCentroids {
			if closestCentroids[highestDistanceIndex].distance < c.distance {
				highestDistanceIndex = i
			}
		}
		if distance < closestCentroids[highestDistanceIndex].distance {
			closestCentroids[highestDistanceIndex] = ClosestCentroids{cluster: cluster, distance: distance}
		}
	}
	var finalDataset []models.DatasetStruct

	for _, c := range closestCentroids {
		finalDataset = append(finalDataset, c.cluster.Lists...)
	}
	return finalDataset
}
