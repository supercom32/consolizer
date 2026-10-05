package consolizer

import (
	"fmt"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/types"
)

/*
BaseControlInstanceType is a structure which provides common functionality for all control instances. It manages the
layer and control aliases and provides common methods that all controls share.

Example:
    var baseControl BaseControlInstanceType
*/
type BaseControlInstanceType struct {
	layerAlias   string
	controlAlias string
	controlType  string
}

/*
GetAlias is a method which allows you to obtain the alias associated with the control.

Example:
    alias := control.GetAlias()
*/
func (shared *BaseControlInstanceType) GetAlias() string {
	return shared.controlAlias
}

/*
GetLayerAlias is a method which allows you to obtain the alias of the layer that the control is associated with.

Example:
    layerAlias := control.GetLayerAlias()
*/
func (shared *BaseControlInstanceType) GetLayerAlias() string {
	return shared.layerAlias
}

/*
getBaseControl is a method which allows you to retrieve the underlying base control type for the control instance.

Example:
    baseControl := shared.getBaseControl()
*/
func (shared *BaseControlInstanceType) getBaseControl() *types.BaseControlType {
	switch shared.controlType {
	case constants.TYPE_BUTTON:
		if entry, isFound := Buttons.Lookup(shared.layerAlias, shared.controlAlias); isFound {
			return &entry.BaseControlType
		}
	case constants.TYPE_CHECKBOX:
		if entry, isFound := Checkboxes.Lookup(shared.layerAlias, shared.controlAlias); isFound {
			return &entry.BaseControlType
		}
	case constants.TYPE_DROPDOWN:
		if entry, isFound := Dropdowns.Lookup(shared.layerAlias, shared.controlAlias); isFound {
			return &entry.BaseControlType
		}
	case constants.TYPE_LABEL:
		if entry, isFound := Labels.Lookup(shared.layerAlias, shared.controlAlias); isFound {
			return &entry.BaseControlType
		}
	case constants.TYPE_PROGRESSBAR:
		if entry, isFound := ProgressBars.Lookup(shared.layerAlias, shared.controlAlias); isFound {
			return &entry.BaseControlType
		}
	case constants.TYPE_SCROLLBAR:
		if entry, isFound := ScrollBars.Lookup(shared.layerAlias, shared.controlAlias); isFound {
			return &entry.BaseControlType
		}
	case constants.TYPE_SELECTOR:
		if entry, isFound := Selectors.Lookup(shared.layerAlias, shared.controlAlias); isFound {
			return &entry.BaseControlType
		}
	case constants.TYPE_TEXTBOX:
		if entry, isFound := Textboxes.Lookup(shared.layerAlias, shared.controlAlias); isFound {
			return &entry.BaseControlType
		}
	case constants.TYPE_TEXTFIELD:
		if entry, isFound := TextFields.Lookup(shared.layerAlias, shared.controlAlias); isFound {
			return &entry.BaseControlType
		}
	case constants.TYPE_TOOLTIP:
		if entry, isFound := Tooltips.Lookup(shared.layerAlias, shared.controlAlias); isFound {
			return &entry.BaseControlType
		}
	case constants.TYPE_RADIOBUTTON:
		if entry, isFound := RadioButtons.Lookup(shared.layerAlias, shared.controlAlias); isFound {
			return &entry.BaseControlType
		}
	}
	return nil
}

/*
GetBounds is a method which allows you to obtain the position and dimensions of the control.

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
SetVisible is a method which allows you to toggle the visibility of the control. In addition, the following should be
noted:

  - Hiding the focused control moves focus to the next stop in its layer's tab order, or clears focus if no other
    stop can take it.

Example:
    control.SetVisible(true)
*/
func (shared *BaseControlInstanceType) SetVisible(visible bool) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.IsVisible = visible
	}
	reconcileFocus()
	return shared
}

/*
SetStyle is a method which allows you to apply a visual style to the control.

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
SetEnabled is a method which allows you to enable or disable user interaction with the control. In addition, the
following should be noted:

  - Disabling the focused control moves focus to the next stop in its layer's tab order, or clears focus if no other
    stop can take it.

Example:
    control.SetEnabled(false)
*/
func (shared *BaseControlInstanceType) SetEnabled(enabled bool) *BaseControlInstanceType {
	if control := shared.getBaseControl(); control != nil {
		control.IsEnabled = enabled
	}
	reconcileFocus()
	return shared
}

