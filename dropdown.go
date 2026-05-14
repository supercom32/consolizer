package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
)

type DropdownInstanceType struct {
	BaseControlInstanceType
}

type dropdownType struct{}

/*
updateKeyboardEvent is a method which allows you to update the state of all dropdowns according to the current keyboard
event. In addition, the following should be noted:

- Handles Enter key to open/close the dropdown.

- Handles Up/Down keys to navigate through dropdown options when open.

- Returns true if the screen needs to be updated due to state changes.

:param keystroke: The keyboard event as a slice of runes.

:return: Whether a screen update is required and whether the keystroke was consumed.

Example:

	isUpdate, isConsumed := Dropdown.updateKeyboardEvent(keystroke)
*/
func (shared *dropdownType) updateKeyboardEvent(keystroke []rune) (bool, bool) {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := false
	isKeystrokeConsumed := false
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType

	// Only process if a dropdown is focused
	if focusedControlType != constants.CellTypeDropdown || !Dropdowns.IsExists(focusedLayerAlias, focusedControlAlias) {
		return isScreenUpdateRequired, isKeystrokeConsumed
	}

	dropdownEntry := Dropdowns.Get(focusedLayerAlias, focusedControlAlias)

	// If dropdown is open but focus is on the dropdown itself (not the selector),
	// move focus to the selector for keyboard navigation
	if dropdownEntry.IsTrayOpen && focusedControlType == constants.CellTypeDropdown {
		isScreenUpdateRequired, isKeystrokeConsumed = Selector.updateKeyboardEventForSelector(focusedLayerAlias, dropdownEntry.SelectorAlias, keystroke)
	}

	// Handle Enter key to open/close dropdown
	if keystrokeAsString == "enter" || keystrokeAsString == "esc" {
		if dropdownEntry.IsTrayOpen {
			// Close dropdown and apply selection
			selectorEntry := Selectors.Get(focusedLayerAlias, dropdownEntry.SelectorAlias)
			scrollBarEntry := ScrollBars.Get(focusedLayerAlias, dropdownEntry.ScrollbarAlias)
			scrollBarEntry.ScrollValue = selectorEntry.ItemSelected
			// Update selected item if changed
			if dropdownEntry.ItemSelected != selectorEntry.ItemSelected {
				dropdownEntry.ItemSelected = selectorEntry.ItemSelected
			}

			// Hide dropdown components
			selectorEntry.IsVisible = false
			scrollBarEntry.IsVisible = false
			dropdownEntry.IsTrayOpen = false

			// Reset focus to the dropdown itself
			setFocusedControl(focusedLayerAlias, focusedControlAlias, constants.CellTypeDropdown)
			eventStateMemory.stateId = constants.EventStateNone
		} else {
			// Open dropdown
			shared.closeAllOpen() // Close any other open dropdowns first
			dropdownEntry.IsTrayOpen = true

			// Show dropdown components
			selectorEntry := Selectors.Get(focusedLayerAlias, dropdownEntry.SelectorAlias)
			selectorEntry.IsVisible = true
			selectorEntry.ItemHighlighted = dropdownEntry.ItemSelected // Highlight current selection

			// Set focus to the selector for keyboard navigation
			//setFocusedControl(focusedLayerAlias, dropdownEntry.SelectorAlias, constants.CellTypeSelectorItem)

			// Show scrollbar if needed
			scrollBarEntry := ScrollBars.Get(focusedLayerAlias, dropdownEntry.ScrollbarAlias)
			if scrollBarEntry.IsEnabled {
				scrollBarEntry.IsVisible = true
			}
		}
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true
	}

	// Handle Escape key to close dropdown without changing selection
	if keystrokeAsString == "escape" && dropdownEntry.IsTrayOpen {
		// Close dropdown without applying selection
		selectorEntry := Selectors.Get(focusedLayerAlias, dropdownEntry.SelectorAlias)
		scrollBarEntry := ScrollBars.Get(focusedLayerAlias, dropdownEntry.ScrollbarAlias)

		// Hide dropdown components
		selectorEntry.IsVisible = false
		scrollBarEntry.IsVisible = false
		dropdownEntry.IsTrayOpen = false

		// Reset focus to the dropdown itself
		setFocusedControl(focusedLayerAlias, focusedControlAlias, constants.CellTypeDropdown)
		eventStateMemory.stateId = constants.EventStateNone
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true
	}

	return isScreenUpdateRequired, isKeystrokeConsumed
}

