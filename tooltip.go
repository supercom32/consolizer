package consolizer

import (
	"time"

	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/internal/memory"
	"supercom32.net/consolizer/internal/stringformat"
	"supercom32.net/consolizer/types"
)

type TooltipInstanceType struct {
	layerAlias   string
	controlAlias string
}

type tooltipType struct{}

var Tooltip tooltipType
var Tooltips = memory.NewControlMemoryManager[types.TooltipEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

func (shared *TooltipInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeTooltip)
}

func (shared *TooltipInstanceType) Delete() string {
	if Tooltips.IsExists(shared.layerAlias, shared.controlAlias) {
		Tooltips.Remove(shared.layerAlias, shared.controlAlias)
	}
	return ""
}

/*
SetTooltipValue allows you to set the value of the tooltip associated with the TooltipInstanceType.
This function updates the value of the tooltip label identified by the layerAlias and tooltipAlias fields.
*/
func (shared *TooltipInstanceType) SetTooltipValue(text string) {
	labelEntry := Labels.Get(shared.layerAlias, shared.controlAlias)
	labelEntry.Text = text
}

func (shared *tooltipType) Add(layerAlias string, tooltipAlias string, tooltipValue string, styleEntry types.TuiStyleEntryType, hotspotXLocation int, hotspotYLocation int, hotspotWidth int, hotspotHeight int, tooltipXLocation int, tooltipYLocation int, tooltipWidth int, tooltipHeight int, isLocationAbsolute bool, isBorderDrawn bool, hoverTime int) TooltipInstanceType {
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
	var tooltipInstance TooltipInstanceType
	tooltipInstance.layerAlias = layerAlias
	tooltipInstance.controlAlias = tooltipAlias
	return tooltipInstance
}

/*
DeleteButton allows you to remove a button from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a button which does not exist, then the request
will simply be ignored.
*/
func (shared *tooltipType) DeleteTooltip(layerAlias string, labelAlias string) {
	Tooltips.Remove(layerAlias, labelAlias)
}

func (shared *tooltipType) DeleteAllTooltips(layerAlias string) {
	Tooltips.RemoveAll(layerAlias)
}

/*
drawButtonsOnLayer allows you to draw all buttons on a given text layer.
*/
func (shared *tooltipType) drawTooltipsOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentTooltipEntry := range Tooltips.GetAllEntries(layerAlias) {
		tooltipEntry := currentTooltipEntry
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
			mouseXLocation, mouseYLocation, _, _ := GetMouseStatus()
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

func (shared *tooltipType) updateMouseEvent() bool {
	isScreenUpdateRequired := false
	var characterEntry types.CharacterEntryType
	mouseXLocation, mouseYLocation, _, _ := GetMouseStatus()
	characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	if characterEntry.AttributeEntry.CellType == constants.CellTypeTooltip && eventStateMemory.stateId == constants.EventStateNone && Tooltips.IsExists(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias) {
		tooltipEntry := Tooltips.Get(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
		mouseXLocation, mouseYLocation, _, _ = GetMouseStatus()
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
			for _, currentTooltipEntry := range Tooltips.GetAllEntriesOverall() {
				currentTooltipEntry.IsDrawn = false
				currentTooltipEntry.HoverStartTime = time.Time{}
			}
			setPreviouslyHighlightedControl("", "", constants.NullControlType)
			isScreenUpdateRequired = true
		}
	}
	return isScreenUpdateRequired
}

// GetBounds returns the position and size of the tooltip
func (shared *TooltipInstanceType) GetBounds() (int, int, int, int) {
	tooltipEntry := Tooltips.Get(shared.layerAlias, shared.controlAlias)
	if tooltipEntry == nil {
		return 0, 0, 0, 0
	}
	return tooltipEntry.XLocation, tooltipEntry.YLocation, tooltipEntry.Width, tooltipEntry.Height
}

// SetPosition sets the position of the tooltip
func (shared *TooltipInstanceType) SetPosition(x, y int) *TooltipInstanceType {
	tooltipEntry := Tooltips.Get(shared.layerAlias, shared.controlAlias)
	if tooltipEntry != nil {
		tooltipEntry.XLocation = x
		tooltipEntry.YLocation = y
	}
	return shared
}

// SetSize sets the dimensions of the tooltip
func (shared *TooltipInstanceType) SetSize(width, height int) *TooltipInstanceType {
	tooltipEntry := Tooltips.Get(shared.layerAlias, shared.controlAlias)
	if tooltipEntry != nil {
		tooltipEntry.Width = width
		tooltipEntry.Height = height
	}
	return shared
}

// SetVisible shows or hides the tooltip
func (shared *TooltipInstanceType) SetVisible(visible bool) *TooltipInstanceType {
	tooltipEntry := Tooltips.Get(shared.layerAlias, shared.controlAlias)
	if tooltipEntry != nil {
		tooltipEntry.IsVisible = visible
	}
	return shared
}

// SetStyle sets the visual style of the tooltip
func (shared *TooltipInstanceType) SetStyle(style types.TuiStyleEntryType) *TooltipInstanceType {
	tooltipEntry := Tooltips.Get(shared.layerAlias, shared.controlAlias)
	if tooltipEntry != nil {
		tooltipEntry.StyleEntry = style
	}
	return shared
}

// SetTabIndex sets the tab order of the tooltip
func (shared *TooltipInstanceType) SetTabIndex(index int) *TooltipInstanceType {
	tooltipEntry := Tooltips.Get(shared.layerAlias, shared.controlAlias)
	if tooltipEntry != nil {
		tooltipEntry.TabIndex = index
	}
	return shared
}

// SetEnabled enables or disables the tooltip
func (shared *TooltipInstanceType) SetEnabled(enabled bool) *TooltipInstanceType {
	tooltipEntry := Tooltips.Get(shared.layerAlias, shared.controlAlias)
	if tooltipEntry != nil {
		tooltipEntry.IsEnabled = enabled
	}
	return shared
}
