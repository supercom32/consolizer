package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
)

type textFieldInstanceType struct {
	layerAlias string
	textFieldAlias string
}

/*
GetValue allows you to get the current value of your text field with.
*/
func (shared *textFieldInstanceType) GetValue() string {
	validatorTextField(shared.layerAlias, shared.textFieldAlias)
	textFieldEntry := memory.GetTextField(shared.layerAlias, shared.textFieldAlias)
	value := textFieldEntry.CurrentValue
	return string(value)
}

/*
SetLocation allows you to set the current location of your text field.
*/
func (shared *textFieldInstanceType) SetLocation(xLocation int, yLocation int) {
	validatorTextField(shared.layerAlias, shared.textFieldAlias)
	validateLayerLocationByLayerAlias(shared.layerAlias, xLocation, yLocation)
	textFieldEntry := memory.GetTextField(shared.layerAlias, shared.textFieldAlias)
	textFieldEntry.XLocation = xLocation
	textFieldEntry.YLocation = yLocation
}

/*
AddTextField allows you to add a text field to a given layer. Once called,
a text field instance is returned which will allow you to read or
manipulate properties of your text field. In addition, the following
information should be noted:

- If the location specified for the input field  falls outside the range
of the text layer, then only the visible portion of your input field will be
drawn.

- If the max length of your input field is less than or equal to 0, a panic
will be generated to fail as fast as possible.

- Password protection will echo back '*' characters to the terminal instead
of the actual characters entered.

- Specifying a default value will simply pre-populate the input field with
the value specified.

- If the cursor position moves outside the visible display area of the
field, then the entire input field will shift to ensure the cursor is always
visible.
*/
func AddTextField(layerAlias string, textFieldAlias string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, width int, maxLengthAllowed int, IsPasswordProtected bool, defaultValue string) textFieldInstanceType {
	validateLayerLocationByLayerAlias(layerAlias, xLocation, yLocation)
	validateTextFieldWidth(width)
	memory.AddTextField(layerAlias, textFieldAlias, styleEntry, xLocation, yLocation, width, maxLengthAllowed, IsPasswordProtected, defaultValue + " ")
	var textFieldInstance textFieldInstanceType
	textFieldInstance.layerAlias = layerAlias
	textFieldInstance.textFieldAlias = textFieldAlias
	return textFieldInstance
}

/*
DeleteTextField allows you to delete a text field on a given layer.
*/
func DeleteTextField(layerAlias string, textFieldAlias string) {
	validatorTextField(layerAlias, textFieldAlias)
	memory.DeleteTextField(layerAlias, textFieldAlias)
}

/*
drawTextFieldOnLayer allows you to draw all text fields on a given text layer
entry.
*/
func drawTextFieldOnLayer(layerEntry memory.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	/*
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	 */
	for currentKey := range memory.TextFieldMemory[layerAlias] {
		textFieldEntry := memory.GetTextField(layerAlias, currentKey)
		drawInputString(&layerEntry, textFieldEntry.StyleEntry, currentKey, textFieldEntry.XLocation, textFieldEntry.YLocation, textFieldEntry.Width, textFieldEntry.ViewportPosition, textFieldEntry.CurrentValue)
		/*
		cursorPosition := textFieldEntry.ViewportPosition + textFieldEntry.CursorPosition
		runeSlice := []rune(textFieldEntry.CurrentValue)
		characterUnderCursor := ' '
		if cursorPosition < len(runeSlice) { // Protect against empty strings
			characterUnderCursor = runeSlice[cursorPosition]
		}

		if focusedControlType == constants.CellTypeTextField && focusedLayerAlias == layerAlias && focusedControlAlias == currentKey {
			drawCursor(&layerEntry, textFieldEntry.StyleEntry, currentKey, characterUnderCursor, textFieldEntry.XLocation, textFieldEntry.YLocation, textFieldEntry.CursorPosition, false)
		}
		 */
	}
}

