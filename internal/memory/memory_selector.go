package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
)

type selectorMemoryType struct {
	sync.Mutex
	Entries map[string]map[string]*types.SelectorEntryType
}

var Selector selectorMemoryType

func InitializeSelectorMemory() {
	Selector.Entries = make(map[string]map[string]*types.SelectorEntryType)
}

func AddSelector(layerAlias string, selectorAlias string, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, itemSelected int, isBorderDrawn bool) {
	Selector.Lock()
	defer func() {
		Selector.Unlock()
	}()
	selectorEntry := types.NewSelectorEntry()
	selectorEntry.StyleEntry = styleEntry
	selectorEntry.SelectionEntry = selectionEntry
	selectorEntry.XLocation = xLocation
	selectorEntry.YLocation = yLocation
	selectorEntry.SelectorHeight = selectorHeight
	selectorEntry.ItemWidth = itemWidth
	selectorEntry.NumberOfColumns = numberOfColumns
	selectorEntry.ViewportPosition = viewportPosition
	selectorEntry.ItemHighlighted = itemSelected
	selectorEntry.IsBorderDrawn = isBorderDrawn
	selectorEntry.IsVisible = true
	if Selector.Entries[layerAlias] == nil {
		Selector.Entries[layerAlias] = make(map[string]*types.SelectorEntryType)
	}
	Selector.Entries[layerAlias][selectorAlias] = &selectorEntry
}

func DeleteSelector(layerAlias string, selectorAlias string) {
	Selector.Lock()
	defer func() {
		Selector.Unlock()
	}()
	delete(Selector.Entries[layerAlias], selectorAlias)
}

func DeleteAllSelectorsFromLayer(layerAlias string) {
	Selector.Lock()
	defer func() {
		Selector.Unlock()
	}()
	for entryToRemove := range Selector.Entries[layerAlias] {
		delete(Selector.Entries[layerAlias], entryToRemove)
	}
}

func IsSelectorExists(layerAlias string, selectorAlias string) bool {
	Selector.Lock()
	defer func() {
		Selector.Unlock()
	}()
	if _, isExist := Selector.Entries[layerAlias][selectorAlias]; isExist {
		return true
	}
	return false
}

func GetSelector(layerAlias string, selectorAlias string) *types.SelectorEntryType {
	Selector.Lock()
	defer func() {
		Selector.Unlock()
	}()
	if _, isExist := Selector.Entries[layerAlias][selectorAlias]; !isExist {
		panic(fmt.Sprintf("The selector '%s' under layer '%s' could not be obtained since it does not exist!", selectorAlias, layerAlias))
	}
	return Selector.Entries[layerAlias][selectorAlias]
}
