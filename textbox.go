package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
	"math"
	"strings"
)

type TextboxInstanceType struct {
	layerAlias  string
	textboxAlias string
}

func (shared *TextboxInstanceType) setText(text string) bool {
	textData := strings.Split(text, "\n")
	textboxEntry := memory.TextboxMemory[shared.layerAlias][shared.textboxAlias]
	for _, text := range textData {
		textboxEntry.TextData = append(textboxEntry.TextData, stringformat.GetRunesFromString(text))
	}
	SetTextboxMaxScrollBarValues(shared.layerAlias, shared.textboxAlias)
	return false
}

func (shared *TextboxInstanceType) setViewport(xLocation int, yLocation int) bool {
	textboxEntry := memory.TextboxMemory[shared.layerAlias][shared.textboxAlias]
	textboxEntry.ViewportXLocation = xLocation
	textboxEntry.ViewportYLocation = yLocation
	return false
}

func getTextboxClickCoordinates(cellId int, tableWidth int) (int, int){
	xLocation := cellId % tableWidth
	yLocation := math.Floor(float64(cellId) / float64(tableWidth))
	return xLocation, int(yLocation)
}

func insertCharacterUsingRelativeCoordinates(textboxEntry *memory.TextboxEntryType, xLocation int, yLocation int, characterToInsert rune){
	stringDataSuffixAfterInsert := textboxEntry.TextData[yLocation][xLocation:len(textboxEntry.TextData[yLocation])]
	textboxEntry.TextData[yLocation] = append([]rune{}, textboxEntry.TextData[yLocation][:xLocation]...)
	textboxEntry.TextData[yLocation] = append(textboxEntry.TextData[yLocation][:xLocation], characterToInsert)
	textboxEntry.TextData[yLocation] = append(textboxEntry.TextData[yLocation], stringDataSuffixAfterInsert...)
	textboxEntry.CursorXLocation++
	// We need to detect if the cursor is at the end of the visible textbox viewing area, so we can move the
	// viewport. However, we only do this check if the remaining amount of text to print on a line is greater
	// than the viewport width.
	if textboxEntry.ViewportXLocation < len(textboxEntry.TextData[yLocation]) - textboxEntry.Width {
		// Get the number of runes we could potentially print.
		arrayOfRunesAvaliableToPrint := textboxEntry.TextData[yLocation][textboxEntry.ViewportXLocation:textboxEntry.ViewportXLocation + textboxEntry.Width]
		// Get the number of runes we can actually print, since some wide runes take up more space on screen.
		arrayOfRunesThatFitStringSize := stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvaliableToPrint, textboxEntry.Width)
		// If amount of runes that fit is equal to width, we just do a normal calculation.
		if len(arrayOfRunesThatFitStringSize) == textboxEntry.Width {
			// If cursor is greater than viewport + printable text, then increase viewport location.
			if textboxEntry.CursorXLocation >= textboxEntry.ViewportXLocation + len(arrayOfRunesThatFitStringSize) {
				textboxEntry.ViewportXLocation++
			}
		} else {
			// If the amount of runes available to print is less than the textbox width, then some runes are wide and take up
			// more space than normal. In this case, we do the same calculation but minus 1 since runes are wide and
			// the last displayable character is a blank (so start checking a space early).
			if textboxEntry.CursorXLocation >= textboxEntry.ViewportXLocation + len(arrayOfRunesThatFitStringSize) - 1 {
				textboxEntry.ViewportXLocation++
			}
		}
	}
}

func deleteCharacterUsingRelativeCoordinates(textboxEntry *memory.TextboxEntryType, xLocation int, yLocation int){
	if yLocation >= len(textboxEntry.TextData){
		return
	}
	if xLocation >= len(textboxEntry.TextData[yLocation]) -1 {
		if len(textboxEntry.TextData) - 1 == yLocation && (len(textboxEntry.TextData[yLocation]) <= 1 || xLocation >= len(textboxEntry.TextData[yLocation])-1) {
			return
		}
		logInfo(fmt.Sprintf("x: %d, len: %d \n", xLocation, len(textboxEntry.TextData[yLocation])-1))
		textboxEntry.TextData[yLocation] = textboxEntry.TextData[yLocation][:len(textboxEntry.TextData[yLocation]) - 1]
		textboxEntry.TextData[yLocation] = append(textboxEntry.TextData[yLocation], textboxEntry.TextData[yLocation + 1]...)
		copy(textboxEntry.TextData[yLocation + 1:], textboxEntry.TextData[yLocation + 2:])
		textboxEntry.TextData = textboxEntry.TextData[:len(textboxEntry.TextData) - 1]
		return
	}
	stringDataSuffixAfterInsert := textboxEntry.TextData[yLocation][xLocation+1:len(textboxEntry.TextData[yLocation])]
	textboxEntry.TextData[yLocation] = append([]rune{}, textboxEntry.TextData[yLocation][:xLocation]...)
	textboxEntry.TextData[yLocation] = append(textboxEntry.TextData[yLocation], stringDataSuffixAfterInsert...)
}

