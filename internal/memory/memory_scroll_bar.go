package memory

import (
	"fmt"
)

var ScrollBarMemory map[string]map[string]*ScrollBarEntryType
// var buttonMutex *sync.Mutex <-- Example of a lock. Just lock this item to block and unblock chunks of code.
// buttonMutex = &sync.Mutex{}
// buttonMutex.Lock()

func InitializeScrollBarMemory() {
	ScrollBarMemory = make(map[string]map[string]*ScrollBarEntryType)
}

func AddScrollBar(layerAlias string, scrollBarAlias string, styleEntry TuiStyleEntryType, xLocation int, yLocation int, length int, MaxScrollValue int, ScrollValue int, isHorizontal bool) {
	scrollBarEntry := NewScrollBarEntry()
	scrollBarEntry.StyleEntry = styleEntry
	scrollBarEntry.XLocation = xLocation
	scrollBarEntry.YLocation = yLocation
	scrollBarEntry.Length = length
	scrollBarEntry.MaxScrollValue = MaxScrollValue -1 // Since scroll values are 0 based, we minus 1 from total.
	scrollBarEntry.ScrollValue = ScrollValue
	scrollBarEntry.IsVisible = true
	scrollBarEntry.IsHorizontal = isHorizontal
	if ScrollBarMemory[layerAlias] == nil {
		ScrollBarMemory[layerAlias] = make(map[string]*ScrollBarEntryType)
	}
	ScrollBarMemory[layerAlias][scrollBarAlias] = &scrollBarEntry
}

func GetScrollBar(layerAlias string, scrollBarAlias string) ScrollBarEntryType {
	if ScrollBarMemory[layerAlias][scrollBarAlias] == nil {
		panic(fmt.Sprintf("The requested scroll bar with alias '%s' on layer '%s' could not be returned since it does not exist.", scrollBarAlias, layerAlias))
	}
	return *ScrollBarMemory[layerAlias][scrollBarAlias]
}

func DeleteScrollBar(layerAlias string, scrollBarAlias string) {
	delete(ScrollBarMemory[layerAlias], scrollBarAlias)
}
