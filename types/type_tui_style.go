package types

import (
	"encoding/json"
	"github.com/supercom32/consolizer/constants"
)

/*
FrameStyle is a structure which contains character definitions for UI elements.

Example:
    var frameStyle FrameStyle
*/
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
	ForegroundColor     constants.ColorType
	BackgroundColor     constants.ColorType
}

/*
LabelStyle is a structure which contains styles for labels.

Example:
    var labelStyle LabelStyle
*/
type LabelStyle struct {
	ForegroundColor constants.ColorType
	BackgroundColor constants.ColorType
}

/*
CheckboxStyle is a structure which contains styles for checkboxes.

Example:
    var checkboxStyle CheckboxStyle
*/
type CheckboxStyle struct {
	ForegroundColor     constants.ColorType
	BackgroundColor     constants.ColorType
	SelectedCharacter   rune
	UnselectedCharacter rune
}

/*
RadioButtonStyle is a structure which contains styles for radio buttons.

Example:
    var radioButtonStyle RadioButtonStyle
*/
type RadioButtonStyle struct {
	ForegroundColor     constants.ColorType
	BackgroundColor     constants.ColorType
	SelectedCharacter   rune
	UnselectedCharacter rune
}

/*
ScrollbarStyle is a structure which contains styles for scrollbars.

Example:
    var scrollbarStyle ScrollbarStyle
*/
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

/*
ProgressBarStyle is a structure which contains styles for progress bars.

Example:
    var progressBarStyle ProgressBarStyle
*/
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

/*
TextFieldStyle is a structure which contains styles for text fields.

Example:
    var textFieldStyle TextFieldStyle
*/
type TextFieldStyle struct {
	ForegroundColor          constants.ColorType
	BackgroundColor          constants.ColorType
	HighlightForegroundColor constants.ColorType
	HighlightBackgroundColor constants.ColorType
	CursorForegroundColor    constants.ColorType
	CursorBackgroundColor    constants.ColorType
}

/*
TextboxStyle is a structure which contains styles for textboxes.

Example:
    var textboxStyle TextboxStyle
*/
type TextboxStyle struct {
	ForegroundColor          constants.ColorType
	BackgroundColor          constants.ColorType
	HighlightForegroundColor constants.ColorType
	HighlightBackgroundColor constants.ColorType
	CursorForegroundColor    constants.ColorType
	CursorBackgroundColor    constants.ColorType
}

/*
SelectorStyle is a structure which contains styles for selectors.

Example:
    var selectorStyle SelectorStyle
*/
type SelectorStyle struct {
	ForegroundColor          constants.ColorType
	BackgroundColor          constants.ColorType
	HighlightForegroundColor constants.ColorType
	HighlightBackgroundColor constants.ColorType
	TextAlignment            int
	IsSelectionCentered      bool
	IsShadowDrawn            bool
}

/*
ButtonStyle is a structure which contains styles for buttons.

Example:
    var buttonStyle ButtonStyle
*/
type ButtonStyle struct {
	RaisedColor        constants.ColorType
	ForegroundColor    constants.ColorType
	BackgroundColor    constants.ColorType
	LabelDisabledColor constants.ColorType
}

/*
TooltipStyle is a structure which contains styles for tooltips.

Example:
    var tooltipStyle TooltipStyle
*/
type TooltipStyle struct {
	ForegroundColor     constants.ColorType
	BackgroundColor     constants.ColorType
	TextForegroundColor constants.ColorType
	TextBackgroundColor constants.ColorType
	DrawWindow          bool
}

/*
WindowStyle is a structure which contains styles for windows.

Example:
    var windowStyle WindowStyle
*/
type WindowStyle struct {
	IsHeaderDrawn                       bool
	IsFooterDrawn                       bool
	LineDrawingTextForegroundColor      constants.ColorType
	LineDrawingTextBackgroundColor      constants.ColorType
	LineDrawingTextLabelColor           constants.ColorType
	LineDrawingSunkenColor              constants.ColorType
	LineDrawingRaisedColor              constants.ColorType
	LineDrawingTextLabelForegroundColor constants.ColorType
	LineDrawingTextLabelBackgroundColor constants.ColorType
}

