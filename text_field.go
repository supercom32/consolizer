package consolizer

import (
	"github.com/atotto/clipboard"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"strings"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

/*
TextFieldInstanceType is a structure which represents an instance of a text field control.
*/
type TextFieldInstanceType struct {
	BaseControlInstanceType
}

type textFieldType struct{}

var textFieldRegistry = memory.NewControlMemoryManager[types.TextFieldEntryType]()
var TextField textFieldType

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
GetValue is a method which gets the current value of your text field. In addition, the following should be noted:

- If the text field is password protected, the actual value will be returned, not the masked characters.

- The returned value will be trimmed of any leading or trailing whitespace.

- If the text field does not exist, an empty string will be returned.

Example:
    value := textField.GetValue()
*/
func (shared *TextFieldInstanceType) GetValue() string {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		value := textFieldEntry.CurrentValue
		if len(value) > 0 {
			value = value[:len(value)-1] // remove one character from the right
		}
		return string(value)
	}
	return ""
}

/*
SetLocation is a method which sets the current location of your text field.

Example:
    textField.SetLocation(10, 5)
*/
func (shared *TextFieldInstanceType) SetLocation(xLocation int, yLocation int) {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		validateLayerLocationByLayerAlias(shared.layerAlias, xLocation, yLocation)
		textFieldEntry.XLocation = xLocation
		textFieldEntry.YLocation = yLocation
	}
}

/*
Add is a method which adds a text field to a given layer. Once called, a text field instance is returned which will allow
you to read or manipulate properties of your text field. In addition, the following should be noted:

- If the location specified for the text field falls outside the range of the text layer, then only the
  visible portion of your text field will be rendered.

- If the max length of your text field is less than or equal to 0, a panic will be generated to fail as fast
  as possible.

- Password protection will echo back '*' characters to the terminal instead of the actual characters entered.

- Specifying a default value will simply pre-populate the text field with the value specified.

- If the cursor position moves outside the visible display area of the field, then the entire text field will
  shift to keep it in view.

Example:
    textFieldInstance := TextField.Add("Layer1", "Text1", style, 0, 0, 20, 100, false, "", true)
*/
func (shared *textFieldType) Add(layerAlias string, textFieldAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, maxLengthAllowed int, IsPasswordProtected bool, defaultValue string, isEnabled bool) TextFieldInstanceType {
	defer lockUIState()()
	return shared.add(layerAlias, textFieldAlias, styleEntry, xLocation, yLocation, width, maxLengthAllowed, IsPasswordProtected, defaultValue, isEnabled)
}

/*
add is a method which allows you to do exactly what Add does, for a caller that already holds uiStateMutex, such as
another entry point of the package that has taken the lock itself.

Example:

	TextField.add(layerAlias, textFieldAlias, styleEntry, xLocation, yLocation, width, maxLengthAllowed, IsPasswordProtected, defaultValue, isEnabled)
*/
func (shared *textFieldType) add(layerAlias string, textFieldAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, maxLengthAllowed int, IsPasswordProtected bool, defaultValue string, isEnabled bool) TextFieldInstanceType {
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
	textFieldEntry.DefaultValue = defaultValue
	textFieldEntry.IsEnabled = isEnabled
	textFieldEntry.TooltipAlias = stringformat.GetLastSortedUUID()

	// Create associated tooltip (always created but disabled by default)
	tooltipInstance := Tooltip.add(layerAlias, textFieldEntry.TooltipAlias, "", styleEntry,
		textFieldEntry.XLocation, textFieldEntry.YLocation,
		textFieldEntry.Width, 1,
		textFieldEntry.XLocation, textFieldEntry.YLocation+1,
		textFieldEntry.Width, 3,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.setEnabled(false)
	tooltipInstance.setParentControlAlias(textFieldAlias)
	// Use the generic memory manager to add the text field entry
	textFieldRegistry.Add(layerAlias, textFieldAlias, &textFieldEntry)

	var textFieldInstance TextFieldInstanceType
	textFieldInstance.layerAlias = layerAlias
	textFieldInstance.controlAlias = textFieldAlias
	textFieldInstance.controlType = constants.TYPE_TEXTFIELD
	return textFieldInstance
}

/*
Delete is a method which deletes a text field on a given layer. In addition, the following should be noted:

- If the text field does not exist, the request will be ignored.

- All memory associated with the text field will be freed.

- The text field will be removed from the tab index if it was added.

Example:
    TextField.Delete("Layer1", "Text1")
*/
func (shared *textFieldType) Delete(layerAlias string, textFieldAlias string) {
	defer lockUIState()()
	validatorTextField(layerAlias, textFieldAlias)
	textFieldRegistry.Remove(layerAlias, textFieldAlias)
	clearStaleControlReferences(layerAlias, textFieldAlias, constants.CellTypeTextField)
}

/*
DeleteAll is a method which deletes all text fields on a given layer. In addition, the following should be noted:

- All text fields on the specified layer will be removed.

- All memory associated with the text fields will be freed.

- The text fields will be removed from the tab index if they were added.

Example:
    TextField.DeleteAll("Layer1")
*/
func (shared *textFieldType) DeleteAll(layerAlias string) {
	defer lockUIState()()
	removeAllControlsAndClearCells(textFieldRegistry, layerAlias, constants.CellTypeTextField, func(entry *types.TextFieldEntryType) string { return entry.Alias })
}

/*
drawOnLayer is a method which draws all text fields on a given text layer entry. In addition, the following should be
noted:

- Text fields are drawn in the order they were created.

- Each text field is drawn with its own style and attributes.

- The cursor is drawn if the text field is currently focused.

Example:
    TextField.drawOnLayer(layerEntry)
*/
func (shared *textFieldType) drawOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentTextFieldEntry := range textFieldRegistry.GetAllEntries(layerAlias) {
		textFieldEntry := currentTextFieldEntry
		shared.drawInputString(&layerEntry, textFieldEntry.StyleEntry, textFieldEntry.Alias, textFieldEntry.XLocation, textFieldEntry.YLocation, textFieldEntry.Width, textFieldEntry.ViewportPosition, textFieldEntry.CurrentValue)
	}
}

