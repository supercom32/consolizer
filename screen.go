// MARKED: Cleaned
package consolizer

import (
	"fmt"
	"sort"
	"supercom32.net/consolizer/internal/memory"
	"supercom32.net/consolizer/types"
)

type layerAliasZOrderPair struct {
	Key   string
	Value int
}
type LayerAliasZOrderPairList []layerAliasZOrderPair

type ScreenInstanceType struct {
	ScreenMemory *memory.MemoryManager[types.LayerEntryType]
}

var Screen ScreenInstanceType

// ReInitializeScreenMemory initializes the screen memory with a new instance of MemoryManager.
func init() {
	Screen.ReInitializeScreenMemory()
}

func (shared *ScreenInstanceType) ReInitializeScreenMemory() {
	shared.ScreenMemory = memory.NewMemoryManager[types.LayerEntryType]() // Initialize MemoryManager
}

// AddLayer adds a new layer to memory using the MemoryManager.
func (shared *ScreenInstanceType) AddLayer(layerAlias string, xLocation int, yLocation int, width int, height int, zOrderPriority int, parentAlias string) {
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
		parentEntry := shared.ScreenMemory.Get(parentAlias)
		if parentEntry != nil {
			parentEntry.IsParent = true
		} else {
			panic(fmt.Sprintf("The layer '%s' could not be created since the parent alias '%s' does not exist!", layerAlias, parentAlias))
		}
	}

	shared.ScreenMemory.Add(layerAlias, &layerEntry)
}

// GetLayer retrieves a layer from memory.
func (shared *ScreenInstanceType) GetLayer(layerAlias string) *types.LayerEntryType {
	layerEntry := shared.ScreenMemory.Get(layerAlias)
	if layerEntry == nil {
		panic(fmt.Sprintf("The layer '%s' could not be obtained since it does not exist!", layerAlias))
	}
	return layerEntry
}

// GetNextLayerAlias retrieves the next available layer alias.
func (shared *ScreenInstanceType) GetNextLayerAlias() string {
	for _, currentEntry := range shared.ScreenMemory.GetAllEntries() {
		return currentEntry.LayerAlias
	}
	return ""
}

func (shared *ScreenInstanceType) DeleteLayer(layerAlias string) {
	screenEntry := shared.ScreenMemory.Get(layerAlias)
	if screenEntry == nil {
		panic(fmt.Sprintf("The layer '%s' could not be deleted since it does not exist!", layerAlias))
	}
	layerEntry := Screen.GetLayer(layerAlias)
	parentAlias := layerEntry.ParentAlias
	isParent := layerEntry.IsParent
	Labels.RemoveAll(layerAlias)
	Buttons.RemoveAll(layerAlias)
	Checkboxes.RemoveAll(layerAlias)
	Dropdowns.RemoveAll(layerAlias)
	ProgressBars.RemoveAll(layerAlias)
	RadioButtons.RemoveAll(layerAlias)
	ScrollBars.RemoveAll(layerAlias)
	Selectors.RemoveAll(layerAlias)
	Textboxes.RemoveAll(layerAlias)
	TextFields.RemoveAll(layerAlias)
	shared.ScreenMemory.Remove(layerAlias)
	if parentAlias != "" {
		parentEntry := shared.ScreenMemory.Get(parentAlias)
		if parentEntry != nil {
			if !shared.IsAParent(parentAlias) {
				layerEntry = shared.ScreenMemory.Get(parentAlias)
				layerEntry.IsParent = false
			}
		}
	}
	if isParent {
		shared.deleteAllChildrenOfParent(layerAlias)
	}
}

func (shared *ScreenInstanceType) deleteAllChildrenOfParent(parentAlias string) {
	for _, currentValue := range shared.ScreenMemory.GetAllEntries() {
		if currentValue.ParentAlias == parentAlias {
			shared.DeleteLayer(currentValue.LayerAlias)
		}
	}
}

func (shared *ScreenInstanceType) IsAParent(parentAlias string) bool {
	isParent := false
	for _, currentValue := range shared.ScreenMemory.GetAllEntries() {
		if currentValue.ParentAlias == parentAlias {
			isParent = true
		}
	}
	return isParent
}

// IsLayerExists checks if a layer exists in memory.
func (shared *ScreenInstanceType) IsLayerExists(layerAlias string) bool {
	layerEntry := shared.ScreenMemory.Get(layerAlias)
	return layerEntry != nil
}

// GetSortedLayerMemoryAliasSlice returns a sorted list of layer aliases based on z-order.
func (shared *ScreenInstanceType) GetSortedLayerMemoryAliasSlice() LayerAliasZOrderPairList {
	pairList := make(LayerAliasZOrderPairList, len(shared.ScreenMemory.GetAllEntries()))
	currentEntry := 0
	for currentKey, currentValue := range shared.ScreenMemory.GetAllEntriesWithKeys() {
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
func (shared *ScreenInstanceType) SetHighestZOrderNumber(layerAlias string, parentAlias string) {
	if shared.IsLayerExists(layerAlias) {
		highestZOrderNumber := shared.getHighestZOrderNumber(parentAlias)
		for _, currentValue := range shared.ScreenMemory.GetAllEntries() {
			if currentValue.ParentAlias == parentAlias && currentValue.ZOrder == highestZOrderNumber {
				currentValue.ZOrder = highestZOrderNumber - 1
				currentValue.IsTopmost = false
			}
		}
		Screen.GetLayer(layerAlias).ZOrder = highestZOrderNumber
		Screen.GetLayer(layerAlias).IsTopmost = true
	}
}

func (shared *ScreenInstanceType) getHighestZOrderNumber(parentAlias string) int {
	highestZOrderNumber := 0
	for _, currentValue := range shared.ScreenMemory.GetAllEntries() {
		if currentValue.ParentAlias == parentAlias && currentValue.ZOrder > highestZOrderNumber {
			highestZOrderNumber = currentValue.ZOrder
		}
	}
	return highestZOrderNumber
}

func (shared *ScreenInstanceType) GetRootParentLayerAlias(layerAlias string, previousChildAlias string) (string, string) {
	layerEntry := shared.GetLayer(layerAlias)
	if layerEntry.ParentAlias != "" {
		return shared.GetRootParentLayerAlias(layerEntry.ParentAlias, layerAlias)
	}
	return layerAlias, previousChildAlias
}
