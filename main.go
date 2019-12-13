package main

import (
	"fmt"
	"math"

	"github.com/gonum/diff/fd"
	"github.com/gonum/optimize"

	"github.com/omri-kornblau/go-path-finder/costcalculator"
	"github.com/omri-kornblau/go-path-finder/pathfinder"
	"github.com/omri-kornblau/go-path-finder/spline"
)

func main() {
	points := []spline.Point{
		spline.Point{0, 1, math.Pi / 2},
		spline.Point{4, 2, math.Pi / 2},
		spline.Point{6, 3, math.Pi / 2}}

	weights := costcalculator.QuinticHermiteWeights{1}
	calculatorInit := costcalculator.NewQuinticHermiteCostCalculator(weights)
	calculator := pathfinder.PointsToPathCostCalculator(points, calculatorInit)
	myPath := calculator.GetSpline()

	costWrapper := func(x []float64) float64 {
		myPath.SetOptimzationParams(x)
		cost := calculator.GetCost()
		return cost
	}

	i := 0
	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			cost := costWrapper(x)
			fmt.Println("iter", i, cost)
			i++
			return cost
		},
		Grad: func(dst, x []float64) {
			fd.Gradient(dst, costWrapper, x, nil)
		},
	}

	_, err := optimize.Local(problem, myPath.GetOptimizationParams(),
		nil, &optimize.BFGS{})

	fmt.Println("E:", err)
}
