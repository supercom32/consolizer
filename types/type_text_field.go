package types

import (
	"encoding/json"
	"reflect"
)

/*
TextFieldEntryType is a structure which contains the properties for a text field control. In addition, the following should be noted:

- This type is used to represent the state of a text field in the TUI.

Example:
    var textField types.TextFieldEntryType
*/
type TextFieldEntryType struct {
	BaseControlType
	MaxLengthAllowed    int
	DefaultValue        string
	CursorPosition      int
	ViewportPosition    int
	IsPasswordProtected bool
	CurrentValue        []rune
	// Highlight positions
	HighlightStart         int
	HighlightEnd           int
	IsHighlightActive      bool
	IsHighlightModeToggled bool
}

/*
GetAlias is a method which allows you to retrieve the alias of a text field control. In addition, the following should be noted:

- Returns the unique identifier for the text field.

- This alias is used to reference the text field in other operations.

- The alias is set when the text field is created.

Example:
    instance.GetAlias()
*/
func (shared TextFieldEntryType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON is a method which allows you to serialize a text field control to JSON. In addition, the following should be noted:

- Converts the text field's state to a JSON representation.

- Includes the base control properties and text field-specific fields.

- Used for saving and loading text field configurations.

Example:
    instance.MarshalJSON()
*/
func (shared TextFieldEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		MaxLengthAllowed       int
		DefaultValue           string
		CursorPosition         int
		ViewportPosition       int
		IsPasswordProtected    bool
		CurrentValue           []rune
		HighlightStart         int
		HighlightEnd           int
		IsHighlightActive      bool
		IsHighlightModeToggled bool
	}{
		BaseControlType:        shared.BaseControlType,
		MaxLengthAllowed:       shared.MaxLengthAllowed,
		DefaultValue:           shared.DefaultValue,
		CursorPosition:         shared.CursorPosition,
		ViewportPosition:       shared.ViewportPosition,
		IsPasswordProtected:    shared.IsPasswordProtected,
		CurrentValue:           shared.CurrentValue,
		HighlightStart:         shared.HighlightStart,
		HighlightEnd:           shared.HighlightEnd,
		IsHighlightActive:      shared.IsHighlightActive,
		IsHighlightModeToggled: shared.IsHighlightModeToggled,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump is a method which allows you to get a JSON string representation of a text field control. In addition, the following should be noted:

- Returns a formatted JSON string of the text field's state.

- Useful for debugging and logging purposes.

- Panics if JSON marshaling fails.

Example:
    instance.GetEntryAsJsonDump()
*/
func (shared TextFieldEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewTextFieldEntry is a constructor which allows you to create a new text field control. In addition, the following should be noted:

- Initializes a text field with default values.

- Can optionally copy properties from an existing text field.

- Sets up the base control properties and text field-specific fields.

Example:
    NewTextFieldEntry(existingTextFieldEntry)
*/
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

/*
IsTextFieldEntryEqual is a method which allows you to compare two text field controls for equality. In addition, the following should be noted:

- Compares all properties of both text fields.

- Returns true if all properties match, false otherwise.

- Used for change detection and state synchronization.

Example:
    IsTextFieldEntryEqual(sourceTextFieldEntry, targetTextFieldEntry)
*/
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

/*
GetTextFieldAlias is a method which allows you to retrieve the alias of a text field control. In addition, the following should be noted:

- Returns the unique identifier for the text field.

- This is a convenience method that delegates to GetAlias.

- The alias is used to reference the text field in other operations.

Example:
    GetTextFieldAlias(entry)
*/
func GetTextFieldAlias(entry *TextFieldEntryType) string {
	return entry.Alias
}
