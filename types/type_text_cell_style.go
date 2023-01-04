package types

import (
	"encoding/json"
	"github.com/supercom32/consolizer/constants"
)

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

func (shared TextCellStyleEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

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
