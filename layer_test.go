package consolizer

import (
	"testing"

	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/types"
)

func TestLayerInitialization(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	// Test initialization
	layer.ReInitializeScreenMemory()
	if Layers == nil {
		t.Error("Layers should not be nil after initialization")
	}
	DeleteAllLayers()
}

func TestLayerAdd(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	// Test valid layer creation
	layer.ReInitializeScreenMemory()
	layer.Add("testLayer", 0, 0, 10, 10, 1, "")

	if !isLayerExists("testLayer") {
		t.Error("Layer should exist after creation")
	}

	// Test invalid width
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid width")
		}
	}()
	layer.Add("invalidLayer", 0, 0, 0, 10, 1, "")
	DeleteAllLayers()
}

func TestLayerDelete(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create test layers
	layer.Add("parentLayer", 0, 0, 10, 10, 1, "")
	layer.Add("childLayer", 0, 0, 5, 5, 2, "parentLayer")

	// Test deletion
	layer.Delete("childLayer")
	if isLayerExists("childLayer") {
		t.Error("Child layer should not exist after deletion")
	}

	// Test parent layer still exists
	if !isLayerExists("parentLayer") {
		t.Error("Parent layer should still exist after child deletion")
	}

	// Test deleting non-existent layer
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for deleting non-existent layer")
		}
	}()
	layer.Delete("nonExistentLayer")
	DeleteAllLayers()
}

func TestLayerParentChild(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create parent-child relationship
	layer.Add("parent", 0, 0, 10, 10, 1, "")
	layer.Add("child", 0, 0, 5, 5, 2, "parent")

	// Test IsAParent
	if !layer.IsAParent("parent") {
		t.Error("Parent layer should be identified as parent")
	}

	// Test GetRootParentLayerAlias
	rootParent, child := layer.GetRootParentLayerAlias("child", "")
	if rootParent != "parent" {
		t.Error("Root parent should be 'parent'")
	}
	if child != "child" {
		t.Error("Child should be 'child'")
	}
	DeleteAllLayers()
}

func TestLayerZOrder(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create layers with different z-orders
	layer.Add("layer1", 0, 0, 10, 10, 1, "")
	layer.Add("layer2", 0, 0, 10, 10, 2, "")
	layer.Add("layer3", 0, 0, 10, 10, 3, "")

	// Test GetSortedLayerMemoryAliasSlice
	sortedLayers := layer.GetSortedLayerMemoryAliasSlice()
	if len(sortedLayers) != 3 {
		t.Error("Should have 3 sorted layers")
	}

	// Verify sorting
	for i := 1; i < len(sortedLayers); i++ {
		if sortedLayers[i].Value <= sortedLayers[i-1].Value {
			t.Error("Layers should be sorted by z-order")
		}
	}

	// Test SetHighestZOrderNumber
	layer.SetHighestZOrderNumber("layer1", "")
	layer1 := Layers.Get("layer1")
	if !layer1.IsTopmost {
		t.Error("Layer1 should be topmost after setting highest z-order")
	}
	DeleteAllLayers()
}

func TestLayerInstanceMethods(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create a test layer
	layerInstance := AddLayer(0, 0, 10, 10, 1, nil)

	// Test Clear
	layerInstance.Clear()

	// Test visibility
	layerInstance.SetIsVisible(false)
	layerEntry := Layers.Get(layerInstance.layerAlias)
	if layerEntry.IsVisible {
		t.Error("Layer should not be visible after SetIsVisible(false)")
	}

	// Test movement
	layerInstance.MoveLayerByAbsoluteValue(5, 5)
	layerEntry = Layers.Get(layerInstance.layerAlias)
	if layerEntry.ScreenXLocation != 5 || layerEntry.ScreenYLocation != 5 {
		t.Error("Layer should be moved to absolute position (5,5)")
	}

	layerInstance.MoveLayerByRelativeValue(1, 1)
	layerEntry = Layers.Get(layerInstance.layerAlias)
	if layerEntry.ScreenXLocation != 6 || layerEntry.ScreenYLocation != 6 {
		t.Error("Layer should be moved by relative position (1,1)")
	}

	// Test deletion
	layerInstance.DeleteLayer()
	if layerInstance.IsLayerExists() {
		t.Error("Layer should not exist after deletion")
	}
	DeleteAllLayers()
}

