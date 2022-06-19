package memory

import (
	"fmt"
)

var TextboxMemory map[string]map[string]*TextboxEntryType
// var buttonMutex *sync.Mutex <-- Example of a lock. Just lock this item to block and unblock chunks of code.
// buttonMutex = &sync.Mutex{}
// buttonMutex.Lock()

func InitializeTextboxMemory() {
	TextboxMemory = make(map[string]map[string]*TextboxEntryType)
}

func AddTextbox(layerAlias string, textboxAlias string, styleEntry TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isBorderDrawn bool) {
	textboxEntry := NewTexboxEntry()
	textboxEntry.StyleEntry = styleEntry
	textboxEntry.XLocation = xLocation
	textboxEntry.YLocation = yLocation
	textboxEntry.Width = width
	textboxEntry.Height = height
	textboxEntry.IsBorderDrawn = isBorderDrawn
	if TextboxMemory[layerAlias] == nil {
		TextboxMemory[layerAlias] = make(map[string]*TextboxEntryType)
	}
	TextboxMemory[layerAlias][textboxAlias] = &textboxEntry
}

func GetTextbox(layerAlias string, textboxAlias string) *TextboxEntryType {
	if TextboxMemory[layerAlias][textboxAlias] == nil {
		panic(fmt.Sprintf("The requested text with alias '%s' on layer '%s' could not be returned since it does not exist.", textboxAlias, layerAlias))
	}
	return TextboxMemory[layerAlias][textboxAlias]
}

func IsTextboxExists(layerAlias string, textboxAlias string) bool {
	if TextboxMemory[layerAlias][textboxAlias] == nil {
		return false
	}
	return true
}

func DeleteTextbox(layerAlias string, textboxAlias string) {
	delete(TextboxMemory[layerAlias], textboxAlias)
}
