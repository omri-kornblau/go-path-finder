package spline

type Path struct {
	splines []Spline
}

func NewPath(splines []Spline) *Path {
	return &Path{splines}
}

func (path *Path) sForSpline(s float64) (splineIndex uint, sSpline float64) {
	amountOfSplines := float64(len(path.splines))
	splineSize := SRange / amountOfSplines
	index := s / splineSize

	sSpline = (s - (index * splineSize)) * amountOfSplines

	return uint(index), sSpline
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
	path.splines[splineIndex].DX(sSpline)
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
	for currentSplineIndex := range path.splines {
		currentLength := path.splines[currentSplineIndex].GetOptimizationParamsLength()
		path.splines[currentSplineIndex].SetOptimzationParams(params[lastIndex : lastIndex+currentLength])
	}
}

func (path *Path) GetOptimizationParams() []float64 {
	var optimizationParams []float64
	for _, currentSpline := range path.splines {
		optimizationParams = append(optimizationParams, currentSpline.GetOptimizationParams()...)
	}
	return optimizationParams
}
