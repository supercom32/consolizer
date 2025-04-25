package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"math"

	"github.com/supercom32/consolizer/types"
)

// ViewportInstanceType represents a viewport instance that can be manipulated.
type ViewportInstanceType struct {
	BaseControlInstanceType
}

// viewportType is the internal type that manages all viewports.
type viewportType struct{}

// viewport is the global instance for managing viewport controls.
var viewport viewportType

// Viewports is the global instance for managing all viewport controls.
var Viewports = memory.NewControlMemoryManager[types.ViewportEntryType]()

// AddViewport creates a new viewport control and adds it to the specified layer.
func AddViewport(layerAlias string, viewportAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isLinesWrapped bool, isBorderDrawn bool, maxHistoryLines int) ViewportInstanceType {
	return viewport.AddViewport(layerAlias, viewportAlias, styleEntry, xLocation, yLocation, width, height, isLinesWrapped, isBorderDrawn, maxHistoryLines)
}

// GetViewport retrieves a viewport entry by its layer and viewport alias.
func GetViewport(layerAlias string, viewportAlias string) *types.ViewportEntryType {
	return Viewports.Get(layerAlias, viewportAlias)
}

// IsViewportExists checks if a viewport exists in the specified layer.
func IsViewportExists(layerAlias string, viewportAlias string) bool {
	return Viewports.IsExists(layerAlias, viewportAlias)
}

// DeleteViewport removes a viewport from the specified layer.
func DeleteViewport(layerAlias string, viewportAlias string) {
	viewport.DeleteViewport(layerAlias, viewportAlias)
}

// DeleteAllViewportsFromLayer removes all viewports from the specified layer.
func DeleteAllViewportsFromLayer(layerAlias string) {
	// Get all entries for the layer
	entries := Viewports.GetAllEntries(layerAlias)
	// Remove each entry
	for _, entry := range entries {
		Viewports.Remove(layerAlias, entry.Alias)
	}
}

