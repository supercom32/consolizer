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

func TestTerminalAddLayer(test *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(20, 20)
	layer1 := AddLayer(0, 0, 20, 20, 1, nil)
	layer2 := AddLayer(0, 0, 20, 20, 2, &layer1)
	layer3 := AddLayer(0, 0, 20, 20, 3, nil)
	layerEntry := Layers.Get(layer1.layerAlias)
	assert.Equalf(test, layer1.layerAlias, layerEntry.LayerAlias, "Failed to get layer entry!")
	layerEntry = Layers.Get(layer2.layerAlias)
	assert.Equalf(test, layer2.layerAlias, layerEntry.LayerAlias, "Failed to get layer entry!")
	layerEntry = Layers.Get(layer3.layerAlias)
	assert.Equalf(test, layer3.layerAlias, layerEntry.LayerAlias, "Failed to get layer entry!")
}

func TestTerminalDefaultLayer(test *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(20, 20)
	layer1 := AddLayer(0, 0, 20, 20, 1, nil)
	assert.Equalf(test, layer1.layerAlias, commonResource.layerInstance.layerAlias, "When creating a new layer, the default layer was not updated correctly!")
	layer2 := AddLayer(0, 0, 20, 20, 2, nil)
	assert.Equalf(test, layer2.layerAlias, commonResource.layerInstance.layerAlias, "When creating a new layer, the default layer was not updated correctly!")
	layer3 := AddLayer(0, 0, 20, 20, 3, nil)
	assert.Equalf(test, layer3.layerAlias, commonResource.layerInstance.layerAlias, "When creating a new layer, the default layer was not updated correctly!")
	Layer(layer1)
	assert.Equalf(test, layer1.layerAlias, commonResource.layerInstance.layerAlias, "When creating a new layer, the default layer was not updated correctly!")
	Layer(layer2)
	assert.Equalf(test, layer2.layerAlias, commonResource.layerInstance.layerAlias, "When creating a new layer, the default layer was not updated correctly!")
}

func TestTerminalSetAlpha(test *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(20, 20)
	layer1 := AddLayer(0, 0, 20, 20, 1, nil)
	SetLayerAlpha(layer1, 50)
	layerEntry := Layers.Get(layer1.layerAlias)
	alphaValue := layerEntry.DefaultAttribute.ForegroundTransformValue
	assert.Equalf(test, float32(50), alphaValue, "Setting the foreground alpha value for a layer failed.")
	alphaValue = layerEntry.DefaultAttribute.BackgroundTransformValue
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
	Color(3, 12)
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
	ColorRGB(foregroundRedIndex, foregroundGreenIndex, foregroundBlueIndex, backgroundRedIndex, backgroundGreenIndex, backgroundBlueIndex)
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
	Locate(xLocation, yLocation)
	layerEntry := Layers.Get(layer1.layerAlias)
	expectedValues := recast.GetArrayOfInterfaces(xLocation, yLocation)
	obtainedValues := recast.GetArrayOfInterfaces(layerEntry.CursorXLocation, layerEntry.CursorYLocation)
	assert.Equalf(test, expectedValues, obtainedValues, "The cursor position did not move to the location specified.")
	xLocation = 15
	yLocation = 15
	Locate(xLocation, yLocation)
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
	Color(10, 7)
	Print("This is a test print on the first line!") // This line will be intentionally scrolled off
	Color(3, 5)
	Print("This is a test print on the second line!") // This line will be intentionally cut at 'print'.
	Locate(7, 7)
	Color(13, 14)
	Print("This is a test print on an arbitrary location!") // This line will be intentionally shifted.
	Color(3, 15)
	Print("This is a test print after printing on an arbitrary location!") // This line will force scroll by 1 line.
	layerEntry := Layers.Get(layer1.layerAlias)
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := "G1szODsyOzEyODsxMjg7MG0bWzQ4OzI7MTI4OzA7MTI4bVRoaXMgaXMgYSB0ZXN0IHByaW50G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTsyNTU7MjU1bSAgICAgICAgICAgICAgICAgICAgG1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTsyNTU7MjU1bSAgICAgICAgICAgICAgICAgICAgG1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTsyNTU7MjU1bSAgICAgICAgICAgICAgICAgICAgG1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTsyNTU7MjU1bSAgICAgICAgICAgICAgICAgICAgG1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTsyNTU7MjU1bSAgICAgICAgICAgICAgICAgICAgG1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTsyNTU7MjU1bSAgICAgICAbWzM4OzI7MjU1OzA7MjU1bRtbNDg7MjswOzI1NTsyNTVtVGhpcyBpcyBhIHRlcxtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjsxMjg7MTI4OzBtG1s0ODsyOzI1NTsyNTU7MjU1bVRoaXMgaXMgYSB0ZXN0IHByaW50G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0K"
	result := assert.Equalf(test, expectedValue, obtainedValue, "The printed screen does not match the master original!")
	if !result {
		dumpScreenComparisons(expectedValue, obtainedValue)
	}
}

