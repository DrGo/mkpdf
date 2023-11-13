package main

import (
	"fmt"
	"strings"
)

type AbsString struct {
	DocObj
	FontIndex int
}

func NewAbstractString(document *Document) AbsString {
	return AbsString{NewDocObj(document), 0}
}

var pdfStringReplacer = strings.NewReplacer("\\", "\\\\", "(", "\\(", ")", "\\)")

// These symbols must be preceded by a backslash:  "(", ")", "\"
func escapePDFString(val string) string {
	return pdfStringReplacer.Replace(val)
}

type String struct {
	AbsString
	Value   string
	CPValue string
}

func NewString(doc *Document, val string) *String {
	if strings.ContainsAny(val, "()\\") {
		val = escapePDFString(val)
	}
	return &String{NewAbstractString(doc), val, ""}
}

func (s *String) GetCPValue() string {
	if s.CPValue == "" {
		s.CPValue = s.Value
		//FIXME:
		// SetCodePage(s.CPValue, 1252);
	}
	return s.CPValue
}

func (s *String) Encode(st PDFWriter) {
	st.WriteString("(" + s.GetCPValue()+ ")")
}

type RawHexString struct {
	DocObj
	Value string
}

func (hs *RawHexString) Encode(st PDFWriter) {
	st.Writef("<%s>", hs.Value)
}

func NewRawHexString(doc *Document, val string) *RawHexString {
	return &RawHexString{NewDocObj(doc), val}
}

//type TPDFUTF8String struct {
//	TPDFAbstractString
//	Value string
//}

type FreeFormString struct {
	AbsString
	Value string
}

func (f *FreeFormString) Encode(st PDFWriter) {
	st.WriteString(f.Value)
}

func NewFreeFormString(doc *Document, val string) *FreeFormString {
	return &FreeFormString{NewAbstractString(doc), val}
}

//type TPDFUTF16String struct {
//	TPDFAbstractString
//	Value string
//}
//func NewTPDFUTF16String(doc *TPDFDocument, val string, AFontIndex int) TPDFUTF16String {
//	utf16str := TPDFUTF16String{}
//	utf16str.FDocument = doc
//	utf16str.Value = val
//	utf16str.FontIndex = AFontIndex
//	return utf16str
//}

////FIXME:

//func (u *TPDFUTF16String) Encode(st TStream) {
//// var
////   i:integer;
////   us:utf8string;
////   s:ansistring;
////   wv:word;
//// begin
////   us := Utf8Encode(FValue);
////   if (length(us)<>length(fValue)) then // quote
////   begin
////     s:='\376\377'; // UTF-16BE BOM
////     for i:=1 to length(fValue) do
////     begin
////       wv:=word(fValue[i]);
////       s:=s+'\'+oct_str(hi(wv));
////       s:=s+'\'+oct_str(lo(wv));
////     end;
////   end else
////   begin
////     if (Pos('(', FValue) > 0) or (Pos(')', FValue) > 0) or (Pos('\', FValue) > 0) then
////       s := InsertEscape(FValue)
////     else
////       s:=fValue;
////   end;

//// 	var s string
//// 	if len(us) != len(u.Value) {
//// 		s = "\\376\\377"
//// 		for _, char := range u.Value {
//// 			wv := int(char)
//// 			s += "\\" + octStr(byte(wv>>8))
//// 			s += "\\" + octStr(byte(wv&0xFF))
//// 		}
//// 	} else {
//// 		if strings.ContainsAny(u.Value, "()\\") {
//// 			s = u.InsertEscape(u.Value)
//// 		} else {
//// 			s = u.Value
//// 		}
//// 	}
//// 	st.WriteString("(")
//// 	st.WriteString(s)
//// 	st.WriteString(")")
//}

//func (u *TPDFUTF8String) RemapedText() string {
//// function TPDFUTF8String.RemapedText: AnsiString;
//// var
////   s: UnicodeString;
//// begin
////   s := UTF8Decode(FValue);
////   Result := Document.Fonts[FontIndex].GetGlyphIndices(s);
//// end;
//	return ""
//}

