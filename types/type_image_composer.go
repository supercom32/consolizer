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
	DrawingStyle                 constants.ImageStyle
	IsHistogramEqualized         bool
	IsGrayscale                  bool
	IsWidthAspectRatioPreserved  bool
	IsHeightAspectRatioPreserved bool
	DitheringStyle               constants.DitheringStyle
	DitheringIntensity           float64
	BlurSigmaIntensity           float64
	TransparencyMode             constants.TransparencyMode
}

func NewImageStyleEntry(existingImageStyleEntry ...*ImageStyleEntryType) ImageStyleEntryType {
	var imageStyleEntry ImageStyleEntryType
	if existingImageStyleEntry != nil {
		imageStyleEntry.DrawingStyle = existingImageStyleEntry[0].DrawingStyle
		imageStyleEntry.IsHistogramEqualized = existingImageStyleEntry[0].IsHistogramEqualized
		imageStyleEntry.IsGrayscale = existingImageStyleEntry[0].IsGrayscale
		imageStyleEntry.DitheringIntensity = existingImageStyleEntry[0].DitheringIntensity
		imageStyleEntry.TransparencyMode = existingImageStyleEntry[0].TransparencyMode
	} else {
		// Default to background mode if not specified
		imageStyleEntry.TransparencyMode = constants.TransparencyModeBackground
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
