package consolizer

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

type controlIdentifierType struct {
	layerAlias   string
	controlAlias string
	controlType  int
}

type eventStateType struct {
	stateId                 int
	currentlyFocusedControl controlIdentifierType
	// This variable is used to keep track of items which were highlighted so that they can be
	// un-highlighted later. Currently, only used by selectors and tooltips
	previouslyHighlightedControl controlIdentifierType
	tabIndexMemory               []controlIdentifierType
	currentTabIndex              int
	// Track modifier key states
	modifierKeys tcell.ModMask
}

var eventStateMemory eventStateType
var eventIntervalTime time.Time

func UpdatePeriodicEvents() {
	elapsedTime := time.Since(eventIntervalTime)
	if elapsedTime >= 500*time.Millisecond {
		eventIntervalTime = time.Now()
		isScreenUpdateRequired := false
		if Tooltip.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		if isScreenUpdateRequired == true {
			UpdateDisplay(false)
		}
	}
}

/*
UpdateEventQueues allows you to update all event queues so that information
such as mouse clicks, keystrokes, and other events are properly registered.
*/
func UpdateEventQueues() {
	event := commonResource.screen.PollEvent()
	switch event := event.(type) {
	case *tcell.EventResize:
		commonResource.screen.Sync()
	case *tcell.EventKey:
		isScreenUpdateRequired := false
		var keystroke []rune

		// Update modifier key state
		eventStateMemory.modifierKeys = event.Modifiers()

		if strings.Contains(event.Name(), "Rune") {
			keystroke = []rune{event.Rune()}
		} else {
			keystroke = []rune(strings.ToLower(event.Name()))
		}
		if string(keystroke) == "tab" {
			nextTabIndex()
			keystroke = nil
			isScreenUpdateRequired = true
		}
		if scrollbar.updateKeyboardEvent(keystroke) {
			isScreenUpdateRequired = true
		}
		if TextField.updateKeyboardEvent(keystroke) {
			isScreenUpdateRequired = true
		}
		if textbox.UpdateKeyboardEvent(keystroke) {
			isScreenUpdateRequired = true
		}
		if Selector.updateKeyboardEvent(keystroke) {
			isScreenUpdateRequired = true
		}
		if Dropdown.updateKeyboardEvent(keystroke) {
			isScreenUpdateRequired = true
		}
		if isScreenUpdateRequired == true {
			UpdateDisplay(false)
		}
		KeyboardMemory.AddKeystrokeToKeyboardBuffer(keystroke)

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
		SetMouseStatus(mouseXLocation, mouseYLocation, mouseButtonNumber, wheelState)
		bringLayerToFrontIfRequired()
		if moveLayerIfRequired() {
			isScreenUpdateRequired = true
		}
		if Tooltip.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		if TextField.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		if Selector.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		if textbox.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		if radioButton.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		// This is done last so that it can update itself if a Selector or scroll bar change was detected.
		if Dropdown.updateDropdownStateMouse() {
			isScreenUpdateRequired = true
		}
		if Checkbox.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		if Button.updateButtonStates(true) {
			isScreenUpdateRequired = true
		}
		if scrollbar.updateMouseEvent() {
			buttonHistory.layerAlias = ""
			buttonHistory.buttonAlias = ""
			isScreenUpdateRequired = true
		}
		// LogInfo("mouse scrollbar" + time.Now().String())
		if Selector.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		// LogInfo("mouse event selector" + time.Now().String())
		if textbox.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		// LogInfo("mouse event textbox" + time.Now().String())
		if radioButton.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		// LogInfo("mouse event radio" + time.Now().String())
		// This is done last so that it can update itself if a Selector or scroll bar change was detected.
		if Dropdown.updateDropdownStateMouse() {
			isScreenUpdateRequired = true
		}
		// LogInfo("mouse event dropdownb")
		if isScreenUpdateRequired {
			UpdateDisplay(false)
		}
	}
}

func ClearTabIndex() {
	eventStateMemory.tabIndexMemory = nil
}

func addTabIndex(layerAlias string, controlAlias string, controlType int) {
	controlEntry := controlIdentifierType{layerAlias: layerAlias, controlAlias: controlAlias, controlType: controlType}
	eventStateMemory.tabIndexMemory = append(eventStateMemory.tabIndexMemory, controlEntry)
}

