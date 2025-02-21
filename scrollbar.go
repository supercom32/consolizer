package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/types"
)

type ScrollbarInstanceType struct {
	layerAlias   string
	controlAlias string
}

type scrollbarType struct{}

var scrollbar scrollbarType

func (shared *ScrollbarInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeScrollbar)
}

func (shared *ScrollbarInstanceType) Delete() *ScrollbarInstanceType {
	if memory.IsScrollbarExists(shared.layerAlias, shared.controlAlias) {
		memory.DeleteScrollbar(shared.layerAlias, shared.controlAlias)
	}
	return nil
}

/*
setScrollValue allows you set the current scroll value of a scroll bar.  In addition, the following
information should be noted:

- The scroll bar handle position will automatically be updated to reflect your new scroll bar value.

- If the scroll bar value specified is out of range, then the closes minimum/maximum value will be selected instead.

- If the scroll bar instance does not exist, then the request is ignored.
*/
func (shared *ScrollbarInstanceType) setScrollValue(value int) {
	// TODO: Add scroll value validation.
	if memory.IsScrollbarExists(shared.layerAlias, shared.controlAlias) {
		scrollbarEntry := memory.GetScrollbar(shared.layerAlias, shared.controlAlias)
		scrollbarEntry.ScrollValue = value
		scrollbar.computeScrollbarHandlePositionByScrollValue(shared.layerAlias, shared.controlAlias)
	}
}

/*
getScrollValue allows you to obtain the scroll bar value for a given scrollbar. If the scroll bar instance
no longer exists, then a result of 0 is always returned.
*/
func (shared *ScrollbarInstanceType) getScrollValue() int {
	if memory.IsScrollbarExists(shared.layerAlias, shared.controlAlias) {
		scrollbarEntry := memory.GetScrollbar(shared.layerAlias, shared.controlAlias)
		return scrollbarEntry.ScrollValue
	}
	return 0
}

/*
setHandlePosition allows you to specify the location of where the scrollbar handle should be.
The scrollbar value is automatically updated to match the location of the scrollbar handle position.
*/
func (shared *ScrollbarInstanceType) setHandlePosition(positionIndex int) {
	if memory.IsScrollbarExists(shared.layerAlias, shared.controlAlias) {
		scrollbarEntry := memory.GetScrollbar(shared.layerAlias, shared.controlAlias)
		scrollbarEntry.HandlePosition = positionIndex
		scrollbar.computeScrollbarValueByHandlePosition(shared.layerAlias, shared.controlAlias)
	}
}

/*
Add allows you to add a scrollbar to a given text layer. Once called, an instance of your
control is returned which will allow you to read or manipulate the properties for it. The Style
of the scrollbar will be determined by the style entry passed in. If you wish to remove a scrollbar
from a text layer, simply call 'DeleteScrollbar'. In addition, the following information should be noted:

- Scrollbars are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create scrollbars without actually overwriting
the text layer data under it.

- If the scrollbar to be drawn falls outside the range of the provided layer,
then only the visible portion of the scrollbar will be drawn.
*/
func (shared *scrollbarType) Add(layerAlias string, scrollbarAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, length int, maxScrollValue int, scrollValue int, scrollIncrement int, isHorizontal bool) ScrollbarInstanceType {
	// TODO: add validation and what happens if failed.
	memory.AddScrollbar(layerAlias, scrollbarAlias, styleEntry, xLocation, yLocation, length, maxScrollValue, scrollValue, scrollIncrement, isHorizontal)
	var ScrollbarInstance ScrollbarInstanceType
	ScrollbarInstance.layerAlias = layerAlias
	ScrollbarInstance.controlAlias = scrollbarAlias
	return ScrollbarInstance
}

/*
DeleteScrollbar allows you to remove a scrollbar from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a scrollbar which does not exist, then the request
will simply be ignored.
*/
func (shared *scrollbarType) DeleteScrollbar(layerAlias string, scrollbarAlias string) {
	memory.DeleteScrollbar(layerAlias, scrollbarAlias)
}

func (shared *scrollbarType) DeleteAllScrollbars(layerAlias string) {
	memory.DeleteAllScrollbarsFromLayer(layerAlias)
}

/*
drawScrollbarsOnLayer allows you to draw all scrollbars on a given text layer.
In addition, the following information should be noted:

- We sort the scroll bars since internally generated scrollbars will have the prefix "zzz" so that
it appears last on the rendering order.
*/
func (shared *scrollbarType) drawScrollbarsOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	compareByAlias := func(a, b *types.ScrollbarEntryType) bool {
		return a.ScrollBarAlias < b.ScrollBarAlias
	}
	for _, currentKey := range memory.ScrollBars.SortEntries(layerAlias, true, compareByAlias) {
		scrollbarEntry := currentKey
		if scrollbarEntry.IsVisible {
			shared.drawScrollbar(&layerEntry, scrollbarEntry.ScrollBarAlias, scrollbarEntry.StyleEntry, scrollbarEntry.XLocation, scrollbarEntry.YLocation, scrollbarEntry.Length, scrollbarEntry.HandlePosition, scrollbarEntry.IsHorizontal)
		}
	}
}