/*
drawInputString is a method which allows you to draw the value of a text field, starting from a given rune position,
in the field's style. Password protected values are masked, the cursor is drawn while the field has focus, and
highlighted text is drawn in the highlight colours. In addition, the following should be noted:

  - A starting position outside the value is clamped into it, so drawing cannot panic on the event goroutine.

Example:

	TextField.drawInputString(&layerEntry, style, "Text1", 0, 0, 20, 0, runes)
*/
func (shared *textFieldType) drawInputString(layerEntry *types.LayerEntryType, styleEntry types.TuiStyleEntryType, textFieldAlias string, xLocation int, yLocation int, width int, stringPosition int, inputValue []rune) {
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.TextField.ForegroundColor
	attributeEntry.BackgroundColor = styleEntry.TextField.BackgroundColor
	attributeEntry.CellType = constants.CellTypeTextField
	attributeEntry.CellControlAlias = textFieldAlias
	// Take a rune-boundary prefix that fits the field width in COLUMNS so a trailing wide rune is
	// never split across the right edge.
	stringPosition = getClampedIndex(stringPosition, 0, len(inputValue)-1)
	runesToDraw, _ := stringformat.GetRunesThatFitInColumnCountFromStart(inputValue[stringPosition:], width)
	textFieldEntry := textFieldRegistry.Get(layerEntry.LayerAlias, textFieldAlias)
	focusedControl := getFocusedControl()
	focusedLayerAlias := focusedControl.layerAlias
	focusedControlAlias := focusedControl.controlAlias
	focusedControlType := focusedControl.controlType
	isFocused := focusedControlType == constants.CellTypeTextField && focusedLayerAlias == layerEntry.LayerAlias && focusedControlAlias == textFieldAlias
	lastRuneIndex := len(textFieldEntry.CurrentValue) - 1
	fillArea(layerEntry, attributeEntry, " ", xLocation, yLocation, width, 1, 0)
	// Here we loop over each character to draw since we need to accommodate for unique
	// cell IDs (if required for mouse location detection). The column offset advances by the
	// printed width of each rune so wide runes stay aligned with the underlying cells.
	xLocationOffset := 0
	for currentRuneIndex := 0; currentRuneIndex < len(runesToDraw); currentRuneIndex++ {
		absolutePosition := stringPosition + currentRuneIndex
		currentCharacter := inputValue[absolutePosition]

		// Handle highlighting in both directions
		isHighlighted := false
		if textFieldEntry.IsHighlightActive {
			highlightStart := textFieldEntry.HighlightStart
			highlightEnd := textFieldEntry.HighlightEnd
			if highlightStart > highlightEnd {
				highlightStart, highlightEnd = highlightEnd, highlightStart
			}
			isHighlighted = absolutePosition >= highlightStart && absolutePosition <= highlightEnd
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
		characterToDraw := currentCharacter
		isMaskedWideRune := false
		if textFieldEntry.IsPasswordProtected && absolutePosition != lastRuneIndex {
			// Every source rune, wide or narrow, is masked by exactly one '*'; the trailing
			// sentinel blank is never masked.
			characterToDraw = '*'
			isMaskedWideRune = stringformat.GetWidthOfRuneWhenPrinted(currentCharacter) == 2
		}
		putRune(layerEntry.CharacterMemory, xLocation+xLocationOffset, yLocation, characterToDraw, attributeEntry, layerEntry.Width, layerEntry.Height)
		if isMaskedWideRune {
			// The narrow '*' that masks a wide rune leaves its second column empty, so pad it
			// with a blank carrying the same attributes to keep the masked field column-aligned
			// with the unmasked one.
			putRune(layerEntry.CharacterMemory, xLocation+xLocationOffset+1, yLocation, ' ', attributeEntry, layerEntry.Width, layerEntry.Height)
		}
		xLocationOffset += stringformat.GetWidthOfRuneWhenPrinted(currentCharacter)
	}
}

/*
updateViewport is a method which updates the current viewport based on the current text and cursor location. The
viewport is scrolled the minimum amount needed to keep the cursor column inside the field width, measuring every
distance in printed columns so a run of wide runes cannot push the cursor off the visible edge. In addition, the
following should be noted:

- ViewportPosition is a rune index; it is only ever moved to a rune boundary so a wide rune is never half shown.

- When the cursor sits left of the viewport the viewport snaps to the cursor; when it sits at or past the right
  edge the viewport is pulled forward so the cursor and everything that fits behind it in Width columns is shown.

- ViewportPosition is clamped into the closed range zero through len(CurrentValue) minus one.

Example:
    TextField.updateViewport(textFieldEntry)
*/
func (shared *textFieldType) updateViewport(textFieldEntry *types.TextFieldEntryType) {
	lastRuneIndex := len(textFieldEntry.CurrentValue) - 1
	if textFieldEntry.ViewportPosition < 0 {
		textFieldEntry.ViewportPosition = 0
	}
	if textFieldEntry.ViewportPosition > lastRuneIndex {
		textFieldEntry.ViewportPosition = lastRuneIndex
	}
	visibleRunes := textFieldEntry.CurrentValue[textFieldEntry.ViewportPosition:]
	visiblePrefix, _ := stringformat.GetRunesThatFitInColumnCountFromStart(visibleRunes, textFieldEntry.Width)
	cursorColumn := stringformat.GetColumnIndexBasedOnRuneIndex(visibleRunes, textFieldEntry.CursorPosition-textFieldEntry.ViewportPosition)
	if textFieldEntry.CursorPosition < textFieldEntry.ViewportPosition {
		textFieldEntry.ViewportPosition = textFieldEntry.CursorPosition
	} else if textFieldEntry.CursorPosition >= textFieldEntry.ViewportPosition+len(visiblePrefix) || cursorColumn >= textFieldEntry.Width {
		trailingSuffix, _ := stringformat.GetRunesThatFitInColumnCountFromEnd(textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition+1], textFieldEntry.Width)
		textFieldEntry.ViewportPosition = textFieldEntry.CursorPosition + 1 - len(trailingSuffix)
	}
	if textFieldEntry.ViewportPosition < 0 {
		textFieldEntry.ViewportPosition = 0
	}
	if textFieldEntry.ViewportPosition > lastRuneIndex {
		textFieldEntry.ViewportPosition = lastRuneIndex
	}
}

