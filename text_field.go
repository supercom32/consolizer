package consolizer

import (
	"strings"

	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/internal/memory"
	"supercom32.net/consolizer/internal/stringformat"
	"supercom32.net/consolizer/types"
)

type textFieldInstanceType struct {
	layerAlias   string
	controlAlias string
}

type textFieldType struct{}

var TextFields = memory.NewControlMemoryManager[types.TextFieldEntryType]()
var TextField textFieldType

// ============================================================================
// REGULAR ENTRY
// ============================================================================

func (shared *textFieldInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeTextField)
}

func (shared *textFieldInstanceType) GetFocus() string {
	if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		setFocusedControl(shared.layerAlias, shared.controlAlias, constants.CellTypeTextField)
	}
	return ""
}

func (shared *textFieldInstanceType) Delete() string {
	if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
		TextFields.Remove(shared.layerAlias, shared.controlAlias)
	}
	return ""
}

/*
GetValue allows you to get the current value of your text field with.
*/
func (shared *textFieldInstanceType) GetValue() string {
	if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry := TextFields.Get(shared.layerAlias, shared.controlAlias)
		value := textFieldEntry.CurrentValue
		return strings.TrimSpace(string(value))
	}
	return ""
}

/*
SetLocation allows you to set the current location of your text field.
*/
func (shared *textFieldInstanceType) SetLocation(xLocation int, yLocation int) {
	if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		validateLayerLocationByLayerAlias(shared.layerAlias, xLocation, yLocation)
		textFieldEntry := TextFields.Get(shared.layerAlias, shared.controlAlias)
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
	// This is required so that the cursor appears at the end of your text field.
	defaultValueWithTrailingBlank := defaultValue + " "
	textFieldEntry := types.NewTextFieldEntry()
	textFieldEntry.Alias = textFieldAlias
	textFieldEntry.StyleEntry = styleEntry
	textFieldEntry.XLocation = xLocation
	textFieldEntry.YLocation = yLocation
	textFieldEntry.Width = width
	textFieldEntry.MaxLengthAllowed = maxLengthAllowed
	textFieldEntry.IsPasswordProtected = IsPasswordProtected
	textFieldEntry.CurrentValue = []rune(defaultValueWithTrailingBlank)
	textFieldEntry.DefaultValue = defaultValueWithTrailingBlank
	textFieldEntry.IsEnabled = isEnabled
	// Use the generic memory manager to add the text field entry
	TextFields.Add(layerAlias, textFieldAlias, &textFieldEntry)

	var textFieldInstance textFieldInstanceType
	textFieldInstance.layerAlias = layerAlias
	textFieldInstance.controlAlias = textFieldAlias
	return textFieldInstance
}

/*
DeleteTextField allows you to delete a text field on a given layer.
*/
func (shared *textFieldType) DeleteTextField(layerAlias string, textFieldAlias string) {
	validatorTextField(layerAlias, textFieldAlias)
	TextFields.Remove(layerAlias, textFieldAlias)
}

func (shared *textFieldType) DeleteAllTextFields(layerAlias string) {
	TextFields.RemoveAll(layerAlias)
}

