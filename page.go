package main

import "fmt"

type Page struct {
	DocObj
	Objects       []Encoder
	Orientation   PaperOrientation
	Paper         Paper
	PaperType     PaperType
	UnitOfMeasure UnitOfMeasure
	Matrix        Matrix
	//FIXME: re-enable
	// Annots            TPDFAnnotList
	LastFont      *EmbeddedFont
	LastFontColor ARGBColor
}

func NewPage(doc *Document) *Page {
	p := &Page{DocObj: NewDocObj(doc)}
	if doc != nil {
		p.PaperType = doc.DefaultPaperType
		p.Orientation = doc.DefaultOrientation
		p.UnitOfMeasure = doc.DefaultUnitOfMeasure
	} else {
		p.PaperType = ptA4
		p.UnitOfMeasure = uomMillimeters
	}
	p.CalcPaperSize()
	p.Matrix._00 = 1
	p.Matrix._20 = 0
	p.AdjustMatrix()
	// p.Annots = p.CreateAnnotList()
	return p
}

// FIXME:
func (p *Page) Encode(st PDFWriter) {}

func (p *Page) AddObject(obj Encoder) {
	p.Objects = append(p.Objects, obj)
}
func (p *Page) GetObject(AIndex int) Encoder {
	return p.Objects[AIndex]
}

func (p *Page) GetObjectCount() int {
	return len(p.Objects)
}

// func (p *TPDFPage) CreateAnnotList() *TPDFAnnotList {
// 	return NewTPDFAnnotList(p.FDocument)
// }

func (p *Page) SetOrientation(val PaperOrientation) {
	if p.Orientation == val {
		return
	}
	p.Orientation = val
	p.CalcPaperSize()
	p.AdjustMatrix()
}

func (p *Page) SetPaper(val Paper) {
	if p.Paper == val {
		return
	}
	p.Paper = val
	p.PaperType = ptCustom
	p.AdjustMatrix()
}

func (p *Page) CalcPaperSize() {
	if p.PaperType == ptCustom {
		return
	}
	O1, O2 := 0, 1
	if p.Orientation == ppoLandscape {
		O1, O2 = 1, 0
	}
	p.Paper = Paper{
		H: PDFPaperDims[p.PaperType][O1],
		W: PDFPaperDims[p.PaperType][O2],
		Printable: Dimensions{
			T: float64(PDFPaperDims[p.PaperType][2+O1]),
			L: float64(PDFPaperDims[p.PaperType][2+O2]),
			R: float64(PDFPaperDims[p.PaperType][2+2+O1]),
			B: float64(PDFPaperDims[p.PaperType][2+2+O2]),
		},
	}
}

func (p *Page) SetPaperType(val PaperType) {
	if p.PaperType == val {
		return
	}
	p.PaperType = val
	if p.PaperType != ptCustom {
		p.CalcPaperSize()
	}
	p.AdjustMatrix()
}

func (p *Page) AddTextToLookupLists(AText string) {
	//FIXME:
	// if AText == "" {
	// 	return
	// }
	// str := utf8.DecodeString(AText)
	// p.FDocument.Fonts[p.FLastFont.FontIndex].AddTextToMappingList(str)
}

func (p *Page) SetUnitOfMeasure(val UnitOfMeasure) {
	if p.UnitOfMeasure == val {
		return
	}
	p.UnitOfMeasure = val
	p.AdjustMatrix()
}

func (p *Page) AdjustMatrix() {
	if p.Document.hasOption(poPageOriginAtTop) {
		p.Matrix._11 = -1
		p.Matrix._21 = p.GetPaperHeight()
	} else {
		p.Matrix._11 = 1
		p.Matrix._21 = 0
	}
}

func (p *Page) DoUnitConversion(coord *Coord) { DoUnitConversion(coord, p.UnitOfMeasure) }

func (p *Page) SetFont(fontIdx int, fontSz float64) {
	p.SetFontEx(fontIdx, fontSz, false, false)
}

