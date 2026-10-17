package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
)

/*
FileMenuInstanceType is a structure which represents an instance of a file menu.
*/
type FileMenuInstanceType struct {
	BaseControlInstanceType
}

/*
fileMenuType is a structure which is the main struct for managing file menus.
*/
type fileMenuType struct {
	previousButtonState uint
}

/*
FileMenu is a structure which represents a singleton instance of fileMenuType.
*/
var FileMenu fileMenuType
var fileMenuRegistry = memory.NewControlMemoryManager[types.FileMenuEntryType]()

/*
Delete is a method which allows you to remove a file menu instance from memory. In addition, the following should be
noted:

- This method is used to clean up resources when a file menu is no longer needed.

- After deletion, the file menu instance should not be used anymore.

Example:
    fileMenu.Delete()
*/
func (shared *FileMenuInstanceType) Delete() *FileMenuInstanceType {
	defer lockUIState()()
	shared.BaseControlInstanceType.delete()
	return nil
}

/*
Add is a method which allows you to add a new file menu to a layer. In addition, the following should be noted:

- The file menu will be drawn at the specified location with the given style.

- Each heading in the menu can have its own dropdown with selectable items.

- The top level headings widths are always dynamic based on how large the heading is.

- When clicking on a heading, the selection includes borders and is as large as the number of items in the menu.

- The dropdown menu appears directly beneath the corresponding header item.

- The dropdown's border aligns neatly under the header and adjusts in width to fit the longest item label.

- The file menu reuses existing selectors for dropdown functionality.

Example:
    FileMenu.Add("main", "fileMenu", style, headings, selections, 0, 0, true)
*/
func (shared *fileMenuType) Add(layerAlias string, menuAlias string, styleEntry types.TuiStyleEntryType,
	menuHeadings []string, menuSelections []types.SelectionEntryType, xLocation int, yLocation int,
	isEnabled bool) FileMenuInstanceType {
	defer lockUIState()()
	return shared.add(layerAlias, menuAlias, styleEntry, menuHeadings, menuSelections, xLocation, yLocation, isEnabled)
}

/*
add is a method which allows you to do exactly what Add does, for a caller that already holds uiStateMutex, such as
another entry point of the package that has taken the lock itself.

Example:

	FileMenu.add(layerAlias, menuAlias, styleEntry, menuHeadings, menuSelections, xLocation, yLocation, isEnabled)
*/
func (shared *fileMenuType) add(layerAlias string, menuAlias string, styleEntry types.TuiStyleEntryType,
	menuHeadings []string, menuSelections []types.SelectionEntryType, xLocation int, yLocation int,
	isEnabled bool) FileMenuInstanceType {

	fileMenuEntry := types.NewFileMenuEntry()
	fileMenuEntry.LayerAlias = layerAlias
	fileMenuEntry.Alias = menuAlias
	fileMenuEntry.StyleEntry = styleEntry
	fileMenuEntry.MenuHeadings = menuHeadings
	fileMenuEntry.MenuSelections = menuSelections
	fileMenuEntry.XLocation = xLocation
	fileMenuEntry.YLocation = yLocation
	fileMenuEntry.DynamicWidth = true
	fileMenuEntry.ActiveHeadingIndex = -1
	fileMenuEntry.IsSubmenuOpen = false
	fileMenuEntry.IsEnabled = isEnabled

	shared.createSubmenus(&fileMenuEntry)

	// Store the file menu entry in memory
	fileMenuRegistry.Add(layerAlias, menuAlias, &fileMenuEntry)

	// Create and return a file menu instance
	var fileMenuInstance FileMenuInstanceType
	fileMenuInstance.layerAlias = layerAlias
	fileMenuInstance.controlAlias = menuAlias
	fileMenuInstance.controlType = constants.TYPE_FILEMENU
	return fileMenuInstance
}

