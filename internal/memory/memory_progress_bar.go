package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
)

type progressBarMemoryType struct {
	sync.Mutex
	Entries map[string]map[string]*types.ProgressBarEntryType
}

var ProgressBar progressBarMemoryType

func InitializeProgressBarMemory() {
	ProgressBar.Entries = make(map[string]map[string]*types.ProgressBarEntryType)
}

func AddProgressBar(layerAlias string, buttonAlias string, buttonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, value int, maxValue int, isBackgroundTransparent bool) {
	ProgressBar.Lock()
	defer func() {
		ProgressBar.Unlock()
	}()
	progressBarEntry := types.NewProgressBarEntry()
	progressBarEntry.StyleEntry = styleEntry
	progressBarEntry.Alias = buttonAlias
	progressBarEntry.Label = buttonLabel
	progressBarEntry.Value = value
	progressBarEntry.MaxValue = maxValue
	progressBarEntry.IsBackgroundTransparent = isBackgroundTransparent
	progressBarEntry.XLocation = xLocation
	progressBarEntry.YLocation = yLocation
	progressBarEntry.Width = width
	progressBarEntry.Height = height
	if ProgressBar.Entries[layerAlias] == nil {
		ProgressBar.Entries[layerAlias] = make(map[string]*types.ProgressBarEntryType)
	}
	ProgressBar.Entries[layerAlias][buttonAlias] = &progressBarEntry
}

func GetProgressBar(layerAlias string, buttonAlias string) *types.ProgressBarEntryType {
	ProgressBar.Lock()
	defer func() {
		ProgressBar.Unlock()
	}()
	if ProgressBar.Entries[layerAlias][buttonAlias] == nil {
		panic(fmt.Sprintf("The requested button with alias '%s' on layer '%s' could not be returned since it does not exist.", buttonAlias, layerAlias))
	}
	return ProgressBar.Entries[layerAlias][buttonAlias]
}

func IsProgressBarExists(layerAlias string, buttonAlias string) bool {
	ProgressBar.Lock()
	defer func() {
		ProgressBar.Unlock()
	}()
	if ProgressBar.Entries[layerAlias][buttonAlias] == nil {
		return false
	}
	return true
}

func DeleteProgressBar(layerAlias string, buttonAlias string) {
	ProgressBar.Lock()
	defer func() {
		ProgressBar.Unlock()
	}()
	delete(ProgressBar.Entries[layerAlias], buttonAlias)
}

func DeleteAllProgressBarsFromLayer(layerAlias string) {
	ProgressBar.Lock()
	defer func() {
		ProgressBar.Unlock()
	}()
	for entryToRemove := range ProgressBar.Entries[layerAlias] {
		delete(ProgressBar.Entries[layerAlias], entryToRemove)
	}
}
