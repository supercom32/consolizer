package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
)

type checkboxMemoryType struct {
	sync.Mutex
	Entries map[string]map[string]*types.CheckboxEntryType
}

var Checkbox checkboxMemoryType

// var buttonMutex *sync.Mutex <-- Example of a lock. Just lock this item to block and unblock chunks of code.
// buttonMutex = &sync.Mutex{}
// buttonMutex.Lock()

func InitializeCheckboxMemory() {
	Checkbox.Entries = make(map[string]map[string]*types.CheckboxEntryType)
}

func AddCheckbox(layerAlias string, checkboxAlias string, checkboxLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool) {
	Checkbox.Lock()
	defer func() {
		Checkbox.Unlock()
	}()
	checkboxEntry := types.NewCheckboxEntry()
	checkboxEntry.StyleEntry = styleEntry
	checkboxEntry.Label = checkboxLabel
	checkboxEntry.XLocation = xLocation
	checkboxEntry.YLocation = yLocation
	checkboxEntry.IsSelected = isSelected
	checkboxEntry.IsEnabled = isEnabled
	if Checkbox.Entries[layerAlias] == nil {
		Checkbox.Entries[layerAlias] = make(map[string]*types.CheckboxEntryType)
	}
	Checkbox.Entries[layerAlias][checkboxAlias] = &checkboxEntry
}

func GetCheckbox(layerAlias string, checkboxAlias string) *types.CheckboxEntryType {
	Checkbox.Lock()
	defer func() {
		Checkbox.Unlock()
	}()
	if Checkbox.Entries[layerAlias][checkboxAlias] == nil {
		panic(fmt.Sprintf("The requested Checkbox with alias '%s' on layer '%s' could not be returned since it does not exist.", checkboxAlias, layerAlias))
	}
	return Checkbox.Entries[layerAlias][checkboxAlias]
}

func IsCheckboxExists(layerAlias string, checkboxAlias string) bool {
	Checkbox.Lock()
	defer func() {
		Checkbox.Unlock()
	}()
	if Checkbox.Entries[layerAlias][checkboxAlias] == nil {
		return false
	}
	return true
}

func DeleteCheckbox(layerAlias string, checkboxAlias string) {
	Checkbox.Lock()
	defer func() {
		Checkbox.Unlock()
	}()
	delete(Checkbox.Entries[layerAlias], checkboxAlias)
}

func DeleteAllCheckboxesFromLayer(layerAlias string) {
	Checkbox.Lock()
	defer func() {
		Checkbox.Unlock()
	}()
	for entryToRemove := range Checkbox.Entries[layerAlias] {
		delete(Checkbox.Entries[layerAlias], entryToRemove)
	}
}
