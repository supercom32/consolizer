package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
	"github.com/supercom32/consolizer/types"
	"strings"
)

type LabelInstanceType struct {
	layerAlias string
	labelAlias string
}

type labelType struct{}

var Label labelType

func (shared *LabelInstanceType) SetLabelValue(value string) {
	labelEntry := memory.GetLabel(shared.layerAlias, shared.labelAlias)
	labelEntry.Value = value
}

func (shared *LabelInstanceType) Delete() *LabelInstanceType {
	if memory.IsLabelExists(shared.layerAlias, shared.labelAlias) {
		memory.DeleteLabel(shared.layerAlias, shared.labelAlias)
	}
	return nil
}

func (shared *labelType) Add(layerAlias string, labelAlias string, labelValue string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int) LabelInstanceType {
	memory.AddLabel(layerAlias, labelAlias, labelValue, styleEntry, xLocation, yLocation, width)
	var labelInstance LabelInstanceType
	labelInstance.layerAlias = layerAlias
	labelInstance.labelAlias = labelAlias
	return labelInstance
}

func (shared *labelType) DeleteLabel(layerAlias string, labelAlias string) {
	memory.DeleteLabel(layerAlias, labelAlias)
}

func (shared *labelType) DeleteAllLabels(layerAlias string) {
	memory.DeleteAllLabelsFromLayer(layerAlias)
}

/*
drawButtonsOnLayer allows you to draw all buttons on a given text layer.
*/
func (shared *labelType) drawLabelsOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for currentKey := range memory.Label.Entries[layerAlias] {
		labelEntry := memory.GetLabel(layerAlias, currentKey)
		drawLabel(&layerEntry, currentKey, labelEntry.Value, labelEntry.StyleEntry, labelEntry.XLocation, labelEntry.YLocation, labelEntry.Width)
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
