package consolizer

import (
	"strings"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"

	"github.com/supercom32/consolizer/types"
)

/*
ViewportInstanceType is a structure which represents a viewport instance that can be manipulated.
*/
type ViewportInstanceType struct {
	BaseControlInstanceType
}

/*
viewportType is a structure which is the internal type that manages all viewports.
*/
type viewportType struct{}

var viewport viewportType
var viewportRegistry = memory.NewControlMemoryManager[types.ViewportEntryType]()

/*
getViewport is a method which allows you to retrieve the entry holding a viewport's state by its layer and alias, or
nil if the viewport does not exist. The caller must hold uiStateMutex.

Example:

	viewportEntry := getViewport("main", "myVP")
*/
func getViewport(layerAlias string, viewportAlias string) *types.ViewportEntryType {
	viewportEntry, _ := viewportRegistry.Lookup(layerAlias, viewportAlias)
	return viewportEntry
}

/*
IsViewportExists is a method which allows you to check if a viewport exists in the specified layer.

Example:
    exists := IsViewportExists("main", "myVP")
*/
func IsViewportExists(layerAlias string, viewportAlias string) bool {
	return viewportRegistry.IsExists(layerAlias, viewportAlias)
}

/*
Println is a method which allows you to add text to the viewport and then end the line, as Println does elsewhere in
Go. The text continues the current line, so it follows anything Print left there, and whatever is printed next starts
on a new line. Println("") ends the current line, or adds a blank line when the current line is already ended. Text
is wrapped to the viewport's width when line wrapping is enabled, and the viewport is then scrolled so its last line
is shown. In addition, the following should be noted:

  - The line break is remembered rather than stored as an empty line, so no blank row is left at the bottom of the
    viewport after the last Println.

Example:

	vp.Println("Hello World!")
*/
func (shared *ViewportInstanceType) Println(text string) *ViewportInstanceType {
	defer lockUIState()()
	viewportEntry := getViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewport.appendText(shared.layerAlias, shared.controlAlias, viewportEntry, text+"\n")
	viewport.scrollToLastLine(shared.layerAlias, shared.controlAlias, viewportEntry)
	return shared
}

/*
Print is a method which allows you to add text to the viewport without ending the line, so it continues the current
line and the next text printed follows it. A "\n" in the text starts a new line, and Print("a\n") is the same as
Println("a"). Text is wrapped to the viewport's width when line wrapping is enabled, and the viewport is then scrolled
so its last line is shown.

Example:

	vp.Print("Appending some text...")
*/
func (shared *ViewportInstanceType) Print(text string) *ViewportInstanceType {
	defer lockUIState()()
	viewportEntry := getViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewport.appendText(shared.layerAlias, shared.controlAlias, viewportEntry, text)
	viewport.scrollToLastLine(shared.layerAlias, shared.controlAlias, viewportEntry)
	return shared
}

/*
SetContent is a method which allows you to replace everything the viewport holds with the given text and show it from
its first line and column, as a document would be opened. The text is split into lines exactly as Print would add it
to an empty viewport, so a final "\n" means the next text printed starts on a new line. Call ScrollToBottom
afterwards to show the end of the text instead. In addition, the following should be noted:

  - When history is disabled the viewport keeps only the lines it can show, so for text longer than that the kept
    lines are the last ones, and those are shown from the top.

Example:

	vp.SetContent("New content starts here.")
*/
func (shared *ViewportInstanceType) SetContent(text string) *ViewportInstanceType {
	defer lockUIState()()
	viewportEntry := getViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewportEntry.TextData = [][]rune{}
	viewportEntry.IsLineBreakPending = false
	viewport.appendText(shared.layerAlias, shared.controlAlias, viewportEntry, text)
	viewportEntry.ViewportXLocation = 0
	viewportEntry.ViewportYLocation = 0
	viewport.setViewportMaxScrollBarValues(shared.layerAlias, shared.controlAlias)
	viewport.updateScrollbarBasedOnViewportViewport(shared.layerAlias, shared.controlAlias)
	return shared
}

/*
SetViewport is a method which allows you to set the viewport scroll position, as the column and line of the content
shown in the viewport's top left corner. A negative column or line is treated as zero, since content cannot be
scrolled to before its start.

Example:

	vp.SetViewport(0, 10)
*/
func (shared *ViewportInstanceType) SetViewport(xLocation int, yLocation int) *ViewportInstanceType {
	defer lockUIState()()
	viewportEntry := getViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewportEntry.ViewportXLocation = getClampedIndex(xLocation, 0, xLocation)
	viewportEntry.ViewportYLocation = getClampedIndex(yLocation, 0, yLocation)

	// Update scrollbars
	viewport.setViewportMaxScrollBarValues(shared.layerAlias, shared.controlAlias)
	viewport.updateScrollbarBasedOnViewportViewport(shared.layerAlias, shared.controlAlias)

	return shared
}

/*
GetViewport is a method which allows you to obtain how far a viewport is scrolled, as the column offset of its first
visible column and the index of its first visible line. If the viewport does not exist, zero and zero are returned.

Example:

	columnOffset, lineIndex := viewportInstance.GetViewport()
*/
func (shared *ViewportInstanceType) GetViewport() (int, int) {
	defer lockUIState()()
	if viewportEntry, isFound := viewportRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return viewportEntry.ViewportXLocation, viewportEntry.ViewportYLocation
	}
	return 0, 0
}

/*
GetLines is a method which allows you to obtain every line a viewport currently holds, in order, as they are stored
for display. When the viewport wraps long lines, each wrapped row is a line of its own, and when history is limited
only the lines it still keeps are returned. If the viewport does not exist, an empty slice is returned. In addition,
the following should be noted:

  - The slice returned is a copy, so changing it does not change the viewport.

Example:

	lines := viewportInstance.GetLines()
*/
func (shared *ViewportInstanceType) GetLines() []string {
	defer lockUIState()()
	viewportEntry, isFound := viewportRegistry.Lookup(shared.layerAlias, shared.controlAlias)
	if !isFound {
		return []string{}
	}
	lines := make([]string, 0, len(viewportEntry.TextData))
	for _, line := range viewportEntry.TextData {
		lines = append(lines, string(line))
	}
	return lines
}

/*
ScrollToBottom is a method which allows you to scroll the viewport to the bottom of its content, so its last line is
shown on its bottom row, or to its first line when all of its lines fit.

Example:

	vp.ScrollToBottom()
*/
func (shared *ViewportInstanceType) ScrollToBottom() *ViewportInstanceType {
	defer lockUIState()()
	viewportEntry := getViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}
	viewport.scrollToLastLine(shared.layerAlias, shared.controlAlias, viewportEntry)
	return shared
}

