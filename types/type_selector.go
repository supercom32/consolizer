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

func (shared SelectorEntryType) GetAlias() string {
	return shared.Alias
}

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

func (shared SelectorEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

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
