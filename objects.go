package main

import (
	"fmt"
	"strconv"
	"strings"
)


type EPDF struct {
	msg string
}

func (e EPDF) Error() string {
	return e.msg
}

type ARGBColor uint32
type PDFFloat float32

// Encoder is able to write some textual representation into a stream 
type Encoder interface {
	Encode(st PDFWriter)
	Name()string 
}

var _ Encoder = (*TPDFObject)(nil)

type TPDFDocumentObject struct {
	TPDFObject
	FLineCapStyle TPDFLineCapStyle
}

func NewTPDFDocumentObject(document *TPDFDocument) TPDFDocumentObject {
	do:= TPDFDocumentObject{
		TPDFObject: *NewTPDFObject(document),
	}
	if do.Document != nil {
		do.FLineCapStyle = do.Document.LineCapStyle
	}
	return do
}
//FIXME: 
  // S:=FloatStr(AWidth)+' w'; // stroke width
  // if (S<>Document.CurrentWidth) then
  //   begin
  //   WriteString(IntToStr(Ord(FLineCapStyle))+' J'+CRLF, AStream); //set line cap
  //   WriteString(S+CRLF, AStream);
  //   Document.CurrentWidth:=S;
  //   end;

func (docObj *TPDFDocumentObject) SetWidth(AWidth PDFFloat, st PDFWriter){
	S := fmt.Sprintf("%f w", AWidth)
	if S != docObj.Document.CurrentWidth {
		st.WriteString(fmt.Sprintf("%d J\n", docObj.FLineCapStyle))
		st.WriteString(S+"\n")
		docObj.Document.CurrentWidth = S
	}
}


type TPDFObject struct{
	Document *TPDFDocument
}

func (pdfObj *TPDFObject) Encode(st PDFWriter){ }
func (pdfObj *TPDFObject) Name()string{return "" }

func NewTPDFObject(ADocument *TPDFDocument) *TPDFObject {
	obj := &TPDFObject{ADocument}
	if ADocument != nil {
		ADocument.ObjectCount++
	}
	return obj
}

func (obj *TPDFObject) FloatStr(F PDFFloat) string {
	if int(F*100)%100 == 0 {
		return strings.TrimSpace(fmt.Sprintf("%.0f", F))
	}
	return strings.TrimSpace(fmt.Sprintf("%.2f", F))
}

type TPDFBoolean struct {
	TPDFDocumentObject
	FValue bool
}

func (b *TPDFBoolean) Enode(st PDFWriter){
	if b.FValue {
		st.WriteString("true")
	} else {
		st.WriteString("false")
	}
}

func NewTPDFBoolean(ADocument *TPDFDocument, AValue bool) *TPDFBoolean {
	b := &TPDFBoolean{FValue: AValue}
	if ADocument != nil {
		ADocument.ObjectCount++
	}
	return b
}

type TPDFInteger struct {
	TPDFDocumentObject
	FInt int
}

func NewTPDFInteger(ADocument *TPDFDocument, AValue int) *TPDFInteger {
	i := &TPDFInteger{FInt: AValue}
	if ADocument != nil {
		ADocument.ObjectCount++
	}
	return i
}

func (i *TPDFInteger) Encode(st PDFWriter){
	st.WriteString(strconv.Itoa(i.FInt))
}

func (i *TPDFInteger) Inc() {
	i.FInt++
}

type TPDFReference struct {
	TPDFDocumentObject
	FValue int
}

func NewTPDFReference(ADocument *TPDFDocument, AValue int) *TPDFReference {
	ref := &TPDFReference{}
	ref.Document = ADocument
	ref.FValue = AValue
	return ref
}
func (r *TPDFReference) Encode(st PDFWriter) {
	st.WriteString(strconv.Itoa(r.FValue)+" 0 R")
}

type TPDFName struct {
	TPDFDocumentObject
	FName       string
	FMustEscape bool
}

func NewPDFName(document *TPDFDocument, value string, mustEscape bool) *TPDFName {
	return &TPDFName{TPDFDocumentObject:  NewTPDFDocumentObject(document), FName: value, FMustEscape: mustEscape}
}

func (n *TPDFName) Encode(st PDFWriter) {
	if n.FName != "" {
		if strings.Contains(n.FName, "Length1") {
			st.WriteString("/Length1")
		} else {
			st.WriteString("/")
			if n.FMustEscape {
				st.WriteString(ConvertCharsToHex(n.FName ))
			} else {
				st.WriteString(n.FName)
			}
		}
	}
}
type TPDFArray struct {
	TPDFDocumentObject
	FArray []Encoder
}
func (a *TPDFArray) Encode(st PDFWriter) {
	st.WriteByte('[')
	for i, obj := range a.FArray {
		if i > 0 {
			st.WriteByte(' ')
		}
		obj.Encode(st)
	}
	st.WriteByte(']')
}

func (a *TPDFArray) AddItem(AValue Encoder) *TPDFArray{
	a.FArray = append(a.FArray, AValue)
	return a
}

// func (a *TPDFArray) AddIntArray(S string) {
// 	parts := strings.Fields(S)
// 	for _, part := range parts {
// 		val, _ := strconv.Atoi(part)
// 		a.AddItem(NewTPDFInteger(val))
// 	}
// }

// func (a *TPDFArray) AddFreeFormArrayValues(S string) {
// // 	a.AddItem(NewTPDFFreeFormString(nil, S))
// }

func NewTPDFArray(ADocument *TPDFDocument) *TPDFArray {
	return &TPDFArray{
		// FArray: make([]*TPDFObject, 0),
	}
}


//FIXME: 
// var PDFFormatSettings FormatSettings

// func init() {
// 	PDFFormatSettings = DefaultFormatSettings()
// 	PDFFormatSettings.DecimalSeparator = '.'
// 	PDFFormatSettings.ThousandSeparator = ','
// 	PDFFormatSettings.DateSeparator = '/'
// }

type EncoderList []Encoder

func (list *EncoderList) Add(e Encoder) int {
    *list = append(*list, e)
    return len(*list)-1
}

func (list *EncoderList) Find(name string ) (idx int, e Encoder) {
	for i, elm := range *list {
		if elm.Name() == name {
			return i, elm 
		}
	} 
	return -1, nil 
}


func eadd[S ~[]E, E Encoder](s *S, v E) int {
	*s = append(*s, v)
	return len(*s)-1 
}

func efind[S ~[]E, E Encoder] (s *S, name string ) (idx int, e Encoder) {
	for i, elm := range *s {
		if elm.Name() == name {
			return i, elm 
		}
	} 
	return -1, nil 
}

