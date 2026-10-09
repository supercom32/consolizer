package consolizer

import (
	"testing"
)

/*
TestViewportOutOfRangeSettings is a test which allows you to verify that a viewport cannot be put into a state that
makes it index outside its lines, where a negative scroll position used to panic on the next screen update and a
negative history limit used to panic while trimming history.

Example:

	Expected Inputs:
		A bordered viewport holding three lines, scrolled with SetViewport(-2, -1), then given SetMaxHistoryLines(-1)
		and one more line, with the screen updated after each step. Then a bordered viewport one row high holding
		two lines, trimmed to the lines it can show.

	Expected Outputs:
		No step panics. The scroll position is 0, 0, the history limit is unchanged, and all four lines are kept.
		Trimming the one row viewport, whose border leaves no room for a line, to its visible lines keeps none.
*/
func TestViewportOutOfRangeSettings(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	viewportInstance := viewport.Add(layerAlias, "view", styleEntry, 2, 2, 20, 6, false, true, 100)
	viewportInstance.SetContent("one\ntwo\nthree")
	viewportInstance.SetViewport(-2, -1)
	UpdateDisplay(false)
	viewportEntry := GetViewport(layerAlias, "view")
	if viewportEntry.ViewportXLocation != 0 || viewportEntry.ViewportYLocation != 0 {
		test.Fatalf("expected scroll position 0, 0, got %d, %d", viewportEntry.ViewportXLocation, viewportEntry.ViewportYLocation)
	}
	viewportInstance.SetMaxHistoryLines(-1)
	viewportInstance.Println("four")
	UpdateDisplay(false)
	if viewportEntry.MaxHistoryLines != 100 || len(viewportEntry.TextData) != 4 {
		test.Fatalf("expected limit 100 and 4 lines, got limit %d and %d lines", viewportEntry.MaxHistoryLines, len(viewportEntry.TextData))
	}

	viewport.Add(layerAlias, "tiny", styleEntry, 2, 12, 20, 1, false, true, 100)
	tinyViewportEntry := GetViewport(layerAlias, "tiny")
	tinyViewportEntry.TextData = [][]rune{[]rune("one"), []rune("two")}
	viewport.trimToVisible(tinyViewportEntry)
	if lineCount := len(tinyViewportEntry.TextData); lineCount != 0 {
		test.Fatalf("expected a viewport with no room for a line to keep none, got %d", lineCount)
	}
}
