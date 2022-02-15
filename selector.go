package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
	"sort"
)

type selectorInstanceType struct {
	layerAlias    string
	selectorAlias string
}

/*
GetValue allows you to get the current value of your text field with.
*/
func (shared *selectorInstanceType) GetSelected() string {
	validatorMenu(shared.layerAlias, shared.selectorAlias)
	menuEntry := memory.SelectorMemory[shared.layerAlias][shared.selectorAlias]
	value := menuEntry.ItemSelected
	return menuEntry.SelectionEntry.SelectionAlias[value]
}

func (shared *selectorInstanceType) setViewport(viewportPosition int) {
	validatorMenu(shared.layerAlias, shared.selectorAlias)
	menuEntry := memory.SelectorMemory[shared.layerAlias][shared.selectorAlias]
	menuEntry.ViewportPosition = viewportPosition
}


// TODO: Protect against viewport out of range errors.
func AddSelector(layerAlias string, selectorAlias string, styleEntry memory.TuiStyleEntryType, selectionEntry memory.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, selectedItem int, isBorderDrawn bool) selectorInstanceType {
	validateLayerLocationByLayerAlias(layerAlias, xLocation, yLocation)
	validateSelectionEntry(selectionEntry)
	// TODO: Add verification to ensure no item can be 0 length/number.
	memory.AddSelector(layerAlias, selectorAlias, styleEntry, selectionEntry, xLocation, yLocation, selectorHeight, itemWidth, numberOfColumns, viewportPosition, selectedItem, isBorderDrawn)
	selectorEntry := memory.SelectorMemory[layerAlias][selectorAlias]
	selectorEntry.ScrollBarAlias = stringformat.GetLastSortedUUID()
	scrollBarMaxValue := len(selectionEntry.SelectionValue) - (selectorHeight * numberOfColumns) + 1
	scrollBarXLocation := xLocation + (itemWidth * numberOfColumns)
	scrollBarYLocation := yLocation
	scrollBarHeight := selectorHeight - 2
	if isBorderDrawn {
		scrollBarXLocation = xLocation + (itemWidth * numberOfColumns) + 1
		scrollBarYLocation = scrollBarYLocation - 1
		scrollBarHeight = selectorHeight
	}
	AddScrollBar(layerAlias, selectorEntry.ScrollBarAlias, styleEntry, scrollBarXLocation, scrollBarYLocation, scrollBarHeight, scrollBarMaxValue, 0, numberOfColumns,false)
	scrollBarEntry := memory.ScrollBarMemory[layerAlias][selectorEntry.ScrollBarAlias]
	selectorWidth := itemWidth
	if len(selectionEntry.SelectionValue) <= selectorHeight * numberOfColumns {
		scrollBarEntry.IsEnabled = false
		scrollBarEntry.IsVisible = false
		selectorWidth = selectorWidth + 1
	}
	var selectorInstance selectorInstanceType
	selectorInstance.layerAlias = layerAlias
	selectorInstance.selectorAlias = selectorAlias
	return selectorInstance
}

