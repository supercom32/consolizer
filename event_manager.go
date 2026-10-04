package consolizer

import (
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
	"strings"
	"time"
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

// tabIndexBeforeFirstEntry marks a tab position that precedes the first registered entry, so that the next Tab press
// lands on entry 0.
const tabIndexBeforeFirstEntry = -1

var eventStateMemory eventStateType
var eventIntervalTime time.Time
var lastMouseMoveTime time.Time

// resizeDebounceDelay is how long EventResize handling waits for resizing to pause before doing its actual work
// (terminal size update, Sync, and UpdateDisplay). A live drag-resize fires many EventResize events in rapid
// succession, and handling each one immediately means a full layer recomposition on every single one of them.
// resizeDebounceTimer is only ever touched from UpdateEventQueues's single dedicated event-polling goroutine, so
// it needs no synchronization of its own.
const resizeDebounceDelay = 150 * time.Millisecond

var resizeDebounceTimer *time.Timer

/*
UpdatePeriodicEvents is a method which triggers periodic events such as tooltip updates and keyboard state clearing.

Example:
    UpdatePeriodicEvents()
*/
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

	// Clear key states periodically to handle key releases
	// This is done more frequently than other periodic events
	// to ensure responsive keyboard input
	KeyboardMemory.ClearKeyStates()
}

