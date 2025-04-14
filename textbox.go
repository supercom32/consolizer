package consolizer

import (
	"fmt"
	"math"
	"strings"

	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/internal/memory"
	"supercom32.net/consolizer/internal/stringformat"
	"supercom32.net/consolizer/types"
)

type TextboxInstanceType struct {
	layerAlias   string
	controlAlias string
}

type textboxType struct{}

var textbox textboxType
var Textboxes = memory.NewControlMemoryManager[types.TextboxEntryType]()

func AddTextbox(layerAlias string, textboxAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isBorderDrawn bool) {
	textboxEntry := types.NewTexboxEntry()
	textboxEntry.Alias = textboxAlias
	textboxEntry.StyleEntry = styleEntry
	textboxEntry.XLocation = xLocation
	textboxEntry.YLocation = yLocation
	textboxEntry.Width = width
	textboxEntry.Height = height
	textboxEntry.IsBorderDrawn = isBorderDrawn
	textboxEntry.TooltipAlias = stringformat.GetLastSortedUUID()

	// Create associated tooltip (always created but disabled by default)
	Tooltip.Add(layerAlias, textboxEntry.TooltipAlias, "", styleEntry,
		textboxEntry.XLocation, textboxEntry.YLocation,
		textboxEntry.Width, textboxEntry.Height,
		textboxEntry.XLocation, textboxEntry.YLocation+textboxEntry.Height+1,
		textboxEntry.Width, 3,
		false, true, constants.DefaultTooltipHoverTime)

	// Use the generic memory manager to add the textbox entry
	Textboxes.Add(layerAlias, textboxAlias, &textboxEntry)
}

func GetTextbox(layerAlias string, textboxAlias string) *types.TextboxEntryType {
	// Use the generic memory manager to retrieve the textbox entry
	textboxEntry := Textboxes.Get(layerAlias, textboxAlias)
	if textboxEntry == nil {
		panic(fmt.Sprintf("The requested text with alias '%s' on layer '%s' could not be returned since it does not exist.", textboxAlias, layerAlias))
	}
	return textboxEntry
}

func IsTextboxExists(layerAlias string, textboxAlias string) bool {
	// Use the generic memory manager to check existence
	return Textboxes.Get(layerAlias, textboxAlias) != nil
}

func DeleteTextbox(layerAlias string, textboxAlias string) {
	// Use the generic memory manager to remove the textbox entry
	Textboxes.Remove(layerAlias, textboxAlias)
}

/*
DeleteAllTextboxesFromLayer allows you to delete all textboxes from a given layer. In addition, the following
information should be noted:

- All textboxes on the specified layer will be removed.
- All memory associated with the textboxes will be freed.
- The textboxes will be removed from the tab index if they were added.
*/
func DeleteAllTextboxesFromLayer(layerAlias string) {
	// Retrieve all textboxes in the specified layer
	textboxes := Textboxes.GetAllEntries(layerAlias)

	// Loop through all entries and delete them
	for _, textbox := range textboxes {
		Textboxes.Remove(layerAlias, textbox.Alias) // Assuming textbox.Alias contains the alias
	}
}

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
AddToTabIndex allows you to add a textbox to the tab index. This enables keyboard navigation
between controls using the tab key. In addition, the following information should be noted:

- The textbox will be added to the tab order based on the order in which it was created.
- The tab index is used to determine which control receives focus when the tab key is pressed.
*/
func (shared *TextboxInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeTextbox)
}

/*
Delete allows you to remove a textbox from a text layer. In addition, the following
information should be noted:

- If you attempt to delete a textbox which does not exist, then the request
will simply be ignored.
- All memory associated with the textbox will be freed.
*/
func (shared *TextboxInstanceType) Delete() string {
	if Textboxes.IsExists(shared.layerAlias, shared.controlAlias) {
		Textboxes.Remove(shared.layerAlias, shared.controlAlias)
	}
	return ""
}

/*
setText allows you to set the text for a textbox. If the textbox instance
no longer exists, then no operation takes place. In addition, the following
information should be noted:

- Text can be broke up into multiple lines by using the '\n' escape sequence.
*/
func (shared *TextboxInstanceType) setText(text string) {
	if Textboxes.IsExists(shared.layerAlias, shared.controlAlias) {
		textData := strings.Split(text, "\n")
		textboxEntry := Textboxes.Get(shared.layerAlias, shared.controlAlias)
		for _, text := range textData {
			textboxEntry.TextData = append(textboxEntry.TextData, stringformat.GetRunesFromString(text))
		}
		textbox.setTextboxMaxScrollBarValues(shared.layerAlias, shared.controlAlias)
	}
}

/*
setViewport allows you to set the current viewport for a textbox. If the textbox instance
no longer exists, then no operation takes place.
*/
func (shared *TextboxInstanceType) setViewport(xLocation int, yLocation int) {
	if Textboxes.IsExists(shared.layerAlias, shared.controlAlias) {
		textboxEntry := Textboxes.Get(shared.layerAlias, shared.controlAlias)
		textboxEntry.ViewportXLocation = xLocation
		textboxEntry.ViewportYLocation = yLocation
	}
}

/*
getTextboxClickCoordinates converts a cell ID to x and y coordinates within a textbox.
The conversion is based on the table width, where:
- x coordinate is calculated as cellId modulo tableWidth
- y coordinate is calculated as cellId divided by tableWidth (rounded down)

Parameters:
  - cellId: The ID of the cell to convert to coordinates
  - tableWidth: The width of the textbox table in cells

Returns:
  - xLocation: The x coordinate of the cell
  - yLocation: The y coordinate of the cell

This function is primarily used to determine cursor position from mouse clicks within a textbox.
*/
func (shared *textboxType) getTextboxClickCoordinates(cellId int, tableWidth int) (int, int) {
	xLocation := cellId % tableWidth
	yLocation := math.Floor(float64(cellId) / float64(tableWidth))
	return xLocation, int(yLocation)
}

/*
insertCharacterUsingAbsoluteCoordinates allows you to insert a character at a specific position in a textbox. In addition,
the following information should be noted:

- The character is inserted at the specified x and y coordinates.
- The text after the insertion point is shifted right.
- The cursor position is updated to after the inserted character.
*/
func (shared *textboxType) insertCharacterUsingAbsoluteCoordinates(textboxEntry *types.TextboxEntryType, xLocation int, yLocation int, characterToInsert rune) {
	stringDataSuffixAfterInsert := textboxEntry.TextData[yLocation][xLocation:len(textboxEntry.TextData[yLocation])]
	textboxEntry.TextData[yLocation] = append([]rune{}, textboxEntry.TextData[yLocation][:xLocation]...)
	textboxEntry.TextData[yLocation] = append(textboxEntry.TextData[yLocation][:xLocation], characterToInsert)
	textboxEntry.TextData[yLocation] = append(textboxEntry.TextData[yLocation], stringDataSuffixAfterInsert...)
	textboxEntry.CursorXLocation++
}

