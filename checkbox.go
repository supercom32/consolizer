package consolizer

import (
	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/internal/memory"
	"supercom32.net/consolizer/internal/stringformat"
	"supercom32.net/consolizer/types"
)

type CheckboxInstanceType struct {
	layerAlias   string
	controlAlias string
}

type checkboxType struct{}

var Checkbox checkboxType

var Checkboxes = memory.NewControlMemoryManager[types.CheckboxEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

func DeleteCheckbox(layerAlias string, checkboxAlias string) {
	// Use ControlMemoryManager to remove the checkbox entry
	Checkboxes.Remove(layerAlias, checkboxAlias)
}

func DeleteAllCheckboxesFromLayer(layerAlias string) {
	// GetLayer all checkbox entries from the layer
	checkboxes := Checkboxes.GetAllEntries(layerAlias)

	// Loop through all entries and delete them
	for _, checkbox := range checkboxes {
		Checkboxes.Remove(layerAlias, checkbox.Label) // Assuming checkbox.Label is used as the alias
	}
}

func (shared *CheckboxInstanceType) Delete() *CheckboxInstanceType {
	if Checkboxes.IsExists(shared.layerAlias, shared.controlAlias) {
		Checkboxes.Remove(shared.layerAlias, shared.controlAlias)
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
	if Checkboxes.IsExists(shared.layerAlias, shared.controlAlias) {
		checkboxEntry := Checkboxes.Get(shared.layerAlias, shared.controlAlias)
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
	checkboxEntry := types.NewCheckboxEntry()
	checkboxEntry.Alias = checkboxAlias
	checkboxEntry.StyleEntry = styleEntry
	checkboxEntry.Label = checkboxLabel
	checkboxEntry.XLocation = xLocation
	checkboxEntry.YLocation = yLocation
	checkboxEntry.IsSelected = isSelected
	checkboxEntry.IsEnabled = isEnabled
	checkboxEntry.TooltipAlias = stringformat.GetLastSortedUUID()

	// Create associated tooltip (always created but disabled by default)
	Tooltip.Add(layerAlias, checkboxEntry.TooltipAlias, "", styleEntry,
		checkboxEntry.XLocation, checkboxEntry.YLocation,
		len(checkboxLabel)+2, 1,
		checkboxEntry.XLocation, checkboxEntry.YLocation+1,
		len(checkboxLabel)+2, 3,
		false, true, constants.DefaultTooltipHoverTime)

	// Use the ControlMemoryManager to add the checkbox entry
	Checkboxes.Add(layerAlias, checkboxAlias, &checkboxEntry)
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
	Checkboxes.Remove(layerAlias, checkboxAlias)
}

/*
DeleteAllCheckboxesFromLayer allows you to delete all checkboxes on a given text layer.
*/
func (shared *checkboxType) DeleteAllCheckboxesFromLayer(layerAlias string) {
	Checkboxes.RemoveAll(layerAlias)
}

/*
drawCheckboxesOnLayer allows you to draw all checkboxes on a given text layer.
*/
func (shared *checkboxType) drawCheckboxesOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, checkboxEntry := range Checkboxes.GetAllEntries(layerAlias) {
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
updateMouseEvent allows you to update the state of all checkboxes according to the current mouse event state.
In the event that a screen update is required this method returns true.
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

// GetBounds returns the position of the checkbox
func (shared *CheckboxInstanceType) GetBounds() (int, int, int, int) {
	checkboxEntry := Checkboxes.Get(shared.layerAlias, shared.controlAlias)
	if checkboxEntry == nil {
		return 0, 0, 0, 0
	}
	// Checkboxes are typically 1x1 in size
	return checkboxEntry.XLocation, checkboxEntry.YLocation, 1, 1
}

// SetPosition sets the position of the checkbox
func (shared *CheckboxInstanceType) SetPosition(x, y int) *CheckboxInstanceType {
	checkboxEntry := Checkboxes.Get(shared.layerAlias, shared.controlAlias)
	if checkboxEntry != nil {
		checkboxEntry.XLocation = x
		checkboxEntry.YLocation = y
	}
	return shared
}

// SetSize is not applicable for checkboxes as they are fixed size
func (shared *CheckboxInstanceType) SetSize(width, height int) *CheckboxInstanceType {
	// Checkboxes are fixed size (1x1)
	return shared
}

// SetVisible shows or hides the checkbox
func (shared *CheckboxInstanceType) SetVisible(visible bool) *CheckboxInstanceType {
	checkboxEntry := Checkboxes.Get(shared.layerAlias, shared.controlAlias)
	if checkboxEntry != nil {
		checkboxEntry.IsVisible = visible
	}
	return shared
}

// SetStyle sets the visual style of the checkbox
func (shared *CheckboxInstanceType) SetStyle(style types.TuiStyleEntryType) *CheckboxInstanceType {
	checkboxEntry := Checkboxes.Get(shared.layerAlias, shared.controlAlias)
	if checkboxEntry != nil {
		checkboxEntry.StyleEntry = style
	}
	return shared
}

// SetTabIndex sets the tab order of the checkbox
func (shared *CheckboxInstanceType) SetTabIndex(index int) *CheckboxInstanceType {
	checkboxEntry := Checkboxes.Get(shared.layerAlias, shared.controlAlias)
	if checkboxEntry != nil {
		checkboxEntry.TabIndex = index
	}
	return shared
}
