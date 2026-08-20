package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/types"
)

/*
ScrollbarInstanceType is a structure which represents an instance of a scrollbar control.
*/
type ScrollbarInstanceType struct {
	BaseControlInstanceType
}

type scrollbarType struct{}

var scrollbar scrollbarType
var ScrollBars = memory.NewControlMemoryManager[types.ScrollbarEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
AddToTabIndex is a method which allows you to add the scroll bar to the tab index. In addition, the following should be
noted:

- This enables the scroll bar to be selected and interacted with when the user cycles through controls using the tab.

Example:
    scrollbar.AddToTabIndex()
*/
func (shared *ScrollbarInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeScrollbar)
}

/*
Delete is a method which allows you to delete the scroll bar instance. In addition, the following should be noted:

- All memory associated with the scroll bar will be freed and it will no longer be rendered.

Example:
    scrollbar = scrollbar.Delete()
*/
func (shared *ScrollbarInstanceType) Delete() *ScrollbarInstanceType {
	shared.BaseControlInstanceType.Delete()
	return nil
}

/*
setScrollValue is a method which allows you to set the current scroll value of a scroll bar. In addition, the following
should be noted:

- The scroll bar handle position will automatically be updated to reflect your new scroll bar value.

- If the scroll bar value specified is out of range, then the closest minimum or maximum value will be selected.

- If the scroll bar instance does not exist, then the request is ignored.

Example:
    scrollbar.setScrollValue(10)
*/
func (shared *ScrollbarInstanceType) setScrollValue(value int) {
	// TODO: AddLayer scroll value validation.
	if ScrollBars.IsExists(shared.layerAlias, shared.controlAlias) {
		scrollbarEntry := ScrollBars.Get(shared.layerAlias, shared.controlAlias)
		scrollbarEntry.ScrollValue = value
		scrollbar.computeHandlePositionByScrollValue(shared.layerAlias, shared.controlAlias)
	}
}

/*
getScrollValue is a method which allows you to obtain the scroll bar value for a given scroll bar. In addition, the
following should be noted:

- If the scroll bar instance no longer exists, then a result of 0 is always returned.

Example:
    value := scrollbar.getScrollValue()
*/
func (shared *ScrollbarInstanceType) getScrollValue() int {
	if ScrollBars.IsExists(shared.layerAlias, shared.controlAlias) {
		scrollbarEntry := ScrollBars.Get(shared.layerAlias, shared.controlAlias)
		return scrollbarEntry.ScrollValue
	}
	return 0
}

/*
setHandlePosition is a method which allows you to specify the location of where the scroll bar handle should be. In
addition, the following should be noted:

- The scroll bar value is automatically updated to match the location of the scroll bar handle position.

Example:
    scrollbar.setHandlePosition(5)
*/
func (shared *ScrollbarInstanceType) setHandlePosition(positionIndex int) {
	if ScrollBars.IsExists(shared.layerAlias, shared.controlAlias) {
		scrollbarEntry := ScrollBars.Get(shared.layerAlias, shared.controlAlias)
		scrollbarEntry.HandlePosition = positionIndex
		scrollbar.computeValueByHandlePosition(shared.layerAlias, shared.controlAlias)
	}
}

/*
Add is a method which allows you to add a scroll bar to a given text layer. Once called, an instance of your control is
returned which will allow you to read or manipulate the properties for it. The style of the scroll bar will be
determined by the style entry passed in. If you wish to remove a scroll bar from a text layer, simply call
DeleteScrollbar. In addition, the following should be noted:

- Scroll bars are not drawn physically to the text layer provided. Instead, they are rendered to the terminal at the
  same time when the text layer is rendered.

- If the scroll bar to be drawn falls outside the range of the provided layer, then only the visible portion of the
  scroll bar will be drawn.

Example:
    scrollbarInstance := scrollbar.Add("Layer1", "Scroll1", style, 0, 0, 10, 100, 0, 1, false)
*/
func (shared *scrollbarType) Add(layerAlias string, scrollbarAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, length int, maxScrollValue int, scrollValue int, scrollIncrement int, isHorizontal bool) ScrollbarInstanceType {
	scrollbarEntry := types.NewScrollbarEntry()
	scrollbarEntry.Alias = scrollbarAlias
	scrollbarEntry.StyleEntry = styleEntry
	scrollbarEntry.XLocation = xLocation
	scrollbarEntry.YLocation = yLocation
	scrollbarEntry.Length = length
	scrollbarEntry.MaxScrollValue = maxScrollValue - 1 // Adjusted for 0-based indexing
	scrollbarEntry.ScrollValue = scrollValue
	scrollbarEntry.IsVisible = true
	scrollbarEntry.IsEnabled = true
	scrollbarEntry.IsHorizontal = isHorizontal
	scrollbarEntry.ScrollIncrement = scrollIncrement
	scrollbarEntry.ParentControlAlias = "" // Empty for standalone scrollbars
	scrollbarEntry.ParentControlType = 0   // No parent control type by default
	// Use the generic memory manager to add the scrollbar entry
	ScrollBars.Add(layerAlias, scrollbarAlias, &scrollbarEntry)
	// TODO: add validation and what happens if failed.
	var scrollbarInstance ScrollbarInstanceType
	scrollbarInstance.layerAlias = layerAlias
	scrollbarInstance.controlAlias = scrollbarAlias
	scrollbarInstance.controlType = constants.TYPE_SCROLLBAR
	scrollbarInstance.setScrollValue(scrollValue)
	return scrollbarInstance
}

