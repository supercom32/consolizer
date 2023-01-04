package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
)

type scrollBarMemoryType struct {
	sync.Mutex
	Entries map[string]map[string]*types.ScrollbarEntryType
}

// var buttonMutex *sync.Mutex <-- Example of a lock. Just lock this item to block and unblock chunks of code.
// buttonMutex = &sync.Mutex{}
// buttonMutex.Lock()

var ScrollBar scrollBarMemoryType

func InitializeScrollbarMemory() {
	ScrollBar.Entries = make(map[string]map[string]*types.ScrollbarEntryType)
}

func AddScrollbar(layerAlias string, scrollBarAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, length int, maxScrollValue int, scrollValue int, ScrollIncrement int, isHorizontal bool) {
	ScrollBar.Lock()
	defer func() {
		ScrollBar.Unlock()
	}()
	scrollbarEntry := types.NewScrollbarEntry()
	scrollbarEntry.StyleEntry = styleEntry
	scrollbarEntry.XLocation = xLocation
	scrollbarEntry.YLocation = yLocation
	scrollbarEntry.Length = length
	scrollbarEntry.MaxScrollValue = maxScrollValue - 1 // Since scroll values are 0 based, we minus 1 from total.
	scrollbarEntry.ScrollValue = scrollValue
	scrollbarEntry.IsVisible = true
	scrollbarEntry.IsEnabled = true
	scrollbarEntry.IsHorizontal = isHorizontal
	scrollbarEntry.ScrollIncrement = ScrollIncrement
	if ScrollBar.Entries[layerAlias] == nil {
		ScrollBar.Entries[layerAlias] = make(map[string]*types.ScrollbarEntryType)
	}
	ScrollBar.Entries[layerAlias][scrollBarAlias] = &scrollbarEntry
}

func GetScrollbar(layerAlias string, scrollBarAlias string) *types.ScrollbarEntryType {
	ScrollBar.Lock()
	defer func() {
		ScrollBar.Unlock()
	}()
	if ScrollBar.Entries[layerAlias][scrollBarAlias] == nil {
		panic(fmt.Sprintf("The requested scroll bar with alias '%s' on layer '%s' could not be returned since it does not exist.", scrollBarAlias, layerAlias))
	}
	return ScrollBar.Entries[layerAlias][scrollBarAlias]
}

func IsScrollbarExists(layerAlias string, scrollBarAlias string) bool {
	ScrollBar.Lock()
	defer func() {
		ScrollBar.Unlock()
	}()
	if ScrollBar.Entries[layerAlias][scrollBarAlias] == nil {
		return false
	}
	return true
}

func DeleteScrollbar(layerAlias string, scrollBarAlias string) {
	ScrollBar.Lock()
	defer func() {
		ScrollBar.Unlock()
	}()
	delete(ScrollBar.Entries[layerAlias], scrollBarAlias)
}

func DeleteAllScrollbarsFromLayer(layerAlias string) {
	ScrollBar.Lock()
	defer func() {
		ScrollBar.Unlock()
	}()
	for entryToRemove := range ScrollBar.Entries[layerAlias] {
		delete(ScrollBar.Entries[layerAlias], entryToRemove)
	}
}
