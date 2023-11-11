package main

import (
	"testing"
	"time"
)

func TestCreatePDF(t *testing.T) {
	pdf := NewTPDFDocument()
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

	// Font1 := pdf.AddFont("cour.ttf", "Courier New")
	// Font2 := pdf.AddFont("arial.ttf", "Arial")
	// Font3 := pdf.AddFont("verdanab.ttf", "Verdana")
	// Font4 := pdf.AddFont("consola.ttf", "Consolas")
	
	pg := sec.Pages[0]
	pg.SetFont(Font1, 10)
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
