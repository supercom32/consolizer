package consolizer

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/recast"
	"github.com/supercom32/consolizer/types"
	"testing"
)

const TUI_ELEMENTS_TEST_SUITE_NAME = "tui_elements"

/*
TestGetDarkenedCharacterEntry is a test which verifies the darkening of a character entry's colors.

Example:
    Expected Inputs:
        A character entry with full white foreground and gray background, darkened by 50%.
    Expected Outputs:
        The resulting character entry has colors with approximately 50% of the original RGB values.
*/
func TestGetDarkenedCharacterEntry(test *testing.T) {
	characterEntry := types.NewCharacterEntry()
	characterEntry.AttributeEntry.ForegroundColor = constants.ColorType(tcell.NewRGBColor(int32(255), int32(255), int32(255)))
	darkenedCharacterEntry := GetDarkenedCharacterEntry(&characterEntry, 0.5)
	redIndex, greenIndex, blueIndex := GetRGBColorComponents(darkenedCharacterEntry.AttributeEntry.ForegroundColor)
	obtainedResult := recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	expectedResult := recast.GetArrayOfInterfaces(int32(127), int32(127), int32(127))
	assert.Equalf(test, obtainedResult, expectedResult, "Darkening Frame entry foreground color to 50 percent failed!")

	characterEntry.AttributeEntry.BackgroundColor = constants.ColorType(tcell.NewRGBColor(int32(127), int32(127), int32(127)))
	darkenedCharacterEntry = GetDarkenedCharacterEntry(&characterEntry, 0.5)
	redIndex, greenIndex, blueIndex = GetRGBColorComponents(darkenedCharacterEntry.AttributeEntry.BackgroundColor)
	obtainedResult = recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	expectedResult = recast.GetArrayOfInterfaces(int32(63), int32(63), int32(63))
	assert.Equalf(test, obtainedResult, expectedResult, "Darkening Frame entry background color to 50 percent failed!")
}

/*
TestColorDarkening is a test which verifies the darkening of individual colors.

Example:
    Expected Inputs:
        White and gray colors subjected to varying darkening percentages (0%, 50%, 100%).
    Expected Outputs:
        RGB components match the expected mathematical reduction for each percentage.
*/
func TestColorDarkening(test *testing.T) {
	testColor := GetDarkenedColor(constants.ColorType(tcell.NewRGBColor(255, 255, 255)), 0.5)
	redIndex, greenIndex, blueIndex := GetRGBColorComponents(testColor)
	obtainedResult := recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	expectedResult := recast.GetArrayOfInterfaces(int32(127), int32(127), int32(127))
	assert.Equalf(test, obtainedResult, expectedResult, "Darkening color to 50 percent failed!")

	testColor = GetDarkenedColor(constants.ColorType(tcell.NewRGBColor(127, 127, 127)), 0.5)
	redIndex, greenIndex, blueIndex = GetRGBColorComponents(testColor)
	obtainedResult = recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	expectedResult = recast.GetArrayOfInterfaces(int32(63), int32(63), int32(63))
	assert.Equalf(test, obtainedResult, expectedResult, "Darkening color to 50 percent failed!")

	testColor = GetDarkenedColor(constants.ColorType(tcell.NewRGBColor(127, 127, 127)), 0.0)
	redIndex, greenIndex, blueIndex = GetRGBColorComponents(testColor)
	obtainedResult = recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	expectedResult = recast.GetArrayOfInterfaces(int32(0), int32(0), int32(0))
	assert.Equalf(test, obtainedResult, expectedResult, "Darkening color to 0 percent failed!")

	testColor = GetDarkenedColor(constants.ColorType(tcell.NewRGBColor(127, 127, 127)), 1.0)
	redIndex, greenIndex, blueIndex = GetRGBColorComponents(testColor)
	obtainedResult = recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	expectedResult = recast.GetArrayOfInterfaces(int32(127), int32(127), int32(127))
	assert.Equalf(test, obtainedResult, expectedResult, "Darkening color to 100 percent failed!")
}

