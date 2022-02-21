package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"sort"
)

type ScrollBarInstanceType struct {
	layerAlias  string
	scrollBarAlias string
}

func (shared *ScrollBarInstanceType) setScrollValue(value int) {
	scrollBarEntry := memory. GetScrollBar(shared.layerAlias, shared.scrollBarAlias)
	scrollBarEntry.ScrollValue = value
	computeScrollBarHandlePositionByScrollValue(shared.layerAlias, shared.scrollBarAlias)
}

func (shared *ScrollBarInstanceType) getScrollValue() int {
	scrollBarEntry := memory. GetScrollBar(shared.layerAlias, shared.scrollBarAlias)
	return scrollBarEntry.ScrollValue
}

func (shared *ScrollBarInstanceType) setHandlePosition(positionIndex int) {
	scrollBarEntry := memory. GetScrollBar(shared.layerAlias, shared.scrollBarAlias)
	scrollBarEntry.HandlePosition = positionIndex
	computeScrollBarValueByHandlePosition(shared.layerAlias, shared.scrollBarAlias)
}

func AddScrollBar(layerAlias string, scrollBarAlias string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, length int, maxScrollValue int, scrollValue int, scrollIncrement int, isHorizontal bool) ScrollBarInstanceType {
	memory.AddScrollBar(layerAlias, scrollBarAlias, styleEntry, xLocation, yLocation, length, maxScrollValue, scrollValue, scrollIncrement, isHorizontal)
	var ScrollBarInstance ScrollBarInstanceType
	ScrollBarInstance.layerAlias = layerAlias
	ScrollBarInstance.scrollBarAlias = scrollBarAlias
	return ScrollBarInstance
}

func DeleteScrollBar(layerAlias string, scrollBarAlias string) {
	memory.DeleteScrollBar(layerAlias, scrollBarAlias)
}

func drawScrollBarsOnLayer(layerEntry memory.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	keyList := make([]string, 0)
	for currentKey := range memory.ScrollBarMemory[layerAlias] {
		keyList = append(keyList, currentKey)
	}
	sort.Strings(keyList)
	for currentKey := range keyList{
		scrollBarEntry := memory.GetScrollBar(layerAlias, keyList[currentKey])
		if scrollBarEntry.IsVisible {
			drawScrollBar(&layerEntry, keyList[currentKey], scrollBarEntry.StyleEntry, scrollBarEntry.XLocation, scrollBarEntry.YLocation, scrollBarEntry.Length, scrollBarEntry.HandlePosition, scrollBarEntry.IsHorizontal)
		}
	}
}

func DrawScrollBar(layerAlias string, scrollBarAlias string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, height int, handlePosition int, isHorizontal bool) {
	layerEntry := memory.GetLayer(layerAlias)
	drawScrollBar(layerEntry, scrollBarAlias, styleEntry, xLocation, yLocation, height, handlePosition, isHorizontal)
}

