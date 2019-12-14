package main

import (
	"math"

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
	calculator := pathfinder.NewPathCostCalculator(points, calculatorInit)
	pathfinder.OptimizePath(calculator, &optimize.BFGS{})
}
