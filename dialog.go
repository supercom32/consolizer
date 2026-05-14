package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/stringformat"
	"strings"

	"github.com/supercom32/consolizer/types"
)

/*
printDialog is a method which allows you to write text to the terminal screen via a typewriter effect. This is useful
for video games or other applications that may require printing text in a dialog box.

:param layerEntry: The layer to print the dialog on.
:param attributeEntry: The default attribute entry to use for printing.
:param xLocation: The X coordinate to start printing at.
:param yLocation: The Y coordinate to start printing at.
:param widthOfLineInCharacters: The width of the line before wrapping occurs.
:param printDelayInMilliseconds: The delay between printing each character.
:param isSkipable: Whether the typewriter effect can be skipped by the user.
:param textToPrint: The text string to print (supports markup).

Example:

	printDialog(layer, attr, 0, 0, 20, 50, true, "Hello {{red}}World{{/}}")
*/
func printDialog(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, widthOfLineInCharacters int, printDelayInMilliseconds int, isSkipable bool, textToPrint string) {
	// If no delay is requested, just print the text normally and return
	if printDelayInMilliseconds <= 0 {
		arrayOfRunes := stringformat.GetRunesFromString(textToPrint)
		layer.print(layerEntry, attributeEntry, xLocation, yLocation, arrayOfRunes, widthOfLineInCharacters, true)
		return
	}

	// Apply typewriter effect with character-by-character printing
	isPrintDelaySkipped := false
	cursorXLocation := xLocation
	cursorYLocation := yLocation
	characterMemory := layerEntry.CharacterMemory

	// For markup
	currentAttributeEntry := attributeEntry

	// Convert runes to string for markup processing
	arrayOfRunes := stringformat.GetRunesFromString(textToPrint)
	textString := textToPrint

	for currentCharacterIndex := 0; currentCharacterIndex < len(arrayOfRunes); currentCharacterIndex++ {
		currentCharacter := arrayOfRunes[currentCharacterIndex]

		// Handle word wrapping if enabled
		if widthOfLineInCharacters > 0 && currentCharacter == ' ' {
			// Check if the word fits within the remaining space on the current line.
			var wordWidth int
			wordWidth = calculateWordWidth(arrayOfRunes, currentCharacterIndex, true)

			if cursorXLocation+wordWidth >= xLocation+widthOfLineInCharacters || cursorXLocation+wordWidth >= layerEntry.Width {
				// Word doesn't fit, move to the next line.
				cursorXLocation = xLocation
				cursorYLocation++
				if cursorYLocation >= layerEntry.Height {
					cursorYLocation--
				}
			}
		}

		// Skip the first blank space at the start of a line if one exists and word wrap is enabled
		if widthOfLineInCharacters > 0 && currentCharacter == ' ' && cursorXLocation == xLocation {
			continue
		}

		// Handle markup
		if currentCharacter == '{' && currentCharacterIndex+1 < len(arrayOfRunes) && arrayOfRunes[currentCharacterIndex+1] == '{' {
			tagStartIndex := currentCharacterIndex + 2
			tagEndRelIndex := strings.Index(textString[tagStartIndex:], "}}")

			if tagEndRelIndex != -1 {
				// Valid tag found
				tagEndIndex := tagStartIndex + tagEndRelIndex
				tagContent := textString[tagStartIndex:tagEndIndex]

				// Apply the new attribute
				currentAttributeEntry = getDialogAttributeEntry(tagContent, attributeEntry)

				// Skip over the full tag, including '{{' and '}}'
				currentCharacterIndex = tagEndIndex + 1 // We find the start of }}, so we need to add 1 to bypass the tag totally.
				continue
			}
		}

		// Print the character
		if cursorXLocation >= 0 && cursorXLocation < layerEntry.Width && cursorYLocation >= 0 && cursorYLocation < layerEntry.Height {
			originalBackgroundColor := characterMemory[cursorYLocation][cursorXLocation].AttributeEntry.BackgroundColor

			characterMemory[cursorYLocation][cursorXLocation].AttributeEntry = types.NewAttributeEntry(&currentAttributeEntry)
			characterMemory[cursorYLocation][cursorXLocation].Character = currentCharacter

			if stringformat.IsRuneCharacterWide(currentCharacter) {
				cursorXLocation++
				if cursorXLocation >= layerEntry.Width {
					if widthOfLineInCharacters > 0 {
						continue
					}
				}
				characterMemory[cursorYLocation][cursorXLocation].AttributeEntry = types.NewAttributeEntry(&currentAttributeEntry)
				characterMemory[cursorYLocation][cursorXLocation].Character = ' '
			}

			if characterMemory[cursorYLocation][cursorXLocation].AttributeEntry.IsBackgroundTransparent {
				characterMemory[cursorYLocation][cursorXLocation].AttributeEntry.BackgroundColor = originalBackgroundColor
			}
		}

		cursorXLocation++

		// Handle line wrapping for basic printing
		if cursorXLocation >= layerEntry.Width {
			if widthOfLineInCharacters > 0 {
				continue
			}
		}

		// Check for skip input
		if isSkipable {
			_, _, mouseButtonPressed, _ := GetMouseStatus()
			keyPressed := Inkey()
			if mouseButtonPressed != 0 || string(keyPressed) == "enter" {
				isPrintDelaySkipped = true
			}
		}

		// Apply delay unless skipped
		if !isPrintDelaySkipped && printDelayInMilliseconds > 0 {
			SleepInMilliseconds(uint(printDelayInMilliseconds))
			UpdateDisplay(false)
		}
	}

	// Final display update
	UpdateDisplay(false)
}

