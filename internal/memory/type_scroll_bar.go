package memory

import (
	"encoding/json"
	"sync"
)

type ScrollbarEntryType struct {
	Mutex                   sync.Mutex
	StyleEntry              TuiStyleEntryType
	ScrollBarAlias          string
	XLocation               int
	YLocation               int
	Length         int
	MaxScrollValue int
	ScrollValue    int
	HandlePosition int
	IsVisible bool
	IsHorizontal            bool
	IsEnabled bool
	ScrollIncrement int
}

func (shared ScrollbarEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		StyleEntry  TuiStyleEntryType
		ScrollBarAlias string
		XLocation   int
		YLocation   int
		Length      int
		MaxScrollValue   int
		ScrollValue int
		ScrollPosition int
		IsVisible bool
		IsHorizontal bool
		IsEnabled bool
		ScrollIncrement int
	}{
		StyleEntry: shared.StyleEntry,
		ScrollBarAlias: shared.ScrollBarAlias,
		XLocation: shared.XLocation,
		YLocation: shared.YLocation,
		Length: shared.Length,
		MaxScrollValue: shared.MaxScrollValue,
		ScrollPosition: shared.HandlePosition,
		IsVisible: shared.IsVisible,
		IsHorizontal: shared.IsHorizontal,
		IsEnabled: shared.IsEnabled,
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

func NewScrollbarEntry(existingScrollbarEntry ...*ScrollbarEntryType) ScrollbarEntryType {
	var scrollbarEntry ScrollbarEntryType
	if existingScrollbarEntry != nil {
		scrollbarEntry.StyleEntry = NewTuiStyleEntry(&existingScrollbarEntry[0].StyleEntry)
		scrollbarEntry.XLocation = existingScrollbarEntry[0].XLocation
		scrollbarEntry.YLocation = existingScrollbarEntry[0].YLocation
		scrollbarEntry.Length = existingScrollbarEntry[0].Length
		scrollbarEntry.MaxScrollValue = existingScrollbarEntry[0].MaxScrollValue
		scrollbarEntry.ScrollValue = existingScrollbarEntry[0].ScrollValue
		scrollbarEntry.HandlePosition = existingScrollbarEntry[0].HandlePosition
		scrollbarEntry.IsVisible = existingScrollbarEntry[0].IsVisible
		scrollbarEntry.IsHorizontal = existingScrollbarEntry[0].IsHorizontal
		scrollbarEntry.IsEnabled = existingScrollbarEntry[0].IsEnabled
		scrollbarEntry.ScrollIncrement = existingScrollbarEntry[0].ScrollIncrement
	}
	return scrollbarEntry
}

func IsScrollbarEntryEqual(sourceScrollbarEntry *ScrollbarEntryType, targetScrollBarEntry *ScrollbarEntryType) bool {
	if sourceScrollbarEntry.StyleEntry == targetScrollBarEntry.StyleEntry &&
		sourceScrollbarEntry.XLocation == targetScrollBarEntry.XLocation &&
		sourceScrollbarEntry.YLocation == targetScrollBarEntry.YLocation &&
		sourceScrollbarEntry.Length == targetScrollBarEntry.Length &&
		sourceScrollbarEntry.MaxScrollValue == targetScrollBarEntry.MaxScrollValue &&
		sourceScrollbarEntry.ScrollValue == targetScrollBarEntry.ScrollValue &&
		sourceScrollbarEntry.HandlePosition == targetScrollBarEntry.HandlePosition &&
		sourceScrollbarEntry.IsVisible == targetScrollBarEntry.IsVisible &&
		sourceScrollbarEntry.IsHorizontal == targetScrollBarEntry.IsHorizontal &&
		sourceScrollbarEntry.IsEnabled == targetScrollBarEntry.IsEnabled &&
		sourceScrollbarEntry.ScrollIncrement == targetScrollBarEntry.ScrollIncrement {
		return true
	}
	return false
}
