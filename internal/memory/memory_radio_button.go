package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
)

type radioButtonMemoryType struct {
	sync.Mutex
	Entries map[string]map[string]*types.RadioButtonEntryType
}

var RadioButton radioButtonMemoryType

func InitializeRadioButtonMemory() {
	RadioButton.Entries = make(map[string]map[string]*types.RadioButtonEntryType)
}

func AddRadioButton(layerAlias string, radioButtonAlias string, radioButtonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, groupId int, isSelected bool) {
	RadioButton.Lock()
	defer func() {
		RadioButton.Unlock()
	}()
	radioButtonEntry := types.NewRadioButtonEntry()
	radioButtonEntry.StyleEntry = styleEntry
	radioButtonEntry.Label = radioButtonLabel
	radioButtonEntry.XLocation = xLocation
	radioButtonEntry.YLocation = yLocation
	radioButtonEntry.GroupId = groupId
	radioButtonEntry.IsSelected = isSelected
	if RadioButton.Entries[layerAlias] == nil {
		RadioButton.Entries[layerAlias] = make(map[string]*types.RadioButtonEntryType)
	}
	RadioButton.Entries[layerAlias][radioButtonAlias] = &radioButtonEntry
}

func GetRadioButton(layerAlias string, radioButtonAlias string) *types.RadioButtonEntryType {
	RadioButton.Lock()
	defer func() {
		RadioButton.Unlock()
	}()
	if RadioButton.Entries[layerAlias][radioButtonAlias] == nil {
		panic(fmt.Sprintf("The requested radio button with alias '%s' on layer '%s' could not be returned since it does not exist.", radioButtonAlias, layerAlias))
	}
	return RadioButton.Entries[layerAlias][radioButtonAlias]
}

func IsRadioButtonExists(layerAlias string, radioButtonAlias string) bool {
	RadioButton.Lock()
	defer func() {
		RadioButton.Unlock()
	}()
	if RadioButton.Entries[layerAlias][radioButtonAlias] == nil {
		return false
	}
	return true
}

func DeleteRadioButton(layerAlias string, radioButtonAlias string) {
	RadioButton.Lock()
	defer func() {
		RadioButton.Unlock()
	}()
	delete(RadioButton.Entries[layerAlias], radioButtonAlias)
}

func DeleteAllRadioButtonsFromLayer(layerAlias string) {
	RadioButton.Lock()
	defer func() {
		RadioButton.Unlock()
	}()
	for entryToRemove := range RadioButton.Entries[layerAlias] {
		delete(RadioButton.Entries[layerAlias], entryToRemove)
	}
}
