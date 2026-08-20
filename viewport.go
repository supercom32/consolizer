package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"math"

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
var Viewports = memory.NewControlMemoryManager[types.ViewportEntryType]()

/*
GetViewport is a method which allows you to retrieve a viewport entry by its layer and viewport alias.

Example:
    entry := GetViewport("main", "myVP")
*/
func GetViewport(layerAlias string, viewportAlias string) *types.ViewportEntryType {
	return Viewports.Get(layerAlias, viewportAlias)
}

/*
IsViewportExists is a method which allows you to check if a viewport exists in the specified layer.

Example:
    exists := IsViewportExists("main", "myVP")
*/
func IsViewportExists(layerAlias string, viewportAlias string) bool {
	return Viewports.IsExists(layerAlias, viewportAlias)
}

/*
Println is a method which allows you to add text to the viewport followed by a new line.

Example:
    vp.Println("Hello World!")
*/
func (shared *ViewportInstanceType) Println(text string) *ViewportInstanceType {
	viewportEntry := GetViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	// Process the text and add it to the viewport
	viewport.printText(viewportEntry, text)

	// Update scrollbars
	viewport.setViewportMaxScrollBarValues(shared.layerAlias, shared.controlAlias)
	viewport.updateScrollbarBasedOnViewportViewport(shared.layerAlias, shared.controlAlias)

	return shared
}

/*
Print is a method which allows you to add text to the viewport by appending to the last line. If the text contains
control characters like \n, it will start a new line.

Example:
    vp.Print("Appending some text...")
*/
func (shared *ViewportInstanceType) Print(text string) *ViewportInstanceType {
	viewportEntry := GetViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	// If there's no text data yet, just use printText
	if len(viewportEntry.TextData) == 0 {
		viewport.printText(viewportEntry, text)
	} else {
		// Process the text to handle control characters
		viewport.appendToLastLine(viewportEntry, text)
	}

	// Update scrollbars
	viewport.setViewportMaxScrollBarValues(shared.layerAlias, shared.controlAlias)
	viewport.updateScrollbarBasedOnViewportViewport(shared.layerAlias, shared.controlAlias)

	return shared
}

/*
SetContent is a method which allows you to clear the viewport and set its content to the given text. This method resets
the viewport scroll position to the top-left.

Example:
    vp.SetContent("New content starts here.")
*/
func (shared *ViewportInstanceType) SetContent(text string) *ViewportInstanceType {
	viewportEntry := GetViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	// Clear existing data and reset viewport position.
	viewportEntry.TextData = [][]rune{}
	viewportEntry.ViewportXLocation = 0
	viewportEntry.ViewportYLocation = 0

	// Process and add the new text.
	viewport.printText(viewportEntry, text, false)

	// Update scrollbar max values to get correct dimensions and max scroll value.
	viewport.setViewportMaxScrollBarValues(shared.layerAlias, shared.controlAlias)

	// Set the viewport's Y location to the maximum possible scroll value.
	vScrollbar := ScrollBars.Get(shared.layerAlias, viewportEntry.VerticalScrollbarAlias)
	if vScrollbar != nil {
		viewportEntry.ViewportYLocation = vScrollbar.MaxScrollValue
	}

	// Update the scrollbar handle to reflect the new viewport position.
	viewport.updateScrollbarBasedOnViewportViewport(shared.layerAlias, shared.controlAlias)
	return shared
}

/*
SetViewport is a method which allows you to set the viewport scroll position.

Example:
    vp.SetViewport(0, 10)
*/
func (shared *ViewportInstanceType) SetViewport(xLocation int, yLocation int) *ViewportInstanceType {
	viewportEntry := GetViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewportEntry.ViewportXLocation = xLocation
	viewportEntry.ViewportYLocation = yLocation

	// Update scrollbars
	viewport.setViewportMaxScrollBarValues(shared.layerAlias, shared.controlAlias)
	viewport.updateScrollbarBasedOnViewportViewport(shared.layerAlias, shared.controlAlias)

	return shared
}

