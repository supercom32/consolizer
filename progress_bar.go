package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
	"github.com/u2takey/go-utils/strings"
)

type ProgressBarInstanceType struct {
	BaseControlInstanceType
}

type progressBarType struct{}

var ProgressBar progressBarType
var ProgressBars = memory.NewControlMemoryManager[types.ProgressBarEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

func (shared *ProgressBarInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeProgressBar)
}

func (shared *ProgressBarInstanceType) Delete() *ProgressBarInstanceType {
	if ProgressBars.IsExists(shared.layerAlias, shared.controlAlias) {
		ProgressBars.Remove(shared.layerAlias, shared.controlAlias)
	}
	return nil
}

func (shared *ProgressBarInstanceType) SetProgressBarValue(value int) {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if value <= progressBarEntry.MaxValue {
		progressBarEntry.Value = value
	} else {
		progressBarEntry.Value = progressBarEntry.MaxValue
	}
}

func (shared *ProgressBarInstanceType) SetProgressBarMaxValue(value int) {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if value > 0 {
		progressBarEntry.MaxValue = value
	}
}

func (shared *ProgressBarInstanceType) SetProgressBarLabel(label string) {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	progressBarEntry.Label = label
}

func (shared *ProgressBarInstanceType) IncrementProgressBarValue() {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	progressBarEntry.Value = progressBarEntry.Value + 1
	if progressBarEntry.Value > progressBarEntry.MaxValue {
		progressBarEntry.Value = progressBarEntry.MaxValue
	}
}

func (shared *ProgressBarInstanceType) GetProgressBarValueAsRatio() string {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	valueAsString := fmt.Sprintf("%d/%d", progressBarEntry.Value, progressBarEntry.MaxValue)
	return valueAsString
}

func (shared *ProgressBarInstanceType) GetProgressBarValue() int {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	return progressBarEntry.Value
}

func (shared *ProgressBarInstanceType) GetProgressBarValueAsPercent() string {
	var returnValue int
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if progressBarEntry.MaxValue > 0 {
		returnValue = (progressBarEntry.Value * 100) / progressBarEntry.MaxValue
	}
	return fmt.Sprintf("%d", returnValue)
}

func (shared *progressBarType) Add(layerAlias string, progressBarAlias string, progressBarLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isVertical bool, value int, maxValue int, isBackgroundTransparent bool) ProgressBarInstanceType {
	progressBarEntry := types.NewProgressBarEntry()
	progressBarEntry.StyleEntry = styleEntry
	progressBarEntry.Alias = progressBarAlias
	progressBarEntry.Label = progressBarLabel
	progressBarEntry.Value = value
	progressBarEntry.MaxValue = maxValue
	progressBarEntry.IsBackgroundTransparent = isBackgroundTransparent
	progressBarEntry.XLocation = xLocation
	progressBarEntry.YLocation = yLocation
	progressBarEntry.Width = width
	progressBarEntry.Height = height
	progressBarEntry.IsVertical = isVertical
	progressBarEntry.TooltipAlias = stringformat.GetLastSortedUUID()

	// Create associated tooltip (always created but disabled by default)
	tooltipInstance := Tooltip.Add(layerAlias, progressBarEntry.TooltipAlias, "", styleEntry,
		progressBarEntry.XLocation, progressBarEntry.YLocation,
		progressBarEntry.Width, progressBarEntry.Height,
		progressBarEntry.XLocation, progressBarEntry.YLocation+progressBarEntry.Height+1,
		progressBarEntry.Width, 3,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.SetEnabled(false)
	tooltipInstance.setParentControlAlias(progressBarAlias)
	// Use the ControlMemoryManager to add the progress bar entry
	ProgressBars.Add(layerAlias, progressBarAlias, &progressBarEntry)
	var progressBarInstance ProgressBarInstanceType
	progressBarInstance.layerAlias = layerAlias
	progressBarInstance.controlAlias = progressBarAlias
	progressBarInstance.controlType = "progressbar"
	return progressBarInstance
}

/*
DeleteButton allows you to remove a button from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a button which does not exist, then the request
will simply be ignored.
*/
func (shared *progressBarType) DeleteProgressBar(layerAlias string, progressBarAlias string) {
	Buttons.Remove(layerAlias, progressBarAlias)
}

func (shared *progressBarType) DeleteAllProgressBars(layerAlias string) {
	ProgressBars.RemoveAll(layerAlias)
}

/*
drawButtonsOnLayer allows you to draw all buttons on a given text layer.
*/
func (shared *progressBarType) drawProgressBarsOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentProgressBarEntry := range ProgressBars.GetAllEntries(layerAlias) {
		progressBarEntry := currentProgressBarEntry
		drawProgressBar(&layerEntry, progressBarEntry.Alias, progressBarEntry.Label, progressBarEntry.StyleEntry, progressBarEntry.XLocation, progressBarEntry.YLocation, progressBarEntry.Width, progressBarEntry.Height, progressBarEntry.Value, progressBarEntry.MaxValue, progressBarEntry.IsBackgroundTransparent, progressBarEntry.IsVertical)
	}
}