func TestLayerDrawingMethods(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create a test layer
	layerInstance := AddLayer(0, 0, 20, 20, 1, nil)
	styleEntry := types.NewTuiStyleEntry()

	// Test drawing methods
	layerInstance.DrawBorder(styleEntry, 0, 0, 10, 10, false)
	layerInstance.DrawHorizontalLine(styleEntry, 0, 0, 10, false)
	layerInstance.DrawVerticalLine(styleEntry, 0, 0, 10, false)
	layerInstance.DrawFrame(styleEntry, true, 0, 0, 10, 10, false)
	layerInstance.DrawWindow(styleEntry, 0, 0, 10, 10, false)
	layerInstance.DrawShadow(0, 0, 10, 10, 0.5)

	// Test filling methods
	layerInstance.FillArea("X", 0, 0, 10, 10)
	layerInstance.FillLayer("Y")
	layerInstance.DrawBar(0, 0, 10, "Z")

	// Test text methods
	layerInstance.Locate(0, 0)
	layerInstance.Print("Test Text")
	layerInstance.PrintDialog(0, 0, 10, 0, true, "Test Dialog")

	// Test color methods
	layerInstance.Color24Bit(constants.ColorWhite, constants.ColorBlack)
	DeleteAllLayers()
}

func TestLayerControlMethods(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create a test layer
	layerInstance := AddLayer(0, 0, 20, 20, 1, nil)
	styleEntry := types.NewTuiStyleEntry()

	// Test adding various controls
	button := layerInstance.AddButton("Test", styleEntry, 0, 0, 10, 1, true)
	if button.layerAlias != layerInstance.layerAlias {
		t.Error("Button should be added to correct layer")
	}

	checkbox := layerInstance.AddCheckbox("Test", styleEntry, 0, 0, false, true)
	if checkbox.layerAlias != layerInstance.layerAlias {
		t.Error("Checkbox should be added to correct layer")
	}

	label := layerInstance.AddLabel("Test", styleEntry, 0, 0, 10)
	if label.layerAlias != layerInstance.layerAlias {
		t.Error("Label should be added to correct layer")
	}

	progressBar := layerInstance.AddProgressBar("Test", styleEntry, 0, 0, 10, 1, 50, 100, false)
	if progressBar.layerAlias != layerInstance.layerAlias {
		t.Error("Progress bar should be added to correct layer")
	}

	radioButton := layerInstance.AddRadioButton("Test", styleEntry, 0, 0, 1, false)
	if radioButton.layerAlias != layerInstance.layerAlias {
		t.Error("Radio button should be added to correct layer")
	}

	scrollbar := layerInstance.AddScrollbar(styleEntry, 0, 0, 10, 100, 0, 1, false)
	if scrollbar.layerAlias != layerInstance.layerAlias {
		t.Error("Scrollbar should be added to correct layer")
	}

	textbox := layerInstance.AddTextbox(styleEntry, 0, 0, 10, 1, true)
	if textbox.layerAlias != layerInstance.layerAlias {
		t.Error("Textbox should be added to correct layer")
	}
	DeleteAllLayers()
}

func TestLayerGlobalMethods(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()
	// Test AddLayer
	layerInstance := AddLayer(0, 0, 10, 10, 1, nil)
	if !isLayerExists(layerInstance.layerAlias) {
		t.Error("Layer should exist after AddLayer")
	}

	// Test MoveLayerByAbsoluteValue
	MoveLayerByAbsoluteValue(layerInstance.layerAlias, 5, 5)
	layerEntry := Layers.Get(layerInstance.layerAlias)
	if layerEntry.ScreenXLocation != 5 || layerEntry.ScreenYLocation != 5 {
		t.Error("Layer should be moved to absolute position (5,5)")
	}

	// Test MoveLayerByRelativeValue
	MoveLayerByRelativeValue(layerInstance.layerAlias, 1, 1)
	layerEntry = Layers.Get(layerInstance.layerAlias)
	if layerEntry.ScreenXLocation != 6 || layerEntry.ScreenYLocation != 6 {
		t.Error("Layer should be moved by relative position (1,1)")
	}

	// Test DeleteLayer
	DeleteLayer(layerInstance)
	if isLayerExists(layerInstance.layerAlias) {
		t.Error("Layer should not exist after DeleteLayer")
	}

	// Test DeleteAllLayers
	layer.Add("test1", 0, 0, 10, 10, 1, "")
	layer.Add("test2", 0, 0, 10, 10, 2, "")
	DeleteAllLayers()
	if len(Layers.GetAllEntries()) != 0 {
		t.Error("All layers should be deleted")
	}

	// Test setLayerIsVisible
	layer.Add("test", 0, 0, 10, 10, 1, "")
	setLayerIsVisible("test", false)
	layerEntry = Layers.Get("test")
	if layerEntry.IsVisible {
		t.Error("Layer should not be visible after setLayerIsVisible(false)")
	}
	DeleteAllLayers()
}