/*
backspaceCharacterUsingRelativeCoordinates allows you to delete the character before the cursor. In addition,
the following information should be noted:

- If at the beginning of a line, moves the cursor to the end of the previous line.
- If at the beginning of the first line, no action is taken.
- The cursor position is updated after the deletion.
*/
func (shared *textboxType) backspaceCharacterUsingRelativeCoordinates(textboxEntry *types.TextboxEntryType) {
	// If nothing left to backspace, do nothing.
	if textboxEntry.CursorXLocation == 0 && textboxEntry.CursorYLocation == 0 {
		return
	} else if textboxEntry.CursorXLocation == 0 {
		// If at the beginning of a line, move cursor the previous line ending.
		textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation-1]) - 1
		shared.moveRemainingTextToPreviousLine(textboxEntry, textboxEntry.CursorYLocation)
		textboxEntry.CursorYLocation--
		return
	}
	// Otherwise, just backspace a single character.
	textboxEntry.CursorXLocation--
	shared.deleteCharacterUsingAbsoluteCoordinates(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
}

/*
deleteCharacterUsingAbsoluteCoordinates allows you to delete a character at a specific position. In addition,
the following information should be noted:

- The character at the specified x and y coordinates is deleted.
- If at the end of a line, moves the next line up if possible.
- The cursor position is updated after the deletion.
*/
func (shared *textboxType) deleteCharacterUsingAbsoluteCoordinates(textboxEntry *types.TextboxEntryType, xLocation int, yLocation int) {
	// If cursor yLocation is out of bounds, do nothing.
	if yLocation >= len(textboxEntry.TextData) {
		return
	}
	// If cursor xLocation is at/out of bounds, move previous line to current line and position.
	if xLocation >= len(textboxEntry.TextData[yLocation])-1 {
		if len(textboxEntry.TextData)-1 == yLocation && (len(textboxEntry.TextData[yLocation]) <= 1 || xLocation >= len(textboxEntry.TextData[yLocation])-1) {
			return
		}
		textboxEntry.TextData[yLocation] = textboxEntry.TextData[yLocation][:len(textboxEntry.TextData[yLocation])-1]
		textboxEntry.TextData[yLocation] = append(textboxEntry.TextData[yLocation], textboxEntry.TextData[yLocation+1]...)
		copy(textboxEntry.TextData[yLocation+1:], textboxEntry.TextData[yLocation+2:])
		textboxEntry.TextData = textboxEntry.TextData[:len(textboxEntry.TextData)-1]
		return
	}
	// Remove the current character.
	stringDataSuffixAfterInsert := textboxEntry.TextData[yLocation][xLocation+1 : len(textboxEntry.TextData[yLocation])]
	textboxEntry.TextData[yLocation] = append([]rune{}, textboxEntry.TextData[yLocation][:xLocation]...)
	textboxEntry.TextData[yLocation] = append(textboxEntry.TextData[yLocation], stringDataSuffixAfterInsert...)
}

/*
moveRemainingTextToPreviousLine allows you to move text from the current line to the previous line. In addition,
the following information should be noted:

- If the cursor is at the beginning of a line, the text after the cursor is moved to the previous line.
- The cursor position is updated to the end of the previous line.
- The current line is removed if it becomes empty after the move.
*/
func (shared *textboxType) moveRemainingTextToPreviousLine(textboxEntry *types.TextboxEntryType, yLocation int) {
	// If the only row of text or the cursor yLocation is out of bounds, then exit.
	if len(textboxEntry.TextData) == 1 || yLocation >= len(textboxEntry.TextData) {
		return
	}
	textboxEntry.TextData[yLocation-1] = textboxEntry.TextData[yLocation-1][:len(textboxEntry.TextData[yLocation-1])-1]
	textboxEntry.TextData[yLocation-1] = append(textboxEntry.TextData[yLocation-1], textboxEntry.TextData[yLocation]...)
	textboxEntry.TextData = shared.removeLine(textboxEntry.TextData, yLocation)
}

/*
removeLine allows you to remove a line from a textbox. In addition, the following
information should be noted:

- The line at the specified index is removed.
- The remaining lines are shifted up to fill the gap.
- Returns the modified text data array.
*/
func (shared *textboxType) removeLine(textData [][]rune, index int) [][]rune {
	return append(textData[:index], textData[index+1:]...)
}

/*
insertLine allows you to insert a new line into a textbox. In addition, the following
information should be noted:

- A new line is inserted at the specified index.
- If the index is out of bounds, the line is appended to the end.
- Returns the modified text data array.
*/
func (shared *textboxType) insertLine(textData [][]rune, lineData []rune, index int) [][]rune {
	// If the index provided is inbounds, insert the line data accordingly.
	if index < len(textData) {
		textData = append(textData[:index+1], textData[index:]...)
		textData[index] = lineData
	} else {
		// Otherwise, append the line data to the end of the array.
		textData = append(textData, []rune{' '})
	}
	return textData
}

/*
moveTextAfterCursorToNextLine allows you to move text after your cursor to a new line underneath it. In addition,
the following information should be noted:

- Creates a new line with a default space character.
- Copies all text after the cursor position to the new line.
- Truncates the current line at the cursor position.
- Updates the cursor position to the start of the new line.
- Maintains proper text formatting and cursor visibility.
*/
func (shared *textboxType) moveTextAfterCursorToNextLine(textboxEntry *types.TextboxEntryType, yLocation int) {
	// Create a new line with our default ' ' rune.
	textboxEntry.TextData = shared.insertLine(textboxEntry.TextData, []rune{' '}, yLocation+1)
	// Copy everything past our cursor on the current line.
	charactersToCopy := textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.CursorXLocation:]
	copyOfCharacters := make([]rune, len(textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.CursorXLocation:]))
	copy(copyOfCharacters, charactersToCopy)
	// Make our current line = everything up to our cursor + ' ' ending.
	textboxEntry.TextData[textboxEntry.CursorYLocation] = append(textboxEntry.TextData[textboxEntry.CursorYLocation][:textboxEntry.CursorXLocation], ' ')
	// Paste the copied text to our new line.
	textboxEntry.CursorYLocation++
	textboxEntry.CursorXLocation = 0
	textboxEntry.TextData[textboxEntry.CursorYLocation] = make([]rune, len(copyOfCharacters))
	copy(textboxEntry.TextData[textboxEntry.CursorYLocation], copyOfCharacters)
}