func insertLine(textData [][]rune, lineData []rune, index int) [][]rune {
	textData = append(textData[:index+1], textData[index:]...)
	textData[index] = lineData
	/*
	textData = append(textData, []rune{'z', 'z'})
	copy(textData[index:], textData[index+2:])
	logInfo(fmt.Sprintf("%d\n", len(textData)-1))
	 */
	return textData

	//textData = append(textData[:index], lineData)
}

func insertLineUsingRelativeCoordinates(textboxEntry *memory.TextboxEntryType, xLocation int, yLocation int){
	textboxEntry.TextData = insertLine(textboxEntry.TextData, []rune{' '}, yLocation)
	textboxEntry.CursorYLocation++
	logInfo(fmt.Sprintf("out %d\n", len(textboxEntry.TextData)))
	//textboxEntry.TextData = append(textboxEntry.TextData, []rune{' '})
	//copy(textboxEntry.TextData[yLocation:], textboxEntry.TextData[yLocation+1:])
}

func setCursorLocation(layerAlias string, textboxAlias string, cellControlId int, cellControlLocation int) {
	textboxEntry := memory.GetTextbox(layerAlias, textboxAlias)
	// If your Y selection is not on a text area.
	if cellControlLocation >= len(textboxEntry.TextData) {
		textboxEntry.CursorYLocation = len(textboxEntry.TextData) - 1
		textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
		// If your cursor xLocation is greater than the remaining viewport width, get how much space it actually
		// takes to print the remaining text and set your viewport accordingly.
		if textboxEntry.CursorXLocation >= textboxEntry.ViewportXLocation + textboxEntry.Width {
			maxVisible := stringformat.GetMaxCharactersThatFitInStringSize(textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.ViewportXLocation:], textboxEntry.Width)
			textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation - len(maxVisible) + 1
		} else if textboxEntry.CursorXLocation < textboxEntry.ViewportXLocation {
			// If your cursor xLocation is less than the current viewport position, calculate how much space it takes
			// to actually print a viewport width of space before your cursor, and set the viewport accordingly.
			var maxVisible []rune
			// If the text to print is less than the textbox width, start at 0.
			if textboxEntry.CursorXLocation - textboxEntry.Width < 0 {
				maxVisible = stringformat.GetMaxCharactersThatFitInStringSize(textboxEntry.TextData[textboxEntry.CursorYLocation][0:], textboxEntry.Width)
			} else {
				 maxVisible = stringformat.GetMaxCharactersThatFitInStringSize(textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.CursorXLocation - textboxEntry.Width:], textboxEntry.Width)
			}
			textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation - len(maxVisible) + 1
			logInfo(fmt.Sprintf("D: %d, x: %d, v: %s\n", textboxEntry.ViewportXLocation, textboxEntry.CursorXLocation, string(maxVisible)))
			if textboxEntry.ViewportXLocation < 0 {
				textboxEntry.ViewportXLocation = 0
			}
		}
		if textboxEntry.CursorYLocation >= textboxEntry.ViewportYLocation + textboxEntry.Height {
			textboxEntry.ViewportYLocation = textboxEntry.CursorYLocation
		}
	} else if cellControlId == constants.NullCellControlId {
		// If your Y location is fine, but x location is not on a text layer
		textboxEntry.CursorYLocation = cellControlLocation
		textboxEntry.CursorXLocation = len(textboxEntry.TextData[cellControlLocation]) - 1
		if textboxEntry.CursorXLocation < textboxEntry.ViewportXLocation {
			textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation - textboxEntry.Width + 1
			if textboxEntry.ViewportXLocation < 0 {
				textboxEntry.ViewportXLocation = 0
			}
		}
	} else {
		textboxEntry.CursorXLocation = cellControlId
		textboxEntry.CursorYLocation = cellControlLocation
		logInfo("HIT3\n")
	}
}
func updateScrollbarBasedOnTextboxViewport(layerAlias string, textboxAlias string) {
	textboxEntry := memory.GetTextbox(layerAlias, textboxAlias)
	horizontalScrollbarEntry := memory.GetScrollBar(layerAlias, textboxEntry.HorizontalScrollbarAlias)
	horizontalScrollbarEntry.ScrollValue = textboxEntry.ViewportXLocation
	computeScrollBarHandlePositionByScrollValue(layerAlias, textboxEntry.HorizontalScrollbarAlias)
	verticalScrollbarEntry := memory.GetScrollBar(layerAlias, textboxEntry.VerticalScrollbarAlias)
	verticalScrollbarEntry.ScrollValue = textboxEntry.ViewportYLocation
	computeScrollBarHandlePositionByScrollValue(layerAlias, textboxEntry.VerticalScrollbarAlias)
}

