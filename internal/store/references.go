package store

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"gin-test/models"
	"os"
)

func LoadReferencesGziped(maxReferences uint32) int8 {
	return LoadReferencesGzipedPath("public/references.json.gz", maxReferences)
}

var References []models.DatasetStruct

func LoadReferencesGzipedPath(path string, maxReferences uint32) int8 {
	f, _ := os.Open(path)
	defer f.Close()
	gr, _ := gzip.NewReader(f)
	defer gr.Close()

	decoder := json.NewDecoder(gr)

	t, err := decoder.Token()
	if err != nil {
		fmt.Println("Error reading token:", err)
		return 0
	}

	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		fmt.Println("Expected JSON array start '['")
		return 0
	}

	count := 0
	for decoder.More() {
		var v models.DatasetStruct
		if count > int(maxReferences) {
			return int8(count)
		}
		if err := decoder.Decode(&v); err != nil {
			fmt.Println("Error decoding obj:", err)
			continue
		}
		References = append(References, v)
		count++
	}
	return int8(count)
}