//func (u *TPDFUTF8String) Encode(st TStream) {
//	st.WriteString("<")
//	st.WriteString(u.RemapedText())
//	st.WriteString(">")
//}

//	func NewTPDFUTF8String(doc *TPDFDocument, val string, AFontIndex int) TPDFUTF8String {
//		utf8str := TPDFUTF8String{}
//		utf8str.FDocument = doc
//		utf8str.Value = val
//		utf8str.FontIndex = AFontIndex
//		return utf8str
//	}
type BaseText struct {
	DocObj
	X, Y          float64
	Font          *EmbeddedFont
	Degrees       float64
	Underline     bool
	Color         ARGBColor
	StrikeThrough bool
}

func NewBaseText(document *Document) *BaseText {
	return &BaseText{
		DocObj:        NewDocObj(document),
		X:             0.0,
		Y:             0.0,
		Font:          nil,
		Degrees:       0.0,
		Underline:     false,
		Color:         clBlack,
		StrikeThrough: false,
	}
}

type Text struct {
	BaseText
	str *String
}

func NewText(document *Document, x, y float64, txt string, font *EmbeddedFont,
	degrees float64, underline, strikeThrough bool) *Text {
	t := &Text{
		BaseText{
			DocObj:        NewDocObj(document),
			X:             x,
			Y:             y,
			Font:          font,
			Degrees:       float64(degrees),
			Underline:     underline,
			StrikeThrough: strikeThrough,
		},
		NewString(document, txt)}
	if t.Font != nil && t.Font.Page != nil {
		t.Color = t.Font.Page.LastFontColor
	}
	return t
}

func (t *Text) GetTextWidth() float64 {
	lFontName := t.Document.Fonts[t.str.FontIndex].Name
	if !IsStandardFont(lFontName) {
		panic(fmt.Sprintf(rsErrUnknownStdFont, lFontName))
	}

	var lWidth int
	CPV := t.str.CPValue
	for i := 0; i < len(CPV); i++ {
		lWidth += GetStdFontCharWidthsArray(lFontName)[CPV[i]]
	}
	return float64(lWidth * t.Font.PointSize() / 1540)
}

func (t *Text) GetTextHeight() float64 {
	lFontName := t.Document.Fonts[t.str.FontIndex].Name
	var result int
	switch lFontName {
	case "Courier", "Courier-Bold", "Courier-Oblique", "Courier-BoldOblique":
		result = FONT_TIMES_COURIER_CAPHEIGHT
	case "Helvetica":
		result = FONT_HELVETICA_ARIAL_CAPHEIGHT
	case "Helvetica-Bold":
		result = FONT_HELVETICA_ARIAL_BOLD_CAPHEIGHT
	case "Helvetica-Oblique":
		result = FONT_HELVETICA_ARIAL_ITALIC_CAPHEIGHT
	case "Helvetica-BoldOblique":
		result = FONT_HELVETICA_ARIAL_BOLD_ITALIC_CAPHEIGHT
	case "Times-Roman":
		result = FONT_TIMES_CAPHEIGHT
	case "Times-Bold":
		result = FONT_TIMES_BOLD_CAPHEIGHT
	case "Times-Italic":
		result = FONT_TIMES_ITALIC_CAPHEIGHT
	case "Times-BoldItalic":
		result = FONT_TIMES_BOLD_ITALIC_CAPHEIGHT
	case "Symbol", "ZapfDingbats":
		result = 300
	default:
		panic(fmt.Sprintf(rsErrUnknownStdFont, lFontName))
	}
	return float64(result * t.Font.PointSize() / 1540)
}