/*
Delete is a method which allows you to remove a scroll bar from a text layer. In addition, the following should be
noted:

- If you attempt to delete a scroll bar which does not exist, then the request will simply be ignored.

Example:
    scrollbar.Delete("Layer1", "Scroll1")
*/
func (shared *scrollbarType) Delete(layerAlias string, scrollbarAlias string) {
	ScrollBars.Remove(layerAlias, scrollbarAlias)
}

/*
DeleteAll is a method which allows you to delete all scroll bars on a specified layer. In addition, the following
should be noted:

- All memory associated with the scroll bars on the layer will be freed.

Example:
    scrollbar.DeleteAll("Layer1")
*/
func (shared *scrollbarType) DeleteAll(layerAlias string) {
	ScrollBars.RemoveAll(layerAlias)
}

/*
drawOnLayer is a method which allows you to draw all scroll bars on a given text layer. In addition, the following should
be noted:

- We sort the scroll bars since internally generated scroll bars will have the prefix "zzz" so that it appears last on
  the display.

Example:
    scrollbar.drawOnLayer(layerEntry)
*/
func (shared *scrollbarType) drawOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	compareByAlias := func(a, b *types.ScrollbarEntryType) bool {
		return a.Alias < b.Alias
	}
	for _, currentKey := range ScrollBars.SortEntries(layerAlias, true, compareByAlias) {
		scrollbarEntry := currentKey
		if scrollbarEntry.IsVisible {
			if scrollbarEntry.ParentControlAlias == "" {
				shared.draw(&layerEntry, scrollbarEntry.Alias, scrollbarEntry.StyleEntry, scrollbarEntry.XLocation, scrollbarEntry.YLocation, scrollbarEntry.Length, scrollbarEntry.HandlePosition, scrollbarEntry.IsHorizontal)
			}
		}
	}
}

/*
drawOnLayerByAlias is a method which allows you to draw a specific scroll bar on a given text layer by its alias. In
addition, the following should be noted:

- If the scroll bar does not exist or is not visible, it will not be drawn.

Example:
    scrollbar.drawOnLayerByAlias(&layerEntry, "Scroll1")
*/
func (shared *scrollbarType) drawOnLayerByAlias(layerEntry *types.LayerEntryType, scrollbarAlias string) {
	layerAlias := layerEntry.LayerAlias
	if ScrollBars.IsExists(layerAlias, scrollbarAlias) {
		scrollbarEntry := ScrollBars.Get(layerAlias, scrollbarAlias)
		if scrollbarEntry.IsVisible {
			shared.draw(layerEntry, scrollbarEntry.Alias, scrollbarEntry.StyleEntry, scrollbarEntry.XLocation, scrollbarEntry.YLocation, scrollbarEntry.Length, scrollbarEntry.HandlePosition, scrollbarEntry.IsHorizontal)
		}
	}
}

