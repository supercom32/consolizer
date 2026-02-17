package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
)

// FileMenuInstanceType represents an instance of a file menu.
type FileMenuInstanceType struct {
	BaseControlInstanceType
}

// fileMenuType is the main struct for managing file menus.
type fileMenuType struct {
	previousButtonState uint
}

// FileMenu is a singleton instance of fileMenuType.
var FileMenu fileMenuType
var FileMenus = memory.NewControlMemoryManager[types.FileMenuEntryType]()

/*
Delete allows you to remove a file menu instance from memory. In addition, the following
information should be noted:

- This method is used to clean up resources when a file menu is no longer needed.
- After deletion, the file menu instance should not be used anymore.
*/
func (shared *FileMenuInstanceType) Delete() *FileMenuInstanceType {
	shared.BaseControlInstanceType.Delete()
	return nil
}

/*
AddToTabIndex allows you to add a file menu to the tab index. In addition, the following
information should be noted:

- This method is used to make the file menu focusable via tab navigation.
*/
func (shared *FileMenuInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeFileMenuHeading)
}

/*
Add allows you to add a new file menu to a layer. In addition, the following
information should be noted:

- The file menu will be drawn at the specified location with the given style.
- Each heading in the menu can have its own dropdown with selectable items.
- The top level headings widths are always dynamic based on how large the heading is.
- When clicking on a heading, the selection includes borders and is as large as the number of items in the menu.
- The dropdown menu appears directly beneath the corresponding header item.
- The dropdown's border aligns neatly under the header and adjusts in width to fit the longest item label.
- The file menu reuses existing selectors for dropdown functionality.
*/
func (shared *fileMenuType) Add(layerAlias string, menuAlias string, styleEntry types.TuiStyleEntryType,
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

	// Create selectors for each menu heading
	for i, selection := range menuSelections {
		// Calculate max width needed for this selector based on the longest dropdown menu item
		maxWidth := 0
		for _, item := range selection.SelectionValue {
			if len(item) > maxWidth {
				maxWidth = len(item)
			}
		}
		maxWidth += 4 // Add padding

		// Calculate position for this selector
		selectorX := xLocation
		// Sum the widths of all headings before this one
		for j := 0; j < i; j++ {
			selectorX += len(menuHeadings[j]) + 2 // Add padding
		}

		// Create a selector for this heading
		selectorAlias := stringformat.GetLastSortedUUID()
		fileMenuEntry.SelectorAliases = append(fileMenuEntry.SelectorAliases, selectorAlias)

		// Add the selector (initially hidden)
		// Position the selector directly beneath the header
		selectorInstance := Selector.Add(layerAlias, selectorAlias, styleEntry, selection,
			selectorX+1, yLocation+2, len(selection.SelectionAlias), maxWidth, 1, 0, 0, false, true)
		selectorInstance.Unselect()
		selectorInstance.SetVisible(false)
	}

	// Store the file menu entry in memory
	FileMenus.Add(layerAlias, menuAlias, &fileMenuEntry)

	// Create and return a file menu instance
	var fileMenuInstance FileMenuInstanceType
	fileMenuInstance.layerAlias = layerAlias
	fileMenuInstance.controlAlias = menuAlias
	fileMenuInstance.controlType = constants.TYPE_FILEMENU
	return fileMenuInstance
}

/*
DeleteFileMenu allows you to delete a file menu from memory. In addition, the following
information should be noted:

- This method is used to clean up resources when a file menu is no longer needed.
- After deletion, the file menu should not be used anymore.
*/
func (shared *fileMenuType) DeleteFileMenu(layerAlias string, menuAlias string) {
	if FileMenus.IsExists(layerAlias, menuAlias) {
		fileMenuEntry := FileMenus.Get(layerAlias, menuAlias)

		// Delete all associated selectors
		for _, selectorAlias := range fileMenuEntry.SelectorAliases {
			Selector.DeleteSelector(layerAlias, selectorAlias)
		}

		// Delete the tooltip
		Tooltip.DeleteTooltip(layerAlias, fileMenuEntry.TooltipAlias)

		// Delete the file menu entry
		FileMenus.Remove(layerAlias, menuAlias)
	}
}

/*
DeleteAllFileMenus allows you to delete all file menus from a layer. In addition, the following
information should be noted:

- This method is used to clean up resources when a layer is no longer needed.
- After deletion, the file menus should not be used anymore.
*/
func (shared *fileMenuType) DeleteAllFileMenus(layerAlias string) {
	fileMenuEntries := FileMenus.GetAllEntries(layerAlias)
	for _, fileMenuEntry := range fileMenuEntries {
		shared.DeleteFileMenu(layerAlias, fileMenuEntry.Alias)
	}
}

/*
drawFileMenusOnLayer allows you to draw all file menus on a layer. In addition, the following
information should be noted:

- This method is called internally by the rendering system.
- It iterates through all file menus and draws them on the specified layer.
*/
func (shared *fileMenuType) drawFileMenusOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	fileMenuEntries := FileMenus.GetAllEntries(layerAlias)
	for _, fileMenuEntry := range fileMenuEntries {
		shared.drawFileMenu(&layerEntry, fileMenuEntry)
	}
}

