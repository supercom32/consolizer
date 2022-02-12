package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
)

type ScrollBarInstanceType struct {
	layerAlias  string
	scrollBarAlias string
}

func (shared *ScrollBarInstanceType) setScrollValue(value int) {
	scrollBarEntry := memory.ScrollBarMemory[shared.layerAlias][shared.scrollBarAlias]
	scrollBarEntry.ScrollValue = value
	computeScrollBarHandlePositionByScrollValue(shared.layerAlias, shared.scrollBarAlias)
}

func (shared *ScrollBarInstanceType) getScrollValue() int {
	scrollBarEntry := memory.ScrollBarMemory[shared.layerAlias][shared.scrollBarAlias]
	return scrollBarEntry.ScrollValue
}

func (shared *ScrollBarInstanceType) setHandlePosition(positionIndex int) {
	scrollBarEntry := memory.ScrollBarMemory[shared.layerAlias][shared.scrollBarAlias]
	scrollBarEntry.HandlePosition = positionIndex
	computeScrollBarValueByHandlePosition(shared.layerAlias, shared.scrollBarAlias)
}

func AddScrollBar(layerAlias string, scrollBarAlias string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, length int, MaxScrollValue int, ScrollValue int, isHorizontal bool) ScrollBarInstanceType {
	memory.AddScrollBar(layerAlias, scrollBarAlias, styleEntry, xLocation, yLocation, length, MaxScrollValue, ScrollValue, isHorizontal)
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
	for currentKey := range memory.ScrollBarMemory[layerAlias] {
		scrollBarEntry := memory.ScrollBarMemory[layerAlias][currentKey]
		if scrollBarEntry.IsVisible {
			drawScrollBar(&layerEntry, currentKey, scrollBarEntry.StyleEntry, scrollBarEntry.XLocation, scrollBarEntry.YLocation, scrollBarEntry.Length, scrollBarEntry.HandlePosition, scrollBarEntry.IsHorizontal)
		}
	}
}

func computeScrollBarValueByHandlePosition(layerAlias string, scrollBarAlias string) {
	scrollBarEntry := memory.ScrollBarMemory[layerAlias][scrollBarAlias]
	// Make sure the handle position is valid first.
	if scrollBarEntry.HandlePosition >= scrollBarEntry.Length {
		scrollBarEntry.HandlePosition = scrollBarEntry.Length - 1
	}
	if scrollBarEntry.HandlePosition < 0 {
		scrollBarEntry.HandlePosition = 0
	}
	// If you scroll to the last square of a scroll bar, set value to max since that's what a user
	// expects to happen.
	if scrollBarEntry.HandlePosition == scrollBarEntry.Length - 1 {
		scrollBarEntry.ScrollValue = scrollBarEntry.MaxScrollValue
		return
	}
	percentScrolled := float64(scrollBarEntry.HandlePosition) / float64(scrollBarEntry.Length)
	scrollBarEntry.ScrollValue = int(float64(scrollBarEntry.MaxScrollValue) * percentScrolled)
}

func computeScrollBarHandlePositionByScrollValue(layerAlias string, scrollBarAlias string) {
	scrollBarEntry := memory.ScrollBarMemory[layerAlias][scrollBarAlias]
	// Make sure the scroll value is valid first.
	if scrollBarEntry.ScrollValue >= scrollBarEntry.MaxScrollValue {
		scrollBarEntry.ScrollValue = scrollBarEntry.MaxScrollValue
	}
	if scrollBarEntry.ScrollValue < 0 {
		scrollBarEntry.ScrollValue = 0
	}
	percentScrolled := float64(scrollBarEntry.ScrollValue) / float64(scrollBarEntry.MaxScrollValue)
	scrollBarEntry.HandlePosition = int(float64(scrollBarEntry.Length) * percentScrolled)
	// Protect in case drawing over the bar limit.
	if scrollBarEntry.HandlePosition >= scrollBarEntry.Length {
		scrollBarEntry.HandlePosition = scrollBarEntry.Length - 1
	}
}

