package types

import (
	"encoding/json"
)

/*
DropdownEntryType is a structure which represents a dropdown control. In addition, the following should be noted:

- It includes base control properties and selection-specific data.

- It manages the state of the dropdown tray and viewport position.

Example:

	var dropdown types.DropdownEntryType
*/
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

/*
GetAlias is a method which retrieves the alias of a dropdown control. In addition, the following should be noted:

- It returns the unique identifier for the dropdown.

- This alias is used to reference the dropdown in other operations.

- The alias is set when the dropdown is created.

Example:

	instance.GetAlias()
*/
func (shared DropdownEntryType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON is a method which serializes a dropdown control to JSON. In addition, the following should be noted:

- It converts the dropdown's state to a JSON representation.

- It includes the base control properties and dropdown-specific fields.

- It is used for saving and loading dropdown configurations.

Example:

	instance.MarshalJSON()
*/
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

/*
GetEntryAsJsonDump is a method which retrieves a JSON string representation of a dropdown control. In addition, the following should be noted:

- It returns a formatted JSON string of the dropdown's state.

- It is useful for debugging and logging purposes.

- It panics if JSON marshaling fails.

Example:

	instance.GetEntryAsJsonDump()
*/
func (shared DropdownEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewDropdownEntry is a constructor which creates a new dropdown control. In addition, the following should be noted:

- It initializes a dropdown with default values.

- It can optionally copy properties from an existing dropdown.

- It sets up the base control properties and dropdown-specific fields.

Example:

	NewDropdownEntry(existingSelectorEntry)
*/
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

/*
IsDropdownEntryEqual is a method which compares two dropdown controls for equality. In addition, the following should be noted:

- It compares all properties of both dropdowns.

- It returns true if all properties match, false otherwise.

- It is used for change detection and state synchronization.

Example:

	IsDropdownEntryEqual(sourceDropdownEntry, targetDropdownEntry)
*/
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

/*
GetDropdownAlias is a method which retrieves the alias of a dropdown control. In addition, the following should be noted:

- It returns the unique identifier for the dropdown.

- This is a convenience method that delegates to GetAlias.

- The alias is used to reference the dropdown in other operations.

Example:

	GetDropdownAlias(entry)
*/
func GetDropdownAlias(entry *DropdownEntryType) string {
	return entry.Alias
}
