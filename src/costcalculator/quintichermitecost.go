package costcalculator

import (
	"github.com/omri-kornblau/go-path-finder/src/spline"
)

type QuinticHermiteWeights struct {
	Radius float64
}

type QuinticHermiteCostCalculator struct {
	spline  *spline.QuinticHermite
	weights QuinticHermiteWeights
}

func NewQuinticHermiteCostCalculator(
	weights QuinticHermiteWeights) CostCalculatorInit {

	return func(points []spline.Point) CostCalculator {
		newSpline := spline.NewQuinticHermite(points)
		return QuinticHermiteCostCalculator{&newSpline, weights}
	}
}

func (calculator QuinticHermiteCostCalculator) GetSpline() spline.Spline {
	return calculator.spline
}

// GetCosts calculates the cost (false measure) of a linear polynom spline
// This measure is used to optimize the spline according to wanted properties
func (calculator QuinticHermiteCostCalculator) GetCost() (cost float64) {
	cost = getRadiusCost(calculator.GetSpline(), radiusCostPower, radiusCostResolution) * calculator.weights.Radius

	return cost
}
