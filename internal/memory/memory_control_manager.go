package memory

import (
	"sort"
)

// ControlMemoryManager is a generic memory manager for handling layer-specific entries.
type ControlMemoryManager[T any] struct {
	MemoryManager map[string]*MemoryManager
}

// NewControlMemoryManager creates a new instance of ControlMemoryManager.
func NewControlMemoryManager[T any]() *ControlMemoryManager[T] {
	return &ControlMemoryManager[T]{
		MemoryManager: make(map[string]*MemoryManager),
	}
}

// Add inserts an entry into the specified layer's memory.
func (l *ControlMemoryManager[T]) Add(layerAlias string, alias string, entry T) {
	if l.MemoryManager[layerAlias] == nil {
		l.MemoryManager[layerAlias] = NewMemoryManager()
	}
	l.MemoryManager[layerAlias].Add(alias, entry)
}

// Remove deletes an entry from the specified layer's memory.
func (l *ControlMemoryManager[T]) Remove(layerAlias string, alias string) {
	if l.MemoryManager[layerAlias] != nil {
		l.MemoryManager[layerAlias].Remove(alias)
	}
}

// RemoveAll deletes all entries from the specified layer's memory.
func (l *ControlMemoryManager[T]) RemoveAll(layerAlias string) {
	if l.MemoryManager[layerAlias] != nil {
		l.MemoryManager[layerAlias].RemoveAll() // Assuming the MemoryManager struct has a RemoveAll method.
	}
}

// Get retrieves an entry from the specified layer's memory.
func (l *ControlMemoryManager[T]) Get(layerAlias string, alias string) T {
	if l.MemoryManager[layerAlias] != nil {
		value := l.MemoryManager[layerAlias].Get(alias)
		if typedValue, ok := value.(T); ok {
			return typedValue
		}
	}
	var zeroValue T
	return zeroValue
}

// GetAllEntries retrieves all entries from the specified layer.
func (l *ControlMemoryManager[T]) GetAllEntries(layerAlias string) []T {
	if l.MemoryManager[layerAlias] == nil {
		return []T{}
	}
	allEntries := l.MemoryManager[layerAlias].GetAllEntries()
	typedEntries := make([]T, 0, len(allEntries))
	for _, entry := range allEntries {
		if typedEntry, ok := entry.(T); ok {
			typedEntries = append(typedEntries, typedEntry)
		}
	}
	return typedEntries
}

func (l *ControlMemoryManager[T]) GetAllEntriesOverall() []T {
	var allEntries []T
	for layerAlias := range l.MemoryManager {
		layerEntries := l.GetAllEntries(layerAlias)
		allEntries = append(allEntries, layerEntries...)
	}
	return allEntries
}

// GetAllEntriesAsAliasList retrieves all aliases from the specified layer.
func (l *ControlMemoryManager[T]) GetAllEntriesAsAliasList(layerAlias string, getAlias func(T) string) []string {
	allEntries := l.GetAllEntries(layerAlias)
	aliases := make([]string, 0, len(allEntries))
	for _, entry := range allEntries {
		aliases = append(aliases, getAlias(entry))
	}
	return aliases
}

// SortEntries sorts entries in the specified layer using a custom comparator.
func (l *ControlMemoryManager[T]) SortEntries(layerAlias string, isAscendingOrder bool, compare func(a, b T) bool) []T {
	allEntries := l.GetAllEntries(layerAlias)
	sortedEntries := append([]T{}, allEntries...) // Make a copy to avoid mutating the original slice
	sort.Slice(sortedEntries, func(i, j int) bool {
		if isAscendingOrder {
			return compare(sortedEntries[i], sortedEntries[j])
		}
		return compare(sortedEntries[j], sortedEntries[i])
	})
	return sortedEntries
}
