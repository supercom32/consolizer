package types

import (
	"encoding/json"
)

/*
LabelEntryType is a structure which represents a label control entry.

Example:

	var labelEntry LabelEntryType
*/
type LabelEntryType struct {
	BaseControlType
	IsBackgroundTransparent bool
}

/*
MarshalJSON is a method which serializes a label control to JSON.

Example:

	MarshalJSON()
*/
func (shared LabelEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		LabelValue string
	}{
		BaseControlType: shared.BaseControlType,
		LabelValue:      shared.Label,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump is a method which returns a JSON string representation of a label control. In addition, the following should be noted:

- Panics if JSON marshaling fails.

Example:

	GetEntryAsJsonDump()
*/
func (shared LabelEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewLabelEntry is a constructor which creates a new label control. In addition, the following should be noted:

- If an existing label entry is provided, the new label entry will be a clone of it.

Example:

	NewLabelEntry()
*/
func NewLabelEntry(existingLabelEntry ...*LabelEntryType) LabelEntryType {
	var labelEntry LabelEntryType
	labelEntry.BaseControlType = NewBaseControl()

	if existingLabelEntry != nil {
		labelEntry.StyleEntry = NewTuiStyleEntry(&existingLabelEntry[0].StyleEntry)
		labelEntry.Alias = existingLabelEntry[0].Alias
		labelEntry.Label = existingLabelEntry[0].Label
		labelEntry.XLocation = existingLabelEntry[0].XLocation
		labelEntry.YLocation = existingLabelEntry[0].YLocation
		labelEntry.Width = existingLabelEntry[0].Width
		labelEntry.IsEnabled = existingLabelEntry[0].IsEnabled
		labelEntry.IsVisible = existingLabelEntry[0].IsVisible
	}
	return labelEntry
}

/*
IsLabelEntryEqual is a method which compares two label controls for equality.

Example:

	IsLabelEntryEqual(&sourceLabelEntry, &targetLabelEntry)
*/
func IsLabelEntryEqual(sourceLabelEntry *LabelEntryType, targetLabelEntry *LabelEntryType) bool {
	return sourceLabelEntry.BaseControlType.IsEqual(&targetLabelEntry.BaseControlType) &&
		sourceLabelEntry.Label == targetLabelEntry.Label &&
		sourceLabelEntry.IsBackgroundTransparent == targetLabelEntry.IsBackgroundTransparent
}
