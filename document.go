package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"time"
)

type TPDFDocument struct {
	Catalogue          int
	CurrentColor       string
	CurrentWidth       string
	LineCapStyle       TPDFLineCapStyle
	DefaultOrientation TPDFPaperOrientation
	DefaultPaperType   TPDFPaperType
	FontDirectory      string
	FontFiles          []string
	Fonts              []*TPDFFont
	// FImages              []*TPDFImage
	Infos                *TPDFInfos
	LineStyleDefs        TPDFLineStyleDef
	ObjectCount          int
	Options              []TPDFOption
	Pages                []*TPDFPage
	Preferences          bool
	PageLayout           TPDFPageLayout
	Sections             TPDFSectionList
	Trailer              TPDFDictionary
	ZoomValue            int
	GlobalXRefs          []*TPDFXRef
	UnitOfMeasure        TPDFUnitOfMeasure
	DefaultUnitOfMeasure TPDFUnitOfMeasure
}

func NewTPDFDocument() *TPDFDocument {
	doc := &TPDFDocument{
		// doc.FontFiles = make([]string, 0)
		Preferences:          true,
		PageLayout:           lSingle,
		DefaultPaperType:     ptA4,
		DefaultOrientation:   ppoPortrait,
		ZoomValue:            100,
		Options:              []TPDFOption{poCompressFonts, poCompressImages},
		DefaultUnitOfMeasure: uomMillimeters,
		LineCapStyle:         plcsRoundCap,
		Sections:             *NewTPDFSectionList(),
		Infos:                NewTPDFInfos(),
	}
	return doc
}

func (d *TPDFDocument) SetOptions(AValue ...TPDFOption) {
	if slices.Contains(AValue, poNoEmbeddedFonts) {
		AValue = append(AValue, poSubsetFont)
	}
	d.Options = AValue
}

func (doc *TPDFDocument) hasOption(opt TPDFOption) bool {
	return slices.Contains(doc.Options, opt)
}

func (doc *TPDFDocument) StartDocument() {
	// doc.Reset()
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
	//FIXME:
	// if doc.FontDirectory == "" {
	// 	doc.FontDirectory = path.Dir(os.Args[0])
	// }
}

func (doc *TPDFDocument) SaveToWriter(st PDFWriter) error {

	doc.CreateSectionsOutLine()
	// doc.CreateFontEntries()
	// doc.CreateImageEntries()
	doc.Trailer.SetElement("Size", NewTPDFInteger(doc, doc.xrefCount()))
	st.Write([]byte(PDF_VERSION + "\n"))
	st.Write([]byte(PDF_BINARY_BLOB + "\n"))
	xRefPos := st.Offset()
	for i := 1; i < doc.xrefCount(); i++ {
		xRefPos = st.Offset()
		doc.WriteObject(i, st)
		doc.GlobalXRefs[i].Offset = xRefPos
	}
	st.Write([]byte(fmt.Sprintf("xref\n0 %d\n", doc.xrefCount())))
	doc.WriteXRefTable(st)
	st.Write([]byte("trailer\n"))
	doc.Trailer.Encode(st)
	st.Write([]byte(fmt.Sprintf("\nstartxref\n%d\n", xRefPos)))
	st.Write([]byte(PDF_FILE_END))
	return st.Err()
}

