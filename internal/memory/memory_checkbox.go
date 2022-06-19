package memory

import (
	"fmt"
)

var CheckboxMemory map[string]map[string]*CheckboxEntryType
// var buttonMutex *sync.Mutex <-- Example of a lock. Just lock this item to block and unblock chunks of code.
// buttonMutex = &sync.Mutex{}
// buttonMutex.Lock()

func InitializeCheckboxMemory() {
	CheckboxMemory = make(map[string]map[string]*CheckboxEntryType)
}

func AddCheckbox(layerAlias string, checkboxAlias string, checkboxLabel string, styleEntry TuiStyleEntryType, xLocation int, yLocation int, isSelected bool) {
	checkboxEntry := NewCheckboxEntry()
	checkboxEntry.StyleEntry = styleEntry
	checkboxEntry.Label = checkboxLabel
	checkboxEntry.XLocation = xLocation
	checkboxEntry.YLocation = yLocation
	checkboxEntry.IsSelected = isSelected
	if CheckboxMemory[layerAlias] == nil {
		CheckboxMemory[layerAlias] = make(map[string]*CheckboxEntryType)
	}
	CheckboxMemory[layerAlias][checkboxAlias] = &checkboxEntry
}

func GetCheckbox(layerAlias string, checkboxAlias string) *CheckboxEntryType {
	if CheckboxMemory[layerAlias][checkboxAlias] == nil {
		panic(fmt.Sprintf("The requested checkbox with alias '%s' on layer '%s' could not be returned since it does not exist.", checkboxAlias, layerAlias))
	}
	return CheckboxMemory[layerAlias][checkboxAlias]
}

func IsCheckboxExists(layerAlias string, checkboxAlias string) bool {
	if CheckboxMemory[layerAlias][checkboxAlias] == nil {
		return false
	}
	return true
}

func DeleteCheckbox(layerAlias string, checkboxAlias string) {
	delete(CheckboxMemory[layerAlias], checkboxAlias)
}
