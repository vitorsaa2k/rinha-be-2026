package store

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"gin-test/models"
	"os"
	"path/filepath"
)

var References = []models.DatasetStruct{}

func LoadReferences() []models.DatasetStruct {
	var references []models.DatasetStruct
	path := filepath.Join("public", "references.json")
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		/* return []models.DatasetStruct{}, err */
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&references); err != nil {
		fmt.Println("Error decoding JSON:", err)
		/* return []models.DatasetStruct{}, err */
	}
	fmt.Printf("Loaded %d references", len(references))
	return references
}

func LoadReferencesStreamed(maxReferences uint32) int8 {
	path := filepath.Join("public", "references.json")
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		/* return []models.DatasetStruct{}, err */
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
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

	/* t, err := decoder.Token()
	if err != nil {
		fmt.Println("Error reading close token:", err)
		return
	} */
	fmt.Printf("Successfully processed %d items with minimal memory usage!\n", count)
	return int8(count)
}

func LoadReferencesGziped(maxReferences uint32) int8 {
	path := filepath.Join("public", "references.json.gz")
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
