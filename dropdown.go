package consolizer

import (
	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/internal/memory"
	"supercom32.net/consolizer/internal/stringformat"
	"supercom32.net/consolizer/types"
)

type DropdownInstanceType struct {
	layerAlias   string
	controlAlias string
}

type dropdownType struct{}

/*
updateKeyboardEventDropdown allows you to update the state of all dropdowns according to the
current keyboard event. In addition, the following information should be noted:

- Handles Enter key to open/close the dropdown.
- Handles Up/Down keys to navigate through dropdown options when open.
- Returns true if the screen needs to be updated due to state changes.
*/
func (shared *dropdownType) updateKeyboardEventDropdown(keystroke []rune) bool {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := false
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType

	// Only process if a dropdown is focused
	if focusedControlType != constants.CellTypeDropdown || !Dropdowns.IsExists(focusedLayerAlias, focusedControlAlias) {
		return isScreenUpdateRequired
	}

	dropdownEntry := Dropdowns.Get(focusedLayerAlias, focusedControlAlias)

	// If dropdown is open but focus is on the dropdown itself (not the selector),
	// move focus to the selector for keyboard navigation
	if dropdownEntry.IsTrayOpen && focusedControlType == constants.CellTypeDropdown {
		isScreenUpdateRequired = Selector.updateKeyboardEventForSelector(focusedLayerAlias, dropdownEntry.SelectorAlias, keystroke)
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
			shared.closeAllOpenDropdowns(focusedLayerAlias) // Close any other open dropdowns first
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
	}

	return isScreenUpdateRequired
}

var Dropdown dropdownType
var Dropdowns = memory.NewControlMemoryManager[types.DropdownEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
Delete allows you to remove a dropdown from a text layer. In addition, the following
information should be noted:

- If you attempt to delete a dropdown which does not exist, then the request
will simply be ignored.
- All memory associated with the dropdown will be freed.
*/
func (shared *DropdownInstanceType) Delete() *DropdownInstanceType {
	if Dropdowns.IsExists(shared.layerAlias, shared.controlAlias) {
		Dropdowns.Remove(shared.layerAlias, shared.controlAlias)
	}
	return nil
}

/*
AddToTabIndex allows you to add a dropdown to the tab index. This enables keyboard navigation
between controls using the tab key. In addition, the following information should be noted:

- The dropdown will be added to the tab order based on the order in which it was created.
- The tab index is used to determine which control receives focus when the tab key is pressed.
*/
func (shared *DropdownInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeDropdown)
}

/*
GetValue allows you to retrieve the currently selected value from a dropdown. In addition,
the following information should be noted:

- Returns the display value of the currently selected item.
- If the dropdown does not exist, returns an empty string.
*/
func (shared *DropdownInstanceType) GetValue() string {
	dropdownEntry := Dropdowns.Get(shared.layerAlias, shared.controlAlias)
	return dropdownEntry.SelectionEntry.SelectionValue[dropdownEntry.ItemSelected]
}

/*
GetAlias allows you to retrieve the currently selected alias from a dropdown. In addition,
the following information should be noted:

- Returns the internal alias of the currently selected item.
- If the dropdown does not exist, returns an empty string.
- The alias is typically used for programmatic access to the selection.
*/
func (shared *DropdownInstanceType) GetAlias() string {
	dropdownEntry := Dropdowns.Get(shared.layerAlias, shared.controlAlias)
	return dropdownEntry.SelectionEntry.SelectionAlias[dropdownEntry.ItemSelected]
}

/*
Add allows you to create a new dropdown control on a text layer. In addition, the following
information should be noted:

- The dropdown consists of a main control and an associated selector for the dropdown tray.
- A scrollbar is automatically added if the number of items exceeds the selector height.
- The dropdown tray is initially hidden and only shown when the dropdown is clicked.
- The default selected item can be specified when creating the dropdown.
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
	Selector.Add(layerAlias, dropdownEntry.SelectorAlias, styleEntry, selectionEntry, xLocation+1, yLocation+1, selectorHeight, selectorWidth, 1, 0, 0, true)
	selectorEntry := Selectors.Get(layerAlias, dropdownEntry.SelectorAlias)
	selectorEntry.IsVisible = false
	dropdownEntry.ScrollbarAlias = selectorEntry.ScrollbarAlias
	scrollBarEntry := ScrollBars.Get(layerAlias, dropdownEntry.ScrollbarAlias)
	scrollBarEntry.IsVisible = false
	if len(selectionEntry.SelectionValue) <= selectorHeight {
		scrollBarEntry.IsEnabled = false
	}
	var dropdownInstance DropdownInstanceType
	dropdownInstance.layerAlias = layerAlias
	dropdownInstance.controlAlias = dropdownAlias
	return dropdownInstance
}

/*
DeleteDropdown allows you to remove a dropdown from a text layer. In addition, the following
information should be noted:

- If you attempt to delete a dropdown which does not exist, then the request
will simply be ignored.
- All memory associated with the dropdown will be freed.
*/
func (shared *dropdownType) DeleteDropdown(layerAlias string, dropdownAlias string) {
	Dropdowns.Remove(layerAlias, dropdownAlias)
}