/*
updateScrollbarBasedOnTextboxViewport allows you to update the scrollbar positions based on the current
viewport position of a textbox. In addition, the following information should be noted:

- Updates both horizontal and vertical scrollbars if they exist.
- Adjusts the scrollbar handle positions to match the viewport position.
- Ensures the scrollbars accurately reflect the current view of the textbox.
*/
func (shared *textboxType) updateScrollbarBasedOnTextboxViewport(layerAlias string, textboxAlias string) {
	textboxEntry := Textboxes.Get(layerAlias, textboxAlias)
	horizontalScrollbarEntry := ScrollBars.Get(layerAlias, textboxEntry.HorizontalScrollbarAlias)
	horizontalScrollbarEntry.ScrollValue = textboxEntry.ViewportXLocation
	scrollbar.computeScrollbarHandlePositionByScrollValue(layerAlias, textboxEntry.HorizontalScrollbarAlias)
	verticalScrollbarEntry := ScrollBars.Get(layerAlias, textboxEntry.VerticalScrollbarAlias)
	verticalScrollbarEntry.ScrollValue = textboxEntry.ViewportYLocation
	scrollbar.computeScrollbarHandlePositionByScrollValue(layerAlias, textboxEntry.VerticalScrollbarAlias)
}

/*
getMaxHorizontalTextValue returns the maximum line length in a textbox. In addition,
the following information should be noted:

- Calculates the maximum width of any line in the textbox.
- Takes into account wide characters that take up multiple spaces.
- Used to determine horizontal scrollbar limits.
*/
func (shared *textboxType) getMaxHorizontalTextValue(layerAlias string, textboxAlias string) int {
	textboxEntry := Textboxes.Get(layerAlias, textboxAlias)
	maxHorizontalValue := 0
	for _, currentLine := range textboxEntry.TextData {
		lengthOfLine := stringformat.GetWidthOfRunesWhenPrinted(currentLine)
		over := lengthOfLine - len(currentLine)
		if lengthOfLine > maxHorizontalValue {
			maxHorizontalValue = lengthOfLine - (over / 2)
		}
	}
	return maxHorizontalValue
}

/*
setTextboxMaxScrollBarValues allows you to update the scrollbar limits based on text content. In addition,
the following information should be noted:

- Updates both horizontal and vertical scrollbar maximum values.
- Disables scrollbars if the content fits within the viewport.
- Ensures scrollbars accurately reflect the text dimensions.
*/
func (shared *textboxType) setTextboxMaxScrollBarValues(layerAlias string, textboxAlias string) {
	textboxEntry := Textboxes.Get(layerAlias, textboxAlias)
	maxVerticalValue := len(textboxEntry.TextData)
	maxHorizontalValue := shared.getMaxHorizontalTextValue(layerAlias, textboxAlias)
	hScrollBarEntry := ScrollBars.Get(layerAlias, textboxEntry.HorizontalScrollbarAlias)
	vScrollBarEntry := ScrollBars.Get(layerAlias, textboxEntry.VerticalScrollbarAlias)
	maxHorizontalValue = maxHorizontalValue - textboxEntry.Width
	// If the max horizontal width is smaller than the textbox width, disable scrolling.
	if maxHorizontalValue <= 0 {
		maxHorizontalValue = 0
		hScrollBarEntry.IsEnabled = false
		hScrollBarEntry.IsVisible = false
	} else {
		hScrollBarEntry.IsEnabled = true
		hScrollBarEntry.IsVisible = true
	}
	maxVerticalValue = maxVerticalValue - textboxEntry.Height
	// If the max vertical height is smaller than the textbox height, disable scrolling.
	if maxVerticalValue <= 0 {
		maxVerticalValue = 0
		vScrollBarEntry.IsEnabled = false
		vScrollBarEntry.IsVisible = false
	} else {
		vScrollBarEntry.IsEnabled = true
		vScrollBarEntry.IsVisible = true
	}
	hScrollBarEntry.MaxScrollValue = maxHorizontalValue
	vScrollBarEntry.MaxScrollValue = maxVerticalValue
}

/*
AddTextbox allows you to add a text box to a text layer. Once called, an instance of your control is
returned which will allow you to read or manipulate the properties for it. The Style of the text box
will be determined by the style entry passed in. If you wish to remove a text box from a text
layer, simply call 'DeleteTextBox'. In addition, the following information should be noted:

- Text boxes are not drawn physically to the text layer provided. Instead
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create text boxes without actually overwriting
the text layer data under it.

- If the text box to be drawn falls outside the range of the provided layer,
then only the visible portion of the text box will be drawn.
*/
func (shared *textboxType) AddTextbox(layerAlias string, textboxAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isBorderDrawn bool) TextboxInstanceType {
	newTextboxEntry := types.NewTexboxEntry()
	newTextboxEntry.Alias = textboxAlias
	newTextboxEntry.StyleEntry = styleEntry
	newTextboxEntry.XLocation = xLocation
	newTextboxEntry.YLocation = yLocation
	newTextboxEntry.Width = width
	newTextboxEntry.Height = height
	newTextboxEntry.IsBorderDrawn = isBorderDrawn
	newTextboxEntry.TooltipAlias = stringformat.GetLastSortedUUID()

	tooltipInstance := Tooltip.Add(layerAlias, newTextboxEntry.TooltipAlias, "", styleEntry,
		newTextboxEntry.XLocation, newTextboxEntry.YLocation,
		newTextboxEntry.Width, newTextboxEntry.Height,
		newTextboxEntry.XLocation, newTextboxEntry.YLocation+1,
		newTextboxEntry.Width, newTextboxEntry.Height,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.SetEnabled(false)
	tooltipInstance.setParentControlAlias(textboxAlias)
	// Use the generic memory manager to add the textbox entry
	Textboxes.Add(layerAlias, textboxAlias, &newTextboxEntry)
	textboxEntry := Textboxes.Get(layerAlias, textboxAlias)
	textboxEntry.TextData = append(textboxEntry.TextData, stringformat.GetRunesFromString(" "))
	textboxEntry.HorizontalScrollbarAlias = stringformat.GetLastSortedUUID()
	textboxEntry.VerticalScrollbarAlias = stringformat.GetLastSortedUUID()
	hScrollbarWidth := width
	hScrollbarXLocation := xLocation
	hScrollbarYLocation := yLocation + height
	vScrollbarHeight := height
	vScrollbarXLocation := xLocation + width
	vScrollbarYLocation := yLocation
	if isBorderDrawn == true {
		hScrollbarYLocation++
		hScrollbarXLocation--
		vScrollbarXLocation++
		vScrollbarYLocation--
		hScrollbarWidth = hScrollbarWidth + 2
		vScrollbarHeight = vScrollbarHeight + 2
	}
	scrollbar.Add(layerAlias, textboxEntry.HorizontalScrollbarAlias, styleEntry, hScrollbarXLocation, hScrollbarYLocation, hScrollbarWidth, 0, 0, 1, true)
	scrollbar.Add(layerAlias, textboxEntry.VerticalScrollbarAlias, styleEntry, vScrollbarXLocation, vScrollbarYLocation, vScrollbarHeight, 0, 0, 1, false)
	shared.setTextboxMaxScrollBarValues(layerAlias, textboxAlias)
	var textboxInstance TextboxInstanceType
	textboxInstance.layerAlias = layerAlias
	textboxInstance.controlAlias = textboxAlias
	return textboxInstance
}

/*
DeleteTextbox allows you to remove a text box from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a text box which does not exist, then the request
will simply be ignored.
*/
func (shared *textboxType) DeleteTextbox(layerAlias string, textboxAlias string) {
	Textboxes.Remove(layerAlias, textboxAlias)
}

/*
DeleteAllTextboxes allows you to delete all textboxes on a given layer. In addition, the following
information should be noted:

- All textboxes on the specified layer will be removed.
- All memory associated with the textboxes will be freed.
- The textboxes will be removed from the tab index if they were added.
*/
func (shared *textboxType) DeleteAllTextboxes(layerAlias string) {
	Textboxes.RemoveAll(layerAlias)
}

/*
drawTextboxesOnLayer allows you to draw all textboxes on a layer. In addition, the following
information should be noted:

- Iterates through all textboxes on the specified layer.
- Draws each textbox with its current content and style.
- Handles cursor and highlight rendering.
*/
func (shared *textboxType) drawTextboxesOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentTextBoxEntry := range Textboxes.GetAllEntries(layerAlias) {
		shared.drawTextbox(&layerEntry, currentTextBoxEntry.Alias)
	}
}

