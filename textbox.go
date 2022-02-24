package consolizer

import (
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
	textboxEntry.TextData = textData
	SetTextboxMaxScrollBarValues(shared.layerAlias, shared.textboxAlias)
	return false
}

func (shared *TextboxInstanceType) setViewport(xLocation int, yLocation int) bool {
	textboxEntry := memory.TextboxMemory[shared.layerAlias][shared.textboxAlias]
	textboxEntry.ViewportXLocation = xLocation
	textboxEntry.ViewportYLocation = yLocation
	return false
}

func getTextboxClickCoordinates(cellId int, tableWidth int, tableHeight int) (int, int){
	xLocation := cellId % tableHeight
	yLocation := math.Floor(float64(cellId) / float64(tableWidth))
	return xLocation, int(yLocation)
}


func SetTextboxMaxScrollBarValues(layerAlias string, textboxAlias string) {
	textboxEntry := memory.TextboxMemory[layerAlias][textboxAlias]
	maxVerticalValue := len(textboxEntry.TextData)
	maxHorizontalValue := 0
	for _, currentLine := range textboxEntry.TextData {
		if len(currentLine) > maxHorizontalValue {
			maxHorizontalValue = len(currentLine)
		}
	}
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
		textboxEntry := memory.GetTextbox(layerAlias, currentKey)
		drawTextbox(&layerEntry, currentKey, textboxEntry.StyleEntry, textboxEntry.TextData, textboxEntry.XLocation, textboxEntry.YLocation, textboxEntry.Width, textboxEntry.Height, textboxEntry.ViewportXLocation, textboxEntry.ViewportYLocation, textboxEntry.IsBorderDrawn)
	}
}

func drawTextbox (layerEntry *memory.LayerEntryType, textboxAlias string, styleEntry memory.TuiStyleEntryType, textData []string, xLocation int, yLocation int, width int, height int, viewportXPosition int, viewportYPosition int, isBorderDrawn bool) {
	localStyleEntry := memory.NewTuiStyleEntry(&styleEntry)
	attributeEntry := memory.NewAttributeEntry()
	attributeEntry.ForegroundColor = localStyleEntry.TextboxForegroundColor
	attributeEntry.BackgroundColor = localStyleEntry.TextboxBackgroundColor
	attributeEntry.CellType = constants.CellTypeTextbox
	attributeEntry.CellControlAlias = textboxAlias
	if isBorderDrawn {
		fillArea(layerEntry, attributeEntry," ", xLocation - 1,yLocation - 1, width + 2, height + 2)
		drawBorder(layerEntry, styleEntry, attributeEntry, xLocation -1, yLocation - 1, width + 2, height + 2, false)
	}
	currentControlId := 0
	for currentLine :=0 ; currentLine < height; currentLine++ {
		var arrayOfRunes []rune
		if viewportYPosition + currentLine < len(textData) && viewportYPosition + currentLine >= 0 {
			arrayOfRunes = stringformat.GetRunesFromString(textData[viewportYPosition+currentLine])
			if viewportXPosition < len(arrayOfRunes) && viewportXPosition >= 0 {
				if viewportXPosition+width < len(arrayOfRunes) {
					arrayOfRunes = arrayOfRunes[viewportXPosition : viewportXPosition+width]
				} else {
					arrayOfRunes = arrayOfRunes[viewportXPosition:]
				}
			} else {
				// If scrolled too far right and there are no column text to print, just show blanks.
				// If scrolled too far left (negative value) then show blanks. Note: This case should never happen really.
				arrayOfRunes = []rune{}
			}
			arrayOfRunes = stringformat.GetFormattedRuneArray(arrayOfRunes, width, constants.AlignmentLeft)
			printControlText(layerEntry, attributeEntry, xLocation, yLocation + currentLine, arrayOfRunes, currentControlId)
		} else {
			// If scrolled too far down and there are no more rows to print, just show blanks.
			// If scrolled too far up and there are no rows to print, just print blanks. Note: This case should never happen really.
			arrayOfRunes = stringformat.GetFormattedRuneArray([]rune{}, width, constants.AlignmentLeft)
			printControlText(layerEntry, attributeEntry, xLocation, yLocation + currentLine, arrayOfRunes, currentControlId)
		}
		currentControlId = currentControlId + len(arrayOfRunes)
	}
}

func printControlText(layerEntry *memory.LayerEntryType, attributeEntry memory.AttributeEntryType, xLocation int, yLocation int, arrayOfRunes []rune, startingControlId int) {
	currentControlId := startingControlId
	currentXOffset := 0
	for _, currentCharacter := range arrayOfRunes {
		attributeEntry.CellControlId = currentControlId
		printLayer(layerEntry, attributeEntry, xLocation + currentXOffset, yLocation, []rune{currentCharacter})
		if stringformat.IsRuneCharacterWide(currentCharacter) {
			// If we find a wide character, we add a blank space with the same ID as the previous
			// character so the next printed character doesn't get covered by the wide one.
			currentXOffset++
			printLayer(layerEntry, attributeEntry, xLocation + currentXOffset, yLocation, []rune{' '})
			currentXOffset++
		} else {
			currentXOffset++
		}
		currentControlId++
	}
	printLayer(layerEntry, attributeEntry, xLocation + currentXOffset + 1, yLocation, []rune{'*'})
}

func updateMouseEventTextbox() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
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
