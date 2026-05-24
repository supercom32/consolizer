package types

import (
	"github.com/supercom32/consolizer/constants"
	"image"
)

/*
ImageComposerEntryType is a structure which represents a collection of images to be composed together.

Example:

	var imageComposer types.ImageComposerEntryType
*/
type ImageComposerEntryType struct {
	images map[string]*ImageComposerImageEntryType
}

/*
ImageComposerImageEntryType is a structure which represents an individual image within an image composer.

Example:

	var imageComposerImage types.ImageComposerImageEntryType
*/
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

/*
NewImageComposerEntry is a constructor which creates a new image composer entry.

Example:

	NewImageComposerEntry()
*/
func NewImageComposerEntry() ImageComposerEntryType {
	var imageComposerEntry ImageComposerEntryType
	imageComposerEntry.images = make(map[string]*ImageComposerImageEntryType)
	return imageComposerEntry
}

/*
NewImageComposerImageEntry is a constructor which creates a new image composer image entry.

Example:

	NewImageComposerImageEntry()
*/
func NewImageComposerImageEntry() ImageComposerImageEntryType {
	var imageComposerImageEntry ImageComposerImageEntryType
	return imageComposerImageEntry
}
