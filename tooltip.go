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

func (shared *TooltipInstanceType) SetEnabled(enabled bool) *TooltipInstanceType {
	tooltipEntry := Tooltips.Get(shared.layerAlias, shared.controlAlias)
	if tooltipEntry != nil {
		tooltipEntry.IsEnabled = enabled
	}
	return shared
}

func (shared *tooltipType) Add(layerAlias string, tooltipAlias string, tooltipText string, styleEntry types.TuiStyleEntryType, hotspotXLocation int, hotspotYLocation int, hotspotWidth int, hotspotHeight int, tooltipXLocation int, tooltipYLocation int, tooltipWidth int, tooltipHeight int, isLocationAbsolute bool, isBorderDrawn bool, hoverTime int) TooltipInstanceType {
	tooltipEntry := types.NewTooltipEntry()
	tooltipEntry.StyleEntry = styleEntry
	tooltipEntry.Alias = tooltipAlias
	tooltipEntry.Text = tooltipText
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
func (shared *tooltipType) drawTooltipHotspotZonesOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentTooltipEntry := range Tooltips.GetAllEntries(layerAlias) {
		tooltipEntry := currentTooltipEntry
		shared.drawTooltipHotspot(&layerEntry, tooltipEntry)
	}
}

// TOOD: This method should really just take in a tooltip entry instead?
func (shared *tooltipType) drawTooltipHotspot(layerEntry *types.LayerEntryType, tooltipEntry *types.TooltipEntryType) {
	if !tooltipEntry.IsEnabled {
		return
	}
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = tooltipEntry.StyleEntry.TooltipForegroundColor
	attributeEntry.BackgroundColor = tooltipEntry.StyleEntry.TooltipBackgroundColor
	attributeEntry.CellType = constants.CellTypeTooltip
	attributeEntry.CellControlAlias = tooltipEntry.Alias
	if tooltipEntry.ParentControlAlias == "" { // If a parent exists, do not overwrite the parent's attributes.
		fillAreaWithControlAlias(layerEntry, attributeEntry.CellType, attributeEntry.CellControlAlias, tooltipEntry.HotspotXLocation, tooltipEntry.HotspotYLocation, tooltipEntry.HotspotWidth, tooltipEntry.HotspotHeight, constants.NullCellControlLocation)
	}
}

func (shared *tooltipType) renderAllTooltips(layerEntry types.LayerEntryType) {
	for _, currentTooltipEntry := range Tooltips.GetAllEntriesOverall() {
		tooltipEntry := currentTooltipEntry
		shared.renderTooltip(&layerEntry, tooltipEntry)
	}
}

func (shared *tooltipType) renderTooltip(layerEntry *types.LayerEntryType, tooltipEntry *types.TooltipEntryType) {
	if !tooltipEntry.IsEnabled {
		return
	}
	if !tooltipEntry.IsDrawn {
		return
	}
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = tooltipEntry.StyleEntry.TooltipForegroundColor
	attributeEntry.BackgroundColor = tooltipEntry.StyleEntry.TooltipBackgroundColor
	attributeEntry.CellType = constants.CellTypeTooltip
	attributeEntry.CellControlAlias = tooltipEntry.Alias
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
	fillArea(layerEntry, attributeEntry, " ", calculatedXLocation, calculatedYLocation, calculatedWidth, calculatedHeight, constants.NullCellControlLocation)
	if tooltipEntry.IsBorderDrawn {
		drawBorder(layerEntry, tooltipEntry.StyleEntry, attributeEntry, calculatedXLocation, calculatedYLocation, calculatedWidth, calculatedHeight, false)
	}
	formattedLabel := tooltipEntry.Text
	arrayOfRunes := stringformat.GetRunesFromString(formattedLabel)
	printLayerWithWordWrap(layerEntry, attributeEntry, calculatedXLocation+2, calculatedYLocation+1, calculatedWidth-1, arrayOfRunes)
}

