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

func GetTextStyle(textStyleAlias string) *types.TextCellStyleEntryType {
	// Use the generic memory manager to retrieve the text style entry
	if !TextStyles.IsExists(textStyleAlias) {
		safeSttyPanic(fmt.Sprintf("The requested text style with alias '%s' could not be returned since it does not exist.", textStyleAlias))
	}
	return TextStyles.Get(textStyleAlias)
}

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

func IsTextStyleExists(textStyleAlias string) bool {
	// Use the generic memory manager to check if the text style exists
	return TextStyles.IsExists(textStyleAlias)
}

/*
CalculateStringLengthWithoutMarkup allows you to calculate the length of a string
without counting markup tags. Markup tags are sequences surrounded by "{{" and "}}".
This is useful for calculating the visual length of a string that contains markup.
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
