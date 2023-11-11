package main

// type TFPCustomImage struct {
// 	Height, Width int
// 	Colors []TFPColor
// }
// type TFPColor struct {
//     red,green,blue,alpha int16
// }

// type TPDFImageStreamOptions []TPDFImageStreamOption
// type TPDFImageItem struct {
// 	*TPDFDocumentObject
// 	FImage           *TFPCustomImage
// 	FOwnsImage       bool
// 	FStreamed        []byte
// 	FCompression     TPDFImageCompression
// 	FStreamedMask    []byte
// 	FCompressionMask TPDFImageCompression
// 	FWidth, FHeight  int

// }

// func (imgItem *TPDFImageItem) SetImage(aValue *TFPCustomImage) {
// 	if imgItem.FImage == aValue {
// 		return
// 	}
// 	imgItem.FImage = aValue
// 	imgItem.FStreamed = []byte{}
// }

// func (imgItem *TPDFImageItem) GetStreamed() []byte {
// 	var opts TPDFImageStreamOptions

// 	if len(imgItem.FStreamed) == 0 {
// 		if imgItem.FDocument != nil  {
// 			opts = imgItem.FDocument.ImageStreamOptions
// 		} else {
// 			opts = TPDFImageStreamOptions{isoCompressed, isoTransparent}
// 		}
// 		imgItem.CreateStreamedData(opts)
// 	}
// 	return imgItem.FStreamed
// }

// func (imgItem *TPDFImageItem) GetStreamedMask() []byte {
// 	imgItem.GetStreamed()
// 	return imgItem.FStreamedMask
// }

// func (imgItem *TPDFImageItem) GetHeight() int {
// 	if imgItem.FImage != nil {
// 		return imgItem.FImage.Height
// 	}
// 	return imgItem.FHeight
// }

// func (imgItem *TPDFImageItem) GetWidth() int {
// 	if imgItem.FImage != nil {
// 		return imgItem.FImage.Width
// 	}
// 	return imgItem.FWidth
// }

// func (imgItem *TPDFImageItem) SetStreamed(aValue []byte) {
// 	if bytes.Equal(aValue, imgItem.FStreamed) {
// 		return
// 	}
// 	imgItem.FStreamed = []byte{}
// 	imgItem.FStreamed = aValue
// }

// func (imgItem *TPDFImageItem) SetStreamedMask(aValue []byte, aCompression TPDFImageCompression) {
// 	if bytes.Equal(aValue, imgItem.FStreamedMask) {
// 		return
// 	}
// 	imgItem.FStreamedMask = []byte{}
// 	imgItem.FStreamedMask = aValue
// 	imgItem.FCompressionMask = aCompression
// }

// func (imgItem *TPDFImageItem) WriteImageStream(aStream TStream) int64 {
// 	return imgItem.WriteStream(imgItem.FStreamed, aStream)
// }

// func (imgItem *TPDFImageItem) WriteMaskStream(aStream TStream) int64 {
// 	return imgItem.WriteStream(imgItem.FStreamedMask, aStream)
// }

// func (imgItem *TPDFImageItem) CreateStreamedData(aOptions TPDFImageStreamOptions) {
// 	needsTransparency := func() bool {
// 		for y := 0; y < imgItem.FHeight; y++ {
// 			for x := 0; x < imgItem.FWidth; x++ {
// 				if imgItem.FImage.Colors[x][y].alpha < 0xFFFF {
// 					return true
// 				}
// 			}
// 		}
// 		return false
// 	}

// 	// var createMask bool
// 	// var cWhite TFPColor
// 	// var x, y int
// 	// var c TFPColor

// 	// cWhite = TFPColor{0xFF, 0xFF, 0xFF, 0}
// 	// imgItem.FWidth = imgItem.FImage.Width
// 	// imgItem.FHeight = imgItem.FImage.Height
// 	// createMask = contains(aOptions, isoTransparent) && needsTransparency()
// }

// func (imgItem *TPDFImageItem) DetachImage() {
// 	imgItem.FImage = nil
// }

// func (imgItem *TPDFImageItem) WriteStream(aStreamedData []byte, aStream TStream) int64 {
// }

// func (imgItem *TPDFImageItem) Equals(aImage TFPCustomImage) bool {
// 	if aImage == nil {
// 		return false
// 	}

