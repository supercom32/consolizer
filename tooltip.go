package consolizer

import (
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"time"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

/*
TooltipInstanceType is a structure which represents a specific instance of a tooltip control.
*/
type TooltipInstanceType struct {
	BaseControlInstanceType
}

type tooltipType struct{}

var Tooltip tooltipType
var Tooltips = memory.NewControlMemoryManager[types.TooltipEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
SetEnabled is a method which enables or disables a tooltip instance.

Example:
    tooltip.SetEnabled(true)
*/
func (shared *TooltipInstanceType) SetEnabled(enabled bool) *TooltipInstanceType {
	tooltipEntry := Tooltips.Get(shared.layerAlias, shared.controlAlias)
	if tooltipEntry != nil {
		tooltipEntry.IsEnabled = enabled
	}
	return shared
}

/*
Add is a method which creates and adds a new tooltip to a given text layer. In addition, the following should be noted:

- The tooltip is defined by a hotspot area and a separate display area for the tooltip text.

Example:
    tooltip.Add("layer1", "tooltip1", "Info", style, 10, 10, 5, 1, 10, 12, 10, 3, true, true, 500)
*/
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
	tooltipInstance.controlType = constants.TYPE_TOOLTIP
	return tooltipInstance
}

/*
Delete is a method which removes a tooltip from a text layer. In addition, the following should be noted:

- If you attempt to delete a tooltip which does not exist, then the request will simply be ignored.

Example:
    tooltip.Delete("layer1", "tooltip1")
*/
func (shared *tooltipType) Delete(layerAlias string, labelAlias string) {
	Tooltips.Remove(layerAlias, labelAlias)
	clearStaleControlReferences(layerAlias, labelAlias, constants.CellTypeTooltip)
}

/*
DeleteAll is a method which removes all tooltips from a specified layer.

Example:
    tooltip.DeleteAll("layer1")
*/
func (shared *tooltipType) DeleteAll(layerAlias string) {
	removeAllControlsAndClearCells(Tooltips, layerAlias, constants.CellTypeTooltip, func(entry *types.TooltipEntryType) string { return entry.Alias })
}

/*
drawHotspotZonesOnLayer is a method which draws all tooltip hotspot zones on a given text layer.

Example:
    tooltip.drawHotspotZonesOnLayer(layerEntry)
*/
func (shared *tooltipType) drawHotspotZonesOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentTooltipEntry := range Tooltips.GetAllEntries(layerAlias) {
		tooltipEntry := currentTooltipEntry
		shared.drawHotspot(&layerEntry, tooltipEntry)
	}
}

/*
drawHotspot is a method which draws a single tooltip hotspot zone. In addition, the following should be noted:

- If a parent exists, do not overwrite the parent's attributes.

Example:
    tooltip.drawHotspot(layerEntry, entry)
*/
func (shared *tooltipType) drawHotspot(layerEntry *types.LayerEntryType, tooltipEntry *types.TooltipEntryType) {
	if !tooltipEntry.IsEnabled {
		return
	}
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = tooltipEntry.StyleEntry.Tooltip.ForegroundColor
	attributeEntry.BackgroundColor = tooltipEntry.StyleEntry.Tooltip.BackgroundColor
	attributeEntry.CellType = constants.CellTypeTooltip
	attributeEntry.CellControlAlias = tooltipEntry.Alias
	if tooltipEntry.ParentControlAlias == "" { // If a parent exists, do not overwrite the parent's attributes.
		fillAreaWithControlAlias(layerEntry, attributeEntry.CellType, attributeEntry.CellControlAlias, tooltipEntry.HotspotXLocation, tooltipEntry.HotspotYLocation, tooltipEntry.HotspotWidth, tooltipEntry.HotspotHeight, constants.NullCellControlLocation)
	}
}

/*
renderAll is a method which renders all tooltips on a given text layer.

Example:
    tooltip.renderAll(layerEntry)
*/
func (shared *tooltipType) renderAll(layerEntry types.LayerEntryType) {
	for _, currentTooltipEntry := range Tooltips.GetAllEntriesOverall() {
		tooltipEntry := currentTooltipEntry
		shared.render(&layerEntry, tooltipEntry)
	}
}

