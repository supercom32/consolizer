package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

/*
ProgressBarInstanceType is a structure which represents an instance of a progress bar control.

Example:
    var progressBar ProgressBarInstanceType
*/
type ProgressBarInstanceType struct {
	BaseControlInstanceType
}

/*
progressBarType is a structure which provides methods for managing progress bar controls.

Example:
    var progressBar progressBarType
*/
type progressBarType struct{}

var ProgressBar progressBarType
var ProgressBars = memory.NewControlMemoryManager[types.ProgressBarEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
AddToTabIndex is a method which adds the progress bar to the tab navigation index.

Example:
    progressBar.AddToTabIndex()
*/
func (shared *ProgressBarInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeProgressBar)
}

/*
Delete is a method which removes the progress bar instance.

Example:
    progressBar = progressBar.Delete()
*/
func (shared *ProgressBarInstanceType) Delete() *ProgressBarInstanceType {
	shared.BaseControlInstanceType.Delete()
	return nil
}

/*
SetValue is a method which sets the current value of the progress bar. In addition, the following should be noted:

- If the value provided is greater than the maximum value, the current value will be set to the maximum value.

Example:
    progressBar.SetValue(50)
*/
func (shared *ProgressBarInstanceType) SetValue(value int) {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if value <= progressBarEntry.MaxValue {
		progressBarEntry.Value = value
	} else {
		progressBarEntry.Value = progressBarEntry.MaxValue
	}
}

/*
SetMaxValue is a method which sets the maximum value of the progress bar. In addition, the following should be noted:

- The maximum value must be greater than zero.

Example:
    progressBar.SetMaxValue(100)
*/
func (shared *ProgressBarInstanceType) SetMaxValue(value int) {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if value > 0 {
		progressBarEntry.MaxValue = value
	}
}

/*
SetLabel is a method which sets the label text displayed on the progress bar.

Example:
    progressBar.SetLabel("Loading...")
*/
func (shared *ProgressBarInstanceType) SetLabel(label string) {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	progressBarEntry.Label = label
}

/*
IncrementValue is a method which increases the current value of the progress bar by one. In addition, the following should be noted:

- The value will not exceed the defined maximum value.

Example:
    progressBar.IncrementValue()
*/
func (shared *ProgressBarInstanceType) IncrementValue() {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	progressBarEntry.Value = progressBarEntry.Value + 1
	if progressBarEntry.Value > progressBarEntry.MaxValue {
		progressBarEntry.Value = progressBarEntry.MaxValue
	}
}

/*
GetValueAsRatio is a method which retrieves the progress bar's current value and maximum value as a formatted ratio
string.

Example:
    ratio := progressBar.GetValueAsRatio()
*/
func (shared *ProgressBarInstanceType) GetValueAsRatio() string {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	valueAsString := fmt.Sprintf("%d/%d", progressBarEntry.Value, progressBarEntry.MaxValue)
	return valueAsString
}

/*
GetValue is a method which retrieves the current value of the progress bar.

Example:
    currentValue := progressBar.GetValue()
*/
func (shared *ProgressBarInstanceType) GetValue() int {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	return progressBarEntry.Value
}

/*
GetValueAsPercent is a method which retrieves the progress bar's current value as a percentage of its maximum value.

Example:
    percent := progressBar.GetValueAsPercent()
*/
func (shared *ProgressBarInstanceType) GetValueAsPercent() string {
	var returnValue int
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if progressBarEntry.MaxValue > 0 {
		returnValue = (progressBarEntry.Value * 100) / progressBarEntry.MaxValue
	}
	return fmt.Sprintf("%d", returnValue)
}

/*
Add is a method which adds a progress bar to a given text layer. In addition, the following should be noted:

- A tooltip is automatically created for the progress bar but is disabled by default.

Example:
    progressBarInstance := ProgressBar.Add("Layer1", "Progress1", "Loading", style, 0, 0, 20, 1, false, 0, 100, false)
*/
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
	progressBarInstance.controlType = constants.TYPE_PROGRESSBAR
	return progressBarInstance
}

/*
Delete is a method which removes a progress bar from a text layer. In addition, the following should be noted:

- If you attempt to delete a progress bar which does not exist, then the request will simply be ignored.

Example:
    ProgressBar.Delete("Layer1", "Progress1")
*/
func (shared *progressBarType) Delete(layerAlias string, progressBarAlias string) {
	Buttons.Remove(layerAlias, progressBarAlias)
}

/*
DeleteAll is a method which removes all progress bars from a specified text layer.

Example:
    ProgressBar.DeleteAll("Layer1")
*/
func (shared *progressBarType) DeleteAll(layerAlias string) {
	ProgressBars.RemoveAll(layerAlias)
}

/*
drawOnLayer is a method which draws all progress bars on a given text layer.

Example:
    ProgressBar.drawOnLayer(layerEntry)
*/
func (shared *progressBarType) drawOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentProgressBarEntry := range ProgressBars.GetAllEntries(layerAlias) {
		progressBarEntry := currentProgressBarEntry
		drawProgressBar(&layerEntry, progressBarEntry.Alias, progressBarEntry.Label, progressBarEntry.StyleEntry, progressBarEntry.XLocation, progressBarEntry.YLocation, progressBarEntry.Width, progressBarEntry.Height, progressBarEntry.Value, progressBarEntry.MaxValue, progressBarEntry.IsBackgroundTransparent, progressBarEntry.IsVertical)
	}
}

