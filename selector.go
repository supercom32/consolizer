package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

type SelectorInstanceType struct {
	BaseControlInstanceType
}

type selectorType struct{}

var Selector selectorType
var Selectors = memory.NewControlMemoryManager[types.SelectorEntryType]()

/*
IsSelectorExists is a method which checks if a selector with the specified alias exists on a given text layer.

In addition, the following should be noted:

- Returns true if the selector exists, false otherwise.

- This method is useful for validating selector existence before performing operations on it.

Example:
    exists := IsSelectorExists("Layer1", "Selector1")
*/
func IsSelectorExists(layerAlias string, selectorAlias string) bool {
	// Use the generic memory manager to check existence
	return Selectors.IsExists(layerAlias, selectorAlias)
}

/*
GetSelector is a method which retrieves a selector entry from a given text layer.

In addition, the following should be noted:

- If the selector does not exist, a panic will be generated to fail as fast as possible.

- The returned selector entry can be used to directly modify the selector's properties.

- Changes made to the returned entry will be reflected when the selector is next drawn.

Example:
    selectorEntry := GetSelector("Layer1", "Selector1")
*/
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
AddToTabIndex is a method which adds a selector to the tab index. This enables keyboard navigation between controls using the tab key. In addition, the following should be noted:

- The selector will be added to the tab order based on the order in which it was created.

- The tab index is used to determine which control receives focus when the tab key is pressed.

Example:
    selector.AddToTabIndex()
*/
func (shared *SelectorInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeSelectorItem)
}

/*
Delete is a method which removes a selector from a text layer.

In addition, the following should be noted:

- If you attempt to delete a selector which does not exist, then the request will simply be ignored.

- All memory associated with the selector will be freed.

Example:
    selector = selector.Delete()
*/
func (shared *SelectorInstanceType) Delete() *SelectorInstanceType {
	shared.BaseControlInstanceType.Delete()
	return nil
}

/*
IsNewItemSelected is a method which checks if a new item has been selected in the selector.

Example:
    if selector.IsNewItemSelected() { ... }
*/
func (shared *SelectorInstanceType) IsNewItemSelected() bool {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		selectorEntry := Selectors.Get(shared.layerAlias, shared.controlAlias)
		return selectorEntry.IsNewItemSelected
	}
	return false
}

/*
GetSelected is a method which retrieves the currently selected item from a selector.

In addition, the following should be noted:

- Returns both the alias and index of the selected item.

- If the selector does not exist, returns an empty string and -1.

- The alias is typically used for display purposes, while the index is used for programmatic access to the selection.

Example:
    alias, index := selector.GetSelected()
*/
func (shared *SelectorInstanceType) GetSelected() (string, int) {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorMenu(shared.layerAlias, shared.controlAlias)
		menuEntry := Selectors.Get(shared.layerAlias, shared.controlAlias)
		menuEntry.IsNewItemSelected = false
		value := menuEntry.ItemSelected
		if value == constants.SELECTED_NONE {
			return "", constants.SELECTED_NONE
		}
		if len(menuEntry.SelectionEntry.SelectionAlias) > value {
			return menuEntry.SelectionEntry.SelectionAlias[value], value
		}
	}
	return "", constants.SELECTED_NONE
}

/*
GetAllItems is a method which retrieves all items from a selector.

In addition, the following should be noted:

- Returns two arrays: one containing all aliases and one containing all values.

- If the selector does not exist, returns empty arrays.

- The arrays are returned in the order they were added to the selector.

Example:
    aliases, values := selector.GetAllItems()
*/
func (shared *SelectorInstanceType) GetAllItems() ([]string, []string) {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorMenu(shared.layerAlias, shared.controlAlias)
		menuEntry := Selectors.Get(shared.layerAlias, shared.controlAlias)
		return menuEntry.SelectionEntry.SelectionAlias, menuEntry.SelectionEntry.SelectionValue
	}
	return []string{}, []string{}
}

/*
Unselect is a method which clears the current selection for a selector.

In addition, the following should be noted:

- The selector's selected item will be set to SELECTED_NONE.

- If the selector does not exist, no operation occurs.

Example:
    selector.Unselect()
*/
func (shared *SelectorInstanceType) Unselect() {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorMenu(shared.layerAlias, shared.controlAlias)
		selectorEntry := Selectors.Get(shared.layerAlias, shared.controlAlias)
		selectorEntry.ItemSelected = constants.SELECTED_NONE
	}
}

