package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/recast"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
	_ "math/rand"
	_ "strconv"
	"testing"
)

const TERMINAL_TEST_SUITE_NAME = "terminal"

/*
TestTerminalAddLayer is a test which verifies that layers can be correctly added to the terminal system,
including nested parent-child relationships.

Example:
    Expected Inputs:
        Multiple layers added with and without parent aliases.
    Expected Outputs:
        All layers are correctly retrieved and their aliases match the input parameters.
*/
func TestTerminalAddLayer(test *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(20, 20)
	layer1 := AddLayer(0, 0, 20, 20, 1, nil)
	layer2 := AddLayer(0, 0, 20, 20, 2, layer1)
	layer3 := AddLayer(0, 0, 20, 20, 3, nil)
	layerEntry := Layers.Get(layer1.layerAlias)
	assert.Equalf(test, layer1.layerAlias, layerEntry.LayerAlias, "Failed to get layer entry!")
	layerEntry = Layers.Get(layer2.layerAlias)
	assert.Equalf(test, layer2.layerAlias, layerEntry.LayerAlias, "Failed to get layer entry!")
	layerEntry = Layers.Get(layer3.layerAlias)
	assert.Equalf(test, layer3.layerAlias, layerEntry.LayerAlias, "Failed to get layer entry!")
}

/*
TestTerminalSetAlpha is a test which verifies that the alpha transparency value for a layer can be correctly
set and retrieved from its default attributes.

Example:
    Expected Inputs:
        A layer instance where SetAlpha is called with a value of 50.0.
    Expected Outputs:
        Both foreground and background alpha values in the layer's default attributes are exactly 50.
*/
func TestTerminalSetAlpha(test *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(20, 20)
	layer1 := AddLayer(0, 0, 20, 20, 1, nil)
	layer1.SetAlpha(50.0)
	layerEntry := Layers.Get(layer1.layerAlias)
	alphaValue := layerEntry.DefaultAttribute.ForegroundAlphaValue
	assert.Equalf(test, float32(50), alphaValue, "Setting the foreground alpha value for a layer failed.")
	alphaValue = layerEntry.DefaultAttribute.BackgroundAlphaValue
	assert.Equalf(test, float32(50), alphaValue, "Setting the background alpha value for a layer failed.")
}

/*
TestTerminalGetColor is a test which verifies that GetColor correctly retrieves ANSI color values from the
predefined color palette.

Example:
    Expected Inputs:
        Iterating through all valid ANSI color indices (0-15).
    Expected Outputs:
        The color values returned match those stored in the constants.AnsiColorByIndex palette.
*/
func TestTerminalGetColor(test *testing.T) {
	for currentColorIndex := 0; currentColorIndex < len(constants.AnsiColorByIndex); currentColorIndex++ {
		colorValue := GetColor(currentColorIndex)
		assert.Equalf(test, constants.AnsiColorByIndex[currentColorIndex], colorValue, "The color returned did not match the color at index '%d'.", currentColorIndex)
	}
}

