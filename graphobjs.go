package main

import (
	"fmt"
)

type TPDFMoveTo struct {
	FPos TPDFCoord
}

func (t TPDFMoveTo) Name() string { return "" }
func (t TPDFMoveTo) CommandXY(x, y PDFFloat) string {
	return FloatStr(x) + " " + FloatStr(y) + " m" + CRLF
}

func (m *TPDFMoveTo) Encode(st PDFWriter) {
	st.Write([]byte(m.Command(m.FPos)))
}

func (m TPDFMoveTo) Command(pos TPDFCoord) string {
	return m.CommandXY(pos.x, pos.y)
}

type TPDFResetPath struct {
}

func (t TPDFResetPath) Name() string { return "" }
func (r *TPDFResetPath) Encode(st PDFWriter) {
	st.Write([]byte("n\n"))
}

type TPDFClosePath struct {
}

func (t TPDFClosePath) Name() string { return "" }
func (c *TPDFClosePath) Encode(st PDFWriter) {
	st.Write([]byte("h\n"))
}

type TPDFStrokePath struct {
}

func (t TPDFStrokePath) Name() string { return "" }
func (s *TPDFStrokePath) Encode(st PDFWriter) {
	st.Write([]byte("S\n"))
}

type TPDFClipPath struct {
}

func (t TPDFClipPath) Name() string { return "" }
func (c *TPDFClipPath) Encode(st PDFWriter) {
	st.Write([]byte("W n\n"))
}

type TPDFPushGraphicsStack struct {
}

func (t TPDFPushGraphicsStack) Name() string { return "" }
func (t *TPDFPushGraphicsStack) Encode(st PDFWriter) {
	st.WriteString(t.Command()) // was , st)
}

func (TPDFPushGraphicsStack) Command() string {
	return "q" + CRLF
}

type TPDFPopGraphicsStack struct {
	TPDFDocumentObject
}

func NewTPDFPopGraphicsStack(ADocument *TPDFDocument) *TPDFPopGraphicsStack {
	return &TPDFPopGraphicsStack{
		TPDFDocumentObject: NewTPDFDocumentObject(ADocument),
	}
}

func (t TPDFPopGraphicsStack) Name() string { return "" }
func (t *TPDFPopGraphicsStack) Encode(st PDFWriter) {
	st.WriteString(t.Command()) // was , st)
	t.Document.CurrentWidth = ""
	t.Document.CurrentColor = ""
}

func (TPDFPopGraphicsStack) Command() string {
	return "Q" + CRLF
}

type TPDFEllipse struct {
	TPDFDocumentObject
	Stroke     bool
	LineWidth  PDFFloat
	Center     TPDFCoord
	Dimensions TPDFCoord
	Fill       bool
}

func (t TPDFEllipse) Name() string { return "" }
func (t *TPDFEllipse) Encode(st PDFWriter) {
	var X, Y, W2, H2, WS, HS PDFFloat
	if t.Stroke {
		t.SetWidth(t.LineWidth, st)
	}
	X = t.Center.x
	Y = t.Center.y
	W2 = t.Dimensions.x / 2
	H2 = t.Dimensions.y / 2
	WS = W2 * BEZIER
	HS = H2 * BEZIER
	st.WriteString(TPDFMoveTo{}.CommandXY(X, Y+H2))
	st.WriteString(TPDFCurveC{}.CommandXY(X, Y+H2-HS, X+W2-WS, Y, X+W2, Y))                // was , st)
	st.WriteString(TPDFCurveC{}.CommandXY(X+W2+WS, Y, X+W2*2, Y+H2-HS, X+W2*2, Y+H2))      // was , st)
	st.WriteString(TPDFCurveC{}.CommandXY(X+W2*2, Y+H2+HS, X+W2+WS, Y+H2*2, X+W2, Y+H2*2)) // was , st)
	st.WriteString(TPDFCurveC{}.CommandXY(X+W2-WS, Y+H2*2, X, Y+H2+HS, X, Y+H2))           // was , st)

	if t.Stroke && t.Fill {
		st.WriteString("b" + CRLF) // was , st)
	} else if t.Fill {
		st.WriteString("f" + CRLF) // was , st)
	} else if t.Stroke {
		st.WriteString("S" + CRLF) // was , st)
	}
}

func NewTPDFEllipse(ADocument *TPDFDocument, APosX, APosY, AWidth, AHeight, ALineWidth PDFFloat, AFill, AStroke bool) *TPDFEllipse {
	return &TPDFEllipse{
		TPDFDocumentObject: NewTPDFDocumentObject(ADocument),
		LineWidth:          ALineWidth,
		Center:             TPDFCoord{x: APosX, y: APosY},
		Dimensions:         TPDFCoord{x: AWidth, y: AHeight},
		Fill:               AFill,
		Stroke:             AStroke,
	}
}