/*
normalizeIndexes is a method which allows you to bring a text field's cursor, viewport and highlight positions back in
line with its current value, so that none of them can be used to index outside it. Each position is clamped between
zero and the index of the value's trailing sentinel blank, and an empty value is given that sentinel first. It is
called before a keystroke is processed and whenever one of the positions is set directly.

Example:

	TextField.normalizeIndexes(textFieldEntry)
*/
func (shared *textFieldType) normalizeIndexes(textFieldEntry *types.TextFieldEntryType) {
	if len(textFieldEntry.CurrentValue) == 0 {
		textFieldEntry.CurrentValue = []rune{' '}
	}
	lastRuneIndex := len(textFieldEntry.CurrentValue) - 1
	textFieldEntry.CursorPosition = getClampedIndex(textFieldEntry.CursorPosition, 0, lastRuneIndex)
	textFieldEntry.ViewportPosition = getClampedIndex(textFieldEntry.ViewportPosition, 0, lastRuneIndex)
	textFieldEntry.HighlightStart = getClampedIndex(textFieldEntry.HighlightStart, 0, lastRuneIndex)
	textFieldEntry.HighlightEnd = getClampedIndex(textFieldEntry.HighlightEnd, 0, lastRuneIndex)
}

/*
insertCharacterAtPosition is a method which inserts a character into a given text field. The location to insert is determined automatically by the current cursor position.

Example:
    TextField.insertCharacterAtPosition(textFieldEntry, 'a')
*/
func (shared *textFieldType) insertCharacterAtPosition(textFieldEntry *types.TextFieldEntryType, characterToInsert rune) {
	textAfterCursor := stringformat.GetRuneArrayCopy(textFieldEntry.CurrentValue[textFieldEntry.CursorPosition:])
	textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition], characterToInsert)
	textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue, textAfterCursor...)
}

/*
deleteCharacterAtPosition is a method which deletes a character at the current cursor position. In addition, the
following should be noted:

- If the cursor is at the end of the text, no character is deleted.

- The cursor position is updated after deletion.

- The text field's current value is updated to reflect the deletion.

Example:
    TextField.deleteCharacterAtPosition(textFieldEntry)
*/
func (shared *textFieldType) deleteCharacterAtPosition(textFieldEntry *types.TextFieldEntryType) {
	if len(textFieldEntry.CurrentValue) != 1 {
		if textFieldEntry.CursorPosition != len(textFieldEntry.CurrentValue)-1 {
			textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition], textFieldEntry.CurrentValue[textFieldEntry.CursorPosition+1:]...)
		}
	}
}