// 	if imgItem.FImage.Width != aImage.Width || imgItem.FImage.Height != aImage.Height {
// 		return false
// 	}

// 	for y := 0; y < imgItem.FImage.Height; y++ {
// 		for x := 0; x < imgItem.FImage.Width; x++ {
// 			if imgItem.FImage.Colors[x][y] != aImage.Colors[x][y] {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }

// func (imgItem *TPDFImageItem) GetHasMask() bool {
// 	return len(imgItem.FStreamedMask) > 0
// }

// func (images *TPDFImages) GetI(index int) *TPDFImageItem {
// 	return images.Items[index].(*TPDFImageItem)
// }

// func (images *TPDFImages) AddImageItem() *TPDFImageItem {
// 	return images.Add().(*TPDFImageItem)
// }

// func (images *TPDFImages) AddJPEGStream(AStream TStream, Width, Height int) int {
// 	IP := images.AddImageItem()
// 	IP.FWidth = Width
// 	IP.FHeight = Height
// 	IP.FCompression = icJPEG
// 	streamData := make([]byte, AStream.Size()-AStream.Position)
// 	if len(streamData) > 0 {
// 		AStream.Read(streamData)
// 		IP.FStreamed = streamData
// 	}
// 	return images.Count() - 1
// }

// func (images *TPDFImages) AddFromFile(AFileName string, KeepImage bool) int {
// 	FS := NewTFileStream(AFileName, fmOpenRead|fmShareDenyNone)
// 	defer FS.Free()
// 	return images.AddFromStream(FS, FindReaderFromFileName(AFileName), KeepImage)
// }

// func (images *TPDFImages) AddFromStream(AStream *TStream, Handler TFPCustomImageReaderClass, KeepImage bool) int {
// 	var result int
// 	if poUseRawJPEGIsSet(Owner.Options) && HandlerIsTFPReaderJPEG(Handler) {
// 		JPEG := NewTFPReaderJPEG()
// 		defer JPEG.Free()

// 		if FPC_FULLVERSION >= 30101 {
// 			size := JPEG.ImageSize(AStream)
// 			result = images.AddJPEGStream(AStream, size.X, size.Y)
// 		} else {
// 			I := NewTFPMemoryImage(0, 0)
// 			defer I.Free()
// 			startPos := AStream.Position
// 			I.LoadFromStream(AStream, JPEG)
// 			AStream.Position = startPos
// 			result = images.AddJPEGStream(AStream, I.Width, I.Height)
// 		}
// 	} else {
// 		IP := images.AddImageItem()
// 		I := NewTFPMemoryImage(0, 0)
// 		if Handler == nil {
// 			panic("Error: No Image Reader")
// 		}
// 		Reader := Handler.Create()
// 		defer Reader.Free()
// 		I.LoadFromStream(AStream, Reader)
// 		IP.Image = I
// 		if KeepImage {
// 			IP.OwnsImage = true
// 		} else {
// 			IP.CreateStreamedData(Owner.ImageStreamOptions)
// 			IP.DetachImage()
// 			I.Free()
// 		}
// 	}
// 	return images.Count() - 1
// }

// type TPDFImage struct {
// 	TPDFDocumentObject
// 	Number int
// 	Pos    TPDFCoord
// 	Size   TPDFCoord
// }

// func (img *TPDFImage) Write(st TStream) {
// 	st.WriteString(PushGraphicsStackCommand())
// 	st.WriteString(fmt.Sprintf("%f 0 0 %f %f %f cm%s", img.Size.x, img.Size.y, img.Pos.x, img.Pos.y, CRLF))
// 	st.WriteString(fmt.Sprintf("/I%d Do%s", img.Number, CRLF))
// 	st.WriteString(PopGraphicsStackCommand())
// }

// func NewPDFImage(document *TPDFDocument, left, bottom, width, height PDFFloat, number int) *TPDFImage {
// 	return &TPDFImage{
// 		TPDFDocumentObject: NewTPDFDocumentObject(document),
// 		Number:           number,
// 		Pos:              TPDFCoord{x: left, y: bottom},
// 		Size:             TPDFCoord{x: width, y: height},
// 	}
// }