// getHorizontalPartialFillChar returns the appropriate block character for a partial horizontal fill
func getHorizontalPartialFillChar(partialFill float64) rune {
	// Map the partial fill (0.0 to 1.0) to one of the block characters
	// For horizontal bars, we use the left-partial block characters
	if partialFill < 0.125 {
		return constants.CharBlockLeftOneEighth
	} else if partialFill < 0.25 {
		return constants.CharBlockLeftOneQuarter
	} else if partialFill < 0.375 {
		return constants.CharBlockLeftThreeEighths
	} else if partialFill < 0.5 {
		return constants.CharBlockLeftHalf
	} else if partialFill < 0.625 {
		return constants.CharBlockLeftFiveEighths
	} else if partialFill < 0.75 {
		return constants.CharBlockLeftThreeQuarters
	} else if partialFill < 0.875 {
		return constants.CharBlockLeftSevenEighths
	} else {
		return constants.CharBlockFull
	}
}

// getVerticalPartialFillChar returns the appropriate block character for a partial vertical fill
func getVerticalPartialFillChar(partialFill float64) rune {
	// Map the partial fill (0.0 to 1.0) to one of the block characters
	// For vertical bars, we use the lower-partial block characters
	if partialFill < 0.125 {
		return constants.CharBlockLowerOneEighth
	} else if partialFill < 0.25 {
		return constants.CharBlockLowerOneQuarter
	} else if partialFill < 0.375 {
		return constants.CharBlockLowerThreeEighths
	} else if partialFill < 0.5 {
		return constants.CharBlockLowerHalf
	} else if partialFill < 0.625 {
		return constants.CharBlockLowerFiveEighths
	} else if partialFill < 0.75 {
		return constants.CharBlockLowerThreeQuarters
	} else if partialFill < 0.875 {
		return constants.CharBlockLowerSevenEighths
	} else {
		return constants.CharBlockFull
	}
}