type TPDFCurveC struct {
	TPDFDocumentObject
	Ctrl1, Ctrl2, To TPDFCoord
	Width            PDFFloat
	Stroke           bool
}

func (t TPDFCurveC) Name() string { return "" }

// FIXME:
func (t TPDFCurveC) CommandXY(xCtrl1, yCtrl1, xCtrl2, yCtrl2, xTo, yTo PDFFloat) string {
	// return FloatStr(xCtrl1) + ' ' + FloatStr(yCtrl1) + ' ' +
	// 	FloatStr(xCtrl2) + ' ' + FloatStr(yCtrl2) + ' ' +
	// 	FloatStr(xTo) + ' ' + FloatStr(yTo) + " c" + CRLF
	return "not implemented"
}

func (t TPDFCurveC) Command(ACtrl1, ACtrl2, ATo3 TPDFCoord) string {
	return t.CommandXY(ACtrl1.x, ACtrl1.y, ACtrl2.x, ACtrl2.y, ATo3.x, ATo3.y)
}

func (t *TPDFCurveC) Encode(st PDFWriter) {
	if t.Stroke {
		t.SetWidth(t.Width, st)
	}
	st.WriteString(FloatStr(t.Ctrl1.x) + " " + FloatStr(t.Ctrl1.y) + " " + FloatStr(t.Ctrl2.x) + " " + FloatStr(t.Ctrl2.y) + " " + FloatStr(t.To.x) + " " + FloatStr(t.To.y) + " c" + CRLF) // was , st)
	if t.Stroke {
		st.WriteString("S" + CRLF) // was , st)
	}
}

func NewTPDFCurveC_EX(ADocument *TPDFDocument, ACtrl1X, ACtrl1Y, ACtrl2X, ACtrl2Y, AToX, AToY, AWidth PDFFloat, AStroke bool) *TPDFCurveC {
	return &TPDFCurveC{
		TPDFDocumentObject: NewTPDFDocumentObject(ADocument),
		Ctrl1:              TPDFCoord{x: ACtrl1X, y: ACtrl1Y},
		Ctrl2:              TPDFCoord{x: ACtrl2X, y: ACtrl2Y},
		To:                 TPDFCoord{x: AToX, y: AToY},
		Width:              AWidth,
		Stroke:             AStroke,
	}
}

func NewTPDFCurveC(ADocument *TPDFDocument, ACtrl1, ACtrl2, ATo TPDFCoord, AWidth PDFFloat, AStroke bool) *TPDFCurveC {
	return &TPDFCurveC{
		Ctrl1:  ACtrl1,
		Ctrl2:  ACtrl2,
		To:     ATo,
		Width:  AWidth,
		Stroke: AStroke,
	}
}

type TPDFCurveY struct {
	TPDFDocumentObject
	Stroke bool
	Width  PDFFloat
	P1     TPDFCoord
	P3     TPDFCoord
}

func (t TPDFCurveY) Name() string { return "" }
func (t *TPDFCurveY) Encode(st PDFWriter) {
	if t.Stroke {
		t.SetWidth(t.Width, st)
	}
	st.WriteString(FloatStr(t.P1.x) + " " + FloatStr(t.P1.y) + " " + FloatStr(t.P3.x) + " " + FloatStr(t.P3.y) + " y" + CRLF) // was , st)
	if t.Stroke {
		st.WriteString("S" + CRLF) // was , st)
	}
}

func NewTPDFCurveY_EX(ADocument *TPDFDocument, X1, Y1, X3, Y3, AWidth PDFFloat, AStroke bool) *TPDFCurveY {
	return &TPDFCurveY{
		TPDFDocumentObject: NewTPDFDocumentObject(ADocument),
		P1:                 TPDFCoord{x: X1, y: Y1},
		P3:                 TPDFCoord{x: X3, y: Y3},
		Width:              AWidth,
		Stroke:             AStroke,
	}
}

func NewTPDFCurveY(ADocument *TPDFDocument, AP1, AP3 TPDFCoord, AWidth PDFFloat, AStroke bool) *TPDFCurveY {
	return &TPDFCurveY{
		TPDFDocumentObject: NewTPDFDocumentObject(ADocument),
		P1:                 AP1,
		P3:                 AP3,
		Width:              AWidth,
		Stroke:             AStroke,
	}
}

type TPDFCurveV struct {
	TPDFDocumentObject
	Stroke bool
	Width  PDFFloat
	P2     TPDFCoord
	P3     TPDFCoord
}

