package types

import (
	"github.com/stretchr/testify/assert"
	"supercom32.net/consolizer/internal/recast"
	"testing"
)

func TestNewSelectionEntry(test *testing.T) {
	selectionEntry := NewSelectionEntry()
	selectionEntry.Add("selectionAlias1", "selectionValue1")
	obtainedValue := recast.GetArrayOfInterfaces(selectionEntry.SelectionAlias[0], selectionEntry.SelectionValue[0])
	expectedValue := recast.GetArrayOfInterfaces("selectionAlias1", "selectionValue1")
	assert.Equalf(test, expectedValue, obtainedValue, "The selection entry obtained does not match what was set!")
	selectionEntry.Clear()
	obtainedValue = recast.GetArrayOfInterfaces(len(selectionEntry.SelectionAlias), len(selectionEntry.SelectionValue))
	expectedValue = recast.GetArrayOfInterfaces(0, 0)
	assert.Equalf(test, expectedValue, obtainedValue, "The number of selection entries does not what was expected!")
}