func updateKeyboardEventScrollBar(keystroke string) bool {
	isScreenUpdateRequired := false
	if eventStateMemory.focusedControlType != constants.CellTypeScrollBar {
		return isScreenUpdateRequired
	}
	scrollBarEntry := memory.ScrollBarMemory[eventStateMemory.focusedLayerAlias][eventStateMemory.focusedControlAlias]
	if keystroke == "up" || keystroke == "left"{
		scrollBarEntry.ScrollValue = scrollBarEntry.ScrollValue - 1
		computeScrollBarHandlePositionByScrollValue(eventStateMemory.focusedLayerAlias, eventStateMemory.focusedControlAlias)
	}
	if keystroke == "down" || keystroke == "right" {
		scrollBarEntry.ScrollValue = scrollBarEntry.ScrollValue + 1
		computeScrollBarHandlePositionByScrollValue(eventStateMemory.focusedLayerAlias, eventStateMemory.focusedControlAlias)
	}
	if keystroke == "pgup" {
		scrollBarEntry.ScrollValue = scrollBarEntry.ScrollValue - 10
		computeScrollBarHandlePositionByScrollValue(eventStateMemory.focusedLayerAlias, eventStateMemory.focusedControlAlias)
	}
	if keystroke == "pgdn" {
		scrollBarEntry.ScrollValue = scrollBarEntry.ScrollValue + 10
		computeScrollBarHandlePositionByScrollValue(eventStateMemory.focusedLayerAlias, eventStateMemory.focusedControlAlias)
	}
	return isScreenUpdateRequired
}

func updateMouseEventScrollBar() bool {
	isScreenUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	previousMouseXLocation, previousMouseYLocation, previousButtonPressed, _ := memory.GetPreviousMouseStatus()
	if buttonPressed != 0 {
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if previousButtonPressed == 0 && characterEntry.AttributeEntry.CellType == constants.CellTypeScrollBar {
			if characterEntry.AttributeEntry.CellControlId == constants.CellControlIdScrollBarHandle {
				// If you click on a scroll bar handle, start the scrolling event.
				eventStateMemory.stateId = constants.EventStateDragAndDropScrollBar
			} else if characterEntry.AttributeEntry.CellControlId == constants.CellControlIdUpScrollArrow {
				// If you click on the up scroll bar button.
				scrollBarEntry := memory.ScrollBarMemory[characterEntry.LayerAlias][characterEntry.AttributeEntry.CellControlAlias]
				scrollBarEntry.ScrollValue = scrollBarEntry.ScrollValue - 1
				computeScrollBarHandlePositionByScrollValue(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			} else if characterEntry.AttributeEntry.CellControlId == constants.CellControlIdDownScrollArrow {
				// If you click on the down scroll bar button.
				scrollBarEntry := memory.ScrollBarMemory[characterEntry.LayerAlias][characterEntry.AttributeEntry.CellControlAlias]
				scrollBarEntry.ScrollValue = scrollBarEntry.ScrollValue + 1
				computeScrollBarHandlePositionByScrollValue( characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			} else {
				// If you click on the scroll bar area itself, jump the scroll bar to it.
				scrollBarEntry := memory.ScrollBarMemory[characterEntry.LayerAlias][characterEntry.AttributeEntry.CellControlAlias]
				scrollBarEntry.HandlePosition = characterEntry.AttributeEntry.CellControlId
				computeScrollBarValueByHandlePosition(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			}
			eventStateMemory.focusedLayerAlias = characterEntry.LayerAlias
			eventStateMemory.focusedControlAlias = characterEntry.AttributeEntry.CellControlAlias
			eventStateMemory.focusedControlType = constants.CellTypeScrollBar
		}
		if previousButtonPressed != 0 && eventStateMemory.stateId == constants.EventStateDragAndDropScrollBar {
			xMove := mouseXLocation - previousMouseXLocation
			yMove := mouseYLocation - previousMouseYLocation
			scrollBarEntry := memory.ScrollBarMemory[eventStateMemory.focusedLayerAlias][eventStateMemory.focusedControlAlias]
			if scrollBarEntry.IsHorizontal {
				scrollBarEntry.HandlePosition = scrollBarEntry.HandlePosition + xMove
			} else {
				scrollBarEntry.HandlePosition = scrollBarEntry.HandlePosition + yMove
			}
			computeScrollBarValueByHandlePosition(eventStateMemory.focusedLayerAlias, eventStateMemory.focusedControlAlias)
			isScreenUpdateRequired = true
		}
	} else {
		eventStateMemory.stateId = 0
	}
	return isScreenUpdateRequired
}