/*
UpdateEventQueues is a method which updates all event queues so that information such as mouse clicks, keystrokes, and
other events are properly registered.

Example:
    UpdateEventQueues()
*/
func UpdateEventQueues() {
	// Skip event processing if screen is not initialized (e.g., in debug mode or tests)
	if commonResource.screen == nil {
		return
	}

	event := commonResource.screen.PollEvent()
	switch event := event.(type) {
	case *tcell.EventResize:
		newWidth, newHeight := event.Size()
		// Each new resize event cancels any pending one and reschedules, so the actual work below only runs
		// once resizing has paused for resizeDebounceDelay, not on every intermediate tick of an active drag.
		if resizeDebounceTimer != nil {
			resizeDebounceTimer.Stop()
		}
		resizeDebounceTimer = time.AfterFunc(resizeDebounceDelay, func() {
			// commonResource.displayUpdate is the same lock UpdateDisplay holds for its entire duration. Taking
			// it here serializes this Sync() (and the terminalWidth/terminalHeight update below) against a
			// concurrent UpdateDisplay call on another goroutine, which would otherwise let two goroutines drive
			// the screen concurrently and corrupt terminal output.
			commonResource.displayUpdate.Lock()
			// A timer scheduled before shutdown can still fire after RestoreTerminalSettings has already torn the
			// screen down, since stopping this timer is not part of that shutdown sequence. Bail out before Sync
			// touches the freed screen; the UpdateDisplay call below is skipped too since it would no-op anyway.
			if commonResource.isTerminated {
				commonResource.displayUpdate.Unlock()
				return
			}
			// Only auto-size sessions track the physical terminal's size. A session initialized with an explicit
			// fixed width and height must keep rendering, mouse hit testing, and layer bounds pinned to that size
			// even after the physical terminal is resized around it.
			if commonResource.isAutoSizeEnabled {
				commonResource.terminalWidth = newWidth
				commonResource.terminalHeight = newHeight
			}
			commonResource.screen.Sync()
			commonResource.displayUpdate.Unlock()
			// Called outside the lock above since the write lock is not reentrant and UpdateDisplay acquires the
			// same lock itself. Without this, the display would stay stale until whatever unrelated event
			// happens to trigger the next UpdateDisplay call, such as a keypress.
			UpdateDisplay(false)
		})
	case *tcell.EventKey:
		isScreenUpdateRequired := false
		isKeystrokeConsumed := false
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
			isKeystrokeConsumed = true
		}
		if updateRequired, consumed := scrollbar.updateKeyboardEvent(keystroke); updateRequired {
			isScreenUpdateRequired = true
			isKeystrokeConsumed = consumed
		}
		if updateRequired, consumed := TextField.updateKeyboardEvent(keystroke); updateRequired {
			isScreenUpdateRequired = true
			isKeystrokeConsumed = consumed
		}
		if updateRequired, consumed := textbox.UpdateKeyboardEvent(keystroke); updateRequired {
			isScreenUpdateRequired = true
			isKeystrokeConsumed = consumed
		}
		if updateRequired, consumed := Selector.updateKeyboardEvent(keystroke); updateRequired {
			isScreenUpdateRequired = true
			isKeystrokeConsumed = consumed
		}
		if updateRequired, consumed := Dropdown.updateKeyboardEvent(keystroke); updateRequired {
			isScreenUpdateRequired = true
			isKeystrokeConsumed = consumed
		}
		if updateRequired, consumed := FileMenu.updateKeyboardEvent(keystroke); updateRequired {
			isScreenUpdateRequired = true
			isKeystrokeConsumed = consumed
		}
		if isScreenUpdateRequired == true {
			UpdateDisplay(false)
		}
		// Only add keystroke to buffer if it wasn't consumed by a control
		if !isKeystrokeConsumed && keystroke != nil {
			KeyboardMemory.AddToBuffer(keystroke)
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

		// Throttle mouse movement events, but never a release: a release is only ever detected here by comparing
		// against the last recorded button (SetMouseStatus has not run yet for this event), and dropping one
		// leaves whatever it should have released (a button, or a scrollbar drag) stuck in its pressed state,
		// since nothing else clears it.
		_, _, lastRecordedButtonNumber, _ := GetMouseStatus()
		isRelease := lastRecordedButtonNumber != 0 && mouseButtonNumber == 0
		isPureMovement := mouseButtonNumber == 0 && lastRecordedButtonNumber == 0 && wheelState == ""
		isScrollbarDragMovement := mouseButtonNumber != 0 && eventStateMemory.stateId == constants.EventStateDragAndDropScrollbar
		if !isRelease && (isPureMovement || isScrollbarDragMovement) {
			elapsedTime := time.Since(lastMouseMoveTime)
			if elapsedTime < 50*time.Millisecond {
				return
			}
			lastMouseMoveTime = time.Now()
		}

		SetMouseStatus(mouseXLocation, mouseYLocation, mouseButtonNumber, wheelState)
		bringLayerToFrontIfRequired()
		if moveLayerIfRequired() {
			isScreenUpdateRequired = true
			// Don't accept any new mouse manipulations if you're in drag-and-drop mode.
			return
		}

		// If a mouse button is pressed, find what control is under the mouse
		// so we can set focus to it.
		if mouseButton&tcell.Button1 != 0 && eventStateMemory.stateId == constants.EventStateNone {
			characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
			setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, characterEntry.AttributeEntry.CellType)
		}
		if Tooltip.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		if TextField.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		if FileMenu.updateStateMouse() {
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
		if Dropdown.updateStateMouse() {
			isScreenUpdateRequired = true
		}
		if Checkbox.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		if Button.updateStates(true) {
			isScreenUpdateRequired = true
		}
		if scrollbar.updateMouseEvent() {
			buttonHistory.clear()
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
		if viewport.updateMouseEvent() {
			isScreenUpdateRequired = true
		}
		// This is done last so that it can update itself if a Selector or scroll bar change was detected.
		if Dropdown.updateStateMouse() {
			isScreenUpdateRequired = true
		}
		// LogInfo("mouse event dropdownb")
		if isScreenUpdateRequired {
			UpdateDisplay(false)
		}
	}
}

/*
ClearTabIndex is a method which allows you to clear all registered tab index entries from memory and reset the tab
position to before the first entry, so the next Tab press after the tab order is rebuilt starts from the beginning of
the new order unless the currently focused control is part of it.

Example:

	ClearTabIndex()
*/
func ClearTabIndex() {
	eventStateMemory.tabIndexMemory = nil
	eventStateMemory.currentTabIndex = tabIndexBeforeFirstEntry
}

/*
addTabIndex is a method which registers a new control in the tab index memory for sequential navigation.

Example:
    addTabIndex("layer1", "button1", constants.CellTypeButton)
*/
func addTabIndex(layerAlias string, controlAlias string, controlType int) {
	controlEntry := controlIdentifierType{layerAlias: layerAlias, controlAlias: controlAlias, controlType: controlType}
	eventStateMemory.tabIndexMemory = append(eventStateMemory.tabIndexMemory, controlEntry)
}

/*
nextTabIndex is a method which allows you to advance the focus to the control registered after the currently focused
control in the tab index sequence, wrapping back to the first entry after the last one. If the focused control is not
part of the tab order, or nothing is focused, focus moves to the first entry. When no tab order is registered, focus is
left unchanged. In addition, the following should be noted:

  - The starting point is always the control that actually has focus rather than the stored tab position, since
    focus can also change through mouse clicks, GetFocus, or Selector.Add without the stored position being
    updated. The stored position is kept in step with the result.

Example:

	nextTabIndex()
*/
func nextTabIndex() {
	tabOrderLength := len(eventStateMemory.tabIndexMemory)
	if tabOrderLength == 0 {
		return
	}
	nextIndex := 0
	focusedIndex := findTabIndexOfControl(eventStateMemory.currentlyFocusedControl)
	if focusedIndex != tabIndexBeforeFirstEntry {
		nextIndex = (focusedIndex + 1) % tabOrderLength
	}
	eventStateMemory.currentTabIndex = nextIndex
	eventStateMemory.currentlyFocusedControl = eventStateMemory.tabIndexMemory[nextIndex]
}

/*
findTabIndexOfControl is a method which allows you to locate a control within the registered tab index sequence by
matching its layer alias, control alias, and control type. It returns the position of the first matching entry, or
tabIndexBeforeFirstEntry if the control is not registered.

Example:

	index := findTabIndexOfControl(eventStateMemory.currentlyFocusedControl)
*/
func findTabIndexOfControl(control controlIdentifierType) int {
	for index, entry := range eventStateMemory.tabIndexMemory {
		if entry == control {
			return index
		}
	}
	return tabIndexBeforeFirstEntry
}

/*
setFocusedControl is a method which explicitly sets which control currently has focus.

Example:
    setFocusedControl("layer1", "textfield1", constants.CellTypeTextField)
*/
func setFocusedControl(layerAlias string, controlAlias string, controlType int) {
	eventStateMemory.currentlyFocusedControl.layerAlias = layerAlias
	eventStateMemory.currentlyFocusedControl.controlAlias = controlAlias
	eventStateMemory.currentlyFocusedControl.controlType = controlType
}

/*
isControlCurrentlyFocused is a method which checks if a specific control is currently the focused control in the
application.

Example:
    isControlCurrentlyFocused("layer1", "textfield1", constants.CellTypeTextField)
*/
func isControlCurrentlyFocused(layerAlias string, controlAlias string, cellType int) bool {
	if eventStateMemory.currentlyFocusedControl.layerAlias == layerAlias &&
		eventStateMemory.currentlyFocusedControl.controlAlias == controlAlias &&
		eventStateMemory.currentlyFocusedControl.controlType == cellType {
		return true
	}
	return false
}

/*
setPreviouslyHighlightedControl is a method which records the control that was previously highlighted by the mouse.

Example:
    setPreviouslyHighlightedControl("layer1", "item1", constants.CellTypeSelector)
*/
func setPreviouslyHighlightedControl(layerAlias string, controlAlias string, controlType int) {
	eventStateMemory.previouslyHighlightedControl.layerAlias = layerAlias
	eventStateMemory.previouslyHighlightedControl.controlAlias = controlAlias
	eventStateMemory.previouslyHighlightedControl.controlType = controlType
}

/*
moveLayerIfRequired is a method which moves any interactive layer that has been captured in a drag and drop action. If
the mouse button is pressed over an interactive part of a layer and not released, this method will move the layer
according to the mouse's new position. In addition, the following should be noted:

- If the layer being moved causes the top row of characters (the interactive title bar of a layer) to fall outside the
  visible screen area, the move is cancelled.

- This is done so that it is impossible to move a window off-screen where it can never be grabbed again.

Example:
    moveLayerIfRequired()
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
			moveLayerByRelativeValue(eventStateMemory.currentlyFocusedControl.layerAlias, xMove, yMove)
			if isInteractiveLayerOffscreen(eventStateMemory.currentlyFocusedControl.layerAlias) {
				moveLayerByRelativeValue(eventStateMemory.currentlyFocusedControl.layerAlias, -xMove, -yMove)
			}
			isScreenUpdateRequired = true
		} else if characterEntry.AttributeEntry.CellType == constants.CellTypeFrameTop && eventStateMemory.stateId != constants.EventStateDragAndDrop {
			// Only set the drag state and focused control if we're not already dragging
			eventStateMemory.stateId = constants.EventStateDragAndDrop
			setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, characterEntry.AttributeEntry.CellType)
		}
	} else {
		eventStateMemory.stateId = constants.EventStateNone
	}
	return isScreenUpdateRequired
}

/*
bringLayerToFrontIfRequired is a method which brings a layer to the front of the visible display area if the layer being
clicked is focusable. In addition, the following should be noted:

  - Any unread button press is cleared here, but only on the mouse-down transition itself (the previous mouse state
    had no button held and the current one does). This runs on every mouse event with a button held, including
    held movement, and clearing unconditionally on those would wipe a press before the application ever reads it.

Example:
    bringLayerToFrontIfRequired()
*/
func bringLayerToFrontIfRequired() {
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	if buttonPressed != 0 {
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if characterEntry.LayerAlias == "" {
			return
		}
		if characterEntry.AttributeEntry.CellType == constants.CellTypeShadow {
			return
		}
		_, _, previousButtonPressed, _ := GetPreviousMouseStatus()
		if previousButtonPressed == 0 {
			buttonHistory.clear()
		}
		// Protect against layer deletions.
		layerEntry, isFound := Layers.Lookup(characterEntry.LayerAlias)
		if !isFound {
			return
		}
		if layerEntry.IsFocusable == true {
			return
		}
		layerAlias, previousLayerAlias := layer.GetRootParentAlias(characterEntry.LayerAlias, "")
		layer.SetHighestZOrderNumber(previousLayerAlias, layerAlias)
	}
}

/*
isInteractiveLayerOffscreen is a method which detects if a layer has been moved off-screen or not. This is useful for
when you want to constrain a window from moving off-screen because it would be impossible for the user to drag it back
to the visible viewing area. In addition, the following should be noted:

- This method only considers a layer off-screen if the top row of characters are not visible (the
  interactive title bar).

- Layers that are moved to the far left are considered off-screen when only two character spaces remain.

- This constraint is triggered two spaces early to account for window drop shadows that are not part of the
  interactive area.

- If a layer has a parent alias, then the constraining area is set to the parent layer dimensions instead of the
  terminal dimensions.

Example:
    isInteractiveLayerOffscreen("layer1")
*/
func isInteractiveLayerOffscreen(layerAlias string) bool {
	layerEntry, isFound := Layers.Lookup(layerAlias)
	if !isFound {
		return false
	}
	viewportWidth := commonResource.terminalWidth
	viewportHeight := commonResource.terminalHeight
	if layerEntry.ParentAlias != "" {
		if parentEntry, isFound := Layers.Lookup(layerEntry.ParentAlias); isFound {
			viewportWidth = parentEntry.Width
			viewportHeight = parentEntry.Height
		}
	}
	if !(layerEntry.ScreenXLocation < viewportWidth && layerEntry.ScreenXLocation+layerEntry.Width-2 > 0) ||
		!(layerEntry.ScreenYLocation >= 0 && layerEntry.ScreenYLocation < viewportHeight) {
		return true
	}
	return false
}

/*
getCellInformationUnderMouseCursor is a method which obtains the layer alias and the buttonType alias for the text cell
currently under the mouse cursor. This is useful for determining which buttonType the user has clicked (if any). In
addition, the following should be noted:

  - The shared screen snapshot is copied out under a read lock before use, so a concurrent UpdateDisplay call on
    another goroutine, such as the one driving the periodic mouse hit-test that calls this method, can never be
    observed mid-write.

Example:

	getCellInformationUnderMouseCursor(10, 20)
*/
func getCellInformationUnderMouseCursor(mouseXLocation int, mouseYLocation int) types.CharacterEntryType {
	var characterEntry types.CharacterEntryType
	commonResource.displayUpdate.RLock()
	layerEntry := commonResource.screenLayer
	commonResource.displayUpdate.RUnlock()
	mouseYLocationOnLayer := mouseYLocation - layerEntry.ScreenYLocation
	mouseXLocationOnLayer := mouseXLocation - layerEntry.ScreenXLocation
	if mouseYLocationOnLayer >= 0 && mouseXLocationOnLayer >= 0 &&
		mouseYLocationOnLayer < len(layerEntry.CharacterMemory) && mouseXLocationOnLayer < len(layerEntry.CharacterMemory[0]) {
		characterEntry = layerEntry.CharacterMemory[mouseYLocation-layerEntry.ScreenYLocation][mouseXLocation-layerEntry.ScreenXLocation]
	}
	return characterEntry
}

/*
IsModifierKeyPressed is a method which checks if a specific modifier key is currently pressed.

Example:
    IsModifierKeyPressed(tcell.ModShift)
*/
func IsModifierKeyPressed(modifier tcell.ModMask) bool {
	return (eventStateMemory.modifierKeys & modifier) != 0
}

/*
IsShiftPressed is a method which checks if the shift key is currently pressed.

Example:
    IsShiftPressed()
*/
func IsShiftPressed() bool {
	return IsModifierKeyPressed(tcell.ModShift)
}

/*
IsCtrlPressed is a method which checks if the control key is currently pressed.

Example:
    IsCtrlPressed()
*/
func IsCtrlPressed() bool {
	return IsModifierKeyPressed(tcell.ModCtrl)
}

/*
IsAltPressed is a method which checks if the alt key is currently pressed.

Example:
    IsAltPressed()
*/
func IsAltPressed() bool {
	return IsModifierKeyPressed(tcell.ModAlt)
}
