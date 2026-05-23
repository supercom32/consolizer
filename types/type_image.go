package types

import (
	"encoding/json"
	"image"
)

/*
ImageEntryType is a structure which represents an image and its associated layer entry.

Example:
    var imageEntry ImageEntryType
*/
type ImageEntryType struct {
	ImageData  image.Image
	LayerEntry LayerEntryType
}

/*
MarshalJSON is a method which marshals the image entry into a JSON byte array.

Example:
    MarshalJSON()
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
GetEntryAsJsonDump is a method which returns a JSON string representation of the image entry.

Example:
    GetEntryAsJsonDump()
*/
func (shared ImageEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewImageEntry is a constructor which creates a new image entry. In addition, the following should be noted:

- If an existing image entry is provided, the new image entry will be a clone of it.

Example:
    NewImageEntry()
*/
func NewImageEntry(existingImageEntry ...*ImageEntryType) ImageEntryType {
	var imageEntry ImageEntryType
	if existingImageEntry != nil {
		imageEntry.ImageData = existingImageEntry[0].ImageData
		imageEntry.LayerEntry = existingImageEntry[0].LayerEntry
	}
	return imageEntry
}
