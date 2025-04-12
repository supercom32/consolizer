package consolizer

import (
	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/internal/memory"
	"supercom32.net/consolizer/internal/stringformat"
	"supercom32.net/consolizer/types"
)

type buttonHistoryType struct {
	buttonAlias string
	layerAlias  string
}

var buttonHistory buttonHistoryType

type ButtonInstanceType struct {
	layerAlias   string
	controlAlias string
}

type buttonType struct{}

var Button buttonType
var Buttons = memory.NewControlMemoryManager[types.ButtonEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

func (shared *ButtonInstanceType) Delete() *ButtonInstanceType {
	if Buttons.IsExists(shared.layerAlias, shared.controlAlias) {
		Buttons.Remove(shared.layerAlias, shared.controlAlias)
	}
	return nil
}

func (shared *ButtonInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeButton)
}

/*
IsButtonPressed allows you to detect if any text button was pressed or not. In
order to obtain the button pressed and to clear this state, you must call the
GetButtonPressed method.
*/
func (shared *ButtonInstanceType) IsButtonPressed() bool {
	if buttonHistory.layerAlias != "" && buttonHistory.buttonAlias != "" {
		if buttonHistory.layerAlias == shared.layerAlias && buttonHistory.buttonAlias == shared.controlAlias {
			for shared.IsButtonStatePressed() {
			}

			buttonHistory.layerAlias = ""
			buttonHistory.buttonAlias = ""
			return true
		}
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

func (shared *ButtonInstanceType) IsButtonStatePressed() bool {
	buttonEntry := Buttons.Get(shared.layerAlias, shared.controlAlias)
	if buttonEntry.IsPressed == true {
		return true
	}
	return false
}

func (shared *ButtonInstanceType) SetEnabled(isEnabled bool) {
	buttonEntry := Buttons.Get(shared.layerAlias, shared.controlAlias)
	buttonEntry.IsEnabled = isEnabled
}

/*
Add allows you to add a button to a text layer. Once called, an instance of your control is
returned which will allow you to read or manipulate the properties for it. The Style of the button
will be determined by the style entry passed in. If you wish to remove a button from a text
layer, simply call 'DeleteButton'. In addition, the following information should be noted:

- Button are not drawn physically to the text layer provided. Instead
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
func (shared *buttonType) Add(layerAlias string, buttonAlias string, buttonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isEnabled bool) ButtonInstanceType {
	buttonEntry := types.NewButtonEntry()
	buttonEntry.StyleEntry = styleEntry
	buttonEntry.Alias = buttonAlias
	buttonEntry.Label = buttonLabel
	buttonEntry.XLocation = xLocation
	buttonEntry.YLocation = yLocation
	buttonEntry.IsEnabled = true
	buttonEntry.Width = width
	buttonEntry.Height = height
	buttonEntry.TooltipAlias = stringformat.GetLastSortedUUID()
	// Use the ControlMemoryManager to handle button entries
	Buttons.Add(layerAlias, buttonAlias, &buttonEntry)

	// Create associated tooltip (always created but disabled by default)
	Tooltip.Add(layerAlias, buttonEntry.TooltipAlias, "", styleEntry,
		buttonEntry.XLocation, buttonEntry.YLocation,
		buttonEntry.Width, buttonEntry.Height,
		buttonEntry.XLocation, buttonEntry.YLocation+buttonEntry.Height+1,
		buttonEntry.Width, 3,
		false, true, constants.DefaultTooltipHoverTime)

	var buttonInstance ButtonInstanceType
	buttonInstance.layerAlias = layerAlias
	buttonInstance.controlAlias = buttonAlias
	return buttonInstance
}

/*
DeleteButton allows you to remove a button from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a button which does not exist, then the request
will simply be ignored.
*/
func (shared *buttonType) DeleteButton(layerAlias string, buttonAlias string) {
	Buttons.Remove(layerAlias, buttonAlias)
}

/*
DeleteAllButtons allows you to delete all buttons on a given text layer.
*/
func (shared *buttonType) DeleteAllButtons(layerAlias string) {
	Buttons.RemoveAll(layerAlias)
}

/*
drawButtonsOnLayer allows you to draw all buttons on a given text layer.
*/
func (shared *buttonType) drawButtonsOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	buttons := Buttons.GetAllEntries(layerAlias)
	for _, buttonEntry := range buttons {
		drawButton(&layerEntry, buttonEntry.Alias, buttonEntry.Label, buttonEntry.StyleEntry, buttonEntry.IsPressed, buttonEntry.IsSelected, buttonEntry.IsEnabled, buttonEntry.XLocation, buttonEntry.YLocation, buttonEntry.Width, buttonEntry.Height)
	}
}

/*
drawButton allows you to draw A button on a given text layer. The
Style of the button will be determined by the style entry passed in. In
addition, the following information should be noted:

- Button are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create buttons without actually overwriting
the text layer data under it.

- If the button to be drawn falls outside the range of the provided layer,
then only the visible portion of the button will be drawn.
*/
func drawButton(layerEntry *types.LayerEntryType, buttonAlias string, buttonLabel string, styleEntry types.TuiStyleEntryType, isPressed bool, isSelected bool, isEnabled bool, xLocation int, yLocation int, width int, height int) {
	localStyleEntry := types.NewTuiStyleEntry(&styleEntry)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.ButtonForegroundColor
	attributeEntry.BackgroundColor = styleEntry.ButtonBackgroundColor
	attributeEntry.CellType = constants.CellTypeButton
	attributeEntry.CellControlAlias = buttonAlias
	if height < 3 {
		height = 3
	}
	if width-2 <= len(buttonLabel) {
		width = len(buttonLabel) + 2
	}
	localStyleEntry.LineDrawingTextForegroundColor = localStyleEntry.ButtonRaisedColor
	localStyleEntry.LineDrawingTextBackgroundColor = localStyleEntry.ButtonBackgroundColor
	fillArea(layerEntry, attributeEntry, " ", xLocation, yLocation, width, height, constants.NullCellControlLocation)
	if isPressed {
		drawFrame(layerEntry, localStyleEntry, attributeEntry, constants.FrameStyleSunken, xLocation, yLocation, width, height, false)
	} else {
		drawFrame(layerEntry, localStyleEntry, attributeEntry, constants.FrameStyleRaised, xLocation, yLocation, width, height, false)
	}
	centerXLocation := (width - len(buttonLabel)) / 2
	centerYLocation := height / 2
	arrayOfRunes := stringformat.GetRunesFromString(buttonLabel)
	if isSelected {
		attributeEntry.IsUnderlined = true
	}
	if !isEnabled {
		attributeEntry.ForegroundColor = styleEntry.ButtonLabelDisabledColor
	}
	printLayer(layerEntry, attributeEntry, xLocation+centerXLocation, yLocation+centerYLocation, arrayOfRunes)
}

/*
updateButtonStates allows you to update the state of all buttons. This needs
to be called when input occurs so that changes in button state are reflected
to the user as quickly as possible. In the event that a screen update is
required this method returns true.
*/
func (shared *buttonType) updateButtonStates(isMouseTriggered bool) bool {
	if isMouseTriggered {
		// Update the button state if a mouse caused a change.
		return shared.updateButtonStateMouse()
	} else {
		// AddLayer code to update when keyboard caused a change.
	}
	return false
}

/*
updateButtonStateMouse allows you to update button states that are triggered
by mouse events. If a screen update is required, then this method returns
true.
*/
func (shared *buttonType) updateButtonStateMouse() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	buttonAlias := characterEntry.AttributeEntry.CellControlAlias

	// If not a button, reset all buttons if needed.
	if characterEntry.AttributeEntry.CellType != constants.CellTypeButton {
		// GetLayer all buttons from all layers using ControlMemoryManager
		for currentLayer := range Buttons.MemoryManager {
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
		}
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

// GetBounds returns the position and size of the button
func (shared *ButtonInstanceType) GetBounds() (int, int, int, int) {
	buttonEntry := Buttons.Get(shared.layerAlias, shared.controlAlias)
	if buttonEntry == nil {
		return 0, 0, 0, 0
	}
	return buttonEntry.XLocation, buttonEntry.YLocation, buttonEntry.Width, buttonEntry.Height
}

// SetPosition sets the position of the button
func (shared *ButtonInstanceType) SetPosition(x, y int) *ButtonInstanceType {
	buttonEntry := Buttons.Get(shared.layerAlias, shared.controlAlias)
	if buttonEntry != nil {
		buttonEntry.XLocation = x
		buttonEntry.YLocation = y
	}
	return shared
}

// SetSize sets the dimensions of the button
func (shared *ButtonInstanceType) SetSize(width, height int) *ButtonInstanceType {
	buttonEntry := Buttons.Get(shared.layerAlias, shared.controlAlias)
	if buttonEntry != nil {
		buttonEntry.Width = width
		buttonEntry.Height = height
	}
	return shared
}

// SetVisible shows or hides the button
func (shared *ButtonInstanceType) SetVisible(visible bool) *ButtonInstanceType {
	buttonEntry := Buttons.Get(shared.layerAlias, shared.controlAlias)
	if buttonEntry != nil {
		buttonEntry.IsVisible = visible
	}
	return shared
}

// SetStyle sets the visual style of the button
func (shared *ButtonInstanceType) SetStyle(style types.TuiStyleEntryType) *ButtonInstanceType {
	buttonEntry := Buttons.Get(shared.layerAlias, shared.controlAlias)
	if buttonEntry != nil {
		buttonEntry.StyleEntry = style
	}
	return shared
}

// SetTabIndex sets the tab order of the button
func (shared *ButtonInstanceType) SetTabIndex(index int) *ButtonInstanceType {
	buttonEntry := Buttons.Get(shared.layerAlias, shared.controlAlias)
	if buttonEntry != nil {
		buttonEntry.TabIndex = index
	}
	return shared
}
