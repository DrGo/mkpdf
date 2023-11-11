package main

import "fmt"

type TPDFPage struct {
	TPDFDocumentObject
	Objects       []Encoder
	Orientation   TPDFPaperOrientation
	Paper         TPDFPaper
	PaperType     TPDFPaperType
	UnitOfMeasure TPDFUnitOfMeasure
	Matrix        TPDFMatrix
	//FIXME: re-enable
	// Annots            TPDFAnnotList
	LastFont      *TPDFEmbeddedFont
	LastFontColor ARGBColor
}

func NewTPDFPage(ADocument *TPDFDocument) *TPDFPage {
	p := &TPDFPage{TPDFDocumentObject: NewTPDFDocumentObject(ADocument)}
	if ADocument != nil {
		p.PaperType = ADocument.DefaultPaperType
		p.Orientation = ADocument.DefaultOrientation
		p.UnitOfMeasure = ADocument.DefaultUnitOfMeasure
	} else {
		p.PaperType = ptA4
		p.CalcPaperSize()
		p.UnitOfMeasure = uomMillimeters
	}
	p.Matrix._00 = 1
	p.Matrix._20 = 0
	p.AdjustMatrix()
	// p.Annots = p.CreateAnnotList()
	return p
}
func (p *TPDFPage) AddObject(AObject Encoder) {
	p.Objects = append(p.Objects, AObject)
}
func (p *TPDFPage) GetO(AIndex int) Encoder {
	return p.Objects[AIndex]
}

func (p *TPDFPage) GetObjectCount() int {
	return len(p.Objects)
}

// func (p *TPDFPage) CreateAnnotList() *TPDFAnnotList {
// 	return NewTPDFAnnotList(p.FDocument)
// }

func (p *TPDFPage) SetOrientation(AValue TPDFPaperOrientation) {
	if p.Orientation == AValue {
		return
	}
	p.Orientation = AValue
	p.CalcPaperSize()
	p.AdjustMatrix()
}

func (p *TPDFPage) SetPaper(AValue TPDFPaper) {
	if p.Paper == AValue {
		return
	}
	p.Paper = AValue
	p.PaperType = ptCustom
	p.AdjustMatrix()
}

func (p *TPDFPage) CalcPaperSize() {
	if p.PaperType == ptCustom {
		return
	}
	O1, O2 := 0, 1
	if p.Orientation == ppoLandscape {
		O1, O2 = 1, 0
	}
	p.Paper = TPDFPaper{
		H: PDFPaperDims[p.PaperType][O1],
		W: PDFPaperDims[p.PaperType][O2],
		Printable: TPDFDimensions{
			T: PDFPaperDims[p.PaperType][2+O1],
			L: PDFPaperDims[p.PaperType][2+O2],
			R: PDFPaperDims[p.PaperType][2+2+O1],
			B: PDFPaperDims[p.PaperType][2+2+O2],
		},
	}
}

func (p *TPDFPage) SetPaperType(AValue TPDFPaperType) {
	if p.PaperType == AValue {
		return
	}
	p.PaperType = AValue
	if p.PaperType != ptCustom {
		p.CalcPaperSize()
	}
	p.AdjustMatrix()
}

func (p *TPDFPage) AddTextToLookupLists(AText string) {
	//FIXME:
	// if AText == "" {
	// 	return
	// }
	// str := utf8.DecodeString(AText)
	// p.FDocument.Fonts[p.FLastFont.FontIndex].AddTextToMappingList(str)
}

func (p *TPDFPage) SetUnitOfMeasure(AValue TPDFUnitOfMeasure) {
	if p.UnitOfMeasure == AValue {
		return
	}
	p.UnitOfMeasure = AValue
	p.AdjustMatrix()
}

func (p *TPDFPage) AdjustMatrix() {
	if p.Document.hasOption(poPageOriginAtTop) {
		p.Matrix._11 = -1
		p.Matrix._21 = p.GetPaperHeight()
	} else {
		p.Matrix._11 = 1
		p.Matrix._21 = 0
	}
}

func (p *TPDFPage) DoUnitConversion(coord *TPDFCoord) {
	DoUnitConversion(coord, p.UnitOfMeasure)
}

func (p *TPDFPage) SetFont(AFontIndex int, AFontSize PDFFloat, ASimulateBold, ASimulateItalic bool) {
	p.LastFont = p.Document.CreateEmbeddedFont(p, AFontIndex, AFontSize, ASimulateBold, ASimulateItalic)
	p.AddObject(p.LastFont)
}

