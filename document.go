package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"time"
)

type Document struct {
	Catalogue            int
	CurrentColor         string
	CurrentWidth         string
	LineCapStyle         LineCapStyle
	DefaultOrientation   PaperOrientation
	DefaultPaperType     PaperType
	FontDir              string
	FontFiles            []string
	Fonts                []*Font
	Images               []*RasterImage
	Infos                *PDFInfo
	LineStyleDefs        LineStyleDef
	Options              []Option
	Pages                []*Page
	Preferences          bool
	PageLayout           PageLayout
	Sections             SectionList
	Trailer              *Dictionary
	ZoomValue            int
	GlobalXRefs          []*XRef
	UnitOfMeasure        UnitOfMeasure
	DefaultUnitOfMeasure UnitOfMeasure
	ImageStreamOptions   []ImageStreamOption
}

func NewDocument() *Document {
	doc := &Document{
		// doc.FontFiles = make([]string, 0)
		Preferences:          true,
		PageLayout:           lSingle,
		DefaultPaperType:     ptA4,
		DefaultOrientation:   ppoPortrait,
		ZoomValue:            100,
		Options:              []Option{poCompressFonts, poCompressImages},
		DefaultUnitOfMeasure: uomMillimeters,
		LineCapStyle:         plcsRoundCap,
		Sections:             *NewTPDFSectionList(),
		Infos:                NewPDFInfo(),
	}
	return doc
}

func (doc *Document) SetOptions(AValue ...Option) {
	if slices.Contains(AValue, poNoEmbeddedFonts) {
		AValue = append(AValue, poSubsetFont)
	}
	doc.Options = AValue
}

func (doc *Document) hasOption(opt Option) bool {
	return slices.Contains(doc.Options, opt)
}

func (doc *Document) StartDocument() {
	doc.CreateRefTable()
	doc.CreateTrailer()
	doc.Catalogue = doc.CreateCatalogEntry()
	doc.CreateInfoEntry(doc.hasOption(poUTF16info))

	if doc.hasOption(poMetadataEntry) {
		doc.CreateMetadataEntry()
	}
	if !doc.hasOption(poNoTrailerID) {
		doc.CreateTrailerID()
	}
	doc.CreatePreferencesEntry()
	if doc.FontDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			doc.FontDir = cwd
		}
	}
}

func (doc *Document) PageCount() int  { return len(doc.Pages) }
func (doc *Document) FontCount() int  { return len(doc.Fonts) }
func (doc *Document) ImageCount() int { return len(doc.Images) }
func (doc *Document) xrefCount() int  { return len(doc.GlobalXRefs) }

func (doc *Document) NewPage() *Page {
	doc.Pages = append(doc.Pages, NewPage(doc))
	return doc.Pages[doc.PageCount()-1]
}

func (doc *Document) AddGlobalXRef(xr *XRef) int {
	doc.GlobalXRefs = append(doc.GlobalXRefs, xr)
	return len(doc.GlobalXRefs) - 1
}

// FIXME: check this is correct
// the first object has id of 0 (not used )
func (doc *Document) CreateRefTable() {
	eadd(&doc.GlobalXRefs, NewXRef(doc))
}

func (doc *Document) CreateGlobalXRef(typ string) *XRef {
	xr := NewXRef(doc)
	xr.Dict.AddName("Type", typ)
	doc.AddGlobalXRef(xr)
	return xr
}

func (doc *Document) GetXref(idx int) *XRef {
	return doc.GlobalXRefs[idx]
}

func (doc *Document) SaveToWriter(st PDFWriter) error {
	doc.CreateSectionsOutLine()
	doc.CreateFontEntries()
	doc.CreateImageEntries()
	doc.Trailer.SetSize(doc.xrefCount())
	st.WriteString(PDF_VERSION + "\n")
	st.WriteString(PDF_BINARY_BLOB + "\n")
	xRefPos := 0
	for i := 1; i < doc.xrefCount(); i++ {
		xRefPos = st.Offset()
		doc.WriteObject(i, st)
		doc.GlobalXRefs[i].Offset = xRefPos
	}
	st.Writef("xref\n0 %d\n", doc.xrefCount())
	doc.WriteXRefTable(st)
	st.Writef("trailer\n")
	doc.Trailer.Encode(st)
	st.Writef("\nstartxref\n%d\n", xRefPos)
	st.WriteString(PDF_FILE_END)
	return st.Err()
}