/*
drawScrollbar allows you to draw A scrollbar on a given text layer. The
Style of the scrollbar will be determined by the style entry passed in. In
addition, the following information should be noted:

- Scrollbars are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create scrollbars without actually overwriting
the text layer data under it.

- If the scrollbar to be drawn falls outside the range of the provided layer,
then only the visible portion of the scrollbar will be drawn.
*/
func (shared *scrollbarType) drawScrollbar(layerEntry *types.LayerEntryType, scrollbarAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, length int, handlePosition int, isHorizontal bool) {
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.CellType = constants.CellTypeScrollbar
	attributeEntry.CellControlAlias = scrollbarAlias
	attributeEntry.ForegroundColor = styleEntry.ScrollbarForegroundColor
	attributeEntry.BackgroundColor = styleEntry.ScrollbarBackgroundColor
	scrollbarEntry := memory.GetScrollbar(layerEntry.LayerAlias, scrollbarAlias)
	// numberOfScrollSegments := length - 2
	// segmentPosition := math.RoundToEven(float64(currentValue) / float64(numberOfTicks) * float64(numberOfScrollSegments))
	if isHorizontal {
		for currentXLocation := 1; currentXLocation < length-1; currentXLocation++ {
			attributeEntry.CellControlId = currentXLocation - 1
			printLayer(layerEntry, attributeEntry, xLocation+currentXLocation, yLocation, []rune{styleEntry.ScrollbarTrackPattern})
		}
		attributeEntry.CellControlId = constants.CellControlIdUpScrollArrow
		printLayer(layerEntry, attributeEntry, xLocation, yLocation, []rune{styleEntry.ScrollbarLeftArrow})
		attributeEntry.CellControlId = constants.CellControlIdDownScrollArrow
		printLayer(layerEntry, attributeEntry, xLocation+length-1, yLocation, []rune{styleEntry.ScrollbarRightArrow})
		attributeEntry.ForegroundColor = styleEntry.ScrollbarHandleColor
		attributeEntry.CellControlId = constants.CellControlIdScrollbarHandle
		// Here we add 1 to the xLocation since handle bars cannot be drawn on scroll arrows.
		if scrollbarEntry.IsEnabled {
			printLayer(layerEntry, attributeEntry, xLocation+1+handlePosition, yLocation, []rune{styleEntry.ScrollbarHandle})
		}
	} else {
		for currentYLocation := 1; currentYLocation < length-1; currentYLocation++ {
			attributeEntry.CellControlId = currentYLocation - 1 // make all Ids 0 based.
			printLayer(layerEntry, attributeEntry, xLocation, yLocation+currentYLocation, []rune{styleEntry.ScrollbarTrackPattern})
		}
		attributeEntry.CellControlId = constants.CellControlIdUpScrollArrow
		printLayer(layerEntry, attributeEntry, xLocation, yLocation, []rune{styleEntry.ScrollbarUpArrow})
		attributeEntry.CellControlId = constants.CellControlIdDownScrollArrow
		printLayer(layerEntry, attributeEntry, xLocation, yLocation+length-1, []rune{styleEntry.ScrollbarDownArrow})
		attributeEntry.ForegroundColor = styleEntry.ScrollbarHandleColor
		attributeEntry.CellControlId = constants.CellControlIdScrollbarHandle
		if scrollbarEntry.IsEnabled {
			printLayer(layerEntry, attributeEntry, xLocation, yLocation+1+handlePosition, []rune{styleEntry.ScrollbarHandle})
		}
	}
}