/*
TestColorTransitions is a test which verifies color transitions between source and target colors.

Example:
    Expected Inputs:
        Specific source and target colors with various percentage transition steps.
    Expected Outputs:
        The resulting transitional colors match expected intermediate RGB values.
*/
func TestColorTransitions(test *testing.T) {
	sourceColor := constants.ColorType(tcell.NewRGBColor(255, 0, 0))
	targetColor := constants.ColorType(tcell.NewRGBColor(0, 255, 255))
	newColor := GetTransitionedColor(sourceColor, targetColor, 0.3)
	redIndex, greenIndex, blueIndex := GetRGBColorComponents(newColor)
	obtainedResult := recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	expectedResult := recast.GetArrayOfInterfaces(int32(178), int32(77), int32(77))
	assert.Equalf(test, expectedResult, obtainedResult, "Transitioning color by 30% failed!")

	sourceColor = constants.ColorType(tcell.NewRGBColor(255, 255, 255))
	targetColor = constants.ColorType(tcell.NewRGBColor(0, 0, 0))
	newColor = GetTransitionedColor(sourceColor, targetColor, 0.5)
	redIndex, greenIndex, blueIndex = GetRGBColorComponents(newColor)
	obtainedResult = recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	expectedResult = recast.GetArrayOfInterfaces(int32(127), int32(127), int32(127))
	assert.Equalf(test, expectedResult, obtainedResult, "Transitioning color by 50% failed!")

	sourceColor = constants.ColorType(tcell.NewRGBColor(255, 255, 255))
	targetColor = constants.ColorType(tcell.NewRGBColor(0, 200, 0))
	newColor = GetTransitionedColor(sourceColor, targetColor, 0.5)
	redIndex, greenIndex, blueIndex = GetRGBColorComponents(newColor)
	obtainedResult = recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	expectedResult = recast.GetArrayOfInterfaces(int32(127), int32(227), int32(127))
	assert.Equalf(test, expectedResult, obtainedResult, "Transitioning color by 50% failed!")

	sourceColor = constants.ColorType(tcell.NewRGBColor(0, 50, 100))
	targetColor = constants.ColorType(tcell.NewRGBColor(150, 200, 255))
	newColor = GetTransitionedColor(sourceColor, targetColor, 0.5)
	redIndex, greenIndex, blueIndex = GetRGBColorComponents(newColor)
	obtainedResult = recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	expectedResult = recast.GetArrayOfInterfaces(int32(75), int32(125), int32(178))
	assert.Equalf(test, expectedResult, obtainedResult, "Transitioning color by 50% failed!")

	sourceColor = constants.ColorType(tcell.NewRGBColor(0, 50, 100))
	targetColor = constants.ColorType(tcell.NewRGBColor(150, 200, 255))
	newColor = GetTransitionedColor(sourceColor, targetColor, 1.2)
	redIndex, greenIndex, blueIndex = GetRGBColorComponents(newColor)
	obtainedResult = recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	expectedResult = recast.GetArrayOfInterfaces(int32(180), int32(230), int32(255))
	assert.Equalf(test, expectedResult, obtainedResult, "Transitioning color by 120% failed!")

	sourceColor = constants.ColorType(tcell.NewRGBColor(0, 50, 100))
	targetColor = constants.ColorType(tcell.NewRGBColor(150, 200, 255))
	newColor = GetTransitionedColor(sourceColor, targetColor, -0.2)
	redIndex, greenIndex, blueIndex = GetRGBColorComponents(newColor)
	obtainedResult = recast.GetArrayOfInterfaces(redIndex, greenIndex, blueIndex)
	expectedResult = recast.GetArrayOfInterfaces(int32(0), int32(20), int32(69))
	assert.Equalf(test, expectedResult, obtainedResult, "Transitioning color by -20% failed!")
}

