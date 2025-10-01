package types

import (
	"encoding/json"
	"github.com/supercom32/consolizer/constants"
)

// FrameStyle contains character definitions for UI elements
type FrameStyle struct {
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
	IsSquareFont        bool
}

// LabelStyle contains styles for labels
type LabelStyle struct {
	ForegroundColor constants.ColorType
	BackgroundColor constants.ColorType
}

// CheckboxStyle contains styles for checkboxes
type CheckboxStyle struct {
	ForegroundColor     constants.ColorType
	BackgroundColor     constants.ColorType
	SelectedCharacter   rune
	UnselectedCharacter rune
}

// RadioButtonStyle contains styles for radio buttons
type RadioButtonStyle struct {
	ForegroundColor     constants.ColorType
	BackgroundColor     constants.ColorType
	SelectedCharacter   rune
	UnselectedCharacter rune
}

// ScrollbarStyle contains styles for scrollbars
type ScrollbarStyle struct {
	TrackPattern    rune
	Handle          rune
	UpArrow         rune
	DownArrow       rune
	LeftArrow       rune
	RightArrow      rune
	ForegroundColor constants.ColorType
	BackgroundColor constants.ColorType
	HandleColor     constants.ColorType
}

// ProgressBarStyle contains styles for progress bars
type ProgressBarStyle struct {
	UnfilledPattern         rune
	FilledPattern           rune
	UnfilledForegroundColor constants.ColorType
	UnfilledBackgroundColor constants.ColorType
	FilledForegroundColor   constants.ColorType
	FilledBackgroundColor   constants.ColorType
	TextForegroundColor     constants.ColorType
	TextBackgroundColor     constants.ColorType
	IsHighResolution        bool
}

// TextFieldStyle contains styles for text fields
type TextFieldStyle struct {
	ForegroundColor          constants.ColorType
	BackgroundColor          constants.ColorType
	HighlightForegroundColor constants.ColorType
	HighlightBackgroundColor constants.ColorType
	CursorForegroundColor    constants.ColorType
	CursorBackgroundColor    constants.ColorType
}

// TextboxStyle contains styles for textboxes
type TextboxStyle struct {
	ForegroundColor          constants.ColorType
	BackgroundColor          constants.ColorType
	HighlightForegroundColor constants.ColorType
	HighlightBackgroundColor constants.ColorType
	CursorForegroundColor    constants.ColorType
	CursorBackgroundColor    constants.ColorType
}

// SelectorStyle contains styles for selectors
type SelectorStyle struct {
	ForegroundColor          constants.ColorType
	BackgroundColor          constants.ColorType
	HighlightForegroundColor constants.ColorType
	HighlightBackgroundColor constants.ColorType
	TextAlignment            int
}

// ButtonStyle contains styles for buttons
type ButtonStyle struct {
	RaisedColor        constants.ColorType
	ForegroundColor    constants.ColorType
	BackgroundColor    constants.ColorType
	LabelDisabledColor constants.ColorType
}

// TooltipStyle contains styles for tooltips
type TooltipStyle struct {
	ForegroundColor     constants.ColorType
	BackgroundColor     constants.ColorType
	TextForegroundColor constants.ColorType
	TextBackgroundColor constants.ColorType
	DrawWindow          bool
}

// WindowStyle contains styles for windows
type WindowStyle struct {
	IsHeaderDrawn                  bool
	IsFooterDrawn                  bool
	LineDrawingTextForegroundColor constants.ColorType
	LineDrawingTextBackgroundColor constants.ColorType
	LineDrawingTextLabelColor      constants.ColorType
	LineDrawingSunkenColor         constants.ColorType
	LineDrawingRaisedColor         constants.ColorType
}

// TuiStyleEntryType contains all UI style definitions
type TuiStyleEntryType struct {
	Frame       FrameStyle
	Label       LabelStyle
	Checkbox    CheckboxStyle
	RadioButton RadioButtonStyle
	Scrollbar   ScrollbarStyle
	ProgressBar ProgressBarStyle
	TextField   TextFieldStyle
	Textbox     TextboxStyle
	Selector    SelectorStyle
	Button      ButtonStyle
	Tooltip     TooltipStyle
	Window      WindowStyle
}

