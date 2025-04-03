package types

import (
	"encoding/json"
	"reflect"
)

type TextboxEntryType struct {
	BaseControlType
	HorizontalScrollbarAlias string
	VerticalScrollbarAlias   string
	TextData                 [][]rune
	ViewportXLocation        int
	ViewportYLocation        int
	CursorXLocation          int
	CursorYLocation          int
}

func (shared TextboxEntryType) GetAlias() string {
	return shared.Alias
}

func (shared TextboxEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		HorizontalScrollbarAlias string
		VerticalScrollbarAlias   string
		TextData                 [][]rune
		ViewportX                int
		ViewportY                int
		CursorX                  int
		CursorY                  int
	}{
		BaseControlType:          shared.BaseControlType,
		HorizontalScrollbarAlias: shared.HorizontalScrollbarAlias,
		VerticalScrollbarAlias:   shared.VerticalScrollbarAlias,
		TextData:                 shared.TextData,
		ViewportX:                shared.ViewportXLocation,
		ViewportY:                shared.ViewportYLocation,
		CursorX:                  shared.CursorXLocation,
		CursorY:                  shared.CursorYLocation,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared TextboxEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewTexboxEntry(existingTextboxEntry ...*TextboxEntryType) TextboxEntryType {
	var textboxEntry TextboxEntryType
	textboxEntry.BaseControlType = NewBaseControl()

	if existingTextboxEntry != nil {
		textboxEntry.BaseControlType = existingTextboxEntry[0].BaseControlType
		textboxEntry.HorizontalScrollbarAlias = existingTextboxEntry[0].HorizontalScrollbarAlias
		textboxEntry.TextData = existingTextboxEntry[0].TextData
		textboxEntry.ViewportXLocation = existingTextboxEntry[0].ViewportXLocation
		textboxEntry.ViewportYLocation = existingTextboxEntry[0].ViewportYLocation
		textboxEntry.CursorXLocation = existingTextboxEntry[0].CursorXLocation
		textboxEntry.CursorYLocation = existingTextboxEntry[0].CursorYLocation
	}
	return textboxEntry
}

func IsTextboxEntryEqual(sourceTextboxEntry *TextboxEntryType, targetTextboxEntry *TextboxEntryType) bool {
	if sourceTextboxEntry.BaseControlType == targetTextboxEntry.BaseControlType &&
		sourceTextboxEntry.HorizontalScrollbarAlias == targetTextboxEntry.HorizontalScrollbarAlias &&
		reflect.DeepEqual(sourceTextboxEntry.TextData, targetTextboxEntry.TextData) &&
		sourceTextboxEntry.ViewportXLocation == targetTextboxEntry.ViewportXLocation &&
		sourceTextboxEntry.ViewportYLocation == targetTextboxEntry.ViewportYLocation &&
		sourceTextboxEntry.CursorXLocation == targetTextboxEntry.CursorXLocation &&
		sourceTextboxEntry.CursorYLocation == targetTextboxEntry.CursorYLocation {
		return true
	}
	return false
}

func GetTextboxAlias(entry *TextboxEntryType) string {
	return entry.Alias
}
