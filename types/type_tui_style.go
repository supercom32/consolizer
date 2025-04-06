package types

import (
	"encoding/json"
	"supercom32.net/consolizer/constants"
)

// TuiStyleEntryType
type TuiStyleEntryType struct {
	// Character definitions.
	UpperLeftCorner     rune
	UpperRightCorner    rune
	HorizontalLine      rune
	LeftSideTConnector  rune
	RightSideTConnector rune
	UpSideTConnector    rune
	DownSideTConnector  rune
	VerticalLine        rune
	LowerRightCorner    rune
	LowerLeftCorner     rune
	CrossConnector      rune
	DesktopPattern      rune
	// Label styles
	LabelForegroundColor constants.ColorType
	LabelBackgroundColor constants.ColorType
	// Checkbox styles.
	CheckboxForegroundColor     constants.ColorType
	CheckboxBackgroundColor     constants.ColorType
	CheckboxSelectedCharacter   rune
	CheckboxUnselectedCharacter rune
	// Radio Button styles.
	RadioButtonForegroundColor     constants.ColorType
	RadioButtonBackgroundColor     constants.ColorType
	RadioButtonSelectedCharacter   rune
	RadioButtonUnselectedCharacter rune
	// Scrollbar styles.
	ScrollbarTrackPattern    rune
	ScrollbarHandle          rune
	ScrollbarUpArrow         rune
	ScrollbarDownArrow       rune
	ScrollbarLeftArrow       rune
	ScrollbarRightArrow      rune
	ScrollbarForegroundColor constants.ColorType
	ScrollbarBackgroundColor constants.ColorType
	ScrollbarHandleColor     constants.ColorType
	// Progress Bar style.
	ProgressBarUnfilledPattern         rune
	ProgressBarFilledPattern           rune
	ProgressBarUnfilledForegroundColor constants.ColorType
	ProgressBarUnfilledBackgroundColor constants.ColorType
	ProgressBarFilledForegroundColor   constants.ColorType
	ProgressBarFilledBackgroundColor   constants.ColorType
	ProgressBarTextForegroundColor     constants.ColorType
	ProgressBarTextBackgroundColor     constants.ColorType
	// Random styles.
	IsSquareFont        bool
	IsWindowHeaderDrawn bool
	IsWindowFooterDrawn bool
	// For removal?
	LineDrawingTextForegroundColor constants.ColorType
	LineDrawingTextBackgroundColor constants.ColorType
	LineDrawingTextLabelColor      constants.ColorType
	//
	TextFieldForegroundColor constants.ColorType
	TextFieldBackgroundColor constants.ColorType
	// CURRENTLY UNUSED.
	TextFieldCursorForegroundColor constants.ColorType
	TextFieldCursorBackgroundColor constants.ColorType
	// Textbox styles.
	TextboxForegroundColor       constants.ColorType
	TextboxBackgroundColor       constants.ColorType
	TextboxCursorForegroundColor constants.ColorType
	TextboxCursorBackgroundColor constants.ColorType
	// Selector styles.
	SelectorForegroundColor  constants.ColorType
	SelectorBackgroundColor  constants.ColorType
	HighlightForegroundColor constants.ColorType
	HighlightBackgroundColor constants.ColorType
	// Button Styles
	ButtonRaisedColor        constants.ColorType
	ButtonForegroundColor    constants.ColorType
	ButtonBackgroundColor    constants.ColorType
	ButtonLabelDisabledColor constants.ColorType
	SelectorTextAlignment    int
	// Tooltip styles.
	TooltipForegroundColor     constants.ColorType
	TooltipBackgroundColor     constants.ColorType
	TooltipTextForegroundColor constants.ColorType
	TooltipTextBackgroundColor constants.ColorType
	TooltipDrawWindow          bool
}

