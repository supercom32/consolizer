package types

import (
	"encoding/json"
)

type DropdownEntryType struct {
	BaseControlType
	SelectionEntry   SelectionEntryType
	ScrollbarAlias   string
	SelectorAlias    string
	ItemWidth        int
	ItemSelected     int
	IsTrayOpen       bool
	ViewportPosition int
}

func (shared DropdownEntryType) GetAlias() string {
	return shared.Alias
}

func (shared DropdownEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		SelectionEntry   SelectionEntryType
		ScrollbarAlias   string
		SelectorAlias    string
		ItemWidth        int
		ItemSelected     int
		IsTrayOpen       bool
		ViewportPosition int
	}{
		BaseControlType:  shared.BaseControlType,
		SelectionEntry:   shared.SelectionEntry,
		ScrollbarAlias:   shared.ScrollbarAlias,
		SelectorAlias:    shared.Alias,
		ItemWidth:        shared.ItemWidth,
		ItemSelected:     shared.ItemSelected,
		IsTrayOpen:       shared.IsTrayOpen,
		ViewportPosition: shared.ViewportPosition,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared DropdownEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewDropdownEntry(existingSelectorEntry ...*DropdownEntryType) DropdownEntryType {
	var dropdownEntry DropdownEntryType
	dropdownEntry.BaseControlType = NewBaseControl()

	if existingSelectorEntry != nil {
		dropdownEntry.BaseControlType = existingSelectorEntry[0].BaseControlType
		dropdownEntry.SelectionEntry = existingSelectorEntry[0].SelectionEntry
		dropdownEntry.ScrollbarAlias = existingSelectorEntry[0].ScrollbarAlias
		dropdownEntry.Alias = existingSelectorEntry[0].Alias
		dropdownEntry.ItemWidth = existingSelectorEntry[0].ItemWidth
		dropdownEntry.IsTrayOpen = existingSelectorEntry[0].IsTrayOpen
		dropdownEntry.ViewportPosition = existingSelectorEntry[0].ViewportPosition
	}
	return dropdownEntry
}

func IsDropdownEntryEqual(sourceDropdownEntry *DropdownEntryType, targetDropdownEntry *DropdownEntryType) bool {
	if sourceDropdownEntry.BaseControlType == targetDropdownEntry.BaseControlType &&
		&sourceDropdownEntry.SelectionEntry == &targetDropdownEntry.SelectionEntry &&
		sourceDropdownEntry.ScrollbarAlias == targetDropdownEntry.ScrollbarAlias &&
		sourceDropdownEntry.Alias == targetDropdownEntry.Alias &&
		sourceDropdownEntry.ItemWidth == targetDropdownEntry.ItemWidth &&
		sourceDropdownEntry.IsTrayOpen == targetDropdownEntry.IsTrayOpen &&
		sourceDropdownEntry.ViewportPosition == targetDropdownEntry.ViewportPosition {
		return true
	}
	return false
}

func GetDropdownAlias(entry *DropdownEntryType) string {
	return entry.Alias
}
