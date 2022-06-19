package memory

import (
	"encoding/json"
)

type RadioButtonEntryType struct {
	StyleEntry  TuiStyleEntryType
	Label string
	XLocation   int
	YLocation   int
	IsSelected  bool
	IsEnabled bool
	IsVisible bool
	GroupId int
	TabIndex    int
}

func (shared RadioButtonEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		StyleEntry  TuiStyleEntryType
		Label string
		XLocation   int
		YLocation   int
		IsSelected  bool
		IsEnabled bool
		IsVisible bool
		GroupId int
		TabIndex    int
	}{
		StyleEntry: shared.StyleEntry,
		Label: shared.Label,
		XLocation: shared.XLocation,
		YLocation: shared.YLocation,
		IsSelected: shared.IsSelected,
		IsEnabled: shared.IsEnabled,
		IsVisible: shared.IsVisible,
		GroupId: shared.GroupId,
		TabIndex: shared.TabIndex,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared RadioButtonEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}
// ☐ ☑ U+2610, U+2611
func NewRadioButtonEntry(existingRadioButtonEntry ...*RadioButtonEntryType) RadioButtonEntryType {
	var radioButtonEntry RadioButtonEntryType
	if existingRadioButtonEntry != nil {
		radioButtonEntry.StyleEntry = NewTuiStyleEntry(&existingRadioButtonEntry[0].StyleEntry)
		radioButtonEntry.Label = existingRadioButtonEntry[0].Label
		radioButtonEntry.XLocation = existingRadioButtonEntry[0].XLocation
		radioButtonEntry.YLocation = existingRadioButtonEntry[0].YLocation
		radioButtonEntry.IsSelected = existingRadioButtonEntry[0].IsSelected
		radioButtonEntry.IsEnabled = existingRadioButtonEntry[0].IsEnabled
		radioButtonEntry.IsVisible = existingRadioButtonEntry[0].IsVisible
		radioButtonEntry.TabIndex = existingRadioButtonEntry[0].TabIndex
	}
	return radioButtonEntry
}

func IsRadioButtonEqual(sourceCheckboxEntry *RadioButtonEntryType, targetCheckboxEntry *RadioButtonEntryType) bool {
	if sourceCheckboxEntry.StyleEntry == targetCheckboxEntry.StyleEntry &&
		sourceCheckboxEntry.Label == targetCheckboxEntry.Label &&
		sourceCheckboxEntry.XLocation == targetCheckboxEntry.XLocation &&
		sourceCheckboxEntry.YLocation == targetCheckboxEntry.YLocation &&
		sourceCheckboxEntry.IsSelected == targetCheckboxEntry.IsSelected &&
		sourceCheckboxEntry.IsEnabled == targetCheckboxEntry.IsEnabled &&
		sourceCheckboxEntry.IsVisible == targetCheckboxEntry.IsVisible &&
		sourceCheckboxEntry.TabIndex == targetCheckboxEntry.TabIndex {
		return true
	}
	return false
}