func (doc *Document) SaveToFile(AFileName string) (err error) {
	f, err := os.Create(AFileName)
	if err != nil {
		return err
	}
	defer func() {
		if errc := f.Close(); errc != nil {
			err = errc
		}
	}()
	w := bufio.NewWriter(f)
	defer func() {
		if errc := w.Flush(); errc != nil {
			err = errc
		}
	}()
	return doc.SaveToWriter(Newfwriter(w))
}

func (doc *Document) WriteObject(obj int, st PDFWriter) {
	st.Writef("%d 0 obj\n", obj)
	x := doc.GlobalXRefs[obj]
	if x.Stream == nil {
		x.Dict.WriteDictionary(obj, st)
	} else {
		doc.CurrentColor = ""
		doc.CurrentWidth = ""
		mem := NewMemWriter(doc, nil)
		x.Stream.Encode(mem)
		d := mem.Size()
		// 		if poCompressText&Options != 0 && !X.FStream.CompressionProhibited {
		// 			MCompressed = &TMemoryStream{}
		// 			CompressStream(M, MCompressed)
		// 			X.Dict.AddName("Filter", "FlateDecode")
		// 			d = MCompressedoc.Size()
		// 		}
		x.Dict.AddInteger("Length", d)
		x.Dict.Encode(st)

		doc.CurrentColor = ""
		doc.CurrentWidth = ""
		st.WriteString("\nstream\n")
		// if poCompressText&Options != 0 && !X.FStream.CompressionProhibited {
		// 	MCompressedoc.Position = 0
		// 	MCompressedoc.SaveToStream(st)
		// 	MCompressed = nil
		// } else {
		mem.Encode(st)
		// }
		st.WriteString("\nendstream")
	}
	st.WriteString(CRLF + "endobj" + CRLF + CRLF)
}

//FIXME:
// func (d *TPDFDocument) AddXObject(AXObject *TPDFXObject) int {
//     doc.FXObjects = append(doc.FXObjects, AXObject)
//     return len(doc.FXObjects) - 1
// }
//FIXME:
// func (d *TPDFDocument) AddPattern(APattern *TPDFPattern) int {
//     doc.FPatterns = append(doc.FPatterns, APattern)
//     return len(doc.FPatterns) - 1
// }

// func (d *TPDFDocument) AddGraphicState(AGraphicState *TPDFGraphicState) int {
//     doc.FGraphicStates = append(doc.FGraphicStates, AGraphicState)
//     return len(doc.FGraphicStates) - 1
// }

// func (d *TPDFDocument) AddAnnotation(AAnnotation *TPDFAnnotation) int {
//     doc.FAnnotations = append(doc.FAnnotations, AAnnotation)
//     return len(doc.FAnnotations) - 1
// }

func (doc *Document) CreateTrailer() {
	doc.Trailer = NewDictionary(doc)
	doc.Trailer.SetSize(doc.xrefCount())
}

func (doc *Document) CreateCatalogEntry() int {
	dict := doc.CreateGlobalXRef("Catalog").Dict
	doc.Trailer.AddReference("Root", doc.xrefCount()-1)
	dict.AddName("PageLayout", PageLayoutNames[doc.PageLayout])
	dict.AddElement("OpenAction", NewArray(doc))
	return doc.xrefCount() - 1
}

func (doc *Document) CreateInfoEntry(UseUTF16 bool) {
	dict := doc.CreateGlobalXRef("Info").Dict
	doc.Trailer.AddReference("Info", doc.xrefCount()-1)
	doc.Trailer.SetSize(doc.xrefCount())
	noUnicode := false
	doEntry := func(name, val string) {
		if val == "" {
			return
		}
		if UseUTF16 && !noUnicode {
			//FIXME:
			// IDict.AddString(aName, utf8Decode(aValue))
		} else {
			dict.AddString(name, val)
		}
	}

	doEntry("Title", doc.Infos.Title)
	doEntry("Author", doc.Infos.Author)
	doEntry("Creator", doc.Infos.ApplicationName)
	doEntry("Producer", doc.Infos.Producer)
	doEntry("Keywords", doc.Infos.Keywords)
	noUnicode = true
	doEntry("CreationDate", DateToPdfDate(doc.Infos.CreationDate))
}

func (doc *Document) CreateMetadataEntry() {
	xr := doc.CreateGlobalXRef("Metadata")
	xr.Dict.AddName("Subtype", "XML")
	xr.Stream = NewPDFStream(doc)
	xr.Stream.AddItem(NewTXMPStream(doc))
	xr.Stream.CompressionProhibited = true

	doc.GlobalXRefs[doc.Catalogue].Dict.AddReference("Metadata", doc.xrefCount()-1)
}

