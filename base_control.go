package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

// BaseControlInstanceType provides common functionality for all control instances.
// It manages the layer and control aliases and provides common methods that all controls share.
type BaseControlInstanceType struct {
	layerAlias   string
	controlAlias string
	controlType  string
}

/*
GetAlias is a method which allows you to obtain the alias associated with the control.

:return: The unique alias string for the control.

Example:

	alias := control.GetAlias()
*/
func (shared *BaseControlInstanceType) GetAlias() string {
	return shared.controlAlias
}

/*
GetLayerAlias is a method which allows you to obtain the alias of the layer that the control is associated with.

:return: The alias string of the associated layer.

Example:

	layerAlias := control.GetLayerAlias()
*/
func (shared *BaseControlInstanceType) GetLayerAlias() string {
	return shared.layerAlias
}

/*
getBaseControl is a method which allows you to retrieve the underlying base control type for the control instance.

:return: A pointer to the BaseControlType structure if found, otherwise nil.

Example:

	baseControl := shared.getBaseControl()
*/
func (shared *BaseControlInstanceType) getBaseControl() *types.BaseControlType {
	switch shared.controlType {
	case constants.TYPE_BUTTON:
		if entry := Buttons.Get(shared.layerAlias, shared.controlAlias); entry != nil {
			return &entry.BaseControlType
		}
	case constants.TYPE_CHECKBOX:
		if entry := Checkboxes.Get(shared.layerAlias, shared.controlAlias); entry != nil {
			return &entry.BaseControlType
		}
	case constants.TYPE_DROPDOWN:
		if entry := Dropdowns.Get(shared.layerAlias, shared.controlAlias); entry != nil {
			return &entry.BaseControlType
		}
	case constants.TYPE_LABEL:
		if entry := Labels.Get(shared.layerAlias, shared.controlAlias); entry != nil {
			return &entry.BaseControlType
		}
	case constants.TYPE_PROGRESSBAR:
		if entry := ProgressBars.Get(shared.layerAlias, shared.controlAlias); entry != nil {
			return &entry.BaseControlType
		}
	case constants.TYPE_SCROLLBAR:
		if entry := ScrollBars.Get(shared.layerAlias, shared.controlAlias); entry != nil {
			return &entry.BaseControlType
		}
	case constants.TYPE_SELECTOR:
		if entry := Selectors.Get(shared.layerAlias, shared.controlAlias); entry != nil {
			return &entry.BaseControlType
		}
	case constants.TYPE_TEXTBOX:
		if entry := Textboxes.Get(shared.layerAlias, shared.controlAlias); entry != nil {
			return &entry.BaseControlType
		}
	case constants.TYPE_TEXTFIELD:
		if entry := TextFields.Get(shared.layerAlias, shared.controlAlias); entry != nil {
			return &entry.BaseControlType
		}
	case constants.TYPE_TOOLTIP:
		if entry := Tooltips.Get(shared.layerAlias, shared.controlAlias); entry != nil {
			return &entry.BaseControlType
		}
	}
	return nil
}

/*
GetBounds is a method which allows you to obtain the position and dimensions of the control.

:return: The X location, Y location, width, and height of the control.

Example:

	x, y, w, h := control.GetBounds()
*/
func (shared *BaseControlInstanceType) GetBounds() (int, int, int, int) {
	if control := shared.getBaseControl(); control != nil {
		return control.XLocation, control.YLocation, control.Width, control.Height
	}
	return 0, 0, 0, 0
}

/*
SetPosition is a method which allows you to set the X and Y coordinates of the control.

:param x: The new X coordinate for the control.
:param y: The new Y coordinate for the control.

:return: The current BaseControlInstanceType instance for method chaining.

Example:

	control.SetPosition(10, 5)
*/
func (shared *BaseControlInstanceType) SetPosition(x, y int) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.XLocation = x
		control.YLocation = y
	}
	return shared
}

/*
GetPosition is a method which allows you to retrieve the current X and Y coordinates of the control.

:return: The current X and Y position of the control.

Example:

	x, y := control.GetPosition()
*/
func (shared *BaseControlInstanceType) GetPosition() (int, int) {
	if control := shared.getBaseControl(); control != nil {
		return control.XLocation, control.YLocation
	}
	return 0, 0
}

/*
GetSize is a method which allows you to retrieve the current width and height of the control.

:return: The current width and height of the control.

Example:

	width, height := control.GetSize()
*/
func (shared *BaseControlInstanceType) GetSize() (int, int) {
	if control := shared.getBaseControl(); control != nil {
		return control.Width, control.Height
	}
	return 0, 0
}

/*
SetSize is a method which allows you to set the width and height of the control.

:param width: The new width for the control.
:param height: The new height for the control.

:return: The current BaseControlInstanceType instance for method chaining.

Example:

	control.SetSize(20, 10)
*/
func (shared *BaseControlInstanceType) SetSize(width, height int) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.Width = width
		control.Height = height
	}
	return shared
}

