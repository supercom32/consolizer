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
	TransparentForegroundPenalty float64
	AggressiveCoverageThreshold  float64
	AggressiveErrorThreshold     float64
}

func NewImageStyleEntry(existingImageStyleEntry ...*ImageStyleEntryType) ImageStyleEntryType {
	var imageStyleEntry ImageStyleEntryType
	if existingImageStyleEntry != nil {
		imageStyleEntry.DrawingStyle = existingImageStyleEntry[0].DrawingStyle
		imageStyleEntry.IsHistogramEqualized = existingImageStyleEntry[0].IsHistogramEqualized
		imageStyleEntry.IsGrayscale = existingImageStyleEntry[0].IsGrayscale
		imageStyleEntry.DitheringIntensity = existingImageStyleEntry[0].DitheringIntensity
		imageStyleEntry.TransparencyMode = existingImageStyleEntry[0].TransparencyMode
		imageStyleEntry.TransparentForegroundPenalty = existingImageStyleEntry[0].TransparentForegroundPenalty
		imageStyleEntry.AggressiveCoverageThreshold = existingImageStyleEntry[0].AggressiveCoverageThreshold
		imageStyleEntry.AggressiveErrorThreshold = existingImageStyleEntry[0].AggressiveErrorThreshold
	} else {
		// Default to background mode if not specified
		imageStyleEntry.TransparencyMode = constants.TransparencyModeBackground
		imageStyleEntry.TransparentForegroundPenalty = 30.0
		imageStyleEntry.AggressiveCoverageThreshold = 0.35
		imageStyleEntry.AggressiveErrorThreshold = 1.5
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
