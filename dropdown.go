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

  - A disabled or hidden dropdown, or one on a hidden layer, ignores every keystroke and leaves it unconsumed, as a
    disabled button does, so it still reaches the keyboard buffer. Such a dropdown cannot normally have focus, so
    this only guards the moment between it being disabled and focus moving on.

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

	// Only process if an enabled, visible dropdown is focused
	if focusedControl.controlType != constants.CellTypeDropdown || !isControlFocusable(focusedControl) {
		return isScreenUpdateRequired, isKeystrokeConsumed
	}
	dropdownEntry, isFound := dropdownRegistry.Lookup(focusedLayerAlias, focusedControlAlias)
	if !isFound {
		return isScreenUpdateRequired, isKeystrokeConsumed
	}

	// While the tray is open, focus stays on the dropdown and the keys are passed on to the tray.
	if dropdownEntry.IsTrayOpen {
		isScreenUpdateRequired, isKeystrokeConsumed = Selector.updateKeyboardEventForSelector(focusedLayerAlias, dropdownEntry.SelectorAlias, keystroke)
	}

	// Handle Enter key to open/close dropdown
	if keystrokeAsString == "enter" {
		if dropdownEntry.IsTrayOpen {
			// Close dropdown and apply selection
			selectorEntry := selectorRegistry.Get(focusedLayerAlias, dropdownEntry.SelectorAlias)
			scrollBarEntry := scrollbarRegistry.Get(focusedLayerAlias, dropdownEntry.ScrollbarAlias)
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
			shared.openTray(focusedLayerAlias, dropdownEntry)
		}
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true
	}

	// Handle Esc key to close dropdown without changing selection. tcell reports the Escape key as "Esc".
	if keystrokeAsString == "esc" && dropdownEntry.IsTrayOpen {
		shared.closeTrayWithoutCommitting(focusedLayerAlias, dropdownEntry)

		// Reset focus to the dropdown itself
		setFocusedControl(focusedLayerAlias, focusedControlAlias, constants.CellTypeDropdown)
		setEventStateId(constants.EventStateNone)
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true
	}

	return isScreenUpdateRequired, isKeystrokeConsumed
}

/*
openTray is a method which allows you to open a dropdown's tray, closing any other open tray first. The tray's
highlight is placed on the dropdown's selected item, and its scroll bar is shown if the items overflow the tray. A
disabled dropdown is left closed, and no other tray is closed for it, since an open tray would let the user change its
selection. The caller must hold uiStateMutex.

Example:

	Dropdown.openTray("layer1", dropdownEntry)
*/
func (shared *dropdownType) openTray(layerAlias string, dropdownEntry *types.DropdownEntryType) {
	if !dropdownEntry.IsEnabled {
		return
	}
	shared.closeAllOpen()
	dropdownEntry.IsTrayOpen = true
	if selectorEntry, isFound := selectorRegistry.Lookup(layerAlias, dropdownEntry.SelectorAlias); isFound {
		selectorEntry.IsVisible = true
		selectorEntry.ItemHighlighted = dropdownEntry.ItemSelected
		Selector.normalizeIndexes(selectorEntry)
	}
	if scrollBarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, dropdownEntry.ScrollbarAlias); isFound && scrollBarEntry.IsEnabled {
		scrollBarEntry.IsVisible = true
	}
}

/*
closeTrayWithoutCommitting is a method which allows you to close a dropdown's tray without changing the dropdown's
selection, as Esc does. The tray is reverted to the dropdown's selected item, so a later close that commits the tray's
selection, such as closeAllOpen, cannot commit anything picked while the tray was open. The caller must hold
uiStateMutex.

Example:

	Dropdown.closeTrayWithoutCommitting("layer1", dropdownEntry)
*/
func (shared *dropdownType) closeTrayWithoutCommitting(layerAlias string, dropdownEntry *types.DropdownEntryType) {
	if selectorEntry, isFound := selectorRegistry.Lookup(layerAlias, dropdownEntry.SelectorAlias); isFound {
		selectorEntry.ItemSelected = dropdownEntry.ItemSelected
		selectorEntry.ItemHighlighted = dropdownEntry.ItemSelected
		Selector.normalizeIndexes(selectorEntry)
		selectorEntry.IsVisible = false
	}
	if scrollBarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, dropdownEntry.ScrollbarAlias); isFound {
		scrollBarEntry.IsVisible = false
	}
	dropdownEntry.IsTrayOpen = false
}

