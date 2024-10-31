package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
	"github.com/supercom32/consolizer/types"
	"time"
)

type TooltipInstanceType struct {
	layerAlias   string
	tooltipAlias string
}

type tooltipType struct{}

var Tooltip tooltipType

func (shared *TooltipInstanceType) Delete() string {
	if memory.IsTooltipExists(shared.layerAlias, shared.tooltipAlias) {
		memory.DeleteTooltip(shared.layerAlias, shared.tooltipAlias)
	}
	return ""
}

/*
SetTooltipValue allows you to set the value of the tooltip associated with the TooltipInstanceType.
This function updates the value of the tooltip label identified by the layerAlias and tooltipAlias fields.
*/
func (shared *TooltipInstanceType) SetTooltipValue(value string) {
	labelEntry := memory.GetLabel(shared.layerAlias, shared.tooltipAlias)
	labelEntry.Value = value
}

func (shared *tooltipType) Add(layerAlias string, tooltipAlias string, tooltipValue string, styleEntry types.TuiStyleEntryType, hotspotXLocation int, hotspotYLocation int, hotspotWidth int, hotspotHeight int, tooltipXLocation int, tooltipYLocation int, tooltipWidth int, tooltipHeight int, isLocationAbsolute bool, isBorderDrawn bool, hoverTime int) TooltipInstanceType {
	memory.AddTooltip(layerAlias, tooltipAlias, tooltipValue, styleEntry, hotspotXLocation, hotspotYLocation, hotspotWidth, hotspotHeight, tooltipXLocation, tooltipYLocation, tooltipWidth, tooltipHeight, isLocationAbsolute, isBorderDrawn, hoverTime)
	var tooltipInstance TooltipInstanceType
	tooltipInstance.layerAlias = layerAlias
	tooltipInstance.tooltipAlias = tooltipAlias
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
	fillAreaWithControlAlias(layerEntry, attributeEntry.CellType, attributeEntry.CellControlAlias, tooltipEntry.HotspotXLocation, tooltipEntry.HotspotYLocation, tooltipEntry.HotspotWidth, tooltipEntry.HotspotHeight, constants.NullCellControlLocation)
	if tooltipEntry.IsDrawn {
		calculatedXLocation := tooltipEntry.TooltipXLocation - 2
		calculatedYLocation := tooltipEntry.TooltipYLocation - 1
		calculatedWidth := tooltipEntry.TooltipWidth + 1
		calculatedHeight := tooltipEntry.TooltipHeight
		if !tooltipEntry.IsLocationAbsolute {
			mouseXLocation, mouseYLocation, _, _ := memory.GetMouseStatus()
			calculatedXLocation = mouseXLocation + tooltipEntry.TooltipXLocation - 2
			calculatedYLocation = mouseYLocation + tooltipEntry.TooltipYLocation - 1
			calculatedWidth = tooltipEntry.TooltipWidth + 1
			calculatedHeight = tooltipEntry.TooltipHeight
		}
		if len(tooltipEntry.Value) > tooltipEntry.TooltipWidth {
			fillArea(layerEntry, attributeEntry, " ", calculatedXLocation, calculatedYLocation, calculatedWidth, calculatedHeight, constants.NullCellControlLocation)
		}
		if tooltipEntry.IsBorderDrawn {
			drawBorder(layerEntry, tooltipEntry.StyleEntry, attributeEntry, calculatedXLocation, calculatedYLocation, calculatedWidth, calculatedHeight, false)
		}
		formattedLabel := tooltipEntry.Value
		arrayOfRunes := stringformat.GetRunesFromString(formattedLabel)
		printLayerWithWordWrap(layerEntry, attributeEntry, calculatedXLocation+2, calculatedYLocation+1, calculatedWidth-1, arrayOfRunes)
	}
}

func (shared *tooltipType) updateMouseEventTooltip() bool {
	isScreenUpdateRequired := false
	var characterEntry types.CharacterEntryType
	mouseXLocation, mouseYLocation, _, _ := memory.GetMouseStatus()
	characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	if characterEntry.AttributeEntry.CellType == constants.CellTypeTooltip && eventStateMemory.stateId == constants.EventStateNone && memory.IsTooltipExists(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias) {
		tooltipEntry := memory.GetTooltip(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
		mouseXLocation, mouseYLocation, _, _ = memory.GetMouseStatus()
		if tooltipEntry.HoverStartTime == (time.Time{}) {
			// If no start time was defined, do it now.
			setPreviouslyHighlightedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeTooltip)
			tooltipEntry.HoverStartTime = time.Now()
			tooltipEntry.HoverXLocation = mouseXLocation
			tooltipEntry.HoverYLocation = mouseYLocation
			return isScreenUpdateRequired
		}
		if tooltipEntry.HoverXLocation != mouseXLocation || tooltipEntry.HoverYLocation != mouseYLocation {
			tooltipEntry.HoverStartTime = time.Time{}
			return isScreenUpdateRequired
		}
		elapsedTime := time.Since(tooltipEntry.HoverStartTime)
		if elapsedTime >= time.Duration(tooltipEntry.HoverDisplayDelay)*time.Millisecond {
			setPreviouslyHighlightedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeTooltip)
			tooltipEntry.IsDrawn = true
			isScreenUpdateRequired = true
		}
	} else {
		if eventStateMemory.previouslyHighlightedControl.controlType == constants.CellTypeTooltip {
			for currentLayer, _ := range memory.Tooltip.Entries {
				for _, currentTooltipEntry := range memory.Tooltip.Entries[currentLayer] {
					currentTooltipEntry.IsDrawn = false
					currentTooltipEntry.HoverStartTime = time.Time{}
				}
			}
			setPreviouslyHighlightedControl("", "", constants.NullControlType)
			isScreenUpdateRequired = true
		}
	}
	return isScreenUpdateRequired
}
