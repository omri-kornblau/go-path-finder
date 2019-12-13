package pathfinder

import (
	"github.com/omri-kornblau/go-path-finder/costcalculator"
	"github.com/omri-kornblau/go-path-finder/spline"
)

func PointsToPathCostCalculator(points []spline.Point,
	init costcalculator.CostCalculatorInit) costcalculator.CostCalculator {

	calculators := costcalculator.PathCostCalculator(
		make([]costcalculator.CostCalculator, len(points)-1))

	for index, point := range points[1:] {
		calculators[index] = init([]spline.Point{points[index], point})
	}
	return calculators
}
