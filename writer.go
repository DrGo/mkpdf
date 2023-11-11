package main

import (
	"fmt"
	"io"
)


type PDFWriter interface {
	WriteString(string)
	Write([]byte)
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

func (s *fwriter) Writef(format string, args ...any) {
	fmt.Fprintf(s.w, format, args...)
}

func (s *fwriter) WriteString(b string) {
	s.Write([]byte(b))
}

func (s *fwriter) WriteByte(b byte) {
	s.Write([]byte{b})
}

func (s *fwriter) Write(b []byte) {
	if s.err != nil {
		return
	}
	var n int
	n, s.err = s.w.Write(b)
	s.offset += n
}

func (s *fwriter) Err() error {
	return s.err
}

func (s *fwriter) Offset() int {
	return s.offset
}