func (doc *Document) AddOutputIntent(subtype, OutCondID, info string, ICCProfile io.Reader) {
	rxNum := doc.xrefCount()
	dict := doc.CreateGlobalXRef("OutputIntent").Dict
	dict.AddName("S", subtype)
	dict.AddString("OutputConditionIdentifier", OutCondID)
	if info != "" {
		dict.AddString("Info", info)
	}
	if ICCProfile != nil {
		profile := doc.CreateGlobalXRef("ICCProfile")
		profile.Dict.AddInteger("N", 3)
		profile.Stream = NewPDFStream(doc)
		profile.Stream.AddItem(NewMemWriter(doc, ICCProfile))
		dict.AddReference("DestOutputProfile", doc.xrefCount()-1)
	}

	OutputIntents := doc.GlobalXRefs[doc.Catalogue].Dict.FindElement("OutputIntents")
	if OutputIntents == nil {
		OutputIntents = doc.GlobalXRefs[doc.Catalogue].Dict.AddElement("OutputIntents", NewArray(doc))
	}
	OutputIntents.Value.(*Array).AddItem(NewReference(rxNum))
}

func (doc *Document) AddPDFA1sRGBOutputIntent() {
	var buf bytes.Buffer
	buf.Grow(len(ICC_sRGB2014) - 1)
	buf.Write(ICC_sRGB2014[1:])
	doc.AddOutputIntent("GTS_PDFA1", "Custom", "sRGB", &buf)
}

func (doc *Document) CreateTrailerID() {
	s := DateToPdfDate(time.Now()) + strconv.Itoa(doc.xrefCount()) +
		doc.Infos.Title + doc.Infos.Author + doc.Infos.ApplicationName + doc.Infos.Producer + DateToPdfDate(doc.Infos.CreationDate)
	s = GetMD5Hash(s)
	id := NewArray(doc)
	id.AddItem(NewRawHexString(doc, s))
	id.AddItem(NewRawHexString(doc, s))
	doc.Trailer.AddElement("ID", id)
}

func (doc *Document) CreatePreferencesEntry() {
	dict := doc.CreateGlobalXRef("ViewerPreferences").Dict
	dict.AddElement("FitWindow", NewBoolean(true))
	doc.GlobalXRefs[doc.Catalogue].Dict.AddReference("ViewerPreferences", doc.xrefCount()-1)
}

func (doc *Document) CreatePagesEntry(Parent int) int {
	pgdict := doc.CreateGlobalXRef("Pages").Dict
	xr := doc.xrefCount() - 1
	pgdict.AddElement("Kids", NewArray(doc))
	pgdict.AddInteger("Count", 0)
	if Parent == 0 {
		doc.GlobalXRefs[doc.Catalogue].Dict.AddReference("Pages", xr)
	} else {
		pgdict.AddReference("Parent", Parent)
		dict := doc.GlobalXRefs[Parent].Dict
		dict.IncCount()
		dict.AddKid(NewReference(xr))
	}
	return xr
}

func (doc *Document) CreatePageEntry(parentNum, pageNum int) int {
	pp := doc.Pages[pageNum]
	pgdict := doc.CreateGlobalXRef("Page").Dict
	pgdict.AddReference("Parent", parentNum)

	parentDict := doc.GlobalXRefs[parentNum].Dict
	parentDict.IncCount()
	parentDict.AddKid(NewReference(doc.xrefCount() - 1))

	arr := NewArray(doc).AddInts(0, 0, pp.Paper.W, pp.Paper.H)
	pgdict.AddElement("MediaBox", arr)
	// doc.CreateAnnotEntries(pageNum, pDict)
	resDict := NewDictionary(doc)
	arr = NewArray(doc) // procset
	arr.AddItem(NewPDFName("PDF"))
	arr.AddItem(NewPDFName("Text"))
	arr.AddItem(NewPDFName("ImageC"))
	resDict.AddElement("ProcSet", arr)
	// fmt.Printf("font count before creating Resource/Font %d\n", doc.FontCount())
	if doc.FontCount() > 0 {
		resDict.AddElement("Font", NewDictionary(doc))
	}
	pgdict.AddElement("Resources", resDict)
	// if pp.HasImages {
	// 	aDict.AddElement("XObject", doc.CreateDictionary())
	// }

	return doc.xrefCount() - 1
}