/*
Select is a method which allows you to select an item by its alias.

In addition, the following should be noted:

- If an item with the matching alias is found, it will be set as the selected and highlighted item.

- If the selector or the item does not exist, no operation occurs.

Example:
    selector.Select("Option1")
*/
func (shared *SelectorInstanceType) Select(selectionAlias string) {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorMenu(shared.layerAlias, shared.controlAlias)
		selectorEntry := Selectors.Get(shared.layerAlias, shared.controlAlias)

		// Find the index of the item with the matching alias
		itemIndex := -1
		for i, alias := range selectorEntry.SelectionEntry.SelectionAlias {
			if alias == selectionAlias {
				itemIndex = i
				break
			}
		}

		// If the item exists, set it as selected.
		if itemIndex != -1 {
			selectorEntry.ItemSelected = itemIndex
			selectorEntry.ItemHighlighted = itemIndex
			selectorEntry.IsNewItemSelected = true
		}
	}
}

/*
FocusSelection is a method which allows you to focus on a specific item in the selector by its alias.

In addition, the following should be noted:

- The item with the matching alias will be scrolled into view.

- If possible, the item will be centered in the visible area.

- If the item is near the top or bottom, the viewport will be adjusted accordingly.

- The scroll bar position will be updated to reflect the new viewport position.

- If the selector or item does not exist, no operation occurs.

Example:
    selector.FocusSelection("Option5")
*/
func (shared *SelectorInstanceType) FocusSelection(selectionAlias string) {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorMenu(shared.layerAlias, shared.controlAlias)
		selectorEntry := Selectors.Get(shared.layerAlias, shared.controlAlias)

		// Find the index of the item with the matching alias
		itemIndex := -1
		for i, alias := range selectorEntry.SelectionEntry.SelectionAlias {
			if alias == selectionAlias {
				itemIndex = i
				break
			}
		}

		// If the item doesn't exist, do nothing
		if itemIndex == -1 {
			return
		}

		// Calculate the number of items visible in the viewport
		visibleItems := selectorEntry.Height * selectorEntry.NumberOfColumns

		// Calculate the ideal viewport position to center the item
		idealPosition := itemIndex - (visibleItems / 2)

		// Adjust the viewport position to ensure it's within valid bounds
		maxPosition := len(selectorEntry.SelectionEntry.SelectionValue) - visibleItems
		if maxPosition < 0 {
			maxPosition = 0
		}

		newPosition := idealPosition
		if newPosition < 0 {
			newPosition = 0
		} else if newPosition > maxPosition {
			newPosition = maxPosition
		}

		// Update the viewport position
		selectorEntry.ViewportPosition = newPosition

		// Update scrollbar if it exists
		if selectorEntry.ScrollbarAlias != "" && ScrollBars.IsExists(shared.layerAlias, selectorEntry.ScrollbarAlias) {
			scrollBarEntry := ScrollBars.Get(shared.layerAlias, selectorEntry.ScrollbarAlias)

			// Update the scroll value based on the new viewport position
			scrollBarEntry.ScrollValue = newPosition

			// Compute and update the handle position
			scrollbar.computeHandlePositionByScrollValue(shared.layerAlias, selectorEntry.ScrollbarAlias)
		}
	}
}

/*
setViewport is a method which allows you to specify the current viewport index for a given selector.

In addition, the following should be noted:

- The viewport determines which items are currently visible in the selector.

- If the selector does not exist, no operation occurs.

- The viewport position is automatically adjusted when navigating through items.

Example:
    selector.setViewport(10)
*/
func (shared *SelectorInstanceType) setViewport(viewportPosition int) {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorMenu(shared.layerAlias, shared.controlAlias)
		menuEntry := Selectors.Get(shared.layerAlias, shared.controlAlias)
		menuEntry.ViewportPosition = viewportPosition
	}
}

