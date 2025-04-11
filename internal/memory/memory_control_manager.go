package memory

import (
	"fmt"
	"reflect"
	"sort"
)

// ControlMemoryManager is a generic memory manager for handling layer-specific entries.
type ControlMemoryManager[T any] struct {
	MemoryManager map[string]*MemoryManager[T] // MemoryManager stores pointers to T
}

// NewControlMemoryManager creates a new instance of ControlMemoryManager.
/*
NewControlMemoryManager allows you to create a new control memory manager. In addition, the following
information should be noted:

- Initializes a new memory manager for handling layer-specific entries.
- Creates an empty map for storing control entries.
- The manager is generic and can handle any type of control.
*/
func NewControlMemoryManager[T any]() *ControlMemoryManager[T] {
	return &ControlMemoryManager[T]{
		MemoryManager: make(map[string]*MemoryManager[T]),
	}
}

// Add inserts a pointer entry into the specified layer's memory.
/*
Add allows you to insert a control entry into the memory manager. In addition, the following
information should be noted:

- Adds a control entry to the specified layer.
- Creates a new layer if it doesn't exist.
- The entry is stored as a pointer to allow for updates.
*/
func (shared *ControlMemoryManager[T]) Add(layerAlias string, alias string, entry *T) {
	// Ensure the layer exists, or create a new one
	if shared.MemoryManager[layerAlias] == nil {
		shared.MemoryManager[layerAlias] = NewMemoryManager[T]()
	}
	// AddLayer the pointer entry to the specified layer
	shared.MemoryManager[layerAlias].Add(alias, entry)
}

// Remove deletes an entry from the specified layer's memory.
/*
Remove allows you to delete a control entry from the memory manager. In addition, the following
information should be noted:

- Removes a control entry from the specified layer.
- Does nothing if the layer or entry doesn't exist.
- The entry's memory is freed when removed.
*/
func (shared *ControlMemoryManager[T]) Remove(layerAlias string, alias string) {
	if shared.MemoryManager[layerAlias] != nil {
		shared.MemoryManager[layerAlias].Remove(alias)
	}
}

// RemoveAll deletes all entries from the specified layer's memory.
/*
RemoveAll allows you to delete all control entries from a layer. In addition, the following
information should be noted:

- Removes all control entries from the specified layer.
- Does nothing if the layer doesn't exist.
- All memory associated with the entries is freed.
*/
func (shared *ControlMemoryManager[T]) RemoveAll(layerAlias string) {
	if shared.MemoryManager[layerAlias] != nil {
		shared.MemoryManager[layerAlias].RemoveAll()
	}
}

// Get retrieves a pointer entry from the specified layer's memory.
/*
Get allows you to retrieve a control entry from the memory manager. In addition, the following
information should be noted:

- Returns a pointer to the control entry if it exists.
- Returns nil if the layer or entry doesn't exist.
- The entry can be modified through the returned pointer.
*/
func (shared *ControlMemoryManager[T]) Get(layerAlias string, alias string) *T {
	typeName := reflect.TypeOf(*new(T)).Name() // GetLayer the type name without pointer
	if shared.MemoryManager[layerAlias] != nil {
		value := shared.MemoryManager[layerAlias].Get(alias)
		if value == nil {
			// Use reflect to get a human-readable type name (without pointer format)
			panic(fmt.Sprintf("The %s '%s' under layer '%s' could not be obtained since it does not exist!", typeName, alias, layerAlias))
		}
		return value
	}
	panic(fmt.Sprintf("The layer '%s' for '%s' could not be found!", layerAlias, typeName))
}

// GetAllEntries retrieves all entries as pointers from the specified layer.
/*
GetAllEntries allows you to retrieve all control entries from a layer. In addition, the following
information should be noted:

- Returns a slice of all control entries in the specified layer.
- Returns an empty slice if the layer doesn't exist.
- The entries are returned in alphabetical order by alias.
*/
func (shared *ControlMemoryManager[T]) GetAllEntries(layerAlias string) []*T {
	if shared.MemoryManager[layerAlias] == nil {
		return []*T{} // Return an empty slice if the layer doesn't exist
	}
	allEntries := shared.MemoryManager[layerAlias].GetAllEntries()
	return allEntries // Return the slice of pointers
}

// GetAllEntriesOverall retrieves all entries from all layers.
func (shared *ControlMemoryManager[T]) GetAllEntriesOverall() []*T {
	var allEntries []*T
	for layerAlias := range shared.MemoryManager {
		layerEntries := shared.GetAllEntries(layerAlias)
		allEntries = append(allEntries, layerEntries...)
	}
	return allEntries
}

// GetAllEntriesAsAliasList retrieves all aliases from the specified layer.
func (shared *ControlMemoryManager[T]) GetAllEntriesAsAliasList(layerAlias string, getAlias func(*T) string) []string {
	allEntries := shared.GetAllEntries(layerAlias)
	aliases := make([]string, 0, len(allEntries))
	for _, entry := range allEntries {
		aliases = append(aliases, getAlias(entry))
	}
	return aliases
}

// SortEntries sorts entries in the specified layer using a custom comparator.
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

// IsExists checks if an entry with the given alias exists in the specified layer.
/*
IsExists allows you to check if a control entry exists in the memory manager. In addition, the following
information should be noted:

- Returns true if the control entry exists in the specified layer.
- Returns false if the layer or entry doesn't exist.
- Useful for validation before performing operations.
*/
func (shared *ControlMemoryManager[T]) IsExists(layerAlias string, alias string) bool {
	if shared.MemoryManager[layerAlias] != nil {
		return shared.MemoryManager[layerAlias].Get(alias) != nil
	}
	return false
}