/*
backspaceCharacterAtPosition is a method which backspaces a character at the current cursor position. In addition, the
following should be noted:

- If the cursor is at the beginning of the text, no character is deleted.

- The cursor position is updated after backspacing.

- The text field's current value is updated to reflect the deletion.

Example:
    TextField.backspaceCharacterAtPosition(textFieldEntry)
*/
func (shared *textFieldType) backspaceCharacterAtPosition(textFieldEntry *types.TextFieldEntryType) {
	if textFieldEntry.CursorPosition >= 0 {
		textFieldEntry.CurrentValue = append(textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition], textFieldEntry.CurrentValue[textFieldEntry.CursorPosition+1:]...)
	}
}

/*
updateCursor is a method which updates a text field cursor's location. In addition, the following should be noted:

- Ensures the cursor position is within valid bounds.

- Updates the viewport position if necessary to keep the cursor visible.

- Handles cases where the cursor position is invalid or out of range.

Example:
    TextField.updateCursor(textFieldEntry)
*/
func (shared *textFieldType) updateCursor(textFieldEntry *types.TextFieldEntryType) {
	if textFieldEntry.CursorPosition == constants.NullCellControlId || textFieldEntry.CursorPosition >= len(textFieldEntry.CurrentValue) {
		textFieldEntry.CursorPosition = len(textFieldEntry.CurrentValue) - 1
	}
	if textFieldEntry.CursorPosition < 0 {
		textFieldEntry.CursorPosition = 0
	}
}

