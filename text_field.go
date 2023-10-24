package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
	"github.com/supercom32/consolizer/types"
	"strings"
)

type textFieldInstanceType struct {
	layerAlias     string
	textFieldAlias string
}

type textFieldType struct{}

var TextField textFieldType

/*
GetValue allows you to get the current value of your text field with.
*/
func (shared *textFieldInstanceType) GetValue() string {
	if memory.IsTextFieldExists(shared.layerAlias, shared.textFieldAlias) {
		validatorTextField(shared.layerAlias, shared.textFieldAlias)
		textFieldEntry := memory.GetTextField(shared.layerAlias, shared.textFieldAlias)
		value := textFieldEntry.CurrentValue
		return strings.TrimSpace(string(value))
	}
	return ""
}

/*
SetLocation allows you to set the current location of your text field.
*/
func (shared *textFieldInstanceType) SetLocation(xLocation int, yLocation int) {
	if memory.IsTextFieldExists(shared.layerAlias, shared.textFieldAlias) {
		validatorTextField(shared.layerAlias, shared.textFieldAlias)
		validateLayerLocationByLayerAlias(shared.layerAlias, xLocation, yLocation)
		textFieldEntry := memory.GetTextField(shared.layerAlias, shared.textFieldAlias)
		textFieldEntry.XLocation = xLocation
		textFieldEntry.YLocation = yLocation
	}
}

/*
Add allows you to add a text field to a given layer. Once called,
a text field instance is returned which will allow you to read or
manipulate properties of your text field. In addition, the following
information should be noted:

- If the location specified for the text field falls outside the range
of the text layer, then only the visible portion of your input field will be
drawn.

- If the max length of your text field is less than or equal to 0, a panic
will be generated to fail as fast as possible.

- Password protection will echo back '*' characters to the terminal instead
of the actual characters entered.

- Specifying a default value will simply pre-populate the text field with
the value specified.

- If the cursor position moves outside the visible display area of the
field, then the entire text field will shift to ensure the cursor is always
visible.
*/
func (shared *textFieldType) Add(layerAlias string, textFieldAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, maxLengthAllowed int, IsPasswordProtected bool, defaultValue string, isEnabled bool) textFieldInstanceType {
	validateLayerLocationByLayerAlias(layerAlias, xLocation, yLocation)
	validateTextFieldWidth(width)
	memory.AddTextField(layerAlias, textFieldAlias, styleEntry, xLocation, yLocation, width, maxLengthAllowed, IsPasswordProtected, defaultValue+" ", isEnabled)
	var textFieldInstance textFieldInstanceType
	textFieldInstance.layerAlias = layerAlias
	textFieldInstance.textFieldAlias = textFieldAlias
	return textFieldInstance
}

/*
DeleteTextField allows you to delete a text field on a given layer.
*/
func (shared *textFieldType) DeleteTextField(layerAlias string, textFieldAlias string) {
	validatorTextField(layerAlias, textFieldAlias)
	memory.DeleteTextField(layerAlias, textFieldAlias)
}

func (shared *textFieldType) DeleteAllTextFields(layerAlias string) {
	memory.DeleteAllTextFieldsFromLayer(layerAlias)
}

/*
drawTextFieldOnLayer allows you to draw all text fields on a given text layer
entry.
*/
func (shared *textFieldType) drawTextFieldOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for currentKey := range memory.TextField.Entries[layerAlias] {
		textFieldEntry := memory.GetTextField(layerAlias, currentKey)
		shared.drawInputString(&layerEntry, textFieldEntry.StyleEntry, currentKey, textFieldEntry.XLocation, textFieldEntry.YLocation, textFieldEntry.Width, textFieldEntry.ViewportPosition, textFieldEntry.CurrentValue)
	}
}