/*
drawTextbox allows you to draw a textbox on a given text layer. In addition, the following
information should be noted:

- Draws the textbox with its current content and style.
- Handles border drawing if enabled.
- Manages cursor and highlight rendering.
- Adjusts the viewport to show the correct portion of text.
*/
func (shared *textboxType) drawTextbox(layerEntry *types.LayerEntryType, textboxAlias string) {
	t := Textboxes.Get(layerEntry.LayerAlias, textboxAlias)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = t.StyleEntry.TextboxForegroundColor
	attributeEntry.BackgroundColor = t.StyleEntry.TextboxBackgroundColor
	attributeEntry.CellControlAlias = textboxAlias
	if t.IsBorderDrawn {
		drawBorder(layerEntry, t.StyleEntry, attributeEntry, t.XLocation-1, t.YLocation-1, t.Width+2, t.Height+2, false)
	}
	attributeEntry.CellType = constants.CellTypeTextbox
	fillArea(layerEntry, attributeEntry, " ", t.XLocation, t.YLocation, t.Width, t.Height, t.ViewportYLocation)
	attributeEntry.CellControlAlias = textboxAlias
	for currentLine := 0; currentLine < t.Height; currentLine++ {
		var arrayOfRunes []rune
		if t.ViewportYLocation+currentLine < len(t.TextData) && t.ViewportYLocation+currentLine >= 0 {
			arrayOfRunes = t.TextData[t.ViewportYLocation+currentLine]
			if t.ViewportXLocation < len(arrayOfRunes) && t.ViewportXLocation >= 0 {
				if t.ViewportXLocation+t.Width < len(arrayOfRunes) {
					arrayOfRunes = arrayOfRunes[t.ViewportXLocation : t.ViewportXLocation+t.Width]
				} else {
					arrayOfRunes = arrayOfRunes[t.ViewportXLocation:]
				}
			} else {
				// If scrolled too far right and there are no column text to print, just show blanks.
				// If scrolled too far left (negative value) then show blanks. Note: This case should never happen really.
				arrayOfRunes = []rune{}
			}
			// arrayOfRunes = stringformat.GetFormattedRuneArray(arrayOfRunes, t.Width, constants.AlignmentLeft)
			arrayOfRunes = stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunes, t.Width)
			shared.printControlText(layerEntry, textboxAlias, t.StyleEntry, attributeEntry, t.XLocation, t.YLocation+currentLine, arrayOfRunes, t.ViewportYLocation+currentLine, t.ViewportXLocation, t.CursorXLocation, t.CursorYLocation)
		} else {
			// If scrolled too far down and there are no more rows to print, just show blanks.
			// If scrolled too far up and there are no rows to print, just print blanks. Note: This case should never happen really.
			// arrayOfRunes = stringformat.GetFormattedRuneArray([]rune{}, t.Width, constants.AlignmentLeft)
			shared.printControlText(layerEntry, textboxAlias, t.StyleEntry, attributeEntry, t.XLocation, t.YLocation+currentLine, arrayOfRunes, t.ViewportYLocation+currentLine, t.ViewportXLocation, t.CursorXLocation, t.CursorYLocation)
		}
	}
}

/*
drawBorder allows you to draw a border around a textbox. In addition, the following
information should be noted:

- Draws a border using the specified style and attributes.
- Handles both single and double line borders.
- Adjusts the border position based on the textbox dimensions.
*/
func (shared *textboxType) drawBorder(layerEntry *types.LayerEntryType, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, height int, isDoubleLine bool) {
	// Implementation of drawBorder method
}

/*
drawTextboxContent allows you to draw the content of a textbox. In addition, the following
information should be noted:

- Draws the text content with proper formatting and alignment.
- Handles wide characters and line wrapping.
- Manages cursor and highlight rendering.
*/
func (shared *textboxType) drawTextboxContent(layerEntry *types.LayerEntryType, textboxAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, height int) {
	// Implementation of drawTextboxContent method
}

/*
drawTextboxCursor allows you to draw the cursor in a textbox. In addition, the following
information should be noted:

- Draws the cursor at the current position.
- Uses the specified cursor style and attributes.
- Handles cursor visibility and blinking.
*/
func (shared *textboxType) drawTextboxCursor(layerEntry *types.LayerEntryType, textboxAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int) {
	// Implementation of drawTextboxCursor method
}

/*
drawTextboxHighlight allows you to draw highlighted text in a textbox. In addition, the following
information should be noted:

- Draws the highlighted text with inverted colors.
- Handles both single-line and multi-line highlights.
- Manages highlight start and end positions.
*/
func (shared *textboxType) drawTextboxHighlight(layerEntry *types.LayerEntryType, textboxAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, height int) {
	// Implementation of drawTextboxHighlight method
}