/*
SetMaxHistoryLines is a method which allows you to set the maximum number of lines to keep in history. A negative
number of lines is ignored, leaving the current limit unchanged.

Example:

	vp.SetMaxHistoryLines(500)
*/
func (shared *ViewportInstanceType) SetMaxHistoryLines(maxLines int) *ViewportInstanceType {
	defer lockUIState()()
	viewportEntry := getViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil || maxLines < 0 {
		return shared
	}

	viewportEntry.MaxHistoryLines = maxLines
	viewport.trimHistory(viewportEntry)
	viewport.setViewportMaxScrollBarValues(shared.layerAlias, shared.controlAlias)
	viewport.updateScrollbarBasedOnViewportViewport(shared.layerAlias, shared.controlAlias)
	return shared
}

/*
EnableHistory is a method which allows you to enable or disable history tracking for the viewport.

Example:
    vp.EnableHistory(true)
*/
func (shared *ViewportInstanceType) EnableHistory(enabled bool) *ViewportInstanceType {
	defer lockUIState()()
	viewportEntry := getViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewportEntry.IsHistoryEnabled = enabled
	if !enabled {
		viewport.setViewportMaxScrollBarValues(shared.layerAlias, shared.controlAlias)
		viewport.trimToVisible(shared.layerAlias, viewportEntry)
	}
	viewport.setViewportMaxScrollBarValues(shared.layerAlias, shared.controlAlias)
	viewport.updateScrollbarBasedOnViewportViewport(shared.layerAlias, shared.controlAlias)
	return shared
}

/*
Clear is a method which allows you to clear all text from the viewport.

Example:
    vp.Clear()
*/
func (shared *ViewportInstanceType) Clear() *ViewportInstanceType {
	defer lockUIState()()
	viewportEntry := getViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewportEntry.TextData = [][]rune{}
	viewportEntry.IsLineBreakPending = false
	viewportEntry.ViewportXLocation = 0
	viewportEntry.ViewportYLocation = 0

	// Update scrollbars
	viewport.setViewportMaxScrollBarValues(shared.layerAlias, shared.controlAlias)
	viewport.updateScrollbarBasedOnViewportViewport(shared.layerAlias, shared.controlAlias)

	return shared
}

/*
SetTransparent is a method which allows you to set whether the viewport background is transparent.

Example:
    vp.SetTransparent(true)
*/
func (shared *ViewportInstanceType) SetTransparent(isTransparent bool) *ViewportInstanceType {
	defer lockUIState()()
	viewportEntry := getViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewportEntry.IsTransparent = isTransparent
	return shared
}

/*
viewportLineBuilderType is a structure which accumulates text into the lines a viewport stores, tracking the printed
column of the line being built and the style tag open on it so that wrapping is measured in visible columns and a
wrapped line keeps its style.
*/
type viewportLineBuilderType struct {
	widthInColumns    int
	isLinesWrapped    bool
	lines             [][]rune
	currentLine       []rune
	currentColumn     int
	activeTag         []rune
	isTagReopenNeeded bool
	isWrapped         bool
}

/*
newViewportLineBuilder is a constructor which allows you to create a line builder that continues from startingLine,
so text added to it is appended to that line. The starting line is copied, so the builder never writes into the slice
it was given, and its printed width and open style tag are measured so wrapping and style carrying continue from
where it ends. Pass nil to start from an empty line.

Example:

	lineBuilder := newViewportLineBuilder(40, true, viewportEntry.TextData[lastLineIndex])
*/
func newViewportLineBuilder(widthInColumns int, isLinesWrapped bool, startingLine []rune) *viewportLineBuilderType {
	lineBuilder := &viewportLineBuilderType{
		widthInColumns: widthInColumns,
		isLinesWrapped: isLinesWrapped,
		currentLine:    append([]rune{}, startingLine...),
	}
	for currentIndex := 0; currentIndex < len(startingLine); currentIndex++ {
		if tagEndIndex := findMarkupTagEnd(startingLine, currentIndex); tagEndIndex >= 0 {
			lineBuilder.setActiveTag(startingLine[currentIndex : tagEndIndex+1])
			currentIndex = tagEndIndex
			continue
		}
		lineBuilder.currentColumn += stringformat.GetWidthOfRuneWhenPrinted(startingLine[currentIndex])
	}
	return lineBuilder
}

/*
addText is a method which allows you to add text to the line builder, starting new lines at every "\n" and, when
wrapping is enabled, wherever the next word or character would not fit in the remaining columns. Markup tags are
copied whole onto the current line and take no columns. An unclosed "{{" is treated as regular text, as it is when
the line is drawn.

Example:

	lineBuilder.addText([]rune("Hello {{red}}World{{/}}"))
*/
func (shared *viewportLineBuilderType) addText(textToAdd []rune) {
	for currentIndex := 0; currentIndex < len(textToAdd); currentIndex++ {
		currentCharacter := textToAdd[currentIndex]
		if currentCharacter == '\n' {
			shared.startNewLine(false)
			continue
		}
		if tagEndIndex := findMarkupTagEnd(textToAdd, currentIndex); tagEndIndex >= 0 {
			markupTag := textToAdd[currentIndex : tagEndIndex+1]
			shared.currentLine = append(shared.currentLine, markupTag...)
			shared.setActiveTag(markupTag)
			shared.isTagReopenNeeded = false
			currentIndex = tagEndIndex
			continue
		}
		columnWidth := stringformat.GetWidthOfRuneWhenPrinted(currentCharacter)
		if shared.isLinesWrapped && shared.currentColumn > 0 && shared.currentColumn+columnWidth > shared.widthInColumns {
			shared.startNewLine(true)
		}
		if currentCharacter == ' ' && shared.isWrapped && shared.currentColumn == 0 {
			continue
		}
		if shared.isTagReopenNeeded {
			shared.currentLine = append(shared.currentLine, shared.activeTag...)
			shared.isTagReopenNeeded = false
		}
		shared.currentLine = append(shared.currentLine, currentCharacter)
		shared.currentColumn += columnWidth
		shared.isWrapped = false
		if shared.isLinesWrapped && currentCharacter == ' ' {
			nextWordWidth := getMarkupWordWidth(textToAdd, currentIndex+1)
			if nextWordWidth > 0 && shared.currentColumn+nextWordWidth > shared.widthInColumns {
				shared.startNewLine(true)
			}
		}
	}
}

