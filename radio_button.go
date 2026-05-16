package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
)

type RadioButtonInstanceType struct {
	BaseControlInstanceType
}

type radioButtonType struct{}

var radioButton radioButtonType
var RadioButtons = memory.NewControlMemoryManager[types.RadioButtonEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
IsRadioButtonExists is a method which allows you to check if a radio button with the specified alias exists on a given
layer.

:param layerAlias: The alias of the layer to check.
:param radioButtonAlias: The alias of the radio button to look for.

:return: A boolean indicating if the radio button exists.
*/
func IsRadioButtonExists(layerAlias string, radioButtonAlias string) bool {
	// Use ControlMemoryManager to check if the radio button exists
	return RadioButtons.Get(layerAlias, radioButtonAlias) != nil
}

/*
AddToTabIndex is a method which allows you to add the radio button to the tab navigation index.
*/
func (shared *RadioButtonInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeRadioButton)
}

/*
Delete is a method which allows you to remove the radio button instance.

:return: A nil pointer of RadioButtonInstanceType.
*/
func (shared *RadioButtonInstanceType) Delete() *RadioButtonInstanceType {
	if RadioButtons.IsExists(shared.layerAlias, shared.controlAlias) {
		RadioButtons.Remove(shared.layerAlias, shared.controlAlias)
	}
	return nil
}

/*
IsSelected is a method which allows you to detect if the given radio button is selected or not. In addition,
the following information should be noted:

- If the radio button instance no longer exists, then a result of false is always returned.

:return: A boolean indicating whether the radio button is selected.
*/
func (shared *RadioButtonInstanceType) IsSelected() bool {
	if RadioButtons.IsExists(shared.layerAlias, shared.controlAlias) {
		selectedRadioButton := getSelectedRadioButton(shared.layerAlias, shared.controlAlias)
		if selectedRadioButton == shared.controlAlias {
			return true
		}
	}
	return false
}

/*
GetSelected is a method which allows you to retrieve the alias of the radio button currently selected within
the same group. In addition, the following information should be noted:

- If the radio button instance no longer exists, then an empty string is returned.

:return: The alias of the currently selected radio button.
*/
func (shared *RadioButtonInstanceType) GetSelected() string {
	if RadioButtons.IsExists(shared.layerAlias, shared.controlAlias) {
		return getSelectedRadioButton(shared.layerAlias, shared.controlAlias)
	}
	return ""
}