/*
drawTextboxScrollbars allows you to draw the scrollbars for a textbox. In addition, the following
information should be noted:

- Draws both horizontal and vertical scrollbars if enabled.
- Updates scrollbar positions based on viewport position.
- Handles scrollbar visibility and interaction.
*/
func (shared *textboxType) drawTextboxScrollbars(layerEntry *types.LayerEntryType, textboxAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, height int) {
	// Implementation of drawTextboxScrollbars method
}

/*
drawTextboxScrollbar allows you to draw a single scrollbar for a textbox. In addition, the following
information should be noted:

- Draws either a horizontal or vertical scrollbar.
- Updates the scrollbar handle position based on the current scroll value.
- Handles scrollbar interaction and updates.
*/
func (shared *textboxType) drawTextboxScrollbar(layerEntry *types.LayerEntryType, textboxAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, length int, isHorizontal bool) {
	// Implementation of drawTextboxScrollbar method
}

/*
drawTextboxScrollbarHandle allows you to draw the handle of a scrollbar. In addition, the following
information should be noted:

- Draws the scrollbar handle at the current position.
- Updates the handle position based on the scroll value.
- Handles handle interaction and updates.
*/
func (shared *textboxType) drawTextboxScrollbarHandle(layerEntry *types.LayerEntryType, textboxAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, length int, isHorizontal bool) {
	// Implementation of drawTextboxScrollbarHandle method
}

/*
drawTextboxScrollbarTrack allows you to draw the track of a scrollbar. In addition, the following
information should be noted:

- Draws the scrollbar track with the specified style.
- Handles both horizontal and vertical tracks.
- Updates the track appearance based on scrollbar state.
*/
func (shared *textboxType) drawTextboxScrollbarTrack(layerEntry *types.LayerEntryType, textboxAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, length int, isHorizontal bool) {
	// Implementation of drawTextboxScrollbarTrack method
}

/*
drawTextboxScrollbarArrows allows you to draw the arrows of a scrollbar. In addition, the following
information should be noted:

- Draws both up/down or left/right arrows for the scrollbar.
- Handles arrow interaction and updates.
- Updates arrow appearance based on scrollbar state.
*/
func (shared *textboxType) drawTextboxScrollbarArrows(layerEntry *types.LayerEntryType, textboxAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, length int, isHorizontal bool) {
	// Implementation of drawTextboxScrollbarArrows method
}

/*
printControlText allows you to print text with control IDs. In addition, the following
information should be noted:

- Prints each character with its associated control ID.
- Handles wide characters that take up multiple spaces.
- Manages cursor and highlight rendering.
*/
func (shared *textboxType) printControlText(layerEntry *types.LayerEntryType, textboxAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, arrayOfRunes []rune, controlYLocation int, startingControlId int, cursorXLocation int, cursorYLocation int) {
	currentControlId := startingControlId
	currentXOffset := 0
	for _, currentCharacter := range arrayOfRunes {
		attributeEntry.CellControlId = currentControlId
		attributeEntry.CellControlLocation = controlYLocation
		// If the textbox being drawn is focused, render the cursor as well.
		if isControlCurrentlyFocused(layerEntry.LayerAlias, textboxAlias, constants.CellTypeTextbox) {
			textboxEntry := Textboxes.Get(layerEntry.LayerAlias, textboxAlias)
			if textboxEntry.IsHighlightActive {
				// Check if current position is within highlight range
				isHighlighted := false

				// Determine the correct start and end positions for highlighting
				var highlightStartX, highlightEndX, highlightStartY, highlightEndY int

				// If cursor is to the left of the start position
				if textboxEntry.CursorYLocation < textboxEntry.HighlightStartY ||
					(textboxEntry.CursorYLocation == textboxEntry.HighlightStartY &&
						textboxEntry.CursorXLocation < textboxEntry.HighlightStartX) {
					// Cursor is before the highlight start, so swap the positions
					highlightStartX = textboxEntry.CursorXLocation
					highlightStartY = textboxEntry.CursorYLocation
					highlightEndX = textboxEntry.HighlightStartX
					highlightEndY = textboxEntry.HighlightStartY
				} else {
					// Cursor is after or at the highlight start
					highlightStartX = textboxEntry.HighlightStartX
					highlightStartY = textboxEntry.HighlightStartY
					highlightEndX = textboxEntry.CursorXLocation
					highlightEndY = textboxEntry.CursorYLocation
				}

				// Check if the current position is within the highlight range
				if controlYLocation >= highlightStartY && controlYLocation <= highlightEndY {
					if controlYLocation == highlightStartY && controlYLocation == highlightEndY {
						// Same line highlight
						isHighlighted = currentControlId >= highlightStartX && currentControlId <= highlightEndX
					} else if controlYLocation == highlightStartY {
						// First line of multi-line highlight
						isHighlighted = currentControlId >= highlightStartX
					} else if controlYLocation == highlightEndY {
						// Last line of multi-line highlight
						isHighlighted = currentControlId <= highlightEndX
					} else {
						// Middle line of multi-line highlight
						isHighlighted = true
					}
				}

				if isHighlighted {
					attributeEntry.ForegroundColor = styleEntry.HighlightForegroundColor
					attributeEntry.BackgroundColor = styleEntry.HighlightBackgroundColor
				} else if cursorXLocation == currentControlId && cursorYLocation == controlYLocation {
					attributeEntry.ForegroundColor = styleEntry.TextboxCursorForegroundColor
					attributeEntry.BackgroundColor = styleEntry.TextboxCursorBackgroundColor
				} else {
					attributeEntry.ForegroundColor = styleEntry.TextboxForegroundColor
					attributeEntry.BackgroundColor = styleEntry.TextboxBackgroundColor
				}
			} else if cursorXLocation == currentControlId && cursorYLocation == controlYLocation {
				attributeEntry.ForegroundColor = styleEntry.TextboxCursorForegroundColor
				attributeEntry.BackgroundColor = styleEntry.TextboxCursorBackgroundColor
			} else {
				attributeEntry.ForegroundColor = styleEntry.TextboxForegroundColor
				attributeEntry.BackgroundColor = styleEntry.TextboxBackgroundColor
			}
		}
		printLayer(layerEntry, attributeEntry, xLocation+currentXOffset, yLocation, []rune{currentCharacter})
		if stringformat.IsRuneCharacterWide(currentCharacter) {
			// If we find a wide character, we add a blank space with the same ID as the previous
			// character so the next printed character doesn't get covered by the wide one.
			currentXOffset++
			printLayer(layerEntry, attributeEntry, xLocation+currentXOffset, yLocation, []rune{' '})
			currentXOffset++
		} else {
			currentXOffset++
		}
		currentControlId++
	}
}

