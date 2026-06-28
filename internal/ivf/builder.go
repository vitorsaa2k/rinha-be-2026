package ivf

import (
	"fmt"
	"gin-test/models"
	"math"
	"math/rand"
)

const (
	SampleSize = 65536
	MaxIters   = 6
)

const QuantScale = 10000

type Cluster struct {
	Centroid []float32
	Lists    []models.DatasetStruct
}

type QuantizedCluster struct {
	Centroid []int16
	Lists    []models.QuantizedData
}

var Clusters = []Cluster{}
var QuantizedClusters = []QuantizedCluster{}

// TrainKMeans runs k-means++ on a random 65K sample and returns K centroids.
func TrainKMeans(data []models.DatasetStruct, k int) []Cluster {
	rng := rand.New(rand.NewSource(0xCAFEBABE))

	sample := drawSample(data, SampleSize, rng)

	centroids := make([]Cluster, k)
	for i := range centroids {
		centroids[i].Centroid = make([]float32, dims)
	}

	initKMeansPlusPlus(sample, centroids, rng)

	for iter := 0; iter < MaxIters; iter++ {
		fmt.Println("K-means iteration:", iter+1)

		counts := make([]uint32, k)
		sums := make([][14]float64, k)

		for _, v := range sample {
			best := 0
			bestD := math.MaxFloat64
			for c, cl := range centroids {
				d := Distance(cl.Centroid, v.Vector)
				if d < bestD {
					bestD = d
					best = c
				}
			}
			for dim := 0; dim < dims; dim++ {
				sums[best][dim] += float64(v.Vector[dim])
			}
			counts[best]++
		}

		for c := 0; c < k; c++ {
			if counts[c] == 0 {
				continue
			}
			for dim := 0; dim < dims; dim++ {
				centroids[c].Centroid[dim] = float32(sums[c][dim] / float64(counts[c]))
			}
		}
	}

	return centroids
}

func drawSample(data []models.DatasetStruct, n int, rng *rand.Rand) []models.DatasetStruct {
	if n > len(data) {
		n = len(data)
	}
	picked := make(map[int]struct{}, n)
	sample := make([]models.DatasetStruct, n)
	for i := 0; i < n; i++ {
		var id int
		for {
			id = rng.Intn(len(data))
			if _, dup := picked[id]; !dup {
				break
			}
		}
		picked[id] = struct{}{}
		vec := make([]float32, dims)
		copy(vec, data[id].Vector)
		sample[i] = models.DatasetStruct{
			Vector: vec,
			Label:  data[id].Label,
		}
	}
	return sample
}

func initKMeansPlusPlus(sample []models.DatasetStruct, centroids []Cluster, rng *rand.Rand) {
	n := len(sample)
	first := rng.Intn(n)
	copy(centroids[0].Centroid, sample[first].Vector)

	d2 := make([]float64, n)
	for i, v := range sample {
		d2[i] = SquaredDistance(centroids[0].Centroid, v.Vector)
	}

	for c := 1; c < len(centroids); c++ {
		total := 0.0
		for _, d := range d2 {
			total += d
		}
		if total <= 0 {
			pick := rng.Intn(n)
			copy(centroids[c].Centroid, sample[pick].Vector)
			continue
		}
		r := rng.Float64() * total
		acc := 0.0
		pick := n - 1
		for i, d := range d2 {
			acc += d
			if acc >= r {
				pick = i
				break
			}
		}
		copy(centroids[c].Centroid, sample[pick].Vector)
		for i, v := range sample {
			d := SquaredDistance(centroids[c].Centroid, v.Vector)
			if d < d2[i] {
				d2[i] = d
			}
		}
	}
}

func Distance(p1, p2 []float32) float64 {
	sum := 0.0
	for i := range p1 {
		diff := p1[i] - p2[i]
		sum += float64(diff) * float64(diff)
	}
	return math.Sqrt(sum)
}

func SquaredDistance(p1, p2 []float32) float64 {
	sum := 0.0
	for i := range p1 {
		diff := p1[i] - p2[i]
		sum += float64(diff) * float64(diff)
	}
	return sum
}

func QuantizedDistance(a, b []int16) uint64 {
	var sum uint64
	for i := range a {
		diff := int64(a[i]) - int64(b[i])
		sum += uint64(diff * diff)
	}
	return sum
}
