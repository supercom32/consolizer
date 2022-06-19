package memory

import (
	"fmt"
)

var RadioButtonMemory map[string]map[string]*RadioButtonEntryType

func InitializeRadioButtonMemory() {
	RadioButtonMemory = make(map[string]map[string]*RadioButtonEntryType)
}

func AddRadioButton(layerAlias string, radioButtonAlias string, radioButtonLabel string, styleEntry TuiStyleEntryType, xLocation int, yLocation int, groupId int, isSelected bool) {
	radioButtonEntry := NewRadioButtonEntry()
	radioButtonEntry.StyleEntry = styleEntry
	radioButtonEntry.Label = radioButtonLabel
	radioButtonEntry.XLocation = xLocation
	radioButtonEntry.YLocation = yLocation
	radioButtonEntry.GroupId = groupId
	radioButtonEntry.IsSelected = isSelected
	if RadioButtonMemory[layerAlias] == nil {
		RadioButtonMemory[layerAlias] = make(map[string]*RadioButtonEntryType)
	}
	RadioButtonMemory[layerAlias][radioButtonAlias] = &radioButtonEntry
}

func GetRadioButton(layerAlias string, radioButtonAlias string) *RadioButtonEntryType {
	if RadioButtonMemory[layerAlias][radioButtonAlias] == nil {
		panic(fmt.Sprintf("The requested radio button with alias '%s' on layer '%s' could not be returned since it does not exist.", radioButtonAlias, layerAlias))
	}
	return RadioButtonMemory[layerAlias][radioButtonAlias]
}

func IsRadioButtonExists(layerAlias string, radioButtonAlias string) bool {
	if RadioButtonMemory[layerAlias][radioButtonAlias] == nil {
		return false
	}
	return true
}

func DeleteRadioButton(layerAlias string, radioButtonAlias string) {
	delete(RadioButtonMemory[layerAlias], radioButtonAlias)
}
