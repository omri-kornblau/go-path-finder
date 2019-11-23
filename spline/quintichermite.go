package spline

import (
	"math"
)

type HermitePoint struct {
	X         float64
	Y         float64
	Dx        float64
	Dy        float64
	Ddx       float64
	Ddy       float64
	magnitude float64
	direction float64
}

func NewHermitePoint(point Point, magnitude float64) *HermitePoint {
	dx := math.Cos(point.Direction) * magnitude
	dy := math.Sin(point.Direction) * magnitude
	return &HermitePoint{0, 0, dx, dy, 0, 0, magnitude, point.Direction}
}

func (hermitePoint *HermitePoint) setMagnitude(magnitude float64) {
	hermitePoint.magnitude = magnitude
	hermitePoint.Dx = math.Cos(hermitePoint.direction) * magnitude
	hermitePoint.Dy = math.Sin(hermitePoint.direction) * magnitude
}

func (hermitePoint *HermitePoint) getMagnitude() float64 {
	return hermitePoint.magnitude
}

func calcQuinticHermitePolynom(p0A, p0dA, p0ddA, p1A, p1dA, p1ddA float64) []float64 {
	factors := make([]float64, 5)
	factors[0] = p0A
	factors[1] = p0dA
	factors[2] = 0.5 * p0ddA
	factors[3] = -10*p0A - 6*p0dA - 1.5*p0ddA + 0.5*p1ddA - 4*p1dA + 10*p1A
	factors[4] = 15*p0A + 8*p0dA + 1.5*p0ddA - p1ddA + 7*p1dA - 15*p1A
	factors[5] = -6*p0A - 3*p0dA - 0.5*p0ddA + 0.5*p1ddA - 3*p1dA + 6*p1A
	return factors
}

type QuinticHermite struct {
	xFactors   []float64
	dXFactors  []float64
	ddXFactors []float64

	yFactors   []float64
	dYFactors  []float64
	ddYFactors []float64

	points        []Point
	hermitePoints [2]*HermitePoint

	optimizationParamsLength uint
}

func NewQuinticHermite(points []Point) *QuinticHermite {
	quinticHermite := &QuinticHermite{points: points, optimizationParamsLength: 6}

	firstPoint, lastPoint := points[0], points[len(points)-1]
	magnitude := 1.2 * pointsDistance(firstPoint, lastPoint)

	p0 := NewHermitePoint(firstPoint, magnitude)
	p1 := NewHermitePoint(lastPoint, magnitude)
	xFactors := calcQuinticHermitePolynom(p0.X, p0.Dx, p0.Ddx, p1.X, p1.Dx, p1.Ddx)
	yFactors := calcQuinticHermitePolynom(p0.Y, p0.Dy, p0.Ddy, p1.Y, p1.Dy, p1.Ddy)

	quinticHermite.setXFactors(xFactors)
	quinticHermite.setYFactors(yFactors)
	quinticHermite.hermitePoints[0], quinticHermite.hermitePoints[1] = p0, p1

	return quinticHermite
}

// setXFactors sets xFactors and calculates a derivative factors slice
// to reduce computations when running functions DX and DDX.
func (quinticHermite *QuinticHermite) setXFactors(xFactors []float64) {
	quinticHermite.xFactors = xFactors
	quinticHermite.dXFactors = getDerivativeFactors(xFactors)
	quinticHermite.ddXFactors = getDerivativeFactors(quinticHermite.dXFactors)
}

// setYFactors sets yFactors and calculates a derivative factors slices
// to reduce computations when running functions DX and DDX.
func (quinticHermite *QuinticHermite) setYFactors(yFactors []float64) {
	quinticHermite.yFactors = yFactors
	quinticHermite.dYFactors = getDerivativeFactors(yFactors)
	quinticHermite.ddYFactors = getDerivativeFactors(quinticHermite.dYFactors)
}

func (quinticHermite *QuinticHermite) X(s float64) float64 {
	return calcPolynom(quinticHermite.xFactors, s)
}

func (quinticHermite *QuinticHermite) Y(s float64) float64 {
	return calcPolynom(quinticHermite.yFactors, s)
}

func (quinticHermite *QuinticHermite) DX(s float64) float64 {
	return calcPolynom(quinticHermite.dXFactors, s)
}

func (quinticHermite *QuinticHermite) DY(s float64) float64 {
	return calcPolynom(quinticHermite.dYFactors, s)
}

func (quinticHermite *QuinticHermite) DDX(s float64) float64 {
	return calcPolynom(quinticHermite.ddXFactors, s)
}

func (quinticHermite *QuinticHermite) DDY(s float64) float64 {
	return calcPolynom(quinticHermite.ddYFactors, s)
}

// SetOptimizationParams recieves a slice of params from the optimizer (see gonum optimize)
// and replaces the corresponding values in a QuinticHermite instance with the values from params.
// This is done to change only the params that should be optimized in QuinticHermite and
// unpack the one dimensional slice from gonum optimizer.
func (quinticHermite *QuinticHermite) SetOptimzationParams(params []float64) {
	for index, hermitePoint := range quinticHermite.hermitePoints {
		hermitePoint.Ddx = params[index]
		hermitePoint.Ddy = params[index+1]
		hermitePoint.setMagnitude(params[index+2])
	}
}

// GetOptimizationParams gets params from QuinticHermite instance that should be
// optimized and turns them into a single one dimensional slice for gonum optimizer.
func (quinticHermite *QuinticHermite) GetOptimizationParams() []float64 {
	params := make([]float64, quinticHermite.optimizationParamsLength)
	for index, hermitePoint := range quinticHermite.hermitePoints {
		params[index] = hermitePoint.Ddx
		params[index+1] = hermitePoint.Ddy
		params[index+2] = hermitePoint.getMagnitude()
	}
	return params
}

func (quinticHermite *QuinticHermite) GetOptimizationLength() uint {
	return quinticHermite.optimizationParamsLength
}