func (p *TPDFPage) SetColor(AColor ARGBColor, AStroke bool) {
	C := NewTPDFColor(p.Document, AColor, AStroke)
	if !AStroke {
		p.LastFontColor = AColor
	}
	p.AddObject(C)
}

//FIXME:
// func (p *TPDFPage) SetPenStyle(AStyle TPDFPenStyle, ALineWidth PDFFloat) {
// 	L := p.Document.CreateLineStyle(AStyle, ALineWidth)
// 	p.AddObject(L)
// }

// func (p *TPDFPage) SetPenStyle(ADashArray TDashArray, ALineWidth PDFFloat) {
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

// // func (p *TPDFPage) SetMiterLimit(AMiterLimit PDFFloat) {
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

// func (p *TPDFPage) DrawImageRawSize(x, Y PDFFloat, APixelWidth, APixelHeight, ANumber int, ADegrees float32) {
// 	p1 := p.Matrix.TransformXY(x, Y)
// 	p.DoUnitConversion(&p1)

// 	if ADegrees != 0.0 {
// p.AddObject(NewPDFFreeFormString(p.Document, fmt.Sprintf("%s %s %s %s %.4f %.4f cm", t1, t2, t3, t1, p1.x, p1.y)))
// 		p.AddObject(p.FDocument.CreateImage(0, 0, APixelWidth, APixelHeight, ANumber))
// 	} else {
// 		p.AddObject(p.FDocument.CreateImage(p1.x, p1.y, APixelWidth, APixelHeight, ANumber))
// 	}

// 	if ADegrees != 0.0 {
// 		p.AddObject(NewTPDFPopGraphicsStack(p.FDocument))
// 	}
// }

// func (p *TPDFPage) DrawImageRawSizePos(APos TPDFCoord, APixelWidth, APixelHeight, ANumber int, ADegrees float32) {
// 	p.DrawImageRawSize(APos.x, APos.y, APixelWidth, APixelHeight, ANumber, ADegrees)
// }

// func (p *TPDFPage) DrawImage(x, Y, AWidth, AHeight PDFFloat, ANumber int, ADegrees float32) {
// 	p1 := p.Matrix.TransformXY(x, Y)
// 	p.DoUnitConversion(&p1)
// 	p2 := TPDFCoord{AWidth, AHeight}
// 	p.DoUnitConversion(&p2)

// 	if ADegrees != 0.0 {
// 		rad := DegToRad(-ADegrees)
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

// 	if ADegrees != 0.0 {
// 		p.AddObject(NewTPDFPopGraphicsStack(p.FDocument))
// 	}
// }

// func (p *TPDFPage) DrawImagePos(APos TPDFCoord, AWidth, AHeight PDFFloat, ANumber int, ADegrees float32) {
// 	p.DrawImage(APos.x, APos.y, AWidth, AHeight, ANumber, ADegrees)
// }

func (p *TPDFPage) addAngle(degrees float32, p1 TPDFCoord) {
	rad := DegToRad(-degrees)
	rads, radc := sincos(rad)
	//FIXME: is per-local formatting needed
	// t1 := fmt.Sprintf(PDF_NUMBER_MASK, radc)
	// t2 := fmt.Sprintf(PDF_NUMBER_MASK, -rads)
	// t3 := fmt.Sprintf(PDF_NUMBER_MASK, rads)
	p.AddObject(&TPDFPushGraphicsStack{})
	//FIXME: this when per-local formatting was used (using tn ); I replaced
	// it with simple %0.4f
	// p.AddObject(NewTPDFFreeFormString(p.Document, fmt.Sprintf("%s %s %s %s %.4f %.4f cm", t1, t2, t3, t1, p1.x, p1.y)))
	p.AddObject(NewTPDFFreeFormString(p.Document, fmt.Sprintf("%.4f %.4f %.4f %.4f %.4f %.4f cm", radc, -rads, rads, radc, p1.x, p1.y)))
}

func (p *TPDFPage) DrawEllipse(APosx, APosY, AWidth, AHeight, ALineWidth PDFFloat, AFill, AStroke bool, ADegrees float32) {
	p1 := p.Matrix.TransformXY(APosx, APosY)
	p.DoUnitConversion(&p1)
	p2 := TPDFCoord{AWidth, AHeight}
	p.DoUnitConversion(&p2)

	if ADegrees != 0.0 {
		p.addAngle(ADegrees, p1)
		p.AddObject(NewTPDFEllipse(p.Document, 0, 0, p2.x, p2.y, ALineWidth, AFill, AStroke))
	} else {
		p.AddObject(NewTPDFEllipse(p.Document, p1.x, p1.y, p2.x, p2.y, ALineWidth, AFill, AStroke))
	}

	if ADegrees != 0.0 {
		p.AddObject(NewTPDFPopGraphicsStack(p.Document))
	}
}

