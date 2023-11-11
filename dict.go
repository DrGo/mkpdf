package main

import (
	"fmt"
)

// In a PDF dictionary, the key is always name obj, the value can be object
type TPDFDictionaryItem struct {
	Key  *TPDFName
	Value Encoder
}

func (item TPDFDictionaryItem) Encode(st PDFWriter) {
	item.Key.Encode(st)
	st.WriteByte(' ')
	item.Value.Encode(st)
	st.WriteString(CRLF)
}

func NewTPDFDictionaryItem(aDocument *TPDFDocument, aKey string, aValue Encoder) TPDFDictionaryItem {
	//FIXME: check that MustEscape should be false 
	return TPDFDictionaryItem{NewPDFName(aDocument, aKey, false), aValue }
}
// A PDF dictionary holds key/value pairs in any order 
type TPDFDictionary struct {
	Document *TPDFDocument
	Elements []TPDFDictionaryItem // list of TPDFDictionaryItem
}

func NewTPDFDictionary(aDocument *TPDFDocument) *TPDFDictionary {
	return &TPDFDictionary{Document: aDocument}
}

func (dict *TPDFDictionary) Name() string {return ""}
func (dict *TPDFDictionary) ElementCount() int {
	return len(dict.Elements)
}

func (dict *TPDFDictionary) WriteDictionary(AObject int, st PDFWriter) {
	var   NumImg, NumFnt  int
	// var D *TPDFDictionary
	// addSize:= func(v int) {
	// 	D.AddElement("Length", NewTPDFInteger(dict.Document,v))
	// }

//FIXME:  what is this about 	
  // if GetE(0).FKey.Name='' then
  //   GetE(0).Write(AStream)  // write a charwidth array of a font
  // else
	if true {
		st.WriteString("<<"+CRLF)
		for _, elem:= range dict.Elements {
			elem.Encode(st)
		}
		NumImg = -1
		NumFnt = -1
		for _, E := range dict.Elements{
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
		if NumImg == -1 && NumFnt == -1 {
			st.WriteString(">>")
		}
	}
}

func (dict *TPDFDictionary) GetE(aIndex int) TPDFDictionaryItem {
	return dict.Elements[aIndex]
}


func (dict *TPDFDictionary) AddElement(aKey string, aValue Encoder) *TPDFDictionaryItem{
	dicElement := NewTPDFDictionaryItem(dict.Document, aKey, aValue)
	dict.Elements = append(dict.Elements, dicElement)
	return &dicElement
}

func (dict *TPDFDictionary) AddName(aKey, aName string) {
		dict.AddElement(aKey, NewPDFName(dict.Document, aName, false) )
}

func (dict *TPDFDictionary) AddNameEscaped(aKey, aName string) {
		dict.AddElement(aKey, NewPDFName(dict.Document, aName, true) )
}
func (dict *TPDFDictionary) AddInteger(aKey string, aInteger int) {
	dict.AddElement(aKey, NewTPDFInteger(dict.Document, aInteger))
}

func (dict *TPDFDictionary) AddReference(aKey string, aReference int) {
	dict.AddElement(aKey, NewTPDFReference(dict.Document, aReference))
}

func (dict *TPDFDictionary) AddString(aKey string, aString string) {
	dict.AddElement(aKey, NewTPDFString(dict.Document,aString))
}

// func (dict *TPDFDictionary) AddUTF16String(aKey string, aString string) {
// 	dict.AddElement(aKey, dict.FDocument.CreateUTF16String(aString, -1))
// }

// SetElement updates or inserts encode if it does not already exist  
func (dict *TPDFDictionary) SetElement(aKey string, aValue Encoder) {
	if idx:= dict.IndexOfKey(aKey); idx > -1 {
		dict.Elements[idx].Value = aValue
		return
	}
	dict.AddElement(aKey, aValue)	
}

// IncCount increases the count entry in the dictionary  
func (dict *TPDFDictionary) IncCount() {
	if idx:= dict.IndexOfKey("Count"); idx > -1 {
		dict.Elements[idx].Value.(*TPDFInteger).Inc()
		return
	}	
	dict.AddInteger("Count", 1)
}
// AddKid adds an entry to the kids entry of a dictionary
func (dict *TPDFDictionary) AddKid(kid Encoder) {
	if idx:= dict.IndexOfKey("Kids"); idx > -1 {
		dict.Elements[idx].Value.(*TPDFArray).AddItem( kid)
		return
	}	

	dict.AddElement("Kids", NewTPDFArray(dict.Document).AddItem(kid))
}
func (dict *TPDFDictionary) IndexOfKey(aValue string) int {
	for i, element := range dict.Elements {
		if element.Key.FName == aValue {
			return i
		}
	}
	return -1
}

func (dict *TPDFDictionary) Encode(s PDFWriter) {
	dict.WriteDictionary(-1, s)
}


func (dict *TPDFDictionary) LastElement() *TPDFDictionaryItem {
	if len(dict.Elements) == 0 {
		return nil
	}
	return &dict.Elements[len(dict.Elements)-1]
}

func (dict *TPDFDictionary) FindElement(aKey string) *TPDFDictionaryItem {
	i := dict.IndexOfKey(aKey)
	if i == -1 {
		return nil
	}
	return &dict.Elements[i]
}

func (dict *TPDFDictionary) ElementByName(aKey string) *TPDFDictionaryItem {
	result := dict.FindElement(aKey)
	if result == nil {
		panic(fmt.Errorf("dictionary element not found: %s", aKey))
	}
	return result
}


type TPDFXRef struct {
	TPDFDocumentObject
	Offset int
	Dict   *TPDFDictionary
	Stream *TPDFStream
}

func (xref *TPDFXRef) Encode(st PDFWriter) {
	st.WriteString(fmt.Sprintf("%010d %05d n"+CRLF, xref.Offset, 0))
}

func NewTPDFXRef(aDocument *TPDFDocument) *TPDFXRef {
	xref := &TPDFXRef{
		Offset: 0,
		Dict:   NewTPDFDictionary(aDocument),
	}
	return xref
}
