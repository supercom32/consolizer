package consolizer

import (
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"strings"
)

type controlIdentifierType struct {
	layerAlias string
	controlAlias string
	controlType int
}

type eventStateType struct {
	stateId int
	currentlyFocusedControl controlIdentifierType
	// This variable is used to keep track of items which were highlighted so that they can be
	// un-highlighted later. Currently, only used by selectors.
	previouslyHighlightedControl controlIdentifierType
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
		var keystroke []rune
		if strings.Contains(event.Name(), "Rune") {
			keystroke = []rune{event.Rune()}
		} else {
			keystroke = []rune(strings.ToLower(event.Name()))
		}
		if scrollbar.updateKeyboardEventScrollbar(keystroke) {
			isScreenUpdateRequired = true
		}
		if textField.updateKeyboardEventTextField(keystroke) {
			isScreenUpdateRequired = true
		}
		if textbox.updateKeyboardEventTextbox(keystroke) {
			isScreenUpdateRequired = true
		}
		if Selector.updateKeyboardEventSelector(keystroke) {
			isScreenUpdateRequired = true
		}
		if isScreenUpdateRequired == true {
			UpdateDisplay()
		} else {
			// Only add keystrokes to the buffer if not already consumed by other controls.
			//fmt.Print(keystroke)
			memory.KeyboardMemory.AddKeystrokeToKeyboardBuffer(keystroke)
		}

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
		isScreenUpdateRequired := false
		memory.SetMouseStatus(mouseXLocation, mouseYLocation, mouseButtonNumber, wheelState)
		bringLayerToFrontIfRequired()
		if moveLayerIfRequired() {
			isScreenUpdateRequired = true
		}
		if textField.updateMouseEventTextField() {
			isScreenUpdateRequired = true
		}
		if checkbox.updateMouseEventCheckbox() {
			isScreenUpdateRequired = true
		}
		if button.updateButtonStates(true) {
			isScreenUpdateRequired = true
		}
		if scrollbar.updateMouseEventScrollbar() {
			isScreenUpdateRequired =  true
		}
		if Selector.updateMouseEventSelector() {
			isScreenUpdateRequired =  true
		}
		if textbox.updateMouseEventTextbox() {
			isScreenUpdateRequired =  true
		}
		if radioButton.updateMouseEventRadioButton() {
			isScreenUpdateRequired = true
		}
		// This is done last so that it can update itself if a Selector or scroll bar change was detected.
		if dropdown.updateDropdownStateMouse() {
			isScreenUpdateRequired =  true
		}
		if isScreenUpdateRequired {
			UpdateDisplay()
		}
	}
}

func setFocusedControl(layerAlias string, controlAlias string, controlType int) {
	eventStateMemory.currentlyFocusedControl.layerAlias = layerAlias
	eventStateMemory.currentlyFocusedControl.controlAlias = controlAlias
	eventStateMemory.currentlyFocusedControl.controlType = controlType
}

func isControlCurrentlyFocused(layerAlias string, controlAlias string, cellType int) bool {
	if eventStateMemory.currentlyFocusedControl.layerAlias == layerAlias &&
		eventStateMemory.currentlyFocusedControl.controlAlias == controlAlias &&
		eventStateMemory.currentlyFocusedControl.controlType == cellType {
		return true
	}
	return false
}

func setPreviouslyHighlightedControl(layerAlias string, controlAlias string, controlType int) {
	eventStateMemory.previouslyHighlightedControl.layerAlias = layerAlias
	eventStateMemory.previouslyHighlightedControl.controlAlias = controlAlias
	eventStateMemory.previouslyHighlightedControl.controlType = controlType
}

/*
moveLayerIfRequired allows you to move any interactive layer that has been
captured in a drag and drop action. If the mouse buttonType is pressed over an
interactive part of a layer and not released, this method will move the
layer according to the mice's new position. In addition, the following
information should be noted:

- If the layer being moved causes the top row of characters (the interactive
title bar of a layer) to fall outside the parent layers visible area, then
no movement is performed. This is done so that it is impossible to move
a window off-screen where it can never be grabbed again.
*/
func moveLayerIfRequired() bool {
	isScreenUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	previousMouseXLocation, previousMouseYLocation, previousButtonPressed, _ := memory.GetPreviousMouseStatus()
	if buttonPressed != 0 {
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if previousButtonPressed != 0 && eventStateMemory.stateId == constants.EventStateDragAndDrop {
			xMove := mouseXLocation - previousMouseXLocation
			yMove := mouseYLocation - previousMouseYLocation
			MoveLayerByRelativeValue(eventStateMemory.currentlyFocusedControl.layerAlias, xMove, yMove)
			if isInteractiveLayerOffscreen(eventStateMemory.currentlyFocusedControl.layerAlias) {
				MoveLayerByRelativeValue(eventStateMemory.currentlyFocusedControl.layerAlias, -xMove, -yMove)
			}
			isScreenUpdateRequired = true
		}
		if characterEntry.AttributeEntry.CellType == constants.CellTypeFrameTop {
			eventStateMemory.stateId = constants.EventStateDragAndDrop
			eventStateMemory.currentlyFocusedControl.layerAlias = characterEntry.LayerAlias
		}
	} else {
		eventStateMemory.stateId = 0
	}
	return isScreenUpdateRequired
}

/*
bringLayerToFrontIfRequired allows you to bring a layer to the front of the
visible display area if the layer being clicked is focusable.
*/
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
getButtonClickIdentifier allows you to obtain the layer alias and the buttonType
alias for the text cell currently under the mouse cursor. This is useful
for determining which buttonType the user has clicked (if any).
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