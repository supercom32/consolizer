package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
	"github.com/supercom32/consolizer/types"
	"strings"
)

type TooltipInstanceType struct {
	layerAlias    string
	tooltiplAlias string
}

type tooltipType struct{}

var Tooltip tooltipType

func (shared *TooltipInstanceType) SetTooltipValue(value string) {
	labelEntry := memory.GetLabel(shared.layerAlias, shared.tooltiplAlias)
	labelEntry.Value = value
}

func (shared *tooltipType) Add(layerAlias string, tooltipAlias string, tooltipValue string, styleEntry types.TuiStyleEntryType, hotspotXLocation int, hotspotYLocation int, hotspotWidth int, hotspotHeight int, tooltipXLocation int, tooltipYLocation int, tooltipWidth int, tooltipHeight int, isLocationAbsolute bool, isBorderDrawn bool) TooltipInstanceType {
	memory.AddTooltip(layerAlias, tooltipAlias, tooltipValue, styleEntry, hotspotXLocation, hotspotYLocation, hotspotWidth, hotspotHeight, tooltipXLocation, tooltipYLocation, tooltipWidth, tooltipHeight, isLocationAbsolute, isBorderDrawn)
	var tooltipInstance TooltipInstanceType
	tooltipInstance.layerAlias = layerAlias
	tooltipInstance.tooltiplAlias = tooltipAlias
	return tooltipInstance
}

/*
DeleteButton allows you to remove a button from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a button which does not exist, then the request
will simply be ignored.
*/
func (shared *tooltipType) DeleteTooltip(layerAlias string, labelAlias string) {
	memory.DeleteTooltip(layerAlias, labelAlias)
}

func (shared *tooltipType) DeleteAllTooltips(layerAlias string) {
	memory.DeleteAllTooltipsFromLayer(layerAlias)
}

/*
drawButtonsOnLayer allows you to draw all buttons on a given text layer.
*/
func (shared *tooltipType) drawTooltipsOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for currentKey := range memory.Tooltip.Entries[layerAlias] {
		tooltipEntry := memory.GetTooltip(layerAlias, currentKey)
		shared.drawTooltip(&layerEntry, tooltipEntry)
	}
}

// TOOD: This method should really just take in a tooltip entry instead?
func (shared *tooltipType) drawTooltip(layerEntry *types.LayerEntryType, tooltipEntry *types.TooltipEntryType) {
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = tooltipEntry.StyleEntry.TooltipForegroundColor
	attributeEntry.BackgroundColor = tooltipEntry.StyleEntry.TooltipBackgroundColor
	attributeEntry.CellType = constants.CellTypeTooltip
	attributeEntry.CellControlAlias = tooltipEntry.Alias
	emptyString := strings.Repeat(string(constants.NullRune), tooltipEntry.HotspotWidth)
	printLayer(layerEntry, attributeEntry, tooltipEntry.HotspotXLocation, tooltipEntry.HotspotYLocation, stringformat.GetRunesFromString(emptyString))
	if tooltipEntry.IsDrawn {
		formattedLabel := tooltipEntry.Value
		if len(tooltipEntry.Value) > tooltipEntry.TooltipWidth {
			formattedLabel = string([]rune(tooltipEntry.Value)[:tooltipEntry.TooltipWidth-3])
			formattedLabel = formattedLabel + "..."
		}
		arrayOfRunes := stringformat.GetRunesFromString(formattedLabel)
		printLayer(layerEntry, attributeEntry, tooltipEntry.TooltipXLocation, tooltipEntry.TooltipYLocation, arrayOfRunes)
	}
}

func (shared *tooltipType) updateMouseEventTooltip() bool {
	isScreenUpdateRequired := false
	var characterEntry types.CharacterEntryType
	mouseXLocation, mouseYLocation, _, _ := memory.GetMouseStatus()
	characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	if characterEntry.AttributeEntry.CellType == constants.CellTypeTooltip && eventStateMemory.stateId == constants.EventStateNone && memory.IsTooltipExists(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias) {
		tooltipEntry := memory.GetTooltip(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
		setPreviouslyHighlightedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeTooltip)
		tooltipEntry.IsDrawn = true
		isScreenUpdateRequired = true
	} else {
		if eventStateMemory.previouslyHighlightedControl.controlType == constants.CellTypeTooltip {
			for currentLayer, _ := range memory.Tooltip.Entries {
				for _, currentTooltipEntry := range memory.Tooltip.Entries[currentLayer] {
					currentTooltipEntry.IsDrawn = false
				}
			}
			setPreviouslyHighlightedControl("", "", constants.NullControlType)
			isScreenUpdateRequired = true
		}
	}
	return isScreenUpdateRequired
}
