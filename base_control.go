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

// GetAlias returns the alias of the control
func (shared *BaseControlInstanceType) GetAlias() string {
	return shared.controlAlias
}

// GetLayerAlias returns the layer alias of the control
func (shared *BaseControlInstanceType) GetLayerAlias() string {
	return shared.layerAlias
}

// GetBaseControl returns the BaseControlType for the control
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

// GetBounds returns the position and size of the control
func (shared *BaseControlInstanceType) GetBounds() (int, int, int, int) {
	if control := shared.getBaseControl(); control != nil {
		return control.XLocation, control.YLocation, control.Width, control.Height
	}
	return 0, 0, 0, 0
}

// SetPosition sets the position of the control
func (shared *BaseControlInstanceType) SetPosition(x, y int) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.XLocation = x
		control.YLocation = y
	}
	return shared
}

// GetPosition returns the X and Y position of the control
func (shared *BaseControlInstanceType) GetPosition() (int, int) {
	if control := shared.getBaseControl(); control != nil {
		return control.XLocation, control.YLocation
	}
	return 0, 0
}

// GetSize returns the width and height of the control
func (shared *BaseControlInstanceType) GetSize() (int, int) {
	if control := shared.getBaseControl(); control != nil {
		return control.Width, control.Height
	}
	return 0, 0
}

// SetSize sets the dimensions of the control
func (shared *BaseControlInstanceType) SetSize(width, height int) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.Width = width
		control.Height = height
	}
	return shared
}

// SetVisible shows or hides the control
func (shared *BaseControlInstanceType) SetVisible(visible bool) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.IsVisible = visible
	}
	return shared
}

// SetStyle sets the visual style of the control
func (shared *BaseControlInstanceType) SetStyle(style types.TuiStyleEntryType) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.StyleEntry = style
	}
	return shared
}

// SetEnabled enables or disables the control
func (shared *BaseControlInstanceType) SetEnabled(enabled bool) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.IsEnabled = enabled
	}
	return shared
}

// SetLabel sets the label text of the control
func (shared *BaseControlInstanceType) SetLabel(label string) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.Label = label
	}
	return shared
}

// GetLabel gets the label text of the control
func (shared *BaseControlInstanceType) GetLabel() string {
	if control := shared.getBaseControl(); control != nil {
		return control.Label
	}
	return ""
}

// SetBorderDrawn controls whether a border is drawn around the control
func (shared *BaseControlInstanceType) SetBorderDrawn(drawn bool) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.IsBorderDrawn = drawn
	}
	return shared
}

// IsBorderDrawn returns whether a border is drawn around the control
func (shared *BaseControlInstanceType) IsBorderDrawn() bool {
	if control := shared.getBaseControl(); control != nil {
		return control.IsBorderDrawn
	}
	return false
}

// SetTooltip sets the tooltip alias for the control
func (shared *BaseControlInstanceType) SetTooltip(tooltipAlias string) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.TooltipAlias = tooltipAlias
	}
	return shared
}

// GetTooltip gets the tooltip alias for the control
func (shared *BaseControlInstanceType) GetTooltip() string {
	if control := shared.getBaseControl(); control != nil {
		return control.TooltipAlias
	}
	return ""
}

// SetTooltipEnabled enables or disables the tooltip for the control
func (shared *BaseControlInstanceType) SetTooltipEnabled(enabled bool) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.IsTooltipEnabled = enabled
	}
	return shared
}

// IsTooltipEnabled returns whether the tooltip is enabled for the control
func (shared *BaseControlInstanceType) IsTooltipEnabled() bool {
	if control := shared.getBaseControl(); control != nil {
		return control.IsTooltipEnabled
	}
	return false
}

// IsVisible returns whether the control is visible
func (shared *BaseControlInstanceType) IsVisible() bool {
	if control := shared.getBaseControl(); control != nil {
		return control.IsVisible
	}
	return false
}

// IsEnabled returns whether the control is enabled
func (shared *BaseControlInstanceType) IsEnabled() bool {
	if control := shared.getBaseControl(); control != nil {
		return control.IsEnabled
	}
	return false
}

// GetStyle returns the visual style of the control
func (shared *BaseControlInstanceType) GetStyle() types.TuiStyleEntryType {
	if control := shared.getBaseControl(); control != nil {
		return control.StyleEntry
	}
	return types.NewTuiStyleEntry()
}

// Delete removes a control from its memory manager. In addition, the following
// information should be noted:
//
// - If you attempt to delete a control which does not exist, then the request
// will simply be ignored.
// - All memory associated with the control will be freed.
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
	}
	return nil
}

// GetFocus updates the event manager to make this control the one in focus
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
