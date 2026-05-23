package types

/*
SelectionEntryType is a structure which represents a list of selectable items. In addition, the following should be noted:

- Used by controls like dropdowns and selectors to manage their options.

Example:
    var selectionEntry types.SelectionEntryType
*/
type SelectionEntryType struct {
	SelectionAlias []string
	SelectionValue []string
}

/*
NewSelectionEntry is a constructor which allows you to create a new selection entry. In addition, the following should be noted:

- Initializes a selection entry with empty arrays for aliases and values.

- Used for managing lists of selectable items in controls like dropdowns and selectors.

- The entry can be populated using the Add method.

Example:
    NewSelectionEntry()
*/
func NewSelectionEntry() SelectionEntryType {
	var selectionEntry SelectionEntryType
	return selectionEntry
}

/*
Add is a method which allows you to add a new selection item to the entry. In addition, the following should be noted:

- Appends a new alias and value pair to the selection entry.

- The alias is used to identify the item, while the value is what's displayed.

- Both arrays (SelectionAlias and SelectionValue) are kept in sync.

Example:
    instance.Add(selectionAlias, selectionValue)
*/
func (shared *SelectionEntryType) Add(selectionAlias string, selectionValue string) {
	shared.SelectionAlias = append(shared.SelectionAlias, selectionAlias)
	shared.SelectionValue = append(shared.SelectionValue, selectionValue)
}

/*
Clear is a method which allows you to remove all items from the selection entry. In addition, the following should be noted:

- Sets both SelectionAlias and SelectionValue arrays to nil.

- Effectively removes all items from the selection.

- The entry can be repopulated using the Add method.

Example:
    instance.Clear()
*/
func (shared *SelectionEntryType) Clear() {
	shared.SelectionAlias = nil
	shared.SelectionValue = nil
}

/*
GetSelectionCount is a method which allows you to get the number of selections currently added.

Example:
    instance.GetSelectionCount()
*/
func (shared *SelectionEntryType) GetSelectionCount() int {
	return len(shared.SelectionAlias)
}
