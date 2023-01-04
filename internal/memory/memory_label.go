package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
)

type labelMemoryType struct {
	sync.Mutex
	Entries map[string]map[string]*types.LabelEntryType
}

var Label labelMemoryType

func InitializeLabelMemory() {
	Label.Entries = make(map[string]map[string]*types.LabelEntryType)
}

func AddLabel(layerAlias string, labelAlias string, labelValue string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int) {
	ProgressBar.Lock()
	defer func() {
		ProgressBar.Unlock()
	}()
	labelEntry := types.NewLabelEntry()
	labelEntry.StyleEntry = styleEntry
	labelEntry.Alias = labelAlias
	labelEntry.Value = labelValue
	labelEntry.XLocation = xLocation
	labelEntry.YLocation = yLocation
	labelEntry.Width = width
	if Label.Entries[layerAlias] == nil {
		Label.Entries[layerAlias] = make(map[string]*types.LabelEntryType)
	}
	Label.Entries[layerAlias][labelAlias] = &labelEntry
}

func GetLabel(layerAlias string, labelAlias string) *types.LabelEntryType {
	Label.Lock()
	defer func() {
		Label.Unlock()
	}()
	if Label.Entries[layerAlias][labelAlias] == nil {
		panic(fmt.Sprintf("The requested label with alias '%s' on layer '%s' could not be returned since it does not exist.", labelAlias, layerAlias))
	}
	return Label.Entries[layerAlias][labelAlias]
}

func IsLabelExists(layerAlias string, labelAlias string) bool {
	Label.Lock()
	defer func() {
		Label.Unlock()
	}()
	if Label.Entries[layerAlias][labelAlias] == nil {
		return false
	}
	return true
}

func DeleteLabel(layerAlias string, labelAlias string) {
	Label.Lock()
	defer func() {
		Label.Unlock()
	}()
	delete(Label.Entries[layerAlias], labelAlias)
}

func DeleteAllLabelsFromLayer(layerAlias string) {
	Label.Lock()
	defer func() {
		Label.Unlock()
	}()
	for entryToRemove := range Label.Entries[layerAlias] {
		delete(Label.Entries[layerAlias], entryToRemove)
	}
}
