package floatutils

import "math"

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))

	return float64(Round(num*output)) / output
}

func BytesToGB(bytes uint64) float64 {
	return float64(bytes) / 1e9
}

func BytesToKB(bytes uint64) float64 {
	return float64(bytes) / 1024
}