/*
startNewLine is a method which allows you to finish the line being built and begin an empty one. When isWrap is set
the break was made by wrapping, so the open style tag, if any, is reopened before the first character of the new line
and leading spaces on it are dropped. A "\n" break closes the open style instead.

Example:

	lineBuilder.startNewLine(true)
*/
func (shared *viewportLineBuilderType) startNewLine(isWrap bool) {
	shared.lines = append(shared.lines, shared.currentLine)
	shared.currentLine = []rune{}
	shared.currentColumn = 0
	shared.isWrapped = isWrap
	if !isWrap {
		shared.activeTag = nil
	}
	shared.isTagReopenNeeded = isWrap && len(shared.activeTag) > 0
}

/*
setActiveTag is a method which allows you to record the style tag now open on the line being built. A closing tag
"{{/}}" or an empty tag "{{}}" leaves no style open, since both return the text to its default style.

Example:

	lineBuilder.setActiveTag([]rune("{{red}}"))
*/
func (shared *viewportLineBuilderType) setActiveTag(markupTag []rune) {
	tagContent := getAttributeTagContent(string(markupTag))
	if tagContent == "" || tagContent == "/" {
		shared.activeTag = nil
		return
	}
	shared.activeTag = markupTag
}

/*
getLines is a method which allows you to obtain every line the builder has produced, followed by the line still
being built, even when that line holds nothing.

Example:

	lines := lineBuilder.getLines()
*/
func (shared *viewportLineBuilderType) getLines() [][]rune {
	return append(shared.lines, shared.currentLine)
}

/*
findMarkupTagEnd is a method which allows you to find the markup tag that opens at a given rune index. If "{{" starts
at currentIndex and a closing "}}" follows it, the index of that final "}" is returned. Otherwise minus one is
returned, which means the text at currentIndex is regular text.

Example:

	tagEndIndex := findMarkupTagEnd([]rune("a{{red}}b"), 1)
*/
func findMarkupTagEnd(textToSearch []rune, currentIndex int) int {
	if currentIndex < 0 || currentIndex+1 >= len(textToSearch) || textToSearch[currentIndex] != '{' || textToSearch[currentIndex+1] != '{' {
		return -1
	}
	for tagEndIndex := currentIndex + 2; tagEndIndex+1 < len(textToSearch); tagEndIndex++ {
		if textToSearch[tagEndIndex] == '}' && textToSearch[tagEndIndex+1] == '}' {
			return tagEndIndex + 1
		}
	}
	return -1
}

/*
getMarkupWordWidth is a method which allows you to obtain the printed width in columns of the word that starts at a
given rune index, ending at the next space, newline or the end of the text. Markup tags inside the word take no
columns, and wide characters take two.

Example:

	wordWidth := getMarkupWordWidth([]rune("Go {{red}}Red{{/}} end"), 3)
*/
func getMarkupWordWidth(textToMeasure []rune, startingIndex int) int {
	wordWidth := 0
	for currentIndex := startingIndex; currentIndex < len(textToMeasure); currentIndex++ {
		currentCharacter := textToMeasure[currentIndex]
		if currentCharacter == ' ' || currentCharacter == '\n' {
			break
		}
		if tagEndIndex := findMarkupTagEnd(textToMeasure, currentIndex); tagEndIndex >= 0 {
			currentIndex = tagEndIndex
			continue
		}
		wordWidth += stringformat.GetWidthOfRuneWhenPrinted(currentCharacter)
	}
	return wordWidth
}

/*
getMarkupLineWidth is a method which allows you to obtain the printed width in columns of a stored viewport line.
Markup tags take no columns, and wide characters take two.

Example:

	lineWidth := getMarkupLineWidth([]rune("Go {{red}}Red{{/}} end"))
*/
func getMarkupLineWidth(line []rune) int {
	lineWidth := 0
	for currentIndex := 0; currentIndex < len(line); currentIndex++ {
		if tagEndIndex := findMarkupTagEnd(line, currentIndex); tagEndIndex >= 0 {
			currentIndex = tagEndIndex
			continue
		}
		lineWidth += stringformat.GetWidthOfRuneWhenPrinted(line[currentIndex])
	}
	return lineWidth
}

/*
getTextWidth is a method which allows you to obtain the number of columns a viewport wraps its text to, which is its
width less the border when one is drawn.

Example:

	textWidth := viewport.getTextWidth(entry)
*/
func (shared *viewportType) getTextWidth(viewportEntry *types.ViewportEntryType) int {
	if viewportEntry.IsBorderDrawn {
		return viewportEntry.Width - 2
	}
	return viewportEntry.Width
}

/*
getVisibleContentSize is a method which allows you to obtain how many columns and lines of text a viewport shows at
once, which is its size less the border when one is drawn and less a row or column for each of its scroll bars that
is visible. Neither value is ever negative. In addition, the following should be noted:

  - Scroll bar visibility is read as it stands, so call setViewportMaxScrollBarValues first when the content has
    just changed.

Example:

	visibleColumns, visibleLines := viewport.getVisibleContentSize("main", entry)
*/
func (shared *viewportType) getVisibleContentSize(layerAlias string, viewportEntry *types.ViewportEntryType) (int, int) {
	visibleColumns := shared.getTextWidth(viewportEntry)
	visibleLines := viewportEntry.Height
	if viewportEntry.IsBorderDrawn {
		visibleLines -= 2
	}
	if scrollbarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, viewportEntry.VerticalScrollbarAlias); isFound && scrollbarEntry.IsVisible {
		visibleColumns--
	}
	if scrollbarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, viewportEntry.HorizontalScrollbarAlias); isFound && scrollbarEntry.IsVisible {
		visibleLines--
	}
	return getClampedIndex(visibleColumns, 0, visibleColumns), getClampedIndex(visibleLines, 0, visibleLines)
}

/*
getMaxScrollPosition is a method which allows you to obtain the furthest column and line a viewport can be scrolled
to while still filling its visible area, measured in printed columns so markup tags and wide characters are counted
as they are drawn. Neither value is ever negative.

Example:

	maxColumn, maxLine := viewport.getMaxScrollPosition("main", entry)
*/
func (shared *viewportType) getMaxScrollPosition(layerAlias string, viewportEntry *types.ViewportEntryType) (int, int) {
	visibleColumns, visibleLines := shared.getVisibleContentSize(layerAlias, viewportEntry)
	maxColumn := shared.getMaxLineWidth(viewportEntry) - visibleColumns
	maxLine := len(viewportEntry.TextData) - visibleLines
	return getClampedIndex(maxColumn, 0, maxColumn), getClampedIndex(maxLine, 0, maxLine)
}

