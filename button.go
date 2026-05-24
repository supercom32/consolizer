package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
)

/*
buttonHistoryType is a structure which allows you to track the history of button presses for specific layers and aliases.
*/
type buttonHistoryType struct {
	buttonAlias string
	layerAlias  string
}

/*
buttonHistory is a variable which stores the last recorded button press information.
*/
var buttonHistory buttonHistoryType

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
*/
type buttonType struct{}

var Button buttonType
var Buttons = memory.NewControlMemoryManager[types.ButtonEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
Delete is a method which removes a button instance from its memory manager.

Example:

	button.Delete()
*/
func (shared *ButtonInstanceType) Delete() *ButtonInstanceType {
	shared.BaseControlInstanceType.Delete()
	return nil
}

/*
AddToTabIndex is a method which adds the button to the tab index of its associated layer.

Example:

	button.AddToTabIndex()
*/
func (shared *ButtonInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeButton)
}

/*
IsPressed is a method which detects if the button was pressed. In order to obtain the button pressed
and to clear this state, you must call the GetPressed method. In addition, the following should be noted:

- If buttonHistory matches the current button, the state is cleared and true is returned.

Example:

	isPressed := button.IsPressed()
*/
func (shared *ButtonInstanceType) IsPressed() bool {
	if buttonHistory.layerAlias != "" && buttonHistory.buttonAlias != "" {
		if buttonHistory.layerAlias == shared.layerAlias && buttonHistory.buttonAlias == shared.controlAlias {
			for shared.IsStatePressed() {
			}

			buttonHistory.layerAlias = ""
			buttonHistory.buttonAlias = ""
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
	if buttonHistory.layerAlias != "" && buttonHistory.buttonAlias != "" {
		layerAlias := buttonHistory.layerAlias
		buttonAlias := buttonHistory.buttonAlias
		buttonHistory.layerAlias = ""
		buttonHistory.buttonAlias = ""
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
	buttonEntry := Buttons.Get(shared.layerAlias, shared.controlAlias)
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

  - If the width of your button is less than the length of your button label, then the width will automatically default to
    the width of your button label.

  - If the height of your button is less than 3 characters high, then the height will automatically default to the minimum
    of 3 characters.

Example:

	buttonInstance := Button.Add("layer1", "btn1", "Submit", style, 10, 5, 20, 3, true)
*/
func (shared *buttonType) Add(layerAlias string, buttonAlias string, buttonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isEnabled bool) ButtonInstanceType {
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
	Buttons.Add(layerAlias, buttonAlias, &buttonEntry)

	// Create associated tooltip (always created but disabled by default)
	tooltipInstance := Tooltip.Add(layerAlias, buttonEntry.TooltipAlias, "", styleEntry,
		buttonEntry.XLocation, buttonEntry.YLocation,
		buttonEntry.Width, buttonEntry.Height,
		buttonEntry.XLocation, buttonEntry.YLocation+buttonEntry.Height+1,
		buttonEntry.Width, 3,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.SetEnabled(false)
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
	Buttons.Remove(layerAlias, buttonAlias)
}

/*
DeleteAll is a method which deletes all buttons on a given text layer.

Example:

	Button.DeleteAll("layer1")
*/
func (shared *buttonType) DeleteAll(layerAlias string) {
	Buttons.RemoveAll(layerAlias)
}

/*
drawOnLayer is a method which draws all buttons on a given text layer.

Example:

	Button.drawOnLayer(myLayer)
*/
func (shared *buttonType) drawOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	buttons := Buttons.GetAllEntries(layerAlias)
	for _, buttonEntry := range buttons {
		shared.draw(&layerEntry, buttonEntry.Alias, buttonEntry.Label, buttonEntry.StyleEntry, buttonEntry.IsPressed, buttonEntry.IsSelected, buttonEntry.IsEnabled, buttonEntry.XLocation, buttonEntry.YLocation, buttonEntry.Width, buttonEntry.Height)
	}
}

/*
draw is a method which draws a button on a given text layer. The style of the button will be
determined by the style entry passed in. In addition, the following should be noted:

  - Buttons are not drawn physically to the text layer provided. Instead, they are rendered to the terminal at the
    same time when the text layer is rendered.

  - If the button to be drawn falls outside the range of the provided layer, then only the visible portion of the
    button will be drawn.

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
	if height < 3 {
		height = 3
	}
	arrayOfRunes := stringformat.GetRunesFromString(buttonLabel)
	labelWidth := stringformat.GetWidthOfRunesWhenPrinted(arrayOfRunes)
	if width-2 <= labelWidth {
		width = labelWidth + 2
	}
	localStyleEntry.Window.LineDrawingTextForegroundColor = localStyleEntry.Button.RaisedColor
	localStyleEntry.Window.LineDrawingTextBackgroundColor = localStyleEntry.Button.BackgroundColor
	fillArea(layerEntry, attributeEntry, " ", xLocation, yLocation, width, height, constants.NullCellControlLocation)
	if isPressed {
		drawFrame(layerEntry, localStyleEntry, attributeEntry, constants.FrameStyleSunken, xLocation, yLocation, width, height, false)
	} else {
		drawFrame(layerEntry, localStyleEntry, attributeEntry, constants.FrameStyleRaised, xLocation, yLocation, width, height, false)
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
updateStateMouse is a method which updates button states that are triggered by mouse events.

Example:

	isUpdateNeeded := Button.updateStateMouse()
*/
func (shared *buttonType) updateStateMouse() bool {
	// If we're currently in a scrollbar drag operation, don't process button clicks
	if eventStateMemory.stateId == constants.EventStateDragAndDropScrollbar {
		return false
	}

	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	buttonAlias := characterEntry.AttributeEntry.CellControlAlias

	// If not a button, reset all buttons if needed.
	if characterEntry.AttributeEntry.CellType != constants.CellTypeButton {
		// GetLayer all buttons from all layers using ControlMemoryManager
		Buttons.MemoryManager.Range(func(key, value interface{}) bool {
			currentLayer := key.(string)
			buttons := Buttons.GetAllEntries(currentLayer)

			for _, buttonEntry := range buttons {
				// In case of delete race condition, we check if button exists
				if !Buttons.IsExists(currentLayer, buttonEntry.Alias) {
					continue
				}

				// If button is pressed, reset it
				if buttonEntry.IsPressed {
					buttonHistory.layerAlias = layerAlias
					buttonHistory.buttonAlias = buttonAlias
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

	if buttonAlias != "" && buttonPressed == 0 && Buttons.IsExists(layerAlias, buttonAlias) {
		buttonEntry := Buttons.Get(layerAlias, buttonAlias)
		if buttonEntry.IsPressed == true {
			buttonEntry.Mutex.Lock()
			buttonEntry.IsPressed = false
			buttonEntry.Mutex.Unlock()
			isUpdateRequired = true
		}
	} else if buttonAlias != "" && buttonPressed != 0 && Buttons.IsExists(layerAlias, buttonAlias) {
		// If button was found and mouse is being pressed, update button only
		// if required.
		buttonEntry := Buttons.Get(layerAlias, buttonAlias)
		if buttonEntry.IsEnabled && buttonEntry.IsPressed == false {
			buttonEntry.Mutex.Lock()
			buttonHistory.layerAlias = layerAlias
			buttonHistory.buttonAlias = buttonAlias
			buttonEntry.IsPressed = true
			buttonEntry.Mutex.Unlock()
			setFocusedControl(layerAlias, buttonAlias, constants.CellTypeButton)
			isUpdateRequired = true
		}
	}
	return isUpdateRequired
}
