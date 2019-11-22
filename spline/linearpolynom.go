package spline

type LinearPolynom struct {
	xFactors   []float64
	dXFactors  []float64
	ddXFactors []float64

	yFactors   []float64
	dYFactors  []float64
	ddYFactors []float64

	points        []Point
	polynomDegree uint
}

func calcLinearPolynom(start, end float64, polynomDegree uint) []float64 {
	m := (end - start) / (SRange)
	b := start

	factors := make([]float64, polynomDegree)
	factors[0] = b
	factors[1] = m

	return factors
}

func NewLinearPolynom(points []Point, polynomDegree uint) *LinearPolynom {
	linearPolynom := LinearPolynom{points: points, polynomDegree: polynomDegree}

	firstPoint, lastPoint := points[0], points[len(points)-1]
	xFactors := calcLinearPolynom(firstPoint.X, lastPoint.X, polynomDegree)
	yFactors := calcLinearPolynom(firstPoint.Y, lastPoint.Y, polynomDegree)

	linearPolynom.setXFactors(xFactors)
	linearPolynom.setYFactors(yFactors)

	return &linearPolynom
}

// setXFactors sets xFactors and calculates a derivative factors slice
// to reduce computations when running functions DX and DDX.
func (linearPolynom *LinearPolynom) setXFactors(xFactors []float64) {
	linearPolynom.xFactors = xFactors
	linearPolynom.dXFactors = getDerivativeFactors(xFactors)
	linearPolynom.ddXFactors = getDerivativeFactors(linearPolynom.dXFactors)
}

// setYFactors sets yFactors and calculates a derivative factors slices
// to reduce computations when running functions DX and DDX.
func (linearPolynom *LinearPolynom) setYFactors(yFactors []float64) {
	linearPolynom.yFactors = yFactors
	linearPolynom.dYFactors = getDerivativeFactors(yFactors)
	linearPolynom.ddYFactors = getDerivativeFactors(linearPolynom.dYFactors)
}

func (linearPolynom *LinearPolynom) X(s float64) float64 {
	return calcPolynom(linearPolynom.xFactors, s)
}

func (linearPolynom *LinearPolynom) Y(s float64) float64 {
	return calcPolynom(linearPolynom.yFactors, s)
}

func (linearPolynom *LinearPolynom) DX(s float64) float64 {
	return calcPolynom(linearPolynom.dXFactors, s)
}

func (linearPolynom *LinearPolynom) DY(s float64) float64 {
	return calcPolynom(linearPolynom.dYFactors, s)
}

func (linearPolynom *LinearPolynom) DDX(s float64) float64 {
	return calcPolynom(linearPolynom.ddXFactors, s)
}

func (linearPolynom *LinearPolynom) DDY(s float64) float64 {
	return calcPolynom(linearPolynom.ddYFactors, s)
}

// SetOptimizationParams recieves a slice of params from the optimizer (see gonum optimize)
// and replaces the corresponding values in a LinearPolynom instance with the values from params.
// This is done to change only the params that should be optimized in LinearPolynom and
// unpack the one dimensional slice from gonum optimizer.
func (linearPolynom *LinearPolynom) SetOptimzationParams(params []float64) {
	// TODO: improve factors slice optimzation
	for i := uint(0); i < linearPolynom.polynomDegree; i++ {
		linearPolynom.xFactors[i] = params[i]
		linearPolynom.yFactors[i] = params[i+linearPolynom.polynomDegree]
	}
}

// GetOptimizationParams gets params from LinearPolynom instance that should be
// optimized and turns them into a single one dimensional slice for gonum optimizer.
func (linearPolynom *LinearPolynom) GetOptimizationParams() []float64 {
	// TODO: improve factors slice optimzation
	return append(linearPolynom.xFactors, linearPolynom.yFactors...)
}

func (linearPolynom *LinearPolynom) GetOptimizationParamsLength() uint {
	return linearPolynom.polynomDegree * 2
}
