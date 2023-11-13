package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)


type PDFWriter interface {
	io.Writer
	WriteString(string)
	Writef(format string, args ...any)
	WriteByte(byte)
	Err() error
	Offset() int
}

var _ PDFWriter = (*fwriter)(nil)

type fwriter struct {
	w      io.Writer
	err    error
	offset int
}

func Newfwriter(w io.Writer) *fwriter {
	return &fwriter{w : w}
}
//Writef writes a formatted string into a a PDFWriter
// errors and offsets are captured through the call to Write 
func (s *fwriter) Writef(format string, args ...any) {
	fmt.Fprintf(s, format, args...)
}

func (s *fwriter) WriteString(b string) {
	s.Write([]byte(b))
}

func (s *fwriter) WriteByte(b byte) {
	s.Write([]byte{b})
}

func (s *fwriter) Write(b []byte) (int, error) {
	if s.err != nil {
		return 0, s.err 
	}
	var n int
	n, s.err = s.w.Write(b)
	s.offset += n
	return n, s.err 	
}

func (s *fwriter) Err() error {
	return s.err
}

func (s *fwriter) Offset() int {
	return s.offset
}

//memWriter implements PDFWriter but keeps data in 
// a memory buffer of a certain size
type memWriter struct {
	DocObj
	PDFWriter
	buffer *bytes.Buffer 		
}

func NewMemWriter(doc *Document, r io.Reader) *memWriter{
	return NewMemWriterSize(doc, r, -1)
}

func NewMemWriterSize(doc *Document, r io.Reader, capacity int) *memWriter{
// TODO: use sync.pool 
	ms:= &memWriter{DocObj: NewDocObj(doc), 
	buffer: bytes.NewBuffer(make([]byte,0,cond(capacity<1,1024,capacity))),}
	ms.PDFWriter=  Newfwriter(ms.buffer)
	if r != nil {
		io.Copy(ms.buffer, r) 
	}	
	return ms 
}

func (s *memWriter) Size() int {
	return s.buffer.Len()
}

func (s *memWriter) Encode(st PDFWriter) {
	st.Write(s.buffer.Bytes())
}	


func NewfTestWriter() *fwriter { return &fwriter{w : os.Stdout} }