func (t *Text) Encode(st PDFWriter) {
	var rad, radc, rads float64
	calcAngle := func(degrees float64) {
		rad = DegToRad(-degrees)
		rads, radc = sincos(rad)
	}
	st.WriteString("BT" + CRLF)
	if t.Degrees != 0.0 {
		calcAngle(t.Degrees)
		st.Writef("%.4f %.4f %.4f %.4f %.4f %.4f Tm%s", radc, -rads, rads, radc, t.X, t.Y, CRLF)
	} else {
		st.WriteString(fmt.Sprintf("%f %f TD%s", t.X, t.Y, CRLF))
	}

	t.str.Encode(st)
	st.WriteString(" Tj" + CRLF)
	st.WriteString("ET" + CRLF)

	if !t.Underline && !t.StrikeThrough {
		return
	}
	// lWidth := t.GetTextWidth()
	// lHeight := t.GetTextHeight()
	// //    convert the Font Units to Millimeters. This is also because fontcache DPI (default 96) could differ from PDF DPI (72).
	// lTextWidthInMM := (lWidth * cInchToMM) / gTTFontCache.DPI
	// lTextHeightInMM := (lHeight * cInchToMM) / gTTFontCache.DPI

	// //   if Degrees <> 0.0 then
	// //     // angled text
	// //     WriteString(Format('q %s %s %s %s %s %s cm', [t1, t2, t3, t1, FloatStr(X), FloatStr(Y)]) + CRLF, AStream)
	// //   else
	// //     // horizontal text
	// //     WriteString(Format('q 1 0 0 1 %s %s cm', [FloatStr(X), FloatStr(Y)]) + CRLF, AStream);

	// //   { set up a pen width and stroke color }
	// //   lColor := TPDFColor.Command(True, Color);
	// //   lLineWidth := FloatStr(mmToPDF(lTextHeightInMM / 12)) + ' w ';
	// //   WriteString(lLineWidth + lColor + CRLF, AStream);

	// //   { line segment is relative to matrix translation coordinate, set above }
	// //   if Underline then
	// //     WriteString(Format('0 -1.5 m %s -1.5 l S', [FloatStr(mmToPDF(lTextWidthInMM))]) + CRLF, AStream);
	// //   if StrikeThrough then
	// //     WriteString(Format('0 %s m %s %0:s l S', [FloatStr(mmToPDF(lTextHeightInMM) / 2), FloatStr(mmToPDF(lTextWidthInMM))]) + CRLF, AStream);

	// //   { restore graphics state to before the translation matrix adjustment }
	// //   WriteString('Q' + CRLF, AStream);

	// if t.Degrees != 0.0 {
	// 	st.WriteString(fmt.Sprintf("%f %f %f %f %f %f cm%s", a1, b1, c1, d1, t.X, t.Y, CRLF))
	// } else {
	// 	st.WriteString(fmt.Sprintf("1 0 0 1 %f %f cm%s", t.X, t.Y, CRLF))
	// }

	// fontData := lFC.FontData

	// if t.Underline {
	// 	lUnderlinePos := PDFTomm(-1.5)
	// 	lUnderlineSize := lTextHeightInMM / 12
	// 	if fontData.PostScript.UnderlinePosition != 0 {
	// 		lUnderlinePos = FontUnitsTomm(fontData.PostScript.UnderlinePosition, t.Font.PointSize, fontData.Head.UnitsPerEm)
	// 	}
	// 	if fontData.PostScript.underlineThickness != 0 {
	// 		lUnderlineSize = FontUnitsTomm(fontData.PostScript.underlineThickness, t.Font.PointSize, fontData.Head.UnitsPerEm)
	// 	}

	// 	lLineWidth := FloatStr(mmToPDF(lUnderlineSize)) + " w "
	// 	st.WriteString(lLineWidth + lColor + CRLF)
	// 	st.WriteString(fmt.Sprintf("0 %s m %s %s l S%s", FloatStr(mmToPDF(lUnderlinePos)), FloatStr(mmToPDF(lTextWidthInMM)), CRLF))
	// }
	// if t.StrikeThrough {
	// 	lStrikeOutPos := lTextHeightInMM / 2
	// 	lStrikeOutSize := lTextHeightInMM / 12
	// 	if fontData.OS2Data.yStrikeoutPosition != 0 {
	// 		lStrikeOutPos = FontUnitsTomm(fontData.OS2Data.yStrikeoutPosition, t.Font.PointSize, fontData.Head.UnitsPerEm)
	// 	}
	// 	if fontData.OS2Data.yStrikeoutSize != 0 {
	// 		lStrikeOutSize = FontUnitsTomm(fontData.OS2Data.yStrikeoutSize, t.Font.PointSize, fontData.Head.UnitsPerEm)
	// 	}

	// 	lLineWidth := FloatStr(mmToPDF(lStrikeOutSize)) + " w "
	// 	st.WriteString(lLineWidth + lColor + CRLF)
	// 	st.WriteString(fmt.Sprintf("0 %s m %s %s l S%s", FloatStr(mmToPDF(lStrikeOutPos)), FloatStr(mmToPDF(lTextWidthInMM)), CRLF))
	// }
}

