package memory

import (
	"fmt"
	"supercom32.net/consolizer/types"
)

var Checkboxes = NewControlMemoryManager[types.CheckboxEntryType]()

func AddCheckbox(layerAlias string, checkboxAlias string, checkboxLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool) {
	checkboxEntry := types.NewCheckboxEntry()
	checkboxEntry.Alias = checkboxAlias
	checkboxEntry.StyleEntry = styleEntry
	checkboxEntry.Label = checkboxLabel
	checkboxEntry.XLocation = xLocation
	checkboxEntry.YLocation = yLocation
	checkboxEntry.IsSelected = isSelected
	checkboxEntry.IsEnabled = isEnabled

	// Use the ControlMemoryManager to add the checkbox entry
	Checkboxes.Add(layerAlias, checkboxAlias, &checkboxEntry)
}

func GetCheckbox(layerAlias string, checkboxAlias string) *types.CheckboxEntryType {
	// Get the checkbox entry using the ControlMemoryManager
	checkboxEntry := Checkboxes.Get(layerAlias, checkboxAlias)
	if checkboxEntry == nil {
		panic(fmt.Sprintf("The requested Checkbox with alias '%s' on layer '%s' could not be returned since it does not exist.", checkboxAlias, layerAlias))
	}
	return checkboxEntry
}

func IsCheckboxExists(layerAlias string, checkboxAlias string) bool {
	// Use ControlMemoryManager to check if the checkbox exists
	return Checkboxes.Get(layerAlias, checkboxAlias) != nil
}

func DeleteCheckbox(layerAlias string, checkboxAlias string) {
	// Use ControlMemoryManager to remove the checkbox entry
	Checkboxes.Remove(layerAlias, checkboxAlias)
}

func DeleteAllCheckboxesFromLayer(layerAlias string) {
	// Get all checkbox entries from the layer
	checkboxes := Checkboxes.GetAllEntries(layerAlias)

	// Loop through all entries and delete them
	for _, checkbox := range checkboxes {
		Checkboxes.Remove(layerAlias, checkbox.Label) // Assuming checkbox.Label is used as the alias
	}
}