/*
render is a method which renders a tooltip on a given text layer. In addition, the following should be noted:

- This method handles both absolute and relative positioning based on the tooltip configuration.

- Rendering always starts at the coordinates specified by the user. However, the dimensions are always for the text
  area.

- If the tooltip is not enabled or not marked as drawn, then no rendering will occur.

- When absolute positioning is not used, the tooltip will be positioned relative to the current mouse cursor location.

- If borders are enabled, they will be drawn around the text area, expanding the total rendered size by 2 characters.

Example:
    tooltip.render(layerEntry, entry)
*/
func (shared *tooltipType) render(layerEntry *types.LayerEntryType, tooltipEntry *types.TooltipEntryType) {
	if !tooltipEntry.IsEnabled || !tooltipEntry.IsDrawn {
		return
	}
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = tooltipEntry.StyleEntry.Tooltip.ForegroundColor
	attributeEntry.BackgroundColor = tooltipEntry.StyleEntry.Tooltip.BackgroundColor
	attributeEntry.CellType = constants.CellTypeTooltip
	attributeEntry.CellControlAlias = tooltipEntry.Alias
	calculatedXLocation := tooltipEntry.TooltipXLocation
	calculatedYLocation := tooltipEntry.TooltipYLocation
	calculatedWidth := tooltipEntry.TooltipWidth
	calculatedHeight := tooltipEntry.TooltipHeight
	if !tooltipEntry.IsLocationAbsolute {
		mouseXLocation, mouseYLocation, _, _ := GetMouseStatus()
		calculatedXLocation = mouseXLocation + tooltipEntry.TooltipXLocation
		calculatedYLocation = mouseYLocation + tooltipEntry.TooltipYLocation
		calculatedWidth = tooltipEntry.TooltipWidth
		calculatedHeight = tooltipEntry.TooltipHeight
	}
	fillStartX := calculatedXLocation
	fillStartY := calculatedYLocation
	fillWidth := calculatedWidth
	fillHeight := calculatedHeight
	xOffset := 2
	yOffset := 0
	// If a height of one, do not add white padding before text.
	if calculatedHeight != 1 {
		yOffset = 1
	}
	if tooltipEntry.IsBorderDrawn {
		calculatedWidth += 2
		calculatedHeight += 2
		xOffset = 2
		yOffset = 0
		fillWidth = calculatedWidth - 2
		fillHeight = calculatedHeight - 2
		fillStartX += 1
		fillStartY += 1
	}
	fillArea(layerEntry, attributeEntry, " ", fillStartX, fillStartY, fillWidth, fillHeight, constants.NullCellControlLocation)
	if tooltipEntry.IsBorderDrawn {
		drawBorder(layerEntry, tooltipEntry.StyleEntry, attributeEntry, calculatedXLocation, calculatedYLocation, calculatedWidth, calculatedHeight, false)
	}
	formattedLabel := " " + tooltipEntry.Text + " "
	arrayOfRunes := stringformat.GetRunesFromString(formattedLabel)
	layer.printWithWordWrap(layerEntry, attributeEntry, fillStartX+xOffset, fillStartY+yOffset, fillWidth-1, arrayOfRunes)
}

