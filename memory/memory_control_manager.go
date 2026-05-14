package memory

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
)

// ControlMemoryManager is a generic memory manager for handling layer-specific entries.
type ControlMemoryManager[T any] struct {
	MemoryManager sync.Map // Keys are strings, values are *MemoryManager[T]
}

/*
NewControlMemoryManager is a method which allows you to create a new control memory manager. In addition, the following
information should be noted:

- Initializes a new memory manager for handling layer-specific entries.

- Creates an empty sync.Map for storing control entries.

- The manager is generic and can handle any type of control.

- Uses sync.Map for thread-safe concurrent access.

:return: A pointer to a new ControlMemoryManager instance.

Example:

	manager := NewControlMemoryManagertypes.ButtonType()
*/
func NewControlMemoryManager[T any]() *ControlMemoryManager[T] {
	return &ControlMemoryManager[T]{
		MemoryManager: sync.Map{},
	}
}

/*
Add is a method which allows you to insert a control entry into the memory manager. In addition, the following
information should be noted:

- Adds a control entry to the specified layer.

- Creates a new layer if it doesn't exist.

- The entry is stored as a pointer to allow for updates.

:param layerAlias: The alias of the layer to add the entry to.
:param alias: The alias of the entry to add.
:param entry: The pointer to the entry to add.

Example:

	manager.Add("layer1", "button1", &button)
*/
func (shared *ControlMemoryManager[T]) Add(layerAlias string, alias string, entry *T) {
	// Ensure the layer exists, or create a new one
	layerManager, ok := shared.MemoryManager.Load(layerAlias)
	if !ok || layerManager == nil {
		layerManager = NewMemoryManager[T]()
		shared.MemoryManager.Store(layerAlias, layerManager)
	}
	// Add the pointer entry to the specified layer
	layerManager.(*MemoryManager[T]).Add(alias, entry)
}

/*
Remove is a method which allows you to delete a control entry from the memory manager. In addition, the following
information should be noted:

- Removes a control entry from the specified layer.

- Does nothing if the layer or entry doesn't exist.

- The entry's memory is freed when removed.

:param layerAlias: The alias of the layer to remove the entry from.
:param alias: The alias of the entry to remove.

Example:

	manager.Remove("layer1", "button1")
*/
func (shared *ControlMemoryManager[T]) Remove(layerAlias string, alias string) {
	layerManager, ok := shared.MemoryManager.Load(layerAlias)
	if ok && layerManager != nil {
		layerManager.(*MemoryManager[T]).Remove(alias)
	}
}

/*
RemoveAll is a method which allows you to delete all control entries from a layer. In addition, the following
information should be noted:

- Removes all control entries from the specified layer.

- Does nothing if the layer doesn't exist.

- All memory associated with the entries is freed.

:param layerAlias: The alias of the layer to clear.

Example:

	manager.RemoveAll("layer1")
*/
func (shared *ControlMemoryManager[T]) RemoveAll(layerAlias string) {
	layerManager, ok := shared.MemoryManager.Load(layerAlias)
	if ok && layerManager != nil {
		layerManager.(*MemoryManager[T]).RemoveAll()
	}
}

/*
Get is a method which allows you to retrieve a control entry from the memory manager. In addition, the following
information should be noted:

- Returns a pointer to the control entry if it exists.

- Returns nil if the layer or entry doesn't exist.

- The entry can be modified through the returned pointer.

:param layerAlias: The alias of the layer to get the entry from.
:param alias: The alias of the entry to get.

:return: A pointer to the requested entry.

Example:

	entry := manager.Get("layer1", "button1")
*/
func (shared *ControlMemoryManager[T]) Get(layerAlias string, alias string) *T {
	typeName := reflect.TypeOf(*new(T)).Name() // Get the type name without pointer
	layerManager, ok := shared.MemoryManager.Load(layerAlias)
	if ok && layerManager != nil {
		value := layerManager.(*MemoryManager[T]).Get(alias)
		if value == nil {
			// Use reflect to get a human-readable type name (without pointer format)
			panic(fmt.Sprintf("The %s '%s' under layer '%s' could not be obtained since it does not exist!", typeName, alias, layerAlias))
		}
		return value
	}
	panic(fmt.Sprintf("The layer '%s' for '%s' could not be found!", layerAlias, typeName))
}

