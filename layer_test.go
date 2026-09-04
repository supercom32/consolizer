package consolizer

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
)

/*
TestLayerInitialization is a test which verifies that the layer memory is correctly reinitialized and that
the Layers manager is not nil after initialization.

Example:

	Expected Inputs:
	    Terminal initialized to 80x25.
	Expected Outputs:
	    Layers manager is not nil, and no errors occur during initialization.
*/
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

/*
TestLayerAdd is a test which verifies that layers can be correctly added to the system and that invalid
dimensions correctly trigger a panic.

Example:

	Expected Inputs:
	    A valid layer with 10x10 dimensions and an invalid layer with 0 width.
	Expected Outputs:
	    The valid layer exists in the system, and adding the invalid layer triggers a panic.
*/
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

/*
TestLayerDelete is a test which verifies that layers can be correctly deleted and that deleting a child does
not affect the parent.

Example:

	Expected Inputs:
	    A parent layer and a child layer.
	Expected Outputs:
	    The child layer is successfully deleted while the parent layer remains.
*/
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

/*
TestLayerParentChild is a test which verifies the hierarchical relationship between parent and child layers.

Example:

	Expected Inputs:
	    A parent layer and a child layer linked to it.
	Expected Outputs:
	    The parent is correctly identified as a parent, and GetRootParentAlias returns the correct root.
*/
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
	rootParent, child := layer.GetRootParentAlias("child", "")
	if rootParent != "parent" {
		t.Error("Root parent should be 'parent'")
	}
	if child != "child" {
		t.Error("Child should be 'child'")
	}
	DeleteAllLayers()
}

/*
TestLayerGetAbsoluteLocation is a test which verifies that layer.GetAbsoluteLocation correctly resolves a
layer's screen position by summing its own offset with every ancestor's offset in the parent chain.

Example:

	Expected Inputs:
	    A three-level hierarchy (root, child, grandchild), each offset from its own parent.
	Expected Outputs:
	    The root's absolute location equals its own offset, and each descendant's absolute location equals
	    the sum of its own offset plus every ancestor's offset.
*/
func TestLayerGetAbsoluteLocation(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	layer.Add("root", 2, 3, 20, 20, 1, "")
	layer.Add("child", 4, 5, 15, 15, 1, "root")
	layer.Add("grandchild", 1, 1, 10, 10, 1, "child")

	rootX, rootY := layer.GetAbsoluteLocation("root")
	if rootX != 2 || rootY != 3 {
		t.Errorf("Root absolute location should be (2, 3), got (%d, %d)", rootX, rootY)
	}

	childX, childY := layer.GetAbsoluteLocation("child")
	if childX != 6 || childY != 8 {
		t.Errorf("Child absolute location should be (6, 8), got (%d, %d)", childX, childY)
	}

	grandchildX, grandchildY := layer.GetAbsoluteLocation("grandchild")
	if grandchildX != 7 || grandchildY != 9 {
		t.Errorf("Grandchild absolute location should be (7, 9), got (%d, %d)", grandchildX, grandchildY)
	}
	DeleteAllLayers()
}

/*
TestLayerInstanceLocationMethods is a test which verifies that GetRelativeLocation, GetAbsoluteLocation, and
the deprecated GetLocation alias each return the correct coordinates for a nested layer instance.

Example:

	Expected Inputs:
	    A parent layer instance and a child layer instance offset from it.
	Expected Outputs:
	    GetRelativeLocation returns the child's parent-local offset, GetAbsoluteLocation returns the child's
	    offset summed with the parent's offset, and GetLocation matches GetRelativeLocation.
*/
func TestLayerInstanceLocationMethods(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	parentInstance := AddLayer(10, 5, 20, 20, 1, nil)
	childInstance := AddLayer(3, 4, 10, 10, 1, parentInstance)

	relativeX, relativeY := childInstance.GetRelativeLocation()
	if relativeX != 3 || relativeY != 4 {
		t.Errorf("Child relative location should be (3, 4), got (%d, %d)", relativeX, relativeY)
	}

	absoluteX, absoluteY := childInstance.GetAbsoluteLocation()
	if absoluteX != 13 || absoluteY != 9 {
		t.Errorf("Child absolute location should be (13, 9), got (%d, %d)", absoluteX, absoluteY)
	}

	legacyX, legacyY := childInstance.GetLocation()
	if legacyX != relativeX || legacyY != relativeY {
		t.Errorf("GetLocation should match GetRelativeLocation (%d, %d), got (%d, %d)", relativeX, relativeY, legacyX, legacyY)
	}

	parentAbsoluteX, parentAbsoluteY := parentInstance.GetAbsoluteLocation()
	if parentAbsoluteX != 10 || parentAbsoluteY != 5 {
		t.Errorf("Root parent absolute location should equal its own offset (10, 5), got (%d, %d)", parentAbsoluteX, parentAbsoluteY)
	}
	DeleteAllLayers()
}

/*
TestLayerZOrder is a test which verifies that layers are correctly sorted by their z-order rendering
priority.

Example:

	Expected Inputs:
	    Three layers with z-orders 1, 2, and 3.
	Expected Outputs:
	    The layers are returned in the correct z-order sequence, and SetHighestZOrderNumber correctly updates the topmost layer.
*/
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

