package ivf

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"gin-test/internal/store"
	"gin-test/models"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

const (
	OutPath                   = "./partitions"
	MaxReferencesPerPartition = 750000
)

func WritePartitions(in string, maxPartitions uint8) {

	centroidsPerPartition := 4096 / uint16(maxPartitions)
	for i := range maxPartitions {
		fmt.Println(string(i))
		partitionOutPath := OutPath + "/idx_" + strconv.Itoa(int(i)) + ".bin"
		LoadReferencesGzipedPathPartitioned(in, uint32(MaxReferencesPerPartition), i)
		fmt.Println(store.References[0])
		centroids := TrainKMeans(store.References, int(centroidsPerPartition))
		idx := BuildIndex(store.References, centroids)
		t3 := time.Now()
		if err := idx.SerializeToFile(partitionOutPath); err != nil {
			log.Fatalf("save: %v", err)
		}
		info, _ := os.Stat(partitionOutPath)
		log.Printf("saved %d bytes (%.1f MB) to %s in %s",
			info.Size(), float64(info.Size())/(1<<20), partitionOutPath, time.Since(t3))
		idx = nil
		store.References = nil
		debug.FreeOSMemory()
	}
}

func WriteQuantizedPartitionsFromIdxFiles() {
	files, err := os.ReadDir("./partitions")
	if err != nil {
		log.Fatal(err)
	}
	for i, file := range files {
		if !file.IsDir() {
			fmt.Println(file.Name())
			//TODO make it load only the files that are NOT part of the quantized index
			partitionPath := "./partitions/" + file.Name()
			if _, err := os.Stat(partitionPath); err == nil {
				idx, err := LoadIndexFromFile(partitionPath)
				if err != nil {
					log.Fatalf("Failed to load index: %v", err)
				}
				quantOut := "./partitions/quant/idxQuantized_" + strconv.Itoa(int(i)) + ".bin"
				quantIdx := QuantizeIVFIndex(idx)
				fmt.Println(quantIdx.Vectors[200+i], i)
				if err := quantIdx.SerializeToFile(quantOut); err != nil {
					log.Fatalf("save quantized: %v", err)
				}
			} else {
				fmt.Println("WARN: index file not found, running without index")
			}
		}
		debug.FreeOSMemory()
	}
}

func LoadQuantizedPartitions() {
	files, err := os.ReadDir("./partitions/quant")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			fmt.Println(file.Name())
			partitionPath := "./partitions/quant/" + file.Name()
			if _, err := os.Stat(partitionPath); err == nil {
				if _, err := LoadQuantizedIndexFromFile(partitionPath); err != nil {
					log.Fatalf("Failed to load index: %v", err)
				}
			} else {
				fmt.Println("WARN: index file not found, running without index")
			}
			for _, c := range QuantizedClusters {
				GlobalQuantizedClusters = append(GlobalQuantizedClusters, c)
			}
		}
		QuantizedClusters = nil
		debug.FreeOSMemory()
	}
}

func LoadPartitions() {
	files, err := os.ReadDir("./partitions")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			fmt.Println(file.Name())
			//TODO make it load only the files that are NOT part of the quantized index
			partitionPath := "./partitions/" + file.Name()
			if _, err := os.Stat(partitionPath); err == nil {
				if _, err := LoadIndexFromFile(partitionPath); err != nil {
					log.Fatalf("Failed to load index: %v", err)
				}
				for _, c := range Clusters {
					GlobalClusters = append(GlobalClusters, c)
				}

			} else {
				fmt.Println("WARN: index file not found, running without index")
			}
		}
		Clusters = nil
		debug.FreeOSMemory()
	}
}

func LoadReferencesGzipedPathPartitioned(path string, maxReferences uint32, partition uint8) []models.DatasetStruct {
	var references []models.DatasetStruct
	f, _ := os.Open(path)
	defer f.Close()
	gr, _ := gzip.NewReader(f)
	defer gr.Close()

	decoder := json.NewDecoder(gr)

	t, err := decoder.Token()
	if err != nil {
		fmt.Println("Error reading token:", err)
		return []models.DatasetStruct{}
	}

	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		fmt.Println("Expected JSON array start '['")
		return []models.DatasetStruct{}
	}
	fmt.Println(partition)
	count := 0
	addedReferencesCount := 0
	for decoder.More() {

		var v models.DatasetStruct
		if addedReferencesCount >= int(MaxReferencesPerPartition) {
			return references
		}
		if err := decoder.Decode(&v); err != nil {
			fmt.Println("Error decoding obj:", err)
			continue
		}
		if count < int(partition)*int(MaxReferencesPerPartition) {
			count++
			continue
		}
		//references = append(references, v)
		//store.References[addedReferencesCount] = v
		store.References = append(store.References, v)
		addedReferencesCount++
		count++
	}
	/* for _, r := range references {
		store.References[addedReferencesCount] = v
	} */
	return references
}
