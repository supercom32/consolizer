package types

import (
	"encoding/json"
)

type CheckboxEntryType struct {
	StyleEntry TuiStyleEntryType
	Label      string
	XLocation  int
	YLocation  int
	IsSelected bool
	IsEnabled  bool
	IsVisible  bool
	TabIndex   int
}

func (shared CheckboxEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		StyleEntry TuiStyleEntryType
		Label      string
		XLocation  int
		YLocation  int
		IsSelected bool
		IsEnabled  bool
		IsVisible  bool
		TabIndex   int
	}{
		StyleEntry: shared.StyleEntry,
		Label:      shared.Label,
		XLocation:  shared.XLocation,
		YLocation:  shared.YLocation,
		IsSelected: shared.IsSelected,
		IsEnabled:  shared.IsEnabled,
		IsVisible:  shared.IsVisible,
		TabIndex:   shared.TabIndex,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared CheckboxEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

// ☐ ☑ U+2610, U+2611
func NewCheckboxEntry(existingButtonEntry ...*CheckboxEntryType) CheckboxEntryType {
	var checkboxEntry CheckboxEntryType
	if existingButtonEntry != nil {
		checkboxEntry.StyleEntry = NewTuiStyleEntry(&existingButtonEntry[0].StyleEntry)
		checkboxEntry.Label = existingButtonEntry[0].Label
		checkboxEntry.XLocation = existingButtonEntry[0].XLocation
		checkboxEntry.YLocation = existingButtonEntry[0].YLocation
		checkboxEntry.IsSelected = existingButtonEntry[0].IsSelected
		checkboxEntry.IsEnabled = existingButtonEntry[0].IsEnabled
		checkboxEntry.IsVisible = existingButtonEntry[0].IsVisible
		checkboxEntry.TabIndex = existingButtonEntry[0].TabIndex
	}
	return checkboxEntry
}

func IsCheckboxEqual(sourceCheckboxEntry *CheckboxEntryType, targetCheckboxEntry *CheckboxEntryType) bool {
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