func getMaxHorizontalTextValue(layerAlias string, textboxAlias string) int {
	textboxEntry := memory.TextboxMemory[layerAlias][textboxAlias]
	maxHorizontalValue := 0
	for _, currentLine := range textboxEntry.TextData {
		lengthOfLine := stringformat.GetWidthOfRunesWhenPrinted(currentLine)
		over := lengthOfLine - len(currentLine)
		// logInfo(fmt.Sprintf("done: %d\n", over))
		if lengthOfLine > maxHorizontalValue {
			maxHorizontalValue = lengthOfLine - (over / 2)
		}
	}
	return maxHorizontalValue
}

func SetTextboxMaxScrollBarValues(layerAlias string, textboxAlias string) {
	textboxEntry := memory.TextboxMemory[layerAlias][textboxAlias]
	maxVerticalValue := len(textboxEntry.TextData)
	maxHorizontalValue := getMaxHorizontalTextValue(layerAlias, textboxAlias)
	hScrollBarEntry := memory.GetScrollBar(layerAlias, textboxEntry.HorizontalScrollbarAlias)
	vScrollBarEntry := memory.GetScrollBar(layerAlias, textboxEntry.VerticalScrollbarAlias)
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

func AddTextbox(layerAlias string, textboxAlias string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isBorderDrawn bool) TextboxInstanceType {
	memory.AddTextbox(layerAlias, textboxAlias, styleEntry, xLocation, yLocation, width, height, isBorderDrawn)
	textboxEntry := memory.GetTextbox(layerAlias, textboxAlias)
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
	memory.AddScrollBar(layerAlias, textboxEntry.HorizontalScrollbarAlias, styleEntry, hScrollbarXLocation, hScrollbarYLocation, hScrollbarWidth, 0, 0, 1, true)
	memory.AddScrollBar(layerAlias, textboxEntry.VerticalScrollbarAlias, styleEntry, vScrollbarXLocation, vScrollbarYLocation, vScrollbarHeight, 0, 0, 1, false)
	var textboxInstance TextboxInstanceType
	textboxInstance.layerAlias = layerAlias
	textboxInstance.textboxAlias = textboxAlias
	return textboxInstance
}

func DeleteTextbox(layerAlias string, textboxAlias string) {
	memory.DeleteTextbox(layerAlias, textboxAlias)
}


func drawTextboxesOnLayer(layerEntry memory.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for currentKey := range memory.TextboxMemory[layerAlias] {
		drawTextbox(&layerEntry, currentKey)
	}
}

func drawTextbox (layerEntry *memory.LayerEntryType, textboxAlias string) {
	t := memory.GetTextbox(layerEntry.LayerAlias, textboxAlias)
	attributeEntry := memory.NewAttributeEntry()
	attributeEntry.ForegroundColor = t.StyleEntry.TextboxForegroundColor
	attributeEntry.BackgroundColor = t.StyleEntry.TextboxBackgroundColor
	attributeEntry.CellControlAlias = textboxAlias
	if t.IsBorderDrawn {
		drawBorder(layerEntry, t.StyleEntry, attributeEntry, t.XLocation -1, t.YLocation - 1, t.Width + 2, t.Height + 2, false)
	}
	attributeEntry.CellType = constants.CellTypeTextbox
	fillArea(layerEntry, attributeEntry," ", t.XLocation, t.YLocation, t.Width, t.Height, t.ViewportYLocation)
	attributeEntry.CellControlAlias = textboxAlias
	for currentLine :=0 ; currentLine < t.Height; currentLine++ {
		var arrayOfRunes []rune
		if t.ViewportYLocation + currentLine < len(t.TextData) && t.ViewportYLocation + currentLine >= 0 {
			arrayOfRunes = t.TextData[t.ViewportYLocation+currentLine]
			if t.ViewportXLocation < len(arrayOfRunes) && t.ViewportXLocation >= 0 {
				if t.ViewportXLocation + t.Width < len(arrayOfRunes) {
					arrayOfRunes = arrayOfRunes[t.ViewportXLocation : t.ViewportXLocation + t.Width]
				} else {
					arrayOfRunes = arrayOfRunes[t.ViewportXLocation:]
				}
			} else {
				// If scrolled too far right and there are no column text to print, just show blanks.
				// If scrolled too far left (negative value) then show blanks. Note: This case should never happen really.
				arrayOfRunes = []rune{}
			}
			//arrayOfRunes = stringformat.GetFormattedRuneArray(arrayOfRunes, t.Width, constants.AlignmentLeft)
			arrayOfRunes = stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunes, t.Width)
			printControlText(layerEntry, textboxAlias, t.StyleEntry, attributeEntry, t.XLocation, t.YLocation + currentLine, arrayOfRunes, t.ViewportYLocation + currentLine, t.ViewportXLocation, t.CursorXLocation, t.CursorYLocation)
		} else {
			// If scrolled too far down and there are no more rows to print, just show blanks.
			// If scrolled too far up and there are no rows to print, just print blanks. Note: This case should never happen really.
			//arrayOfRunes = stringformat.GetFormattedRuneArray([]rune{}, t.Width, constants.AlignmentLeft)
			printControlText(layerEntry, textboxAlias, t.StyleEntry, attributeEntry, t.XLocation, t.YLocation + currentLine, arrayOfRunes,  t.ViewportYLocation + currentLine, t.ViewportXLocation, t.CursorXLocation, t.CursorYLocation)
		}
	}
}