/*
updateKeyboardEventManually is a method which allows you to update the state of a text field according to a keystroke,
editing its value, moving its cursor, or changing its highlight. It returns whether the screen needs updating and
whether the keystroke was consumed. In addition, the following should be noted:

  - The cursor, viewport and highlight positions are clamped to the current value before the keystroke is applied,
    so positions left behind by a shorter value, or set directly, cannot be used to index outside it.

Example:

	isUpdateRequired, isConsumed := TextField.updateKeyboardEventManually("Layer1", "Text1", []rune("a"))
*/
func (shared *textFieldType) updateKeyboardEventManually(layerAlias string, textFieldAlias string, keystroke []rune) (bool, bool) {
	keystrokeAsString := string(keystroke)
	isScreenUpdateRequired := false
	isKeystrokeConsumed := false
	textFieldEntry, isFound := textFieldRegistry.Lookup(layerAlias, textFieldAlias)
	if !isFound || !textFieldEntry.IsEnabled {
		return false, false
	}
	shared.normalizeIndexes(textFieldEntry)

	// Windows Quirk: On Linux, the Shift key is only reported as "pressed" when used with non-character keys
	// (e.g., Shift+Delete). For regular character input like capital letters or symbols (e.g., Shift+A or Shift+;),
	// only the resulting character is passed, not the Shift key state.
	//
	// On Windows, however, the Shift key is reported as pressed with every key event, including single characters.
	// To prevent unintended behaviors like text highlighting, we ignore the Shift state on Windows when the key event
	// corresponds to a single printable character.
	if IsShiftPressed() && len(keystroke) != 1 {
		if !textFieldEntry.IsHighlightModeToggled {
			// Start new highlight when toggling on
			textFieldEntry.IsHighlightModeToggled = true
			textFieldEntry.IsHighlightActive = true
			textFieldEntry.HighlightStart = textFieldEntry.CursorPosition
			textFieldEntry.HighlightEnd = textFieldEntry.CursorPosition
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
		isKeystrokeConsumed = true

	case "ctrl+c", "ctrl+insert": // Copy
		// Copy highlighted text
		if textFieldEntry.IsHighlightActive {
			err := clipboard.WriteAll(shared.getSelectedText(textFieldEntry))
			if err != nil {
				panic(1)
				// Handle clipboard write error (could log it)
				// fmt.Println("Clipboard write error:", err)
			}
		}
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true

	case "ctrl+x": // Cut
		// Cut highlighted text
		if textFieldEntry.IsHighlightActive {
			start := textFieldEntry.HighlightStart
			end := textFieldEntry.HighlightEnd
			if start > end {
				start, end = end, start
			}

			// Ensure we don't exceed array bounds
			if end >= len(textFieldEntry.CurrentValue) {
				end = len(textFieldEntry.CurrentValue) - 1
			}
			if start < 0 {
				start = 0
			}

			// Copy to clipboard before deleting
			err := clipboard.WriteAll(shared.getSelectedText(textFieldEntry))
			if err != nil {
				// Handle clipboard write error (could log it)
				// fmt.Println("Clipboard write error:", err)
			}

			// Include the cursor position in the deletion
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
			isKeystrokeConsumed = true
		}

	case "ctrl+v", "shift+insert": // Paste
		// If there's highlighted text, delete it first
		if textFieldEntry.IsHighlightActive {
			start := textFieldEntry.HighlightStart
			end := textFieldEntry.HighlightEnd
			if start > end {
				start, end = end, start
			}

			// Include the cursor position in the deletion
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

		// Get text from clipboard
		clipboardText, err := clipboard.ReadAll()
		if err != nil {
			// Handle clipboard read error (could log it)
			// fmt.Println("Clipboard read error:", err)
		} else {
			// Insert clipboard text at cursor position
			// For text fields, we only insert the first line (no newlines)
			clipboardLine := strings.Split(clipboardText, "\n")[0]

			// Check if adding the clipboard text would exceed the max length
			if len(textFieldEntry.CurrentValue)+len(clipboardLine) <= textFieldEntry.MaxLengthAllowed+1 {
				// Insert the clipboard text at the cursor position
				newValue := append([]rune{}, textFieldEntry.CurrentValue[:textFieldEntry.CursorPosition]...)
				newValue = append(newValue, []rune(clipboardLine)...)
				newValue = append(newValue, textFieldEntry.CurrentValue[textFieldEntry.CursorPosition:]...)
				textFieldEntry.CurrentValue = newValue

				// Move cursor to after the inserted text
				textFieldEntry.CursorPosition += len(clipboardLine)
				shared.updateCursor(textFieldEntry)
				shared.updateViewport(textFieldEntry)
			}
		}
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true

	case "delete", "shift+delete":
		if textFieldEntry.IsHighlightActive {
			// Delete highlighted text
			start := textFieldEntry.HighlightStart
			end := textFieldEntry.HighlightEnd
			if start > end {
				start, end = end, start
			}
			// Include the cursor position in the deletion
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
		isKeystrokeConsumed = true

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
		isKeystrokeConsumed = true

	case "end", "shift+end":
		if textFieldEntry.IsHighlightModeToggled == false {
			textFieldEntry.IsHighlightActive = false
		}
		textFieldEntry.CursorPosition = len(textFieldEntry.CurrentValue) - 1
		shared.updateViewport(textFieldEntry)
		if textFieldEntry.IsHighlightActive {
			textFieldEntry.HighlightEnd = textFieldEntry.CursorPosition
		}
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true

	case "backspace", "backspace2", "shift+backspace", "shift+backspace2":
		if textFieldEntry.IsHighlightActive {
			// Delete highlighted text
			start := textFieldEntry.HighlightStart
			end := textFieldEntry.HighlightEnd
			if start > end {
				start, end = end, start
			}
			// Include the cursor position in the deletion
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
		shared.updateCursor(textFieldEntry)
		shared.updateViewport(textFieldEntry)
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true

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
		shared.updateCursor(textFieldEntry)
		shared.updateViewport(textFieldEntry)
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true

	case "right", "shift+right":
		if textFieldEntry.IsHighlightModeToggled == false {
			textFieldEntry.IsHighlightActive = false
		}
		if textFieldEntry.IsHighlightActive {
			textFieldEntry.HighlightEnd = textFieldEntry.CursorPosition
		}
		textFieldEntry.CursorPosition++
		shared.updateCursor(textFieldEntry)
		shared.updateViewport(textFieldEntry)
		isScreenUpdateRequired = true
		isKeystrokeConsumed = true

	default:
		// Handle regular character input
		if len(keystroke) == 1 {
			// Check if character limit is under max length allowed
			if len(textFieldEntry.CurrentValue) < textFieldEntry.MaxLengthAllowed+1 {
				if textFieldEntry.IsHighlightActive && !IsShiftPressed() {
					// Delete highlighted text before inserting a new character
					start := textFieldEntry.HighlightStart
					end := textFieldEntry.HighlightEnd
					if start > end {
						start, end = end, start
					}
					// Include the cursor position in the deletion
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
				shared.updateCursor(textFieldEntry)
				shared.updateViewport(textFieldEntry)
				isScreenUpdateRequired = true
				isKeystrokeConsumed = true
			}
		}
		if textFieldEntry.IsHighlightModeToggled == false {
			textFieldEntry.IsHighlightActive = false
		}
		// Handle Shift+Arrow keys for highlighting
		if strings.HasPrefix(keystrokeAsString, "shift+") {
			if !textFieldEntry.IsHighlightActive {
				textFieldEntry.HighlightStart = textFieldEntry.CursorPosition
				textFieldEntry.IsHighlightActive = true
			}
			textFieldEntry.HighlightEnd = textFieldEntry.CursorPosition
			isScreenUpdateRequired = true
			isKeystrokeConsumed = true
		}
	}
	return isScreenUpdateRequired, isKeystrokeConsumed
}

/*
updateKeyboardEvent is a method which updates the state of the currently focused text field according to the current
keystroke event. Only the currently focused text field will process the event. In addition, the following should be
noted:

  - If the focused field has an OnValueChanged hook set, the hook is invoked with the field's trimmed current value
    after the keystroke has been processed, and if it returns a different value, CurrentValue is replaced and the
    cursor is re-clamped before this method returns, so the correction is drawn in the same screen update as the
    keystroke.

  - The hook is application code, so it is called with uiStateMutex released and may call any exported function.
    The field is looked up again afterwards, and the correction is applied only if the field still exists and
    still holds the value the hook was given, so a hook can never overwrite a newer value with a stale correction.

Example:

	updateRequired, consumed := TextField.updateKeyboardEvent([]rune("a"))
*/
func (shared *textFieldType) updateKeyboardEvent(keystroke []rune) (bool, bool) {
	focusedControl := getFocusedControl()
	focusedLayerAlias := focusedControl.layerAlias
	focusedControlAlias := focusedControl.controlAlias
	focusedControlType := focusedControl.controlType
	textFieldEntry, isFound := textFieldRegistry.Lookup(focusedLayerAlias, focusedControlAlias)
	if focusedControlType != constants.CellTypeTextField || !isFound {
		return false, false
	}
	isScreenUpdateRequired, isKeystrokeConsumed := shared.updateKeyboardEventManually(focusedLayerAlias, focusedControlAlias, keystroke)
	onValueChanged := textFieldEntry.OnValueChanged
	if onValueChanged == nil {
		return isScreenUpdateRequired, isKeystrokeConsumed
	}
	currentValueAsString := shared.getTrimmedValue(textFieldEntry)
	var correctedValue string
	callWithUIStateUnlocked(func() { correctedValue = onValueChanged(currentValueAsString) })
	textFieldEntry, isFound = textFieldRegistry.Lookup(focusedLayerAlias, focusedControlAlias)
	if !isFound || correctedValue == currentValueAsString || shared.getTrimmedValue(textFieldEntry) != currentValueAsString {
		return isScreenUpdateRequired, isKeystrokeConsumed
	}
	textFieldEntry.CurrentValue = []rune(correctedValue + " ")
	shared.updateCursor(textFieldEntry)
	shared.updateViewport(textFieldEntry)
	return true, isKeystrokeConsumed
}

/*
getTrimmedValue is a method which allows you to obtain a text field's current value as a string without the trailing
blank sentinel rune its value is stored with.

Example:

	value := TextField.getTrimmedValue(textFieldEntry)
*/
func (shared *textFieldType) getTrimmedValue(textFieldEntry *types.TextFieldEntryType) string {
	currentValue := textFieldEntry.CurrentValue
	if len(currentValue) > 0 {
		currentValue = currentValue[:len(currentValue)-1]
	}
	return string(currentValue)
}

/*
updateMouseEvent is a method which updates the state of all text fields according to the current mouse event. In
addition, the following should be noted:

- Handles mouse clicks and drags for text selection and cursor positioning.

Example:
    updateRequired := TextField.updateMouseEvent()
*/
func (shared *textFieldType) updateMouseEvent() bool {
	isScreenUpdateRequired := false
	var characterEntry types.CharacterEntryType
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	if buttonPressed != 0 && getEventStateId() != constants.EventStateDragAndDropScrollbar &&
		getEventStateId() != constants.EventStateDragAndDrop {
		characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if textFieldEntry, isFound := textFieldRegistry.Lookup(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias); characterEntry.AttributeEntry.CellType == constants.CellTypeTextField && isFound {
			if !textFieldEntry.IsEnabled {
				return isScreenUpdateRequired
			}
			textFieldEntry.CursorPosition = characterEntry.AttributeEntry.CellControlId
			shared.updateCursor(textFieldEntry)
			setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeTextField)
			isScreenUpdateRequired = true
		}
	}

	// Handle mouse drag for text selection. A window being dragged by its title bar also uses EventStateDragAndDrop,
	// and since that drag leaves focus where it was, it must be excluded here so it cannot select text in a focused
	// field it happens to pass over.
	if getEventStateId() == constants.EventStateDragAndDrop && getDraggedLayerAlias() == "" && getFocusedControl().controlType == constants.CellTypeTextField {
		characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if textFieldEntry, isFound := textFieldRegistry.Lookup(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias); characterEntry.AttributeEntry.CellType == constants.CellTypeTextField && isFound {
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
updateKeyboardEventTextboxWithString is a method which updates a text field with a string of characters. In addition, the
following should be noted:

- Processes each character in the string as a separate keystroke.

- Useful for programmatically inserting text into a text field.

Example:
    TextField.updateKeyboardEventTextboxWithString("Hello")
*/
func (shared *textFieldType) updateKeyboardEventTextboxWithString(keystroke string) {
	for _, currentCharacter := range keystroke {
		shared.updateKeyboardEvent([]rune{currentCharacter})
	}
}

/*
updateKeyboardEventTextboxWithCommands is a method which updates a text field with a list of command strings. In
addition, the following should be noted:

- Processes each command string as a separate keystroke.

- Useful for programmatically executing commands in a text field.

Example:
    TextField.updateKeyboardEventTextboxWithCommands("ctrl+a", "delete")
*/
func (shared *textFieldType) updateKeyboardEventTextboxWithCommands(keystroke ...string) {
	for _, currentCommand := range keystroke {
		shared.updateKeyboardEvent([]rune(currentCommand))
	}
}

/*
SetValue is a method which allows you to set the current value of a text field. The cursor is placed after the last
character, any highlight is cleared, and the viewport scrolls to show the cursor. If the text field is password
protected, the value is stored as given but displayed masked. If the text field does not exist, the request is
ignored.

Example:

	textField.SetValue("New Value")
*/
func (shared *TextFieldInstanceType) SetValue(value string) *TextFieldInstanceType {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry.CurrentValue = []rune(value + " ")
		textFieldEntry.CursorPosition = len(textFieldEntry.CurrentValue) - 1
		textFieldEntry.IsHighlightActive = false
		TextField.normalizeIndexes(textFieldEntry)
		TextField.updateViewport(textFieldEntry)
	}
	return shared
}

/*
GetCursorPosition is a method which allows you to obtain the zero based rune position of a text field's cursor within
its value. The position may equal the length of the value, which places the cursor just after its last character. If
the text field does not exist, 0 is returned.

Example:

	cursorPosition := textField.GetCursorPosition()
*/
func (shared *TextFieldInstanceType) GetCursorPosition() int {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return textFieldEntry.CursorPosition
	}
	return 0
}

/*
GetViewportPosition is a method which allows you to obtain the zero based rune position of the first character shown
in a text field, which moves as the field scrolls to keep its cursor visible. If the text field does not exist, 0 is
returned.

Example:

	viewportPosition := textField.GetViewportPosition()
*/
func (shared *TextFieldInstanceType) GetViewportPosition() int {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return textFieldEntry.ViewportPosition
	}
	return 0
}

/*
GetMaxLength is a method which allows you to obtain the maximum number of characters a text field accepts, as given
when it was added or set with SetMaxLength. If the text field does not exist, 0 is returned.

Example:

	maxLength := textField.GetMaxLength()
*/
func (shared *TextFieldInstanceType) GetMaxLength() int {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return textFieldEntry.MaxLengthAllowed
	}
	return 0
}

/*
GetDefaultValue is a method which allows you to obtain a text field's default value, as given when it was added or
set with SetDefaultValue. If the text field does not exist, an empty string is returned.

Example:

	defaultValue := textField.GetDefaultValue()
*/
func (shared *TextFieldInstanceType) GetDefaultValue() string {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return textFieldEntry.DefaultValue
	}
	return ""
}

/*
IsPasswordProtected is a method which allows you to check whether a text field masks its contents, as given when it
was added or set with SetPasswordProtected. If the text field does not exist, false is returned.

Example:

	isMasked := textField.IsPasswordProtected()
*/
func (shared *TextFieldInstanceType) IsPasswordProtected() bool {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return textFieldEntry.IsPasswordProtected
	}
	return false
}

/*
GetSelectedText is a method which allows you to obtain the text currently selected in a text field, whether the user
selected it or it was set with SetSelection. The blank sentinel the value is stored with is never included. If nothing
is selected, or the text field does not exist, an empty string is returned.

Example:

	selectedText := textField.GetSelectedText()
*/
func (shared *TextFieldInstanceType) GetSelectedText() string {
	defer lockUIState()()
	textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias)
	if !isFound {
		return ""
	}
	return TextField.getSelectedText(textFieldEntry)
}

