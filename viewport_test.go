package consolizer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

/*
TestViewportOutOfRangeSettings is a test which allows you to verify that a viewport cannot be put into a state that
makes it index outside its lines, where a negative scroll position used to panic on the next screen update and a
negative history limit used to panic while trimming history.

Example:

	Expected Inputs:
		A bordered viewport given "one\ntwo\nthree\n", scrolled with SetViewport(-2, -1), then given
		SetMaxHistoryLines(-1) and Println("four"), with the screen updated after each step. Then a bordered viewport one row high holding
		two lines, trimmed to the lines it can show.

	Expected Outputs:
		No step panics. The scroll position is 0, 0, the history limit is unchanged, and all four lines are kept.
		Trimming the one row viewport, whose border leaves no room for a line, to its visible lines keeps none.
*/
func TestViewportOutOfRangeSettings(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	viewportInstance := viewport.Add(layerAlias, "view", styleEntry, 2, 2, 20, 6, false, true, 100)
	viewportInstance.SetContent("one\ntwo\nthree\n")
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
	viewport.trimToVisible(layerAlias, tinyViewportEntry)
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

/*
TestViewportAppliesMarkupStyles is a test which verifies that a viewport colours text inside a style tag while every
cell, styled or not, still belongs to the viewport. Before the fix the tag was looked up with its braces, as
"{{viewportWarn}}", so no style ever matched and the text stayed uncoloured.

Example:

	Expected Inputs:
	    A style "viewportWarn" (red) and an unbordered viewport "view" at (2, 2), 20 columns wide and 4 rows high,
	    holding the line "Go {{viewportWarn}}Red{{/}} end". The line is 31 runes long but prints in 10 columns, so it
	    only fits because tags take no columns.
	Expected Outputs:
	    Row 2 shows "Go Red end" from column 2. All ten cells have CellType CellTypeTextbox and CellControlAlias
	    "view", and only the three cells of "Red" are red.
*/
func TestViewportAppliesMarkupStyles(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	warnColor := GetRGBColor(255, 0, 0)
	addTestTextStyle(test, "viewportWarn", warnColor)
	viewportInstance := viewport.Add(layerAlias, "view", styleEntry, 2, 2, 20, 4, false, false, 100)
	viewportInstance.SetContent("Go {{viewportWarn}}Red{{/}} end")
	UpdateDisplay(false)

	getCell := func(column int) types.CharacterEntryType { return getCellInformationUnderMouseCursor(column, 2) }
	assertMarkupRowCells(test, "viewport", getCell, 2, constants.CellTypeTextbox, "view", warnColor)
}

/*
TestViewportWrapsTextByPrintedColumns is a test which verifies that text given to a viewport is split into lines by
printed columns: markup tags take no columns and are never split, wide characters take two, no line is wider than the
viewport, a wrapped line keeps its open style, and a "\n" does not carry a style onto the next line.

Example:

	Expected Inputs:
	    Each case gives SetContent a text for an unbordered viewport of the given width, where "{{vw}}" names a
	    style: "{{vw}}abcdefghij{{/}}" at 10, "abcdefghijk" at 10, "汉字汉字汉字" at 10, "a汉字汉字汉" at 10,
	    "{{vw}}one two three{{/}}" at 8, "abcdefghi{{vw}}jk" at 10, "abcdefghij klm" at 10, "{{vw}}a\nb" at 10,
	    "a\n" at 10, "ab{{cd" at 4, all wrapped, and "abcdefghijkl" at 5 unwrapped.
	Expected Outputs:
	    GetLines returns ["{{vw}}abcdefghij{{/}}"], ["abcdefghij", "k"], ["汉字汉字汉", "字"], ["a汉字汉字", "汉"],
	    ["{{vw}}one two ", "{{vw}}three{{/}}"], ["abcdefghi{{vw}}j", "{{vw}}k"], ["abcdefghij", "klm"],
	    ["{{vw}}a", "b"], ["a"], ["ab{{", "cd"] and ["abcdefghijkl"] respectively.
*/
func TestViewportWrapsTextByPrintedColumns(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	testCases := []struct {
		name           string
		text           string
		width          int
		isLinesWrapped bool
		expectedLines  []string
	}{
		{"tags take no columns", "{{vw}}abcdefghij{{/}}", 10, true, []string{"{{vw}}abcdefghij{{/}}"}},
		{"no line wider than the viewport", "abcdefghijk", 10, true, []string{"abcdefghij", "k"}},
		{"wide characters take two columns", "汉字汉字汉字", 10, true, []string{"汉字汉字汉", "字"}},
		{"wide character moved whole to the next line", "a汉字汉字汉", 10, true, []string{"a汉字汉字", "汉"}},
		{"wrapped word keeps its style", "{{vw}}one two three{{/}}", 8, true, []string{"{{vw}}one two ", "{{vw}}three{{/}}"}},
		{"tag never split", "abcdefghi{{vw}}jk", 10, true, []string{"abcdefghi{{vw}}j", "{{vw}}k"}},
		{"leading space dropped after a wrap", "abcdefghij klm", 10, true, []string{"abcdefghij", "klm"}},
		{"newline does not carry a style", "{{vw}}a\nb", 10, true, []string{"{{vw}}a", "b"}},
		{"trailing newline adds no line", "a\n", 10, true, []string{"a"}},
		{"unclosed tag is regular text", "ab{{cd", 4, true, []string{"ab{{", "cd"}},
		{"unwrapped lines are never split", "abcdefghijkl", 5, false, []string{"abcdefghijkl"}},
	}
	for caseIndex, testCase := range testCases {
		viewportInstance := viewport.Add(layerAlias, fmt.Sprintf("wrap%d", caseIndex), styleEntry, 0, 0, testCase.width, 4, testCase.isLinesWrapped, false, 100)
		viewportInstance.SetContent(testCase.text)
		assert.Equal(test, testCase.expectedLines, viewportInstance.GetLines(), testCase.name)
	}
}

/*
TestViewportPrintContinuesLastLine is a test which verifies that Print wraps text appended to the last line from the
column where that line already ends, keeps the style still open on it, and that text printed after a "\n" starts on
a new line even when the viewport was empty.

Example:

	Expected Inputs:
	    A style "vw" and an unbordered, wrapping viewport 10 columns wide and 4 rows high. Print("{{vw}}abcde") then
	    Print("fghijkl"). A second such viewport gets Print("a\n") then Print("b").
	Expected Outputs:
	    The first viewport holds ["{{vw}}abcdefghij", "{{vw}}kl"] and the second holds ["a", "b"].
*/
func TestViewportPrintContinuesLastLine(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	addTestTextStyle(test, "vw", GetRGBColor(255, 0, 0))

	styledViewport := viewport.Add(layerAlias, "styled", styleEntry, 2, 2, 10, 4, true, false, 100)
	styledViewport.Print("{{vw}}abcde")
	styledViewport.Print("fghijkl")
	assert.Equal(test, []string{"{{vw}}abcdefghij", "{{vw}}kl"}, styledViewport.GetLines(), "appended text should wrap by printed columns and keep its style")

	newlineViewport := viewport.Add(layerAlias, "newline", styleEntry, 2, 8, 10, 4, true, false, 100)
	newlineViewport.Print("a\n")
	newlineViewport.Print("b")
	assert.Equal(test, []string{"a", "b"}, newlineViewport.GetLines(), "text printed after a newline should start a new line")
}

/*
TestViewportHorizontalScrollByPrintedColumns is a test which verifies that horizontal scrolling and clipping are
measured in printed columns: a style opened before the first visible column still applies, a wide character that
would straddle either edge is left blank, and the widest line is measured without its tags.

Example:

	Expected Inputs:
	    A style "viewportScrollWarn" (red). An unbordered, unwrapped viewport "styled" at (2, 2), 10 columns wide,
	    holding "Go {{viewportScrollWarn}}Red{{/}} end" and scrolled to column 3. An unbordered, unwrapped viewport
	    "wide" at (2, 6), 4 columns wide, holding "a汉字", drawn at column 0 and then at column 2.
	Expected Outputs:
	    "styled" shows "Red end" from column 2 with "Red" red and " end" not, and reports a widest line of 10
	    columns. "wide" shows "a汉" then a blank column 5 at column 0, since "字" would cross the right edge, and a
	    blank column 2 then "字" at column 3 when scrolled to column 2, since "汉" would cross the left edge. It
	    reports a widest line of 5 columns.
*/
func TestViewportHorizontalScrollByPrintedColumns(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	warnColor := GetRGBColor(255, 0, 0)
	addTestTextStyle(test, "viewportScrollWarn", warnColor)
	getRowCell := func(row int) func(column int) types.CharacterEntryType {
		return func(column int) types.CharacterEntryType { return getCellInformationUnderMouseCursor(column, row) }
	}

	styledViewport := viewport.Add(layerAlias, "styled", styleEntry, 2, 2, 10, 3, false, false, 100)
	styledViewport.SetContent("Go {{viewportScrollWarn}}Red{{/}} end")
	styledViewport.SetViewport(3, 0)
	UpdateDisplay(false)
	getStyledCell := getRowCell(2)
	assertCellsShowText(test, "scrolled styled line", getStyledCell, 2, "Red end")
	for column := 2; column < 9; column++ {
		if column < 5 {
			assert.Equalf(test, warnColor, getStyledCell(column).AttributeEntry.ForegroundColor, "column %d should keep the style opened before the scroll position", column)
		} else {
			assert.NotEqualf(test, warnColor, getStyledCell(column).AttributeEntry.ForegroundColor, "column %d should not use the style", column)
		}
	}
	assert.Equal(test, 10, viewport.getMaxHorizontalTextValue(layerAlias, "styled"), "the widest line should be measured without its tags")

	wideViewport := viewport.Add(layerAlias, "wide", styleEntry, 2, 6, 4, 3, false, false, 100)
	wideViewport.SetContent("a汉字")
	UpdateDisplay(false)
	getWideCell := getRowCell(6)
	assertCellsShowText(test, "wide line at column 0", getWideCell, 2, "a汉 ")
	wideViewport.SetViewport(2, 0)
	UpdateDisplay(false)
	assertCellsShowText(test, "wide line at column 2", getWideCell, 2, " 字")
	assert.Equal(test, 5, viewport.getMaxHorizontalTextValue(layerAlias, "wide"), "wide characters should count as two columns")
}

/*
TestViewportScrollLimitsMatchVisibleArea is a test which verifies that auto-scrolling, scrolling to the bottom and
the mouse wheel all stop where the last line reaches the bottom row of the visible area, whether or not a border is
drawn, and that trimming history moves the scroll position up by the lines removed.

Example:

	Expected Inputs:
	    An unbordered viewport 10 by 4 given six lines through Println. A bordered viewport 10 by 6, which shows four
	    lines, given six lines through Print, then scrolled to line 0 and to the bottom. A bordered viewport 10 by 7,
	    which shows five lines, holding six lines, scrolled to line 0 and moved by three wheel down notches. An
	    unbordered viewport 10 by 4 holding ten lines, scrolled to line 6 and given a history limit of five lines.
	Expected Outputs:
	    The unbordered viewport is at line 2. The bordered viewport is at line 2 after printing and after
	    ScrollToBottom. The wheel stops at line 1. The trimmed viewport holds lines "5" to "9" and is at line 1.
*/
func TestViewportScrollLimitsMatchVisibleArea(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	assertLine := func(context string, viewportInstance ViewportInstanceType, expectedLine int) {
		test.Helper()
		if _, line := viewportInstance.GetViewport(); line != expectedLine {
			test.Fatalf("%s: expected line %d, got %d", context, expectedLine, line)
		}
	}

	unborderedViewport := viewport.Add(layerAlias, "unbordered", styleEntry, 0, 0, 10, 4, false, false, 100)
	for lineNumber := 0; lineNumber < 6; lineNumber++ {
		unborderedViewport.Println("line")
	}
	assertLine("Println into an unbordered viewport", unborderedViewport, 2)

	borderedViewport := viewport.Add(layerAlias, "bordered", styleEntry, 12, 0, 10, 6, false, true, 100)
	borderedViewport.Print("1\n2\n3\n4\n5\n6")
	assertLine("Print into a bordered viewport", borderedViewport, 2)
	borderedViewport.SetViewport(0, 0)
	borderedViewport.ScrollToBottom()
	assertLine("ScrollToBottom", borderedViewport, 2)

	wheelViewport := viewport.Add(layerAlias, "wheel", styleEntry, 24, 0, 10, 7, false, true, 100)
	wheelViewport.SetContent("1\n2\n3\n4\n5\n6")
	wheelViewport.SetViewport(0, 0)
	simScreen := startInputSimulation(test)
	for notch := 0; notch < 3; notch++ {
		injectWheel(simScreen, 26, 2, tcell.WheelDown)
	}
	assertLine("wheel down", wheelViewport, 1)

	trimmedViewport := viewport.Add(layerAlias, "trimmed", styleEntry, 0, 10, 10, 4, false, false, 100)
	trimmedViewport.SetContent("0\n1\n2\n3\n4\n5\n6\n7\n8\n9")
	trimmedViewport.SetViewport(0, 6)
	trimmedViewport.SetMaxHistoryLines(5)
	assert.Equal(test, []string{"5", "6", "7", "8", "9"}, trimmedViewport.GetLines(), "the oldest lines should be trimmed")
	assertLine("after trimming history", trimmedViewport, 1)
}

/*
TestViewportPrintlnEndsTheLine is a test which verifies that Println prints its text and then ends the line, as
Println does elsewhere in Go, that Print continues the current line, and that "\n" at the end of Print, SetContent
or Println("") behaves as a line break without leaving a blank row behind.

Example:

	Expected Inputs:
	    Each case runs a sequence of calls on a new, unbordered viewport 20 columns wide and 6 rows high:
	    Println("a") Print("b"); Print("a") Println("b"); Println("") Println(""); Println("a") Println("") Print("b");
	    Print("a\n") Print("b"); Print(""); Println("a") Print(""); SetContent("x\n") Print("y");
	    Println("a") Clear() Print("b"); Println("a") Println("b").
	Expected Outputs:
	    GetLines returns ["a", "b"], ["ab"], ["", ""], ["a", "", "b"], ["a", "b"], [], ["a"], ["x", "y"], ["b"] and
	    ["a", "b"] respectively.
*/
func TestViewportPrintlnEndsTheLine(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	testCases := []struct {
		name          string
		calls         func(viewportInstance ViewportInstanceType)
		expectedLines []string
	}{
		{"Println then Print", func(v ViewportInstanceType) { v.Println("a"); v.Print("b") }, []string{"a", "b"}},
		{"Print then Println", func(v ViewportInstanceType) { v.Print("a"); v.Println("b") }, []string{"ab"}},
		{"two empty Println calls", func(v ViewportInstanceType) { v.Println(""); v.Println("") }, []string{"", ""}},
		{"blank line between lines", func(v ViewportInstanceType) { v.Println("a"); v.Println(""); v.Print("b") }, []string{"a", "", "b"}},
		{"Print with a final newline", func(v ViewportInstanceType) { v.Print("a\n"); v.Print("b") }, []string{"a", "b"}},
		{"empty Print adds nothing", func(v ViewportInstanceType) { v.Print("") }, []string{}},
		{"empty Print does not start the pending line", func(v ViewportInstanceType) { v.Println("a"); v.Print("") }, []string{"a"}},
		{"SetContent with a final newline", func(v ViewportInstanceType) { v.SetContent("x\n"); v.Print("y") }, []string{"x", "y"}},
		{"Clear forgets the pending line break", func(v ViewportInstanceType) { v.Println("a"); v.Clear(); v.Print("b") }, []string{"b"}},
		{"Println lines", func(v ViewportInstanceType) { v.Println("a"); v.Println("b") }, []string{"a", "b"}},
	}
	for caseIndex, testCase := range testCases {
		viewportInstance := viewport.Add(layerAlias, fmt.Sprintf("println%d", caseIndex), styleEntry, 0, 0, 20, 6, true, false, 100)
		testCase.calls(viewportInstance)
		assert.Equal(test, testCase.expectedLines, viewportInstance.GetLines(), testCase.name)
	}
}

/*
TestViewportSetContentShowsTheTop is a test which verifies that SetContent shows new text from its first line and
column even when the viewport was scrolled, and that ScrollToBottom then shows its end.

Example:

	Expected Inputs:
	    An unbordered, unwrapped viewport 10 columns wide and 4 rows high holding twenty lines, scrolled to (3, 12).
	    SetContent is then given the ten lines "0" to "9", followed by a screen update and ScrollToBottom.
	Expected Outputs:
	    After SetContent the viewport is at (0, 0), its vertical scroll bar value is 0 and row 0 shows "0". After
	    ScrollToBottom it is at line 6.
*/
func TestViewportSetContentShowsTheTop(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	viewportInstance := viewport.Add(layerAlias, "document", styleEntry, 0, 0, 10, 4, false, false, 100)
	viewportInstance.SetContent(strings.Repeat("old line\n", 20))
	viewportInstance.SetViewport(3, 12)

	viewportInstance.SetContent("0\n1\n2\n3\n4\n5\n6\n7\n8\n9")
	UpdateDisplay(false)
	if column, line := viewportInstance.GetViewport(); column != 0 || line != 0 {
		test.Fatalf("expected SetContent to show (0, 0), got (%d, %d)", column, line)
	}
	viewportEntry := getViewport(layerAlias, "document")
	assert.Equal(test, 0, scrollbarRegistry.Get(layerAlias, viewportEntry.VerticalScrollbarAlias).ScrollValue, "the scroll bar should follow the viewport to the top")
	assert.Equal(test, '0', getCellInformationUnderMouseCursor(0, 0).Character, "the first line should be shown on the top row")

	viewportInstance.ScrollToBottom()
	if _, line := viewportInstance.GetViewport(); line != 6 {
		test.Fatalf("expected ScrollToBottom to show line 6, got %d", line)
	}
}

/*
TestViewportHistoryDisabledKeepsOnlyVisibleLines is a test which verifies that a viewport with history disabled keeps
only the newest lines it can show, whether text arrives through Println, Print or SetContent, and that a visible
horizontal scroll bar takes a row from that count.

Example:

	Expected Inputs:
	    Unbordered, unwrapped viewports 10 columns wide and 4 rows high with history disabled. One is given Println
	    for "0" to "5". One is given Print("0\n1\n2\n3\n4\n5"). One is given SetContent of "0" to "9". One is given
	    SetContent of "0" to "4" followed by a 16 column line, which needs the horizontal scroll bar.
	Expected Outputs:
	    ["2", "3", "4", "5"] for the first two, ["6", "7", "8", "9"] shown from line 0 for the third, and
	    ["3", "4", "abcdefghijklmnop"] for the fourth, as its scroll bar leaves three rows.
*/
func TestViewportHistoryDisabledKeepsOnlyVisibleLines(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	addViewport := func(viewportAlias string) ViewportInstanceType {
		viewportInstance := viewport.Add(layerAlias, viewportAlias, styleEntry, 0, 0, 10, 4, false, false, 100)
		viewportInstance.EnableHistory(false)
		return viewportInstance
	}

	printlnViewport := addViewport("println")
	for lineNumber := 0; lineNumber < 6; lineNumber++ {
		printlnViewport.Println(fmt.Sprint(lineNumber))
	}
	assert.Equal(test, []string{"2", "3", "4", "5"}, printlnViewport.GetLines(), "Println should keep only the visible lines")

	printViewport := addViewport("print")
	printViewport.Print("0\n1\n2\n3\n4\n5")
	assert.Equal(test, []string{"2", "3", "4", "5"}, printViewport.GetLines(), "Print should keep only the visible lines")

	contentViewport := addViewport("content")
	contentViewport.SetContent("0\n1\n2\n3\n4\n5\n6\n7\n8\n9")
	assert.Equal(test, []string{"6", "7", "8", "9"}, contentViewport.GetLines(), "SetContent should keep only the visible lines")
	if _, line := contentViewport.GetViewport(); line != 0 {
		test.Fatalf("expected SetContent to show the kept lines from line 0, got %d", line)
	}

	scrollbarViewport := addViewport("scrollbar")
	scrollbarViewport.SetContent("0\n1\n2\n3\n4\nabcdefghijklmnop")
	assert.Equal(test, []string{"3", "4", "abcdefghijklmnop"}, scrollbarViewport.GetLines(), "the horizontal scroll bar should take a row")
}
