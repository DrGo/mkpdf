package main

import (
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

// Encoder is able to write some textual representation into a stream
type Encoder interface {
	Encode(st PDFWriter)
}

type Boolean struct {
	val bool
}

func NewBoolean(val bool) *Boolean { return &Boolean{val: val} }
func (b *Boolean) Encode(st PDFWriter) {
	if b.val {
		st.WriteString("true")
	}
	st.WriteString("false")
}

// Integer signed decimal integer; exponential notation not allowed
type Integer struct {
	val int
}

func NewInteger(val int) *Integer { return &Integer{val: val} }
func (i *Integer) Encode(st PDFWriter) {
	st.WriteString(strconv.Itoa(i.val))
}
func (i *Integer) Inc() { i.val++ }

type Reference struct {
	val int
}

func NewReference(val int) *Reference {
	return &Reference{val}
}
func (r *Reference) Encode(st PDFWriter) {
	st.WriteString(strconv.Itoa(r.val) + " 0 R")
}

type PDFName struct {
	Name       string
	MustEscape bool
}

func NewPDFName(value string) *PDFName {
	return NewPDFNameEx(value, true)
}

func NewPDFNameEx(value string, mustEscape bool) *PDFName {
	return &PDFName{Name: value, MustEscape: mustEscape}
}
func (n *PDFName) Encode(st PDFWriter) {
	if n.Name != "" {
		if strings.Contains(n.Name, "Length1") {
			st.WriteString("/Length1")
		} else {
			st.WriteString("/")
			if n.MustEscape {
				st.WriteString(ConvertCharsToHex(n.Name))
			} else {
				st.WriteString(n.Name)
			}
		}
	}
}

type Array struct {
	Elements []Encoder
}

func NewArray(doc *Document) *Array { return &Array{} }
func (a *Array) Encode(st PDFWriter) {
	st.WriteByte('[')
	for i, obj := range a.Elements {
		if i > 0 {
			st.WriteByte(' ')
		}
		obj.Encode(st)
	}
	st.WriteByte(']')
}

func (a *Array) AddItem(val Encoder) *Array {
	a.Elements = append(a.Elements, val)
	return a
}

func (a *Array) AddInts(vals ...int) *Array {
	for _, val := range vals {
		a.AddItem(NewInteger(val))
	}
	return a
}

// func (a *TPDFArray) AddFreeFormArrayValues(S string) {
// // 	a.AddItem(NewTPDFFreeFormString(nil, S))
// }

// type EncoderList []Encoder

// func (list *EncoderList) Add(e Encoder) int {
// 	*list = append(*list, e)
// 	return len(*list) - 1
// }

// func (list *EncoderList) Find(name string) (idx int, e Encoder) {
// 	for i, elm := range *list {
// 		if elm.Name() == name {
// 			return i, elm
// 		}
// 	}
// 	return -1, nil
// }

func eadd[S ~[]E, E Encoder](s *S, v E) int {
	*s = append(*s, v)
	return len(*s) - 1
}

// func efind[S ~[]E, E Encoder](s *S, name string) (idx int, e Encoder) {
// 	for i, elm := range *s {
// 		if elm.Name() == name {
// 			return i, elm
// 		}
// 	}
// 	return -1, nil
// }