/*
DeleteAllDropdowns allows you to remove all dropdowns from a text layer. In addition, the following
information should be noted:

- This operation cannot be undone.
- All memory associated with the dropdowns will be freed.
*/
func (shared *dropdownType) DeleteAllDropdowns(layerAlias string) {
	Dropdowns.RemoveAll(layerAlias)
}

/*
drawDropdownsOnLayer allows you to draw all dropdowns on a given text layer. In addition,
the following information should be noted:

- Dropdowns are drawn in alphabetical order by their alias.
- This ensures consistent rendering order across multiple frames.
- The dropdown tray (selector) is only drawn when the dropdown is open.
*/
func (shared *dropdownType) drawDropdownsOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentDropdownEntry := range Dropdowns.GetAllEntries(layerAlias) {
		shared.drawDropdown(&layerEntry, currentDropdownEntry.Alias)
	}
}

/*
drawDropdown allows you to draw a single dropdown on a given text layer. In addition, the following
information should be noted:

- The dropdown is drawn with a border and a down arrow indicator.
- The selected item text is formatted according to the specified width and alignment.
- The dropdown uses the style entry's foreground and background colors for rendering.
*/
func (shared *dropdownType) drawDropdown(layerEntry *types.LayerEntryType, dropdownAlias string) {
	layerAlias := layerEntry.LayerAlias
	dropdownEntry := Dropdowns.Get(layerAlias, dropdownAlias)
	localStyleEntry := types.NewTuiStyleEntry(&dropdownEntry.StyleEntry)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = localStyleEntry.SelectorForegroundColor
	attributeEntry.BackgroundColor = localStyleEntry.SelectorBackgroundColor
	attributeEntry.CellType = constants.CellTypeDropdown
	attributeEntry.CellControlAlias = dropdownAlias
	itemSelected := dropdownEntry.SelectionEntry.SelectionValue[dropdownEntry.ItemSelected]
	// We add +2 to account for the Dropdown border window which will appear. Otherwise, the item name
	// will appear 2 characters smaller than the popup Dropdown window.
	formattedItemName := stringformat.GetFormattedString(itemSelected, dropdownEntry.ItemWidth+2, localStyleEntry.SelectorTextAlignment)
	arrayOfRunes := stringformat.GetRunesFromString(formattedItemName)
	printLayer(layerEntry, attributeEntry, dropdownEntry.XLocation, dropdownEntry.YLocation, arrayOfRunes)
	attributeEntry.ForegroundColor = localStyleEntry.SelectorBackgroundColor
	attributeEntry.BackgroundColor = localStyleEntry.SelectorForegroundColor
	printLayer(layerEntry, attributeEntry, dropdownEntry.XLocation+len(arrayOfRunes), dropdownEntry.YLocation, []rune{constants.CharTriangleDown})
}

/*
updateDropdownStateMouse allows you to update the state of all dropdowns according to the current mouse event state.
In the event that a screen update is required this method returns true. In addition, the following information should be noted:

- Handles mouse clicks to open/close dropdowns.
- Manages scrollbar synchronization for dropdowns with many items.
- Returns true if the screen needs to be updated due to state changes.
*/
func (shared *dropdownType) updateDropdownStateMouse() bool {
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
		shared.closeAllOpenDropdowns(layerAlias)
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
			shared.closeAllOpenDropdowns(layerAlias)
		}
	}
	return isUpdateRequired
}

/*
closeAllOpenDropdowns allows you to close all dropdowns for a given layer alias. In addition,
the following information should be noted:

- This method is called when clicking outside of any dropdown.
- All open dropdown trays are closed and their scrollbars are hidden.
- The selected item is updated if it was changed while the dropdown was open.
*/
func (shared *dropdownType) closeAllOpenDropdowns(layerAlias string) {
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
		}
	}
}

/*
Get allows you to retrieve a dropdown entry from the control memory manager. In addition, the following
information should be noted:

- Returns a pointer to the dropdown entry if it exists, nil otherwise.
- The dropdown entry contains all properties and state information for the control.
- This method is used internally by other dropdown methods to access control data.
*/
func (shared *dropdownType) Get(layerAlias string, dropdownAlias string) *types.DropdownEntryType {
	return Dropdowns.Get(layerAlias, dropdownAlias)
}

/*
IsExists allows you to check if a dropdown exists in the control memory manager. In addition, the following
information should be noted:

- Returns true if the dropdown exists, false otherwise.
- This method is used to validate dropdown existence before performing operations.
- Useful for preventing null pointer exceptions when accessing dropdown properties.
*/
func (shared *dropdownType) IsExists(layerAlias string, dropdownAlias string) bool {
	return Dropdowns.IsExists(layerAlias, dropdownAlias)
}

/*
GetAllEntries allows you to retrieve all dropdown entries for a given layer. In addition, the following
information should be noted:

- Returns a slice of all dropdown entries for the specified layer.
- The entries are returned in alphabetical order by their alias.
- This method is useful for iterating over all dropdowns on a layer.
*/
func (shared *dropdownType) GetAllEntries(layerAlias string) []*types.DropdownEntryType {
	return Dropdowns.GetAllEntries(layerAlias)
}
