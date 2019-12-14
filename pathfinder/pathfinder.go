package pathfinder

import (
	"github.com/gonum/diff/fd"
	"github.com/gonum/optimize"
	"github.com/omri-kornblau/go-path-finder/costcalculator"
	"github.com/omri-kornblau/go-path-finder/spline"
)

func NewPathCostCalculator(points []spline.Point,
	init costcalculator.CostCalculatorInit) costcalculator.CostCalculator {

	calculators := costcalculator.PathCostCalculator(
		make([]costcalculator.CostCalculator, len(points)-1))

	for index, point := range points[1:] {
		calculators[index] = init([]spline.Point{points[index], point})
	}
	return calculators
}

func OptimizePath(calculator costcalculator.CostCalculator,
	method optimize.Method) error {

	calcSpline := calculator.GetSpline()

	costWrapper := func(x []float64) float64 {
		calcSpline.SetOptimzationParams(x)
		cost := calculator.GetCost()
		return cost
	}

	problem := optimize.Problem{
		Func: costWrapper,
		Grad: func(dst, x []float64) {
			fd.Gradient(dst, costWrapper, x, nil)
		},
	}

	_, err := optimize.Local(problem, calcSpline.GetOptimizationParams(),
		nil, &optimize.BFGS{})

	return err
}