/*
closeTrayIfDisabled is a method which allows you to close a disabled dropdown's tray without changing its selection, as
Esc does, so that nothing picked in the tray before the dropdown was disabled can be committed. Nothing happens if the
dropdown is enabled, its tray is already closed, or it does not exist. The caller must hold uiStateMutex.

Example:

	Dropdown.closeTrayIfDisabled("layer1", "myDropdown")
*/
func (shared *dropdownType) closeTrayIfDisabled(layerAlias string, dropdownAlias string) {
	dropdownEntry, isFound := dropdownRegistry.Lookup(layerAlias, dropdownAlias)
	if !isFound || dropdownEntry.IsEnabled || !dropdownEntry.IsTrayOpen {
		return
	}
	shared.closeTrayWithoutCommitting(layerAlias, dropdownEntry)
	setEventStateId(constants.EventStateNone)
}

var Dropdown dropdownType
var dropdownRegistry = memory.NewControlMemoryManager[types.DropdownEntryType]()

// dropdownPressOrigin is the dropdown the held mouse press began on, or nil when the press began anywhere else or no
// button is held. It is only read and written while the UI state lock is held, so it needs no synchronization of its
// own.
var dropdownPressOrigin *types.DropdownEntryType

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
	defer lockUIState()()
	shared.BaseControlInstanceType.delete()
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
	defer lockUIState()()
	dropdownEntry := dropdownRegistry.Get(shared.layerAlias, shared.controlAlias)
	if len(dropdownEntry.SelectionEntry.SelectionValue) != 0 &&
		dropdownEntry.ItemSelected >= 0 && dropdownEntry.ItemSelected < len(dropdownEntry.SelectionEntry.SelectionValue) {
		return dropdownEntry.SelectionEntry.SelectionValue[dropdownEntry.ItemSelected]
	}
	return ""
}

/*
GetSelected is a method which allows you to retrieve the currently selected item of a dropdown, as its alias and its
zero based index. If nothing is selected, the recorded selection is out of range, or the dropdown does not exist, an
empty string and constants.SELECTED_NONE are returned. Use GetAlias to obtain the alias of the dropdown control itself.

Example:

	alias, index := dropdown.GetSelected()
*/
func (shared *DropdownInstanceType) GetSelected() (string, int) {
	defer lockUIState()()
	if dropdownEntry, isFound := dropdownRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		itemSelected := dropdownEntry.ItemSelected
		if itemSelected >= 0 && itemSelected < Dropdown.getItemCount(dropdownEntry) {
			return dropdownEntry.SelectionEntry.SelectionAlias[itemSelected], itemSelected
		}
	}
	return "", constants.SELECTED_NONE
}

/*
GetAllItems is a method which allows you to retrieve every item a dropdown offers, as one slice of aliases and one slice
of display values, in the order the items appear. If the dropdown does not exist, two empty slices are returned. In
addition, the following should be noted:

  - The slices returned are copies, so changing them does not change the dropdown. Use SetSelectionEntry to change
    its items.

Example:

	aliases, values := dropdown.GetAllItems()
*/
func (shared *DropdownInstanceType) GetAllItems() ([]string, []string) {
	defer lockUIState()()
	if dropdownEntry, isFound := dropdownRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		selectionEntryCopy := Selector.getSelectionEntryCopy(dropdownEntry.SelectionEntry)
		return selectionEntryCopy.SelectionAlias, selectionEntryCopy.SelectionValue
	}
	return []string{}, []string{}
}

/*
IsOpen is a method which allows you to check whether a dropdown's tray is open, whether the user opened it or it was
opened with Open. If the dropdown does not exist, false is returned.

Example:

	isOpen := dropdown.IsOpen()
*/
func (shared *DropdownInstanceType) IsOpen() bool {
	defer lockUIState()()
	if dropdownEntry, isFound := dropdownRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return dropdownEntry.IsTrayOpen
	}
	return false
}

