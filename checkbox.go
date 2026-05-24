package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
)

/*
CheckboxInstanceType is a structure which represents an instance of a checkbox control.

Example:

	var checkboxInstance CheckboxInstanceType
*/
type CheckboxInstanceType struct {
	BaseControlInstanceType
}

type checkboxType struct{}

var Checkbox checkboxType

var Checkboxes = memory.NewControlMemoryManager[types.CheckboxEntryType]()

/*
Delete is a method which removes a checkbox instance from its memory manager.

Example:

	checkbox.Delete()
*/
func (shared *CheckboxInstanceType) Delete() *CheckboxInstanceType {
	shared.BaseControlInstanceType.Delete()
	return nil
}

/*
AddToTabIndex is a method which adds the checkbox to the tab index of its associated layer.

Example:

	checkbox.AddToTabIndex()
*/
func (shared *CheckboxInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeCheckbox)
}

/*
IsSelected is a method which detects if the checkbox is selected. If the checkbox instance no
longer exists, then a result of false is always returned.

Example:

	isSelected := checkbox.IsSelected()
*/
func (shared *CheckboxInstanceType) IsSelected() bool {
	if Checkboxes.IsExists(shared.layerAlias, shared.controlAlias) {
		checkboxEntry := Checkboxes.Get(shared.layerAlias, shared.controlAlias)
		if checkboxEntry.IsSelected == true {
			return true
		}
	}
	return false
}

/*
SetState is a method which sets the selection state of the checkbox.

Example:

	checkbox.SetState(true)
*/
func (shared *CheckboxInstanceType) SetState(isChecked bool) {
	if Checkboxes.IsExists(shared.layerAlias, shared.controlAlias) {
		checkboxEntry := Checkboxes.Get(shared.layerAlias, shared.controlAlias)
		checkboxEntry.IsSelected = isChecked
	}
}

/*
Add is a method which adds a checkbox to a given text layer. Once called, an instance of your control is
returned which will allow you to read or manipulate the properties for it. The style of the checkbox will be determined
by the style entry passed in. If you wish to remove a checkbox from a text layer, simply call 'DeleteCheckbox'. In
addition, the following should be noted:

  - Checkboxes are not drawn physically to the text layer provided. Instead, they are rendered to the terminal at the same
    time when the text layer is rendered. This allows you to create checkboxes without actually overwriting the text layer
    data under it.

  - If the checkbox to be drawn falls outside the range of the provided layer, then only the visible portion of the
    checkbox will be drawn.

Example:

	checkboxInstance := Checkbox.Add("layer1", "cb1", "Enable Feature", style, 5, 5, false, true)
*/
func (shared *checkboxType) Add(layerAlias string, checkboxAlias string, checkboxLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool) CheckboxInstanceType {
	checkboxEntry := types.NewCheckboxEntry()
	checkboxEntry.Alias = checkboxAlias
	checkboxEntry.StyleEntry = styleEntry
	checkboxEntry.Label = checkboxLabel
	checkboxEntry.XLocation = xLocation
	checkboxEntry.YLocation = yLocation
	checkboxEntry.IsSelected = isSelected
	checkboxEntry.IsEnabled = isEnabled
	checkboxEntry.TooltipAlias = stringformat.GetLastSortedUUID()

	arrayOfRunes := stringformat.GetRunesFromString(checkboxLabel)
	labelWidth := stringformat.GetWidthOfRunesWhenPrinted(arrayOfRunes)
	// Create associated tooltip (always created but disabled by default)
	tooltipInstance := Tooltip.Add(layerAlias, checkboxEntry.TooltipAlias, "", styleEntry,
		checkboxEntry.XLocation, checkboxEntry.YLocation,
		labelWidth+2, 1,
		checkboxEntry.XLocation, checkboxEntry.YLocation+1,
		labelWidth+2, 3,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.SetEnabled(false)
	tooltipInstance.setParentControlAlias(checkboxAlias)
	// Use the ControlMemoryManager to add the checkbox entry
	Checkboxes.Add(layerAlias, checkboxAlias, &checkboxEntry)
	var checkboxInstance CheckboxInstanceType
	checkboxInstance.layerAlias = layerAlias
	checkboxInstance.controlAlias = checkboxAlias
	checkboxInstance.controlType = constants.TYPE_CHECKBOX
	return checkboxInstance
}

/*
Delete is a method which removes a checkbox from a text layer. In addition, the following should be
noted:

- If you attempt to delete a checkbox which does not exist, then the request will simply be ignored.

Example:

	Checkbox.Delete("layer1", "cb1")
*/
func (shared *checkboxType) Delete(layerAlias string, checkboxAlias string) {
	Checkboxes.Remove(layerAlias, checkboxAlias)
}

/*
DeleteAll is a method which deletes all checkboxes on a given text layer.

Example:

	Checkbox.DeleteAll("layer1")
*/
func (shared *checkboxType) DeleteAll(layerAlias string) {
	Checkboxes.RemoveAll(layerAlias)
}

/*
drawOnLayer is a method which draws all checkboxes on a given text layer.

Example:

	Checkbox.drawOnLayer(myLayer)
*/
func (shared *checkboxType) drawOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, checkboxEntry := range Checkboxes.GetAllEntries(layerAlias) {
		shared.draw(&layerEntry, checkboxEntry.Alias, checkboxEntry.Label, checkboxEntry.StyleEntry, checkboxEntry.XLocation, checkboxEntry.YLocation, checkboxEntry.IsSelected, checkboxEntry.IsEnabled)
	}
}