/*
SetLabel is a method which allows you to set the display text for the control's label.

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
Delete is a method which removes a control from its memory manager. In addition, the following should be
noted:

- If you attempt to delete a control which does not exist, then the request will simply be ignored.

- All memory associated with the control will be freed.

Example:
    control.Delete()
*/
func (shared *BaseControlInstanceType) Delete() *BaseControlInstanceType {
	switch shared.controlType {
	case constants.TYPE_BUTTON:
		if Buttons.IsExists(shared.layerAlias, shared.controlAlias) {
			Buttons.Remove(shared.layerAlias, shared.controlAlias)
			clearStaleControlReferences(shared.layerAlias, shared.controlAlias, constants.CellTypeButton)
		}
	case constants.TYPE_CHECKBOX:
		if Checkboxes.IsExists(shared.layerAlias, shared.controlAlias) {
			Checkboxes.Remove(shared.layerAlias, shared.controlAlias)
			clearStaleControlReferences(shared.layerAlias, shared.controlAlias, constants.CellTypeCheckbox)
		}
	case constants.TYPE_DROPDOWN:
		if Dropdowns.IsExists(shared.layerAlias, shared.controlAlias) {
			Dropdowns.Remove(shared.layerAlias, shared.controlAlias)
			clearStaleControlReferences(shared.layerAlias, shared.controlAlias, constants.CellTypeDropdown)
		}
	case constants.TYPE_LABEL:
		if Labels.IsExists(shared.layerAlias, shared.controlAlias) {
			Labels.Remove(shared.layerAlias, shared.controlAlias)
			clearStaleControlReferences(shared.layerAlias, shared.controlAlias, constants.CellTypeLabel)
		}
	case constants.TYPE_PROGRESSBAR:
		if ProgressBars.IsExists(shared.layerAlias, shared.controlAlias) {
			ProgressBars.Remove(shared.layerAlias, shared.controlAlias)
			clearStaleControlReferences(shared.layerAlias, shared.controlAlias, constants.CellTypeProgressBar)
		}
	case constants.TYPE_SCROLLBAR:
		if ScrollBars.IsExists(shared.layerAlias, shared.controlAlias) {
			ScrollBars.Remove(shared.layerAlias, shared.controlAlias)
			clearStaleControlReferences(shared.layerAlias, shared.controlAlias, constants.CellTypeScrollbar)
		}
	case constants.TYPE_SELECTOR:
		if Selectors.IsExists(shared.layerAlias, shared.controlAlias) {
			Selectors.Remove(shared.layerAlias, shared.controlAlias)
			clearStaleControlReferences(shared.layerAlias, shared.controlAlias, constants.CellTypeSelectorItem)
		}
	case constants.TYPE_TEXTBOX:
		if Textboxes.IsExists(shared.layerAlias, shared.controlAlias) {
			Textboxes.Remove(shared.layerAlias, shared.controlAlias)
			clearStaleControlReferences(shared.layerAlias, shared.controlAlias, constants.CellTypeTextbox)
		}
	case constants.TYPE_TEXTFIELD:
		if TextFields.IsExists(shared.layerAlias, shared.controlAlias) {
			TextFields.Remove(shared.layerAlias, shared.controlAlias)
			clearStaleControlReferences(shared.layerAlias, shared.controlAlias, constants.CellTypeTextField)
		}
	case constants.TYPE_TOOLTIP:
		if Tooltips.IsExists(shared.layerAlias, shared.controlAlias) {
			Tooltips.Remove(shared.layerAlias, shared.controlAlias)
			clearStaleControlReferences(shared.layerAlias, shared.controlAlias, constants.CellTypeTooltip)
		}
	case constants.TYPE_RADIOBUTTON:
		if RadioButtons.IsExists(shared.layerAlias, shared.controlAlias) {
			RadioButtons.Remove(shared.layerAlias, shared.controlAlias)
			clearStaleControlReferences(shared.layerAlias, shared.controlAlias, constants.CellTypeRadioButton)
		}
	case constants.TYPE_VIEWPORT:
		if Viewports.IsExists(shared.layerAlias, shared.controlAlias) {
			Viewports.Remove(shared.layerAlias, shared.controlAlias)
			clearStaleControlReferences(shared.layerAlias, shared.controlAlias, constants.CellTypeTextbox)
		}
	}
	return nil
}

/*
removeAllControlsAndClearCells is a method which allows you to remove every control of one type from a layer while
also resetting the screen's stale hit-testing metadata for each one, so a control created afterward cannot be
reached by a leftover reference to a deleted control that reused its alias. In addition, the following should be
noted:

  - getAlias is called once per surviving entry, before manager.RemoveAll runs, since the entries are no longer
    reachable through manager afterward. The cells and focus state are only cleared after the removal, so a
    deleted control that had focus passes it on to a control that still exists rather than to one about to be
    removed in the same call.

  - This is the shared implementation behind every control type's DeleteAll method and every DeleteAllX helper on
    LayerInstanceType.

Example:

	removeAllControlsAndClearCells(Buttons, "layer1", constants.CellTypeButton,
		func(entry *types.ButtonEntryType) string { return entry.Alias })
*/
func removeAllControlsAndClearCells[T any](manager *memory.ControlMemoryManager[T], layerAlias string, cellType int, getAlias func(*T) string) {
	var aliases []string
	for _, entry := range manager.GetAllEntries(layerAlias) {
		aliases = append(aliases, getAlias(entry))
	}
	manager.RemoveAll(layerAlias)
	for _, alias := range aliases {
		clearStaleControlReferences(layerAlias, alias, cellType)
	}
}

