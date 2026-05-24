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

func TestTerminalGetColor(test *testing.T) {
	for currentColorIndex := 0; currentColorIndex < len(constants.AnsiColorByIndex); currentColorIndex++ {
		colorValue := GetColor(currentColorIndex)
		assert.Equalf(test, constants.AnsiColorByIndex[currentColorIndex], colorValue, "The color returned did not match the color at index '%d'.", currentColorIndex)
	}
}

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