func (doc *Document) CreateOutlines() int {
	dict := doc.CreateGlobalXRef("Outlines").Dict
	dict.AddInteger("Count", 0)
	return doc.xrefCount() - 1
}

func (doc *Document) CreateOutlineEntry(parent, sectNo, pageNo int, aTitle string) int {
	dict := doc.CreateGlobalXRef("OutlineEntry").Dict
	s := aTitle
	if s == "" {
		s = fmt.Sprintf("Section %d", sectNo)
	}
	if pageNo > -1 {
		s = fmt.Sprintf("%s Page %d", s, pageNo)
	}
	dict.AddString("Title", s)
	dict.AddReference("Parent", parent)
	dict.AddInteger("Count", 0)
	dict.AddElement("Dest", NewArray(doc))
	return doc.xrefCount() - 1
}

// for i:=1 to GLobalXRefCount-1 do
//
//	begin
//	ADict:=GlobalXRefs[i].Dict;
//	if (ADict.ElementCount>0) then
//	  if (ADict.Values[0] is TPDFName) and ((ADict.Values[0] as TPDFName).Name= 'Page') then
//	    begin
//	    ADict:=ADict.ValueByName('Resources') as TPDFDictionary;
//	    ADict:=ADict.ValueByName('Font') as TPDFDictionary;
//	    ADict.AddReference(AName,ANum);
//	    end;
//	end;
func debug(e Encoder, msg string) {
	fmt.Print(msg)
	if e != nil {
		e.Encode(NewfTestWriter())
		fmt.Println("")
	}
}

// AddFontNameToPages for each page adds a ref to its Resources.Font dict
func (doc *Document) AddFontNameToPages(aName string, aNum int) {
	for _, xr := range doc.GlobalXRefs {
		dict := xr.Dict
		if dict.ElementCount() > 0 {
			debug(dict.Elements[0].Value, "value of first element of dict: ")
			if v, ok := dict.Elements[0].Value.(*PDFName); ok && v.Name == "Page" {
				fmt.Println("adding font to page:", aName, aNum)
				res := dict.FindElement("Resources").Value.(*Dictionary)
				res = res.FindElement("Font").Value.(*Dictionary)
				res.AddReference(aName, aNum)
			}
		}
	}
}

func (doc *Document) CreateStdFont(embeddedFontName string, embeddedFontNum int) {
	lFontXRef := doc.xrefCount()
	dict := doc.CreateGlobalXRef("Font").Dict
	dict.AddName("Subtype", "Type1")
	dict.AddName("Encoding", "WinAnsiEncoding")
	dict.AddInteger("FirstChar", 32)
	dict.AddInteger("LastChar", 255)
	dict.AddName("BaseFont", embeddedFontName)
	n := NewPDFName(fmt.Sprintf("F%d", embeddedFontNum))
	dict.AddElement("Name", n)
	doc.AddFontNameToPages(n.Name, lFontXRef)
	doc.FontFiles = append(doc.FontFiles, "")
}

// func (doc *TPDFDocument) LoadFont(aFont *TPDFFont) bool {
// 	lFName := ""
// 	if path.Dir(aFont.FontFile) != "" {
// 		lFName = aFont.FontFile
// 	} else {
// 		lFName = path.Join(doc.FontDirectory, aFont.FontFile)
// 	}

// 	if _, err := os.Stat(lFName); err == nil {
// 		s := strings.ToLower(path.Ext(lFName))
// 		return s == ".ttf" || s == ".otf"
// 	} else {
// 		panic(fmt.Sprintf("File missing: %s", lFName))
// 	}
// }

// func (doc *TPDFDocument) CreateTTFFont(embeddedFontNum int) {
//     lFontXRef := doc.xrefCount()
//     dict := doc.CreateGlobalXRef().Dict
//     dict.AddName("Type", "Font")
//     dict.AddName("Subtype", "Type0")
//     if doc.Options&poSubsetFont != 0 {
//         dict.AddName("BaseFont", doc.GetFontNamePrefix(embeddedFontNum)+doc.Fonts[embeddedFontNum].Name)
//     } else {
//         dict.AddName("BaseFont", doc.Fonts[embeddedFontNum].Name)
//     }
//     dict.AddName("Encoding", "Identity-H")
//     n := NewPDFName(fmt.Sprintf("F%d", embeddedFontNum))
//     dict.AddElement("Name", n)
//     doc.AddFontNameToPages(n.Name, lFontXRef)
//     arr := NewTPDFArray(doc)
//     arr.AddItem(TPDFReference{doc, doc.xrefCount()})
//     dict.AddElement("DescendantFonts", arr)
//     doc.CreateTTFDescendantFont(embeddedFontNum)
//     if doc.Options&poNoEmbeddedFonts == 0 {
//         dict.AddReference("ToUnicode", doc.xrefCount())
//         doc.CreateToUnicode(embeddedFontNum)
//     }
//     doc.FontFiles = append(doc.FontFiles, doc.Fonts[embeddedFontNum].FTrueTypeFile.Filename)
// }