/*
drawInputString allows you to draw a string for an input field. This is
different from regular printing, since input fields are usually
constrained for space and have the possibility of not being able to
show the entire string. In addition, the following information should be
noted:

- If the location specified for the input string falls outside the range of
the text layer, then only the visible portion will be displayed.

- HotspotWidth indicates how large the visible area of your input string should be.

- String position indicates the location in your string where printing should
start. If the remainder of your string is too long for the specified width,
then only the visible portion will be displayed.
*/
func (shared *textFieldType) drawInputString(layerEntry *types.LayerEntryType, styleEntry types.TuiStyleEntryType, textFieldAlias string, xLocation int, yLocation int, width int, stringPosition int, inputValue []rune) {
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.TextFieldForegroundColor
	attributeEntry.BackgroundColor = styleEntry.TextFieldBackgroundColor
	attributeEntry.CellType = constants.CellTypeTextField
	attributeEntry.CellControlAlias = textFieldAlias
	numberOfCharactersToSafelyPrint := stringformat.GetMaxCharactersThatFitInStringSize(inputValue[stringPosition:], width)
	textFieldEntry := memory.GetTextField(layerEntry.LayerAlias, textFieldAlias)
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	fillArea(layerEntry, attributeEntry, " ", xLocation, yLocation, width, 1, 0)
	// Here we loop over each character to draw since we need to accommodate for unique
	// cell IDs (if required for mouse location detection).
	xLocationOffset := 0
	for currentRuneIndex := 0; currentRuneIndex < len(numberOfCharactersToSafelyPrint); currentRuneIndex++ {
		if focusedControlType == constants.CellTypeTextField && focusedLayerAlias == layerEntry.LayerAlias && focusedControlAlias == textFieldAlias {
			if stringPosition+currentRuneIndex == textFieldEntry.CursorPosition {
				attributeEntry.ForegroundColor = styleEntry.TextFieldCursorForegroundColor
				attributeEntry.BackgroundColor = styleEntry.TextFieldCursorBackgroundColor
			} else {
				attributeEntry.ForegroundColor = styleEntry.TextFieldForegroundColor
				attributeEntry.BackgroundColor = styleEntry.TextFieldBackgroundColor
			}
		}
		attributeEntry.CellControlId = stringPosition + currentRuneIndex
		if textFieldEntry.IsPasswordProtected {
			// If the field is password protected, then do not print the terminating ' ' character with an *.
			if xLocation+xLocationOffset == len(textFieldEntry.CurrentValue)-1 {
				printLayer(layerEntry, attributeEntry, xLocation+xLocationOffset, yLocation, []rune{' '})
			} else {
				printLayer(layerEntry, attributeEntry, xLocation+xLocationOffset, yLocation, []rune{'*'})
			}
		} else {
			printLayer(layerEntry, attributeEntry, xLocation+xLocationOffset, yLocation, []rune{inputValue[stringPosition+currentRuneIndex]})
		}
		xLocationOffset++
		if stringformat.IsRuneCharacterWide(inputValue[stringPosition+currentRuneIndex]) {
			xLocationOffset++
			if textFieldEntry.IsPasswordProtected {
				// If the field is password protected, then do not print the terminating ' ' character with an *.
				if xLocation+xLocationOffset == len(textFieldEntry.CurrentValue)-1 {
					printLayer(layerEntry, attributeEntry, xLocation+xLocationOffset, yLocation, []rune{' '})
				} else {
					printLayer(layerEntry, attributeEntry, xLocation+xLocationOffset, yLocation, []rune{'*'})
				}
			} else {
				printLayer(layerEntry, attributeEntry, xLocation+xLocationOffset, yLocation, []rune{' '})
			}
		}
	}
}

