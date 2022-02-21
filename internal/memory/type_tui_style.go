package memory

import (
	"encoding/json"
	"github.com/supercom32/consolizer/constants"
)

// TuiStyleEntryType adsas
type TuiStyleEntryType struct {
	UpperLeftCorner              rune
	UpperRightCorner             rune
	HorizontalLine               rune
	LeftSideTConnector           rune
	RightSideTConnector          rune
	UpSideTConnector             rune
	DownSideTConnector           rune
	VerticalLine                 rune
	LowerRightCorner             rune
	LowerLeftCorner              rune
	CrossConnector               rune
	DesktopPattern          rune
	CheckboxForegroundColor   constants.ColorType
	CheckboxBackgroundColor   constants.ColorType
	CheckboxSelectedCharacter rune
	CheckboxUnselectedCharacter rune
	ScrollBarTrackPattern rune
	ScrollBarHandle rune
	ScrollBarUpArrow      rune
	ScrollBarDownArrow			 rune
	ScrollBarLeftArrow      rune
	ScrollBarRightArrow			 rune
	ScrollBarForegroundColor	constants.ColorType
	ScrollBarBackgroundColor	constants.ColorType
	ScrollBarHandleColor	 	constants.ColorType
	ProgressBarBackgroundPattern rune
	ProgressBarForegroundPattern rune
	IsSquareFont                 bool
	IsWindowHeaderDrawn          bool
	IsWindowFooterDrawn          bool
	TextForegroundColor          constants.ColorType
	TextBackgroundColor          constants.ColorType
	TextLabelColor				 constants.ColorType
	TextInputForegroundColor     constants.ColorType
	TextInputBackgroundColor constants.ColorType
	TextboxForegroundColor constants.ColorType
	TextboxBackgroundColor constants.ColorType
	CursorForegroundColor    constants.ColorType
	CursorBackgroundColor   constants.ColorType
	SelectorForegroundColor  constants.ColorType
	SelectorBackgroundColor  constants.ColorType
	HighlightForegroundColor constants.ColorType
	HighlightBackgroundColor constants.ColorType
	ButtonRaisedColor        constants.ColorType
	ButtonForegroundColor    constants.ColorType
	ButtonBackgroundColor constants.ColorType
	SelectorTextAlignment int
}

