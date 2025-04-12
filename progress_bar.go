package consolizer

import (
	"fmt"

	"github.com/u2takey/go-utils/strings"
	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/internal/memory"
	"supercom32.net/consolizer/internal/stringformat"
	"supercom32.net/consolizer/types"
)

type ProgressBarInstanceType struct {
	layerAlias   string
	controlAlias string
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

func (shared *ProgressBarInstanceType) GetProgressBarValueAsPercent() string {
	var returnValue int
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if progressBarEntry.MaxValue > 0 {
		returnValue = (progressBarEntry.Value * 100) / progressBarEntry.MaxValue
	}
	return fmt.Sprintf("%d", returnValue)
}

func (shared *progressBarType) Add(layerAlias string, progressBarAlias string, progressBarLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, value int, maxValue int, isBackgroundTransparent bool) ProgressBarInstanceType {
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
	progressBarEntry.TooltipAlias = stringformat.GetLastSortedUUID()

	// Create associated tooltip (always created but disabled by default)
	Tooltip.Add(layerAlias, progressBarEntry.TooltipAlias, "", styleEntry,
		progressBarEntry.XLocation, progressBarEntry.YLocation,
		progressBarEntry.Width, progressBarEntry.Height,
		progressBarEntry.XLocation, progressBarEntry.YLocation+progressBarEntry.Height+1,
		progressBarEntry.Width, 3,
		false, true, constants.DefaultTooltipHoverTime)

	// Use the ControlMemoryManager to add the progress bar entry
	ProgressBars.Add(layerAlias, progressBarAlias, &progressBarEntry)
	var progressBarInstance ProgressBarInstanceType
	progressBarInstance.layerAlias = layerAlias
	progressBarInstance.controlAlias = progressBarAlias
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
		drawProgressBar(&layerEntry, progressBarEntry.Alias, progressBarEntry.Label, progressBarEntry.StyleEntry, progressBarEntry.XLocation, progressBarEntry.YLocation, progressBarEntry.Width, progressBarEntry.Height, progressBarEntry.Value, progressBarEntry.MaxValue, progressBarEntry.IsBackgroundTransparent)
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

// GetBounds returns the position and size of the progress bar
func (shared *ProgressBarInstanceType) GetBounds() (int, int, int, int) {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if progressBarEntry == nil {
		return 0, 0, 0, 0
	}
	return progressBarEntry.XLocation, progressBarEntry.YLocation, progressBarEntry.Width, progressBarEntry.Height
}

// SetPosition sets the position of the progress bar
func (shared *ProgressBarInstanceType) SetPosition(x, y int) *ProgressBarInstanceType {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if progressBarEntry != nil {
		progressBarEntry.XLocation = x
		progressBarEntry.YLocation = y
	}
	return shared
}

// SetSize sets the dimensions of the progress bar
func (shared *ProgressBarInstanceType) SetSize(width, height int) *ProgressBarInstanceType {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if progressBarEntry != nil {
		progressBarEntry.Width = width
		progressBarEntry.Height = height
	}
	return shared
}

// SetVisible shows or hides the progress bar
func (shared *ProgressBarInstanceType) SetVisible(visible bool) *ProgressBarInstanceType {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if progressBarEntry != nil {
		progressBarEntry.IsVisible = visible
	}
	return shared
}

// SetStyle sets the visual style of the progress bar
func (shared *ProgressBarInstanceType) SetStyle(style types.TuiStyleEntryType) *ProgressBarInstanceType {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if progressBarEntry != nil {
		progressBarEntry.StyleEntry = style
	}
	return shared
}

// SetTabIndex sets the tab order of the progress bar
func (shared *ProgressBarInstanceType) SetTabIndex(index int) *ProgressBarInstanceType {
	progressBarEntry := ProgressBars.Get(shared.layerAlias, shared.controlAlias)
	if progressBarEntry != nil {
		progressBarEntry.TabIndex = index
	}
	return shared
}
