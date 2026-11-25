package consolizer

import (
	"fmt"
	"strings"

	"github.com/supercom32/consolizer/stringformat"

	"github.com/supercom32/consolizer/types"
)

/*
printDialog is a method which writes text to the terminal screen via a typewriter effect. This is useful for video
games or other applications that may require printing text in a dialog box. With a delay, one character is printed
at a time, the screen is redrawn and the delay is waited out, and the remaining characters are printed at once if
the dialog is skippable and Enter is pressed or a mouse button is held. In addition, the following should be noted:

  - The caller must not hold uiStateMutex. Each character is printed under the lock, and the lock is released while
    waiting out the delay and checking for the skip key, so input keeps being handled while the dialog prints and
    the key that skips it can be delivered at all.

Example:

	printDialog(layer, attr, 0, 0, 20, 50, true, "Hello {{red}}World{{/}}")
*/
func printDialog(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, widthOfLineInCharacters int, printDelayInMilliseconds int, isSkipable bool, textToPrint string) {
	// If no delay is requested, just print the text normally and return
	if printDelayInMilliseconds <= 0 {
		func() {
			defer lockUIState()()
			layer.print(layerEntry, attributeEntry, xLocation, yLocation, stringformat.GetRunesFromString(textToPrint), widthOfLineInCharacters, true)
		}()
		return
	}

	// Apply typewriter effect with character-by-character printing
	isPrintDelaySkipped := false
	cursorXLocation := xLocation
	cursorYLocation := yLocation

	// For markup
	currentAttributeEntry := attributeEntry

	// Convert runes to string for markup processing
	arrayOfRunes := stringformat.GetRunesFromString(textToPrint)

	for currentCharacterIndex := 0; currentCharacterIndex < len(arrayOfRunes); currentCharacterIndex++ {
		// The character is printed under the UI lock. False means nothing was printed that needs a pause.
		isPauseRequired := func() bool {
			defer lockUIState()()
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
				return false
			}

			// Handle markup
			if currentCharacter == '{' && currentCharacterIndex+1 < len(arrayOfRunes) && arrayOfRunes[currentCharacterIndex+1] == '{' {
				tagStartIndex := currentCharacterIndex + 2
				foundClosing := false
				for j := tagStartIndex; j < len(arrayOfRunes)-1; j++ {
					if arrayOfRunes[j] == '}' && arrayOfRunes[j+1] == '}' {
						tagContent := string(arrayOfRunes[tagStartIndex:j])
						currentAttributeEntry = getDialogAttributeEntry(tagContent, attributeEntry)
						currentCharacterIndex = j + 1
						foundClosing = true
						break
					}
				}
				if foundClosing {
					return false
				}
			}

			// Print the character
			if cursorXLocation >= 0 && cursorXLocation < layerEntry.Width && cursorYLocation >= 0 && cursorYLocation < layerEntry.Height {
				columnsConsumed := putRune(layerEntry.CharacterMemory, cursorXLocation, cursorYLocation, currentCharacter, currentAttributeEntry, layerEntry.Width, layerEntry.Height)
				if columnsConsumed == 1 && stringformat.IsRuneCharacterWide(currentCharacter) {
					// A wide rune whose placeholder would fall past the right edge.
					cursorXLocation++
					if widthOfLineInCharacters > 0 {
						return false
					}
				} else {
					cursorXLocation += columnsConsumed - 1
				}
			}

			cursorXLocation++

			// Handle line wrapping for basic printing
			if cursorXLocation >= layerEntry.Width {
				if widthOfLineInCharacters > 0 {
					return false
				}
			}

			return true
		}()
		if !isPauseRequired {
			continue
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
getAttributeTagContent is a method which allows you to obtain the style name held inside a full markup tag, such as
"red" from "{{red}}" or "/" from "{{/}}", so it can be passed to getDialogAttributeEntry. A tag that is not wrapped in
"{{" and "}}", including an empty string, yields an empty string, which getDialogAttributeEntry treats as the
default attribute entry.

Example:

	tagContent := getAttributeTagContent("{{red}}")
*/
func getAttributeTagContent(attributeTag string) string {
	if len(attributeTag) < 4 || !strings.HasPrefix(attributeTag, "{{") || !strings.HasSuffix(attributeTag, "}}") {
		return ""
	}
	return attributeTag[2 : len(attributeTag)-2]
}

/*
getDialogAttributeEntry is a method which allows you to obtain the attribute entry a markup tag selects, given the
default attribute entry of the text being printed. A tag naming a registered text style yields the default entry with
that style's visual fields applied, while the closing tag "/", an empty tag and any tag naming an unknown text style
yield the default entry unchanged. In addition, the following should be noted:

  - Only the visual fields are taken from the text style. The cell and control fields of the default entry, such as
    CellType, CellControlAlias and CellControlId, are always kept, so styled cells still belong to the control that
    printed them and remain reachable by mouse hit-testing.

Example:

	attributeEntry := getDialogAttributeEntry("red", defaultAttributeEntry)
*/
func getDialogAttributeEntry(attributeTag string, defaultAttributeEntry types.AttributeEntryType) types.AttributeEntryType {
	if attributeTag == "" || attributeTag == "/" || !IsTextStyleExists(attributeTag) {
		return defaultAttributeEntry
	}
	return applyTextStyleVisuals(defaultAttributeEntry, GetTextStyleAsAttributeEntry(attributeTag))
}

/*
applyTextStyleVisuals is a method which allows you to combine two attribute entries, returning a copy of the base
entry whose colours, bold, underline, reverse, blink, italic, transparency and alpha fields are replaced by those of
the style entry. Every other field, namely CellType, CellControlAlias, CellControlId, CellControlLocation, CellUserId
and CellUserAlias, is kept from the base entry.

Example:

	attributeEntry := applyTextStyleVisuals(defaultAttributeEntry, GetTextStyleAsAttributeEntry("red"))
*/
func applyTextStyleVisuals(baseAttributeEntry types.AttributeEntryType, styleAttributeEntry types.AttributeEntryType) types.AttributeEntryType {
	attributeEntry := baseAttributeEntry
	attributeEntry.ForegroundColor = styleAttributeEntry.ForegroundColor
	attributeEntry.BackgroundColor = styleAttributeEntry.BackgroundColor
	attributeEntry.IsBold = styleAttributeEntry.IsBold
	attributeEntry.IsUnderlined = styleAttributeEntry.IsUnderlined
	attributeEntry.IsReversed = styleAttributeEntry.IsReversed
	attributeEntry.IsBlinking = styleAttributeEntry.IsBlinking
	attributeEntry.IsItalic = styleAttributeEntry.IsItalic
	attributeEntry.IsBackgroundTransparent = styleAttributeEntry.IsBackgroundTransparent
	attributeEntry.IsForegroundTransparent = styleAttributeEntry.IsForegroundTransparent
	attributeEntry.ForegroundAlphaValue = styleAttributeEntry.ForegroundAlphaValue
	attributeEntry.BackgroundAlphaValue = styleAttributeEntry.BackgroundAlphaValue
	return attributeEntry
}

/*
GetNonMarkupText is a method which gets a string without {{...}} markup control characters in it. This is
useful for calculating words and word wrapping without control characters messing it up. If no terminating }} can be
found, then the {{someTagText is printed out as regular text. In addition, the following should be noted:

- Used by word width calculation functions to get accurate text length.

- Handles nested tags and unclosed tags appropriately.

Example:
    plainText := GetNonMarkupText("Hello {{red}}World{{/}}")
*/
func GetNonMarkupText(textString string) string {
	var result []rune
	runes := []rune(textString)
	var i int

	for i < len(runes) {
		// Check for opening tag
		if i+1 < len(runes) && runes[i] == '{' && runes[i+1] == '{' {
			// Look for closing tag
			tagStart := i
			i += 2 // Skip past "{{"
			foundClosing := false

			for j := i; j < len(runes)-1; j++ {
				if runes[j] == '}' && runes[j+1] == '}' {
					// Found closing tag, skip the entire tag
					i = j + 2
					foundClosing = true
					break
				}
			}

			// If no closing tag found, treat opening tag as regular text
			if !foundClosing {
				result = append(result, '{', '{')
				i = tagStart + 2
			}
		} else {
			// Regular character, add to result
			result = append(result, runes[i])
			i++
		}
	}

	return string(result)
}

/*
printMarkup is a method which allows you to write text to the terminal screen with word wrapping and attribute tags,
like printDialog without the typewriter effect and printing delay, and then redraw the screen. A
widthOfLineInCharacters of zero or less turns word wrapping off. The text is printed by the same routine as the
layer's own markup printing, so tags, wide characters and wrapping behave identically on both paths.

Example:

	printMarkup(layerEntry, attributeEntry, 0, 0, 20, "Hello {{red}}World{{/}}")
*/
func printMarkup(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, widthOfLineInCharacters int, stringToPrint string) {
	if xLocation < 0 || xLocation > layerEntry.Width || yLocation < 0 || yLocation > layerEntry.Height {
		safeSttyPanic(fmt.Sprintf("The specified location (%d, %d) is out of bounds for the layer with a size of (%d, %d).", xLocation, yLocation, layerEntry.Width, layerEntry.Height))
	}
	layer.printMarkup(layerEntry, attributeEntry, xLocation, yLocation, widthOfLineInCharacters, stringToPrint)
	UpdateDisplay(false)
}
