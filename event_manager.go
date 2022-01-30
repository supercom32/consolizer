package consolizer

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"strings"
)
// TODO: Make memory for events? Or at least a structure for it?

type eventStateType struct {
	stateId int
	focusedLayerAlias string
	focusedControlAlias string
	focusedControlType int
}

var eventStateMemory eventStateType

/*
updateEventQueues allows you to update all event queues so that information
such as mouse clicks, keystrokes, and other events are properly registered.
*/
func updateEventQueues() {
	event := commonResource.screen.PollEvent()
	switch event := event.(type) {
	case *tcell.EventResize:
		commonResource.screen.Sync()
	case *tcell.EventKey:
		isScreenUpdateRequired := false
		keystroke := ""
		if strings.Contains(event.Name(), "Rune") {
			keystroke = fmt.Sprintf("%c", event.Rune())
		} else {
			keystroke = strings.ToLower(event.Name())
		}
		if updateKeyboardEventTextField(keystroke) {
			isScreenUpdateRequired = true
		}
		if isScreenUpdateRequired == true {
			UpdateDisplay()
		}
		memory.KeyboardMemory.AddKeystrokeToKeyboardBuffer(keystroke)
	case *tcell.EventMouse:
		mouseXLocation, mouseYLocation := event.Position()
		var mouseButtonNumber uint
		mouseButton := event.Buttons()
		for index := uint(0); index < 8; index++ {
			if int(mouseButton)&(1<<index) != 0 {
				mouseButtonNumber = index + 1
			}
		}
		wheelState := ""
		if mouseButton&tcell.WheelUp != 0 {
			wheelState = "Up"
		} else if mouseButton&tcell.WheelDown != 0 {
			wheelState = "Down"
		} else if mouseButton&tcell.WheelLeft != 0 {
			wheelState = "Left"
		} else if mouseButton&tcell.WheelRight != 0 {
			wheelState = "Right"
		}
		memory.SetMouseStatus(mouseXLocation, mouseYLocation, mouseButtonNumber, wheelState)
		setFocusedControl()
		bringLayerToFrontIfRequired()
		moveLayerIfRequired()
		updateMouseEventTextField()
		updateButtonStates(true)
	}
}

func setFocusedControl() {
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	if buttonPressed != 0 {
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if characterEntry.AttributeEntry.CellType == constants.CellTypeTextField {
			eventStateMemory.focusedControlType = constants.CellTypeTextField
			eventStateMemory.focusedControlAlias = characterEntry.AttributeEntry.CellAlias
			eventStateMemory.focusedLayerAlias = characterEntry.LayerAlias
		}
	}
}

