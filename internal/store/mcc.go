package store

import (
	"encoding/json"
	"fmt"
	"gin-test/models"
	"os"
	"path/filepath"
)

var MccMap = map[string]float64{
	"5411": 0.15,
	"5812": 0.30,
	"5912": 0.20,
	"5944": 0.45,
	"7801": 0.80,
	"7802": 0.75,
	"7995": 0.85,
	"4511": 0.35,
	"5311": 0.25,
	"5999": 0.50}

func LoadMccRisk() models.NormalizerStruct {
	path := filepath.Join("public", "mcc_risk.json")
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&MccMap); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	return Normalizer
}