/*
getControlIdentifier is a method which allows you to obtain the identifier the focus manager uses for this control: its
layer alias, control alias, and cell type. It is what makes every control instance usable with SetFocus.

Example:

	control := button.getControlIdentifier()
*/
func (shared *BaseControlInstanceType) getControlIdentifier() controlIdentifierType {
	return controlIdentifierType{layerAlias: shared.layerAlias, controlAlias: shared.controlAlias, controlType: getFocusStopCellType(shared.controlType)}
}

/*
validateTabStop is a method which allows you to check that this control can be registered as a tab stop, returning an
error that names the control when it is not an interactive type or does not exist.

Example:

	err := control.validateTabStop()
*/
func (shared *BaseControlInstanceType) validateTabStop() error {
	control := shared.getControlIdentifier()
	if !isTabStopCellType(control.controlType) {
		return fmt.Errorf("the control '%s' on layer '%s' cannot be added to the tab order since it is not an interactive control", shared.controlAlias, shared.layerAlias)
	}
	if !isControlExists(control) {
		return fmt.Errorf("the control '%s' on layer '%s' cannot be added to the tab order since it does not exist", shared.controlAlias, shared.layerAlias)
	}
	return nil
}

/*
AddToTabIndex is a method which allows you to add the control to its layer's tab order, so that Tab and Shift+Tab move
focus to it. Stops are visited in the order they were added, except that stops given an explicit order with
SetTabIndex come first. An error is returned if the control is not an interactive type, such as a label or progress
bar, or does not exist. Adding a control that is already in the tab order does nothing. In addition, the following
should be noted:

  - Each layer has its own tab order, so Tab never moves focus from one layer to another.

  - Radio buttons of the same group count as a single stop, entered at the selected button. The arrow keys then move
    and select within the group.

  - While a control is disabled or hidden, or its layer is hidden, Tab skips it. A deleted control is removed from
    the tab order automatically.

Example:

	err := okButton.AddToTabIndex()
*/
func (shared *BaseControlInstanceType) AddToTabIndex() error {
	if err := shared.validateTabStop(); err != nil {
		return err
	}
	defer focusManager.beginChange()()
	focusManager.registerStopLocked(shared.getControlIdentifier(), noExplicitTabIndex)
	return nil
}

/*
SetTabIndex is a method which allows you to give the control an explicit position in its layer's tab order, adding it
to the tab order if it is not already there. Stops with an explicit position come before every other stop, in
ascending order of position, and stops with the same position keep the order they were added in. An error is returned
if the position is negative, or the control is not an interactive type or does not exist.

Example:

	err := cancelButton.SetTabIndex(2)
*/
func (shared *BaseControlInstanceType) SetTabIndex(tabIndex int) error {
	if tabIndex < 0 {
		return fmt.Errorf("the tab index for control '%s' on layer '%s' must not be negative, but %d was given", shared.controlAlias, shared.layerAlias, tabIndex)
	}
	if err := shared.validateTabStop(); err != nil {
		return err
	}
	defer focusManager.beginChange()()
	focusManager.registerStopLocked(shared.getControlIdentifier(), tabIndex)
	return nil
}

/*
RemoveFromTabIndex is a method which allows you to remove the control from its layer's tab order, so that Tab no longer
moves focus to it. The control keeps focus if it already has it, and can still be focused by clicking it or with
SetFocus.

Example:

	helpButton.RemoveFromTabIndex()
*/
func (shared *BaseControlInstanceType) RemoveFromTabIndex() {
	defer focusManager.beginChange()()
	scope, isFound := focusManager.scopes[shared.layerAlias]
	if !isFound {
		return
	}
	control := shared.getControlIdentifier()
	remainingStops := scope.stops[:0]
	for _, stop := range scope.stops {
		if stop.control != control {
			remainingStops = append(remainingStops, stop)
		}
	}
	scope.stops = remainingStops
}

/*
IsFocused is a method which allows you to check whether the control currently has keyboard focus.

Example:

	if okButton.IsFocused() { ... }
*/
func (shared *BaseControlInstanceType) IsFocused() bool {
	control := shared.getControlIdentifier()
	return isControlCurrentlyFocused(control.layerAlias, control.controlAlias, control.controlType)
}

/*
GetFocus is a method which allows you to give this control keyboard focus. It is shorthand for SetFocus with this
control, and does nothing if SetFocus would return an error, such as for a disabled or hidden control.

Example:

	control.GetFocus()
*/
func (shared *BaseControlInstanceType) GetFocus() *BaseControlInstanceType {
	_ = SetFocus(shared)
	return shared
}