// Print adds text to the viewport.
func (shared *ViewportInstanceType) Print(text string) *ViewportInstanceType {
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

// SetViewport sets the viewport position.
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

// SetMaxHistoryLines sets the maximum number of lines to keep in history.
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

// EnableHistory enables or disables history tracking.
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

// Clear clears all text from the viewport.
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

// SetTransparent sets whether the viewport background is transparent.
func (shared *ViewportInstanceType) SetTransparent(isTransparent bool) *ViewportInstanceType {
	viewportEntry := GetViewport(shared.layerAlias, shared.controlAlias)
	if viewportEntry == nil {
		return shared
	}

	viewportEntry.IsTransparent = isTransparent
	return shared
}

// printText processes and adds text to the viewport.
func (shared *viewportType) printText(viewportEntry *types.ViewportEntryType, textToPrint string) {
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

// processTextWithMarkup processes text with markup and word wrapping.
func (shared *viewportType) processTextWithMarkup(textToPrint string, widthOfLineInCharacters int, isLinesWrapped bool) [][]rune {
	var result [][]rune
	var currentLine []rune

	arrayOfRunes := stringformat.GetRunesFromString(textToPrint)
	cursorXLocation := 0

	for currentCharacterIndex := 0; currentCharacterIndex < len(arrayOfRunes); currentCharacterIndex++ {
		currentCharacter := stringformat.GetSubString(textToPrint, currentCharacterIndex, 1)

		// Add the character to the current line
		currentLine = append(currentLine, arrayOfRunes[currentCharacterIndex])
		cursorXLocation++

		// Check for word wrapping if lines should be wrapped
		lengthOfNextWord := 0
		if isLinesWrapped && currentCharacter == " " {
			lengthOfNextWord = getLengthOfNextWord(textToPrint, currentCharacterIndex+1)
		}

		// Handle newlines
		if currentCharacter == "\n" {
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

// trimHistory trims the history to the maximum number of lines.
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

// trimToVisible trims the text data to only what's visible.
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

// updateScrollbarBasedOnViewportViewport updates scrollbar positions based on viewport position.
func (shared *viewportType) updateScrollbarBasedOnViewportViewport(layerAlias string, viewportAlias string) {
	viewportEntry := GetViewport(layerAlias, viewportAlias)
	if viewportEntry == nil {
		return
	}

	// Update vertical scrollbar if it exists
	if viewportEntry.VerticalScrollbarAlias != "" && ScrollBars.IsExists(layerAlias, viewportEntry.VerticalScrollbarAlias) {
		scrollbarEntry := ScrollBars.Get(layerAlias, viewportEntry.VerticalScrollbarAlias)
		scrollbarEntry.ScrollValue = viewportEntry.ViewportYLocation
		scrollbar.computeScrollbarHandlePositionByScrollValue(layerAlias, viewportEntry.VerticalScrollbarAlias)
	}

	// Update horizontal scrollbar if it exists
	if viewportEntry.HorizontalScrollbarAlias != "" && ScrollBars.IsExists(layerAlias, viewportEntry.HorizontalScrollbarAlias) {
		scrollbarEntry := ScrollBars.Get(layerAlias, viewportEntry.HorizontalScrollbarAlias)
		scrollbarEntry.ScrollValue = viewportEntry.ViewportXLocation
		scrollbar.computeScrollbarHandlePositionByScrollValue(layerAlias, viewportEntry.HorizontalScrollbarAlias)
	}
}

// getMaxHorizontalTextValue gets the maximum width of text in the viewport.
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

// setViewportMaxScrollBarValues sets the maximum scroll values for scrollbars.
func (shared *viewportType) setViewportMaxScrollBarValues(layerAlias string, viewportAlias string) {
	viewportEntry := GetViewport(layerAlias, viewportAlias)
	if viewportEntry == nil {
		return
	}

	// Calculate effective width and height (accounting for border if present)
	effectiveWidth := viewportEntry.Width
	effectiveHeight := viewportEntry.Height
	if viewportEntry.IsBorderDrawn {
		effectiveWidth -= 2
		effectiveHeight -= 2
	}

	// Set vertical scrollbar max value
	if viewportEntry.VerticalScrollbarAlias != "" && ScrollBars.IsExists(layerAlias, viewportEntry.VerticalScrollbarAlias) {
		scrollbarEntry := ScrollBars.Get(layerAlias, viewportEntry.VerticalScrollbarAlias)
		maxScrollValue := 0
		if len(viewportEntry.TextData) > effectiveHeight {
			maxScrollValue = len(viewportEntry.TextData) - effectiveHeight
			scrollbarEntry.IsEnabled = true
			scrollbarEntry.IsVisible = true
		} else {
			scrollbarEntry.IsEnabled = false
			scrollbarEntry.IsVisible = false
		}
		scrollbarEntry.MaxScrollValue = maxScrollValue
	}

	// Set horizontal scrollbar max value
	if viewportEntry.HorizontalScrollbarAlias != "" && ScrollBars.IsExists(layerAlias, viewportEntry.HorizontalScrollbarAlias) {
		scrollbarEntry := ScrollBars.Get(layerAlias, viewportEntry.HorizontalScrollbarAlias)
		maxWidth := shared.getMaxHorizontalTextValue(layerAlias, viewportAlias)
		maxScrollValue := 0

		// If lines are wrapped, always disable horizontal scrollbar
		if viewportEntry.IsLinesWrapped {
			scrollbarEntry.IsEnabled = false
			scrollbarEntry.IsVisible = false
		} else if maxWidth > effectiveWidth {
			// Only show horizontal scrollbar if lines are not wrapped and content is wider than viewport
			maxScrollValue = maxWidth - effectiveWidth
			scrollbarEntry.IsEnabled = true
			scrollbarEntry.IsVisible = true
		} else {
			scrollbarEntry.IsEnabled = false
			scrollbarEntry.IsVisible = false
		}
		scrollbarEntry.MaxScrollValue = maxScrollValue
	}
}

// AddViewport adds a new viewport to the specified layer.
func (shared *viewportType) AddViewport(layerAlias string, viewportAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isLinesWrapped bool, isBorderDrawn bool, maxHistoryLines int) ViewportInstanceType {
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

	// Calculate scrollbar position and length
	scrollbarXLocation := xLocation
	scrollbarYLocation := yLocation + height
	scrollbarLength := width

	// Add horizontal scrollbar
	scrollbar.Add(layerAlias, horizontalScrollbarAlias, scrollbarStyleEntry, scrollbarXLocation, scrollbarYLocation, scrollbarLength, 0, 0, 1, true)
	viewportEntry.HorizontalScrollbarAlias = horizontalScrollbarAlias

	// For vertical scrollbar
	verticalScrollbarAlias := viewportAlias + "_vscroll"

	// Calculate scrollbar position and length
	scrollbarXLocation = xLocation + width
	scrollbarYLocation = yLocation
	scrollbarLength = height

	// Add vertical scrollbar
	scrollbar.Add(layerAlias, verticalScrollbarAlias, scrollbarStyleEntry, scrollbarXLocation, scrollbarYLocation, scrollbarLength, 0, 0, 1, false)
	viewportEntry.VerticalScrollbarAlias = verticalScrollbarAlias

	// Set up the viewport instance
	viewportInstance.layerAlias = layerAlias
	viewportInstance.controlAlias = viewportAlias
	viewportInstance.controlType = "viewport"

	return viewportInstance
}

// DeleteViewport removes a viewport from the specified layer.
func (shared *viewportType) DeleteViewport(layerAlias string, viewportAlias string) {
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

// DeleteAllViewports removes all viewports from the specified layer.
func (shared *viewportType) DeleteAllViewports(layerAlias string) {
	// Get all entries for the layer
	entries := Viewports.GetAllEntries(layerAlias)
	// Remove each entry
	for _, entry := range entries {
		Viewports.Remove(layerAlias, entry.Alias)
	}
}

// drawViewportsOnLayer draws all viewports on the specified layer.
func (shared *viewportType) drawViewportsOnLayer(layerEntry types.LayerEntryType) {
	// Get all entries for the layer
	entries := Viewports.GetAllEntries(layerEntry.LayerAlias)
	// Draw each viewport
	for _, entry := range entries {
		shared.drawViewport(&layerEntry, entry.Alias)
	}
}

// drawViewport draws a single viewport on the specified layer.
func (shared *viewportType) drawViewport(layerEntry *types.LayerEntryType, viewportAlias string) {
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
	shared.drawViewportContent(layerEntry, viewportAlias, styleEntry, attributeEntry, viewportEntry.XLocation, viewportEntry.YLocation, viewportEntry.Width, viewportEntry.Height)

	// Draw scrollbars if needed
	shared.drawViewportScrollbars(layerEntry, viewportAlias, styleEntry, attributeEntry, viewportEntry.XLocation, viewportEntry.YLocation, viewportEntry.Width, viewportEntry.Height)
}

// drawBorder draws a border around the viewport.
func (shared *viewportType) drawBorder(layerEntry *types.LayerEntryType, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, height int, isDoubleLine bool) {
	// Create a temporary layer instance to call DrawBorder on
	var layerInstance LayerInstanceType
	layerInstance.layerAlias = layerEntry.LayerAlias
	layerInstance.DrawBorder(styleEntry, xLocation, yLocation, width, height, false)
}

// drawViewportContent draws the content of the viewport.
func (shared *viewportType) drawViewportContent(layerEntry *types.LayerEntryType, viewportAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, height int) {
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
				printLayer(layerEntry, localAttributeEntry, contentXLocation+x, contentYLocation+y, spaceChar)
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
					printLayer(layerEntry, currentAttributeEntry, cursorXLocation, contentYLocation+y, []rune{line[currentCharacterIndex]})
					cursorXLocation++
				}
			}
		}
	}
}

// drawViewportScrollbars draws scrollbars for the viewport.
func (shared *viewportType) drawViewportScrollbars(layerEntry *types.LayerEntryType, viewportAlias string, styleEntry types.TuiStyleEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, height int) {
	viewportEntry := GetViewport(layerEntry.LayerAlias, viewportAlias)
	if viewportEntry == nil {
		return
	}

	// Draw horizontal scrollbar if it exists
	if viewportEntry.HorizontalScrollbarAlias != "" && ScrollBars.IsExists(layerEntry.LayerAlias, viewportEntry.HorizontalScrollbarAlias) {
		// The scrollbar will be drawn by the scrollbar drawing code
	}

	// Draw vertical scrollbar if it exists
	if viewportEntry.VerticalScrollbarAlias != "" && ScrollBars.IsExists(layerEntry.LayerAlias, viewportEntry.VerticalScrollbarAlias) {
		// The scrollbar will be drawn by the scrollbar drawing code
	}
}

// Note: We don't need to draw scrollbars manually as they are drawn by the scrollbar drawing code

// updateViewport updates the viewport position based on scrollbar values.
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

// updateMouseEvent handles mouse events for viewports.
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
					if viewportEntry.ViewportXLocation < maxWidth - viewportEntry.Width {
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

// checkScrollbarEvents checks for scrollbar events and updates all viewports accordingly.
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

// Note: To fully integrate this viewport control into the system,
// you would need to add the drawViewportsOnLayer function to the renderControls
// function in terminal.go. However, since we can't modify existing files,
// this control will need to be manually added to the rendering pipeline.
