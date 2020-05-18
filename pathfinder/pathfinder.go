package pathfinder

import (
	"github.com/gonum/diff/fd"
	"github.com/gonum/optimize"
	"github.com/omri-kornblau/go-path-finder/costcalculator"
)

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
