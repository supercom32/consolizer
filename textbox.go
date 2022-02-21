package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
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
	hScrollBarEntry.MaxScrollValue = maxHorizontalValue
	vScrollBarEntry := memory.GetScrollBar(layerAlias, textboxEntry.VerticalScrollbarAlias)
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
	for currentLine :=0 ; currentLine < height; currentLine++ {
		var arrayOfRunes []rune
		if viewportYPosition + currentLine < len(textData) && viewportYPosition + currentLine >= 0 {
			arrayOfRunes = stringformat.GetRunesFromString(textData[viewportYPosition+currentLine])
			if viewportXPosition < len(arrayOfRunes) && viewportXPosition >= 0 {
				if viewportXPosition+width < len(arrayOfRunes) {
					arrayOfRunes = arrayOfRunes[viewportXPosition : viewportXPosition+width]
				} else {
					arrayOfRunes = arrayOfRunes[viewportXPosition:]
					arrayOfRunes = stringformat.GetRunesFromString(stringformat.GetFormattedString(string(arrayOfRunes), width, constants.AlignmentLeft))
				}
			} else {
				// If scrolled too far right and there are no column text to print, just show blanks.
				// If scrolled too far left (negative value) then show blanks. Note: This case should never happen really.
				arrayOfRunes = stringformat.GetRunesFromString(stringformat.GetFormattedString("", width, constants.AlignmentLeft))
			}
			printLayer(layerEntry, attributeEntry, xLocation, yLocation + currentLine, arrayOfRunes)
		} else {
			// If scrolled too far down and there are no more rows to print, just show blanks.
			// If scrolled too far up and there are no rows to print, just print blanks. Note: This case should never happen really.
			arrayOfRunes = stringformat.GetRunesFromString(stringformat.GetFormattedString("", width, constants.AlignmentLeft))
			printLayer(layerEntry, attributeEntry, xLocation, yLocation + currentLine, arrayOfRunes)
		}
	}
}

func updateMouseEventTextbox() bool {
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