var Dropdown dropdownType
var Dropdowns = memory.NewControlMemoryManager[types.DropdownEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
Delete is a method which allows you to remove a dropdown from a text layer. In addition, the following should be noted:

- If you attempt to delete a dropdown which does not exist, then the request will simply be ignored.

- All memory associated with the dropdown will be freed.

:return: A pointer to the dropdown instance (always nil).

Example:

	dropdown.Delete()
*/
func (shared *DropdownInstanceType) Delete() *DropdownInstanceType {
	shared.BaseControlInstanceType.Delete()
	return nil
}

/*
AddToTabIndex is a method which allows you to add a dropdown to the tab index. This enables keyboard navigation between
controls using the tab key. In addition, the following should be noted:

- The dropdown will be added to the tab order based on the order in which it was created.

- The tab index is used to determine which control receives focus when the tab key is pressed.

Example:

	dropdown.AddToTabIndex()
*/
func (shared *DropdownInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeDropdown)
}

/*
GetValue is a method which allows you to retrieve the currently selected value from a dropdown. In addition, the
following should be noted:

- Returns the display value of the currently selected item.

- If the dropdown does not exist, returns an empty string.

:return: The display value of the selected item.

Example:

	val := dropdown.GetValue()
*/
func (shared *DropdownInstanceType) GetValue() string {
	dropdownEntry := Dropdowns.Get(shared.layerAlias, shared.controlAlias)
	if len(dropdownEntry.SelectionEntry.SelectionValue) != 0 &&
		dropdownEntry.ItemSelected >= 0 && dropdownEntry.ItemSelected < len(dropdownEntry.SelectionEntry.SelectionValue) {
		return dropdownEntry.SelectionEntry.SelectionValue[dropdownEntry.ItemSelected]
	}
	return ""
}

/*
GetAlias is a method which allows you to retrieve the currently selected alias from a dropdown. In addition, the
following should be noted:

- Returns the internal alias of the currently selected item.

- If the dropdown does not exist, returns an empty string.

- The alias is typically used for programmatic access to the selection.

:return: The alias of the selected item.

Example:

	alias := dropdown.GetAlias()
*/
func (shared *DropdownInstanceType) GetAlias() string {
	dropdownEntry := Dropdowns.Get(shared.layerAlias, shared.controlAlias)
	if len(dropdownEntry.SelectionEntry.SelectionAlias) != 0 &&
		dropdownEntry.ItemSelected >= 0 && dropdownEntry.ItemSelected < len(dropdownEntry.SelectionEntry.SelectionAlias) {
		return dropdownEntry.SelectionEntry.SelectionAlias[dropdownEntry.ItemSelected]
	}
	return ""
}

/*
GetSelectedItemIndex is a method which allows you to retrieve the index of the currently selected item in the dropdown.

:return: The index of the selected item.

Example:

	index := dropdown.GetSelectedItemIndex()
*/
func (shared *DropdownInstanceType) GetSelectedItemIndex() int {
	dropdownEntry := Dropdowns.Get(shared.layerAlias, shared.controlAlias)
	return dropdownEntry.ItemSelected
}

/*
SetSelectedItemIndex is a method which allows you to set the currently selected item in the dropdown by its index.

:param itemIndex: The index of the item to select.

Example:

	dropdown.SetSelectedItemIndex(2)
*/
func (shared *DropdownInstanceType) SetSelectedItemIndex(itemIndex int) {
	dropdownEntry := Dropdowns.Get(shared.layerAlias, shared.controlAlias)
	if itemIndex < len(dropdownEntry.SelectionEntry.SelectionValue) {
		dropdownEntry.ItemSelected = itemIndex
		selectorEntry := Selectors.Get(shared.layerAlias, dropdownEntry.SelectorAlias)
		selectorEntry.ItemSelected = itemIndex
	}
}

