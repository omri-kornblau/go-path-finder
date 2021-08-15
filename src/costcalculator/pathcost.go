package costcalculator

import (
	"math"

	"github.com/omri-kornblau/go-path-finder/src/spline"
	"github.com/omri-kornblau/go-path-finder/src/utils"
)

var x int = 0

type PathWeights struct {
	RadiusContinuity float64
}

type PathCostCalculator struct {
	calculators []CostCalculator
	weights     PathWeights
}

func NewPathCostCalculator(points []spline.Point,
	init CostCalculatorInit,
	pathWeights PathWeights) PathCostCalculator {

	calculators := make([]CostCalculator, len(points)-1)
	pathCostCalculator := PathCostCalculator{calculators, pathWeights}

	for index, point := range points[1:] {
		calculators[index] = init([]spline.Point{points[index], point})
	}

	return pathCostCalculator
}

func (calculator PathCostCalculator) GetSpline() spline.Spline {
	splines := make([]spline.Spline, len(calculator.calculators))
	for index, subCalculator := range calculator.calculators {
		splines[index] = subCalculator.GetSpline()
	}

	return spline.NewPath(splines)
}

func (calculator PathCostCalculator) GetCost() (cost float64) {
	x += 1
	println(x)
	cost = float64(0)
	for _, subCalculator := range calculator.calculators {
		cost += subCalculator.GetCost()
	}
	cost += calculator.weights.RadiusContinuity *
		calculator.getRadiusContinuityCost(radiusContinuityPower)

	return cost
}

func (calculator PathCostCalculator) getRadiusContinuityCost(
	costPower float64) (cost float64) {
	cost = float64(0)
	for index, currCalculator := range calculator.calculators[1:] {
		prevCalculator := calculator.calculators[index]
		// Subtract epsilon (tiny number) to stay away from straight line in
		// the end of each spline
		firstRadius := utils.SetDefault(1/spline.Radius(prevCalculator.GetSpline(),
			spline.MaxSRange-spline.Epsilon), 0)
		secondRadius := utils.SetDefault(1/spline.Radius(currCalculator.GetSpline(),
			spline.MinSRange), 0)
		cost += math.Pow(secondRadius-firstRadius, costPower)
	}
	return cost
}
