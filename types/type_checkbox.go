package types

import (
	"encoding/json"
)

type CheckboxEntryType struct {
	BaseControlType
	IsSelected bool
}

func (shared CheckboxEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		IsSelected bool
	}{
		BaseControlType: shared.BaseControlType,
		IsSelected:      shared.IsSelected,
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
func NewCheckboxEntry(existingCheckboxEntry ...*CheckboxEntryType) CheckboxEntryType {
	var checkboxEntry CheckboxEntryType
	checkboxEntry.BaseControlType = NewBaseControl()
	checkboxEntry.IsSelected = false

	if existingCheckboxEntry != nil {
		checkboxEntry.BaseControlType = existingCheckboxEntry[0].BaseControlType
		checkboxEntry.IsSelected = existingCheckboxEntry[0].IsSelected
	}
	return checkboxEntry
}

func IsCheckboxEqual(sourceCheckboxEntry *CheckboxEntryType, targetCheckboxEntry *CheckboxEntryType) bool {
	if sourceCheckboxEntry.BaseControlType == targetCheckboxEntry.BaseControlType &&
		sourceCheckboxEntry.IsSelected == targetCheckboxEntry.IsSelected {
		return true
	}
	return false
}