/*
updateTextFieldViewport allows you to update the current viewport based on the current text and cursor
location.
*/
func (shared *textFieldType) updateTextFieldViewport(textFieldEntry *types.TextFieldEntryType) {
	// If cursor xLocation is lower than the viewport window
	if textFieldEntry.CursorPosition <= textFieldEntry.ViewportPosition {
		maxViewportWidthAvaliable := textFieldEntry.Width
		if textFieldEntry.CursorPosition-textFieldEntry.Width < 0 {
			maxViewportWidthAvaliable = textFieldEntry.CursorPosition
		}
		arrayOfRunesAvailableToPrint := textFieldEntry.CurrentValue[textFieldEntry.CursorPosition-maxViewportWidthAvaliable : textFieldEntry.CursorPosition]
		numberOfRunesThatFitStringSize := stringformat.GetMaxCharactersThatFitInStringSizeReverse(arrayOfRunesAvailableToPrint, textFieldEntry.Width)
		textFieldEntry.ViewportPosition = textFieldEntry.CursorPosition - numberOfRunesThatFitStringSize
	}
	// Figure out how much displayable space is in our current viewport window.
	arrayOfRunesAvailableToPrint := textFieldEntry.CurrentValue[textFieldEntry.ViewportPosition:]
	arrayOfRunesThatFitStringSize := stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvailableToPrint, textFieldEntry.Width)
	// If the cursor xLocation is equal or greater than the visible viewport window width.
	if textFieldEntry.CursorPosition >= textFieldEntry.ViewportPosition+len(arrayOfRunesThatFitStringSize) {
		// Then make the viewport xLocation equal to the visible viewport width behind it.
		maxViewportWidthAvaliable := textFieldEntry.Width
		if textFieldEntry.CursorPosition-textFieldEntry.Width < 0 {
			maxViewportWidthAvaliable = textFieldEntry.CursorPosition
		}
		arrayOfRunesAvailableToPrint = textFieldEntry.CurrentValue[textFieldEntry.CursorPosition-maxViewportWidthAvaliable : textFieldEntry.CursorPosition]
		numberOfRunesThatFitStringSize := stringformat.GetMaxCharactersThatFitInStringSizeReverse(arrayOfRunesAvailableToPrint, textFieldEntry.Width)
		// LogInfo(fmt.Sprintf("v: %d x: %d off: %d fit: %d, aval: %s", textFieldEntry.ViewportPosition, textFieldEntry.CursorPosition, maxViewportWidthAvaliable, numberOfRunesThatFitStringSize, string(arrayOfRunesAvailableToPrint)))
		textFieldEntry.ViewportPosition = textFieldEntry.CursorPosition - numberOfRunesThatFitStringSize + 1
	}
}

/*
insertCharacterAtPosition allows you to insert a character into a given text field. The location to insert is
determined automatically by the current cursor position.
*/
func (shared *textFieldType) insertCharacterAtPosition(textFieldEntry *types.TextFieldEntryType, characterToInsert rune) {
	textAfterCursor := stringformat.GetRuneArrayCopy(textFieldEntry.CurrentValue[textFieldEntry.CursorPosition:])
	textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition], characterToInsert)
	textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue, textAfterCursor...)
}

/*
deleteCharacterAtPosition allows you to delete a character from a given text field. The location to delete is
determined automatically by the current cursor position.
*/
func (shared *textFieldType) deleteCharacterAtPosition(textFieldEntry *types.TextFieldEntryType) {
	if len(textFieldEntry.CurrentValue) != 1 {
		if textFieldEntry.CursorPosition != len(textFieldEntry.CurrentValue)-1 {
			textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition], textFieldEntry.CurrentValue[textFieldEntry.CursorPosition+1:]...)
		}
	}
}

/*
backspaceCharacterAtPosition allows you to backspace a character from a given text field. The location to backspace is
determined automatically by the current cursor position.
*/
func (shared *textFieldType) backspaceCharacterAtPosition(textFieldEntry *types.TextFieldEntryType) {
	if textFieldEntry.CursorPosition >= 0 {
		textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition], textFieldEntry.CurrentValue[textFieldEntry.CursorPosition+1:]...)
	}
}

