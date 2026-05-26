package utils

func GetLimit(value float64, maxValue float64) float64 {
	division := float64(value) / float64(maxValue)
	if division > 1 {
		return 1.0
	}
	if division < 0 {
		return 0
	}
	return float64(division)
}
