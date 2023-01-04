package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
)

type textStyleMemoryType struct {
	sync.Mutex
	Entries map[string]*types.TextCellStyleEntryType
}

var TextStyle textStyleMemoryType

func InitializeTextStyleMemory() {
	TextStyle.Entries = make(map[string]*types.TextCellStyleEntryType)
}

func AddTextStyle(textStyleAlias string, attributeEntry types.TextCellStyleEntryType) {
	TextStyle.Lock()
	defer func() {
		TextStyle.Unlock()
	}()
	TextStyle.Entries[textStyleAlias] = &attributeEntry
}

func GetTextStyle(textStyleAlias string) *types.TextCellStyleEntryType {
	TextStyle.Lock()
	defer func() {
		TextStyle.Unlock()
	}()
	if TextStyle.Entries[textStyleAlias] == nil {
		panic(fmt.Sprintf("The requested text style with alias '%s' could not be returned since it does not exist.", textStyleAlias))
	}
	return TextStyle.Entries[textStyleAlias]
}

func GetTextStyleAsAttributeEntry(textStyleAlias string) types.AttributeEntryType {
	TextStyle.Lock()
	defer func() {
		TextStyle.Unlock()
	}()
	textStyleEntry := TextStyle.Entries[textStyleAlias]
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

// DeleteTextStyle
func DeleteTextStyle(textStyleAlias string) {
	TextStyle.Lock()
	defer func() {
		TextStyle.Unlock()
	}()
	delete(TextStyle.Entries, textStyleAlias)
}

func IsTextStyleExists(textStyleAlias string) bool {
	TextStyle.Lock()
	defer func() {
		TextStyle.Unlock()
	}()
	if _, isExist := TextStyle.Entries[textStyleAlias]; isExist {
		return true
	}
	return false
}
