package types

import (
	"encoding/json"
)

type ScrollbarEntryType struct {
	BaseControlType
	Length             int
	MaxScrollValue     int
	ScrollValue        int
	HandlePosition     int
	IsHorizontal       bool
	ScrollIncrement    int
	ParentControlAlias string // Empty for standalone scrollbars
	ParentControlType  int    // Constants.CellTypeTextbox, etc.
}

/*
GetAlias is a method which allows you to retrieve the alias of a scroll bar control. In addition, the following
information should be noted:

- Returns the unique identifier for the scroll bar.

- This alias is used to reference the scroll bar in other operations.

- The alias is set when the scroll bar is created.

:return: string

Example:

	instance.GetAlias()
*/
func (shared ScrollbarEntryType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON is a method which allows you to serialize a scroll bar control to JSON. In addition, the following
information should be noted:

- Converts the scroll bar's state to a JSON representation.

- Includes the base control properties and scroll bar-specific fields.

- Used for saving and loading scroll bar configurations.

Example:

	instance.MarshalJSON()
*/
func (shared ScrollbarEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		Length             int
		MaxScrollValue     int
		ScrollValue        int
		HandlePosition     int
		IsHorizontal       bool
		ScrollIncrement    int
		ParentControlAlias string
		ParentControlType  int
	}{
		BaseControlType:    shared.BaseControlType,
		Length:             shared.Length,
		MaxScrollValue:     shared.MaxScrollValue,
		ScrollValue:        shared.ScrollValue,
		HandlePosition:     shared.HandlePosition,
		IsHorizontal:       shared.IsHorizontal,
		ScrollIncrement:    shared.ScrollIncrement,
		ParentControlAlias: shared.ParentControlAlias,
		ParentControlType:  shared.ParentControlType,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump is a method which allows you to get a JSON string representation of a scroll bar control. In
addition, the following information should be noted:

- Returns a formatted JSON string of the scroll bar's state.

- Useful for debugging and logging purposes.

- Panics if JSON marshaling fails.

:return: string

Example:

	instance.GetEntryAsJsonDump()
*/
func (shared ScrollbarEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
GetScrollBarAlias is a method which allows you to retrieve the alias of a scroll bar control. In addition, the following
information should be noted:

- Returns the unique identifier for the scroll bar.

- This is a convenience method that delegates to GetAlias.

- The alias is used to reference the scroll bar in other operations.

:param entry: The entry parameter.

:return: string

Example:

	GetScrollBarAlias(entry)
*/
func GetScrollBarAlias(entry *ScrollbarEntryType) string {
	return entry.Alias
}

/*
NewScrollbarEntry is a constructor which allows you to create a new scroll bar control. In addition, the following
information should be noted:

- Initializes a scroll bar with default values.

- Can optionally copy properties from an existing scroll bar.

- Sets up the base control properties and scroll bar-specific fields.

:param existingScrollbarEntry: The existingScrollbarEntry parameter.

:return: ScrollbarEntryType

Example:

	NewScrollbarEntry(existingScrollbarEntry)
*/
func NewScrollbarEntry(existingScrollbarEntry ...*ScrollbarEntryType) ScrollbarEntryType {
	var scrollbarEntry ScrollbarEntryType
	if existingScrollbarEntry != nil {
		scrollbarEntry.Alias = existingScrollbarEntry[0].Alias
		scrollbarEntry.Length = existingScrollbarEntry[0].Length
		scrollbarEntry.MaxScrollValue = existingScrollbarEntry[0].MaxScrollValue
		scrollbarEntry.ScrollValue = existingScrollbarEntry[0].ScrollValue
		scrollbarEntry.HandlePosition = existingScrollbarEntry[0].HandlePosition
		scrollbarEntry.IsHorizontal = existingScrollbarEntry[0].IsHorizontal
		scrollbarEntry.ScrollIncrement = existingScrollbarEntry[0].ScrollIncrement
		scrollbarEntry.ParentControlAlias = existingScrollbarEntry[0].ParentControlAlias
		scrollbarEntry.ParentControlType = existingScrollbarEntry[0].ParentControlType
	}
	return scrollbarEntry
}

/*
IsScrollbarEntryEqual is a method which allows you to isscrollbarentryequal.

:param sourceScrollbarEntry: The sourceScrollbarEntry parameter.
:param targetScrollBarEntry: The targetScrollBarEntry parameter.

:return: bool

Example:

	IsScrollbarEntryEqual(sourceScrollbarEntry, targetScrollBarEntry)
*/
func IsScrollbarEntryEqual(sourceScrollbarEntry *ScrollbarEntryType, targetScrollBarEntry *ScrollbarEntryType) bool {
	if sourceScrollbarEntry.Length == targetScrollBarEntry.Length &&
		sourceScrollbarEntry.MaxScrollValue == targetScrollBarEntry.MaxScrollValue &&
		sourceScrollbarEntry.ScrollValue == targetScrollBarEntry.ScrollValue &&
		sourceScrollbarEntry.HandlePosition == targetScrollBarEntry.HandlePosition &&
		sourceScrollbarEntry.IsHorizontal == targetScrollBarEntry.IsHorizontal &&
		sourceScrollbarEntry.ScrollIncrement == targetScrollBarEntry.ScrollIncrement &&
		sourceScrollbarEntry.ParentControlAlias == targetScrollBarEntry.ParentControlAlias &&
		sourceScrollbarEntry.ParentControlType == targetScrollBarEntry.ParentControlType {
		return true
	}
	return false
}
