package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/types"
)

var TextStyles *memory.MemoryManager[types.TextCellStyleEntryType]

func init() {
	TextStyles = memory.NewMemoryManager[types.TextCellStyleEntryType]()
}

func GetTextStyle(textStyleAlias string) *types.TextCellStyleEntryType {
	// Use the generic memory manager to retrieve the text style entry
	if !TextStyles.IsExists(textStyleAlias) {
		panic(fmt.Sprintf("The requested text style with alias '%s' could not be returned since it does not exist.", textStyleAlias))
	}
	return TextStyles.Get(textStyleAlias)
}

func GetTextStyleAsAttributeEntry(textStyleAlias string) types.AttributeEntryType {
	// Use the generic memory manager to retrieve the text style entry
	if !TextStyles.IsExists(textStyleAlias) {
		panic(fmt.Sprintf("The requested text style with alias '%s' could not be returned since it does not exist.", textStyleAlias))
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