/*
scrollToLastLine is a method which allows you to scroll a viewport so its last line is shown on its bottom row, or to
its first line when all of its lines fit, and to update its scroll bars to match. The column it is scrolled to is
left unchanged.

Example:

	viewport.scrollToLastLine("main", "myVP", entry)
*/
func (shared *viewportType) scrollToLastLine(layerAlias string, viewportAlias string, viewportEntry *types.ViewportEntryType) {
	shared.setViewportMaxScrollBarValues(layerAlias, viewportAlias)
	_, maxLine := shared.getMaxScrollPosition(layerAlias, viewportEntry)
	viewportEntry.ViewportYLocation = maxLine
	shared.updateScrollbarBasedOnViewportViewport(layerAlias, viewportAlias)
}

/*
trimHistory is a method which allows you to trim the viewport's history to the maximum number of lines allowed. A
negative maximum, which can only be given when the viewport is added, is treated as no limit rather than used to slice
the history. When lines are removed, the scroll position moves up by the same number of lines, stopping at the first
line, so the lines that were shown stay in view.

Example:

	viewport.trimHistory(entry)
*/
func (shared *viewportType) trimHistory(viewportEntry *types.ViewportEntryType) {
	if !viewportEntry.IsHistoryEnabled || viewportEntry.MaxHistoryLines < 0 || len(viewportEntry.TextData) <= viewportEntry.MaxHistoryLines {
		return
	}
	removedLineCount := len(viewportEntry.TextData) - viewportEntry.MaxHistoryLines
	viewportEntry.TextData = viewportEntry.TextData[removedLineCount:]
	viewportEntry.ViewportYLocation = getClampedIndex(viewportEntry.ViewportYLocation-removedLineCount, 0, viewportEntry.ViewportYLocation)
}

/*
trimToVisible is a method which allows you to keep only the newest lines a viewport can show at once, measured
inside its border and any visible scroll bar, and to scroll it to its first line when lines are removed. A viewport
too small to show any line, such as a bordered one less than three rows high, keeps no lines.

Example:

	viewport.trimToVisible("main", entry)
*/
func (shared *viewportType) trimToVisible(layerAlias string, viewportEntry *types.ViewportEntryType) {
	_, visibleLines := shared.getVisibleContentSize(layerAlias, viewportEntry)
	if len(viewportEntry.TextData) > visibleLines {
		viewportEntry.TextData = viewportEntry.TextData[len(viewportEntry.TextData)-visibleLines:]
		viewportEntry.ViewportYLocation = 0
	}
}

/*
appendText is a method which allows you to add text to a viewport, continuing its last line unless that line has
been ended, in which case the text starts a new one. Every "\n" in the text ends a line, and when line wrapping is
enabled the text is wrapped exactly as the line builder wraps it, measured in printed columns from where the last
line already ends and keeping any style tag still open on it. The lines are then trimmed to the history limit, and
when history is disabled only the lines the viewport can show are kept. Empty text changes nothing. In addition, the
following should be noted:

  - A final "\n" is not stored as an empty line. It sets IsLineBreakPending instead, so the next text starts a new
    line without a blank row being shown first.

Example:

	viewport.appendText("main", "myVP", entry, "More text\n")
*/
func (shared *viewportType) appendText(layerAlias string, viewportAlias string, viewportEntry *types.ViewportEntryType, textToPrint string) {
	if textToPrint == "" {
		return
	}
	if len(viewportEntry.TextData) == 0 || viewportEntry.IsLineBreakPending {
		viewportEntry.TextData = append(viewportEntry.TextData, []rune{})
		viewportEntry.IsLineBreakPending = false
	}
	lastLineIndex := len(viewportEntry.TextData) - 1
	lineBuilder := newViewportLineBuilder(shared.getTextWidth(viewportEntry), viewportEntry.IsLinesWrapped, viewportEntry.TextData[lastLineIndex])
	lineBuilder.addText(stringformat.GetRunesFromString(textToPrint))
	lines := lineBuilder.getLines()
	if strings.HasSuffix(textToPrint, "\n") {
		lines = lines[:len(lines)-1]
		viewportEntry.IsLineBreakPending = true
	}
	viewportEntry.TextData = append(viewportEntry.TextData[:lastLineIndex], lines...)
	shared.trimHistory(viewportEntry)
	if !viewportEntry.IsHistoryEnabled {
		shared.setViewportMaxScrollBarValues(layerAlias, viewportAlias)
		shared.trimToVisible(layerAlias, viewportEntry)
	}
}

/*
updateScrollbarBasedOnViewportViewport is a method which allows you to update scrollbar positions based on the current
viewport scroll position.

Example:
    viewport.updateScrollbarBasedOnViewportViewport("main", "myVP")
*/
func (shared *viewportType) updateScrollbarBasedOnViewportViewport(layerAlias string, viewportAlias string) {
	viewportEntry := getViewport(layerAlias, viewportAlias)
	if viewportEntry == nil {
		return
	}

	// Update vertical scrollbar if it exists
	if viewportEntry.VerticalScrollbarAlias != "" {
		if scrollbarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, viewportEntry.VerticalScrollbarAlias); isFound {
			scrollbarEntry.ScrollValue = viewportEntry.ViewportYLocation
			scrollbar.computeHandlePositionByScrollValue(layerAlias, viewportEntry.VerticalScrollbarAlias)
		}
	}

	// Update horizontal scrollbar if it exists
	if viewportEntry.HorizontalScrollbarAlias != "" {
		if scrollbarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, viewportEntry.HorizontalScrollbarAlias); isFound {
			scrollbarEntry.ScrollValue = viewportEntry.ViewportXLocation
			scrollbar.computeHandlePositionByScrollValue(layerAlias, viewportEntry.HorizontalScrollbarAlias)
		}
	}
}

/*
getMaxHorizontalTextValue is a method which allows you to obtain the printed width in columns of the widest line in
the viewport, or zero when the viewport does not exist.

Example:

	width := viewport.getMaxHorizontalTextValue("main", "myVP")
*/
func (shared *viewportType) getMaxHorizontalTextValue(layerAlias string, viewportAlias string) int {
	viewportEntry := getViewport(layerAlias, viewportAlias)
	if viewportEntry == nil {
		return 0
	}
	return shared.getMaxLineWidth(viewportEntry)
}

