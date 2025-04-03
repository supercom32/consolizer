package memory

import (
	"fmt"
	"supercom32.net/consolizer/types"
)

var RadioButtons = NewControlMemoryManager[types.RadioButtonEntryType]()

func AddRadioButton(layerAlias string, radioButtonAlias string, radioButtonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, groupId int, isSelected bool) {
	radioButtonEntry := types.NewRadioButtonEntry()
	radioButtonEntry.Alias = radioButtonAlias
	radioButtonEntry.StyleEntry = styleEntry
	radioButtonEntry.Label = radioButtonLabel
	radioButtonEntry.XLocation = xLocation
	radioButtonEntry.YLocation = yLocation
	radioButtonEntry.GroupId = groupId
	radioButtonEntry.IsSelected = isSelected

	// Use the ControlMemoryManager to add the radio button entry
	RadioButtons.Add(layerAlias, radioButtonAlias, &radioButtonEntry)
}

func GetRadioButton(layerAlias string, radioButtonAlias string) *types.RadioButtonEntryType {
	// Get the radio button entry using ControlMemoryManager
	radioButtonEntry := RadioButtons.Get(layerAlias, radioButtonAlias)
	if radioButtonEntry == nil {
		panic(fmt.Sprintf("The requested radio button with alias '%s' on layer '%s' could not be returned since it does not exist.", radioButtonAlias, layerAlias))
	}
	return radioButtonEntry
}

func IsRadioButtonExists(layerAlias string, radioButtonAlias string) bool {
	// Use ControlMemoryManager to check if the radio button exists
	return RadioButtons.Get(layerAlias, radioButtonAlias) != nil
}

func DeleteRadioButton(layerAlias string, radioButtonAlias string) {
	// Use ControlMemoryManager to remove the radio button entry
	RadioButtons.Remove(layerAlias, radioButtonAlias)
}

func DeleteAllRadioButtonsFromLayer(layerAlias string) {
	// Get all radio button entries from the layer
	radioButtons := RadioButtons.GetAllEntries(layerAlias)

	// Loop through all entries and delete them
	for _, radioButton := range radioButtons {
		RadioButtons.Remove(layerAlias, radioButton.Alias) // Assuming radioButton.Alias is used as the alias
	}
}
