package consolizer

import (
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"strings"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

type textFieldInstanceType struct {
	BaseControlInstanceType
}

type textFieldType struct{}

var TextFields = memory.NewControlMemoryManager[types.TextFieldEntryType]()
var TextField textFieldType

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
GetValue allows you to get the current value of your text field. In addition, the following
information should be noted:

- If the text field is password protected, the actual value will be returned, not the masked characters.
- The returned value will be trimmed of any leading or trailing whitespace.
- If the text field does not exist, an empty string will be returned.
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
	textFieldEntry.TooltipAlias = stringformat.GetLastSortedUUID()

	// Create associated tooltip (always created but disabled by default)
	tooltipInstance := Tooltip.Add(layerAlias, textFieldEntry.TooltipAlias, "", styleEntry,
		textFieldEntry.XLocation, textFieldEntry.YLocation,
		textFieldEntry.Width, 1,
		textFieldEntry.XLocation, textFieldEntry.YLocation+1,
		textFieldEntry.Width, 3,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.SetEnabled(false)
	tooltipInstance.setParentControlAlias(textFieldAlias)
	// Use the generic memory manager to add the text field entry
	TextFields.Add(layerAlias, textFieldAlias, &textFieldEntry)

	var textFieldInstance textFieldInstanceType
	textFieldInstance.layerAlias = layerAlias
	textFieldInstance.controlAlias = textFieldAlias
	textFieldInstance.controlType = "textField"
	return textFieldInstance
}

/*
DeleteTextField allows you to delete a text field on a given layer. In addition, the following
information should be noted:

- If the text field does not exist, the request will be ignored.
- All memory associated with the text field will be freed.
- The text field will be removed from the tab index if it was added.
*/
func (shared *textFieldType) DeleteTextField(layerAlias string, textFieldAlias string) {
	validatorTextField(layerAlias, textFieldAlias)
	TextFields.Remove(layerAlias, textFieldAlias)
}

/*
DeleteAllTextFields allows you to delete all text fields on a given layer. In addition, the following
information should be noted:

- All text fields on the specified layer will be removed.
- All memory associated with the text fields will be freed.
- The text fields will be removed from the tab index if they were added.
*/
func (shared *textFieldType) DeleteAllTextFields(layerAlias string) {
	TextFields.RemoveAll(layerAlias)
}

/*
drawTextFieldOnLayer allows you to draw all text fields on a given text layer entry. In addition,
the following information should be noted:

- Text fields are drawn in the order they were created.
- Each text field is drawn with its own style and attributes.
- The cursor is drawn if the text field is currently focused.
*/
func (shared *textFieldType) drawTextFieldOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentTextFieldEntry := range TextFields.GetAllEntries(layerAlias) {
		textFieldEntry := currentTextFieldEntry
		shared.drawInputString(&layerEntry, textFieldEntry.StyleEntry, textFieldEntry.Alias, textFieldEntry.XLocation, textFieldEntry.YLocation, textFieldEntry.Width, textFieldEntry.ViewportPosition, textFieldEntry.CurrentValue)
	}
}

/*
drawInputString allows you to draw the input string for a text field. In addition, the following
information should be noted:

- The input string is drawn with the specified style and attributes.
- If the text field is password protected, characters are masked.
- The cursor is drawn if the text field is currently focused.
- Highlighted text is drawn with inverted colors if active.
*/
func (shared *textFieldType) drawInputString(layerEntry *types.LayerEntryType, styleEntry types.TuiStyleEntryType, textFieldAlias string, xLocation int, yLocation int, width int, stringPosition int, inputValue []rune) {
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.TextField.ForegroundColor
	attributeEntry.BackgroundColor = styleEntry.TextField.BackgroundColor
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
			attributeEntry.ForegroundColor = styleEntry.TextField.CursorForegroundColor
			attributeEntry.BackgroundColor = styleEntry.TextField.CursorBackgroundColor
		} else if isHighlighted {
			attributeEntry.ForegroundColor = styleEntry.TextField.HighlightForegroundColor
			attributeEntry.BackgroundColor = styleEntry.TextField.HighlightBackgroundColor
		} else {
			attributeEntry.ForegroundColor = styleEntry.TextField.ForegroundColor
			attributeEntry.BackgroundColor = styleEntry.TextField.BackgroundColor
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
location. In addition, the following information should be noted:

- Adjusts the viewport to ensure the cursor remains visible within the text field's width.
- Handles cases where the cursor moves outside the current viewport window.
- Automatically scrolls the text left or right to keep the cursor in view.
- Maintains proper text alignment and visibility when the cursor moves.
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
deleteCharacterAtPosition allows you to delete a character at the current cursor position. In addition,
the following information should be noted:

- If the cursor is at the end of the text, no character is deleted.
- The cursor position is updated after deletion.
- The text field's current value is updated to reflect the deletion.
*/
func (shared *textFieldType) deleteCharacterAtPosition(textFieldEntry *types.TextFieldEntryType) {
	if len(textFieldEntry.CurrentValue) != 1 {
		if textFieldEntry.CursorPosition != len(textFieldEntry.CurrentValue)-1 {
			textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition], textFieldEntry.CurrentValue[textFieldEntry.CursorPosition+1:]...)
		}
	}
}

