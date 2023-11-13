package main

import (
	"fmt"
)
type DocObj struct {
	Document      *Document
	FLineCapStyle LineCapStyle
}

func NewDocObj(doc *Document) DocObj {
	do := DocObj{Document: doc}
	if do.Document != nil {
		do.FLineCapStyle = do.Document.LineCapStyle
	}
	return do
}

func (docObj *DocObj) SetWidth(width float64, st PDFWriter) {
	S := fmt.Sprintf("%f w", width)
	if S != docObj.Document.CurrentWidth {
		st.Writef("%d J\n", docObj.FLineCapStyle)
		st.WriteString(S + "\n")
		docObj.Document.CurrentWidth = S
	}
}


type MoveTo struct {
	FPos Coord
}

func (t MoveTo) CommandXY(x, y float64) string {
	return FloatStr(x) + " " + FloatStr(y) + " m" + CRLF
}

func (m *MoveTo) Encode(st PDFWriter) {
	st.Write([]byte(m.Command(m.FPos)))
}

func (m MoveTo) Command(pos Coord) string {
	return m.CommandXY(pos.x, pos.y)
}

type ResetPath struct{}

func (r *ResetPath) Encode(st PDFWriter) {
	st.Write([]byte("n\n"))
}

type ClosePath struct{}

func (c *ClosePath) Encode(st PDFWriter) {
	st.Write([]byte("h\n"))
}

type StrokePath struct{}

func (s *StrokePath) Encode(st PDFWriter) {
	st.Write([]byte("S\n"))
}

type ClipPath struct{}

func (c *ClipPath) Encode(st PDFWriter) {
	st.Write([]byte("W n\n"))
}

type PushGraphicsStack struct{}

func (t *PushGraphicsStack) Encode(st PDFWriter) {
	st.WriteString(t.Command()) 
}

func (PushGraphicsStack) Command() string {
	return "q" + CRLF
}

type PopGraphicsStack struct {
	DocObj
}

func NewTPDFPopGraphicsStack(ADocument *Document) *PopGraphicsStack {
	return &PopGraphicsStack{
		DocObj: NewDocObj(ADocument),
	}
}

func (t *PopGraphicsStack) Encode(st PDFWriter) {
	st.WriteString(t.Command()) 
	t.Document.CurrentWidth = ""
	t.Document.CurrentColor = ""
}

func (PopGraphicsStack) Command() string {
	return "Q" + CRLF
}

type Ellipse struct {
	DocObj
	Stroke     bool
	LineWidth  float64
	Center     Coord
	Dimensions Coord
	Fill       bool
}

func (t *Ellipse) Encode(st PDFWriter) {
	var X, Y, W2, H2, WS, HS float64
	if t.Stroke {
		t.SetWidth(t.LineWidth, st)
	}
	X = t.Center.x
	Y = t.Center.y
	W2 = t.Dimensions.x / 2
	H2 = t.Dimensions.y / 2
	WS = W2 * BEZIER
	HS = H2 * BEZIER
	st.WriteString(MoveTo{}.CommandXY(X, Y+H2))
	st.WriteString(CurveC{}.CommandXY(X, Y+H2-HS, X+W2-WS, Y, X+W2, Y))                
	st.WriteString(CurveC{}.CommandXY(X+W2+WS, Y, X+W2*2, Y+H2-HS, X+W2*2, Y+H2))      
	st.WriteString(CurveC{}.CommandXY(X+W2*2, Y+H2+HS, X+W2+WS, Y+H2*2, X+W2, Y+H2*2)) 
	st.WriteString(CurveC{}.CommandXY(X+W2-WS, Y+H2*2, X, Y+H2+HS, X, Y+H2))           

	if t.Stroke && t.Fill {
		st.WriteString("b" + CRLF) 
	} else if t.Fill {
		st.WriteString("f" + CRLF) 
	} else if t.Stroke {
		st.WriteString("S" + CRLF) 
	}
}

func NewEllipse(ADocument *Document, APosX, APosY, AWidth, AHeight, ALineWidth float64, AFill, AStroke bool) *Ellipse {
	return &Ellipse{
		DocObj:     NewDocObj(ADocument),
		LineWidth:  ALineWidth,
		Center:     Coord{x: APosX, y: APosY},
		Dimensions: Coord{x: AWidth, y: AHeight},
		Fill:       AFill,
		Stroke:     AStroke,
	}
}

type CurveC struct {
	DocObj
	Ctrl1, Ctrl2, To Coord
	Width            float64
	Stroke           bool
}

// FIXME:
func (t CurveC) CommandXY(xCtrl1, yCtrl1, xCtrl2, yCtrl2, xTo, yTo float64) string {
	// return FloatStr(xCtrl1) + ' ' + FloatStr(yCtrl1) + ' ' +
	// 	FloatStr(xCtrl2) + ' ' + FloatStr(yCtrl2) + ' ' +
	// 	FloatStr(xTo) + ' ' + FloatStr(yTo) + " c" + CRLF
	return "not implemented"
}