// procedure TPDFText.Write(const AStream: TStream);
//   // result is in Font Units
//   lWidth := GetTextWidth;
//   lHeight := GetTextHeight;
//   { convert the Font Units to Millimeters. This is also because fontcache DPI (default 96) could differ from PDF DPI (72). }
//   lTextWidthInMM := (lWidth * cInchToMM) / gTTFontCache.DPI;
//   lTextHeightInMM := (lHeight * cInchToMM) / gTTFontCache.DPI;

//   if Degrees <> 0.0 then
//     // angled text
//     WriteString(Format('q %s %s %s %s %s %s cm', [t1, t2, t3, t1, FloatStr(X), FloatStr(Y)]) + CRLF, AStream)
//   else
//     // horizontal text
//     WriteString(Format('q 1 0 0 1 %s %s cm', [FloatStr(X), FloatStr(Y)]) + CRLF, AStream);

//   { set up a pen width and stroke color }
//   lColor := TPDFColor.Command(True, Color);
//   lLineWidth := FloatStr(mmToPDF(lTextHeightInMM / 12)) + ' w ';
//   WriteString(lLineWidth + lColor + CRLF, AStream);

//   { line segment is relative to matrix translation coordinate, set above }
//   if Underline then
//     WriteString(Format('0 -1.5 m %s -1.5 l S', [FloatStr(mmToPDF(lTextWidthInMM))]) + CRLF, AStream);
//   if StrikeThrough then
//     WriteString(Format('0 %s m %s %0:s l S', [FloatStr(mmToPDF(lTextHeightInMM) / 2), FloatStr(mmToPDF(lTextWidthInMM))]) + CRLF, AStream);

//   { restore graphics state to before the translation matrix adjustment }
//   WriteString('Q' + CRLF, AStream);
// end;

// type TPDFUTF8Text struct {
// 	PDFBaseText
// 	string *PDFUTF8String
// }

// type TPDFUTF16Text struct {
// 	PDFBaseText
// 	string *PDFUTF16String
// }

// type TPDFUTF8Text struct {
// 	X, Y, Degrees float64
// 	Font          *TPDFEmbeddedFont
// 	Underline     bool
// 	StrikeThrough bool
// }

// func (t *TPDFUTF8Text) Write(st *TStream) {
// 	st.WriteString("q"+CRLF)

// 	st.WriteString("BT"+CRLF)

// 	a1, b1, c1, d1 := 1.0, 0.0, 0.0, 1.0
// 	if t.Degrees != 0.0 {
// 		rad := DegToRad(-t.Degrees)
// 		a1 = math.Cos(rad)
// 		b1 = -math.Sin(rad)
// 		c1 = math.Sin(rad)
// 		d1 = a1
// 	} else {
// 		st.WriteString(fmt.Sprintf("%f %f TD%s", t.X, t.Y, CRLF))
// 	}

// 	lFC := gTTFontCache.Find(t.FDocument.Fonts[t.Font.FontIndex].Name)

// 	lColor := TPDFColor.Command(true, t.Color)

// 	if lFC != nil {
// 		if t.Font.SimulateBold && !lFC.IsBold {
// 			st.WriteString(lColor+CRLF)
// 			st.WriteString(fmt.Sprintf("2 Tr %s w", FloatStr(t.Font.PointSize/30))+CRLF)
// 		}
// 		if t.Font.SimulateItalic && !lFC.IsItalic {
// 			a2, b2 := 1.0, 0.0
// 			c2 := math.Tan(DegToRad(12))
// 			d2 := 1.0
// 			a1, b1, c1, d1 = a2*a1+b2*c1, a2*b1+b2*d1, c2*a1+d2*c1, c2*b1+d2*d1
// 		}
// 	}
// 	if t.Degrees != 0.0 || (t.Font.SimulateItalic && !lFC.IsItalic) {
// 		st.WriteString(fmt.Sprintf("%f %f %f %f %f %f Tm%s", a1, b1, c1, d1, t.X, t.Y, CRLF))
// 	}

