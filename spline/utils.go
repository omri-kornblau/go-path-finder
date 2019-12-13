package spline

import (
	"math"
)

func PointsDistance(first, second Point) float64 {
	return math.Sqrt(math.Pow(first.Y-second.Y, 2) + math.Pow(first.X-second.X, 2))
}

func calcPolynom(factors []float64, s float64) float64 {
	var result float64 = 0
	for index, factor := range factors {
		result += factor * math.Pow(s, float64(index))
	}

	return result
}

func getDerivativeFactors(factors []float64) []float64 {
	newFactors := make([]float64, len(factors)-1)
	for index, factor := range factors[1:] {
		newFactors[index] = factor * float64(index+1)
	}

	return newFactors
}
