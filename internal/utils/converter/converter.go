package converter

const multiplier = 100

func FloatToInt(n float64) int {
	return int(n * multiplier)
}

func IntToFloat(n int) float64 {
	return float64(n) / multiplier
}