func updateKeyboardEventTextField(keystroke string) bool {
	isScreenUpdateRequired := true
	if eventStateMemory.focusedControlType != constants.CellTypeTextField {
		return false
	}
	textFieldEntry := memory.GetTextField(eventStateMemory.focusedLayerAlias, eventStateMemory.focusedControlAlias)
	viewportPosition := textFieldEntry.ViewportPosition
	cursorPosition := textFieldEntry.CursorPosition
	currentValue := textFieldEntry.CurrentValue
	viewportWidth := textFieldEntry.Width
	maxLengthAllowed := textFieldEntry.MaxLengthAllowed
	if len(keystroke) == 1 { // If a regular char is entered.
		if len(currentValue) < maxLengthAllowed {
			currentValue = currentValue[:viewportPosition+cursorPosition] + keystroke + currentValue[viewportPosition+cursorPosition:]
			if cursorPosition < viewportWidth {
				cursorPosition++
			} else {
				viewportPosition++
			}
			isScreenUpdateRequired = true
		}
	}
	if keystroke == "delete" {
		if currentValue != "" {
			// Protect if nothing else to delete left of string
			if viewportPosition+cursorPosition+1 <= len(currentValue) {
				currentValue = currentValue[:viewportPosition+cursorPosition] + currentValue[viewportPosition+cursorPosition+1:]
				if viewportPosition+cursorPosition == len(currentValue) {
					cursorPosition--
					if cursorPosition < 0 {
						cursorPosition = 0
					}
				}
				isScreenUpdateRequired = true
			}
		}
	}
	if keystroke == "home" {
		cursorPosition = 0
		viewportPosition = 0
		isScreenUpdateRequired = true
	}
	if keystroke == "end" {
		// If your current viewport shows the end of the input string, just move the cursor to the end of the string.
		if viewportPosition > len(currentValue)- viewportWidth {
			cursorPosition = len(currentValue) - viewportPosition
		} else {
			// Otherwise advance viewport to end of input string.
			viewportPosition = len(currentValue) - viewportWidth
			if viewportPosition < 0 {
				// If input string is smaller than even one viewport block, just set cursor to end.
				viewportPosition = 0
				cursorPosition = len(currentValue)
			} else {
				// Otherwise place cursor at end of viewport / string
				cursorPosition = viewportWidth
			}
		}
		isScreenUpdateRequired = true
	}
	if keystroke == "backspace" || keystroke == "backspace2" {
		if currentValue == "" {
			return false
		}
		// Protect if nothing else to delete left of string
		if viewportPosition + cursorPosition - 1 >= 0 {
			currentValue = currentValue[:viewportPosition+cursorPosition-1] + currentValue[viewportPosition+cursorPosition:]
			cursorPosition--
			if cursorPosition < 1 {
				if len(currentValue) < viewportWidth {
					cursorPosition = viewportPosition + cursorPosition
					viewportPosition = 0
				} else {
					if viewportPosition != 0 { // If your not at the start of an input string
						if cursorPosition == 0 {
							viewportPosition = viewportPosition - viewportWidth + 1
						} else {
							viewportPosition = viewportPosition - viewportWidth
						}
						if viewportPosition < 0 {
							viewportPosition = 0
						}
						cursorPosition = viewportWidth - 1
					}
				}
			}
			isScreenUpdateRequired = true
		}
	}
	if keystroke == "left" {
		cursorPosition--
		if cursorPosition < 0 {
			if viewportPosition == 0 {
				cursorPosition = 0
			} else {
				viewportPosition =- viewportWidth
				if viewportPosition < 0 {
					viewportPosition = 0
				}
				cursorPosition = viewportWidth
			}
		}
		isScreenUpdateRequired = true
	}
	if keystroke == "right" {
		cursorPosition++
		if viewportPosition + cursorPosition > len(currentValue){
			cursorPosition--
		} else {
			if cursorPosition >= viewportWidth {
				viewportPosition++
				cursorPosition = viewportWidth - 1
			}
		}
		isScreenUpdateRequired = true
	}
	if isScreenUpdateRequired {
		textFieldEntry.ViewportPosition = viewportPosition
		textFieldEntry.CursorPosition = cursorPosition
		textFieldEntry.CurrentValue = currentValue
		textFieldEntry.Width = viewportWidth
		textFieldEntry.MaxLengthAllowed = maxLengthAllowed
	}
	return isScreenUpdateRequired
}

func updateMouseEventTextField() {
	var characterEntry memory.CharacterEntryType
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	//previousMouseXLocation, previousMouseYLocation, previousButtonPressed, _ := memory.GetPreviousMouseStatus()
	if buttonPressed != 0 {
		characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if characterEntry.AttributeEntry.CellType == constants.CellTypeTextField {
			textFieldEntry := memory.GetTextField(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellAlias)
			textFieldEntry.CursorPosition = characterEntry.AttributeEntry.CellTypeId
			UpdateDisplay()
		}
	}
}