func printControlText(layerEntry *memory.LayerEntryType, textboxAlias string, styleEntry memory.TuiStyleEntryType, attributeEntry memory.AttributeEntryType, xLocation int, yLocation int, arrayOfRunes []rune, controlYLocation int, startingControlId int, cursorXLocation int, cursorYLocation int) {
	currentControlId := startingControlId
	currentXOffset := 0
	for _, currentCharacter := range arrayOfRunes {
		attributeEntry.CellControlId = currentControlId
		attributeEntry.CellControlLocation = controlYLocation
		if isControlCurrentlyFocused(layerEntry.LayerAlias, textboxAlias, constants.CellTypeTextbox) {
			if cursorXLocation == currentControlId && cursorYLocation == controlYLocation {
				attributeEntry.ForegroundColor = styleEntry.CursorForegroundColor
				attributeEntry.BackgroundColor = styleEntry.CursorBackgroundColor
			} else {
				attributeEntry.ForegroundColor = styleEntry.TextboxForegroundColor
				attributeEntry.BackgroundColor = styleEntry.TextboxBackgroundColor
			}
		}
		printLayer(layerEntry, attributeEntry, xLocation + currentXOffset, yLocation, []rune{currentCharacter})
		if stringformat.IsRuneCharacterWide(currentCharacter) {
			// If we find a wide character, we add a blank space with the same ID as the previous
			// character so the next printed character doesn't get covered by the wide one.
			currentXOffset++
			// we also increment the control ID for logInfo(fmt.Sprintf("zzz %d, %d \n", textboxEntry.CursorXLocation, textboxEntry.CursorYLocation))the blank space, so we can keep the control ID (location of screen)
			// accurate. TODO: We need a new variable to keep track of proper indexes so when they select a cell the right char is picked.
			printLayer(layerEntry, attributeEntry, xLocation + currentXOffset, yLocation, []rune{' '})
			currentXOffset++
		} else {
			currentXOffset++
		}
		currentControlId++
	}
}
func updateViewport2(textboxEntry *memory.TextboxEntryType) {
	// If cursor xLocation is greater than the line, make it max line length
	if textboxEntry.CursorXLocation >= len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
		textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
		textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation - textboxEntry.Width
		if textboxEntry.ViewportXLocation < 0 {
			textboxEntry.ViewportXLocation = 0
		}
	}
	// if cursor yLocation is less than viewport y, then update viewport yLocation
	if textboxEntry.CursorYLocation < textboxEntry.ViewportYLocation {
		textboxEntry.ViewportYLocation = textboxEntry.CursorYLocation
	}
	// We need to detect if the cursor is at the end of the visible textbox viewing area, so we can move the
	// viewport. However, we only do this check if the remaining amount of text to print on a line is greater
	// than the viewport width.
	if textboxEntry.ViewportXLocation < len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - textboxEntry.Width {
		// Get the number of runes we could potentially print.
		arrayOfRunesAvaliableToPrint := textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.ViewportXLocation:textboxEntry.ViewportXLocation + textboxEntry.Width]
		// Get the number of runes we can actually print, since some wide runes take up more space on screen.
		arrayOfRunesThatFitStringSize := stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvaliableToPrint, textboxEntry.Width)
		// If amount of runes that fit is equal to width, we just do a normal calculation.
		if len(arrayOfRunesThatFitStringSize) == textboxEntry.Width {
			// If cursor is greater than viewport + printable text, then increase viewport location.
			if textboxEntry.CursorXLocation >= textboxEntry.ViewportXLocation + len(arrayOfRunesThatFitStringSize) {
				textboxEntry.ViewportXLocation++
			}
		} else {
			// If the amount of runes available to print is less than the textbox width, then some runes are wide and take up
			// more space than normal. In this case, we do the same calculation but minus 1 since runes are wide and
			// the last displayable character is a blank (so start checking a space early).
			if textboxEntry.CursorXLocation >= textboxEntry.ViewportXLocation + len(arrayOfRunesThatFitStringSize) - 1 {
				textboxEntry.ViewportXLocation++
			}
		}
	}
	if textboxEntry.CursorYLocation >= textboxEntry.ViewportYLocation + textboxEntry.Height {
		textboxEntry.ViewportYLocation++
	}
	if textboxEntry.ViewportXLocation > textboxEntry.CursorXLocation {
		textboxEntry.ViewportXLocation--
	}
	if textboxEntry.ViewportYLocation > textboxEntry.CursorYLocation {
		textboxEntry.ViewportXLocation--
	}
}

