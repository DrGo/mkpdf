package main

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func GetLocalTZD(ADate time.Time, ISO8601 bool) string {
	offset := ADate.Sub(ADate.UTC()).Minutes()
	sign := "+"
	if offset < 0 {
		sign = "-"
		offset = -offset
	} else if offset == 0 {
		return "Z"
	}

	var format string
	if ISO8601 {
		format = "%02d:%02d"
	} else {
		format = "%02d'%02d'"
	}
	return sign + fmt.Sprintf(format, int(offset/60), int(math.Mod(offset, 60)))
}

func DateToPdfDate(ADate time.Time) string {
	return ADate.Format("D:20060102150405") + GetLocalTZD(ADate, false)
}

func FormatPDFInt(value, padLen int) string {
	result := strconv.Itoa(value)
	padLen -= len(result)
	if padLen > 0 {
		result = fmt.Sprintf("%0*s", padLen, "") + result
	}
	return result
}

func CompressStream(AFrom, ATo *bytes.Buffer, ACompressLevel TCompressionLevel, ASkipHeader bool) {
	if AFrom.Len() == 0 {
		ATo.Reset()
		return
	}

	writer := zlib.NewWriter(ATo)
	_, _ = AFrom.WriteTo(writer)
	_ = writer.Close()
}

func CompressString(AFrom string) string {
	lStreamFrom := bytes.NewBufferString(AFrom)
	lStreamTo := &bytes.Buffer{}
	CompressStream(lStreamFrom, lStreamTo, clDefault, false)
	return lStreamTo.String()
}

func DecompressStream(AFrom, ATo *bytes.Buffer) {
	if AFrom.Len() == 0 {
		ATo.Reset()
		return
	}

	reader, _ := zlib.NewReader(AFrom)
	// FIXME	
	// _, _ = reader.WriteTo(ATo)
	_ = reader.Close()
}

func mmToPDF(mm PDFFloat) PDFFloat {
	return PDFFloat(mm * (cDefaultDPI / cInchToMM))
}

func PDFTomm(APixels PDFFloat) PDFFloat {
	return PDFFloat(APixels*cInchToMM) / cDefaultDPI
}

func cmToPDF(cm PDFFloat) PDFFloat {
	return PDFFloat(cm * (cDefaultDPI / cInchToCM))
}

func PDFtoCM(APixels PDFFloat) PDFFloat {
	return PDFFloat(APixels*cInchToCM) / cDefaultDPI
}

func InchesToPDF(Inches PDFFloat) PDFFloat {
	return PDFFloat(Inches * cDefaultDPI)
}
func PDFtoInches(APixels PDFFloat) PDFFloat {
	return PDFFloat(APixels) / cDefaultDPI
}

// func FontUnitsTomm(AUnits, APointSize PDFFloat, AUnitsPerEm int) PDFFloat {
// 	return PDFFloat(AUnits * APointSize * dpi / (72 * PDFFloat(AUnitsPerEm)) * cInchToMM / dpi)
// }

// func XMLEscape(data string) string {
// 	result := make([]rune, len(data)*6)
// 	iPos := 0
// 	for _, r := range data {
// 		switch r {
// 		case '<':
// 			iPos += copy(result[iPos:], "&lt;")
// 		case '>':
// 			iPos += copy(result[iPos:], "&gt;")
// 		case '&':
// 			iPos += copy(result[iPos:], "&amp;")
// 		case '"':
// 			iPos += copy(result[iPos:], "&quot;")
// 		default:
// 			result[iPos] = r
// 			iPos++
// 		}
// 	}
// 	return string(result[:iPos])
// }

func ExtractBaseFontName(AValue string) string {
	FontName := strings.TrimRight(AValue, "-")
	S1 := strings.Split(AValue, ":")[1]
	S1 = strings.Title(S1)
	S2 := ""
	if strings.Contains(S1, ":") {
		parts := strings.Split(S1, ":")
		S1 = parts[0]
		S2 = strings.Title(parts[1])
	}
	return FontName + "-" + S1 + S2
}

func FloatStr(f PDFFloat) string {
	if f == PDFFloat(int64(f)) { //is a whole number
		return strconv.FormatFloat(float64(f), 'f', 0, 32)
	}
	return strconv.FormatFloat(float64(f), 'f', 2, 32)
}

// FIXME: this mess
func ConvertCharsToHex(schar string) string {
	var s string
	for _, char := range schar {
		d := int(char)
		if d < 33 || d > 126 {
			s += "#" + strconv.FormatInt(int64(d), 16)
		} else {
			s += string(char)
		}
	}
	return s
}

func octStr(b byte) string {
	var result string
	for b != 0 {
		result = strconv.Itoa(int(b&7)) + result
		b = b >> 3
	}
	return result
}

func DoUnitConversion(APoint *TPDFCoord, UnitOfMeasure TPDFUnitOfMeasure) {
	switch UnitOfMeasure {
	case uomMillimeters:
		APoint.x = mmToPDF(APoint.x)
		APoint.y = mmToPDF(APoint.y)
	case uomCentimeters:
		APoint.x = cmToPDF(APoint.x)
		APoint.y = cmToPDF(APoint.y)
	case uomInches:
		APoint.x = InchesToPDF(APoint.x)
		APoint.y = InchesToPDF(APoint.y)
	}
}

// func DegToRad(a PDFFloat) PDFFloat {
// 	return PDFFloat(a * _DegToRad)
// }


func DegToRad(a float32) float32 {
	return float32(a * _DegToRad)
}

func sincos(a float32) (sin,cost float64) {
	return math.Sincos(float64(a))
}

func GetMD5Hash(text string) string {
   hash := md5.Sum([]byte(text))
   return hex.EncodeToString(hash[:])
}
