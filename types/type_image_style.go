package types

import (
	"github.com/supercom32/consolizer/constants"
	"time"
)

/*
ImageStyleEntryType is a structure which represents the styling configuration for an image.

Example:
    var imageStyle types.ImageStyleEntryType
*/
type ImageStyleEntryType struct {
	DrawingStyle                 constants.ImageStyle
	IsHistogramEqualized         bool
	IsGrayscale                  bool
	IsWidthAspectRatioPreserved  bool
	IsHeightAspectRatioPreserved bool
	DitheringStyle               constants.DitheringStyle
	DitheringIntensity           float64
	BlurSigmaIntensity           float64
	TransparentForegroundPenalty float64
	AggressiveCoverageThreshold  float64
	AggressiveErrorThreshold     float64
	RandomSeed                   int64
}

/*
NewImageStyleEntry is a constructor which creates a new image style entry.

Example:
    NewImageStyleEntry(existingImageStyleEntry)
*/
func NewImageStyleEntry(existingImageStyleEntry ...*ImageStyleEntryType) ImageStyleEntryType {
	var imageStyleEntry ImageStyleEntryType
	if existingImageStyleEntry != nil {
		imageStyleEntry.DrawingStyle = existingImageStyleEntry[0].DrawingStyle
		imageStyleEntry.IsHistogramEqualized = existingImageStyleEntry[0].IsHistogramEqualized
		imageStyleEntry.IsGrayscale = existingImageStyleEntry[0].IsGrayscale
		imageStyleEntry.DitheringIntensity = existingImageStyleEntry[0].DitheringIntensity
		imageStyleEntry.TransparentForegroundPenalty = existingImageStyleEntry[0].TransparentForegroundPenalty
		imageStyleEntry.AggressiveCoverageThreshold = existingImageStyleEntry[0].AggressiveCoverageThreshold
		imageStyleEntry.AggressiveErrorThreshold = existingImageStyleEntry[0].AggressiveErrorThreshold
		imageStyleEntry.RandomSeed = existingImageStyleEntry[0].RandomSeed
	} else {
		// Default to background mode if not specified
		imageStyleEntry.TransparentForegroundPenalty = 30.0
		imageStyleEntry.AggressiveCoverageThreshold = 0.35
		imageStyleEntry.AggressiveErrorThreshold = 1.5
	}
	if imageStyleEntry.RandomSeed == 0 {
		imageStyleEntry.RandomSeed = time.Now().UnixNano()
	}
	imageStyleEntry.DitheringIntensity = 1 // We set it to 1, which is no change.
	return imageStyleEntry
}