/*
ScrollToBottom is a method which allows you to scroll the viewport to the bottom of its content.

Example:
    vp.ScrollToBottom()
*/
func (shared *ViewportInstanceType) ScrollToBottom() *ViewportInstanceType {
	viewportEntry := GetViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}
	viewport.setViewportMaxScrollBarValues(shared.layerAlias, shared.controlAlias)
	vScrollbar := ScrollBars.Get(shared.layerAlias, viewportEntry.VerticalScrollbarAlias)
	viewportEntry.ViewportYLocation = vScrollbar.MaxScrollValue
	viewport.updateScrollbarBasedOnViewportViewport(shared.layerAlias, shared.controlAlias)
	return shared
}

/*
SetMaxHistoryLines is a method which allows you to set the maximum number of lines to keep in history.

Example:
    vp.SetMaxHistoryLines(500)
*/
func (shared *ViewportInstanceType) SetMaxHistoryLines(maxLines int) *ViewportInstanceType {
	viewportEntry := GetViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewportEntry.MaxHistoryLines = maxLines

	// Trim history if needed
	viewport.trimHistory(viewportEntry)

	return shared
}

/*
EnableHistory is a method which allows you to enable or disable history tracking for the viewport.

Example:
    vp.EnableHistory(true)
*/
func (shared *ViewportInstanceType) EnableHistory(enabled bool) *ViewportInstanceType {
	viewportEntry := GetViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewportEntry.IsHistoryEnabled = enabled

	// If history is disabled, trim to just what's visible
	if !enabled {
		viewport.trimToVisible(viewportEntry)
	}

	return shared
}

/*
Clear is a method which allows you to clear all text from the viewport.

Example:
    vp.Clear()
*/
func (shared *ViewportInstanceType) Clear() *ViewportInstanceType {
	viewportEntry := GetViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewportEntry.TextData = [][]rune{}
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
	viewportEntry := GetViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewportEntry.IsTransparent = isTransparent
	return shared
}

/*
printText is a method which allows you to process and add text to the viewport.

Example:
    viewport.printText(entry, "New line of text", true)
*/
func (shared *viewportType) printText(viewportEntry *types.ViewportEntryType, textToPrint string, autoScrollToBottom ...bool) {
	// Get the effective width (accounting for border if present)
	effectiveWidth := viewportEntry.Width
	if viewportEntry.IsBorderDrawn {
		effectiveWidth -= 2
	}

	// Process the text with markup and word wrapping
	lines := shared.processTextWithMarkup(textToPrint, effectiveWidth, viewportEntry.IsLinesWrapped)

	// Add the processed lines to the viewport's text data
	for _, line := range lines {
		viewportEntry.TextData = append(viewportEntry.TextData, line)
	}

	// Trim history if needed
	shared.trimHistory(viewportEntry)

	// By default, scroll to the bottom. This is the behavior for Print/Println.
	shouldAutoScroll := true
	if len(autoScrollToBottom) > 0 {
		shouldAutoScroll = autoScrollToBottom[0]
	}

	if shouldAutoScroll {
		// Update viewport position to show the latest content.
		if len(viewportEntry.TextData) > viewportEntry.Height-2 { // -2 for border
			viewportEntry.ViewportYLocation = len(viewportEntry.TextData) - (viewportEntry.Height - 2)
		}
	}
}