func updateCursor(textboxEntry *memory.TextboxEntryType, xLocation int, yLocation int) {
	textboxEntry.CursorXLocation = xLocation
	textboxEntry.CursorYLocation = yLocation
	// If yLocation is less than text data bounds.
	if textboxEntry.CursorYLocation < 0 {
		textboxEntry.CursorYLocation = 0
	}
	// If yLocation is greater than column data bounds.
	if textboxEntry.CursorYLocation > len(textboxEntry.TextData) - 1 {
		textboxEntry.CursorYLocation = len(textboxEntry.TextData) - 1
	}
	// If x is less than or greater than row data bounds, set it to max length of current line
	if textboxEntry.CursorXLocation < 0 || textboxEntry.CursorXLocation > len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1 {
		textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
	}
	//logInfo(fmt.Sprintf("ix: %d, iy: %d cx: %d, cy: %d  \n", xLocation, yLocation, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation))
}

func updateViewport(textboxEntry *memory.TextboxEntryType) {
	// If cursor yLocation is greater than viewport width, set viewport to one width before cursor.
	if textboxEntry.CursorYLocation >= textboxEntry.ViewportYLocation + textboxEntry.Height {
		textboxEntry.ViewportYLocation = textboxEntry.CursorYLocation - textboxEntry.Height + 1
	}
	// If cursor yLocation is less than viewport, set viewport to one width before cursor.
	if textboxEntry.CursorYLocation < textboxEntry.ViewportYLocation {
		textboxEntry.ViewportYLocation = textboxEntry.CursorYLocation
	}
	// If cursor yLocation is less <= 0, set the viewport to zero.
	if textboxEntry.CursorYLocation <= 0 {
		textboxEntry.ViewportYLocation = 0
	}

	// If cursor xLocation is less than viewport.
	if textboxEntry.CursorXLocation < textboxEntry.ViewportXLocation {
		if textboxEntry.CursorXLocation - textboxEntry.Width < 0 {
			//arrayOfRunesAvailableToPrint = textboxEntry.TextData[textboxEntry.CursorYLocation][:textboxEntry.CursorXLocation]
			//arrayOfRunesThatFitStringSize = stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvailableToPrint, textboxEntry.Width)
			textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation
		} else {
			//arrayOfRunesAvailableToPrint = textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.CursorXLocation:textboxEntry.CursorXLocation + textboxEntry.Width]
			//arrayOfRunesThatFitStringSize = stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvailableToPrint, textboxEntry.Width)
			textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation
		}
	}

	var arrayOfRunesAvailableToPrint []rune
	var arrayOfRunesThatFitStringSize []rune

	if textboxEntry.ViewportXLocation + textboxEntry.Width >= len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
		arrayOfRunesAvailableToPrint = textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.ViewportXLocation:]
		arrayOfRunesThatFitStringSize = stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvailableToPrint, len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - textboxEntry.ViewportXLocation)
	} else {
		arrayOfRunesAvailableToPrint = textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.ViewportXLocation:textboxEntry.ViewportXLocation + textboxEntry.Width]
		arrayOfRunesThatFitStringSize = stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvailableToPrint, textboxEntry.Width)
	}

	if textboxEntry.CursorXLocation >= textboxEntry.ViewportXLocation + len(arrayOfRunesThatFitStringSize){
		textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation - len(arrayOfRunesThatFitStringSize) + 1
	}

	/*
	if textboxEntry.ViewportXLocation + textboxEntry.Width < len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
		arrayOfRunesAvailableToPrint = textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.ViewportXLocation:textboxEntry.ViewportXLocation + textboxEntry.Width]
	} else {
		arrayOfRunesAvailableToPrint = textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.ViewportXLocation:]
	}
	arrayOfRunesThatFitStringSize = stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvailableToPrint, textboxEntry.Width)
	if textboxEntry.CursorXLocation >= textboxEntry.ViewportXLocation + len(arrayOfRunesThatFitStringSize) {
		textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation - len(arrayOfRunesThatFitStringSize)
	}
	*/
	//logInfo(fmt.Sprint("ZZZZZZZZZZZZZZZZZZZZZZZZZ"))

	/*
	// If cursor xLocation is greater than the line, make it max line length
	if textboxEntry.CursorXLocation >= len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
		// If the text for the line is shorter than width of viewport, just set viewport to 0.
		if len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1 < textboxEntry.Width {
			textboxEntry.ViewportXLocation = 0
			return
		}
		textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
		textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation - textboxEntry.Width
		// Get the number of runes we could potentially print.
		arrayOfRunesAvailableToPrint := textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.ViewportXLocation:textboxEntry.ViewportXLocation + textboxEntry.Width]
		// Get the number of runes we can actually print, since some wide runes take up more space on screen.
		arrayOfRunesThatFitStringSize := stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvailableToPrint, textboxEntry.Width)
		// If the size are the same, no wide runes exist.
		if len(arrayOfRunesThatFitStringSize) != textboxEntry.Width {
			textboxEntry.ViewportXLocation = textboxEntry.ViewportXLocation - len(arrayOfRunesThatFitStringSize)
		}
		if textboxEntry.ViewportXLocation < 0 {
			textboxEntry.ViewportXLocation = 0
		}
	}
	*/

}


