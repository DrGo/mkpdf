package main

import (
	"compress/flate"
	"io"
)

type TPDFStream struct {
	DocObj
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

func NewPDFStream(doc *Document) *TPDFStream{
	return &TPDFStream{DocObj: NewDocObj(doc)}
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
	DocObj
}

func NewTXMPStream(doc *Document) *TXMPStream{
	return &TXMPStream{DocObj: NewDocObj(doc)}
}

func compressStream(dest io.Writer, src io.Reader /*, compressionLevel int*/) error {
	zw, err := flate.NewWriter(dest, flate.DefaultCompression)
	if err!= nil {
		return err
	}
	_, err= io.Copy(zw, src)
	return err 
}

// procedure CompressString(const AFrom: rawbytestring; var ATo: rawbytestring);
// var
//   lStreamFrom : TStringStream;
//   lStreamTo  : TStringStream;
// begin
//   { TODO : Possible improvement would be to perform this compression directly on
//            the string as a buffer, and not go through the stream stage. }
//   lStreamFrom := TStringStream.Create(AFrom);
//   try
//     lStreamTo  := TStringStream.Create('');
//     try
//       lStreamFrom.Position := 0;
//       lStreamTo.Size := 0;
//       CompressStream(lStreamFrom, lStreamTo);
//       ATo  := lStreamTo.DataString;
//     finally
//       lStreamTo.Free;
//     end;
//   finally
//     lStreamFrom.Free;
//   end;
// end;

// procedure DecompressStream(AFrom: TStream; ATo: TStream);
// Const
//   BufSize = 1024; // 1K
// Type
//   TBuffer = Array[0..BufSize-1] of byte;
// var
//   d: TDecompressionStream;
//   Count : Integer;
//   Buffer : TBuffer;

// begin
//   if AFrom.Size = 0 then
//   begin
//     ATo.Size := 0;
//     Exit; //==>
//   end;
//   FillMem(@Buffer, SizeOf(TBuffer), 0);

//   AFrom.Position := 0;
//   AFrom.Seek(0,soFromEnd);
//   D:=TDecompressionStream.Create(AFrom, False);
//   try
//     repeat
//        Count:=D.Read(Buffer,BufSize);
//        ATo.WriteBuffer(Buffer,Count);
//      until (Count<BufSize);
//   finally
//     d.Free;
//   end;
// end;


