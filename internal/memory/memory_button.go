package memory

import (
	"fmt"
)

var ButtonMemory map[string]map[string]*ButtonEntryType
// var buttonMutex *sync.Mutex <-- Example of a lock. Just lock this item to block and unblock chunks of code.
// buttonMutex = &sync.Mutex{}
// buttonMutex.Lock()

func InitializeButtonMemory() {
	ButtonMemory = make(map[string]map[string]*ButtonEntryType)
}

func AddButton(layerAlias string, buttonAlias string, buttonLabel string, styleEntry TuiStyleEntryType, xLocation int, yLocation int, width int, height int) {
	buttonEntry := NewButtonEntry()
	buttonEntry.StyleEntry = styleEntry
	buttonEntry.ButtonAlias = buttonAlias
	buttonEntry.ButtonLabel = buttonLabel
	buttonEntry.XLocation = xLocation
	buttonEntry.YLocation = yLocation
	buttonEntry.Width = width
	buttonEntry.Height = height
	if ButtonMemory[layerAlias] == nil {
		ButtonMemory[layerAlias] = make(map[string]*ButtonEntryType)
	}
	ButtonMemory[layerAlias][buttonAlias] = &buttonEntry
}

func GetButton(layerAlias string, buttonAlias string) *ButtonEntryType {
	if ButtonMemory[layerAlias][buttonAlias] == nil {
		panic(fmt.Sprintf("The requested button with alias '%s' on layer '%s' could not be returned since it does not exist.", buttonAlias, layerAlias))
	}
	return ButtonMemory[layerAlias][buttonAlias]
}

func IsButtonExists(layerAlias string, buttonAlias string) bool {
	if ButtonMemory[layerAlias][buttonAlias] == nil {
		return false
	}
	return true
}

func DeleteButton(layerAlias string, buttonAlias string) {
	delete(ButtonMemory[layerAlias], buttonAlias)
}