func updateKeyboardEventTextbox(keystroke string) bool {
	isScreenUpdateRequired := true
	focusedLayerAlias := eventStateMemory.currentlyFocusedControl.layerAlias
	focusedControlAlias := eventStateMemory.currentlyFocusedControl.controlAlias
	focusedControlType := eventStateMemory.currentlyFocusedControl.controlType
	if focusedControlType != constants.CellTypeTextbox {
		return false
	}
	textboxEntry := memory.GetTextbox(focusedLayerAlias, focusedControlAlias)
	if len(keystroke) == 1 { // If a regular char is entered.
		insertCharacterUsingRelativeCoordinates(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation, []rune(keystroke)[0])
		SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "delete" {
		deleteCharacterUsingRelativeCoordinates(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "enter" {
		insertLineUsingRelativeCoordinates(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "home" {

	}
	if keystroke == "end" {

	}
	if keystroke == "backspace" || keystroke == "backspace2" {

	}
	if keystroke == "left" {
		textboxEntry.CursorXLocation--
		if textboxEntry.CursorXLocation < 0 {
			textboxEntry.CursorXLocation = 0
		}
		if textboxEntry.CursorXLocation >= len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
			textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation])
		}
		updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateViewport(textboxEntry)
		//setCursorLocation(focusedLayerAlias, focusedControlAlias, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "right" {
		textboxEntry.CursorXLocation++
		if textboxEntry.CursorXLocation >= len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
			textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
		}
		if textboxEntry.CursorXLocation >= len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
			textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation])
		}
		updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateViewport(textboxEntry)
		//setCursorLocation(focusedLayerAlias, focusedControlAlias, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "up" {
		textboxEntry.CursorYLocation--
		if textboxEntry.CursorYLocation < 0 {
			textboxEntry.CursorYLocation = 0
		}
		updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateViewport(textboxEntry)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "down" {
		textboxEntry.CursorYLocation++
		if textboxEntry.CursorYLocation >= len(textboxEntry.TextData) {
			textboxEntry.CursorYLocation = len(textboxEntry.TextData) - 1
		}
		updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateViewport(textboxEntry)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	return isScreenUpdateRequired
}

