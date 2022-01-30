package memory

import (
	"encoding/json"
	"github.com/supercom32/consolizer/constants"
)

type AttributeEntryType struct {
	ForegroundColor          int32
	BackgroundColor          int32
	IsBold                   bool
	IsUnderlined             bool
	IsReversed               bool
	IsBlinking               bool
	IsItalic                 bool
	ForegroundTransformValue float32
	BackgroundTransformValue float32
	CellUserId               int
	CellTypeId				 int // The unique ID of a type instance
	CellType                 int // The type of control a cell belongs to
	CellAlias                string // The alias of the cell control.
}

func (shared AttributeEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		ForegroundColor int32
		BackgroundColor int32
		IsBold bool
		IsUnderlined bool
		IsReversed bool
		IsBlinking bool
		IsItalic bool
		ForegroundTransformValue float32
		BackgroundTransformValue float32
		CellUserId int
		CellTypeId int
		CellType int
		CellAlias string
	}{
		ForegroundColor: shared.ForegroundColor,
		BackgroundColor: shared.BackgroundColor,
		IsBold: shared.IsBold,
		IsUnderlined: shared.IsUnderlined,
		IsReversed: shared.IsReversed,
		IsBlinking: shared.IsBlinking,
		IsItalic: shared.IsItalic,
		ForegroundTransformValue: shared.ForegroundTransformValue,
		BackgroundTransformValue: shared.ForegroundTransformValue,
		CellUserId: shared.CellUserId, // An identifier for the instance of a cell type (Ie. button instance, etc).
		CellTypeId: shared.CellTypeId,
		// Need an attribute for sub cell type.
		CellType: shared.CellType, // Type of control the cell belongs to
		CellAlias: shared.CellAlias, // A string that represents some kind of string id.
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared AttributeEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewAttributeEntry(existingAttributeEntry ...*AttributeEntryType) AttributeEntryType {
	var attributeEntry AttributeEntryType
	if existingAttributeEntry != nil {
		attributeEntry.ForegroundColor = existingAttributeEntry[0].ForegroundColor
		attributeEntry.BackgroundColor = existingAttributeEntry[0].BackgroundColor
		attributeEntry.IsBold = existingAttributeEntry[0].IsBold
		attributeEntry.IsUnderlined = existingAttributeEntry[0].IsUnderlined
		attributeEntry.IsReversed = existingAttributeEntry[0].IsReversed
		attributeEntry.IsBlinking = existingAttributeEntry[0].IsBlinking
		attributeEntry.IsItalic = existingAttributeEntry[0].IsItalic
		attributeEntry.ForegroundTransformValue = existingAttributeEntry[0].ForegroundTransformValue
		attributeEntry.BackgroundTransformValue = existingAttributeEntry[0].BackgroundTransformValue
		attributeEntry.CellUserId = existingAttributeEntry[0].CellUserId
		attributeEntry.CellType = existingAttributeEntry[0].CellType
		attributeEntry.CellTypeId = existingAttributeEntry[0].CellTypeId
		attributeEntry.CellAlias = existingAttributeEntry[0].CellAlias
	} else {
		attributeEntry.ForegroundTransformValue = 1
		attributeEntry.BackgroundTransformValue = 1
		attributeEntry.ForegroundColor = constants.AnsiColorByIndex[15]
		attributeEntry.BackgroundColor = constants.AnsiColorByIndex[0]
		attributeEntry.CellUserId = constants.NullCellId
		attributeEntry.CellType = constants.NullCellType
		attributeEntry.CellTypeId = constants.NullCellId
	}
	return attributeEntry
}