/*
SetSelectionEntry is a method which allows you to overwrite the selection entry for a dropdown. In addition, the
following should be noted:

- This replaces the list of items that can be selected in the dropdown.

- The currently selected item index will be reset to -1.

- The associated selector and scrollbar will be updated to reflect the new items.

:param selectionEntry: The new selection entry containing the items to be displayed in the dropdown.

Example:

	dropdown.SetSelectionEntry(newSelection)
*/
func (shared *DropdownInstanceType) SetSelectionEntry(selectionEntry types.SelectionEntryType) {
	dropdownEntry := Dropdowns.Get(shared.layerAlias, shared.controlAlias)
	dropdownEntry.SelectionEntry = selectionEntry
	dropdownEntry.ItemSelected = -1

	selectorEntry := Selectors.Get(shared.layerAlias, dropdownEntry.SelectorAlias)
	selectorEntry.SelectionEntry = selectionEntry
	selectorEntry.ItemSelected = -1

	scrollBarEntry := ScrollBars.Get(shared.layerAlias, dropdownEntry.ScrollbarAlias)
	scrollBarEntry.MaxScrollValue = len(selectionEntry.SelectionValue) - selectorEntry.Height
}

/*
Add is a method which allows you to create a new dropdown control on a text layer. In addition, the following should be
noted:

- The dropdown consists of a main control and an associated selector for the dropdown tray.

- A scrollbar is automatically added if the number of items exceeds the selector height.

- The dropdown tray is initially hidden and only shown when the dropdown is clicked.

- The default selected item can be specified when creating the dropdown.

:param layerAlias: The alias of the layer to add the dropdown to.
:param dropdownAlias: The unique alias for the new dropdown.
:param styleEntry: The TUI style entry to use for rendering.
:param selectionEntry: The selection entry containing the items for the dropdown.
:param xLocation: The X coordinate for the dropdown.
:param yLocation: The Y coordinate for the dropdown.
:param selectorHeight: The height of the dropdown tray when open.
:param itemWidth: The width of the dropdown items.
:param defaultItemSelected: The index of the item to be selected by default.

:return: A new dropdown instance.

Example:

	dropdown := Dropdown.Add("layer1", "myDropdown", style, items, 10, 10, 5, 15, 0)
*/
func (shared *dropdownType) Add(layerAlias string, dropdownAlias string, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, defaultItemSelected int) DropdownInstanceType {
	// TODO: AddLayer validation to the default item selected.
	newDropdownEntry := types.NewDropdownEntry()
	newDropdownEntry.Alias = dropdownAlias
	newDropdownEntry.StyleEntry = styleEntry
	newDropdownEntry.SelectionEntry = selectionEntry
	newDropdownEntry.XLocation = xLocation
	newDropdownEntry.YLocation = yLocation
	newDropdownEntry.ItemWidth = itemWidth
	newDropdownEntry.ItemSelected = defaultItemSelected
	newDropdownEntry.TooltipAlias = stringformat.GetLastSortedUUID()

	// Use the ControlMemoryManager to add the dropdown entry
	Dropdowns.Add(layerAlias, dropdownAlias, &newDropdownEntry)

	dropdownEntry := Dropdowns.Get(layerAlias, dropdownAlias)
	dropdownEntry.ScrollbarAlias = stringformat.GetLastSortedUUID()
	// Here we add +2 to x to account for the scroll bar being outside the Selector border on ether side. Also, we
	// minus the scroll bar max selection size by the height of the Selector, so we don't scroll over values
	// which do not change viewport.
	selectorWidth := itemWidth
	if len(selectionEntry.SelectionValue) <= selectorHeight {
		selectorWidth = selectorWidth + 1
	}
	dropdownEntry.SelectorAlias = stringformat.GetLastSortedUUID()
	// Here we add +1 to x and y to account for borders around the selection.
	Selector.Add(layerAlias, dropdownEntry.SelectorAlias, styleEntry, selectionEntry, xLocation+1, yLocation+1, selectorHeight, selectorWidth, 1, 0, 0, false, true)
	selectorEntry := Selectors.Get(layerAlias, dropdownEntry.SelectorAlias)
	selectorEntry.IsVisible = false
	dropdownEntry.ScrollbarAlias = selectorEntry.ScrollbarAlias
	scrollBarEntry := ScrollBars.Get(layerAlias, dropdownEntry.ScrollbarAlias)
	scrollBarEntry.IsVisible = false
	if len(selectionEntry.SelectionValue) <= selectorHeight {
		scrollBarEntry.IsEnabled = false
	}

	// Create associated tooltip (always created but disabled by default)
	tooltipInstance := Tooltip.Add(layerAlias, dropdownEntry.TooltipAlias, "", styleEntry,
		dropdownEntry.XLocation, dropdownEntry.YLocation,
		dropdownEntry.ItemWidth, 1,
		dropdownEntry.XLocation, dropdownEntry.YLocation+2,
		dropdownEntry.ItemWidth, 3,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.SetEnabled(false)
	tooltipInstance.setParentControlAlias(dropdownAlias)
	var dropdownInstance DropdownInstanceType
	dropdownInstance.layerAlias = layerAlias
	dropdownInstance.controlAlias = dropdownAlias
	dropdownInstance.controlType = constants.TYPE_DROPDOWN
	return dropdownInstance
}

/*
Delete is a method which allows you to remove a dropdown from a text layer. In addition, the following should be
noted:

- If you attempt to delete a dropdown which does not exist, then the request will simply be ignored.

- All memory associated with the dropdown will be freed.

:param layerAlias: The alias of the layer containing the dropdown.
:param dropdownAlias: The alias of the dropdown to delete.

Example:

	Dropdown.Delete("layer1", "myDropdown")
*/
func (shared *dropdownType) Delete(layerAlias string, dropdownAlias string) {
	Dropdowns.Remove(layerAlias, dropdownAlias)
}

/*
DeleteAll is a method which allows you to delete all dropdowns from a text layer. In addition, the following
should be noted:

- This operation cannot be undone.

- All memory associated with the dropdowns will be freed.

:param layerAlias: The alias of the layer from which to delete all dropdowns.

Example:

	Dropdown.DeleteAll("layer1")
*/
func (shared *dropdownType) DeleteAll(layerAlias string) {
	Dropdowns.RemoveAll(layerAlias)
}

/*
drawOnLayer is a method which allows you to draw all dropdowns on a given text layer. In addition, the
following should be noted:

- Dropdowns are drawn in alphabetical order by their alias.

- This ensures consistent rendering order across multiple frames.

- The dropdown tray (selector) is only drawn when the dropdown is open.

:param layerEntry: The layer entry on which to draw the dropdowns.

Example:

	Dropdown.drawOnLayer(layer)
*/
func (shared *dropdownType) drawOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentDropdownEntry := range Dropdowns.GetAllEntries(layerAlias) {
		shared.draw(&layerEntry, currentDropdownEntry.Alias)
	}
}