/*
getMaxLineWidth is a method which allows you to obtain the printed width in columns of the widest line a viewport
holds. Markup tags take no columns, and wide characters take two.

Example:

	width := viewport.getMaxLineWidth(entry)
*/
func (shared *viewportType) getMaxLineWidth(viewportEntry *types.ViewportEntryType) int {
	maxWidth := 0
	for _, line := range viewportEntry.TextData {
		if lineWidth := getMarkupLineWidth(line); lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}
	return maxWidth
}

/*
setViewportMaxScrollBarValues is a method which allows you to calculate and set the maximum scroll values for the
viewport's scrollbars.

Example:
    viewport.setViewportMaxScrollBarValues("main", "myVP")
*/
func (shared *viewportType) setViewportMaxScrollBarValues(layerAlias string, viewportAlias string) {
	viewportEntry := getViewport(layerAlias, viewportAlias)
	if viewportEntry == nil {
		return
	}

	// Calculate initial effective width and height (accounting for border if present)
	effectiveWidth := viewportEntry.Width
	effectiveHeight := viewportEntry.Height
	if viewportEntry.IsBorderDrawn {
		effectiveWidth -= 2
		effectiveHeight -= 2
	}

	// Get scrollbar entries if they exist
	var vScrollBarEntry *types.ScrollbarEntryType
	var hScrollBarEntry *types.ScrollbarEntryType

	if viewportEntry.VerticalScrollbarAlias != "" {
		vScrollBarEntry, _ = scrollbarRegistry.Lookup(layerAlias, viewportEntry.VerticalScrollbarAlias)
	}

	if viewportEntry.HorizontalScrollbarAlias != "" {
		hScrollBarEntry, _ = scrollbarRegistry.Lookup(layerAlias, viewportEntry.HorizontalScrollbarAlias)
	}

	// Initial state - assume no scrollbars are visible
	isVerticalScrollbarVisible := false
	isHorizontalScrollbarVisible := false

	// Get content dimensions
	contentHeight := len(viewportEntry.TextData)
	contentWidth := shared.getMaxHorizontalTextValue(layerAlias, viewportAlias)

	// First pass - check if scrollbars are needed based on full content area
	if vScrollBarEntry != nil {
		if contentHeight > effectiveHeight {
			isVerticalScrollbarVisible = true
		}
	}

	if hScrollBarEntry != nil && !viewportEntry.IsLinesWrapped {
		if contentWidth > effectiveWidth {
			isHorizontalScrollbarVisible = true
		}
	}

	// Second pass - adjust effective area if scrollbars are visible and recheck
	if isVerticalScrollbarVisible {
		// Vertical scrollbar takes up 1 column of width
		effectiveWidth--

		// Recheck horizontal scrollbar with reduced width
		if hScrollBarEntry != nil && !viewportEntry.IsLinesWrapped && !isHorizontalScrollbarVisible {
			if contentWidth > effectiveWidth {
				isHorizontalScrollbarVisible = true
			}
		}
	}

	if isHorizontalScrollbarVisible {
		// Horizontal scrollbar takes up 1 row of height
		effectiveHeight--

		// Recheck vertical scrollbar with reduced height
		if vScrollBarEntry != nil && !isVerticalScrollbarVisible {
			if contentHeight > effectiveHeight {
				isVerticalScrollbarVisible = true
				// Vertical scrollbar now appears, reduce width
				effectiveWidth--
			}
		}
	}

	// Set vertical scrollbar properties
	if vScrollBarEntry != nil {
		maxScrollValue := 0
		if isVerticalScrollbarVisible {
			maxScrollValue = contentHeight - effectiveHeight
			vScrollBarEntry.IsEnabled = true
			vScrollBarEntry.IsVisible = true
		} else {
			vScrollBarEntry.IsEnabled = false
			vScrollBarEntry.IsVisible = false
		}
		vScrollBarEntry.MaxScrollValue = maxScrollValue
	}

	// Set horizontal scrollbar properties
	if hScrollBarEntry != nil {
		maxScrollValue := 0
		if viewportEntry.IsLinesWrapped {
			hScrollBarEntry.IsEnabled = false
			hScrollBarEntry.IsVisible = false
		} else if isHorizontalScrollbarVisible {
			maxScrollValue = contentWidth - effectiveWidth
			hScrollBarEntry.IsEnabled = true
			hScrollBarEntry.IsVisible = true
		} else {
			hScrollBarEntry.IsEnabled = false
			hScrollBarEntry.IsVisible = false
		}
		hScrollBarEntry.MaxScrollValue = maxScrollValue
	}
}

/*
Add is a method which allows you to add a new viewport to the specified layer.

Example:
    vp := viewport.Add("main", "myVP", style, 0, 0, 40, 10, true, true, 100)
*/
func (shared *viewportType) Add(layerAlias string, viewportAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isLinesWrapped bool, isBorderDrawn bool, maxHistoryLines int) ViewportInstanceType {
	var viewportInstance ViewportInstanceType

	// Create a new viewport entry
	viewportEntry := types.NewViewportEntry()
	viewportEntry.Alias = viewportAlias
	viewportEntry.StyleEntry = styleEntry
	viewportEntry.XLocation = xLocation
	viewportEntry.YLocation = yLocation
	viewportEntry.Width = width
	viewportEntry.Height = height
	viewportEntry.IsLinesWrapped = isLinesWrapped
	viewportEntry.IsBorderDrawn = isBorderDrawn
	viewportEntry.MaxHistoryLines = maxHistoryLines

	// Add the viewport to memory
	viewportRegistry.Add(layerAlias, viewportAlias, &viewportEntry)

	// Create scrollbars if needed
	// For horizontal scrollbar
	horizontalScrollbarAlias := viewportAlias + "_hscroll"
	scrollbarStyleEntry := styleEntry

	// Calculate scrollbar positions and lengths
	// For horizontal scrollbar - position it inside the frame at the bottom
	scrollbarXLocation := xLocation
	if isBorderDrawn {
		scrollbarXLocation++
	}
	scrollbarYLocation := yLocation + height - 1
	if isBorderDrawn {
		scrollbarYLocation--
	}
	horizontalScrollbarLength := width
	if isBorderDrawn {
		horizontalScrollbarLength -= 2
	}
	// Add horizontal scrollbar
	scrollbar.Add(layerAlias, horizontalScrollbarAlias, scrollbarStyleEntry, scrollbarXLocation, scrollbarYLocation, horizontalScrollbarLength, 0, 0, 1, true)
	viewportEntry.HorizontalScrollbarAlias = horizontalScrollbarAlias

	// Set parent control information for horizontal scrollbar
	if hScrollbarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, horizontalScrollbarAlias); isFound {
		hScrollbarEntry.ParentControlAlias = viewportAlias
		hScrollbarEntry.ParentControlType = constants.CellTypeTextbox // Viewports use textbox cell type
	}

	// For vertical scrollbar - position it inside the frame on the right
	verticalScrollbarAlias := viewportAlias + "_vscroll"

	scrollbarXLocation = xLocation + width - 1
	if isBorderDrawn {
		scrollbarXLocation--
	}
	scrollbarYLocation = yLocation
	if isBorderDrawn {
		scrollbarYLocation++
	}
	verticalScrollbarLength := height
	if isBorderDrawn {
		verticalScrollbarLength -= 2
	}
	// Add vertical scrollbar
	scrollbar.Add(layerAlias, verticalScrollbarAlias, scrollbarStyleEntry, scrollbarXLocation, scrollbarYLocation, verticalScrollbarLength, 0, 0, 1, false)
	viewportEntry.VerticalScrollbarAlias = verticalScrollbarAlias

	// Set parent control information for vertical scrollbar
	if vScrollbarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, verticalScrollbarAlias); isFound {
		vScrollbarEntry.ParentControlAlias = viewportAlias
		vScrollbarEntry.ParentControlType = constants.CellTypeTextbox // Viewports use textbox cell type
	}

	// Set up the viewport instance
	viewportInstance.layerAlias = layerAlias
	viewportInstance.controlAlias = viewportAlias
	viewportInstance.controlType = constants.TYPE_VIEWPORT

	return viewportInstance
}

