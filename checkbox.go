package consolizer

import (
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
)

type CheckboxInstanceType struct {
	layerAlias  string
	checkboxAlias string
}

type checkboxType struct {}

var checkbox checkboxType

/*
IsCheckboxSelected allows you to detect if the given checkbox is selected or not. If the checkbox instance
no longer exists, then a result of false is always returned.
*/
func (shared *CheckboxInstanceType) IsCheckboxSelected() bool {
	if memory.IsCheckboxExists(shared.layerAlias, shared.checkboxAlias){
		checkboxEntry := memory.GetCheckbox(shared.layerAlias, shared.checkboxAlias)
		if checkboxEntry.IsSelected == true {
			return true
		}
	}
	return false
}

/*
AddCheckbox allows you to add a checkbox to a given text layer. Once called, an instance
of your control is returned which will allow you to read or manipulate the properties for it.
The Style of the checkbox will be determined by the style entry passed in. If you wish to
remove a checkbox from a text layer, simply call 'DeleteCheckbox'. In addition, the following
information should be noted:

- Checkboxes are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create checkboxes without actually overwriting
the text layer data under it.

- If the checkbox to be drawn falls outside the range of the provided layer,
then only the visible portion of the checkbox will be drawn.
*/
 func (shared *checkboxType) AddCheckbox(layerAlias string, checkboxAlias string, checkboxLabel string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool) CheckboxInstanceType {
	memory.AddCheckbox(layerAlias, checkboxAlias, checkboxLabel, styleEntry, xLocation, yLocation, isSelected)
	var checkboxInstance CheckboxInstanceType
	 checkboxInstance.layerAlias = layerAlias
	 checkboxInstance.checkboxAlias = checkboxAlias
	return checkboxInstance
}

/*
DeleteCheckbox allows you to remove a checkbox from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a checkbox which does not exist, then the request
will simply be ignored.
*/
func (shared *checkboxType) DeleteCheckbox(layerAlias string, checkboxAlias string) {
	memory.DeleteCheckbox(layerAlias, checkboxAlias)
}

/*
drawCheckboxesOnLayer allows you to draw all checkboxes on a given text layer.
*/
func (shared *checkboxType) drawCheckboxesOnLayer(layerEntry memory.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for currentKey := range memory.CheckboxMemory[layerAlias] {
		checkboxEntry := memory.GetCheckbox(layerAlias, currentKey)
		shared.drawCheckbox(&layerEntry, currentKey, checkboxEntry.Label, checkboxEntry.StyleEntry, checkboxEntry.XLocation, checkboxEntry.YLocation, checkboxEntry.IsSelected, checkboxEntry.IsEnabled)
	}
}

/*
drawCheckbox allows you to draw A checkbox on a given text layer. The
Style of the checkbox will be determined by the style entry passed in. In
addition, the following information should be noted:

- Checkboxes are not drawn physically to the text layer provided. Instead,
they are rendered to the terminal at the same time when the text layer is
rendered. This allows you to create checkboxes without actually overwriting
the text layer data under it.

- If the checkbox to be drawn falls outside the range of the provided layer,
then only the visible portion of the checkbox will be drawn.
*/
func (shared *checkboxType) drawCheckbox (layerEntry *memory.LayerEntryType, checkboxAlias string, checkboxLabel string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool ) {
	localStyleEntry := memory.NewTuiStyleEntry(&styleEntry)
	attributeEntry := memory.NewAttributeEntry()
	attributeEntry.ForegroundColor = localStyleEntry.CheckboxForegroundColor
	attributeEntry.BackgroundColor = localStyleEntry.CheckboxBackgroundColor
	attributeEntry.CellType = constants.CellTypeCheckbox
	attributeEntry.CellControlAlias = checkboxAlias
	firstArrayOfRunes := stringformat.GetRunesFromString(checkboxLabel)
	firstArrayOfRunes = append(firstArrayOfRunes, ' ')
	numberOfSpacesUsed := stringformat.GetWidthOfRunesWhenPrinted(firstArrayOfRunes)
	printLayer(layerEntry, attributeEntry, xLocation, yLocation, firstArrayOfRunes)
	var secondArrayOfRunes []rune
	if isSelected {
		secondArrayOfRunes = []rune{localStyleEntry.CheckboxSelectedCharacter}
		attributeEntry.CellControlId = constants.CellControlIdChecked
	} else {
		secondArrayOfRunes = []rune{localStyleEntry.CheckboxUnselectedCharacter}
		attributeEntry.CellControlId = constants.CellControlIdUnchecked
	}
	printLayer(layerEntry, attributeEntry, xLocation + numberOfSpacesUsed, yLocation, secondArrayOfRunes)
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
		if buttonPressed != 0 && previousButtonPressed == 0 {
			checkboxEntry := memory.GetCheckbox(layerAlias, controlAlias)
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