func (shared TuiStyleEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

// NewTuiStyleEntry asdasd existingCharacterObject ...*CharacterEntryType) CharacterEntryType
func NewTuiStyleEntry(existingStyleEntry ...*TuiStyleEntryType) TuiStyleEntryType {
	var styleEntry TuiStyleEntryType
	if existingStyleEntry != nil {
		styleEntry.UpperLeftCorner = existingStyleEntry[0].UpperLeftCorner
		styleEntry.UpperRightCorner = existingStyleEntry[0].UpperRightCorner
		styleEntry.HorizontalLine = existingStyleEntry[0].HorizontalLine
		styleEntry.LeftSideTConnector = existingStyleEntry[0].LeftSideTConnector
		styleEntry.RightSideTConnector = existingStyleEntry[0].RightSideTConnector
		styleEntry.UpSideTConnector = existingStyleEntry[0].UpSideTConnector
		styleEntry.DownSideTConnector = existingStyleEntry[0].DownSideTConnector
		styleEntry.VerticalLine = existingStyleEntry[0].VerticalLine
		styleEntry.LowerRightCorner = existingStyleEntry[0].LowerRightCorner
		styleEntry.LowerLeftCorner = existingStyleEntry[0].LowerLeftCorner
		styleEntry.CrossConnector = existingStyleEntry[0].CrossConnector
		styleEntry.DesktopPattern = existingStyleEntry[0].DesktopPattern
		styleEntry.CheckboxForegroundColor = existingStyleEntry[0].CheckboxForegroundColor
		styleEntry.CheckboxBackgroundColor = existingStyleEntry[0].CheckboxBackgroundColor
		styleEntry.CheckboxSelectedCharacter = existingStyleEntry[0].CheckboxSelectedCharacter
		styleEntry.CheckboxUnselectedCharacter = existingStyleEntry[0].CheckboxUnselectedCharacter
		styleEntry.ScrollBarTrackPattern = existingStyleEntry[0].ScrollBarTrackPattern
		styleEntry.ScrollBarHandle = existingStyleEntry[0].ScrollBarHandle
		styleEntry.ScrollBarUpArrow = existingStyleEntry[0].ScrollBarUpArrow
		styleEntry.ScrollBarDownArrow = existingStyleEntry[0].ScrollBarDownArrow
		styleEntry.ScrollBarLeftArrow = existingStyleEntry[0].ScrollBarLeftArrow
		styleEntry.ScrollBarRightArrow = existingStyleEntry[0].ScrollBarRightArrow
		styleEntry.ScrollBarForegroundColor = existingStyleEntry[0].ScrollBarForegroundColor
		styleEntry.ScrollBarBackgroundColor = existingStyleEntry[0].ScrollBarBackgroundColor
		styleEntry.ScrollBarHandleColor = existingStyleEntry[0].ScrollBarHandleColor
		styleEntry.ProgressBarBackgroundPattern = existingStyleEntry[0].ProgressBarBackgroundPattern
		styleEntry.ProgressBarForegroundPattern = existingStyleEntry[0].ProgressBarForegroundPattern
		styleEntry.TextForegroundColor = existingStyleEntry[0].TextForegroundColor
		styleEntry.TextBackgroundColor = existingStyleEntry[0].TextBackgroundColor
		styleEntry.TextLabelColor = existingStyleEntry[0].TextLabelColor
		styleEntry.TextInputForegroundColor = existingStyleEntry[0].TextInputForegroundColor
		styleEntry.TextInputBackgroundColor = existingStyleEntry[0].TextInputBackgroundColor
		styleEntry.TextboxForegroundColor = existingStyleEntry[0].TextboxForegroundColor
		styleEntry.TextboxBackgroundColor = existingStyleEntry[0].TextboxBackgroundColor
		styleEntry.CursorForegroundColor = existingStyleEntry[0].CursorForegroundColor
		styleEntry.CursorBackgroundColor = existingStyleEntry[0].CursorBackgroundColor
		styleEntry.SelectorForegroundColor = existingStyleEntry[0].SelectorForegroundColor
		styleEntry.SelectorBackgroundColor = existingStyleEntry[0].SelectorBackgroundColor
		styleEntry.HighlightForegroundColor = existingStyleEntry[0].HighlightForegroundColor
		styleEntry.HighlightBackgroundColor = existingStyleEntry[0].HighlightBackgroundColor
		styleEntry.ButtonRaisedColor = existingStyleEntry[0].ButtonRaisedColor
		styleEntry.ButtonForegroundColor = existingStyleEntry[0].ButtonForegroundColor
		styleEntry.ButtonBackgroundColor = existingStyleEntry[0].ButtonBackgroundColor
		styleEntry.IsSquareFont = existingStyleEntry[0].IsSquareFont
		styleEntry.IsWindowFooterDrawn = existingStyleEntry[0].IsWindowFooterDrawn
		styleEntry.IsWindowHeaderDrawn = existingStyleEntry[0].IsWindowHeaderDrawn
		styleEntry.SelectorTextAlignment = existingStyleEntry[0].SelectorTextAlignment
	} else {
		styleEntry.UpperLeftCorner = constants.CharULCorner
		styleEntry.UpperRightCorner = constants.CharURCorner
		styleEntry.HorizontalLine = constants.CharHLine
		styleEntry.LeftSideTConnector = constants.CharSingleLineTLeft
		styleEntry.RightSideTConnector = constants.CharSingleLineTRight
		styleEntry.UpSideTConnector = constants.CharSingleLineTUp
		styleEntry.DownSideTConnector = constants.CharSingleLineTDown
		styleEntry.VerticalLine = constants.CharSingleLineVertical
		styleEntry.LowerRightCorner = constants.CharSingleLineLowerRightCorner
		styleEntry.LowerLeftCorner = constants.CharSingleLineLowerLeftCorner
		styleEntry.CrossConnector = constants.CharSingleLineCross
		styleEntry.DesktopPattern = constants.CharBlockSparce
		styleEntry.ScrollBarTrackPattern = constants.CharBlockSparce
		styleEntry.ScrollBarHandle = constants.CharBlockSolid
		styleEntry.ScrollBarUpArrow = constants.CharTriangleUp
		styleEntry.ScrollBarDownArrow = constants.CharTriangleDown
		styleEntry.ScrollBarLeftArrow = constants.CharTriangleLeft
		styleEntry.ScrollBarRightArrow = constants.CharTriangleRight
		styleEntry.ProgressBarBackgroundPattern = constants.CharBlockSparce
		styleEntry.ProgressBarForegroundPattern = constants.CharBlockSolid
		styleEntry.TextForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.TextBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.TextLabelColor = constants.AnsiColorByIndex[15]
		styleEntry.TextInputForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.TextInputBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.TextboxForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.TextboxBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.CheckboxForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.CheckboxBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.CheckboxSelectedCharacter = constants.CharUncheckedBox
		styleEntry.CheckboxUnselectedCharacter = constants.CharCheckedBox
		styleEntry.ScrollBarForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.ScrollBarBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.ScrollBarHandleColor = constants.AnsiColorByIndex[15]
		styleEntry.CursorForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.CursorBackgroundColor = constants.AnsiColorByIndex[15]
		styleEntry.SelectorForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.SelectorBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.HighlightForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.HighlightBackgroundColor = constants.AnsiColorByIndex[15]
		styleEntry.ButtonRaisedColor = constants.AnsiColorByIndex[15]
		styleEntry.ButtonForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.ButtonBackgroundColor = constants.AnsiColorByIndex[7]
		styleEntry.IsSquareFont = false
		styleEntry.IsWindowFooterDrawn = false
		styleEntry.IsWindowHeaderDrawn = false
		styleEntry.SelectorTextAlignment = constants.AlignmentLeft
	}
	return styleEntry
}