/*
updateCursor allows you to update the cursor position in a textbox. In addition, the following
information should be noted:

- Ensures the cursor stays within valid bounds.
- Handles cases where the cursor moves outside the text.
- Updates the cursor position in the textbox entry.
*/
func (shared *textboxType) updateCursor(textboxEntry *types.TextboxEntryType, xLocation int, yLocation int) {
	textboxEntry.CursorXLocation = xLocation
	textboxEntry.CursorYLocation = yLocation
	// If yLocation is less than text data bounds.
	if textboxEntry.CursorYLocation < 0 {
		textboxEntry.CursorYLocation = 0
	}
	// If yLocation is greater than column data bounds.
	if textboxEntry.CursorYLocation > len(textboxEntry.TextData)-1 {
		textboxEntry.CursorYLocation = len(textboxEntry.TextData) - 1
	}
	// If our cursor xLocation was jumped (due to NullCellControlId) or placed in an invalid xLocation spot greater than the length of our text line.
	// Move it to the end of the line.
	if textboxEntry.CursorXLocation == constants.NullCellControlId || textboxEntry.CursorXLocation > len(textboxEntry.TextData[textboxEntry.CursorYLocation])-1 {
		textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
	} else if textboxEntry.CursorXLocation < 0 {
		// If the cursor xLocation was scrolled to be out of lower bounds, just set the location to 0.
		textboxEntry.CursorXLocation = 0
	}
}

/*
updateViewport allows you to update the visible portion of a textbox. In addition, the following
information should be noted:

- Adjusts the viewport to keep the cursor visible.
- Handles cases where the cursor moves outside the viewport.
- Updates the viewport position in the textbox entry.
*/
func (shared *textboxType) updateViewport(textboxEntry *types.TextboxEntryType) {
	// If cursor yLocation is higher than the viewport window, move the window to make the cursor appear at the end.
	if textboxEntry.CursorYLocation >= textboxEntry.ViewportYLocation+textboxEntry.Height {
		textboxEntry.ViewportYLocation = textboxEntry.CursorYLocation - textboxEntry.Height + 1
	}
	// If cursor yLocation is lower than viewport window, make the viewport window start at yLocation.
	if textboxEntry.CursorYLocation < textboxEntry.ViewportYLocation {
		textboxEntry.ViewportYLocation = textboxEntry.CursorYLocation
	}
	// If cursor yLocation is less than 0 (Out of range), just set viewport window to 0.
	if textboxEntry.CursorYLocation <= 0 {
		textboxEntry.ViewportYLocation = 0
	}

	// If cursor xLocation is lower than the viewport window
	if textboxEntry.CursorXLocation < textboxEntry.ViewportXLocation {
		// LogInfo("YES1 " + fmt.Sprintf("%d", time.Now().Unix()))
		isCursorJumped := false
		// Detect if the cursor xLocation was jumped or if it was scrolled.
		if textboxEntry.ViewportXLocation-textboxEntry.CursorXLocation > 2 || textboxEntry.CursorXLocation-textboxEntry.ViewportXLocation > 2 {
			isCursorJumped = true
		}
		// If the cursor xLocation is less than the size of our viewport and was jumped, just set the viewport to 0.
		if isCursorJumped && textboxEntry.CursorXLocation-textboxEntry.Width < 0 {
			textboxEntry.ViewportXLocation = 0
		} else {
			// Otherwise, this is a normal backwards scroll so make viewport equal to our cursor location.
			textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation
		}
	}
	// Figure out how much displayable space is in our current viewport window.
	arrayOfRunesAvailableToPrint := textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.ViewportXLocation:]
	arrayOfRunesThatFitStringSize := stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvailableToPrint, textboxEntry.Width)
	// If the cursor xLocation is equal or greater than the visible viewport window width.
	if textboxEntry.CursorXLocation >= textboxEntry.ViewportXLocation+len(arrayOfRunesThatFitStringSize) {
		// Then make the viewport xLocation equal to the visible viewport width behind it.
		maxViewportWidthAvaliable := textboxEntry.Width
		if textboxEntry.CursorXLocation-textboxEntry.Width < 0 {
			maxViewportWidthAvaliable = textboxEntry.CursorXLocation
		}
		arrayOfRunesAvailableToPrint = textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.CursorXLocation-maxViewportWidthAvaliable : textboxEntry.CursorXLocation]
		numberOfRunesThatFitStringSize := stringformat.GetMaxCharactersThatFitInStringSizeReverse(arrayOfRunesAvailableToPrint, textboxEntry.Width)
		// LogInfo(fmt.Sprintf("v: %d x: %d off: %d fit: %d, aval: %s", textboxEntry.ViewportXLocation, textboxEntry.CursorXLocation, maxViewportWidthAvaliable, numberOfRunesThatFitStringSize, string(arrayOfRunesAvailableToPrint)))
		textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation - numberOfRunesThatFitStringSize + 1
	}
}

/*
UpdateKeyboardEventTextboxWithString allows you to process a string of characters as keyboard input. In addition,
the following information should be noted:

- Processes each character in the string as a separate keystroke.
- Maintains all textbox functionality like highlighting and cursor movement.
- Returns true if a screen update is required.
*/
func (shared *textboxType) UpdateKeyboardEventTextboxWithString(keystroke string) {
	for _, currentCharacter := range keystroke {
		shared.UpdateKeyboardEvent([]rune{currentCharacter})
	}
}

/*
UpdateKeyboardEventTextboxWithCommands allows you to process a list of command strings. In addition,
the following information should be noted:

- Processes each command string as a separate keystroke.
- Useful for programmatically inserting text or executing commands.
- Returns true if a screen update is required.
*/
func (shared *textboxType) UpdateKeyboardEventTextboxWithCommands(keystroke ...string) {
	for _, currentCommand := range keystroke {
		shared.UpdateKeyboardEvent([]rune(currentCommand))
	}
}

