package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
)

type tooltipMemoryType struct {
	sync.Mutex
	Entries map[string]map[string]*types.TooltipEntryType
}

var Tooltip tooltipMemoryType

func InitializeTooltipMemory() {
	Tooltip.Entries = make(map[string]map[string]*types.TooltipEntryType)
}

func AddTooltip(layerAlias string, tooltipAlias string, tooltipValue string, styleEntry types.TuiStyleEntryType, hotspotXLocation int, hotspotYLocation int, hotspotWidth int, hotspotHeight int, tooltipXLocation int, tooltipYLocation int, tooltipWidth int, tooltipHeight int, isLocationAbsolute bool, isBorderDrawn bool) {
	Tooltip.Lock()
	defer func() {
		Tooltip.Unlock()
	}()
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
	if Tooltip.Entries[layerAlias] == nil {
		Tooltip.Entries[layerAlias] = make(map[string]*types.TooltipEntryType)
	}
	Tooltip.Entries[layerAlias][tooltipAlias] = &tooltipEntry
}

func GetTooltip(layerAlias string, tooltipAlias string) *types.TooltipEntryType {
	Tooltip.Lock()
	defer func() {
		Tooltip.Unlock()
	}()
	if Tooltip.Entries[layerAlias][tooltipAlias] == nil {
		panic(fmt.Sprintf("The requested label with alias '%s' on layer '%s' could not be returned since it does not exist.", tooltipAlias, layerAlias))
	}
	return Tooltip.Entries[layerAlias][tooltipAlias]
}

func IsTooltipExists(layerAlias string, tooltipAlias string) bool {
	Tooltip.Lock()
	defer func() {
		Tooltip.Unlock()
	}()
	if Tooltip.Entries[layerAlias][tooltipAlias] == nil {
		return false
	}
	return true
}

func DeleteTooltip(layerAlias string, tooltipAlias string) {
	Tooltip.Lock()
	defer func() {
		Tooltip.Unlock()
	}()
	delete(Tooltip.Entries[layerAlias], tooltipAlias)
}

func DeleteAllTooltipsFromLayer(layerAlias string) {
	Tooltip.Lock()
	defer func() {
		Tooltip.Unlock()
	}()
	for entryToRemove := range Tooltip.Entries[layerAlias] {
		delete(Tooltip.Entries[layerAlias], entryToRemove)
	}
}