/*
drawFileMenu allows you to draw a file menu on a layer. In addition, the following
information should be noted:

- This method is called internally by drawFileMenusOnLayer.
- It draws the menu bar and manages the visibility of associated selectors.
*/
func (shared *fileMenuType) drawFileMenu(layerEntry *types.LayerEntryType, fileMenuEntry *types.FileMenuEntryType) {
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
		selectorEntry := Selectors.Get(fileMenuEntry.LayerAlias, selectorAlias)
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
updateKeyboardEvent allows you to update the state of file menus based on keyboard events.
In addition, the following information should be noted:

- This method is called internally by the input handling system.
- It handles keyboard navigation for file menus.
*/
func (shared *fileMenuType) updateKeyboardEvent(keystroke []rune) (bool, bool) {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := false
	isKeystrokeConsumed := false

	// Handle escape key to close open menus
	if keystrokeAsString == "escape" {
		// Get all layer entries
		for _, layerEntry := range Layers.GetAllEntries() {
			layerAlias := layerEntry.LayerAlias
			fileMenuEntries := FileMenus.GetAllEntries(layerAlias)
			for _, fileMenuEntry := range fileMenuEntries {
				if fileMenuEntry.IsSubmenuOpen {
					fileMenuEntry.IsSubmenuOpen = false
					fileMenuEntry.ActiveHeadingIndex = -1
					isScreenUpdateRequired = true
					isKeystrokeConsumed = true
				}
			}
		}
	}

	return isScreenUpdateRequired, isKeystrokeConsumed
}

/*
updateFileMenuStateMouse allows you to update the state of file menus based on mouse events.
In addition, the following information should be noted:

- This method is called internally by the input handling system.
- It handles mouse clicks on menu headings and delegates to selectors for menu items.
*/
func (shared *fileMenuType) updateFileMenuStateMouse() bool {
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
		if cellType == constants.CellTypeFileMenuHeading && FileMenus.IsExists(layerAlias, cellControlAlias) {
			fileMenuEntry := FileMenus.Get(layerAlias, cellControlAlias)

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
closeAllOpenMenus allows you to close all open file menu submenus. In addition, the following
information should be noted:

- This method is called when clicking outside any menu or when selecting a menu item.
- Returns true if any menu was closed, false otherwise.
*/
func (shared *fileMenuType) closeAllOpenMenus() bool {
	menuClosed := false
	// Get all layer entries
	for _, layerEntry := range Layers.GetAllEntries() {
		layerAlias := layerEntry.LayerAlias
		fileMenuEntries := FileMenus.GetAllEntries(layerAlias)
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

func (shared *FileMenuInstanceType) GetSelectedItem() (int, int, string, string) {
	if !FileMenus.IsExists(shared.layerAlias, shared.controlAlias) {
		return -1, -1, "", ""
	}
	fileMenuEntry := FileMenus.Get(shared.layerAlias, shared.controlAlias)
	// Iterate through selectors (one per heading)
	for headingIndex, selectorAlias := range fileMenuEntry.SelectorAliases {
		selectorEntry := Selectors.Get(shared.layerAlias, selectorAlias)
		if selectorEntry.ItemSelected >= 0 {
			// Store the values to be returned
			itemIndex := selectorEntry.ItemSelected
			itemAlias := selectorEntry.SelectionEntry.SelectionAlias[itemIndex]
			itemValue := selectorEntry.SelectionEntry.SelectionValue[itemIndex]
			shared.Unselect()
			return headingIndex, itemIndex, itemAlias, itemValue
		}
	}
	// Nothing selected
	return -1, -1, "", ""
}

/*
IsFileMenuOpen allows you to check if a file menu is currently open.
In addition, the following information should be noted:

- This method returns true if the file menu's dropdown is visible,
  and false otherwise.
- If the file menu instance does not exist, it will return false.
*/
func (shared *FileMenuInstanceType) IsFileMenuOpen() bool {
	if !FileMenus.IsExists(shared.layerAlias, shared.controlAlias) {
		return false
	}
	fileMenuEntry := FileMenus.Get(shared.layerAlias, shared.controlAlias)
	return fileMenuEntry.IsSubmenuOpen
}

/*
Unselect allows you to clear the current selection for a file menu. In addition,
the following information should be noted:

  - This method iterates through all submenus (selectors) of the file menu and
    unselects any selected item.
  - If the file menu does not exist, no operation occurs.
*/
func (shared *FileMenuInstanceType) Unselect() {
	if !FileMenus.IsExists(shared.layerAlias, shared.controlAlias) {
		return
	}
	fileMenuEntry := FileMenus.Get(shared.layerAlias, shared.controlAlias)
	// Iterate through selectors (one per heading) and unselect them.
	for _, selectorAlias := range fileMenuEntry.SelectorAliases {
		selectorEntry := Selectors.Get(shared.layerAlias, selectorAlias)
		selectorEntry.ItemSelected = constants.SELECTED_NONE
	}
}
