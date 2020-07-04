// TODO:
// 1: type Bezier (points and control points) (done)
// 2: open to factors 5th-order Bezier (not needed because 5 is done)
// 3: derivatives (done)
// 4: get/set optimization (done)

// 5: general n-th order open to factors (done)
// 6: smoothing with other Bezier (not done yet)

package spline

import (
	"math"
)

//this includes edge points and control points
//when smoothing with other bezier we might need to seperate edge from control points
type BezierPoint struct {
	X float64
	Y float64
}

type Bezier struct {
	points []BezierPoint

	xFactors   []float64
	dXFactors  []float64
	ddXFactors []float64

	yFactors   []float64
	dYFactors  []float64
	ddYFactors []float64

	optimizationParamsLength uint
}

func newBezierCurve(points []BezierPoint) Bezier {
	bezier := Bezier{points: points, optimizationParamsLength: uint((len(points) - 2) * 2)}
	xFactors := bezier.GetXFactors()
	yFactors := bezier.GetYFactors()

	bezier.SetXFactors(xFactors)
	bezier.SetYFactors(yFactors)

	return bezier
}

func (bezier Bezier) SetXFactors(xFactors []float64) {
	bezier.xFactors = xFactors
	bezier.dXFactors = getDerivativeFactors(xFactors)
	bezier.ddXFactors = getDerivativeFactors(bezier.dXFactors)
}

func (bezier Bezier) SetYFactors(yFactors []float64) {
	bezier.yFactors = yFactors
	bezier.dYFactors = getDerivativeFactors(yFactors)
	bezier.ddYFactors = getDerivativeFactors(bezier.dYFactors)
}

//GetXFactors and GetYFactors are the exact same but the only difference is in lines 70 and 86
//in points[i].X/Y
//these functions take a Bezier curve with any number of points and turns it into factors
func (bezier Bezier) GetXFactors() []float64 {
	points := bezier.points
	n := len(points) - 1
	xFactors := []float64{}
	for j := 0; j <= n; j++ {
		a := 0.0
		for i := 0; i <= j; i++ {
			a += (math.Pow(-1, float64(i+j)) * points[i].X) / float64((Factorial(i) * Factorial(j-i)))
		}
		b := float64(Factorial(n) / Factorial(n-j))

		xFactors = append(xFactors, a*b)
	}
	return xFactors
}

func (bezier Bezier) GetYFactors() []float64 {
	points := bezier.points
	n := len(points) - 1
	yFactors := []float64{}
	for j := 0; j <= n; j++ {
		a := 0.0
		for i := 0; i <= j; i++ {
			a += (math.Pow(-1, float64(i+j)) * points[i].Y) / float64((Factorial(i) * Factorial(j-i)))
		}
		b := float64(Factorial(n) / Factorial(n-j))

		yFactors = append(yFactors, a*b)
	}
	return yFactors
}

func (bezier Bezier) X(s float64) float64 {
	return calcPolynom(bezier.xFactors, s)
}

func (bezier Bezier) Y(s float64) float64 {
	return calcPolynom(bezier.yFactors, s)
}

func (bezier Bezier) DX(s float64) float64 {
	return calcPolynom(bezier.dXFactors, s)
}

func (bezier Bezier) DY(s float64) float64 {
	return calcPolynom(bezier.dYFactors, s)
}

func (bezier Bezier) DDX(s float64) float64 {
	return calcPolynom(bezier.ddXFactors, s)
}

func (bezier Bezier) DDY(s float64) float64 {
	return calcPolynom(bezier.ddYFactors, s)
}

//the optimzation parameters we work with in the case of Bezier are the X,Y coordinates of the control points
func (bezier Bezier) SetOptimzationParams(params []float64) {
	controlPoints := bezier.points[1 : len(bezier.points)-1]
	singleOptimizeParamsLength :=
		bezier.optimizationParamsLength / uint(len(controlPoints))

	for index, point := range controlPoints {
		offset := singleOptimizeParamsLength * uint(index)
		point.X = params[offset]
		point.Y = params[offset+1]
		bezier.points[index+1] = point
	}

	xFactors := bezier.GetXFactors()
	yFactors := bezier.GetYFactors()

	bezier.SetXFactors(xFactors)
	bezier.SetYFactors(yFactors)
}

func (bezier Bezier) GetOptimizationParams() []float64 {
	params := make([]float64, bezier.optimizationParamsLength)
	controlPoints := bezier.points[1 : len(bezier.points)-1]
	singleOptimizeParamsLength :=
		bezier.optimizationParamsLength /
			uint(len(controlPoints))

	for index, controlPoint := range controlPoints {
		offset := singleOptimizeParamsLength * uint(index)
		params[offset] = controlPoint.X
		params[offset+1] = controlPoint.Y
	}

	return params
}

//explicit definition
// func calcBezier(points BezierPoint[], t float64) float64[] {
// 	n := len(points) - 1
// 	result BezierPoint;
// 	for i := 0; i <= n; i++ {
// 		point := points[i]
// 		binomFactor := binomialCoefficient(n, i)
// 		result.X += binomFactor * point.X * math.Pow(t, i) * math.Pow(1-t, n-i)
// 		result.Y += binomFactor * point.Y * math.Pow(t, i) * math.Pow(1-t, n-i)
// 	}
// 	return result
// }

//utils
//binomial coefficients is used just for explicit definition
// func binomialCoefficient(n int, k int) {
// 	return Factorial(n) / ( Factorial(k) * Factorial(n-k) )
// }
