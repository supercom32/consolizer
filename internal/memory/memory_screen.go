// MARKED: Cleaned
package memory

import (
	"fmt"
	"sort"
	"supercom32.net/consolizer/types"
)

type layerAliasZOrderPair struct {
	Key   string
	Value int
}
type LayerAliasZOrderPairList []layerAliasZOrderPair

var Screen *MemoryManager[types.LayerEntryType]

// ReInitializeScreenMemory initializes the screen memory with a new instance of MemoryManager.
func init() {
	Screen = NewMemoryManager[types.LayerEntryType]() // Initialize MemoryManager
}

func ReInitializeScreenMemory() {
	Screen = NewMemoryManager[types.LayerEntryType]() // Initialize MemoryManager
}

// AddLayer adds a new layer to memory using the MemoryManager.
func AddLayer(layerAlias string, xLocation int, yLocation int, width int, height int, zOrderPriority int, parentAlias string) {
	if width <= 0 {
		panic(fmt.Sprintf("The layer '%s' could not be created since a HotspotWidth of '%d' was specified!", layerAlias, width))
	}
	if height <= 0 {
		panic(fmt.Sprintf("The layer '%s' could not be created since a Length of '%d' was specified!", layerAlias, height))
	}

	layerEntry := types.NewLayerEntry(layerAlias, parentAlias, width, height)
	layerEntry.LayerAlias = layerAlias
	layerEntry.ScreenXLocation = xLocation
	layerEntry.ScreenYLocation = yLocation
	layerEntry.ZOrder = zOrderPriority
	layerEntry.ParentAlias = parentAlias

	if parentAlias != "" {
		parentEntry := Screen.Get(parentAlias)
		if parentEntry != nil {
			parentEntry.IsParent = true
		} else {
			panic(fmt.Sprintf("The layer '%s' could not be created since the parent alias '%s' does not exist!", layerAlias, parentAlias))
		}
	}

	Screen.Add(layerAlias, &layerEntry)
}

// GetLayer retrieves a layer from memory.
func GetLayer(layerAlias string) *types.LayerEntryType {
	layerEntry := Screen.Get(layerAlias)
	if layerEntry == nil {
		panic(fmt.Sprintf("The layer '%s' could not be obtained since it does not exist!", layerAlias))
	}
	return layerEntry
}

// GetNextLayerAlias retrieves the next available layer alias.
func GetNextLayerAlias() string {
	for _, currentEntry := range Screen.memoryItems {
		return currentEntry.LayerAlias
	}
	return ""
}

func DeleteLayer(layerAlias string) {
	screenEntry := Screen.Get(layerAlias)
	if screenEntry == nil {
		panic(fmt.Sprintf("The layer '%s' could not be deleted since it does not exist!", layerAlias))
	}
	layerEntry := Screen.Get(layerAlias)
	parentAlias := layerEntry.ParentAlias
	isParent := layerEntry.IsParent

	DeleteAllButtonsFromLayer(layerAlias)
	DeleteAllCheckboxesFromLayer(layerAlias)
	DeleteAllDropdownsFromLayer(layerAlias)
	DeleteAllProgressBarsFromLayer(layerAlias)
	DeleteAllRadioButtonsFromLayer(layerAlias)
	DeleteAllScrollbarsFromLayer(layerAlias)
	DeleteAllSelectorsFromLayer(layerAlias)
	DeleteAllTextboxesFromLayer(layerAlias)
	DeleteAllTextFieldsFromLayer(layerAlias)
	Screen.Remove(layerAlias)
	if parentAlias != "" {
		parentEntry := Screen.Get(parentAlias)
		if parentEntry != nil {
			if !IsAParent(parentAlias) {
				layerEntry = GetLayer(parentAlias)
				layerEntry.IsParent = false
			}
		}
	}
	if isParent {
		deleteAllChildrenOfParent(layerAlias)
	}
}

func deleteAllChildrenOfParent(parentAlias string) {
	for _, currentValue := range Screen.GetAllEntries() {
		if currentValue.ParentAlias == parentAlias {
			DeleteLayer(currentValue.LayerAlias)
		}
	}
}

func IsAParent(parentAlias string) bool {
	isParent := false
	for _, currentValue := range Screen.GetAllEntries() {
		if currentValue.ParentAlias == parentAlias {
			isParent = true
		}
	}
	return isParent
}

// IsLayerExists checks if a layer exists in memory.
func IsLayerExists(layerAlias string) bool {
	layerEntry := Screen.Get(layerAlias)
	return layerEntry != nil
}

// GetSortedLayerMemoryAliasSlice returns a sorted list of layer aliases based on z-order.
func GetSortedLayerMemoryAliasSlice() LayerAliasZOrderPairList {
	pairList := make(LayerAliasZOrderPairList, len(Screen.memoryItems))
	currentEntry := 0
	for currentKey, currentValue := range Screen.memoryItems {
		pairList[currentEntry].Key = currentKey
		pairList[currentEntry].Value = currentValue.ZOrder
		currentEntry++
	}
	sort.SliceStable(pairList, func(firstIndex, secondIndex int) bool {
		return pairList[firstIndex].Value < pairList[secondIndex].Value
	})
	return pairList
}

// SetHighestZOrderNumber sets the highest z-order number for the given layer.
func SetHighestZOrderNumber(layerAlias string, parentAlias string) {
	if IsLayerExists(layerAlias) {
		highestZOrderNumber := getHighestZOrderNumber(parentAlias)
		for _, currentValue := range Screen.memoryItems {
			if currentValue.ParentAlias == parentAlias && currentValue.ZOrder == highestZOrderNumber {
				currentValue.ZOrder = highestZOrderNumber - 1
				currentValue.IsTopmost = false
			}
		}
		Screen.Get(layerAlias).ZOrder = highestZOrderNumber
		Screen.Get(layerAlias).IsTopmost = true
	}
}

func getHighestZOrderNumber(parentAlias string) int {
	highestZOrderNumber := 0
	for _, currentValue := range Screen.memoryItems {
		if currentValue.ParentAlias == parentAlias && currentValue.ZOrder > highestZOrderNumber {
			highestZOrderNumber = currentValue.ZOrder
		}
	}
	return highestZOrderNumber
}

func GetRootParentLayerAlias(layerAlias string, previousChildAlias string) (string, string) {
	layerEntry := GetLayer(layerAlias)
	if layerEntry.ParentAlias != "" {
		return GetRootParentLayerAlias(layerEntry.ParentAlias, layerAlias)
	}
	return layerAlias, previousChildAlias
}