/*
GetAllEntries is a method which allows you to retrieve all control entries from a layer. In addition, the following
information should be noted:

- Returns a slice of all control entries in the specified layer.

- Returns an empty slice if the layer doesn't exist.

- The entries are returned in alphabetical order by alias.

:param layerAlias: The alias of the layer to get entries from.

:return: A slice of pointers to all entries in the layer.

Example:

	entries := manager.GetAllEntries("layer1")
*/
func (shared *ControlMemoryManager[T]) GetAllEntries(layerAlias string) []*T {
	layerManager, ok := shared.MemoryManager.Load(layerAlias)
	if !ok || layerManager == nil {
		return []*T{} // Return an empty slice if the layer doesn't exist
	}
	allEntries := layerManager.(*MemoryManager[T]).GetAllEntries()
	return allEntries // Return the slice of pointers
}

/*
GetAllEntriesOverall is a method which allows you to retrieve all control entries from all layers in the memory manager.

:return: A slice of pointers to all entries from all layers.

Example:

	allEntries := manager.GetAllEntriesOverall()
*/
func (shared *ControlMemoryManager[T]) GetAllEntriesOverall() []*T {
	var allEntries []*T
	shared.MemoryManager.Range(func(key, value interface{}) bool {
		layerAlias := key.(string)
		layerEntries := shared.GetAllEntries(layerAlias)
		allEntries = append(allEntries, layerEntries...)
		return true
	})
	return allEntries
}

/*
GetAllEntriesAsAliasList is a method which allows you to retrieve all aliases from a specific layer in the memory
manager.

:param layerAlias: The alias of the layer to get aliases from.
:param getAlias: A function that extracts the alias from an entry.

:return: A slice of strings containing all aliases in the layer.

Example:

	aliases := manager.GetAllEntriesAsAliasList("layer1", func(e *Button) string { return e.Alias })
*/
func (shared *ControlMemoryManager[T]) GetAllEntriesAsAliasList(layerAlias string, getAlias func(*T) string) []string {
	allEntries := shared.GetAllEntries(layerAlias)
	aliases := make([]string, 0, len(allEntries))
	for _, entry := range allEntries {
		aliases = append(aliases, getAlias(entry))
	}
	return aliases
}

/*
SortEntries is a method which allows you to sort control entries in a layer using a custom comparator.

:param layerAlias: The alias of the layer to sort.
:param isAscendingOrder: Whether to sort in ascending or descending order.
:param compare: A function that defines the sorting logic between two entries.

:return: A sorted slice of pointers to the control entries.

Example:

	sorted := manager.SortEntries("layer1", true, func(a, b *Button) bool { return a.ID < b.ID })
*/
func (shared *ControlMemoryManager[T]) SortEntries(layerAlias string, isAscendingOrder bool, compare func(a, b *T) bool) []*T {
	allEntries := shared.GetAllEntries(layerAlias)
	sortedEntries := append([]*T{}, allEntries...) // Make a copy to avoid mutating the original slice
	sort.Slice(sortedEntries, func(i, j int) bool {
		if isAscendingOrder {
			return compare(sortedEntries[i], sortedEntries[j])
		}
		return compare(sortedEntries[j], sortedEntries[i])
	})
	return sortedEntries
}

/*
IsExists is a method which allows you to check if a control entry exists in the memory manager. In addition, the
following information should be noted:

- Returns true if the control entry exists in the specified layer.

- Returns false if the layer or entry doesn't exist.

- Useful for validation before performing operations.

:param layerAlias: The alias of the layer to check.
:param alias: The alias of the entry to check.

:return: True if the entry exists, false otherwise.

Example:

	exists := manager.IsExists("layer1", "button1")
*/
func (shared *ControlMemoryManager[T]) IsExists(layerAlias string, alias string) bool {
	layerManager, ok := shared.MemoryManager.Load(layerAlias)
	if ok && layerManager != nil {
		return layerManager.(*MemoryManager[T]).Get(alias) != nil
	}
	return false
}
