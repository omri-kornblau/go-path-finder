package main

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"github.com/gonum/optimize"

	"github.com/omri-kornblau/go-path-finder/src/costcalculator"
	"github.com/omri-kornblau/go-path-finder/src/pathfinder"
	"github.com/omri-kornblau/go-path-finder/src/spline"
)

func main() {
	points := []spline.Point{
		{X: 0, Y: 1, Direction: math.Pi / 2},
		{X: 4, Y: 2, Direction: math.Pi / 2},
		{X: 6, Y: 3, Direction: math.Pi / 2},
	}

	// hermiteWeights := costcalculator.QuinticHermiteWeights{Radius: 5}
	lpWeights := costcalculator.LinearPolynomWeights{
		Position: 0,
		Angle:    100000000,
		Radius:   0,
	}
	pathWeights := costcalculator.PathWeights{RadiusContinuity: 0}

	// calculatorInit := costcalculator.NewQuinticHermiteCostCalculator(
	// 	hermiteWeights)
	calculatorInit := costcalculator.NewLinearPolynomCostCalculator(lpWeights, 5)
	path := costcalculator.NewPathCostCalculator(points,
		calculatorInit,
		pathWeights)
	pathfinder.OptimizePath(path, &optimize.BFGS{})

	pts := generatePoints(path, 100)
	savePointsPlot(pts, convertPathPoints(points))

	println("Cost: ", path.GetCost())
}

func savePointsPlot(pts plotter.XYs, points plotter.XYs) {
	p := plot.New()

	p.Title.Text = "Path Plot"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	s, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}

	err = plotutil.AddLines(p, "Path", pts)
	if err != nil {
		panic(err)
	}

	err = plotutil.AddScatters(p, "Path", points)
	if err != nil {
		panic(err)
	}

	if err = p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

func generatePoints(path costcalculator.CostCalculator, pointAmount int) plotter.XYs {
	pts := make(plotter.XYs, pointAmount)

	inc := 1 / float64(pointAmount-1)

	for i := range pts {
		s := float64(i) * inc
		pts[i].X = path.GetSpline().X(s)
		pts[i].Y = path.GetSpline().Y(s)
	}

	return pts
}

func convertPathPoints(pathPoints []spline.Point) plotter.XYs {
	pts := make(plotter.XYs, len(pathPoints))

	for i, pathPoint := range pathPoints {
		pts[i].X = pathPoint.X
		pts[i].Y = pathPoint.Y
	}

	return pts
}