/*
backspaceCharacterAtPosition allows you to backspace a character at the current cursor position. In addition,
the following information should be noted:

- If the cursor is at the beginning of the text, no character is deleted.
- The cursor position is updated after backspacing.
- The text field's current value is updated to reflect the deletion.
*/
func (shared *textFieldType) backspaceCharacterAtPosition(textFieldEntry *types.TextFieldEntryType) {
	if textFieldEntry.CursorPosition >= 0 {
		textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition], textFieldEntry.CurrentValue[textFieldEntry.CursorPosition+1:]...)
	}
}

/*
updateTextFieldCursor allows you to update a text field cursor's location. In addition, the following
information should be noted:

- Ensures the cursor position is within valid bounds.
- Updates the viewport position if necessary to keep the cursor visible.
- Handles cases where the cursor position is invalid or out of range.
*/
func (shared *textFieldType) updateTextFieldCursor(textFieldEntry *types.TextFieldEntryType) {
	if textFieldEntry.CursorPosition == constants.NullCellControlId || textFieldEntry.CursorPosition >= len(textFieldEntry.CurrentValue) {
		textFieldEntry.CursorPosition = len(textFieldEntry.CurrentValue) - 1
	}
	if textFieldEntry.CursorPosition < 0 {
		textFieldEntry.CursorPosition = 0
	}
}

func (shared *textFieldType) updateKeyboardEventManually(layerAlias string, textFieldAlias string, keystroke []rune) bool {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := false
	textFieldEntry := TextFields.Get(layerAlias, textFieldAlias)
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
		// Only process specific shift+key combinations for highlighting
		// This prevents shift+character combinations from being treated as special commands
		if strings.HasPrefix(keystrokeAsString, "shift+") && 
		   (strings.Contains(keystrokeAsString, "left") || 
		    strings.Contains(keystrokeAsString, "right") || 
		    strings.Contains(keystrokeAsString, "up") || 
		    strings.Contains(keystrokeAsString, "down") || 
		    strings.Contains(keystrokeAsString, "home") || 
		    strings.Contains(keystrokeAsString, "end") || 
		    strings.Contains(keystrokeAsString, "pgup") || 
		    strings.Contains(keystrokeAsString, "pgdn")) {
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
updateKeyboardEvent allows you to update the state of all text fields according to the current
keystroke event. In addition, the following information should be noted:

- Handles all keyboard input for text fields.
- Manages cursor movement, text insertion, and deletion.
- Handles special keys like Home, End, Delete, and Backspace.
- Manages text highlighting and selection.
- Returns true if a screen update is required.
*/
func (shared *textFieldType) updateKeyboardEvent(keystroke []rune) bool {
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	if focusedControlType != constants.CellTypeTextField || !TextFields.IsExists(focusedLayerAlias, focusedControlAlias) {
		return false
	}
	return shared.updateKeyboardEventManually(focusedLayerAlias, focusedControlAlias, keystroke)
}

/*
updateMouseEvent allows you to update the state of all text fields according to the current
mouse event. In addition, the following information should be noted:

- Handles mouse clicks and drags for text selection.
- Manages cursor positioning based on mouse clicks.
- Returns true if a screen update is required.
*/
func (shared *textFieldType) updateMouseEvent() bool {
	isScreenUpdateRequired := false
	var characterEntry types.CharacterEntryType
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	if buttonPressed != 0 && eventStateMemory.stateId != constants.EventStateDragAndDropScrollbar &&
		eventStateMemory.stateId != constants.EventStateDragAndDrop {
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

/*
updateKeyboardEventTextboxWithString allows you to update a text field with a string of characters. In addition,
the following information should be noted:

- Processes each character in the string as a separate keystroke.
- Useful for programmatically inserting text into a text field.
- Maintains all text field functionality like highlighting and cursor movement.
*/
func (shared *textFieldType) updateKeyboardEventTextboxWithString(keystroke string) {
	for _, currentCharacter := range keystroke {
		shared.updateKeyboardEvent([]rune{currentCharacter})
	}
}

/*
updateKeyboardEventTextboxWithCommands allows you to update a text field with a list of command strings. In addition,
the following information should be noted:

- Processes each command string as a separate keystroke.
- Useful for programmatically inserting text or executing commands in a text field.
- Maintains all text field functionality like highlighting and cursor movement.
*/
func (shared *textFieldType) updateKeyboardEventTextboxWithCommands(keystroke ...string) {
	for _, currentCommand := range keystroke {
		shared.updateKeyboardEvent([]rune(currentCommand))
	}
}

/*
SetValue allows you to set the current value of your text field. In addition, the following
information should be noted:

- If the text field is password protected, the value will be stored but displayed as masked characters.
- The value will be trimmed of any leading or trailing whitespace before being set.
- If the text field does not exist, the request will be ignored.
*/
func (shared *textFieldInstanceType) SetValue(value string) *textFieldInstanceType {
	if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry := TextFields.Get(shared.layerAlias, shared.controlAlias)
		textFieldEntry.CurrentValue = []rune(value)
	}
	return shared
}

/*
SetDefaultValue allows you to set the default value of your text field. In addition, the following
information should be noted:

- The default value will be used when the text field is first created or reset.
- If the text field is password protected, the default value will be stored but displayed as masked characters.
- The default value will be trimmed of any leading or trailing whitespace before being set.
*/
func (shared *textFieldInstanceType) SetDefaultValue(value string) *textFieldInstanceType {
	if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry := TextFields.Get(shared.layerAlias, shared.controlAlias)
		textFieldEntry.DefaultValue = value
	}
	return shared
}

/*
SetMaxLength allows you to set the maximum number of characters allowed in the text field. In addition,
the following information should be noted:

- If the current value exceeds the new maximum length, it will be truncated.
- Setting the maximum length to 0 or a negative number will remove the length restriction.
- The maximum length applies to the actual text, not including any masked characters for password fields.
*/
func (shared *textFieldInstanceType) SetMaxLength(length int) *textFieldInstanceType {
	if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry := TextFields.Get(shared.layerAlias, shared.controlAlias)
		textFieldEntry.MaxLengthAllowed = length
	}
	return shared
}

/*
SetPasswordProtected allows you to specify whether the text field should mask its contents. In addition,
the following information should be noted:

- When enabled, all characters will be displayed as asterisks (*) or another specified mask character.
- The actual value is still stored and can be retrieved using GetValue.
- This setting can be changed at any time, and the display will update accordingly.
*/
func (shared *textFieldInstanceType) SetPasswordProtected(isProtected bool) *textFieldInstanceType {
	if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry := TextFields.Get(shared.layerAlias, shared.controlAlias)
		textFieldEntry.IsPasswordProtected = isProtected
	}
	return shared
}