/*
getHorizontalPartialFillChar is a method which retrieves the appropriate block character for a partial horizontal fill.

Example:
    char := getHorizontalPartialFillChar(0.5)
*/
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

/*
getVerticalPartialFillChar is a method which retrieves the appropriate block character for a partial vertical fill.

Example:
    char := getVerticalPartialFillChar(0.5)
*/
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

/*
drawProgressBar is a method which renders a progress bar onto a layer. In addition, the following should be noted:

- The progress bar can be rendered horizontally or vertically.

- Labels are automatically centered and truncated if they exceed the available dimensions.

Example:
    drawProgressBar(&layerEntry, "Progress1", "Loading", style, 0, 0, 20, 1, 50, 100, false, false)
*/
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
	unfilledPattern := styleEntry.ProgressBar.UnfilledPattern
	// If in high-res mode do not use a pattern and simply use a block character since it is the only valid one.
	if styleEntry.ProgressBar.IsHighResolution {
		unfilledPattern = constants.CharBlockSolid
	}
	fillArea(layerEntry, attributeEntry, string(unfilledPattern), xLocation, yLocation, width, height, constants.CellTypeProgressBar)

	// Calculate and draw the filled portion based on orientation
	attributeEntry.ForegroundColor = styleEntry.ProgressBar.FilledForegroundColor
	attributeEntry.BackgroundColor = styleEntry.ProgressBar.FilledBackgroundColor

	// Create a separate attribute entry for partial fill characters
	partialFillAttributeEntry := types.NewAttributeEntry(&attributeEntry)
	partialFillAttributeEntry.ForegroundColor = styleEntry.ProgressBar.FilledForegroundColor
	// Use the unfilled background color to avoid visible color differences for partial fill characters
	partialFillAttributeEntry.BackgroundColor = styleEntry.ProgressBar.UnfilledBackgroundColor

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
					layer.printLayer(layerEntry, partialFillAttributeEntry, xLocation+fullCellsCount, yLocation+h, []rune{partialChar})
				}
			} else {
				// For low resolution, only fill whole blocks
				// If partialFill is significant enough (e.g., > 0.5), consider it a full block
				if partialFill > 0.5 && fullCellsCount < width {
					for h := 0; h < height; h++ {
						layer.printLayer(layerEntry, attributeEntry, xLocation+fullCellsCount, yLocation+h, []rune{styleEntry.ProgressBar.FilledPattern})
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
						layer.printLayer(layerEntry, partialFillAttributeEntry, xLocation+w, partialCellY, []rune{partialChar})
					}
				}
			} else {
				// For low resolution, only fill whole blocks
				// If partialFill is significant enough (e.g., > 0.5), consider it a full block
				if partialFill > 0.5 && fullCellsCount < height {
					partialCellY := fillStartY - 1
					if partialCellY >= yLocation {
						for w := 0; w < width; w++ {
							layer.printLayer(layerEntry, attributeEntry, xLocation+w, partialCellY, []rune{styleEntry.ProgressBar.FilledPattern})
						}
					}
				}
			}
		}
	}

	// Handle label display
	arrayOfRunes := stringformat.GetRunesFromString(progressBarLabel)
	labelWidth := stringformat.GetWidthOfRunesWhenPrinted(arrayOfRunes)
	attributeEntry.ForegroundColor = styleEntry.ProgressBar.TextForegroundColor
	attributeEntry.BackgroundColor = styleEntry.ProgressBar.TextBackgroundColor
	attributeEntry.IsBackgroundTransparent = isBackgroundTransparent

	if !isVertical {
		// For horizontal progress bars, check if label is too long for width
		if labelWidth > width {
			if width <= 3 {
				arrayOfRunes = stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunes, width)
			} else {
				arrayOfRunes = stringformat.GetMaxCharactersThatFitInStringSize(arrayOfRunes, width-3)
				arrayOfRunes = append(arrayOfRunes, '.', '.', '.')
			}
			labelWidth = stringformat.GetWidthOfRunesWhenPrinted(arrayOfRunes)
		}
		// Center the label for horizontal progress bars
		centerXLocation := (width - labelWidth) / 2
		centerYLocation := height / 2
		layer.printLayer(layerEntry, attributeEntry, xLocation+centerXLocation, yLocation+centerYLocation, arrayOfRunes)
	} else {
		// For vertical progress bars, check if label is too long for height
		numberOfRunes := len(arrayOfRunes)
		if numberOfRunes > height {
			arrayOfRunes = arrayOfRunes[:height]
			numberOfRunes = len(arrayOfRunes)
		}
		// For vertical progress bars, print each character vertically and center it
		centerXLocation := width / 2
		// Calculate vertical starting position to center the label
		centerYLocation := (height - numberOfRunes) / 2

		// Print each character of the label vertically
		for i, char := range arrayOfRunes {
			layer.printLayer(layerEntry, attributeEntry, xLocation+centerXLocation, yLocation+centerYLocation+i, []rune{char})
		}
	}
}