func (p *TPDFPage) DrawEllipsePos(APos TPDFCoord, AWidth, AHeight, ALineWidth PDFFloat, AFill, AStroke bool, ADegrees float32) {
	p.DrawEllipse(APos.x, APos.y, AWidth, AHeight, ALineWidth, AFill, AStroke, ADegrees)
}

// func (p *TPDFPage) DrawLine(APos1, APos2 PDFCoord, ALineWidth PDFFloat, AStroke bool) {
// 	p.DrawLine(APos1.x, APos1.y, APos2.x, APos2.y, ALineWidth, AStroke)
// }

// func (p *TPDFPage) DrawLineStyle(X1, Y1, X2, Y2 PDFFloat, AStyle int) {
// 	S := p.Document.LineStyles[AStyle]
// 	p.SetLineStyle(S)
// 	p.DrawLine(X1, Y1, X2, Y2, S.LineWidth)
// }

//	func (p *TPDFPage) DrawLineStyleCoord(APos1, APos2 PDFCoord, AStyle int) {
//		p.DrawLineStyle(APos1.x, APos1.y, APos2.x, APos2.y, AStyle)
//	}
func (p *TPDFPage) DrawRect(x, Y, W, H, ALineWidth PDFFloat, AFill, AStroke bool, ADegrees float32) {
	p1 := p.Matrix.TransformXY(x, Y)
	p.DoUnitConversion(&p1)
	p2 := TPDFCoord{x: W, y: H}
	p.DoUnitConversion(&p2)

	var R *TPDFRectangle
	if ADegrees != 0.0 {
		p.addAngle(ADegrees, p1)
		R = NewTPDFRectangle(p.Document, 0, 0, p2.x, p2.y, ALineWidth, AFill, AStroke)
	} else {
		R = NewTPDFRectangle(p.Document, p1.x, p1.y, p2.x, p2.y, ALineWidth, AFill, AStroke)
	}

	p.AddObject(R)

	if ADegrees != 0.0 {
		p.AddObject(NewTPDFPopGraphicsStack(p.Document))
	}
}

func (p *TPDFPage) DrawRectCoord(APos TPDFCoord, W, H, ALineWidth PDFFloat, AFill, AStroke bool, ADegrees float32) {
	p.DrawRect(APos.x, APos.y, W, H, ALineWidth, AFill, AStroke, ADegrees)
}

func (p *TPDFPage) DrawRoundedRect(x, Y, W, H, ARadius, ALineWidth PDFFloat, AFill, AStroke bool, ADegrees float32) {
	p1 := p.Matrix.TransformXY(x, Y)
	p.DoUnitConversion(&p1)
	p2 := TPDFCoord{x: W, y: H}
	p.DoUnitConversion(&p2)
	p3 := TPDFCoord{x: ARadius}
	p.DoUnitConversion(&p3)

	var R *PDFRoundedRectangle
	if ADegrees != 0.0 {
		p.addAngle(ADegrees, p1)
		R = NewPDFRoundedRectangle(p.Document, 0, 0, p2.x, p2.y, p3.x, ALineWidth, AFill, AStroke)
	} else {
		R = NewPDFRoundedRectangle(p.Document, p1.x, p1.y, p2.x, p2.y, p3.x, ALineWidth, AFill, AStroke)
	}

	p.AddObject(R)

	if ADegrees != 0.0 {
		p.AddObject(NewTPDFPopGraphicsStack(p.Document))
	}
}

func (p *TPDFPage) DrawPolygon(APoints []TPDFCoord, ALineWidth PDFFloat) {
	p.DrawPolyLine(APoints, ALineWidth)
	p.ClosePath()
}

func (p *TPDFPage) DrawPolyLine(APoints []TPDFCoord, ALineWidth PDFFloat) {
	if len(APoints) < 2 {
		return
	}
	p.MoveTo(APoints[0].x, APoints[0].y)
	for i := 1; i < len(APoints); i++ {
		//FIXME:
		// p.DrawLine(APoints[i-1].x, APoints[i-1].y, APoints[i].x, APoints[i].y, ALineWidth, false)
	}
}