func (p *Page) SetFontEx(fontIdx int, fontSz float64, ASimulateBold, ASimulateItalic bool) {
	p.LastFont = p.Document.CreateEmbeddedFont(p, fontIdx, fontSz, ASimulateBold, ASimulateItalic)
	p.AddObject(p.LastFont)
}

func (p *Page) SetColor(AColor ARGBColor, AStroke bool) {
	C := NewColor(p.Document, AColor, AStroke)
	if !AStroke {
		p.LastFontColor = AColor
	}
	p.AddObject(C)
}

//FIXME:
// func (p *TPDFPage) SetPenStyle(AStyle TPDFPenStyle, ALineWidth float64) {
// 	L := p.Document.CreateLineStyle(AStyle, ALineWidth)
// 	p.AddObject(L)
// }

// func (p *TPDFPage) SetPenStyle(ADashArray TDashArray, ALineWidth float64) {
// 	L := p.Document.CreateLineStyleFromDashArray(ADashArray, ALineWidth)
// 	p.AddObject(L)
// }

// func (p *TPDFPage) SetLineCapStyle(AStyle TPDFLineCapStyle) {
// 	p.Document.LineCapStyle = AStyle
// 	C := p.Document.CreateLineCapStyle(AStyle)
// 	p.AddObject(C)
// }

// func (p *TPDFPage) SetLineJoinStyle(AStyle TPDFLineJoinStyle) {
// 	J := p.Document.CreateLineJoinStyle(AStyle)
// 	p.AddObject(J)
// }

// // func (p *TPDFPage) SetMiterLimit(AMiterLimit float64) {
// // 	M := p.Document.CreateMiterLimit(AMiterLimit)
// // 	p.AddObject(M)
// // }

// func (p *TPDFPage) SetLineStyle(AIndex int, AStroke bool) {
// 	p.SetLineStyleFromDef(p.Document.LineStyles[AIndex], AStroke)
// }

// func (p *TPDFPage) SetLineStyleFromDef(ALineStyle *TPDFLineStyle, AStroke bool) {
// 	if ALineStyle != nil {
// 		if ALineStyle.DashArray == nil {
// 			p.SetPenStyle(ALineStyle.Style, ALineStyle.LineWidth)
// 		} else {
// 			p.SetPenStyle(*ALineStyle.DashArray, ALineStyle.LineWidth)
// 		}
// 		p.SetLineCapStyle(ALineStyle.LineCap)
// 		p.SetLineJoinStyle(ALineStyle.LineJoin)
// 		p.SetMiterLimit(ALineStyle.MiterLimit)
// 		if AStroke {
// 			p.SetColor(ALineStyle.LineColor, true)
// 		} else {
// 			p.SetColor(ALineStyle.FillColor, false)
// 		}
// 	}
// }

// func (p *TPDFPage) SetFillStyle(AIndex int) {
// 	p.SetFillStyleFromDef(p.Document.FillStyles[AIndex])
// }

// func (p *TPDFPage) SetFillStyleFromDef(AFillStyle *TPDFFillStyle) {
// 	if AFillStyle != nil {
// 		if AFillStyle.Pattern == nil {
// 			p.SetColor(AFillStyle.Color, false)
// 		} else {
// 			p.Document.PatternStyle = *AFillStyle.Pattern
// 			P := p.Document.CreatePattern(*AFillStyle.Pattern)
// 			p.AddObject(P)
// 		}
// 	}
// }

// func (p *TPDFPage) DrawImageRawSize(x, Y float64, APixelWidth, APixelHeight, ANumber int, degs float64) {
// 	p1 := p.Matrix.TransformXY(x, Y)
// 	p.DoUnitConversion(&p1)

// 	if degs != 0.0 {
// p.AddObject(NewPDFFreeFormString(p.Document, fmt.Sprintf("%s %s %s %s %.4f %.4f cm", t1, t2, t3, t1, p1.x, p1.y)))
// 		p.AddObject(p.FDocument.CreateImage(0, 0, APixelWidth, APixelHeight, ANumber))
// 	} else {
// 		p.AddObject(p.FDocument.CreateImage(p1.x, p1.y, APixelWidth, APixelHeight, ANumber))
// 	}