/*
BarStyle is a structure which represents styles for bars.

Example:
    var barStyle BarStyle
*/
type BarStyle struct {
	ForegroundColor constants.ColorType
	BackgroundColor constants.ColorType
}

/*
FileMenuStyle is a structure which contains styles for file menus.

Example:
    var fileMenuStyle FileMenuStyle
*/
type FileMenuStyle struct {
	ForegroundColor          constants.ColorType
	BackgroundColor          constants.ColorType
	HighlightForegroundColor constants.ColorType
	HighlightBackgroundColor constants.ColorType
}

/*
DropdownStyle is a structure which contains styles for dropdowns.

Example:
    var dropdownStyle DropdownStyle
*/
type DropdownStyle struct {
	ForegroundColor constants.ColorType
	BackgroundColor constants.ColorType
	TextAlignment   int
}

/*
TextStyle is a structure which represents styles for text.

Example:
    var textStyle TextStyle
*/
type TextStyle struct {
	ForegroundColor constants.ColorType
	BackgroundColor constants.ColorType
}

/*
TuiStyleEntryType is a structure which contains all UI style definitions.

Example:
    var tuiStyle TuiStyleEntryType
*/
type TuiStyleEntryType struct {
	Text        TextStyle
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
	Bar         BarStyle
	FileMenu    FileMenuStyle
	Dropdown    DropdownStyle
}

/*
GetEntryAsJsonDump is a method which allows you to get a JSON string representation of the TUI style entry.

Example:
    instance.GetEntryAsJsonDump()
*/
func (shared TuiStyleEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewTuiStyleEntry is a constructor which allows you to create a new TuiStyleEntryType with default values or copy from an
existing one.

Example:
    NewTuiStyleEntry(existingStyleEntry)
*/
func NewTuiStyleEntry(existingStyleEntry ...*TuiStyleEntryType) TuiStyleEntryType {
	var styleEntry TuiStyleEntryType
	if existingStyleEntry != nil {
		styleEntry = *existingStyleEntry[0]
	} else {
		styleEntry.Text.BackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Text.ForegroundColor = constants.AnsiColorByIndex[15]
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
		styleEntry.Frame.ForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Frame.BackgroundColor = constants.AnsiColorByIndex[15]

		styleEntry.Label.ForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Label.BackgroundColor = constants.AnsiColorByIndex[0]

		styleEntry.Checkbox.ForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Checkbox.BackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Checkbox.SelectedCharacter = constants.CharCheckedBox
		styleEntry.Checkbox.UnselectedCharacter = constants.CharUncheckedBox

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
		styleEntry.Selector.IsShadowDrawn = false

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
		styleEntry.Window.LineDrawingTextLabelForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Window.LineDrawingTextLabelBackgroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Window.LineDrawingTextLabelColor = constants.AnsiColorByIndex[15]
		styleEntry.Window.LineDrawingRaisedColor = constants.AnsiColorByIndex[15]
		styleEntry.Window.LineDrawingSunkenColor = constants.AnsiColorByIndex[0]
		styleEntry.Frame.IsSquareFont = false

		styleEntry.FileMenu.ForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.FileMenu.BackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.FileMenu.HighlightForegroundColor = constants.AnsiColorByIndex[0]
		styleEntry.FileMenu.HighlightBackgroundColor = constants.AnsiColorByIndex[15]

		styleEntry.Dropdown.ForegroundColor = constants.AnsiColorByIndex[15]
		styleEntry.Dropdown.BackgroundColor = constants.AnsiColorByIndex[0]
		styleEntry.Dropdown.TextAlignment = constants.AlignmentLeft
	}

	return styleEntry
}