/*
Open is a method which allows you to open a dropdown's tray, as clicking it or pressing Enter on it does. Any other
open tray is closed first, the tray's highlight is placed on the selected item, and its scroll bar is shown if the
items overflow it. Nothing happens if the tray is already open, the dropdown is disabled, or the dropdown does not
exist. In addition, the following should be noted:

  - Focus is not moved. The arrow keys and Enter only reach the tray while the dropdown has focus, so give it focus
    with SetFocus as well if the user should be able to pick from the tray with the keyboard.

Example:

	dropdown.Open()
*/
func (shared *DropdownInstanceType) Open() {
	defer lockUIState()()
	if dropdownEntry, isFound := dropdownRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound && !dropdownEntry.IsTrayOpen {
		Dropdown.openTray(shared.layerAlias, dropdownEntry)
	}
}

/*
Close is a method which allows you to close a dropdown's tray without changing its selection, as Esc does, so any
item highlighted in the tray is discarded. Nothing happens if the tray is already closed or the dropdown does not
exist.

Example:

	dropdown.Close()
*/
func (shared *DropdownInstanceType) Close() {
	defer lockUIState()()
	if dropdownEntry, isFound := dropdownRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound && dropdownEntry.IsTrayOpen {
		Dropdown.closeTrayWithoutCommitting(shared.layerAlias, dropdownEntry)
		setEventStateId(constants.EventStateNone)
	}
}

/*
GetSelectedItemIndex is a method which retrieves the index of the currently selected item in the dropdown.

Example:
    index := dropdown.GetSelectedItemIndex()
*/
func (shared *DropdownInstanceType) GetSelectedItemIndex() int {
	defer lockUIState()()
	dropdownEntry := dropdownRegistry.Get(shared.layerAlias, shared.controlAlias)
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
	defer lockUIState()()
	dropdownEntry := dropdownRegistry.Get(shared.layerAlias, shared.controlAlias)
	if itemIndex < constants.SELECTED_NONE || itemIndex >= Dropdown.getItemCount(dropdownEntry) {
		return
	}
	dropdownEntry.ItemSelected = itemIndex
	if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, dropdownEntry.SelectorAlias); isFound {
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
tray is open. The tray is widened by one column when its scroll bar is not needed, and narrowed again when it is, so
the tray always lines up with the dropdown. This is safe to call whether the tray is open or closed, and with an
empty list. In addition, the following should be noted:

  - The dropdown keeps its own copy of the item slices, so changing the slices of the entry passed in afterwards
    does not change the dropdown.

Example:

	dropdown.SetSelectionEntry(newSelection)
*/
func (shared *DropdownInstanceType) SetSelectionEntry(selectionEntry types.SelectionEntryType) {
	defer lockUIState()()
	dropdownEntry := dropdownRegistry.Get(shared.layerAlias, shared.controlAlias)
	dropdownEntry.SelectionEntry = Selector.getSelectionEntryCopy(selectionEntry)
	dropdownEntry.ItemSelected = constants.SELECTED_NONE

	selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, dropdownEntry.SelectorAlias)
	if !isFound {
		return
	}
	selectorEntry.SelectionEntry = Selector.getSelectionEntryCopy(selectionEntry)
	selectorEntry.ItemSelected = constants.SELECTED_NONE
	selectorEntry.ItemHighlighted = constants.NullItemSelection
	selectorEntry.ViewportPosition = 0
	Dropdown.updateLayout(shared.layerAlias, dropdownEntry)
}

// dropdownBoxExtraWidth is the number of columns a dropdown's box is wider than its item width: one column of padding
// on each side of the selected item, where the tray's border sits when it opens, plus one for the arrow.
const dropdownBoxExtraWidth = 3

/*
getTrayGeometry is a method which allows you to obtain where a dropdown's tray belongs, as the column and row of the
tray's first item and the tray's item width. The tray opens one row below the dropdown, with its border lined up with
the dropdown's left edge, so its items start one column in. While the items overflow the tray, its scroll bar sits
in the column under the dropdown's arrow and the items are as wide as the dropdown's item width. When every item
fits, there is no scroll bar, so the items are one column wider to fill that column instead.

Example:

	trayXLocation, trayYLocation, trayItemWidth := Dropdown.getTrayGeometry(dropdownEntry)
*/
func (shared *dropdownType) getTrayGeometry(dropdownEntry *types.DropdownEntryType) (int, int, int) {
	trayItemWidth := dropdownEntry.ItemWidth
	if shared.getItemCount(dropdownEntry) <= dropdownEntry.Height {
		trayItemWidth++
	}
	return dropdownEntry.XLocation + 1, dropdownEntry.YLocation + 1, trayItemWidth
}

