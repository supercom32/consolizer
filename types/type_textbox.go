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
	// Highlight positions
	HighlightStartX        int
	HighlightStartY        int
	HighlightEndX          int
	HighlightEndY          int
	IsHighlightActive      bool
	IsHighlightModeToggled bool
}

/*
GetAlias allows you to retrieve the alias of a textbox control. In addition, the following
information should be noted:

- Returns the unique identifier for the textbox.
- This alias is used to reference the textbox in other operations.
- The alias is set when the textbox is created.
*/
func (shared TextboxEntryType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON allows you to serialize a textbox control to JSON. In addition, the following
information should be noted:

- Converts the textbox's state to a JSON representation.
- Includes the base control properties and textbox-specific fields.
- Used for saving and loading textbox configurations.
*/
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
		HighlightStartX          int
		HighlightStartY          int
		HighlightEndX            int
		HighlightEndY            int
		IsHighlightActive        bool
		IsHighlightModeToggled   bool
	}{
		BaseControlType:          shared.BaseControlType,
		HorizontalScrollbarAlias: shared.HorizontalScrollbarAlias,
		VerticalScrollbarAlias:   shared.VerticalScrollbarAlias,
		TextData:                 shared.TextData,
		ViewportX:                shared.ViewportXLocation,
		ViewportY:                shared.ViewportYLocation,
		CursorX:                  shared.CursorXLocation,
		CursorY:                  shared.CursorYLocation,
		HighlightStartX:          shared.HighlightStartX,
		HighlightStartY:          shared.HighlightStartY,
		HighlightEndX:            shared.HighlightEndX,
		HighlightEndY:            shared.HighlightEndY,
		IsHighlightActive:        shared.IsHighlightActive,
		IsHighlightModeToggled:   shared.IsHighlightModeToggled,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump allows you to get a JSON string representation of a textbox control. In addition,
the following information should be noted:

- Returns a formatted JSON string of the textbox's state.
- Useful for debugging and logging purposes.
- Panics if JSON marshaling fails.
*/
func (shared TextboxEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewTexboxEntry allows you to create a new textbox control. In addition, the following
information should be noted:

- Initializes a textbox with default values.
- Can optionally copy properties from an existing textbox.
- Sets up the base control properties and textbox-specific fields.
*/
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
		textboxEntry.HighlightStartX = existingTextboxEntry[0].HighlightStartX
		textboxEntry.HighlightStartY = existingTextboxEntry[0].HighlightStartY
		textboxEntry.HighlightEndX = existingTextboxEntry[0].HighlightEndX
		textboxEntry.HighlightEndY = existingTextboxEntry[0].HighlightEndY
		textboxEntry.IsHighlightActive = existingTextboxEntry[0].IsHighlightActive
		textboxEntry.IsHighlightModeToggled = existingTextboxEntry[0].IsHighlightModeToggled
	}
	return textboxEntry
}

/*
IsTextboxEntryEqual allows you to compare two textbox controls for equality. In addition, the following
information should be noted:

- Compares all properties of both textboxes.
- Returns true if all properties match, false otherwise.
- Used for change detection and state synchronization.
*/
func IsTextboxEntryEqual(sourceTextboxEntry *TextboxEntryType, targetTextboxEntry *TextboxEntryType) bool {
	return sourceTextboxEntry.BaseControlType.IsEqual(&targetTextboxEntry.BaseControlType) &&
		sourceTextboxEntry.HorizontalScrollbarAlias == targetTextboxEntry.HorizontalScrollbarAlias &&
		sourceTextboxEntry.VerticalScrollbarAlias == targetTextboxEntry.VerticalScrollbarAlias &&
		reflect.DeepEqual(sourceTextboxEntry.TextData, targetTextboxEntry.TextData) &&
		sourceTextboxEntry.ViewportXLocation == targetTextboxEntry.ViewportXLocation &&
		sourceTextboxEntry.ViewportYLocation == targetTextboxEntry.ViewportYLocation &&
		sourceTextboxEntry.CursorXLocation == targetTextboxEntry.CursorXLocation &&
		sourceTextboxEntry.CursorYLocation == targetTextboxEntry.CursorYLocation &&
		sourceTextboxEntry.HighlightStartX == targetTextboxEntry.HighlightStartX &&
		sourceTextboxEntry.HighlightStartY == targetTextboxEntry.HighlightStartY &&
		sourceTextboxEntry.HighlightEndX == targetTextboxEntry.HighlightEndX &&
		sourceTextboxEntry.HighlightEndY == targetTextboxEntry.HighlightEndY &&
		sourceTextboxEntry.IsHighlightActive == targetTextboxEntry.IsHighlightActive &&
		sourceTextboxEntry.IsHighlightModeToggled == targetTextboxEntry.IsHighlightModeToggled
}

/*
GetTextboxAlias allows you to retrieve the alias of a textbox control. In addition, the following
information should be noted:

- Returns the unique identifier for the textbox.
- This is a convenience method that delegates to GetAlias.
- The alias is used to reference the textbox in other operations.
*/
func GetTextboxAlias(entry *TextboxEntryType) string {
	return entry.Alias
}
