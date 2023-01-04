package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
	"github.com/supercom32/consolizer/types"
	"github.com/u2takey/go-utils/strings"
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

func (shared *labelType) Add(layerAlias string, labelAlias string, labelValue string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int) LabelInstanceType {
	memory.AddLabel(layerAlias, labelAlias, labelValue, styleEntry, xLocation, yLocation, width)
	var labelInstance LabelInstanceType
	labelInstance.layerAlias = layerAlias
	labelInstance.labelAlias = labelAlias
	return labelInstance
}

/*
DeleteButton allows you to remove a button from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a button which does not exist, then the request
will simply be ignored.
*/
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
	if len(labelValue) > width {
		labelValue = strings.ShortenString(labelValue, width-3)
		labelValue = labelValue + "..."
	}
	arrayOfRunes := stringformat.GetRunesFromString(labelValue)
	printLayer(layerEntry, attributeEntry, xLocation, yLocation, arrayOfRunes)
}
