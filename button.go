package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
	"sync"
)

/*
buttonHistoryType is a structure which allows you to track the history of button presses for specific layers and aliases.
In addition, the following should be noted:

  - Every field access goes through the embedded mutex, since the event goroutine, the periodic-event goroutine, and
    the application's own goroutine can all read or write it concurrently.

Example:
    var history buttonHistoryType
*/
type buttonHistoryType struct {
	mutex       sync.Mutex
	buttonAlias string
	layerAlias  string
}

/*
set is a method which records a new button press location, replacing whatever was previously recorded. In addition,
the following should be noted:

  - Callers must not hold a button entry's own Mutex when calling this, since that lock ordering is also taken in
    the reverse direction elsewhere and would risk a deadlock.

Example:
    buttonHistory.set("layer1", "myButton")
*/
func (shared *buttonHistoryType) set(layerAlias string, buttonAlias string) {
	shared.mutex.Lock()
	defer shared.mutex.Unlock()
	shared.layerAlias = layerAlias
	shared.buttonAlias = buttonAlias
}

/*
clear is a method which erases the currently recorded button press, if any.

Example:
    buttonHistory.clear()
*/
func (shared *buttonHistoryType) clear() {
	shared.set("", "")
}

/*
get is a method which returns the currently recorded button press without clearing it.

Example:
    layerAlias, buttonAlias := buttonHistory.get()
*/
func (shared *buttonHistoryType) get() (string, string) {
	shared.mutex.Lock()
	defer shared.mutex.Unlock()
	return shared.layerAlias, shared.buttonAlias
}

/*
getAndClear is a method which returns the currently recorded button press and clears it in a single locked step, so
that no other goroutine can observe or overwrite the value in between the read and the clear.

Example:
    layerAlias, buttonAlias := buttonHistory.getAndClear()
*/
func (shared *buttonHistoryType) getAndClear() (string, string) {
	shared.mutex.Lock()
	defer shared.mutex.Unlock()
	layerAlias := shared.layerAlias
	buttonAlias := shared.buttonAlias
	shared.layerAlias = ""
	shared.buttonAlias = ""
	return layerAlias, buttonAlias
}

/*
buttonHistory is a variable which stores the last recorded button press information.

Example:
    buttonHistory.set("layer1", "myButton")
*/
var buttonHistory buttonHistoryType

// buttonPressOrigin is the button the held mouse press began on, or nil when the press began anywhere else or no
// button is held. It is only read and written while the UI state lock is held, so it needs no synchronization of its
// own.
var buttonPressOrigin *types.ButtonEntryType

/*
ButtonInstanceType is a structure which represents an instance of a button control.

Example:
    var buttonInstance ButtonInstanceType
*/
type ButtonInstanceType struct {
	BaseControlInstanceType
}

/*
buttonType is a structure which provides the global namespace for button management operations.

Example:
    var button buttonType
*/
type buttonType struct{}