/*
getSelectedText is a method which allows you to obtain the text selected in a text field exactly as GetSelectedText
returns it, without the blank sentinel the value is stored with, for a caller that already holds uiStateMutex. Copy
and cut use it for the text they put on the clipboard. If nothing is selected, an empty string is returned.

Example:

	selectedText := TextField.getSelectedText(textFieldEntry)
*/
func (shared *textFieldType) getSelectedText(textFieldEntry *types.TextFieldEntryType) string {
	if !textFieldEntry.IsHighlightActive {
		return ""
	}
	start, end := textFieldEntry.HighlightStart, textFieldEntry.HighlightEnd
	if start > end {
		start, end = end, start
	}
	// The last rune of the value is its sentinel, which is never part of the text.
	lastCharacterIndex := len(textFieldEntry.CurrentValue) - 2
	start = getClampedIndex(start, 0, start)
	end = getClampedIndex(end, -1, lastCharacterIndex)
	if start > end {
		return ""
	}
	return string(textFieldEntry.CurrentValue[start : end+1])
}

/*
SetSelection is a method which allows you to select the characters of a text field from one zero based rune position
to another, both included. The positions may be given in either order, and are clamped to the field's value. The
cursor is placed at the second position given and the field scrolls to show it, and the selection is drawn and then
used by copy, cut, and deletion exactly as one the user made. If the text field does not exist, no operation occurs.

Example:

	textField.SetSelection(0, 4)
*/
func (shared *TextFieldInstanceType) SetSelection(startPosition int, endPosition int) {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		textFieldEntry.HighlightStart = startPosition
		textFieldEntry.HighlightEnd = endPosition
		textFieldEntry.CursorPosition = endPosition
		textFieldEntry.IsHighlightActive = true
		textFieldEntry.IsHighlightModeToggled = false
		TextField.normalizeIndexes(textFieldEntry)
		TextField.updateViewport(textFieldEntry)
	}
}