/*
draw is a method which allows you to draw a scroll bar on a given text layer. The style of the scroll bar will be
determined by the style entry passed in. In addition, the following should be noted:

- Scroll bars are not drawn physically to the text layer provided. Instead, they are rendered to the terminal at the
  same time when the text layer is rendered.

- If the scroll bar to be drawn falls outside the range of the provided layer, then only the visible portion of the
  scroll bar will be drawn.

Example:
    scrollbar.draw(&layerEntry, "Scroll1", style, 0, 0, 10, 5, false)
*/
func (shared *scrollbarType) draw(layerEntry *types.LayerEntryType, scrollbarAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, length int, handlePosition int, isHorizontal bool) {
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.CellType = constants.CellTypeScrollbar
	attributeEntry.CellControlAlias = scrollbarAlias
	attributeEntry.ForegroundColor = styleEntry.Scrollbar.ForegroundColor
	attributeEntry.BackgroundColor = styleEntry.Scrollbar.BackgroundColor
	scrollbarEntry := ScrollBars.Get(layerEntry.LayerAlias, scrollbarAlias)
	// numberOfScrollSegments := length - 2
	// segmentPosition := math.RoundToEven(float64(currentValue) / float64(numberOfTicks) * float64(numberOfScrollSegments))
	if isHorizontal {
		for currentXLocation := 1; currentXLocation < length-1; currentXLocation++ {
			attributeEntry.CellControlId = currentXLocation - 1
			layer.printLayer(layerEntry, attributeEntry, xLocation+currentXLocation, yLocation, []rune{styleEntry.Scrollbar.TrackPattern})
		}
		attributeEntry.CellControlId = constants.CellControlIdUpScrollArrow
		layer.printLayer(layerEntry, attributeEntry, xLocation, yLocation, []rune{styleEntry.Scrollbar.LeftArrow})
		attributeEntry.CellControlId = constants.CellControlIdDownScrollArrow
		layer.printLayer(layerEntry, attributeEntry, xLocation+length-1, yLocation, []rune{styleEntry.Scrollbar.RightArrow})
		attributeEntry.ForegroundColor = styleEntry.Scrollbar.HandleColor
		attributeEntry.CellControlId = constants.CellControlIdScrollbarHandle
		// Here we add 1 to the xLocation since handle bars cannot be drawn on scroll arrows.
		if scrollbarEntry.IsEnabled {
			layer.printLayer(layerEntry, attributeEntry, xLocation+1+handlePosition, yLocation, []rune{styleEntry.Scrollbar.Handle})
		}
	} else {
		for currentYLocation := 1; currentYLocation < length-1; currentYLocation++ {
			attributeEntry.CellControlId = currentYLocation - 1 // make all Ids 0 based.
			layer.printLayer(layerEntry, attributeEntry, xLocation, yLocation+currentYLocation, []rune{styleEntry.Scrollbar.TrackPattern})
		}
		attributeEntry.CellControlId = constants.CellControlIdUpScrollArrow
		layer.printLayer(layerEntry, attributeEntry, xLocation, yLocation, []rune{styleEntry.Scrollbar.UpArrow})
		attributeEntry.CellControlId = constants.CellControlIdDownScrollArrow
		layer.printLayer(layerEntry, attributeEntry, xLocation, yLocation+length-1, []rune{styleEntry.Scrollbar.DownArrow})
		attributeEntry.ForegroundColor = styleEntry.Scrollbar.HandleColor
		attributeEntry.CellControlId = constants.CellControlIdScrollbarHandle
		if scrollbarEntry.IsEnabled {
			layer.printLayer(layerEntry, attributeEntry, xLocation, yLocation+1+handlePosition, []rune{styleEntry.Scrollbar.Handle})
		}
	}
}

/*
computeValueByHandlePosition is a method which allows you to compute the scroll bar value based on the position of the
scroll bar handle. In addition, the following should be noted:

- If the scroll bar is disabled, no computation occurs.

Example:
    scrollbar.computeValueByHandlePosition("Layer1", "Scroll1")
*/
func (shared *scrollbarType) computeValueByHandlePosition(layerAlias string, scrollbarAlias string) {
	scrollbarEntry := ScrollBars.Get(layerAlias, scrollbarAlias)
	// If instructed not to draw scroll bars, do not compute values.
	if scrollbarEntry.IsEnabled == false {
		return
	}
	// Make sure the handle position is valid first. We minus 3 to the length to account for the two arrows
	// and the handle itself.
	if scrollbarEntry.HandlePosition >= scrollbarEntry.Length-3 {
		scrollbarEntry.HandlePosition = scrollbarEntry.Length - 3
	}
	if scrollbarEntry.HandlePosition < 0 {
		scrollbarEntry.HandlePosition = 0
	}
	// If you scroll to the last square of a scroll bar, set value to max since that's what a user
	// expects to happen.
	if scrollbarEntry.HandlePosition == scrollbarEntry.Length-3 {
		scrollbarEntry.ScrollValue = scrollbarEntry.MaxScrollValue
		return
	}
	percentScrolled := float64(scrollbarEntry.HandlePosition) / float64(scrollbarEntry.Length)
	scrollbarEntry.ScrollValue = int(float64(scrollbarEntry.MaxScrollValue) * percentScrolled)
}