/*
getAttributeTag is a method which allows you to obtain an attribute tag from a given text string. Attributes are always
surrounded by "{{" and "}}" characters. In addition, the following should be noted:

- If no attribute tag could be detected at the given string location, then an empty string will be returned instead.

:param stringToParse: The text string to parse for an attribute tag.
:param startingCharacterIndex: The index in the string where the search should begin.

:return: The detected attribute tag, or an empty string if none is found.

Example:

	tag := getAttributeTag("Hello {{red}}World", 6)
*/
func getAttributeTag(stringToParse string, startingCharacterIndex int) string {
	var lengthOfAttributeTag int
	for currentCharacterIndex := startingCharacterIndex; currentCharacterIndex < len(stringToParse)-1; currentCharacterIndex++ {
		lengthOfAttributeTag++
		if stringformat.GetSubString(stringToParse, currentCharacterIndex, 2) == "}}" {
			return stringformat.GetSubString(stringToParse, startingCharacterIndex, lengthOfAttributeTag+1)
		}
	}
	// If we reach here, we didn't find the closing tag.
	// Return an empty string, but we'll handle this case differently in the print function.
	return ""
}

/*
getDialogAttributeEntry is a method which allows you to obtain an attribute entry based on the text style detected in
your attribute tag. In addition, the following should be noted:

- If no text style could be found that matches the attribute tag provided, then the default attribute entry will be.

:param attributeTag: The attribute tag string to evaluate.
:param defaultAttributeEntry: The attribute entry to return if no match is found.

:return: The matching attribute entry or the default attribute entry.

Example:

	attr := getDialogAttributeEntry("red", defaultAttr)
*/
func getDialogAttributeEntry(attributeTag string, defaultAttributeEntry types.AttributeEntryType) types.AttributeEntryType {
	var attributeEntry types.AttributeEntryType
	if attributeTag != "" {
		// Special case for the closing tag "/"
		if attributeTag == "/" {
			return defaultAttributeEntry
		}
		// Normal case for style tags
		if IsTextStyleExists(attributeTag) {
			attributeEntry = GetTextStyleAsAttributeEntry(attributeTag)
			return attributeEntry
		}
	}
	return defaultAttributeEntry
}