/*
ClearSelection is a method which allows you to clear any text selected in a text field, leaving its value and cursor
where they are. If the text field does not exist, no operation occurs.

Example:

	textField.ClearSelection()
*/
func (shared *TextFieldInstanceType) ClearSelection() {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		textFieldEntry.IsHighlightActive = false
		textFieldEntry.IsHighlightModeToggled = false
	}
}

/*
SetDefaultValue is a method which sets the default value of your text field. In addition, the following should be noted:

- The default value will be used when the text field is reset.

Example:
    textField.SetDefaultValue("Default")
*/
func (shared *TextFieldInstanceType) SetDefaultValue(value string) *TextFieldInstanceType {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry.DefaultValue = value
	}
	return shared
}

/*
SetMaxLength is a method which sets the maximum number of characters allowed in the text field.

Example:
    textField.SetMaxLength(50)
*/
func (shared *TextFieldInstanceType) SetMaxLength(length int) *TextFieldInstanceType {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry.MaxLengthAllowed = length
	}
	return shared
}

/*
SetPasswordProtected is a method which specifies whether the text field should mask its contents. In addition, the
following should be noted:

- When enabled, all characters will be displayed as asterisks (*).

- The actual value is still stored and can be retrieved using GetValue.

Example:
    textField.SetPasswordProtected(true)
*/
func (shared *TextFieldInstanceType) SetPasswordProtected(isProtected bool) *TextFieldInstanceType {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry.IsPasswordProtected = isProtected
	}
	return shared
}

