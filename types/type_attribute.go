package types

import (
	"encoding/json"

	"github.com/supercom32/consolizer/constants"
)

/*
AttributeEntryType is a structure which contains the visual and metadata properties for a single terminal cell.

Example:

	attributeEntry := AttributeEntryType{
	    ForegroundColor: constants.ColorType{R: 255, G: 255, B: 255},
	    BackgroundColor: constants.ColorType{R: 0, G: 0, B: 0},
	}
*/
type AttributeEntryType struct {
	ForegroundColor         constants.ColorType
	BackgroundColor         constants.ColorType
	IsBold                  bool
	IsUnderlined            bool
	IsReversed              bool
	IsBlinking              bool
	IsItalic                bool
	IsBackgroundTransparent bool
	IsForegroundTransparent bool
	ForegroundAlphaValue    float32
	BackgroundAlphaValue    float32
	CellUserId              int
	CellControlId           int // The unique ID of a control type
	CellControlLocation     int // The relative location of a cell from within a control
	CellControlAlias        string
	CellType                int    // The type of control a cell belongs to
	CellUserAlias           string // The alias of the cell control.
}

/*
MarshalJSON is a method which serializes an attribute entry to JSON. In addition, the following should be noted:

- Converts the attribute entry's state to a JSON representation.

- Includes all visual and control-specific attributes.

- Used for saving and loading attribute configurations.

Example:

	jsonData, err := attributeEntry.MarshalJSON()
*/
func (shared AttributeEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		ForegroundColor          constants.ColorType
		BackgroundColor          constants.ColorType
		IsBold                   bool
		IsUnderlined             bool
		IsReversed               bool
		IsBlinking               bool
		IsItalic                 bool
		IsBackgroundTransparent  bool
		IsForegroundTransparent  bool
		ForegroundTransformValue float32
		BackgroundTransformValue float32
		CellUserAlias            string
		CellUserId               int
		CellControlAlias         string
		CellControlId            int
		CellControlLocation      int
		CellType                 int
	}{
		ForegroundColor:          shared.ForegroundColor,
		BackgroundColor:          shared.BackgroundColor,
		IsBold:                   shared.IsBold,
		IsUnderlined:             shared.IsUnderlined,
		IsReversed:               shared.IsReversed,
		IsBlinking:               shared.IsBlinking,
		IsItalic:                 shared.IsItalic,
		IsBackgroundTransparent:  shared.IsBackgroundTransparent,
		IsForegroundTransparent:  shared.IsForegroundTransparent,
		ForegroundTransformValue: shared.ForegroundAlphaValue,
		BackgroundTransformValue: shared.BackgroundAlphaValue,
		CellUserAlias:            shared.CellUserAlias, // A string that represents some kind of string id.
		CellUserId:               shared.CellUserId,    // An identifier for the instance of a cell type (Ie. button instance, etc).
		CellControlAlias:         shared.CellControlAlias,
		CellControlId:            shared.CellControlId,
		CellControlLocation:      shared.CellControlLocation,
		// Need an attribute for sub cell type.
		CellType: shared.CellType, // Type of control the cell belongs to

	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump is a method which returns a JSON string representation of an attribute entry. In addition, the following should be noted:

- Returns a formatted JSON string of the attribute entry's state.

- Useful for debugging and logging purposes.

- Panics if JSON marshaling fails.

Example:

	jsonString := attributeEntry.GetEntryAsJsonDump()
*/
func (shared AttributeEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewAttributeEntry is a constructor which creates a new attribute entry. In addition, the following should be noted:

- Initializes an attribute entry with default values.

- Can optionally copy properties from an existing attribute entry.

- Sets up all visual and control-specific attributes.

Example:

	attributeEntry := NewAttributeEntry(&existingAttributeEntry)
*/
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
		attributeEntry.IsBackgroundTransparent = existingAttributeEntry[0].IsBackgroundTransparent
		attributeEntry.IsForegroundTransparent = existingAttributeEntry[0].IsForegroundTransparent
		attributeEntry.ForegroundAlphaValue = existingAttributeEntry[0].ForegroundAlphaValue
		attributeEntry.BackgroundAlphaValue = existingAttributeEntry[0].BackgroundAlphaValue
		attributeEntry.CellUserAlias = existingAttributeEntry[0].CellUserAlias
		attributeEntry.CellUserId = existingAttributeEntry[0].CellUserId
		attributeEntry.CellType = existingAttributeEntry[0].CellType
		attributeEntry.CellControlId = existingAttributeEntry[0].CellControlId
		attributeEntry.CellControlAlias = existingAttributeEntry[0].CellControlAlias
		attributeEntry.CellControlLocation = existingAttributeEntry[0].CellControlLocation
	} else {
		attributeEntry.ForegroundAlphaValue = 1
		attributeEntry.BackgroundAlphaValue = 1
		attributeEntry.ForegroundColor = constants.AnsiColorByIndex[15]
		attributeEntry.BackgroundColor = constants.AnsiColorByIndex[0]
		attributeEntry.CellUserId = constants.NullCellId
		attributeEntry.CellType = constants.NullCellType
		attributeEntry.CellControlId = constants.NullCellControlId
		attributeEntry.CellControlLocation = constants.NullCellControlLocation
	}
	return attributeEntry
}
