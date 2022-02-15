package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
)

type DropdownInstanceType struct {
	layerAlias  string
	dropdownAlias string
}

func AddDropdown(layerAlias string, dropdownAlias string, styleEntry memory.TuiStyleEntryType, selectionEntry memory.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, itemSelected int) DropdownInstanceType {
	memory.AddDropdown(layerAlias, dropdownAlias, styleEntry, selectionEntry, xLocation, yLocation, itemWidth, itemSelected)
	dropdownEntry := memory.GetDropdown(layerAlias, dropdownAlias)
	dropdownEntry.ScrollBarAlias = stringformat.GetLastSortedUUID()
	// Here we add +2 to x to account for the scroll bar being outside the selector border on ether side. Also, we
	// minus the scroll bar max size by the height of the selector so we don't scroll over values which do not change viewport.
	selectorWidth := itemWidth
	if len(selectionEntry.SelectionValue) <= selectorHeight {
		selectorWidth = selectorWidth + 1
	}
	dropdownEntry.SelectorAlias = stringformat.GetLastSortedUUID()
	// Here we add +1 to x and y to account for borders around the selection.
	AddSelector(layerAlias, dropdownEntry.SelectorAlias, styleEntry, selectionEntry, xLocation + 1, yLocation + 1, selectorHeight, selectorWidth, 1, 0, 0, true)
	selectorEntry := memory.SelectorMemory[layerAlias][dropdownEntry.SelectorAlias]
	selectorEntry.IsVisible = false
	dropdownEntry.ScrollBarAlias = selectorEntry.ScrollBarAlias
	scrollBarEntry := memory.ScrollBarMemory[layerAlias][dropdownEntry.ScrollBarAlias]
	scrollBarEntry.IsVisible = false
	if len(selectionEntry.SelectionValue) <= selectorHeight {
		scrollBarEntry.IsEnabled = false
	}
	var dropdownInstance DropdownInstanceType
	dropdownInstance.layerAlias = layerAlias
	dropdownInstance.dropdownAlias = dropdownAlias
	return dropdownInstance
}

/*
drawButtonsOnLayer allows you to draw all buttons on a given text layer
entry.
*/
func drawDropdownsOnLayer(layerEntry memory.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for currentKey := range memory.DropdownMemory[layerAlias] {
		drawDropdown(&layerEntry, currentKey)
	}
}

func drawDropdown(layerEntry *memory.LayerEntryType, dropdownAlias string) {
	layerAlias := layerEntry.LayerAlias
	dropdownEntry := memory.DropdownMemory[layerAlias][dropdownAlias]
	localStyleEntry := memory.NewTuiStyleEntry(&dropdownEntry.StyleEntry)
	attributeEntry := memory.NewAttributeEntry()
	attributeEntry.ForegroundColor = localStyleEntry.SelectorForegroundColor
	attributeEntry.BackgroundColor = localStyleEntry.SelectorBackgroundColor
	attributeEntry.CellType = constants.CellTypeDropdown
	attributeEntry.CellAlias = dropdownAlias
	itemSelected := dropdownEntry.SelectionEntry.SelectionValue[dropdownEntry.ItemSelected]
	// We add +2 to account for the dropdown border window which will appear. Otherwise, the item name
	// will appear 2 characters smaller than the popup dropdown window.
	formattedItemName := stringformat.GetFormattedString(itemSelected, dropdownEntry.ItemWidth + 2, localStyleEntry.SelectorTextAlignment)
	arrayOfRunes := stringformat.GetRunesFromString(formattedItemName)
	printLayer(layerEntry, attributeEntry, dropdownEntry.XLocation, dropdownEntry.YLocation, arrayOfRunes)
	attributeEntry.ForegroundColor = localStyleEntry.SelectorBackgroundColor
	attributeEntry.BackgroundColor = localStyleEntry.SelectorForegroundColor
	printLayer(layerEntry, attributeEntry, dropdownEntry.XLocation + len(arrayOfRunes), dropdownEntry.YLocation, []rune{constants.CharTriangleDown})
}

func updateDropdownStateMouse() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	cellAlias := characterEntry.AttributeEntry.CellAlias

	// If a button is pressed AND (you are in a drag and drop event OR the cell type is scroll bar), then
	// sync all dropdown selectors with their appropriate scroll bars. If the control under focus
	// matches a control that belongs to a dropdown list, then stop processing (Do not attempt to close dropdown).
	if buttonPressed != 0 && (eventStateMemory.stateId == constants.EventStateDragAndDropScrollBar ||
		characterEntry.AttributeEntry.CellType == constants.CellTypeScrollBar) {
		isMatchFound := false
		for currentKey := range memory.DropdownMemory[layerAlias] {
			dropdownEntry := memory.DropdownMemory[layerAlias][currentKey]
			selectorEntry := memory.SelectorMemory[layerAlias][dropdownEntry.SelectorAlias]
			scrollBarEntry := memory.ScrollBarMemory[layerAlias][dropdownEntry.ScrollBarAlias]
			if selectorEntry.ViewportPosition != scrollBarEntry.ScrollValue {
				selectorEntry.ViewportPosition = scrollBarEntry.ScrollValue
				isUpdateRequired = true
			}
			if isControlFocusMatch(layerAlias, dropdownEntry.ScrollBarAlias, constants.CellTypeScrollBar) {
				isMatchFound = true
				break; // If the current scrollbar being dragged and dropped matches, don't process more dropdowns.
			}
		}
		if isMatchFound {
			return isUpdateRequired
		}
	}

	// If our dropdown alias is not empty, then open our dropdown.
	if buttonPressed != 0 && cellAlias != "" && characterEntry.AttributeEntry.CellType == constants.CellTypeDropdown {
		closeAllOpenDropdowns(layerAlias)
		dropdownEntry := memory.DropdownMemory[layerAlias][cellAlias]
		dropdownEntry.IsTrayOpen = true
		selectorEntry := memory.SelectorMemory[layerAlias][dropdownEntry.SelectorAlias]
		selectorEntry.IsVisible = true
		scrollBarEntry := memory.ScrollBarMemory[layerAlias][dropdownEntry.ScrollBarAlias]
		if scrollBarEntry.IsEnabled {
			scrollBarEntry.IsVisible = true
		}
		isUpdateRequired = true
		return isUpdateRequired
	}

	if buttonPressed != 0 && characterEntry.AttributeEntry.CellType != constants.CellTypeDropdown {
		closeAllOpenDropdowns(layerAlias)
	}
	return isUpdateRequired
}

func closeAllOpenDropdowns(layerAlias string) {
	for currentKey := range memory.DropdownMemory[layerAlias] {
		dropdownEntry := memory.DropdownMemory[layerAlias][currentKey]
		if dropdownEntry.IsTrayOpen == true {
			selectorEntry := memory.SelectorMemory[layerAlias][dropdownEntry.SelectorAlias]
			selectorEntry.IsVisible = false
			scrollBarEntry := memory.ScrollBarMemory[layerAlias][dropdownEntry.ScrollBarAlias]
			scrollBarEntry.IsVisible = false
			dropdownEntry.IsTrayOpen = false
			if dropdownEntry.ItemSelected != selectorEntry.ItemSelected {
				dropdownEntry.ItemSelected = selectorEntry.ItemSelected
			}
			setFocusedControl("", "", constants.NullCellType)
		}
	}
}