/*
TestTerminalGetRGBColor is a test which verifies that GetRGBColor correctly constructs 24-bit color values
and that GetRGBColorComponents can successfully deconstruct them.

Example:
    Expected Inputs:
        Specific RGB component triplets (e.g., 0,0,0 and 20,50,75).
    Expected Outputs:
        Correct ColorType values are created and deconstructed back into their original red, green, and blue components.
*/
func TestTerminalGetRGBColor(test *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(20, 20)
	inputRedIndex := int32(0)
	inputGreenIndex := int32(0)
	inputBlueIndex := int32(0)
	color := GetRGBColor(inputRedIndex, inputGreenIndex, inputBlueIndex)
	assert.Equalf(test, constants.ColorType(0x300000000), color, "The color returned for '%d, %d, %d' was not correct. ", inputRedIndex, inputGreenIndex, inputBlueIndex)

	redIndex, greenIndex, blueIndex := GetRGBColorComponents(color)
	expectedValues := recast.GetArrayOfInterfaces(inputRedIndex, inputGreenIndex, inputBlueIndex)
	obtainedValues := recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	assert.Equalf(test, expectedValues, obtainedValues, "The color components returned for '16777216' was not correct. ")

	inputRedIndex = int32(20)
	inputGreenIndex = int32(50)
	inputBlueIndex = int32(75)
	color = GetRGBColor(inputRedIndex, inputGreenIndex, inputBlueIndex)
	assert.Equalf(test, constants.ColorType(0x30014324b), color, "The color returned for '%d, %d, %d' was not correct. ", inputRedIndex, inputGreenIndex, inputBlueIndex)

	redIndex, greenIndex, blueIndex = GetRGBColorComponents(color)
	expectedValues = recast.GetArrayOfInterfaces(inputRedIndex, inputGreenIndex, inputBlueIndex)
	obtainedValues = recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	assert.Equalf(test, expectedValues, obtainedValues, "The color components returned for '18100811' was not correct. ")
}

/*
TestTerminalColor is a test which verifies that setting a layer's default colors using ANSI indices correctly
updates the layer's default attributes.

Example:
    Expected Inputs:
        A layer where Color(3, 12) is called.
    Expected Outputs:
        The layer's default foreground and background colors match ANSI indices 3 and 12 respectively.
*/
func TestTerminalColor(test *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(20, 20)
	layer1 := AddLayer(0, 0, 20, 20, 1, nil)
	layer1.Color(3, 12)
	layerEntry := Layers.Get(layer1.layerAlias)
	expectedValues := recast.GetArrayOfInterfaces(constants.AnsiColorByIndex[3], constants.AnsiColorByIndex[12])
	obtainedValues := recast.GetArrayOfInterfaces(layerEntry.DefaultAttribute.ForegroundColor, layerEntry.DefaultAttribute.BackgroundColor)
	assert.Equalf(test, expectedValues, obtainedValues, "The default specified layer color does not match what was set.")
}

/*
TestTerminalColorRGB is a test which verifies that setting a layer's default colors using RGB component
indices correctly updates the layer's default attributes.

Example:
    Expected Inputs:
        A layer where ColorRGB is called with specific 24-bit component values.
    Expected Outputs:
        The layer's default foreground and background colors match the constructed 24-bit ColorType values.
*/
func TestTerminalColorRGB(test *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(20, 20)
	foregroundRedIndex := int32(75)
	foregroundGreenIndex := int32(101)
	foregroundBlueIndex := int32(249)
	backgroundRedIndex := int32(123)
	backgroundGreenIndex := int32(145)
	backgroundBlueIndex := int32(192)
	layer1 := AddLayer(0, 0, 20, 20, 1, nil)
	layer1.ColorRGB(foregroundRedIndex, foregroundGreenIndex, foregroundBlueIndex, backgroundRedIndex, backgroundGreenIndex, backgroundBlueIndex)
	layerEntry := Layers.Get(layer1.layerAlias)
	expectedValues := recast.GetArrayOfInterfaces(constants.ColorType(0x3004b65f9), constants.ColorType(0x3007b91c0))
	obtainedValues := recast.GetArrayOfInterfaces(layerEntry.DefaultAttribute.ForegroundColor, layerEntry.DefaultAttribute.BackgroundColor)
	assert.Equalf(test, expectedValues, obtainedValues, "The default specified layer color does not match what was set.")
}

