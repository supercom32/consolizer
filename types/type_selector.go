package types

import (
	"encoding/json"
	"reflect"
)

/*
SelectorEntryType is a structure which contains the properties for a selector control. In addition, the following should be noted:

- This type is used to represent the state of a selector in the TUI.

Example:
    var selector types.SelectorEntryType
*/
type SelectorEntryType struct {
	BaseControlType
	ScrollbarAlias       string
	SelectionEntry       SelectionEntryType
	ItemWidth            int
	ColumnCount          int
	NumberOfColumns      int
	ViewportX            int
	ViewportY            int
	ViewportPosition     int
	ItemHighlighted      int
	ItemSelected         int
	HighlightOnClickOnly bool
	IsNewItemSelected    bool
}

/*
GetAlias is a method which allows you to retrieve the alias of a selector control. In addition, the following should be noted:

- Returns the unique identifier for the selector.

- This alias is used to reference the selector in other operations.

- The alias is set when the selector is created.

Example:
    instance.GetAlias()
*/
func (shared SelectorEntryType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON is a method which allows you to serialize a selector control to JSON. In addition, the following should be noted:

- Converts the selector's state to a JSON representation.

- Includes the base control properties and selector-specific fields.

- Used for saving and loading selector configurations.

Example:
    instance.MarshalJSON()
*/
func (shared SelectorEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		ScrollbarAlias       string
		SelectionEntry       SelectionEntryType
		ItemWidth            int
		ColumnCount          int
		NumberOfColumns      int
		ViewportX            int
		ViewportY            int
		ViewportPosition     int
		ItemHighlighted      int
		ItemSelected         int
		HighlightOnClickOnly bool
		IsNewItemSelected    bool
	}{
		BaseControlType:      shared.BaseControlType,
		ScrollbarAlias:       shared.ScrollbarAlias,
		SelectionEntry:       shared.SelectionEntry,
		ItemWidth:            shared.ItemWidth,
		ColumnCount:          shared.ColumnCount,
		NumberOfColumns:      shared.NumberOfColumns,
		ViewportX:            shared.ViewportX,
		ViewportY:            shared.ViewportY,
		ViewportPosition:     shared.ViewportPosition,
		ItemHighlighted:      shared.ItemHighlighted,
		ItemSelected:         shared.ItemSelected,
		HighlightOnClickOnly: shared.HighlightOnClickOnly,
		IsNewItemSelected:    shared.IsNewItemSelected,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump is a method which allows you to get a JSON string representation of a selector control. In addition, the following should be noted:

- Returns a formatted JSON string of the selector's state.

- Useful for debugging and logging purposes.

- Panics if JSON marshaling fails.

Example:
    instance.GetEntryAsJsonDump()
*/
func (shared SelectorEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewSelectorEntry is a constructor which allows you to create a new selector control. In addition, the following should be noted:

- Initializes a selector with default values.

- Can optionally copy properties from an existing selector.

- Sets up the base control properties and selector-specific fields.

Example:
    NewSelectorEntry(existingSelectorEntry)
*/
func NewSelectorEntry(existingSelectorEntry ...*SelectorEntryType) SelectorEntryType {
	var selectorEntry SelectorEntryType
	selectorEntry.BaseControlType = NewBaseControl()
	selectorEntry.IsNewItemSelected = false

	if existingSelectorEntry != nil {
		selectorEntry.BaseControlType = existingSelectorEntry[0].BaseControlType
		selectorEntry.ScrollbarAlias = existingSelectorEntry[0].ScrollbarAlias
		selectorEntry.SelectionEntry = existingSelectorEntry[0].SelectionEntry
		selectorEntry.ItemWidth = existingSelectorEntry[0].ItemWidth
		selectorEntry.ColumnCount = existingSelectorEntry[0].ColumnCount
		selectorEntry.NumberOfColumns = existingSelectorEntry[0].NumberOfColumns
		selectorEntry.ViewportX = existingSelectorEntry[0].ViewportX
		selectorEntry.ViewportY = existingSelectorEntry[0].ViewportY
		selectorEntry.ViewportPosition = existingSelectorEntry[0].ViewportPosition
		selectorEntry.ItemHighlighted = existingSelectorEntry[0].ItemHighlighted
		selectorEntry.HighlightOnClickOnly = existingSelectorEntry[0].HighlightOnClickOnly
		selectorEntry.ItemSelected = existingSelectorEntry[0].ItemSelected
		selectorEntry.IsNewItemSelected = existingSelectorEntry[0].IsNewItemSelected
	}
	return selectorEntry
}

/*
IsSelectorEntryEqual is a method which allows you to compare two selector controls for equality. In addition, the following should be noted:

- Compares all properties of both selectors.

Example:
    IsSelectorEntryEqual(sourceSelectorEntry, targetSelectorEntry)
*/
func IsSelectorEntryEqual(sourceSelectorEntry *SelectorEntryType, targetSelectorEntry *SelectorEntryType) bool {
	if sourceSelectorEntry.BaseControlType == targetSelectorEntry.BaseControlType &&
		sourceSelectorEntry.ScrollbarAlias == targetSelectorEntry.ScrollbarAlias &&
		reflect.DeepEqual(sourceSelectorEntry.SelectionEntry, targetSelectorEntry.SelectionEntry) &&
		sourceSelectorEntry.ItemWidth == targetSelectorEntry.ItemWidth &&
		sourceSelectorEntry.ColumnCount == targetSelectorEntry.ColumnCount &&
		sourceSelectorEntry.NumberOfColumns == targetSelectorEntry.NumberOfColumns &&
		sourceSelectorEntry.ViewportX == targetSelectorEntry.ViewportX &&
		sourceSelectorEntry.ViewportY == targetSelectorEntry.ViewportY &&
		sourceSelectorEntry.ViewportPosition == targetSelectorEntry.ViewportPosition &&
		sourceSelectorEntry.ItemHighlighted == targetSelectorEntry.ItemHighlighted &&
		sourceSelectorEntry.HighlightOnClickOnly == targetSelectorEntry.HighlightOnClickOnly &&
		sourceSelectorEntry.ItemSelected == targetSelectorEntry.ItemSelected &&
		sourceSelectorEntry.IsNewItemSelected == targetSelectorEntry.IsNewItemSelected {
		return true
	}
	return false
}

/*
GetSelectorAlias is a method which allows you to retrieve the alias of a selector control. In addition, the following should be noted:

- Returns the unique identifier for the selector.

Example:
    GetSelectorAlias(entry)
*/
func GetSelectorAlias(entry *SelectorEntryType) string {
	return entry.Alias
}