func (shared TuiStyleEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

// NewTuiStyleEntry creates a new TuiStyleEntryType with default values or copies from an existing one
func NewTuiStyleEntry(existingStyleEntry ...*TuiStyleEntryType) TuiStyleEntryType {
	var styleEntry TuiStyleEntryType
	if existingStyleEntry != nil {
		styleEntry = *existingStyleEntry[0]
	} else {
		styleEntry.Frame.UpperLeftCorner = constants.CharULCorner
		styleEntry.Frame.UpperRightCorner = constants.CharURCorner
		styleEntry.Frame.HorizontalLine = constants.CharHLine
		styleEntry.Frame.LeftSideTConnector = constants.CharSingleLineTLeft
		styleEntry.Frame.RightSideTConnector = constants.CharSingleLineTRight
		styleEntry.Frame.UpSideTConnector = constants.CharSingleLineTUp
		styleEntry.Frame.DownSideTConnector = constants.CharSingleLineTDown
		styleEntry.Frame.VerticalLine = constants.CharSingleLineVertical
		styleEntry.Frame.LowerRightCorner = constants.CharSingleLineLowerRightCorner
		styleEntry.Frame.LowerLeftCorner = constants.CharSingleLineLowerLeftCorner
		styleEntry.Frame.CrossConnector = constants.CharSingleLineCross
		styleEntry.Frame.DesktopPattern = constants.CharBlockSparce

		styleEntry.Label.ForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Label.BackgroundColor = constants.AnsiColorByIndex[0]

		styleEntry.Checkbox.ForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Checkbox.BackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Checkbox.SelectedCharacter = constants.CharUncheckedBox
		styleEntry.Checkbox.UnselectedCharacter = constants.CharCheckedBox

		styleEntry.RadioButton.ForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.RadioButton.BackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.RadioButton.SelectedCharacter = constants.CharUncheckedRadioButton
		styleEntry.RadioButton.UnselectedCharacter = constants.CharCheckedRadioButton

		styleEntry.Scrollbar.TrackPattern = constants.CharBlockSparce
		styleEntry.Scrollbar.Handle = constants.CharBlockSolid
		styleEntry.Scrollbar.UpArrow = constants.CharTriangleUp
		styleEntry.Scrollbar.DownArrow = constants.CharTriangleDown
		styleEntry.Scrollbar.LeftArrow = constants.CharTriangleLeft
		styleEntry.Scrollbar.RightArrow = constants.CharTriangleRight
		styleEntry.Scrollbar.ForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Scrollbar.BackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Scrollbar.HandleColor = constants.AnsiColorByIndex[15]

		styleEntry.ProgressBar.UnfilledPattern = ' '
		styleEntry.ProgressBar.FilledPattern = constants.CharBlockSolid
		styleEntry.ProgressBar.UnfilledForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.ProgressBar.UnfilledBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.ProgressBar.FilledForegroundColor = constants.AnsiColorByIndex[3]
		styleEntry.ProgressBar.FilledBackgroundColor = constants.AnsiColorByIndex[3]
		styleEntry.ProgressBar.TextForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.ProgressBar.TextBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.ProgressBar.IsHighResolution = true

		styleEntry.TextField.ForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.TextField.BackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.TextField.CursorForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.TextField.CursorBackgroundColor = constants.AnsiColorByIndex[15]

		styleEntry.Textbox.ForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Textbox.BackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Textbox.CursorForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Textbox.CursorBackgroundColor = constants.AnsiColorByIndex[15]

		styleEntry.Selector.ForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Selector.BackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Selector.HighlightForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Selector.HighlightBackgroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Selector.TextAlignment = constants.AlignmentLeft

		styleEntry.Button.RaisedColor = constants.AnsiColorByIndex[15]
		styleEntry.Button.ForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Button.BackgroundColor = constants.AnsiColorByIndex[7]
		styleEntry.Button.LabelDisabledColor = constants.AnsiColorByIndex[15]

		styleEntry.Tooltip.ForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Tooltip.BackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Tooltip.TextForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Tooltip.TextBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Tooltip.DrawWindow = false

		styleEntry.Window.IsHeaderDrawn = false
		styleEntry.Window.IsFooterDrawn = false
		styleEntry.Window.LineDrawingTextForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Window.LineDrawingTextBackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Window.LineDrawingTextLabelColor = constants.AnsiColorByIndex[15]
		styleEntry.Window.LineDrawingRaisedColor = constants.AnsiColorByIndex[0]
		styleEntry.Window.LineDrawingSunkenColor = constants.AnsiColorByIndex[0]
		styleEntry.Frame.IsSquareFont = false
	}

	return styleEntry
}
