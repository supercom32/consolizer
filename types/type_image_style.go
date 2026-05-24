package types

import (
	"github.com/supercom32/consolizer/constants"
	"time"
)

/*
ImageStyleEntryType is a structure which represents the styling configuration for an image. In addition, the following
should be noted:

- `DitheringIntensity`: This is a multiplier for the dithering effect. A value of 1.0 means no change.

- `BlurSigmaIntensity`: This is the sigma value used for blurring after resizing. A value of 0.0 means no blurring.

  - `TransparentForegroundPenalty`: This is a penalty value used when selecting the best block element. A higher value
    makes it less likely that a transparent foreground will be chosen if it would result in a loss of detail.

  - `AggressiveCoverageThreshold`: This is the minimum percentage of opaque pixels (0.0 to 1.0) required within a cell to
    consider it for rendering. Cells with coverage below this threshold may be culled unless they are a very good fit.

  - `AggressiveErrorThreshold`: This is the maximum allowed error for a low-coverage cell to survive culling. Lower
    values make culling more aggressive.

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
