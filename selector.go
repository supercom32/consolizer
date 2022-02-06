package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
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


func AddSelector(layerAlias string, selectorAlias string, styleEntry memory.TuiStyleEntryType, selectionEntry memory.SelectionEntryType, xLocation int, yLocation int, itemWidth int, numberOfColumns int, selectedItem int) selectorInstanceType {
	validateLayerLocationByLayerAlias(layerAlias, xLocation, yLocation)
	// TODO: Add verification to ensure no item can be 0 length/number.
	memory.AddSelector(layerAlias, selectorAlias, styleEntry, selectionEntry, xLocation, yLocation, itemWidth, numberOfColumns, selectedItem)
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
func DrawSelector(selectorAlias string, layerEntry *memory.LayerEntryType, styleEntry memory.TuiStyleEntryType, selectionEntry memory.SelectionEntryType, xLocation int, yLocation int, itemWidth int, numberOfColumns int, itemHighlighted int) {
	menuAttributeEntry := memory.NewAttributeEntry()
	menuAttributeEntry.ForegroundColor = styleEntry.MenuForegroundColor
	menuAttributeEntry.BackgroundColor = styleEntry.MenuBackgroundColor
	highlightAttributeEntry := memory.NewAttributeEntry()
	highlightAttributeEntry.ForegroundColor = styleEntry.HighlightForegroundColor
	highlightAttributeEntry.BackgroundColor = styleEntry.HighlightBackgroundColor
	currentYLocation := yLocation
	currentMenuItemIndex := 0
	currentColumn := 0
	for currentMenuItemIndex < len(selectionEntry.SelectionValue) {
		attributeEntry := menuAttributeEntry
		if currentMenuItemIndex == itemHighlighted {
			attributeEntry = highlightAttributeEntry
		}
		menuItemName := stringformat.GetFormattedString(selectionEntry.SelectionValue[currentMenuItemIndex], itemWidth, styleEntry.MenuTextAlignment)
		arrayOfRunes := stringformat.GetRunesFromString(menuItemName)
		attributeEntry.CellControlId = currentMenuItemIndex
		attributeEntry.CellControlAlias = selectorAlias
		attributeEntry.CellType = constants.CellTypeMenuItem
		printLayer(layerEntry, attributeEntry, xLocation + (currentColumn * itemWidth), currentYLocation, arrayOfRunes)
		currentMenuItemIndex++
		currentColumn++
		if currentColumn >= numberOfColumns {
			currentColumn = 0
			currentYLocation++
		}
	}
}

func drawSelectorsOnLayer(layerEntry memory.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for currentKey := range memory.SelectorMemory[layerAlias] {
		selectorEntry := memory.SelectorMemory[layerAlias][currentKey]
		DrawSelector(currentKey, &layerEntry, selectorEntry.StyleEntry, selectorEntry.SelectionEntry, selectorEntry.XLocation, selectorEntry.YLocation, selectorEntry.ItemWidth, selectorEntry.NumberOfColumns, selectorEntry.ItemHighlighted)
	}
}

func updateKeyboardEventSelector(keystroke string) bool {
	isScreenUpdateRequired := true
	if eventStateMemory.focusedControlType != constants.CellTypeMenuItem {
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
	if characterEntry.AttributeEntry.CellType == constants.CellTypeMenuItem {
		if buttonPressed != 0 {
			selectorEntry := memory.GetSelector(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			selectorEntry.ItemHighlighted = characterEntry.AttributeEntry.CellControlId
			selectorEntry.ItemSelected = characterEntry.AttributeEntry.CellControlId
			eventStateMemory.focusedControlAlias = characterEntry.AttributeEntry.CellControlAlias
			eventStateMemory.focusedLayerAlias = characterEntry.LayerAlias
			eventStateMemory.focusedControlType = constants.CellTypeMenuItem
			eventStateMemory.previousHighlightedLayerAlias = characterEntry.LayerAlias
			eventStateMemory.previousHighlightedControlAlias = characterEntry.AttributeEntry.CellControlAlias
			eventStateMemory.previousHighlightedControlType = constants.CellTypeMenuItem
			isScreenUpdateRequired = true
		} else {
			selectorEntry := memory.GetSelector(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			selectorEntry.ItemHighlighted = characterEntry.AttributeEntry.CellControlId
			eventStateMemory.focusedControlAlias = characterEntry.AttributeEntry.CellControlAlias
			eventStateMemory.focusedLayerAlias = characterEntry.LayerAlias
			eventStateMemory.focusedControlType = constants.CellTypeMenuItem
			eventStateMemory.previousHighlightedLayerAlias = characterEntry.LayerAlias
			eventStateMemory.previousHighlightedControlAlias = characterEntry.AttributeEntry.CellControlAlias
			eventStateMemory.previousHighlightedControlType = constants.CellTypeMenuItem
			isScreenUpdateRequired = true
		}
	} else {
		if eventStateMemory.previousHighlightedControlType == constants.CellTypeMenuItem && memory.IsSelectorExists(eventStateMemory.previousHighlightedLayerAlias, eventStateMemory.previousHighlightedControlAlias) {
			selectorEntry := memory.GetSelector(characterEntry.LayerAlias, eventStateMemory.previousHighlightedControlAlias)
			selectorEntry.ItemHighlighted = constants.NullItemSelection
			eventStateMemory.focusedControlAlias = ""
			eventStateMemory.focusedLayerAlias = ""
			eventStateMemory.focusedControlType = constants.NullControlType
			eventStateMemory.previousHighlightedLayerAlias = ""
			eventStateMemory.previousHighlightedControlAlias = ""
			eventStateMemory.previousHighlightedControlType = 0
			isScreenUpdateRequired = true
		}
	}
	return isScreenUpdateRequired
}