/*
Add is a method which allows you to add a radio button to a given text layer. Once called, an instance of your control
is returned which will allow you to read or manipulate the properties for it. In addition, the following information
should be noted:

- The Style of the radio button will be determined by the style entry passed in.

- If you wish to remove a radio button from a text layer, simply call 'DeleteRadioButton'.

- Radio buttons are not drawn physically to the text layer provided.

- Instead, they are rendered to the terminal at the same time when the text layer is rendered.

- This allows you to create radio buttons without actually overwriting the text layer data under it.

  - If the radio button to be drawn falls outside the range of the provided layer, then only the visible portion of the
    radio button will be drawn.

- The group ID allows you to specify which collection of radio buttons belong together.

- Only one radio button may be selected at any given time for a particular group.

  - If the radio button being created is marked as being selected, then any previously selected radio button with the same
    group ID becomes unselected.

:param layerAlias: The alias of the layer to which the radio button will be added.
:param radioButtonAlias: The unique alias for the radio button control.
:param radioButtonLabel: The label text for the radio button.
:param styleEntry: The style configuration for the radio button.
:param xLocation: The x coordinate of the position.
:param yLocation: The y coordinate of the position.
:param groupId: The ID of the group this radio button belongs to.
:param isSelected: A boolean indicating if the radio button should be initially selected.

:return: An instance of the created radio button.
*/
func (shared *radioButtonType) Add(layerAlias string, radioButtonAlias string, radioButtonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, groupId int, isSelected bool) RadioButtonInstanceType {
	radioButtonEntry := types.NewRadioButtonEntry()
	radioButtonEntry.Alias = radioButtonAlias
	radioButtonEntry.StyleEntry = styleEntry
	radioButtonEntry.Label = radioButtonLabel
	radioButtonEntry.XLocation = xLocation
	radioButtonEntry.YLocation = yLocation
	radioButtonEntry.GroupId = groupId
	radioButtonEntry.IsSelected = isSelected
	radioButtonEntry.TooltipAlias = stringformat.GetLastSortedUUID()

	// Create associated tooltip (always created but disabled by default)
	tooltipInstance := Tooltip.Add(layerAlias, radioButtonEntry.TooltipAlias, "", styleEntry,
		radioButtonEntry.XLocation, radioButtonEntry.YLocation,
		len(radioButtonLabel)+2, 1,
		radioButtonEntry.XLocation, radioButtonEntry.YLocation+1,
		len(radioButtonLabel)+2, 3,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.SetEnabled(false)
	tooltipInstance.setParentControlAlias(radioButtonAlias)
	// Use the ControlMemoryManager to add the radio button entry
	RadioButtons.Add(layerAlias, radioButtonAlias, &radioButtonEntry)
	var radioButtonInstance RadioButtonInstanceType
	radioButtonInstance.layerAlias = layerAlias
	radioButtonInstance.controlAlias = radioButtonAlias
	radioButtonInstance.controlType = constants.TYPE_RADIOBUTTON
	if isSelected {
		selectRadioButton(layerAlias, radioButtonAlias)
	}
	return radioButtonInstance
}

/*
drawOnLayer is a method which allows you to draw all radio buttons on a given text layer.

:param layerEntry: The layer entry on which to draw the radio buttons.
*/
func (shared *radioButtonType) drawOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentRadioButtonEntry := range RadioButtons.GetAllEntries(layerAlias) {
		radioButtonEntry := currentRadioButtonEntry
		shared.draw(&layerEntry, radioButtonEntry.Alias, radioButtonEntry.Label, radioButtonEntry.StyleEntry, radioButtonEntry.XLocation, radioButtonEntry.YLocation, radioButtonEntry.IsSelected, radioButtonEntry.IsEnabled)
	}
}

/*
draw is a method which allows you to draw a radio button on a given text layer. In addition, the following
information should be noted:

- The Style of the radio button will be determined by the style entry passed in.

- Radio buttons are not drawn physically to the text layer provided.

- Instead, they are rendered to the terminal at the same time when the text layer is rendered.

- This allows you to create radio buttons without actually overwriting the text layer data under it.

  - If the radio button to be drawn falls outside the range of the provided layer, then only the visible portion of the
    radio button will be drawn.

:param layerEntry: The layer on which to draw the radio button.
:param radioButtonAlias: The alias of the radio button.
:param radioButtonLabel: The label text to display.
:param styleEntry: The style configuration to use.
:param xLocation: The x coordinate of the position.
:param yLocation: The y coordinate of the position.
:param isSelected: Whether the radio button is selected.
:param isEnabled: Whether the radio button is enabled.
*/
func (shared *radioButtonType) draw(layerEntry *types.LayerEntryType, radioButtonAlias string, radioButtonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool) {
	localStyleEntry := types.NewTuiStyleEntry(&styleEntry)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = localStyleEntry.RadioButton.ForegroundColor
	attributeEntry.BackgroundColor = localStyleEntry.RadioButton.BackgroundColor
	attributeEntry.CellType = constants.CellTypeRadioButton
	attributeEntry.CellControlAlias = radioButtonAlias
	firstArrayOfRunes := stringformat.GetRunesFromString(radioButtonLabel)
	firstArrayOfRunes = append(firstArrayOfRunes, ' ')
	numberOfSpacesUsed := stringformat.GetWidthOfRunesWhenPrinted(firstArrayOfRunes)
	layer.printLayer(layerEntry, attributeEntry, xLocation, yLocation, firstArrayOfRunes)
	var secondArrayOfRunes []rune
	if isSelected {
		secondArrayOfRunes = []rune{localStyleEntry.RadioButton.SelectedCharacter}
		attributeEntry.CellControlId = constants.CellControlIdChecked
	} else {
		secondArrayOfRunes = []rune{localStyleEntry.RadioButton.UnselectedCharacter}
		attributeEntry.CellControlId = constants.CellControlIdUnchecked
	}
	layer.printLayer(layerEntry, attributeEntry, xLocation+numberOfSpacesUsed, yLocation, secondArrayOfRunes)
}