/*
Delete is a method which allows you to remove a viewport from the specified layer.

Example:
    viewport.Delete("main", "myVP")
*/
func (shared *viewportType) Delete(layerAlias string, viewportAlias string) {
	viewportEntry := getViewport(layerAlias, viewportAlias)
	if viewportEntry == nil {
		return
	}

	// Delete associated scrollbars
	if viewportEntry.HorizontalScrollbarAlias != "" {
		var scrollbarInstance ScrollbarInstanceType
		scrollbarInstance.layerAlias = layerAlias
		scrollbarInstance.controlAlias = viewportEntry.HorizontalScrollbarAlias
		scrollbarInstance.delete()
	}

	if viewportEntry.VerticalScrollbarAlias != "" {
		var scrollbarInstance ScrollbarInstanceType
		scrollbarInstance.layerAlias = layerAlias
		scrollbarInstance.controlAlias = viewportEntry.VerticalScrollbarAlias
		scrollbarInstance.delete()
	}

	// Delete the viewport
	viewportRegistry.Remove(layerAlias, viewportAlias)
	clearStaleControlReferences(layerAlias, viewportAlias, constants.CellTypeTextbox)
}

/*
DeleteAll is a method which allows you to remove all viewports from the specified layer.

Example:
    viewport.DeleteAll("main")
*/
func (shared *viewportType) DeleteAll(layerAlias string) {
	// Get all entries for the layer
	entries := viewportRegistry.GetAllEntries(layerAlias)
	// Remove each entry
	for _, entry := range entries {
		viewportRegistry.Remove(layerAlias, entry.Alias)
		clearStaleControlReferences(layerAlias, entry.Alias, constants.CellTypeTextbox)
	}
}

/*
drawOnLayer is a method which allows you to draw all viewports on the specified layer.

Example:
    viewport.drawOnLayer(layer)
*/
func (shared *viewportType) drawOnLayer(layerEntry types.LayerEntryType) {
	// Get all entries for the layer
	entries := viewportRegistry.GetAllEntries(layerEntry.LayerAlias)
	// Draw each viewport
	for _, entry := range entries {
		shared.draw(&layerEntry, entry.Alias)
	}
}

/*
draw is a method which allows you to draw a single viewport on the specified layer.

Example:
    viewport.drawViewport(layer, "myVP")
*/
func (shared *viewportType) draw(layerEntry *types.LayerEntryType, viewportAlias string) {
	viewportEntry := getViewport(layerEntry.LayerAlias, viewportAlias)
	if viewportEntry == nil || !viewportEntry.IsVisible {
		return
	}

	// Get the viewport's style and attribute
	styleEntry := viewportEntry.StyleEntry
	attributeEntry := types.NewAttributeEntry()

	// Draw border if needed
	if viewportEntry.IsBorderDrawn {
		shared.drawBorder(layerEntry, styleEntry, attributeEntry, viewportEntry.XLocation, viewportEntry.YLocation, viewportEntry.Width, viewportEntry.Height, false)
	}

	// Draw viewport content
	shared.drawContent(layerEntry, viewportAlias, styleEntry, attributeEntry, viewportEntry.XLocation, viewportEntry.YLocation, viewportEntry.Width, viewportEntry.Height)

	shared.updateScrollbars(layerEntry, viewportAlias, styleEntry, attributeEntry, viewportEntry.XLocation, viewportEntry.YLocation, viewportEntry.Width, viewportEntry.Height)
	scrollbar.drawOnLayerByAlias(layerEntry, viewportEntry.HorizontalScrollbarAlias)
	scrollbar.drawOnLayerByAlias(layerEntry, viewportEntry.VerticalScrollbarAlias)
}

/*
drawBorder is a method which allows you to draw a border around the viewport.

Example:
    viewport.drawBorder(layer, style, attr, 0, 0, 40, 10, false)
*/
func (shared *viewportType) drawBorder(layerEntry *types.LayerEntryType, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, height int, isDoubleLine bool) {
	// This draws exactly what LayerInstanceType.DrawBorder would, without taking the UI lock already held for drawing.
	drawBorder(layerRegistry.Get(layerEntry.LayerAlias), styleEntry, types.NewAttributeEntry(), xLocation, yLocation, width, height, false)
}