func TestComplexLayerHierarchy(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create a complex layer hierarchy:
	// root1
	//   ├── child1
	//   │     ├── grandchild1
	//   │     └── grandchild2
	//   └── child2
	//         └── grandchild3
	// root2
	//   └── child3
	//         └── grandchild4

	// Create root layers
	layer.Add("root1", 0, 0, 20, 20, 1, "")
	layer.Add("root2", 0, 0, 20, 20, 2, "")

	// Create child layers
	layer.Add("child1", 0, 0, 15, 15, 1, "root1")
	layer.Add("child2", 0, 0, 15, 15, 2, "root1")
	layer.Add("child3", 0, 0, 15, 15, 1, "root2")

	// Create grandchild layers
	layer.Add("grandchild1", 0, 0, 10, 10, 1, "child1")
	layer.Add("grandchild2", 0, 0, 10, 10, 2, "child1")
	layer.Add("grandchild3", 0, 0, 10, 10, 1, "child2")
	layer.Add("grandchild4", 0, 0, 10, 10, 1, "child3")

	// Verify parent-child relationships
	if !layer.IsAParent("root1") {
		t.Error("root1 should be identified as parent")
	}
	if !layer.IsAParent("root2") {
		t.Error("root2 should be identified as parent")
	}
	if !layer.IsAParent("child1") {
		t.Error("child1 should be identified as parent")
	}
	if !layer.IsAParent("child2") {
		t.Error("child2 should be identified as parent")
	}
	if !layer.IsAParent("child3") {
		t.Error("child3 should be identified as parent")
	}

	// Verify root parent relationships
	rootParent, child := layer.GetRootParentLayerAlias("grandchild1", "")
	if rootParent != "root1" {
		t.Error("grandchild1 should have root1 as root parent")
	}
	if child != "grandchild1" {
		t.Error("child should be grandchild1")
	}

	rootParent, child = layer.GetRootParentLayerAlias("grandchild4", "")
	if rootParent != "root2" {
		t.Error("grandchild4 should have root2 as root parent")
	}
	if child != "grandchild4" {
		t.Error("child should be grandchild4")
	}

	// Delete a middle layer and verify children are properly deleted
	layer.Delete("child1")
	if isLayerExists("grandchild1") {
		t.Error("grandchild1 should be deleted when child1 is deleted")
	}
	if isLayerExists("grandchild2") {
		t.Error("grandchild2 should be deleted when child1 is deleted")
	}
	if !isLayerExists("root1") {
		t.Error("root1 should still exist after deleting child1")
	}
	if !isLayerExists("child2") {
		t.Error("child2 should still exist after deleting child1")
	}

	// Delete a root layer and verify all descendants are deleted
	layer.Delete("root2")
	if isLayerExists("child3") {
		t.Error("child3 should be deleted when root2 is deleted")
	}
	if isLayerExists("grandchild4") {
		t.Error("grandchild4 should be deleted when root2 is deleted")
	}

	// Verify remaining layers
	if !isLayerExists("root1") {
		t.Error("root1 should still exist after deleting root2")
	}
	if !isLayerExists("child2") {
		t.Error("child2 should still exist after deleting root2")
	}
	if !isLayerExists("grandchild3") {
		t.Error("grandchild3 should still exist after deleting root2")
	}
	DeleteAllLayers()
}