/*
draw is a method which allows you to draw a single dropdown on a given text layer. In addition, the following
should be noted:

- The dropdown is drawn with a border and a down arrow indicator.

- The selected item text is formatted according to the specified width and alignment.

- The dropdown uses the style entry's foreground and background colors for rendering.

:param layerEntry: The layer entry on which to draw the dropdown.
:param dropdownAlias: The alias of the dropdown to draw.

Example:

	Dropdown.draw(layer, "myDropdown")
*/
func (shared *dropdownType) draw(layerEntry *types.LayerEntryType, dropdownAlias string) {
	layerAlias := layerEntry.LayerAlias
	dropdownEntry := Dropdowns.Get(layerAlias, dropdownAlias)
	localStyleEntry := types.NewTuiStyleEntry(&dropdownEntry.StyleEntry)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = localStyleEntry.Dropdown.ForegroundColor
	attributeEntry.BackgroundColor = localStyleEntry.Dropdown.BackgroundColor
	attributeEntry.CellType = constants.CellTypeDropdown
	attributeEntry.CellControlAlias = dropdownAlias

	var itemSelected string
	if len(dropdownEntry.SelectionEntry.SelectionValue) != 0 &&
		dropdownEntry.ItemSelected >= 0 && dropdownEntry.ItemSelected < len(dropdownEntry.SelectionEntry.SelectionValue) {
		itemSelected = dropdownEntry.SelectionEntry.SelectionValue[dropdownEntry.ItemSelected]
	}

	// We add +2 to account for the Dropdown border window which will appear. Otherwise, the item name
	// will appear 2 characters smaller than the popup Dropdown window.
	formattedItemName := stringformat.GetFormattedString(itemSelected, dropdownEntry.ItemWidth+2, localStyleEntry.Dropdown.TextAlignment)
	arrayOfRunes := stringformat.GetRunesFromString(formattedItemName)
	printLayer(layerEntry, attributeEntry, dropdownEntry.XLocation, dropdownEntry.YLocation, arrayOfRunes)
	// Invert colors for the dropdown arrow
	attributeEntry.ForegroundColor = localStyleEntry.Dropdown.BackgroundColor
	attributeEntry.BackgroundColor = localStyleEntry.Dropdown.ForegroundColor
	printLayer(layerEntry, attributeEntry, dropdownEntry.XLocation+len(arrayOfRunes), dropdownEntry.YLocation, []rune{constants.CharTriangleDown})
}

