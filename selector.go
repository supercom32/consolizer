package consolizer

import (
	"fmt"

	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/internal/memory"
	"supercom32.net/consolizer/internal/stringformat"
	"supercom32.net/consolizer/types"
)

type selectorInstanceType struct {
	layerAlias   string
	controlAlias string
}

type selectorType struct{}

var Selector selectorType
var Selectors = memory.NewControlMemoryManager[types.SelectorEntryType]()

func AddSelector(layerAlias string, selectorAlias string, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, itemSelected int, isBorderDrawn bool) {
	selectorEntry := types.NewSelectorEntry()
	selectorEntry.Alias = selectorAlias
	selectorEntry.StyleEntry = styleEntry
	selectorEntry.SelectionEntry = selectionEntry
	selectorEntry.XLocation = xLocation
	selectorEntry.YLocation = yLocation
	selectorEntry.Height = selectorHeight
	selectorEntry.ItemWidth = itemWidth
	selectorEntry.NumberOfColumns = numberOfColumns
	selectorEntry.ViewportPosition = viewportPosition
	selectorEntry.ItemHighlighted = itemSelected
	selectorEntry.IsBorderDrawn = isBorderDrawn
	selectorEntry.IsVisible = true

	// Use the generic memory manager to add the selector entry
	Selectors.Add(layerAlias, selectorAlias, &selectorEntry)
}

func DeleteSelector(layerAlias string, selectorAlias string) {
	// Use the generic memory manager to remove the selector entry
	Selectors.Remove(layerAlias, selectorAlias)
}

func DeleteAllSelectorsFromLayer(layerAlias string) {
	// Retrieve all selectors in the specified layer
	selectors := Selectors.GetAllEntries(layerAlias)

	// Loop through all entries and delete them
	for _, selector := range selectors {
		Selectors.Remove(layerAlias, selector.Alias) // Assuming selector.Alias contains the alias
	}
}

func IsSelectorExists(layerAlias string, selectorAlias string) bool {
	// Use the generic memory manager to check existence
	return Selectors.IsExists(layerAlias, selectorAlias)
}

func GetSelector(layerAlias string, selectorAlias string) *types.SelectorEntryType {
	// Use the generic memory manager to retrieve the selector entry
	selectorEntry := Selectors.Get(layerAlias, selectorAlias)
	if selectorEntry == nil {
		panic(fmt.Sprintf("The selector '%s' under layer '%s' could not be obtained since it does not exist!", selectorAlias, layerAlias))
	}
	return selectorEntry
}

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
AddToTabIndex allows you to add a selector to the tab index. This enables keyboard navigation
between controls using the tab key. In addition, the following information should be noted:

- The selector will be added to the tab order based on the order in which it was created.
- The tab index is used to determine which control receives focus when the tab key is pressed.
*/
func (shared *selectorInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeSelectorItem)
}

/*
Delete allows you to remove a selector from a text layer. In addition, the following
information should be noted:

- If you attempt to delete a selector which does not exist, then the request
will simply be ignored.
- All memory associated with the selector will be freed.
*/
func (shared *selectorInstanceType) Delete() *selectorInstanceType {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		Selectors.Remove(shared.layerAlias, shared.controlAlias)
	}
	return nil
}

/*
GetSelected allows you to retrieve the currently selected item from a selector. In addition,
the following information should be noted:

  - Returns both the alias and index of the selected item.
  - If the selector does not exist, returns an empty string and -1.
  - The alias is typically used for display purposes, while the index is used for
    programmatic access to the selection.
*/
func (shared *selectorInstanceType) GetSelected() (string, int) {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorMenu(shared.layerAlias, shared.controlAlias)
		menuEntry := Selectors.Get(shared.layerAlias, shared.controlAlias)
		value := menuEntry.ItemSelected
		return menuEntry.SelectionEntry.SelectionAlias[value], value
	}
	return "", -1
}

/*
setViewport allows you to specify the current viewport index for a given selector. In addition,
the following information should be noted:

- The viewport determines which items are currently visible in the selector.
- If the selector does not exist, no operation occurs.
- The viewport position is automatically adjusted when navigating through items.
*/
func (shared *selectorInstanceType) setViewport(viewportPosition int) {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorMenu(shared.layerAlias, shared.controlAlias)
		menuEntry := Selectors.Get(shared.layerAlias, shared.controlAlias)
		menuEntry.ViewportPosition = viewportPosition
	}
}

