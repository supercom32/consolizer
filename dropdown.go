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

type dropdownType struct {}
var dropdown dropdownType


/*
AddDropdown allows you to add a dropdown to a given text layer. Once called, an instance of
your control is returned which will allow you to read or manipulate the properties for it.
The Style of the dropdown will be determined by the style entry passed in. If you wish to
remove a dropdown from a text layer, simply call 'DeleteDropdown'. In addition, the
following information should be noted:

- Dropdowns are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create dropdowns without actually overwriting
the text layer data under it.

- If the dropdown to be drawn falls outside the range of the provided layer,
then only the visible portion of the checkbox will be drawn.

- If the number of selections available is smaller or equal to the selector height,
then no scrollbars will be drawn.
*/
func (shared *dropdownType) AddDropdown(layerAlias string, dropdownAlias string, styleEntry memory.TuiStyleEntryType, selectionEntry memory.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, defaultItemSelected int) DropdownInstanceType {
	// TODO: Add validation to the default item selected.
	memory.AddDropdown(layerAlias, dropdownAlias, styleEntry, selectionEntry, xLocation, yLocation, itemWidth, defaultItemSelected)
	dropdownEntry := memory.GetDropdown(layerAlias, dropdownAlias)
	dropdownEntry.ScrollBarAlias = stringformat.GetLastSortedUUID()
	// Here we add +2 to x to account for the scroll bar being outside the selector border on ether side. Also, we
	// minus the scroll bar max selection size by the height of the selector, so we don't scroll over values
	// which do not change viewport.
	selectorWidth := itemWidth
	if len(selectionEntry.SelectionValue) <= selectorHeight {
		selectorWidth = selectorWidth + 1
	}
	dropdownEntry.SelectorAlias = stringformat.GetLastSortedUUID()
	// Here we add +1 to x and y to account for borders around the selection.
	selector.AddSelector(layerAlias, dropdownEntry.SelectorAlias, styleEntry, selectionEntry, xLocation + 1, yLocation + 1, selectorHeight, selectorWidth, 1, 0, 0, true)
	selectorEntry := memory.GetSelector(layerAlias, dropdownEntry.SelectorAlias)
	selectorEntry.IsVisible = false
	dropdownEntry.ScrollBarAlias = selectorEntry.ScrollBarAlias
	scrollBarEntry := memory.GetScrollbar(layerAlias, dropdownEntry.ScrollBarAlias)
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
drawDropdownsOnLayer allows you to draw all dropdowns on a given text layer.
*/
func (shared *dropdownType) drawDropdownsOnLayer(layerEntry memory.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for currentKey := range memory.DropdownMemory[layerAlias] {
		shared.drawDropdown(&layerEntry, currentKey)
	}
}

/*
drawDropdown allows you to draw A dropdown on a given text layer. The
Style of the dropdown will be determined by the style entry passed in. In
addition, the following information should be noted:

- dropdowns are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create dropdowns without actually overwriting
the text layer data under it.

- If the dropdown to be drawn falls outside the range of the provided layer,
then only the visible portion of the dropdown will be drawn.
*/
func (shared *dropdownType) drawDropdown(layerEntry *memory.LayerEntryType, dropdownAlias string) {
	layerAlias := layerEntry.LayerAlias
	dropdownEntry := memory.GetDropdown(layerAlias, dropdownAlias)
	localStyleEntry := memory.NewTuiStyleEntry(&dropdownEntry.StyleEntry)
	attributeEntry := memory.NewAttributeEntry()
	attributeEntry.ForegroundColor = localStyleEntry.SelectorForegroundColor
	attributeEntry.BackgroundColor = localStyleEntry.SelectorBackgroundColor
	attributeEntry.CellType = constants.CellTypeDropdown
	attributeEntry.CellControlAlias = dropdownAlias
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

/*
updateDropdownStateMouse allows you to update the state of all dropdowns according to the current mouse event state.
In the event that a screen update is required this method returns true.
*/
func (shared *dropdownType) updateDropdownStateMouse() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	cellControlAlias := characterEntry.AttributeEntry.CellControlAlias
	// If a buttonType is pressed AND (you are in a drag and drop event OR the cell type is scroll bar), then
	// sync all dropdown selectors with their appropriate scroll bars. If the control under focus
	// matches a control that belongs to a dropdown list, then stop processing (Do not attempt to close dropdown).
	if buttonPressed != 0 && (eventStateMemory.stateId == constants.EventStateDragAndDropScrollbar ||
		characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar) {
		isMatchFound := false
		for currentKey := range memory.DropdownMemory[layerAlias] {
			dropdownEntry := memory.GetDropdown(layerAlias, currentKey)
			selectorEntry := memory.GetSelector(layerAlias, dropdownEntry.SelectorAlias)
			scrollBarEntry := memory.GetScrollbar(layerAlias, dropdownEntry.ScrollBarAlias)
			if selectorEntry.ViewportPosition != scrollBarEntry.ScrollValue {
				selectorEntry.ViewportPosition = scrollBarEntry.ScrollValue
				isUpdateRequired = true
			}
			if isControlCurrentlyFocused(layerAlias, dropdownEntry.ScrollBarAlias, constants.CellTypeScrollbar) {
				isMatchFound = true
				break; // If the current scrollbar being dragged and dropped matches, don't process more dropdowns.
			}
		}
		if isMatchFound {
			return isUpdateRequired
		}
	}
	// If our dropdown alias is not empty, then open our dropdown.
	if buttonPressed != 0 && cellControlAlias != "" && characterEntry.AttributeEntry.CellType == constants.CellTypeDropdown {
		shared.closeAllOpenDropdowns(layerAlias)
		dropdownEntry := memory.GetDropdown(layerAlias, cellControlAlias)
		dropdownEntry.IsTrayOpen = true
		selectorEntry := memory.GetSelector(layerAlias, dropdownEntry.SelectorAlias)
		selectorEntry.IsVisible = true
		scrollBarEntry := memory.GetScrollbar(layerAlias, dropdownEntry.ScrollBarAlias)
		if scrollBarEntry.IsEnabled {
			scrollBarEntry.IsVisible = true
			setFocusedControl(layerAlias, selectorEntry.ScrollBarAlias, constants.CellTypeScrollbar)
		}
		isUpdateRequired = true
		return isUpdateRequired
	}
	_, _, previousButtonPress, _ := memory.GetPreviousMouseStatus()
	if buttonPressed != 0 && previousButtonPress == 0 && characterEntry.AttributeEntry.CellType != constants.CellTypeDropdown {
		shared.closeAllOpenDropdowns(layerAlias)
	}
	return isUpdateRequired
}

/*
closeAllOpenDropdowns allows you to close all dropdowns for a given layer alias.
*/
func (shared *dropdownType) closeAllOpenDropdowns(layerAlias string) {
	for currentKey := range memory.DropdownMemory[layerAlias] {
		dropdownEntry := memory.GetDropdown(layerAlias, currentKey)
		if dropdownEntry.IsTrayOpen == true {
			selectorEntry := memory.GetSelector(layerAlias, dropdownEntry.SelectorAlias)
			selectorEntry.IsVisible = false
			scrollBarEntry := memory.GetScrollbar(layerAlias, dropdownEntry.ScrollBarAlias)
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