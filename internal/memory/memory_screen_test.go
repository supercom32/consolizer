package memory

import (
	"github.com/stretchr/testify/assert"
	"supercom32.net/consolizer"
	"supercom32.net/consolizer/internal/recast"
	"testing"
)

func TestScreenLayerCreation(test *testing.T) {
	layerAlias := "MyAlias"
	layerWidth := 20
	layerHeight := 10
	layerXLocation := 1
	layerYLocation := 2
	layerZOrderPriority := 1
	layerParentAlias := ""

	InitializeScreenMemory()
	consolizer.AddLayer(layerAlias, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, layerParentAlias)
	layer := consolizer.GetLayer(layerAlias)
	obtainedResult := recast.GetArrayOfInterfaces(layer.Width, layer.Height, layer.ScreenXLocation, layer.ScreenYLocation, layer.ZOrder, layer.ParentAlias)
	expectedResult := recast.GetArrayOfInterfaces(20, 10, 1, 2, 1, "")
	assert.Equalf(test, expectedResult, obtainedResult, "The created layer was not added correctly!")
}

func TestScreenLayerParentIsCorrectlyLinked(test *testing.T) {
	layerWidth := 20
	layerHeight := 10
	layerXLocation := 1
	layerYLocation := 2
	layerZOrderPriority := 1
	parentAlias := "ParentAlias"
	childAlias := "ChildAlias"
	InitializeScreenMemory()
	consolizer.AddLayer(parentAlias, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, "")
	parentLayer := consolizer.GetLayer(parentAlias)
	consolizer.AddLayer(childAlias, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, "ParentAlias")
	childLayer := consolizer.GetLayer(childAlias)
	if parentLayer.IsParent != true {
		test.Errorf("Creating a child layer failed to update the 'IsParent' flag on the parent layer!")
	}
	if childLayer.ParentAlias != parentAlias {
		test.Errorf("Creating a child layer did not update itself with the correct parent alias!")
	}
}

func TestScreenLayerInvalidParent(test *testing.T) {
	layerAlias := "MyAlias"
	layerWidth := 20
	layerHeight := 10
	layerXLocation := 1
	layerYLocation := 2
	layerZOrderPriority := 1
	layerParentAlias := "BadParent"
	defer func() {
		if r := recover(); r == nil {
			test.Errorf("Creating a layer with a bad parent should have thrown a panic!")
		}
	}()
	InitializeScreenMemory()
	consolizer.AddLayer(layerAlias, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, layerParentAlias)
}

func TestScreenLayerInvalidWidth(test *testing.T) {
	layerAlias := "MyAlias"
	layerWidth := -20
	layerHeight := 10
	layerXLocation := 1
	layerYLocation := 2
	layerZOrderPriority := 1
	layerParentAlias := ""
	defer func() {
		if r := recover(); r == nil {
			test.Errorf("Creating a layer with an invalid HotspotWidth should have thrown a panic!")
		}
	}()
	InitializeScreenMemory()
	consolizer.AddLayer(layerAlias, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, layerParentAlias)
}

func TestScreenLayerInvalidHeight(test *testing.T) {
	layerAlias := "MyAlias"
	layerWidth := 20
	layerHeight := -10
	layerXLocation := 1
	layerYLocation := 2
	layerZOrderPriority := 1
	layerParentAlias := ""
	defer func() {
		if r := recover(); r == nil {
			test.Errorf("Creating a layer with an invalid Length should have thrown a panic!")
		}
	}()
	InitializeScreenMemory()
	consolizer.AddLayer(layerAlias, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, layerParentAlias)
}

func TestScreenLayerSimpleDelete(test *testing.T) {
	layerAlias := "MyAlias"
	layerWidth := 20
	layerHeight := 10
	layerXLocation := 1
	layerYLocation := 2
	layerZOrderPriority := 1
	layerParentAlias := ""
	InitializeScreenMemory()
	consolizer.AddLayer(layerAlias, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, layerParentAlias)
	consolizer.DeleteLayer(layerAlias, false)
	defer func() {
		if r := recover(); r == nil {
			test.Errorf("Obtaining a layer that has already been deleted should throw a panic!")
		}
	}()
	_ = consolizer.GetLayer(layerAlias)
}

func TestScreenLayerChildrenDelete(test *testing.T) {
	layerWidth := 20
	layerHeight := 10
	layerXLocation := 1
	layerYLocation := 2
	layerZOrderPriority := 1
	layerParentAlias := "ParentAlias"

	childAlias1 := "ChildAlias1"
	childAlias2 := "ChildAlias2"
	childAlias3 := "ChildAlias3"

	InitializeScreenMemory()
	consolizer.AddLayer(layerParentAlias, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, "")
	consolizer.AddLayer(childAlias1, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, layerParentAlias)
	consolizer.AddLayer(childAlias2, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, layerParentAlias)
	consolizer.AddLayer(childAlias3, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, layerParentAlias)
	var layerEntry = consolizer.GetLayer(layerParentAlias)

	consolizer.DeleteLayer(childAlias2, false)
	if layerEntry.IsParent != true {
		test.Errorf("The parent layer is no longer marked as a parent when it should have two children remaining!")
	}
	consolizer.DeleteLayer(childAlias1, false)
	if layerEntry.IsParent != true {
		test.Errorf("The parent layer is no longer marked as a parent when it should have one children remaining!")
	}
	consolizer.DeleteLayer(childAlias3, false)
	if layerEntry.IsParent == true {
		test.Errorf("The parent layer is no longer a parent, but is still marked as one!")
	}
}