func (shared *tooltipType) getTooltipFromCharacterEntry(entry types.CharacterEntryType) *types.TooltipEntryType {
	layer := entry.LayerAlias
	alias := entry.AttributeEntry.CellControlAlias

	switch entry.AttributeEntry.CellType {
	case constants.CellTypeButton:
		if Buttons.IsExists(layer, alias) {
			button := Buttons.Get(layer, alias)
			if button.TooltipAlias != "" {
				return Tooltips.Get(layer, button.TooltipAlias)
			}
		}
	case constants.CellTypeLabel:
		if Labels.IsExists(layer, alias) {
			label := Labels.Get(layer, alias)
			if label.TooltipAlias != "" {
				return Tooltips.Get(layer, label.TooltipAlias)
			}
		}
	case constants.CellTypeCheckbox:
		if Checkboxes.IsExists(layer, alias) {
			checkbox := Checkboxes.Get(layer, alias)
			if checkbox.TooltipAlias != "" {
				return Tooltips.Get(layer, checkbox.TooltipAlias)
			}
		}
	case constants.CellTypeRadioButton:
		if RadioButtons.IsExists(layer, alias) {
			radio := RadioButtons.Get(layer, alias)
			if radio.TooltipAlias != "" {
				return Tooltips.Get(layer, radio.TooltipAlias)
			}
		}
	case constants.CellTypeTextField:
		if TextFields.IsExists(layer, alias) {
			textField := TextFields.Get(layer, alias)
			if textField.TooltipAlias != "" {
				return Tooltips.Get(layer, textField.TooltipAlias)
			}
		}
	case constants.CellTypeTextbox:
		if Textboxes.IsExists(layer, alias) {
			textbox := Textboxes.Get(layer, alias)
			if textbox.TooltipAlias != "" {
				return Tooltips.Get(layer, textbox.TooltipAlias)
			}
		}
	case constants.CellTypeProgressBar:
		if ProgressBars.IsExists(layer, alias) {
			progressBar := ProgressBars.Get(layer, alias)
			if progressBar.TooltipAlias != "" {
				return Tooltips.Get(layer, progressBar.TooltipAlias)
			}
		}
	case constants.CellTypeSelectorItem:
		if Selectors.IsExists(layer, alias) {
			selector := Selectors.Get(layer, alias)
			if selector.TooltipAlias != "" {
				return Tooltips.Get(layer, selector.TooltipAlias)
			}
		}
	case constants.CellTypeTooltip:
		if Tooltips.IsExists(layer, alias) {
			return Tooltips.Get(layer, alias)
		}
	}
	return nil
}

func (shared *tooltipType) updateMouseEvent() bool {
	isScreenUpdateRequired := false
	mouseXLocation, mouseYLocation, _, _ := GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)

	if eventStateMemory.stateId != constants.EventStateNone {
		return false
	}

	tooltipEntry := shared.getTooltipFromCharacterEntry(characterEntry)

	if tooltipEntry != nil {
		mouseXLocation, mouseYLocation, _, _ = GetMouseStatus()
		if tooltipEntry.HoverStartTime.IsZero() {
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
		if time.Since(tooltipEntry.HoverStartTime) >= time.Duration(tooltipEntry.HoverDisplayDelay)*time.Millisecond {
			setPreviouslyHighlightedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeTooltip)
			tooltipEntry.IsDrawn = true
			isScreenUpdateRequired = true
		}
	} else {
		for _, currentTooltipEntry := range Tooltips.GetAllEntriesOverall() {
			currentTooltipEntry.IsDrawn = false
			currentTooltipEntry.HoverStartTime = time.Time{}
		}
		if eventStateMemory.previouslyHighlightedControl.controlType == constants.CellTypeTooltip {
			setPreviouslyHighlightedControl("", "", constants.NullControlType)
		}
		isScreenUpdateRequired = true
	}
	return isScreenUpdateRequired
}

func (shared *TooltipInstanceType) setParentControlAlias(parentControlAlias string) *TooltipInstanceType {
	tooltipEntry := Tooltips.Get(shared.layerAlias, shared.controlAlias)
	if tooltipEntry != nil {
		tooltipEntry.ParentControlAlias = parentControlAlias
	}
	return shared
}