/*
TestLayerInstanceMethods is a test which verifies the behavior of various methods on a layer instance, such
as visibility, movement, and deletion.

Example:

	Expected Inputs:
	    A layer instance subjected to Clear, SetIsVisible, and movement commands.
	Expected Outputs:
	    The layer's state (visibility, position) matches the values specified in the commands.
*/
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
	layerInstance.Delete()
	if layerInstance.IsExists() {
		t.Error("Layer should not exist after deletion")
	}
	DeleteAllLayers()
}

/*
TestLayerDrawingMethods is a test which verifies that various drawing methods (borders, lines, frames,
windows, etc.) on a layer instance can be called without errors.

Example:

	Expected Inputs:
	    A sequence of drawing commands (DrawBorder, DrawFrame, etc.) on a layer instance.
	Expected Outputs:
	    All drawing commands complete without errors or panics.
*/
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
	layerInstance.DrawBar(styleEntry, 0, 0, 10, "Z")

	// Test text methods
	layerInstance.Locate(0, 0)
	layerInstance.Print("Test Text")
	layerInstance.PrintDialog(0, 0, 10, 0, true, "Test Dialog")

	// Test color methods
	layerInstance.Color24Bit(constants.ColorWhite, constants.ColorBlack)
	DeleteAllLayers()
}

/*
TestPrintMethod is a test which verifies that text printing, space preservation, markup handling, and word
wrapping are all functioning correctly.

Example:

	Expected Inputs:
	    Strings containing markup tags, multiple spaces, and long sentences for wrapping.
	Expected Outputs:
	    The character memory reflects the correctly rendered text with applied colors and proper line wrapping.
*/
func TestPrintMethod(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create a test layer
	layerInstance := AddLayer(0, 0, 40, 20, 1, nil)

	// Create a test style for markup testing
	redStyle := NewTextStyle()
	redStyle.ForegroundColor = GetRGBColor(255, 0, 0) // Red
	AddTextStyle("red", redStyle)

	blueStyle := NewTextStyle()
	blueStyle.ForegroundColor = GetRGBColor(0, 0, 255) // Blue
	AddTextStyle("blue", blueStyle)

	// Test 1: Basic printing
	layerInstance.Locate(0, 0)
	layerInstance.Print("Basic text")
	layerEntry := Layers.Get(layerInstance.layerAlias)

	// Verify basic text was printed correctly
	for i := 0; i < 9; i++ {
		expectedChar := []rune("Basic text")[i]
		actualChar := layerEntry.CharacterMemory[0][i].Character
		if actualChar != expectedChar {
			t.Errorf("Basic text printing failed. Expected '%c' at position %d, got '%c'", expectedChar, i, actualChar)
		}
	}

	// Test 2: Space preservation
	layerInstance.Locate(0, 1)
	layerInstance.Print("Text with   multiple    spaces")

	// Verify spaces are preserved
	expectedText := "Text with   multiple    spaces"
	for i := 0; i < len(expectedText); i++ {
		expectedChar := []rune(expectedText)[i]
		actualChar := layerEntry.CharacterMemory[1][i].Character
		if actualChar != expectedChar {
			t.Errorf("Space preservation failed. Expected '%c' at position %d, got '%c'", expectedChar, i, actualChar)
		}
	}

	// Test 3: Markup handling with PrintMarkup
	layerInstance.PrintMarkup(0, 2, 40, "Normal {{red}}Red{{}} and {{blue}}Blue{{}} text")

	// Verify normal text color
	if layerEntry.CharacterMemory[2][0].AttributeEntry.ForegroundColor != layerEntry.DefaultAttribute.ForegroundColor {
		t.Errorf("Markup handling failed. Normal text should have default color")
	}

	// Verify red text color (character 'R' at position 7)
	redColor := GetRGBColor(255, 0, 0)
	if layerEntry.CharacterMemory[2][7].AttributeEntry.ForegroundColor != redColor {
		t.Errorf("Markup handling failed. Text should be red")
	}

	// Verify blue text color (character 'B' at position 16)
	blueColor := GetRGBColor(0, 0, 255)
	if layerEntry.CharacterMemory[2][16].AttributeEntry.ForegroundColor != blueColor {
		t.Errorf("Markup handling failed. Text should be blue")
	}

	// Test 4: Missing/broken tags with PrintMarkup
	layerInstance.PrintMarkup(0, 3, 40, "Text with {{redincomplete tag")
	layerEntry.GetBasicAnsiStringAsBase64()

	// Verify text with incomplete tag is printed correctly
	expectedText = "Text with {{redincomplete tag"
	for i := 0; i < len(expectedText); i++ {
		expectedChar := []rune(expectedText)[i]
		actualChar := layerEntry.CharacterMemory[3][i].Character
		if actualChar != expectedChar {
			t.Errorf("Missing tag handling failed. Expected '%c' at position %d, got '%c'", expectedChar, i, actualChar)
		}
	}

	// Test 5: Empty tags with PrintMarkup
	layerInstance.PrintMarkup(0, 4, 40, "Text with {{}}empty tag")

	// Verify text with empty tag is printed with default style
	if layerEntry.CharacterMemory[4][10].AttributeEntry.ForegroundColor != layerEntry.DefaultAttribute.ForegroundColor {
		t.Errorf("Empty tag handling failed. Text should have default color")
	}

	// Test 6: Long text with word wrapping
	longText := "This is a very long text that should wrap to the next line when printed with word wrapping enabled"
	layerInstance.PrintMarkup(0, 5, 20, longText)

	// Verify text wrapping (check if text continues on next line)
	if layerEntry.CharacterMemory[6][0].Character == 0 {
		t.Errorf("Word wrapping failed. Text should continue on next line")
	}

	DeleteAllLayers()
}

