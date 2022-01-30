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

func (shared *textFieldInstanceType) GetValue() string {
	mapEntry, isFound := memory.TextFieldMemory[shared.layerAlias][shared.textFieldAlias]
	if !isFound {
		panic(fmt.Sprintf("The text field '%s' under the text layer '%s' does not exist!", shared.textFieldAlias, shared.layerAlias))
	}
	return mapEntry.CurrentValue
}

func AddTextField(layerAlias string, textFieldAlias string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, width int, maxLengthAllowed int, IsPasswordProtected bool, defaultValue string) textFieldInstanceType {
	memory.AddTextField(layerAlias, textFieldAlias, styleEntry, xLocation, yLocation, width, maxLengthAllowed, IsPasswordProtected, defaultValue)
	var textFieldInstance textFieldInstanceType
	textFieldInstance.layerAlias = layerAlias
	textFieldInstance.textFieldAlias = textFieldAlias
	return textFieldInstance
}

func DeleteTextField(layerAlias string, textFieldAlias string) {
	memory.DeleteTextField(layerAlias, textFieldAlias)
}

/*
drawButtonsOnLayer allows you to draw all buttons on a given text layer
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
		Locate(10,10)
		PrintLayer("Layer1","   ")
		Locate(10,10)
		PrintLayer("Layer1", textFieldEntry.CurrentValue)
		drawCursor(&layerEntry, textFieldEntry.StyleEntry, currentKey, characterUnderCursor, textFieldEntry.XLocation, textFieldEntry.YLocation, textFieldEntry.CursorPosition, false)
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
	attributeEntry.CellAlias = textFieldAlias
	attributeEntry.CellTypeId = cursorPosition
	var arrayOfRunes []rune
	arrayOfRunes = append(arrayOfRunes, characterUnderCursor)
	printLayer(layerEntry, attributeEntry, xLocation+cursorPosition, yLocation, arrayOfRunes)
	if isMovementBackwards {
		printLayer(layerEntry, attributeEntry, xLocation+cursorPosition+1, yLocation, arrayOfRunes)
	}
}

/*
drawInputString allows you to draw a string for an input field. This is
different than regular printing, since input fields are usually
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
	attributeEntry.CellAlias = textFieldAlias
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
		attributeEntry.CellTypeId = currentRuneIndex
		printLayer(layerEntry, attributeEntry, xLocation + currentRuneIndex, yLocation, []rune{arrayOfRunes[currentRuneIndex]})
	}
}

