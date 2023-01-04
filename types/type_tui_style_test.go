package types

import (
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/recast"
	"testing"
)

func TestStyleTypeCreation(test *testing.T) {
	firstStyleEntry := NewTuiStyleEntry()
	firstStyleEntry.UpperLeftCorner = 'a'
	firstStyleEntry.UpperRightCorner = 'b'
	firstStyleEntry.HorizontalLine = 'c'
	firstStyleEntry.LeftSideTConnector = 'd'
	firstStyleEntry.RightSideTConnector = 'e'
	firstStyleEntry.UpSideTConnector = 'f'
	firstStyleEntry.DownSideTConnector = 'g'
	firstStyleEntry.VerticalLine = 'g'
	firstStyleEntry.LowerRightCorner = 'h'
	firstStyleEntry.LowerLeftCorner = 'i'
	firstStyleEntry.CrossConnector = 'j'
	firstStyleEntry.DesktopPattern = 'k'
	firstStyleEntry.ProgressBarUnfilledPattern = 'l'
	firstStyleEntry.ProgressBarFilledPattern = 'm'
	firstStyleEntry.LineDrawingTextForegroundColor = constants.AnsiColorByIndex[1]
	firstStyleEntry.LineDrawingTextBackgroundColor = constants.AnsiColorByIndex[2]
	firstStyleEntry.TextFieldForegroundColor = constants.AnsiColorByIndex[3]
	firstStyleEntry.TextFieldBackgroundColor = constants.AnsiColorByIndex[4]
	firstStyleEntry.TextboxCursorForegroundColor = constants.AnsiColorByIndex[5]
	firstStyleEntry.TextboxCursorBackgroundColor = constants.AnsiColorByIndex[6]
	firstStyleEntry.SelectorForegroundColor = constants.AnsiColorByIndex[7]
	firstStyleEntry.SelectorBackgroundColor = constants.AnsiColorByIndex[8]
	firstStyleEntry.HighlightForegroundColor = constants.AnsiColorByIndex[9]
	firstStyleEntry.HighlightBackgroundColor = constants.AnsiColorByIndex[10]
	firstStyleEntry.ButtonRaisedColor = constants.AnsiColorByIndex[11]
	firstStyleEntry.ButtonForegroundColor = constants.AnsiColorByIndex[12]
	firstStyleEntry.ButtonBackgroundColor = constants.AnsiColorByIndex[13]
	firstStyleEntry.IsSquareFont = true
	firstStyleEntry.IsWindowFooterDrawn = true
	firstStyleEntry.IsWindowHeaderDrawn = true
	firstStyleEntry.SelectorTextAlignment = constants.AlignmentLeft
	secondStyleEntry := NewTuiStyleEntry()
	firstResult := recast.GetArrayOfInterfaces(firstStyleEntry)
	secondResult := recast.GetArrayOfInterfaces(secondStyleEntry)
	assert.NotEqualf(test, secondResult, firstResult, "The first style entry is the same as the second, even though it should be different.")

	secondStyleEntry = NewTuiStyleEntry(&firstStyleEntry)
	secondResult = recast.GetArrayOfInterfaces(secondStyleEntry)
	assert.Equalf(test, secondResult, firstResult, "The first style entry is not the same as the second, even though it should be an identical clone.")
}