/*
createSubmenus is a method which allows you to create the hidden selector that serves as the submenu of each of a file
menu's headings, from its headings and selections, and record their aliases in the entry. Each submenu is placed
directly under its heading and made wide enough for its longest item. The caller must hold uiStateMutex, and the entry
must have exactly one selection for each heading.

Example:

	FileMenu.createSubmenus(fileMenuEntry)
*/
func (shared *fileMenuType) createSubmenus(fileMenuEntry *types.FileMenuEntryType) {
	fileMenuEntry.SelectorAliases = nil
	for i, selection := range fileMenuEntry.MenuSelections {
		// Calculate max width needed for this selector based on the longest dropdown menu item
		maxWidth := 0
		for _, item := range selection.SelectionValue {
			if len(item) > maxWidth {
				maxWidth = len(item)
			}
		}
		maxWidth += 4 // Add padding

		// Sum the widths of all headings before this one to place the selector under its heading
		selectorX := fileMenuEntry.XLocation
		for j := 0; j < i; j++ {
			selectorX += len(fileMenuEntry.MenuHeadings[j]) + 2 // Add padding
		}

		selectorAlias := stringformat.GetLastSortedUUID()
		fileMenuEntry.SelectorAliases = append(fileMenuEntry.SelectorAliases, selectorAlias)
		// The selector starts hidden, directly beneath the heading.
		selectorInstance := Selector.add(fileMenuEntry.LayerAlias, selectorAlias, fileMenuEntry.StyleEntry, selection,
			selectorX+1, fileMenuEntry.YLocation+2, len(selection.SelectionAlias), maxWidth, 1, 0, 0, false, true)
		selectorInstance.unselect()
		selectorInstance.setVisible(false)
	}
}

/*
SetMenus is a method which allows you to replace a file menu's headings and the items under each of them. Any open
submenu is closed and any pending selection is discarded. An error is returned, and the file menu is left unchanged,
if the number of headings and the number of selections differ. If the file menu does not exist, nothing happens and
no error is returned.

Example:

	err := fileMenu.SetMenus([]string{"File", "Edit"}, []types.SelectionEntryType{fileItems, editItems})
*/
func (shared *FileMenuInstanceType) SetMenus(menuHeadings []string, menuSelections []types.SelectionEntryType) error {
	defer lockUIState()()
	if len(menuHeadings) != len(menuSelections) {
		return fmt.Errorf("a file menu needs one selection for each heading, but %d headings and %d selections were given", len(menuHeadings), len(menuSelections))
	}
	fileMenuEntry, isFound := fileMenuRegistry.Lookup(shared.layerAlias, shared.controlAlias)
	if !isFound {
		return nil
	}
	for _, selectorAlias := range fileMenuEntry.SelectorAliases {
		Selector.delete(shared.layerAlias, selectorAlias)
	}
	fileMenuEntry.MenuHeadings = append([]string(nil), menuHeadings...)
	fileMenuEntry.MenuSelections = append([]types.SelectionEntryType(nil), menuSelections...)
	fileMenuEntry.ActiveHeadingIndex = -1
	fileMenuEntry.IsSubmenuOpen = false
	FileMenu.createSubmenus(fileMenuEntry)
	return nil
}

/*
Open is a method which allows you to open the submenu under one of a file menu's headings, given by its zero based
index, as clicking the heading does. Any other open submenu is closed first. A heading index that does not exist, a
disabled file menu, or a file menu that does not exist leaves everything unchanged.

Example:

	fileMenu.Open(0)
*/
func (shared *FileMenuInstanceType) Open(headingIndex int) {
	defer lockUIState()()
	fileMenuEntry, isFound := fileMenuRegistry.Lookup(shared.layerAlias, shared.controlAlias)
	if !isFound || !fileMenuEntry.IsEnabled || headingIndex < 0 || headingIndex >= len(fileMenuEntry.SelectorAliases) {
		return
	}
	FileMenu.closeAllOpenMenus()
	fileMenuEntry.ActiveHeadingIndex = headingIndex
	fileMenuEntry.IsSubmenuOpen = true
}

/*
Close is a method which allows you to close a file menu's open submenu, if it has one. If the file menu does not
exist, no operation occurs.

Example:

	fileMenu.Close()
*/
func (shared *FileMenuInstanceType) Close() {
	defer lockUIState()()
	if fileMenuEntry, isFound := fileMenuRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		fileMenuEntry.IsSubmenuOpen = false
		fileMenuEntry.ActiveHeadingIndex = -1
	}
}

/*
SetEnabled is a method which allows you to enable or disable a file menu. A disabled file menu is not drawn, cannot be
opened, and cannot take focus, so disabling it closes its open submenu and moves focus away from it if it had focus.
If the file menu does not exist, no operation occurs.

Example:

	fileMenu.SetEnabled(false)
*/
func (shared *FileMenuInstanceType) SetEnabled(isEnabled bool) *FileMenuInstanceType {
	defer lockUIState()()
	if fileMenuEntry, isFound := fileMenuRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		fileMenuEntry.IsEnabled = isEnabled
		if !isEnabled {
			fileMenuEntry.IsSubmenuOpen = false
			fileMenuEntry.ActiveHeadingIndex = -1
		}
	}
	reconcileFocus()
	return shared
}

