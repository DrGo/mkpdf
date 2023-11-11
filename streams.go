package main

import (
	"bytes"
	"io"
)

type TPDFStream struct {
	TPDFDocumentObject
	CompressionProhibited bool
	Items                 []Encoder
}

func (s *TPDFStream) AddItem(AValue Encoder) {
	s.Items = append(s.Items, AValue)
}

func (s *TPDFStream) Encode(st PDFWriter) {
	for _, item := range s.Items {
		item.Encode(st)
	}
}

func NewPDFStream(doc *TPDFDocument) *TPDFStream{
	return &TPDFStream{TPDFDocumentObject: NewTPDFDocumentObject(doc)}
}

type TPDFMemoryStream struct {
	TPDFDocumentObject
	buffer bytes.Buffer 		
}

func NewPDFMemoryStream(doc *TPDFDocument, r io.Reader) *TPDFMemoryStream{
	ms:= &TPDFMemoryStream{TPDFDocumentObject: NewTPDFDocumentObject(doc)}
	if r != nil {
		io.Copy(&ms.buffer, r) 
	}	
	return ms 
}

func (s *TPDFMemoryStream) Encode(st PDFWriter) {
	st.Write(s.buffer.Bytes())
}	


//FIXME: check implementation 
func (xs *TXMPStream) Encode(st PDFWriter) {
	// add := func(tag, value string) {
	// 	fmt.Fprintf(st, "<%s>%s</%s>\n", tag, value, tag)
	// }
	// dateToISO8601Date := func(t time.Time) string {
	// 	return t.Format(time.RFC3339)
	// }
	// nbsp := '\uFEFF'
	// fmt.Fprintf(st, `<?xpacket begin="%c" id="W5M0MpCehiHzreSzNTczkc9d"?>`+"\n", nbsp)

	// for i := 0; i < 21; i++ {
	// 	st.Write([]byte(strings.Repeat(" ", 99) + "\n"))
	// }
	// st.Write([]byte(`<?xpacket end="w"?>`))
}

type TXMPStream struct {
	TPDFDocumentObject
}

func NewTXMPStream(doc *TPDFDocument) *TXMPStream{
	return &TXMPStream{TPDFDocumentObject: NewTPDFDocumentObject(doc)}
}
