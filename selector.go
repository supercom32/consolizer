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

type selectorType struct {}
var Selector selectorType

/*
GetSelected allows you to retrieve the currently selected item. If the Selector instance
no longer exists, then an empty result is always returned.
*/
func (shared *selectorInstanceType) GetSelected() string {
	if memory.IsSelectorExists(shared.layerAlias, shared.selectorAlias) {
		validatorMenu(shared.layerAlias, shared.selectorAlias)
		menuEntry := memory.GetSelector(shared.layerAlias, shared.selectorAlias)
		value := menuEntry.ItemSelected
		return menuEntry.SelectionEntry.SelectionAlias[value]
	}
	return ""
}

/*
setViewport allows you to specify the current viewport index for a given selector. If the selector instance
no longer exists, then no operation occurs.
*/
func (shared *selectorInstanceType) setViewport(viewportPosition int) {
	if memory.IsSelectorExists(shared.layerAlias, shared.selectorAlias) {
		validatorMenu(shared.layerAlias, shared.selectorAlias)
		menuEntry := memory.GetSelector(shared.layerAlias, shared.selectorAlias)
		menuEntry.ViewportPosition = viewportPosition
	}
}

/*
AddSelector allows you to add a Selector to a given text layer. Once called, an instance of your control is returned
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
func (shared *selectorType) AddSelector(layerAlias string, selectorAlias string, styleEntry memory.TuiStyleEntryType, selectionEntry memory.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, selectedItem int, isBorderDrawn bool) selectorInstanceType {
	validateSelectionEntry(selectionEntry)
	// TODO: Add verification to ensure no item can be 0 length/number.
	memory.AddSelector(layerAlias, selectorAlias, styleEntry, selectionEntry, xLocation, yLocation, selectorHeight, itemWidth, numberOfColumns, viewportPosition, selectedItem, isBorderDrawn)
	selectorEntry := memory.GetSelector(layerAlias, selectorAlias)
	selectorEntry.ScrollBarAlias = stringformat.GetLastSortedUUID()
	scrollBarMaxValue := len(selectionEntry.SelectionValue) - (selectorHeight * numberOfColumns) + 1
	scrollBarXLocation := xLocation + (itemWidth * numberOfColumns)
	scrollBarYLocation := yLocation
	scrollBarHeight := selectorHeight
	if isBorderDrawn {
		scrollBarXLocation = xLocation + (itemWidth * numberOfColumns) + 1
		scrollBarYLocation = scrollBarYLocation - 1
		scrollBarHeight = selectorHeight + 2
	}
	scrollbar.AddScrollbar(layerAlias, selectorEntry.ScrollBarAlias, styleEntry, scrollBarXLocation, scrollBarYLocation, scrollBarHeight, scrollBarMaxValue, 0, numberOfColumns,false)
	scrollBarEntry := memory.GetScrollbar(layerAlias, selectorEntry.ScrollBarAlias)
	selectorWidth := itemWidth
	if len(selectionEntry.SelectionValue) <= selectorHeight * numberOfColumns || styleEntry.SelectorTextAlignment == constants.AlignmentNoPadding {
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
DrawSelector allows you to draw a Selector on a given text layer. The
Style of the Selector will be determined by the style entry passed in. In
addition, the following information should be noted:

- Selectors are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create selectors without actually overwriting
the text layer data under it.

- If the Selector to be drawn falls outside the range of the provided layer,
then only the visible portion of the Selector will be drawn.
*/
func (shared *selectorType) DrawSelector(selectorAlias string, layerEntry *memory.LayerEntryType, styleEntry memory.TuiStyleEntryType, selectionEntry memory.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, itemHighlighted int) {
	selectorEntry := memory.GetSelector(layerEntry.LayerAlias, selectorAlias)
	if selectorEntry.IsVisible == false {
		return
	}
	menuAttributeEntry := memory.NewAttributeEntry()
	menuAttributeEntry.ForegroundColor = styleEntry.SelectorForegroundColor
	menuAttributeEntry.BackgroundColor = styleEntry.SelectorBackgroundColor
	if selectorEntry.IsBorderDrawn {
		fillArea(layerEntry, menuAttributeEntry," ", xLocation - 1,yLocation - 1, itemWidth + 2, selectorHeight + 2, constants.NullCellControlLocation)
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
		attributeEntry.CellType = constants.CellTypeSelectorItem
		printLayer(layerEntry, attributeEntry, xLocation + (currentXOffset), currentYLocation, arrayOfRunes)
		currentMenuItemIndex++
		currentXOffset = currentXOffset + stringformat.GetWidthOfRunesWhenPrinted(arrayOfRunes)//len(arrayOfRunes)
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
drawSelectorsOnLayer allows you to draw all selectors on a given text layer.
*/
func (shared *selectorType) drawSelectorsOnLayer(layerEntry memory.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	// Range over all our Selector aliases and add them to an array.
	keyList := make([]string, 0)
	for currentKey := range memory.SelectorMemory[layerAlias] {
		keyList = append(keyList, currentKey)
	}
	// Sort array so internally generated selectors appear last (Since sorted by name, and
	// UUID generates "zzz" prefixes). This prevents dropdown selectors from appearing under
	// user created selectors, when they should always be on top.
	sort.Strings(keyList)
	for currentKey := range keyList {
		selectorEntry := memory.GetSelector(layerAlias, keyList[currentKey])
		shared.DrawSelector(keyList[currentKey], &layerEntry, selectorEntry.StyleEntry, selectorEntry.SelectionEntry, selectorEntry.XLocation, selectorEntry.YLocation, selectorEntry.SelectorHeight, selectorEntry.ItemWidth, selectorEntry.NumberOfColumns, selectorEntry.ViewportPosition, selectorEntry.ItemHighlighted)
	}
}

/*
updateKeyboardEventSelector allows you to update the state of all selectors according to the current keystroke event.
In the event that a screen update is required this method returns true.
*/
func (shared *selectorType) updateKeyboardEventSelector(keystroke []rune) bool {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := false
	if eventStateMemory.currentlyFocusedControl.controlType != constants.CellTypeSelectorItem {
		return isScreenUpdateRequired
	}
	selectorEntry := memory.GetSelector(eventStateMemory.currentlyFocusedControl.layerAlias, eventStateMemory.currentlyFocusedControl.controlAlias)
	if keystrokeAsString == "down" {
		//remainder := selectorEntry.ItemHighlighted % selectorEntry.NumberOfColumns
		selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + selectorEntry.NumberOfColumns
		if selectorEntry.ItemHighlighted >= len(selectorEntry.SelectionEntry.SelectionAlias) {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - selectorEntry.NumberOfColumns
		}
		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "up" {
		selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - selectorEntry.NumberOfColumns
		if selectorEntry.ItemHighlighted < 0 {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + selectorEntry.NumberOfColumns
		}
		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "left" {
		if selectorEntry.ItemHighlighted % selectorEntry.NumberOfColumns != 0 {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - 1
			if selectorEntry.ItemHighlighted < 0 {
				selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + 1
			}
			isScreenUpdateRequired = true
		}
	}
	if keystrokeAsString == "right" {
		if selectorEntry.ItemHighlighted % selectorEntry.NumberOfColumns != selectorEntry.NumberOfColumns - 1 {
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
updateMouseEventSelector allows you to update the state of all selectors according to the current mouse event state.
In the event that a screen update is required this method returns true.
*/
func (shared *selectorType) updateMouseEventSelector() bool {
	isScreenUpdateRequired := false
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	var characterEntry memory.CharacterEntryType
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	if characterEntry.AttributeEntry.CellType == constants.CellTypeSelectorItem && eventStateMemory.stateId == constants.EventStateNone {
		if buttonPressed != 0 {
			selectorEntry := memory.GetSelector(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			selectorEntry.ItemHighlighted = characterEntry.AttributeEntry.CellControlId
			selectorEntry.ItemSelected = characterEntry.AttributeEntry.CellControlId
		} else {
			selectorEntry := memory.GetSelector(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			selectorEntry.ItemHighlighted = characterEntry.AttributeEntry.CellControlId
		}
		setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeSelectorItem)
		setPreviouslyHighlightedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeSelectorItem)
		isScreenUpdateRequired = true
	} else {
		if eventStateMemory.previouslyHighlightedControl.controlType == constants.CellTypeSelectorItem && memory.IsSelectorExists(eventStateMemory.previouslyHighlightedControl.layerAlias, eventStateMemory.previouslyHighlightedControl.controlAlias) {
			selectorEntry := memory.GetSelector(eventStateMemory.previouslyHighlightedControl.layerAlias, eventStateMemory.previouslyHighlightedControl.controlAlias)
			selectorEntry.ItemHighlighted = constants.NullItemSelection
			setFocusedControl("", "", constants.NullControlType)
			setPreviouslyHighlightedControl("", "", constants.NullControlType)
			isScreenUpdateRequired = true
		}
	}

	// --- SCROLL BAR SYNC CODE ---
	layerAlias := characterEntry.LayerAlias

	// If a buttonType is pressed AND (you are in a drag and drop event OR the cell type is scroll bar), then
	// sync all dropdown selectors with their appropriate scroll bars. If the control under focus
	// matches a control that belongs to a dropdown list, then stop processing (Do not attempt to close dropdown).
	if buttonPressed != 0 && (eventStateMemory.stateId == constants.EventStateDragAndDropScrollbar ||
		characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar) {
		for currentKey := range memory.SelectorMemory[focusedLayerAlias] {
			selectorEntry := memory.GetSelector(focusedLayerAlias, currentKey)
			scrollBarEntry := memory.GetScrollbar(focusedLayerAlias, selectorEntry.ScrollBarAlias)
			if selectorEntry.ViewportPosition != scrollBarEntry.ScrollValue {
				selectorEntry.ViewportPosition = scrollBarEntry.ScrollValue
				isScreenUpdateRequired = true
			}
		}
	}
	// If a Selector is no longer visible, then make the scroll bars associated with it invisible as well.
	for currentKey := range memory.SelectorMemory[layerAlias] {
		selectorEntry := memory.GetSelector(layerAlias, currentKey)
		scrollBarEntry := memory.GetScrollbar(layerAlias, selectorEntry.ScrollBarAlias)
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