package main

import (
	"testing"
	"time"
)

func TestStdFontPDF(t *testing.T) {
	pdf := NewDocument() 
	pdf.SetOptions(poPageOriginAtTop, poNoEmbeddedFonts, poSubsetFont, poCompressFonts, poCompressImages)
	pdf.Infos.Producer = "Test"
	pdf.Infos.CreationDate = time.Now()
	pdf.DefaultOrientation = ppoPortrait
	pdf.DefaultPaperType = ptA4
	pdf.DefaultUnitOfMeasure = uomMillimeters

	pdf.StartDocument()
	sec := pdf.Sections.NewSection("section 1")
	sec.AddPage(pdf.NewPage())
	stdFtHelvetica := pdf.AddFont("Helvetica", "")

	pg := sec.Pages[0]
	pg.DrawRect(10,15,10,10,5,true, true, 0)
	pg.SetFont(stdFtHelvetica, 14)
	pg.WriteTextEx(100, 12, "FPC Demo: PDF öäü ÖÄÜ Test", 0, true, true)
	pg.WriteText(10, 10, "AEIOU-ÁÉÍÓÚ-ČŠŇŽ")
	pg.WriteText(10, 15, "----------------")

	// pg.SetFont(Font2, 10)
	pg.WriteText(10, 30, "AEIOU-ÁÉÍÓÚ-ČŠŇŽ")
	pg.WriteText(10, 35, "----------------")

	// pg.SetFont(Font3, 10)
	pg.WriteText(10, 40, "AEIOU-ÁÉÍÓÚ-ČŠŇŽ")
	pg.WriteText(10, 45, "----------------")

	// pg.SetFont(Font4, 10)
	pg.WriteText(10, 50, "AEIOU-ÁÉÍÓÚ-ČŠŇŽ")
	pg.WriteText(10, 55, "----------------")


	pdf.SaveToFile("test-stdfont.pdf")
}



func TestCreatePDF(t *testing.T) {
	pdf := NewDocument()
	pdf.Infos.Producer = ""
	pdf.Infos.CreationDate = time.Now()
	pdf.SetOptions(poPageOriginAtTop, poNoEmbeddedFonts, poSubsetFont, poCompressFonts, poCompressImages)
	pdf.DefaultOrientation = ppoPortrait
	pdf.DefaultPaperType = ptA4
	pdf.DefaultUnitOfMeasure = uomMillimeters
	pdf.FontDirectory = ""
	pdf.StartDocument()
	sec := pdf.Sections.NewSection("section 1")
	sec.AddPage(pdf.NewPage())

	// Font1 := pdf.AddFont("Courier New","cour.ttf")
	// Font2 := pdf.AddFont("Arial",    "arial.ttf")
	// Font3 := pdf.AddFont("Verdana","verdanab.ttf")
	// Font4 := pdf.AddFont("Consolas","consola.ttf")
	
	pg := sec.Pages[0]
	// pg.SetFont(Font1, 10)
	pg.WriteText(10, 10, "AEIOU-ÁÉÍÓÚ-ČŠŇŽ")
	pg.WriteText(10, 15, "----------------")

	// pg.SetFont(Font2, 10)
	pg.WriteText(10, 30, "AEIOU-ÁÉÍÓÚ-ČŠŇŽ")
	pg.WriteText(10, 35, "----------------")

	// pg.SetFont(Font3, 10)
	pg.WriteText(10, 40, "AEIOU-ÁÉÍÓÚ-ČŠŇŽ")
	pg.WriteText(10, 45, "----------------")

	// pg.SetFont(Font4, 10)
	pg.WriteText(10, 50, "AEIOU-ÁÉÍÓÚ-ČŠŇŽ")
	pg.WriteText(10, 55, "----------------")

	pdf.SaveToFile("test.pdf")
}
