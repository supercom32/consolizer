package types

import (
	"encoding/json"
)

// func DrawButton(LayerAlias string, Label string, StyleEntry TuiStyleEntryType, IsPressed bool, IsSelected bool, XLocation int, YLocation int, Width int, Height int) {
type ProgressBarEntryType struct {
	BaseControlType
	Label                   string
	Value                   int
	MaxValue                int
	IsBackgroundTransparent bool
	Length                  int
	IsVertical              bool
}

/*
GetAlias is a method which allows you to retrieve the alias of a progress bar control. In addition, the following
information should be noted:

- Returns the unique identifier for the progress bar.

- This alias is used to reference the progress bar in other operations.

- The alias is set when the progress bar is created.

:return: string

Example:

	instance.GetAlias()
*/
func (shared ProgressBarEntryType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON is a method which allows you to serialize a progress bar control to JSON. In addition, the following
information should be noted:

- Converts the progress bar's state to a JSON representation.

- Includes the base control properties and progress bar-specific fields.

- Used for saving and loading progress bar configurations.

Example:

	instance.MarshalJSON()
*/
func (shared ProgressBarEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		Label                   string
		Value                   int
		MaxValue                int
		IsBackgroundTransparent bool
		Length                  int
		CurrentValue            int
		IsVertical              bool
	}{
		BaseControlType:         shared.BaseControlType,
		Label:                   shared.Label,
		Value:                   shared.Value,
		MaxValue:                shared.MaxValue,
		IsBackgroundTransparent: shared.IsBackgroundTransparent,
		Length:                  shared.Length,
		IsVertical:              shared.IsVertical,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump is a method which allows you to get a JSON string representation of a progress bar control. In
addition, the following information should be noted:

- Returns a formatted JSON string of the progress bar's state.

- Useful for debugging and logging purposes.

- Panics if JSON marshaling fails.

:return: string

Example:

	instance.GetEntryAsJsonDump()
*/
func (shared ProgressBarEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewProgressBarEntry is a constructor which allows you to create a new progress bar control. In addition, the following
information should be noted:

- Initializes a progress bar with default values.

- Can optionally copy properties from an existing progress bar.

- Sets up the base control properties and progress bar-specific fields.

:param ExistingProgressBarEntry: The ExistingProgressBarEntry parameter.

:return: ProgressBarEntryType

Example:

	NewProgressBarEntry(ExistingProgressBarEntry)
*/
func NewProgressBarEntry(ExistingProgressBarEntry ...*ProgressBarEntryType) ProgressBarEntryType {
	var progressBarEntry ProgressBarEntryType
	if ExistingProgressBarEntry != nil {
		progressBarEntry.BaseControlType = ExistingProgressBarEntry[0].BaseControlType
		progressBarEntry.Label = ExistingProgressBarEntry[0].Label
		progressBarEntry.Value = ExistingProgressBarEntry[0].Value
		progressBarEntry.MaxValue = ExistingProgressBarEntry[0].MaxValue
		progressBarEntry.IsBackgroundTransparent = ExistingProgressBarEntry[0].IsBackgroundTransparent
		progressBarEntry.Length = ExistingProgressBarEntry[0].Length
		progressBarEntry.IsVertical = ExistingProgressBarEntry[0].IsVertical
	}
	return progressBarEntry
}

/*
IsProgressBarEntryEqual is a method which allows you to compare two progress bar controls for equality. In addition, the
following information should be noted:

- Compares all properties of both progress bars.

- Returns true if all properties match, false otherwise.

- Used for change detection and state synchronization.

:param sourceProgressBarEntry: The sourceProgressBarEntry parameter.
:param targetProgressBarEntry: The targetProgressBarEntry parameter.

:return: bool

Example:

	IsProgressBarEntryEqual(sourceProgressBarEntry, targetProgressBarEntry)
*/
func IsProgressBarEntryEqual(sourceProgressBarEntry *ProgressBarEntryType, targetProgressBarEntry *ProgressBarEntryType) bool {
	return sourceProgressBarEntry.BaseControlType.IsEqual(&targetProgressBarEntry.BaseControlType) &&
		sourceProgressBarEntry.Label == targetProgressBarEntry.Label &&
		sourceProgressBarEntry.Value == targetProgressBarEntry.Value &&
		sourceProgressBarEntry.MaxValue == targetProgressBarEntry.MaxValue &&
		sourceProgressBarEntry.IsBackgroundTransparent == targetProgressBarEntry.IsBackgroundTransparent &&
		sourceProgressBarEntry.Length == targetProgressBarEntry.Length &&
		sourceProgressBarEntry.IsVertical == targetProgressBarEntry.IsVertical
}

/*
GetProgressBarAlias is a method which allows you to retrieve the alias of a progress bar control. In addition, the
following information should be noted:

- Returns the unique identifier for the progress bar.

- This is a convenience method that delegates to GetAlias.

- The alias is used to reference the progress bar in other operations.

:param entry: The entry parameter.

:return: string

Example:

	GetProgressBarAlias(entry)
*/
func GetProgressBarAlias(entry *ProgressBarEntryType) string {
	return entry.Alias
}
