package main

// Some popular predefined colors. Channel format is: RRGGBB
const (
	clBlack     = 0x000000
	clWhite     = 0xFFFFFF
	clBlue      = 0x0000FF
	clGreen     = 0x008000
	clRed       = 0xFF0000
	clAqua      = 0x00FFFF
	clMagenta   = 0xFF00FF
	clYellow    = 0xFFFF00
	clLtGray    = 0xC0C0C0
	clMaroon    = 0x800000
	clOlive     = 0x808000
	clDkGray    = 0x808080
	clTeal      = 0x008080
	clNavy      = 0x000080
	clPurple    = 0x800080
	clLime      = 0x00FF00
	clWaterMark = 0xF0F0F0
)

type PaperType int

const (
	ptCustom PaperType = iota
	ptA4
	ptA5
	ptLetter
	ptLegal
	ptExecutive
	ptComm10
	ptMonarch
	ptDL
	ptC5
	ptB5
)

type PaperOrientation int

const (
	ppoPortrait PaperOrientation = iota
	ppoLandscape
)

type PenStyle int

const (
	ppsSolid PenStyle = iota
	ppsDash
	ppsDot
	ppsDashDot
	ppsDashDotDot
)

type LineCapStyle int

const (
	plcsButtCap LineCapStyle = iota
	plcsRoundCap
	plcsProjectingSquareCap
)

type LineJoinStyle int

const (
	pljsMiterJoin LineJoinStyle = iota
	pljsRoundJoin
	pljsBevelJoin
)

type PageLayout int

const (
	lSingle PageLayout = iota
	lTwo
	lContinuous
)

type UnitOfMeasure int

const (
	uomInches UnitOfMeasure = iota
	uomMillimeters
	uomCentimeters
	uomPixels
)

type Option int

const (
	poOutLine Option = iota
	poCompressText
	poCompressFonts
	poCompressImages
	poUseRawJPEG
	poNoEmbeddedFonts
	poPageOriginAtTop
	poSubsetFont
	poMetadataEntry
	poNoTrailerID
	poUseImageTransparency
	poUTF16info
)

type ImageCompression int

const (
	icNone ImageCompression = iota
	icDeflate
	icJPEG
)

type ImageStreamOption int

const (
	isoCompressed ImageStreamOption = iota
	isoTransparent
)

// uses flate pkg constants
type CompressionLevel int

const (
	clDefault CompressionLevel = iota
)

const (
	cInchToMM                  = 25.4
	cInchToCM                  = 2.54
	cDefaultDPI                = 72
	BEZIER                     = 0.5522847498 // 4/3 * (sqrt(2) - 1)
	CRLF                       = "\r\n"       //  EOL marker regardless of OS
	PDF_VERSION                = "%PDF-1.3"
	PDF_BINARY_BLOB            = "%" + "\xC3\xA4" + "\xC3\xBC" + "\xC3\xB6" + "\xC3\x9F"
	PDF_FILE_END               = "%%EOF"
	PDF_MAX_GEN_NUM            = 65535
	PDF_UNICODE_HEADER         = "FEFF001B%s001B"
	PDF_LANG_STRING            = "en"
	PDF_NUMBER_MASK            = "0.####"
	rsErrReportFontFileMissing = `Font File "%s" does not exist.`
	rsErrDictElementNotFound   = `Error: Dictionary element "%s" not found.`
	rsErrInvalidSectionPage    = `Error: Invalid section page index.`
	rsErrNoGlobalDict          = `Error: no global XRef named "%s".`
	rsErrInvalidPageIndex      = `Invalid page index: %d`
	rsErrInvalidAnnotIndex     = `Invalid annot index: %d`
	rsErrNoFontDefined         = `No Font was set - please use SetFont() first.`
	rsErrNoImageReader         = `Unsupported image format - no image reader available.`
	rsErrUnknownStdFont        = `Unknown standard PDF font name <%s>.`

	ProducerID = "mkPDF"
)

type Dimensions struct {
	B, T, L, R float64
}

type Paper struct {
	H, W      int
	Printable Dimensions
}

func (a Dimensions) Equal(b Dimensions) bool {
	return a.B == b.B && a.T == b.T && a.L == b.L && a.R == b.R
}

func (a Paper) Equal(b Paper) bool {
	return a.H == b.H && a.W == b.W && a.Printable.Equal(b.Printable)
}

// Height,Width,Top,Left,Right,Bottom (units in pixels)
var PDFPaperDims = map[PaperType][6]int{
	ptCustom:    {0, 0, 0, 0, 0, 0},
	ptA4:        {842, 595, 10, 11, 586, 822},
	ptA5:        {595, 420, 10, 11, 407, 588},
	ptLetter:    {792, 612, 13, 13, 599, 780},
	ptLegal:     {1008, 61, 13, 13, 599, 996},
	ptExecutive: {756, 522, 14, 13, 508, 744},
	ptComm10:    {684, 297, 13, 13, 284, 672},
	ptMonarch:   {540, 279, 13, 13, 266, 528},
	ptDL:        {624, 312, 14, 13, 297, 611},
	ptC5:        {649, 459, 13, 13, 446, 637},
	ptB5:        {709, 499, 14, 13, 485, 696},
}
var PageLayoutNames = []string{"SinglePage", "TwoColumnLeft", "OneColumn"}

// not using pi to avoid importing math
const (
	_DegToRad = 0.017453292519943295769236907684886127134428718885417 // N[Pi/180, 50]
	_RadToDeg = 57.295779513082320876798154814105170332405472466564   // N[180/Pi, 50]

	// GradToRad = 0.015707963267948966192313216916397514420985846996876 // N[Pi/200, 50]
	// RadToGrad = 63.661977236758134307553505349005744813783858296183   // N[200/Pi, 50]
)
