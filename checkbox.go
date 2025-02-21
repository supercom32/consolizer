package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
	"github.com/supercom32/consolizer/types"
)

type CheckboxInstanceType struct {
	layerAlias   string
	controlAlias string
}

type checkboxType struct{}

var Checkbox checkboxType

func (shared *CheckboxInstanceType) Delete() *CheckboxInstanceType {
	if memory.IsCheckboxExists(shared.layerAlias, shared.controlAlias) {
		memory.DeleteCheckbox(shared.layerAlias, shared.controlAlias)
	}
	return nil
}

func (shared *CheckboxInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeCheckbox)
}

/*
IsCheckboxSelected allows you to detect if the given Checkbox is selected or not. If the Checkbox instance
no longer exists, then a result of false is always returned.
*/
func (shared *CheckboxInstanceType) IsCheckboxSelected() bool {
	if memory.IsCheckboxExists(shared.layerAlias, shared.controlAlias) {
		checkboxEntry := memory.GetCheckbox(shared.layerAlias, shared.controlAlias)
		if checkboxEntry.IsSelected == true {
			return true
		}
	}
	return false
}

/*
Add allows you to add a Checkbox to a given text layer. Once called, an instance
of your control is returned which will allow you to read or manipulate the properties for it.
The Style of the Checkbox will be determined by the style entry passed in. If you wish to
remove a Checkbox from a text layer, simply call 'DeleteCheckbox'. In addition, the following
information should be noted:

- Checkboxes are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create checkboxes without actually overwriting
the text layer data under it.

- If the Checkbox to be drawn falls outside the range of the provided layer,
then only the visible portion of the Checkbox will be drawn.
*/
func (shared *checkboxType) Add(layerAlias string, checkboxAlias string, checkboxLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool) CheckboxInstanceType {
	memory.AddCheckbox(layerAlias, checkboxAlias, checkboxLabel, styleEntry, xLocation, yLocation, isSelected, isEnabled)
	var checkboxInstance CheckboxInstanceType
	checkboxInstance.layerAlias = layerAlias
	checkboxInstance.controlAlias = checkboxAlias
	return checkboxInstance
}

/*
DeleteCheckbox allows you to remove a Checkbox from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a Checkbox which does not exist, then the request
will simply be ignored.
*/
func (shared *checkboxType) DeleteCheckbox(layerAlias string, checkboxAlias string) {
	memory.DeleteCheckbox(layerAlias, checkboxAlias)
}

/*
DeleteAllCheckboxesFromLayer allows you to delete all checkboxes on a given text layer.
*/
func (shared *checkboxType) DeleteAllCheckboxesFromLayer(layerAlias string) {
	memory.DeleteAllCheckboxesFromLayer(layerAlias)
}

/*
drawCheckboxesOnLayer allows you to draw all checkboxes on a given text layer.
*/
func (shared *checkboxType) drawCheckboxesOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, checkboxEntry := range memory.Checkboxes.GetAllEntries(layerAlias) {
		shared.drawCheckbox(&layerEntry, checkboxEntry.Alias, checkboxEntry.Label, checkboxEntry.StyleEntry, checkboxEntry.XLocation, checkboxEntry.YLocation, checkboxEntry.IsSelected, checkboxEntry.IsEnabled)
	}
}

/*
drawCheckbox allows you to draw A Checkbox on a given text layer. The
Style of the Checkbox will be determined by the style entry passed in. In
addition, the following information should be noted:

- Checkboxes are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create checkboxes without actually overwriting
the text layer data under it.

- If the Checkbox to be drawn falls outside the range of the provided layer,
then only the visible portion of the Checkbox will be drawn.
*/
func (shared *checkboxType) drawCheckbox(layerEntry *types.LayerEntryType, checkboxAlias string, checkboxLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool) {
	localStyleEntry := types.NewTuiStyleEntry(&styleEntry)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = localStyleEntry.CheckboxForegroundColor
	attributeEntry.BackgroundColor = localStyleEntry.CheckboxBackgroundColor
	attributeEntry.CellType = constants.CellTypeCheckbox
	attributeEntry.CellControlAlias = checkboxAlias

	var secondArrayOfRunes []rune
	if isSelected {
		secondArrayOfRunes = []rune{localStyleEntry.CheckboxSelectedCharacter}
		attributeEntry.CellControlId = constants.CellControlIdChecked
	} else {
		secondArrayOfRunes = []rune{localStyleEntry.CheckboxUnselectedCharacter}
		attributeEntry.CellControlId = constants.CellControlIdUnchecked
	}
	printLayer(layerEntry, attributeEntry, xLocation, yLocation, secondArrayOfRunes)
	firstArrayOfRunes := stringformat.GetRunesFromString(checkboxLabel)
	printLayer(layerEntry, attributeEntry, xLocation+2, yLocation, firstArrayOfRunes)
}

/*
updateMouseEventCheckbox allows you to update the state of all checkboxes according to the current mouse event state.
In the event that a screen update is required this method returns true.
*/
func (shared *checkboxType) updateMouseEventCheckbox() bool {
	isUpdateRequired := false
	mouseXLocation, mouseYLocation, buttonPressed, _ := memory.GetMouseStatus()
	characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	layerAlias := characterEntry.LayerAlias
	controlAlias := characterEntry.AttributeEntry.CellControlAlias
	if characterEntry.AttributeEntry.CellType == constants.CellTypeCheckbox && characterEntry.AttributeEntry.CellControlId != constants.NullCellId {
		_, _, previousButtonPressed, _ := memory.GetPreviousMouseStatus()
		if buttonPressed != 0 && previousButtonPressed == 0 && memory.IsCheckboxExists(layerAlias, controlAlias) {
			eventStateMemory.currentlyFocusedControl.layerAlias = layerAlias
			eventStateMemory.currentlyFocusedControl.controlAlias = controlAlias
			eventStateMemory.currentlyFocusedControl.controlType = constants.CellTypeCheckbox
			checkboxEntry := memory.GetCheckbox(layerAlias, controlAlias)
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
