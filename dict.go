package main

import (
	"fmt"
)

// In a PDF dictionary, the key is always name obj, the value can be any object
type DictionaryItem struct {
	Key   *PDFName
	Value Encoder
}

func (item DictionaryItem) Name() string { return item.Key.Name }
func (item DictionaryItem) Encode(st PDFWriter) {
	item.Key.Encode(st)
	st.WriteByte(' ')
	item.Value.Encode(st)
	st.WriteString(CRLF)
}

func NewDictionaryItem(doc *Document, key string, aValue Encoder) DictionaryItem {
	//FIXME: check that MustEscape should be false
	return DictionaryItem{NewPDFNameEx( key, false), aValue}
}

// A PDF dictionary holds key/value pairs in any order
type Dictionary struct {
	Document *Document
	Elements []DictionaryItem // list of TPDFDictionaryItem
}

func NewDictionary(doc *Document) *Dictionary { return &Dictionary{Document: doc} }

func (dict *Dictionary) Name() string       { return "" }
func (dict *Dictionary) ElementCount() int  { return len(dict.Elements) }
func (dict *Dictionary) Encode(s PDFWriter) { dict.WriteDictionary(-1, s) }
func (dict *Dictionary) WriteDictionary(AObject int, st PDFWriter) {
	var imgCo, fontCo int
	// var D *TPDFDictionary
	// addSize:= func(v int) {
	// 	D.AddElement("Length", NewTPDFInteger(dict.Document,v))
	// }

	//FIXME:  what is this about
	if dict.GetE(0).Key.Name == "" {
		dict.GetE(0).Encode(st) // write a charwidth array of a font
	} else {
		st.WriteString("<<" + CRLF)
		for _, elem := range dict.Elements {
			elem.Encode(st)
		}
		imgCo = -1
		fontCo = -1
		for _, E := range dict.Elements {
			_ = E
			if AObject > -1 {
				//if E.Key.FName == "Name" {
				//	if obj, ok := E.Value.(*TPDFName); ok && obj.FName[0] == 'M' {
				//		//FIXME: check error
				//		NumImg, _ = strconv.Atoi(obj.FName[1:])
				//		ISize = len(dict.Document.FImages[NumImg].StreamedMask)
				//		D = dict.Document.GlobalXRefs[AObject].Dict
				//		addSize( ISize)
				//		dict.LastElement().Encode(st)
				//		switch dict.Document.FImages[NumImg].FCompressionMask {
				//		case icJPEG:
				//			st.WriteString("/Filter /DCTDecode"+CRLF)
				//		case icDeflate:
				//			st.WriteString("/Filter /FlateDecode"+CRLF)
				//		}
				//		st.WriteString(">>")
				//		dict.Document.FImages[NumImg].WriteMaskStream(st)
				//	} else if obj, ok := E.Value.(*TPDFName); ok && obj.FName[0] == 'I' {
				//		NumImg, _ = strconv.Atoi(obj.FName[1:])
				//		ISize = len(dict.Document.FImages[NumImg].StreamedData)
				//		D = dict.Document.GlobalXRefs[AObject].Dict
				//		addSize( ISize)
				//		dict.LastElement().Encode(st)
				//		switch dict.Document.FImages[NumImg].FCompression {
				//		case icJPEG:
				//			st.WriteString("/Filter /DCTDecode"+CRLF)
				//		case icDeflate:
				//			st.WriteString("/Filter /FlateDecode"+CRLF)
				//		}
				//		st.WriteString(">>")
				//		dict.Document.FImages[NumImg].WriteImageStream(st)
				//	}
				//}
				// if strings.Contains(E.Key.FName, "Length1") {
				// 	Value = E.Key.FName
				// 	pos := strings.Index(Value, " ")
				// 	NumFnt, _ = strconv.Atoi(Value[pos+1:])
				// 	if dict.Document.hasOption(poSubsetFont) {
				// 		var Buf bytes.Buffer
				// BufSize = TPDFEmbeddedFont{}.WriteEmbeddedSubsetFont(dict.Document, NumFnt, Buf)
				// Buf.SetPosition(0)
				// D = dict.Document.GlobalXRefs[AObject].Dict
				// addSize( BufSize)
				// dict.LastElement().Encode(st)
				// st.WriteString(">>")
				// Buf.SaveToStream(st)
				// } else {
				// M = &TMemoryStream{}
				// M.LoadFromFile(dict.Document.FontFiles[NumFnt])
				// Buf = &TMemoryStream{}
				// BufSize = TPDFEmbeddedFont{}.WriteEmbeddedFont(dict.Document, M, Buf)
				// Buf.SetPosition(0)
				// D = dict.Document.GlobalXRefs[AObject].Dict
				// addSize( BufSize)
				// dict.LastElement().Encode(st)
				// st.WriteString(">>")
				// Buf.SaveToStream(st)
				// }
				// }
			}
		}
		if imgCo == -1 && fontCo == -1 {
			st.WriteString(">>")
		}
	}
}