/*
computeHandlePositionByScrollValue is a method which allows you to calculate the position of the scroll bar handle based
on the current scroll bar value. In addition, the following should be noted:

- If the scroll bar is disabled, no computation occurs.

Example:
    scrollbar.computeHandlePositionByScrollValue("Layer1", "Scroll1")
*/
func (shared *scrollbarType) computeHandlePositionByScrollValue(layerAlias string, scrollbarAlias string) {
	scrollbarEntry := ScrollBars.Get(layerAlias, scrollbarAlias)
	// If instructed not to draw scroll bars, do not compute values.
	if scrollbarEntry.IsEnabled == false {
		return
	}
	// Make sure the scroll value is valid first.
	if scrollbarEntry.ScrollValue >= scrollbarEntry.MaxScrollValue {
		scrollbarEntry.ScrollValue = scrollbarEntry.MaxScrollValue
	}
	if scrollbarEntry.ScrollValue < 0 {
		scrollbarEntry.ScrollValue = 0
	}
	percentScrolled := float64(0)
	// Protect against divide by zero cases.
	if scrollbarEntry.MaxScrollValue != 0 {
		percentScrolled = float64(scrollbarEntry.ScrollValue) / float64(scrollbarEntry.MaxScrollValue)
	}
	scrollbarEntry.HandlePosition = int(float64(scrollbarEntry.Length-3) * percentScrolled)
	// Protect in case drawing over the bar limit.
	if scrollbarEntry.HandlePosition >= scrollbarEntry.Length {
		scrollbarEntry.HandlePosition = scrollbarEntry.Length - 3
	}
}

/*
updateKeyboardEventManually is a method which allows you to manually update the state of a scroll bar according to a
keystroke event. In addition, the following should be noted:

- This is used when the scroll bar is not the primary focused control but still needs to react to input.

Example:
    updateRequired, consumed := scrollbar.updateKeyboardEventManually("Layer1", "Scroll1", rune("up"))
*/
func (shared *scrollbarType) updateKeyboardEventManually(layerAlias string, scrollbarAlias string, keystroke []rune) (bool, bool) {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := false
	isKeystrokeConsumed := false
	// Check for scrollbar input only if the scroll bar is not disabled (not null).
	scrollbarEntry := ScrollBars.Get(layerAlias, scrollbarAlias)
	if scrollbarEntry.IsEnabled {
		if keystrokeAsString == "up" || keystrokeAsString == "left" {
			scrollbarEntry.ScrollValue = scrollbarEntry.ScrollValue - scrollbarEntry.ScrollIncrement
			shared.computeHandlePositionByScrollValue(layerAlias, scrollbarAlias)
			// Update selector viewport position
			for _, currentSelectorEntry := range Selectors.GetAllEntries(layerAlias) {
				selectorEntry := currentSelectorEntry
				if selectorEntry.ScrollbarAlias == scrollbarAlias {
					selectorEntry.ViewportPosition = scrollbarEntry.ScrollValue
					isScreenUpdateRequired = true
					isKeystrokeConsumed = true
					break
				}
			}
		}
		if keystrokeAsString == "down" || keystrokeAsString == "right" {
			scrollbarEntry.ScrollValue = scrollbarEntry.ScrollValue + scrollbarEntry.ScrollIncrement
			shared.computeHandlePositionByScrollValue(layerAlias, scrollbarAlias)
			// Update selector viewport position
			for _, currentSelectorEntry := range Selectors.GetAllEntries(layerAlias) {
				selectorEntry := currentSelectorEntry
				if selectorEntry.ScrollbarAlias == scrollbarAlias {
					selectorEntry.ViewportPosition = scrollbarEntry.ScrollValue
					isScreenUpdateRequired = true
					isKeystrokeConsumed = true
					break
				}
			}
		}
		if keystrokeAsString == "pgup" {
			scrollbarEntry.ScrollValue = scrollbarEntry.ScrollValue - (scrollbarEntry.ScrollIncrement * 3)
			shared.computeHandlePositionByScrollValue(layerAlias, scrollbarAlias)
			// Update selector viewport position
			for _, currentSelectorEntry := range Selectors.GetAllEntries(layerAlias) {
				selectorEntry := currentSelectorEntry
				if selectorEntry.ScrollbarAlias == scrollbarAlias {
					selectorEntry.ViewportPosition = scrollbarEntry.ScrollValue
					isScreenUpdateRequired = true
					isKeystrokeConsumed = true
					break
				}
			}
		}
		if keystrokeAsString == "pgdn" {
			scrollbarEntry.ScrollValue = scrollbarEntry.ScrollValue + (scrollbarEntry.ScrollIncrement * 3)
			shared.computeHandlePositionByScrollValue(layerAlias, scrollbarAlias)
			// Update selector viewport position
			for _, currentSelectorEntry := range Selectors.GetAllEntries(layerAlias) {
				selectorEntry := currentSelectorEntry
				if selectorEntry.ScrollbarAlias == scrollbarAlias {
					selectorEntry.ViewportPosition = scrollbarEntry.ScrollValue
					isScreenUpdateRequired = true
					isKeystrokeConsumed = true
					break
				}
			}
		}
	}
	return isScreenUpdateRequired, isKeystrokeConsumed
}