// func (doc *TPDFDocument) CreateTTFDescendantFont(embeddedFontNum int) {
//     dict := doc.CreateGlobalXRef().Dict
//     dict.AddName("Type", "Font")
//     dict.AddName("Subtype", "CIDFontType2")
//     if doc.Options&poSubsetFont != 0 {
//         dict.AddName("BaseFont", doc.GetFontNamePrefix(embeddedFontNum)+doc.Fonts[embeddedFontNum].Name)
//     } else {
//         dict.AddName("BaseFont", doc.Fonts[embeddedFontNum].Name)
//     }
//     dict.AddReference("CIDSystemInfo", doc.xrefCount())
//     doc.CreateTTFCIDSystemInfo()
//     dict.AddReference("FontDescriptor", doc.xrefCount())
//     doc.CreateFontDescriptor(embeddedFontNum)
//     arr := NewTPDFArray(doc)
//     dict.AddElement("W", arr)
//     arr.AddItem(TPDFTrueTypeCharWidths{doc, embeddedFontNum})
//     if doc.Options&poSubsetFont != 0 {
//         dict.AddReference("CIDToGIDMap", doc.CreateCIDToGIDMap(embeddedFontNum))
//     }
// }

// func (doc *TPDFDocument) CreateTTFCIDSystemInfo() {
//     dict := doc.CreateGlobalXRef().Dict
//     dict.AddString("Registry", "Adobe")
//     dict.AddString("Ordering", "Identity")
//     dict.AddInteger("Supplement", 0)
// }

// func (doc *TPDFDocument) CreateFontDescriptor(embeddedFontNum int) {
//     dict := doc.CreateGlobalXRef().Dict
//     dict.AddName("Type", "FontDescriptor")
//     if doc.Options&poSubsetFont != 0 {
//         dict.AddName("FontName", doc.GetFontNamePrefix(embeddedFontNum)+doc.Fonts[embeddedFontNum].Name)
//     } else {
//         dict.AddName("FontName", doc.Fonts[embeddedFontNum].Name)
//     }
//     dict.AddInteger("Flags", doc.Fonts[embeddedFontNum].Flags)
//     dict.AddInteger("ItalicAngle", 0)
//     dict.AddInteger("Ascent", doc.Fonts[embeddedFontNum].Ascent)
//     dict.AddInteger("Descent", doc.Fonts[embeddedFontNum].Descent)
//     dict.AddInteger("CapHeight", doc.Fonts[embeddedFontNum].CapHeight)
//     dict.AddInteger("StemV", 80)
//     arr := NewTPDFArray(doc)
//     arr.AddItem(NewTPDFInteger(doc,doc.Fonts[embeddedFontNum].FTrueTypeFile.FontBox[0]))
//     arr.AddItem(NewTPDFInteger(doc,doc.Fonts[embeddedFontNum].FTrueTypeFile.FontBox[1]))
//     arr.AddItem(NewTPDFInteger(doc,doc.Fonts[embeddedFontNum].FTrueTypeFile.FontBox[2]))
//     arr.AddItem(NewTPDFInteger(doc,doc.Fonts[embeddedFontNum].FTrueTypeFile.FontBox[3]))
//     dict.AddElement("FontBBox", arr)
//     if doc.Options&poNoEmbeddedFonts == 0 {
//         dict.AddReference("FontFile2", doc.xrefCount())
//         doc.CreateFontFile(embeddedFontNum)
//     }
// }

// func (doc *TPDFDocument) CreateTp1Font(EmbeddedFontNum int) {
//     if EmbeddedFontNum == -1 {
//         panic("Assertion failed: EmbeddedFontNum is -1")
//     }
// }

// func (doc *TPDFDocument) CreateFontDescriptor(EmbeddedFontNum int) {
//     dict := doc.CreateGlobalXRef().Dict
//     dict.AddName("Type", "FontDescriptor")