/*
SetVisible is a method which allows you to toggle the visibility of the control.

:param visible: Set to true to make the control visible, or false to hide it.

:return: The current BaseControlInstanceType instance for method chaining.

Example:

	control.SetVisible(true)
*/
func (shared *BaseControlInstanceType) SetVisible(visible bool) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.IsVisible = visible
	}
	return shared
}

/*
SetStyle is a method which allows you to apply a visual style to the control.

:param style: The TuiStyleEntryType structure defining the new visual style.

:return: The current BaseControlInstanceType instance for method chaining.

Example:

	control.SetStyle(newStyle)
*/
func (shared *BaseControlInstanceType) SetStyle(style types.TuiStyleEntryType) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.StyleEntry = style
	}
	return shared
}

/*
SetEnabled is a method which allows you to enable or disable user interaction with the control.

:param enabled: Set to true to enable the control, or false to disable it.

:return: The current BaseControlInstanceType instance for method chaining.

Example:

	control.SetEnabled(false)
*/
func (shared *BaseControlInstanceType) SetEnabled(enabled bool) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.IsEnabled = enabled
	}
	return shared
}

/*
SetLabel is a method which allows you to set the display text for the control's label.

:param label: The new label text for the control.

:return: The current BaseControlInstanceType instance for method chaining.

Example:

	control.SetLabel("Click Me")
*/
func (shared *BaseControlInstanceType) SetLabel(label string) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.Label = label
	}
	return shared
}

/*
GetLabel is a method which allows you to retrieve the current label text of the control.

:return: The current label text of the control.

Example:

	labelText := control.GetLabel()
*/
func (shared *BaseControlInstanceType) GetLabel() string {
	if control := shared.getBaseControl(); control != nil {
		return control.Label
	}
	return ""
}

/*
SetBorderDrawn is a method which allows you to specify whether a border should be drawn around the control.

:param drawn: Set to true to draw a border, or false to omit it.

:return: The current BaseControlInstanceType instance for method chaining.

Example:

	control.SetBorderDrawn(true)
*/
func (shared *BaseControlInstanceType) SetBorderDrawn(drawn bool) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.IsBorderDrawn = drawn
	}
	return shared
}

/*
IsBorderDrawn is a method which allows you to check if a border is currently being drawn around the control.

:return: True if a border is drawn, otherwise false.

Example:

	isDrawn := control.IsBorderDrawn()
*/
func (shared *BaseControlInstanceType) IsBorderDrawn() bool {
	if control := shared.getBaseControl(); control != nil {
		return control.IsBorderDrawn
	}
	return false
}

/*
SetTooltip is a method which allows you to associate a tooltip with the control.

:param tooltipAlias: The alias string of the tooltip to be associated with this control.

:return: The current BaseControlInstanceType instance for method chaining.

Example:

	control.SetTooltip("myTooltip")
*/
func (shared *BaseControlInstanceType) SetTooltip(tooltipAlias string) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.TooltipAlias = tooltipAlias
	}
	return shared
}

/*
GetTooltip is a method which allows you to retrieve the alias of the tooltip associated with the control.

:return: The alias string of the associated tooltip.

Example:

	tooltip := control.GetTooltip()
*/
func (shared *BaseControlInstanceType) GetTooltip() string {
	if control := shared.getBaseControl(); control != nil {
		return control.TooltipAlias
	}
	return ""
}

/*
SetTooltipEnabled is a method which allows you to enable or disable the display of the control's tooltip.

:param enabled: Set to true to enable the tooltip, or false to disable it.

:return: The current BaseControlInstanceType instance for method chaining.

Example:

	control.SetTooltipEnabled(true)
*/
func (shared *BaseControlInstanceType) SetTooltipEnabled(enabled bool) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.IsTooltipEnabled = enabled
	}
	return shared
}

/*
IsTooltipEnabled is a method which allows you to check if the tooltip for the control is currently enabled.

:return: True if the tooltip is enabled, otherwise false.

Example:

	isEnabled := control.IsTooltipEnabled()
*/
func (shared *BaseControlInstanceType) IsTooltipEnabled() bool {
	if control := shared.getBaseControl(); control != nil {
		return control.IsTooltipEnabled
	}
	return false
}

/*
IsVisible is a method which allows you to check if the control is currently set to be visible.

:return: True if the control is visible, otherwise false.

Example:

	isVisible := control.IsVisible()
*/
func (shared *BaseControlInstanceType) IsVisible() bool {
	if control := shared.getBaseControl(); control != nil {
		return control.IsVisible
	}
	return false
}

/*
IsEnabled is a method which allows you to check if the control is currently enabled for user interaction.

:return: True if the control is enabled, otherwise false.

Example:

	isEnabled := control.IsEnabled()
*/
func (shared *BaseControlInstanceType) IsEnabled() bool {
	if control := shared.getBaseControl(); control != nil {
		return control.IsEnabled
	}
	return false
}