func TestTerminalClear(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 14
	layerHeight := 8
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	Color(13, 14)
	layer1.FillLayer("0123456789")
	layerEntry := Layers.Get(layer1.layerAlias)
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := "G1szODsyOzI1NTswOzI1NW0bWzQ4OzI7MDsyNTU7MjU1bTAxMjM0NTY3ODkwMTIzG1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTswOzI1NW0bWzQ4OzI7MDsyNTU7MjU1bTQ1Njc4OTAxMjM0NTY3G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTswOzI1NW0bWzQ4OzI7MDsyNTU7MjU1bTg5MDEyMzQ1Njc4OTAxG1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTswOzI1NW0bWzQ4OzI7MDsyNTU7MjU1bTIzNDU2Nzg5MDEyMzQ1G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTswOzI1NW0bWzQ4OzI7MDsyNTU7MjU1bTY3ODkwMTIzNDU2Nzg5G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTswOzI1NW0bWzQ4OzI7MDsyNTU7MjU1bTAxMjM0NTY3ODkwMTIzG1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTswOzI1NW0bWzQ4OzI7MDsyNTU7MjU1bTQ1Njc4OTAxMjM0NTY3G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzI1NTswOzI1NW0bWzQ4OzI7MDsyNTU7MjU1bTg5MDEyMzQ1Njc4OTAxG1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0K"
	assert.Equalf(test, expectedValue, obtainedValue, "The filled layer does not match the expected result")
	Clear()
	obtainedValue = layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue = "G1szODsyOzI1NTsyNTU7MjU1bRtbNDg7MjswOzA7MG0gICAgICAgICAgICAgIBtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjsyNTU7MjU1OzI1NW0gICAgICAgICAgICAgIBtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjsyNTU7MjU1OzI1NW0gICAgICAgICAgICAgIBtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjsyNTU7MjU1OzI1NW0gICAgICAgICAgICAgIBtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjsyNTU7MjU1OzI1NW0gICAgICAgICAgICAgIBtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjsyNTU7MjU1OzI1NW0gICAgICAgICAgICAgIBtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjsyNTU7MjU1OzI1NW0gICAgICAgICAgICAgIBtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjsyNTU7MjU1OzI1NW0gICAgICAgICAgICAgIBtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtCg=="
	result := assert.Equalf(test, expectedValue, obtainedValue, "The filled layer does not match the expected result")
	if !result {
		dumpScreenComparisons(expectedValue, obtainedValue)
	}
}