/*
updateKeyboardEvent is a method which allows you to update the state of all scroll bars according to the current
keystroke event. In addition, the following should be noted:

- Only the currently focused scroll bar will process the event.

Example:
    updateRequired, consumed := scrollbar.updateKeyboardEvent(rune("down"))
*/
func (shared *scrollbarType) updateKeyboardEvent(keystroke []rune) (bool, bool) {
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	if focusedControlType != constants.CellTypeScrollbar || !ScrollBars.IsExists(focusedLayerAlias, focusedControlAlias) {
		return false, false
	}
	return shared.updateKeyboardEventManually(focusedLayerAlias, focusedControlAlias, keystroke)
}

/*
updateMouseEvent is a method which allows you to update the state of all scroll bars according to the current mouse
event state. In addition, the following should be noted:

- Handles clicking and dragging of the scroll bar handle and arrows.

Example:
    updateRequired := scrollbar.updateMouseEvent()
*/
func (shared *scrollbarType) updateMouseEvent() bool {
	isScreenUpdateRequired := false
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	previousMouseXLocation, previousMouseYLocation, previousButtonPressed, _ := GetPreviousMouseStatus()
	if buttonPressed != 0 {
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if previousButtonPressed == 0 && characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar {
			scrollbarEntry := ScrollBars.Get(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			// Check for scrollbar input only if the scroll bar is not disabled (not null).
			if scrollbarEntry.IsEnabled {
				if characterEntry.AttributeEntry.CellControlId == constants.CellControlIdScrollbarHandle {
					// If you click on a scroll bar handle, start the scrolling event.
					eventStateMemory.stateId = constants.EventStateDragAndDropScrollbar
				} else if characterEntry.AttributeEntry.CellControlId == constants.CellControlIdUpScrollArrow {
					// If you click on the up scroll bar buttonType.
					scrollbarEntry.ScrollValue = scrollbarEntry.ScrollValue - scrollbarEntry.ScrollIncrement
					shared.computeHandlePositionByScrollValue(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
				} else if characterEntry.AttributeEntry.CellControlId == constants.CellControlIdDownScrollArrow {
					// If you click on the down scroll bar buttonType.
					scrollbarEntry.ScrollValue = scrollbarEntry.ScrollValue + scrollbarEntry.ScrollIncrement
					shared.computeHandlePositionByScrollValue(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
				} else {
					// If you click on the scroll bar area itself,  jump the scroll bar to it.
					scrollbarEntry.HandlePosition = characterEntry.AttributeEntry.CellControlId
					shared.computeValueByHandlePosition(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
				}
			}
			setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeScrollbar)
			isScreenUpdateRequired = true
		} else if previousButtonPressed != 0 && eventStateMemory.stateId == constants.EventStateDragAndDropScrollbar {
			xMove := mouseXLocation - previousMouseXLocation
			yMove := mouseYLocation - previousMouseYLocation
			if focusedControlType == constants.CellTypeScrollbar {
				scrollbarEntry := ScrollBars.Get(focusedLayerAlias, focusedControlAlias)
				if scrollbarEntry.IsHorizontal {
					scrollbarEntry.HandlePosition = scrollbarEntry.HandlePosition + xMove
				} else {
					scrollbarEntry.HandlePosition = scrollbarEntry.HandlePosition + yMove
				}
				shared.computeValueByHandlePosition(focusedLayerAlias, focusedControlAlias)
				isScreenUpdateRequired = true
			}
		}
	} else {
		eventStateMemory.stateId = constants.EventStateNone
	}
	return isScreenUpdateRequired
}
