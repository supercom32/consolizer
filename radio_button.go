package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
	"github.com/supercom32/consolizer/types"
)

type RadioButtonInstanceType struct {
	layerAlias       string
	radioButtonAlias string
}

type radioButtonType struct{}

var radioButton radioButtonType

func (shared *RadioButtonInstanceType) Delete() *RadioButtonInstanceType {
	if memory.IsRadioButtonExists(shared.layerAlias, shared.radioButtonAlias) {
		memory.DeleteRadioButton(shared.layerAlias, shared.radioButtonAlias)
	}
	return nil
}

/*
IsRadioButtonSelected allows you to detect if the given radio button is selected or not. If the radio button instance
no longer exists, then a result of false is always returned.
*/
func (shared *RadioButtonInstanceType) IsRadioButtonSelected() bool {
	if memory.IsRadioButtonExists(shared.layerAlias, shared.radioButtonAlias) {
		selectedRadioButton := getSelectedRadioButton(shared.layerAlias, shared.radioButtonAlias)
		if selectedRadioButton == shared.radioButtonAlias {
			return true
		}
	}
	return false
}

/*
GetSelectedRadioButton allows you to retrieve the radio button currently selected. If the radio button instance
no longer exists, then a result of false is always returned.
*/
func (shared *RadioButtonInstanceType) GetSelectedRadioButton() string {
	if memory.IsRadioButtonExists(shared.layerAlias, shared.radioButtonAlias) {
		return getSelectedRadioButton(shared.layerAlias, shared.radioButtonAlias)
	}
	return ""
}

/*
SetIsVisible allows you to change if a radio button currently visible or not. If the radio button instance
no longer exists, then nothing is done.
*/
func (shared *RadioButtonInstanceType) SetIsVisible(isVisible bool) {
	if memory.IsRadioButtonExists(shared.layerAlias, shared.radioButtonAlias) {
		radioButtonEntry := memory.GetRadioButton(shared.layerAlias, shared.radioButtonAlias)
		radioButtonEntry.IsVisible = isVisible
	}
}

/*
Add allows you to add a radio button to a given text layer. Once called, an instance
of your control is returned which will allow you to read or manipulate the properties for it.
The Style of the radio button will be determined by the style entry passed in. If you wish to
remove a radio button from a text layer, simply call 'DeleteRadioButton'. In addition, the
following information should be noted:

- Radio buttons are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create radio buttons without actually overwriting
the text layer data under it.

- If the radio button to be drawn falls outside the range of the provided layer,
then only the visible portion of the radio button will be drawn.

- The group ID allows you to specify which collection of radio buttons belong together. Only one radio
button may be selected at any given time for a particular group.

- If the radio button being created is marked as being selected, then any previously selected radio button
with the same group ID becomes unselected.
*/
func (shared *radioButtonType) Add(layerAlias string, radioButtonAlias string, radioButtonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, groupId int, isSelected bool) RadioButtonInstanceType {
	memory.AddRadioButton(layerAlias, radioButtonAlias, radioButtonLabel, styleEntry, xLocation, yLocation, groupId, isSelected)
	var radioButtonInstance RadioButtonInstanceType
	radioButtonInstance.layerAlias = layerAlias
	radioButtonInstance.radioButtonAlias = radioButtonAlias
	if isSelected {
		selectRadioButton(layerAlias, radioButtonAlias)
	}
	return radioButtonInstance
}

/*
drawRadioButtonsOnLayer allows you to draw all radio buttons on a given text layer.
*/
func (shared *radioButtonType) drawRadioButtonsOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for currentKey := range memory.RadioButton.Entries[layerAlias] {
		radioButtonEntry := memory.GetRadioButton(layerAlias, currentKey)
		shared.drawRadioButton(&layerEntry, currentKey, radioButtonEntry.Label, radioButtonEntry.StyleEntry, radioButtonEntry.XLocation, radioButtonEntry.YLocation, radioButtonEntry.IsSelected, radioButtonEntry.IsEnabled)
	}
}

