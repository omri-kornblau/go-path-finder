package costcalculator

import (
	"github.com/omri-kornblau/go-path-finder/src/spline"
)

type LinearPolynomWeights struct {
	Position float64
	Angle    float64
	Radius   float64
}

type LinearPolynomCostCalculator struct {
	spline  *spline.LinearPolynom
	weights LinearPolynomWeights
}

func NewLinearPolynomCostCalculator(
	weights LinearPolynomWeights, polynomDegree uint) CostCalculatorInit {

	return func(points []spline.Point) CostCalculator {
		newSpline := spline.NewLinearPolynom(points, polynomDegree)
		return LinearPolynomCostCalculator{&newSpline, weights}
	}
}

func (calculator LinearPolynomCostCalculator) GetSpline() spline.Spline {
	return calculator.spline
}

// GetCosts calculates the cost (false measure) of a linear polynom spline
// This measure is used to optimize the spline according to wanted properties
func (calculator LinearPolynomCostCalculator) GetCost() (cost float64) {
	currentSpline := calculator.GetSpline()
	cost = getPositionCost(currentSpline, positionCostPower)*calculator.weights.Position +
		getAngleCost(currentSpline, angleCostPower)*calculator.weights.Angle +
		getRadiusCost(currentSpline, radiusCostPower, radiusCostResolution)*calculator.weights.Radius

	return cost
}