func updateMouseEventTextbox() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	if buttonPressed != 0 && characterEntry.AttributeEntry.CellType == constants.CellTypeTextbox && eventStateMemory.stateId != constants.EventStateDragAndDropScrollBar {
		textboxEntry := memory.GetTextbox(layerAlias, characterEntry.AttributeEntry.CellControlAlias)
		updateCursor(textboxEntry, characterEntry.AttributeEntry.CellControlId, characterEntry.AttributeEntry.CellControlLocation)
		updateViewport(textboxEntry)
		//setCursorLocation(layerAlias, characterEntry.AttributeEntry.CellControlAlias, characterEntry.AttributeEntry.CellControlId, characterEntry.AttributeEntry.CellControlLocation)
		SetTextboxMaxScrollBarValues(layerAlias, characterEntry.AttributeEntry.CellControlAlias)
		updateScrollbarBasedOnTextboxViewport(layerAlias, characterEntry.AttributeEntry.CellControlAlias)
		setFocusedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, characterEntry.AttributeEntry.CellType)
		isUpdateRequired = true
	}
	if buttonPressed != 0 && (eventStateMemory.stateId == constants.EventStateDragAndDropScrollBar ||
		characterEntry.AttributeEntry.CellType == constants.CellTypeScrollBar) {
			isMatchFound := false
			for currentKey := range memory.TextboxMemory[layerAlias] {
				textboxEntry := memory.GetTextbox(layerAlias, currentKey)
				hScrollBarEntry := memory.GetScrollBar(layerAlias, textboxEntry.HorizontalScrollbarAlias)
				vScrollBarEntry := memory.GetScrollBar(layerAlias, textboxEntry.VerticalScrollbarAlias)
				if textboxEntry.ViewportXLocation != hScrollBarEntry.ScrollValue {
					textboxEntry.ViewportXLocation = hScrollBarEntry.ScrollValue
					isUpdateRequired = true
				}
				if textboxEntry.ViewportYLocation != vScrollBarEntry.ScrollValue {
					textboxEntry.ViewportYLocation = vScrollBarEntry.ScrollValue
					isUpdateRequired = true
				}
				if isControlCurrentlyFocused(layerAlias, textboxEntry.HorizontalScrollbarAlias, constants.CellTypeScrollBar) ||
					isControlCurrentlyFocused(layerAlias, textboxEntry.VerticalScrollbarAlias, constants.CellTypeScrollBar) {
					isMatchFound = true
					break; // If the current scrollbar being dragged and dropped matches, don't process more dropdowns.
				}
			}
			if isMatchFound {
				return isUpdateRequired
			}
	}
	/*
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	controlAlias := characterEntry.AttributeEntry.CellControlAlias
	if characterEntry.AttributeEntry.CellType == constants.CellTypeTextbox && characterEntry.AttributeEntry.CellControlId != constants.NullCellId {
		_, _, previousButtonPressed, _ := memory.GetPreviousMouseStatus()
		if buttonPressed != 0 && previousButtonPressed == 0 {
			checkboxEntry := memory.GetTextbox(layerAlias, controlAlias)
			if checkboxEntry.IsSelected {
				checkboxEntry.IsSelected = false
			} else {
				checkboxEntry.IsSelected = true
			}
			return isUpdateRequired
		}
	}
	return isUpdateRequired
	*/
	return false
}