/*
processTextWithMarkup is a method which allows you to process text with markup and word wrapping.

Example:
    lines := viewport.processTextWithMarkup("Wrapped text content", 40, true)
*/
func (shared *viewportType) processTextWithMarkup(textToPrint string, widthOfLineInCharacters int, isLinesWrapped bool) [][]rune {
	var result [][]rune
	var currentLine []rune

	arrayOfRunes := stringformat.GetRunesFromString(textToPrint)
	cursorXLocation := 0

	for currentCharacterIndex := 0; currentCharacterIndex < len(arrayOfRunes); currentCharacterIndex++ {
		currentCharacter := arrayOfRunes[currentCharacterIndex]

		// Add the character to the current line
		currentLine = append(currentLine, currentCharacter)
		cursorXLocation++

		// Check for word wrapping if lines should be wrapped
		lengthOfNextWord := 0
		if isLinesWrapped && currentCharacter == ' ' {
			lengthOfNextWord = getLengthOfNextWord(textToPrint, currentCharacterIndex+1)
		}

		// Handle newlines
		if currentCharacter == '\n' {
			// Remove the newline character
			currentLine = currentLine[:len(currentLine)-1]
			// Add the current line to the result
			result = append(result, currentLine)
			// Start a new line
			currentLine = []rune{}
			cursorXLocation = 0
			continue
		}

		// Check if we need to wrap to the next line (only if lines should be wrapped)
		if isLinesWrapped && cursorXLocation+lengthOfNextWord > widthOfLineInCharacters {
			// Add the current line to the result
			result = append(result, currentLine)
			// Start a new line
			currentLine = []rune{}
			cursorXLocation = 0
		}
	}

	// Add the last line if it's not empty
	if len(currentLine) > 0 {
		result = append(result, currentLine)
	}

	return result
}

/*
trimHistory is a method which allows you to trim the viewport's history to the maximum number of lines allowed.

Example:
    viewport.trimHistory(entry)
*/
func (shared *viewportType) trimHistory(viewportEntry *types.ViewportEntryType) {
	if viewportEntry.IsHistoryEnabled && len(viewportEntry.TextData) > viewportEntry.MaxHistoryLines {
		// Keep only the most recent MaxHistoryLines
		viewportEntry.TextData = viewportEntry.TextData[len(viewportEntry.TextData)-viewportEntry.MaxHistoryLines:]
		// Adjust viewport position
		if viewportEntry.ViewportYLocation > 0 {
			viewportEntry.ViewportYLocation = int(math.Max(0, float64(viewportEntry.ViewportYLocation-(len(viewportEntry.TextData)-viewportEntry.MaxHistoryLines))))
		}
	}
}

/*
trimToVisible is a method which allows you to trim the viewport's text data to only what is currently visible.

Example:
    viewport.trimToVisible(entry)
*/
func (shared *viewportType) trimToVisible(viewportEntry *types.ViewportEntryType) {
	effectiveHeight := viewportEntry.Height
	if viewportEntry.IsBorderDrawn {
		effectiveHeight -= 2
	}

	if len(viewportEntry.TextData) > effectiveHeight {
		// Keep only the most recent lines that fit in the viewport
		viewportEntry.TextData = viewportEntry.TextData[len(viewportEntry.TextData)-effectiveHeight:]
		viewportEntry.ViewportYLocation = 0
	}
}

/*
appendToLastLine is a method which allows you to append text to the last line of the viewport. If the text contains
control characters like \n, it will start a new line.

Example:
    viewport.appendToLastLine(entry, " More text")
*/
func (shared *viewportType) appendToLastLine(viewportEntry *types.ViewportEntryType, textToPrint string) {
	// Get the effective width (accounting for border if present)
	effectiveWidth := viewportEntry.Width
	if viewportEntry.IsBorderDrawn {
		effectiveWidth -= 2
	}

	// Get the runes from the string
	arrayOfRunes := stringformat.GetRunesFromString(textToPrint)

	// Process each character
	currentLine := viewportEntry.TextData[len(viewportEntry.TextData)-1]

	for i := 0; i < len(arrayOfRunes); i++ {
		currentChar := string(arrayOfRunes[i])

		// Handle newline character
		if currentChar == "\n" {
			// Add the current line to the viewport
			viewportEntry.TextData[len(viewportEntry.TextData)-1] = currentLine
			// Start a new line
			currentLine = []rune{}
			viewportEntry.TextData = append(viewportEntry.TextData, currentLine)
		} else {
			// Append the character to the current line
			currentLine = append(currentLine, arrayOfRunes[i])

			// Check if we need to wrap to the next line (only if lines should be wrapped)
			if viewportEntry.IsLinesWrapped && len(currentLine) >= effectiveWidth {
				// Update the last line
				viewportEntry.TextData[len(viewportEntry.TextData)-1] = currentLine
				// Start a new line
				currentLine = []rune{}
				viewportEntry.TextData = append(viewportEntry.TextData, currentLine)
			}
		}
	}

	// Update the last line with any remaining characters
	if len(currentLine) > 0 {
		viewportEntry.TextData[len(viewportEntry.TextData)-1] = currentLine
	}

	// Trim history if needed
	shared.trimHistory(viewportEntry)

	// Update viewport position to show the latest content
	if len(viewportEntry.TextData) > viewportEntry.Height {
		if viewportEntry.IsHistoryEnabled {
			viewportEntry.ViewportYLocation = len(viewportEntry.TextData) - viewportEntry.Height
		} else {
			// If history is disabled, keep only what's visible
			shared.trimToVisible(viewportEntry)
		}
	}
}