// 	if degs != 0.0 {
// 		p.AddObject(NewTPDFPopGraphicsStack(p.FDocument))
// 	}
// }

// func (p *TPDFPage) DrawImageRawSizePos(pos TPDFCoord, APixelWidth, APixelHeight, ANumber int, degs float64) {
// 	p.DrawImageRawSize(pos.x, pos.y, APixelWidth, APixelHeight, ANumber, degs)
// }

// func (p *TPDFPage) DrawImage(x, Y, AWidth, AHeight float64, ANumber int, degs float64) {
// 	p1 := p.Matrix.TransformXY(x, Y)
// 	p.DoUnitConversion(&p1)
// 	p2 := TPDFCoord{AWidth, AHeight}
// 	p.DoUnitConversion(&p2)

// 	if degs != 0.0 {
// 		rad := DegToRad(-degs)
// 		rads, radc := sincos(rad)
// 		t1 := fmt.Sprintf(PDF_NUMBER_MASK, radc)
// 		t2 := fmt.Sprintf(PDF_NUMBER_MASK, -rads)
// 		t3 := fmt.Sprintf(PDF_NUMBER_MASK, rads)

// 		p.AddObject(NewTPDFPushGraphicsStack(p.FDocument))
// 		p.AddObject(NewTPDFFreeFormString(p.FDocument, fmt.Sprintf("%s %s %s %s %.4f %.4f cm", t1, t2, t3, t1, p1.x, p1.y)))
// 		p.AddObject(p.FDocument.CreateImage(0, 0, p2.x, p2.y, ANumber))
// 	} else {
// 		p.AddObject(p.FDocument.CreateImage(p1.x, p1.y, p2.x, p2.y, ANumber))
// 	}

// 	if degs != 0.0 {
// 		p.AddObject(NewTPDFPopGraphicsStack(p.FDocument))
// 	}
// }

// func (p *TPDFPage) DrawImagePos(pos TPDFCoord, AWidth, AHeight float64, ANumber int, degs float64) {
// 	p.DrawImage(pos.x, pos.y, AWidth, AHeight, ANumber, degs)
// }

func (p *Page) addAngle(degrees float64, p1 Coord) {
	rad := DegToRad(-degrees)
	rads, radc := sincos(rad)
	//FIXME: is per-local formatting needed
	// t1 := fmt.Sprintf(PDF_NUMBER_MASK, radc)
	// t2 := fmt.Sprintf(PDF_NUMBER_MASK, -rads)
	// t3 := fmt.Sprintf(PDF_NUMBER_MASK, rads)
	p.AddObject(&PushGraphicsStack{})
	//FIXME: this when per-local formatting was used (using tn ); I replaced
	// it with simple %0.4f
	// p.AddObject(NewTPDFFreeFormString(p.Document, fmt.Sprintf("%s %s %s %s %.4f %.4f cm", t1, t2, t3, t1, p1.x, p1.y)))
	p.AddObject(NewFreeFormString(p.Document, fmt.Sprintf("%.4f %.4f %.4f %.4f %.4f %.4f cm", radc, -rads, rads, radc, p1.x, p1.y)))
}

func (p *Page) DrawEllipse(posx, posY, AWidth, AHeight, ALineWidth float64, AFill, AStroke bool, degs float64) {
	p1 := p.Matrix.TransformXY(posx, posY)
	p.DoUnitConversion(&p1)
	p2 := Coord{AWidth, AHeight}
	p.DoUnitConversion(&p2)

	if degs != 0.0 {
		p.addAngle(degs, p1)
		p.AddObject(NewEllipse(p.Document, 0, 0, p2.x, p2.y, ALineWidth, AFill, AStroke))
	} else {
		p.AddObject(NewEllipse(p.Document, p1.x, p1.y, p2.x, p2.y, ALineWidth, AFill, AStroke))
	}

	if degs != 0.0 {
		p.AddObject(NewTPDFPopGraphicsStack(p.Document))
	}
}

