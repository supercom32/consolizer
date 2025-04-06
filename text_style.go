package consolizer

import (
	"fmt"
	"supercom32.net/consolizer/internal/memory"
	"supercom32.net/consolizer/types"
)

var TextStyles = memory.NewControlMemoryManager[types.TextCellStyleEntryType]()

func GetTextStyle(textStyleAlias string) *types.TextCellStyleEntryType {
	// Use the generic memory manager to retrieve the text style entry
	textStyleEntry := TextStyles.Get("", textStyleAlias)
	if textStyleEntry == nil {
		panic(fmt.Sprintf("The requested text style with alias '%s' could not be returned since it does not exist.", textStyleAlias))
	}
	return textStyleEntry
}

func GetTextStyleAsAttributeEntry(textStyleAlias string) types.AttributeEntryType {
	// Use the generic memory manager to retrieve the text style entry
	textStyleEntry := TextStyles.Get("", textStyleAlias)
	if textStyleEntry == nil {
		panic(fmt.Sprintf("The requested text style with alias '%s' could not be returned since it does not exist.", textStyleAlias))
	}

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
	return TextStyles.Get("", textStyleAlias) != nil
}
