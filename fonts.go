package main

import (
	"bytes"
	"fmt"
	"io"
)

type TPDFEmbeddedFont struct {
	TPDFDocumentObject
	TxtFont        int
	TxtSize        PDFFloat
	Page           *TPDFPage
	SimulateBold   bool
	SimulateItalic bool
}

func NewTPDFEmbeddedFont(ADocument *TPDFDocument, APage *TPDFPage, AFont int, ASize PDFFloat) *TPDFEmbeddedFont {
	return &TPDFEmbeddedFont{
		TxtFont: AFont,
		TxtSize: ASize,
		Page:    APage,
	}
}

func NewTPDFEmbeddedFontAdvanced(ADocument *TPDFDocument, APage *TPDFPage, AFont int, ASize PDFFloat, ASimulateBold, ASimulateItalic bool) *TPDFEmbeddedFont {
	return &TPDFEmbeddedFont{
		TxtFont:        AFont,
		TxtSize:        ASize,
		Page:           APage,
		SimulateBold:   ASimulateBold,
		SimulateItalic: ASimulateItalic,
	}
}

func (f *TPDFEmbeddedFont) PointSize() int {
	return int(f.TxtSize)
	// val, _ := strconv.ParseFloat(f.FTxtSize, 64)
	// return int(math.Round(val))
}

func (f *TPDFEmbeddedFont) FontSize() PDFFloat {
	return f.TxtSize
	// val, _ := strconv.ParseFloat(f.FTxtSize, 64)
	// return val
}

// WriteString('/F'+IntToStr(FTxtFont)+' '+FTxtSize+' Tf'+CRLF, AStream);
func (f *TPDFEmbeddedFont) Encode(st PDFWriter) {
	st.Writef("/F%d %f Tf%s", f.TxtFont, f.TxtSize, CRLF)
}

func (f *TPDFEmbeddedFont) WriteEmbeddedFont(ADocument *TPDFDocument, Src io.Reader, st PDFWriter) int {
	st.Writef("%sstream%s", CRLF, CRLF)
	PS := st.Offset()
	if f.Document.hasOption(poCompressFonts) {
		var CompressedStream bytes.Buffer
		//FIXME: which CompressedStream to use
		// CompressStream(Src, CompressedStream)
		st.Write(CompressedStream.Bytes())
	} else {
		//FIXME: handle errors
		buf, _ := io.ReadAll(Src)
		st.Write(buf)
	}

	st.WriteString(CRLF)
	st.WriteString("endstream")
	return st.Offset() - PS
}

// func (TPDFEmbeddedFont) WriteEmbeddedSubsetFont(ADocument *TPDFDocument, AFontNum int, AOutStream io.Writer) int64 {
// 	if ADocument.Fonts[AFontNum].SubsetFont == nil {
// 		panic("WriteEmbeddedSubsetFont: SubsetFont stream was not initialised.")
// 	}
// 	fmt.Fprintf(AOutStream, "%sstream%s", CRLF, CRLF)
// 	PS := int64(AOutStream.(io.Seeker).Seek(0, io.SeekCurrent))
// 	if ADocument.Options.Contains(poCompressFonts) {
// 		CompressedStream := new(bytes.Buffer)
// 		CompressStream(ADocument.Fonts[AFontNum].SubsetFont, CompressedStream)
// 		CompressedStream.WriteTo(AOutStream)
// 	} else {
// 		io.Copy(AOutStream, ADocument.Fonts[AFontNum].SubsetFont)
// 	}
// 	return int64(AOutStream.(io.Seeker).Seek(0, io.SeekCurrent)) - PS
// }

type TPDFFont struct {
	FIsStdFont       bool
	FName            string
	FFontFilename    string
	// FTrueTypeFile    TTFFileInfo
	FTextMappingList TTextMappingList
	FSubsetFont      io.Reader
}

func (f *TPDFFont) SetFontFilename(filename string) {
	if f.FFontFilename == filename {
		return
	}
	f.FFontFilename = filename
	f.PrepareTextMapping()
}

func (f *TPDFFont) PrepareTextMapping() {
	if f.FFontFilename != "" {
		f.FTextMappingList = TTextMappingList{}
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

func NewTPDFFont() *TPDFFont {
	return &TPDFFont{
	}
}
func (f *TPDFFont) GetGlyphIndices(AText string) string {
	result := ""
	if len(AText) == 0 {
		return result
	}
	for _, char := range AText {
		c := uint16(char)
		matched := false
		for _, mapping := range f.FTextMappingList{
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

func (f *TPDFFont) AddTextToMappingList(AText string) {
	// for _, char := range AText {
	// 	c := uint16(char)
	// 	gid := f.FTrueTypeFile.GetGlyphIndex(c)
	// 	f.FTextMappingList=append(f.FTextMappingList, )
	// 	f.FTextMappingList.Add(c, gid)
	// }
}

type TPDFTrueTypeCharWidths struct {
	Document        *TPDFDocument
	EmbeddedFontNum int
}

func (w *TPDFTrueTypeCharWidths) Encode(st PDFWriter) {
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

func IsStandardPDFFont(AFontName string) bool {
	switch AFontName {
	case "Courier", "Courier-Bold", "Courier-Oblique", "Courier-BoldOblique",
		"Helvetica", "Helvetica-Bold", "Helvetica-Oblique", "Helvetica-BoldOblique",
		"Times-Roman", "Times-Bold", "Times-Italic", "Times-BoldItalic",
		"Symbol", "ZapfDingbats":
		return true
	}
	return false
}