/*
draw is a method which draws a checkbox on a given text layer. The style of the checkbox will be
determined by the style entry passed in. In addition, the following should be noted:

  - Checkboxes are not drawn physically to the text layer provided. Instead, they are rendered to the terminal at the same
    time when the text layer is rendered. This allows you to create checkboxes without actually overwriting the text layer
    data under it.

  - If the checkbox to be drawn falls outside the range of the provided layer, then only the visible portion of the
    checkbox will be drawn.

Example:

	Checkbox.draw(&myLayer, "cb1", "Enable Feature", style, 0, 0, false, true)
*/
func (shared *checkboxType) draw(layerEntry *types.LayerEntryType, checkboxAlias string, checkboxLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool) {
	localStyleEntry := types.NewTuiStyleEntry(&styleEntry)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = localStyleEntry.Checkbox.ForegroundColor
	attributeEntry.BackgroundColor = localStyleEntry.Checkbox.BackgroundColor
	attributeEntry.CellType = constants.CellTypeCheckbox
	attributeEntry.CellControlAlias = checkboxAlias

	var secondArrayOfRunes []rune
	if isSelected {
		secondArrayOfRunes = []rune{localStyleEntry.Checkbox.SelectedCharacter}
		attributeEntry.CellControlId = constants.CellControlIdChecked
	} else {
		secondArrayOfRunes = []rune{localStyleEntry.Checkbox.UnselectedCharacter}
		attributeEntry.CellControlId = constants.CellControlIdUnchecked
	}
	layer.printLayer(layerEntry, attributeEntry, xLocation, yLocation, secondArrayOfRunes)
	firstArrayOfRunes := stringformat.GetRunesFromString(checkboxLabel)
	layer.printLayer(layerEntry, attributeEntry, xLocation+2, yLocation, firstArrayOfRunes)
}

/*
updateMouseEvent is a method which updates the state of all checkboxes according to the current mouse event
state.

Example:

	isUpdateNeeded := Checkbox.updateMouseEvent()
*/
func (shared *checkboxType) updateMouseEvent() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	controlAlias := characterEntry.AttributeEntry.CellControlAlias
	if characterEntry.AttributeEntry.CellType == constants.CellTypeCheckbox && characterEntry.AttributeEntry.CellControlId != constants.NullCellId {
		_, _, previousButtonPressed, _ := GetPreviousMouseStatus()
		if buttonPressed != 0 && previousButtonPressed == 0 && Checkboxes.IsExists(layerAlias, controlAlias) {
			eventStateMemory.currentlyFocusedControl.layerAlias = layerAlias
			eventStateMemory.currentlyFocusedControl.controlAlias = controlAlias
			eventStateMemory.currentlyFocusedControl.controlType = constants.CellTypeCheckbox
			checkboxEntry := Checkboxes.Get(layerAlias, controlAlias)
			if !checkboxEntry.IsEnabled {
				return isUpdateRequired
			}
			if checkboxEntry.IsSelected {
				checkboxEntry.IsSelected = false
			} else {
				checkboxEntry.IsSelected = true
			}
			return isUpdateRequired
		}
	}
	return isUpdateRequired
}