/*
Delete is a method which allows you to remove a radio button from a text layer. In addition, the following
information should be noted:

- If you attempt to delete a radio button which does not exist, then the request will simply be ignored.

:param layerAlias: The alias of the layer from which to remove the radio button.
:param radioButtonAlias: The alias of the radio button to be removed.
*/
func (shared *radioButtonType) Delete(layerAlias string, radioButtonAlias string) {
	RadioButtons.Remove(layerAlias, radioButtonAlias)
}

/*
DeleteAll is a method which allows you to remove all radio buttons from a specified text layer.

:param layerAlias: The alias of the layer from which all radio buttons will be removed.
*/
func (shared *radioButtonType) DeleteAll(layerAlias string) {
	RadioButtons.RemoveAll(layerAlias)
}

/*
updateMouseEvent is a method which allows you to update the state of all radio buttons according to the current mouse
event state. In addition, the following information should be noted:

- In the event that a screen update is required this method returns true.

:return: A boolean indicating if a screen update is required.
*/
func (shared *radioButtonType) updateMouseEvent() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	controlAlias := characterEntry.AttributeEntry.CellControlAlias
	if characterEntry.AttributeEntry.CellType == constants.CellTypeRadioButton && characterEntry.AttributeEntry.CellControlId != constants.NullCellId {
		_, _, previousButtonPressed, _ := GetPreviousMouseStatus()
		if buttonPressed != 0 && previousButtonPressed == 0 && RadioButtons.IsExists(layerAlias, controlAlias) {
			selectRadioButton(layerAlias, controlAlias)
			isUpdateRequired = true
			return isUpdateRequired
		}
	}
	return isUpdateRequired
}

/*
selectRadioButton is a method which allows you to select a radio button on a given text layer. In addition, the
following information should be noted:

- Since only one radio button may be selected at a time, any previously selected radio button becomes unselected.

:param layerAlias: The alias of the layer.
:param radioButtonAlias: The alias of the radio button to select.
*/
func selectRadioButton(layerAlias string, radioButtonAlias string) {
	radioButtonSelectedEntry := RadioButtons.Get(layerAlias, radioButtonAlias)
	for _, currentRadioButtonEntry := range RadioButtons.GetAllEntries(layerAlias) {
		if currentRadioButtonEntry.Alias == radioButtonAlias {
			currentRadioButtonEntry.IsSelected = true
		} else if currentRadioButtonEntry.GroupId == radioButtonSelectedEntry.GroupId {
			if currentRadioButtonEntry.IsSelected == true {
				currentRadioButtonEntry.IsSelected = false
			}
		}
	}
}

/*
getSelectedRadioButton is a method which allows you to obtain the selected radio button for a particular group ID. In
addition, the following information should be noted:

- The group ID used is automatically determined based on the radio button alias given.

:param layerAlias: The alias of the layer.
:param radioButtonAlias: An alias of a radio button within the target group.

:return: The alias of the currently selected radio button in the group.
*/
func getSelectedRadioButton(layerAlias string, radioButtonAlias string) string {
	selectedItem := ""
	radioButtonEntry := RadioButtons.Get(layerAlias, radioButtonAlias)
	for _, currentRadioButtonEntry := range RadioButtons.GetAllEntries(layerAlias) {
		if currentRadioButtonEntry.GroupId == radioButtonEntry.GroupId {
			if currentRadioButtonEntry.IsSelected {
				return currentRadioButtonEntry.Alias
			}
		}
	}
	return selectedItem
}
