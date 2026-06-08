package ivf

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"gin-test/internal/store"
	"gin-test/models"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Cluster struct {
	Centroid []float32
	Lists    []models.DatasetStruct
}

var Clusters = []Cluster{}

func BuildIVFIndex(data []models.DatasetStruct, k int, maxIterations uint8) []Cluster {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	clusters := make([]Cluster, k)
	perm := r.Perm(len(data))
	for i := range k {
		clusters[i].Centroid = data[perm[i]].Vector
	}
	for intr := range maxIterations {
		fmt.Println("Iteration number:", intr+1)
		// Reset references
		for i := range clusters {
			clusters[i].Lists = nil
		}

		for _, reference := range data {
			bestIndex := 0
			minDist := math.MaxFloat64
			for i, cluster := range clusters {
				dist := Distance(cluster.Centroid, reference.Vector)
				if dist < minDist {
					minDist = dist
					bestIndex = i
				}
			}
			clusters[bestIndex].Lists = append(clusters[bestIndex].Lists, reference)
		}

		converged := true
		for i := range clusters {
			if len(clusters[i].Lists) == 0 {
				continue
			}
			//Compute dimensional averages
			newCentroid := make([]float32, 14)
			for _, reference := range clusters[i].Lists {
				for dim := range reference.Vector {
					newCentroid[dim] += reference.Vector[dim]
				}
			}
			for dim := range newCentroid {
				newCentroid[dim] /= float32(len(clusters[i].Lists))
			}
			if Distance(clusters[i].Centroid, newCentroid) > 1e-4 {
				converged = false
			}
			clusters[i].Centroid = newCentroid
		}
		if converged {
			break
		}
	}
	store.References = nil
	return clusters
}

func BuildIVFIndexStreamed(data []models.DatasetStruct, k int, maxIterations uint8) []Cluster {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	clusters := make([]Cluster, k)
	perm := r.Perm(len(data))
	for i := range k {
		clusters[i].Centroid = data[perm[i]].Vector
	}

	for intr := range maxIterations {
		path := filepath.Join("public", "references.json.gz")
		f, _ := os.Open(path)
		gr, _ := gzip.NewReader(f)

		decoder := json.NewDecoder(gr)
		t, err := decoder.Token()
		if err != nil {
			fmt.Println("Error reading token:", err)
			return []Cluster{}
		}

		if delim, ok := t.(json.Delim); !ok || delim != '[' {
			fmt.Println("Expected JSON array start '['")
			return []Cluster{}
		}

		fmt.Println("Iteration number:", intr+1)
		// Reset references
		for i := range clusters {
			clusters[i].Lists = nil
		}

		for decoder.More() {
			var reference models.DatasetStruct
			if err := decoder.Decode(&reference); err != nil {
				fmt.Println("Error decoding obj:", err)
				continue
			}
			bestIndex := 0
			minDist := math.MaxFloat64
			for i, cluster := range clusters {
				dist := Distance(cluster.Centroid, reference.Vector)
				if dist < minDist {
					minDist = dist
					bestIndex = i
				}
			}
			clusters[bestIndex].Lists = append(clusters[bestIndex].Lists, reference)
		}

		converged := true
		for i := range clusters {
			if len(clusters[i].Lists) == 0 {
				continue
			}
			//Compute dimensional averages
			newCentroid := make([]float32, 14)
			for _, reference := range clusters[i].Lists {
				for dim := range reference.Vector {
					newCentroid[dim] += reference.Vector[dim]
				}
			}
			for dim := range newCentroid {
				newCentroid[dim] /= float32(len(clusters[i].Lists))
			}
			if Distance(clusters[i].Centroid, newCentroid) > 1e-4 {
				converged = false
			}
			clusters[i].Centroid = newCentroid
		}
		gr.Close()
		f.Close()
		if converged {
			break
		}
	}
	store.References = nil
	return clusters
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