// 	t.FString.Write(st)
// 	st.WriteString(" Tj"+CRLF)
// 	st.WriteString("ET"+CRLF)

// 	if !t.Underline && !t.StrikeThrough {
// 		return
// 	}

// 	if lFC == nil {
// 		return
// 	}

// 	lWidth := lFC.TextWidth(t.FString.Value, t.Font.PointSize)
// 	lHeight, lDescender := lFC.TextHeight(t.FString.Value, t.Font.PointSize)
// 	lTextWidthInMM := (lWidth * cInchToMM) / gTTFontCache.DPI
// 	lTextHeightInMM := (lHeight * cInchToMM) / gTTFontCache.DPI

// 	if t.Degrees != 0.0 {
// 		st.WriteString(fmt.Sprintf("%f %f %f %f %f %f cm%s", a1, b1, c1, d1, t.X, t.Y, CRLF))
// 	} else {
// 		st.WriteString(fmt.Sprintf("1 0 0 1 %f %f cm%s", t.X, t.Y, CRLF))
// 	}

// 	fontData := lFC.FontData

// 	if t.Underline {
// 		lUnderlinePos := PDFTomm(-1.5)
// 		lUnderlineSize := lTextHeightInMM / 12
// 		if fontData.PostScript.UnderlinePosition != 0 {
// 			lUnderlinePos = FontUnitsTomm(fontData.PostScript.UnderlinePosition, t.Font.PointSize, fontData.Head.UnitsPerEm)
// 		}
// 		if fontData.PostScript.underlineThickness != 0 {
// 			lUnderlineSize = FontUnitsTomm(fontData.PostScript.underlineThickness, t.Font.PointSize, fontData.Head.UnitsPerEm)
// 		}

// 		lLineWidth := FloatStr(mmToPDF(lUnderlineSize)) + " w "
// 		st.WriteString(lLineWidth+lColor+CRLF)
// 		st.WriteString(fmt.Sprintf("0 %s m %s %s l S%s", FloatStr(mmToPDF(lUnderlinePos)), FloatStr(mmToPDF(lTextWidthInMM)), CRLF))
// 	}
// 	if t.StrikeThrough {
// 		lStrikeOutPos := lTextHeightInMM / 2
// 		lStrikeOutSize := lTextHeightInMM / 12
// 		if fontData.OS2Data.yStrikeoutPosition != 0 {
// 			lStrikeOutPos = FontUnitsTomm(fontData.OS2Data.yStrikeoutPosition, t.Font.PointSize, fontData.Head.UnitsPerEm)
// 		}
// 		if fontData.OS2Data.yStrikeoutSize != 0 {
// 			lStrikeOutSize = FontUnitsTomm(fontData.OS2Data.yStrikeoutSize, t.Font.PointSize, fontData.Head.UnitsPerEm)
// 		}

// 		lLineWidth := FloatStr(mmToPDF(lStrikeOutSize)) + " w "
// 		st.WriteString(lLineWidth+lColor+CRLF)
// 		st.WriteString(fmt.Sprintf("0 %s m %s %s l S%s", FloatStr(mmToPDF(lStrikeOutPos)), FloatStr(mmToPDF(lTextWidthInMM)), CRLF))
// 	}
// }

// func NewTPDFUTF8Text(doc *TPDFDocument, AX, AY float64, AText string, AFont *TPDFEmbeddedFont, ADegrees Tfloat64, AUnderline, AStrikeThrough bool) *TPDFUTF8Text {
// 	t := &TPDFUTF8Text{
// 		X:             AX,
// 		Y:             AY,
// 		Font:          AFont,
// 		Degrees:       ADegrees,
// 		Underline:     AUnderline,
// 		StrikeThrough: AStrikeThrough,
// 		FString:       doc.CreateUTF8String(AText, AFont.FontIndex),
// 	}
// 	if AFont != nil && AFont.Page != nil {
// 		t.Color = AFont.Page.FLastFontColor
// 	}
// 	return t
// }

