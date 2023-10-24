package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
)

type textFieldMemoryType struct {
	sync.Mutex
	Entries map[string]map[string]*types.TextFieldEntryType
}

var TextField textFieldMemoryType

func InitializeTextFieldMemory() {
	TextField.Entries = make(map[string]map[string]*types.TextFieldEntryType)
}

func AddTextField(layerAlias string, textFieldAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, maxLengthAllowed int, IsPasswordProtected bool, defaultValue string, isEnabled bool) {
	TextField.Lock()
	defer func() {
		TextField.Unlock()
	}()
	textFieldEntry := types.NewTextFieldEntry()
	textFieldEntry.StyleEntry = styleEntry
	textFieldEntry.XLocation = xLocation
	textFieldEntry.YLocation = yLocation
	textFieldEntry.Width = width
	textFieldEntry.MaxLengthAllowed = maxLengthAllowed
	textFieldEntry.IsPasswordProtected = IsPasswordProtected
	textFieldEntry.CurrentValue = []rune(defaultValue)
	textFieldEntry.DefaultValue = defaultValue
	textFieldEntry.IsEnabled = isEnabled
	if TextField.Entries[layerAlias] == nil {
		TextField.Entries[layerAlias] = make(map[string]*types.TextFieldEntryType)
	}
	TextField.Entries[layerAlias][textFieldAlias] = &textFieldEntry
}

func DeleteTextField(layerAlias string, textFieldAlias string) {
	TextField.Lock()
	defer func() {
		TextField.Unlock()
	}()
	delete(TextField.Entries[layerAlias], textFieldAlias)
}

func DeleteAllTextFieldsFromLayer(layerAlias string) {
	TextField.Lock()
	defer func() {
		TextField.Unlock()
	}()
	for entryToRemove := range TextField.Entries[layerAlias] {
		delete(TextField.Entries[layerAlias], entryToRemove)
	}
}

func IsTextFieldExists(layerAlias string, textFieldAlias string) bool {
	TextField.Lock()
	defer func() {
		TextField.Unlock()
	}()
	if _, isExist := TextField.Entries[layerAlias][textFieldAlias]; isExist {
		return true
	}
	return false
}

func GetTextField(layerAlias string, textFieldAlias string) *types.TextFieldEntryType {
	TextField.Lock()
	defer func() {
		TextField.Unlock()
	}()
	if _, isExist := TextField.Entries[layerAlias][textFieldAlias]; !isExist {
		panic(fmt.Sprintf("The text field '%s' under layer '%s' could not be obtained since it does not exist!", textFieldAlias, layerAlias))
	}
	return TextField.Entries[layerAlias][textFieldAlias]
}