/*
GetStyle is a method which allows you to retrieve the current visual style of the control.

:return: The current TuiStyleEntryType structure representing the control's style.

Example:

	style := control.GetStyle()
*/
func (shared *BaseControlInstanceType) GetStyle() types.TuiStyleEntryType {
	if control := shared.getBaseControl(); control != nil {
		return control.StyleEntry
	}
	return types.NewTuiStyleEntry()
}

/*
Delete is a method which allows you to remove a control from its memory manager. In addition, the following should be
noted:

- If you attempt to delete a control which does not exist, then the request will simply be ignored.

- All memory associated with the control will be freed.

:return: Nil value after the control has been removed.

Example:

	control.Delete()
*/
func (shared *BaseControlInstanceType) Delete() *BaseControlInstanceType {
	switch shared.controlType {
	case constants.TYPE_BUTTON:
		if Buttons.IsExists(shared.layerAlias, shared.controlAlias) {
			Buttons.Remove(shared.layerAlias, shared.controlAlias)
		}
	case constants.TYPE_CHECKBOX:
		if Checkboxes.IsExists(shared.layerAlias, shared.controlAlias) {
			Checkboxes.Remove(shared.layerAlias, shared.controlAlias)
		}
	case constants.TYPE_DROPDOWN:
		if Dropdowns.IsExists(shared.layerAlias, shared.controlAlias) {
			Dropdowns.Remove(shared.layerAlias, shared.controlAlias)
		}
	case constants.TYPE_LABEL:
		if Labels.IsExists(shared.layerAlias, shared.controlAlias) {
			Labels.Remove(shared.layerAlias, shared.controlAlias)
		}
	case constants.TYPE_PROGRESSBAR:
		if ProgressBars.IsExists(shared.layerAlias, shared.controlAlias) {
			ProgressBars.Remove(shared.layerAlias, shared.controlAlias)
		}
	case constants.TYPE_SCROLLBAR:
		if ScrollBars.IsExists(shared.layerAlias, shared.controlAlias) {
			ScrollBars.Remove(shared.layerAlias, shared.controlAlias)
		}
	case constants.TYPE_SELECTOR:
		if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
			Selectors.Remove(shared.layerAlias, shared.controlAlias)
		}
	case constants.TYPE_TEXTBOX:
		if Textboxes.IsExists(shared.layerAlias, shared.controlAlias) {
			Textboxes.Remove(shared.layerAlias, shared.controlAlias)
		}
	case constants.TYPE_TEXTFIELD:
		if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
			TextFields.Remove(shared.layerAlias, shared.controlAlias)
		}
	case constants.TYPE_TOOLTIP:
		if Tooltips.IsExists(shared.layerAlias, shared.controlAlias) {
			Tooltips.Remove(shared.layerAlias, shared.controlAlias)
		}
	case constants.TYPE_RADIOBUTTON:
		if RadioButtons.IsExists(shared.layerAlias, shared.controlAlias) {
			RadioButtons.Remove(shared.layerAlias, shared.controlAlias)
		}
	case constants.TYPE_VIEWPORT:
		if Viewports.IsExists(shared.layerAlias, shared.controlAlias) {
			Viewports.Remove(shared.layerAlias, shared.controlAlias)
		}
	}
	return nil
}

/*
GetFocus is a method which allows you to update the event manager to set this control as the one currently in focus.

:return: The current BaseControlInstanceType instance for method chaining.

Example:

	control.GetFocus()
*/
func (shared *BaseControlInstanceType) GetFocus() *BaseControlInstanceType {
	controlTypeInt := constants.NullControlType
	switch shared.controlType {
	case constants.TYPE_BUTTON:
		controlTypeInt = constants.CellTypeButton
	case constants.TYPE_CHECKBOX:
		controlTypeInt = constants.CellTypeCheckbox
	case constants.TYPE_DROPDOWN:
		controlTypeInt = constants.CellTypeDropdown
	case constants.TYPE_LABEL:
		controlTypeInt = constants.CellTypeLabel
	case constants.TYPE_PROGRESSBAR:
		controlTypeInt = constants.CellTypeProgressBar
	case constants.TYPE_SCROLLBAR:
		controlTypeInt = constants.CellTypeScrollbar
	case constants.TYPE_SELECTOR:
		controlTypeInt = constants.CellTypeSelectorItem
	case constants.TYPE_TEXTBOX:
		controlTypeInt = constants.CellTypeTextbox
	case constants.TYPE_TEXTFIELD:
		controlTypeInt = constants.CellTypeTextField
	case constants.TYPE_TOOLTIP:
		controlTypeInt = constants.CellTypeTooltip
	case constants.TYPE_RADIOBUTTON:
		controlTypeInt = constants.CellTypeRadioButton
	}
	setFocusedControl(shared.layerAlias, shared.controlAlias, controlTypeInt)
	return shared
}