func (shared *textboxType) UpdateKeyboardEventManually(layerAlias string, textboxAlias string, keystroke []rune) bool {
	isScreenUpdateRequired := false
	keystrokeAsString := string(keystroke)
	textboxEntry := Textboxes.Get(layerAlias, textboxAlias)

	// Store old cursor position for highlight updates
	oldCursorX := textboxEntry.CursorXLocation
	oldCursorY := textboxEntry.CursorYLocation

	if IsShiftPressed() {
		if !textboxEntry.IsHighlightModeToggled {
			// Start new highlight when toggling on
			textboxEntry.IsHighlightModeToggled = true
			textboxEntry.IsHighlightActive = true
			textboxEntry.HighlightStartX = oldCursorX
			textboxEntry.HighlightStartY = oldCursorY

		}
	} else {
		textboxEntry.IsHighlightModeToggled = false
	}

	// Handle cursor movement and text modification
	switch keystrokeAsString {
	case "left", "shift+left":
		if textboxEntry.IsHighlightModeToggled == false {
			textboxEntry.IsHighlightActive = false
		}
		textboxEntry.CursorXLocation--
		if textboxEntry.CursorXLocation < 0 {
			if textboxEntry.CursorYLocation > 0 {
				textboxEntry.CursorYLocation--
				textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
			} else {
				textboxEntry.CursorXLocation = 0
			}
		}
		isScreenUpdateRequired = true

	case "right", "shift+right":
		if textboxEntry.IsHighlightModeToggled == false {
			textboxEntry.IsHighlightActive = false
		}
		textboxEntry.CursorXLocation++
		if textboxEntry.CursorXLocation >= len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
			if textboxEntry.CursorYLocation < len(textboxEntry.TextData)-1 {
				textboxEntry.CursorYLocation++
				textboxEntry.CursorXLocation = 0
			} else {
				textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
			}
		}
		isScreenUpdateRequired = true

	case "up", "shift+up":
		if textboxEntry.IsHighlightModeToggled == false {
			textboxEntry.IsHighlightActive = false
		}
		textboxEntry.CursorYLocation--
		if textboxEntry.CursorYLocation < 0 {
			textboxEntry.CursorYLocation = 0
		}
		if textboxEntry.CursorXLocation >= len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
			textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
		}
		isScreenUpdateRequired = true

	case "down", "shift+down":
		if textboxEntry.IsHighlightModeToggled == false {
			textboxEntry.IsHighlightActive = false
		}
		textboxEntry.CursorYLocation++
		if textboxEntry.CursorYLocation >= len(textboxEntry.TextData) {
			textboxEntry.CursorYLocation = len(textboxEntry.TextData) - 1
		}
		if textboxEntry.CursorXLocation >= len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
			textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
		}
		isScreenUpdateRequired = true

	case "home", "shift+home":
		if textboxEntry.IsHighlightModeToggled == false {
			textboxEntry.IsHighlightActive = false
		}
		textboxEntry.CursorXLocation = 0
		isScreenUpdateRequired = true

	case "end", "shift+end":
		if textboxEntry.IsHighlightModeToggled == false {
			textboxEntry.IsHighlightActive = false
		}
		textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
		isScreenUpdateRequired = true

	case "pgup", "shift+pgup":
		if textboxEntry.IsHighlightModeToggled == false {
			textboxEntry.IsHighlightActive = false
		}
		textboxEntry.CursorYLocation = textboxEntry.CursorYLocation - textboxEntry.Height
		if textboxEntry.CursorYLocation < 0 {
			textboxEntry.CursorYLocation = 0
		}
		isScreenUpdateRequired = true

	case "pgdn", "shift+pgdn":
		if textboxEntry.IsHighlightModeToggled == false {
			textboxEntry.IsHighlightActive = false
		}
		textboxEntry.CursorYLocation = textboxEntry.CursorYLocation + textboxEntry.Height
		if textboxEntry.CursorYLocation >= len(textboxEntry.TextData) {
			textboxEntry.CursorYLocation = len(textboxEntry.TextData) - 1
		}
		isScreenUpdateRequired = true

	case "delete", "shift+delete":
		if textboxEntry.IsHighlightActive {
			// Delete all highlighted text
			shared.deleteHighlightedText(textboxEntry)
			textboxEntry.IsHighlightActive = false
		} else {
			// Normal delete behavior
			shared.deleteCharacterUsingAbsoluteCoordinates(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		}
		isScreenUpdateRequired = true

	case "backspace", "backspace2", "shift+backspace", "shift+backspace2":
		if textboxEntry.IsHighlightActive {
			// Delete all highlighted text
			shared.deleteHighlightedText(textboxEntry)
			textboxEntry.IsHighlightActive = false
		} else {
			// Normal backspace behavior
			shared.backspaceCharacterUsingRelativeCoordinates(textboxEntry)
		}
		isScreenUpdateRequired = true
	case "enter":
		if textboxEntry.IsHighlightModeToggled == false {
			textboxEntry.IsHighlightActive = false
		}
		shared.moveTextAfterCursorToNextLine(textboxEntry, textboxEntry.CursorYLocation)
		isScreenUpdateRequired = true

	default:
		if len(keystroke) == 1 { // If a regular char is entered
			shared.insertCharacterUsingAbsoluteCoordinates(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation, []rune(keystrokeAsString)[0])
			isScreenUpdateRequired = true
		}
	}

	// Update highlight end position if highlight mode is toggled on
	if textboxEntry.IsHighlightActive {
		textboxEntry.HighlightEndX = textboxEntry.CursorXLocation
		textboxEntry.HighlightEndY = textboxEntry.CursorYLocation
		isScreenUpdateRequired = true
	}

	// Update cursor position and viewport
	shared.updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
	shared.updateViewport(textboxEntry)
	shared.setTextboxMaxScrollBarValues(layerAlias, textboxAlias)
	shared.updateScrollbarBasedOnTextboxViewport(layerAlias, textboxAlias)
	return isScreenUpdateRequired
}

/*
UpdateKeyboardEvent allows you to process keyboard input for a textbox. In addition,
the following information should be noted:

- Handles all keyboard events including cursor movement and text editing.
- Manages text highlighting and selection.
- Returns true if a screen update is required.
*/
func (shared *textboxType) UpdateKeyboardEvent(keystroke []rune) bool {
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	if focusedControlType != constants.CellTypeTextbox || !Textboxes.IsExists(focusedLayerAlias, focusedControlAlias) {
		return false
	}
	return shared.UpdateKeyboardEventManually(focusedLayerAlias, focusedControlAlias, keystroke)
}

