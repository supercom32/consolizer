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
	if index < len(textData) {
		textData = append(textData[:index+1], textData[index:]...)
		textData[index] = lineData
	} else {
		textData = append(textData, []rune {' '})
	}
	/*
	textData = append(textData, []rune{'z', 'z'})
	copy(textData[index:], textData[index+2:])
	logInfo(fmt.Sprintf("%d\n", len(textData)-1))
	 */
	return textData

	//textData = append(textData[:index], lineData)
}

func insertLineUsingRelativeCoordinates(textboxEntry *memory.TextboxEntryType, xLocation int, yLocation int){
	textboxEntry.TextData = insertLine(textboxEntry.TextData, []rune{' '}, yLocation+1)
	charactersToCopy := textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.CursorXLocation:]
	copyOfCharacters := make([]rune, len(textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.CursorXLocation:]))
	copy(copyOfCharacters, charactersToCopy)


	logInfo(fmt.Sprintf("1: %s\n", string(copyOfCharacters)))
	offsetValue := 0
	if textboxEntry.CursorXLocation == 0 {
		offsetValue = 0
	}
	textboxEntry.TextData[textboxEntry.CursorYLocation] = append(textboxEntry.TextData[textboxEntry.CursorYLocation][:textboxEntry.CursorXLocation - offsetValue], ' ')
	logInfo(fmt.Sprintf("1: %s\n", string(copyOfCharacters)))
	textboxEntry.CursorYLocation++
	textboxEntry.CursorXLocation = 0
	textboxEntry.TextData[textboxEntry.CursorYLocation] = make([]rune, len(copyOfCharacters))
	copy(textboxEntry.TextData[textboxEntry.CursorYLocation], copyOfCharacters)


	//textboxEntry.TextData = append(textboxEntry.TextData, []rune{' '})
	//copy(textboxEntry.TextData[yLocation:], textboxEntry.TextData[yLocation+1:])
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
	// If our cursor xLocation was jumped (due to NullCellControlId) or placed in an invalid xLocation spot greater than the length of our text line.
	// Move it to the end of the line.
	if textboxEntry.CursorXLocation == constants.NullCellControlId || textboxEntry.CursorXLocation > len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1 {
		textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - 1
	} else if textboxEntry.CursorXLocation < 0 {
		// If the cursor xLocation was scrolled to be out of lower bounds, just set the location to 0.
		textboxEntry.CursorXLocation = 0
	}
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
		isCursorJumped := false
		if textboxEntry.ViewportXLocation - textboxEntry.CursorXLocation > 2 {
			isCursorJumped = true
		}
		if textboxEntry.CursorXLocation - textboxEntry.Width < 0 {
			if isCursorJumped {
				textboxEntry.ViewportXLocation = 0
			} else {
				textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation
			}
		} else {
			if isCursorJumped {
				textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation - textboxEntry.Width
			} else {
				// If the cursor xLocation was just scrolled, scroll the viewport accordingly.
				textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation
			}
		}
	}

	var arrayOfRunesAvailableToPrint []rune
	var arrayOfRunesThatFitStringSize []rune

	if textboxEntry.ViewportXLocation + textboxEntry.Width >= len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
		arrayOfRunesAvailableToPrint = textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.ViewportXLocation:]
		arrayOfRunesThatFitStringSize = stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvailableToPrint, len(textboxEntry.TextData[textboxEntry.CursorYLocation]) - textboxEntry.ViewportXLocation)
		logInfo("HIT1\n")
	} else {
		arrayOfRunesAvailableToPrint = textboxEntry.TextData[textboxEntry.CursorYLocation][textboxEntry.ViewportXLocation:textboxEntry.ViewportXLocation + textboxEntry.Width]
		arrayOfRunesThatFitStringSize = stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunesAvailableToPrint, textboxEntry.Width)
		logInfo("HIT2\n")
	}
	// If the cursor xLocation >= the visible viewport xLocation + amount of text that can be visibly drawn, and if the cursor xLocation is != 1
	// (For the case of a single wide character on a line), then move the viewport, so it makes your cursor appear at the end of it.
	if textboxEntry.CursorXLocation >= textboxEntry.ViewportXLocation + len(arrayOfRunesThatFitStringSize) && textboxEntry.CursorXLocation != 1 {
		logInfo(fmt.Sprintf("c: %d, v: %d, l: %d\n", textboxEntry.CursorXLocation, textboxEntry.ViewportXLocation, len(arrayOfRunesThatFitStringSize)))
		textboxEntry.ViewportXLocation = textboxEntry.CursorXLocation - len(arrayOfRunesThatFitStringSize) + 1
	}
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
		updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateViewport(textboxEntry)
		SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "delete" {
		deleteCharacterUsingRelativeCoordinates(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateViewport(textboxEntry)
		SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "enter" {
		insertLineUsingRelativeCoordinates(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateViewport(textboxEntry)
		SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "home" {
		textboxEntry.CursorXLocation = 0
		updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateViewport(textboxEntry)
		SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "end" {
		textboxEntry.CursorXLocation = len(textboxEntry.TextData[textboxEntry.CursorYLocation])
		updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateViewport(textboxEntry)
		SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "pgup" {
		textboxEntry.CursorYLocation = textboxEntry.CursorYLocation - textboxEntry.Width
		updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateViewport(textboxEntry)
		SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "pgdn" {
		textboxEntry.CursorYLocation = textboxEntry.CursorYLocation + textboxEntry.Width
		updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateViewport(textboxEntry)
		SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
		updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
	}
	if keystroke == "backspace" || keystroke == "backspace2" {
		if textboxEntry.CursorXLocation != 0 {
			textboxEntry.CursorXLocation--
			updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
			deleteCharacterUsingRelativeCoordinates(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
			updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
			updateViewport(textboxEntry)
			SetTextboxMaxScrollBarValues(focusedLayerAlias, focusedControlAlias)
			updateScrollbarBasedOnTextboxViewport(focusedLayerAlias, focusedControlAlias)
		}
	}
	if keystroke == "left" {
		textboxEntry.CursorXLocation--
		updateCursor(textboxEntry, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
		updateViewport(textboxEntry)
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