/*
updateLayout is a method which allows you to bring everything derived from a dropdown's position, size and items back
in line with them, and is called whenever any of those change after the dropdown is added. The tray is moved, given
the dropdown's number of rows, and widened or narrowed to fit its scroll bar, and the tray's own scroll bar, tooltip
and viewport follow it. The dropdown's tooltip is moved and resized to match. In addition, the following should be
noted:

  - This is safe whether the tray is open or closed. An open tray is redrawn in its new place on the next screen
    update, and a closed tray keeps its scroll bar hidden.

  - The caller must hold uiStateMutex.

Example:

	Dropdown.updateLayout("layer1", dropdownEntry)
*/
func (shared *dropdownType) updateLayout(layerAlias string, dropdownEntry *types.DropdownEntryType) {
	if selectorEntry, isFound := selectorRegistry.Lookup(layerAlias, dropdownEntry.SelectorAlias); isFound {
		selectorEntry.XLocation, selectorEntry.YLocation, selectorEntry.ItemWidth = shared.getTrayGeometry(dropdownEntry)
		selectorEntry.Height = dropdownEntry.Height
		Selector.updateLayout(layerAlias, selectorEntry)
	}
	if tooltipEntry, isFound := tooltipRegistry.Lookup(layerAlias, dropdownEntry.TooltipAlias); isFound {
		tooltipEntry.HotspotXLocation, tooltipEntry.HotspotYLocation, tooltipEntry.HotspotWidth = dropdownEntry.XLocation, dropdownEntry.YLocation, dropdownEntry.ItemWidth
		tooltipEntry.TooltipXLocation, tooltipEntry.TooltipYLocation, tooltipEntry.TooltipWidth = dropdownEntry.XLocation, dropdownEntry.YLocation+2, dropdownEntry.ItemWidth
	}
}

/*
setSize is a method which allows you to resize a dropdown. The width is that of the whole box shown while the dropdown
is closed, including its arrow, and the tray opens at the same width. The height is the number of item rows the tray
shows when it opens, since the closed box is always one row tall. The tray, its scroll bar and both tooltips are then
updated to match. A width too narrow to leave at least one column for the selected item, a height below one, or a
dropdown that does not exist leaves everything unchanged. The caller must hold uiStateMutex.

Example:

	Dropdown.setSize("layer1", "myDropdown", 20, 6)
*/
func (shared *dropdownType) setSize(layerAlias string, dropdownAlias string, width int, height int) {
	dropdownEntry, isFound := dropdownRegistry.Lookup(layerAlias, dropdownAlias)
	if !isFound || width <= dropdownBoxExtraWidth || height < 1 {
		return
	}
	dropdownEntry.ItemWidth = width - dropdownBoxExtraWidth
	dropdownEntry.Width = width
	dropdownEntry.Height = height
	shared.updateLayout(layerAlias, dropdownEntry)
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
	defer lockUIState()()
	return shared.add(layerAlias, dropdownAlias, styleEntry, selectionEntry, xLocation, yLocation, selectorHeight, itemWidth, defaultItemSelected)
}

