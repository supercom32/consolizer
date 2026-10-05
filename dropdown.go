package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
)

/*
DropdownInstanceType is a structure which represents an instance of a dropdown control.
*/
type DropdownInstanceType struct {
	BaseControlInstanceType
}

type dropdownType struct{}

/*
updateKeyboardEvent is a method which allows you to update the focused dropdown according to the current keystroke.
Enter opens the tray, or closes it and commits the tray's selection, the arrow keys move through the items while the
tray is open, and Esc closes an open tray without changing the selection. It returns whether the screen needs
updating and whether the keystroke was consumed. In addition, the following should be noted:

  - Esc on a closed dropdown is left unconsumed so the application can treat it as a Back or Cancel action.

  - Every index copied between the dropdown and its tray is brought back into range for the tray's items, so a
    dropdown with no items, or with no selection, can be opened, navigated and closed safely.

Example:

	isUpdateRequired, isConsumed := Dropdown.updateKeyboardEvent(keystroke)
*/
func (shared *dropdownType) updateKeyboardEvent(keystroke []rune) (bool, bool) {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := false
	isKeystrokeConsumed := false
	focusedControl := getFocusedControl()
	focusedLayerAlias := focusedControl.layerAlias
	focusedControlAlias := focusedControl.controlAlias
	focusedControlType := focusedControl.controlType

	// Only process if a dropdown is focused
	dropdownEntry, isFound := Dropdowns.Lookup(focusedLayerAlias, focusedControlAlias)
	if focusedControlType != constants.CellTypeDropdown || !isFound {
		return isScreenUpdateRequired, isKeystrokeConsumed
	}

	// If dropdown is open but focus is on the dropdown itself (not the selector),
	// move focus to the selector for keyboard navigation
	if dropdownEntry.IsTrayOpen && focusedControlType == constants.CellTypeDropdown {
		isScreenUpdateRequired, isKeystrokeConsumed = Selector.updateKeyboardEventForSelector(focusedLayerAlias, dropdownEntry.SelectorAlias, keystroke)
	}

	// Handle Enter key to open/close dropdown
	if keystrokeAsString == "enter" {
		if dropdownEntry.IsTrayOpen {
			// Close dropdown and apply selection
			selectorEntry := Selectors.Get(focusedLayerAlias, dropdownEntry.SelectorAlias)
			scrollBarEntry := ScrollBars.Get(focusedLayerAlias, dropdownEntry.ScrollbarAlias)
			scrollBarEntry.ScrollValue = selectorEntry.ItemSelected
			scrollbar.computeHandlePositionByScrollValue(focusedLayerAlias, dropdownEntry.ScrollbarAlias)
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
			setEventStateId(constants.EventStateNone)
		} else {
			// Open dropdown
			shared.closeAllOpen() // Close any other open dropdowns first
			dropdownEntry.IsTrayOpen = true

			// Show dropdown components
			selectorEntry := Selectors.Get(focusedLayerAlias, dropdownEntry.SelectorAlias)
			selectorEntry.IsVisible = true
			selectorEntry.ItemHighlighted = dropdownEntry.ItemSelected // Highlight current selection
			Selector.normalizeIndexes(selectorEntry)

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

	// Handle Esc key to close dropdown without changing selection. tcell reports the Escape key as "Esc".
	if keystrokeAsString == "esc" && dropdownEntry.IsTrayOpen {
		selectorEntry := Selectors.Get(focusedLayerAlias, dropdownEntry.SelectorAlias)
		scrollBarEntry := ScrollBars.Get(focusedLayerAlias, dropdownEntry.ScrollbarAlias)

		// Revert the tray to the committed selection, so a later close that applies the tray's selection (such as
		// closeAllOpen) cannot commit anything picked before Esc was pressed.
		selectorEntry.ItemSelected = dropdownEntry.ItemSelected
		selectorEntry.ItemHighlighted = dropdownEntry.ItemSelected
		Selector.normalizeIndexes(selectorEntry)

		// Hide dropdown components
		selectorEntry.IsVisible = false
		scrollBarEntry.IsVisible = false
		dropdownEntry.IsTrayOpen = false

		// Reset focus to the dropdown itself
		setFocusedControl(focusedLayerAlias, focusedControlAlias, constants.CellTypeDropdown)
		setEventStateId(constants.EventStateNone)
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
Delete is a method which removes a dropdown from a text layer. In addition, the following should be noted:

- If you attempt to delete a dropdown which does not exist, then the request will simply be ignored.

- All memory associated with the dropdown will be freed.

Example:
    dropdown.Delete()
*/
func (shared *DropdownInstanceType) Delete() *DropdownInstanceType {
	shared.BaseControlInstanceType.Delete()
	return nil
}

/*
GetValue is a method which retrieves the currently selected value from a dropdown. In addition, the
following should be noted:

- Returns the display value of the currently selected item.

- If the dropdown does not exist, returns an empty string.

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
GetAlias is a method which retrieves the currently selected alias from a dropdown. In addition, the
following should be noted:

- Returns the internal alias of the currently selected item.

- If the dropdown does not exist, returns an empty string.

- The alias is typically used for programmatic access to the selection.

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
GetSelectedItemIndex is a method which retrieves the index of the currently selected item in the dropdown.

Example:
    index := dropdown.GetSelectedItemIndex()
*/
func (shared *DropdownInstanceType) GetSelectedItemIndex() int {
	dropdownEntry := Dropdowns.Get(shared.layerAlias, shared.controlAlias)
	return dropdownEntry.ItemSelected
}

/*
SetSelectedItemIndex is a method which allows you to set the currently selected item of a dropdown by its zero based
index. Passing constants.SELECTED_NONE, which is minus one, clears the selection. Any other index that is negative or
past the last item is ignored, leaving the selection unchanged.

Example:

	dropdown.SetSelectedItemIndex(2)
*/
func (shared *DropdownInstanceType) SetSelectedItemIndex(itemIndex int) {
	dropdownEntry := Dropdowns.Get(shared.layerAlias, shared.controlAlias)
	if itemIndex < constants.SELECTED_NONE || itemIndex >= Dropdown.getItemCount(dropdownEntry) {
		return
	}
	dropdownEntry.ItemSelected = itemIndex
	if selectorEntry, isFound := Selectors.Lookup(shared.layerAlias, dropdownEntry.SelectorAlias); isFound {
		selectorEntry.ItemSelected = itemIndex
		Selector.normalizeIndexes(selectorEntry)
	}
}

/*
getItemCount is a method which allows you to obtain the number of items of a dropdown that can be safely indexed. It is
the length of the shorter of the alias and value slices, so any index below it is valid for both.

Example:

	itemCount := Dropdown.getItemCount(dropdownEntry)
*/
func (shared *dropdownType) getItemCount(dropdownEntry *types.DropdownEntryType) int {
	return getClampedIndex(len(dropdownEntry.SelectionEntry.SelectionAlias), 0, len(dropdownEntry.SelectionEntry.SelectionValue))
}

/*
SetSelectionEntry is a method which allows you to replace the list of items a dropdown offers. The selected item is
cleared, and the dropdown's tray is reset to match the new list: its highlight is cleared, it scrolls back to the
first item, and its scroll bar is resized, enabled only while the items overflow the tray, and shown only while the
tray is open. This is safe to call whether the tray is open or closed, and with an empty list. In addition, the
following should be noted:

  - The dropdown keeps its own copy of the item slices, so changing the slices of the entry passed in afterwards
    does not change the dropdown.

Example:

	dropdown.SetSelectionEntry(newSelection)
*/
func (shared *DropdownInstanceType) SetSelectionEntry(selectionEntry types.SelectionEntryType) {
	dropdownEntry := Dropdowns.Get(shared.layerAlias, shared.controlAlias)
	dropdownEntry.SelectionEntry = Selector.getSelectionEntryCopy(selectionEntry)
	dropdownEntry.ItemSelected = constants.SELECTED_NONE

	selectorEntry, isFound := Selectors.Lookup(shared.layerAlias, dropdownEntry.SelectorAlias)
	if !isFound {
		return
	}
	selectorEntry.SelectionEntry = Selector.getSelectionEntryCopy(selectionEntry)
	selectorEntry.ItemSelected = constants.SELECTED_NONE
	selectorEntry.ItemHighlighted = constants.NullItemSelection
	selectorEntry.ViewportPosition = 0
	Selector.updateScrollbar(shared.layerAlias, selectorEntry)
}

/*
Add is a method which allows you to create a new dropdown control on a text layer. The dropdown is made of a main
control showing the selected item and a hidden selector for its tray, which is shown when the dropdown is opened, plus
a scroll bar that is enabled when the items exceed the tray height. The default selected item is given by its zero
based index. In addition, the following should be noted:

  - A default item index outside the items given is stored as constants.SELECTED_NONE, so the dropdown starts with
    no selection rather than an invalid one.

  - The tray starts with the same selected item as the dropdown, so opening the tray and closing it again without
    picking anything leaves the selection unchanged.

  - The dropdown keeps its own copy of the item slices.

Example:

	dropdown := Dropdown.Add("layer1", "myDropdown", style, items, 10, 10, 5, 15, 0)
*/
func (shared *dropdownType) Add(layerAlias string, dropdownAlias string, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, defaultItemSelected int) DropdownInstanceType {
	newDropdownEntry := types.NewDropdownEntry()
	newDropdownEntry.Alias = dropdownAlias
	newDropdownEntry.StyleEntry = styleEntry
	newDropdownEntry.SelectionEntry = Selector.getSelectionEntryCopy(selectionEntry)
	newDropdownEntry.XLocation = xLocation
	newDropdownEntry.YLocation = yLocation
	newDropdownEntry.ItemWidth = itemWidth
	newDropdownEntry.ItemSelected = defaultItemSelected
	if defaultItemSelected < 0 || defaultItemSelected >= shared.getItemCount(&newDropdownEntry) {
		newDropdownEntry.ItemSelected = constants.SELECTED_NONE
	}
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
	// The tray starts out agreeing with the dropdown, so closing it without a pick cannot change the selection.
	selectorEntry.ItemSelected = dropdownEntry.ItemSelected
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
Delete is a method which removes a dropdown from a text layer. In addition, the following should be
noted:

- If you attempt to delete a dropdown which does not exist, then the request will simply be ignored.

- All memory associated with the dropdown will be freed.

Example:
    Dropdown.Delete("layer1", "myDropdown")
*/
func (shared *dropdownType) Delete(layerAlias string, dropdownAlias string) {
	Dropdowns.Remove(layerAlias, dropdownAlias)
	clearStaleControlReferences(layerAlias, dropdownAlias, constants.CellTypeDropdown)
}

/*
DeleteAll is a method which deletes all dropdowns from a text layer. In addition, the following
should be noted:

- This operation cannot be undone.

- All memory associated with the dropdowns will be freed.

Example:
    Dropdown.DeleteAll("layer1")
*/
func (shared *dropdownType) DeleteAll(layerAlias string) {
	removeAllControlsAndClearCells(Dropdowns, layerAlias, constants.CellTypeDropdown, func(entry *types.DropdownEntryType) string { return entry.Alias })
}

/*
drawOnLayer is a method which draws all dropdowns on a given text layer. In addition, the
following should be noted:

- Dropdowns are drawn in alphabetical order by their alias.

- This ensures consistent rendering order across multiple frames.

- The dropdown tray (selector) is only drawn when the dropdown is open.

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
draw is a method which allows you to draw a single dropdown on a given text layer, showing its selected item formatted
to the dropdown's width and alignment, in the style's colours, followed by a down arrow drawn in inverted colours. In
addition, the following should be noted:

  - While the dropdown has focus, it is drawn with the style's focused colours, resolved by getFocusedColors so that
    focus stays visible even when the style does not set them, and the arrow inverts those instead. As described for
    setFocusIndicatorVisible, these colours are only shown while the user's last input came from the keyboard.

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
	if isFocusIndicatorShown(layerAlias, dropdownAlias, constants.CellTypeDropdown) {
		attributeEntry.ForegroundColor, attributeEntry.BackgroundColor = getFocusedColors(localStyleEntry.Dropdown.ForegroundColor,
			localStyleEntry.Dropdown.BackgroundColor, localStyleEntry.Dropdown.FocusedForegroundColor, localStyleEntry.Dropdown.FocusedBackgroundColor)
	}
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
	attributeEntry.ForegroundColor, attributeEntry.BackgroundColor = attributeEntry.BackgroundColor, attributeEntry.ForegroundColor
	printLayer(layerEntry, attributeEntry, dropdownEntry.XLocation+stringformat.GetWidthOfRunesWhenPrinted(arrayOfRunes), dropdownEntry.YLocation, []rune{constants.CharTriangleDown})
}

/*
updateStateMouse is a method which allows you to update the state of all dropdowns according to the current mouse
event state. Clicking a dropdown opens its tray, clicking elsewhere closes every open tray, and dragging a tray's
scroll bar moves the tray's viewport. It returns true if the screen needs to be updated. In addition, the following
should be noted:

  - A viewport position copied from a scroll bar is clamped into the range the tray can scroll to.

Example:

	isUpdateRequired := Dropdown.updateStateMouse()
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
	if buttonPressed != 0 && (getEventStateId() == constants.EventStateDragAndDropScrollbar ||
		characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar) {
		isMatchFound := false
		for _, currentDropdownEntry := range Dropdowns.GetAllEntries(layerAlias) {
			dropdownEntry := currentDropdownEntry
			selectorEntry := Selectors.Get(layerAlias, dropdownEntry.SelectorAlias)
			scrollBarEntry := ScrollBars.Get(layerAlias, dropdownEntry.ScrollbarAlias)
			if selectorEntry.ViewportPosition != scrollBarEntry.ScrollValue {
				selectorEntry.ViewportPosition = scrollBarEntry.ScrollValue
				Selector.normalizeIndexes(selectorEntry)
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
	if dropdownEntry, isFound := Dropdowns.Lookup(layerAlias, cellControlAlias); buttonPressed != 0 && cellControlAlias != "" &&
		characterEntry.AttributeEntry.CellType == constants.CellTypeDropdown && isFound {
		shared.closeAllOpen()
		dropdownEntry.IsTrayOpen = true
		selectorEntry := Selectors.Get(layerAlias, dropdownEntry.SelectorAlias)
		selectorEntry.IsVisible = true
		scrollBarEntry := ScrollBars.Get(layerAlias, dropdownEntry.ScrollbarAlias)
		if scrollBarEntry.IsEnabled {
			scrollBarEntry.IsVisible = true
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
closeAllOpenOnLayer is a method which closes all dropdowns for a given layer alias. In addition,
the following should be noted:

- This method is called when clicking outside of any dropdown.

- All open dropdown trays are closed and their scrollbars are hidden.

- The selected item is updated if it was changed while the dropdown was open.

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
			// Focus is left where it is: on the dropdown itself, or on whatever control the click that closed the
			// tray landed on.
			// Reset the event state only if a tray is closed.
			setEventStateId(constants.EventStateNone)
			isAnyDropdownClosed = true
		}
	}
	return isAnyDropdownClosed
}

/*
closeAllOpen is a method which closes all dropdowns on all layers. In addition, the following
should be noted:

- This method is useful for ensuring no dropdowns remain open when changing application state.

- All open dropdown trays are closed and their scrollbars are hidden.

- The selected item is updated if it was changed while the dropdown was open.

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
Get is a method which retrieves a dropdown entry from the control memory manager. In addition, the
following should be noted:

- Returns a pointer to the dropdown entry if it exists, nil otherwise.

- The dropdown entry contains all properties and state information for the control.

- This method is used internally by other dropdown methods to access control data.

Example:
    entry := Dropdown.Get("layer1", "myDropdown")
*/
func (shared *dropdownType) Get(layerAlias string, dropdownAlias string) *types.DropdownEntryType {
	return Dropdowns.Get(layerAlias, dropdownAlias)
}

/*
IsExists is a method which checks if a dropdown exists in the control memory manager. In addition, the
following should be noted:

- Returns true if the dropdown exists, false otherwise.

- This method is used to validate dropdown existence before performing operations.

- Useful for preventing null pointer exceptions when accessing dropdown properties.

Example:
    exists := Dropdown.IsExists("layer1", "myDropdown")
*/
func (shared *dropdownType) IsExists(layerAlias string, dropdownAlias string) bool {
	return Dropdowns.IsExists(layerAlias, dropdownAlias)
}

/*
GetAllEntries is a method which retrieves all dropdown entries for a given layer. In addition, the
following should be noted:

- Returns a slice of all dropdown entries for the specified layer.

- The entries are returned in alphabetical order by their alias.

- This method is useful for iterating over all dropdowns on a layer.

Example:
    entries := Dropdown.GetAllEntries("layer1")
*/
func (shared *dropdownType) GetAllEntries(layerAlias string) []*types.DropdownEntryType {
	return Dropdowns.GetAllEntries(layerAlias)
}