/*
SetSelectionEntry is a method which allows you to overwrite the current selection entry with a new one.

In addition, the following should be noted:

- The selector's selected and highlighted item will be reset.

- The viewport will be reset to the beginning.

- The associated scroll bar will be updated to reflect the new item list.

- If the selector does not exist, no operation occurs.

Example:
    selector.SetSelectionEntry(newSelection)
*/
func (shared *SelectorInstanceType) SetSelectionEntry(selectionEntry types.SelectionEntryType) {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		selectorEntry := Selectors.Get(shared.layerAlias, shared.controlAlias)
		selectorEntry.SelectionEntry = selectionEntry
		selectorEntry.ItemSelected = constants.SELECTED_NONE
		selectorEntry.ItemHighlighted = constants.NullItemSelection
		selectorEntry.ViewportPosition = 0

		// Update scrollbar if it exists
		if selectorEntry.ScrollbarAlias != "" && ScrollBars.IsExists(shared.layerAlias, selectorEntry.ScrollbarAlias) {
			scrollBarEntry := ScrollBars.Get(shared.layerAlias, selectorEntry.ScrollbarAlias)

			// Calculate max scroll value
			scrollBarMaxValue := len(selectorEntry.SelectionEntry.SelectionValue) - (selectorEntry.Height * selectorEntry.NumberOfColumns) + 1

			// Enable or disable scrollbar based on whether items overflow
			if len(selectorEntry.SelectionEntry.SelectionValue) > selectorEntry.Height*selectorEntry.NumberOfColumns &&
				selectorEntry.StyleEntry.Selector.TextAlignment != constants.AlignmentNoPadding {
				scrollBarEntry.IsEnabled = true
				scrollBarEntry.IsVisible = true
			} else {
				scrollBarEntry.IsEnabled = false
				scrollBarEntry.IsVisible = false
			}

			if scrollBarMaxValue < 0 {
				scrollBarMaxValue = 0
			}
			scrollBarEntry.MaxScrollValue = scrollBarMaxValue
		}
	}
}

/*
AddItem is a method which allows you to add a new selector item to the already loaded list of selector items.

In addition, the following should be noted:

- The new item is added to the end of the list.

- Both the alias and value for the item must be provided.

- If the selector does not exist, no operation occurs.

- Scroll bars are automatically enabled if items overflow the visible area.

Example:
    selector.AddItem("NewOpt", "New Option")
*/
func (shared *SelectorInstanceType) AddItem(selectionAlias string, selectionValue string) {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorMenu(shared.layerAlias, shared.controlAlias)
		menuEntry := Selectors.Get(shared.layerAlias, shared.controlAlias)
		menuEntry.SelectionEntry.Add(selectionAlias, selectionValue)

		// Update scrollbar if it exists
		if menuEntry.ScrollbarAlias != "" && ScrollBars.IsExists(shared.layerAlias, menuEntry.ScrollbarAlias) {
			scrollBarEntry := ScrollBars.Get(shared.layerAlias, menuEntry.ScrollbarAlias)
			scrollBarMaxValue := len(menuEntry.SelectionEntry.SelectionValue) - (menuEntry.Height * menuEntry.NumberOfColumns) + 1

			// Enable scrollbar if items overflow
			if len(menuEntry.SelectionEntry.SelectionValue) > menuEntry.Height*menuEntry.NumberOfColumns &&
				menuEntry.StyleEntry.Selector.TextAlignment != constants.AlignmentNoPadding {
				scrollBarEntry.IsEnabled = true
				scrollBarEntry.IsVisible = true
			}

			// Update max value
			if scrollBarMaxValue < 0 {
				scrollBarMaxValue = 0
			}
			scrollBarEntry.MaxScrollValue = scrollBarMaxValue
		}
	}
}