/*
DrawSelector allows you to obtain a user selection from a horizontal menu.
In addition, the following information should be noted:

- If the location to draw a menu item falls outside the range of the text
layer, then only the visible portion of your menu item will be drawn.
*/
func DrawSelector(selectorAlias string, layerEntry *memory.LayerEntryType, styleEntry memory.TuiStyleEntryType, selectionEntry memory.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, itemHighlighted int) {
	selectorEntry := memory.SelectorMemory[layerEntry.LayerAlias][selectorAlias]
	if selectorEntry.IsVisible == false {
		return
	}
	menuAttributeEntry := memory.NewAttributeEntry()
	menuAttributeEntry.ForegroundColor = styleEntry.SelectorForegroundColor
	menuAttributeEntry.BackgroundColor = styleEntry.SelectorBackgroundColor
	if selectorEntry.IsBorderDrawn {
		fillArea(layerEntry, menuAttributeEntry," ", xLocation - 1,yLocation - 1, itemWidth + 2, selectorHeight + 2)
		drawBorder(layerEntry, styleEntry, menuAttributeEntry, xLocation -1, yLocation - 1, itemWidth + 2, selectorHeight + 2, false)
	}
	highlightAttributeEntry := memory.NewAttributeEntry()
	highlightAttributeEntry.ForegroundColor = styleEntry.HighlightForegroundColor
	highlightAttributeEntry.BackgroundColor = styleEntry.HighlightBackgroundColor
	currentYLocation := yLocation
	currentMenuItemIndex := viewportPosition
	currentXOffset := 0
	currentColumn := 0
	currentRow :=0
	for currentMenuItemIndex < len(selectionEntry.SelectionValue) && currentRow < selectorHeight {
		attributeEntry := menuAttributeEntry
		if currentMenuItemIndex == itemHighlighted {
			attributeEntry = highlightAttributeEntry
		}
		menuItemName := stringformat.GetFormattedString(selectionEntry.SelectionValue[currentMenuItemIndex], itemWidth, styleEntry.SelectorTextAlignment)
		arrayOfRunes := stringformat.GetRunesFromString(menuItemName)
		attributeEntry.CellControlId = currentMenuItemIndex
		attributeEntry.CellControlAlias = selectorAlias
		attributeEntry.CellType = constants.CellTypeSelectiorItem
		printLayer(layerEntry, attributeEntry, xLocation + (currentXOffset), currentYLocation, arrayOfRunes)
		currentMenuItemIndex++
		currentXOffset = currentXOffset + len(arrayOfRunes)
		currentColumn++
		if currentColumn >= numberOfColumns {
			currentXOffset = 0
			currentColumn = 0
			currentYLocation++
			currentRow++
		}
	}
}

func drawSelectorsOnLayer(layerEntry memory.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	// Range over all our selector aliases and add them to an array.
	keyList := make([]string, 0)
	for currentKey := range memory.SelectorMemory[layerAlias] {
		keyList = append(keyList, currentKey)
	}
	// Sort array so internally generated selectors appear last (Since sorted by name, and
	// UUID generates "zzz" prefixes). This prevents dropdown selectors from appearing under
	// user created selectors, when they should always be on top.
	sort.Strings(keyList)
	for currentKey := range keyList {
		selectorEntry := memory.SelectorMemory[layerAlias][keyList[currentKey]]
		DrawSelector(keyList[currentKey], &layerEntry, selectorEntry.StyleEntry, selectorEntry.SelectionEntry, selectorEntry.XLocation, selectorEntry.YLocation, selectorEntry.SelectorHeight, selectorEntry.ItemWidth, selectorEntry.NumberOfColumns, selectorEntry.ViewportPosition, selectorEntry.ItemHighlighted)
	}
}

func updateKeyboardEventSelector(keystroke string) bool {
	isScreenUpdateRequired := false
	if eventStateMemory.focusedControlType != constants.CellTypeSelectiorItem {
		return isScreenUpdateRequired
	}
	selectorEntry := memory.GetSelector(eventStateMemory.focusedLayerAlias, eventStateMemory.focusedControlAlias)
	if keystroke == "down" {
		//remainder := selectorEntry.ItemHighlighted % selectorEntry.NumberOfColumns
		selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + selectorEntry.NumberOfColumns
		if selectorEntry.ItemHighlighted >= len(selectorEntry.SelectionEntry.SelectionAlias) {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - selectorEntry.NumberOfColumns
		}
		isScreenUpdateRequired = true
	}
	if keystroke == "up" {
		selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - selectorEntry.NumberOfColumns
		if selectorEntry.ItemHighlighted < 0 {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + selectorEntry.NumberOfColumns
		}
		isScreenUpdateRequired = true
	}
	if keystroke == "left" {
		if selectorEntry.ItemHighlighted % selectorEntry.NumberOfColumns != 0 {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - 1
			if selectorEntry.ItemHighlighted < 0 {
				selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + 1
			}
			isScreenUpdateRequired = true
		}
	}
	if keystroke == "right" {
		if selectorEntry.ItemHighlighted % selectorEntry.NumberOfColumns != selectorEntry.NumberOfColumns - 1 {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + 1
			if selectorEntry.ItemHighlighted >= len(selectorEntry.SelectionEntry.SelectionAlias) {
				selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - 1
			}
			isScreenUpdateRequired = true
		}
	}
	if keystroke == "enter" {
		selectorEntry.ItemSelected = selectorEntry.ItemHighlighted
		isScreenUpdateRequired = true
	}
	return isScreenUpdateRequired
}

