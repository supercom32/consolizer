package types

import (
	"encoding/json"
)

/*
CheckboxEntryType is a structure which represents a checkbox control. In addition, the following should be noted:

- It includes base control properties and the selection state of the checkbox.

Example:
    var checkbox types.CheckboxEntryType
*/
type CheckboxEntryType struct {
	BaseControlType
	IsSelected bool
}

/*
MarshalJSON is a method which serializes a checkbox control to JSON. In addition, the following should be noted:

- It converts the checkbox's state to a JSON representation.

Example:
    instance.MarshalJSON()
*/
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

/*
GetEntryAsJsonDump is a method which retrieves a JSON string representation of a checkbox control. In addition, the following should be noted:

- It returns a formatted JSON string of the checkbox's state.

Example:
    instance.GetEntryAsJsonDump()
*/
func (shared CheckboxEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewCheckboxEntry is a constructor which creates a new checkbox control. In addition, the following should be noted:

- It initializes a checkbox with default values.

- It can optionally copy properties from an existing checkbox.

- It uses Unicode characters U+2610 and U+2611 for rendering states.

Example:
    NewCheckboxEntry(existingCheckboxEntry)
*/
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

/*
IsCheckboxEqual is a method which compares two checkbox controls for equality. In addition, the following should be noted:

- It compares both the base control properties and the selection state.

Example:
    IsCheckboxEqual(sourceCheckboxEntry, targetCheckboxEntry)
*/
func IsCheckboxEqual(sourceCheckboxEntry *CheckboxEntryType, targetCheckboxEntry *CheckboxEntryType) bool {
	return sourceCheckboxEntry.BaseControlType.IsEqual(&targetCheckboxEntry.BaseControlType) &&
		sourceCheckboxEntry.IsSelected == targetCheckboxEntry.IsSelected
}