/*
TestTerminalMoveLayerByAbsoluteValue is a test which verifies that a layer can be moved to an absolute screen
position, including negative coordinates.

Example:
    Expected Inputs:
        A layer moved to (9, 8) and then to (-10, -13).
    Expected Outputs:
        The layer's screen coordinates match the specified absolute positions after each move.
*/
func TestTerminalMoveLayerByAbsoluteValue(test *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(20, 20)
	layer1 := AddLayer(0, 0, 20, 20, 1, nil)
	xLocation := 9
	yLocation := 8
	layer1.MoveLayerByAbsoluteValue(xLocation, yLocation)
	layerEntry := Layers.Get(layer1.layerAlias)
	expectedValues := recast.GetArrayOfInterfaces(xLocation, yLocation)
	obtainedValues := recast.GetArrayOfInterfaces(layerEntry.ScreenXLocation, layerEntry.ScreenYLocation)
	assert.Equalf(test, expectedValues, obtainedValues, "The layer did not move by the absolute value specified.")
	xLocation = -10
	yLocation = -13
	layer1.MoveLayerByAbsoluteValue(xLocation, yLocation)
	expectedValues = recast.GetArrayOfInterfaces(xLocation, yLocation)
	obtainedValues = recast.GetArrayOfInterfaces(layerEntry.ScreenXLocation, layerEntry.ScreenYLocation)
	assert.Equalf(test, expectedValues, obtainedValues, "The layer did not move by the absolute value specified.")
}

/*
TestTerminalMoveLayerByRelativeValue is a test which verifies that a layer can be moved relative to its current
position using specified offsets.

Example:
    Expected Inputs:
        A layer at (0, 0) moved by (9, -8) and then further by (10, -13).
    Expected Outputs:
        The final screen coordinates (19, -21) correctly reflect the sum of all relative movements.
*/
func TestTerminalMoveLayerByRelativeValue(test *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(20, 20)
	layer1 := AddLayer(0, 0, 20, 20, 1, nil)
	xLocation := 9
	yLocation := -8
	layer1.MoveLayerByRelativeValue(xLocation, yLocation)
	layerEntry := Layers.Get(layer1.layerAlias)
	expectedValues := recast.GetArrayOfInterfaces(xLocation, yLocation)
	obtainedValues := recast.GetArrayOfInterfaces(layerEntry.ScreenXLocation, layerEntry.ScreenYLocation)
	assert.Equalf(test, expectedValues, obtainedValues, "The layer did not move by the relative value specified.")
	xLocation = +10
	yLocation = -13
	layer1.MoveLayerByRelativeValue(xLocation, yLocation)
	expectedValues = recast.GetArrayOfInterfaces(9+xLocation, (-8)+yLocation)
	obtainedValues = recast.GetArrayOfInterfaces(layerEntry.ScreenXLocation, layerEntry.ScreenYLocation)
	assert.Equalf(test, expectedValues, obtainedValues, "The layer did not move by the relative value specified.")
}

/*
TestTerminalLocate is a test which verifies that the layer's internal cursor position can be correctly set
to absolute coordinates.

Example:
    Expected Inputs:
        A layer instance where Locate(9, 10) and then Locate(15, 15) are called.
    Expected Outputs:
        The layer's CursorXLocation and CursorYLocation match the specified coordinates after each call.
*/
func TestTerminalLocate(test *testing.T) {
	commonResource.isDebugEnabled = true
	xLocation := 9
	yLocation := 10
	InitializeTerminal(20, 20)
	layer1 := AddLayer(0, 0, 20, 20, 1, nil)
	layer1.Locate(xLocation, yLocation)
	layerEntry := Layers.Get(layer1.layerAlias)
	expectedValues := recast.GetArrayOfInterfaces(xLocation, yLocation)
	obtainedValues := recast.GetArrayOfInterfaces(layerEntry.CursorXLocation, layerEntry.CursorYLocation)
	assert.Equalf(test, expectedValues, obtainedValues, "The cursor position did not move to the location specified.")
	xLocation = 15
	yLocation = 15
	layer1.Locate(xLocation, yLocation)
	expectedValues = recast.GetArrayOfInterfaces(xLocation, yLocation)
	obtainedValues = recast.GetArrayOfInterfaces(layerEntry.CursorXLocation, layerEntry.CursorYLocation)
	assert.Equalf(test, expectedValues, obtainedValues, "The cursor position did not move to the location specified.")
}