/*
getLengthOfNextWord is a method which allows you to get the length of the next word at a given position of a text
string. It ignores markup tags when calculating the length.

:param stringToParse: The text string to analyze.
:param startingCharacterIndex: The index in the string where the next word begins.

:return: The length of the next word in characters.

Example:

	length := getLengthOfNextWord("Hello {{red}}World", 0)
*/
func getLengthOfNextWord(stringToParse string, startingCharacterIndex int) int {
	// First, get the substring from the starting index to the end
	substring := stringToParse[startingCharacterIndex:]

	// Get the text without markup
	textWithoutMarkup := GetNonMarkupText(substring)

	// Now calculate the length of the next word in the text without markup
	var lengthOfNextWord int
	for currentCharacterIndex := 0; currentCharacterIndex < len(textWithoutMarkup); currentCharacterIndex++ {
		if stringformat.GetSubString(textWithoutMarkup, currentCharacterIndex, 1) == " " {
			return lengthOfNextWord
		}
		lengthOfNextWord++
	}
	return lengthOfNextWord
}

/*
GetNonMarkupText is a method which allows you to get a string without {{...}} markup control characters in it. This is
useful for calculating words and word wrapping without control characters messing it up. If no terminating }} can be
found, then the {{someTagText is printed out as regular text. In addition, the following should be noted:

- Used by word width calculation functions to get accurate text length.

- Handles nested tags and unclosed tags appropriately.

:param textString: The text string containing markup to be stripped.

:return: The text string without any markup tags.

Example:

	plainText := GetNonMarkupText("Hello {{red}}World{{/}}")
*/
func GetNonMarkupText(textString string) string {
	var result string
	var i int

	for i < len(textString) {
		// Check for opening tag
		if i+1 < len(textString) && textString[i:i+2] == "{{" {
			// Look for closing tag
			tagStart := i
			i += 2 // Skip past "{{"
			foundClosing := false

			for j := i; j < len(textString)-1; j++ {
				if textString[j:j+2] == "}}" {
					// Found closing tag, skip the entire tag
					i = j + 2
					foundClosing = true
					break
				}
			}

			// If no closing tag found, treat opening tag as regular text
			if !foundClosing {
				result += "{{"
				i = tagStart + 2
			}
		} else {
			// Regular character, add to result
			result += string(textString[i])
			i++
		}
	}

	return result
}

/*
printMarkup is a method which allows you to write text to the terminal screen with word wrapping and attribute tags.
This is similar to printDialog but without the typewriter effect and printing delay.

:param layerEntry: The layer to print the markup on.
:param attributeEntry: The default attribute entry to use for printing.
:param xLocation: The X coordinate to start printing at.
:param yLocation: The Y coordinate to start printing at.
:param widthOfLineInCharacters: The width of the line before wrapping occurs.
:param stringToPrint: The text string to print (supports markup).

Example:

	printMarkup(layer, attr, 0, 0, 20, "Hello {{red}}World{{/}}")
*/
func printMarkup(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, widthOfLineInCharacters int, stringToPrint string) {
	if xLocation < 0 || xLocation > layerEntry.Width || yLocation < 0 || yLocation > layerEntry.Height {
		safeSttyPanic(fmt.Sprintf("The specified location (%d, %d) is out of bounds for the layer with a size of (%d, %d).", xLocation, yLocation, layerEntry.Width, layerEntry.Height))
	}
	arrayOfRunes := stringformat.GetRunesFromString(stringToPrint)
	layerWidth := layerEntry.Width
	layerHeight := layerEntry.Height
	cursorXLocation := xLocation
	cursorYLocation := yLocation
	currentAttributeEntry := attributeEntry
	for currentCharacterIndex := 0; currentCharacterIndex < len(arrayOfRunes); currentCharacterIndex++ {
		currentCharacter := stringformat.GetSubString(stringToPrint, currentCharacterIndex, 1)
		printLayer(layerEntry, currentAttributeEntry, cursorXLocation, cursorYLocation, []rune{arrayOfRunes[currentCharacterIndex]})
		cursorXLocation++
		lengthOfNextWord := 0
		if currentCharacter == " " {
			lengthOfNextWord = getLengthOfNextWord(stringToPrint, currentCharacterIndex+1)
		}
		nextCharacter := stringformat.GetSubString(stringToPrint, currentCharacterIndex+1, 2)
		if nextCharacter == "{{" {
			attributeTag := getAttributeTag(stringToPrint, currentCharacterIndex+1)
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
	}
	UpdateDisplay(false)
}