func (p *Page) DrawEllipsePos(pos Coord, AWidth, AHeight, ALineWidth float64, AFill, AStroke bool, degs float64) {
	p.DrawEllipse(pos.x, pos.y, AWidth, AHeight, ALineWidth, AFill, AStroke, degs)
}

// func (p *TPDFPage) DrawLine(pos1, pos2 PDFCoord, ALineWidth float64, AStroke bool) {
// 	p.DrawLine(pos1.x, pos1.y, pos2.x, pos2.y, ALineWidth, AStroke)
// }

// func (p *TPDFPage) DrawLineStyle(X1, Y1, X2, Y2 float64, AStyle int) {
// 	S := p.Document.LineStyles[AStyle]
// 	p.SetLineStyle(S)
// 	p.DrawLine(X1, Y1, X2, Y2, S.LineWidth)
// }

//	func (p *TPDFPage) DrawLineStyleCoord(pos1, pos2 PDFCoord, AStyle int) {
//		p.DrawLineStyle(pos1.x, pos1.y, pos2.x, pos2.y, AStyle)
//	}
func (p *Page) DrawRect(x, Y, W, H, ALineWidth float64, AFill, AStroke bool, degs float64) {
	p1 := p.Matrix.TransformXY(x, Y)
	p.DoUnitConversion(&p1)
	p2 := Coord{x: W, y: H}
	p.DoUnitConversion(&p2)

	var R *Rectangle
	if degs != 0.0 {
		p.addAngle(degs, p1)
		R = NewRectangle(p.Document, 0, 0, p2.x, p2.y, ALineWidth, AFill, AStroke)
	} else {
		R = NewRectangle(p.Document, p1.x, p1.y, p2.x, p2.y, ALineWidth, AFill, AStroke)
	}

	p.AddObject(R)

	if degs != 0.0 {
		p.AddObject(NewTPDFPopGraphicsStack(p.Document))
	}
}

func (p *Page) DrawRectCoord(pos Coord, W, H, ALineWidth float64, AFill, AStroke bool, degs float64) {
	p.DrawRect(pos.x, pos.y, W, H, ALineWidth, AFill, AStroke, degs)
}

func (p *Page) DrawRoundedRect(x, Y, W, H, ARadius, ALineWidth float64, AFill, AStroke bool, degs float64) {
	p1 := p.Matrix.TransformXY(x, Y)
	p.DoUnitConversion(&p1)
	p2 := Coord{x: W, y: H}
	p.DoUnitConversion(&p2)
	p3 := Coord{x: ARadius}
	p.DoUnitConversion(&p3)

	var R *RoundedRectangle
	if degs != 0.0 {
		p.addAngle(degs, p1)
		R = NewRoundedRectangle(p.Document, 0, 0, p2.x, p2.y, p3.x, ALineWidth, AFill, AStroke)
	} else {
		R = NewRoundedRectangle(p.Document, p1.x, p1.y, p2.x, p2.y, p3.x, ALineWidth, AFill, AStroke)
	}

	p.AddObject(R)

	if degs != 0.0 {
		p.AddObject(NewTPDFPopGraphicsStack(p.Document))
	}
}

func (p *Page) DrawPolygon(APoints []Coord, ALineWidth float64) {
	p.DrawPolyLine(APoints, ALineWidth)
	p.ClosePath()
}

func (p *Page) DrawPolyLine(APoints []Coord, ALineWidth float64) {
	if len(APoints) < 2 {
		return
	}
	p.MoveTo(APoints[0].x, APoints[0].y)
	for i := 1; i < len(APoints); i++ {
		//FIXME:
		// p.DrawLine(APoints[i-1].x, APoints[i-1].y, APoints[i].x, APoints[i].y, ALineWidth, false)
	}
}

