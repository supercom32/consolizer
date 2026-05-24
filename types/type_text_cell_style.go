package types

import (
	"encoding/json"
	"github.com/supercom32/consolizer/constants"
)

/*
TextCellStyleEntryType is a structure which represents the visual style of a single text cell. In addition, the following should be noted:

- Contains color and formatting information like bold, italic, and underline.

Example:

	var textCellStyle types.TextCellStyleEntryType
*/
type TextCellStyleEntryType struct {
	ForegroundColor          constants.ColorType
	BackgroundColor          constants.ColorType
	IsBold                   bool
	IsUnderlined             bool
	IsReversed               bool
	IsBlinking               bool
	IsItalic                 bool
	ForegroundTransformValue float32
	BackgroundTransformValue float32
}

/*
MarshalJSON is a method which allows you to serialize a text cell style to JSON. In addition, the following should be noted:

- Converts the style properties into a JSON format for persistence or transmission.

Example:

	instance.MarshalJSON()
*/
func (shared TextCellStyleEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		ForegroundColor          constants.ColorType
		BackgroundColor          constants.ColorType
		IsBold                   bool
		IsUnderlined             bool
		IsReversed               bool
		IsBlinking               bool
		IsItalic                 bool
		ForegroundTransformValue float32
		BackgroundTransformValue float32
	}{
		ForegroundColor:          shared.ForegroundColor,
		BackgroundColor:          shared.BackgroundColor,
		IsBold:                   shared.IsBold,
		IsUnderlined:             shared.IsUnderlined,
		IsReversed:               shared.IsReversed,
		IsBlinking:               shared.IsBlinking,
		IsItalic:                 shared.IsItalic,
		ForegroundTransformValue: shared.ForegroundTransformValue,
		BackgroundTransformValue: shared.ForegroundTransformValue,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump is a method which allows you to get a JSON string representation of a text cell style. In addition, the following should be noted:

- Panics if the marshaling process fails.

Example:

	instance.GetEntryAsJsonDump()
*/
func (shared TextCellStyleEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewTextCellStyleEntry is a constructor which allows you to create a new text cell style entry. In addition, the following should be noted:

- Can optionally initialize properties from an existing style entry.

Example:

	NewTextCellStyleEntry(existingAttributeEntry)
*/
func NewTextCellStyleEntry(existingAttributeEntry ...*TextCellStyleEntryType) TextCellStyleEntryType {
	var attributeEntry TextCellStyleEntryType
	if existingAttributeEntry != nil {
		attributeEntry.ForegroundColor = existingAttributeEntry[0].ForegroundColor
		attributeEntry.BackgroundColor = existingAttributeEntry[0].BackgroundColor
		attributeEntry.IsBold = existingAttributeEntry[0].IsBold
		attributeEntry.IsUnderlined = existingAttributeEntry[0].IsUnderlined
		attributeEntry.IsReversed = existingAttributeEntry[0].IsReversed
		attributeEntry.IsBlinking = existingAttributeEntry[0].IsBlinking
		attributeEntry.IsItalic = existingAttributeEntry[0].IsItalic
	} else {
		attributeEntry.ForegroundColor = constants.AnsiColorByIndex[15]
		attributeEntry.BackgroundColor = constants.AnsiColorByIndex[0]
	}
	return attributeEntry
}