/*
add is a method which allows you to do exactly what Add does, for a caller that already holds uiStateMutex, such as
another entry point of the package that has taken the lock itself.

Example:

	Dropdown.add(layerAlias, dropdownAlias, styleEntry, selectionEntry, xLocation, yLocation, selectorHeight, itemWidth, defaultItemSelected)
*/
func (shared *dropdownType) add(layerAlias string, dropdownAlias string, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, defaultItemSelected int) DropdownInstanceType {
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
	dropdownEntry := &newDropdownEntry
	dropdownEntry.Width = itemWidth + dropdownBoxExtraWidth
	dropdownEntry.Height = selectorHeight
	dropdownEntry.SelectorAlias = stringformat.GetLastSortedUUID()

	trayXLocation, trayYLocation, trayItemWidth := shared.getTrayGeometry(dropdownEntry)
	Selector.add(layerAlias, dropdownEntry.SelectorAlias, styleEntry, selectionEntry, trayXLocation, trayYLocation, selectorHeight, trayItemWidth, 1, 0, 0, false, true)
	// Create associated tooltip (always created but disabled by default)
	tooltipInstance := Tooltip.add(layerAlias, dropdownEntry.TooltipAlias, "", styleEntry,
		dropdownEntry.XLocation, dropdownEntry.YLocation,
		dropdownEntry.ItemWidth, 1,
		dropdownEntry.XLocation, dropdownEntry.YLocation+2,
		dropdownEntry.ItemWidth, 3,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.setEnabled(false)
	tooltipInstance.setParentControlAlias(dropdownAlias)

	// The dropdown entry is complete before it is added to dropdownRegistry, so nothing can find it half set up.
	selectorEntry := selectorRegistry.Get(layerAlias, dropdownEntry.SelectorAlias)
	selectorEntry.IsVisible = false
	// The tray starts out agreeing with the dropdown, so closing it without a pick cannot change the selection.
	selectorEntry.ItemSelected = dropdownEntry.ItemSelected
	dropdownEntry.ScrollbarAlias = selectorEntry.ScrollbarAlias
	scrollBarEntry := scrollbarRegistry.Get(layerAlias, dropdownEntry.ScrollbarAlias)
	scrollBarEntry.IsVisible = false
	if shared.getItemCount(dropdownEntry) <= selectorHeight {
		scrollBarEntry.IsEnabled = false
	}
	dropdownRegistry.Add(layerAlias, dropdownAlias, dropdownEntry)
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
	defer lockUIState()()
	dropdownRegistry.Remove(layerAlias, dropdownAlias)
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
	defer lockUIState()()
	removeAllControlsAndClearCells(dropdownRegistry, layerAlias, constants.CellTypeDropdown, func(entry *types.DropdownEntryType) string { return entry.Alias })
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
	for _, currentDropdownEntry := range dropdownRegistry.GetAllEntries(layerAlias) {
		dropdownAlias := currentDropdownEntry.Alias
		shared.draw(&layerEntry, dropdownAlias, isFocusIndicatorShown(layerAlias, dropdownAlias, constants.CellTypeDropdown))
	}
}

/*
draw is a method which allows you to draw a single dropdown on a given text layer, showing its selected item formatted
to the dropdown's width and alignment, in the style's colours, followed by a down arrow drawn in inverted colours. In
addition, the following should be noted:

  - A disabled dropdown still shows its selected item, but its text and arrow are drawn in the disabled colour
    given by getDisabledForegroundColor, and never with the focused colours.

  - When isFocusShown is true, the dropdown is drawn with the style's focused colours, resolved by getFocusedColors
    so that focus stays visible even when the style does not set them, and the arrow inverts those instead. The
    caller passes isFocusIndicatorShown for the dropdown, which as described for setFocusIndicatorVisible is only
    true while the user's last input came from the keyboard.

  - The caller must hold uiStateMutex.

Example:

	Dropdown.draw(layer, "myDropdown", false)
*/
func (shared *dropdownType) draw(layerEntry *types.LayerEntryType, dropdownAlias string, isFocusShown bool) {
	layerAlias := layerEntry.LayerAlias
	dropdownEntry, isFound := dropdownRegistry.Lookup(layerAlias, dropdownAlias)
	if !isFound {
		return
	}
	localStyleEntry := types.NewTuiStyleEntry(&dropdownEntry.StyleEntry)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = localStyleEntry.Dropdown.ForegroundColor
	attributeEntry.BackgroundColor = localStyleEntry.Dropdown.BackgroundColor
	if !dropdownEntry.IsEnabled {
		attributeEntry.ForegroundColor = shared.getDisabledForegroundColor(localStyleEntry)
	} else if isFocusShown {
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
getDisabledForegroundColor is a method which allows you to obtain the foreground colour a disabled dropdown is drawn
with, so that it reads as muted against the dropdown's normal colours. The style's Button.LabelDisabledColor, the
colour disabled buttons draw their label in, is used when it is set and differs from both the dropdown's foreground and
background colours. Otherwise the first of grey and silver that differs from both is used, and grey when neither does,
so the result is only ever invisible against the dropdown's background when the style leaves no other choice.

Example:

	foregroundColor := Dropdown.getDisabledForegroundColor(styleEntry)
*/
func (shared *dropdownType) getDisabledForegroundColor(styleEntry types.TuiStyleEntryType) constants.ColorType {
	foregroundColor := styleEntry.Dropdown.ForegroundColor
	backgroundColor := styleEntry.Dropdown.BackgroundColor
	candidateColors := []constants.ColorType{
		styleEntry.Button.LabelDisabledColor,
		constants.AnsiColorByIndex[constants.ColorBrightBlack],
		constants.AnsiColorByIndex[constants.ColorWhite],
	}
	for _, candidateColor := range candidateColors {
		if candidateColor != 0 && candidateColor != foregroundColor && candidateColor != backgroundColor {
			return candidateColor
		}
	}
	return constants.AnsiColorByIndex[constants.ColorBrightBlack]
}

/*
updatePressOrigin is a method which allows you to record which dropdown the current mouse press began on. A fresh press
over a dropdown records that dropdown, a fresh press anywhere else records that no dropdown owns the press, and a
release clears the record. Movement with the button still held leaves the record unchanged. In addition, the following
should be noted:

  - It must run for every mouse event, before anything that can stop the event from reaching the controls, such as a
    modal layer rejecting it, so that a press which was never delivered to a dropdown cannot be mistaken for one that
    was.

  - The dropdown is recorded by its entry rather than its alias, so a dropdown deleted during the press, and one added
    later under the same alias, never inherit the press.

  - A disabled dropdown is recorded like any other, as a disabled button is, since the press still belongs to it.
    It is updateStateMouse that refuses to open it.

Example:

	Dropdown.updatePressOrigin(characterEntry, buttonPressed, previousButtonPressed)
*/
func (shared *dropdownType) updatePressOrigin(characterEntry types.CharacterEntryType, buttonPressed uint, previousButtonPressed uint) {
	if buttonPressed == 0 {
		dropdownPressOrigin = nil
		return
	}
	if previousButtonPressed != 0 {
		return
	}
	dropdownPressOrigin = nil
	if characterEntry.AttributeEntry.CellType != constants.CellTypeDropdown {
		return
	}
	if dropdownEntry, isFound := dropdownRegistry.Lookup(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias); isFound {
		dropdownPressOrigin = dropdownEntry
	}
}

/*
updateStateMouse is a method which allows you to update the state of all dropdowns according to the current mouse
event state. Clicking a dropdown opens its tray, clicking elsewhere closes every open tray, and dragging a tray's
scroll bar moves the tray's viewport. It returns true if the screen needs to be updated. In addition, the following
should be noted:

  - Only a fresh press that began on a dropdown, as recorded by updatePressOrigin, opens its tray. A held button
    dragged onto a dropdown, whether the press began elsewhere or was already held when the dropdown was added, never
    opens it.

  - A press on a disabled dropdown never opens it, and is treated as a press outside every dropdown, so it closes
    any tray that is open, as a press on a disabled button does.

  - A viewport position copied from a scroll bar is clamped into the range the tray can scroll to.

Example:

	isUpdateRequired := Dropdown.updateStateMouse()
*/
func (shared *dropdownType) updateStateMouse() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	_, _, previousButtonPress, _ := GetPreviousMouseStatus()
	isFreshPress := buttonPressed != 0 && previousButtonPress == 0
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	cellControlAlias := characterEntry.AttributeEntry.CellControlAlias

	// If a buttonType is pressed AND (you are in a drag and drop event OR the cell type is scroll bar), then
	// sync all Dropdown selectors with their appropriate scroll bars. If the control under focus
	// matches a control that belongs to a Dropdown list, then stop processing (Do not attempt to close Dropdown).
	if buttonPressed != 0 && (getEventStateId() == constants.EventStateDragAndDropScrollbar ||
		characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar) {
		isMatchFound := false
		for _, currentDropdownEntry := range dropdownRegistry.GetAllEntries(layerAlias) {
			dropdownEntry := currentDropdownEntry
			selectorEntry := selectorRegistry.Get(layerAlias, dropdownEntry.SelectorAlias)
			scrollBarEntry := scrollbarRegistry.Get(layerAlias, dropdownEntry.ScrollbarAlias)
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

	// A fresh press that began on this dropdown opens its tray, unless the dropdown is disabled.
	dropdownEntry, isFound := dropdownRegistry.Lookup(layerAlias, cellControlAlias)
	isOverDropdown := cellControlAlias != "" && characterEntry.AttributeEntry.CellType == constants.CellTypeDropdown && isFound
	isOverDisabledDropdown := isOverDropdown && !dropdownEntry.IsEnabled
	if isFreshPress && isOverDropdown && !isOverDisabledDropdown && dropdownEntry == dropdownPressOrigin {
		shared.closeAllOpen()
		dropdownEntry.IsTrayOpen = true
		selectorEntry := selectorRegistry.Get(layerAlias, dropdownEntry.SelectorAlias)
		selectorEntry.IsVisible = true
		scrollBarEntry := scrollbarRegistry.Get(layerAlias, dropdownEntry.ScrollbarAlias)
		if scrollBarEntry.IsEnabled {
			scrollBarEntry.IsVisible = true
		}
		isUpdateRequired = true
		return isUpdateRequired
	}

	// Only close dropdowns if clicking outside both the dropdown and its scrollbar
	if isFreshPress {
		// Check if we're clicking on a scrollbar that belongs to an open dropdown
		isScrollbarOfOpenDropdown := false
		if characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar {
			for _, currentDropdownEntry := range dropdownRegistry.GetAllEntries(layerAlias) {
				dropdownEntry := currentDropdownEntry
				if dropdownEntry.IsTrayOpen && dropdownEntry.ScrollbarAlias == cellControlAlias {
					isScrollbarOfOpenDropdown = true
					break
				}
			}
		}

		// Only close if not clicking on an enabled dropdown or the scrollbar of an open one
		if (characterEntry.AttributeEntry.CellType != constants.CellTypeDropdown || isOverDisabledDropdown) && !isScrollbarOfOpenDropdown {
			isUpdateRequired = shared.closeAllOpen()
		}
	}
	return isUpdateRequired
}

/*
closeAllOpenOnLayer is a method which closes every open dropdown on a given text layer. Each open tray and its scroll
bar are hidden, and the tray's selected item is committed to the dropdown if it changed while the tray was open. It
returns true if any dropdown was closed. In addition, the following should be noted:

  - The caller must hold uiStateMutex.

Example:

	isClosed := Dropdown.closeAllOpenOnLayer("layer1")
*/
func (shared *dropdownType) closeAllOpenOnLayer(layerAlias string) bool {
	isAnyDropdownClosed := false
	for _, currentDropdownEntry := range dropdownRegistry.GetAllEntries(layerAlias) {
		dropdownEntry := currentDropdownEntry
		if dropdownEntry.IsTrayOpen == true {
			selectorEntry := selectorRegistry.Get(layerAlias, dropdownEntry.SelectorAlias)
			selectorEntry.IsVisible = false
			scrollBarEntry := scrollbarRegistry.Get(layerAlias, dropdownEntry.ScrollbarAlias)
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
closeAllOpen is a method which allows you to close every open dropdown on every layer, as described for
closeAllOpenOnLayer, such as before focus moves to another control. It returns true if any dropdown was closed.

Example:

	isClosed := Dropdown.closeAllOpen()
*/
func (shared *dropdownType) closeAllOpen() bool {
	var wasAnyDropdownClosed bool
	dropdownRegistry.MemoryManager.Range(func(key, value interface{}) bool {
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
get is a method which allows you to retrieve the entry holding a dropdown's state from the control memory manager. If
the dropdown does not exist, a panic is generated to fail as fast as possible. The caller must hold uiStateMutex.

Example:

	entry := Dropdown.get("layer1", "myDropdown")
*/
func (shared *dropdownType) get(layerAlias string, dropdownAlias string) *types.DropdownEntryType {
	return dropdownRegistry.Get(layerAlias, dropdownAlias)
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
	defer lockUIState()()
	return dropdownRegistry.IsExists(layerAlias, dropdownAlias)
}