/*
TestTerminalPrint is a test which verifies that text printing correctly updates the character memory and
handles scrolling when text exceeds the layer height.

Example:
    Expected Inputs:
        A sequence of Print commands with different colors and coordinates.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) reflecting the final state of the layer memory.
*/
func TestTerminalPrint(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 20
	layerHeight := 8
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer1.Color(10, 7)
	layer1.Print("This is a test print on the first line!") // This line will be intentionally scrolled off
	layer1.Color(3, 5)
	layer1.Print("This is a test print on the second line!") // This line will be intentionally cut at 'print'.
	layer1.Locate(7, 7)
	layer1.Color(13, 14)
	layer1.Print("This is a test print on an arbitrary location!") // This line will be intentionally shifted.
	layer1.Color(3, 15)
	layer1.Print("This is a test print after printing on an arbitrary location!") // This line will force scroll by 1 line.
	layerEntry := Layers.Get(layer1.layerAlias)
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TERMINAL_TEST_SUITE_NAME, "TestTerminalPrint", obtainedValue)
	expectedValue := LoadMasterImage(TERMINAL_TEST_SUITE_NAME, "TestTerminalPrint")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The printed screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTerminalClear is a test which verifies that the Clear method successfully resets all characters and
attributes in a layer back to their default state.

Example:
    Expected Inputs:
        A layer filled with text and then cleared.
    Expected Outputs:
        Screen content matches the master ANSI string for a filled layer (before) and a blank layer (after).
*/
func TestTerminalClear(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 14
	layerHeight := 8
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer1.Color(13, 14)
	layer1.FillLayer("0123456789")
	layerEntry := Layers.Get(layer1.layerAlias)
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TERMINAL_TEST_SUITE_NAME, "TestTerminalClear_Before", obtainedValue)
	expectedValue := LoadMasterImage(TERMINAL_TEST_SUITE_NAME, "TestTerminalClear_Before")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The filled layer does not match the expected result") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
	layer1.Clear()
	obtainedValue = layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TERMINAL_TEST_SUITE_NAME, "TestTerminalClear_After", obtainedValue)
	expectedValue = LoadMasterImage(TERMINAL_TEST_SUITE_NAME, "TestTerminalClear_After")

	obtainedValueBase64 = layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 = layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The filled layer does not match the expected result") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTerminalScrollCharacterMemory is a test which verifies that character memory correctly scrolls when the
number of lines printed exceeds the layer height.

Example:
    Expected Inputs:
        13 lines of text printed to a layer with a height of 8.
    Expected Outputs:
        The final screen content matches the master ANSI string showing only the last few lines printed.
*/
func TestTerminalScrollCharacterMemory(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 8
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer1.Color(10, 7)
	for lineIndex := 0; lineIndex < 13; lineIndex++ {
		layer1.Print(fmt.Sprintf("This is the '%d' line of text printed!", lineIndex))
	}
	layerEntry := Layers.Get(layer1.layerAlias)
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TERMINAL_TEST_SUITE_NAME, "TestTerminalScrollCharacterMemory", obtainedValue)
	expectedValue := LoadMasterImage(TERMINAL_TEST_SUITE_NAME, "TestTerminalScrollCharacterMemory")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The printed screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTerminalGetRuneOnLayer is a test which verifies that individual cell IDs and character information can
be accurately retrieved from specific layer locations.

Example:
    Expected Inputs:
        Retrieve cell ID from (3, 7) on a layer where a character with ID 999 was printed.
    Expected Outputs:
        Correct cell ID (999) and NullCellId for adjacent unprinted locations.
