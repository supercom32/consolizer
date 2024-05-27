package types

import (
	"github.com/supercom32/consolizer/constants"
	"image"
)

type ImageComposerEntryType struct {
	images map[string]*ImageComposerImageEntryType
}

type ImageComposerImageEntryType struct {
	ImageData   image.Image
	XLocation   int
	YLocation   int
	Width       int
	Height      int
	ImageStyle  ImageStyleEntryType
	EffectStyle constants.EffectStyle
	EffectStep  float64
	ZOrder      int
	IsVisible   bool
	AlphaValue  float32
}

type ImageStyleEntryType struct {
	DrawingStyle         constants.ImageStyle
	IsHistogramEqualized bool
	IsGrayscale          bool
	DitheringStyle       constants.DitheringStyle
	DitheringIntensity   float64
}

func NewImageStyleEntry(existingImageStyleEntry ...*ImageStyleEntryType) ImageStyleEntryType {
	var imageStyleEntry ImageStyleEntryType
	if existingImageStyleEntry != nil {
		imageStyleEntry.DrawingStyle = existingImageStyleEntry[0].DrawingStyle
		imageStyleEntry.IsHistogramEqualized = existingImageStyleEntry[0].IsHistogramEqualized
		imageStyleEntry.IsGrayscale = existingImageStyleEntry[0].IsGrayscale
		imageStyleEntry.DitheringIntensity = existingImageStyleEntry[0].DitheringIntensity
	}
	imageStyleEntry.DitheringIntensity = 1 // We set it to 1, which is no change.
	return imageStyleEntry
}

func NewImageComposerEntry() ImageComposerEntryType {
	var imageComposerEntry ImageComposerEntryType
	imageComposerEntry.images = make(map[string]*ImageComposerImageEntryType)
	return imageComposerEntry
}

func NewImageComposerImageEntry() ImageComposerImageEntryType {
	var imageComposerImageEntry ImageComposerImageEntryType
	return imageComposerImageEntry
}

/*func (shared *ImageComposerEntryType) Add(imageAlias string, xLocation int, yLocation int, imageStyle ImageStyleEntryType, zOrder int) *ImageComposerImageEntryType {
	imageComposerImage := NewImageComposerImageEntry()
	imageComposerImage.ZOrder = zOrder
	imageComposerImage.HotspotXLocation = xLocation
	imageComposerImage.HotspotXLocation = yLocation
	imageComposerImage.ImageStyle = imageStyle
	imageComposerImage.IsVisible = true
	imageEntry := memory.GetImage(imageAlias)
	imageComposerImage.ImageData = imageEntry.ImageData
	shared.images[imageAlias] = &imageComposerImage
	return &imageComposerImage
}

func (shared *ImageComposerEntryType) Clear() {
	shared.images = make(map[string]*ImageComposerImageEntryType)
}*/
