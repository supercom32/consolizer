package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
)

type CheckboxInstanceType struct {
	BaseControlInstanceType
}

type checkboxType struct{}

var Checkbox checkboxType

var Checkboxes = memory.NewControlMemoryManager[types.CheckboxEntryType]()

/*
Delete is a method which allows you to remove a checkbox instance from its memory manager.

:return: A nil pointer of type CheckboxInstanceType.

Example:

	checkbox.Delete()
*/
func (shared *CheckboxInstanceType) Delete() *CheckboxInstanceType {
	shared.BaseControlInstanceType.Delete()
	return nil
}

/*
AddToTabIndex is a method which allows you to add the checkbox to the tab index of its associated layer.

Example:

	checkbox.AddToTabIndex()
*/
func (shared *CheckboxInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeCheckbox)
}

/*
IsSelected is a method which allows you to detect if the checkbox is selected. If the checkbox instance no
longer exists, then a result of false is always returned.

:return: True if the checkbox is selected, otherwise false.

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
SetState is a method which allows you to set the selection state of the checkbox.

:param isChecked: Set to true to select the checkbox, or false to deselect it.

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
Add is a method which allows you to add a checkbox to a given text layer. Once called, an instance of your control is
returned which will allow you to read or manipulate the properties for it. The style of the checkbox will be determined
by the style entry passed in. If you wish to remove a checkbox from a text layer, simply call 'DeleteCheckbox'. In
addition, the following should be noted:

  - Checkboxes are not drawn physically to the text layer provided. Instead, they are rendered to the terminal at the same
    time when the text layer is rendered. This allows you to create checkboxes without actually overwriting the text layer
    data under it.

  - If the checkbox to be drawn falls outside the range of the provided layer, then only the visible portion of the
    checkbox will be drawn.

:param layerAlias: The alias of the layer to add the checkbox to.
:param checkboxAlias: A unique alias for the checkbox.
:param checkboxLabel: The text label to display next to the checkbox.
:param styleEntry: The visual style to apply to the checkbox.
:param xLocation: The X coordinate of the checkbox.
:param yLocation: The Y coordinate of the checkbox.
:param isSelected: Set to true if the checkbox should be initially selected.
:param isEnabled: Set to true if the checkbox should be initially enabled.

:return: An instance of the newly created checkbox.

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
Delete is a method which allows you to remove a checkbox from a text layer. In addition, the following should be
noted:

- If you attempt to delete a checkbox which does not exist, then the request will simply be ignored.

:param layerAlias: The alias of the layer the checkbox is on.
:param checkboxAlias: The alias of the checkbox to remove.

Example:

	Checkbox.Delete("layer1", "cb1")
*/
func (shared *checkboxType) Delete(layerAlias string, checkboxAlias string) {
	Checkboxes.Remove(layerAlias, checkboxAlias)
}

/*
DeleteAll is a method which allows you to delete all checkboxes on a given text layer.

:param layerAlias: The alias of the layer to remove all checkboxes from.

Example:

	Checkbox.DeleteAll("layer1")
*/
func (shared *checkboxType) DeleteAll(layerAlias string) {
	Checkboxes.RemoveAll(layerAlias)
}

/*
drawOnLayer is a method which allows you to draw all checkboxes on a given text layer.

:param layerEntry: The LayerEntryType structure representing the layer to draw on.

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
draw is a method which allows you to draw a checkbox on a given text layer. The style of the checkbox will be
determined by the style entry passed in. In addition, the following should be noted:

- Checkboxes are not drawn physically to the text layer provided. Instead, they are rendered to the terminal at the.

- If the checkbox to be drawn falls outside the range of the provided layer, then only the visible portion of the.

:param layerEntry: A pointer to the LayerEntryType to draw the checkbox on.
:param checkboxAlias: The unique alias of the checkbox.
:param checkboxLabel: The text label to display next to the checkbox.
:param styleEntry: The visual style to apply to the checkbox.
:param xLocation: The X coordinate for the checkbox.
:param yLocation: The Y coordinate for the checkbox.
:param isSelected: Set to true if the checkbox is currently selected.
:param isEnabled: Set to true if the checkbox is currently enabled.

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
updateMouseEvent is a method which allows you to update the state of all checkboxes according to the current mouse event
state.

:return: True if a screen update is required, otherwise false.

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