// type TPDFUTF8Text struct {
// 	TPDFText
// }

// type TPDFUTF16Text struct {
// 	X, Y, Degrees            float64
// 	Font                     *TPDFEmbeddedFont
// 	Underline, StrikeThrough bool
// 	Color                    string
// }

// func (txt *TPDFUTF16Text) Write(st *TStream) {

// 	txst.WriteString("q"+CRLF)
// 	defer txst.WriteString("Q"+CRLF)

// 	txt.st.WriteString("BT"+CRLF)

// 	a1, b1, c1, d1 := 1.0, 0.0, 0.0, 1.0

// 	if txt.Degrees != 0.0 {
// 		rad := DegToRad(-txt.Degrees)
// 		a1, b1, c1, d1 = math.Cos(rad), -math.Sin(rad), math.Sin(rad), math.Cos(rad)
// 	} else {
// 		txt.st.WriteString(fmt.Sprintf("%f %f TD%s", txt.X, txt.Y, CRLF))
// 	}

// 	lFC := gTTFontCache.Find(txt.FDocument.Fonts[txt.Font.FontIndex].Name)
// 	lColor := TPDFColor.Command(true, txt.Color)

// 	if lFC != nil {
// 		if txt.Font.SimulateBold && !lFC.IsBold {
// 			txt.st.WriteString(lColor+CRLF)
// 			txt.st.WriteString(fmt.Sprintf("2 Tr %f w%s", txt.Font.PointSize/30, CRLF))
// 		}
// 		if txt.Font.SimulateItalic && !lFC.IsItalic {
// 			a2, b2 := 1.0, 0.0
// 			c2, d2 := math.Tan(DegToRad(12)), 1.0
// 			a1, b1, c1, d1 = a2*a1+b2*c1, a2*b1+b2*d1, c2*a1+d2*c1, c2*b1+d2*d1
// 		}
// 	}

// 	if txt.Degrees != 0.0 || (txt.Font.SimulateItalic && !lFC.IsItalic) {
// 		txt.st.WriteString(fmt.Sprintf("%f %f %f %f %f %f Tm%s", a1, b1, c1, d1, txt.X, txt.Y, CRLF))
// 	}

// 	txt.FString.Write(st)
// 	txt.st.WriteString(" Tj"+CRLF)
// 	txt.st.WriteString("ET"+CRLF)

// 	if !txt.Underline && !txt.StrikeThrough {
// 		return
// 	}

// 	if lFC == nil {
// 		return
// 	}

// 	v := UTF8Encode(txt.FString.Value)
// 	lWidth := lFC.TextWidth(v, txt.Font.PointSize)
// 	lHeight, lDescender := lFC.TextHeight(v, txt.Font.PointSize)
// 	lTextWidthInMM := (lWidth * cInchToMM) / gTTFontCache.DPI
// 	lTextHeightInMM := (lHeight * cInchToMM) / gTTFontCache.DPI

// 	if txt.Degrees != 0.0 {
// 		txt.st.WriteString(fmt.Sprintf("%f %f %f %f %f %f cm%s", a1, b1, c1, d1, txt.X, txt.Y, CRLF))
// 	} else {
// 		txt.st.WriteString(fmt.Sprintf("1 0 0 1 %f %f cm%s", txt.X, txt.Y, CRLF))
// 	}

// 	lFD := lFC.FontData
// 	if txt.Underline {
// 		lUnderlinePos := PDFTomm(-1.5)
// 		lUnderlineSize := lTextHeightInMM / 12
// 		if lFD.PostScript.UnderlinePosition != 0 {
// 			lUnderlinePos = FontUnitsTomm(lFD.PostScript.UnderlinePosition, txt.Font.PointSize, lFD.Head.UnitsPerEm)
// 		}
// 		if lFD.PostScript.underlineThickness != 0 {
// 			lUnderlineSize = FontUnitsTomm(lFD.PostScript.underlineThickness, txt.Font.PointSize, lFD.Head.UnitsPerEm)
// 		}
// 		lLineWidth := fmt.Sprintf("%f w ", mmToPDF(lUnderlineSize))
// 		txt.st.WriteString(lLineWidth+lColor+CRLF)
// 		txt.st.WriteString(fmt.Sprintf("0 %f m %f 0 l S%s", mmToPDF(lUnderlinePos), mmToPDF(lTextWidthInMM), CRLF))
// 	}