var Button buttonType
var buttonRegistry = memory.NewControlMemoryManager[types.ButtonEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
Delete is a method which removes a button instance from its memory manager.

Example:
    button.Delete()
*/
func (shared *ButtonInstanceType) Delete() *ButtonInstanceType {
	defer lockUIState()()
	shared.BaseControlInstanceType.delete()
	return nil
}

/*
IsPressed is a method which detects if the button was pressed. In order to obtain the button pressed
and to clear this state, you must call the GetPressed method. In addition, the following should be noted:

- If buttonHistory matches the current button, the state is cleared and true is returned.

Example:
    isPressed := button.IsPressed()
*/
func (shared *ButtonInstanceType) IsPressed() bool {
	defer lockUIState()()
	layerAlias, buttonAlias := buttonHistory.get()
	if layerAlias != "" && buttonAlias != "" {
		if layerAlias == shared.layerAlias && buttonAlias == shared.controlAlias {
			for shared.isStatePressed() {
			}

			buttonHistory.clear()
			return true
		}
	}
	return false
}

/*
GetPressed is a method which detects which button was pressed. In the event no button was pressed,
empty values for the layer and button alias are returned instead. In addition, the following should be noted:

- If any button is successfully returned, the pressed state is automatically cleared.

Example:
    layerAlias, buttonAlias := button.GetPressed()
*/
func (shared *ButtonInstanceType) GetPressed() (string, string) {
	defer lockUIState()()
	layerAlias, buttonAlias := buttonHistory.getAndClear()
	if layerAlias != "" && buttonAlias != "" {
		return layerAlias, buttonAlias
	}
	return "", ""
}

/*
IsStatePressed is a method which checks the current internal pressed state of the button.

Example:
    isStatePressed := button.IsStatePressed()
*/
func (shared *ButtonInstanceType) IsStatePressed() bool {
	defer lockUIState()()
	return shared.isStatePressed()
}

/*
isStatePressed is a method which allows you to do exactly what IsStatePressed does, for a caller that already holds
uiStateMutex, such as another entry point of the package that has taken the lock itself.

Example:

	button.isStatePressed()
*/
func (shared *ButtonInstanceType) isStatePressed() bool {
	buttonEntry := buttonRegistry.Get(shared.layerAlias, shared.controlAlias)
	if buttonEntry.IsPressed == true {
		return true
	}
	return false
}

/*
Add is a method which adds a button to a text layer. Once called, an instance of your control is returned
which will allow you to read or manipulate the properties for it. The Style of the button will be determined by the
style entry passed in. If you wish to remove a button from a text layer, simply call 'DeleteButton'. In addition, the
following should be noted:

- Buttons are not drawn physically to the text layer provided. Instead they are rendered to the terminal at the same
  time when the text layer is rendered. This allows you to create buttons without actually overwriting the text layer
  data under it.

- If the button to be drawn falls outside the range of the provided layer, then only the visible portion of the button
  will be drawn.

- If the width of your button is less than the length of your button label, then the width will automatically
  default to the width of your button label.

- If the height of your button is less than 3 characters high, then the height will automatically default to the
  minimum of 3 characters.

Example:
    buttonInstance := Button.Add("layer1", "btn1", "Submit", style, 10, 5, 20, 3, true)
*/
func (shared *buttonType) Add(layerAlias string, buttonAlias string, buttonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isEnabled bool) ButtonInstanceType {
	defer lockUIState()()
	return shared.add(layerAlias, buttonAlias, buttonLabel, styleEntry, xLocation, yLocation, width, height, isEnabled)
}

/*
add is a method which allows you to do exactly what Add does, for a caller that already holds uiStateMutex, such as
another entry point of the package that has taken the lock itself.

Example:

	Button.add(layerAlias, buttonAlias, buttonLabel, styleEntry, xLocation, yLocation, width, height, isEnabled)
*/
func (shared *buttonType) add(layerAlias string, buttonAlias string, buttonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isEnabled bool) ButtonInstanceType {
	buttonEntry := types.NewButtonEntry()
	buttonEntry.StyleEntry = styleEntry
	buttonEntry.Alias = buttonAlias
	buttonEntry.Label = buttonLabel
	buttonEntry.XLocation = xLocation
	buttonEntry.YLocation = yLocation
	buttonEntry.IsEnabled = isEnabled
	buttonEntry.Width = width
	buttonEntry.Height = height
	buttonEntry.TooltipAlias = stringformat.GetLastSortedUUID()
	// Use the ControlMemoryManager to handle button entries
	buttonRegistry.Add(layerAlias, buttonAlias, &buttonEntry)

	// Create associated tooltip (always created but disabled by default)
	tooltipInstance := Tooltip.add(layerAlias, buttonEntry.TooltipAlias, "", styleEntry,
		buttonEntry.XLocation, buttonEntry.YLocation,
		buttonEntry.Width, buttonEntry.Height,
		buttonEntry.XLocation, buttonEntry.YLocation+buttonEntry.Height+1,
		buttonEntry.Width, 3,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.setEnabled(false)
	tooltipInstance.setParentControlAlias(buttonAlias)
	var buttonInstance ButtonInstanceType
	buttonInstance.layerAlias = layerAlias
	buttonInstance.controlAlias = buttonAlias
	buttonInstance.controlType = constants.TYPE_BUTTON
	return buttonInstance
}

/*
Delete is a method which removes a button from a text layer. In addition, the following should be noted:

- If you attempt to delete a button which does not exist, then the request will simply be ignored.

Example:
    Button.Delete("layer1", "btn1")
*/
func (shared *buttonType) Delete(layerAlias string, buttonAlias string) {
	defer lockUIState()()
	buttonRegistry.Remove(layerAlias, buttonAlias)
	clearStaleControlReferences(layerAlias, buttonAlias, constants.CellTypeButton)
}

/*
DeleteAll is a method which deletes all buttons on a given text layer.

Example:
    Button.DeleteAll("layer1")
*/
func (shared *buttonType) DeleteAll(layerAlias string) {
	defer lockUIState()()
	removeAllControlsAndClearCells(buttonRegistry, layerAlias, constants.CellTypeButton, func(entry *types.ButtonEntryType) string { return entry.Alias })
}

/*
drawOnLayer is a method which draws all buttons on a given text layer.

Example:
    Button.drawOnLayer(myLayer)
*/
func (shared *buttonType) drawOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	buttons := buttonRegistry.GetAllEntries(layerAlias)
	for _, buttonEntry := range buttons {
		shared.draw(&layerEntry, buttonEntry.Alias, buttonEntry.Label, buttonEntry.StyleEntry, buttonEntry.IsPressed, buttonEntry.IsSelected, buttonEntry.IsEnabled, buttonEntry.XLocation, buttonEntry.YLocation, buttonEntry.Width, buttonEntry.Height)
	}
}

/*
draw is a method which allows you to draw a button on a given text layer, in the style given by the style entry passed
in. Buttons are not drawn physically to the text layer provided. Instead, they are rendered to the terminal at the
same time as the text layer, and if the button falls outside the range of the layer, only its visible portion is
drawn. In addition, the following should be noted:

  - styleEntry.Button.StyleMode selects the look: beveled (default, two-tone 3D frame that flips when pressed),
    flat (single-colour frame, pressed state shown with the pressed colours), or borderless (no frame at all,
    pressed state shown with the pressed colours). Flat and borderless keep the three-row minimum only for the
    beveled and flat frames; a borderless button may be a single row.

  - While the button has focus and is not pressed, its face and label are drawn with the style's focused colours,
    resolved by getFocusedColors so that focus stays visible even when the style does not set them. As described
    for setFocusIndicatorVisible, these colours are only shown while the user's last input came from the keyboard,
    so clicking a button does not flash them over it.

Example:

	Button.draw(&myLayer, "btn1", "OK", style, false, false, true, 0, 0, 10, 3)
*/
func (shared *buttonType) draw(layerEntry *types.LayerEntryType, buttonAlias string, buttonLabel string, styleEntry types.TuiStyleEntryType, isPressed bool, isSelected bool, isEnabled bool, xLocation int, yLocation int, width int, height int) {
	localStyleEntry := types.NewTuiStyleEntry(&styleEntry)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.Button.ForegroundColor
	attributeEntry.BackgroundColor = styleEntry.Button.BackgroundColor
	attributeEntry.CellType = constants.CellTypeButton
	attributeEntry.CellControlAlias = buttonAlias

	styleMode := styleEntry.Button.StyleMode
	isBorderless := styleMode == constants.ButtonStyleBorderless

	minimumHeight := 3
	horizontalPadding := 2
	if isBorderless {
		minimumHeight = 1
		horizontalPadding = 0
	}
	if height < minimumHeight {
		height = minimumHeight
	}
	arrayOfRunes := stringformat.GetRunesFromString(buttonLabel)
	labelWidth := stringformat.GetWidthOfRunesWhenPrinted(arrayOfRunes)
	if width-horizontalPadding <= labelWidth {
		width = labelWidth + horizontalPadding
	}

	// A pressed flat or borderless button has no bevel to flip, so it signals the press through its colours.
	if isPressed && styleMode != constants.ButtonStyleBeveled {
		attributeEntry.ForegroundColor = styleEntry.Button.PressedForegroundColor
		attributeEntry.BackgroundColor = styleEntry.Button.PressedBackgroundColor
	}
	if !isPressed && isFocusIndicatorShown(layerEntry.LayerAlias, buttonAlias, constants.CellTypeButton) {
		attributeEntry.ForegroundColor, attributeEntry.BackgroundColor = getFocusedColors(styleEntry.Button.ForegroundColor,
			styleEntry.Button.BackgroundColor, styleEntry.Button.FocusedForegroundColor, styleEntry.Button.FocusedBackgroundColor)
	}

	fillArea(layerEntry, attributeEntry, " ", xLocation, yLocation, width, height, constants.NullCellControlLocation)

	switch styleMode {
	case constants.ButtonStyleBorderless:
		// No frame is drawn.
	case constants.ButtonStyleFlat:
		localStyleEntry.Window.LineDrawingTextForegroundColor = attributeEntry.ForegroundColor
		localStyleEntry.Window.LineDrawingTextBackgroundColor = attributeEntry.BackgroundColor
		drawFrame(layerEntry, localStyleEntry, attributeEntry, constants.FrameStyleNormal, xLocation, yLocation, width, height, false)
	default:
		localStyleEntry.Window.LineDrawingTextForegroundColor = localStyleEntry.Button.RaisedColor
		localStyleEntry.Window.LineDrawingTextBackgroundColor = localStyleEntry.Button.BackgroundColor
		if isPressed {
			drawFrame(layerEntry, localStyleEntry, attributeEntry, constants.FrameStyleSunken, xLocation, yLocation, width, height, false)
		} else {
			drawFrame(layerEntry, localStyleEntry, attributeEntry, constants.FrameStyleRaised, xLocation, yLocation, width, height, false)
		}
	}

	centerXLocation := (width - labelWidth) / 2
	centerYLocation := height / 2
	if isSelected {
		attributeEntry.IsUnderlined = true
	}
	if !isEnabled {
		attributeEntry.ForegroundColor = styleEntry.Button.LabelDisabledColor
	}
	layer.printLayer(layerEntry, attributeEntry, xLocation+centerXLocation, yLocation+centerYLocation, arrayOfRunes)
}

/*
updateStates is a method which updates the state of all buttons. This needs to be called when input
occurs so that changes in button state are reflected to the user as quickly as possible.

Example:
    isUpdateNeeded := Button.updateStates(true)
*/
func (shared *buttonType) updateStates(isMouseTriggered bool) bool {
	if isMouseTriggered {
		// Update the button state if a mouse caused a change.
		return shared.updateStateMouse()
	} else {
		// AddLayer code to update when keyboard caused a change.
	}
	return false
}

/*
updateKeyboardEvent is a method which allows you to press the currently focused button from the keyboard. When a button
has focus and the keystroke is Enter or Space, the press is recorded exactly as a completed mouse click would be, so
GetPressed and IsPressed report it. It returns whether a screen update is required and whether the keystroke was
consumed, and both are true only when a press was recorded. In addition, the following should be noted:

  - A disabled or hidden button, or one on a hidden layer, ignores the keystroke and leaves it unconsumed, so it
    still reaches the keyboard buffer.

  - A consumed Enter or Space is never added to the keyboard buffer, so the application does not see the same
    action twice.

Example:

	updateRequired, consumed := Button.updateKeyboardEvent([]rune("enter"))
*/
func (shared *buttonType) updateKeyboardEvent(keystroke []rune) (bool, bool) {
	keystrokeAsString := string(keystroke)
	if keystrokeAsString != "enter" && keystrokeAsString != " " {
		return false, false
	}
	focusedControl := getFocusedControl()
	if focusedControl.controlType != constants.CellTypeButton || !isControlFocusable(focusedControl) {
		return false, false
	}
	buttonHistory.set(focusedControl.layerAlias, focusedControl.controlAlias)
	return true, true
}

/*
clearAllPressedButtons is a method which resets the pressed state of every button across every layer. In addition,
the following should be noted:

  - This is used whenever a mouse release is observed, since the press may have started on a different button (or
    on no button at all) than the one currently under the cursor, and every stale pressed state must be cleared so
    that button does not get stuck ignoring its next click.

  - Buttons removed by a concurrent delete are skipped rather than acted on, since buttonRegistry.IsExists is checked
    immediately before each entry is touched.

Example:
    isUpdateRequired := Button.clearAllPressedButtons()
*/
func (shared *buttonType) clearAllPressedButtons() bool {
	isUpdateRequired := false
	buttonRegistry.MemoryManager.Range(func(key, value interface{}) bool {
		currentLayer := key.(string)
		buttons := buttonRegistry.GetAllEntries(currentLayer)

		for _, buttonEntry := range buttons {
			// In case of delete race condition, we check if button exists
			if !buttonRegistry.IsExists(currentLayer, buttonEntry.Alias) {
				continue
			}

			// If button is pressed, reset it
			if buttonEntry.IsPressed {
				buttonEntry.Mutex.Lock()
				buttonEntry.IsPressed = false
				buttonEntry.Mutex.Unlock()
				isUpdateRequired = true
			}
		}
		return true // continue iteration
	})
	return isUpdateRequired
}

/*
updatePressOrigin is a method which allows you to record which button the current mouse press began on. A fresh press
over a button records that button, a fresh press anywhere else records that no button owns the press, and a release
clears the record. Movement with the button still held leaves the record unchanged. In addition, the following should
be noted:

  - It must run for every mouse event, before anything that can stop the event from reaching the controls, such as a
    modal layer rejecting it, so that a press which was never delivered to a button cannot be mistaken for one that
    was.

  - The button is recorded by its entry rather than its alias, so a button deleted during the press, and one added
    later under the same alias, never inherit the press.

Example:

	Button.updatePressOrigin(characterEntry, buttonPressed, previousButtonPressed)
*/
func (shared *buttonType) updatePressOrigin(characterEntry types.CharacterEntryType, buttonPressed uint, previousButtonPressed uint) {
	if buttonPressed == 0 {
		buttonPressOrigin = nil
		return
	}
	if previousButtonPressed != 0 {
		return
	}
	buttonPressOrigin = nil
	if characterEntry.AttributeEntry.CellType != constants.CellTypeButton {
		return
	}
	if buttonEntry, isFound := buttonRegistry.Lookup(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias); isFound {
		buttonPressOrigin = buttonEntry
	}
}

/*
updateStateMouse is a method which allows you to update button states according to the current mouse event. It returns
true if the screen needs to be updated. In addition, the following should be noted:

  - Only the button a press began on, as recorded by updatePressOrigin, can be armed by it. A press that began
    elsewhere, or that was already held when the button was added, never arms a button it is dragged onto.

  - The visual pressed state (IsPressed, the sunken frame) is set the instant the press lands on an enabled button,
    so feedback stays immediate. Dragging off the button, onto any other cell or another button, disarms it, and
    dragging back onto it while still held re-arms it.

  - The click itself, recorded into buttonHistory for GetPressed and IsPressed to report, only fires on a release
    over the armed button. Any release clears every armed button, so none can get stuck pressed.

Example:

	isUpdateNeeded := Button.updateStateMouse()
*/
func (shared *buttonType) updateStateMouse() bool {
	// If we're currently in a scrollbar drag operation, don't process button clicks
	if getEventStateId() == constants.EventStateDragAndDropScrollbar {
		return false
	}

	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	buttonAlias := characterEntry.AttributeEntry.CellControlAlias
	var buttonEntry *types.ButtonEntryType
	isFound := false
	if characterEntry.AttributeEntry.CellType == constants.CellTypeButton {
		buttonEntry, isFound = buttonRegistry.Lookup(layerAlias, buttonAlias)
	}

	if buttonPressed == 0 {
		if isFound && buttonEntry.IsPressed {
			buttonHistory.set(layerAlias, buttonAlias)
		}
		return shared.clearAllPressedButtons()
	}
	if !isFound || buttonEntry != buttonPressOrigin {
		return shared.clearAllPressedButtons()
	}
	if !buttonEntry.IsEnabled || buttonEntry.IsPressed {
		return false
	}
	buttonEntry.Mutex.Lock()
	buttonEntry.IsPressed = true
	buttonEntry.Mutex.Unlock()
	setFocusedControl(layerAlias, buttonAlias, constants.CellTypeButton)
	return true
}
