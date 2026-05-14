package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/types"
	"strings"
)

var TextStyles *memory.MemoryManager[types.TextCellStyleEntryType]

func init() {
	TextStyles = memory.NewMemoryManager[types.TextCellStyleEntryType]()
}

/*
GetTextStyle is a method which allows you to retrieve a text style entry by its alias.

:param textStyleAlias: The alias of the text style to retrieve.

:return: A pointer to the text style entry.

Example:

	style := GetTextStyle("myStyle")
*/
func GetTextStyle(textStyleAlias string) *types.TextCellStyleEntryType {
	// Use the generic memory manager to retrieve the text style entry
	if !TextStyles.IsExists(textStyleAlias) {
		safeSttyPanic(fmt.Sprintf("The requested text style with alias '%s' could not be returned since it does not exist.", textStyleAlias))
	}
	return TextStyles.Get(textStyleAlias)
}

/*
GetTextStyleAsAttributeEntry is a method which allows you to retrieve a text style entry and convert it to an attribute
entry.

:param textStyleAlias: The alias of the text style to retrieve.

:return: An attribute entry containing the text style's properties.

Example:

	attribute := GetTextStyleAsAttributeEntry("myStyle")
*/
func GetTextStyleAsAttributeEntry(textStyleAlias string) types.AttributeEntryType {
	// Use the generic memory manager to retrieve the text style entry
	if !TextStyles.IsExists(textStyleAlias) {
		safeSttyPanic(fmt.Sprintf("The requested text style with alias '%s' could not be returned since it does not exist.", textStyleAlias))
	}
	textStyleEntry := TextStyles.Get(textStyleAlias)

	// Convert to AttributeEntryType
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = textStyleEntry.ForegroundColor
	attributeEntry.BackgroundColor = textStyleEntry.BackgroundColor
	attributeEntry.IsBlinking = textStyleEntry.IsBlinking
	attributeEntry.IsItalic = textStyleEntry.IsItalic
	attributeEntry.IsReversed = textStyleEntry.IsReversed
	attributeEntry.IsUnderlined = textStyleEntry.IsUnderlined
	attributeEntry.IsBold = textStyleEntry.IsBold
	return attributeEntry
}

/*
IsTextStyleExists is a method which allows you to check if a text style exists in memory.

:param textStyleAlias: The alias of the text style to check.

:return: True if the text style exists, false otherwise.

Example:

	exists := IsTextStyleExists("myStyle")
*/
func IsTextStyleExists(textStyleAlias string) bool {
	// Use the generic memory manager to check if the text style exists
	return TextStyles.IsExists(textStyleAlias)
}

/*
CalculateStringLengthWithoutMarkup is a method which allows you to calculate the length of a string without counting
markup tags. In addition, the following should be noted:

- Markup tags are sequences surrounded by "{{" and "}}".

- This is useful for calculating the visual length of a string that contains markup.

:param text: The string to calculate the length of.

:return: The length of the string without markup tags.

Example:

	length := CalculateStringLengthWithoutMarkup("{{red}}Hello{{white}} World")
*/
func CalculateStringLengthWithoutMarkup(text string) int {
	length := 0
	i := 0

	for i < len(text) {
		// Look for a possible markup opening
		if i+1 < len(text) && text[i] == '{' && text[i+1] == '{' {
			// Look ahead for closing }}
			end := strings.Index(text[i+2:], "}}")
			if end != -1 {
				// Found closing tag — skip entire markup
				i += 2 + end + 2 // Skip {{tag}}
				continue
			} else {
				// No closing }} — treat "{{" as normal text
				length += 2
				i += 2
				continue
			}
		}

		// Regular character
		length++
		i++
	}

	return length
}
