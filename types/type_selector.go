package types

import (
	"encoding/json"
)

type SelectorEntryType struct {
	ScrollBarAlias   string
	StyleEntry       TuiStyleEntryType
	SelectionEntry   SelectionEntryType
	XLocation        int
	YLocation        int
	SelectorHeight   int
	ItemWidth        int
	NumberOfColumns  int
	ViewportPosition int
	ItemHighlighted  int
	ItemSelected     int
	IsVisible        bool
	IsBorderDrawn    bool
}

func (shared SelectorEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		ScrollBarAlias   string
		StyleEntry       TuiStyleEntryType
		SelectionEntry   SelectionEntryType
		XLocation        int
		YLocation        int
		ItemWidth        int
		NumberOfColumns  int
		ViewportPosition int
		ItemHighlighted  int
		ItemSelected     int
		IsVisible        bool
		IsBorderDrawn    bool
	}{
		ScrollBarAlias:   shared.ScrollBarAlias,
		StyleEntry:       shared.StyleEntry,
		SelectionEntry:   shared.SelectionEntry,
		XLocation:        shared.XLocation,
		YLocation:        shared.YLocation,
		ItemWidth:        shared.ItemWidth,
		NumberOfColumns:  shared.NumberOfColumns,
		ViewportPosition: shared.ViewportPosition,
		ItemHighlighted:  shared.ItemHighlighted,
		ItemSelected:     shared.ItemSelected,
		IsVisible:        shared.IsVisible,
		IsBorderDrawn:    shared.IsBorderDrawn,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared SelectorEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewSelectorEntry(existingSelectorEntry ...*SelectorEntryType) SelectorEntryType {
	var selectorEntry SelectorEntryType
	if existingSelectorEntry != nil {
		selectorEntry.ScrollBarAlias = existingSelectorEntry[0].ScrollBarAlias
		selectorEntry.StyleEntry = existingSelectorEntry[0].StyleEntry
		selectorEntry.SelectionEntry = existingSelectorEntry[0].SelectionEntry
		selectorEntry.XLocation = existingSelectorEntry[0].XLocation
		selectorEntry.YLocation = existingSelectorEntry[0].YLocation
		selectorEntry.ItemWidth = existingSelectorEntry[0].ItemWidth
		selectorEntry.NumberOfColumns = existingSelectorEntry[0].NumberOfColumns
		selectorEntry.ViewportPosition = existingSelectorEntry[0].ViewportPosition
		selectorEntry.ItemHighlighted = existingSelectorEntry[0].ItemHighlighted
		selectorEntry.ItemSelected = existingSelectorEntry[0].ItemSelected
		selectorEntry.IsVisible = existingSelectorEntry[0].IsVisible
		selectorEntry.IsBorderDrawn = existingSelectorEntry[0].IsBorderDrawn
	}
	return selectorEntry
}