func TestScreenLayerParentDelete(test *testing.T) {
	var layerWidth = 20
	var layerHeight = 10
	var layerXLocation = 1
	var layerYLocation = 2
	var layerZOrderPriority = 1

	var parentAlias = "ParentAlias"
	var childAlias1 = "ChildAlias1"
	var childAlias2 = "ChildAlias2"
	var childAlias3 = "ChildAlias3"

	InitializeScreenMemory()
	consolizer.AddLayer(parentAlias, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, "")
	consolizer.AddLayer(childAlias1, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, parentAlias)
	consolizer.AddLayer(childAlias2, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, parentAlias)
	consolizer.AddLayer(childAlias3, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, parentAlias)
	consolizer.DeleteLayer(parentAlias, false)
	if consolizer.IsLayerExists(parentAlias) {
		test.Errorf("The parent layer exists when it should have been deleted!")
	}
	if consolizer.IsLayerExists(childAlias1) {
		test.Errorf("The first child layer should have been deleted since the parent no longer exists!")
	}
	if consolizer.IsLayerExists(childAlias2) {
		test.Errorf("The second child layer should have been deleted since the parent no longer exists!")
	}
	if consolizer.IsLayerExists(childAlias3) {
		test.Errorf("The third child layer should have been deleted since the parent no longer exists!")
	}
}

func TestScreenLayerSubParentDelete(test *testing.T) {
	var layerWidth = 20
	var layerHeight = 10
	var layerXLocation = 1
	var layerYLocation = 2
	var layerZOrderPriority = 1

	var parentAlias = "ParentAlias"
	var childAlias1 = "ChildAlias1"
	var childAlias3 = "ChildAlias3"
	var subParent1 = "SubParent1"
	var subChild1 = "SubChild1"
	var subChild2 = "SubChild2"

	InitializeScreenMemory()
	consolizer.AddLayer(parentAlias, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, "")
	consolizer.AddLayer(childAlias1, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, parentAlias)
	consolizer.AddLayer(subParent1, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, parentAlias)
	consolizer.AddLayer(childAlias3, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, parentAlias)
	consolizer.AddLayer(subChild1, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, subParent1)
	consolizer.AddLayer(subChild2, layerXLocation, layerYLocation, layerWidth, layerHeight, layerZOrderPriority, subParent1)
	consolizer.DeleteLayer(subParent1, false)
	if consolizer.IsLayerExists(subChild1) {
		test.Errorf("Deleting the sub-parent layer did not delete the first sub-child layer as expected!")
	}
	if consolizer.IsLayerExists(subChild2) {
		test.Errorf("Deleting the sub-parent layer did not delete the second sub-child layer as expected!")
	}
	if consolizer.IsLayerExists(subParent1) {
		test.Errorf("The sub-parent was supposed to be deleted, but it still exists!")
	}
}

func TestScreenLayerSorting(test *testing.T) {
	var layerWidth = 20
	var layerHeight = 10
	var layerXLocation = 1
	var layerYLocation = 2
	var layerParentAlias = ""
	InitializeScreenMemory()
	consolizer.AddLayer("Alias1", layerXLocation, layerYLocation, layerWidth, layerHeight, 6, layerParentAlias)
	consolizer.AddLayer("Alias2", layerXLocation, layerYLocation, layerWidth, layerHeight, 4, layerParentAlias)
	consolizer.AddLayer("Alias3", layerXLocation, layerYLocation, layerWidth, layerHeight, 1, layerParentAlias)
	consolizer.AddLayer("Alias4", layerXLocation, layerYLocation, layerWidth, layerHeight, 3, layerParentAlias)
	consolizer.AddLayer("Alias5", layerXLocation, layerYLocation, layerWidth, layerHeight, 9, layerParentAlias)
	consolizer.AddLayer("Alias6", layerXLocation, layerYLocation, layerWidth, layerHeight, 8, layerParentAlias)
	var pairList consolizer.LayerAliasZOrderPairList = consolizer.GetSortedLayerMemoryAliasSlice()
	if pairList[0].Key != "Alias3" || pairList[1].Key != "Alias4" || pairList[2].Key != "Alias2" ||
		pairList[3].Key != "Alias1" || pairList[4].Key != "Alias6" || pairList[5].Key != "Alias5" {
		test.Errorf("The sorted screen layer pair list is not correct")
	}
}
