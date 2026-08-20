package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
TestLayerTypeCreation is a test which verifies the creation and cloning of layer entries.

Example:

	Expected Inputs:
	    None

	Expected Outputs:
	    None
*/
func TestLayerTypeCreation(test *testing.T) {
	layerAlias := "MyAlias"
	parentAlias := "MyParentAlias"
	firstLayerEntry := NewLayerEntry(layerAlias, parentAlias, 20, 20)
	firstLayerEntry.IsParent = false
	firstLayerEntry.ScreenXLocation = 1
	firstLayerEntry.ScreenYLocation = 2
	firstLayerEntry.CursorXLocation = 3
	firstLayerEntry.CursorYLocation = 4
	firstLayerEntry.ZOrder = 1
	firstLayerEntry.AlphaValue = 0.5
	firstLayerEntry.IsVisible = true
	secondLayerEntry := NewLayerEntry(layerAlias, parentAlias, 20, 20)
	assert.NotEqualf(test, secondLayerEntry, firstLayerEntry, "The first layer entry is the same as the second, even though it should be different.")

	secondLayerEntry = NewLayerEntry(layerAlias, parentAlias, 0, 0, &firstLayerEntry)
	assert.Equalf(test, secondLayerEntry, firstLayerEntry, "The first layer is not the same as the second, even though it should be an identical clone.")
	assert.Equalf(test, float32(0.5), secondLayerEntry.AlphaValue, "The alpha value was not cloned correctly.")
}
