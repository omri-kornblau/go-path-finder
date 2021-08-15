package utils

import "math"

func Sign(num float64) float64 {
	if num > 0 {
		return 1
	} else if num < 0 {
		return -1
	}
	return 0
}

func AngleDiff(firstAngle, secondAngle float64) float64 {
	diff := secondAngle - firstAngle
	for diff < -2*math.Pi || diff > 2*math.Pi {
		diff += 2*math.Pi - Sign(diff)
	}
	return diff
}

func SetDefault(value, defaultValue float64) float64 {
	if math.IsInf(value, 0) || math.IsNaN(value) {
		return defaultValue
	}
	return value
}
