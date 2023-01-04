package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
)

type buttonMemoryType struct {
	sync.Mutex
	Entries map[string]map[string]*types.ButtonEntryType
}

var Button buttonMemoryType

// var buttonMutex *sync.Mutex <-- Example of a lock. Just lock this item to block and unblock chunks of code.
// buttonMutex = &sync.Mutex{}
// buttonMutex.Lock()

func InitializeButtonMemory() {
	Button.Entries = make(map[string]map[string]*types.ButtonEntryType)
}

func AddButton(layerAlias string, buttonAlias string, buttonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int) {
	Button.Lock()
	defer func() {
		Button.Unlock()
	}()
	buttonEntry := types.NewButtonEntry()
	buttonEntry.StyleEntry = styleEntry
	buttonEntry.ButtonAlias = buttonAlias
	buttonEntry.ButtonLabel = buttonLabel
	buttonEntry.XLocation = xLocation
	buttonEntry.YLocation = yLocation
	buttonEntry.IsEnabled = true
	buttonEntry.Width = width
	buttonEntry.Height = height
	if Button.Entries[layerAlias] == nil {
		Button.Entries[layerAlias] = make(map[string]*types.ButtonEntryType)
	}
	Button.Entries[layerAlias][buttonAlias] = &buttonEntry
}

func GetButton(layerAlias string, buttonAlias string) *types.ButtonEntryType {
	Button.Lock()
	defer func() {
		Button.Unlock()
	}()
	if Button.Entries[layerAlias][buttonAlias] == nil {
		panic(fmt.Sprintf("The requested button with alias '%s' on layer '%s' could not be returned since it does not exist.", buttonAlias, layerAlias))
	}
	return Button.Entries[layerAlias][buttonAlias]
}

func IsButtonExists(layerAlias string, buttonAlias string) bool {
	/*	Button.Lock()
		defer func() {
			Button.Unlock()
		}()*/
	if Button.Entries[layerAlias][buttonAlias] == nil {
		return false
	}
	return true
}

func DeleteButton(layerAlias string, buttonAlias string) {
	/*	Button.Lock()
		defer func() {
			Button.Unlock()
		}()*/
	delete(Button.Entries[layerAlias], buttonAlias)
}

func DeleteAllButtonsFromLayer(layerAlias string) {
	/*	Button.Lock()
		defer func() {
			Button.Unlock()
		}()*/
	for buttonEntry := range Button.Entries[layerAlias] {
		delete(Button.Entries[layerAlias], buttonEntry)
	}
}
