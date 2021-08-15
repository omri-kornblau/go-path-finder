package pathfinder

import (
	"github.com/omri-kornblau/go-path-finder/src/costcalculator"

	"github.com/gonum/diff/fd"
	"github.com/gonum/optimize"
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

	gradientDescent := &optimize.GradientDescent{}
	// temp, err := optimize.Local(problem, calcSpline.GetOptimizationParams(),
	// 	nil, &optimize.BFGS{})
	temp, err := optimize.Local(problem, calcSpline.GetOptimizationParams(),
		&optimize.Settings{
			GradientThreshold: 1e-7,
			FunctionThreshold: 1e-7,
			FunctionConverge: &optimize.FunctionConverge{
				Absolute:   0,
				Relative:   0,
				Iterations: 100,
			},
		}, gradientDescent)

	println(temp)
	return err
}
