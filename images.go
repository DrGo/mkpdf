package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"io"
	"os"
)

type TFPColor struct {
	red, green, blue, alpha int16
}
type ImageCompression byte

const (
	icNone ImageCompression = iota
	icDeflate
	icJPEG
)

type ImageStreamOption byte

const (
	isoCompressed ImageStreamOption = iota
	isoTransparent
)

// uses flate pkg constants
type CompressionLevel int

const (
	clDefault CompressionLevel = iota
)

type RasterImage struct {
	*DocObj
	Image            image.Image
	Width, Height    int
	ColorSpace       color.Model
	Format           string
	Filter           string
	BitsPerComponent byte
	OwnsImage        bool
	Stream           []byte
	Compression      ImageCompression
	MaskStream       []byte
	MaskCompression  ImageCompression
}

func (img *RasterImage) PDFColorSpace() string {
	switch img.ColorSpace {
	case color.GrayModel, color.Gray16Model:
		return "DeviceGray"
	case color.YCbCrModel, color.NRGBAModel, color.NRGBA64Model, color.RGBAModel, color.RGBA64Model:
		return "DeviceRGB"
	case color.CMYKModel:
		return "DeviceCMYK"
	default:
		return ""
	}
}

func (img *RasterImage) GetStreamed() []byte {
	var opts []ImageStreamOption
	if len(img.Stream) == 0 {
		if img.Document != nil {
			opts = img.Document.ImageStreamOptions
		} else {
			opts = []ImageStreamOption{isoCompressed, isoTransparent}
		}
		img.CreateStreamedData(opts)
	}
	return img.Stream
}

func (img *RasterImage) GetStreamedMask() []byte {
	img.GetStreamed()
	return img.MaskStream
}

func (img *RasterImage) SetStreamed(aValue []byte) {
	if bytes.Equal(aValue, img.Stream) {
		return
	}
	img.Stream = []byte{}
	img.Stream = aValue
}

func (img *RasterImage) SetStreamedMask(aValue []byte, aCompression ImageCompression) {
	if bytes.Equal(aValue, img.MaskStream) {
		return
	}
	img.MaskStream = []byte{}
	img.MaskStream = aValue
	img.MaskCompression = aCompression
}

func (img *RasterImage) WriteImageStream(aStream PDFWriter) int64 {
	return img.WriteStream(img.Stream, aStream)
}

func (img *RasterImage) WriteMaskStream(aStream PDFWriter) int64 {
	return img.WriteStream(img.MaskStream, aStream)
}

func (img *RasterImage) CreateStreamedData(aOptions []ImageStreamOption) {
	// needsTransparency := func() bool {
	// 	for y := 0; y < img.Height; y++ {
	// 		for x := 0; x < img.Width; x++ {
	// 			if img.Image.Colors[x][y].alpha < 0xFFFF {
	// 				return true
	// 			}
	// 		}
	// 	}
	// 	return false
	// }

	// var createMask bool
	// var cWhite TFPColor
	// var x, y int
	// var c TFPColor

	// cWhite = TFPColor{0xFF, 0xFF, 0xFF, 0}
	// img.FWidth = img.FImage.Width
	// img.FHeight = img.FImage.Height
	// createMask = contains(aOptions, isoTransparent) && needsTransparency()
}

func (img *RasterImage) DetachImage() {
	img.Image = nil
}

// FIXME:
func (img *RasterImage) WriteStream(aStreamedData []byte, aStream PDFWriter) int64 {
	return 0
}

// func (img *ImageItem) Equals(aImage *PDFImageConfig) bool {
// 	if aImage == nil {
// 		return false
// 	}

// 	if img.Image.Width != aImage.Width || img.Image.Height != aImage.Height {
// 		return false
// 	}

// 	for y := 0; y < img.Image.Height; y++ {
// 		for x := 0; x < img.Image.Width; x++ {
// 			if img.Image.Colors[x][y] != aImage.Colors[x][y] {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }

func (img *RasterImage) HasMask() bool {
	return len(img.MaskStream) > 0
}

func (doc *Document) GetRasterImage(index int) *RasterImage {
	return doc.Images[index]
}

func (doc *Document) AddImageItem(img *RasterImage) int {
	doc.Images = append(doc.Images, img)
	return doc.ImageCount() - 1
}

func (doc *Document) AddImageFromFile(AFileName string) (int, error) {
	f, err := os.Open(AFileName)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	return doc.AddImageFromReader(f)
}

func NewRasterImage(config image.Config, format string) *RasterImage {
	return &RasterImage{
		Width:  config.Width,
		Height: config.Height,
		Format: format,
	}
}

func (doc *Document) AddImageFromReader(r io.Reader) (int, error) {
	config, format, err := image.DecodeConfig(r)
	if err != nil {
		return -1, err
	}
	img := NewRasterImage(config, format)
	if format == "jpeg" {
		img.Compression = icJPEG
		img.ColorSpace = config.ColorModel
		img.BitsPerComponent = 8
		img.Filter = "DCTDecode"
	}
	pixels, format, err := image.Decode(r)
	if err != nil {
		return -1, err
	}
	img.Image = pixels
	return doc.AddImageItem(img), nil
}

func (doc *Document) SetImageStreamOptions() []ImageStreamOption {
	var result []ImageStreamOption
	if doc.hasOption(poCompressImages) {
		result = append(result, isoCompressed)
	}
	if doc.hasOption(poUseImageTransparency) {
		result = append(result, isoTransparent)
	}
	return result
}

func (doc *Document) CreateImageEntries() {
	for i, image := range doc.Images {
		doc.CreateImageEntry(image, i)
		// if image.HasMask() {
		// 	doc.CreateImageMaskEntry(image.Width, image.Height, i)
		// }
	}
}

func ImgToRawDate(img image.Image) []byte {
	b := img.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
	return m.Pix
}

func (doc *Document) CreateImageEntry(img *RasterImage, num int) *Dictionary {
	rx := doc.xrefCount() // reference to be used later
	imgDict := doc.CreateGlobalXRef("XObject").Dict
	imgDict.AddName("Subtype", "Image")
	imgDict.AddInteger("Width", img.Width)
	imgDict.AddInteger("Height", img.Height)
	imgDict.AddName("ColorSpace", img.PDFColorSpace())
	imgDict.AddInteger("BitsPerComponent", int(img.BitsPerComponent))
	N := NewPDFName(fmt.Sprintf("I%d", num)) // Needed later
	imgDict.AddElement("Name", N)

	// now find where we must add the image xref - we are looking for "Resources"
	
	for _, xr := range doc.GlobalXRefs {
		if dict := xr.Dict; dict.ElementCount() > 0 {
			debug(dict.Elements[0].Value, "value of first element of dict: ")
			if v, ok := dict.Elements[0].Value.(*PDFName); ok && v.Name == "Page" {
				res := dict.FindElement("Resources").Value.(*Dictionary)
				res = res.FindElement("XObject").Value.(*Dictionary)
				res.AddReference(N.Name, rx)
			}
		}
	}
	return imgDict
}

func (doc *Document) CreateImageMaskEntry(ImgWidth, ImgHeight, NumImg int, ImageDict *Dictionary) {
	lXRef := doc.xrefCount() // reference to be used later
	MDict := doc.CreateGlobalXRef("XObject").Dict
	MDict.AddName("Subtype", "Image")
	MDict.AddInteger("Width", ImgWidth)
	MDict.AddInteger("Height", ImgHeight)
	MDict.AddName("ColorSpace", "DeviceGray")
	MDict.AddInteger("BitsPerComponent", 8)
	N := NewPDFName(fmt.Sprintf("M%d", NumImg)) // Needed later
	MDict.AddElement("Name", N)
	ImageDict.AddReference("SMask", lXRef)
}
