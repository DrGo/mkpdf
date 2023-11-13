package main

type FontWidthArray [256]int

type DashArray []float64

type Coord struct {
	x, y float64
}

func NewCoord(x, y float64) Coord {
	return Coord{x, y}
}

type CoordArray []Coord

type Matrix struct {
	_00, _11, _20, _21 float64
}

func (m *Matrix) Transform(point Coord) Coord {
	return Coord{x: m._00*point.x + m._20, y: m._11*point.y + m._21}
}

func (m *Matrix) TransformXY(x, y float64) Coord {
	return Coord{x: m._00*x + m._20, y: m._11*y + m._21}
}

func (m *Matrix) ReverseTransform(point Coord) Coord {
	return Coord{x: (point.x - m._20) / m._00, y: (point.y - m._21) / m._11}
}

func (m *Matrix) SetXScalation(value float64) {
	m._00 = value
}

func (m *Matrix) SetYScalation(value float64) {
	m._11 = value
}

func (m *Matrix) SetXTranslation(value float64) {
	m._20 = value
}

func (m *Matrix) SetYTranslation(value float64) {
	m._21 = value
}
