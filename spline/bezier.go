package spline

import (
	"math"
)

//this includes edge points and control points
//when smoothing with other bezier we might need to seperate edge from control points
type ControlPoint struct {
	X float64
	Y float64
}

type BezierPoint struct {
	X              float64
	Y              float64
	VelDierction   float64
	VelMag         float64
	AccelDierction float64
	AccelMag       float64
}

type Bezier struct {
	BezierPoints  [2]BezierPoint
	ControlPoints [4]ControlPoint

	xFactors   []float64
	dXFactors  []float64
	ddXFactors []float64

	yFactors   []float64
	dYFactors  []float64
	ddYFactors []float64

	optimizationParamsLength uint
}

// func NewBezierCurve(points []ControlPoint) Bezier {
// 	bezier := Bezier{ControlPoints: points, optimizationParamsLength: uint((len(points) - 2) * 2)}
// 	xFactors := bezier.GetXFactors()
// 	yFactors := bezier.GetYFactors()

// 	bezier.setXFactors(xFactors)
// 	bezier.setYFactors(yFactors)

// 	return bezier
// }

func NewBezierCurve(points [2]BezierPoint) Bezier {
	firstPoint := points[0]
	lastPoint := points[len(points)-1]
	bezier := Bezier{BezierPoints: points, optimizationParamsLength: 8}
	controlPoints := [4]ControlPoint{}

	controlPoints[0].X = math.Cos(firstPoint.VelDierction)*firstPoint.VelMag + firstPoint.X
	controlPoints[0].Y = math.Sin(firstPoint.VelDierction)*firstPoint.VelMag + firstPoint.Y

	controlPoints[1].X = math.Cos(firstPoint.AccelDierction)*firstPoint.AccelMag + firstPoint.X
	controlPoints[1].Y = math.Sin(firstPoint.AccelDierction)*firstPoint.AccelMag + firstPoint.Y

	controlPoints[2].X = math.Cos(lastPoint.VelDierction)*lastPoint.VelMag + lastPoint.X
	controlPoints[2].Y = math.Sin(lastPoint.VelDierction)*lastPoint.VelMag + lastPoint.Y

	controlPoints[3].X = math.Cos(lastPoint.AccelDierction)*lastPoint.AccelMag + lastPoint.X
	controlPoints[3].Y = math.Sin(lastPoint.AccelDierction)*lastPoint.AccelMag + lastPoint.Y

	bezier.ControlPoints = controlPoints

	xFactors := bezier.GetXFactors()
	yFactors := bezier.GetYFactors()

	bezier.setXFactors(xFactors)
	bezier.setYFactors(yFactors)

	return bezier
}

func (bezier Bezier) setXFactors(xFactors []float64) {
	bezier.xFactors = xFactors
	bezier.dXFactors = getDerivativeFactors(xFactors)
	bezier.ddXFactors = getDerivativeFactors(bezier.dXFactors)
}

func (bezier Bezier) setYFactors(yFactors []float64) {
	bezier.yFactors = yFactors
	bezier.dYFactors = getDerivativeFactors(yFactors)
	bezier.ddYFactors = getDerivativeFactors(bezier.dYFactors)
}

func (bezier Bezier) getPoints() []ControlPoint {
	points := []ControlPoint{}
	points[0].X = bezier.BezierPoints[0].X
	points[0].Y = bezier.BezierPoints[0].Y
	for index := range bezier.ControlPoints {
		points[index+1] = bezier.ControlPoints[index]
	}
	points[len(points)].X = bezier.BezierPoints[1].X
	points[len(points)].Y = bezier.BezierPoints[1].Y

	return points
}

// GetXFactors and GetYFactors are the exact same but the only difference is in lines 70 and 86
// in points[i].X/Y
// these functions take a Bezier curve with any number of points and turns it into factors

// func (bezier Bezier) GetXFactors() []float64 {
// 	points := bezier.getPoints()
// 	n := len(points) - 1
// 	xFactors := []float64{}
// 	for j := 0; j <= n; j++ {
// 		a := 0.0
// 		for i := 0; i <= j; i++ {
// 			a += (math.Pow(-1, float64(i+j)) * points[i].X) / float64((Factorial(i) * Factorial(j-i)))
// 		}
// 		b := float64(Factorial(n) / Factorial(n-j))

// 		xFactors = append(xFactors, a*b)
// 	}
// 	return xFactors
// }

// func (bezier Bezier) GetYFactors() []float64 {
// 	points := bezier.getPoints()
// 	n := len(points) - 1
// 	yFactors := []float64{}
// 	for j := 0; j <= n; j++ {
// 		a := 0.0
// 		for i := 0; i <= j; i++ {
// 			a += (math.Pow(-1, float64(i+j)) * points[i].Y) / float64((Factorial(i) * Factorial(j-i)))
// 		}
// 		b := float64(Factorial(n) / Factorial(n-j))

// 		yFactors = append(yFactors, a*b)
// 	}
// 	return yFactors
// }

func (bezier Bezier) getFactors(getPointCoordinate func(point ControlPoint) float64) []float64 {
	points := bezier.getPoints()
	n := len(points) - 1
	factors := []float64{}
	for j := 0; j <= n; j++ {
		a := 0.0
		for i := 0; i <= j; i++ {
			a += (math.Pow(-1, float64(i+j)) * getPointCoordinate(points[i])) / float64((Factorial(i) * Factorial(j-i)))
		}
		b := float64(Factorial(n) / Factorial(n-j))

		factors = append(factors, a*b)
	}
	return factors
}

func (bezier Bezier) GetXFactors() []float64 {
	return bezier.getFactors(func(point ControlPoint) float64 { return point.X })
}

func (bezier Bezier) GetYFactors() []float64 {
	return bezier.getFactors(func(point ControlPoint) float64 { return point.Y })
}

func (bezier Bezier) X(s float64) float64 { return calcPolynom(bezier.xFactors, s) }

func (bezier Bezier) Y(s float64) float64 { return calcPolynom(bezier.yFactors, s) }

func (bezier Bezier) DX(s float64) float64 { return calcPolynom(bezier.dXFactors, s) }

func (bezier Bezier) DY(s float64) float64 { return calcPolynom(bezier.dYFactors, s) }

func (bezier Bezier) DDX(s float64) float64 { return calcPolynom(bezier.ddXFactors, s) }

func (bezier Bezier) DDY(s float64) float64 { return calcPolynom(bezier.ddYFactors, s) }

//the optimzation parameters we work with in the case of Bezier are the X,Y coordinates of the control points
func (bezier Bezier) SetOptimzationParams(params []float64) {
	controlPoints := bezier.ControlPoints
	singleOptimizeParamsLength :=
		bezier.optimizationParamsLength / uint(len(controlPoints))

	for index, point := range controlPoints {
		offset := singleOptimizeParamsLength * uint(index)
		point.X = params[offset]
		point.Y = params[offset+1]
		bezier.ControlPoints[index+1] = point
	}

	xFactors := bezier.GetXFactors()
	yFactors := bezier.GetYFactors()

	bezier.setXFactors(xFactors)
	bezier.setYFactors(yFactors)
}

func (bezier Bezier) GetOptimizationParams() []float64 {
	params := make([]float64, bezier.optimizationParamsLength)
	controlPoints := bezier.ControlPoints
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
// func calcBezier(points ControlPoint[], t float64) float64[] {
// 	n := len(points) - 1
// 	result ControlPoint;
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

func Factorial(x int) int {
	if x == 0 {
		return 1
	}

	return x * Factorial(x-1)
}