/*
drawCursor allows you to draw a cursor at the appropriate location for a
text field. In addition, the following information should be noted:

- If the location specified for the cursor range falls outside of the text
layer, then the cursor will only be rendered on the visible portion.

- The cursor position indicates how many spaces to the right of the starting x
and y location your cursor should be drawn at.

- If it is indicated that your cursor is moving backwards, then the space in
which the cursor was previously located will be automatically cleared.
*/
func drawCursor(layerEntry *memory.LayerEntryType, styleEntry memory.TuiStyleEntryType, textFieldAlias string, characterUnderCursor rune, xLocation int, yLocation int, cursorPosition int, isMovementBackwards bool) {
	attributeEntry := memory.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.CursorForegroundColor
	attributeEntry.BackgroundColor = styleEntry.CursorBackgroundColor
	attributeEntry.CellType = constants.CellTypeTextField
	attributeEntry.CellControlAlias = textFieldAlias
	attributeEntry.CellControlId = cursorPosition
	var arrayOfRunes []rune
	arrayOfRunes = append(arrayOfRunes, characterUnderCursor)
	printLayer(layerEntry, attributeEntry, xLocation+cursorPosition, yLocation, arrayOfRunes)
	if isMovementBackwards {
		printLayer(layerEntry, attributeEntry, xLocation+cursorPosition+1, yLocation, arrayOfRunes)
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

- Width indicates how large the visible area of your input string should be.

- String position indicates the location in your string where printing should
start. If the remainder of your string is too long for the specified width,
then only the visible portion will be displayed.
*/
func drawInputString(layerEntry *memory.LayerEntryType, styleEntry memory.TuiStyleEntryType, textFieldAlias string, xLocation int, yLocation int, width int, stringPosition int, inputValue []rune) {
	attributeEntry := memory.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.TextInputForegroundColor
	attributeEntry.BackgroundColor = styleEntry.TextInputBackgroundColor
	attributeEntry.CellType = constants.CellTypeTextField
	attributeEntry.CellControlAlias = textFieldAlias
	numberOfCharactersToSafelyPrint := stringformat.GetMaxCharactersThatFitInStringSize(inputValue[stringPosition:], width)
	textFieldEntry := memory.GetTextField(layerEntry.LayerAlias, textFieldAlias)
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType

	fillArea(layerEntry, attributeEntry," ", xLocation, yLocation, width, 1, 0)
	// Here we loop over each character to draw since we need to accommodate for unique
	// cell IDs (if required for mouse location detection).
	xLocationOffset := 0
	for currentRuneIndex := 0; currentRuneIndex < len(numberOfCharactersToSafelyPrint); currentRuneIndex ++ {
		if focusedControlType == constants.CellTypeTextField && focusedLayerAlias == layerEntry.LayerAlias && focusedControlAlias == textFieldAlias {
			if stringPosition+currentRuneIndex == textFieldEntry.CursorPosition {
				attributeEntry.ForegroundColor = styleEntry.CursorForegroundColor
				attributeEntry.BackgroundColor = styleEntry.CursorBackgroundColor
			} else {
				attributeEntry.ForegroundColor = styleEntry.TextInputForegroundColor
				attributeEntry.BackgroundColor = styleEntry.TextInputBackgroundColor
			}
		}
		attributeEntry.CellControlId = stringPosition + currentRuneIndex
		printLayer(layerEntry, attributeEntry, xLocation + xLocationOffset, yLocation, []rune{inputValue[stringPosition + currentRuneIndex]})
		xLocationOffset++
		if stringformat.IsRuneCharacterWide(inputValue[stringPosition + currentRuneIndex]) {
			xLocationOffset++
			printLayer(layerEntry, attributeEntry, xLocation + xLocationOffset, yLocation, []rune{' '})
		}
	}
}

func updateCursor2(textFieldEntry *memory.TextFieldEntryType, xLocation int) {
	textFieldEntry.CursorPosition = xLocation
	// If our cursor xLocation was jumped (due to NullCellControlId) or placed in an invalid xLocation spot greater than the length of our text line.
	// Move it to the end of the line.
	if textFieldEntry.CursorPosition == constants.NullCellControlId || textFieldEntry.CursorPosition > len(textFieldEntry.CurrentValue) - 1 {
		textFieldEntry.CursorPosition = 0
	}
}

func updateTextfieldViewport2(textFieldEntry *memory.TextFieldEntryType) {
	// If cursor xLocation is lower than the viewport window
	if textFieldEntry.CursorPosition <= textFieldEntry.ViewportPosition {
		maxViewportWidthAvaliable := textFieldEntry.Width
		if textFieldEntry.CursorPosition - textFieldEntry.Width < 0 {
			maxViewportWidthAvaliable = textFieldEntry.CursorPosition
		}
		arrayOfRunesAvailableToPrint := textFieldEntry.CurrentValue[textFieldEntry.CursorPosition - maxViewportWidthAvaliable : textFieldEntry.CursorPosition]
		numberOfRunesThatFitStringSize := stringformat.GetMaxCharactersThatFitInStringSizeReverse(arrayOfRunesAvailableToPrint, textFieldEntry.Width)
		textFieldEntry.ViewportPosition = textFieldEntry.CursorPosition - numberOfRunesThatFitStringSize
	}
	// Figure out how much displayable space is in our current viewport window.
	arrayOfRunesAvailableToPrint := textFieldEntry.CurrentValue[textFieldEntry.ViewportPosition:]
	arrayOfRunesThatFitStringSize := stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvailableToPrint, textFieldEntry.Width)
	// If the cursor xLocation is equal or greater than the visible viewport window width.
	if textFieldEntry.CursorPosition >= textFieldEntry.ViewportPosition + len(arrayOfRunesThatFitStringSize) {
		// Then make the viewport xLocation equal to the visible viewport width behind it.
		maxViewportWidthAvaliable := textFieldEntry.Width
		if textFieldEntry.CursorPosition - textFieldEntry.Width < 0 {
			maxViewportWidthAvaliable = textFieldEntry.CursorPosition
		}
		arrayOfRunesAvailableToPrint = textFieldEntry.CurrentValue[textFieldEntry.CursorPosition - maxViewportWidthAvaliable : textFieldEntry.CursorPosition]
		numberOfRunesThatFitStringSize := stringformat.GetMaxCharactersThatFitInStringSizeReverse(arrayOfRunesAvailableToPrint, textFieldEntry.Width)
		logInfo(fmt.Sprintf("v: %d x: %d off: %d fit: %d, aval: %s", textFieldEntry.ViewportPosition, textFieldEntry.CursorPosition, maxViewportWidthAvaliable, numberOfRunesThatFitStringSize, string(arrayOfRunesAvailableToPrint)))
		textFieldEntry.ViewportPosition = textFieldEntry.CursorPosition - numberOfRunesThatFitStringSize + 1
	}
}

func insertCharacterAtPosition(textFieldEntry *memory.TextFieldEntryType, characterToInsert rune) {
	textAfterCursor := stringformat.GetRuneArrayCopy(textFieldEntry.CurrentValue[textFieldEntry.CursorPosition:])
	textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition], characterToInsert)
	textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue, textAfterCursor...)
}

