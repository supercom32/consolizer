package consolizer

import (
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"time"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

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
SetEnabled is a method which allows you to enable or disable a tooltip.

:param enabled: Whether the tooltip should be enabled.

:return: The tooltip instance for method chaining.

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
Add is a method which allows you to create and add a new tooltip to a layer.

:param layerAlias: The alias of the layer to add the tooltip to.
:param tooltipAlias: The alias of the tooltip to create.
:param tooltipText: The text to display in the tooltip.
:param styleEntry: The TUI style entry to use for the tooltip.
:param hotspotXLocation: The X coordinate of the hotspot area.
:param hotspotYLocation: The Y coordinate of the hotspot area.
:param hotspotWidth: The width of the hotspot area.
:param hotspotHeight: The height of the hotspot area.
:param tooltipXLocation: The X coordinate of the tooltip (absolute or relative).
:param tooltipYLocation: The Y coordinate of the tooltip (absolute or relative).
:param tooltipWidth: The width of the tooltip.
:param tooltipHeight: The height of the tooltip.
:param isLocationAbsolute: Whether the tooltip coordinates are absolute.
:param isBorderDrawn: Whether to draw a border around the tooltip.
:param hoverTime: The delay in milliseconds before the tooltip appears.

:return: An instance of the created tooltip.

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
Delete is a method which allows you to remove a tooltip from a text layer. In addition, the following should be
noted:

- If you attempt to delete a tooltip which does not exist, then the request will simply be ignored.

:param layerAlias: The alias of the layer containing the tooltip.
:param labelAlias: The alias of the tooltip to delete.

Example:

	tooltip.Delete("layer1", "tooltip1")
*/
func (shared *tooltipType) Delete(layerAlias string, labelAlias string) {
	Tooltips.Remove(layerAlias, labelAlias)
}

/*
DeleteAll is a method which allows you to remove all tooltips from a layer.

:param layerAlias: The alias of the layer to remove all tooltips from.

Example:

	tooltip.DeleteAll("layer1")
*/
func (shared *tooltipType) DeleteAll(layerAlias string) {
	Tooltips.RemoveAll(layerAlias)
}

/*
drawHotspotZonesOnLayer is a method which allows you to draw all tooltip hotspot zones on a given text layer.

:param layerEntry: The entry for the layer to draw hotspots on.

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
drawHotspot is a method which allows you to draw a single tooltip hotspot zone. In addition, the following should
be noted:

- If a parent exists, do not overwrite the parent's attributes.

:param layerEntry: The entry for the layer to draw the hotspot on.
:param tooltipEntry: The entry for the tooltip whose hotspot is being drawn.

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
renderAll is a method which allows you to render all tooltips on a given text layer.

:param layerEntry: The entry for the layer to render tooltips on.

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
render is a method which allows you to render a tooltip on a given text layer. In addition, the following should
be noted:

- This method handles both absolute and relative positioning based on the tooltip configuration.

  - Rendering always starts at the coordinates specified by the user. However, the dimensions are always for the text
    area.

- If the tooltip is not enabled or not marked as drawn, then no rendering will occur.

- When absolute positioning is not used, the tooltip will be positioned relative to the current mouse cursor location.

- If borders are enabled, they will be drawn around the text area, expanding the total rendered size by 2 characters in.

:param layerEntry: The entry for the layer to render the tooltip on.
:param tooltipEntry: The entry for the tooltip to render.

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
getFromCharacterEntry is a method which allows you to retrieve a tooltip entry associated with a given character
entry. In addition, the following should be noted:

- It checks various control types (buttons, labels, checkboxes, etc.) for associated tooltips.

:param entry: The character entry to check for an associated tooltip.

:return: The associated tooltip entry, or nil if none found.

Example:

	entry := tooltip.getFromCharacterEntry(charEntry)
*/
func (shared *tooltipType) getFromCharacterEntry(entry types.CharacterEntryType) *types.TooltipEntryType {
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

/*
updateMouseEvent is a method which allows you to process mouse events for tooltips. In addition, the following should be
noted:

- Handles hover detection and timing.

- Manages showing and hiding of tooltips.

:return: True if a screen update is required.

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
setParentControlAlias is a method which allows you to associate a tooltip with a parent control.

:param parentControlAlias: The alias of the parent control.

:return: The tooltip instance for method chaining.

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
SetValue is a method which allows you to set the value of the tooltip associated with the TooltipInstanceType. In
addition, the following should be noted:

- This function updates the value of the tooltip label identified by the layerAlias and tooltipAlias fields.

:param text: The text to set for the tooltip.

:return: The tooltip instance for method chaining.

Example:

	tooltip.SetValue("New value")
*/
func (shared *TooltipInstanceType) SetValue(text string) *TooltipInstanceType {
	labelEntry := Labels.Get(shared.layerAlias, shared.controlAlias)
	labelEntry.Label = text
	return shared
}

/*
SetText is a method which allows you to set the text of the tooltip. In addition, the following should be noted:

- This is an alias for SetTooltipValue for consistency with other controls.

:param text: The text to set for the tooltip.

:return: The tooltip instance for method chaining.

Example:

	tooltip.SetText("New text")
*/
func (shared *TooltipInstanceType) SetText(text string) *TooltipInstanceType {
	return shared.SetValue(text)
}