/*
SetOnValueChanged is a method which sets a hook that is invoked synchronously whenever the text field processes a
keystroke, allowing a caller to correct the field's value before the pending screen update is drawn. The hook
receives the field's current value in its trimmed form, the same form GetValue returns, and must return the value
the field should actually hold; returning the same value is a no-op. In addition, the following should be noted:

  - The hook runs inside the same keystroke handling call that changed the value, so any correction it returns
    replaces CurrentValue and is drawn in the same screen update as the keystroke, eliminating the one-frame flicker
    that would otherwise occur if the correction were applied on a later frame.

  - The hook runs on consolizer's event goroutine, not on the goroutine that called SetOnValueChanged, so it runs
    concurrently with the application's own code. Keep it free of shared state, ideally a pure function of the
    value it receives, or guard any state it must touch with a lock. It should also return quickly, since no further
    input is processed until it does.

  - The hook is called with the package's UI lock released, so it may call any exported function, including ones
    on this text field. Its correction is applied only if the field still holds the value the hook was given.

  - Passing nil removes any previously set hook.

Example:

	textField.SetOnValueChanged(func(current string) string {
		return current
	})
*/
func (shared *TextFieldInstanceType) SetOnValueChanged(fn func(current string) string) *TextFieldInstanceType {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry.OnValueChanged = fn
	}
	return shared
}

/*
SetCursorPosition is a method which allows you to set the zero based rune position of the cursor within a text field,
and scrolls the viewport to keep the cursor visible. A position before the start or past the end of the value is
clamped to the nearest valid position.

Example:

	textField.SetCursorPosition(10)
*/
func (shared *TextFieldInstanceType) SetCursorPosition(position int) *TextFieldInstanceType {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry.CursorPosition = position
		TextField.normalizeIndexes(textFieldEntry)
		TextField.updateViewport(textFieldEntry)
	}
	return shared
}

/*
SetViewportPosition is a method which allows you to set the zero based rune position of the first character shown in
a text field. A position before the start or past the end of the value is clamped to the nearest valid position. The
cursor is not moved, so it may lie outside the visible portion until the next keystroke scrolls it back into view.

Example:

	textField.SetViewportPosition(5)
*/
func (shared *TextFieldInstanceType) SetViewportPosition(position int) *TextFieldInstanceType {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		validatorTextField(shared.layerAlias, shared.controlAlias)
		textFieldEntry.ViewportPosition = position
		TextField.normalizeIndexes(textFieldEntry)
	}
	return shared
}

/*
GetTooltip is a method which retrieves the tooltip associated with this text field and returns the text field instance
for method chaining. In addition, the following should be noted:

- The tooltip is automatically created when the text field is added.

Example:
    textField.GetTooltip()
*/
func (shared *TextFieldInstanceType) GetTooltip() *TextFieldInstanceType {
	defer lockUIState()()
	// No need to retrieve the tooltip, just return self for chaining
	return shared
}

/*
SetTooltipText is a method which sets the text of the tooltip associated with the text field.

Example:
    textField.SetTooltipText("Helpful info")
*/
func (shared *TextFieldInstanceType) SetTooltipText(text string) *TextFieldInstanceType {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		var tooltipInstance TooltipInstanceType
		tooltipInstance.layerAlias = shared.layerAlias
		tooltipInstance.controlAlias = textFieldEntry.TooltipAlias
		tooltipInstance.setValue(text)
	}
	return shared
}

/*
EnableTooltip is a method which enables or disables the tooltip associated with the text field.

Example:
    textField.EnableTooltip(true)
*/
func (shared *TextFieldInstanceType) EnableTooltip(enabled bool) *TextFieldInstanceType {
	defer lockUIState()()
	if textFieldEntry, isFound := textFieldRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		var tooltipInstance TooltipInstanceType
		tooltipInstance.layerAlias = shared.layerAlias
		tooltipInstance.controlAlias = textFieldEntry.TooltipAlias
		tooltipInstance.setEnabled(enabled)
	}
	return shared
}
