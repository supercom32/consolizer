package types

import (
	"encoding/json"
)

// func DrawButton(LayerAlias string, Text string, StyleEntry TuiStyleEntryType, IsPressed bool, IsSelected bool, XLocation int, YLocation int, Width int, Height int) {
type LabelEntryType struct {
	BaseControlType
	IsBackgroundTransparent bool
}

/*
MarshalJSON is a method which allows you to serialize a label control to JSON. In addition, the following information
should be noted:

- Converts the label's state to a JSON representation.

- Includes the base control properties and label-specific fields.

- Used for saving and loading label configurations.

Example:

	instance.MarshalJSON()
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
GetEntryAsJsonDump is a method which allows you to get a JSON string representation of a label control. In addition, the
following information should be noted:

- Returns a formatted JSON string of the label's state.

- Useful for debugging and logging purposes.

- Panics if JSON marshaling fails.

:return: string

Example:

	instance.GetEntryAsJsonDump()
*/
func (shared LabelEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewLabelEntry is a constructor which allows you to create a new label control. In addition, the following information
should be noted:

- Initializes a label with default values.

- Can optionally copy properties from an existing label.

- Sets up the base control properties and label-specific fields.

:param existingLabelEntry: The existingLabelEntry parameter.

:return: LabelEntryType

Example:

	NewLabelEntry(existingLabelEntry)
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
IsLabelEntryEqual is a method which allows you to compare two label controls for equality. In addition, the following
information should be noted:

- Compares all properties of both labels.

- Returns true if all properties match, false otherwise.

- Used for change detection and state synchronization.

:param sourceLabelEntry: The sourceLabelEntry parameter.
:param targetLabelEntry: The targetLabelEntry parameter.

:return: bool

Example:

	IsLabelEntryEqual(sourceLabelEntry, targetLabelEntry)
*/
func IsLabelEntryEqual(sourceLabelEntry *LabelEntryType, targetLabelEntry *LabelEntryType) bool {
	return sourceLabelEntry.BaseControlType.IsEqual(&targetLabelEntry.BaseControlType) &&
		sourceLabelEntry.Label == targetLabelEntry.Label &&
		sourceLabelEntry.IsBackgroundTransparent == targetLabelEntry.IsBackgroundTransparent
}
