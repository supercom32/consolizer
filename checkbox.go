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

/*
IsButtonPressed allows you to detect if any text button was pressed or not. In
order to obtain the button pressed and clear this state, you must call the
GetButtonPressed method.
*/
func (shared *CheckboxInstanceType) IsCheckboxSelected() bool {
	if buttonHistory.layerAlias != "" && buttonHistory.buttonAlias != "" {
		return true
	}
	return false
}

 func AddCheckbox(layerAlias string, checkboxAlias string, checkboxLabel string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool) CheckboxInstanceType {
	memory.AddCheckbox(layerAlias, checkboxAlias, checkboxLabel, styleEntry, xLocation, yLocation, isSelected)
	var checkboxInstance CheckboxInstanceType
	 checkboxInstance.layerAlias = layerAlias
	 checkboxInstance.checkboxAlias = checkboxAlias
	return checkboxInstance
}

/*
DeleteButton allows you to remove a button from a text layer. In addition,
the following information should be noted:

- If you attempt to delete a button which does not exist, then the request
will simply be ignored.
*/
func DeleteCheckbox(layerAlias string, checkboxAlias string) {
	memory.DeleteCheckbox(layerAlias, checkboxAlias)
}

/*
drawButtonsOnLayer allows you to draw all buttons on a given text layer
entry.
*/
func drawCheckboxesOnLayer(layerEntry memory.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for currentKey := range memory.CheckboxMemory[layerAlias] {
		checkboxEntry := memory.GetCheckbox(layerAlias, currentKey)
		drawCheckbox(&layerEntry, currentKey, checkboxEntry.Label, checkboxEntry.StyleEntry, checkboxEntry.XLocation, checkboxEntry.YLocation, checkboxEntry.IsSelected, checkboxEntry.IsEnabled)
	}
}

func drawCheckbox (layerEntry *memory.LayerEntryType, checkboxAlias string, checkboxLabel string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool ) {
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

func updateMouseEventCheckbox() bool {
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