func drawProgressBar(layerEntry *types.LayerEntryType, progressBarAlias string, progressBarLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, currentValue int, maxValue int, isBackgroundTransparent bool, isVertical bool) {
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.ProgressBar.UnfilledForegroundColor
	attributeEntry.BackgroundColor = styleEntry.ProgressBar.UnfilledBackgroundColor
	attributeEntry.CellType = constants.CellTypeProgressBar
	attributeEntry.CellControlAlias = progressBarAlias
	if height < 1 {
		height = 1
	}

	// Fill the entire area with the unfilled pattern
	fillArea(layerEntry, attributeEntry, string(styleEntry.ProgressBar.UnfilledPattern), xLocation, yLocation, width, height, constants.CellTypeProgressBar)

	// Calculate and draw the filled portion based on orientation
	attributeEntry.ForegroundColor = styleEntry.ProgressBar.FilledForegroundColor
	attributeEntry.BackgroundColor = styleEntry.ProgressBar.FilledBackgroundColor

	// Create a separate attribute entry for partial fill characters
	partialFillAttributeEntry := types.NewAttributeEntry(&attributeEntry)
	partialFillAttributeEntry.ForegroundColor = styleEntry.ProgressBar.FilledForegroundColor
	// Use the unfilled background color to avoid visible color differences for partial fill characters
	partialFillAttributeEntry.BackgroundColor = styleEntry.ProgressBar.FilledBackgroundColor

	if !isVertical {
		// Horizontal progress bar (left to right)
		progressBarWidthExact := float64(width) * float64(currentValue) / float64(maxValue)
		fullCellsCount := int(progressBarWidthExact)

		// Fill the complete cells
		if fullCellsCount > 0 {
			fillArea(layerEntry, attributeEntry, string(styleEntry.ProgressBar.FilledPattern), xLocation, yLocation, fullCellsCount, height, constants.CellTypeProgressBar)
		}

		// Calculate the partial fill for the last cell
		partialFill := progressBarWidthExact - float64(fullCellsCount)
		if partialFill > 0 && fullCellsCount < width {
			if styleEntry.ProgressBar.IsHighResolution {
				// Map the partial fill to one of the block characters
				partialChar := getHorizontalPartialFillChar(partialFill)

				// Draw the partial fill character
				for h := 0; h < height; h++ {
					printLayer(layerEntry, partialFillAttributeEntry, xLocation+fullCellsCount, yLocation+h, []rune{partialChar})
				}
			} else {
				// For low resolution, only fill whole blocks
				// If partialFill is significant enough (e.g., > 0.5), consider it a full block
				if partialFill > 0.5 && fullCellsCount < width {
					for h := 0; h < height; h++ {
						printLayer(layerEntry, attributeEntry, xLocation+fullCellsCount, yLocation+h, []rune{styleEntry.ProgressBar.FilledPattern})
					}
				}
			}
		}
	} else {
		// Vertical progress bar (bottom to top)
		progressBarHeightExact := float64(height) * float64(currentValue) / float64(maxValue)
		fullCellsCount := int(progressBarHeightExact)

		// Calculate the starting position for the filled area
		fillStartY := yLocation + height - fullCellsCount

		// Fill the complete cells
		if fullCellsCount > 0 {
			fillArea(layerEntry, attributeEntry, string(styleEntry.ProgressBar.FilledPattern), xLocation, fillStartY, width, fullCellsCount, constants.CellTypeProgressBar)
		}

		// Calculate the partial fill for the last cell
		partialFill := progressBarHeightExact - float64(fullCellsCount)
		if partialFill > 0 && fullCellsCount < height {
			if styleEntry.ProgressBar.IsHighResolution {
				// Map the partial fill to one of the block characters
				partialChar := getVerticalPartialFillChar(partialFill)

				// Draw the partial fill character
				partialCellY := fillStartY - 1
				if partialCellY >= yLocation {
					for w := 0; w < width; w++ {
						printLayer(layerEntry, partialFillAttributeEntry, xLocation+w, partialCellY, []rune{partialChar})
					}
				}
			} else {
				// For low resolution, only fill whole blocks
				// If partialFill is significant enough (e.g., > 0.5), consider it a full block
				if partialFill > 0.5 && fullCellsCount < height {
					partialCellY := fillStartY - 1
					if partialCellY >= yLocation {
						for w := 0; w < width; w++ {
							printLayer(layerEntry, attributeEntry, xLocation+w, partialCellY, []rune{styleEntry.ProgressBar.FilledPattern})
						}
					}
				}
			}
		}
	}

	// Handle label display
	arrayOfRunes := stringformat.GetRunesFromString(progressBarLabel)
	attributeEntry.ForegroundColor = styleEntry.ProgressBar.TextForegroundColor
	attributeEntry.BackgroundColor = styleEntry.ProgressBar.TextBackgroundColor
	attributeEntry.IsBackgroundTransparent = isBackgroundTransparent

	if !isVertical {
		// For horizontal progress bars, check if label is too long for width
		if len(progressBarLabel) > width {
			if width <= 3 {
				// If width is too small, just show what we can or nothing
				if width > 0 {
					progressBarLabel = progressBarLabel[:width]
				} else {
					progressBarLabel = ""
				}
			} else {
				progressBarLabel = strings.ShortenString(progressBarLabel, width-3)
				progressBarLabel = progressBarLabel + "..."
			}
			arrayOfRunes = stringformat.GetRunesFromString(progressBarLabel)
		}
		// Center the label for horizontal progress bars
		centerXLocation := (width - len(progressBarLabel)) / 2
		centerYLocation := height / 2
		printLayer(layerEntry, attributeEntry, xLocation+centerXLocation, yLocation+centerYLocation, arrayOfRunes)
	} else {
		// For vertical progress bars, check if label is too long for height
		if len(progressBarLabel) > height {
			if height <= 3 {
				// If height is too small, just show what we can or nothing
				if height > 0 {
					progressBarLabel = progressBarLabel[:height]
				} else {
					progressBarLabel = ""
				}
			} else {
				progressBarLabel = strings.ShortenString(progressBarLabel, height-3)
				progressBarLabel = progressBarLabel + "..."
			}
			arrayOfRunes = stringformat.GetRunesFromString(progressBarLabel)
		}
		// For vertical progress bars, print each character vertically and center it
		centerXLocation := width / 2
		// Calculate vertical starting position to center the label
		centerYLocation := (height - len(progressBarLabel)) / 2

		// Print each character of the label vertically
		for i, char := range arrayOfRunes {
			printLayer(layerEntry, attributeEntry, xLocation+centerXLocation, yLocation+centerYLocation+i, []rune{char})
		}
	}
}