/*
IsEnabled is a method which allows you to check whether a file menu is enabled. If the file menu does not exist, false
is returned.

Example:

	isEnabled := fileMenu.IsEnabled()
*/
func (shared *FileMenuInstanceType) IsEnabled() bool {
	defer lockUIState()()
	if fileMenuEntry, isFound := fileMenuRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return fileMenuEntry.IsEnabled
	}
	return false
}

/*
Delete is a method which allows you to delete a file menu from memory. In addition, the following should be
noted:

- This method is used to clean up resources when a file menu is no longer needed.

- After deletion, the file menu should not be used anymore.

Example:
    FileMenu.Delete("main", "fileMenu")
*/
func (shared *fileMenuType) Delete(layerAlias string, menuAlias string) {
	defer lockUIState()()
	shared.delete(layerAlias, menuAlias)
}

/*
delete is a method which allows you to do exactly what Delete does, for a caller that already holds uiStateMutex, such
as another entry point of the package that has taken the lock itself.

Example:

	FileMenu.delete(layerAlias, menuAlias)
*/
func (shared *fileMenuType) delete(layerAlias string, menuAlias string) {
	if fileMenuEntry, isFound := fileMenuRegistry.Lookup(layerAlias, menuAlias); isFound {
		// Delete all associated selectors
		for _, selectorAlias := range fileMenuEntry.SelectorAliases {
			Selector.delete(layerAlias, selectorAlias)
		}

		// Delete the tooltip
		Tooltip.delete(layerAlias, fileMenuEntry.TooltipAlias)

		// Delete the file menu entry
		fileMenuRegistry.Remove(layerAlias, menuAlias)
		clearStaleControlReferences(layerAlias, menuAlias, constants.CellTypeFileMenuHeading)
	}
}

/*
DeleteAll is a method which allows you to delete all file menus from a layer. In addition, the following should
be noted:

- This method is used to clean up resources when a layer is no longer needed.

- After deletion, the file menus should not be used anymore.

Example:
    FileMenu.DeleteAll("main")
*/
func (shared *fileMenuType) DeleteAll(layerAlias string) {
	defer lockUIState()()
	fileMenuEntries := fileMenuRegistry.GetAllEntries(layerAlias)
	for _, fileMenuEntry := range fileMenuEntries {
		shared.delete(layerAlias, fileMenuEntry.Alias)
	}
}

/*
drawOnLayer is a method which allows you to draw all file menus on a layer. In addition, the following should
be noted:

- This method is called internally by the rendering system.

- It iterates through all file menus and draws them on the specified layer.

Example:
    FileMenu.drawOnLayer(layer)
*/
func (shared *fileMenuType) drawOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	fileMenuEntries := fileMenuRegistry.GetAllEntries(layerAlias)
	for _, fileMenuEntry := range fileMenuEntries {
		shared.draw(&layerEntry, fileMenuEntry)
	}
}