/*
deleteHighlightedText allows you to delete the currently highlighted text. In addition,
the following information should be noted:

- Removes all text within the highlight range.
- Handles both single-line and multi-line highlights.
- Updates the cursor position to the start of the deleted text.
*/
func (shared *textboxType) deleteHighlightedText(textboxEntry *types.TextboxEntryType) {
	// Determine the correct start and end positions for highlighting
	var highlightStartX, highlightEndX, highlightStartY, highlightEndY int

	// If cursor is to the left of the start position
	if textboxEntry.CursorYLocation < textboxEntry.HighlightStartY ||
		(textboxEntry.CursorYLocation == textboxEntry.HighlightStartY &&
			textboxEntry.CursorXLocation < textboxEntry.HighlightStartX) {
		// Cursor is before the highlight start, so swap the positions
		highlightStartX = textboxEntry.CursorXLocation
		highlightStartY = textboxEntry.CursorYLocation
		highlightEndX = textboxEntry.HighlightStartX
		highlightEndY = textboxEntry.HighlightStartY
	} else {
		// Cursor is after or at the highlight start
		highlightStartX = textboxEntry.HighlightStartX
		highlightStartY = textboxEntry.HighlightStartY
		highlightEndX = textboxEntry.CursorXLocation
		highlightEndY = textboxEntry.CursorYLocation
	}

	// Ensure we don't exceed array bounds
	if highlightStartY >= len(textboxEntry.TextData) {
		highlightStartY = len(textboxEntry.TextData) - 1
	}
	if highlightEndY >= len(textboxEntry.TextData) {
		highlightEndY = len(textboxEntry.TextData) - 1
	}

	// If the highlight is on a single line
	if highlightStartY == highlightEndY {
		// Ensure we don't exceed line bounds
		if highlightStartX >= len(textboxEntry.TextData[highlightStartY]) {
			highlightStartX = len(textboxEntry.TextData[highlightStartY]) - 1
		}
		if highlightEndX >= len(textboxEntry.TextData[highlightStartY]) {
			highlightEndX = len(textboxEntry.TextData[highlightStartY]) - 1
		}

		// Delete the highlighted portion of the line
		line := textboxEntry.TextData[highlightStartY]
		if highlightStartX < len(line) {
			if highlightEndX+1 < len(line) {
				textboxEntry.TextData[highlightStartY] = append(line[:highlightStartX], line[highlightEndX+1:]...)
			} else {
				textboxEntry.TextData[highlightStartY] = line[:highlightStartX]
			}
		}
	} else {
		// Multi-line highlight
		// Create a new slice to hold the result
		newTextData := make([][]rune, 0, len(textboxEntry.TextData))

		// Add lines before the highlight
		if highlightStartY > 0 {
			newTextData = append(newTextData, textboxEntry.TextData[:highlightStartY]...)
		}

		// Handle the first line of the highlight
		if highlightStartX > 0 && highlightStartX < len(textboxEntry.TextData[highlightStartY]) {
			newTextData = append(newTextData, textboxEntry.TextData[highlightStartY][:highlightStartX])
		}

		// Handle the last line of the highlight
		if highlightEndY < len(textboxEntry.TextData) && highlightEndX+1 < len(textboxEntry.TextData[highlightEndY]) {
			newTextData = append(newTextData, textboxEntry.TextData[highlightEndY][highlightEndX+1:])
		}

		// Add lines after the highlight
		if highlightEndY+1 < len(textboxEntry.TextData) {
			newTextData = append(newTextData, textboxEntry.TextData[highlightEndY+1:]...)
		}

		// If we ended up with no lines, add a blank line
		if len(newTextData) == 0 {
			newTextData = append(newTextData, []rune{' '})
		}

		textboxEntry.TextData = newTextData
	}

	// Move cursor to the start of the deleted text
	textboxEntry.CursorXLocation = highlightStartX
	textboxEntry.CursorYLocation = highlightStartY

	// Ensure we have at least one line with a space character
	if len(textboxEntry.TextData) == 0 {
		textboxEntry.TextData = append(textboxEntry.TextData, []rune{' '})
	}

	// Ensure the cursor position is valid
	if textboxEntry.CursorYLocation >= len(textboxEntry.TextData) {
		textboxEntry.CursorYLocation = len(textboxEntry.TextData) - 1
	}
	if textboxEntry.CursorXLocation >= len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
		textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
	}

	// Ensure each line ends with a space character for the cursor
	for i := 0; i < len(textboxEntry.TextData); i++ {
		if len(textboxEntry.TextData[i]) == 0 || textboxEntry.TextData[i][len(textboxEntry.TextData[i])-1] != ' ' {
			textboxEntry.TextData[i] = append(textboxEntry.TextData[i], ' ')
		}
	}

	// Turn off highlighting mode
	textboxEntry.IsHighlightActive = false
}

/*
updateMouseEvent allows you to process mouse events for a textbox. In addition,
the following information should be noted:

- Handles mouse clicks for cursor positioning.
- Manages text selection with mouse drag.
- Returns true if a screen update is required.
*/
func (shared *textboxType) updateMouseEvent() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	// If your clicking on a text box and not in the drag and drop event state.
	if buttonPressed != 0 && characterEntry.AttributeEntry.CellType == constants.CellTypeTextbox &&
		eventStateMemory.stateId != constants.EventStateDragAndDropScrollbar &&
		eventStateMemory.stateId != constants.EventStateDragAndDrop && // Add check for layer drag and drop
		Textboxes.IsExists(layerAlias, characterEntry.AttributeEntry.CellControlAlias) {
		textboxEntry := Textboxes.Get(layerAlias, characterEntry.AttributeEntry.CellControlAlias)
		shared.updateCursor(textboxEntry, characterEntry.AttributeEntry.CellControlId, characterEntry.AttributeEntry.CellControlLocation)
		shared.updateViewport(textboxEntry)
		shared.setTextboxMaxScrollBarValues(layerAlias, characterEntry.AttributeEntry.CellControlAlias)
		shared.updateScrollbarBasedOnTextboxViewport(layerAlias, characterEntry.AttributeEntry.CellControlAlias)
		setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, characterEntry.AttributeEntry.CellType)
		isUpdateRequired = true
		return isUpdateRequired
	}
	// If you are dragging and dropping, then update the scroll bars as needed.
	if buttonPressed != 0 && (eventStateMemory.stateId == constants.EventStateDragAndDropScrollbar ||
		characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar) {
		isMatchFound := false
		for _, currentTextBoxEntry := range Textboxes.GetAllEntries(layerAlias) {
			textboxEntry := currentTextBoxEntry
			hScrollBarEntry := ScrollBars.Get(layerAlias, textboxEntry.HorizontalScrollbarAlias)
			vScrollBarEntry := ScrollBars.Get(layerAlias, textboxEntry.VerticalScrollbarAlias)
			if textboxEntry.ViewportXLocation != hScrollBarEntry.ScrollValue {
				textboxEntry.ViewportXLocation = hScrollBarEntry.ScrollValue
				isUpdateRequired = true
			}
			if textboxEntry.ViewportYLocation != vScrollBarEntry.ScrollValue {
				textboxEntry.ViewportYLocation = vScrollBarEntry.ScrollValue
				isUpdateRequired = true
			}
			if isControlCurrentlyFocused(layerAlias, textboxEntry.HorizontalScrollbarAlias, constants.CellTypeScrollbar) ||
				isControlCurrentlyFocused(layerAlias, textboxEntry.VerticalScrollbarAlias, constants.CellTypeScrollbar) {
				isMatchFound = true
				break // If the current scrollbar being dragged and dropped matches, don't process more dropdowns.
			}
		}
		if isMatchFound {
			return isUpdateRequired
		}
	}
	return false
}
