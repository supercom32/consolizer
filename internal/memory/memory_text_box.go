package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
)

type textboxMemoryType struct {
	sync.Mutex
	Entries map[string]map[string]*types.TextboxEntryType
}

var Textbox textboxMemoryType

// var buttonMutex *sync.Mutex <-- Example of a lock. Just lock this item to block and unblock chunks of code.
// buttonMutex = &sync.Mutex{}
// buttonMutex.Lock()

func InitializeTextboxMemory() {
	Textbox.Entries = make(map[string]map[string]*types.TextboxEntryType)
}

func AddTextbox(layerAlias string, textboxAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isBorderDrawn bool) {
	Textbox.Lock()
	defer func() {
		Textbox.Unlock()
	}()
	textboxEntry := types.NewTexboxEntry()
	textboxEntry.StyleEntry = styleEntry
	textboxEntry.XLocation = xLocation
	textboxEntry.YLocation = yLocation
	textboxEntry.Width = width
	textboxEntry.Height = height
	textboxEntry.IsBorderDrawn = isBorderDrawn
	if Textbox.Entries[layerAlias] == nil {
		Textbox.Entries[layerAlias] = make(map[string]*types.TextboxEntryType)
	}
	Textbox.Entries[layerAlias][textboxAlias] = &textboxEntry
}

func GetTextbox(layerAlias string, textboxAlias string) *types.TextboxEntryType {
	Textbox.Lock()
	defer func() {
		Textbox.Unlock()
	}()
	if Textbox.Entries[layerAlias][textboxAlias] == nil {
		panic(fmt.Sprintf("The requested text with alias '%s' on layer '%s' could not be returned since it does not exist.", textboxAlias, layerAlias))
	}
	return Textbox.Entries[layerAlias][textboxAlias]
}

func IsTextboxExists(layerAlias string, textboxAlias string) bool {
	Textbox.Lock()
	defer func() {
		Textbox.Unlock()
	}()
	if Textbox.Entries[layerAlias][textboxAlias] == nil {
		return false
	}
	return true
}

func DeleteTextbox(layerAlias string, textboxAlias string) {
	Textbox.Lock()
	defer func() {
		Textbox.Unlock()
	}()
	delete(Textbox.Entries[layerAlias], textboxAlias)
}

func DeleteAllTextboxesFromLayer(layerAlias string) {
	Textbox.Lock()
	defer func() {
		Textbox.Unlock()
	}()
	for entryToRemove := range Textbox.Entries[layerAlias] {
		delete(Textbox.Entries[layerAlias], entryToRemove)
	}
}
