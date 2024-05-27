package types

import (
	"encoding/json"
	"sync"
	"time"
)

// func DrawButton(LayerAlias string, Value string, StyleEntry TuiStyleEntryType, IsPressed bool, IsSelected bool, HotspotXLocation int, HotspotYLocation int, HotspotWidth int, Height int) {
type TooltipEntryType struct {
	Mutex              sync.Mutex
	StyleEntry         TuiStyleEntryType
	Alias              string
	Value              string
	HotspotXLocation   int
	HotspotYLocation   int
	HotspotWidth       int
	HotspotHeight      int
	TooltipXLocation   int
	TooltipYLocation   int
	TooltipWidth       int
	TooltipHeight      int
	IsLocationAbsolute bool
	IsBorderDrawn      bool
	IsDrawn            bool
	HoverDisplayDelay  int
	HoverStartTime     time.Time
	HoverXLocation     int
	HoverYLocation     int
}

func (shared TooltipEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		StyleEntry         TuiStyleEntryType
		LabelAlias         string
		LabelValue         string
		HotspotXLocation   int
		HotspotYLocation   int
		HotspotWidth       int
		HotspotHeight      int
		TooltipXLocation   int
		TooltipYLocation   int
		TooltipWidth       int
		TooltipHeight      int
		IsLocationAbsolute bool
		IsBorderDrawn      bool
		IsDrawn            bool
		HoverTime          int
		HoverStartTime     time.Time
		HoverXLocation     int
		HoverYLocation     int
	}{
		StyleEntry:         shared.StyleEntry,
		LabelAlias:         shared.Alias,
		LabelValue:         shared.Value,
		HotspotXLocation:   shared.HotspotXLocation,
		HotspotYLocation:   shared.HotspotYLocation,
		HotspotWidth:       shared.HotspotWidth,
		HotspotHeight:      shared.HotspotWidth,
		TooltipXLocation:   shared.TooltipXLocation,
		TooltipYLocation:   shared.TooltipYLocation,
		TooltipWidth:       shared.TooltipWidth,
		TooltipHeight:      shared.TooltipHeight,
		IsLocationAbsolute: shared.IsLocationAbsolute,
		IsBorderDrawn:      shared.IsBorderDrawn,
		IsDrawn:            shared.IsDrawn,
		HoverTime:          shared.HoverDisplayDelay,
		HoverStartTime:     shared.HoverStartTime,
		HoverXLocation:     shared.HoverXLocation,
		HoverYLocation:     shared.HoverYLocation,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared TooltipEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewTooltipEntry(existingButtonEntry ...*TooltipEntryType) TooltipEntryType {
	var tooltipEntry TooltipEntryType
	if existingButtonEntry != nil {
		tooltipEntry.StyleEntry = NewTuiStyleEntry(&existingButtonEntry[0].StyleEntry)
		tooltipEntry.Alias = existingButtonEntry[0].Alias
		tooltipEntry.Value = existingButtonEntry[0].Value
		tooltipEntry.HotspotXLocation = existingButtonEntry[0].HotspotXLocation
		tooltipEntry.HotspotYLocation = existingButtonEntry[0].HotspotYLocation
		tooltipEntry.HotspotWidth = existingButtonEntry[0].HotspotWidth
		tooltipEntry.HotspotHeight = existingButtonEntry[0].HotspotHeight
		tooltipEntry.TooltipXLocation = existingButtonEntry[0].TooltipXLocation
		tooltipEntry.TooltipYLocation = existingButtonEntry[0].TooltipYLocation
		tooltipEntry.TooltipWidth = existingButtonEntry[0].TooltipWidth
		tooltipEntry.TooltipHeight = existingButtonEntry[0].TooltipHeight
		tooltipEntry.IsLocationAbsolute = existingButtonEntry[0].IsLocationAbsolute
		tooltipEntry.IsBorderDrawn = existingButtonEntry[0].IsBorderDrawn
		tooltipEntry.HoverDisplayDelay = existingButtonEntry[0].HoverDisplayDelay
		tooltipEntry.HoverStartTime = existingButtonEntry[0].HoverStartTime
	}
	return tooltipEntry
}

func IsTooltipEntryEqual(sourceTooltipEntry *TooltipEntryType, targetTooltipEntry *TooltipEntryType) bool {
	if sourceTooltipEntry.StyleEntry == targetTooltipEntry.StyleEntry &&
		sourceTooltipEntry.Alias == targetTooltipEntry.Alias &&
		sourceTooltipEntry.Value == targetTooltipEntry.Value &&
		sourceTooltipEntry.HotspotXLocation == targetTooltipEntry.HotspotXLocation &&
		sourceTooltipEntry.HotspotYLocation == targetTooltipEntry.HotspotYLocation &&
		sourceTooltipEntry.HotspotWidth == targetTooltipEntry.HotspotWidth &&
		sourceTooltipEntry.TooltipXLocation == targetTooltipEntry.TooltipXLocation &&
		sourceTooltipEntry.TooltipYLocation == targetTooltipEntry.TooltipYLocation &&
		sourceTooltipEntry.TooltipWidth == targetTooltipEntry.TooltipWidth &&
		sourceTooltipEntry.IsBorderDrawn == targetTooltipEntry.IsBorderDrawn &&
		sourceTooltipEntry.IsLocationAbsolute == targetTooltipEntry.IsLocationAbsolute &&
		sourceTooltipEntry.HoverDisplayDelay == targetTooltipEntry.HoverDisplayDelay &&
		sourceTooltipEntry.HoverStartTime == targetTooltipEntry.HoverStartTime {
		return true
	}
	return false
}