*/
func TestTerminalGetRuneOnLayer(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 8
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer1.Color(10, 7)
	layerEntry := Layers.Get(layer1.layerAlias)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.CellUserId = 999
	arrayOfRunes := stringformat.GetRunesFromString("T")
	layer.printLayer(layerEntry, attributeEntry, 3, 7, arrayOfRunes)
	obtainedValue := getCellIdByLayerAlias(layer1.layerAlias, 3, 7)
	expectedValue := 999
	assert.Equalf(test, expectedValue, obtainedValue, "The expected cell ID was not found at the specified location!")
	obtainedValue = getCellIdByLayerAlias(layer1.layerAlias, 2, 7)
	expectedValue = constants.NullCellId
	assert.Equalf(test, expectedValue, obtainedValue, "The expected cell ID was not found at the specified location!")
	obtainedValue = getCellIdByLayerAlias(layer1.layerAlias, 4, 7)
	expectedValue = constants.NullCellId
	assert.Equalf(test, expectedValue, obtainedValue, "The expected cell ID was not found at the specified location!")
}

/*
TestTerminalUpdateDisplay is a test which verifies that multiple layers with different priorities and fills
are correctly composited onto the final screen buffer.

Example:
    Expected Inputs:
        Three layers with unique background colors and overlapping positions.
    Expected Outputs:
        Screen content matches the master ANSI string showing correctly layered and overlapping content.
*/
func TestTerminalUpdateDisplay(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 8
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(3, 2, layerWidth, layerHeight, 2, nil)
	layer3 := AddLayer(6, 4, layerWidth, layerHeight, 3, nil)
	layer1.Color(4, 6)
	layer1.FillLayer("a1a2a3a4a5")
	layer2.Color(3, 11)
	layer2.FillLayer("b1b2b3b4b5")
	layer3.Color(12, 13)
	layer3.FillLayer("c1c2c3c4c5")

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TERMINAL_TEST_SUITE_NAME, "TestTerminalUpdateDisplay", obtainedValue)
	expectedValue := LoadMasterImage(TERMINAL_TEST_SUITE_NAME, "TestTerminalUpdateDisplay")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTerminalRenderParentLayer is a test which verifies the rendering of deep layer hierarchies involving multiple
levels of nesting and offsets.

Example:
    Expected Inputs:
        Two distinct sets of nested layers (up to 5 levels deep) with unique fill patterns.
    Expected Outputs:
        Screen content matches the master ANSI string with all parent-child relationships correctly rendered.
*/
func TestTerminalRenderParentLayer(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 80
	layerHeight := 20
	InitializeTerminal(layerWidth, layerHeight)
	// First set of nested text layers.
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(3, 2, 15, 15, 2, layer1)
	layer3 := AddLayer(6, 4, 20, 20, 3, layer2)
	layer4 := AddLayer(-6, -4, 10, 10, 3, layer3)
	// Second set of nested text layers, on top of first root parent.
	layer5 := AddLayer(15, 0, 40, 20, 2, nil)
	layer6 := AddLayer(3, 2, 15, 15, 3, layer5)
	layer7 := AddLayer(6, 4, 20, 20, 4, layer6)
	layer8 := AddLayer(0, -4, 10, 10, 5, layer7)
	layer9 := AddLayer(20, 3, 10, 10, 3, layer5)
	layer10 := AddLayer(50, 3, 10, 10, 3, layer1)

	layer1.Color(4, 6)
	layer1.FillLayer("a1a2a3a4a5")
	layer2.Color(3, 11)
	layer2.FillLayer("b1b2b3b4b5")
	layer3.Color(12, 13)
	layer3.FillLayer("c1c2c3c4c5")
	layer4.Color(1, 2)
	layer4.FillLayer("c1c2c3c4c5")
	layer5.Color(6, 7)
	layer5.FillLayer("a1a2a3a4a5")
	layer6.Color(4, 12)
	layer6.FillLayer("b1b2b3b4b5")
	layer7.Color(13, 14)
	layer7.FillLayer("c1c2c3c4c5")
	layer8.Color(2, 3)
	layer8.FillLayer("c1c2c3c4c5")
	layer9.Color(7, 4)
	layer9.FillLayer("c1c2c3c4c5")
	layer10.Color(9, 12)
	layer10.FillLayer("c1c2c3c4c5")

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TERMINAL_TEST_SUITE_NAME, "TestTerminalRenderParentLayer", obtainedValue)
	expectedValue := LoadMasterImage(TERMINAL_TEST_SUITE_NAME, "TestTerminalRenderParentLayer")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDeleteLayer is a test which verifies that deleting a parent layer correctly removes all its descendants
