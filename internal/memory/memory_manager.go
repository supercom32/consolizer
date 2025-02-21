package memory

import (
	"sync"
)

type MemoryManager struct {
	muxtex      sync.RWMutex
	memoryItems map[string]interface{} // Store values as interface{}
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		memoryItems: make(map[string]interface{}),
	}
}

// Add inserts a value into the memoryItems map under the specified key.
func (shared *MemoryManager) Add(key string, value interface{}) {
	shared.muxtex.Lock()
	defer shared.muxtex.Unlock()
	shared.memoryItems[key] = value // Directly store the interface{}
}

// Remove deletes a key-value pair from the memoryItems map.
func (shared *MemoryManager) Remove(key string) {
	shared.muxtex.Lock()
	defer shared.muxtex.Unlock()
	delete(shared.memoryItems, key)
}

// Get retrieves the value stored at the specified key.
func (shared *MemoryManager) Get(key string) interface{} {
	shared.muxtex.RLock()
	defer shared.muxtex.RUnlock()
	return shared.memoryItems[key] // Return the value directly
}

// GetAllEntries retrieves all values stored in memoryItems.
func (shared *MemoryManager) GetAllEntries() []interface{} {
	shared.muxtex.RLock()
	defer shared.muxtex.RUnlock()

	var entries []interface{}
	for _, value := range shared.memoryItems {
		entries = append(entries, value) // Append the value directly
	}
	return entries
}

// RemoveAll clears all entries from the memoryItems map.
func (shared *MemoryManager) RemoveAll() {
	shared.muxtex.Lock()
	defer shared.muxtex.Unlock()
	shared.memoryItems = make(map[string]interface{}) // Reinitialize the map to reset it
}