func drawScrollBar(layerEntry *memory.LayerEntryType, scrollBarAlias string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, length int, handlePosition int, isHorizontal bool) {
	attributeEntry := memory.NewAttributeEntry()
	attributeEntry.CellType = constants.CellTypeScrollBar
	attributeEntry.CellControlAlias = scrollBarAlias
	attributeEntry.ForegroundColor = styleEntry.ScrollBarForegroundColor
	attributeEntry.BackgroundColor = styleEntry.ScrollBarBackgroundColor
	scrollBarEntry := memory.GetScrollBar(layerEntry.LayerAlias, scrollBarAlias)
	//numberOfScrollSegments := length - 2
	//segmentPosition := math.RoundToEven(float64(currentValue) / float64(numberOfTicks) * float64(numberOfScrollSegments))
	if isHorizontal {
		for currentXLocation := 1; currentXLocation < length - 1; currentXLocation++ {
			attributeEntry.CellControlId = currentXLocation - 1
			printLayer(layerEntry, attributeEntry, xLocation + currentXLocation, yLocation, []rune{styleEntry.ScrollBarTrackPattern})
		}
		attributeEntry.CellControlId = constants.CellControlIdUpScrollArrow
		printLayer(layerEntry, attributeEntry, xLocation, yLocation, []rune{styleEntry.ScrollBarLeftArrow})
		attributeEntry.CellControlId = constants.CellControlIdDownScrollArrow
		printLayer(layerEntry, attributeEntry, xLocation + length - 1, yLocation, []rune{styleEntry.ScrollBarRightArrow})
		attributeEntry.ForegroundColor = styleEntry.ScrollBarHandleColor
		attributeEntry.CellControlId = constants.CellControlIdScrollBarHandle
		// Here we add 1 to the xLocation since handle bars cannot be drawn on scroll arrows.
		if scrollBarEntry.IsEnabled {
			printLayer(layerEntry, attributeEntry, xLocation + 1 + handlePosition, yLocation, []rune{styleEntry.ScrollBarHandle})
		}
	} else {
		for currentYLocation := 1; currentYLocation < length - 1; currentYLocation++ {
			attributeEntry.CellControlId = currentYLocation - 1 // make all Ids 0 based.
			printLayer(layerEntry, attributeEntry, xLocation, yLocation + currentYLocation, []rune{styleEntry.ScrollBarTrackPattern})
		}
		attributeEntry.CellControlId = constants.CellControlIdUpScrollArrow
		printLayer(layerEntry, attributeEntry, xLocation, yLocation, []rune{styleEntry.ScrollBarUpArrow})
		attributeEntry.CellControlId = constants.CellControlIdDownScrollArrow
		printLayer(layerEntry, attributeEntry, xLocation, yLocation + length - 1, []rune{styleEntry.ScrollBarDownArrow})
		attributeEntry.ForegroundColor = styleEntry.ScrollBarHandleColor
		attributeEntry.CellControlId = constants.CellControlIdScrollBarHandle
		if scrollBarEntry.IsEnabled {
			printLayer(layerEntry, attributeEntry, xLocation, yLocation+1+handlePosition, []rune{styleEntry.ScrollBarHandle})
		}
	}
}

func computeScrollBarValueByHandlePosition(layerAlias string, scrollBarAlias string) {
	scrollBarEntry := memory.GetScrollBar(layerAlias, scrollBarAlias)
	// If instructed not to draw scroll bars, do not compute values.
	if scrollBarEntry.IsEnabled == false {
		return
	}
	// Make sure the handle position is valid first. We minus 3 to the length to account for the two arrows
	// and the handle itself.
	if scrollBarEntry.HandlePosition >= scrollBarEntry.Length - 3 {
		scrollBarEntry.HandlePosition = scrollBarEntry.Length - 3
	}
	if scrollBarEntry.HandlePosition < 0 {
		scrollBarEntry.HandlePosition = 0
	}
	// If you scroll to the last square of a scroll bar, set value to max since that's what a user
	// expects to happen.
	if scrollBarEntry.HandlePosition == scrollBarEntry.Length - 3 {
		scrollBarEntry.ScrollValue = scrollBarEntry.MaxScrollValue
		return
	}
	percentScrolled := float64(scrollBarEntry.HandlePosition) / float64(scrollBarEntry.Length)
	scrollBarEntry.ScrollValue = int(float64(scrollBarEntry.MaxScrollValue) * percentScrolled)
}

func computeScrollBarHandlePositionByScrollValue(layerAlias string, scrollBarAlias string) {
	scrollBarEntry := memory.GetScrollBar(layerAlias, scrollBarAlias)
	// If instructed not to draw scroll bars, do not compute values.
	if scrollBarEntry.IsEnabled == false {
		return
	}
	// Make sure the scroll value is valid first.
	if scrollBarEntry.ScrollValue >= scrollBarEntry.MaxScrollValue {
		scrollBarEntry.ScrollValue = scrollBarEntry.MaxScrollValue
	}
	if scrollBarEntry.ScrollValue < 0 {
		scrollBarEntry.ScrollValue = 0
	}
	percentScrolled := float64(scrollBarEntry.ScrollValue) / float64(scrollBarEntry.MaxScrollValue)
	scrollBarEntry.HandlePosition = int(float64(scrollBarEntry.Length - 3) * percentScrolled)
	// Protect in case drawing over the bar limit.
	if scrollBarEntry.HandlePosition >= scrollBarEntry.Length {
		scrollBarEntry.HandlePosition = scrollBarEntry.Length - 3
	}
}