func (shared TuiStyleEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

// NewTuiStyleEntry (existingCharacterObject ...*CharacterEntryType) CharacterEntryType
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
		styleEntry.LabelForegroundColor = existingStyleEntry[0].LabelForegroundColor
		styleEntry.LabelBackgroundColor = existingStyleEntry[0].LabelBackgroundColor
		styleEntry.CheckboxForegroundColor = existingStyleEntry[0].CheckboxForegroundColor
		styleEntry.CheckboxBackgroundColor = existingStyleEntry[0].CheckboxBackgroundColor
		styleEntry.CheckboxSelectedCharacter = existingStyleEntry[0].CheckboxSelectedCharacter
		styleEntry.CheckboxUnselectedCharacter = existingStyleEntry[0].CheckboxUnselectedCharacter
		styleEntry.RadioButtonForegroundColor = existingStyleEntry[0].RadioButtonForegroundColor
		styleEntry.RadioButtonBackgroundColor = existingStyleEntry[0].RadioButtonBackgroundColor
		styleEntry.RadioButtonSelectedCharacter = existingStyleEntry[0].RadioButtonSelectedCharacter
		styleEntry.RadioButtonUnselectedCharacter = existingStyleEntry[0].RadioButtonUnselectedCharacter
		styleEntry.ScrollbarTrackPattern = existingStyleEntry[0].ScrollbarTrackPattern
		styleEntry.ScrollbarHandle = existingStyleEntry[0].ScrollbarHandle
		styleEntry.ScrollbarUpArrow = existingStyleEntry[0].ScrollbarUpArrow
		styleEntry.ScrollbarDownArrow = existingStyleEntry[0].ScrollbarDownArrow
		styleEntry.ScrollbarLeftArrow = existingStyleEntry[0].ScrollbarLeftArrow
		styleEntry.ScrollbarRightArrow = existingStyleEntry[0].ScrollbarRightArrow
		styleEntry.ScrollbarForegroundColor = existingStyleEntry[0].ScrollbarForegroundColor
		styleEntry.ScrollbarBackgroundColor = existingStyleEntry[0].ScrollbarBackgroundColor
		styleEntry.ScrollbarHandleColor = existingStyleEntry[0].ScrollbarHandleColor
		styleEntry.ProgressBarUnfilledPattern = existingStyleEntry[0].ProgressBarUnfilledPattern
		styleEntry.ProgressBarFilledPattern = existingStyleEntry[0].ProgressBarFilledPattern
		styleEntry.ProgressBarUnfilledForegroundColor = existingStyleEntry[0].ProgressBarUnfilledForegroundColor
		styleEntry.ProgressBarUnfilledBackgroundColor = existingStyleEntry[0].ProgressBarUnfilledBackgroundColor
		styleEntry.ProgressBarFilledForegroundColor = existingStyleEntry[0].ProgressBarFilledForegroundColor
		styleEntry.ProgressBarFilledBackgroundColor = existingStyleEntry[0].ProgressBarFilledBackgroundColor
		styleEntry.ProgressBarTextForegroundColor = existingStyleEntry[0].ProgressBarTextForegroundColor
		styleEntry.ProgressBarTextBackgroundColor = existingStyleEntry[0].ProgressBarTextBackgroundColor
		styleEntry.LineDrawingTextForegroundColor = existingStyleEntry[0].LineDrawingTextForegroundColor
		styleEntry.LineDrawingTextBackgroundColor = existingStyleEntry[0].LineDrawingTextBackgroundColor
		styleEntry.LineDrawingTextLabelColor = existingStyleEntry[0].LineDrawingTextLabelColor
		styleEntry.TextFieldForegroundColor = existingStyleEntry[0].TextFieldForegroundColor
		styleEntry.TextFieldBackgroundColor = existingStyleEntry[0].TextFieldBackgroundColor
		styleEntry.TextFieldCursorForegroundColor = existingStyleEntry[0].TextFieldCursorForegroundColor
		styleEntry.TextFieldCursorBackgroundColor = existingStyleEntry[0].TextFieldCursorBackgroundColor
		styleEntry.TextboxForegroundColor = existingStyleEntry[0].TextboxForegroundColor
		styleEntry.TextboxBackgroundColor = existingStyleEntry[0].TextboxBackgroundColor
		styleEntry.TextboxCursorForegroundColor = existingStyleEntry[0].TextboxCursorForegroundColor
		styleEntry.TextboxCursorBackgroundColor = existingStyleEntry[0].TextboxCursorBackgroundColor
		styleEntry.SelectorForegroundColor = existingStyleEntry[0].SelectorForegroundColor
		styleEntry.SelectorBackgroundColor = existingStyleEntry[0].SelectorBackgroundColor
		styleEntry.HighlightForegroundColor = existingStyleEntry[0].HighlightForegroundColor
		styleEntry.HighlightBackgroundColor = existingStyleEntry[0].HighlightBackgroundColor
		styleEntry.ButtonRaisedColor = existingStyleEntry[0].ButtonRaisedColor
		styleEntry.ButtonForegroundColor = existingStyleEntry[0].ButtonForegroundColor
		styleEntry.ButtonBackgroundColor = existingStyleEntry[0].ButtonBackgroundColor
		styleEntry.ButtonLabelDisabledColor = existingStyleEntry[0].ButtonLabelDisabledColor
		styleEntry.IsSquareFont = existingStyleEntry[0].IsSquareFont
		styleEntry.IsWindowFooterDrawn = existingStyleEntry[0].IsWindowFooterDrawn
		styleEntry.IsWindowHeaderDrawn = existingStyleEntry[0].IsWindowHeaderDrawn
		styleEntry.SelectorTextAlignment = existingStyleEntry[0].SelectorTextAlignment
		styleEntry.TooltipForegroundColor = existingStyleEntry[0].TooltipForegroundColor
		styleEntry.TooltipBackgroundColor = existingStyleEntry[0].TooltipBackgroundColor
		styleEntry.TooltipTextForegroundColor = existingStyleEntry[0].TooltipTextForegroundColor
		styleEntry.TooltipTextBackgroundColor = existingStyleEntry[0].TooltipTextBackgroundColor
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
		styleEntry.ScrollbarTrackPattern = constants.CharBlockSparce
		styleEntry.ScrollbarHandle = constants.CharBlockSolid
		styleEntry.ScrollbarUpArrow = constants.CharTriangleUp
		styleEntry.ScrollbarDownArrow = constants.CharTriangleDown
		styleEntry.ScrollbarLeftArrow = constants.CharTriangleLeft
		styleEntry.ScrollbarRightArrow = constants.CharTriangleRight
		styleEntry.LabelForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.LabelBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.ProgressBarUnfilledPattern = constants.CharBlockSparce
		styleEntry.ProgressBarFilledPattern = constants.CharBlockSolid
		styleEntry.ProgressBarUnfilledForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.ProgressBarUnfilledBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.ProgressBarFilledForegroundColor = constants.AnsiColorByIndex[3]
		styleEntry.ProgressBarFilledBackgroundColor = constants.AnsiColorByIndex[3]
		styleEntry.LineDrawingTextForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.LineDrawingTextBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.LineDrawingTextLabelColor = constants.AnsiColorByIndex[15]
		styleEntry.TextFieldForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.TextFieldBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.TextFieldCursorForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.TextFieldCursorBackgroundColor = constants.AnsiColorByIndex[15]
		styleEntry.TextboxForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.TextboxBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.CheckboxForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.CheckboxBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.CheckboxSelectedCharacter = constants.CharUncheckedBox
		styleEntry.CheckboxUnselectedCharacter = constants.CharCheckedBox
		styleEntry.RadioButtonForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.RadioButtonBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.RadioButtonSelectedCharacter = constants.CharUncheckedRadioButton
		styleEntry.RadioButtonUnselectedCharacter = constants.CharCheckedRadioButton
		styleEntry.ScrollbarForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.ScrollbarBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.ScrollbarHandleColor = constants.AnsiColorByIndex[15]
		styleEntry.TextboxCursorForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.TextboxCursorBackgroundColor = constants.AnsiColorByIndex[15]
		styleEntry.SelectorForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.SelectorBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.HighlightForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.HighlightBackgroundColor = constants.AnsiColorByIndex[15]
		styleEntry.ButtonRaisedColor = constants.AnsiColorByIndex[15]
		styleEntry.ButtonForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.ButtonBackgroundColor = constants.AnsiColorByIndex[7]
		styleEntry.ButtonLabelDisabledColor = constants.AnsiColorByIndex[15]
		styleEntry.IsSquareFont = false
		styleEntry.IsWindowFooterDrawn = false
		styleEntry.IsWindowHeaderDrawn = false
		styleEntry.SelectorTextAlignment = constants.AlignmentLeft
		styleEntry.TooltipForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.TooltipBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.TooltipTextForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.TooltipTextBackgroundColor = constants.AnsiColorByIndex[0]
	}
	return styleEntry
}