/*
TestDrawButton is a test which verifies the correct rendering of buttons on a layer.

Example:
    Expected Inputs:
        A layer hierarchy with multiple buttons, one of which is programmatically set to a pressed state.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing both normal and pressed button visuals.
*/
func TestDrawButton(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(3, 2, 30, 20, 2, layer1)
	layer1.Color(4, 6)
	layer1.FillLayer("a1a2a3a4a5")
	layer2.Color(3, 11)
	layer2.FillLayer("b1b2b3b4b5")
	layer2.AddButton("ButtonText", styleEntry, 2, 2, 20, 7, true)
	button2 := layer2.AddButton("ButtonText", styleEntry, 2, 10, 20, 7, true)
	Buttons.Get(layer2.layerAlias, button2.controlAlias).IsPressed = true
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawButton", obtainedValue)
	expectedValue := LoadMasterImage(TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawButton")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDrawFrameLabel is a test which verifies the correct rendering of frame labels on a layer.

Example:
    Expected Inputs:
        Multiple frame labels drawn at various coordinates including off-screen and edge positions.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing correctly formatted and clipped labels.
*/
func TestDrawFrameLabel(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer1.Color(4, 6)
	layer1.FillLayer("a1a2a3a4a5")
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(62, 128, 128)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 200)
	styleEntry.Window.LineDrawingTextLabelForegroundColor = GetRGBColor(255, 255, 255)
	styleEntry.Window.LineDrawingTextLabelBackgroundColor = GetRGBColor(0, 0, 200)

	layer1.DrawFrameLabel(styleEntry, "My Frame", 1, 1)
	layer1.DrawFrameLabel(styleEntry, "My Frame", -7, 3)
	layer1.DrawFrameLabel(styleEntry, "My Frame", layerWidth-7, 3)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawFrameLabel", obtainedValue)
	expectedValue := LoadMasterImage(TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawFrameLabel")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDrawFrame is a test which verifies the correct rendering of frames on a layer.

Example:
    Expected Inputs:
        Raised and sunken frames drawn with specific line-drawing colors.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) with correct frame highlights and shadows.
*/
func TestDrawFrame(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	styleEntry.Window.LineDrawingRaisedColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingSunkenColor = GetRGBColor(0, 0, 0)

	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer1.Color(4, 6)
	layer1.FillLayer("a1a2a3a4a5")
	layer1.DrawFrame(styleEntry, false, 2, 2, 10, 10, false)
	layer1.DrawFrame(styleEntry, true, 15, 2, 10, 10, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawFrame", obtainedValue)
	expectedValue := LoadMasterImage(TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawFrame")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDrawBorder is a test which verifies the correct rendering of borders on a layer.

Example:
    Expected Inputs:
        A simple flat border drawn with specific foreground and background colors.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the single-line border.
*/
func TestDrawBorder(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer1.Color(4, 6)
	layer1.FillLayer("a1a2a3a4a5")
	layer1.DrawBorder(styleEntry, 2, 2, 10, 10, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawBorder", obtainedValue)
	expectedValue := LoadMasterImage(TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawBorder")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDrawWindow is a test which verifies the correct rendering of windows on multiple layers.

Example:
    Expected Inputs:
        A multi-layered setup with a window drawn on the topmost layer.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) with correct transparency and overlap.
*/
func TestDrawWindow(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(2, 2, layerWidth, layerHeight, 2, nil)
	layer3 := AddLayer(4, 4, layerWidth, layerHeight, 3, nil)
	layer4 := AddLayer(0, 0, layerWidth, layerHeight, 4, nil)
	layer1.Color(4, 13)
	layer1.FillLayer("a1a2a3a4a5")
	layer2.Color(8, 10)
	layer2.FillLayer("a1a2a3a4a5")
	layer3.Color(2, 11)
	layer3.FillLayer("a1a2a3a4a5")
	layer4.Color(14, 15)
	layer4.DrawWindow(styleEntry, 0, 0, 10, 10, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawWindow", obtainedValue)
	expectedValue := LoadMasterImage(TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawWindow")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDrawShadow is a test which verifies the correct rendering of shadows on a layer.

Example:
    Expected Inputs:
        A multi-layered setup with a transparent shadow area drawn on the topmost layer.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the darkening effect on underlying layers.
*/
func TestDrawShadow(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(2, 2, layerWidth, layerHeight, 2, nil)
	layer3 := AddLayer(4, 4, layerWidth, layerHeight, 3, nil)
	layer4 := AddLayer(0, 0, layerWidth, layerHeight, 4, nil)
	layer1.Color(4, 13)
	layer1.FillLayer("a1a2a3a4a5")
	layer2.Color(8, 10)
	layer2.FillLayer("a1a2a3a4a5")
	layer3.Color(2, 11)
	layer3.FillLayer("a1a2a3a4a5")
	layer4.Color(14, 15)
	layer4.DrawShadow(0, 0, 10, 10, 0.5)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawShadow", obtainedValue)
	expectedValue := LoadMasterImage(TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawShadow")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDrawLines is a test which verifies the correct rendering of horizontal and vertical lines on a layer.

Example:
    Expected Inputs:
        A variety of horizontal and vertical lines, some with connectors enabled and some without.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing all lines and intersections.
*/
func TestDrawLines(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer1.Color(4, 13)
	layer1.FillLayer("a1a2a3a4a5")
	layer1.DrawBorder(styleEntry, 10, 4, 10, 10, false)
	layer1.DrawHorizontalLine(styleEntry, 7, 6, 15, true)
	layer1.DrawHorizontalLine(styleEntry, 7, 8, 15, false)
	layer1.DrawHorizontalLine(styleEntry, 10, 10, 10, true)
	layer1.DrawVerticalLine(styleEntry, 12, 2, 15, false)
	layer1.DrawVerticalLine(styleEntry, 14, 2, 15, true)
	layer1.DrawVerticalLine(styleEntry, 16, 4, 10, true)
	layer1.DrawVerticalLine(styleEntry, 5, -5, 10, true)
	layer1.DrawVerticalLine(styleEntry, 5, 15, 10, true)
	layer1.DrawHorizontalLine(styleEntry, -5, 5, 10, true)
	layer1.DrawHorizontalLine(styleEntry, 35, 5, 10, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawLines", obtainedValue)
	expectedValue := LoadMasterImage(TUI_ELEMENTS_TEST_SUITE_NAME, "TestDrawLines")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}