/*
computeScrollbarValueByHandlePosition allows you to compute the scrollbar value based on the position of the
scrollbar handle.
*/
func (shared *scrollbarType) computeScrollbarValueByHandlePosition(layerAlias string, scrollbarAlias string) {
	scrollbarEntry := memory.GetScrollbar(layerAlias, scrollbarAlias)
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
computeScrollbarHandlePositionByScrollValue allows you to calculate the position of the scrollbar handle
based on the current scrollbar value.
*/
func (shared *scrollbarType) computeScrollbarHandlePositionByScrollValue(layerAlias string, scrollbarAlias string) {
	scrollbarEntry := memory.GetScrollbar(layerAlias, scrollbarAlias)
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
updateKeyboardEventScrollbar allows you to update the state of all scrollbars according to the current keystroke event.
In the event that a screen update is required this method returns true.
*/
func (shared *scrollbarType) updateKeyboardEventScrollbar(keystroke []rune) bool {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := false
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	if focusedControlType != constants.CellTypeScrollbar || !memory.IsScrollbarExists(focusedLayerAlias, focusedControlAlias) {
		return isScreenUpdateRequired
	}
	// Check for scrollbar input only if the scroll bar is not disabled (not null).
	scrollbarEntry := memory.GetScrollbar(focusedLayerAlias, focusedControlAlias)
	if scrollbarEntry.IsEnabled {
		if keystrokeAsString == "up" || keystrokeAsString == "left" {
			scrollbarEntry.ScrollValue = scrollbarEntry.ScrollValue - scrollbarEntry.ScrollIncrement
			shared.computeScrollbarHandlePositionByScrollValue(focusedLayerAlias, focusedControlAlias)
		}
		if keystrokeAsString == "down" || keystrokeAsString == "right" {
			scrollbarEntry.ScrollValue = scrollbarEntry.ScrollValue + scrollbarEntry.ScrollIncrement
			shared.computeScrollbarHandlePositionByScrollValue(focusedLayerAlias, focusedControlAlias)
		}
		if keystrokeAsString == "pgup" {
			scrollbarEntry.ScrollValue = scrollbarEntry.ScrollValue - (scrollbarEntry.ScrollIncrement * 3)
			shared.computeScrollbarHandlePositionByScrollValue(focusedLayerAlias, focusedControlAlias)
		}
		if keystrokeAsString == "pgdn" {
			scrollbarEntry.ScrollValue = scrollbarEntry.ScrollValue + (scrollbarEntry.ScrollIncrement * 3)
			shared.computeScrollbarHandlePositionByScrollValue(focusedLayerAlias, focusedControlAlias)
		}
	}
	return isScreenUpdateRequired
}

/*
updateMouseEventScrollbar allows you to update the state of all scrollbars according to the current mouse event state.
In the event that a screen update is required this method returns true.
*/
func (shared *scrollbarType) updateMouseEventScrollbar() bool {
	isScreenUpdateRequired := false
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	previousMouseXLocation, previousMouseYLocation, previousButtonPressed, _ := memory.GetPreviousMouseStatus()
	if buttonPressed != 0 {
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if previousButtonPressed == 0 && characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar {
			scrollbarEntry := memory.GetScrollbar(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			// Check for scrollbar input only if the scroll bar is not disabled (not null).
			if scrollbarEntry.IsEnabled {
				if characterEntry.AttributeEntry.CellControlId == constants.CellControlIdScrollbarHandle {
					// If you click on a scroll bar handle, start the scrolling event.
					eventStateMemory.stateId = constants.EventStateDragAndDropScrollbar
				} else if characterEntry.AttributeEntry.CellControlId == constants.CellControlIdUpScrollArrow {
					// If you click on the up scroll bar buttonType.
					scrollbarEntry.ScrollValue = scrollbarEntry.ScrollValue - scrollbarEntry.ScrollIncrement
					shared.computeScrollbarHandlePositionByScrollValue(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
				} else if characterEntry.AttributeEntry.CellControlId == constants.CellControlIdDownScrollArrow {
					// If you click on the down scroll bar buttonType.
					scrollbarEntry.ScrollValue = scrollbarEntry.ScrollValue + scrollbarEntry.ScrollIncrement
					shared.computeScrollbarHandlePositionByScrollValue(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
				} else {
					// If you click on the scroll bar area itself,  jump the scroll bar to it.
					scrollbarEntry.HandlePosition = characterEntry.AttributeEntry.CellControlId
					shared.computeScrollbarValueByHandlePosition(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
				}
			}
			setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeScrollbar)
			isScreenUpdateRequired = true
		} else if previousButtonPressed != 0 && eventStateMemory.stateId == constants.EventStateDragAndDropScrollbar {
			xMove := mouseXLocation - previousMouseXLocation
			yMove := mouseYLocation - previousMouseYLocation
			if focusedControlType == constants.CellTypeScrollbar {
				scrollbarEntry := memory.GetScrollbar(focusedLayerAlias, focusedControlAlias)
				if scrollbarEntry.IsHorizontal {
					scrollbarEntry.HandlePosition = scrollbarEntry.HandlePosition + xMove
				} else {
					scrollbarEntry.HandlePosition = scrollbarEntry.HandlePosition + yMove
				}
				shared.computeScrollbarValueByHandlePosition(focusedLayerAlias, focusedControlAlias)
				isScreenUpdateRequired = true
			}
		}
	} else {
		eventStateMemory.stateId = constants.EventStateNone
	}
	return isScreenUpdateRequired
}