/*
drawContent is a method which allows you to draw the lines of a viewport that fall within its scroll position, inside
its border and scroll bars. In addition, the following should be noted:

  - A negative scroll position is drawn as zero rather than used as an index, so drawing cannot panic on the event
    goroutine.

Example:

	viewport.drawContent(layer, "myVP", style, attr, 1, 1, 38, 8)
*/
func (shared *viewportType) drawContent(layerEntry *types.LayerEntryType, viewportAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, height int) {
	viewportEntry := getViewport(layerEntry.LayerAlias, viewportAlias)
	if viewportEntry == nil {
		return
	}

	// Calculate content area (accounting for border if present)
	contentXLocation := xLocation
	contentYLocation := yLocation
	contentWidth := width
	contentHeight := height

	if viewportEntry.IsBorderDrawn {
		contentXLocation++
		contentYLocation++
		contentWidth -= 2
		contentHeight -= 2
	}

	// Further adjust content area for scrollbars inside the frame
	if viewportEntry.VerticalScrollbarAlias != "" {
		if scrollbarEntry, isFound := scrollbarRegistry.Lookup(layerEntry.LayerAlias, viewportEntry.VerticalScrollbarAlias); isFound && scrollbarEntry.IsVisible {
			contentWidth-- // Reduce width by 1 for vertical scrollbar
		}
	}

	if viewportEntry.HorizontalScrollbarAlias != "" {
		if scrollbarEntry, isFound := scrollbarRegistry.Lookup(layerEntry.LayerAlias, viewportEntry.HorizontalScrollbarAlias); isFound && scrollbarEntry.IsVisible {
			contentHeight-- // Reduce height by 1 for horizontal scrollbar
		}
	}

	// Set the attribute entry for the viewport
	localAttributeEntry := types.NewAttributeEntry(&attributeEntry)
	localAttributeEntry.CellType = constants.CellTypeTextbox // Use textbox type for viewport
	localAttributeEntry.CellControlAlias = viewportAlias

	// Create a style attribute entry from styleEntry.Textbox for default text attributes
	styleAttributeEntry := types.NewAttributeEntry(&localAttributeEntry)
	styleAttributeEntry.ForegroundColor = styleEntry.Textbox.ForegroundColor
	styleAttributeEntry.BackgroundColor = styleEntry.Textbox.BackgroundColor

	// Set background color based on transparency
	if !viewportEntry.IsTransparent {
		// Fill the background with the style's background color
		for y := 0; y < contentHeight; y++ {
			for x := 0; x < contentWidth; x++ {
				// Create a space character with the background color
				spaceChar := []rune{' '}
				layer.printLayer(layerEntry, localAttributeEntry, contentXLocation+x, contentYLocation+y, spaceChar)
			}
		}
	}

	// Draw each line of text in the viewport
	firstVisibleColumn := getClampedIndex(viewportEntry.ViewportXLocation, 0, viewportEntry.ViewportXLocation)
	for y := 0; y < contentHeight; y++ {
		textDataY := y + viewportEntry.ViewportYLocation
		if textDataY >= 0 && textDataY < len(viewportEntry.TextData) {
			shared.drawLine(layerEntry, viewportEntry.TextData[textDataY], styleAttributeEntry, contentXLocation, contentYLocation+y, firstVisibleColumn, contentWidth)
		}
	}
}

/*
drawLine is a method which allows you to draw one stored viewport line, scrolled so that printed column
firstVisibleColumn of the line lands at xLocation, showing at most visibleColumns columns. Markup tags in the line
switch the style the following characters are drawn with and take no columns, so scrolling and clipping are measured
in what is actually printed. In addition, the following should be noted:

  - A wide character that would straddle the left or right edge of the visible area is not drawn, leaving that
    column to the viewport's background rather than drawing half a character or drawing past the edge.

  - Tags before the first visible column still apply, so text scrolled into view keeps the style it was given.

Example:

	viewport.drawLine(layerEntry, []rune("Go {{red}}Red{{/}} end"), styleAttributeEntry, 1, 1, 0, 38)
*/
func (shared *viewportType) drawLine(layerEntry *types.LayerEntryType, line []rune, defaultAttributeEntry types.AttributeEntryType, xLocation int, yLocation int, firstVisibleColumn int, visibleColumns int) {
	lastVisibleColumn := firstVisibleColumn + visibleColumns
	currentAttributeEntry := defaultAttributeEntry
	lineColumn := 0
	for currentIndex := 0; currentIndex < len(line) && lineColumn < lastVisibleColumn; currentIndex++ {
		if tagEndIndex := findMarkupTagEnd(line, currentIndex); tagEndIndex >= 0 {
			currentAttributeEntry = getDialogAttributeEntry(getAttributeTagContent(string(line[currentIndex:tagEndIndex+1])), defaultAttributeEntry)
			currentIndex = tagEndIndex
			continue
		}
		columnWidth := stringformat.GetWidthOfRuneWhenPrinted(line[currentIndex])
		if lineColumn >= firstVisibleColumn && lineColumn+columnWidth <= lastVisibleColumn {
			layer.printLayer(layerEntry, currentAttributeEntry, xLocation+lineColumn-firstVisibleColumn, yLocation, []rune{line[currentIndex]})
		}
		lineColumn += columnWidth
	}
}

/*
updateScrollbars is a method which allows you to update and draw scrollbars for the viewport.

Example:
    viewport.updateScrollbars(layer, "myVP", style, attr, 0, 0, 40, 10)
*/
func (shared *viewportType) updateScrollbars(layerEntry *types.LayerEntryType, viewportAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, height int) {
	viewportEntry := getViewport(layerEntry.LayerAlias, viewportAlias)
	if viewportEntry == nil {
		return
	}

	// Check if scrollbars exist and are visible
	var hScrollbarEntry, vScrollbarEntry *types.ScrollbarEntryType
	var isHorizontalScrollbarVisible, isVerticalScrollbarVisible bool

	if viewportEntry.HorizontalScrollbarAlias != "" {
		if entry, isFound := scrollbarRegistry.Lookup(layerEntry.LayerAlias, viewportEntry.HorizontalScrollbarAlias); isFound {
			hScrollbarEntry = entry
			isHorizontalScrollbarVisible = hScrollbarEntry.IsVisible
		}
	}

	if viewportEntry.VerticalScrollbarAlias != "" {
		if entry, isFound := scrollbarRegistry.Lookup(layerEntry.LayerAlias, viewportEntry.VerticalScrollbarAlias); isFound {
			vScrollbarEntry = entry
			isVerticalScrollbarVisible = vScrollbarEntry.IsVisible
		}
	}

	// If both scrollbars are visible, adjust their lengths and fill the corner
	if isHorizontalScrollbarVisible && isVerticalScrollbarVisible {
		// Adjust horizontal scrollbar length
		if hScrollbarEntry != nil {
			if viewportEntry.IsBorderDrawn {
				hScrollbarEntry.Length = viewportEntry.Width - 3 // -2 for border, -1 for corner
			} else {
				hScrollbarEntry.Length = viewportEntry.Width - 1 // -1 for corner
			}
		}

		// Adjust vertical scrollbar length
		if vScrollbarEntry != nil {
			if viewportEntry.IsBorderDrawn {
				vScrollbarEntry.Length = viewportEntry.Height - 3 // -2 for border, -1 for corner
			} else {
				vScrollbarEntry.Length = viewportEntry.Height - 1 // -1 for corner
			}
		}

		// Fill the corner with the background color
		// Calculate the corner position
		cornerX := xLocation + width - 1
		cornerY := yLocation + height - 1
		if viewportEntry.IsBorderDrawn {
			cornerX--
			cornerY--
		}

		// Draw a space character with the background color in the corner
		localAttributeEntry := types.NewAttributeEntry()
		localAttributeEntry.BackgroundColor = styleEntry.Textbox.BackgroundColor
		spaceChar := []rune{' '}
		layer.printLayer(layerEntry, localAttributeEntry, cornerX, cornerY, spaceChar)
	} else {
		// Restore original lengths if only one scrollbar is visible
		if hScrollbarEntry != nil && isHorizontalScrollbarVisible && !isVerticalScrollbarVisible {
			if viewportEntry.IsBorderDrawn {
				hScrollbarEntry.Length = viewportEntry.Width - 2 // -2 for border
			} else {
				hScrollbarEntry.Length = viewportEntry.Width
			}
		}
		if vScrollbarEntry != nil && isVerticalScrollbarVisible && !isHorizontalScrollbarVisible {
			if viewportEntry.IsBorderDrawn {
				vScrollbarEntry.Length = viewportEntry.Height - 2 // -2 for border
			} else {
				vScrollbarEntry.Length = viewportEntry.Height
			}
		}
	}
}

