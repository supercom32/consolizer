package types

import (
	"encoding/json"
	"image"
)

type ImageEntryType struct {
	ImageData  image.Image
	LayerEntry LayerEntryType
}

/*
MarshalJSON is a method which allows you to marshaljson.

Example:

	instance.MarshalJSON()
*/
func (shared ImageEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		ImageData  image.Image
		LayerEntry LayerEntryType
	}{
		ImageData:  shared.ImageData,
		LayerEntry: shared.LayerEntry,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump is a method which allows you to getentryasjsondump.

:return: string

Example:

	instance.GetEntryAsJsonDump()
*/
func (shared ImageEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewImageEntry is a constructor which allows you to newimageentry.

:param existingImageEntry: The existingImageEntry parameter.

:return: ImageEntryType

Example:

	NewImageEntry(existingImageEntry)
*/
func NewImageEntry(existingImageEntry ...*ImageEntryType) ImageEntryType {
	var imageEntry ImageEntryType
	if existingImageEntry != nil {
		imageEntry.ImageData = existingImageEntry[0].ImageData
		imageEntry.LayerEntry = existingImageEntry[0].LayerEntry
	}
	return imageEntry
}
