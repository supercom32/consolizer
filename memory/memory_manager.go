package memory

import (
	"fmt"
	"sync"
)

type MemoryManager[T any] struct {
	muxtex      sync.RWMutex
	memoryItems map[string]*T // Store pointers to values of type T
}

func NewMemoryManager[T any]() *MemoryManager[T] {
	return &MemoryManager[T]{
		memoryItems: make(map[string]*T),
	}
}

// Add inserts a pointer value into the memoryItems map under the specified key.
// Only pointers to T are allowed; if a value is not a pointer, it panics.
func (shared *MemoryManager[T]) Add(key string, value *T) {
	if value == nil {
		panic(fmt.Sprintf("Cannot add nil value for key %s", key))
	}

	shared.muxtex.Lock()
	defer shared.muxtex.Unlock()
	shared.memoryItems[key] = value // Store the pointer to T directly
}

// Remove deletes a key-value pair from the memoryItems map.
func (shared *MemoryManager[T]) Remove(key string) {
	shared.muxtex.Lock()
	defer shared.muxtex.Unlock()
	delete(shared.memoryItems, key)
}

// Get retrieves the value stored at the specified key.
// Returns a pointer to T.
func (shared *MemoryManager[T]) Get(key string) *T {
	shared.muxtex.RLock()
	defer shared.muxtex.RUnlock()
	return shared.memoryItems[key] // Return the pointer directly
}

// GetAllEntries retrieves all values stored in memoryItems as pointers.
func (shared *MemoryManager[T]) GetAllEntries() []*T {
	shared.muxtex.RLock()
	defer shared.muxtex.RUnlock()

	var entries []*T
	for _, value := range shared.memoryItems {
		entries = append(entries, value) // Append the pointer directly
	}
	return entries
}

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

// RemoveAll clears all entries from the memoryItems map.
func (shared *MemoryManager[T]) RemoveAll() {
	shared.muxtex.Lock()
	defer shared.muxtex.Unlock()
	shared.memoryItems = make(map[string]*T) // Reinitialize the map to reset it
}

// IsExists checks if a value with the given key exists in memoryItems.
func (shared *MemoryManager[T]) IsExists(key string) bool {
	shared.muxtex.RLock()
	defer shared.muxtex.RUnlock()
	_, exists := shared.memoryItems[key]
	return exists
}