func (p *Page) ResetPath()             { p.AddObject(&ResetPath{}) }
func (p *Page) ClipPath()              { p.AddObject(&ClipPath{}) }
func (p *Page) ClosePath()             { p.AddObject(&ClosePath{}) }
func (p *Page) StrokePath()            { p.AddObject(&StrokePath{}) }
func (p *Page) ClosePathStroke()       { p.AddObject(NewFreeFormString(p.Document, "s\n")) }
func (p *Page) FillStrokePath()        { p.AddObject(NewFreeFormString(p.Document, "B\n")) }
func (p *Page) FillEvenOddStrokePath() { p.AddObject(NewFreeFormString(p.Document, "B*\n")) }
func (p *Page) PushGraphicsStack()     { p.AddObject(&PushGraphicsStack{}) }
func (p *Page) PopGraphicsStack()      { p.AddObject(&PopGraphicsStack{}) }
func (p *Page) MoveToPoint(pos Coord)  { p.MoveTo(pos.x, pos.y) }

func (p *Page) MoveTo(x, y float64) {
	p1 := p.Matrix.TransformXY(x, y)
	DoUnitConversion(&p1, p.UnitOfMeasure)
	p.AddObject(&MoveTo{p1})
}

func (p *Page) GetPaperHeight() float64 {
	switch p.UnitOfMeasure {
	case uomMillimeters:
		return PDFTomm(float64(p.Paper.H))
	case uomCentimeters:
		return PDFtoCM(float64(p.Paper.H))
	case uomInches:
		return PDFtoInches(float64(p.Paper.H))
	case uomPixels:
		return float64(p.Paper.H)
	default:
		return 0
	}
}

// func (p *TPDFPage) HasImages() bool {
// 	for _, obj := range p.Objects {
// 		if _, ok := obj.(TPDFImage); ok {
// 			return true
// 		}
// 	}
// 	return false
// }

func (p *Page) CreateStdFontText(x, Y float64, AText string, AFont *EmbeddedFont, degs float64, AUnderline, AStrikethrough bool) {
	p.AddObject(NewText(p.Document, x, Y, AText, AFont, degs, AUnderline, AStrikethrough))
}

// func (p *TPDFPage) CreateTTFFontText(x, Y float64, AText string, AFont *TPDFEmbeddedFont, degs float64, AUnderline, AStrikethrough bool) {
// 	p.AddTextToLookupLists(AText)
// 	T := p.Document.CreateText(x, Y, AText, AFont, degs, AUnderline, AStrikethrough)
// 	p.AddObject(T)
// }

// procedure TPDFPage.WriteText(X, Y: Tfloat64; AText: UTF8String; const degs: single;
//     const AUnderline: boolean; const AStrikethrough: boolean);
// var
//   p: TPDFCoord;
// begin
//   if not Assigned(FLastFont) then
//     raise EPDF.Create(rsErrNoFontDefined);
//   p := Matrix.Transform(X, Y);
//   DoUnitConversion(p);
//   if Document.Fonts[FLastFont.FontIndex].IsStdFont then
//     CreateStdFontText(p.X, p.Y, AText, FLastFont, degs, AUnderline, AStrikeThrough)
//   else
//     CreateTTFFontText(p.X, p.Y, AText, FLastFont, degs, AUnderline, AStrikeThrough);
// end;

// procedure TPDFPage.WriteText(pos: TPDFCoord; AText: UTF8String; const degs: single;
//
//	const AUnderline: boolean; const AStrikethrough: boolean);
//
// begin
//
//	WriteText(pos.X, pos.Y, AText, degs, AUnderline, AStrikeThrough);
//
// end;
// FIXME: complete implementation
func (p *Page) WriteTextEx(x, y float64, txt string, degs float64, underline, strikethro bool) {
	p1 := p.Matrix.TransformXY(x, y)
	p.DoUnitConversion(&p1)
	p.CreateStdFontText(p1.x, p1.y, txt, p.LastFont, degs, underline, strikethro)
}

func (p *Page) WriteText(x, y float64, txt string) {
	p.WriteTextEx(x, y, txt, 0, false, false)
}
