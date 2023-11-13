package main

import (
	"time"
)

type PDFInfo struct {
	ApplicationName string
	Author          string
	CreationDate    time.Time
	Producer        string
	Title           string
	Keywords        string
}

func NewPDFInfo() *PDFInfo {
	return &PDFInfo{Producer: ProducerID}
}

type Section struct {
	Title string
	Pages []*Page
}

func (section *Section) PageCount() int     { return len(section.Pages) }
func (section *Section) Get(idx int) *Page  { return section.Pages[idx] }
func (section *Section) Count() int         { return len(section.Pages) }
func (section *Section) AddPage(page *Page) { eadd(&section.Pages, page) }

type SectionList struct {
	Sections []*Section
}

func NewTPDFSectionList() *SectionList         { return &SectionList{} }
func (list *SectionList) Get(idx int) *Section { return list.Sections[idx] }
func (list *SectionList) Count() int           { return len(list.Sections) }
func (list *SectionList) NewSection(title string) *Section {
	list.Sections = append(list.Sections, &Section{Title: title})
	return list.Sections[list.Count()-1]
}

// type TPDFAnnot struct {
// 	TPDFObject
// 	FLeft, FBottom, FWidth, FHeight float64
// 	FURI                            string
// 	FBorder, FExternalLink          bool
// }

// func NewTPDFAnnotDetailed(ADocument *TPDFDocument, ALeft, ABottom, AWidth, AHeight float64, AURI string, ABorder, AExternalLink bool) *TPDFAnnot {
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