func (doc *TPDFDocument) SaveToFile(AFileName string) error {
	f, err := os.Create(AFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	st := Newfwriter(f)
	return doc.SaveToWriter(st)
}

func (d *TPDFDocument) PageCount() int {
	return len(d.Pages)
}
func (d *TPDFDocument) NewPage() *TPDFPage {
	d.Pages = append(d.Pages, NewTPDFPage(d))
	return d.Pages[d.PageCount()-1]
}

func (d *TPDFDocument) xrefCount() int {
	return len(d.GlobalXRefs)
}

func (d *TPDFDocument) AddGlobalXRef(AXRef *TPDFXRef) int {
	d.GlobalXRefs = append(d.GlobalXRefs, AXRef)
	return len(d.GlobalXRefs) - 1
}

// FIXME: check this is correct
// the first object has id of 0 (not used )
func (doc *TPDFDocument) CreateRefTable() {
	eadd(&doc.GlobalXRefs, NewTPDFXRef(doc))
}

func (doc *TPDFDocument) CreateGlobalXRef() *TPDFXRef {
	xref := NewTPDFXRef(doc)
	doc.AddGlobalXRef(xref)
	return xref
}

func (d *TPDFDocument) GetX(AIndex int) *TPDFXRef {
	return d.GlobalXRefs[AIndex]
}

//FIXME:
// func (d *TPDFDocument) AddXObject(AXObject *TPDFXObject) int {
//     d.FXObjects = append(d.FXObjects, AXObject)
//     return len(d.FXObjects) - 1
// }
//FIXME:
// func (d *TPDFDocument) AddPattern(APattern *TPDFPattern) int {
//     d.FPatterns = append(d.FPatterns, APattern)
//     return len(d.FPatterns) - 1
// }

// func (d *TPDFDocument) AddImage(AImage *TPDFImage) int {
//     d.FImages = append(d.FImages, AImage)
//     return len(d.FImages) - 1
// }

// func (d *TPDFDocument) AddGraphicState(AGraphicState *TPDFGraphicState) int {
//     d.FGraphicStates = append(d.FGraphicStates, AGraphicState)
//     return len(d.FGraphicStates) - 1
// }

// func (d *TPDFDocument) AddAnnotation(AAnnotation *TPDFAnnotation) int {
//     d.FAnnotations = append(d.FAnnotations, AAnnotation)
//     return len(d.FAnnotations) - 1
// }

func (doc *TPDFDocument) CreateTrailer() {
	doc.Trailer = *NewTPDFDictionary(doc)
	doc.Trailer.AddInteger("Size", doc.xrefCount())
}

func (doc *TPDFDocument) CreateCatalogEntry() int {
	CDict := doc.CreateGlobalXRef().Dict
	doc.Trailer.AddReference("Root", doc.xrefCount()-1)
	CDict.AddName("Type", "Catalog")
	CDict.AddName("PageLayout", PageLayoutNames[doc.PageLayout])
	CDict.AddElement("OpenAction", NewTPDFArray(doc))
	return doc.xrefCount() - 1
}

func (doc *TPDFDocument) CreateInfoEntry(UseUTF16 bool) {
	IDict := doc.CreateGlobalXRef().Dict
	doc.Trailer.AddReference("Info", doc.xrefCount()-1)
	doc.Trailer.SetElement("Size", NewTPDFInteger(doc, doc.xrefCount()))
	noUnicode := false
	doEntry := func(aName, aValue string) {
		if aValue == "" {
			return
		}
		if UseUTF16 && !noUnicode {
			//FIXME:
			// IDict.AddString(aName, utf8Decode(aValue))
		} else {
			IDict.AddString(aName, aValue)
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

func (doc *TPDFDocument) CreateMetadataEntry() {
	lXRef := doc.CreateGlobalXRef()
	lXRef.Dict.AddName("Type", "Metadata")
	lXRef.Dict.AddName("Subtype", "XML")
	lXRef.Stream = NewPDFStream(doc)
	lXRef.Stream.AddItem(NewTXMPStream(doc))
	lXRef.Stream.CompressionProhibited = true

	doc.GlobalXRefs[doc.Catalogue].Dict.AddReference("Metadata", doc.xrefCount()-1)
}

func (doc *TPDFDocument) AddOutputIntent(Subtype, OutputConditionIdentifier, Info string,
	ICCProfile io.Reader) {
	OIRef := doc.xrefCount()
	OIDict := doc.CreateGlobalXRef().Dict
	OIDict.AddName("Type", "OutputIntent")
	OIDict.AddName("S", Subtype)
	OIDict.AddString("OutputConditionIdentifier", OutputConditionIdentifier)
	if Info != "" {
		OIDict.AddString("Info", Info)
	}
	if ICCProfile != nil {
		Profile := doc.CreateGlobalXRef()
		Profile.Dict.AddInteger("N", 3)
		Profile.Stream = NewPDFStream(doc)
		Profile.Stream.AddItem(NewPDFMemoryStream(doc, ICCProfile))
		OIDict.AddReference("DestOutputProfile", doc.xrefCount()-1)
	}

	OutputIntents := doc.GlobalXRefs[doc.Catalogue].Dict.ElementByName("OutputIntents")
	if OutputIntents == nil {
		OutputIntents = doc.GlobalXRefs[doc.Catalogue].Dict.AddElement("OutputIntents", NewTPDFArray(doc))
	}
	OutputIntents.Value.(*TPDFArray).AddItem(NewTPDFReference(doc, OIRef))
}

func (doc *TPDFDocument) AddPDFA1sRGBOutputIntent() {
	var buf bytes.Buffer
	buf.Grow(len(ICC_sRGB2014) - 1)
	buf.Write(ICC_sRGB2014[1:])
	doc.AddOutputIntent("GTS_PDFA1", "Custom", "sRGB", &buf)
}

func (doc *TPDFDocument) CreateTrailerID() {
	s := DateToPdfDate(time.Now()) + strconv.Itoa(doc.xrefCount()) +
		doc.Infos.Title + doc.Infos.Author + doc.Infos.ApplicationName + doc.Infos.Producer + DateToPdfDate(doc.Infos.CreationDate)
	s = GetMD5Hash(s)
	ID := NewTPDFArray(doc)
	ID.AddItem(NewTPDFRawHexString(doc, s))
	ID.AddItem(NewTPDFRawHexString(doc, s))
	doc.Trailer.AddElement("ID", ID)
}

func (doc *TPDFDocument) CreatePreferencesEntry() {
	VDict := doc.CreateGlobalXRef().Dict
	VDict.AddName("Type", "ViewerPreferences")
	VDict.AddElement("FitWindow", NewTPDFBoolean(doc, true))
	doc.GlobalXRefs[doc.Catalogue].Dict.AddReference("ViewerPreferences", doc.xrefCount()-1)
	//FIXME: confirm that the lien above repalces the two below
	// VDict = doc.GlobalXRefByName("Catalog").Dict
	// VDict.AddReference("ViewerPreferences", doc.xrefCount()-1)
}

func (doc *TPDFDocument) CreatePagesEntry(Parent int) int {
	EDict := doc.CreateGlobalXRef().Dict
	result := doc.xrefCount() - 1
	EDict.AddName("Type", "Pages")
	EDict.AddElement("Kids", NewTPDFArray(doc))
	EDict.AddInteger("Count", 0)
	if Parent == 0 {
		doc.GlobalXRefs[doc.Catalogue].Dict.AddReference("Pages", result)
		//FIXME: confirm that the lien above repalces the line below
		// GlobalXRefByName("Catalog").Dict.AddReference("Pages", result)
	} else {
		EDict.AddReference("Parent", Parent)
		ADict := doc.GlobalXRefs[Parent].Dict
		ADict.IncCount()
		ADict.AddKid(NewTPDFReference(doc, result))
	}
	return result
}

func (doc *TPDFDocument) CreatePageEntry(parent, pageNum int) int {
	pp := doc.Pages[pageNum]
	pDict := doc.CreateGlobalXRef().Dict
	pDict.AddName("Type", "Page")
	pDict.AddReference("Parent", parent)
	aDict := doc.GlobalXRefs[parent].Dict
	aDict.IncCount()
	aDict.AddKid(NewTPDFReference(doc, doc.xrefCount()-1))
	arr := NewTPDFArray(doc)
	arr.AddItem(NewTPDFInteger(doc, 0))
	arr.AddItem(NewTPDFInteger(doc, 0))
	arr.AddItem(NewTPDFInteger(doc, pp.Paper.W))
	arr.AddItem(NewTPDFInteger(doc, pp.Paper.H))
	pDict.AddElement("MediaBox", arr)
	// doc.CreateAnnotEntries(pageNum, pDict)
	aDict = NewTPDFDictionary(doc)
	pDict.AddElement("Resources", aDict)
	arr = NewTPDFArray(doc) // procset
	aDict.AddElement("ProcSet", arr)
	arr.AddItem(NewPDFName(doc, "PDF", true))
	arr.AddItem(NewPDFName(doc, "Text", true))
	arr.AddItem(NewPDFName(doc, "ImageC", true))
	if len(doc.Fonts) > 0 {
		aDict.AddElement("Font", NewTPDFDictionary(doc))
	}
	// if pp.HasImages {
	// 	aDict.AddElement("XObject", doc.CreateDictionary())
	// }

	return doc.xrefCount() - 1
}

func (doc *TPDFDocument) CreateOutlines() int {
	oDict := doc.CreateGlobalXRef().Dict
	oDict.AddName("Type", "Outlines")
	oDict.AddInteger("Count", 0)
	return doc.xrefCount() - 1
}

func (doc *TPDFDocument) CreateOutlineEntry(parent, sectNo, pageNo int, aTitle string) int {
	oDict := doc.CreateGlobalXRef().Dict
	s := aTitle
	if s == "" {
		s = fmt.Sprintf("Section %d", sectNo)
	}
	if pageNo > -1 {
		s = fmt.Sprintf("%s Page %d", s, pageNo)
	}
	oDict.AddString("Title", s)
	oDict.AddReference("Parent", parent)
	oDict.AddInteger("Count", 0)
	oDict.AddElement("Dest", NewTPDFArray(doc))
	return doc.xrefCount() - 1
}

// func (doc *TPDFDocument) AddFontNameToPages(aName string, aNum int) {
// 	for i := 1; i < doc.xrefCount(); i++ {
// 		aDict := doc.GlobalXRefs[i].Dict
// 		if aDict.ElementCount() > 0 {
// 			if v, ok := aDict.Values[0].(*TPDFName); ok && v.Name == "Page" {
// 				aDict = aDict.ValueByName("Resources").(*TPDFDictionary)
// 				aDict = aDict.ValueByName("Font").(*TPDFDictionary)
// 				aDict.AddReference(aName, aNum)
// 			}
// 		}
// 	}
// }

// func (doc *TPDFDocument) CreateStdFont(embeddedFontName string, embeddedFontNum int) {
// 	lFontXRef := doc.xrefCount()
// 	fDict := doc.CreateGlobalXRef().Dict
// 	fDict.AddName("Type", "Font")
// 	fDict.AddName("Subtype", "Type1")
// 	fDict.AddName("Encoding", "WinAnsiEncoding")
// 	fDict.AddInteger("FirstChar", 32)
// 	fDict.AddInteger("LastChar", 255)
// 	fDict.AddName("BaseFont", embeddedFontName)
// 	n := NewPDFName(doc, fmt.Sprintf("F%d", embeddedFontNum))
// 	fDict.AddElement("Name", n)
// 	doc.AddFontNameToPages(n.Name, lFontXRef)

// 	doc.FontFiles = append(doc.FontFiles, "")
// }

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
//     fDict := doc.CreateGlobalXRef().Dict
//     fDict.AddName("Type", "Font")
//     fDict.AddName("Subtype", "Type0")
//     if doc.Options&poSubsetFont != 0 {
//         fDict.AddName("BaseFont", doc.GetFontNamePrefix(embeddedFontNum)+doc.Fonts[embeddedFontNum].Name)
//     } else {
//         fDict.AddName("BaseFont", doc.Fonts[embeddedFontNum].Name)
//     }
//     fDict.AddName("Encoding", "Identity-H")
//     n := NewPDFName(doc,fmt.Sprintf("F%d", embeddedFontNum))
//     fDict.AddElement("Name", n)
//     doc.AddFontNameToPages(n.Name, lFontXRef)
//     arr := NewTPDFArray(doc)
//     arr.AddItem(TPDFReference{doc, doc.xrefCount()})
//     fDict.AddElement("DescendantFonts", arr)
//     doc.CreateTTFDescendantFont(embeddedFontNum)
//     if doc.Options&poNoEmbeddedFonts == 0 {
//         fDict.AddReference("ToUnicode", doc.xrefCount())
//         doc.CreateToUnicode(embeddedFontNum)
//     }
//     doc.FontFiles = append(doc.FontFiles, doc.Fonts[embeddedFontNum].FTrueTypeFile.Filename)
// }

// func (doc *TPDFDocument) CreateTTFDescendantFont(embeddedFontNum int) {
//     fDict := doc.CreateGlobalXRef().Dict
//     fDict.AddName("Type", "Font")
//     fDict.AddName("Subtype", "CIDFontType2")
//     if doc.Options&poSubsetFont != 0 {
//         fDict.AddName("BaseFont", doc.GetFontNamePrefix(embeddedFontNum)+doc.Fonts[embeddedFontNum].Name)
//     } else {
//         fDict.AddName("BaseFont", doc.Fonts[embeddedFontNum].Name)
//     }
//     fDict.AddReference("CIDSystemInfo", doc.xrefCount())
//     doc.CreateTTFCIDSystemInfo()
//     fDict.AddReference("FontDescriptor", doc.xrefCount())
//     doc.CreateFontDescriptor(embeddedFontNum)
//     arr := NewTPDFArray(doc)
//     fDict.AddElement("W", arr)
//     arr.AddItem(TPDFTrueTypeCharWidths{doc, embeddedFontNum})
//     if doc.Options&poSubsetFont != 0 {
//         fDict.AddReference("CIDToGIDMap", doc.CreateCIDToGIDMap(embeddedFontNum))
//     }
// }

// func (doc *TPDFDocument) CreateTTFCIDSystemInfo() {
//     fDict := doc.CreateGlobalXRef().Dict
//     fDict.AddString("Registry", "Adobe")
//     fDict.AddString("Ordering", "Identity")
//     fDict.AddInteger("Supplement", 0)
// }

// func (doc *TPDFDocument) CreateFontDescriptor(embeddedFontNum int) {
//     fDict := doc.CreateGlobalXRef().Dict
//     fDict.AddName("Type", "FontDescriptor")
//     if doc.Options&poSubsetFont != 0 {
//         fDict.AddName("FontName", doc.GetFontNamePrefix(embeddedFontNum)+doc.Fonts[embeddedFontNum].Name)
//     } else {
//         fDict.AddName("FontName", doc.Fonts[embeddedFontNum].Name)
//     }
//     fDict.AddInteger("Flags", doc.Fonts[embeddedFontNum].Flags)
//     fDict.AddInteger("ItalicAngle", 0)
//     fDict.AddInteger("Ascent", doc.Fonts[embeddedFontNum].Ascent)
//     fDict.AddInteger("Descent", doc.Fonts[embeddedFontNum].Descent)
//     fDict.AddInteger("CapHeight", doc.Fonts[embeddedFontNum].CapHeight)
//     fDict.AddInteger("StemV", 80)
//     arr := NewTPDFArray(doc)
//     arr.AddItem(NewTPDFInteger(doc,doc.Fonts[embeddedFontNum].FTrueTypeFile.FontBox[0]))
//     arr.AddItem(NewTPDFInteger(doc,doc.Fonts[embeddedFontNum].FTrueTypeFile.FontBox[1]))
//     arr.AddItem(NewTPDFInteger(doc,doc.Fonts[embeddedFontNum].FTrueTypeFile.FontBox[2]))
//     arr.AddItem(NewTPDFInteger(doc,doc.Fonts[embeddedFontNum].FTrueTypeFile.FontBox[3]))
//     fDict.AddElement("FontBBox", arr)
//     if doc.Options&poNoEmbeddedFonts == 0 {
//         fDict.AddReference("FontFile2", doc.xrefCount())
//         doc.CreateFontFile(embeddedFontNum)
//     }
// }

// func (doc *TPDFDocument) CreateTp1Font(EmbeddedFontNum int) {
//     if EmbeddedFontNum == -1 {
//         panic("Assertion failed: EmbeddedFontNum is -1")
//     }
// }

// func (doc *TPDFDocument) CreateFontDescriptor(EmbeddedFontNum int) {
//     FDict := doc.CreateGlobalXRef().Dict
//     FDict.AddName("Type", "FontDescriptor")

//     if doc.Options&poSubsetFont != 0 {
//         FDict.AddName("FontName", doc.GetFontNamePrefix(EmbeddedFontNum)+doc.Fonts[EmbeddedFontNum].Name)
//         FDict.AddInteger("Flags", 4)
//     } else {
//         FDict.AddName("FontName", doc.Fonts[EmbeddedFontNum].Name)
//         FDict.AddName("FontFamily", doc.Fonts[EmbeddedFontNum].FTrueTypeFile.FamilyName)
//         FDict.AddInteger("Flags", 32)
//     }

//     FDict.AddInteger("Ascent", doc.Fonts[EmbeddedFontNum].FTrueTypeFile.Ascender)
//     FDict.AddInteger("Descent", doc.Fonts[EmbeddedFontNum].FTrueTypeFile.Descender)
//     FDict.AddInteger("CapHeight", doc.Fonts[EmbeddedFontNum].FTrueTypeFile.CapHeight)
//     Arr := NewTPDFArray(doc)
//     FDict.AddElement("FontBBox", Arr)
//     Arr.AddIntArray(doc.Fonts[EmbeddedFontNum].FTrueTypeFile.BBox)
//     FDict.AddInteger("ItalicAngle", int(doc.Fonts[EmbeddedFontNum].FTrueTypeFile.ItalicAngle))
//     FDict.AddInteger("StemV", doc.Fonts[EmbeddedFontNum].FTrueTypeFile.StemV)
//     FDict.AddInteger("MissingWidth", doc.Fonts[EmbeddedFontNum].FTrueTypeFile.MissingWidth)

//     if doc.Options&poNoEmbeddedFonts == 0 {
//         FDict.AddReference("FontFile2", doc.xrefCount())
//         doc.CreateFontFileEntry(EmbeddedFontNum)

//         if doc.Options&poSubsetFont != 0 {
//             // todo /CIDSet reference
//             FDict.AddReference("CIDSet", doc.xrefCount())
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
//     FDict := doc.CreateGlobalXRef().Dict
//     if doc.Options&poCompressFonts != 0 {
//         FDict.AddName("Filter", "FlateDecode")
//     }
//     var Len int
//     if doc.Options&poSubsetFont != 0 {
//         Len = doc.Fonts[AFontNum].SubsetFont.Size
//     } else {
//         Len = doc.Fonts[AFontNum].FTrueTypeFile.OriginalSize
//     }
//     FDict.AddInteger(fmt.Sprintf("Length1 %d", AFontNum), Len)
// }

// func (doc *TPDFDocument) CreateCIDSet(AFontNum int) {
//     lXRef := doc.CreateGlobalXRef()
//     lXRef.FStream = doc.NewPDFStream(doc)
//     lXRef.FStream.AddItem(NewTPDFCIDSet(doc, AFontNum))
// }

// func (doc *TPDFDocument) CreateImageEntry(ImgWidth, ImgHeight, NumImg int) *TPDFDictionary {
//     lXRef := doc.xrefCount()// reference to be used later

//     ImageDict := doc.CreateGlobalXRef().Dict
//     ImageDict.AddName("Type", "XObject")
//     ImageDict.AddName("Subtype", "Image")
//     ImageDict.AddInteger("Width", ImgWidth)
//     ImageDict.AddInteger("Height", ImgHeight)
//     ImageDict.AddName("ColorSpace", "DeviceRGB")
//     ImageDict.AddInteger("BitsPerComponent", 8)
//     N := NewPDFName(doc,fmt.Sprintf("I%d", NumImg)) // Needed later
//     ImageDict.AddElement("Name", N)

//     // now find where we must add the image xref - we are looking for "Resources"
//     for i := 1; i < doc.xrefCount(); i++ {
//         ADict := doc.GlobalXRefs[i].Dict
//         if len(ADict.Values) > 0 {
//             if val, ok := ADict.Values[0].(*TPDFName); ok && val.Name == "Page" {
//                 ADict = ADict.ValueByName("Resources").(*TPDFDictionary)
//                 ADict = ADict.FindValue("XObject").(*TPDFDictionary)
//                 if ADict != nil {
//                     ADict.AddReference(N.Name, lXRef)
//                 }
//             }
//         }
//     }

//     return ImageDict
// }

// func (doc *TPDFDocument) CreateImageMaskEntry(ImgWidth, ImgHeight, NumImg int, ImageDict *TPDFDictionary) {
//     lXRef := doc.xrefCount()// reference to be used later

//     MDict := doc.CreateGlobalXRef().Dict
//     MDict.AddName("Type", "XObject")
//     MDict.AddName("Subtype", "Image")
//     MDict.AddInteger("Width", ImgWidth)
//     MDict.AddInteger("Height", ImgHeight)
//     MDict.AddName("ColorSpace", "DeviceGray")
//     MDict.AddInteger("BitsPerComponent", 8)
//     N := NewPDFName(doc,fmt.Sprintf("M%d", NumImg)) // Needed later
//     MDict.AddElement("Name", N)
//     ImageDict.AddReference("SMask", lXRef)
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

func (doc *TPDFDocument) CreateContentsEntry(APageNum int) int {
	contents := doc.CreateGlobalXRef()
	contents.Stream = NewPDFStream(doc)
	result := doc.xrefCount() - 1
	//FIXME:
	i := 2 + 0 // doc.Pages[APageNum].Annots.Count() // + GetTotalAnnotsCount()
	doc.GlobalXRefs[doc.xrefCount()-i].Dict.AddReference("Contents", result)
	return result
}

func (doc *TPDFDocument) CreatePageStream(APage *TPDFPage, PageNum int) {
	pageStream := doc.GlobalXRefs[PageNum].Stream
	for i := 0; i < len(APage.Objects); i++ {
		pageStream.AddItem(APage.Objects[i])
	}
}

// func (doc *TPDFDocument) ImageStreamOptions() TPDFImageStreamOptions {
// 	var result TPDFImageStreamOptions
// 	if doc.Options.contains(poCompressImages) {
// 		result = append(result, isoCompressed)
// 	}
// 	if doc.Options.contains(poUseImageTransparency) {
// 		result = append(result, isoTransparent)
// 	}
// 	return result
// }

func (doc *TPDFDocument) CreateSectionPageOutLine(S *TPDFSection, PageOutLine, PageIndex, NewPage, ParentOutline, NextOutline, PrevOutLine int) int {
	aDict := doc.GlobalXRefs[ParentOutline].Dict
	aDict.IncCount()
	aDict = doc.GlobalXRefs[PageOutLine].Dict
	arr := aDict.FindElement("Dest").Value.(*TPDFArray)
	arr.AddItem(NewTPDFReference(doc, NewPage))
	arr.AddItem(NewPDFName(doc, "Fit", true))
	result := PrevOutLine
	if PageIndex == 0 {
		doc.GlobalXRefs[ParentOutline].Dict.AddReference("First", doc.xrefCount()-1)
		result = doc.xrefCount() - 1
		aDict = doc.GlobalXRefs[ParentOutline].Dict
		arr = aDict.FindElement("Dest").Value.(*TPDFArray)
		arr.AddItem(NewTPDFReference(doc, NewPage))
		arr.AddItem(NewPDFName(doc, "Fit", true))
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

func (doc *TPDFDocument) CreateSectionOutLine(SectionIndex, OutLineRoot, ParentOutLine, NextSect, PrevSect int) int {
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

func (doc *TPDFDocument) CreateSectionsOutLine() int {
	var result, treeRoot, outlineRoot, pc, j, parentOutline, pageNum, pageOutline, nextOutline, nextSect, newPage, prevOutline, prevSect int
	var aDict *TPDFDictionary
	var arr *TPDFArray

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
				arr = aDict.FindElement("OpenAction").Value.(*TPDFArray)
				arr.AddItem(NewTPDFReference(doc, doc.xrefCount()-1))
				arr.AddItem(NewPDFName(doc, fmt.Sprintf("XYZ null null %f", PDFFloat(doc.ZoomValue/100)), false))
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
	aDict.FindElement("Count").Value.(*TPDFInteger).FInt = pc
	return result
}

// func (doc *TPDFDocument) CreateFontEntries() {
// 	numFont := 0
// 	for _, font := range Fonts {
// 		fontName := font.Name
// 		if IsStandardPDFFont(fontName) {
// 			doc.CreateStdFont(fontName, &numFont)
// 		} else if doc.LoadFont(font) {
// 			if poSubsetFont&Options != 0 {
// 				font.GenerateSubsetFont()
// 			}
// 			doc.CreateTtfFont(&numFont)
// 		} else {
// 			doc.CreateTp1Font(&numFont)
// 		}
// 		numFont++
// 	}
// }

// func (doc *TPDFDocument) CreateImageEntries() {
// 	for i, image := range Images {
// 		doc.CreateImageEntry(image.Width, image.Height, i)
// 		if image.HasMask {
// 			doc.CreateImageMaskEntry(image.Width, image.Height, i)
// 		}
// 	}
// }

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

func (doc *TPDFDocument) WriteXRefTable(st PDFWriter) {
	for _, xr := range doc.GlobalXRefs {
		xr.Encode(st)
	}
}

func (doc *TPDFDocument) CreateEmbeddedFont(APage *TPDFPage, AFontIndex int, AFontSize PDFFloat, ASimulateBold, ASimulateItalic bool) *TPDFEmbeddedFont {
	return NewTPDFEmbeddedFontAdvanced(doc, APage, AFontIndex, AFontSize, ASimulateBold, ASimulateItalic)
}

// func (doc *TPDFDocument) CreateText(X, Y PDFFloat, AText string, AFont *TPDFEmbeddedFont, ADegrees float32, AUnderline, AStrikethrough bool) *TPDFText {
// 	return NewTPDFText(doc, X, Y, AText, AFont, ADegrees, AUnderline, AStrikethrough)
// }

// func (doc *TPDFDocument) CreateUTF8Text(X, Y PDFFloat, AText string, AFont *TPDFEmbeddedFont, ADegrees float32, AUnderline, AStrikethrough bool) *TPDFUTF8Text {
// 	return NewTPDFUTF8Text(doc, X, Y, AText, AFont, ADegrees, AUnderline, AStrikethrough)
// }

// func (doc *TPDFDocument) CreateUTF16Text(X, Y PDFFloat, AText string, AFont *TPDFEmbeddedFont, ADegrees float32, AUnderline, AStrikethrough bool) *TPDFUTF16Text {
// 	return NewTPDFUTF16Text(doc, X, Y, AText, AFont, ADegrees, AUnderline, AStrikethrough)
// }

// func (doc *TPDFDocument) CreateRectangle(X, Y, W, H, ALineWidth PDFFloat, AFill, AStroke bool) *TPDFRectangle {
// 	return NewTPDFRectangle(doc, X, Y, W, H, ALineWidth, AFill, AStroke)
// }

// func (doc *TPDFDocument) CreateRoundedRectangle(X, Y, W, H, ARadius, ALineWidth PDFFloat, AFill, AStroke bool) *TPDFRoundedRectangle {
// 	return NewTPDFRoundedRectangle(doc, X, Y, W, H, ARadius, ALineWidth, AFill, AStroke)
// }

func (doc *TPDFDocument) WriteObject(AObject int, st PDFWriter) {
	st.Writef("%d 0 obj\n", AObject)
	X := doc.GlobalXRefs[AObject]
	if X.Stream == nil {
		X.Dict.WriteDictionary(AObject, st)
	} else {
		// doc.FCurrentColor = ""
		// doc.FCurrentWidth = ""

		// M = NewTPDFStream()
		// X.FStream.Write(M)
		// d = M.Size()

		// if poCompressText&Options != 0 && !X.FStream.CompressionProhibited {
		// 	MCompressed = &TMemoryStream{}
		// 	CompressStream(M, MCompressed)
		// 	X.Dict.AddName("Filter", "FlateDecode")
		// 	d = MCompressed.Size()
		// }
		// X.Dict.AddInteger("Length", d)

		// X.Dict.Write(st)

		// CurrentColor = ""
		// CurrentWidth = ""
		// st.WriteString("\nstream\n", st)
		// if poCompressText&Options != 0 && !X.FStream.CompressionProhibited {
		// 	MCompressed.Position = 0
		// 	MCompressed.SaveToStream(st)
		// 	MCompressed = nil
		// } else {
		// 	M.Position = 0
		// 	M.SaveToStream(st)
		// }

		// M = nil
		// st.WriteString("\n", st)
		// st.WriteString("endstream", st)
	}

	st.WriteString(CRLF + "endobj" + CRLF + CRLF)
}

func (d *TPDFDocument) FindFont(AName string) *TPDFFont {
	for _, f := range d.Fonts {
		if f.FName == AName {
			return f
		}
	}
	return nil
}
