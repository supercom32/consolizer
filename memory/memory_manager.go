package memory

import (
	"fmt"
	"sync"
)

type MemoryManager[T any] struct {
	muxtex      sync.RWMutex
	memoryItems map[string]*T // Store pointers to values of type T
}

/*
NewMemoryManager is a method which allows you to create a new generic memory manager.

:return: A pointer to a new MemoryManager instance.

Example:

	manager := NewMemoryManagertypes.ButtonType()
*/
func NewMemoryManager[T any]() *MemoryManager[T] {
	return &MemoryManager[T]{
		memoryItems: make(map[string]*T),
	}
}

/*
Add is a method which allows you to insert a pointer value into the memory manager under the specified key.

:param key: The unique key to store the value under.
:param value: A pointer to the value of type T to store.

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
Remove is a method which allows you to delete a key-value pair from the memory manager.

:param key: The key of the entry to remove.

Example:

	manager.Remove("button1")
*/
func (shared *MemoryManager[T]) Remove(key string) {
	shared.muxtex.Lock()
	defer shared.muxtex.Unlock()
	delete(shared.memoryItems, key)
}

/*
Get is a method which allows you to retrieve the value stored at the specified key.

:param key: The key of the entry to retrieve.

:return: A pointer to the value stored at the key, or nil if not found.

Example:

	value := manager.Get("button1")
*/
func (shared *MemoryManager[T]) Get(key string) *T {
	shared.muxtex.RLock()
	defer shared.muxtex.RUnlock()
	return shared.memoryItems[key] // Return the pointer directly
}

/*
GetAllEntries is a method which allows you to retrieve all values stored in the memory manager as a slice of pointers.

:return: A slice of pointers to all stored values of type T.

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
GetAllEntriesWithKeys is a method which allows you to retrieve all values stored in the memory manager as a map.

:return: A map where keys are strings and values are pointers of type T.

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
RemoveAll is a method which allows you to clear all entries from the memory manager.

Example:

	manager.RemoveAll()
*/
func (shared *MemoryManager[T]) RemoveAll() {
	shared.muxtex.Lock()
	defer shared.muxtex.Unlock()
	shared.memoryItems = make(map[string]*T) // Reinitialize the map to reset it
}

/*
IsExists is a method which allows you to check if a value with the given key exists in the memory manager.

:param key: The key to check for existence.

:return: True if the key exists, false otherwise.

Example:

	exists := manager.IsExists("button1")
*/
func (shared *MemoryManager[T]) IsExists(key string) bool {
	shared.muxtex.RLock()
	defer shared.muxtex.RUnlock()
	_, exists := shared.memoryItems[key]
	return exists
}
