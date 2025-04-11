package types

import (
	"encoding/json"
	"time"
)

// func DrawButton(LayerAlias string, Value string, StyleEntry TuiStyleEntryType, IsPressed bool, IsSelected bool, HotspotXLocation int, HotspotYLocation int, HotspotWidth int, Height int) {
type TooltipEntryType struct {
	BaseControlType
	Text               string
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
	IsDrawn            bool
	HoverStartTime     time.Time
	HoverXLocation     int
	HoverYLocation     int
	HoverDisplayDelay  int
}

/*
GetAlias allows you to retrieve the alias of a tooltip control. In addition, the following
information should be noted:

- Returns the unique identifier for the tooltip.
- This alias is used to reference the tooltip in other operations.
- The alias is set when the tooltip is created.
*/
func (shared TooltipEntryType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON allows you to serialize a tooltip control to JSON. In addition, the following
information should be noted:

- Converts the tooltip's state to a JSON representation.
- Includes the base control properties and tooltip-specific fields.
- Used for saving and loading tooltip configurations.
*/
func (shared TooltipEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		Text               string
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
		IsDrawn            bool
		HoverStartTime     time.Time
		HoverXLocation     int
		HoverYLocation     int
		HoverDisplayDelay  int
	}{
		BaseControlType:    shared.BaseControlType,
		Text:               shared.Text,
		Value:              shared.Value,
		HotspotXLocation:   shared.HotspotXLocation,
		HotspotYLocation:   shared.HotspotYLocation,
		HotspotWidth:       shared.HotspotWidth,
		HotspotHeight:      shared.HotspotHeight,
		TooltipXLocation:   shared.TooltipXLocation,
		TooltipYLocation:   shared.TooltipYLocation,
		TooltipWidth:       shared.TooltipWidth,
		TooltipHeight:      shared.TooltipHeight,
		IsLocationAbsolute: shared.IsLocationAbsolute,
		IsDrawn:            shared.IsDrawn,
		HoverStartTime:     shared.HoverStartTime,
		HoverXLocation:     shared.HoverXLocation,
		HoverYLocation:     shared.HoverYLocation,
		HoverDisplayDelay:  shared.HoverDisplayDelay,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump allows you to get a JSON string representation of a tooltip control. In addition,
the following information should be noted:

- Returns a formatted JSON string of the tooltip's state.
- Useful for debugging and logging purposes.
- Panics if JSON marshaling fails.
*/
func (shared TooltipEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
GetTooltipAlias allows you to retrieve the alias of a tooltip control. In addition, the following
information should be noted:

- Returns the unique identifier for the tooltip.
- This is a convenience method that delegates to GetAlias.
- The alias is used to reference the tooltip in other operations.
*/
func GetTooltipAlias(entry *TooltipEntryType) string {
	return entry.Alias
}

/*
NewTooltipEntry allows you to create a new tooltip control. In addition, the following
information should be noted:

- Initializes a tooltip with default values.
- Can optionally copy properties from an existing tooltip.
- Sets up the base control properties and tooltip-specific fields.
*/
func NewTooltipEntry(existingButtonEntry ...*TooltipEntryType) TooltipEntryType {
	var tooltipEntry TooltipEntryType
	tooltipEntry.BaseControlType = NewBaseControl()

	if existingButtonEntry != nil {
		tooltipEntry.BaseControlType = existingButtonEntry[0].BaseControlType
		tooltipEntry.Text = existingButtonEntry[0].Text
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
		tooltipEntry.HoverDisplayDelay = existingButtonEntry[0].HoverDisplayDelay
	}
	return tooltipEntry
}

/*
IsTooltipEntryEqual allows you to compare two tooltip controls for equality. In addition, the following
information should be noted:

- Compares all properties of both tooltips.
- Returns true if all properties match, false otherwise.
- Used for change detection and state synchronization.
*/
func IsTooltipEntryEqual(sourceTooltipEntry *TooltipEntryType, targetTooltipEntry *TooltipEntryType) bool {
	if sourceTooltipEntry.BaseControlType == targetTooltipEntry.BaseControlType &&
		sourceTooltipEntry.Text == targetTooltipEntry.Text &&
		sourceTooltipEntry.Value == targetTooltipEntry.Value &&
		sourceTooltipEntry.HotspotXLocation == targetTooltipEntry.HotspotXLocation &&
		sourceTooltipEntry.HotspotYLocation == targetTooltipEntry.HotspotYLocation &&
		sourceTooltipEntry.HotspotWidth == targetTooltipEntry.HotspotWidth &&
		sourceTooltipEntry.HotspotHeight == targetTooltipEntry.HotspotHeight &&
		sourceTooltipEntry.TooltipXLocation == targetTooltipEntry.TooltipXLocation &&
		sourceTooltipEntry.TooltipYLocation == targetTooltipEntry.TooltipYLocation &&
		sourceTooltipEntry.TooltipWidth == targetTooltipEntry.TooltipWidth &&
		sourceTooltipEntry.TooltipHeight == targetTooltipEntry.TooltipHeight &&
		sourceTooltipEntry.IsLocationAbsolute == targetTooltipEntry.IsLocationAbsolute &&
		sourceTooltipEntry.HoverDisplayDelay == targetTooltipEntry.HoverDisplayDelay {
		return true
	}
	return false
}
