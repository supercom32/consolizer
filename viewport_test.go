package consolizer

import (
	"strings"
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
	viewportEntry := getViewport(layerAlias, "view")
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
	tinyViewportEntry := getViewport(layerAlias, "tiny")
	tinyViewportEntry.TextData = [][]rune{[]rune("one"), []rune("two")}
	viewport.trimToVisible(tinyViewportEntry)
	if lineCount := len(tinyViewportEntry.TextData); lineCount != 0 {
		test.Fatalf("expected a viewport with no room for a line to keep none, got %d", lineCount)
	}
}

/*
TestViewportGetViewportAndLines is a test which allows you to verify that a viewport's scroll position and lines can
be read through its methods, that the lines are a copy, and that a deleted viewport is ignored.

Example:

	Expected Inputs:
		A bordered viewport 20 wide and 4 tall given SetContent("one\ntwo\nthree\nfour"), then SetViewport(0, 1),
		then the returned line slice changed, then the viewport deleted.

	Expected Outputs:
		GetLines reads [one, two, three, four], GetViewport reads (0, 1), changing the returned slice does not change
		the viewport, and a deleted viewport reads (0, 0) and no lines.
*/
func TestViewportGetViewportAndLines(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	viewportInstance := viewport.Add(layerAlias, "view", styleEntry, 2, 2, 20, 4, false, true, 100)
	viewportInstance.SetContent("one\ntwo\nthree\nfour")
	lines := viewportInstance.GetLines()
	if strings.Join(lines, ",") != "one,two,three,four" {
		test.Fatalf("expected the lines [one two three four], got %q", lines)
	}
	viewportInstance.SetViewport(0, 1)
	if columnOffset, lineIndex := viewportInstance.GetViewport(); columnOffset != 0 || lineIndex != 1 {
		test.Fatalf("expected the viewport at (0, 1), got (%d, %d)", columnOffset, lineIndex)
	}
	lines[0] = "changed"
	if viewportInstance.GetLines()[0] != "one" {
		test.Fatalf("expected GetLines to return a copy")
	}
	viewportInstance.Delete()
	if columnOffset, lineIndex := viewportInstance.GetViewport(); columnOffset != 0 || lineIndex != 0 || len(viewportInstance.GetLines()) != 0 {
		test.Fatalf("expected a deleted viewport to read (0, 0) and no lines")
	}
}