/*
drawTextFieldOnLayer allows you to draw all text fields on a given text layer
entry.
*/
func (shared *textFieldType) drawTextFieldOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentTextFieldEntry := range TextFields.GetAllEntries(layerAlias) {
		textFieldEntry := currentTextFieldEntry
		shared.drawInputString(&layerEntry, textFieldEntry.StyleEntry, textFieldEntry.Alias, textFieldEntry.XLocation, textFieldEntry.YLocation, textFieldEntry.Width, textFieldEntry.ViewportPosition, textFieldEntry.CurrentValue)
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
	textFieldEntry := TextFields.Get(layerEntry.LayerAlias, textFieldAlias)
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	fillArea(layerEntry, attributeEntry, " ", xLocation, yLocation, width, 1, 0)
	// Here we loop over each character to draw since we need to accommodate for unique
	// cell IDs (if required for mouse location detection).
	xLocationOffset := 0
	for currentRuneIndex := 0; currentRuneIndex < len(numberOfCharactersToSafelyPrint); currentRuneIndex++ {
		absolutePosition := stringPosition + currentRuneIndex
		isFocused := focusedControlType == constants.CellTypeTextField && focusedLayerAlias == layerEntry.LayerAlias && focusedControlAlias == textFieldAlias

		// Handle highlighting in both directions
		isHighlighted := false
		if textFieldEntry.IsHighlightActive {
			start := textFieldEntry.HighlightStart
			end := textFieldEntry.HighlightEnd
			if start > end {
				start, end = end, start
			}
			isHighlighted = absolutePosition >= start && absolutePosition <= end
		}

		isCursor := isFocused && absolutePosition == textFieldEntry.CursorPosition

		// Set colors based on state
		if isCursor {
			attributeEntry.ForegroundColor = styleEntry.TextFieldCursorForegroundColor
			attributeEntry.BackgroundColor = styleEntry.TextFieldCursorBackgroundColor
		} else if isHighlighted {
			attributeEntry.ForegroundColor = styleEntry.HighlightForegroundColor
			attributeEntry.BackgroundColor = styleEntry.HighlightBackgroundColor
		} else {
			attributeEntry.ForegroundColor = styleEntry.TextFieldForegroundColor
			attributeEntry.BackgroundColor = styleEntry.TextFieldBackgroundColor
		}

		attributeEntry.CellControlId = absolutePosition
		if textFieldEntry.IsPasswordProtected {
			// If the field is password protected, then do not print the terminating ' ' character with an *.
			if xLocationOffset == len(textFieldEntry.CurrentValue)-1 {
				printLayer(layerEntry, attributeEntry, xLocation+xLocationOffset, yLocation, []rune{' '})
			} else {
				printLayer(layerEntry, attributeEntry, xLocation+xLocationOffset, yLocation, []rune{'*'})
			}
		} else {
			printLayer(layerEntry, attributeEntry, xLocation+xLocationOffset, yLocation, []rune{inputValue[absolutePosition]})
		}
		xLocationOffset++
		if stringformat.IsRuneCharacterWide(inputValue[absolutePosition]) {
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
		textFieldEntry.CursorPosition = len(textFieldEntry.CurrentValue) - 1
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
	isScreenUpdateRequired := false
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	if focusedControlType != constants.CellTypeTextField || !TextFields.IsExists(focusedLayerAlias, focusedControlAlias) {
		return false
	}
	textFieldEntry := TextFields.Get(focusedLayerAlias, focusedControlAlias)
	if !textFieldEntry.IsEnabled {
		return false
	}

	if IsShiftPressed() {
		if !textFieldEntry.IsHighlightModeToggled {
			// Start new highlight when toggling on
			textFieldEntry.IsHighlightModeToggled = true
			textFieldEntry.IsHighlightActive = true
			textFieldEntry.HighlightStart = textFieldEntry.CursorPosition
		}
	} else {
		textFieldEntry.IsHighlightModeToggled = false
	}

	switch keystrokeAsString {
	case "ctrl+a":
		// Select all text
		textFieldEntry.HighlightStart = 0
		textFieldEntry.HighlightEnd = len(textFieldEntry.CurrentValue) - 1
		textFieldEntry.IsHighlightActive = true
		isScreenUpdateRequired = true

	case "ctrl+c":
		// Copy highlighted text
		if textFieldEntry.IsHighlightActive {
			// TODO: Implement clipboard functionality
			// highlightedText := textFieldEntry.CurrentValue[textFieldEntry.HighlightStart:textFieldEntry.HighlightEnd+1]
		}

	case "ctrl+x":
		// Cut highlighted text
		if textFieldEntry.IsHighlightActive {
			// TODO: Implement clipboard functionality
			// highlightedText := textFieldEntry.CurrentValue[textFieldEntry.HighlightStart:textFieldEntry.HighlightEnd+1]
			start := textFieldEntry.HighlightStart
			end := textFieldEntry.HighlightEnd
			if start > end {
				start, end = end, start
			}
			// Include cursor position in the deletion
			if textFieldEntry.CursorPosition > end {
				end = textFieldEntry.CursorPosition
			} else if textFieldEntry.CursorPosition < start {
				start = textFieldEntry.CursorPosition
			}
			// Preserve the trailing blank character
			if end == len(textFieldEntry.CurrentValue)-1 {
				// If we're deleting up to the end, keep the trailing blank
				textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:start], ' ')
			} else {
				// Otherwise, delete the highlighted text
				textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:start], textFieldEntry.CurrentValue[end+1:]...)
			}
			textFieldEntry.CursorPosition = start
			textFieldEntry.IsHighlightActive = false
			isScreenUpdateRequired = true
		}

	case "delete", "shift+delete":
		if textFieldEntry.IsHighlightActive {
			// Delete highlighted text
			start := textFieldEntry.HighlightStart
			end := textFieldEntry.HighlightEnd
			if start > end {
				start, end = end, start
			}
			// Include cursor position in the deletion
			if textFieldEntry.CursorPosition > end {
				end = textFieldEntry.CursorPosition
			} else if textFieldEntry.CursorPosition < start {
				start = textFieldEntry.CursorPosition
			}
			// Preserve the trailing blank character
			if end == len(textFieldEntry.CurrentValue)-1 {
				// If we're deleting up to the end, keep the trailing blank
				textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:start], ' ')
			} else {
				// Otherwise, delete the highlighted text
				textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:start], textFieldEntry.CurrentValue[end+1:]...)
			}
			textFieldEntry.CursorPosition = start
			textFieldEntry.IsHighlightActive = false
		} else {
			shared.deleteCharacterAtPosition(textFieldEntry)
		}
		isScreenUpdateRequired = true

	case "home", "shift+home":
		if textFieldEntry.IsHighlightModeToggled == false {
			textFieldEntry.IsHighlightActive = false
		}
		textFieldEntry.CursorPosition = 0
		textFieldEntry.ViewportPosition = 0
		if textFieldEntry.IsHighlightActive {
			textFieldEntry.HighlightEnd = textFieldEntry.CursorPosition
		}
		isScreenUpdateRequired = true

	case "end", "shift+end":
		if textFieldEntry.IsHighlightModeToggled == false {
			textFieldEntry.IsHighlightActive = false
		}
		textFieldEntry.CursorPosition = len(textFieldEntry.CurrentValue) - 1
		shared.updateTextFieldViewport(textFieldEntry)
		if textFieldEntry.IsHighlightActive {
			textFieldEntry.HighlightEnd = textFieldEntry.CursorPosition
		}
		isScreenUpdateRequired = true

	case "backspace", "backspace2", "shift+backspace", "shift+backspace2":
		if textFieldEntry.IsHighlightActive {
			// Delete highlighted text
			start := textFieldEntry.HighlightStart
			end := textFieldEntry.HighlightEnd
			if start > end {
				start, end = end, start
			}
			// Include cursor position in the deletion
			if textFieldEntry.CursorPosition > end {
				end = textFieldEntry.CursorPosition
			} else if textFieldEntry.CursorPosition < start {
				start = textFieldEntry.CursorPosition
			}
			// Preserve the trailing blank character
			if end == len(textFieldEntry.CurrentValue)-1 {
				// If we're deleting up to the end, keep the trailing blank
				textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:start], ' ')
			} else {
				// Otherwise, delete the highlighted text
				textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:start], textFieldEntry.CurrentValue[end+1:]...)
			}
			textFieldEntry.CursorPosition = start
			textFieldEntry.IsHighlightActive = false
		} else {
			textFieldEntry.CursorPosition--
			shared.backspaceCharacterAtPosition(textFieldEntry)
		}
		shared.updateTextFieldCursor(textFieldEntry)
		shared.updateTextFieldViewport(textFieldEntry)
		isScreenUpdateRequired = true

	case "left", "shift+left":
		if textFieldEntry.IsHighlightModeToggled == false {
			textFieldEntry.IsHighlightActive = false
		}
		if textFieldEntry.IsHighlightActive {
			// When highlighting to the left, exclude the character under the cursor
			if textFieldEntry.CursorPosition > textFieldEntry.HighlightStart {
				textFieldEntry.HighlightEnd = textFieldEntry.CursorPosition - 2
			} else {
				textFieldEntry.HighlightEnd = textFieldEntry.CursorPosition
			}
		}
		textFieldEntry.CursorPosition--
		shared.updateTextFieldCursor(textFieldEntry)
		shared.updateTextFieldViewport(textFieldEntry)
		isScreenUpdateRequired = true

	case "right", "shift+right":
		if textFieldEntry.IsHighlightModeToggled == false {
			textFieldEntry.IsHighlightActive = false
		}
		if textFieldEntry.IsHighlightActive {
			textFieldEntry.HighlightEnd = textFieldEntry.CursorPosition
		}
		textFieldEntry.CursorPosition++
		shared.updateTextFieldCursor(textFieldEntry)
		shared.updateTextFieldViewport(textFieldEntry)
		isScreenUpdateRequired = true

	default:
		if textFieldEntry.IsHighlightModeToggled == false {
			textFieldEntry.IsHighlightActive = false
		}
		// Handle regular character input
		if len(keystroke) == 1 {
			// Check if character limit is under max length allowed
			if len(textFieldEntry.CurrentValue) < textFieldEntry.MaxLengthAllowed+1 {
				if textFieldEntry.IsHighlightActive {
					// Delete highlighted text before inserting new character
					start := textFieldEntry.HighlightStart
					end := textFieldEntry.HighlightEnd
					if start > end {
						start, end = end, start
					}
					// Include cursor position in the deletion
					if textFieldEntry.CursorPosition > end {
						end = textFieldEntry.CursorPosition
					} else if textFieldEntry.CursorPosition < start {
						start = textFieldEntry.CursorPosition
					}
					// Preserve the trailing blank character
					if end == len(textFieldEntry.CurrentValue)-1 {
						// If we're deleting up to the end, keep the trailing blank
						textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:start], ' ')
					} else {
						// Otherwise, delete the highlighted text
						textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:start], textFieldEntry.CurrentValue[end+1:]...)
					}
					textFieldEntry.CursorPosition = start
					textFieldEntry.IsHighlightActive = false
				}
				shared.insertCharacterAtPosition(textFieldEntry, keystroke[0])
				textFieldEntry.CursorPosition++
				shared.updateTextFieldCursor(textFieldEntry)
				shared.updateTextFieldViewport(textFieldEntry)
				isScreenUpdateRequired = true
			}
		}

		// Handle Shift+Arrow keys for highlighting
		if strings.HasPrefix(keystrokeAsString, "shift+") {
			if !textFieldEntry.IsHighlightActive {
				textFieldEntry.HighlightStart = textFieldEntry.CursorPosition
				textFieldEntry.IsHighlightActive = true
			}
			textFieldEntry.HighlightEnd = textFieldEntry.CursorPosition
			isScreenUpdateRequired = true
		}
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
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	if buttonPressed != 0 && eventStateMemory.stateId != constants.EventStateDragAndDropScrollbar {
		characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if characterEntry.AttributeEntry.CellType == constants.CellTypeTextField && TextFields.IsExists(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias) {
			textFieldEntry := TextFields.Get(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			if !textFieldEntry.IsEnabled {
				return isScreenUpdateRequired
			}
			textFieldEntry.CursorPosition = characterEntry.AttributeEntry.CellControlId
			shared.updateTextFieldCursor(textFieldEntry)
			setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeTextField)
			isScreenUpdateRequired = true
		}
	}

	// Handle mouse drag for text selection
	if eventStateMemory.stateId == constants.EventStateDragAndDrop && eventStateMemory.currentlyFocusedControl.controlType == constants.CellTypeTextField {
		characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if characterEntry.AttributeEntry.CellType == constants.CellTypeTextField && TextFields.IsExists(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias) {
			textFieldEntry := TextFields.Get(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
			if !textFieldEntry.IsEnabled {
				return isScreenUpdateRequired
			}
			if !textFieldEntry.IsHighlightActive {
				textFieldEntry.HighlightStart = textFieldEntry.CursorPosition
				textFieldEntry.IsHighlightActive = true
			}
			textFieldEntry.HighlightEnd = characterEntry.AttributeEntry.CellControlId
			isScreenUpdateRequired = true
		}
	}

	return isScreenUpdateRequired
}

func (shared *textFieldType) updateKeyboardEventTextboxWithString(keystroke string) {
	for _, currentCharacter := range keystroke {
		shared.updateKeyboardEventTextField([]rune{currentCharacter})
	}
}

func (shared *textFieldType) updateKeyboardEventTextboxWithCommands(keystroke ...string) {
	for _, currentCommand := range keystroke {
		shared.updateKeyboardEventTextField([]rune(currentCommand))
	}
}