/*
updateTextFieldCursor allows you to update a text field cursors location to ensure that no values are out of bounds.
*/
func (shared *textFieldType) updateTextFieldCursor(textFieldEntry *types.TextFieldEntryType) {
	if textFieldEntry.CursorPosition == constants.NullCellControlId || textFieldEntry.CursorPosition >= len(textFieldEntry.CurrentValue) {
		textFieldEntry.CursorPosition = len(textFieldEntry.CurrentValue) - 1 // used to be -1 here.
	}
	if textFieldEntry.CursorPosition < 0 {
		textFieldEntry.CursorPosition = 0
	}
}

/*
updateKeyboardEventTextField allows you to update the state of all text fields according to the current keystroke event.
In the event that a screen update is required this method returns true.
*/
func (shared *textFieldType) updateKeyboardEventTextField(keystroke []rune) bool {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := true
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	if focusedControlType != constants.CellTypeTextField || !memory.IsTextFieldExists(focusedLayerAlias, focusedControlAlias) {
		return false
	}
	textFieldEntry := memory.GetTextField(focusedLayerAlias, focusedControlAlias)
	if !textFieldEntry.IsEnabled {
		return false
	}
	if len(keystroke) == 1 { // If a regular char is entered.
		if len(textFieldEntry.CurrentValue) < textFieldEntry.MaxLengthAllowed {
			// LogInfo(fmt.Sprintf("cur: %d view: %d", textFieldEntry.CursorPosition, textFieldEntry.ViewportPosition))
			shared.insertCharacterAtPosition(textFieldEntry, keystroke[0])
			textFieldEntry.CursorPosition++
			shared.updateTextFieldCursor(textFieldEntry)
			shared.updateTextFieldViewport(textFieldEntry)
			isScreenUpdateRequired = true
		}
	}
	if keystrokeAsString == "delete" {
		shared.deleteCharacterAtPosition(textFieldEntry)
		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "home" {
		textFieldEntry.CursorPosition = 0
		textFieldEntry.ViewportPosition = 0
		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "end" {
		textFieldEntry.CursorPosition = len(textFieldEntry.CurrentValue) - 1
		shared.updateTextFieldViewport(textFieldEntry)
		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "backspace" || keystrokeAsString == "backspace2" {
		textFieldEntry.CursorPosition--
		shared.backspaceCharacterAtPosition(textFieldEntry)
		shared.updateTextFieldCursor(textFieldEntry)
		shared.updateTextFieldViewport(textFieldEntry)

		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "left" {
		textFieldEntry.CursorPosition--
		shared.updateTextFieldCursor(textFieldEntry)
		shared.updateTextFieldViewport(textFieldEntry)

		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "right" {
		textFieldEntry.CursorPosition++
		shared.updateTextFieldCursor(textFieldEntry)
		shared.updateTextFieldViewport(textFieldEntry)
		isScreenUpdateRequired = true
	}
	return isScreenUpdateRequired
}

/*
updateMouseEventTextField allows you to update the state of all text fields according to the current mouse event state.
In the event that a screen update is required this method returns true.
*/
func (shared *textFieldType) updateMouseEventTextField() bool {
	isScreenUpdateRequired := false
	var characterEntry types.CharacterEntryType
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	if buttonPressed != 0 && eventStateMemory.stateId != constants.EventStateDragAndDropScrollbar {
		characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if characterEntry.AttributeEntry.CellType == constants.CellTypeTextField && memory.IsTextFieldExists(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias) {
			textFieldEntry := memory.GetTextField(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			if !textFieldEntry.IsEnabled {
				return isScreenUpdateRequired
			}
			textFieldEntry.CursorPosition = characterEntry.AttributeEntry.CellControlId
			// LogInfo(strconv.Itoa(textFieldEntry.CursorPosition))
			shared.updateTextFieldCursor(textFieldEntry)
			setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeTextField)
			isScreenUpdateRequired = true
		}
	}
	return isScreenUpdateRequired
}