/*
Add allows you to add a selector to a given text layer. Once called, an instance of your control is returned
which will allow you to read or manipulate the properties for it. The Style of the Selector
will be determined by the style entry passed in. If you wish to remove a Selector from a text layer, simply
call 'DeleteSelector'. In addition, the following information should be noted:

- Selectors are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create selectors without actually overwriting
the text layer data under it.

- If the Selector to be drawn falls outside the range of the provided layer,
then only the visible portion of the radio button will be drawn.

- If the Selector height is greater than the number of selections available, then no scroll bars are drawn.
*/
// TODO: Protect against viewport out of range errors.
func (shared *selectorType) Add(layerAlias string, selectorAlias string, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, selectedItem int, isBorderDrawn bool) selectorInstanceType {
	validateSelectionEntry(selectionEntry)
	newSelectorEntry := types.NewSelectorEntry()
	newSelectorEntry.Alias = selectorAlias
	newSelectorEntry.StyleEntry = styleEntry
	newSelectorEntry.SelectionEntry = selectionEntry
	newSelectorEntry.XLocation = xLocation
	newSelectorEntry.YLocation = yLocation
	newSelectorEntry.Height = selectorHeight
	newSelectorEntry.ItemWidth = itemWidth
	newSelectorEntry.NumberOfColumns = numberOfColumns
	newSelectorEntry.ViewportPosition = viewportPosition
	newSelectorEntry.ItemHighlighted = selectedItem
	newSelectorEntry.IsBorderDrawn = isBorderDrawn
	newSelectorEntry.IsVisible = true

	// Use the generic memory manager to add the selector entry
	Selectors.Add(layerAlias, selectorAlias, &newSelectorEntry)
	// TODO: AddLayer verification to ensure no item can be 0 length/number.

	selectorEntry := Selectors.Get(layerAlias, selectorAlias)
	selectorEntry.ScrollbarAlias = stringformat.GetLastSortedUUID()
	scrollBarMaxValue := len(selectionEntry.SelectionValue) - (selectorHeight * numberOfColumns) + 1
	scrollBarXLocation := xLocation + (itemWidth * numberOfColumns)
	scrollBarYLocation := yLocation
	scrollBarHeight := selectorHeight
	if isBorderDrawn {
		scrollBarXLocation = xLocation + (itemWidth * numberOfColumns) + 1
		scrollBarYLocation = scrollBarYLocation - 1
		scrollBarHeight = selectorHeight + 2
	}
	scrollbar.Add(layerAlias, selectorEntry.ScrollbarAlias, styleEntry, scrollBarXLocation, scrollBarYLocation, scrollBarHeight, scrollBarMaxValue, 0, numberOfColumns, false)
	scrollBarEntry := ScrollBars.Get(layerAlias, selectorEntry.ScrollbarAlias)
	selectorWidth := itemWidth
	if len(selectionEntry.SelectionValue) <= selectorHeight*numberOfColumns || styleEntry.SelectorTextAlignment == constants.AlignmentNoPadding {
		scrollBarEntry.IsEnabled = false
		scrollBarEntry.IsVisible = false
		selectorWidth = selectorWidth + 1
	}
	var selectorInstance selectorInstanceType
	selectorInstance.layerAlias = layerAlias
	selectorInstance.controlAlias = selectorAlias
	setFocusedControl(layerAlias, selectorAlias, constants.CellTypeSelectorItem)
	return selectorInstance
}

/*
DeleteSelector allows you to remove a selector from a text layer. In addition, the following
information should be noted:

- If you attempt to delete a selector which does not exist, then the request
will simply be ignored.
- All memory associated with the selector will be freed.
*/
func (shared *selectorType) DeleteSelector(layerAlias string, selectorAlias string) {
	Selectors.Remove(layerAlias, selectorAlias)
}

/*
DeleteAllSelectors allows you to remove all selectors from a text layer. In addition, the following
information should be noted:

- This operation cannot be undone.
- All memory associated with the selectors will be freed.
*/
func (shared *selectorType) DeleteAllSelectors(layerAlias string) {
	Selectors.RemoveAll(layerAlias)
}

