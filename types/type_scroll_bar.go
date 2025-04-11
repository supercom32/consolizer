package types

import (
	"encoding/json"
)

type ScrollbarEntryType struct {
	BaseControlType
	Length          int
	MaxScrollValue  int
	ScrollValue     int
	HandlePosition  int
	IsHorizontal    bool
	ScrollIncrement int
}

/*
GetAlias allows you to retrieve the alias of a scroll bar control. In addition, the following
information should be noted:

- Returns the unique identifier for the scroll bar.
- This alias is used to reference the scroll bar in other operations.
- The alias is set when the scroll bar is created.
*/
func (shared ScrollbarEntryType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON allows you to serialize a scroll bar control to JSON. In addition, the following
information should be noted:

- Converts the scroll bar's state to a JSON representation.
- Includes the base control properties and scroll bar-specific fields.
- Used for saving and loading scroll bar configurations.
*/
func (shared ScrollbarEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		Length          int
		MaxScrollValue  int
		ScrollValue     int
		HandlePosition  int
		IsHorizontal    bool
		ScrollIncrement int
	}{
		BaseControlType: shared.BaseControlType,
		Length:          shared.Length,
		MaxScrollValue:  shared.MaxScrollValue,
		ScrollValue:     shared.ScrollValue,
		HandlePosition:  shared.HandlePosition,
		IsHorizontal:    shared.IsHorizontal,
		ScrollIncrement: shared.ScrollIncrement,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump allows you to get a JSON string representation of a scroll bar control. In addition,
the following information should be noted:

- Returns a formatted JSON string of the scroll bar's state.
- Useful for debugging and logging purposes.
- Panics if JSON marshaling fails.
*/
func (shared ScrollbarEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
GetScrollBarAlias allows you to retrieve the alias of a scroll bar control. In addition, the following
information should be noted:

- Returns the unique identifier for the scroll bar.
- This is a convenience method that delegates to GetAlias.
- The alias is used to reference the scroll bar in other operations.
*/
func GetScrollBarAlias(entry *ScrollbarEntryType) string {
	return entry.Alias
}

/*
NewScrollbarEntry allows you to create a new scroll bar control. In addition, the following
information should be noted:

- Initializes a scroll bar with default values.
- Can optionally copy properties from an existing scroll bar.
- Sets up the base control properties and scroll bar-specific fields.
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
	}
	return scrollbarEntry
}

func IsScrollbarEntryEqual(sourceScrollbarEntry *ScrollbarEntryType, targetScrollBarEntry *ScrollbarEntryType) bool {
	if sourceScrollbarEntry.Length == targetScrollBarEntry.Length &&
		sourceScrollbarEntry.MaxScrollValue == targetScrollBarEntry.MaxScrollValue &&
		sourceScrollbarEntry.ScrollValue == targetScrollBarEntry.ScrollValue &&
		sourceScrollbarEntry.HandlePosition == targetScrollBarEntry.HandlePosition &&
		sourceScrollbarEntry.IsHorizontal == targetScrollBarEntry.IsHorizontal &&
		sourceScrollbarEntry.ScrollIncrement == targetScrollBarEntry.ScrollIncrement {
		return true
	}
	return false
}