func TestTerminalScrollCharacterMemory(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 8
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	Color(10, 7)
	for lineIndex := 0; lineIndex < 13; lineIndex++ {
		Print(fmt.Sprintf("This is the '%d' line of text printed!", lineIndex))
	}
	layerEntry := Layers.Get(layer1.layerAlias)
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := "G1szODsyOzA7MjU1OzBtG1s0ODsyOzE5MjsxOTI7MTkybVRoaXMgaXMgdGhlICc1JyBsaW5lIG9mIHRleHQgcHJpbnRlZCEbWzM4OzI7MjU1OzI1NTsyNTVtG1s0ODsyOzA7MDswbSAgIBtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzI1NTswbRtbNDg7MjsxOTI7MTkyOzE5Mm1UaGlzIGlzIHRoZSAnNicgbGluZSBvZiB0ZXh0IHByaW50ZWQhG1szODsyOzI1NTsyNTU7MjU1bRtbNDg7MjswOzA7MG0gICAbWzM4OzI7MDswOzBtG1s0ODsyOzA7MDswbQobWzM4OzI7MDsyNTU7MG0bWzQ4OzI7MTkyOzE5MjsxOTJtVGhpcyBpcyB0aGUgJzcnIGxpbmUgb2YgdGV4dCBwcmludGVkIRtbMzg7MjsyNTU7MjU1OzI1NW0bWzQ4OzI7MDswOzBtICAgG1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzA7MjU1OzBtG1s0ODsyOzE5MjsxOTI7MTkybVRoaXMgaXMgdGhlICc4JyBsaW5lIG9mIHRleHQgcHJpbnRlZCEgICAbWzM4OzI7MDswOzBtG1s0ODsyOzA7MDswbQobWzM4OzI7MDsyNTU7MG0bWzQ4OzI7MTkyOzE5MjsxOTJtVGhpcyBpcyB0aGUgJzknIGxpbmUgb2YgdGV4dCBwcmludGVkISAgIBtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzI1NTswbRtbNDg7MjsxOTI7MTkyOzE5Mm1UaGlzIGlzIHRoZSAnMTAnIGxpbmUgb2YgdGV4dCBwcmludGVkISAgG1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzA7MjU1OzBtG1s0ODsyOzE5MjsxOTI7MTkybVRoaXMgaXMgdGhlICcxMScgbGluZSBvZiB0ZXh0IHByaW50ZWQhICAbWzM4OzI7MDswOzBtG1s0ODsyOzA7MDswbQobWzM4OzI7MDsyNTU7MG0bWzQ4OzI7MTkyOzE5MjsxOTJtVGhpcyBpcyB0aGUgJzEyJyBsaW5lIG9mIHRleHQgcHJpbnRlZCEgIBtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtCg=="
	assert.Equalf(test, expectedValue, obtainedValue, "The printed screen does not match the master original!")
}

