package costcalculator

import (
	"math"

	"github.com/omri-kornblau/go-path-finder/spline"
	"github.com/omri-kornblau/go-path-finder/utils"
)

const (
	positionCostPower     = 1
	angleCostPower        = 2
	radiusCostPower       = 4
	radiusContinuityPower = 2
	radiusCostResolution  = 0.1
)

type CostCalculator interface {
	GetCost() (cost float64)
	GetSpline() (spline spline.Spline)
}

type CostCalculatorInit func(points []spline.Point) CostCalculator

func powerCosts(power float64, costs ...float64) float64 {
	finalCost := float64(0)
	for _, cost := range costs {
		finalCost += math.Pow(cost, power)
	}
	return finalCost
}

func getPositionCost(currentSpline spline.Spline, costPower float64) float64 {
	points := currentSpline.GetPoints()
	firstPoint, lastPoint := points[0], points[len(points)-1]

	actualFirstPoint := spline.Point{
		X: currentSpline.X(spline.MinSRange),
		Y: currentSpline.Y(spline.MinSRange)}

	actualLastPoint := spline.Point{
		X: currentSpline.X(spline.MaxSRange),
		Y: currentSpline.Y(spline.MaxSRange)}

	firstCost := spline.PointsDistance(firstPoint, actualFirstPoint)
	lastCost := spline.PointsDistance(lastPoint, actualLastPoint)

	return powerCosts(costPower, firstCost, lastCost)
}

func getAngleCost(currentSpline spline.Spline, costPower float64) float64 {
	points := currentSpline.GetPoints()
	firstPoint, lastPoint := points[0], points[len(points)-1]

	actualFirstDirection := spline.Angle(currentSpline, spline.MinSRange)
	actualLastDirection := spline.Angle(currentSpline, spline.MaxSRange)

	firstCost := utils.AngleDiff(actualFirstDirection, firstPoint.Direction)
	lastCost := utils.AngleDiff(actualLastDirection, lastPoint.Direction)

	return powerCosts(costPower, firstCost, lastCost)
}

func getRadiusCost(currentSpline spline.Spline, costPower float64,
	calcResolution float64) float64 {

	radiusCost := float64(0)
	for s := spline.MinSRange; s < spline.MaxSRange; s += calcResolution {
		radiusAtPoint := spline.Radius(currentSpline, s)
		radiusCostAtPoint := utils.SetDefault(1/radiusAtPoint, 0)
		radiusCost += math.Pow(radiusCostAtPoint, costPower)
	}
	return radiusCost
}
