package types

import (
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/recast"
	"testing"
)

func TestStyleTypeCreation(test *testing.T) {
	firstStyleEntry := NewTuiStyleEntry()
	firstStyleEntry.Frame.UpperLeftCorner = 'a'
	firstStyleEntry.Frame.UpperRightCorner = 'b'
	firstStyleEntry.Frame.HorizontalLine = 'c'
	firstStyleEntry.Frame.LeftSideTConnector = 'd'
	firstStyleEntry.Frame.RightSideTConnector = 'e'
	firstStyleEntry.Frame.UpSideTConnector = 'f'
	firstStyleEntry.Frame.DownSideTConnector = 'g'
	firstStyleEntry.Frame.VerticalLine = 'g'
	firstStyleEntry.Frame.LowerRightCorner = 'h'
	firstStyleEntry.Frame.LowerLeftCorner = 'i'
	firstStyleEntry.Frame.CrossConnector = 'j'
	firstStyleEntry.Frame.DesktopPattern = 'k'
	firstStyleEntry.ProgressBar.UnfilledPattern = 'l'
	firstStyleEntry.ProgressBar.FilledPattern = 'm'
	firstStyleEntry.Window.LineDrawingTextForegroundColor = constants.AnsiColorByIndex[1]
	firstStyleEntry.Window.LineDrawingTextBackgroundColor = constants.AnsiColorByIndex[2]
	firstStyleEntry.TextField.ForegroundColor = constants.AnsiColorByIndex[3]
	firstStyleEntry.TextField.BackgroundColor = constants.AnsiColorByIndex[4]
	firstStyleEntry.Textbox.CursorForegroundColor = constants.AnsiColorByIndex[5]
	firstStyleEntry.Textbox.CursorBackgroundColor = constants.AnsiColorByIndex[6]
	firstStyleEntry.Selector.ForegroundColor = constants.AnsiColorByIndex[7]
	firstStyleEntry.Selector.BackgroundColor = constants.AnsiColorByIndex[8]
	firstStyleEntry.Textbox.HighlightForegroundColor = constants.AnsiColorByIndex[9]
	firstStyleEntry.Textbox.HighlightBackgroundColor = constants.AnsiColorByIndex[10]
	firstStyleEntry.Button.RaisedColor = constants.AnsiColorByIndex[11]
	firstStyleEntry.Button.ForegroundColor = constants.AnsiColorByIndex[12]
	firstStyleEntry.Button.BackgroundColor = constants.AnsiColorByIndex[13]
	firstStyleEntry.Misc.IsSquareFont = true
	firstStyleEntry.Window.IsFooterDrawn = true
	firstStyleEntry.Window.IsHeaderDrawn = true
	firstStyleEntry.Selector.TextAlignment = constants.AlignmentLeft
	secondStyleEntry := NewTuiStyleEntry()
	firstResult := recast.GetArrayOfInterfaces(firstStyleEntry)
	secondResult := recast.GetArrayOfInterfaces(secondStyleEntry)
	assert.NotEqualf(test, secondResult, firstResult, "The first style entry is the same as the second, even though it should be different.")

	secondStyleEntry = NewTuiStyleEntry(&firstStyleEntry)
	secondResult = recast.GetArrayOfInterfaces(secondStyleEntry)
	assert.Equalf(test, secondResult, firstResult, "The first style entry is not the same as the second, even though it should be an identical clone.")
}