func TestTerminalGetRuneOnLayer(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 8
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	Color(10, 7)
	layerEntry := Layers.Get(layer1.layerAlias)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.CellUserId = 999
	arrayOfRunes := stringformat.GetRunesFromString("T")
	printLayer(layerEntry, attributeEntry, 3, 7, arrayOfRunes)
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
	Layer(layer1)
	Color(4, 6)
	layer1.FillLayer("a1a2a3a4a5")
	Layer(layer2)
	Color(3, 11)
	layer2.FillLayer("b1b2b3b4b5")
	Layer(layer3)
	Color(12, 13)
	layer3.FillLayer("c1c2c3c4c5")
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := "G1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEyYTNhNGE1YTFhMmEzYTRhNWExYTJhM2E0YTVhMWEyYTNhNGE1G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEyYTNhNGE1YTFhMmEzYTRhNWExYTJhM2E0YTVhMWEyYTNhNGE1G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEbWzM4OzI7MTI4OzEyODswbRtbNDg7MjsyNTU7MjU1OzBtYjFiMmIzYjRiNWIxYjJiM2I0YjViMWIyYjNiNGI1YjFiMmIzYhtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhG1szODsyOzEyODsxMjg7MG0bWzQ4OzI7MjU1OzI1NTswbWIxYjJiM2I0YjViMWIyYjNiNGI1YjFiMmIzYjRiNWIxYjJiM2IbWzM4OzI7MDswOzBtG1s0ODsyOzA7MDswbQobWzM4OzI7MDswOzEyOG0bWzQ4OzI7MDsxMjg7MTI4bWExYRtbMzg7MjsxMjg7MTI4OzBtG1s0ODsyOzI1NTsyNTU7MG1iMWIbWzM4OzI7MDswOzI1NW0bWzQ4OzI7MjU1OzA7MjU1bWMxYzJjM2M0YzVjMWMyYzNjNGM1YzFjMmMzYzRjNWMxYzIbWzM4OzI7MDswOzBtG1s0ODsyOzA7MDswbQobWzM4OzI7MDswOzEyOG0bWzQ4OzI7MDsxMjg7MTI4bWExYRtbMzg7MjsxMjg7MTI4OzBtG1s0ODsyOzI1NTsyNTU7MG1iMWIbWzM4OzI7MDswOzI1NW0bWzQ4OzI7MjU1OzA7MjU1bWMxYzJjM2M0YzVjMWMyYzNjNGM1YzFjMmMzYzRjNWMxYzIbWzM4OzI7MDswOzBtG1s0ODsyOzA7MDswbQobWzM4OzI7MDswOzEyOG0bWzQ4OzI7MDsxMjg7MTI4bWExYRtbMzg7MjsxMjg7MTI4OzBtG1s0ODsyOzI1NTsyNTU7MG1iMWIbWzM4OzI7MDswOzI1NW0bWzQ4OzI7MjU1OzA7MjU1bWMxYzJjM2M0YzVjMWMyYzNjNGM1YzFjMmMzYzRjNWMxYzIbWzM4OzI7MDswOzBtG1s0ODsyOzA7MDswbQobWzM4OzI7MDswOzEyOG0bWzQ4OzI7MDsxMjg7MTI4bWExYRtbMzg7MjsxMjg7MTI4OzBtG1s0ODsyOzI1NTsyNTU7MG1iMWIbWzM4OzI7MDswOzI1NW0bWzQ4OzI7MjU1OzA7MjU1bWMxYzJjM2M0YzVjMWMyYzNjNGM1YzFjMmMzYzRjNWMxYzIbWzM4OzI7MDswOzBtG1s0ODsyOzA7MDswbQo="
	assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!")
}

