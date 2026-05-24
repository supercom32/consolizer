package types

import (
	"encoding/json"
)

/*
RadioButtonEntryType is a structure which represents a radio button control entry.

Example:

	var radioButton types.RadioButtonEntryType
*/
type RadioButtonEntryType struct {
	BaseControlType
	IsSelected bool
	GroupId    int
}

/*
MarshalJSON is a method which serializes the radio button entry to JSON and returns the resulting byte array.

Example:

	instance.MarshalJSON()
*/
func (shared RadioButtonEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		IsSelected bool
		GroupId    int
	}{
		BaseControlType: shared.BaseControlType,
		IsSelected:      shared.IsSelected,
		GroupId:         shared.GroupId,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump is a method which returns a JSON string representation of the radio button entry.

Example:

	instance.GetEntryAsJsonDump()
*/
func (shared RadioButtonEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewRadioButtonEntry is a constructor which creates and returns a new radio button entry instance.

Example:

	NewRadioButtonEntry(existingRadioButtonEntry)
*/
func NewRadioButtonEntry(existingRadioButtonEntry ...*RadioButtonEntryType) RadioButtonEntryType {
	var radioButtonEntry RadioButtonEntryType
	radioButtonEntry.BaseControlType = NewBaseControl()
	radioButtonEntry.IsSelected = false
	radioButtonEntry.GroupId = 0

	if existingRadioButtonEntry != nil {
		radioButtonEntry.BaseControlType = existingRadioButtonEntry[0].BaseControlType
		radioButtonEntry.IsSelected = existingRadioButtonEntry[0].IsSelected
		radioButtonEntry.GroupId = existingRadioButtonEntry[0].GroupId
	}
	return radioButtonEntry
}

/*
IsRadioButtonEqual is a method which compares two radio button entries for equality and returns true if they are equal.

Example:

	IsRadioButtonEqual(sourceRadioButtonEntry, targetRadioButtonEntry)
*/
func IsRadioButtonEqual(sourceRadioButtonEntry *RadioButtonEntryType, targetRadioButtonEntry *RadioButtonEntryType) bool {
	return sourceRadioButtonEntry.BaseControlType.IsEqual(&targetRadioButtonEntry.BaseControlType) &&
		sourceRadioButtonEntry.IsSelected == targetRadioButtonEntry.IsSelected &&
		sourceRadioButtonEntry.GroupId == targetRadioButtonEntry.GroupId
}