/*
updateScrollbarBasedOnViewportViewport is a method which allows you to update scrollbar positions based on the current
viewport scroll position.

Example:
    viewport.updateScrollbarBasedOnViewportViewport("main", "myVP")
*/
func (shared *viewportType) updateScrollbarBasedOnViewportViewport(layerAlias string, viewportAlias string) {
	viewportEntry := GetViewport(layerAlias, viewportAlias)
	if viewportEntry == nil {
		return
	}

	// Update vertical scrollbar if it exists
	if viewportEntry.VerticalScrollbarAlias != "" && ScrollBars.IsExists(layerAlias, viewportEntry.VerticalScrollbarAlias) {
		scrollbarEntry := ScrollBars.Get(layerAlias, viewportEntry.VerticalScrollbarAlias)
		scrollbarEntry.ScrollValue = viewportEntry.ViewportYLocation
		scrollbar.computeHandlePositionByScrollValue(layerAlias, viewportEntry.VerticalScrollbarAlias)
	}

	// Update horizontal scrollbar if it exists
	if viewportEntry.HorizontalScrollbarAlias != "" && ScrollBars.IsExists(layerAlias, viewportEntry.HorizontalScrollbarAlias) {
		scrollbarEntry := ScrollBars.Get(layerAlias, viewportEntry.HorizontalScrollbarAlias)
		scrollbarEntry.ScrollValue = viewportEntry.ViewportXLocation
		scrollbar.computeHandlePositionByScrollValue(layerAlias, viewportEntry.HorizontalScrollbarAlias)
	}
}