func (t CurveC) Command(ACtrl1, ACtrl2, ATo3 Coord) string {
	return t.CommandXY(ACtrl1.x, ACtrl1.y, ACtrl2.x, ACtrl2.y, ATo3.x, ATo3.y)
}

func (t *CurveC) Encode(st PDFWriter) {
	if t.Stroke {
		t.SetWidth(t.Width, st)
	}
	st.WriteString(FloatStr(t.Ctrl1.x) + " " + FloatStr(t.Ctrl1.y) + " " + FloatStr(t.Ctrl2.x) + " " + FloatStr(t.Ctrl2.y) + " " + FloatStr(t.To.x) + " " + FloatStr(t.To.y) + " c" + CRLF) 
	if t.Stroke {
		st.WriteString("S" + CRLF) 
	}
}

func NewCurveC_EX(ADocument *Document, ACtrl1X, ACtrl1Y, ACtrl2X, ACtrl2Y, AToX, AToY, AWidth float64, AStroke bool) *CurveC {
	return &CurveC{
		DocObj: NewDocObj(ADocument),
		Ctrl1:  Coord{x: ACtrl1X, y: ACtrl1Y},
		Ctrl2:  Coord{x: ACtrl2X, y: ACtrl2Y},
		To:     Coord{x: AToX, y: AToY},
		Width:  AWidth,
		Stroke: AStroke,
	}
}

func NewCurveC(ADocument *Document, ACtrl1, ACtrl2, ATo Coord, AWidth float64, AStroke bool) *CurveC {
	return &CurveC{
		Ctrl1:  ACtrl1,
		Ctrl2:  ACtrl2,
		To:     ATo,
		Width:  AWidth,
		Stroke: AStroke,
	}
}

type CurveY struct {
	DocObj
	Stroke bool
	Width  float64
	P1     Coord
	P3     Coord
}

func (t *CurveY) Encode(st PDFWriter) {
	if t.Stroke {
		t.SetWidth(t.Width, st)
	}
	st.WriteString(FloatStr(t.P1.x) + " " + FloatStr(t.P1.y) + " " + FloatStr(t.P3.x) + " " + FloatStr(t.P3.y) + " y" + CRLF) 
	if t.Stroke {
		st.WriteString("S" + CRLF) 
	}
}

func NewCurveY_EX(ADocument *Document, X1, Y1, X3, Y3, AWidth float64, AStroke bool) *CurveY {
	return &CurveY{
		DocObj: NewDocObj(ADocument),
		P1:     Coord{x: X1, y: Y1},
		P3:     Coord{x: X3, y: Y3},
		Width:  AWidth,
		Stroke: AStroke,
	}
}

func NewCurveY(ADocument *Document, AP1, AP3 Coord, AWidth float64, AStroke bool) *CurveY {
	return &CurveY{
		DocObj: NewDocObj(ADocument),
		P1:     AP1,
		P3:     AP3,
		Width:  AWidth,
		Stroke: AStroke,
	}
}

type CurveV struct {
	DocObj
	Stroke bool
	Width  float64
	P2     Coord
	P3     Coord
}

func (t *CurveV) Encode(st PDFWriter) {
	if t.Stroke {
		t.SetWidth(t.Width, st)
	}
	st.WriteString(FloatStr(t.P2.x) + " " + FloatStr(t.P2.y) + " " + FloatStr(t.P3.x) + " " + FloatStr(t.P3.y) + " v" + CRLF) 
	if t.Stroke {
		st.WriteString("S" + CRLF) 
	}
}

func NewCurveV_EX(ADocument *Document, X2, Y2, X3, Y3, AWidth float64, AStroke bool) *CurveV {
	return &CurveV{
		DocObj: NewDocObj(ADocument),
		P2:     Coord{x: X2, y: Y2},
		P3:     Coord{x: X3, y: Y3},
		Width:  AWidth,
		Stroke: AStroke,
	}
}

func NewTPDFCurveV(ADocument *Document, AP2, AP3 Coord, AWidth float64, AStroke bool) *CurveV {
	return &CurveV{
		P2:     AP2,
		P3:     AP3,
		Width:  AWidth,
		Stroke: AStroke,
	}
}

type LineSegment struct {
	DocObj
	Width  float64
	Stroke bool
	P1, P2 Coord
}

func (ls *LineSegment) Encode(st PDFWriter) {
	ls.SetWidth(ls.Width, st)
	if ls.Stroke {
		st.WriteString(MoveTo{}.Command(ls.P1)) 
	}
	st.WriteString(ls.Command(ls.P2)) 
	if ls.Stroke {
		st.WriteString("S" + CRLF) 
	}
}

func (ls *LineSegment) Command(APos Coord) string {
	return FloatStr(APos.x) + " " + FloatStr(APos.y) + " l" + CRLF
}

