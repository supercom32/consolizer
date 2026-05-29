package types

/*
FileMenuEntryType is a structure which represents a file menu entry in memory.
*/
type FileMenuEntryType struct {
	LayerAlias         string
	Alias              string
	StyleEntry         TuiStyleEntryType
	MenuHeadings       []string
	MenuSelections     []SelectionEntryType
	XLocation          int
	YLocation          int
	DynamicWidth       bool
	HeadingWidth       int
	ActiveHeadingIndex int
	IsSubmenuOpen      bool
	IsEnabled          bool
	// Selectors for each menu heading
	SelectorAliases []string
	// Tooltip for the file menu
	TooltipAlias string
}

/*
NewFileMenuEntry is a constructor which creates a new file menu entry. In addition, the following should be noted:

- Initializes a file menu entry with default values.

- Used for managing file menus in the TUI.

- Sets up arrays for selector aliases.

Example:
    NewFileMenuEntry()
*/
func NewFileMenuEntry() FileMenuEntryType {
	var fileMenuEntry FileMenuEntryType
	fileMenuEntry.ActiveHeadingIndex = -1
	fileMenuEntry.IsSubmenuOpen = false
	fileMenuEntry.IsEnabled = true
	fileMenuEntry.DynamicWidth = true
	fileMenuEntry.SelectorAliases = []string{}
	return fileMenuEntry
}