/*
drawSelector allows you to draw a selector on a given text layer. The
Style of the Selector will be determined by the style entry passed in. In
addition, the following information should be noted:

- Selectors are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create selectors without actually overwriting
the text layer data under it.

- If the Selector to be drawn falls outside the range of the provided layer,
then only the visible portion of the Selector will be drawn.
*/
func (shared *selectorType) drawSelector(selectorAlias string, layerEntry *types.LayerEntryType, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, itemHighlighted int) {
	selectorEntry := Selectors.Get(layerEntry.LayerAlias, selectorAlias)
	if selectorEntry.IsVisible == false {
		return
	}
	menuAttributeEntry := types.NewAttributeEntry()
	menuAttributeEntry.ForegroundColor = styleEntry.SelectorForegroundColor
	menuAttributeEntry.BackgroundColor = styleEntry.SelectorBackgroundColor
	if selectorEntry.IsBorderDrawn {
		fillArea(layerEntry, menuAttributeEntry, " ", xLocation-1, yLocation-1, itemWidth+2, selectorHeight+2, constants.NullCellControlLocation)
		drawBorder(layerEntry, styleEntry, menuAttributeEntry, xLocation-1, yLocation-1, itemWidth+2, selectorHeight+2, false)
	}
	highlightAttributeEntry := types.NewAttributeEntry()
	highlightAttributeEntry.ForegroundColor = styleEntry.HighlightForegroundColor
	highlightAttributeEntry.BackgroundColor = styleEntry.HighlightBackgroundColor
	currentYLocation := yLocation
	currentMenuItemIndex := viewportPosition
	currentXOffset := 0
	currentColumn := 0
	currentRow := 0
	for currentMenuItemIndex < len(selectionEntry.SelectionValue) && currentRow < selectorHeight {
		attributeEntry := menuAttributeEntry
		if currentMenuItemIndex == itemHighlighted {
			attributeEntry = highlightAttributeEntry
		}
		menuItemName := stringformat.GetFormattedString(selectionEntry.SelectionValue[currentMenuItemIndex], itemWidth, styleEntry.SelectorTextAlignment)
		arrayOfRunes := stringformat.GetRunesFromString(menuItemName)
		attributeEntry.CellControlId = currentMenuItemIndex
		attributeEntry.CellControlAlias = selectorAlias
		attributeEntry.CellType = constants.CellTypeSelectorItem
		printLayer(layerEntry, attributeEntry, xLocation+(currentXOffset), currentYLocation, arrayOfRunes)
		currentMenuItemIndex++
		currentXOffset = currentXOffset + stringformat.GetWidthOfRunesWhenPrinted(arrayOfRunes) // len(arrayOfRunes)
		currentColumn++
		if currentColumn >= numberOfColumns {
			currentXOffset = 0
			currentColumn = 0
			currentYLocation++
			currentRow++
		}
	}
}

/*
drawSelectorsOnLayer allows you to draw all selectors on a given text layer. In addition,
the following information should be noted:

- Selectors are drawn in alphabetical order by their alias.
- This ensures consistent rendering order across multiple frames.
- Internally generated selectors (like those used by dropdowns) are drawn last.
*/
func (shared *selectorType) drawSelectorsOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	compareByAlias := func(a, b *types.SelectorEntryType) bool {
		return a.Alias < b.Alias
	}
	// Sort array so internally generated selectors appear last (Since sorted by name, and
	// UUID generates "zzz" prefixes). This prevents Dropdown selectors from appearing under
	// user created selectors, when they should always be on top.
	for _, currentSelectorEntry := range Selectors.SortEntries(layerAlias, true, compareByAlias) {
		selectorEntry := currentSelectorEntry
		shared.drawSelector(selectorEntry.Alias, &layerEntry, selectorEntry.StyleEntry, selectorEntry.SelectionEntry, selectorEntry.XLocation, selectorEntry.YLocation, selectorEntry.Height, selectorEntry.ItemWidth, selectorEntry.NumberOfColumns, selectorEntry.ViewportPosition, selectorEntry.ItemHighlighted)
	}
}