/*
DeleteItem is a method which allows you to delete a selector item at a specified index from the list of selector items.

In addition, the following should be noted:

- The index is zero-based.

- If the index is out of range, no operation occurs.

- If the selector does not exist, no operation occurs.

- If the currently highlighted or selected item is deleted, the highlight or selection is adjusted.

- Scroll bars are automatically disabled if items no longer overflow the visible area.

Example:
    selector.DeleteItem(2)
*/
func (shared *SelectorInstanceType) DeleteItem(index int) {
	if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorMenu(shared.layerAlias, shared.controlAlias)
		menuEntry := Selectors.Get(shared.layerAlias, shared.controlAlias)

		// Check if index is valid
		if index < 0 || index >= len(menuEntry.SelectionEntry.SelectionAlias) {
			return
		}

		// Create new slices without the item at the specified index
		newAliases := append(menuEntry.SelectionEntry.SelectionAlias[:index], menuEntry.SelectionEntry.SelectionAlias[index+1:]...)
		newValues := append(menuEntry.SelectionEntry.SelectionValue[:index], menuEntry.SelectionEntry.SelectionValue[index+1:]...)

		// Update the selection entry
		menuEntry.SelectionEntry.SelectionAlias = newAliases
		menuEntry.SelectionEntry.SelectionValue = newValues

		// Adjust highlighted and selected items if necessary
		if menuEntry.ItemHighlighted >= index {
			if menuEntry.ItemHighlighted > 0 {
				menuEntry.ItemHighlighted--
			}
		}
		if menuEntry.ItemSelected >= index {
			if menuEntry.ItemSelected > 0 {
				menuEntry.ItemSelected--
			}
		}

		// Update scrollbar if it exists
		if menuEntry.ScrollbarAlias != "" && ScrollBars.IsExists(shared.layerAlias, menuEntry.ScrollbarAlias) {
			scrollBarEntry := ScrollBars.Get(shared.layerAlias, menuEntry.ScrollbarAlias)

			// Disable scrollbar if items no longer overflow
			if len(menuEntry.SelectionEntry.SelectionValue) <= menuEntry.Height*menuEntry.NumberOfColumns ||
				menuEntry.StyleEntry.Selector.TextAlignment == constants.AlignmentNoPadding {
				scrollBarEntry.IsEnabled = false
				scrollBarEntry.IsVisible = false
			} else {
				scrollBarEntry.IsEnabled = true
				scrollBarEntry.IsVisible = true
			}

			// Calculate max scroll value
			scrollBarMaxValue := len(menuEntry.SelectionEntry.SelectionValue) - (menuEntry.Height * menuEntry.NumberOfColumns) + 1

			// Update max value
			if scrollBarMaxValue < 0 {
				scrollBarMaxValue = 0
			}
			scrollBarEntry.MaxScrollValue = scrollBarMaxValue
		}
	}
}

