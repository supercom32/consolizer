package consolizer

import (
	"fmt"

	"supercom32.net/consolizer/internal/stringformat"
	"supercom32.net/consolizer/types"
)

/*
printDialog allows you to write text to the terminal screen via a typewriter
effect. This is useful for video games or other applications that may
require printing text in a dialog box.
*/
func printDialog(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, widthOfLineInCharacters int, printDelayInMilliseconds int, isSkipable bool, textToPrint string) {
	if xLocation < 0 || xLocation > layerEntry.Width || yLocation < 0 || yLocation > layerEntry.Height {
		panic(fmt.Sprintf("The specified location (%d, %d) is out of bounds for the layer with a size of (%d, %d).", xLocation, yLocation, layerEntry.Width, layerEntry.Height))
	}
	var isPrintDelaySkipped bool
	arrayOfRunes := stringformat.GetRunesFromString(textToPrint)
	layerWidth := layerEntry.Width
	layerHeight := layerEntry.Height
	cursorXLocation := xLocation
	cursorYLocation := yLocation
	currentAttributeEntry := attributeEntry
	for currentCharacterIndex := 0; currentCharacterIndex < len(arrayOfRunes); currentCharacterIndex++ {
		currentCharacter := stringformat.GetSubString(textToPrint, currentCharacterIndex, 1)
		printLayer(layerEntry, currentAttributeEntry, cursorXLocation, cursorYLocation, []rune{arrayOfRunes[currentCharacterIndex]})
		cursorXLocation++
		lengthOfNextWord := 0
		if currentCharacter == " " {
			lengthOfNextWord = getLengthOfNextWord(textToPrint, currentCharacterIndex+1)
		}
		nextCharacter := stringformat.GetSubString(textToPrint, currentCharacterIndex+1, 1)
		if nextCharacter == "{" {
			attributeTag := getAttributeTag(textToPrint, currentCharacterIndex+1)
			currentAttributeEntry = getDialogAttributeEntry(attributeTag, attributeEntry)
			currentCharacterIndex += len(attributeTag)
		}
		if cursorXLocation+lengthOfNextWord-xLocation >= widthOfLineInCharacters || cursorXLocation+lengthOfNextWord >= layerWidth {
			cursorXLocation = xLocation
			cursorYLocation++
			if cursorYLocation >= layerHeight {
				cursorYLocation--
			}
		}
		if isSkipable == true {
			_, _, mouseButtonPressed, _ := GetMouseStatus()
			keyPressed := Inkey()
			if mouseButtonPressed != 0 || string(keyPressed) == "enter" {
				isPrintDelaySkipped = true
			}
		}
		if isPrintDelaySkipped == false {
			if printDelayInMilliseconds != 0 {
				SleepInMilliseconds(uint(printDelayInMilliseconds))
				UpdateDisplay(false)
			}
		}
	}
	UpdateDisplay(false)
}

/*
getAttributeTag allows you to obtain an attribute tag from a given text string.
Attributes are always surrounded by "{" and "}" characters.  In addition, the
following information should be noted:

- If no attribute tag could be detected at the given string location, then
an empty string will be returned instead.
*/
func getAttributeTag(stringToParse string, startingCharacterIndex int) string {
	var lengthOfAttributeTag int
	for currentCharacterIndex := startingCharacterIndex; currentCharacterIndex < len(stringToParse); currentCharacterIndex++ {
		lengthOfAttributeTag++
		if stringformat.GetSubString(stringToParse, currentCharacterIndex, 1) == "}" {
			return stringformat.GetSubString(stringToParse, startingCharacterIndex, lengthOfAttributeTag)
		}
	}
	return ""
}

/*
getDialogAttributeEntry allows you to obtain an attribute entry based on the
text style detected in your attribute tag. If no text style could be found
that matches the attribute tag provided, then the default attribute entry
will be returned instead.
*/
func getDialogAttributeEntry(attributeTag string, defaultAttributeEntry types.AttributeEntryType) types.AttributeEntryType {
	var attributeEntry types.AttributeEntryType
	if attributeTag != "" {
		textStyleAlias := stringformat.GetSubString(attributeTag, 1, len(attributeTag)-2)
		if IsTextStyleExists(textStyleAlias) {
			attributeEntry = GetTextStyleAsAttributeEntry(textStyleAlias)
			return attributeEntry
		}
	}
	return defaultAttributeEntry
}

/*
getLengthOfNextWord allows you to get the length of the next word at a given
position of a text string.
*/
func getLengthOfNextWord(stringToParse string, startingCharacterIndex int) int {
	var lengthOfNextWord int
	for currentCharacterIndex := startingCharacterIndex; currentCharacterIndex < len(stringToParse); currentCharacterIndex++ {
		if stringformat.GetSubString(stringToParse, currentCharacterIndex, 1) == " " {
			return lengthOfNextWord
		}
		lengthOfNextWord++
	}
	return lengthOfNextWord
}