//     if doc.Options&poSubsetFont != 0 {
//         dict.AddName("FontName", doc.GetFontNamePrefix(EmbeddedFontNum)+doc.Fonts[EmbeddedFontNum].Name)
//         dict.AddInteger("Flags", 4)
//     } else {
//         dict.AddName("FontName", doc.Fonts[EmbeddedFontNum].Name)
//         dict.AddName("FontFamily", doc.Fonts[EmbeddedFontNum].FTrueTypeFile.FamilyName)
//         dict.AddInteger("Flags", 32)
//     }

//     dict.AddInteger("Ascent", doc.Fonts[EmbeddedFontNum].FTrueTypeFile.Ascender)
//     dict.AddInteger("Descent", doc.Fonts[EmbeddedFontNum].FTrueTypeFile.Descender)
//     dict.AddInteger("CapHeight", doc.Fonts[EmbeddedFontNum].FTrueTypeFile.CapHeight)
//     Arr := NewTPDFArray(doc)
//     dict.AddElement("FontBBox", Arr)
//     Arr.AddIntArray(doc.Fonts[EmbeddedFontNum].FTrueTypeFile.BBox)
//     dict.AddInteger("ItalicAngle", int(doc.Fonts[EmbeddedFontNum].FTrueTypeFile.ItalicAngle))
//     dict.AddInteger("StemV", doc.Fonts[EmbeddedFontNum].FTrueTypeFile.StemV)
//     dict.AddInteger("MissingWidth", doc.Fonts[EmbeddedFontNum].FTrueTypeFile.MissingWidth)

//     if doc.Options&poNoEmbeddedFonts == 0 {
//         dict.AddReference("FontFile2", doc.xrefCount())
//         doc.CreateFontFileEntry(EmbeddedFontNum)

//         if doc.Options&poSubsetFont != 0 {
//             // todo /CIDSet reference
//             dict.AddReference("CIDSet", doc.xrefCount())
//             doc.CreateCIDSet(EmbeddedFontNum)
//         }
//     }
// }

// func (doc *TPDFDocument) CreateToUnicode(AFontNum int) {
//     lXRef := doc.CreateGlobalXRef()
//     lXRef.FStream = doc.NewPDFStream(doc)
//     lXRef.FStream.AddItem(NewTPDFToUnicode(doc, AFontNum))
// }

// func (doc *TPDFDocument) CreateFontFileEntry(AFontNum int) {
//     dict := doc.CreateGlobalXRef().Dict
//     if doc.Options&poCompressFonts != 0 {
//         dict.AddName("Filter", "FlateDecode")
//     }
//     var Len int
//     if doc.Options&poSubsetFont != 0 {
//         Len = doc.Fonts[AFontNum].SubsetFont.Size
//     } else {
//         Len = doc.Fonts[AFontNum].FTrueTypeFile.OriginalSize
//     }
//     dict.AddInteger(fmt.Sprintf("Length1 %d", AFontNum), Len)
// }

// func (doc *TPDFDocument) CreateCIDSet(AFontNum int) {
//     lXRef := doc.CreateGlobalXRef()
//     lXRef.FStream = doc.NewPDFStream(doc)
//     lXRef.FStream.AddItem(NewTPDFCIDSet(doc, AFontNum))
// }

// func (doc *TPDFDocument) CreateAnnotEntry(APageNum, AnnotNum int) int {
//     an := doc.Pages[APageNum].Annots[AnnotNum]
//     lXRef := doc.CreateGlobalXRef()
//     lDict := lXRef.Dict
//     lDict.AddName("Type", "Annot")
//     lDict.AddName("Subtype", "Link")
//     lDict.AddName("H", "I")

//     ar := NewTPDFArray(doc)
//     lDict.AddElement("Border", ar)
//     if an.FBorder {
//         ar.AddFreeFormArrayValues("0 0 1")
//     } else {
//         ar.AddFreeFormArrayValues("0 0 0")
//     }

//     ar = NewTPDFArray(doc)
//     lDict.AddElement("Rect", ar)
//     s := fmt.Sprintf("%f %f %f %f", an.FLeft, an.FBottom, an.FLeft+an.FWidth, an.FBottom+an.FHeight)
//     ar.AddFreeFormArrayValues(s)

