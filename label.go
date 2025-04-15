package consolizer

import (
	"strings"

	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/internal/memory"
	"supercom32.net/consolizer/internal/stringformat"
	"supercom32.net/consolizer/types"
)

type LabelInstanceType struct {
	BaseControlInstanceType
}

type labelType struct{}

var Label labelType
var Labels = memory.NewControlMemoryManager[types.LabelEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

func (shared *LabelInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeLabel)
}

func (shared *LabelInstanceType) SetLabelValue(value string) {
	labelEntry := Labels.Get(shared.layerAlias, shared.controlAlias)
	labelEntry.Text = value
}

func (shared *LabelInstanceType) Delete() *LabelInstanceType {
	if Labels.IsExists(shared.layerAlias, shared.controlAlias) {
		Labels.Remove(shared.layerAlias, shared.controlAlias)
	}
	return nil
}

func (shared *labelType) Add(layerAlias string, labelAlias string, labelValue string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int) LabelInstanceType {
	labelEntry := types.NewLabelEntry()
	labelEntry.StyleEntry = styleEntry
	labelEntry.Alias = labelAlias
	labelEntry.Text = labelValue
	labelEntry.XLocation = xLocation
	labelEntry.YLocation = yLocation
	labelEntry.Width = width
	labelEntry.TooltipAlias = stringformat.GetLastSortedUUID()
	// Use the ControlMemoryManager to add the label entry
	Labels.Add(layerAlias, labelAlias, &labelEntry)

	tooltipInstance := Tooltip.Add(layerAlias, labelEntry.TooltipAlias, "", styleEntry,
		labelEntry.XLocation, labelEntry.YLocation,
		labelEntry.Width+2, 1,
		labelEntry.XLocation, labelEntry.YLocation+1,
		labelEntry.Width+2, 3,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.SetEnabled(false)
	tooltipInstance.setParentControlAlias(labelAlias)
	var labelInstance LabelInstanceType
	labelInstance.layerAlias = layerAlias
	labelInstance.controlAlias = labelAlias
	labelInstance.controlType = "label"
	return labelInstance
}

func (shared *labelType) DeleteLabel(layerAlias string, labelAlias string) {
	Labels.Remove(layerAlias, labelAlias)
}

func (shared *labelType) DeleteAllLabels(layerAlias string) {
	Labels.RemoveAll(layerAlias)
}

/*
drawButtonsOnLayer allows you to draw all buttons on a given text layer.
*/
func (shared *labelType) drawLabelsOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentLabelEntry := range Labels.GetAllEntries(layerAlias) {
		labelEntry := currentLabelEntry
		drawLabel(&layerEntry, labelEntry.Alias, labelEntry.Text, labelEntry.StyleEntry, labelEntry.XLocation, labelEntry.YLocation, labelEntry.Width)
	}
}

func drawLabel(layerEntry *types.LayerEntryType, labelAlias string, labelValue string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int) {
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.LabelForegroundColor
	attributeEntry.BackgroundColor = styleEntry.LabelBackgroundColor
	attributeEntry.CellType = constants.CellTypeLabel
	attributeEntry.CellControlAlias = labelAlias
	emptyString := strings.Repeat(" ", width)
	printLayer(layerEntry, attributeEntry, xLocation, yLocation, stringformat.GetRunesFromString(emptyString))
	if len(labelValue) > width {
		labelValue = string([]rune(labelValue)[:width-3])
		labelValue = labelValue + "..."
	}
	arrayOfRunes := stringformat.GetRunesFromString(labelValue)
	printLayer(layerEntry, attributeEntry, xLocation, yLocation, arrayOfRunes)
}