/*
getMaxHorizontalTextValue is a method which allows you to obtain the maximum width of text in the viewport.

Example:
    width := viewport.getMaxHorizontalTextValue("main", "myVP")
*/
func (shared *viewportType) getMaxHorizontalTextValue(layerAlias string, viewportAlias string) int {
	viewportEntry := GetViewport(layerAlias, viewportAlias)
	if viewportEntry == nil {
		return 0
	}

	maxWidth := 0
	for _, line := range viewportEntry.TextData {
		if len(line) > maxWidth {
			maxWidth = len(line)
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
	viewportEntry := GetViewport(layerAlias, viewportAlias)
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

	if viewportEntry.VerticalScrollbarAlias != "" && ScrollBars.IsExists(layerAlias, viewportEntry.VerticalScrollbarAlias) {
		vScrollBarEntry = ScrollBars.Get(layerAlias, viewportEntry.VerticalScrollbarAlias)
	}

	if viewportEntry.HorizontalScrollbarAlias != "" && ScrollBars.IsExists(layerAlias, viewportEntry.HorizontalScrollbarAlias) {
		hScrollBarEntry = ScrollBars.Get(layerAlias, viewportEntry.HorizontalScrollbarAlias)
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
	Viewports.Add(layerAlias, viewportAlias, &viewportEntry)

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
	hScrollbarEntry := ScrollBars.Get(layerAlias, horizontalScrollbarAlias)
	if hScrollbarEntry != nil {
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
	vScrollbarEntry := ScrollBars.Get(layerAlias, verticalScrollbarAlias)
	if vScrollbarEntry != nil {
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
	viewportEntry := GetViewport(layerAlias, viewportAlias)
	if viewportEntry == nil {
		return
	}

	// Delete associated scrollbars
	if viewportEntry.HorizontalScrollbarAlias != "" {
		var scrollbarInstance ScrollbarInstanceType
		scrollbarInstance.layerAlias = layerAlias
		scrollbarInstance.controlAlias = viewportEntry.HorizontalScrollbarAlias
		scrollbarInstance.Delete()
	}

	if viewportEntry.VerticalScrollbarAlias != "" {
		var scrollbarInstance ScrollbarInstanceType
		scrollbarInstance.layerAlias = layerAlias
		scrollbarInstance.controlAlias = viewportEntry.VerticalScrollbarAlias
		scrollbarInstance.Delete()
	}

	// Delete the viewport
	Viewports.Remove(layerAlias, viewportAlias)
}

/*
DeleteAll is a method which allows you to remove all viewports from the specified layer.

Example:
    viewport.DeleteAll("main")
*/
func (shared *viewportType) DeleteAll(layerAlias string) {
	// Get all entries for the layer
	entries := Viewports.GetAllEntries(layerAlias)
	// Remove each entry
	for _, entry := range entries {
		Viewports.Remove(layerAlias, entry.Alias)
	}
}

/*
drawOnLayer is a method which allows you to draw all viewports on the specified layer.

Example:
    viewport.drawOnLayer(layer)
*/
func (shared *viewportType) drawOnLayer(layerEntry types.LayerEntryType) {
	// Get all entries for the layer
	entries := Viewports.GetAllEntries(layerEntry.LayerAlias)
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
	viewportEntry := GetViewport(layerEntry.LayerAlias, viewportAlias)
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
	// Create a temporary layer instance to call DrawBorder on
	var layerInstance LayerInstanceType
	layerInstance.layerAlias = layerEntry.LayerAlias
	layerInstance.DrawBorder(styleEntry, xLocation, yLocation, width, height, false)
}

/*
drawContent is a method which allows you to draw the content of the viewport.

Example:
    viewport.drawContent(layer, "myVP", style, attr, 1, 1, 38, 8)
*/
func (shared *viewportType) drawContent(layerEntry *types.LayerEntryType, viewportAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, height int) {
	viewportEntry := GetViewport(layerEntry.LayerAlias, viewportAlias)
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
	if viewportEntry.VerticalScrollbarAlias != "" && ScrollBars.IsExists(layerEntry.LayerAlias, viewportEntry.VerticalScrollbarAlias) {
		scrollbarEntry := ScrollBars.Get(layerEntry.LayerAlias, viewportEntry.VerticalScrollbarAlias)
		if scrollbarEntry.IsVisible {
			contentWidth-- // Reduce width by 1 for vertical scrollbar
		}
	}

	if viewportEntry.HorizontalScrollbarAlias != "" && ScrollBars.IsExists(layerEntry.LayerAlias, viewportEntry.HorizontalScrollbarAlias) {
		scrollbarEntry := ScrollBars.Get(layerEntry.LayerAlias, viewportEntry.HorizontalScrollbarAlias)
		if scrollbarEntry.IsVisible {
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
	for y := 0; y < contentHeight; y++ {
		textDataY := y + viewportEntry.ViewportYLocation
		if textDataY < len(viewportEntry.TextData) {
			line := viewportEntry.TextData[textDataY]

			// Apply horizontal scrolling
			startX := viewportEntry.ViewportXLocation
			endX := startX + contentWidth
			if endX > len(line) {
				endX = len(line)
			}

			if startX < len(line) {
				// Convert the line to a string for markup processing
				lineString := string(line)

				// Process the line with markup
				cursorXLocation := contentXLocation
				currentAttributeEntry := types.NewAttributeEntry(&styleAttributeEntry)

				for currentCharacterIndex := startX; currentCharacterIndex < endX; currentCharacterIndex++ {
					// Check for markup tags
					if currentCharacterIndex+1 < len(lineString) {
						nextCharacter := stringformat.GetSubString(lineString, currentCharacterIndex, 2)
						if nextCharacter == "{{" {
							attributeTag := getAttributeTag(lineString, currentCharacterIndex)
							if attributeTag != "" {
								// Apply the style from the markup tag
								currentAttributeEntry = getDialogAttributeEntry(attributeTag, styleAttributeEntry)
								// Skip the markup tag
								currentCharacterIndex += len(attributeTag) - 1
								continue
							}
						}
					}

					// Print the character with the current attribute entry
					layer.printLayer(layerEntry, currentAttributeEntry, cursorXLocation, contentYLocation+y, []rune{line[currentCharacterIndex]})
					cursorXLocation++
				}
			}
		}
	}
}

/*
updateScrollbars is a method which allows you to update and draw scrollbars for the viewport.

Example:
    viewport.updateScrollbars(layer, "myVP", style, attr, 0, 0, 40, 10)
*/
func (shared *viewportType) updateScrollbars(layerEntry *types.LayerEntryType, viewportAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, height int) {
	viewportEntry := GetViewport(layerEntry.LayerAlias, viewportAlias)
	if viewportEntry == nil {
		return
	}

	// Check if scrollbars exist and are visible
	var hScrollbarEntry, vScrollbarEntry *types.ScrollbarEntryType
	var isHorizontalScrollbarVisible, isVerticalScrollbarVisible bool

	if viewportEntry.HorizontalScrollbarAlias != "" && ScrollBars.IsExists(layerEntry.LayerAlias, viewportEntry.HorizontalScrollbarAlias) {
		hScrollbarEntry = ScrollBars.Get(layerEntry.LayerAlias, viewportEntry.HorizontalScrollbarAlias)
		isHorizontalScrollbarVisible = hScrollbarEntry.IsVisible
	}

	if viewportEntry.VerticalScrollbarAlias != "" && ScrollBars.IsExists(layerEntry.LayerAlias, viewportEntry.VerticalScrollbarAlias) {
		vScrollbarEntry = ScrollBars.Get(layerEntry.LayerAlias, viewportEntry.VerticalScrollbarAlias)
		isVerticalScrollbarVisible = vScrollbarEntry.IsVisible
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
	if viewportEntry.HorizontalScrollbarAlias != "" && ScrollBars.IsExists(layerAlias, viewportEntry.HorizontalScrollbarAlias) {
		scrollbarEntry := ScrollBars.Get(layerAlias, viewportEntry.HorizontalScrollbarAlias)
		viewportEntry.ViewportXLocation = scrollbarEntry.ScrollValue
	}

	if viewportEntry.VerticalScrollbarAlias != "" && ScrollBars.IsExists(layerAlias, viewportEntry.VerticalScrollbarAlias) {
		scrollbarEntry := ScrollBars.Get(layerAlias, viewportEntry.VerticalScrollbarAlias)
		viewportEntry.ViewportYLocation = scrollbarEntry.ScrollValue
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
			viewportEntry := GetViewport(layerAlias, viewportAlias)

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
					maxScroll := 0
					if len(viewportEntry.TextData) > viewportEntry.Height {
						maxScroll = len(viewportEntry.TextData) - viewportEntry.Height
						if viewportEntry.IsBorderDrawn {
							maxScroll += 2
						}
					}

					if viewportEntry.ViewportYLocation < maxScroll {
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
					maxWidth := shared.getMaxHorizontalTextValue(layerAlias, viewportAlias)
					if viewportEntry.ViewportXLocation < maxWidth-viewportEntry.Width {
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
		viewports := Viewports.GetAllEntries(layerAlias)

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
