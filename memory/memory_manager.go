package memory

import (
	"fmt"
	"sync"
)

/*
MemoryManager is a structure which handles generic memory management for storing pointer values.
*/
type MemoryManager[T any] struct {
	muxtex      sync.RWMutex
	memoryItems map[string]*T // Store pointers to values of type T
}

/*
NewMemoryManager is a method which creates a new generic memory manager.

Example:
    manager := NewMemoryManagertypes.ButtonType()
*/
func NewMemoryManager[T any]() *MemoryManager[T] {
	return &MemoryManager[T]{
		memoryItems: make(map[string]*T),
	}
}

/*
Add is a method which inserts a pointer value into the memory manager under the specified key.

Example:
    manager.Add("button1", &button)
*/
func (shared *MemoryManager[T]) Add(key string, value *T) {
	if value == nil {
		panic(fmt.Sprintf("Cannot add nil value for key %s", key))
	}

	shared.muxtex.Lock()
	defer shared.muxtex.Unlock()
	shared.memoryItems[key] = value // Store the pointer to T directly
}

/*
Remove is a method which deletes a key-value pair from the memory manager.

Example:
    manager.Remove("button1")
*/
func (shared *MemoryManager[T]) Remove(key string) {
	shared.muxtex.Lock()
	defer shared.muxtex.Unlock()
	delete(shared.memoryItems, key)
}

/*
Get is a method which retrieves the value stored at the specified key.

Example:
    value := manager.Get("button1")
*/
func (shared *MemoryManager[T]) Get(key string) *T {
	shared.muxtex.RLock()
	defer shared.muxtex.RUnlock()
	return shared.memoryItems[key] // Return the pointer directly
}

/*
GetAllEntries is a method which retrieves all values stored in the memory manager as a slice of pointers.

Example:
    entries := manager.GetAllEntries()
*/
func (shared *MemoryManager[T]) GetAllEntries() []*T {
	shared.muxtex.RLock()
	defer shared.muxtex.RUnlock()

	var entries []*T
	for _, value := range shared.memoryItems {
		entries = append(entries, value) // Append the pointer directly
	}
	return entries
}

/*
GetAllEntriesWithKeys is a method which retrieves all values stored in the memory manager as a map.

Example:
    entriesMap := manager.GetAllEntriesWithKeys()
*/
func (shared *MemoryManager[T]) GetAllEntriesWithKeys() map[string]*T {
	shared.muxtex.RLock()
	defer shared.muxtex.RUnlock()

	// Create a copy of the map to avoid race conditions
	entries := make(map[string]*T, len(shared.memoryItems))
	for key, value := range shared.memoryItems {
		entries[key] = value
	}
	return entries
}

/*
RemoveAll is a method which clears all entries from the memory manager.

Example:
    manager.RemoveAll()
*/
func (shared *MemoryManager[T]) RemoveAll() {
	shared.muxtex.Lock()
	defer shared.muxtex.Unlock()
	shared.memoryItems = make(map[string]*T) // Reinitialize the map to reset it
}

/*
IsExists is a method which checks if a value with the given key exists in the memory manager.

Example:
    exists := manager.IsExists("button1")
*/
func (shared *MemoryManager[T]) IsExists(key string) bool {
	shared.muxtex.RLock()
	defer shared.muxtex.RUnlock()
	_, exists := shared.memoryItems[key]
	return exists
}