func (t TPDFCurveV) Name() string { return "" }
func (t *TPDFCurveV) Encode(st PDFWriter) {
	if t.Stroke {
		t.SetWidth(t.Width, st)
	}
	st.WriteString(FloatStr(t.P2.x) + " " + FloatStr(t.P2.y) + " " + FloatStr(t.P3.x) + " " + FloatStr(t.P3.y) + " v" + CRLF) // was , st)
	if t.Stroke {
		st.WriteString("S" + CRLF) // was , st)
	}
}

func NewTPDFCurveV_EX(ADocument *TPDFDocument, X2, Y2, X3, Y3, AWidth PDFFloat, AStroke bool) *TPDFCurveV {
	return &TPDFCurveV{
		TPDFDocumentObject: NewTPDFDocumentObject(ADocument),
		P2:                 TPDFCoord{x: X2, y: Y2},
		P3:                 TPDFCoord{x: X3, y: Y3},
		Width:              AWidth,
		Stroke:             AStroke,
	}
}

func NewTPDFCurveV(ADocument *TPDFDocument, AP2, AP3 TPDFCoord, AWidth PDFFloat, AStroke bool) *TPDFCurveV {
	return &TPDFCurveV{
		P2:     AP2,
		P3:     AP3,
		Width:  AWidth,
		Stroke: AStroke,
	}
}

type TPDFLineSegment struct {
	TPDFDocumentObject
	Width  PDFFloat
	Stroke bool
	P1, P2 TPDFCoord
}

func (t TPDFLineSegment) Name() string { return "" }
func (ls *TPDFLineSegment) Encode(st PDFWriter) {
	ls.SetWidth(ls.Width, st)
	if ls.Stroke {
		st.WriteString(TPDFMoveTo{}.Command(ls.P1)) // was , st)
	}
	st.WriteString(ls.Command(ls.P2)) // was , st)
	if ls.Stroke {
		st.WriteString("S" + CRLF) // was , st)
	}
}

func (ls *TPDFLineSegment) Command(APos TPDFCoord) string {
	return FloatStr(APos.x) + " " + FloatStr(APos.y) + " l" + CRLF
}

func (ls *TPDFLineSegment) CommandXY(x1, y1 PDFFloat) string {
	return FloatStr(x1) + " " + FloatStr(y1) + " l" + CRLF
}

func (ls *TPDFLineSegment) CommandCoords(APos1, APos2 TPDFCoord) string {
	return TPDFMoveTo{}.Command(APos1) + ls.Command(APos2)
}

func NewTPDFLineSegment(ADocument *TPDFDocument, AWidth, X1, Y1, X2, Y2 PDFFloat, AStroke bool) *TPDFLineSegment {
	return &TPDFLineSegment{
		TPDFDocumentObject: NewTPDFDocumentObject(ADocument),
		Width:              AWidth,
		P1:                 TPDFCoord{x: X1, y: Y1},
		P2:                 TPDFCoord{x: X2, y: Y2},
		Stroke:             AStroke,
	}
}

type TPDFRectangle struct {
	TPDFDocumentObject
	Width      PDFFloat
	TopLeft    TPDFCoord
	Dimensions TPDFCoord
	Fill       bool
	Stroke     bool
}

func (t TPDFRectangle) Name() string { return "" }
func (r *TPDFRectangle) Encode(st PDFWriter) {
	if r.Stroke {
		r.SetWidth(r.Width, st)
	}
	st.WriteString(FloatStr(r.TopLeft.x) + " " + FloatStr(r.TopLeft.y) + " " + FloatStr(r.Dimensions.x) + " " + FloatStr(r.Dimensions.y) + " re" + CRLF) // was , st)
	if r.Stroke && r.Fill {
		st.WriteString("b" + CRLF) // was , st)
	} else if r.Fill {
		st.WriteString("f" + CRLF) // was , st)
	} else if r.Stroke {
		st.WriteString("S" + CRLF) // was , st)
	}
}

func NewTPDFRectangle(ADocument *TPDFDocument, APosX, APosY, AWidth, AHeight, ALineWidth PDFFloat, AFill, AStroke bool) *TPDFRectangle {
	return &TPDFRectangle{
		TPDFDocumentObject: NewTPDFDocumentObject(ADocument),
		TopLeft:            TPDFCoord{x: APosX, y: APosY},
		Dimensions:         TPDFCoord{x: AWidth, y: AHeight},
		Width:              ALineWidth,
		Fill:               AFill,
		Stroke:             AStroke,
	}
}

type TPDFRoundedRectangle struct {
	TPDFDocumentObject
	FBottomLeft, FDimensions TPDFCoord
	FWidth, FRadius          PDFFloat
	FFill, FStroke           bool
}

func (t TPDFRoundedRectangle) Name() string { return "" }
func (rr *TPDFRoundedRectangle) Encode(st PDFWriter) {}

