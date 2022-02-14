package consolizer

import (
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
	textFieldEntry := memory.TextFieldMemory[shared.layerAlias][shared.textFieldAlias]
	value := textFieldEntry.CurrentValue
	return value
}

/*
SetLocation allows you to set the current location of your text field.
*/
func (shared *textFieldInstanceType) SetLocation(xLocation int, yLocation int) {
	validatorTextField(shared.layerAlias, shared.textFieldAlias)
	validateLayerLocationByLayerAlias(shared.layerAlias, xLocation, yLocation)
	textFieldEntry := memory.TextFieldMemory[shared.layerAlias][shared.textFieldAlias]
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
	memory.AddTextField(layerAlias, textFieldAlias, styleEntry, xLocation, yLocation, width, maxLengthAllowed, IsPasswordProtected, defaultValue)
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
	for currentKey := range memory.TextFieldMemory[layerAlias] {
		textFieldEntry := memory.TextFieldMemory[layerAlias][currentKey]
		drawInputString(&layerEntry, textFieldEntry.StyleEntry, currentKey, textFieldEntry.XLocation, textFieldEntry.YLocation, textFieldEntry.Width, textFieldEntry.ViewportPosition, textFieldEntry.CurrentValue)
		cursorPosition := textFieldEntry.ViewportPosition + textFieldEntry.CursorPosition
		runeSlice := []rune(textFieldEntry.CurrentValue)
		characterUnderCursor := ' '
		if cursorPosition < len(runeSlice) { // Protect against empty strings
			characterUnderCursor = runeSlice[cursorPosition]
		}
		if eventStateMemory.focusedControlType == constants.CellTypeTextField && eventStateMemory.focusedLayerAlias == layerAlias && eventStateMemory.focusedControlAlias == currentKey {
			drawCursor(&layerEntry, textFieldEntry.StyleEntry, currentKey, characterUnderCursor, textFieldEntry.XLocation, textFieldEntry.YLocation, textFieldEntry.CursorPosition, false)
		}
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
func drawInputString(layerEntry *memory.LayerEntryType, styleEntry memory.TuiStyleEntryType, textFieldAlias string, xLocation int, yLocation int, width int, stringPosition int, inputString string) {
	attributeEntry := memory.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.TextInputForegroundColor
	attributeEntry.BackgroundColor = styleEntry.TextInputBackgroundColor
	attributeEntry.CellType = constants.CellTypeTextField
	attributeEntry.CellControlAlias = textFieldAlias
	runeSlice := []rune(inputString)
	var safeSubstring string
	if stringPosition+width <= len(inputString) {
		safeSubstring = string(runeSlice[stringPosition : stringPosition+width])
	} else {
		safeSubstring = string(runeSlice[stringPosition : stringPosition+len(inputString)-stringPosition])
	}
	arrayOfRunes := stringformat.GetRunesFromString(safeSubstring)
	// Here we loop over each character to draw since we need to accommodate for unique
	// cell IDs (if required for mouse location detection).
	for currentRuneIndex := 0; currentRuneIndex < len(arrayOfRunes); currentRuneIndex ++ {
		attributeEntry.CellControlId = currentRuneIndex
		printLayer(layerEntry, attributeEntry, xLocation + currentRuneIndex, yLocation, []rune{arrayOfRunes[currentRuneIndex]})
	}
}

func updateKeyboardEventTextField(keystroke string) bool {
	isScreenUpdateRequired := true
	if eventStateMemory.focusedControlType != constants.CellTypeTextField {
		return false
	}
	textFieldEntry := memory.GetTextField(eventStateMemory.focusedLayerAlias, eventStateMemory.focusedControlAlias)
	viewportPosition := textFieldEntry.ViewportPosition
	cursorPosition := textFieldEntry.CursorPosition
	currentValue := textFieldEntry.CurrentValue
	viewportWidth := textFieldEntry.Width
	maxLengthAllowed := textFieldEntry.MaxLengthAllowed
	if len(keystroke) == 1 { // If a regular char is entered.
		if len(currentValue) < maxLengthAllowed {
			currentValue = currentValue[:viewportPosition+cursorPosition] + keystroke + currentValue[viewportPosition+cursorPosition:]
			if cursorPosition < viewportWidth {
				cursorPosition++
			} else {
				viewportPosition++
			}
			isScreenUpdateRequired = true
		}
	}
	if keystroke == "delete" {
		if currentValue != "" {
			// Protect if nothing else to delete left of string
			if viewportPosition+cursorPosition+1 <= len(currentValue) {
				currentValue = currentValue[:viewportPosition+cursorPosition] + currentValue[viewportPosition+cursorPosition+1:]
				if viewportPosition+cursorPosition == len(currentValue) {
					cursorPosition--
					if cursorPosition < 0 {
						cursorPosition = 0
					}
				}
				isScreenUpdateRequired = true
			}
		}
	}
	if keystroke == "home" {
		cursorPosition = 0
		viewportPosition = 0
		isScreenUpdateRequired = true
	}
	if keystroke == "end" {
		// If your current viewport shows the end of the input string, just move the cursor to the end of the string.
		if viewportPosition > len(currentValue)- viewportWidth {
			cursorPosition = len(currentValue) - viewportPosition
		} else {
			// Otherwise advance viewport to end of input string.
			viewportPosition = len(currentValue) - viewportWidth
			if viewportPosition < 0 {
				// If input string is smaller than even one viewport block, just set cursor to end.
				viewportPosition = 0
				cursorPosition = len(currentValue)
			} else {
				// Otherwise place cursor at end of viewport / string
				cursorPosition = viewportWidth
			}
		}
		isScreenUpdateRequired = true
	}
	if keystroke == "backspace" || keystroke == "backspace2" {
		if currentValue == "" {
			return false
		}
		// Protect if nothing else to delete left of string
		if viewportPosition + cursorPosition - 1 >= 0 {
			currentValue = currentValue[:viewportPosition+cursorPosition-1] + currentValue[viewportPosition+cursorPosition:]
			cursorPosition--
			if cursorPosition < 1 {
				if len(currentValue) < viewportWidth {
					cursorPosition = viewportPosition + cursorPosition
					viewportPosition = 0
				} else {
					if viewportPosition != 0 { // If your not at the start of an input string
						if cursorPosition == 0 {
							viewportPosition = viewportPosition - viewportWidth + 1
						} else {
							viewportPosition = viewportPosition - viewportWidth
						}
						if viewportPosition < 0 {
							viewportPosition = 0
						}
						cursorPosition = viewportWidth - 1
					}
				}
			}
			isScreenUpdateRequired = true
		}
	}
	if keystroke == "left" {
		cursorPosition--
		if cursorPosition < 0 {
			if viewportPosition == 0 {
				cursorPosition = 0
			} else {
				viewportPosition =- viewportWidth
				if viewportPosition < 0 {
					viewportPosition = 0
				}
				cursorPosition = viewportWidth
			}
		}
		isScreenUpdateRequired = true
	}
	if keystroke == "right" {
		cursorPosition++
		if viewportPosition + cursorPosition > len(currentValue){
			cursorPosition--
		} else {
			if cursorPosition >= viewportWidth {
				viewportPosition++
				cursorPosition = viewportWidth - 1
			}
		}
		isScreenUpdateRequired = true
	}
	if isScreenUpdateRequired {
		textFieldEntry.ViewportPosition = viewportPosition
		textFieldEntry.CursorPosition = cursorPosition
		textFieldEntry.CurrentValue = currentValue
		textFieldEntry.Width = viewportWidth
		textFieldEntry.MaxLengthAllowed = maxLengthAllowed
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