func TestTerminalRenderParentLayer(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 80
	layerHeight := 20
	InitializeTerminal(layerWidth, layerHeight)
	// First set of nested text layers.
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(3, 2, 15, 15, 2, &layer1)
	layer3 := AddLayer(6, 4, 20, 20, 3, &layer2)
	layer4 := AddLayer(-6, -4, 10, 10, 3, &layer3)
	// Second set of nested text layers, on top of first root parent.
	layer5 := AddLayer(15, 0, 40, 20, 2, nil)
	layer6 := AddLayer(3, 2, 15, 15, 3, &layer5)
	layer7 := AddLayer(6, 4, 20, 20, 4, &layer6)
	layer8 := AddLayer(0, -4, 10, 10, 5, &layer7)
	layer9 := AddLayer(20, 3, 10, 10, 3, &layer5)
	layer10 := AddLayer(50, 3, 10, 10, 3, &layer1)

	Layer(layer1)
	Color(4, 6)
	layer1.FillLayer("a1a2a3a4a5")
	Layer(layer2)
	Color(3, 11)
	layer2.FillLayer("b1b2b3b4b5")
	Layer(layer3)
	Color(12, 13)
	layer3.FillLayer("c1c2c3c4c5")
	Layer(layer4)
	Color(1, 2)
	layer4.FillLayer("c1c2c3c4c5")
	Layer(layer5)
	Color(6, 7)
	layer5.FillLayer("a1a2a3a4a5")
	Layer(layer6)
	Color(4, 12)
	layer6.FillLayer("b1b2b3b4b5")
	Layer(layer7)
	Color(13, 14)
	layer7.FillLayer("c1c2c3c4c5")
	Layer(layer8)
	Color(2, 3)
	layer8.FillLayer("c1c2c3c4c5")
	Layer(layer9)
	Color(7, 4)
	layer9.FillLayer("c1c2c3c4c5")
	Layer(layer10)
	Color(9, 12)
	layer10.FillLayer("c1c2c3c4c5")

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := "G1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEyYTNhNGE1YTFhMmEbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hMWEyYTNhNGE1YTFhMmEzYTRhNWExYTJhM2E0YTVhMWEyYTNhNGE1G1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG0zYTRhNWExYTJhM2E0YTVhMWEyYTNhNGE1G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEyYTNhNGE1YTFhMmEbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hMWEyYTNhNGE1YTFhMmEzYTRhNWExYTJhM2E0YTVhMWEyYTNhNGE1G1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG0zYTRhNWExYTJhM2E0YTVhMWEyYTNhNGE1G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEbWzM4OzI7MTI4OzEyODswbRtbNDg7MjsyNTU7MjU1OzBtYjFiMmIzYjRiNWIxG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MDsyNTVtYjFiMmIzYjRiNWIxYjJiG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTVhMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtM2E0YTVhMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhG1szODsyOzEyODsxMjg7MG0bWzQ4OzI7MjU1OzI1NTswbTNiNGI1YjFiMmIzYhtbMzg7MjswOzEyODsxMjhtG1s0ODsyOzE5MjsxOTI7MTkybWExYRtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzA7MjU1bTNiNGI1YjFiMmIzYjRiNRtbMzg7MjswOzEyODsxMjhtG1s0ODsyOzE5MjsxOTI7MTkybWE1G1szODsyOzE5MjsxOTI7MTkybRtbNDg7MjswOzA7MTI4bWMxYzJjM2M0YzUbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hMWEyYTNhNGE1G1szODsyOzI1NTswOzBtG1s0ODsyOzA7MDsyNTVtM2M0YzUbWzM4OzI7MDswOzEyOG0bWzQ4OzI7MDsxMjg7MTI4bWExYTJhM2E0YTVhMWEyYTNhNGE1G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEbWzM4OzI7MTI4OzEyODswbRtbNDg7MjsyNTU7MjU1OzBtYjFiMmIzYjRiNWIxG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MDsyNTVtYjFiMmIzYjRiNWIxYjJiG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTUbWzM4OzI7MTkyOzE5MjsxOTJtG1s0ODsyOzA7MDsxMjhtYzFjMmMzYzRjNRtbMzg7MjswOzEyODsxMjhtG1s0ODsyOzE5MjsxOTI7MTkybWExYTJhM2E0YTUbWzM4OzI7MjU1OzA7MG0bWzQ4OzI7MDswOzI1NW0zYzRjNRtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhMmEzYTRhNWExYTJhM2E0YTUbWzM4OzI7MDswOzBtG1s0ODsyOzA7MDswbQobWzM4OzI7MDswOzEyOG0bWzQ4OzI7MDsxMjg7MTI4bWExYRtbMzg7MjsxMjg7MTI4OzBtG1s0ODsyOzI1NTsyNTU7MG0zYjRiNWIxYjJiM2IbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hMWEbWzM4OzI7MDswOzEyOG0bWzQ4OzI7MDswOzI1NW0zYjRiNWIxYjJiM2I0YjUbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hNRtbMzg7MjsxOTI7MTkyOzE5Mm0bWzQ4OzI7MDswOzEyOG1jMWMyYzNjNGM1G1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhMmEzYTRhNRtbMzg7MjsyNTU7MDswbRtbNDg7MjswOzA7MjU1bTNjNGM1G1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhG1szODsyOzEyODsxMjg7MG0bWzQ4OzI7MjU1OzI1NTswbWIxYjJiMxtbMzg7MjsxMjg7MDswbRtbNDg7MjswOzEyODswbWM0YzUbWzM4OzI7MDswOzI1NW0bWzQ4OzI7MjU1OzA7MjU1bWMzG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MDsyNTVtYjFiMmIzG1szODsyOzA7MTI4OzBtG1s0ODsyOzEyODsxMjg7MG1jMWMyYzNjNGMbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hNRtbMzg7MjsxOTI7MTkyOzE5Mm0bWzQ4OzI7MDswOzEyOG1jMWMyYzNjNGM1G1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhMmEzYTRhNRtbMzg7MjsyNTU7MDswbRtbNDg7MjswOzA7MjU1bTNjNGM1G1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhG1szODsyOzEyODsxMjg7MG0bWzQ4OzI7MjU1OzI1NTswbTNiNGI1YhtbMzg7MjsxMjg7MDswbRtbNDg7MjswOzEyODswbWM0YzUbWzM4OzI7MDswOzI1NW0bWzQ4OzI7MjU1OzA7MjU1bWMzG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MDsyNTVtM2I0YjViG1szODsyOzA7MTI4OzBtG1s0ODsyOzEyODsxMjg7MG1jMWMyYzNjNGMbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hNRtbMzg7MjsxOTI7MTkyOzE5Mm0bWzQ4OzI7MDswOzEyOG1jMWMyYzNjNGM1G1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhMmEzYTRhNRtbMzg7MjsyNTU7MDswbRtbNDg7MjswOzA7MjU1bTNjNGM1G1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhG1szODsyOzEyODsxMjg7MG0bWzQ4OzI7MjU1OzI1NTswbWIxYjJiMxtbMzg7MjsxMjg7MDswbRtbNDg7MjswOzEyODswbWM0YzUbWzM4OzI7MDswOzI1NW0bWzQ4OzI7MjU1OzA7MjU1bWMzG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MDsyNTVtYjFiMmIzG1szODsyOzA7MTI4OzBtG1s0ODsyOzEyODsxMjg7MG1jMWMyYzNjNGMbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hNRtbMzg7MjsxOTI7MTkyOzE5Mm0bWzQ4OzI7MDswOzEyOG1jMWMyYzNjNGM1G1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhMmEzYTRhNRtbMzg7MjsyNTU7MDswbRtbNDg7MjswOzA7MjU1bTNjNGM1G1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhG1szODsyOzEyODsxMjg7MG0bWzQ4OzI7MjU1OzI1NTswbTNiNGI1YhtbMzg7MjsxMjg7MDswbRtbNDg7MjswOzEyODswbWM0YzUbWzM4OzI7MDswOzI1NW0bWzQ4OzI7MjU1OzA7MjU1bWMzG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MDsyNTVtM2I0YjViG1szODsyOzA7MTI4OzBtG1s0ODsyOzEyODsxMjg7MG1jMWMyYzNjNGMbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hNRtbMzg7MjsxOTI7MTkyOzE5Mm0bWzQ4OzI7MDswOzEyOG1jMWMyYzNjNGM1G1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhMmEzYTRhNRtbMzg7MjsyNTU7MDswbRtbNDg7MjswOzA7MjU1bTNjNGM1G1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhG1szODsyOzEyODsxMjg7MG0bWzQ4OzI7MjU1OzI1NTswbWIxYjJiMxtbMzg7MjsxMjg7MDswbRtbNDg7MjswOzEyODswbWM0YzUbWzM4OzI7MDswOzI1NW0bWzQ4OzI7MjU1OzA7MjU1bWMzG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MDsyNTVtYjFiMmIzG1szODsyOzA7MTI4OzBtG1s0ODsyOzEyODsxMjg7MG1jMWMyYzNjNGMbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hNRtbMzg7MjsxOTI7MTkyOzE5Mm0bWzQ4OzI7MDswOzEyOG1jMWMyYzNjNGM1G1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhMmEzYTRhNRtbMzg7MjsyNTU7MDswbRtbNDg7MjswOzA7MjU1bTNjNGM1G1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhG1szODsyOzEyODsxMjg7MG0bWzQ4OzI7MjU1OzI1NTswbTNiNGI1YhtbMzg7MjsxMjg7MDswbRtbNDg7MjswOzEyODswbWM0YzUbWzM4OzI7MDswOzI1NW0bWzQ4OzI7MjU1OzA7MjU1bWMzG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MDsyNTVtM2I0YjViG1szODsyOzA7MTI4OzBtG1s0ODsyOzEyODsxMjg7MG1jMWMyYzNjNGMbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hNRtbMzg7MjsxOTI7MTkyOzE5Mm0bWzQ4OzI7MDswOzEyOG1jMWMyYzNjNGM1G1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhMmEzYTRhNRtbMzg7MjsyNTU7MDswbRtbNDg7MjswOzA7MjU1bTNjNGM1G1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhG1szODsyOzEyODsxMjg7MG0bWzQ4OzI7MjU1OzI1NTswbWIxYjJiMxtbMzg7MjswOzA7MjU1bRtbNDg7MjsyNTU7MDsyNTVtYzFjMmMzG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MDsyNTVtYjFiMmIzG1szODsyOzI1NTswOzI1NW0bWzQ4OzI7MDsyNTU7MjU1bWMxYzJjM2M0YxtbMzg7MjswOzEyODsxMjhtG1s0ODsyOzE5MjsxOTI7MTkybWE1G1szODsyOzE5MjsxOTI7MTkybRtbNDg7MjswOzA7MTI4bWMxYzJjM2M0YzUbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hMWEyYTNhNGE1G1szODsyOzI1NTswOzBtG1s0ODsyOzA7MDsyNTVtM2M0YzUbWzM4OzI7MDswOzEyOG0bWzQ4OzI7MDsxMjg7MTI4bWExYTJhM2E0YTVhMWEyYTNhNGE1G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEbWzM4OzI7MTI4OzEyODswbRtbNDg7MjsyNTU7MjU1OzBtM2I0YjViG1szODsyOzA7MDsyNTVtG1s0ODsyOzI1NTswOzI1NW1jMWMyYzMbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hMWEbWzM4OzI7MDswOzEyOG0bWzQ4OzI7MDswOzI1NW0zYjRiNWIbWzM4OzI7MjU1OzA7MjU1bRtbNDg7MjswOzI1NTsyNTVtYzFjMmMzYzRjG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTVhMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtM2E0YTVhMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhG1szODsyOzEyODsxMjg7MG0bWzQ4OzI7MjU1OzI1NTswbWIxYjJiMxtbMzg7MjswOzA7MjU1bRtbNDg7MjsyNTU7MDsyNTVtYzFjMmMzG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MDsyNTVtYjFiMmIzG1szODsyOzI1NTswOzI1NW0bWzQ4OzI7MDsyNTU7MjU1bWMxYzJjM2M0YxtbMzg7MjswOzEyODsxMjhtG1s0ODsyOzE5MjsxOTI7MTkybWE1YTFhMmEzYTRhNWExYTJhM2E0YTUbWzM4OzI7MDswOzEyOG0bWzQ4OzI7MDsxMjg7MTI4bTNhNGE1YTFhMmEzYTRhNWExYTJhM2E0YTUbWzM4OzI7MDswOzBtG1s0ODsyOzA7MDswbQobWzM4OzI7MDswOzEyOG0bWzQ4OzI7MDsxMjg7MTI4bWExYRtbMzg7MjsxMjg7MTI4OzBtG1s0ODsyOzI1NTsyNTU7MG0zYjRiNWIbWzM4OzI7MDswOzI1NW0bWzQ4OzI7MjU1OzA7MjU1bWMxYzJjMxtbMzg7MjswOzEyODsxMjhtG1s0ODsyOzE5MjsxOTI7MTkybWExYRtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzA7MjU1bTNiNGI1YhtbMzg7MjsyNTU7MDsyNTVtG1s0ODsyOzA7MjU1OzI1NW1jMWMyYzNjNGMbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hNWExYTJhM2E0YTVhMWEyYTNhNGE1G1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG0zYTRhNWExYTJhM2E0YTVhMWEyYTNhNGE1G1szODsyOzA7MDswbRtbNDg7MjswOzA7MG0KG1szODsyOzA7MDsxMjhtG1s0ODsyOzA7MTI4OzEyOG1hMWEbWzM4OzI7MTI4OzEyODswbRtbNDg7MjsyNTU7MjU1OzBtYjFiMmIzG1szODsyOzA7MDsyNTVtG1s0ODsyOzI1NTswOzI1NW1jMWMyYzMbWzM4OzI7MDsxMjg7MTI4bRtbNDg7MjsxOTI7MTkyOzE5Mm1hMWEbWzM4OzI7MDswOzEyOG0bWzQ4OzI7MDswOzI1NW1iMWIyYjMbWzM4OzI7MjU1OzA7MjU1bRtbNDg7MjswOzI1NTsyNTVtYzFjMmMzYzRjG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTVhMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtM2E0YTVhMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhMmEzYTRhNWExYTJhG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhMmEzYTRhNWExYTJhM2E0YTVhMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtM2E0YTVhMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhMmEzYTRhNWExYTJhG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhMmEzYTRhNWExYTJhM2E0YTVhMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtM2E0YTVhMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtChtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtYTFhMmEzYTRhNWExYTJhG1szODsyOzA7MTI4OzEyOG0bWzQ4OzI7MTkyOzE5MjsxOTJtYTFhMmEzYTRhNWExYTJhM2E0YTVhMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MTI4bRtbNDg7MjswOzEyODsxMjhtM2E0YTVhMWEyYTNhNGE1YTFhMmEzYTRhNRtbMzg7MjswOzA7MG0bWzQ4OzI7MDswOzBtCg=="
	assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!")
}