type PDFRoundedRectangle struct {
	TPDFDocumentObject
	Width      PDFFloat
	BottomLeft TPDFCoord
	Dimensions TPDFCoord
	Fill       bool
	Stroke     bool
	Radius     PDFFloat
}

func NewPDFRoundedRectangle(document *TPDFDocument, posX, posY, width, height, radius, lineWidth PDFFloat, fill, stroke bool) *PDFRoundedRectangle {
	return &PDFRoundedRectangle{
		TPDFDocumentObject: NewTPDFDocumentObject(document),
		BottomLeft:         TPDFCoord{x: posX, y: posY},
		Dimensions:         TPDFCoord{x: width, y: height},
		Width:              lineWidth,
		Fill:               fill,
		Stroke:             stroke,
		Radius:             radius,
	}
}

type PDFSurface struct {
	TPDFDocumentObject
	Points []TPDFCoord
	Close  bool
	Fill   bool
}

func NewPDFSurface(document *TPDFDocument, points []TPDFCoord, close, fill bool) *PDFSurface {
	return &PDFSurface{
		TPDFDocumentObject: NewTPDFDocumentObject(document),
		Points:             points,
		Close:              close,
		Fill:               fill,
	}
}
func (t PDFSurface) Name() string { return "" }
func (s *PDFSurface) Encode(st PDFWriter) {
	st.WriteString(TPDFMoveTo{}.CommandXY(s.Points[0].x, s.Points[0].y))
	for i := 1; i < len(s.Points); i++ {
		st.WriteString(fmt.Sprintf("%f %f l%s", s.Points[i].x, s.Points[i].y, CRLF))
	}
	if s.Close {
		st.WriteString("h" + CRLF)
	}
	if s.Fill {
		st.WriteString("f" + CRLF)
	}
}

type TPDFLineStyle struct {
	TPDFDocumentObject
	Style     TPDFPenStyle
	Phase     int
	LineWidth PDFFloat
	LineMask  string
}

func (t TPDFLineStyle) Name() string { return "" }
func NewPDFLineStyle(document *TPDFDocument, style TPDFPenStyle, phase int, lineWidth PDFFloat) *TPDFLineStyle {
	return &TPDFLineStyle{
		TPDFDocumentObject: NewTPDFDocumentObject(document),
		Style:              style,
		Phase:              phase,
		LineWidth:          lineWidth,
	}
}
func (ls *TPDFLineStyle) Encode(st PDFWriter) {
	var lMask string
	w := ls.LineWidth
	if ls.LineMask != "" {
		lMask = ls.LineMask
	} else {
		switch ls.Style {
		case ppsSolid:
			lMask = ""
		case ppsDash:
			lMask = fmt.Sprintf("%f %f", 5*w, 5*w)
		case ppsDot:
			lMask = fmt.Sprintf("%f %f", 0.8*w, 4*w)
		case ppsDashDot:
			lMask = fmt.Sprintf("%f %f %f %f", 5*w, 3*w, 0.8*w, 3*w)
		case ppsDashDotDot:
			lMask = fmt.Sprintf("%f %f %f %f %f %f", 5*w, 3*w, 0.8*w, 3*w, 0.8*w, 3*w)
		}
	}
	st.WriteString(fmt.Sprintf("[%s] %d d%s", lMask, ls.Phase, CRLF))
}

func NewPDFLineStyleCustom(document *TPDFDocument, dashArray []PDFFloat, phase int, lineWidth PDFFloat) *TPDFLineStyle {
	ls := NewPDFLineStyle(document, ppsSolid, phase, lineWidth)
	for _, d := range dashArray {
		if ls.LineMask != "" {
			ls.LineMask += " "
		}
		ls.LineMask += fmt.Sprintf("%f", d*lineWidth)
	}
	return ls
}

type PDFCapStyle struct {
	Style TPDFLineCapStyle
}

func (t PDFCapStyle) Name() string { return "" }

func NewPDFCapStyle(document *TPDFDocument, style TPDFLineCapStyle) *PDFCapStyle {
	return &PDFCapStyle{
		Style: style,
	}
}

func (cs *PDFCapStyle) Encode(st PDFWriter) {
	st.WriteString(fmt.Sprintf("%d J%s", cs.Style, CRLF))
}

type PDFJoinStyle struct {
	Style TPDFLineJoinStyle
}

func (t PDFJoinStyle) Name() string { return "" }
func NewPDFJoinStyle(document *TPDFDocument, style TPDFLineJoinStyle) *PDFJoinStyle {
	return &PDFJoinStyle{
		Style: style,
	}
}

func (js *PDFJoinStyle) Encode(st PDFWriter) {
	st.WriteString(fmt.Sprintf("%d j%s", js.Style, CRLF))
}
