package store

import (
	"encoding/json"
	"fmt"
	"gin-test/models"
	"os"
	"path/filepath"
)

var Normalizer models.NormalizerStruct

func LoadNormalizer() models.NormalizerStruct {
	path := filepath.Join("public", "normalization.json")
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&Normalizer); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	return Normalizer
}