// 	if txt.StrikeThrough {
// 		lStrikeOutPos := lTextHeightInMM / 2
// 		lStrikeOutSize := lTextHeightInMM / 12
// 		if lFD.OS2Data.yStrikeoutPosition != 0 {
// 			lStrikeOutPos = FontUnitsTomm(lFD.OS2Data.yStrikeoutPosition, txt.Font.PointSize, lFD.Head.UnitsPerEm)
// 		}
// 		if lFD.OS2Data.yStrikeoutSize != 0 {
// 			lStrikeOutSize = FontUnitsTomm(lFD.OS2Data.yStrikeoutSize, txt.Font.PointSize, lFD.Head.UnitsPerEm)
// 		}
// 		lLineWidth := fmt.Sprintf("%f w ", mmToPDF(lStrikeOutSize))
// 		txt.st.WriteString(lLineWidth+lColor+CRLF)
// 		txt.st.WriteString(fmt.Sprintf("0 %f m %f 0 l S%s", mmToPDF(lStrikeOutPos), mmToPDF(lTextWidthInMM), CRLF))
// 	}
// }

//	func NewTPDFUTF16Text(doc *TPDFDocument, AX, AY float64, AText string, AFont *TPDFEmbeddedFont, ADegrees Tfloat64, AUnderline, AStrikeThrough bool) *TPDFUTF16Text {
//		txt := &TPDFUTF16Text{
//			X:             AX,
//			Y:             AY,
//			Font:          AFont,
//			Degrees:       ADegrees,
//			Underline:     AUnderline,
//			StrikeThrough: AStrikeThrough,
//			FString:       doc.CreateUTF16String(AText, AFont.FontIndex),
//		}
//		if AFont != nil && AFont.Page != nil {
//			txt.Color = AFont.Page.FLastFontColor
//		}
//		return txt
//	}

type TTextMapping struct {
	CharID, GlyphID uint16
}
type TTextMappingList []TTextMapping

func NewTPDFTrueTypeCharWidths(doc *Document, AEmbeddedFontNum int) *TrueTypeCharWidths {
	return &TrueTypeCharWidths{
		Document:        doc,
		EmbeddedFontNum: AEmbeddedFontNum,
	}
}

// type PDFMiterLimit struct {
// 	PDFGraphicObject
// 	MiterLimit float64
// }

// func (ml *PDFMiterLimit) Write(st TStream) {
// 	st.WriteString(fmt.Sprintf("%f M%s", ml.MiterLimit, CRLF), stream)
// }

// func NewPDFMiterLimit(document *TPDFDocument, miterLimit float64) *PDFMiterLimit {
// 	return &PDFMiterLimit{
// 		PDFGraphicObject: NewPDFGraphicObject(document),
// 		MiterLimit:       miterLimit,
// 	}
// }

// func NewTPDFFontNumBaseObject(doc *TPDFDocument, aFontNum int) *TPDFFontNumBaseObject {
// 	return &TPDFFontNumBaseObject{
// 		TPDFObject: *NewTPDFObject(doc),
// 		FFontNum:   aFontNum,
// 	}
// }

// func (p *TPDFToUnicode) Write(st TStream){
// 	lst := p.FDocument.Fonts[p.FontNum].TextMapping

// 	st.WriteString("/CIDInit /ProcSet findresource begin"+CRLF)
// 	st.WriteString("12 dict begin"+CRLF)
// 	st.WriteString("begincmap"+CRLF)
// 	st.WriteString("/CIDSystemInfo"+CRLF)
// 	st.WriteString("<</Registry (Adobe)"+CRLF)

// 	if p.FDocument.Options&poSubsetFont != 0 {
// 		st.WriteString("/Ordering (UCS)"+CRLF)
// 	} else {
// 		st.WriteString("/Ordering (Identity)"+CRLF)
// 	}

