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
	IsHorizontal            bool
}

/*
GetAlias allows you to retrieve the alias of a progress bar control. In addition, the following
information should be noted:

- Returns the unique identifier for the progress bar.
- This alias is used to reference the progress bar in other operations.
- The alias is set when the progress bar is created.
*/
func (shared ProgressBarEntryType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON allows you to serialize a progress bar control to JSON. In addition, the following
information should be noted:

- Converts the progress bar's state to a JSON representation.
- Includes the base control properties and progress bar-specific fields.
- Used for saving and loading progress bar configurations.
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
		IsHorizontal            bool
	}{
		BaseControlType:         shared.BaseControlType,
		Label:                   shared.Label,
		Value:                   shared.Value,
		MaxValue:                shared.MaxValue,
		IsBackgroundTransparent: shared.IsBackgroundTransparent,
		Length:                  shared.Length,
		IsHorizontal:            shared.IsHorizontal,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump allows you to get a JSON string representation of a progress bar control. In addition,
the following information should be noted:

- Returns a formatted JSON string of the progress bar's state.
- Useful for debugging and logging purposes.
- Panics if JSON marshaling fails.
*/
func (shared ProgressBarEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewProgressBarEntry allows you to create a new progress bar control. In addition, the following
information should be noted:

- Initializes a progress bar with default values.
- Can optionally copy properties from an existing progress bar.
- Sets up the base control properties and progress bar-specific fields.
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
		progressBarEntry.IsHorizontal = ExistingProgressBarEntry[0].IsHorizontal
	}
	return progressBarEntry
}

/*
IsProgressBarEntryEqual allows you to compare two progress bar controls for equality. In addition, the following
information should be noted:

- Compares all properties of both progress bars.
- Returns true if all properties match, false otherwise.
- Used for change detection and state synchronization.
*/
func IsProgressBarEntryEqual(sourceProgressBarEntry *ProgressBarEntryType, targetProgressBarEntry *ProgressBarEntryType) bool {
	if sourceProgressBarEntry.BaseControlType == targetProgressBarEntry.BaseControlType &&
		sourceProgressBarEntry.Label == targetProgressBarEntry.Label &&
		sourceProgressBarEntry.Value == targetProgressBarEntry.Value &&
		sourceProgressBarEntry.MaxValue == targetProgressBarEntry.MaxValue &&
		sourceProgressBarEntry.IsBackgroundTransparent == targetProgressBarEntry.IsBackgroundTransparent &&
		sourceProgressBarEntry.Length == targetProgressBarEntry.Length &&
		sourceProgressBarEntry.IsHorizontal == targetProgressBarEntry.IsHorizontal {
		return true
	}
	return false
}

/*
GetProgressBarAlias allows you to retrieve the alias of a progress bar control. In addition, the following
information should be noted:

- Returns the unique identifier for the progress bar.
- This is a convenience method that delegates to GetAlias.
- The alias is used to reference the progress bar in other operations.
*/
func GetProgressBarAlias(entry *ProgressBarEntryType) string {
	return entry.Alias
}