//     ADict := doc.CreateDictionary()
//     lDict.AddElement("A", ADict)
//     if an.FExternalLink {
//         ADict.AddName("Type", "Action")
//         ADict.AddName("S", "URI")
//         ADict.AddString("URI", an.FURI)
//     } else {
//         ADict.AddName("Type", "Action")
//         ADict.AddName("S", "GoTo")
//         ADict.AddReference("D", doc.Pages[an.FPageDest].XRef)
//     }
//     return lXRef
// }

// func (doc *TPDFDocument) CreateCIDToGIDMap(AFontNum int) int {
//     lXRef := doc.CreateGlobalXRef()
//     lXRef.FStream = doc.NewPDFStream(doc)
//     lXRef.FStream.AddItem(TCIDToGIDMap{doc, AFontNum})
//     return doc.xrefCount()- 1
// }

func (doc *Document) CreateContentsEntry(APageNum int) int {
	contents := doc.CreateGlobalXRef("ContentEntry")
	contents.Stream = NewPDFStream(doc)
	result := doc.xrefCount() - 1
	//FIXME:
	i := 2 + 0 // doc.Pages[APageNum].Annots.Count() // + GetTotalAnnotsCount()
	doc.GlobalXRefs[doc.xrefCount()-i].Dict.AddReference("Contents", result)
	return result
}

func (doc *Document) CreatePageStream(APage *Page, PageNum int) {
	pageStream := doc.GlobalXRefs[PageNum].Stream
	for i := 0; i < len(APage.Objects); i++ {
		pageStream.AddItem(APage.Objects[i])
	}
}

func (doc *Document) CreateSectionPageOutLine(S *Section, PageOutLine, PageIndex, NewPage, ParentOutline, NextOutline, PrevOutLine int) int {
	aDict := doc.GlobalXRefs[ParentOutline].Dict
	aDict.IncCount()
	aDict = doc.GlobalXRefs[PageOutLine].Dict
	arr := aDict.FindElement("Dest").Value.(*Array)
	arr.AddItem(NewReference(NewPage))
	arr.AddItem(NewPDFName("Fit"))
	result := PrevOutLine
	if PageIndex == 0 {
		doc.GlobalXRefs[ParentOutline].Dict.AddReference("First", doc.xrefCount()-1)
		result = doc.xrefCount() - 1
		aDict = doc.GlobalXRefs[ParentOutline].Dict
		arr = aDict.FindElement("Dest").Value.(*Array)
		arr.AddItem(NewReference(NewPage))
		arr.AddItem(NewPDFName("Fit"))
	} else {
		doc.GlobalXRefs[NextOutline].Dict.AddReference("Next", doc.xrefCount()-1)
		doc.GlobalXRefs[PageOutLine].Dict.AddReference("Prev", PrevOutLine)
		if PageIndex < S.PageCount() {
			result = doc.xrefCount() - 1
		}
	}
	if PageIndex == S.PageCount()-1 {
		doc.GlobalXRefs[ParentOutline].Dict.AddReference("Last", doc.xrefCount()-1)
	}
	return result
}

func (doc *Document) CreateSectionOutLine(SectionIndex, OutLineRoot, ParentOutLine, NextSect, PrevSect int) int {
	aDict := doc.GlobalXRefs[OutLineRoot].Dict
	aDict.IncCount()

	if SectionIndex == 0 {
		doc.GlobalXRefs[OutLineRoot].Dict.AddReference("First", doc.xrefCount()-1)
		return doc.xrefCount() - 1
	}

	doc.GlobalXRefs[NextSect].Dict.AddReference("Next", doc.xrefCount()-1)
	doc.GlobalXRefs[ParentOutLine].Dict.AddReference("Prev", PrevSect)
	if SectionIndex < doc.Sections.Count()-1 {
		return doc.xrefCount() - 1
	}

	if SectionIndex == doc.Sections.Count()-1 {
		doc.GlobalXRefs[OutLineRoot].Dict.AddReference("Last", doc.xrefCount()-1)
	}
	return PrevSect
}