/*
updateStateMouse is a method which allows you to update the state of all dropdowns according to the current
mouse event state. In addition, the following should be noted:

- Handles mouse clicks to open/close dropdowns.

- Manages scrollbar synchronization for dropdowns with many items.

- Returns true if the screen needs to be updated due to state changes.

:return: Whether a screen update is required.

Example:

	isUpdate := Dropdown.updateStateMouse()
*/
func (shared *dropdownType) updateStateMouse() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	cellControlAlias := characterEntry.AttributeEntry.CellControlAlias

	// If a buttonType is pressed AND (you are in a drag and drop event OR the cell type is scroll bar), then
	// sync all Dropdown selectors with their appropriate scroll bars. If the control under focus
	// matches a control that belongs to a Dropdown list, then stop processing (Do not attempt to close Dropdown).
	if buttonPressed != 0 && (eventStateMemory.stateId == constants.EventStateDragAndDropScrollbar ||
		characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar) {
		isMatchFound := false
		for _, currentDropdownEntry := range Dropdowns.GetAllEntries(layerAlias) {
			dropdownEntry := currentDropdownEntry
			selectorEntry := Selectors.Get(layerAlias, dropdownEntry.SelectorAlias)
			scrollBarEntry := ScrollBars.Get(layerAlias, dropdownEntry.ScrollbarAlias)
			if selectorEntry.ViewportPosition != scrollBarEntry.ScrollValue {
				selectorEntry.ViewportPosition = scrollBarEntry.ScrollValue
				isUpdateRequired = true
			}
			if isControlCurrentlyFocused(layerAlias, dropdownEntry.Alias, constants.CellTypeDropdown) {
				isMatchFound = true
				break // If the current scrollbar being dragged and dropped matches, don't process more dropdowns.
			}
		}
		if isMatchFound {
			return isUpdateRequired
		}
	}

	// If our Dropdown alias is not empty, then open our Dropdown.
	if buttonPressed != 0 && cellControlAlias != "" && characterEntry.AttributeEntry.CellType == constants.CellTypeDropdown &&
		Dropdowns.IsExists(layerAlias, cellControlAlias) {
		shared.closeAllOpen()
		dropdownEntry := Dropdowns.Get(layerAlias, cellControlAlias)
		dropdownEntry.IsTrayOpen = true
		selectorEntry := Selectors.Get(layerAlias, dropdownEntry.SelectorAlias)
		selectorEntry.IsVisible = true
		scrollBarEntry := ScrollBars.Get(layerAlias, dropdownEntry.ScrollbarAlias)
		if scrollBarEntry.IsEnabled {
			scrollBarEntry.IsVisible = true
			setFocusedControl(layerAlias, dropdownEntry.Alias, constants.CellTypeDropdown)
		}
		isUpdateRequired = true
		return isUpdateRequired
	}

	// Only close dropdowns if clicking outside both the dropdown and its scrollbar
	_, _, previousButtonPress, _ := GetPreviousMouseStatus()
	if buttonPressed != 0 && previousButtonPress == 0 {
		// Check if we're clicking on a scrollbar that belongs to an open dropdown
		isScrollbarOfOpenDropdown := false
		if characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar {
			for _, currentDropdownEntry := range Dropdowns.GetAllEntries(layerAlias) {
				dropdownEntry := currentDropdownEntry
				if dropdownEntry.IsTrayOpen && dropdownEntry.ScrollbarAlias == cellControlAlias {
					isScrollbarOfOpenDropdown = true
					break
				}
			}
		}

		// Only close if not clicking on a dropdown or its scrollbar
		if characterEntry.AttributeEntry.CellType != constants.CellTypeDropdown && !isScrollbarOfOpenDropdown {
			isUpdateRequired = shared.closeAllOpen()
		}
	}
	return isUpdateRequired
}