func (ls *LineSegment) CommandXY(x1, y1 float64) string {
	return FloatStr(x1) + " " + FloatStr(y1) + " l" + CRLF
}

func (ls *LineSegment) CommandCoords(APos1, APos2 Coord) string {
	return MoveTo{}.Command(APos1) + ls.Command(APos2)
}

func NewLineSegment(ADocument *Document, AWidth, X1, Y1, X2, Y2 float64, AStroke bool) *LineSegment {
	return &LineSegment{
		DocObj: NewDocObj(ADocument),
		Width:  AWidth,
		P1:     Coord{x: X1, y: Y1},
		P2:     Coord{x: X2, y: Y2},
		Stroke: AStroke,
	}
}

type Rectangle struct {
	DocObj
	Width      float64
	TopLeft    Coord
	Dimensions Coord
	Fill       bool
	Stroke     bool
}

func (r *Rectangle) Encode(st PDFWriter) {
	if r.Stroke {
		r.SetWidth(r.Width, st)
	}
	st.WriteString(FloatStr(r.TopLeft.x) + " " + FloatStr(r.TopLeft.y) + " " + FloatStr(r.Dimensions.x) + " " + FloatStr(r.Dimensions.y) + " re" + CRLF) 
	if r.Stroke && r.Fill {
		st.WriteString("b" + CRLF) 
	} else if r.Fill {
		st.WriteString("f" + CRLF) 
	} else if r.Stroke {
		st.WriteString("S" + CRLF) 
	}
}

func NewRectangle(ADocument *Document, APosX, APosY, AWidth, AHeight, ALineWidth float64, AFill, AStroke bool) *Rectangle {
	return &Rectangle{
		DocObj:     NewDocObj(ADocument),
		TopLeft:    Coord{x: APosX, y: APosY},
		Dimensions: Coord{x: AWidth, y: AHeight},
		Width:      ALineWidth,
		Fill:       AFill,
		Stroke:     AStroke,
	}
}


type RoundedRectangle struct {
	DocObj
	Width      float64
	BottomLeft Coord
	Dimensions Coord
	Fill       bool
	Stroke     bool
	Radius     float64
}

func NewRoundedRectangle(document *Document, posX, posY, width, height, radius, lineWidth float64, fill, stroke bool) *RoundedRectangle {
	return &RoundedRectangle{
		DocObj:     NewDocObj(document),
		BottomLeft: Coord{x: posX, y: posY},
		Dimensions: Coord{x: width, y: height},
		Width:      lineWidth,
		Fill:       fill,
		Stroke:     stroke,
		Radius:     radius,
	}
}

func (rr *RoundedRectangle) Encode(st PDFWriter) {}

type Surface struct {
	DocObj
	Points []Coord
	Close  bool
	Fill   bool
}

func NewSurface(document *Document, points []Coord, close, fill bool) *Surface {
	return &Surface{
		DocObj: NewDocObj(document),
		Points: points,
		Close:  close,
		Fill:   fill,
	}
}
func (s *Surface) Encode(st PDFWriter) {
	st.WriteString(MoveTo{}.CommandXY(s.Points[0].x, s.Points[0].y))
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

type LineStyle struct {
	DocObj
	Style     PenStyle
	Phase     int
	LineWidth float64
	LineMask  string
}

func NewPDFLineStyle(document *Document, style PenStyle, phase int, lineWidth float64) *LineStyle {
	return &LineStyle{
		DocObj:    NewDocObj(document),
		Style:     style,
		Phase:     phase,
		LineWidth: lineWidth,
	}
}
func (ls *LineStyle) Encode(st PDFWriter) {
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

func NewLineStyleCustom(document *Document, dashArray []float64, phase int, lineWidth float64) *LineStyle {
	ls := NewPDFLineStyle(document, ppsSolid, phase, lineWidth)
	for _, d := range dashArray {
		if ls.LineMask != "" {
			ls.LineMask += " "
		}
		ls.LineMask += fmt.Sprintf("%f", d*lineWidth)
	}
	return ls
}

type CapStyle struct {
	Style LineCapStyle
}


func NewCapStyle(document *Document, style LineCapStyle) *CapStyle {
	return &CapStyle{
		Style: style,
	}
}

func (cs *CapStyle) Encode(st PDFWriter) {
	st.WriteString(fmt.Sprintf("%d J%s", cs.Style, CRLF))
}

type JoinStyle struct {
	Style LineJoinStyle
}

func NewPDFJoinStyle(document *Document, style LineJoinStyle) *JoinStyle {
	return &JoinStyle{
		Style: style,
	}
}

func (js *JoinStyle) Encode(st PDFWriter) {
	st.WriteString(fmt.Sprintf("%d j%s", js.Style, CRLF))
}

type LineStyleDef struct {
	FColor     ARGBColor
	FLineWidth float64
	FPenStyle  PenStyle
	FDashArray DashArray
}
