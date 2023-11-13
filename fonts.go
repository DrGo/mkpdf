package main

import (
	"fmt"
	"io"
)

type EmbeddedFont struct {
	DocObj
	TxtFont        int
	TxtSize        float64
	Page           *Page
	SimulateBold   bool
	SimulateItalic bool
}

func NewEmbeddedFont(doc *Document, page *Page, fontNum int, size float64) *EmbeddedFont {
	return NewEmbeddedFontEx(doc, page, fontNum, size, false, false)
}

func NewEmbeddedFontEx(doc *Document, page *Page, fontNum int, size float64, ASimulateBold, ASimulateItalic bool) *EmbeddedFont {
	return &EmbeddedFont{
		TxtFont:        fontNum,
		TxtSize:        size,
		Page:           page,
		SimulateBold:   ASimulateBold,
		SimulateItalic: ASimulateItalic,
	}
}

func (f *EmbeddedFont) PointSize() int {
	return int(f.TxtSize)
	// val, _ := strconv.ParseFloat(f.FTxtSize, 64)
	// return int(math.Round(val))
}

func (f *EmbeddedFont) FontSize() float64 {
	return f.TxtSize
	// val, _ := strconv.ParseFloat(f.FTxtSize, 64)
	// return val
}

// WriteString('/F'+IntToStr(FTxtFont)+' '+FTxtSize+' Tf'+CRLF, AStream);
func (f *EmbeddedFont) Encode(st PDFWriter) {
	st.Writef("/F%d %f Tf%s", f.TxtFont, f.TxtSize, CRLF)
}

func (f *EmbeddedFont) WriteEmbeddedFont(doc *Document, src io.Reader, st PDFWriter) int {
	st.Writef("%sstream%s", CRLF, CRLF)
	PS := st.Offset()
	if f.Document.hasOption(poCompressFonts) {
		_ = compressStream(st, src)		
	} else {
		//FIXME: handle errors
		_, _ = io.Copy(st,  src)
	}
	st.WriteString(CRLF)
	st.WriteString("endstream")
	return st.Offset() - PS
}

// func (TPDFEmbeddedFont) WriteEmbeddedSubsetFont(doc *TPDFDocument, AFontNum int, AOutStream io.Writer) int64 {
// 	if doc.Fonts[AFontNum].SubsetFont == nil {
// 		panic("WriteEmbeddedSubsetFont: SubsetFont stream was not initialised.")
// 	}
// 	fmt.Fprintf(AOutStream, "%sstream%s", CRLF, CRLF)
// 	PS := int64(AOutStream.(io.Seeker).Seek(0, io.SeekCurrent))
// 	if doc.Options.Contains(poCompressFonts) {
// 		CompressedStream := new(bytes.Buffer)
// 		CompressStream(doc.Fonts[AFontNum].SubsetFont, CompressedStream)
// 		CompressedStream.WriteTo(AOutStream)
// 	} else {
// 		io.Copy(AOutStream, doc.Fonts[AFontNum].SubsetFont)
// 	}
// 	return int64(AOutStream.(io.Seeker).Seek(0, io.SeekCurrent)) - PS
// }

type Font struct {
	IsStdFont    bool
	Name         string
	FontFilename string
	// FTrueTypeFile    TTFFileInfo
	TextMappingList TTextMappingList
	SubsetFont      io.Reader
}

func NewFont(name   string, isStd bool) *Font { return &Font{Name: name, IsStdFont: isStd } }
func (f *Font) SetFontFilename(filename string) {
	if f.FontFilename == filename {
		return
	}
	f.FontFilename = filename
	f.PrepareTextMapping()
}

func (f *Font) PrepareTextMapping() {
	if f.FontFilename != "" {
		f.TextMappingList = TTextMappingList{}
		// f.FTrueTypeFile = &TTFFileInfo{}
		// f.FTrueTypeFile.LoadFromFile(f.FFontFilename)
		// f.FTrueTypeFile.PrepareFontDefinition("cp1252", true)
	}
}

// func (f *TPDFFont) GenerateSubsetFont() {
// 	if f.FSubsetFont != nil {
// 		f.FSubsetFont = nil
// 	}
// 	fontSubsetter := &TFontSubsetter{
// 		TrueTypeFile:    f.FTrueTypeFile,
// 		TextMappingList: f.FTextMappingList,
// 	}
// 	defer fontSubsetter.Free()
// 	f.FSubsetFont = &bytes.Buffer{}
// 	fontSubsetter.SaveToStream(f.FSubsetFont)
// 	fs, _ := os.Create(f.FTrueTypeFile.PostScriptName + "-subset.ttf")
// 	f.FSubsetFont.WriteTo(fs)
// 	fs.Close()
// }

func (f *Font) GetGlyphIndices(txt string) string {
	result := ""
	if len(txt) == 0 {
		return result
	}
	for _, char := range txt {
		c := uint16(char)
		matched := false
		for _, mapping := range f.TextMappingList {
			if mapping.CharID == c {
				result += fmt.Sprintf("%04X", mapping.GlyphID)
				c = 0
				matched = true
				break
			}
		}
		if !matched {
			result += fmt.Sprintf("%04X", c)
		}
	}
	return result
}

func (f *Font) AddTextToMappingList(txt string) {
	// for _, char := range txt {
	// 	c := uint16(char)
	// 	gid := f.FTrueTypeFile.GetGlyphIndex(c)
	// 	f.FTextMappingList=append(f.FTextMappingList, )
	// 	f.FTextMappingList.Add(c, gid)
	// }
}

type TrueTypeCharWidths struct {
	Document        *Document
	EmbeddedFontNum int
}

func (w *TrueTypeCharWidths) Encode(st PDFWriter) {
	// var s string
	// lst := w.Document.Fonts[w.EmbeddedFontNum].TextMapping
	// lst.Sort()
	// lFont := w.Document.Fonts[w.EmbeddedFontNum].FTrueTypeFile

	// for _, item := range lst.Items {
	// 	var lWidthIndex int
	// 	if item.GlyphID < lFont.HHead.numberOfHMetrics {
	// 		lWidthIndex = item.GlyphID
	// 	} else {
	// 		lWidthIndex = lFont.HHead.numberOfHMetrics - 1
	// 	}
	// 	s += fmt.Sprintf(" %d [%d]", item.GlyphID, TTTFFriendClass(lFont).ToNatural(lFont.Widths[lWidthIndex].AdvanceWidth))
	// }

	// st.Write([]byte(s))
}

func IsStandardFont(fontname string) bool {
	switch fontname {
	case "Courier", "Courier-Bold", "Courier-Oblique", "Courier-BoldOblique",
		"Helvetica", "Helvetica-Bold", "Helvetica-Oblique", "Helvetica-BoldOblique",
		"Times-Roman", "Times-Bold", "Times-Italic", "Times-BoldItalic",
		"Symbol", "ZapfDingbats":
		return true
	}
	return false
}