func (doc *Document) CreateSectionsOutLine() int {
	var result, treeRoot, outlineRoot, pc, j, parentOutline, pageNum, pageOutline, nextOutline, nextSect, newPage, prevOutline, prevSect int
	var aDict *Dictionary
	var arr *Array

	if doc.Sections.Count() > 1 {
		if doc.hasOption(poOutLine) {
			outlineRoot = doc.CreateOutlines()
			doc.GlobalXRefs[doc.Catalogue].Dict.AddReference("Outlines", doc.xrefCount()-1)
			doc.GlobalXRefs[doc.Catalogue].Dict.AddName("PageMode", "UseOutlines")
		}
		treeRoot = doc.CreatePagesEntry(result)
	} else {
		result = doc.CreatePagesEntry(result)
		treeRoot = result
	}

	for j = 0; j < doc.Sections.Count(); j++ {
		s := doc.Sections.Get(j)
		if doc.hasOption(poOutLine) {
			parentOutline = doc.CreateOutlineEntry(outlineRoot, j+1, -1, s.Title)
			prevSect = doc.CreateSectionOutLine(j, outlineRoot, parentOutline, nextSect, prevSect)
			nextSect = parentOutline
			result = doc.CreatePagesEntry(treeRoot)
		}
		for k := 0; k < s.PageCount(); k++ {
			newPage = doc.CreatePageEntry(result, k)
			if j == 0 && k == 0 {
				aDict = doc.GlobalXRefs[doc.Catalogue].Dict
				arr = aDict.FindElement("OpenAction").Value.(*Array)
				arr.AddItem(NewReference(doc.xrefCount() - 1))
				arr.AddItem(NewPDFNameEx(fmt.Sprintf("XYZ null null %f", float64(doc.ZoomValue/100)), false))
			}
			pageNum = doc.CreateContentsEntry(k)
			doc.CreatePageStream(s.Pages[k], pageNum)
			if doc.Sections.Count() > 1 && doc.hasOption(poOutLine) {
				pageOutline = doc.CreateOutlineEntry(parentOutline, j+1, k+1, s.Title)
				doc.CreateSectionPageOutLine(s, pageOutline, k, newPage, parentOutline, nextOutline, prevOutline)
				nextOutline = pageOutline
			}
		}
	}

	aDict = doc.GlobalXRefs[treeRoot].Dict
	pc = 0
	for j := 0; j < doc.Sections.Count(); j++ {
		pc += doc.Sections.Get(j).PageCount()
	}
	aDict.FindElement("Count").Value.(*Integer).val = pc
	return result
}

func (doc *Document) CreateFontEntries() {
	numFont := 0
	for _, font := range doc.Fonts {
		fontName := font.Name
		if IsStandardFont(fontName) {
			doc.CreateStdFont(fontName, numFont)
		}
		// } else if doc.LoadFont(font) {
		// 	if poSubsetFont&Options != 0 {
		// 		font.GenerateSubsetFont()
		// 	}
		// 	doc.CreateTtfFont(&numFont)
		// } else {
		// 	doc.CreateTp1Font(&numFont)
		// }
		numFont++
	}
}

// func (doc *TPDFDocument) CreateAnnotEntries(APageNum int, APageDict *TPDFDictionary) {
// 	if doc.GetTotalAnnotsCount() == 0 {
// 		return
// 	}
// 	ar := NewTPDFArray(doc)
// 	APageDict.AddElement("Annots", ar)
// 	for i := range doc.Pages[APageNum].Annots {
// 		refnum := doc.CreateAnnotEntry(APageNum, i)
// 		ar.AddItem(NewTPDFReference(doc, refnum))
// 	}
// }

func (doc *Document) WriteXRefTable(st PDFWriter) {
	for _, xr := range doc.GlobalXRefs {
		xr.Encode(st)
	}
}

func (doc *Document) CreateEmbeddedFont(APage *Page, AFontIndex int, AFontSize float64, ASimulateBold, ASimulateItalic bool) *EmbeddedFont {
	return NewEmbeddedFontEx(doc, APage, AFontIndex, AFontSize, ASimulateBold, ASimulateItalic)
}

func (doc *Document) FindFont(AName string) int {
	for i, f := range doc.Fonts {
		if f.Name == AName {
			return i
		}
	}
	return -1
}

// AddFont adds a  font to the PDF document and returns its index.
// If the font already exists, it returns the existing font's index.
// If filename is empty, font is assumed to be a standard font
func (doc *Document) AddFont(name, filename string) int {
	idx := doc.FindFont(name)
	if idx > -1 {
		return idx
	}
	f := NewFont(name, filename == "") //std font
	doc.Fonts = append(doc.Fonts, f)
	fontNum := doc.FontCount() - 1
	if f.IsStdFont {
		return fontNum
	}
	// non-standard font
	if path.Dir(filename) == "" { //just a filename
		filename = filepath.Join(doc.FontDir, filename)
	}
	f.FontFilename = filename
	return fontNum
}