func (dict *Dictionary) AddElement(key string, aValue Encoder) *DictionaryItem {
	dicElement := NewDictionaryItem(dict.Document, key, aValue)
	dict.Elements = append(dict.Elements, dicElement)
	return &dicElement
}

func (dict *Dictionary) AddName(key, aName string) {
	dict.AddElement(key, NewPDFNameEx( aName, false))
}

func (dict *Dictionary) AddNameEscaped(key, aName string) {
	dict.AddElement(key, NewPDFName( aName))
}
func (dict *Dictionary) AddInteger(key string, aInteger int) {
	dict.AddElement(key, NewInteger( aInteger))
}

func (dict *Dictionary) AddReference(key string, aReference int) {
	dict.AddElement(key, NewReference( aReference))
}

func (dict *Dictionary) AddString(key string, aString string) {
	dict.AddElement(key, NewString(dict.Document, aString))
}

// func (dict *TPDFDictionary) AddUTF16String(key string, aString string) {
// 	dict.AddElement(key, dict.FDocument.CreateUTF16String(aString, -1))
// }

// SetElement updates or inserts encoder if it does not already exist
func (dict *Dictionary) SetElement(key string, aValue Encoder) {
	if idx := dict.IndexOfKey(key); idx > -1 {
		dict.Elements[idx].Value = aValue
		return
	}
	dict.AddElement(key, aValue)
}

func (dict *Dictionary) SetSize(size int) {
	dict.SetElement("Size", NewInteger(size))
}

// IncCount increases the count entry in the dictionary
// or create a count entry with value=1 if none exists
func (dict *Dictionary) IncCount() {
	if idx := dict.IndexOfKey("Count"); idx > -1 {
		dict.Elements[idx].Value.(*Integer).Inc()
		return
	}
	dict.AddInteger("Count", 1)
}

// AddKid adds an entry to the kids entry of a dictionary
func (dict *Dictionary) AddKid(kid Encoder) {
	if idx := dict.IndexOfKey("Kids"); idx > -1 {
		dict.Elements[idx].Value.(*Array).AddItem(kid)
		return
	}

	dict.AddElement("Kids", NewArray(dict.Document).AddItem(kid))
}

func (dict *Dictionary) GetE(idx int) DictionaryItem {
	return dict.Elements[idx]
}

func (dict *Dictionary) IndexOfKey(key string) int {
	for i, element := range dict.Elements {
		if element.Key.Name == key {
			return i
		}
	}
	return -1
}

func (dict *Dictionary) LastElement() *DictionaryItem {
	if len(dict.Elements) == 0 {
		return nil
	}
	return &dict.Elements[len(dict.Elements)-1]
}

func (dict *Dictionary) FindElement(key string) *DictionaryItem {
	i := dict.IndexOfKey(key)
	if i == -1 {
		return nil
	}
	return &dict.Elements[i]
}

type XRef struct {
	DocObj
	Offset int
	Dict   *Dictionary
	Stream *TPDFStream
}

func (xref *XRef) Encode(st PDFWriter) {
	st.WriteString(fmt.Sprintf("%010d %05d n"+CRLF, xref.Offset, 0))
}

func NewXRef(doc *Document) *XRef {
	xref := &XRef{
		Offset: 0,
		Dict:   NewDictionary(doc),
	}
	return xref
}