func deleteCharacterAtPosition(textFieldEntry *memory.TextFieldEntryType) {
	if len(textFieldEntry.CurrentValue) != 1 {
		if textFieldEntry.CursorPosition != len(textFieldEntry.CurrentValue) - 1 {
			textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition], textFieldEntry.CurrentValue[textFieldEntry.CursorPosition+1:]...)
		}
	}
}

func backspaceCharacterAtPosition(textFieldEntry *memory.TextFieldEntryType) {
	if textFieldEntry.CursorPosition >= 0 {
		textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition], textFieldEntry.CurrentValue[textFieldEntry.CursorPosition+1:]...)
	}
}

func updateTextFieldCursor(textFieldEntry *memory.TextFieldEntryType) {
	if textFieldEntry.CursorPosition < 0 {
		textFieldEntry.CursorPosition = 0
	}
	if textFieldEntry.CursorPosition >= len(textFieldEntry.CurrentValue) {
		textFieldEntry.CursorPosition = len(textFieldEntry.CurrentValue) - 1
	}
}

func updateKeyboardEventTextField(keystroke []rune) bool {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := true
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	if focusedControlType != constants.CellTypeTextField {
		return false
	}
	textFieldEntry := memory.GetTextField(focusedLayerAlias, focusedControlAlias)
	if len(keystroke) == 1 { // If a regular char is entered.
		if len(textFieldEntry.CurrentValue) < textFieldEntry.MaxLengthAllowed {
			logInfo(fmt.Sprintf("cur: %d view: %d", textFieldEntry.CursorPosition, textFieldEntry.ViewportPosition))
			insertCharacterAtPosition(textFieldEntry, keystroke[0])
			textFieldEntry.CursorPosition++
			updateTextFieldCursor(textFieldEntry)
			updateTextfieldViewport2(textFieldEntry)
			isScreenUpdateRequired = true
		}
	}
	if keystrokeAsString == "delete" {
		deleteCharacterAtPosition(textFieldEntry)
		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "home" {
		textFieldEntry.CursorPosition = 0
		textFieldEntry.ViewportPosition = 0
		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "end" {
		textFieldEntry.CursorPosition = len(textFieldEntry.CurrentValue) - 1
		updateTextfieldViewport2(textFieldEntry)
		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "backspace" || keystrokeAsString == "backspace2" {
		textFieldEntry.CursorPosition--
		backspaceCharacterAtPosition(textFieldEntry)
		updateTextFieldCursor(textFieldEntry)
		updateTextfieldViewport2(textFieldEntry)

		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "left" {
		textFieldEntry.CursorPosition--
		updateTextFieldCursor(textFieldEntry)
		updateTextfieldViewport2(textFieldEntry)

		isScreenUpdateRequired = true
	}
	if keystrokeAsString == "right" {
		textFieldEntry.CursorPosition++
		updateTextFieldCursor(textFieldEntry)
		updateTextfieldViewport2(textFieldEntry)
		isScreenUpdateRequired = true
	}
	return isScreenUpdateRequired
}

func updateMouseEventTextField() bool {
	isScreenUpdateRequired := false
	var characterEntry memory.CharacterEntryType
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	if buttonPressed != 0 {
		characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if characterEntry.AttributeEntry.CellType == constants.CellTypeTextField {
			textFieldEntry := memory.GetTextField(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			textFieldEntry.CursorPosition = characterEntry.AttributeEntry.CellControlId
			setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeTextField)
			isScreenUpdateRequired = true
		}
	}
	return isScreenUpdateRequired
}