/*
closeAllOpenOnLayer is a method which allows you to close all dropdowns for a given layer alias. In addition,
the following should be noted:

- This method is called when clicking outside of any dropdown.

- All open dropdown trays are closed and their scrollbars are hidden.

- The selected item is updated if it was changed while the dropdown was open.

:param layerAlias: The alias of the layer for which to close all open dropdowns.

:return: Whether any dropdown was actually closed.

Example:

	isClosed := Dropdown.closeAllOpenOnLayer("layer1")
*/
func (shared *dropdownType) closeAllOpenOnLayer(layerAlias string) bool {
	isAnyDropdownClosed := false
	for _, currentDropdownEntry := range Dropdowns.GetAllEntries(layerAlias) {
		dropdownEntry := currentDropdownEntry
		if dropdownEntry.IsTrayOpen == true {
			selectorEntry := Selectors.Get(layerAlias, dropdownEntry.SelectorAlias)
			selectorEntry.IsVisible = false
			scrollBarEntry := ScrollBars.Get(layerAlias, dropdownEntry.ScrollbarAlias)
			scrollBarEntry.IsVisible = false
			dropdownEntry.IsTrayOpen = false
			if dropdownEntry.ItemSelected != selectorEntry.ItemSelected {
				dropdownEntry.ItemSelected = selectorEntry.ItemSelected
			}
			setFocusedControl("", "", constants.NullCellType)
			// Reset the event state only if a tray is closed.
			eventStateMemory.stateId = constants.EventStateNone
			isAnyDropdownClosed = true
		}
	}
	return isAnyDropdownClosed
}

/*
closeAllOpen is a method which allows you to close all dropdowns on all layers. In addition, the following
should be noted:

- This method is useful for ensuring no dropdowns remain open when changing application state.

- All open dropdown trays are closed and their scrollbars are hidden.

- The selected item is updated if it was changed while the dropdown was open.

:return: Whether any dropdown was actually closed.

Example:

	isClosed := Dropdown.closeAllOpen()
*/
func (shared *dropdownType) closeAllOpen() bool {
	var wasAnyDropdownClosed bool
	Dropdowns.MemoryManager.Range(func(key, value interface{}) bool {
		layerAlias := key.(string)
		isDropdownClosed := shared.closeAllOpenOnLayer(layerAlias)
		if isDropdownClosed == true {
			wasAnyDropdownClosed = true
		}
		return true
	})
	return wasAnyDropdownClosed
}

/*
Get is a method which allows you to retrieve a dropdown entry from the control memory manager. In addition, the
following should be noted:

- Returns a pointer to the dropdown entry if it exists, nil otherwise.

- The dropdown entry contains all properties and state information for the control.

- This method is used internally by other dropdown methods to access control data.

:param layerAlias: The alias of the layer containing the dropdown.
:param dropdownAlias: The alias of the dropdown to retrieve.

:return: A pointer to the dropdown entry.

Example:

	entry := Dropdown.Get("layer1", "myDropdown")
*/
func (shared *dropdownType) Get(layerAlias string, dropdownAlias string) *types.DropdownEntryType {
	return Dropdowns.Get(layerAlias, dropdownAlias)
}

/*
IsExists is a method which allows you to check if a dropdown exists in the control memory manager. In addition, the
following should be noted:

- Returns true if the dropdown exists, false otherwise.

- This method is used to validate dropdown existence before performing operations.

- Useful for preventing null pointer exceptions when accessing dropdown properties.

:param layerAlias: The alias of the layer to check.
:param dropdownAlias: The alias of the dropdown to check.

:return: Whether the dropdown exists.

Example:

	exists := Dropdown.IsExists("layer1", "myDropdown")
*/
func (shared *dropdownType) IsExists(layerAlias string, dropdownAlias string) bool {
	return Dropdowns.IsExists(layerAlias, dropdownAlias)
}

/*
GetAllEntries is a method which allows you to retrieve all dropdown entries for a given layer. In addition, the
following should be noted:

- Returns a slice of all dropdown entries for the specified layer.

- The entries are returned in alphabetical order by their alias.

- This method is useful for iterating over all dropdowns on a layer.

:param layerAlias: The alias of the layer to retrieve entries for.

:return: A slice of pointers to dropdown entries.

Example:

	entries := Dropdown.GetAllEntries("layer1")
*/
func (shared *dropdownType) GetAllEntries(layerAlias string) []*types.DropdownEntryType {
	return Dropdowns.GetAllEntries(layerAlias)
}