/*
TestLayerControlMethods is a test which verifies that various controls (buttons, checkboxes, labels, etc.)
can be correctly added to a layer and assigned the correct layer alias.

Example:

	Expected Inputs:
	    A series of Add commands for different UI controls.
	Expected Outputs:
	    The created controls contain the correct layer alias and match the input parameters.
*/
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

	progressBar := layerInstance.AddProgressBar("Test", styleEntry, 0, 0, 10, 1, false, 50, 100, false)
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

/*
TestLayerGlobalMethods is a test which verifies the behavior of global layer management functions like
AddLayer, MoveLayerByAbsoluteValue, MoveLayerByRelativeValue, and DeleteAllLayers.

Example:

	Expected Inputs:
	    Global commands for adding, moving, and deleting layers.
	Expected Outputs:
	    The system-wide layer memory correctly reflects the addition, movement, and eventual clearing of all layers.
*/
func TestLayerGlobalMethods(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()
	// Test AddLayer
	layerInstance := AddLayer(0, 0, 10, 10, 1, nil)
	layerAlias := layerInstance.layerAlias
	if !isLayerExists(layerAlias) {
		t.Error("Layer should exist after AddLayer")
	}

	// Test layerInstance.MoveLayerByAbsoluteValue
	layerInstance.MoveLayerByAbsoluteValue(5, 5)
	layerEntry := Layers.Get(layerAlias)
	if layerEntry.ScreenXLocation != 5 || layerEntry.ScreenYLocation != 5 {
		t.Error("Layer should be moved to absolute position (5,5)")
	}

	// Test layerInstance.MoveLayerByRelativeValue
	layerInstance.MoveLayerByRelativeValue(1, 1)
	layerEntry = Layers.Get(layerAlias)
	if layerEntry.ScreenXLocation != 6 || layerEntry.ScreenYLocation != 6 {
		t.Error("Layer should be moved by relative position (1,1)")
	}

	// Test layerInstance.Delete()
	layerInstance.Delete()
	if isLayerExists(layerAlias) {
		t.Error("Layer should not exist after Delete()")
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

/*
TestComplexLayerHierarchy is a test which verifies that complex parent-child relationships between layers
are correctly managed, especially during deletion of intermediate layers.

Example:

	Expected Inputs:
	    A deep hierarchy of root, child, and grandchild layers.
	Expected Outputs:
	    Deleting an intermediate parent correctly removes all its descendants while leaving root and sibling branches intact.
*/
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
	rootParent, child := layer.GetRootParentAlias("grandchild1", "")
	if rootParent != "root1" {
		t.Error("grandchild1 should have root1 as root parent")
	}
	if child != "grandchild1" {
		t.Error("child should be grandchild1")
	}

	rootParent, child = layer.GetRootParentAlias("grandchild4", "")
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

/*
TestComplexControlManagement is a test which verifies that controls on different layers are correctly
isolated and that deleting a layer correctly removes only its associated controls.

Example:

	Expected Inputs:
	    A hierarchy of layers each containing unique UI controls.
	Expected Outputs:
	    Deleting a parent layer correctly cleans up all its child layers and their respective controls.
*/
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
	child1Layer := AddLayer(0, 0, 15, 15, 1, rootLayer)
	child2Layer := AddLayer(0, 0, 15, 15, 2, rootLayer)

	// Create style entries
	styleEntry := types.NewTuiStyleEntry()

	// Add controls to root layer
	rootButton := rootLayer.AddButton("Root Button", styleEntry, 0, 0, 10, 1, true)
	rootLabel := rootLayer.AddLabel("Root Label", styleEntry, 0, 1, 10)
	rootCheckbox := rootLayer.AddCheckbox("Root Checkbox", styleEntry, 0, 2, false, true)

	// Add controls to child1 layer
	child1Button := child1Layer.AddButton("Child1 Button", styleEntry, 0, 0, 10, 1, true)
	child1Label := child1Layer.AddLabel("Child1 Label", styleEntry, 0, 1, 10)
	child1ProgressBar := child1Layer.AddProgressBar("Child1 Progress", styleEntry, 0, 2, 10, 1, false, 50, 100, false)

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
	child1Layer.Delete()
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
	rootLayer.Delete()
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

/*
TestLayerAndControlPropertyStability is a test which verifies that layer and control properties (enabled
status, position, z-order) remain stable across various operations.

Example:

	Expected Inputs:
	    Modifier commands (SetEnabled, movement, z-order) on layers and controls.
	Expected Outputs:
	    The properties are correctly maintained and correctly inherited by new layers after deletion and recreation.
*/
func TestLayerAndControlPropertyStability(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create a layer hierarchy
	rootLayer := AddLayer(0, 0, 20, 20, 1, nil)
	childLayer := AddLayer(0, 0, 15, 15, 1, rootLayer)

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
	rootLayer.Delete()

	newRootLayer := AddLayer(0, 0, 20, 20, 1, nil)
	newChildLayer := AddLayer(0, 0, 15, 15, 1, newRootLayer)

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

/*
TestLayerAndControlMemoryLeaks is a test which verifies that layers and controls are properly cleaned up
from memory after deletion, preventing memory leaks.

Example:

	Expected Inputs:
	    Creation and deletion of 100 root layers each with several controls and child layers.
	Expected Outputs:
	    No layer or control entries remain in the global memory managers after all deletions.
*/
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
		childLayer := AddLayer(0, 0, 15, 15, i, rootLayer)

		// Add controls to child layer
		childLayer.AddButton("Child Button", styleEntry, 0, 0, 10, 1, true)
		childLayer.AddLabel("Child Label", styleEntry, 0, 1, 10)
		childLayer.AddProgressBar("Progress", styleEntry, 0, 2, 10, 1, false, 50, 100, false)

		// Delete layers and verify all controls are deleted
		rootLayer.Delete()

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

/*
TestComplexInterleavedOperations is a test which verifies that multiple layers with unique styles and
controls can be interleaved, moved, and deleted without affecting each other's state.

Example:

	Expected Inputs:
	    A complex series of interleaved layer creation, control assignment, movement, and selective deletion.
	Expected Outputs:
	    Each branch of the layer hierarchy maintains its specific control properties and styles throughout the operations.
*/
func TestComplexInterleavedOperations(t *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(80, 25)
	layer.ReInitializeScreenMemory()

	// Create a layer hierarchy
	root1Layer := AddLayer(0, 0, 20, 20, 1, nil)
	root2Layer := AddLayer(0, 0, 20, 20, 2, nil)
	child1Layer := AddLayer(0, 0, 15, 15, 1, root1Layer)
	child2Layer := AddLayer(0, 0, 15, 15, 2, root1Layer)
	child3Layer := AddLayer(0, 0, 15, 15, 1, root2Layer)
	grandchild1Layer := AddLayer(0, 0, 10, 10, 1, child1Layer)

	// Create style entries with unique properties
	styleEntry1 := types.NewTuiStyleEntry()
	styleEntry1.Button.ForegroundColor = constants.ColorRed
	styleEntry1.Button.BackgroundColor = constants.ColorBlue
	styleEntry1.Button.RaisedColor = constants.ColorGreen

	styleEntry2 := types.NewTuiStyleEntry()
	styleEntry2.Button.ForegroundColor = constants.ColorYellow
	styleEntry2.Button.BackgroundColor = constants.ColorMagenta
	styleEntry2.Button.RaisedColor = constants.ColorCyan

	selectionEntry1 := types.NewSelectionEntry()
	selectionEntry1.SelectionValue = []string{"Option A", "Option B", "Option C"}
	selectionEntry1.SelectionAlias = []string{"optA", "optB", "optC"}

	selectionEntry2 := types.NewSelectionEntry()
	selectionEntry2.SelectionValue = []string{"Choice 1", "Choice 2", "Choice 3"}
	selectionEntry2.SelectionAlias = []string{"ch1", "ch2", "ch3"}

	// Phase 1: Add controls to root1 layer with unique properties
	root1Button := root1Layer.AddButton("Root1 Button", styleEntry1, 0, 0, 10, 1, true)
	root1Checkbox := root1Layer.AddCheckbox("Root1 Checkbox", styleEntry1, 0, 2, true, true)
	root1ProgressBar := root1Layer.AddProgressBar("Root1 Progress", styleEntry1, 0, 4, 10, 1, false, 75, 100, false)
	root1Layer.AddDropdown(styleEntry1, selectionEntry1, 0, 3, 3, 10, 1)
	root1Layer.AddRadioButton("Root1 Radio", styleEntry1, 0, 5, 2, true)
	root1Layer.AddScrollbar(styleEntry1, 0, 6, 10, 100, 25, 1, false)
	root1Layer.AddSelector(styleEntry1, selectionEntry1, 0, 7, 3, 10, 2, 0, 0, true, true)
	root1Layer.AddTextField(styleEntry1, 0, 8, 10, 20, true, "Root1 Text", true)
	root1Layer.AddTextbox(styleEntry1, 0, 9, 10, 3, true)
	root1Layer.AddTooltip("Root1 Tooltip", styleEntry1, 0, 0, 10, 1, 0, 0, 10, 1, true, true, 2000)

	// Phase 2: Delete child1 layer and verify
	child1Layer.Delete()
	if isLayerExists(child1Layer.layerAlias) {
		t.Error("Child1 layer should be deleted")
	}
	if isLayerExists(grandchild1Layer.layerAlias) {
		t.Error("Grandchild1 layer should be deleted when child1 is deleted")
	}

	// Phase 3: Add new controls to remaining layers with unique properties
	child2Button := child2Layer.AddButton("Child2 Button", styleEntry2, 0, 0, 10, 1, true)
	child2Checkbox := child2Layer.AddCheckbox("Child2 Checkbox", styleEntry2, 0, 2, false, true)
	child2ProgressBar := child2Layer.AddProgressBar("Child2 Progress", styleEntry2, 0, 4, 10, 1, false, 25, 100, false)
	child2Layer.AddDropdown(styleEntry2, selectionEntry2, 0, 3, 3, 10, 2)
	child2Layer.AddRadioButton("Child2 Radio", styleEntry2, 0, 5, 1, false)
	child2Layer.AddScrollbar(styleEntry2, 0, 6, 10, 100, 50, 1, false)
	child2Layer.AddSelector(styleEntry2, selectionEntry2, 0, 7, 3, 10, 1, 0, 0, true, true)
	child2Layer.AddTextField(styleEntry2, 0, 8, 10, 20, false, "Child2 Text", true)
	child2Layer.AddTextbox(styleEntry2, 0, 9, 10, 3, true)
	child2Layer.AddTooltip("Child2 Tooltip", styleEntry2, 0, 0, 10, 1, 0, 0, 10, 1, false, true, 1500)

	// Verify child2 controls after creation
	child2ButtonEntry := Buttons.Get(child2Layer.layerAlias, child2Button.controlAlias)
	if child2ButtonEntry.StyleEntry.Button.ForegroundColor != constants.ColorYellow {
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
	newChildLayer := AddLayer(0, 0, 15, 15, 3, root1Layer)
	newChildLayer.AddButton("New Child Button", styleEntry1, 0, 0, 10, 1, true)
	newChildLayer.AddCheckbox("New Child Checkbox", styleEntry1, 0, 2, true, true)
	newChildLayer.AddProgressBar("New Child Progress", styleEntry1, 0, 4, 10, 1, false, 50, 100, false)
	newChildLayer.AddDropdown(styleEntry1, selectionEntry1, 0, 3, 3, 10, 0)
	newChildLayer.AddRadioButton("New Child Radio", styleEntry1, 0, 5, 3, true)
	newChildLayer.AddScrollbar(styleEntry1, 0, 6, 10, 100, 75, 1, false)
	newChildLayer.AddSelector(styleEntry1, selectionEntry1, 0, 7, 3, 10, 1, 0, 0, true, true)
	newChildLayer.AddTextField(styleEntry1, 0, 8, 10, 20, true, "New Child Text", true)
	newChildLayer.AddTextbox(styleEntry1, 0, 9, 10, 3, true)
	newChildLayer.AddTooltip("New Child Tooltip", styleEntry1, 0, 0, 10, 1, 0, 0, 10, 1, true, true, 3000)

	// Phase 5: Delete root2 and verify
	root2Layer.Delete()
	if isLayerExists(root2Layer.layerAlias) {
		t.Error("Root2 layer should be deleted")
	}
	if isLayerExists(child3Layer.layerAlias) {
		t.Error("Child3 layer should be deleted when root2 is deleted")
	}

	// Phase 6: Add new root layer and children with unique properties
	root3Layer := AddLayer(0, 0, 20, 20, 3, nil)
	root3Child1 := AddLayer(0, 0, 15, 15, 1, root3Layer)
	root3Child2 := AddLayer(0, 0, 15, 15, 2, root3Layer)

	// Add controls to root3 layer with unique properties
	root3Button := root3Layer.AddButton("Root3 Button", styleEntry2, 0, 0, 10, 1, true)
	root3Checkbox := root3Layer.AddCheckbox("Root3 Checkbox", styleEntry2, 0, 2, true, true)
	root3ProgressBar := root3Layer.AddProgressBar("Root3 Progress", styleEntry2, 0, 4, 10, 1, false, 90, 100, false)
	root3Layer.AddDropdown(styleEntry2, selectionEntry2, 0, 3, 3, 10, 1)
	root3Layer.AddRadioButton("Root3 Radio", styleEntry2, 0, 5, 2, true)
	root3Layer.AddScrollbar(styleEntry2, 0, 6, 10, 100, 60, 1, false)
	root3Layer.AddSelector(styleEntry2, selectionEntry2, 0, 7, 3, 10, 2, 0, 0, true, true)
	root3Layer.AddTextField(styleEntry2, 0, 8, 10, 20, true, "Root3 Text", true)
	root3Layer.AddTextbox(styleEntry2, 0, 9, 10, 3, true)
	root3Layer.AddTooltip("Root3 Tooltip", styleEntry2, 0, 0, 10, 1, 0, 0, 10, 1, true, true, 2500)

	// Add controls to root3Child1 with unique properties
	root3Child1.AddButton("Root3 Child1 Button", styleEntry1, 0, 0, 10, 1, true)
	root3Child1.AddCheckbox("Root3 Child1 Checkbox", styleEntry1, 0, 2, false, true)
	root3Child1.AddProgressBar("Root3 Child1 Progress", styleEntry1, 0, 4, 10, 1, false, 40, 100, false)
	root3Child1.AddDropdown(styleEntry1, selectionEntry1, 0, 3, 3, 10, 2)
	root3Child1.AddRadioButton("Root3 Child1 Radio", styleEntry1, 0, 5, 1, false)
	root3Child1.AddScrollbar(styleEntry1, 0, 6, 10, 100, 30, 1, false)
	root3Child1.AddSelector(styleEntry1, selectionEntry1, 0, 7, 3, 10, 1, 0, 0, true, true)
	root3Child1.AddTextField(styleEntry1, 0, 8, 10, 20, false, "Root3 Child1 Text", true)
	root3Child1.AddTextbox(styleEntry1, 0, 9, 10, 3, true)
	root3Child1.AddTooltip("Root3 Child1 Tooltip", styleEntry1, 0, 0, 10, 1, 0, 0, 10, 1, false, true, 3500)

	// Phase 7: Move some layers and verify control properties remain intact
	root1Layer.MoveLayerByAbsoluteValue(5, 5)
	newChildLayer.MoveLayerByAbsoluteValue(2, 2)
	root3Layer.MoveLayerByAbsoluteValue(10, 10)

	// Verify control properties after moving layers
	root1ButtonEntry := Buttons.Get(root1Layer.layerAlias, root1Button.controlAlias)
	if root1ButtonEntry.StyleEntry.Button.ForegroundColor != constants.ColorRed {
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
	root1Layer.Delete()
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
	if root3ButtonEntry.StyleEntry.Button.ForegroundColor != constants.ColorYellow {
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
	root3Child2ProgressBar := root3Child2.AddProgressBar("Root3 Child2 Progress", styleEntry2, 0, 4, 10, 1, false, 60, 100, false)
	root3Child2.AddDropdown(styleEntry2, selectionEntry2, 0, 3, 3, 10, 0)
	root3Child2.AddRadioButton("Root3 Child2 Radio", styleEntry2, 0, 5, 1, true)
	root3Child2.AddScrollbar(styleEntry2, 0, 6, 10, 100, 40, 1, false)
	root3Child2.AddSelector(styleEntry2, selectionEntry2, 0, 7, 3, 10, 1, 0, 0, true, true)
	root3Child2.AddTextField(styleEntry2, 0, 8, 10, 20, true, "Root3 Child2 Text", true)
	root3Child2.AddTextbox(styleEntry2, 0, 9, 10, 3, true)
	root3Child2.AddTooltip("Root3 Child2 Tooltip", styleEntry2, 0, 0, 10, 1, 0, 0, 10, 1, true, true, 4000)

	// Verify root3Child2 controls after creation
	root3Child2ButtonEntry := Buttons.Get(root3Child2.layerAlias, root3Child2Button.controlAlias)
	if root3Child2ButtonEntry.StyleEntry.Button.ForegroundColor != constants.ColorYellow {
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
	root3Child1.Delete()
	if isLayerExists(root3Child1.layerAlias) {
		t.Error("Root3 child1 layer should be deleted")
	}
	if !isLayerExists(root3Child2.layerAlias) {
		t.Error("Root3 child2 layer should still exist")
	}

	// Verify root3 and root3Child2 controls still have correct properties
	root3ButtonEntry = Buttons.Get(root3Layer.layerAlias, root3Button.controlAlias)
	if root3ButtonEntry.StyleEntry.Button.ForegroundColor != constants.ColorYellow {
		t.Error("Root3 button should maintain its yellow foreground color after child deletion")
	}
	root3Child2ButtonEntry = Buttons.Get(root3Child2.layerAlias, root3Child2Button.controlAlias)
	if root3Child2ButtonEntry.StyleEntry.Button.ForegroundColor != constants.ColorYellow {
		t.Error("Root3 Child2 button should maintain its yellow foreground color after sibling deletion")
	}

	DeleteAllLayers()
}

/*
TestCalculateWordWidthWithRunes is a test which verifies that calculateWordWidth measures a word in printed
COLUMNS, skipping markup tags, so a word of wide runes reports two columns per rune.

Example:

	Expected Inputs:
	    Rune arrays containing Japanese characters and {{red}} markup tags.
	Expected Outputs:
	    The word "世界" (two wide runes wrapped in {{red}}...{{/}}) reports width 4; the ASCII word "World"
	    reports width 5.
*/
func TestCalculateWordWidthWithRunes(test *testing.T) {
	input := []rune(" こんにちは {{red}}世界{{/}} test")

	// "世界" is two wide runes, so its printed width is four columns.
	width := calculateWordWidth(input, 6, true)
	assert.Equal(test, 4, width, "calculateWordWidth failed to calculate correct column width for '世界'")

	// Test width without markup
	inputNoMarkup := []rune(" Hello World")
	widthNoMarkup := calculateWordWidth(inputNoMarkup, 6, false)
	assert.Equal(test, 5, widthNoMarkup, "calculateWordWidth failed for standard ASCII")
}

/*
TestPutRuneNarrow is a test which verifies that putRune writes a single narrow rune into the character grid,
applies the supplied attributes to that cell, reports one column consumed, and leaves the following cell
untouched.

Example:

	Expected Inputs:
	    A 10x3 grid; putRune at (2,1) with 'A' and an attribute whose CellControlId is 7.
	Expected Outputs:
	    Return value 1; cell (2,1) holds 'A' with CellControlId 7; cell (3,1) still holds the zero rune.
*/
func TestPutRuneNarrow(test *testing.T) {
	layerEntry := types.NewLayerEntry("test", "", 10, 3)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.CellControlId = 7
	attributeEntry.CellType = constants.CellTypeTextbox

	columnsConsumed := putRune(layerEntry.CharacterMemory, 2, 1, 'A', attributeEntry, layerEntry.Width, layerEntry.Height)

	assert.Equal(test, 1, columnsConsumed, "A narrow rune must consume exactly one column.")
	assert.Equal(test, rune('A'), layerEntry.CharacterMemory[1][2].Character, "The lead cell must hold the written rune.")
	assert.Equal(test, 7, layerEntry.CharacterMemory[1][2].AttributeEntry.CellControlId, "The lead cell must carry the supplied attributes.")
	assert.Equal(test, rune(0), layerEntry.CharacterMemory[1][3].Character, "The following cell must be left untouched.")
}

/*
TestPutRuneWide is a test which verifies that putRune writes a wide rune plus a trailing blank placeholder cell,
reports two columns consumed, and gives the placeholder its own copy of the same attribute values as the lead
cell so a hit test on either half resolves alike.

Example:

	Expected Inputs:
	    A 10x3 grid; putRune at (2,1) with '中' and an attribute whose CellControlId is 5, CellType is textbox.
	Expected Outputs:
	    Return value 2; cell (2,1) holds '中'; cell (3,1) holds ' '; both cells share CellControlId 5,
	    the same foreground and background colour, and the same CellType.
*/
func TestPutRuneWide(test *testing.T) {
	layerEntry := types.NewLayerEntry("test", "", 10, 3)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.CellControlId = 5
	attributeEntry.CellType = constants.CellTypeTextbox
	attributeEntry.ForegroundColor = constants.ColorType(1234)
	attributeEntry.BackgroundColor = constants.ColorType(5678)

	columnsConsumed := putRune(layerEntry.CharacterMemory, 2, 1, '中', attributeEntry, layerEntry.Width, layerEntry.Height)

	assert.Equal(test, 2, columnsConsumed, "A wide rune whose placeholder fits must consume two columns.")
	leadCell := layerEntry.CharacterMemory[1][2]
	placeholderCell := layerEntry.CharacterMemory[1][3]
	assert.Equal(test, rune('中'), leadCell.Character, "The lead cell must hold the wide rune.")
	assert.Equal(test, rune(' '), placeholderCell.Character, "The placeholder cell must hold a blank space.")
	assert.Equal(test, leadCell.AttributeEntry.CellControlId, placeholderCell.AttributeEntry.CellControlId, "The placeholder must share the lead cell's CellControlId.")
	assert.Equal(test, leadCell.AttributeEntry.ForegroundColor, placeholderCell.AttributeEntry.ForegroundColor, "The placeholder must share the lead cell's foreground colour.")
	assert.Equal(test, leadCell.AttributeEntry.BackgroundColor, placeholderCell.AttributeEntry.BackgroundColor, "The placeholder must share the lead cell's background colour.")
	assert.Equal(test, leadCell.AttributeEntry.CellType, placeholderCell.AttributeEntry.CellType, "The placeholder must share the lead cell's CellType.")
}

/*
TestPutRuneWideRightEdgeClip is a test which verifies that when a wide rune is written to the final column its
placeholder is dropped rather than written past the right edge, only the lead cell is set, and one column is
reported as consumed.

Example:

	Expected Inputs:
	    A 4x2 grid; putRune at (3,0) with '中'.
	Expected Outputs:
	    Return value 1; cell (3,0) holds '中'; the grid keeps its 4x2 shape and nothing is written past it.
*/
func TestPutRuneWideRightEdgeClip(test *testing.T) {
	layerEntry := types.NewLayerEntry("test", "", 4, 2)
	attributeEntry := types.NewAttributeEntry()

	columnsConsumed := putRune(layerEntry.CharacterMemory, 3, 0, '中', attributeEntry, layerEntry.Width, layerEntry.Height)

	assert.Equal(test, 1, columnsConsumed, "A wide rune clipped by the right edge must report one column consumed.")
	assert.Equal(test, rune('中'), layerEntry.CharacterMemory[0][3].Character, "The lead cell must still hold the wide rune.")
	assert.Len(test, layerEntry.CharacterMemory[0], 4, "The row must keep its original width.")
}

/*
TestPutRuneOutOfBounds is a test which verifies that putRune writes nothing and reports zero columns consumed
when the target position lies outside the grid.

Example:

	Expected Inputs:
	    A 5x5 grid; putRune calls at (5,0), (-1,0), (0,5) and (0,-1).
	Expected Outputs:
	    Every call returns 0 and the grid is left entirely as zero runes.
*/
func TestPutRuneOutOfBounds(test *testing.T) {
	layerEntry := types.NewLayerEntry("test", "", 5, 5)
	attributeEntry := types.NewAttributeEntry()

	for _, position := range [][2]int{{5, 0}, {-1, 0}, {0, 5}, {0, -1}} {
		columnsConsumed := putRune(layerEntry.CharacterMemory, position[0], position[1], 'A', attributeEntry, layerEntry.Width, layerEntry.Height)
		assert.Equalf(test, 0, columnsConsumed, "An out-of-bounds write at (%d,%d) must consume no columns.", position[0], position[1])
	}
	for _, row := range layerEntry.CharacterMemory {
		for _, cell := range row {
			assert.Equal(test, rune(0), cell.Character, "No cell may be modified by an out-of-bounds write.")
		}
	}
}

/*
TestPutRuneOverWidePlaceholderDestroysLead is a test which verifies that when a later write lands on the
placeholder cell of a wide rune, putRune blanks that rune's now-orphaned lead cell to a space so no corrupted
half-wide rune is left in the grid.

Example:

	Expected Inputs:
	    A 10x3 grid; putRune at (2,1) with '中' (consumes columns 2 and 3), then putRune at (3,1) with 'A'.
	Expected Outputs:
	    The second call returns 1; cell (2,1) holds ' '; cell (3,1) holds 'A'.
*/
func TestPutRuneOverWidePlaceholderDestroysLead(test *testing.T) {
	layerEntry := types.NewLayerEntry("test", "", 10, 3)
	attributeEntry := types.NewAttributeEntry()

	putRune(layerEntry.CharacterMemory, 2, 1, '中', attributeEntry, layerEntry.Width, layerEntry.Height)
	columnsConsumed := putRune(layerEntry.CharacterMemory, 3, 1, 'A', attributeEntry, layerEntry.Width, layerEntry.Height)

	assert.Equal(test, 1, columnsConsumed, "The overwriting narrow rune must consume exactly one column.")
	assert.Equal(test, rune(' '), layerEntry.CharacterMemory[1][2].Character, "The orphaned wide-rune lead cell must be blanked to a space.")
	assert.Equal(test, rune('A'), layerEntry.CharacterMemory[1][3].Character, "The target cell must hold the newly written rune.")
}

/*
TestNoLayerOverflowFuzz is a test which drives putRune with randomised runes, positions and grid sizes drawn
from an ASCII pool and a CJK pool and checks two invariants on every iteration: the grid never changes shape
(an out-of-range write would panic first), and a wide rune whose placeholder fits is always followed by a blank
placeholder cell carrying the same CellControlId, colours and CellType as the lead cell.

Example:

	Expected Inputs:
	    2000 iterations, grid width 1..40, grid height 1..5, start x in 0..width+1, runes from
	    "Hello, World! 0123456789 abcXYZ" and "中文字漢字あいうえおカタカナ한국어".
	Expected Outputs:
	    Out-of-range start x returns 0; an in-range narrow rune returns 1; an in-range wide rune with room
	    returns 2 and its placeholder matches the lead cell's id, colours and CellType.
*/
func TestNoLayerOverflowFuzz(test *testing.T) {
	asciiPool := []rune("Hello, World! 0123456789 abcXYZ")
	cjkPool := []rune("中文字漢字あいうえおカタカナ한국어")
	randomSource := rand.New(rand.NewSource(1))

	for iteration := 0; iteration < 2000; iteration++ {
		layerWidth := 1 + randomSource.Intn(40)
		layerHeight := 1 + randomSource.Intn(5)
		layerEntry := types.NewLayerEntry("fuzz", "", layerWidth, layerHeight)

		attributeEntry := types.NewAttributeEntry()
		attributeEntry.CellControlId = randomSource.Intn(1000)
		attributeEntry.CellType = constants.CellTypeTextbox
		attributeEntry.ForegroundColor = constants.ColorType(randomSource.Int63())
		attributeEntry.BackgroundColor = constants.ColorType(randomSource.Int63())

		xLocation := randomSource.Intn(layerWidth + 2)
		yLocation := randomSource.Intn(layerHeight)
		var character rune
		if randomSource.Intn(2) == 0 {
			character = asciiPool[randomSource.Intn(len(asciiPool))]
		} else {
			character = cjkPool[randomSource.Intn(len(cjkPool))]
		}

		columnsConsumed := putRune(layerEntry.CharacterMemory, xLocation, yLocation, character, attributeEntry, layerWidth, layerHeight)

		assert.Len(test, layerEntry.CharacterMemory, layerHeight, "The grid height must never change.")
		for _, row := range layerEntry.CharacterMemory {
			assert.Len(test, row, layerWidth, "The grid width must never change.")
		}

		if xLocation >= layerWidth {
			assert.Equal(test, 0, columnsConsumed, "A start column past the right edge must consume no columns.")
			continue
		}

		leadCell := layerEntry.CharacterMemory[yLocation][xLocation]
		assert.Equal(test, character, leadCell.Character, "The lead cell must hold the written rune.")
		if stringformat.GetWidthOfRuneWhenPrinted(character) == 2 && xLocation+1 < layerWidth {
			assert.Equal(test, 2, columnsConsumed, "A wide rune with room for its placeholder must consume two columns.")
			placeholderCell := layerEntry.CharacterMemory[yLocation][xLocation+1]
			assert.Equal(test, rune(' '), placeholderCell.Character, "The placeholder cell must be a blank space.")
			assert.Equal(test, leadCell.AttributeEntry.CellControlId, placeholderCell.AttributeEntry.CellControlId, "The placeholder must share the lead cell's CellControlId.")
			assert.Equal(test, leadCell.AttributeEntry.ForegroundColor, placeholderCell.AttributeEntry.ForegroundColor, "The placeholder must share the lead cell's foreground colour.")
			assert.Equal(test, leadCell.AttributeEntry.BackgroundColor, placeholderCell.AttributeEntry.BackgroundColor, "The placeholder must share the lead cell's background colour.")
			assert.Equal(test, leadCell.AttributeEntry.CellType, placeholderCell.AttributeEntry.CellType, "The placeholder must share the lead cell's CellType.")
		} else {
			assert.Equal(test, 1, columnsConsumed, "A narrow rune, or a wide rune clipped by the right edge, must consume one column.")
		}
	}
}