func (shared *selectorType) updateKeyboardEventForSelector(layerAlias string, selectorAlias string, keystroke []rune) bool {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := false
	selectorEntry := Selectors.Get(layerAlias, selectorAlias)
	if keystrokeAsString == "down" {
		// remainder := selectorEntry.ItemHighlighted % selectorEntry.NumberOfColumns
		selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + selectorEntry.NumberOfColumns
		if selectorEntry.ItemHighlighted >= len(selectorEntry.SelectionEntry.SelectionAlias) {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - selectorEntry.NumberOfColumns
		}
		// Adjust viewport if highlighted item is outside visible range
		if selectorEntry.ItemHighlighted >= selectorEntry.ViewportPosition+(selectorEntry.Height*selectorEntry.NumberOfColumns) {
			selectorEntry.ViewportPosition = selectorEntry.ItemHighlighted - (selectorEntry.Height * selectorEntry.NumberOfColumns) + selectorEntry.NumberOfColumns
			// Update associated scrollbar
			if scrollBarEntry := ScrollBars.Get(layerAlias, selectorEntry.ScrollbarAlias); scrollBarEntry != nil {
				scrollBarEntry.ScrollValue = selectorEntry.ViewportPosition
				scrollbar.computeScrollbarHandlePositionByScrollValue(layerAlias, selectorEntry.ScrollbarAlias)
			}
		}
		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "up" {
		selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - selectorEntry.NumberOfColumns
		if selectorEntry.ItemHighlighted < 0 {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + selectorEntry.NumberOfColumns
		}
		// Adjust viewport if highlighted item is outside visible range
		if selectorEntry.ItemHighlighted < selectorEntry.ViewportPosition {
			selectorEntry.ViewportPosition = selectorEntry.ItemHighlighted
			// Update associated scrollbar
			if scrollBarEntry := ScrollBars.Get(layerAlias, selectorEntry.ScrollbarAlias); scrollBarEntry != nil {
				scrollBarEntry.ScrollValue = selectorEntry.ViewportPosition
				scrollbar.computeScrollbarHandlePositionByScrollValue(layerAlias, selectorEntry.ScrollbarAlias)
			}
		}
		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "left" {
		if selectorEntry.ItemHighlighted%selectorEntry.NumberOfColumns != 0 {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - 1
			if selectorEntry.ItemHighlighted < 0 {
				selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + 1
			}
			isScreenUpdateRequired = true
		}
	}
	if keystrokeAsString == "right" {
		if selectorEntry.ItemHighlighted%selectorEntry.NumberOfColumns != selectorEntry.NumberOfColumns-1 {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + 1
			if selectorEntry.ItemHighlighted >= len(selectorEntry.SelectionEntry.SelectionAlias) {
				selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - 1
			}
			isScreenUpdateRequired = true
		}
	}
	if keystrokeAsString == "enter" {
		selectorEntry.ItemSelected = selectorEntry.ItemHighlighted
		isScreenUpdateRequired = true
	}
	return isScreenUpdateRequired
}

/*
updateKeyboardEvent allows you to update the state of all selectors according to the current keystroke event.
In the event that a screen update is required this method returns true. In addition, the following information should be noted:

- Handles navigation keys (up, down, left, right) to move between items.
- Enter key selects the currently highlighted item.
- Returns true if the screen needs to be updated due to state changes.
*/
func (shared *selectorType) updateKeyboardEvent(keystroke []rune) bool {
	isScreenUpdateRequired := false
	if eventStateMemory.currentlyFocusedControl.controlType != constants.CellTypeSelectorItem || !Selectors.IsExists(eventStateMemory.currentlyFocusedControl.layerAlias, eventStateMemory.currentlyFocusedControl.controlAlias) {
		return isScreenUpdateRequired
	}
	return shared.updateKeyboardEventForSelector(eventStateMemory.currentlyFocusedControl.layerAlias, eventStateMemory.currentlyFocusedControl.controlAlias, keystroke)
}

/*
updateMouseEvent allows you to update the state of all selectors according to the current mouse event state.
In the event that a screen update is required this method returns true. In addition, the following information should be noted:

- Handles mouse clicks to select items.
- Manages scrollbar synchronization for selectors with many items.
- Returns true if the screen needs to be updated due to state changes.
*/
func (shared *selectorType) updateMouseEvent() bool {
	isScreenUpdateRequired := false
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	var characterEntry types.CharacterEntryType
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	if characterEntry.AttributeEntry.CellType == constants.CellTypeSelectorItem && eventStateMemory.stateId == constants.EventStateNone && Selectors.IsExists(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias) {
		if buttonPressed != 0 {
			selectorEntry := Selectors.Get(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			selectorEntry.ItemHighlighted = characterEntry.AttributeEntry.CellControlId
			selectorEntry.ItemSelected = characterEntry.AttributeEntry.CellControlId
		} else {
			selectorEntry := Selectors.Get(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			selectorEntry.ItemHighlighted = characterEntry.AttributeEntry.CellControlId
		}
		// Check if this selector belongs to a dropdown
		for _, currentDropdownEntry := range Dropdowns.GetAllEntries(characterEntry.LayerAlias) {
			dropdownEntry := currentDropdownEntry
			if dropdownEntry.SelectorAlias == characterEntry.AttributeEntry.CellControlAlias {
				// If it belongs to a dropdown, set the dropdown as the focused control
				setFocusedControl(characterEntry.LayerAlias, dropdownEntry.Alias, constants.CellTypeDropdown)
				setPreviouslyHighlightedControl(characterEntry.LayerAlias, dropdownEntry.Alias, constants.CellTypeDropdown)
				isScreenUpdateRequired = true
				return isScreenUpdateRequired
			}
		}
		// If not part of a dropdown, set the selector as the focused control
		setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeSelectorItem)
		setPreviouslyHighlightedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeSelectorItem)
		isScreenUpdateRequired = true
	} else {
		if eventStateMemory.previouslyHighlightedControl.controlType == constants.CellTypeSelectorItem && Selectors.IsExists(eventStateMemory.previouslyHighlightedControl.layerAlias, eventStateMemory.previouslyHighlightedControl.controlAlias) &&
			Selectors.IsExists(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias) {
			selectorEntry := Selectors.Get(eventStateMemory.previouslyHighlightedControl.layerAlias, eventStateMemory.previouslyHighlightedControl.controlAlias)
			selectorEntry.ItemHighlighted = constants.NullItemSelection
			setFocusedControl("", "", constants.NullControlType)
			setPreviouslyHighlightedControl("", "", constants.NullControlType)
			isScreenUpdateRequired = true
		}
	}

	// --- SCROLL BAR SYNC CODE ---
	layerAlias := characterEntry.LayerAlias

	// If a buttonType is pressed AND (you are in a drag and drop event OR the cell type is scroll bar), then
	// sync all Dropdown selectors with their appropriate scroll bars. If the control under focus
	// matches a control that belongs to a Dropdown list, then stop processing (Do not attempt to close Dropdown).
	if buttonPressed != 0 && (eventStateMemory.stateId == constants.EventStateDragAndDropScrollbar ||
		characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar) {
		for _, currentSelectorEntry := range Selectors.GetAllEntries(focusedLayerAlias) {
			selectorEntry := currentSelectorEntry
			// TODO: Here we don't need to protect this since it is not user controlled?
			scrollBarEntry := ScrollBars.Get(focusedLayerAlias, selectorEntry.ScrollbarAlias)
			if selectorEntry.ViewportPosition != scrollBarEntry.ScrollValue {
				selectorEntry.ViewportPosition = scrollBarEntry.ScrollValue
				isScreenUpdateRequired = true
			}
		}
	}
	// If a Selector is no longer visible, then make the scroll bars associated with it invisible as well.
	for _, currentSelectorEntry := range Selectors.GetAllEntries(layerAlias) {
		selectorEntry := currentSelectorEntry
		// TODO: Here we don't need to protect this since it is not user controlled?
		scrollBarEntry := ScrollBars.Get(layerAlias, selectorEntry.ScrollbarAlias)
		if !selectorEntry.IsVisible {
			scrollBarEntry.IsVisible = false
		} else {
			if scrollBarEntry.IsEnabled {
				scrollBarEntry.IsVisible = true
			}
		}
	}
	return isScreenUpdateRequired
}