/*
updateViewport is a method which allows you to update the viewport scroll position based on current scrollbar values.

Example:
    viewport.updateViewport(entry, "main", "myVP")
*/
func (shared *viewportType) updateViewport(viewportEntry *types.ViewportEntryType, layerAlias string, viewportAlias string) {
	// This method is called when scrollbars are moved
	// It updates the viewport's position based on scrollbar values

	// Update viewport position based on scrollbar values
	if viewportEntry.HorizontalScrollbarAlias != "" {
		if scrollbarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, viewportEntry.HorizontalScrollbarAlias); isFound {
			viewportEntry.ViewportXLocation = scrollbarEntry.ScrollValue
		}
	}

	if viewportEntry.VerticalScrollbarAlias != "" {
		if scrollbarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, viewportEntry.VerticalScrollbarAlias); isFound {
			viewportEntry.ViewportYLocation = scrollbarEntry.ScrollValue
		}
	}
}

/*
updateMouseEvent is a method which allows you to handle mouse events for viewports, including wheel scrolling.

Example:
    updateRequired := viewport.updateMouseEvent()
*/
func (shared *viewportType) updateMouseEvent() bool {
	// Handle mouse wheel events for scrolling
	isScreenUpdateRequired := false
	mouseXLocation, mouseYLocation, _, wheelState := GetMouseStatus()

	// Check for scrollbar events and update viewports
	isScreenUpdateRequired = shared.checkScrollbarEvents() || isScreenUpdateRequired

	if wheelState != "" {
		// Find the viewport under the mouse cursor
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		if characterEntry.AttributeEntry.CellType == constants.CellTypeTextbox {
			layerAlias := characterEntry.LayerAlias
			viewportAlias := characterEntry.AttributeEntry.CellControlAlias
			viewportEntry := getViewport(layerAlias, viewportAlias)

			if viewportEntry != nil {
				// Handle wheel scrolling
				if wheelState == "Up" && viewportEntry.VerticalScrollbarAlias != "" {
					// Scroll up
					if viewportEntry.ViewportYLocation > 0 {
						viewportEntry.ViewportYLocation--
						shared.updateScrollbarBasedOnViewportViewport(layerAlias, viewportAlias)
						isScreenUpdateRequired = true
					}
				} else if wheelState == "Down" && viewportEntry.VerticalScrollbarAlias != "" {
					// Scroll down
					_, maxLine := shared.getMaxScrollPosition(layerAlias, viewportEntry)
					if viewportEntry.ViewportYLocation < maxLine {
						viewportEntry.ViewportYLocation++
						shared.updateScrollbarBasedOnViewportViewport(layerAlias, viewportAlias)
						isScreenUpdateRequired = true
					}
				} else if wheelState == "Left" && viewportEntry.HorizontalScrollbarAlias != "" {
					// Scroll left
					if viewportEntry.ViewportXLocation > 0 {
						viewportEntry.ViewportXLocation--
						shared.updateScrollbarBasedOnViewportViewport(layerAlias, viewportAlias)
						isScreenUpdateRequired = true
					}
				} else if wheelState == "Right" && viewportEntry.HorizontalScrollbarAlias != "" {
					// Scroll right
					maxColumn, _ := shared.getMaxScrollPosition(layerAlias, viewportEntry)
					if viewportEntry.ViewportXLocation < maxColumn {
						viewportEntry.ViewportXLocation++
						shared.updateScrollbarBasedOnViewportViewport(layerAlias, viewportAlias)
						isScreenUpdateRequired = true
					}
				}
			}
		}
	}

	return isScreenUpdateRequired
}

/*
checkScrollbarEvents is a method which allows you to check for scrollbar events and update all viewports accordingly.

Example:
    updateRequired := viewport.checkScrollbarEvents()
*/
func (shared *viewportType) checkScrollbarEvents() bool {
	isScreenUpdateRequired := false

	// Get all layers
	sortedLayerAliasSlice := layer.GetSortedLayerMemoryAliasSlice()

	// Iterate through all layers
	for _, layerAliasPair := range sortedLayerAliasSlice {
		layerAlias := layerAliasPair.Key

		// Get all viewports in the layer
		viewports := viewportRegistry.GetAllEntries(layerAlias)

		// Iterate through all viewports
		for _, viewportEntry := range viewports {
			// Update the viewport position based on scrollbar values
			oldXLocation := viewportEntry.ViewportXLocation
			oldYLocation := viewportEntry.ViewportYLocation

			// Update the viewport
			shared.updateViewport(viewportEntry, layerAlias, viewportEntry.Alias)

			// Check if the viewport position changed
			if oldXLocation != viewportEntry.ViewportXLocation || oldYLocation != viewportEntry.ViewportYLocation {
				isScreenUpdateRequired = true
			}
		}
	}

	return isScreenUpdateRequired
}
