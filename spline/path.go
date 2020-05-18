package spline

import (
	"math"
)

type Path struct {
	splines []Spline
	points  []Point
}

func NewPath(splines []Spline) *Path {
	points := []Point{}
	for index, spline := range splines {
		if index == 0 {
			points = append(points, spline.GetPoints()[0])
		}
		// TODO: make this handle more than two points
		points = append(points, spline.GetPoints()[1])
	}
	return &Path{splines, points}
}

func (path *Path) GetPoints() []Point {
	return path.points
}

func (path *Path) GetSplines() []Spline {
	return path.splines
}

func (path *Path) X(s float64) float64 {
	splineIndex, sSpline := path.sForSpline(s)
	return path.splines[splineIndex].X(sSpline)
}

func (path *Path) Y(s float64) float64 {
	splineIndex, sSpline := path.sForSpline(s)
	return path.splines[splineIndex].Y(sSpline)
}

func (path *Path) DX(s float64) float64 {
	splineIndex, sSpline := path.sForSpline(s)
	return path.splines[splineIndex].DX(sSpline)
}

func (path *Path) DY(s float64) float64 {
	splineIndex, sSpline := path.sForSpline(s)
	return path.splines[splineIndex].DY(sSpline)
}

func (path *Path) DDX(s float64) float64 {
	splineIndex, sSpline := path.sForSpline(s)
	return path.splines[splineIndex].DDX(sSpline)
}

func (path *Path) DDY(s float64) float64 {
	splineIndex, sSpline := path.sForSpline(s)
	return path.splines[splineIndex].DDY(sSpline)
}

func (path *Path) SetOptimzationParams(params []float64) {
	lastIndex := uint(0)
	for _, spline := range path.splines {
		currentLength := spline.GetOptimizationParamsLength()
		spline.SetOptimzationParams(params[lastIndex : lastIndex+currentLength])
		lastIndex += currentLength
	}
}

func (path *Path) GetOptimizationParams() []float64 {
	var optimizationParams []float64
	for _, currentSpline := range path.splines {
		optimizationParams =
			append(optimizationParams, currentSpline.GetOptimizationParams()...)
	}
	return optimizationParams
}

func (path *Path) GetOptimizationParamsLength() uint {
	totalLength := uint(0)
	for _, spline := range path.splines {
		totalLength += spline.GetOptimizationParamsLength()
	}
	return totalLength
}

func (path *Path) sForSpline(s float64) (splineIndex uint, sSpline float64) {
	amountOfSplines := float64(len(path.splines))
	splineSize := SRange / amountOfSplines
	index := uint(s / splineSize)
	if s == 1.0 {
		return uint(amountOfSplines) - 1, 1.0
	}
	sSpline = (math.Mod(s, splineSize) * amountOfSplines)

	return index, sSpline
}