/*
draw is a method which allows you to draw a file menu on a layer. In addition, the following should be noted:

- This method is called internally by drawFileMenusOnLayer.

- It draws the menu bar and manages the visibility of associated selectors.

Example:
    FileMenu.draw(layer, entry)
*/
func (shared *fileMenuType) draw(layerEntry *types.LayerEntryType, fileMenuEntry *types.FileMenuEntryType) {
	if !fileMenuEntry.IsEnabled {
		return
	}

	// Draw the menu bar
	currentX := fileMenuEntry.XLocation
	for index, heading := range fileMenuEntry.MenuHeadings {
		// Determine the width of this heading
		headingWidth := len(heading) + 2 // Add padding

		// Set up attributes for drawing
		attributeEntry := types.NewAttributeEntry()
		attributeEntry.ForegroundColor = fileMenuEntry.StyleEntry.FileMenu.ForegroundColor
		attributeEntry.BackgroundColor = fileMenuEntry.StyleEntry.FileMenu.BackgroundColor

		// Highlight active heading if its submenu is open
		if index == fileMenuEntry.ActiveHeadingIndex && fileMenuEntry.IsSubmenuOpen {
			attributeEntry.ForegroundColor = fileMenuEntry.StyleEntry.FileMenu.HighlightForegroundColor
			attributeEntry.BackgroundColor = fileMenuEntry.StyleEntry.FileMenu.HighlightBackgroundColor
		}

		// Set cell type for mouse interaction
		attributeEntry.CellType = constants.CellTypeFileMenuHeading
		attributeEntry.CellControlAlias = fileMenuEntry.Alias
		attributeEntry.CellControlId = index

		// Draw the heading with padding
		paddedHeading := " " + heading + " "
		if len(paddedHeading) < headingWidth {
			paddedHeading = paddedHeading + stringformat.GetFilledString(headingWidth-len(paddedHeading), " ")
		} else if len(paddedHeading) > headingWidth {
			paddedHeading = paddedHeading[:headingWidth]
		}
		printLayer(layerEntry, attributeEntry, currentX, fileMenuEntry.YLocation, []rune(paddedHeading))

		// Move to the next heading position
		currentX += headingWidth
	}

	// Manage selector visibility based on active heading
	for i, selectorAlias := range fileMenuEntry.SelectorAliases {
		selectorEntry := selectorRegistry.Get(fileMenuEntry.LayerAlias, selectorAlias)
		if i == fileMenuEntry.ActiveHeadingIndex && fileMenuEntry.IsSubmenuOpen {
			// Show the selector for the active heading
			selectorEntry.IsVisible = true
		} else {
			// Hide selectors for inactive headings
			selectorEntry.IsVisible = false
		}
	}
}

/*
updateKeyboardEvent is a method which allows you to update the state of file menus based on keyboard events. In
addition, the following should be noted:

- This method is called internally by the input handling system.

- It handles keyboard navigation for file menus.

- Esc closes every open file menu and is consumed. When no file menu is open, Esc is left unconsumed so the
  application can treat it as a Back or Cancel action.

Example:
    updateRequired, consumed := FileMenu.updateKeyboardEvent(keystroke)
*/
func (shared *fileMenuType) updateKeyboardEvent(keystroke []rune) (bool, bool) {
	// tcell reports the Escape key as "Esc", which is lowercased to "esc" before reaching the controls.
	if string(keystroke) == "esc" && shared.closeAllOpenMenus() {
		return true, true
	}
	return false, false
}

/*
updateStateMouse is a method which allows you to update the state of file menus based on mouse events. In
addition, the following should be noted:

- This method is called internally by the input handling system.

- It handles mouse clicks on menu headings and delegates to selectors for menu items.

Example:
    updateRequired := FileMenu.updateStateMouse()
*/
func (shared *fileMenuType) updateStateMouse() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	cellControlAlias := characterEntry.AttributeEntry.CellControlAlias
	cellControlId := characterEntry.AttributeEntry.CellControlId
	cellType := characterEntry.AttributeEntry.CellType

	// Detect a new click (button state changed from released to pressed)
	isNewClick := buttonPressed != 0 && shared.previousButtonState == 0

	// Update the previous button state for the next call
	defer func() { shared.previousButtonState = buttonPressed }()

	// Only process if a button is pressed and it's a new click (not a continued press)
	if isNewClick {
		// Check if the mouse clicked on a file menu heading
		if fileMenuEntry, isFound := fileMenuRegistry.Lookup(layerAlias, cellControlAlias); cellType == constants.CellTypeFileMenuHeading && isFound {
			// If clicking on the active heading, close the submenu
			if fileMenuEntry.ActiveHeadingIndex == cellControlId && fileMenuEntry.IsSubmenuOpen {
				fileMenuEntry.IsSubmenuOpen = false
				fileMenuEntry.ActiveHeadingIndex = -1
			} else {
				// Close any other open menus first
				shared.closeAllOpenMenus()

				// Open the submenu for this heading
				fileMenuEntry.ActiveHeadingIndex = cellControlId
				fileMenuEntry.IsSubmenuOpen = true
			}
			isUpdateRequired = true
		} else if cellType == constants.CellTypeSelectorItem {
			// If a selector item was clicked, check if it belongs to a file menu
			// This is handled by the selector's own event handling

			// After a selection, close all menus
			shared.closeAllOpenMenus()
			isUpdateRequired = true
		} else {
			// If clicking outside any menu, close all open menus
			if shared.closeAllOpenMenus() {
				isUpdateRequired = true
			}
		}
	}

	return isUpdateRequired
}

