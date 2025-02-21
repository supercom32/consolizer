package types

import (
	"encoding/json"
)

type TextboxEntryType struct {
	Alias                    string
	HorizontalScrollbarAlias string
	VerticalScrollbarAlias   string
	StyleEntry               TuiStyleEntryType
	TextData                 [][]rune
	XLocation                int
	YLocation                int
	Width                    int
	Height                   int
	ViewportXLocation        int
	ViewportYLocation        int
	CursorXLocation          int
	CursorYLocation          int
	IsEnabled                bool
	IsVisible                bool
	IsBorderDrawn            bool
}

func (shared TextboxEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		Alias                    string
		HorizontalScrollbarAlias string
		VerticalScrollbarAlias   string
		StyleEntry               TuiStyleEntryType
		TextData                 [][]rune
		XLocation                int
		YLocation                int
		Width                    int
		Height                   int
		ViewportXLocation        int
		ViewportYLocation        int
		CursorXLocation          int
		CursorYLocation          int
		IsEnabled                bool
		IsVisible                bool
		IsBorderDrawn            bool
	}{
		Alias:                    shared.Alias,
		HorizontalScrollbarAlias: shared.HorizontalScrollbarAlias,
		VerticalScrollbarAlias:   shared.VerticalScrollbarAlias,
		StyleEntry:               shared.StyleEntry,
		TextData:                 shared.TextData,
		XLocation:                shared.XLocation,
		YLocation:                shared.YLocation,
		Width:                    shared.Width,
		Height:                   shared.Height,
		ViewportXLocation:        shared.ViewportXLocation,
		ViewportYLocation:        shared.ViewportYLocation,
		CursorXLocation:          shared.CursorXLocation,
		CursorYLocation:          shared.CursorYLocation,
		IsEnabled:                shared.IsEnabled,
		IsVisible:                shared.IsVisible,
		IsBorderDrawn:            shared.IsBorderDrawn,
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
	if existingTextboxEntry != nil {
		textboxEntry.Alias = existingTextboxEntry[0].Alias
		textboxEntry.HorizontalScrollbarAlias = existingTextboxEntry[0].HorizontalScrollbarAlias
		textboxEntry.VerticalScrollbarAlias = existingTextboxEntry[0].VerticalScrollbarAlias
		textboxEntry.StyleEntry = NewTuiStyleEntry(&existingTextboxEntry[0].StyleEntry)
		textboxEntry.TextData = existingTextboxEntry[0].TextData
		textboxEntry.XLocation = existingTextboxEntry[0].XLocation
		textboxEntry.YLocation = existingTextboxEntry[0].YLocation
		textboxEntry.Width = existingTextboxEntry[0].Width
		textboxEntry.Height = existingTextboxEntry[0].Height
		textboxEntry.ViewportXLocation = existingTextboxEntry[0].ViewportXLocation
		textboxEntry.ViewportYLocation = existingTextboxEntry[0].ViewportYLocation
		textboxEntry.CursorXLocation = existingTextboxEntry[0].CursorXLocation
		textboxEntry.CursorYLocation = existingTextboxEntry[0].CursorYLocation
		textboxEntry.IsEnabled = existingTextboxEntry[0].IsEnabled
		textboxEntry.IsVisible = existingTextboxEntry[0].IsVisible
		textboxEntry.IsBorderDrawn = existingTextboxEntry[0].IsBorderDrawn
	}
	return textboxEntry
}

func IsTexboxEqual(sourceTextboxEntry *TextboxEntryType, targetTextboxEntry *TextboxEntryType) bool {
	if sourceTextboxEntry.StyleEntry == targetTextboxEntry.StyleEntry &&
		sourceTextboxEntry.XLocation == targetTextboxEntry.XLocation &&
		sourceTextboxEntry.YLocation == targetTextboxEntry.YLocation &&
		sourceTextboxEntry.Width == targetTextboxEntry.Width &&
		sourceTextboxEntry.Height == targetTextboxEntry.Height &&
		sourceTextboxEntry.IsEnabled == targetTextboxEntry.IsEnabled &&
		sourceTextboxEntry.IsVisible == targetTextboxEntry.IsVisible &&
		sourceTextboxEntry.IsBorderDrawn == targetTextboxEntry.IsBorderDrawn {
		return true
	}
	return false
}