func TestComplexControlManagement(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create a layer hierarchy
	layer.Add("root", 0, 0, 20, 20, 1, "")
	layer.Add("child1", 0, 0, 15, 15, 1, "root")
	layer.Add("child2", 0, 0, 15, 15, 2, "root")

	// Create layer instances
	rootLayer := AddLayer(0, 0, 20, 20, 1, nil)
	child1Layer := AddLayer(0, 0, 15, 15, 1, &rootLayer)
	child2Layer := AddLayer(0, 0, 15, 15, 2, &rootLayer)

	// Create style entries
	styleEntry := types.NewTuiStyleEntry()

	// Add controls to root layer
	rootButton := rootLayer.AddButton("Root Button", styleEntry, 0, 0, 10, 1, true)
	rootLabel := rootLayer.AddLabel("Root Label", styleEntry, 0, 1, 10)
	rootCheckbox := rootLayer.AddCheckbox("Root Checkbox", styleEntry, 0, 2, false, true)

	// Add controls to child1 layer
	child1Button := child1Layer.AddButton("Child1 Button", styleEntry, 0, 0, 10, 1, true)
	child1Label := child1Layer.AddLabel("Child1 Label", styleEntry, 0, 1, 10)
	child1ProgressBar := child1Layer.AddProgressBar("Child1 Progress", styleEntry, 0, 2, 10, 1, 50, 100, false)

	// Add controls to child2 layer
	child2Button := child2Layer.AddButton("Child2 Button", styleEntry, 0, 0, 10, 1, true)
	child2RadioButton := child2Layer.AddRadioButton("Child2 Radio", styleEntry, 0, 1, 1, false)
	child2Scrollbar := child2Layer.AddScrollbar(styleEntry, 0, 2, 10, 100, 0, 1, false)

	// Verify controls are added to correct layers
	if rootButton.layerAlias != rootLayer.layerAlias {
		t.Error("Root button should be added to root layer")
	}
	if child1Button.layerAlias != child1Layer.layerAlias {
		t.Error("Child1 button should be added to child1 layer")
	}
	if child2Button.layerAlias != child2Layer.layerAlias {
		t.Error("Child2 button should be added to child2 layer")
	}

	// Delete child1 layer and verify its controls are deleted
	child1Layer.DeleteLayer()
	if Buttons.IsExists(child1Layer.layerAlias, child1Button.controlAlias) {
		t.Error("Child1 button should be deleted when child1 layer is deleted")
	}
	if Labels.IsExists(child1Layer.layerAlias, child1Label.controlAlias) {
		t.Error("Child1 label should be deleted when child1 layer is deleted")
	}
	if ProgressBars.IsExists(child1Layer.layerAlias, child1ProgressBar.controlAlias) {
		t.Error("Child1 progress bar should be deleted when child1 layer is deleted")
	}

	// Verify controls on other layers still exist
	if !Buttons.IsExists(rootLayer.layerAlias, rootButton.controlAlias) {
		t.Error("Root button should still exist after deleting child1 layer")
	}
	if !Labels.IsExists(rootLayer.layerAlias, rootLabel.controlAlias) {
		t.Error("Root label should still exist after deleting child1 layer")
	}
	if !Checkboxes.IsExists(rootLayer.layerAlias, rootCheckbox.controlAlias) {
		t.Error("Root checkbox should still exist after deleting child1 layer")
	}
	if !Buttons.IsExists(child2Layer.layerAlias, child2Button.controlAlias) {
		t.Error("Child2 button should still exist after deleting child1 layer")
	}
	if !RadioButtons.IsExists(child2Layer.layerAlias, child2RadioButton.controlAlias) {
		t.Error("Child2 radio button should still exist after deleting child1 layer")
	}
	if !ScrollBars.IsExists(child2Layer.layerAlias, child2Scrollbar.controlAlias) {
		t.Error("Child2 scrollbar should still exist after deleting child1 layer")
	}

	// Delete root layer and verify all controls are deleted
	rootLayer.DeleteLayer()
	if Buttons.IsExists(rootLayer.layerAlias, rootButton.controlAlias) {
		t.Error("Root button should be deleted when root layer is deleted")
	}
	if Labels.IsExists(rootLayer.layerAlias, rootLabel.controlAlias) {
		t.Error("Root label should be deleted when root layer is deleted")
	}
	if Checkboxes.IsExists(rootLayer.layerAlias, rootCheckbox.controlAlias) {
		t.Error("Root checkbox should be deleted when root layer is deleted")
	}
	if Buttons.IsExists(child2Layer.layerAlias, child2Button.controlAlias) {
		t.Error("Child2 button should be deleted when root layer is deleted")
	}
	if RadioButtons.IsExists(child2Layer.layerAlias, child2RadioButton.controlAlias) {
		t.Error("Child2 radio button should be deleted when root layer is deleted")
	}
	if ScrollBars.IsExists(child2Layer.layerAlias, child2Scrollbar.controlAlias) {
		t.Error("Child2 scrollbar should be deleted when root layer is deleted")
	}
	DeleteAllLayers()
}

func TestLayerAndControlPropertyStability(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create a layer hierarchy
	rootLayer := AddLayer(0, 0, 20, 20, 1, nil)
	childLayer := AddLayer(0, 0, 15, 15, 1, &rootLayer)

	// Create style entries
	styleEntry := types.NewTuiStyleEntry()

	// Add controls to layers
	rootButton := rootLayer.AddButton("Root Button", styleEntry, 0, 0, 10, 1, true)
	childButton := childLayer.AddButton("Child Button", styleEntry, 0, 0, 10, 1, true)

	// Modify control properties
	rootButton.SetEnabled(false)
	childButton.SetEnabled(false)

	// Verify control properties are maintained
	rootButtonEntry := Buttons.Get(rootLayer.layerAlias, rootButton.controlAlias)
	if rootButtonEntry.IsEnabled {
		t.Error("Root button should be disabled")
	}

	childButtonEntry := Buttons.Get(childLayer.layerAlias, childButton.controlAlias)
	if childButtonEntry.IsEnabled {
		t.Error("Child button should be disabled")
	}

	// Move layers and verify control positions are maintained
	rootLayer.MoveLayerByAbsoluteValue(5, 5)
	childLayer.MoveLayerByAbsoluteValue(2, 2)

	// Verify layer positions
	rootLayerEntry := Layers.Get(rootLayer.layerAlias)
	if rootLayerEntry.ScreenXLocation != 5 || rootLayerEntry.ScreenYLocation != 5 {
		t.Error("Root layer should be moved to position (5,5)")
	}

	childLayerEntry := Layers.Get(childLayer.layerAlias)
	if childLayerEntry.ScreenXLocation != 2 || childLayerEntry.ScreenYLocation != 2 {
		t.Error("Child layer should be moved to position (2,2)")
	}

	// Change z-order and verify
	layer.SetHighestZOrderNumber(childLayer.layerAlias, rootLayer.layerAlias)
	childLayerEntry = Layers.Get(childLayer.layerAlias)
	if !childLayerEntry.IsTopmost {
		t.Error("Child layer should be topmost after setting highest z-order")
	}

	// Delete and recreate layers with same aliases
	rootLayer.DeleteLayer()

	newRootLayer := AddLayer(0, 0, 20, 20, 1, nil)
	newChildLayer := AddLayer(0, 0, 15, 15, 1, &newRootLayer)

	// Verify new layers have correct properties
	newRootLayerEntry := Layers.Get(newRootLayer.layerAlias)
	if newRootLayerEntry.ScreenXLocation != 0 || newRootLayerEntry.ScreenYLocation != 0 {
		t.Error("New root layer should have default position (0,0)")
	}

	newChildLayerEntry := Layers.Get(newChildLayer.layerAlias)
	if newChildLayerEntry.ParentAlias != newRootLayer.layerAlias {
		t.Error("New child layer should have correct parent")
	}
	DeleteAllLayers()
}