/*
getFromCharacterEntry is a method which retrieves a tooltip entry associated with a given character entry. In addition,
the following should be noted:

- It checks various control types (buttons, labels, checkboxes, etc.) for associated tooltips.

Example:
    entry := tooltip.getFromCharacterEntry(charEntry)
*/
func (shared *tooltipType) getFromCharacterEntry(entry types.CharacterEntryType) *types.TooltipEntryType {
	layer := entry.LayerAlias
	alias := entry.AttributeEntry.CellControlAlias

	switch entry.AttributeEntry.CellType {
	case constants.CellTypeButton:
		if button, isFound := Buttons.Lookup(layer, alias); isFound && button.TooltipAlias != "" {
			if tooltipEntry, isFound := Tooltips.Lookup(layer, button.TooltipAlias); isFound {
				return tooltipEntry
			}
		}
	case constants.CellTypeLabel:
		if label, isFound := Labels.Lookup(layer, alias); isFound && label.TooltipAlias != "" {
			if tooltipEntry, isFound := Tooltips.Lookup(layer, label.TooltipAlias); isFound {
				return tooltipEntry
			}
		}
	case constants.CellTypeCheckbox:
		if checkbox, isFound := Checkboxes.Lookup(layer, alias); isFound && checkbox.TooltipAlias != "" {
			if tooltipEntry, isFound := Tooltips.Lookup(layer, checkbox.TooltipAlias); isFound {
				return tooltipEntry
			}
		}
	case constants.CellTypeRadioButton:
		if radio, isFound := RadioButtons.Lookup(layer, alias); isFound && radio.TooltipAlias != "" {
			if tooltipEntry, isFound := Tooltips.Lookup(layer, radio.TooltipAlias); isFound {
				return tooltipEntry
			}
		}
	case constants.CellTypeTextField:
		if textField, isFound := TextFields.Lookup(layer, alias); isFound && textField.TooltipAlias != "" {
			if tooltipEntry, isFound := Tooltips.Lookup(layer, textField.TooltipAlias); isFound {
				return tooltipEntry
			}
		}
	case constants.CellTypeTextbox:
		if textbox, isFound := Textboxes.Lookup(layer, alias); isFound && textbox.TooltipAlias != "" {
			if tooltipEntry, isFound := Tooltips.Lookup(layer, textbox.TooltipAlias); isFound {
				return tooltipEntry
			}
		}
	case constants.CellTypeProgressBar:
		if progressBar, isFound := ProgressBars.Lookup(layer, alias); isFound && progressBar.TooltipAlias != "" {
			if tooltipEntry, isFound := Tooltips.Lookup(layer, progressBar.TooltipAlias); isFound {
				return tooltipEntry
			}
		}
	case constants.CellTypeSelectorItem:
		if selector, isFound := Selectors.Lookup(layer, alias); isFound && selector.TooltipAlias != "" {
			if tooltipEntry, isFound := Tooltips.Lookup(layer, selector.TooltipAlias); isFound {
				return tooltipEntry
			}
		}
	case constants.CellTypeTooltip:
		if tooltipEntry, isFound := Tooltips.Lookup(layer, alias); isFound {
			return tooltipEntry
		}
	}
	return nil
}

/*
updateMouseEvent is a method which processes mouse events for tooltips. In addition, the following should be noted:

- Handles hover detection and timing.

- Manages showing and hiding of tooltips based on mouse movement and position.

Example:
    update := tooltip.updateMouseEvent()
*/
func (shared *tooltipType) updateMouseEvent() bool {
	isScreenUpdateRequired := false
	mouseXLocation, mouseYLocation, _, _ := GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)

	if eventStateMemory.stateId != constants.EventStateNone {
		return false
	}

	tooltipEntry := shared.getFromCharacterEntry(characterEntry)
	if tooltipEntry != nil && !tooltipEntry.IsEnabled {
		tooltipEntry = nil
	}

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
			if !tooltipEntry.IsDrawn {
				isScreenUpdateRequired = true
			}
			tooltipEntry.IsDrawn = true
		}
	} else {
		for _, currentTooltipEntry := range Tooltips.GetAllEntriesOverall() {
			if currentTooltipEntry.IsDrawn == true {
				// Only update if a change was detected.
				isScreenUpdateRequired = true
			}
			currentTooltipEntry.IsDrawn = false
			currentTooltipEntry.HoverStartTime = time.Time{}
		}
		if eventStateMemory.previouslyHighlightedControl.controlType == constants.CellTypeTooltip {
			setPreviouslyHighlightedControl("", "", constants.NullControlType)
		}
	}
	return isScreenUpdateRequired
}

/*
setParentControlAlias is a method which associates a tooltip with a parent control.

Example:
    tooltip.setParentControlAlias("parent1")
*/
func (shared *TooltipInstanceType) setParentControlAlias(parentControlAlias string) *TooltipInstanceType {
	tooltipEntry := Tooltips.Get(shared.layerAlias, shared.controlAlias)
	if tooltipEntry != nil {
		tooltipEntry.ParentControlAlias = parentControlAlias
	}
	return shared
}

/*
SetValue is a method which sets the value of the tooltip associated with the TooltipInstanceType. In addition, the
following should be noted:

- This function updates the value of the tooltip label identified by the layerAlias and tooltipAlias fields.

Example:
    tooltip.SetValue("New value")
*/
func (shared *TooltipInstanceType) SetValue(text string) *TooltipInstanceType {
	tooltipEntry := Tooltips.Get(shared.layerAlias, shared.controlAlias)
	tooltipEntry.Text = text
	return shared
}

/*
SetText is a method which sets the text of the tooltip. In addition, the following should be noted:

- This is an alias for SetTooltipValue for consistency with other controls.

Example:
    tooltip.SetText("New text")
*/
func (shared *TooltipInstanceType) SetText(text string) *TooltipInstanceType {
	return shared.SetValue(text)
}
