package spline

type Point struct {
	X         float64
	Y         float64
	Direction float64
}

type Spline interface {
	X(s float64) float64
	Y(s float64) float64
	DX(s float64) float64
	DY(s float64) float64
	DDX(s float64) float64
	DDY(s float64) float64
	SetOptimzationParams(params []float64)
	GetOptimizationParams() []float64
	GetOptimizationParamsLength() uint
}

const (
	MinSRange float64 = 0
	MaxSRange float64 = 1
	SRange    float64 = MaxSRange - MinSRange
)
