package types

import (
	"encoding/json"
	"reflect"
)

type TextFieldEntryType struct {
	BaseControlType
	MaxLengthAllowed    int
	DefaultValue        string
	CursorPosition      int
	ViewportPosition    int
	IsPasswordProtected bool
	CurrentValue        []rune
}

func (shared TextFieldEntryType) GetAlias() string {
	return shared.Alias
}

func (shared TextFieldEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		MaxLengthAllowed    int
		DefaultValue        string
		CursorPosition      int
		ViewportPosition    int
		IsPasswordProtected bool
		CurrentValue        []rune
	}{
		BaseControlType:     shared.BaseControlType,
		MaxLengthAllowed:    shared.MaxLengthAllowed,
		DefaultValue:        shared.DefaultValue,
		CursorPosition:      shared.CursorPosition,
		ViewportPosition:    shared.ViewportPosition,
		IsPasswordProtected: shared.IsPasswordProtected,
		CurrentValue:        shared.CurrentValue,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared TextFieldEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewTextFieldEntry(existingTextFieldEntry ...*TextFieldEntryType) TextFieldEntryType {
	var textFieldEntry TextFieldEntryType
	textFieldEntry.BaseControlType = NewBaseControl()

	if existingTextFieldEntry != nil {
		textFieldEntry.BaseControlType = existingTextFieldEntry[0].BaseControlType
		textFieldEntry.MaxLengthAllowed = existingTextFieldEntry[0].MaxLengthAllowed
		textFieldEntry.DefaultValue = existingTextFieldEntry[0].DefaultValue
		textFieldEntry.CursorPosition = existingTextFieldEntry[0].CursorPosition
		textFieldEntry.ViewportPosition = existingTextFieldEntry[0].ViewportPosition
		textFieldEntry.IsPasswordProtected = existingTextFieldEntry[0].IsPasswordProtected
		textFieldEntry.CurrentValue = existingTextFieldEntry[0].CurrentValue
	}
	textFieldEntry.CurrentValue = []rune{' '}
	return textFieldEntry
}

func IsTextFieldEntryEqual(sourceTextFieldEntry *TextFieldEntryType, targetTextFieldEntry *TextFieldEntryType) bool {
	if sourceTextFieldEntry.BaseControlType == targetTextFieldEntry.BaseControlType &&
		sourceTextFieldEntry.MaxLengthAllowed == targetTextFieldEntry.MaxLengthAllowed &&
		sourceTextFieldEntry.DefaultValue == targetTextFieldEntry.DefaultValue &&
		sourceTextFieldEntry.CursorPosition == targetTextFieldEntry.CursorPosition &&
		sourceTextFieldEntry.ViewportPosition == targetTextFieldEntry.ViewportPosition &&
		sourceTextFieldEntry.IsPasswordProtected == targetTextFieldEntry.IsPasswordProtected &&
		reflect.DeepEqual(sourceTextFieldEntry.CurrentValue, targetTextFieldEntry.CurrentValue) {
		return true
	}
	return false
}

func GetTextFieldAlias(entry *TextFieldEntryType) string {
	return entry.Alias
}
