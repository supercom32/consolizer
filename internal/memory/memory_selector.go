package memory

import "fmt"

var SelectorMemory map[string]map[string]*SelectorEntryType

func InitializeSelectorMemory() {
	SelectorMemory = make(map[string]map[string]*SelectorEntryType)
}

func AddSelector(layerAlias string, menuBarAlias string, styleEntry TuiStyleEntryType, selectionEntry SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, itemSelected int) {
	selectorEntry := NewSelectorEntry()
	selectorEntry.StyleEntry = styleEntry
	selectorEntry.SelectionEntry = selectionEntry
	selectorEntry.XLocation = xLocation
	selectorEntry.YLocation = yLocation
	selectorEntry.SelectorHeight = selectorHeight
	selectorEntry.ItemWidth = itemWidth
	selectorEntry.NumberOfColumns = numberOfColumns
	selectorEntry.ViewportPosition = viewportPosition
	selectorEntry.ItemHighlighted = itemSelected
	if SelectorMemory[layerAlias] == nil {
		SelectorMemory[layerAlias] = make(map[string]*SelectorEntryType)
	}
	SelectorMemory[layerAlias][menuBarAlias] = &selectorEntry
}

func DeleteSelector(layerAlias string, menuBarAlias string) {
	delete(SelectorMemory[layerAlias], menuBarAlias)
}

func IsSelectorExists(layerAlias string, menuBarAlias string) bool {
	if _, isExist := SelectorMemory[layerAlias][menuBarAlias]; isExist {
		return true
	}
	return false
}

func GetSelector(layerAlias string, selectorAlias string) *SelectorEntryType {
	if !IsSelectorExists(layerAlias, selectorAlias) {
		panic(fmt.Sprintf("The selector '%s' under layer '%s' could not be obtained since it does not exist!", selectorAlias,  layerAlias))
	}
	return SelectorMemory[layerAlias][selectorAlias]
}