// 	st.WriteString("/Supplement 0"+CRLF)
// 	st.WriteString(">> def"+CRLF)

// 	if p.FDocument.Options&poSubsetFont != 0 {
// 		st.WriteString(fmt.Sprintf("/CMapName /Adobe-Identity-UCS def"+CRLF))
// 	} else {
// 		st.WriteString(fmt.Sprintf("/CMapName /%s def", p.FDocument.Fonts[p.FontNum].FTrueTypeFile.PostScriptName)+CRLF)
// 	}

// 	st.WriteString("1 begincodespacerange"+CRLF)
// 	st.WriteString("<0000> <FFFF>"+CRLF)
// 	st.WriteString("endcodespacerange"+CRLF)

// 	if p.FDocument.Options&poSubsetFont != 0 {
// 		st.WriteString(fmt.Sprintf("%d beginbfrange", len(lst)-1)+CRLF)
// 		for i := 0; i < len(lst)-1; i++ {
// 			st.WriteString(fmt.Sprintf("<%s> <%s> <%s>", IntToHex(lst[i].GlyphID, 4), IntToHex(lst[i].GlyphID, 4), IntToHex(lst[i].CharID, 4))+CRLF)
// 		}
// 		st.WriteString("endbfrange"+CRLF)
// 	} else {
// 		st.WriteString(fmt.Sprintf("%d beginbfchar", len(lst))+CRLF)
// 		for i := 0; i < len(lst); i++ {
// 			st.WriteString(fmt.Sprintf("<%s> <%s>", IntToHex(lst[i].GlyphID, 4), IntToHex(lst[i].CharID, 4))+CRLF)
// 		}
// 		st.WriteString("endbfchar"+CRLF)
// 	}
// 	st.WriteString("endcmap"+CRLF)
// 	st.WriteString("CMapName currentdict /CMap defineresource pop"+CRLF)
// 	st.WriteString("end"+CRLF)
// 	st.WriteString("end"+CRLF)
// }

// func (c *TCIDToGIDMap) Write(st TStream){
// 	lst := c.FDocument.Fonts[c.FontNum].TextMapping
// 	sort.Slice(lst, func(i, j int) bool { return lst[i].GlyphID < lst[j].GlyphID })
// 	lMaxCID := lst[len(lst)-1].GlyphID
// 	ba := make([]byte, (lMaxCID+1)*2)
// 	for i := 0; i < len(lst); i++ {
// 		cid := lst[i].GlyphID
// 		gid := lst[i].NewGlyphID

// 		ba[2*cid] = byte(gid >> 8)
// 		ba[2*cid+1] = byte(gid)
// 	}

// 	st.WriteBuffer(ba, len(ba))
// }

// func (c *TPDFCIDSet) Encode(st TStream){
// 	lst := c.FDocument.Fonts[c.FontNum].TextMapping
// 	sort.Slice(lst, func(i, j int) bool { return lst[i].GlyphID < lst[j].GlyphID })
// 	lSize := (lst[len(lst)-1].GlyphID / 8) + 1
// 	ba := make([]byte, lSize)
// 	for i := 0; i < len(lst); i++ {
// 		cid := lst[i].GlyphID
// 		mask := uint8(1 << (7 - (cid % 8)))
// 		gid := cid / 8
// 		ba[gid] = ba[gid] | mask
// 	}
// 	st.WriteBuffer(ba, len(ba))
// }

// type TPDFFontNumBaseObject struct {
// 	TPDFDocumentObject
// 	FFontNum int
// }

// func NewTPDFFontNumBaseObject(doc *TPDFDocument, AFontNum int) *TPDFFontNumBaseObject {
// 	return &TPDFFontNumBaseObject{TPDFDocumentObject: NewTPDFDocumentObject(doc), FFontNum: AFontNum}
// }

// type TPDFToUnicode struct {
// 	TPDFFontNumBaseObject
// }

// type TCIDToGIDMap struct {
// 	TPDFFontNumBaseObject
// }

// type TPDFCIDSet struct {
// 	TPDFFontNumBaseObject
// }
