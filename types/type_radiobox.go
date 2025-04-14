package types

import (
	"encoding/json"
)

type RadioButtonEntryType struct {
	BaseControlType
	IsSelected bool
	GroupId    int
}

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

func (shared RadioButtonEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

// ○ ● U+25CB, U+25CF
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

func IsRadioButtonEqual(sourceRadioButtonEntry *RadioButtonEntryType, targetRadioButtonEntry *RadioButtonEntryType) bool {
	return sourceRadioButtonEntry.BaseControlType.IsEqual(&targetRadioButtonEntry.BaseControlType) &&
		sourceRadioButtonEntry.IsSelected == targetRadioButtonEntry.IsSelected &&
		sourceRadioButtonEntry.GroupId == targetRadioButtonEntry.GroupId
}
