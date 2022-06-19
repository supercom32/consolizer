package memory

import (
	"fmt"
)

var ScrollbarMemory map[string]map[string]*ScrollbarEntryType
// var buttonMutex *sync.Mutex <-- Example of a lock. Just lock this item to block and unblock chunks of code.
// buttonMutex = &sync.Mutex{}
// buttonMutex.Lock()

func InitializeScrollbarMemory() {
	ScrollbarMemory = make(map[string]map[string]*ScrollbarEntryType)
}

func AddScrollbar(layerAlias string, scrollBarAlias string, styleEntry TuiStyleEntryType, xLocation int, yLocation int, length int, maxScrollValue int, scrollValue int, ScrollIncrement int, isHorizontal bool) {
	scrollbarEntry := NewScrollbarEntry()
	scrollbarEntry.StyleEntry = styleEntry
	scrollbarEntry.XLocation = xLocation
	scrollbarEntry.YLocation = yLocation
	scrollbarEntry.Length = length
	scrollbarEntry.MaxScrollValue = maxScrollValue -1 // Since scroll values are 0 based, we minus 1 from total.
	scrollbarEntry.ScrollValue = scrollValue
	scrollbarEntry.IsVisible = true
	scrollbarEntry.IsEnabled = true
	scrollbarEntry.IsHorizontal = isHorizontal
	scrollbarEntry.ScrollIncrement = ScrollIncrement
	if ScrollbarMemory[layerAlias] == nil {
		ScrollbarMemory[layerAlias] = make(map[string]*ScrollbarEntryType)
	}
	ScrollbarMemory[layerAlias][scrollBarAlias] = &scrollbarEntry
}

func GetScrollbar(layerAlias string, scrollBarAlias string) *ScrollbarEntryType {
	if ScrollbarMemory[layerAlias][scrollBarAlias] == nil {
		panic(fmt.Sprintf("The requested scroll bar with alias '%s' on layer '%s' could not be returned since it does not exist.", scrollBarAlias, layerAlias))
	}
	return ScrollbarMemory[layerAlias][scrollBarAlias]
}

func IsScrollbarExists(layerAlias string, scrollBarAlias string) bool {
	if ScrollbarMemory[layerAlias][scrollBarAlias] == nil {
		return false
	}
	return true
}

func DeleteScrollbar(layerAlias string, scrollBarAlias string) {
	delete(ScrollbarMemory[layerAlias], scrollBarAlias)
}