/*
moveLayerIfRequired allows you to move any interactive layer that has been
captured in a drag and drop action. If the mouse button is pressed over an
interactive part of a layer and not released, this method will move the
layer according to the mice's new position. In addition, the following
information should be noted:

- If the layer being moved causes the top row of characters (the interactive
title bar of a layer) to fall outside the parent layers visible area, then
no movement is performed. This is done so that it is impossible to move
a window off-screen where it can never be grabbed again.
*/
func moveLayerIfRequired() {
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	previousMouseXLocation, previousMouseYLocation, previousButtonPressed, _ := memory.GetPreviousMouseStatus()
	if buttonPressed != 0 {
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if previousButtonPressed != 0 && eventStateMemory.stateId == constants.EventStateDragAndDrop {
			xMove := mouseXLocation - previousMouseXLocation
			yMove := mouseYLocation - previousMouseYLocation
			MoveLayerByRelativeValue(eventStateMemory.focusedLayerAlias, xMove, yMove)
			if isInteractiveLayerOffscreen(eventStateMemory.focusedLayerAlias) {
				MoveLayerByRelativeValue(eventStateMemory.focusedLayerAlias, -xMove, -yMove)
			}
			UpdateDisplay()
		}
		if characterEntry.AttributeEntry.CellType == constants.CellTypeFrameTop {
			eventStateMemory.stateId = constants.EventStateDragAndDrop
			eventStateMemory.focusedLayerAlias = characterEntry.LayerAlias
		}
	} else {
		eventStateMemory.stateId = 0
	}
}

func bringLayerToFrontIfRequired() {
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	if buttonPressed != 0 {
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if characterEntry.LayerAlias == "" {
			return
		}
		layerEntry := memory.GetLayer(characterEntry.LayerAlias)
		if layerEntry.IsFocusable == true {
			return
		}
		layerAlias, previousLayerAlias := memory.GetRootParentLayerAlias(characterEntry.LayerAlias, "")
		memory.SetHighestZOrderNumber(previousLayerAlias, layerAlias)
	}
}

func getHighestZOrderNumberFromBaseLayers() {

}


/*
isInteractiveLayerOffscreen allows you to detect if a layer has been moved
off-screen or not. This is useful for when you want to constrain a window
from moving off-screen because it would be impossible for the user to drag
it back to the visible viewing area. In addition, the following
information should be noted:

- This method only considers a layer off-screen if the top row of characters
are not visible (the interactive title bar of a layer).

- Layers that are moved to the far left are considered off-screen when only
two character spaces remain. This constraint is triggered two spaces early
to account for window drop shadows that are not part of the interactive
area.

- If a layer has a parent alias, then the constraining area is set to the
parent layer dimensions instead of the terminal window dimensions.
*/
func isInteractiveLayerOffscreen(layerAlias string) bool {
	layerEntry := memory.GetLayer(layerAlias)
	viewportWidth := commonResource.terminalWidth
	viewportHeight := commonResource.terminalHeight
	if layerEntry.ParentAlias != "" {
		parentEntry := memory.GetLayer(layerEntry.ParentAlias)
		viewportWidth = parentEntry.Width
		viewportHeight = parentEntry.Height
	}
	if !(layerEntry.ScreenXLocation < viewportWidth && layerEntry.ScreenXLocation + layerEntry.Width - 2 > 0) ||
		!(layerEntry.ScreenYLocation >= 0 && layerEntry.ScreenYLocation < viewportHeight) {
		return true
	}
	return false
}

/*
getButtonClickIdentifier allows you to obtain the layer alias and the button
alias for the text cell currently under the mouse cursor. This is useful
for determining which button the user has clicked (if any).
*/
func getCellInformationUnderMouseCursor(mouseXLocation int, mouseYLocation int) memory.CharacterEntryType{
	var characterEntry memory.CharacterEntryType
	layerEntry := commonResource.screenLayer
	mouseYLocationOnLayer := mouseYLocation - layerEntry.ScreenYLocation
	mouseXLocationOnLayer := mouseXLocation - layerEntry.ScreenXLocation
	if mouseYLocationOnLayer >= 0 && mouseXLocationOnLayer >= 0 &&
		mouseYLocationOnLayer < len(layerEntry.CharacterMemory) && mouseXLocationOnLayer < len(layerEntry.CharacterMemory[0]) {
		characterEntry = layerEntry.CharacterMemory[mouseYLocation-layerEntry.ScreenYLocation][mouseXLocation-layerEntry.ScreenXLocation]
	}
	return characterEntry
}