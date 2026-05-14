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

/*
GetAlias is a method which allows you to retrieve the alias of a dropdown control. In addition, the following
information should be noted:

- Returns the unique identifier for the dropdown.

- This alias is used to reference the dropdown in other operations.

- The alias is set when the dropdown is created.

:return: string

Example:

	instance.GetAlias()
*/
func (shared DropdownEntryType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON is a method which allows you to serialize a dropdown control to JSON. In addition, the following information
should be noted:

- Converts the dropdown's state to a JSON representation.

- Includes the base control properties and dropdown-specific fields.

- Used for saving and loading dropdown configurations.

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
GetEntryAsJsonDump is a method which allows you to get a JSON string representation of a dropdown control. In addition,
the following information should be noted:

- Returns a formatted JSON string of the dropdown's state.

- Useful for debugging and logging purposes.

- Panics if JSON marshaling fails.

:return: string

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
NewDropdownEntry is a constructor which allows you to create a new dropdown control. In addition, the following
information should be noted:

- Initializes a dropdown with default values.

- Can optionally copy properties from an existing dropdown.

- Sets up the base control properties and dropdown-specific fields.

:param existingSelectorEntry: The existingSelectorEntry parameter.

:return: DropdownEntryType

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
IsDropdownEntryEqual is a method which allows you to compare two dropdown controls for equality. In addition, the
following information should be noted:

- Compares all properties of both dropdowns.

- Returns true if all properties match, false otherwise.

- Used for change detection and state synchronization.

:param sourceDropdownEntry: The sourceDropdownEntry parameter.
:param targetDropdownEntry: The targetDropdownEntry parameter.

:return: bool

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
GetDropdownAlias is a method which allows you to retrieve the alias of a dropdown control. In addition, the following
information should be noted:

- Returns the unique identifier for the dropdown.

- This is a convenience method that delegates to GetAlias.

- The alias is used to reference the dropdown in other operations.

:param entry: The entry parameter.

:return: string

Example:

	GetDropdownAlias(entry)
*/
func GetDropdownAlias(entry *DropdownEntryType) string {
	return entry.Alias
}