func (p *TPDFPage) ResetPath() {
	p.AddObject(&TPDFResetPath{})
}

func (p *TPDFPage) ClipPath() {
	p.AddObject(&TPDFClipPath{})
}

func (p *TPDFPage) ClosePath() {
	p.AddObject(&TPDFClosePath{})
}

func (p *TPDFPage) StrokePath() {
	p.AddObject(&TPDFStrokePath{})
}
func (p *TPDFPage) ClosePathStroke() {
	p.AddObject(NewTPDFFreeFormString(p.Document, "s\n"))
}

func (p *TPDFPage) FillStrokePath() {
	p.AddObject(NewTPDFFreeFormString(p.Document, "B\n"))
}

func (p *TPDFPage) FillEvenOddStrokePath() {
	p.AddObject(NewTPDFFreeFormString(p.Document, "B*\n"))
}

func (p *TPDFPage) PushGraphicsStack() {
	p.AddObject(&TPDFPushGraphicsStack{})
}

func (p *TPDFPage) PopGraphicsStack() {
	p.AddObject(&TPDFPopGraphicsStack{})
}

func (p *TPDFPage) MoveTo(x, y PDFFloat) {
	p1 := p.Matrix.TransformXY(x, y)
	DoUnitConversion(&p1, p.UnitOfMeasure)
	p.AddObject(&TPDFMoveTo{p1})
}

func (p *TPDFPage) MoveToPoint(pos TPDFCoord) {
	p.MoveTo(pos.x, pos.y)
}

func (p *TPDFPage) GetPaperHeight() PDFFloat {
	switch p.UnitOfMeasure {
	case uomMillimeters:
		return PDFTomm(PDFFloat(p.Paper.H))
	case uomCentimeters:
		return PDFtoCM(PDFFloat(p.Paper.H))
	case uomInches:
		return PDFtoInches(PDFFloat(p.Paper.H))
	case uomPixels:
		return PDFFloat(p.Paper.H)
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

func (p *TPDFPage) CreateStdFontText(x, Y PDFFloat, AText string, AFont *TPDFEmbeddedFont, ADegrees float32, AUnderline, AStrikethrough bool) {
	p.AddObject(NewTPDFText(p.Document, x, Y, AText, AFont, ADegrees, AUnderline, AStrikethrough))
}

// func (p *TPDFPage) CreateTTFFontText(x, Y PDFFloat, AText string, AFont *TPDFEmbeddedFont, ADegrees float32, AUnderline, AStrikethrough bool) {
// 	p.AddTextToLookupLists(AText)
// 	T := p.Document.CreateText(x, Y, AText, AFont, ADegrees, AUnderline, AStrikethrough)
// 	p.AddObject(T)
// }

// procedure TPDFPage.WriteText(X, Y: TPDFFloat; AText: UTF8String; const ADegrees: single;
//     const AUnderline: boolean; const AStrikethrough: boolean);
// var
//   p: TPDFCoord;
// begin
//   if not Assigned(FLastFont) then
//     raise EPDF.Create(rsErrNoFontDefined);
//   p := Matrix.Transform(X, Y);
//   DoUnitConversion(p);
//   if Document.Fonts[FLastFont.FontIndex].IsStdFont then
//     CreateStdFontText(p.X, p.Y, AText, FLastFont, ADegrees, AUnderline, AStrikeThrough)
//   else
//     CreateTTFFontText(p.X, p.Y, AText, FLastFont, ADegrees, AUnderline, AStrikeThrough);
// end;

// procedure TPDFPage.WriteText(APos: TPDFCoord; AText: UTF8String; const ADegrees: single;
//
//	const AUnderline: boolean; const AStrikethrough: boolean);
//
// begin
//
//	WriteText(APos.X, APos.Y, AText, ADegrees, AUnderline, AStrikeThrough);
//
// end;
// FIXME: complete implementation
func (p *TPDFPage) WriteTextEx(x, y PDFFloat, txt string, ADegrees float32, underline, strikethro bool) {
	p1 := p.Matrix.TransformXY(x, y)
	p.DoUnitConversion(&p1)
	p.CreateStdFontText(p1.x, p1.y, txt, p.LastFont, ADegrees, underline, strikethro)
}

func (p *TPDFPage) WriteText(x, y PDFFloat, txt string) {
	p.WriteTextEx(x, y, txt, 0, false, false)
}