/*
drawRadioButton allows you to draw A radio button on a given text layer. The
Style of the radio button will be determined by the style entry passed in. In
addition, the following information should be noted:

- Radio buttons are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create radio buttons without actually overwriting
the text layer data under it.

- If the radio button to be drawn falls outside the range of the provided layer,
then only the visible portion of the radio button will be drawn.
*/
func (shared *radioButtonType) drawRadioButton(layerEntry *types.LayerEntryType, radioButtonAlias string, radioButtonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool) {
	localStyleEntry := types.NewTuiStyleEntry(&styleEntry)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = localStyleEntry.RadioButtonForegroundColor
	attributeEntry.BackgroundColor = localStyleEntry.RadioButtonBackgroundColor
	attributeEntry.CellType = constants.CellTypeRadioButton
	attributeEntry.CellControlAlias = radioButtonAlias
	firstArrayOfRunes := stringformat.GetRunesFromString(radioButtonLabel)
	firstArrayOfRunes = append(firstArrayOfRunes, ' ')
	numberOfSpacesUsed := stringformat.GetWidthOfRunesWhenPrinted(firstArrayOfRunes)
	printLayer(layerEntry, attributeEntry, xLocation, yLocation, firstArrayOfRunes)
	var secondArrayOfRunes []rune
	if isSelected {
		secondArrayOfRunes = []rune{localStyleEntry.RadioButtonSelectedCharacter}
		attributeEntry.CellControlId = constants.CellControlIdChecked
	} else {
		secondArrayOfRunes = []rune{localStyleEntry.RadioButtonUnselectedCharacter}
		attributeEntry.CellControlId = constants.CellControlIdUnchecked
	}
	printLayer(layerEntry, attributeEntry, xLocation+numberOfSpacesUsed, yLocation, secondArrayOfRunes)
}

/*
DeleteRadioButton allows you to remove a radio button from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a radio button which does not exist, then the request
will simply be ignored.
*/
func (shared *radioButtonType) DeleteRadioButton(layerAlias string, radioButtonAlias string) {
	memory.DeleteRadioButton(layerAlias, radioButtonAlias)
}

func (shared *radioButtonType) DeleteAllRadioButtons(layerAlias string) {
	memory.DeleteAllRadioButtonsFromLayer(layerAlias)
}

/*
updateMouseEventRadioButton allows you to update the state of all radio buttons according to the current mouse event state.
In the event that a screen update is required this method returns true.
*/
func (shared *radioButtonType) updateMouseEventRadioButton() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	controlAlias := characterEntry.AttributeEntry.CellControlAlias
	if characterEntry.AttributeEntry.CellType == constants.CellTypeRadioButton && characterEntry.AttributeEntry.CellControlId != constants.NullCellId {
		_, _, previousButtonPressed, _ := memory.GetPreviousMouseStatus()
		if buttonPressed != 0 && previousButtonPressed == 0 && memory.IsRadioButtonExists(layerAlias, controlAlias) {
			selectRadioButton(layerAlias, controlAlias)
			isUpdateRequired = true
			return isUpdateRequired
		}
	}
	return isUpdateRequired
}

/*
selectRadioButton allows you to select a radio button on a given text layer. Since only one radio button may
be selected at a time, any previously selected radio button becomes unselected.
*/
func selectRadioButton(layerAlias string, radioButtonAlias string) {
	radioButtonSelectedEntry := memory.GetRadioButton(layerAlias, radioButtonAlias)
	for currentRadioButtonKey := range memory.RadioButton.Entries[layerAlias] {
		currentRadioButtonEntry := memory.GetRadioButton(layerAlias, currentRadioButtonKey)
		if currentRadioButtonKey == radioButtonAlias {
			currentRadioButtonEntry.IsSelected = true
		} else if currentRadioButtonEntry.GroupId == radioButtonSelectedEntry.GroupId {
			if currentRadioButtonEntry.IsSelected == true {
				currentRadioButtonEntry.IsSelected = false
			}
		}
	}
}

/*
getSelectedRadioButton allows you to obtain the selected radio button for a particular group ID.
The group ID used is automatically determined based on the radio button alias given.
*/
func getSelectedRadioButton(layerAlias string, radioButtonAlias string) string {
	selectedItem := ""
	radioButtonEntry := memory.GetRadioButton(layerAlias, radioButtonAlias)
	for currentRadioButtonKey := range memory.RadioButton.Entries[layerAlias] {
		currentRadioButtonEntry := memory.GetRadioButton(layerAlias, currentRadioButtonKey)
		if currentRadioButtonEntry.GroupId == radioButtonEntry.GroupId {
			if currentRadioButtonEntry.IsSelected {
				return currentRadioButtonKey
			}
		}
	}
	return selectedItem
}
