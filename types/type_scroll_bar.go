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

func (shared ScrollbarEntryType) GetAlias() string {
	return shared.Alias
}

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

func (shared ScrollbarEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func GetScrollBarAlias(entry *ScrollbarEntryType) string {
	return entry.Alias
}

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
