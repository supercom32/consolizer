package consolizer

import (
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
	"strings"
	"sync"
	"time"
)

/*
controlIdentifierType is a structure which identifies a single control by the alias of the layer it belongs to, its
own alias, and its cell type, which is one of the constants.CellType values.

Example:

	control := controlIdentifierType{layerAlias: "main", controlAlias: "ok", controlType: constants.CellTypeButton}
*/
type controlIdentifierType struct {
	layerAlias   string
	controlAlias string
	controlType  int
}

/*
eventStateType is a structure which holds the shared interaction state of the event manager: the drag state, the layer
being dragged, which control was last highlighted, and the current modifier keys. Focus is not kept here; the focus
manager in focus_manager.go is its only owner. In addition, the following should be noted:

  - Every field is guarded by mutex, since the event goroutine, the periodic-event goroutine, and the application's
    own goroutine all read and write this state. Fields must only be accessed through the accessor functions in this
    file, and no accessor may call another while holding the lock.

Example:

	var eventState eventStateType
*/
type eventStateType struct {
	mutex   sync.Mutex
	stateId int
	// draggedLayerAlias is the layer being moved by its title bar while stateId is EventStateDragAndDrop.
	draggedLayerAlias string
	// This variable is used to keep track of items which were highlighted so that they can be
	// un-highlighted later. Currently, only used by selectors and tooltips
	previouslyHighlightedControl controlIdentifierType
	// Track modifier key states
	modifierKeys tcell.ModMask
}

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
UpdateEventQueues is a method which allows you to process the next pending terminal event, so that information such as
mouse clicks, keystrokes, and resizes are registered and passed on to the controls, redrawing the screen when they
change. In addition, the following should be noted:

  - Every keystroke shows the keyboard focus indicator and every mouse button press hides it, as described for
    setFocusIndicatorVisible, and the screen is redrawn when that changes how the focused control looks.

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
		reconcileFocus()

		// Update modifier key state
		setModifierKeys(event.Modifiers())
		if setFocusIndicatorVisible(true) {
			isScreenUpdateRequired = true
		}

		if strings.Contains(event.Name(), "Rune") {
			keystroke = []rune{event.Rune()}
		} else {
			keystroke = []rune(strings.ToLower(event.Name()))
		}
		switch string(keystroke) {
		case "tab":
			// An open dropdown tray belongs to the control losing focus, so it is closed before focus moves on.
			Dropdown.closeAllOpen()
			FocusNext()
			keystroke = nil
			isScreenUpdateRequired = true
			isKeystrokeConsumed = true
		case "backtab", "shift+backtab", "shift+tab":
			Dropdown.closeAllOpen()
			FocusPrevious()
			keystroke = nil
			isScreenUpdateRequired = true
			isKeystrokeConsumed = true
		}
		updateRequired, consumed := runKeyboardEventHandlers(keystroke, getKeyboardEventHandlers())
		isScreenUpdateRequired = isScreenUpdateRequired || updateRequired
		isKeystrokeConsumed = isKeystrokeConsumed || consumed
		// Enter and Esc that no control used fall through to the active scope's default and cancel buttons.
		if !isKeystrokeConsumed && keystroke != nil {
			if scopeButton := getDefaultButtonForKey(string(keystroke)); scopeButton.controlAlias != "" {
				buttonHistory.set(scopeButton.layerAlias, scopeButton.controlAlias)
				isScreenUpdateRequired = true
				isKeystrokeConsumed = true
			}
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
		isScrollbarDragMovement := mouseButtonNumber != 0 && getEventStateId() == constants.EventStateDragAndDropScrollbar
		if !isRelease && (isPureMovement || isScrollbarDragMovement) {
			elapsedTime := time.Since(lastMouseMoveTime)
			if elapsedTime < 50*time.Millisecond {
				return
			}
			lastMouseMoveTime = time.Now()
		}

		SetMouseStatus(mouseXLocation, mouseYLocation, mouseButtonNumber, wheelState)
		reconcileFocus()
		// While a modal layer is active, mouse input outside it is ignored, unless a drag that started inside it is
		// still in progress. A release outside still un-presses any armed button, so it cannot get stuck pressed.
		if getEventStateId() == constants.EventStateNone {
			characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
			if !isLayerAcceptingInput(characterEntry.LayerAlias) {
				if isRelease && Button.clearAllPressedButtons() {
					UpdateDisplay(false)
				}
				return
			}
		}
		// Any mouse button press, but not movement or the wheel, hides the keyboard focus indicator.
		if mouseButtonNumber != 0 && lastRecordedButtonNumber == 0 && setFocusIndicatorVisible(false) {
			isScreenUpdateRequired = true
		}
		bringLayerToFrontIfRequired()
		if moveLayerIfRequired() {
			isScreenUpdateRequired = true
			// Don't accept any new mouse manipulations if you're in drag-and-drop mode.
			return
		}

		// A press of the left button focuses the control under it. Only the press itself does, not held movement,
		// so dragging out of a control does not take focus away from it.
		if mouseButton&tcell.Button1 != 0 && lastRecordedButtonNumber == 0 && getEventStateId() == constants.EventStateNone {
			focusControlFromClick(getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation))
			isScreenUpdateRequired = true
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
keyboardEventHandlerType is a function type which allows you to describe a control's keyboard handler. It receives the
keystroke being processed and returns whether a screen update is required and whether the keystroke was consumed.

Example:

	var handler keyboardEventHandlerType = Button.updateKeyboardEvent
*/
type keyboardEventHandlerType func(keystroke []rune) (bool, bool)

/*
getKeyboardEventHandlers is a method which allows you to obtain every control keyboard handler, in the order
UpdateEventQueues runs them.

Example:

	handlers := getKeyboardEventHandlers()
*/
func getKeyboardEventHandlers() []keyboardEventHandlerType {
	return []keyboardEventHandlerType{
		Button.updateKeyboardEvent,
		Checkbox.updateKeyboardEvent,
		radioButton.updateKeyboardEvent,
		scrollbar.updateKeyboardEvent,
		TextField.updateKeyboardEvent,
		textbox.UpdateKeyboardEvent,
		Selector.updateKeyboardEvent,
		Dropdown.updateKeyboardEvent,
		FileMenu.updateKeyboardEvent,
	}
}

/*
runKeyboardEventHandlers is a method which allows you to pass a keystroke to every handler given, in order, and combine
their results. A screen update is required if any handler asks for one, and the keystroke is consumed if any handler
consumes it. In addition, the following should be noted:

  - Every handler runs even after one has consumed the keystroke, matching the behavior each control already
    relies on, since each handler only acts when its own control has focus or is open.

  - The two results are combined independently, so a later handler that requests a redraw without consuming the
    keystroke cannot undo an earlier handler's consumption, and a handler that consumes a keystroke without needing
    a redraw is still honored.

Example:

	updateRequired, consumed := runKeyboardEventHandlers([]rune("esc"), getKeyboardEventHandlers())
*/
func runKeyboardEventHandlers(keystroke []rune, handlers []keyboardEventHandlerType) (bool, bool) {
	isScreenUpdateRequired := false
	isKeystrokeConsumed := false
	for _, handler := range handlers {
		updateRequired, consumed := handler(keystroke)
		isScreenUpdateRequired = isScreenUpdateRequired || updateRequired
		isKeystrokeConsumed = isKeystrokeConsumed || consumed
	}
	return isScreenUpdateRequired, isKeystrokeConsumed
}

/*
setPreviouslyHighlightedControl is a method which records the control that was previously highlighted by the mouse.

Example:
    setPreviouslyHighlightedControl("layer1", "item1", constants.CellTypeSelector)
*/
func setPreviouslyHighlightedControl(layerAlias string, controlAlias string, controlType int) {
	eventStateMemory.mutex.Lock()
	defer eventStateMemory.mutex.Unlock()
	eventStateMemory.previouslyHighlightedControl.layerAlias = layerAlias
	eventStateMemory.previouslyHighlightedControl.controlAlias = controlAlias
	eventStateMemory.previouslyHighlightedControl.controlType = controlType
}

/*
getPreviouslyHighlightedControl is a method which allows you to obtain the control that was last recorded as
highlighted by the mouse, read as a single consistent snapshot.

Example:

	highlightedControl := getPreviouslyHighlightedControl()
*/
func getPreviouslyHighlightedControl() controlIdentifierType {
	eventStateMemory.mutex.Lock()
	defer eventStateMemory.mutex.Unlock()
	return eventStateMemory.previouslyHighlightedControl
}

/*
getEventStateId is a method which allows you to obtain the current interaction state, one of the constants.EventState
values, such as whether a layer or scroll bar is being dragged.

Example:

	isDragging := getEventStateId() == constants.EventStateDragAndDrop
*/
func getEventStateId() int {
	eventStateMemory.mutex.Lock()
	defer eventStateMemory.mutex.Unlock()
	return eventStateMemory.stateId
}

/*
setEventStateId is a method which allows you to set the current interaction state to one of the constants.EventState
values.

Example:

	setEventStateId(constants.EventStateNone)
*/
func setEventStateId(stateId int) {
	eventStateMemory.mutex.Lock()
	defer eventStateMemory.mutex.Unlock()
	eventStateMemory.stateId = stateId
}

/*
getDraggedLayerAlias is a method which allows you to obtain the alias of the layer currently being moved by its title
bar, or an empty string when no layer is being dragged.

Example:

	layerAlias := getDraggedLayerAlias()
*/
func getDraggedLayerAlias() string {
	eventStateMemory.mutex.Lock()
	defer eventStateMemory.mutex.Unlock()
	return eventStateMemory.draggedLayerAlias
}

/*
setDraggedLayerAlias is a method which allows you to record the alias of the layer being moved by its title bar.
Passing an empty string records that no layer is being dragged.

Example:

	setDraggedLayerAlias("window1")
*/
func setDraggedLayerAlias(layerAlias string) {
	eventStateMemory.mutex.Lock()
	defer eventStateMemory.mutex.Unlock()
	eventStateMemory.draggedLayerAlias = layerAlias
}

/*
setModifierKeys is a method which allows you to record the modifier keys held during the most recent key event, so that
IsModifierKeyPressed and its helpers can report them.

Example:

	setModifierKeys(tcell.ModShift)
*/
func setModifierKeys(modifierKeys tcell.ModMask) {
	eventStateMemory.mutex.Lock()
	defer eventStateMemory.mutex.Unlock()
	eventStateMemory.modifierKeys = modifierKeys
}

/*
moveLayerIfRequired is a method which moves any interactive layer that has been captured in a drag and drop action. If
the mouse button is pressed over an interactive part of a layer and not released, this method will move the layer
according to the mouse's new position. In addition, the following should be noted:

- If the layer being moved causes the top row of characters (the interactive title bar of a layer) to fall outside the
  visible screen area, the move is cancelled.

- This is done so that it is impossible to move a window off-screen where it can never be grabbed again.

- The layer being dragged is tracked separately from focus, so dragging a window by its title bar never takes focus
  away from the control that has it.

Example:
    moveLayerIfRequired()
*/
func moveLayerIfRequired() bool {
	isScreenUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	previousMouseXLocation, previousMouseYLocation, previousButtonPressed, _ := GetPreviousMouseStatus()
	if buttonPressed != 0 {
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		draggedLayerAlias := getDraggedLayerAlias()
		if previousButtonPressed != 0 && getEventStateId() == constants.EventStateDragAndDrop && isLayerExists(draggedLayerAlias) {
			xMove := mouseXLocation - previousMouseXLocation
			yMove := mouseYLocation - previousMouseYLocation
			moveLayerByRelativeValue(draggedLayerAlias, xMove, yMove)
			if isInteractiveLayerOffscreen(draggedLayerAlias) {
				moveLayerByRelativeValue(draggedLayerAlias, -xMove, -yMove)
			}
			isScreenUpdateRequired = true
		} else if characterEntry.AttributeEntry.CellType == constants.CellTypeFrameTop && getEventStateId() != constants.EventStateDragAndDrop {
			// Only start a drag if one is not already in progress.
			setEventStateId(constants.EventStateDragAndDrop)
			setDraggedLayerAlias(characterEntry.LayerAlias)
		}
	} else {
		setEventStateId(constants.EventStateNone)
		setDraggedLayerAlias("")
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
	eventStateMemory.mutex.Lock()
	defer eventStateMemory.mutex.Unlock()
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
