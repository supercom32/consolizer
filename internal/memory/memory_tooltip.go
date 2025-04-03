package memory

import (
	"fmt"
	"supercom32.net/consolizer/types"
)

var Tooltips = NewControlMemoryManager[types.TooltipEntryType]()

func AddTooltip(layerAlias string, tooltipAlias string, tooltipValue string, styleEntry types.TuiStyleEntryType, hotspotXLocation int, hotspotYLocation int, hotspotWidth int, hotspotHeight int, tooltipXLocation int, tooltipYLocation int, tooltipWidth int, tooltipHeight int, isLocationAbsolute bool, isBorderDrawn bool, hoverTime int) {
	tooltipEntry := types.NewTooltipEntry()
	tooltipEntry.StyleEntry = styleEntry
	tooltipEntry.Alias = tooltipAlias
	tooltipEntry.Value = tooltipValue
	tooltipEntry.HotspotXLocation = hotspotXLocation
	tooltipEntry.HotspotYLocation = hotspotYLocation
	tooltipEntry.HotspotWidth = hotspotWidth
	tooltipEntry.HotspotHeight = hotspotHeight
	tooltipEntry.TooltipXLocation = tooltipXLocation
	tooltipEntry.TooltipYLocation = tooltipYLocation
	tooltipEntry.TooltipWidth = tooltipWidth
	tooltipEntry.TooltipHeight = tooltipHeight
	tooltipEntry.IsLocationAbsolute = isLocationAbsolute
	tooltipEntry.IsBorderDrawn = isBorderDrawn
	tooltipEntry.HoverDisplayDelay = hoverTime

	Tooltips.Add(layerAlias, tooltipAlias, &tooltipEntry)
}

func GetTooltip(layerAlias string, tooltipAlias string) *types.TooltipEntryType {
	tooltipEntry := Tooltips.Get(layerAlias, tooltipAlias)
	if tooltipEntry == nil {
		panic(fmt.Sprintf("The requested tooltip with alias '%s' on layer '%s' could not be returned since it does not exist.", tooltipAlias, layerAlias))
	}
	return tooltipEntry
}

func IsTooltipExists(layerAlias string, tooltipAlias string) bool {
	return Tooltips.Get(layerAlias, tooltipAlias) != nil
}

func DeleteTooltip(layerAlias string, tooltipAlias string) {
	Tooltips.Remove(layerAlias, tooltipAlias)
}

func DeleteAllTooltipsFromLayer(layerAlias string) {
	Tooltips.RemoveAll(layerAlias)
}