/*
closeAllOpenMenus is a method which allows you to close all open file menu submenus. In addition, the following should
be noted:

- This method is called when clicking outside any menu or when selecting a menu item.

- Returns true if any menu was closed, false otherwise.

Example:
    FileMenu.closeAllOpenMenus()
*/
func (shared *fileMenuType) closeAllOpenMenus() bool {
	menuClosed := false
	// Get all layer entries
	for _, layerEntry := range layerRegistry.GetAllEntries() {
		layerAlias := layerEntry.LayerAlias
		fileMenuEntries := fileMenuRegistry.GetAllEntries(layerAlias)
		for _, fileMenuEntry := range fileMenuEntries {
			if fileMenuEntry.IsSubmenuOpen {
				fileMenuEntry.IsSubmenuOpen = false
				fileMenuEntry.ActiveHeadingIndex = -1
				menuClosed = true
			}
		}
	}
	return menuClosed
}

/*
GetSelectedItem is a method which allows you to retrieve the currently selected item from the file menu, as the index
of its heading, its index within that heading's submenu, its alias, and its display value. Reading the selection clears
it. If nothing is selected, or the file menu does not exist, minus one, minus one and two empty strings are returned.
In addition, the following should be noted:

  - A recorded selection that no longer fits its submenu's items is ignored rather than used to index them, and the
    submenu of a heading whose selector no longer exists is skipped.

Example:

	heading, item, alias, value := fileMenu.GetSelectedItem()
*/
func (shared *FileMenuInstanceType) GetSelectedItem() (int, int, string, string) {
	defer lockUIState()()
	fileMenuEntry, isFound := fileMenuRegistry.Lookup(shared.layerAlias, shared.controlAlias)
	if !isFound {
		return -1, -1, "", ""
	}
	// Iterate through selectors (one per heading)
	for headingIndex, selectorAlias := range fileMenuEntry.SelectorAliases {
		selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, selectorAlias)
		if !isFound {
			continue
		}
		if selectorEntry.ItemSelected >= 0 && selectorEntry.ItemSelected < Selector.getItemCount(selectorEntry) {
			// Store the values to be returned
			itemIndex := selectorEntry.ItemSelected
			itemAlias := selectorEntry.SelectionEntry.SelectionAlias[itemIndex]
			itemValue := selectorEntry.SelectionEntry.SelectionValue[itemIndex]
			shared.unselectLocked(fileMenuEntry)
			return headingIndex, itemIndex, itemAlias, itemValue
		}
	}
	// Nothing selected
	return -1, -1, "", ""
}

/*
IsOpen is a method which allows you to check if a file menu is currently open. In addition, the following should
be noted:

- This method returns true if the file menu's dropdown is visible, and false otherwise.

- If the file menu instance does not exist, it will return false.

Example:
    if fileMenu.IsOpen() {
        fmt.Println("Menu is open")
    }
*/
func (shared *FileMenuInstanceType) IsOpen() bool {
	defer lockUIState()()
	fileMenuEntry, isFound := fileMenuRegistry.Lookup(shared.layerAlias, shared.controlAlias)
	if !isFound {
		return false
	}
	return fileMenuEntry.IsSubmenuOpen
}

/*
Unselect is a method which allows you to clear the current selection of a file menu, by clearing the selected item of
every one of its submenus. If the file menu does not exist, no operation occurs, and a submenu whose selector no longer
exists is skipped.

Example:

	fileMenu.Unselect()
*/
func (shared *FileMenuInstanceType) Unselect() {
	defer lockUIState()()
	fileMenuEntry, isFound := fileMenuRegistry.Lookup(shared.layerAlias, shared.controlAlias)
	if !isFound {
		return
	}
	shared.unselectLocked(fileMenuEntry)
}

/*
unselectLocked is a method which allows you to clear the selected item of every submenu of a file menu, as Unselect
does, for a caller that already holds uiStateMutex. A submenu whose selector no longer exists is skipped.

Example:

	fileMenu.unselectLocked(fileMenuEntry)
*/
func (shared *FileMenuInstanceType) unselectLocked(fileMenuEntry *types.FileMenuEntryType) {
	for _, selectorAlias := range fileMenuEntry.SelectorAliases {
		if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, selectorAlias); isFound {
			selectorEntry.ItemSelected = constants.SELECTED_NONE
		}
	}
}