func updateKeyboardEventScrollBar(keystroke string) bool {
	isScreenUpdateRequired := false
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	if focusedControlType != constants.CellTypeScrollBar {
		return isScreenUpdateRequired
	}
	// Check for scrollbar input only if the scroll bar is not disabled (not null).
	scrollBarEntry := memory.GetScrollBar(focusedLayerAlias, focusedControlAlias)
	if scrollBarEntry.IsEnabled {
		if keystroke == "up" || keystroke == "left"{
			scrollBarEntry.ScrollValue = scrollBarEntry.ScrollValue - scrollBarEntry.ScrollIncrement
			computeScrollBarHandlePositionByScrollValue(focusedLayerAlias, focusedControlAlias)
		}
		if keystroke == "down" || keystroke == "right" {
			scrollBarEntry.ScrollValue = scrollBarEntry.ScrollValue + scrollBarEntry.ScrollIncrement
			computeScrollBarHandlePositionByScrollValue(focusedLayerAlias, focusedControlAlias)
		}
		if keystroke == "pgup" {
			scrollBarEntry.ScrollValue = scrollBarEntry.ScrollValue - (scrollBarEntry.ScrollIncrement * 3)
			computeScrollBarHandlePositionByScrollValue(focusedLayerAlias, focusedControlAlias)
		}
		if keystroke == "pgdn" {
			scrollBarEntry.ScrollValue = scrollBarEntry.ScrollValue + (scrollBarEntry.ScrollIncrement * 3)
			computeScrollBarHandlePositionByScrollValue(focusedLayerAlias, focusedControlAlias)
		}
	}
	return isScreenUpdateRequired
}

func updateMouseEventScrollBar() bool {
	isScreenUpdateRequired := false
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	previousMouseXLocation, previousMouseYLocation, previousButtonPressed, _ := memory.GetPreviousMouseStatus()
	if buttonPressed != 0 {
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if previousButtonPressed == 0 && characterEntry.AttributeEntry.CellType == constants.CellTypeScrollBar {
			scrollBarEntry := memory.GetScrollBar(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			// Check for scrollbar input only if the scroll bar is not disabled (not null).
			if scrollBarEntry.IsEnabled {
				if characterEntry.AttributeEntry.CellControlId == constants.CellControlIdScrollBarHandle {
					// If you click on a scroll bar handle, start the scrolling event.
					eventStateMemory.stateId = constants.EventStateDragAndDropScrollBar
				} else if characterEntry.AttributeEntry.CellControlId == constants.CellControlIdUpScrollArrow {
					// If you click on the up scroll bar button.
					scrollBarEntry.ScrollValue = scrollBarEntry.ScrollValue - scrollBarEntry.ScrollIncrement
					computeScrollBarHandlePositionByScrollValue(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
				} else if characterEntry.AttributeEntry.CellControlId == constants.CellControlIdDownScrollArrow {
					// If you click on the down scroll bar button.
					scrollBarEntry.ScrollValue = scrollBarEntry.ScrollValue + scrollBarEntry.ScrollIncrement
					computeScrollBarHandlePositionByScrollValue( characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
				} else {
					// If you click on the scroll bar area itself,  jump the scroll bar to it.
					scrollBarEntry.HandlePosition = characterEntry.AttributeEntry.CellControlId
					computeScrollBarValueByHandlePosition(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
				}
			}
			setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeScrollBar)
			isScreenUpdateRequired = true
		} else if previousButtonPressed != 0 && eventStateMemory.stateId == constants.EventStateDragAndDropScrollBar {
			xMove := mouseXLocation - previousMouseXLocation
			yMove := mouseYLocation - previousMouseYLocation
			scrollBarEntry := memory.GetScrollBar(focusedLayerAlias, focusedControlAlias)
			if scrollBarEntry.IsHorizontal {
				scrollBarEntry.HandlePosition = scrollBarEntry.HandlePosition + xMove
			} else {
				scrollBarEntry.HandlePosition = scrollBarEntry.HandlePosition + yMove
			}
			computeScrollBarValueByHandlePosition(focusedLayerAlias, focusedControlAlias)
			isScreenUpdateRequired = true
		}
	} else {
		eventStateMemory.stateId = constants.EventStateNone
	}
	return isScreenUpdateRequired
}