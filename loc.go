package main

type TPDFFontWidthArray [256]int

type TDashArray []PDFFloat

type TPDFCoord struct {
	x, y PDFFloat
}
type TPDFCoordArray []TPDFCoord

type TPDFLineStyleDef struct {
	FColor     ARGBColor
	FLineWidth PDFFloat
	FPenStyle  TPDFPenStyle
	FDashArray TDashArray
}

func NewPDFCoord(x, y PDFFloat) TPDFCoord {
	return TPDFCoord{x, y}
}

type TPDFMatrix struct {
	_00, _11, _20, _21 PDFFloat
}

func (m *TPDFMatrix) Transform(point TPDFCoord) TPDFCoord {
	return TPDFCoord{x: m._00*point.x + m._20, y: m._11*point.y + m._21}
}

func (m *TPDFMatrix) TransformXY(x, y PDFFloat) TPDFCoord {
	return TPDFCoord{x: m._00*x + m._20, y: m._11*y + m._21}
}

func (m *TPDFMatrix) ReverseTransform(point TPDFCoord) TPDFCoord {
	return TPDFCoord{x: (point.x - m._20) / m._00, y: (point.y - m._21) / m._11}
}

func (m *TPDFMatrix) SetXScalation(value PDFFloat) {
	m._00 = value
}

func (m *TPDFMatrix) SetYScalation(value PDFFloat) {
	m._11 = value
}

func (m *TPDFMatrix) SetXTranslation(value PDFFloat) {
	m._20 = value
}

func (m *TPDFMatrix) SetYTranslation(value PDFFloat) {
	m._21 = value
}
