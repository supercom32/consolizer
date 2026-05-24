package types

/*
FileMenuEntryType is a structure which represents a file menu entry in memory. In addition, the following should be noted:

- The file menu will be drawn at the specified location with the given style.

- Each heading in the menu can have its own dropdown with selectable items.

- The top level headings widths are always dynamic based on how large the heading is.

- The file menu reuses existing selectors for dropdown functionality.

Example:

	var fileMenu types.FileMenuEntryType
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
