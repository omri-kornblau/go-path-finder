package costcalculator

import (
	"math"

	"github.com/omri-kornblau/go-path-finder/utils"

	"github.com/omri-kornblau/go-path-finder/spline"
)

type PathCostCalculator []CostCalculator

func (calculators PathCostCalculator) GetSpline() spline.Spline {
	splines := make([]spline.Spline, len(calculators))
	for index, subCalculator := range calculators {
		splines[index] = subCalculator.GetSpline()
	}

	return spline.NewPath(splines)
}

func (caculators PathCostCalculator) GetCost() (cost float64) {
	cost = float64(0)
	for _, subCalculator := range caculators {
		cost += subCalculator.GetCost()
	}
	cost += caculators.getRadiusContinuityCost(2)
	return cost
}

func (calculators PathCostCalculator) getRadiusContinuityCost(
	costPower float64) (cost float64) {
	cost = float64(0)
	for index, currCalculator := range calculators[1:] {
		prevCalculator := calculators[index]
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