/*
Add is a method which allows you to add a selector to a given text layer. Once called, an instance of your control is
returned which will allow you to read or manipulate the properties for it. The style of the selector will be determined
by the style entry passed in. If you wish to remove a selector from a text layer, simply call DeleteSelector.

In addition, the following should be noted:

- Selectors are not drawn physically to the text layer provided. Instead, they are rendered to the terminal at the same
  time when the text layer is rendered. This allows you to create selectors without actually overwriting the text layer
  data under it.

- If the selector to be drawn falls outside the range of the provided layer, then only the visible portion of the
  selector will be drawn.

- If the selector height is greater than the number of selections available, then no scroll bars are drawn.

Example:
    selectorInstance := Selector.Add("Layer1", "Selector1", style, selection, 0, 0, 10, 20, 1, 0, 0, false, true)
*/
// TODO: Protect against viewport out of range errors.
func (shared *selectorType) Add(layerAlias string, selectorAlias string, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, selectedItem int, highlightOnClickOnly bool, isBorderDrawn bool) SelectorInstanceType {
	newSelectorEntry := types.NewSelectorEntry()
	newSelectorEntry.Alias = selectorAlias
	newSelectorEntry.StyleEntry = styleEntry
	newSelectorEntry.SelectionEntry = selectionEntry
	newSelectorEntry.XLocation = xLocation
	newSelectorEntry.YLocation = yLocation
	newSelectorEntry.Height = selectorHeight
	newSelectorEntry.ItemWidth = itemWidth
	newSelectorEntry.NumberOfColumns = numberOfColumns
	newSelectorEntry.HighlightOnClickOnly = highlightOnClickOnly
	newSelectorEntry.ViewportPosition = viewportPosition
	newSelectorEntry.ItemHighlighted = selectedItem
	//newSelectorEntry.ItemSelected = constants.SELECTED_NONE
	newSelectorEntry.IsBorderDrawn = isBorderDrawn
	newSelectorEntry.IsVisible = true

	// Use the generic memory manager to add the selector entry
	Selectors.Add(layerAlias, selectorAlias, &newSelectorEntry)
	// TODO: AddLayer verification to ensure no item can be 0 length/number.

	tooltipInstance := Tooltip.Add(layerAlias, newSelectorEntry.TooltipAlias, "", styleEntry,
		newSelectorEntry.XLocation, newSelectorEntry.YLocation,
		newSelectorEntry.ItemWidth*newSelectorEntry.NumberOfColumns+2, 1,
		newSelectorEntry.XLocation, newSelectorEntry.YLocation+1,
		newSelectorEntry.ItemWidth*newSelectorEntry.NumberOfColumns+2, 3,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.SetEnabled(false)
	tooltipInstance.setParentControlAlias(selectorAlias)
	selectorEntry := Selectors.Get(layerAlias, selectorAlias)
	selectorEntry.ScrollbarAlias = stringformat.GetLastSortedUUID()

	// Calculate max scroll value
	scrollBarMaxValue := len(selectionEntry.SelectionValue) - (selectorHeight * numberOfColumns) + 1

	// Position scrollbar at the edge of the selector area
	scrollBarXLocation := xLocation + (itemWidth * numberOfColumns) - 1
	scrollBarYLocation := yLocation
	scrollBarHeight := selectorHeight

	if isBorderDrawn {
		scrollBarXLocation = xLocation + (itemWidth * numberOfColumns) + 1
		scrollBarYLocation = scrollBarYLocation - 1
		scrollBarHeight = selectorHeight + 2
	}

	scrollbar.Add(layerAlias, selectorEntry.ScrollbarAlias, styleEntry, scrollBarXLocation, scrollBarYLocation, scrollBarHeight, scrollBarMaxValue, 0, numberOfColumns, false)
	scrollBarEntry := ScrollBars.Get(layerAlias, selectorEntry.ScrollbarAlias)

	// Set parent control information for scrollbar
	if scrollBarEntry != nil {
		scrollBarEntry.ParentControlAlias = selectorAlias
		scrollBarEntry.ParentControlType = constants.CellTypeSelectorItem
	}

	if len(selectionEntry.SelectionValue) <= selectorHeight*numberOfColumns || styleEntry.Selector.TextAlignment == constants.AlignmentNoPadding {
		scrollBarEntry.IsEnabled = false
		scrollBarEntry.IsVisible = false
	}
	var selectorInstance SelectorInstanceType
	selectorInstance.layerAlias = layerAlias
	selectorInstance.controlAlias = selectorAlias
	selectorInstance.controlType = constants.TYPE_SELECTOR
	setFocusedControl(layerAlias, selectorAlias, constants.CellTypeSelectorItem)
	return selectorInstance
}

/*
Delete is a method which allows you to remove a selector from a text layer.

In addition, the following should be noted:

- If you attempt to delete a selector which does not exist, then the request will simply be ignored.

- All memory associated with the selector will be freed.

Example:
    Selector.Delete("Layer1", "Selector1")
*/
func (shared *selectorType) Delete(layerAlias string, selectorAlias string) {
	Selectors.Remove(layerAlias, selectorAlias)
}

/*
DeleteAll is a method which allows you to remove all selectors from a text layer.

In addition, the following should be noted:

- This operation cannot be undone.

- All memory associated with the selectors will be freed.

Example:
    Selector.DeleteAll("Layer1")
*/
func (shared *selectorType) DeleteAll(layerAlias string) {
	Selectors.RemoveAll(layerAlias)
}

/*
setupSelectorAttributes is a method which allows you to create and configure the standard and highlight attribute
entries for a selector based on the provided style entry.

In addition, the following should be noted:

- Returns two attribute entries: one for normal menu items and one for highlighted items.

- The attributes are configured based on the colors specified in the style entry.

- These attributes control the visual appearance of selector items when drawn.

Example:
    attr, highAttr := Selector.setupSelectorAttributes(style)
*/
func (shared *selectorType) setupSelectorAttributes(styleEntry types.TuiStyleEntryType) (types.AttributeEntryType, types.AttributeEntryType) {
	menuAttributeEntry := types.NewAttributeEntry()
	menuAttributeEntry.ForegroundColor = styleEntry.Selector.ForegroundColor
	menuAttributeEntry.BackgroundColor = styleEntry.Selector.BackgroundColor

	highlightAttributeEntry := types.NewAttributeEntry()
	highlightAttributeEntry.ForegroundColor = styleEntry.Selector.HighlightForegroundColor
	highlightAttributeEntry.BackgroundColor = styleEntry.Selector.HighlightBackgroundColor

	return menuAttributeEntry, highlightAttributeEntry
}

/*
drawSelectorBorder is a method which allows you to draw a border around a selector on the specified text layer.

In addition, the following should be noted:

- The border is drawn using the border characters defined in the style entry.

- The border is drawn one character outside the selector area.

- The background of the border area is filled with spaces using the provided attribute entry.

- If IsShadowDrawn is enabled, a shadow is drawn using drawWindow, otherwise a border is drawn using drawBorder.

Example:
    Selector.drawSelectorBorder(&layerEntry, style, attr, 0, 0, 20, 10)
*/
func (shared *selectorType) drawSelectorBorder(layerEntry *types.LayerEntryType, styleEntry types.TuiStyleEntryType,
	attributeEntry types.AttributeEntryType, xLocation int, yLocation int, itemWidth int, selectorHeight int) {
	fillArea(layerEntry, attributeEntry, " ", xLocation-1, yLocation-1, itemWidth+2, selectorHeight+2, constants.NullCellControlLocation)

	if styleEntry.Selector.IsShadowDrawn {
		drawWindow(layerEntry, styleEntry, attributeEntry, xLocation-1, yLocation-1, itemWidth+2, selectorHeight+2, false)
	} else {
		drawBorder(layerEntry, styleEntry, attributeEntry, xLocation-1, yLocation-1, itemWidth+2, selectorHeight+2, false)
	}
}

/*
formatSelectorItemText is a method which allows you to format the text for a selector item based on the provided style
and whether the item is highlighted.

In addition, the following should be noted:

- Handles text alignment according to the style entry's text alignment setting.

- For centered text, special handling is applied to ensure proper centering with markup.

- For highlighted items, markup tags are stripped to ensure consistent highlighting.

- Returns a string formatted to the specified width with appropriate padding.

Example:
    formattedText := Selector.formatSelectorItemText("Option 1", 20, style, false)
*/
func (shared *selectorType) formatSelectorItemText(menuItemText string, itemWidth int, styleEntry types.TuiStyleEntryType, isHighlighted bool) string {
	var menuItemName string

	if styleEntry.Selector.IsSelectionCentered {
		// If centered, we need to handle markup specially
		if isHighlighted {
			// For highlighted items, strip markup tags to ensure they don't apply
			menuItemText = GetNonMarkupText(menuItemText)
			menuItemName = stringformat.GetFormattedString(menuItemText, itemWidth, constants.AlignmentCenter)
		} else {
			// For non-highlighted items, calculate centering that accounts for markup
			textLength := CalculateStringLengthWithoutMarkup(menuItemText)
			padding := (itemWidth - textLength) / 2
			if padding < 0 {
				padding = 0
			}

			// Create padding
			leftPadding := stringformat.GetFilledString(padding, " ")
			rightPadding := stringformat.GetFilledString(itemWidth-textLength-padding, " ")

			// Combine with original text (preserving markup)
			menuItemName = leftPadding + menuItemText
			if len(menuItemName) < itemWidth {
				menuItemName += rightPadding
			}
		}
	} else {
		// Use standard formatting for non-centered items
		if isHighlighted {
			// For highlighted items, strip markup tags
			menuItemText = GetNonMarkupText(menuItemText)
		}
		menuItemName = stringformat.GetFormattedString(menuItemText, itemWidth, styleEntry.Selector.TextAlignment)
	}

	return menuItemName
}

/*
drawSelectorItem is a method which allows you to draw a single selector item on the specified text layer.

In addition, the following should be noted:

- For non-highlighted items, markup processing is applied to support styled text.

- For highlighted items, standard printing is used with the highlight attributes.

- The cell control ID and alias are set to enable mouse and keyboard interaction.

- Returns the width of the drawn item, which may vary based on the content.

Example:
    width := Selector.drawSelectorItem(&layerEntry, attr, "Option 1", 0, 0, 0, false)
*/
func (shared *selectorType) drawSelectorItem(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType,
	menuItemName string, xLocation int, currentXOffset int, currentYLocation int, isHighlighted bool) int {

	arrayOfRunes := stringformat.GetRunesFromString(menuItemName)

	// Use printMarkup for non-highlighted items, otherwise use printLayer
	if !isHighlighted {
		layer.printMarkup(layerEntry, attributeEntry, xLocation+(currentXOffset), currentYLocation, 0, menuItemName)
	} else {
		layer.printLayer(layerEntry, attributeEntry, xLocation+(currentXOffset), currentYLocation, arrayOfRunes)
	}

	return stringformat.GetWidthOfRunesWhenPrinted(arrayOfRunes)
}

/*
drawSelector is a method which allows you to draw a selector on a given text layer. The style of the selector will be
determined by the style entry passed in.

In addition, the following should be noted:

- Selectors are not drawn physically to the text layer provided. Instead, they are rendered to the terminal at the
  same time when the text layer is rendered.

- If the selector to be drawn falls outside the range of the provided layer, then only the visible portion of the
  selector will be drawn.

Example:
    Selector.drawSelector("Sel1", &layerEntry, style, selection, 0, 0, 10, 20, 1, 0, 0)
*/
func (shared *selectorType) drawSelector(selectorAlias string, layerEntry *types.LayerEntryType, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, itemHighlighted int) {
	selectorEntry := Selectors.Get(layerEntry.LayerAlias, selectorAlias)
	if selectorEntry.IsVisible == false {
		return
	}

	menuAttributeEntry, highlightAttributeEntry := shared.setupSelectorAttributes(styleEntry)

	if selectorEntry.IsBorderDrawn {
		shared.drawSelectorBorder(layerEntry, styleEntry, menuAttributeEntry, xLocation, yLocation, itemWidth, selectorHeight)
	}

	currentYLocation := yLocation
	currentMenuItemIndex := viewportPosition
	currentXOffset := 0
	currentColumn := 0
	currentRow := 0
	for currentMenuItemIndex < len(selectionEntry.SelectionValue) && currentRow < selectorHeight {
		isHighlighted := currentMenuItemIndex == itemHighlighted
		attributeEntry := menuAttributeEntry
		if isHighlighted {
			attributeEntry = highlightAttributeEntry
		}
		menuItemName := shared.formatSelectorItemText(selectionEntry.SelectionValue[currentMenuItemIndex], itemWidth, styleEntry, isHighlighted)
		attributeEntry.CellControlId = currentMenuItemIndex
		attributeEntry.CellControlAlias = selectorAlias
		attributeEntry.CellType = constants.CellTypeSelectorItem
		itemWidth := shared.drawSelectorItem(layerEntry, attributeEntry, menuItemName, xLocation, currentXOffset, currentYLocation, isHighlighted)
		currentMenuItemIndex++
		currentXOffset = currentXOffset + itemWidth
		currentColumn++
		if currentColumn >= numberOfColumns {
			currentXOffset = 0
			currentColumn = 0
			currentYLocation++
			currentRow++
		}
	}
	scrollbar.drawOnLayerByAlias(layerEntry, selectorEntry.ScrollbarAlias)
}

/*
drawSelectorsOnLayer is a method which allows you to draw all selectors on a given text layer.

In addition, the following should be noted:

- Selectors are drawn in alphabetical order by their alias.

- This ensures consistent rendering order across multiple frames.

- Internally generated selectors (like those used by dropdowns) are drawn last.

Example:
    Selector.drawSelectorsOnLayer(layerEntry)
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

/*
updateKeyboardEventForSelector is a method which allows you to process keyboard events for a specific selector.

In addition, the following should be noted:

- Handles navigation keys (up, down, left, right) to move between items.

- Automatically adjusts the viewport position when navigating to items outside the visible area.

- Updates the associated scroll bar position when the viewport changes.

- Returns true if the screen needs to be updated due to state changes.

Example:
    updateRequired, consumed := Selector.updateKeyboardEventForSelector("Layer1", "Sel1", rune("down"))
*/
func (shared *selectorType) updateKeyboardEventForSelector(layerAlias string, selectorAlias string, keystroke []rune) (bool, bool) {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := false
	isKeystrokeConsumed := false
	selectorEntry := Selectors.Get(layerAlias, selectorAlias)

	// Use full number of columns
	effectiveColumns := selectorEntry.NumberOfColumns

	if keystrokeAsString == "down" {
		// remainder := selectorEntry.ItemHighlighted % effectiveColumns
		selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + effectiveColumns
		if selectorEntry.ItemHighlighted >= len(selectorEntry.SelectionEntry.SelectionAlias) {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - effectiveColumns
		}
		// Adjust viewport if highlighted item is outside visible range
		if selectorEntry.ItemHighlighted >= selectorEntry.ViewportPosition+(selectorEntry.Height*effectiveColumns) {
			selectorEntry.ViewportPosition = selectorEntry.ItemHighlighted - (selectorEntry.Height * effectiveColumns) + effectiveColumns
			// Update associated scrollbar
			if scrollBarEntry := ScrollBars.Get(layerAlias, selectorEntry.ScrollbarAlias); scrollBarEntry != nil {
				scrollBarEntry.ScrollValue = selectorEntry.ViewportPosition
				scrollbar.computeHandlePositionByScrollValue(layerAlias, selectorEntry.ScrollbarAlias)
			}
		}
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true
	}
	if keystrokeAsString == "up" {
		selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - effectiveColumns
		if selectorEntry.ItemHighlighted < 0 {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + effectiveColumns
		}
		// Adjust viewport if highlighted item is outside visible range
		if selectorEntry.ItemHighlighted < selectorEntry.ViewportPosition {
			selectorEntry.ViewportPosition = selectorEntry.ItemHighlighted
			// Update associated scrollbar
			if scrollBarEntry := ScrollBars.Get(layerAlias, selectorEntry.ScrollbarAlias); scrollBarEntry != nil {
				scrollBarEntry.ScrollValue = selectorEntry.ViewportPosition
				scrollbar.computeHandlePositionByScrollValue(layerAlias, selectorEntry.ScrollbarAlias)
			}
		}
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true
	}
	if keystrokeAsString == "left" {
		if selectorEntry.ItemHighlighted%effectiveColumns != 0 {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - 1
			if selectorEntry.ItemHighlighted < 0 {
				selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + 1
			}
			isScreenUpdateRequired = true
			isKeystrokeConsumed = true
		}
	}
	if keystrokeAsString == "right" {
		if selectorEntry.ItemHighlighted%effectiveColumns != effectiveColumns-1 {
			selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted + 1
			if selectorEntry.ItemHighlighted >= len(selectorEntry.SelectionEntry.SelectionAlias) {
				selectorEntry.ItemHighlighted = selectorEntry.ItemHighlighted - 1
			}
			isScreenUpdateRequired = true
			isKeystrokeConsumed = true
		}
	}
	if keystrokeAsString == "enter" {
		selectorEntry.ItemSelected = selectorEntry.ItemHighlighted
		selectorEntry.IsNewItemSelected = true
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true
	}
	return isScreenUpdateRequired, isKeystrokeConsumed
}

/*
updateKeyboardEvent is a method which allows you to update the state of all selectors according to the current keystroke
event.

In addition, the following should be noted:

- Handles navigation keys (up, down, left, right) to move between items.

- Enter key selects the currently highlighted item.

- Returns true if the screen needs to be updated due to state changes.

Example:
    updateRequired, consumed := Selector.updateKeyboardEvent(rune("up"))
*/
func (shared *selectorType) updateKeyboardEvent(keystroke []rune) (bool, bool) {
	isScreenUpdateRequired := false
	isKeystrokeConsumed := false
	if eventStateMemory.currentlyFocusedControl.controlType != constants.CellTypeSelectorItem || !Selectors.IsExists(eventStateMemory.currentlyFocusedControl.layerAlias, eventStateMemory.currentlyFocusedControl.controlAlias) {
		return isScreenUpdateRequired, isKeystrokeConsumed
	}
	return shared.updateKeyboardEventForSelector(eventStateMemory.currentlyFocusedControl.layerAlias, eventStateMemory.currentlyFocusedControl.controlAlias, keystroke)
}

/*
updateMouseEvent is a method which allows you to update the state of all selectors according to the current mouse event
state.

In addition, the following should be noted:

- Handles mouse clicks to select items.

- Manages scroll bar synchronization for selectors with many items.

- Returns true if the screen needs to be updated due to state changes.

Example:
    updateRequired := Selector.updateMouseEvent()
*/
func (shared *selectorType) updateMouseEvent() bool {
	isScreenUpdateRequired := false
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	var characterEntry types.CharacterEntryType
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	if characterEntry.AttributeEntry.CellType == constants.CellTypeSelectorItem && eventStateMemory.stateId == constants.EventStateNone && Selectors.IsExists(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias) {
		selectorEntry := Selectors.Get(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
		if buttonPressed != 0 {
			selectorEntry.ItemHighlighted = characterEntry.AttributeEntry.CellControlId
			selectorEntry.ItemSelected = characterEntry.AttributeEntry.CellControlId
			selectorEntry.IsNewItemSelected = true
		} else if !selectorEntry.HighlightOnClickOnly {
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
			// Only clear highlighting if HighlightOnClickOnly is false
			if !selectorEntry.HighlightOnClickOnly {
				selectorEntry.ItemHighlighted = constants.NullItemSelection
			}
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