func TestLayerAndControlMemoryLeaks(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create a large number of layers and controls
	for i := 0; i < 100; i++ {
		// Create root layer
		rootLayer := AddLayer(0, 0, 20, 20, i, nil)
		styleEntry := types.NewTuiStyleEntry()

		// Add controls to root layer
		rootLayer.AddButton("Button", styleEntry, 0, 0, 10, 1, true)
		rootLayer.AddLabel("Label", styleEntry, 0, 1, 10)
		rootLayer.AddCheckbox("Checkbox", styleEntry, 0, 2, false, true)

		// Create child layer
		childLayer := AddLayer(0, 0, 15, 15, i, &rootLayer)

		// Add controls to child layer
		childLayer.AddButton("Child Button", styleEntry, 0, 0, 10, 1, true)
		childLayer.AddLabel("Child Label", styleEntry, 0, 1, 10)
		childLayer.AddProgressBar("Progress", styleEntry, 0, 2, 10, 1, 50, 100, false)

		// Delete layers and verify all controls are deleted
		rootLayer.DeleteLayer()

		// Verify layers are deleted
		if isLayerExists(rootLayer.layerAlias) {
			t.Errorf("Root layer %d should be deleted", i)
		}
		if isLayerExists(childLayer.layerAlias) {
			t.Errorf("Child layer %d should be deleted", i)
		}

		// Verify no controls remain
		if len(Buttons.GetAllEntriesOverall()) > 0 {
			t.Errorf("Buttons should be deleted for layer %d", i)
		}
		if len(Labels.GetAllEntriesOverall()) > 0 {
			t.Errorf("Labels should be deleted for layer %d", i)
		}
		if len(Checkboxes.GetAllEntriesOverall()) > 0 {
			t.Errorf("Checkboxes should be deleted for layer %d", i)
		}
		if len(ProgressBars.GetAllEntriesOverall()) > 0 {
			t.Errorf("Progress bars should be deleted for layer %d", i)
		}
	}

	// Verify no layers or controls remain
	if len(Layers.GetAllEntries()) > 0 {
		t.Errorf("No layers should remain after deletion")
	}
	if len(Buttons.GetAllEntriesOverall()) > 0 {
		t.Errorf("No buttons should remain after deletion")
	}
	if len(Labels.GetAllEntriesOverall()) > 0 {
		t.Errorf("No labels should remain after deletion")
	}
	if len(Checkboxes.GetAllEntriesOverall()) > 0 {
		t.Errorf("No checkboxes should remain after deletion")
	}
	if len(ProgressBars.GetAllEntriesOverall()) > 0 {
		t.Errorf("No progress bars should remain after deletion")
	}
	DeleteAllLayers()
}