func TestDeleteLayer(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 20
	layerHeight := 20
	InitializeTerminal(layerWidth, layerHeight)
	p1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	p2 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	p3 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	p4 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	p1c1 := AddLayer(0, 0, layerWidth, layerHeight, 1, &p1)
	p1c2 := AddLayer(0, 0, layerWidth, layerHeight, 1, &p1c1)
	p1c3 := AddLayer(0, 0, layerWidth, layerHeight, 1, &p1c2)
	AddLayer(0, 0, layerWidth, layerHeight, 1, &p1c3)
	p3c1 := AddLayer(0, 0, layerWidth, layerHeight, 1, &p3)
	p3c2 := AddLayer(0, 0, layerWidth, layerHeight, 1, &p3c1)
	p3c3 := AddLayer(0, 0, layerWidth, layerHeight, 1, &p3c2)
	AddLayer(0, 0, layerWidth, layerHeight, 1, &p3c3)
	sortedLayerAliasSlice := layer.GetSortedLayerMemoryAliasSlice()
	obtainedValue := len(sortedLayerAliasSlice)
	expectedValue := 12
	assert.Equalf(test, expectedValue, obtainedValue, "The number of layers created does not match!")

	Layer(p1)
	p3.DeleteLayer()
	sortedLayerAliasSlice = layer.GetSortedLayerMemoryAliasSlice()
	obtainedValue = len(sortedLayerAliasSlice)
	expectedValue = 7
	assert.Equalf(test, expectedValue, obtainedValue, "The number of layers created does not match!")

	Layer(p1)
	p4.DeleteLayer()
	sortedLayerAliasSlice = layer.GetSortedLayerMemoryAliasSlice()
	obtainedValue = len(sortedLayerAliasSlice)
	expectedValue = 6
	assert.Equalf(test, expectedValue, obtainedValue, "The number of layers created does not match!")

	Layer(p2)
	p1.DeleteLayer()
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
