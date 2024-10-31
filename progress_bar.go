package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
	"github.com/supercom32/consolizer/types"
	"github.com/u2takey/go-utils/strings"
)

type ProgressBarInstanceType struct {
	layerAlias       string
	progressBarAlias string
}

type progressBarType struct{}

var ProgressBar progressBarType

func (shared *ProgressBarInstanceType) Delete() *ProgressBarInstanceType {
	if memory.IsProgressBarExists(shared.layerAlias, shared.progressBarAlias) {
		memory.DeleteProgressBar(shared.layerAlias, shared.progressBarAlias)
	}
	return nil
}

func (shared *ProgressBarInstanceType) SetProgressBarValue(value int) {
	progressBarEntry := memory.GetProgressBar(shared.layerAlias, shared.progressBarAlias)
	if value <= progressBarEntry.MaxValue {
		progressBarEntry.Value = value
	} else {
		progressBarEntry.Value = progressBarEntry.MaxValue
	}
}

func (shared *ProgressBarInstanceType) SetProgressBarMaxValue(value int) {
	progressBarEntry := memory.GetProgressBar(shared.layerAlias, shared.progressBarAlias)
	if value > 0 {
		progressBarEntry.MaxValue = value
	}
}

func (shared *ProgressBarInstanceType) SetProgressBarLabel(label string) {
	progressBarEntry := memory.GetProgressBar(shared.layerAlias, shared.progressBarAlias)
	progressBarEntry.Label = label
}

func (shared *ProgressBarInstanceType) IncrementProgressBarValue() {
	progressBarEntry := memory.GetProgressBar(shared.layerAlias, shared.progressBarAlias)
	progressBarEntry.Value = progressBarEntry.Value + 1
	if progressBarEntry.Value > progressBarEntry.MaxValue {
		progressBarEntry.Value = progressBarEntry.MaxValue
	}
}

func (shared *ProgressBarInstanceType) GetProgressBarValueAsRatio() string {
	progressBarEntry := memory.GetProgressBar(shared.layerAlias, shared.progressBarAlias)
	valueAsString := fmt.Sprintf("%d/%d", progressBarEntry.Value, progressBarEntry.MaxValue)
	return valueAsString
}

func (shared *ProgressBarInstanceType) GetProgressBarValueAsPercent() string {
	var returnValue int
	progressBarEntry := memory.GetProgressBar(shared.layerAlias, shared.progressBarAlias)
	if progressBarEntry.MaxValue > 0 {
		returnValue = (progressBarEntry.Value * 100) / progressBarEntry.MaxValue
	}
	return fmt.Sprintf("%d", returnValue)
}

func (shared *progressBarType) Add(layerAlias string, progressBarAlias string, progressBarLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, value int, maxValue int, isBackgroundTransparent bool) ProgressBarInstanceType {
	memory.AddProgressBar(layerAlias, progressBarAlias, progressBarLabel, styleEntry, xLocation, yLocation, width, height, value, maxValue, isBackgroundTransparent)
	var progressBarInstance ProgressBarInstanceType
	progressBarInstance.layerAlias = layerAlias
	progressBarInstance.progressBarAlias = progressBarAlias
	return progressBarInstance
}

/*
DeleteButton allows you to remove a button from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a button which does not exist, then the request
will simply be ignored.
*/
func (shared *progressBarType) DeleteProgressBar(layerAlias string, progressBarAlias string) {
	memory.DeleteButton(layerAlias, progressBarAlias)
}

func (shared *progressBarType) DeleteAllProgressBars(layerAlias string) {
	memory.DeleteAllProgressBarsFromLayer(layerAlias)
}

/*
drawButtonsOnLayer allows you to draw all buttons on a given text layer.
*/
func (shared *progressBarType) drawProgressBarsOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for currentKey := range memory.ProgressBar.Entries[layerAlias] {
		progressBarEntry := memory.GetProgressBar(layerAlias, currentKey)
		drawProgressBar(&layerEntry, currentKey, progressBarEntry.Label, progressBarEntry.StyleEntry, progressBarEntry.XLocation, progressBarEntry.YLocation, progressBarEntry.Width, progressBarEntry.Height, progressBarEntry.Value, progressBarEntry.MaxValue, progressBarEntry.IsBackgroundTransparent)
	}
}

func drawProgressBar(layerEntry *types.LayerEntryType, progressBarAlias string, progressBarLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, currentValue int, maxValue int, isBackgroundTransparent bool) {
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.ProgressBarUnfilledForegroundColor
	attributeEntry.BackgroundColor = styleEntry.ProgressBarUnfilledBackgroundColor
	attributeEntry.CellType = constants.CellTypeProgressBar
	attributeEntry.CellControlAlias = progressBarAlias
	if height < 1 {
		height = 1
	}
	fillArea(layerEntry, attributeEntry, string(styleEntry.ProgressBarUnfilledPattern), xLocation, yLocation, width, height, constants.CellTypeProgressBar)
	progressBarWidth := float64(width) * float64(currentValue) / float64(maxValue)
	attributeEntry.ForegroundColor = styleEntry.ProgressBarFilledForegroundColor
	attributeEntry.BackgroundColor = styleEntry.ProgressBarFilledBackgroundColor
	fillArea(layerEntry, attributeEntry, string(styleEntry.ProgressBarFilledPattern), xLocation, yLocation, int(progressBarWidth), height, constants.CellTypeProgressBar)

	if len(progressBarLabel) > width {
		progressBarLabel = strings.ShortenString(progressBarLabel, width-3)
		progressBarLabel = progressBarLabel + "..."
	}
	centerXLocation := (width - len(progressBarLabel)) / 2
	centerYLocation := height / 2
	arrayOfRunes := stringformat.GetRunesFromString(progressBarLabel)
	attributeEntry.ForegroundColor = styleEntry.ProgressBarTextForegroundColor
	attributeEntry.BackgroundColor = styleEntry.ProgressBarTextBackgroundColor
	attributeEntry.IsBackgroundTransparent = isBackgroundTransparent
	printLayer(layerEntry, attributeEntry, xLocation+centerXLocation, yLocation+centerYLocation, arrayOfRunes)
}