func nextTabIndex() {
	eventStateMemory.currentTabIndex++
	if eventStateMemory.currentTabIndex >= len(eventStateMemory.tabIndexMemory) {
		eventStateMemory.currentTabIndex = 0
	}
	if len(eventStateMemory.tabIndexMemory) != 0 {
		eventStateMemory.currentlyFocusedControl = eventStateMemory.tabIndexMemory[eventStateMemory.currentTabIndex]
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
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	previousMouseXLocation, previousMouseYLocation, previousButtonPressed, _ := GetPreviousMouseStatus()
	if buttonPressed != 0 {
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if previousButtonPressed != 0 && eventStateMemory.stateId == constants.EventStateDragAndDrop && isLayerExists(eventStateMemory.currentlyFocusedControl.layerAlias) {
			xMove := mouseXLocation - previousMouseXLocation
			yMove := mouseYLocation - previousMouseYLocation
			MoveLayerByRelativeValue(eventStateMemory.currentlyFocusedControl.layerAlias, xMove, yMove)
			if isInteractiveLayerOffscreen(eventStateMemory.currentlyFocusedControl.layerAlias) {
				MoveLayerByRelativeValue(eventStateMemory.currentlyFocusedControl.layerAlias, -xMove, -yMove)
			}
			isScreenUpdateRequired = true
		} else if characterEntry.AttributeEntry.CellType == constants.CellTypeFrameTop && eventStateMemory.stateId != constants.EventStateDragAndDrop {
			// Only set the drag state and focused control if we're not already dragging
			eventStateMemory.stateId = constants.EventStateDragAndDrop
			eventStateMemory.currentlyFocusedControl.layerAlias = characterEntry.LayerAlias
		}
	} else {
		eventStateMemory.stateId = constants.EventStateNone
	}
	return isScreenUpdateRequired
}

/*
bringLayerToFrontIfRequired allows you to bring a layer to the front of the
visible display area if the layer being clicked is focusable.
*/
func bringLayerToFrontIfRequired() {
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	if buttonPressed != 0 {
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if characterEntry.LayerAlias == "" {
			return
		}
		buttonHistory.layerAlias = ""
		buttonHistory.buttonAlias = ""
		// Protect against layer deletions.
		if !Layers.IsExists(characterEntry.LayerAlias) {
			return
		}
		layerEntry := Layers.Get(characterEntry.LayerAlias)
		if layerEntry.IsFocusable == true {
			return
		}
		layerAlias, previousLayerAlias := layer.GetRootParentLayerAlias(characterEntry.LayerAlias, "")
		layer.SetHighestZOrderNumber(previousLayerAlias, layerAlias)
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
	layerEntry := Layers.Get(layerAlias)
	viewportWidth := commonResource.terminalWidth
	viewportHeight := commonResource.terminalHeight
	if layerEntry.ParentAlias != "" {
		parentEntry := Layers.Get(layerEntry.ParentAlias)
		viewportWidth = parentEntry.Width
		viewportHeight = parentEntry.Height
	}
	if !(layerEntry.ScreenXLocation < viewportWidth && layerEntry.ScreenXLocation+layerEntry.Width-2 > 0) ||
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
func getCellInformationUnderMouseCursor(mouseXLocation int, mouseYLocation int) types.CharacterEntryType {
	var characterEntry types.CharacterEntryType
	layerEntry := commonResource.screenLayer
	mouseYLocationOnLayer := mouseYLocation - layerEntry.ScreenYLocation
	mouseXLocationOnLayer := mouseXLocation - layerEntry.ScreenXLocation
	if mouseYLocationOnLayer >= 0 && mouseXLocationOnLayer >= 0 &&
		mouseYLocationOnLayer < len(layerEntry.CharacterMemory) && mouseXLocationOnLayer < len(layerEntry.CharacterMemory[0]) {
		characterEntry = layerEntry.CharacterMemory[mouseYLocation-layerEntry.ScreenYLocation][mouseXLocation-layerEntry.ScreenXLocation]
	}
	return characterEntry
}

// IsModifierKeyPressed checks if a specific modifier key is currently pressed
func IsModifierKeyPressed(modifier tcell.ModMask) bool {
	return (eventStateMemory.modifierKeys & modifier) != 0
}

// IsShiftPressed checks if the shift key is currently pressed
func IsShiftPressed() bool {
	return IsModifierKeyPressed(tcell.ModShift)
}

// IsCtrlPressed checks if the control key is currently pressed
func IsCtrlPressed() bool {
	return IsModifierKeyPressed(tcell.ModCtrl)
}

// IsAltPressed checks if the alt key is currently pressed
func IsAltPressed() bool {
	return IsModifierKeyPressed(tcell.ModAlt)
}
