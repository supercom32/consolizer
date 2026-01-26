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

func NewImageComposerEntry() ImageComposerEntryType {
	var imageComposerEntry ImageComposerEntryType
	imageComposerEntry.images = make(map[string]*ImageComposerImageEntryType)
	return imageComposerEntry
}

func NewImageComposerImageEntry() ImageComposerImageEntryType {
	var imageComposerImageEntry ImageComposerImageEntryType
	return imageComposerImageEntry
}
