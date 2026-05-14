package types

import (
	"encoding/json"
)

type RadioButtonEntryType struct {
	BaseControlType
	IsSelected bool
	GroupId    int
}

/*
MarshalJSON is a method which allows you to marshaljson.

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
GetEntryAsJsonDump is a method which allows you to getentryasjsondump.

:return: string

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
NewRadioButtonEntry is a constructor which allows you to ○ ● U+25CB, U+25CF.

:param existingRadioButtonEntry: The existingRadioButtonEntry parameter.

:return: RadioButtonEntryType

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
IsRadioButtonEqual is a method which allows you to isradiobuttonequal.

:param sourceRadioButtonEntry: The sourceRadioButtonEntry parameter.
:param targetRadioButtonEntry: The targetRadioButtonEntry parameter.

:return: bool

Example:

	IsRadioButtonEqual(sourceRadioButtonEntry, targetRadioButtonEntry)
*/
func IsRadioButtonEqual(sourceRadioButtonEntry *RadioButtonEntryType, targetRadioButtonEntry *RadioButtonEntryType) bool {
	return sourceRadioButtonEntry.BaseControlType.IsEqual(&targetRadioButtonEntry.BaseControlType) &&
		sourceRadioButtonEntry.IsSelected == targetRadioButtonEntry.IsSelected &&
		sourceRadioButtonEntry.GroupId == targetRadioButtonEntry.GroupId
}
