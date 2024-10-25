package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sort"
	"sync"
)

var screenWidth int
var screenHeight int

type screenMemoryType struct {
	sync.Mutex
	Entries map[string]*types.LayerEntryType
}

var Screen screenMemoryType

func InitializeScreenMemory() {
	Screen.Entries = make(map[string]*types.LayerEntryType)
}

func AddLayer(layerAlias string, xLocation int, yLocation int, width int, height int, zOrderPriority int, parentAlias string) {
	Screen.Mutex.Lock()
	defer func() {
		Screen.Mutex.Unlock()
	}()
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
		if _, isExist := Screen.Entries[parentAlias]; isExist {
			parentEntry := Screen.Entries[parentAlias]
			parentEntry.IsParent = true
		} else {
			panic(fmt.Sprintf("The layer '%s' could not be created since the parent alias '%s' does not exist!", layerAlias, parentAlias))
		}

	}
	// consolizer.commonResource.layerAlias = layerAlias
	Screen.Entries[layerAlias] = &layerEntry
}

func getLayer(layerAlias string) *types.LayerEntryType {
	if _, isExist := Screen.Entries[layerAlias]; !isExist {
		panic(fmt.Sprintf("The layer '%s' could not be obtained since it does not exist!", layerAlias))
	}
	return Screen.Entries[layerAlias]
}

func GetLayer(layerAlias string) *types.LayerEntryType {
	Screen.Mutex.Lock()
	defer func() {
		Screen.Mutex.Unlock()
	}()
	if _, isExist := Screen.Entries[layerAlias]; !isExist {
		panic(fmt.Sprintf("The layer '%s' could not be obtained since it does not exist!", layerAlias))
	}
	return Screen.Entries[layerAlias]
}

func GetNextLayerAlias() string {
	Screen.Mutex.Lock()
	defer func() {
		Screen.Mutex.Unlock()
	}()
	for _, currentEntry := range Screen.Entries {
		return currentEntry.LayerAlias
	}
	return ""
}

func DeleteLayer(layerAlias string, noLocking bool) {
	if !noLocking {
		Screen.Mutex.Lock()
		defer func() {
			Screen.Mutex.Unlock()
		}()
	}
	if _, isExist := Screen.Entries[layerAlias]; !isExist {
		panic(fmt.Sprintf("The layer '%s' could not be deleted since it does not exist!", layerAlias))
	}
	layerEntry := Screen.Entries[layerAlias]
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
	delete(Screen.Entries, layerAlias)
	if parentAlias != "" {
		if _, isExist := Screen.Entries[parentAlias]; isExist {
			if !IsAParent(parentAlias) {
				layerEntry = getLayer(parentAlias)
				layerEntry.IsParent = false
			}
		}
	}
	if isParent {
		deleteAllChildrenOfParent(layerAlias)
	}
}

func deleteAllChildrenOfParent(parentAlias string) {
	for currentKey, currentValue := range Screen.Entries {
		if currentValue.ParentAlias == parentAlias {
			DeleteLayer(currentKey, true)
		}
	}
}
func IsLayerExists(layerAlias string) bool {
	Screen.Mutex.Lock()
	defer func() {
		Screen.Mutex.Unlock()
	}()
	if _, isExist := Screen.Entries[layerAlias]; isExist {
		return true
	}
	return false
}

func isLayerExists(layerAlias string) bool {
	if _, isExist := Screen.Entries[layerAlias]; isExist {
		return true
	}
	return false
}

func IsAParent(parentAlias string) bool {
	isParent := false
	for _, currentValue := range Screen.Entries {
		if currentValue.ParentAlias == parentAlias {
			isParent = true
		}
	}
	return isParent
}

type layerAliasZOrderPair struct {
	Key   string
	Value int
}
type LayerAliasZOrderPairList []layerAliasZOrderPair

func GetSortedLayerMemoryAliasSlice() LayerAliasZOrderPairList {
	Screen.Mutex.Lock()
	defer func() {
		Screen.Mutex.Unlock()
	}()
	pairList := make(LayerAliasZOrderPairList, len(Screen.Entries))
	currentEntry := 0
	for currentKey, currentValue := range Screen.Entries {
		pairList[currentEntry].Key = currentKey
		pairList[currentEntry].Value = currentValue.ZOrder
		currentEntry++
	}
	sort.SliceStable(pairList, func(firstIndex, secondIndex int) bool {
		return pairList[firstIndex].Value < pairList[secondIndex].Value
	})
	return pairList
}

func getHighestZOrderNumber(parentAlias string) int {
	highestZOrderNumber := 0
	for _, currentValue := range Screen.Entries {
		if currentValue.ParentAlias == parentAlias && currentValue.ZOrder > highestZOrderNumber {
			highestZOrderNumber = currentValue.ZOrder
		}
	}
	return highestZOrderNumber
}

func SetHighestZOrderNumber(layerAlias string, parentAlias string) {
	if IsLayerExists(layerAlias) {
		highestZOrderNumber := getHighestZOrderNumber(parentAlias)
		for _, currentValue := range Screen.Entries {
			if currentValue.ParentAlias == parentAlias && currentValue.ZOrder == highestZOrderNumber {
				currentValue.ZOrder = highestZOrderNumber - 1
				currentValue.IsTopmost = false
			}
		}
		Screen.Entries[layerAlias].ZOrder = highestZOrderNumber
		Screen.Entries[layerAlias].IsTopmost = true

	}
}

func GetRootParentLayerAlias(layerAlias string, previousChildAlias string) (string, string) {
	layerEntry := GetLayer(layerAlias)
	if layerEntry.ParentAlias != "" {
		return GetRootParentLayerAlias(layerEntry.ParentAlias, layerAlias)
	}
	return layerAlias, previousChildAlias
}