/*
SetCursorPosition allows you to set the position of the cursor within the text field. In addition,
the following information should be noted:

- The cursor position is zero-based, where 0 represents the start of the text.
- If the specified position is beyond the length of the text, the cursor will be placed at the end.
- The viewport will automatically adjust to keep the cursor visible.
*/
func (shared *textFieldInstanceType) SetCursorPosition(position int) *textFieldInstanceType {
	if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry := TextFields.Get(shared.layerAlias, shared.controlAlias)
		textFieldEntry.CursorPosition = position
	}
	return shared
}

/*
SetViewportPosition allows you to set the starting position of the visible portion of the text field.
In addition, the following information should be noted:

- The viewport position is zero-based, where 0 represents the start of the text.
- If the specified position is beyond the length of the text, the viewport will be set to the end.
- The cursor will remain visible within the viewport.
*/
func (shared *textFieldInstanceType) SetViewportPosition(position int) *textFieldInstanceType {
	if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry := TextFields.Get(shared.layerAlias, shared.controlAlias)
		textFieldEntry.ViewportPosition = position
	}
	return shared
}

/*
GetTooltip retrieves the tooltip associated with this text field and returns the text field instance
for method chaining. In addition, the following information should be noted:

- The tooltip is automatically created when the text field is added.
- Use the returned instance to continue chaining method calls.
- Follow the same pattern as other controls for consistency.
*/
func (shared *textFieldInstanceType) GetTooltip() *textFieldInstanceType {
	// No need to retrieve the tooltip, just return self for chaining
	return shared
}

// Add a helper method to set tooltip text
func (shared *textFieldInstanceType) SetTooltipText(text string) *textFieldInstanceType {
	if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
		textFieldEntry := TextFields.Get(shared.layerAlias, shared.controlAlias)
		var tooltipInstance TooltipInstanceType
		tooltipInstance.layerAlias = shared.layerAlias
		tooltipInstance.controlAlias = textFieldEntry.TooltipAlias
		tooltipInstance.SetTooltipValue(text)
	}
	return shared
}

// Add a helper method to enable/disable the tooltip
func (shared *textFieldInstanceType) EnableTooltip(enabled bool) *textFieldInstanceType {
	if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
		textFieldEntry := TextFields.Get(shared.layerAlias, shared.controlAlias)
		var tooltipInstance TooltipInstanceType
		tooltipInstance.layerAlias = shared.layerAlias
		tooltipInstance.controlAlias = textFieldEntry.TooltipAlias
		tooltipInstance.SetEnabled(enabled)
	}
	return shared
}
