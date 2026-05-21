package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
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

/*
AddToTabIndex is a method which allows you to add the progress bar to the tab navigation index.
*/
func (shared *ProgressBarInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeProgressBar)
}

/*
Delete is a method which allows you to remove the progress bar instance.

:return: A nil pointer of ProgressBarInstanceType.
*/
func (shared *ProgressBarInstanceType) Delete() *ProgressBarInstanceType {
	shared.BaseControlInstanceType.Delete()
	return nil
}

/*
SetValue is a method which allows you to set the current value of the progress bar. In addition, the following
information should be noted:

- If the value provided is greater than the maximum value, the current value will be set to the maximum value.

:param value: The new value to set for the progress bar.
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
SetMaxValue is a method which allows you to set the maximum value of the progress bar. In addition, the following
information should be noted:

- The maximum value must be greater than zero.

:param value: The new maximum value to set for the progress bar.
*/
func (shared *ProgressBarInstanceType) SetMaxValue(value int) {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if value > 0 {
		progressBarEntry.MaxValue = value
	}
}

/*
SetLabel is a method which allows you to set the label text displayed on the progress bar.

:param label: The new label text to set for the progress bar.
*/
func (shared *ProgressBarInstanceType) SetLabel(label string) {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	progressBarEntry.Label = label
}

/*
IncrementValue is a method which allows you to increase the current value of the progress bar by one. In addition, the
following information should be noted:

- The value will not exceed the defined maximum value.
*/
func (shared *ProgressBarInstanceType) IncrementValue() {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	progressBarEntry.Value = progressBarEntry.Value + 1
	if progressBarEntry.Value > progressBarEntry.MaxValue {
		progressBarEntry.Value = progressBarEntry.MaxValue
	}
}

/*
GetValueAsRatio is a method which allows you to retrieve the progress bar's current value and maximum value as a
formatted ratio string.

:return: A string representing the progress as "value/maxValue".
*/
func (shared *ProgressBarInstanceType) GetValueAsRatio() string {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	valueAsString := fmt.Sprintf("%d/%d", progressBarEntry.Value, progressBarEntry.MaxValue)
	return valueAsString
}

/*
GetValue is a method which allows you to retrieve the current value of the progress bar.

:return: The current integer value of the progress bar.
*/
func (shared *ProgressBarInstanceType) GetValue() int {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	return progressBarEntry.Value
}

/*
GetValueAsPercent is a method which allows you to retrieve the progress bar's current value as a percentage of its
maximum value.

:return: A string representing the progress percentage.
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
Add is a method which allows you to add a progress bar to a given text layer. In addition, the following information
should be noted:

- A tooltip is automatically created for the progress bar but is disabled by default.

:param layerAlias: The alias of the layer to which the progress bar will be added.
:param progressBarAlias: The unique alias for the progress bar control.
:param progressBarLabel: The label text to be displayed on the progress bar.
:param styleEntry: The style configuration for the progress bar.
:param xLocation: The x coordinate of the progress bar's position.
:param yLocation: The y coordinate of the progress bar's position.
:param width: The width of the progress bar.
:param height: The height of the progress bar.
:param isVertical: A boolean indicating if the progress bar should be rendered vertically.
:param value: The initial value of the progress bar.
:param maxValue: The maximum value of the progress bar.
:param isBackgroundTransparent: A boolean indicating if the background should be transparent.

:return: An instance of the created progress bar.
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
Delete is a method which allows you to remove a progress bar from a text layer. In addition, the following
information should be noted:

- If you attempt to delete a progress bar which does not exist, then the request will simply be ignored.

:param layerAlias: The alias of the layer from which to remove the progress bar.
:param progressBarAlias: The alias of the progress bar to be removed.
*/
func (shared *progressBarType) Delete(layerAlias string, progressBarAlias string) {
	Buttons.Remove(layerAlias, progressBarAlias)
}

/*
DeleteAll is a method which allows you to remove all progress bars from a specified text layer.

:param layerAlias: The alias of the layer from which all progress bars will be removed.
*/
func (shared *progressBarType) DeleteAll(layerAlias string) {
	ProgressBars.RemoveAll(layerAlias)
}

/*
drawOnLayer is a method which allows you to draw all progress bars on a given text layer.

:param layerEntry: The layer entry on which to draw the progress bars.
*/
func (shared *progressBarType) drawOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentProgressBarEntry := range ProgressBars.GetAllEntries(layerAlias) {
		progressBarEntry := currentProgressBarEntry
		drawProgressBar(&layerEntry, progressBarEntry.Alias, progressBarEntry.Label, progressBarEntry.StyleEntry, progressBarEntry.XLocation, progressBarEntry.YLocation, progressBarEntry.Width, progressBarEntry.Height, progressBarEntry.Value, progressBarEntry.MaxValue, progressBarEntry.IsBackgroundTransparent, progressBarEntry.IsVertical)
	}
}

/*
getHorizontalPartialFillChar is a method which allows you to retrieve the appropriate block character for a partial
horizontal fill.

:param partialFill: The partial fill amount (0.0 to 1.0).

:return: The rune representing the partial fill character.
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
getVerticalPartialFillChar is a method which allows you to retrieve the appropriate block character for a partial
vertical fill.

:param partialFill: The partial fill amount (0.0 to 1.0).

:return: The rune representing the partial fill character.
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
drawProgressBar is a method which allows you to render a progress bar onto a layer. In addition, the following
information should be noted:

- The progress bar can be rendered horizontally or vertically.

- Labels are automatically centered and truncated if they exceed the available dimensions.

:param layerEntry: The layer on which to draw the progress bar.
:param progressBarAlias: The alias of the progress bar.
:param progressBarLabel: The label text to display.
:param styleEntry: The style configuration to use.
:param xLocation: The x coordinate of the position.
:param yLocation: The y coordinate of the position.
:param width: The width of the bar.
:param height: The height of the bar.
:param currentValue: The current value.
:param maxValue: The maximum value.
:param isBackgroundTransparent: Whether the background is transparent.
:param isVertical: Whether the bar is vertical.
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
