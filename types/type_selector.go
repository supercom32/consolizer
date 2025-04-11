package types

import (
	"encoding/json"
	"reflect"
)

type SelectorEntryType struct {
	BaseControlType
	ScrollbarAlias   string
	SelectionEntry   SelectionEntryType
	ItemWidth        int
	ColumnCount      int
	NumberOfColumns  int
	ViewportX        int
	ViewportY        int
	ViewportPosition int
	ItemHighlighted  int
	ItemSelected     int
}

/*
GetAlias allows you to retrieve the alias of a selector control. In addition, the following
information should be noted:

- Returns the unique identifier for the selector.
- This alias is used to reference the selector in other operations.
- The alias is set when the selector is created.
*/
func (shared SelectorEntryType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON allows you to serialize a selector control to JSON. In addition, the following
information should be noted:

- Converts the selector's state to a JSON representation.
- Includes the base control properties and selector-specific fields.
- Used for saving and loading selector configurations.
*/
func (shared SelectorEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		ScrollbarAlias   string
		SelectionEntry   SelectionEntryType
		ItemWidth        int
		ColumnCount      int
		NumberOfColumns  int
		ViewportX        int
		ViewportY        int
		ViewportPosition int
		ItemHighlighted  int
		ItemSelected     int
	}{
		BaseControlType:  shared.BaseControlType,
		ScrollbarAlias:   shared.ScrollbarAlias,
		SelectionEntry:   shared.SelectionEntry,
		ItemWidth:        shared.ItemWidth,
		ColumnCount:      shared.ColumnCount,
		NumberOfColumns:  shared.NumberOfColumns,
		ViewportX:        shared.ViewportX,
		ViewportY:        shared.ViewportY,
		ViewportPosition: shared.ViewportPosition,
		ItemHighlighted:  shared.ItemHighlighted,
		ItemSelected:     shared.ItemSelected,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump allows you to get a JSON string representation of a selector control. In addition,
the following information should be noted:

- Returns a formatted JSON string of the selector's state.
- Useful for debugging and logging purposes.
- Panics if JSON marshaling fails.
*/
func (shared SelectorEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewSelectorEntry allows you to create a new selector control. In addition, the following
information should be noted:

- Initializes a selector with default values.
- Can optionally copy properties from an existing selector.
- Sets up the base control properties and selector-specific fields.
*/
func NewSelectorEntry(existingSelectorEntry ...*SelectorEntryType) SelectorEntryType {
	var selectorEntry SelectorEntryType
	selectorEntry.BaseControlType = NewBaseControl()

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
		selectorEntry.ItemSelected = existingSelectorEntry[0].ItemSelected
	}
	return selectorEntry
}

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
		sourceSelectorEntry.ItemSelected == targetSelectorEntry.ItemSelected {
		return true
	}
	return false
}

func GetSelectorAlias(entry *SelectorEntryType) string {
	return entry.Alias
}