from the layer system.

Example:
    Expected Inputs:
        A hierarchy of 12 layers where three primary branches are selectively deleted.
    Expected Outputs:
        The total number of active layers in the system decreases correctly after each branch deletion.
*/
func TestDeleteLayer(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 20
	layerHeight := 20
	InitializeTerminal(layerWidth, layerHeight)
	p1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	p3 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	p4 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	p1c1 := AddLayer(0, 0, layerWidth, layerHeight, 1, p1)
	p1c2 := AddLayer(0, 0, layerWidth, layerHeight, 1, p1c1)
	p1c3 := AddLayer(0, 0, layerWidth, layerHeight, 1, p1c2)
	AddLayer(0, 0, layerWidth, layerHeight, 1, p1c3)
	p3c1 := AddLayer(0, 0, layerWidth, layerHeight, 1, p3)
	p3c2 := AddLayer(0, 0, layerWidth, layerHeight, 1, p3c1)
	p3c3 := AddLayer(0, 0, layerWidth, layerHeight, 1, p3c2)
	AddLayer(0, 0, layerWidth, layerHeight, 1, p3c3)
	sortedLayerAliasSlice := layer.GetSortedLayerMemoryAliasSlice()
	obtainedValue := len(sortedLayerAliasSlice)
	expectedValue := 12
	assert.Equalf(test, expectedValue, obtainedValue, "The number of layers created does not match!")

	p3.Delete()
	sortedLayerAliasSlice = layer.GetSortedLayerMemoryAliasSlice()
	obtainedValue = len(sortedLayerAliasSlice)
	expectedValue = 7
	assert.Equalf(test, expectedValue, obtainedValue, "The number of layers created does not match!")

	p4.Delete()
	sortedLayerAliasSlice = layer.GetSortedLayerMemoryAliasSlice()
	obtainedValue = len(sortedLayerAliasSlice)
	expectedValue = 6
	assert.Equalf(test, expectedValue, obtainedValue, "The number of layers created does not match!")

	p1.Delete()
	sortedLayerAliasSlice = layer.GetSortedLayerMemoryAliasSlice()
	obtainedValue = len(sortedLayerAliasSlice)
	expectedValue = 1
	assert.Equalf(test, expectedValue, obtainedValue, "The number of layers created does not match!")
}

/*
TestNewAssetList is a test which verifies that asset lists correctly store and clear image asset metadata.

Example:
    Expected Inputs:
        An asset list where one image metadata entry is added and then cleared.
    Expected Outputs:
        The metadata entry correctly reflects the input parameters, and the list becomes empty after clearing.
*/
func TestNewAssetList(test *testing.T) {
	imageFileList := NewAssetList()
	imageStyle := NewImageStyle()
	imageFileList.AddPreloadedImage("fileName1", imageStyle, 10, 11, 0.6)
	obtainedValue := recast.GetArrayOfInterfaces(imageFileList.PreloadedImageList[0].FileName, imageFileList.PreloadedImageList[0].FileAlias, imageFileList.PreloadedImageList[0].WidthInCharacters, imageFileList.PreloadedImageList[0].HeightInCharacters, imageFileList.PreloadedImageList[0].BlurSigma)
	expectedValue := recast.GetArrayOfInterfaces("fileName1", "fileName1", 10, 11, 0.6)
	assert.Equalf(test, expectedValue, obtainedValue, "The file entry obtained does not match what was set!")
	imageFileList.Clear()
	obtainedValue = recast.GetArrayOfInterfaces(len(imageFileList.PreloadedImageList))
	expectedValue = recast.GetArrayOfInterfaces(0)
	assert.Equalf(test, expectedValue, obtainedValue, "The number of file entries does not what was expected!")
}