func TestComplexInterleavedOperations(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create a layer hierarchy
	root1Layer := AddLayer(0, 0, 20, 20, 1, nil)
	root2Layer := AddLayer(0, 0, 20, 20, 2, nil)
	child1Layer := AddLayer(0, 0, 15, 15, 1, &root1Layer)
	child2Layer := AddLayer(0, 0, 15, 15, 2, &root1Layer)
	child3Layer := AddLayer(0, 0, 15, 15, 1, &root2Layer)
	grandchild1Layer := AddLayer(0, 0, 10, 10, 1, &child1Layer)

	// Create style entries with unique properties
	styleEntry1 := types.NewTuiStyleEntry()
	styleEntry1.ButtonForegroundColor = constants.ColorRed
	styleEntry1.ButtonBackgroundColor = constants.ColorBlue
	styleEntry1.ButtonRaisedColor = constants.ColorGreen

	styleEntry2 := types.NewTuiStyleEntry()
	styleEntry2.ButtonForegroundColor = constants.ColorYellow
	styleEntry2.ButtonBackgroundColor = constants.ColorMagenta
	styleEntry2.ButtonRaisedColor = constants.ColorCyan

	selectionEntry1 := types.NewSelectionEntry()
	selectionEntry1.SelectionValue = []string{"Option A", "Option B", "Option C"}
	selectionEntry1.SelectionAlias = []string{"optA", "optB", "optC"}

	selectionEntry2 := types.NewSelectionEntry()
	selectionEntry2.SelectionValue = []string{"Choice 1", "Choice 2", "Choice 3"}
	selectionEntry2.SelectionAlias = []string{"ch1", "ch2", "ch3"}

	// Phase 1: Add controls to root1 layer with unique properties
	root1Button := root1Layer.AddButton("Root1 Button", styleEntry1, 0, 0, 10, 1, true)
	root1Checkbox := root1Layer.AddCheckbox("Root1 Checkbox", styleEntry1, 0, 2, true, true)
	root1ProgressBar := root1Layer.AddProgressBar("Root1 Progress", styleEntry1, 0, 4, 10, 1, 75, 100, false)
	root1Layer.AddDropdown(styleEntry1, selectionEntry1, 0, 3, 3, 10, 1)
	root1Layer.AddRadioButton("Root1 Radio", styleEntry1, 0, 5, 2, true)
	root1Layer.AddScrollbar(styleEntry1, 0, 6, 10, 100, 25, 1, false)
	root1Layer.AddSelector(styleEntry1, selectionEntry1, 0, 7, 3, 10, 2, 0, 0, true)
	root1Layer.AddTextField(styleEntry1, 0, 8, 10, 20, true, "Root1 Text", true)
	root1Layer.AddTextbox(styleEntry1, 0, 9, 10, 3, true)
	root1Layer.AddTooltip("Root1 Tooltip", styleEntry1, 0, 0, 10, 1, 0, 0, 10, 1, true, true, 2000)

	// Phase 2: Delete child1 layer and verify
	child1Layer.DeleteLayer()
	if isLayerExists(child1Layer.layerAlias) {
		t.Error("Child1 layer should be deleted")
	}
	if isLayerExists(grandchild1Layer.layerAlias) {
		t.Error("Grandchild1 layer should be deleted when child1 is deleted")
	}

	// Phase 3: Add new controls to remaining layers with unique properties
	child2Button := child2Layer.AddButton("Child2 Button", styleEntry2, 0, 0, 10, 1, true)
	child2Checkbox := child2Layer.AddCheckbox("Child2 Checkbox", styleEntry2, 0, 2, false, true)
	child2ProgressBar := child2Layer.AddProgressBar("Child2 Progress", styleEntry2, 0, 4, 10, 1, 25, 100, false)
	child2Layer.AddDropdown(styleEntry2, selectionEntry2, 0, 3, 3, 10, 2)
	child2Layer.AddRadioButton("Child2 Radio", styleEntry2, 0, 5, 1, false)
	child2Layer.AddScrollbar(styleEntry2, 0, 6, 10, 100, 50, 1, false)
	child2Layer.AddSelector(styleEntry2, selectionEntry2, 0, 7, 3, 10, 1, 0, 0, true)
	child2Layer.AddTextField(styleEntry2, 0, 8, 10, 20, false, "Child2 Text", true)
	child2Layer.AddTextbox(styleEntry2, 0, 9, 10, 3, true)
	child2Layer.AddTooltip("Child2 Tooltip", styleEntry2, 0, 0, 10, 1, 0, 0, 10, 1, false, true, 1500)

	// Verify child2 controls after creation
	child2ButtonEntry := Buttons.Get(child2Layer.layerAlias, child2Button.controlAlias)
	if child2ButtonEntry.StyleEntry.ButtonForegroundColor != constants.ColorYellow {
		t.Error("Child2 button should maintain its yellow foreground color")
	}
	child2CheckboxEntry := Checkboxes.Get(child2Layer.layerAlias, child2Checkbox.controlAlias)
	if child2CheckboxEntry.IsSelected {
		t.Error("Child2 checkbox should be unchecked")
	}
	child2ProgressBarEntry := ProgressBars.Get(child2Layer.layerAlias, child2ProgressBar.controlAlias)
	if child2ProgressBarEntry.Value != 25 {
		t.Error("Child2 progress bar should be at 25%")
	}

	// Phase 4: Create new child layer under root1 with unique properties
	newChildLayer := AddLayer(0, 0, 15, 15, 3, &root1Layer)
	newChildLayer.AddButton("New Child Button", styleEntry1, 0, 0, 10, 1, true)
	newChildLayer.AddCheckbox("New Child Checkbox", styleEntry1, 0, 2, true, true)
	newChildLayer.AddProgressBar("New Child Progress", styleEntry1, 0, 4, 10, 1, 50, 100, false)
	newChildLayer.AddDropdown(styleEntry1, selectionEntry1, 0, 3, 3, 10, 0)
	newChildLayer.AddRadioButton("New Child Radio", styleEntry1, 0, 5, 3, true)
	newChildLayer.AddScrollbar(styleEntry1, 0, 6, 10, 100, 75, 1, false)
	newChildLayer.AddSelector(styleEntry1, selectionEntry1, 0, 7, 3, 10, 1, 0, 0, true)
	newChildLayer.AddTextField(styleEntry1, 0, 8, 10, 20, true, "New Child Text", true)
	newChildLayer.AddTextbox(styleEntry1, 0, 9, 10, 3, true)
	newChildLayer.AddTooltip("New Child Tooltip", styleEntry1, 0, 0, 10, 1, 0, 0, 10, 1, true, true, 3000)

	// Phase 5: Delete root2 and verify
	root2Layer.DeleteLayer()
	if isLayerExists(root2Layer.layerAlias) {
		t.Error("Root2 layer should be deleted")
	}
	if isLayerExists(child3Layer.layerAlias) {
		t.Error("Child3 layer should be deleted when root2 is deleted")
	}

	// Phase 6: Add new root layer and children with unique properties
	root3Layer := AddLayer(0, 0, 20, 20, 3, nil)
	root3Child1 := AddLayer(0, 0, 15, 15, 1, &root3Layer)
	root3Child2 := AddLayer(0, 0, 15, 15, 2, &root3Layer)

	// Add controls to root3 layer with unique properties
	root3Button := root3Layer.AddButton("Root3 Button", styleEntry2, 0, 0, 10, 1, true)
	root3Checkbox := root3Layer.AddCheckbox("Root3 Checkbox", styleEntry2, 0, 2, true, true)
	root3ProgressBar := root3Layer.AddProgressBar("Root3 Progress", styleEntry2, 0, 4, 10, 1, 90, 100, false)
	root3Layer.AddDropdown(styleEntry2, selectionEntry2, 0, 3, 3, 10, 1)
	root3Layer.AddRadioButton("Root3 Radio", styleEntry2, 0, 5, 2, true)
	root3Layer.AddScrollbar(styleEntry2, 0, 6, 10, 100, 60, 1, false)
	root3Layer.AddSelector(styleEntry2, selectionEntry2, 0, 7, 3, 10, 2, 0, 0, true)
	root3Layer.AddTextField(styleEntry2, 0, 8, 10, 20, true, "Root3 Text", true)
	root3Layer.AddTextbox(styleEntry2, 0, 9, 10, 3, true)
	root3Layer.AddTooltip("Root3 Tooltip", styleEntry2, 0, 0, 10, 1, 0, 0, 10, 1, true, true, 2500)

	// Add controls to root3Child1 with unique properties
	root3Child1.AddButton("Root3 Child1 Button", styleEntry1, 0, 0, 10, 1, true)
	root3Child1.AddCheckbox("Root3 Child1 Checkbox", styleEntry1, 0, 2, false, true)
	root3Child1.AddProgressBar("Root3 Child1 Progress", styleEntry1, 0, 4, 10, 1, 40, 100, false)
	root3Child1.AddDropdown(styleEntry1, selectionEntry1, 0, 3, 3, 10, 2)
	root3Child1.AddRadioButton("Root3 Child1 Radio", styleEntry1, 0, 5, 1, false)
	root3Child1.AddScrollbar(styleEntry1, 0, 6, 10, 100, 30, 1, false)
	root3Child1.AddSelector(styleEntry1, selectionEntry1, 0, 7, 3, 10, 1, 0, 0, true)
	root3Child1.AddTextField(styleEntry1, 0, 8, 10, 20, false, "Root3 Child1 Text", true)
	root3Child1.AddTextbox(styleEntry1, 0, 9, 10, 3, true)
	root3Child1.AddTooltip("Root3 Child1 Tooltip", styleEntry1, 0, 0, 10, 1, 0, 0, 10, 1, false, true, 3500)

	// Phase 7: Move some layers and verify control properties remain intact
	root1Layer.MoveLayerByAbsoluteValue(5, 5)
	newChildLayer.MoveLayerByAbsoluteValue(2, 2)
	root3Layer.MoveLayerByAbsoluteValue(10, 10)

	// Verify control properties after moving layers
	root1ButtonEntry := Buttons.Get(root1Layer.layerAlias, root1Button.controlAlias)
	if root1ButtonEntry.StyleEntry.ButtonForegroundColor != constants.ColorRed {
		t.Error("Root1 button should maintain its red foreground color after move")
	}
	root1CheckboxEntry := Checkboxes.Get(root1Layer.layerAlias, root1Checkbox.controlAlias)
	if !root1CheckboxEntry.IsSelected {
		t.Error("Root1 checkbox should remain checked after move")
	}
	root1ProgressBarEntry := ProgressBars.Get(root1Layer.layerAlias, root1ProgressBar.controlAlias)
	if root1ProgressBarEntry.Value != 75 {
		t.Error("Root1 progress bar should maintain 75% after move")
	}

	// Phase 8: Delete root1 and verify
	root1Layer.DeleteLayer()
	if isLayerExists(root1Layer.layerAlias) {
		t.Error("Root1 layer should be deleted")
	}
	if isLayerExists(child2Layer.layerAlias) {
		t.Error("Child2 layer should be deleted when root1 is deleted")
	}
	if isLayerExists(newChildLayer.layerAlias) {
		t.Error("New child layer should be deleted when root1 is deleted")
	}

	// Phase 9: Verify root3 and its children still exist with correct properties
	if !isLayerExists(root3Layer.layerAlias) {
		t.Error("Root3 layer should still exist")
	}
	if !isLayerExists(root3Child1.layerAlias) {
		t.Error("Root3 child1 layer should still exist")
	}
	if !isLayerExists(root3Child2.layerAlias) {
		t.Error("Root3 child2 layer should still exist")
	}

	// Verify root3 layer position and control properties
	root3Entry := Layers.Get(root3Layer.layerAlias)
	if root3Entry.ScreenXLocation != 10 || root3Entry.ScreenYLocation != 10 {
		t.Error("Root3 layer should maintain its position")
	}

	root3ButtonEntry := Buttons.Get(root3Layer.layerAlias, root3Button.controlAlias)
	if root3ButtonEntry.StyleEntry.ButtonForegroundColor != constants.ColorYellow {
		t.Error("Root3 button should maintain its yellow foreground color")
	}
	root3CheckboxEntry := Checkboxes.Get(root3Layer.layerAlias, root3Checkbox.controlAlias)
	if !root3CheckboxEntry.IsSelected {
		t.Error("Root3 checkbox should remain checked")
	}
	root3ProgressBarEntry := ProgressBars.Get(root3Layer.layerAlias, root3ProgressBar.controlAlias)
	if root3ProgressBarEntry.Value != 90 {
		t.Error("Root3 progress bar should maintain 90%")
	}

	// Phase 10: Add new controls to root3's children with unique properties
	root3Child2Button := root3Child2.AddButton("Root3 Child2 Button", styleEntry2, 0, 0, 10, 1, true)
	root3Child2Checkbox := root3Child2.AddCheckbox("Root3 Child2 Checkbox", styleEntry2, 0, 2, true, true)
	root3Child2ProgressBar := root3Child2.AddProgressBar("Root3 Child2 Progress", styleEntry2, 0, 4, 10, 1, 60, 100, false)
	root3Child2.AddDropdown(styleEntry2, selectionEntry2, 0, 3, 3, 10, 0)
	root3Child2.AddRadioButton("Root3 Child2 Radio", styleEntry2, 0, 5, 1, true)
	root3Child2.AddScrollbar(styleEntry2, 0, 6, 10, 100, 40, 1, false)
	root3Child2.AddSelector(styleEntry2, selectionEntry2, 0, 7, 3, 10, 1, 0, 0, true)
	root3Child2.AddTextField(styleEntry2, 0, 8, 10, 20, true, "Root3 Child2 Text", true)
	root3Child2.AddTextbox(styleEntry2, 0, 9, 10, 3, true)
	root3Child2.AddTooltip("Root3 Child2 Tooltip", styleEntry2, 0, 0, 10, 1, 0, 0, 10, 1, true, true, 4000)

	// Verify root3Child2 controls after creation
	root3Child2ButtonEntry := Buttons.Get(root3Child2.layerAlias, root3Child2Button.controlAlias)
	if root3Child2ButtonEntry.StyleEntry.ButtonForegroundColor != constants.ColorYellow {
		t.Error("Root3 Child2 button should have yellow foreground color")
	}
	root3Child2CheckboxEntry := Checkboxes.Get(root3Child2.layerAlias, root3Child2Checkbox.controlAlias)
	if !root3Child2CheckboxEntry.IsSelected {
		t.Error("Root3 Child2 checkbox should be checked")
	}
	root3Child2ProgressBarEntry := ProgressBars.Get(root3Child2.layerAlias, root3Child2ProgressBar.controlAlias)
	if root3Child2ProgressBarEntry.Value != 60 {
		t.Error("Root3 Child2 progress bar should be at 60%")
	}

	// Phase 11: Delete root3's first child and verify remaining controls
	root3Child1.DeleteLayer()
	if isLayerExists(root3Child1.layerAlias) {
		t.Error("Root3 child1 layer should be deleted")
	}
	if !isLayerExists(root3Child2.layerAlias) {
		t.Error("Root3 child2 layer should still exist")
	}

	// Verify root3 and root3Child2 controls still have correct properties
	root3ButtonEntry = Buttons.Get(root3Layer.layerAlias, root3Button.controlAlias)
	if root3ButtonEntry.StyleEntry.ButtonForegroundColor != constants.ColorYellow {
		t.Error("Root3 button should maintain its yellow foreground color after child deletion")
	}
	root3Child2ButtonEntry = Buttons.Get(root3Child2.layerAlias, root3Child2Button.controlAlias)
	if root3Child2ButtonEntry.StyleEntry.ButtonForegroundColor != constants.ColorYellow {
		t.Error("Root3 Child2 button should maintain its yellow foreground color after sibling deletion")
	}

	DeleteAllLayers()
}