func updateMouseEventSelector() bool {
	isScreenUpdateRequired := false
	var characterEntry memory.CharacterEntryType
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	if characterEntry.AttributeEntry.CellType == constants.CellTypeSelectiorItem && eventStateMemory.stateId == constants.EventStateNone {
		if buttonPressed != 0 {
			selectorEntry := memory.GetSelector(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			selectorEntry.ItemHighlighted = characterEntry.AttributeEntry.CellControlId
			selectorEntry.ItemSelected = characterEntry.AttributeEntry.CellControlId
			/*
			eventStateMemory.focusedControlAlias = characterEntry.AttributeEntry.CellControlAlias
			eventStateMemory.focusedLayerAlias = characterEntry.LayerAlias
			eventStateMemory.focusedControlType = constants.CellTypeSelectiorItem
			eventStateMemory.previousHighlightedLayerAlias = characterEntry.LayerAlias
			eventStateMemory.previousHighlightedControlAlias = characterEntry.AttributeEntry.CellControlAlias
			eventStateMemory.previousHighlightedControlType = constants.CellTypeSelectiorItem
			 */
			isScreenUpdateRequired = true
		} else {
			selectorEntry := memory.GetSelector(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			selectorEntry.ItemHighlighted = characterEntry.AttributeEntry.CellControlId
			/*
			eventStateMemory.focusedControlAlias = characterEntry.AttributeEntry.CellControlAlias
			eventStateMemory.focusedLayerAlias = characterEntry.LayerAlias
			eventStateMemory.focusedControlType = constants.CellTypeSelectiorItem
			eventStateMemory.previousHighlightedLayerAlias = characterEntry.LayerAlias
			eventStateMemory.previousHighlightedControlAlias = characterEntry.AttributeEntry.CellControlAlias
			eventStateMemory.previousHighlightedControlType = constants.CellTypeSelectiorItem

			 */
			isScreenUpdateRequired = true
		}
	} else {
		if eventStateMemory.previousHighlightedControlType == constants.CellTypeSelectiorItem && memory.IsSelectorExists(eventStateMemory.previousHighlightedLayerAlias, eventStateMemory.previousHighlightedControlAlias) {
			selectorEntry := memory.GetSelector(eventStateMemory.previousHighlightedLayerAlias, eventStateMemory.previousHighlightedControlAlias)
			selectorEntry.ItemHighlighted = constants.NullItemSelection
			/*
			eventStateMemory.focusedControlAlias = ""
			eventStateMemory.focusedLayerAlias = ""
			eventStateMemory.focusedControlType = constants.NullControlType
			eventStateMemory.previousHighlightedLayerAlias = ""
			eventStateMemory.previousHighlightedControlAlias = ""
			eventStateMemory.previousHighlightedControlType = 0

			 */
			isScreenUpdateRequired = true
		}
	}

	// --- SCROLL BAR SYNC CODE ---
	layerAlias := characterEntry.LayerAlias

	// If a button is pressed AND (you are in a drag and drop event OR the cell type is scroll bar), then
	// sync all dropdown selectors with their appropriate scroll bars. If the control under focus
	// matches a control that belongs to a dropdown list, then stop processing (Do not attempt to close dropdown).
	if buttonPressed != 0 && (eventStateMemory.stateId == constants.EventStateDragAndDropScrollBar ||
		characterEntry.AttributeEntry.CellType == constants.CellTypeScrollBar) {
		isMatchFound := false
		for currentKey := range memory.SelectorMemory[layerAlias] {
			selectorEntry := memory.SelectorMemory[layerAlias][currentKey]
			scrollBarEntry := memory.ScrollBarMemory[layerAlias][selectorEntry.ScrollBarAlias]
			if selectorEntry.ViewportPosition != scrollBarEntry.ScrollValue {
				selectorEntry.ViewportPosition = scrollBarEntry.ScrollValue
				isScreenUpdateRequired = true
			}
			if isControlFocusMatch(layerAlias, selectorEntry.ScrollBarAlias, constants.CellTypeScrollBar) {
				isMatchFound = true
				break; // If the current scrollbar being dragged and dropped matches, don't process more dropdowns.
			}
		}
		if isMatchFound {
			return isScreenUpdateRequired
		}
	}
	for currentKey := range memory.SelectorMemory[layerAlias] {
		selectorEntry := memory.SelectorMemory[layerAlias][currentKey]
		scrollBarEntry := memory.ScrollBarMemory[layerAlias][selectorEntry.ScrollBarAlias]
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