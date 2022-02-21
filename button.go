package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
)

type buttonHistoryType struct {
	buttonAlias string
	layerAlias  string
}

var buttonHistory buttonHistoryType

type ButtonInstanceType struct {
	layerAlias  string
	buttonAlias string
}

/*
IsButtonPressed allows you to detect if any text button was pressed or not. In
order to obtain the button pressed and clear this state, you must call the
GetButtonPressed method.
*/
func (shared *ButtonInstanceType) IsButtonPressed() bool {
	if buttonHistory.layerAlias != "" && buttonHistory.buttonAlias != "" {
		return true
	}
	return false
}

/*
GetButtonPressed allows you to detect which text button was pressed or not. In
the event no button was pressed, empty values for the layer and button
alias are returned instead. In addition, the following information should be
noted:

- If any button is successfully returned, the pressed state is automatically
cleared.
*/
func (shared *ButtonInstanceType) GetButtonPressed() (string, string) {
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
AddButton allows you to add a button to a text layer. The Style of the button
will be determined by the style entry passed in. If you wish to remove a
button from a text layer, simply call 'DeleteButton'. In addition, the
following information should be noted:

- Buttons are not drawn physically to the text layer provided. Instead
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create buttons without actually overwriting
the text layer data under it.

- If the button to be drawn falls outside the range of the provided layer,
then only the visible portion of the button will be drawn.

- If the width of your button is less than the length of your button label,
then the width will automatically default to the width of your button label.

- If the height of your button is less than 3 characters high, then the height
will automatically default to the minimum of 3 characters.
*/
 func AddButton(layerAlias string, buttonAlias string, buttonLabel string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, width int, height int) ButtonInstanceType {
	memory.AddButton(layerAlias, buttonAlias, buttonLabel, styleEntry, xLocation, yLocation, width, height)
	var buttonInstance ButtonInstanceType
	buttonInstance.layerAlias = layerAlias
	buttonInstance.buttonAlias = buttonAlias
	return buttonInstance
}

/*
DeleteButton allows you to remove a button from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a button which does not exist, then the request
will simply be ignored.
*/
func DeleteButton(layerAlias string, buttonAlias string) {
	memory.DeleteButton(layerAlias, buttonAlias)
}

/*
drawButtonsOnLayer allows you to draw all buttons on a given text layer
entry.
*/
func drawButtonsOnLayer(layerEntry memory.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for currentKey := range memory.ButtonMemory[layerAlias] {
		buttonEntry := memory.GetButton(layerAlias, currentKey)
		drawButton(&layerEntry, currentKey, buttonEntry.ButtonLabel, buttonEntry.StyleEntry, buttonEntry.IsPressed, buttonEntry.IsSelected, buttonEntry.XLocation, buttonEntry.YLocation, buttonEntry.Width, buttonEntry.Height)
	}
}

/*
updateButtonStates allows you to update the state of all buttons. This needs
to be called when input occurs so that changes in button state are reflected
to the user as quickly as possible. In the event that a screen update is
required this method returns true.
*/
func updateButtonStates(isMouseTriggered bool) bool {
	if isMouseTriggered {
		// Update the button state if a mouse caused a change.
		return updateButtonStateMouse()
	} else {
		// Add code to update when keyboard caused a change.
	}
	return false
}

/*
updateButtonStateMouse allows you to update button states that are triggered
by mouse events. If a screen update is required, then this method returns
true.
*/
func updateButtonStateMouse() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	buttonAlias := characterEntry.AttributeEntry.CellControlAlias
	// If not a button, reset all buttons if needed.
	if characterEntry.AttributeEntry.CellType != constants.CellTypeButton {
		for currentLayer, _ := range memory.ButtonMemory {
			for currentButton, _ := range memory.ButtonMemory[currentLayer] {
				// If button not found, reset all buttons if necessary...
				buttonEntry := memory.GetButton(currentLayer,currentButton)
				if buttonEntry.IsPressed == true {
					buttonHistory.layerAlias = layerAlias
					buttonHistory.buttonAlias = buttonAlias
					buttonEntry.Mutex.Lock()
					buttonEntry.IsPressed = false
					buttonEntry.Mutex.Unlock()
					isUpdateRequired = true
				}
			}
		}
		return isUpdateRequired
	}
	if buttonAlias != "" && buttonPressed == 0 {
		// If button was found, but mouse is not being pressed, update button
		// only if required.
		buttonEntry := memory.GetButton(layerAlias, buttonAlias)
		if buttonEntry.IsPressed == true {
			buttonEntry.Mutex.Lock()
			//memory.scrollBarMemory[layerAlias][buttonAlias].IsSelected = true
			buttonEntry.IsPressed = false
			buttonEntry.Mutex.Unlock()
			isUpdateRequired = true
		}
	} else if buttonAlias != "" && buttonPressed != 0 {
		// If button was found and mouse is being pressed, update button only
		// if required.
		buttonEntry := memory.GetButton(layerAlias, buttonAlias)
		if  buttonEntry.IsPressed == false {
			buttonEntry.Mutex.Lock()
			//memory.scrollBarMemory[layerAlias][buttonAlias].IsSelected = true
			buttonEntry.IsPressed = true
			buttonEntry.Mutex.Unlock()
			setFocusedControl(layerAlias, buttonAlias, constants.CellTypeButton)
			isUpdateRequired = true
		}
	}
	return isUpdateRequired
}
