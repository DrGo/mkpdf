package main

import (
	"time"
)

type TPDFInfos struct {
	ApplicationName string
	Author          string
	CreationDate    time.Time
	Producer        string
	Title           string
	Keywords        string
}

func NewTPDFInfos() *TPDFInfos {
	return &TPDFInfos{
		Producer: "fpGUI Toolkit 1.4",
		Keywords: "",
	}
}

type TPDFSection struct {
	Title string
	Pages []*TPDFPage
}

func (section *TPDFSection) PageCount() int {
	return len(section.Pages)
}

func (section *TPDFSection) Get(AIndex int) *TPDFPage {
	return section.Pages[AIndex]
}

func (section *TPDFSection) Count() int {
	return len(section.Pages)
}

func (section *TPDFSection) AddPage(APage *TPDFPage) {
	eadd(&section.Pages, APage)
}

type TPDFSectionList struct {
	Sections []*TPDFSection
}

func (list *TPDFSectionList) Get(AIndex int) *TPDFSection {
	return list.Sections[AIndex]
}

func (list *TPDFSectionList) Count() int {
	return len(list.Sections)
}
func (list *TPDFSectionList) NewSection(title string) *TPDFSection {
	list.Sections = append(list.Sections, &TPDFSection{Title: title})
	return list.Sections[list.Count()-1]
}

func NewTPDFSectionList() *TPDFSectionList {
	return &TPDFSectionList{}
}

// type TPDFAnnot struct {
// 	TPDFObject
// 	FLeft, FBottom, FWidth, FHeight PDFFloat
// 	FURI                            string
// 	FBorder, FExternalLink          bool
// }

// func NewTPDFAnnotDetailed(ADocument *TPDFDocument, ALeft, ABottom, AWidth, AHeight PDFFloat, AURI string, ABorder, AExternalLink bool) *TPDFAnnot {
// 	return &TPDFAnnot{
// 		TPDFObject:    NewTPDFObject(ADocument),
// 		FLeft:         ALeft,
// 		FBottom:       ABottom,
// 		FWidth:        AWidth,
// 		FHeight:       AHeight,
// 		FURI:          AURI,
// 		FBorder:       ABorder,
// 		FExternalLink: AExternalLink,
// 	}
// }

// type TPDFAnnotList struct {
// 	TPDFDocumentObject
// 	FList []*TPDFAnnot
// }

// func (l *TPDFAnnotList) Count() int {
// 	return len(l.FList)
// }

// func (l *TPDFAnnotList) Add(AAnnot *TPDFAnnot) {
// 	l.FList = append(l.FList, AAnnot)
// }
