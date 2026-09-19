package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

const BUTTON_TEST_SUITE_NAME = "button"

/*
TestButtonDefaultState is a test which verifies that a button is rendered correctly in its default state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button is rendered at (2,2) with label "Test" and width 10.
*/
func TestButtonDefaultState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonDefaultState", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonDefaultState") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonDisabledState is a test which verifies that a button is rendered correctly when it is in a disabled state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button is rendered with a disabled label color.
*/
func TestButtonDisabledState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonDisabledState", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonDisabledState") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonPressedState is a test which verifies that a button is rendered correctly when it is in a pressed state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button is rendered with a sunken frame style.
*/
func TestButtonPressedState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	buttonObject := layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)

	// Simulate a button press by directly setting the button state
	buttonEntry := Buttons.Get(layer1.layerAlias, buttonObject.controlAlias)
	buttonEntry.IsPressed = true

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonPressedState", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonPressedState") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonSelectedState is a test which verifies that a button is rendered correctly when it is in a selected state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button label is underlined.
*/
func TestButtonSelectedState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	buttonObject := layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)

	// Simulate a button selection by directly setting the button state
	buttonEntry := Buttons.Get(layer1.layerAlias, buttonObject.controlAlias)
	buttonEntry.IsSelected = true

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonSelectedState", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonSelectedState") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonCustomDimensions is a test which verifies that a button is rendered correctly with custom width and height.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button is rendered with width 15 and height 5.
*/
func TestButtonCustomDimensions(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddButton("Test", styleEntry, 2, 2, 15, 5, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonCustomDimensions", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonCustomDimensions")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonLongLabel is a test which verifies that a button is rendered correctly when it has a label longer than its width.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button width expands to fit the long label.
*/
func TestButtonLongLabel(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddButton("This is a long button label", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonLongLabel", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonLongLabel") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonDelete is a test which verifies that a button is successfully removed when its Delete method is called.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button is absent from the rendered output after deletion.
*/
func TestButtonDelete(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	buttonInstance := layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)

	// Delete the button
	buttonInstance.Delete()

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonDelete", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonDelete") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonDeleteAll is a test which verifies that all buttons are successfully removed from a layer.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        All buttons are absent from the rendered output after calling DeleteAllButtons.
*/
func TestButtonDeleteAll(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddButton("Test 1", styleEntry, 2, 2, 10, 3, true)
	layer1.AddButton("Test 2", styleEntry, 2, 6, 10, 3, true)
	UpdateDisplay(false)

	// Delete all buttons
	Button.DeleteAll(layer1.layerAlias)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonDeleteAll", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonDeleteAll") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonFocus is a test which verifies that a button correctly handles focus state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button focus state is reflected in the system.
*/
func TestButtonFocus(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)

	// Set focus to the button
	setFocusedControl(layer1.layerAlias, "testButton", constants.CellTypeButton)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonFocus", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonFocus") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonFlatStyle is a test which verifies that a flat button renders its frame without the two-tone bevel.

Example:
    Expected Inputs:
        styleEntry.Button.StyleMode set to constants.ButtonStyleFlat; a button "Test" at (2,2) width 10.

    Expected Outputs:
        Screen content matches the committed master image; the frame is a single uniform colour.
*/
func TestButtonFlatStyle(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	styleEntry.Button.StyleMode = constants.ButtonStyleFlat
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonFlatStyle", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonFlatStyle")
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", layerEntry.GetAnsiStringFromBase64(expectedValue))
		fmt.Println("Obtained:\n", layerEntry.GetAnsiStringFromBase64(obtainedValue))
	}
}

/*
TestButtonFlatFrameHasNoBevel is a test which verifies that a flat button's opposite frame corners share one
foreground colour, so there is no raised or sunken bevel, and that the colour is the button foreground.

Example:
    Expected Inputs:
        A flat button at (2,2) width 10 height 3 with a distinct grey Button.ForegroundColor.

    Expected Outputs:
        The top-left corner cell (2,2) and the bottom-right corner cell (11,4) both carry that grey foreground
        colour.
*/
func TestButtonFlatFrameHasNoBevel(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	frameColor := GetRGBColor(210, 210, 210)
	styleEntry.Button.StyleMode = constants.ButtonStyleFlat
	styleEntry.Button.ForegroundColor = frameColor
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)
	characterMemory := commonResource.screenLayer.CharacterMemory

	topLeftForeground := characterMemory[2][2].AttributeEntry.ForegroundColor
	bottomRightForeground := characterMemory[4][11].AttributeEntry.ForegroundColor
	assert.Equal(test, topLeftForeground, bottomRightForeground, "a flat frame has no bevel, so opposite corners match")
	assert.Equal(test, frameColor, topLeftForeground)
}

/*
TestButtonFlatPressedUsesPressedColors is a test which verifies that a pressed flat button, which has no bevel to
flip, shows the press through Button.PressedBackgroundColor on its cells.

Example:
    Expected Inputs:
        A flat button at (2,2) width 10 with Button.PressedBackgroundColor a distinct green; its entry
        IsPressed is set to true before drawing.

    Expected Outputs:
        A fill cell inside the button (3,3) carries the green pressed background colour.
*/
func TestButtonFlatPressedUsesPressedColors(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	pressedBackgroundColor := GetRGBColor(0, 150, 0)
	styleEntry.Button.StyleMode = constants.ButtonStyleFlat
	styleEntry.Button.PressedBackgroundColor = pressedBackgroundColor
	buttonInstance := layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	buttonEntry := Buttons.Get(buttonInstance.layerAlias, buttonInstance.controlAlias)
	buttonEntry.IsPressed = true
	UpdateDisplay(false)
	characterMemory := commonResource.screenLayer.CharacterMemory

	assert.Equal(test, pressedBackgroundColor, characterMemory[3][3].AttributeEntry.BackgroundColor)
}

/*
TestButtonBorderlessStyle is a test which verifies that a borderless button renders its label with no frame at
all.

Example:
    Expected Inputs:
        styleEntry.Button.StyleMode set to constants.ButtonStyleBorderless; a button "Test" at (2,2) width 10.

    Expected Outputs:
        Screen content matches the committed master image.
*/
func TestButtonBorderlessStyle(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	styleEntry.Button.StyleMode = constants.ButtonStyleBorderless
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonBorderlessStyle", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonBorderlessStyle")
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", layerEntry.GetAnsiStringFromBase64(expectedValue))
		fmt.Println("Obtained:\n", layerEntry.GetAnsiStringFromBase64(obtainedValue))
	}
}

/*
TestButtonBorderlessHasNoFrame is a test which verifies that a borderless button draws no box-drawing characters
anywhere in its rectangle and still centres its label.

Example:
    Expected Inputs:
        A borderless button "Test" at (2,2) width 10 height 3.

    Expected Outputs:
        No cell in rows 2..4, columns 2..11 holds a box-drawing rune; the centred label "Test" starts at
        column 5 of row 3.
*/
func TestButtonBorderlessHasNoFrame(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	styleEntry.Button.StyleMode = constants.ButtonStyleBorderless
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)
	characterMemory := commonResource.screenLayer.CharacterMemory

	boxDrawingRunes := "┌┐└┘─│├┤┬┴┼╔╗╚╝═║"
	for rowIndex := 2; rowIndex <= 4; rowIndex++ {
		for columnIndex := 2; columnIndex <= 11; columnIndex++ {
			currentRune := characterMemory[rowIndex][columnIndex].Character
			assert.NotContainsf(test, boxDrawingRunes, string(currentRune), "borderless button drew a frame rune at (%d,%d)", columnIndex, rowIndex)
		}
	}
	assert.Equal(test, []rune("Test"), []rune{characterMemory[3][5].Character, characterMemory[3][6].Character, characterMemory[3][7].Character, characterMemory[3][8].Character})
}

/*
TestButtonBorderlessPressedUsesPressedColors is a test which verifies that a pressed borderless button shows the
press through Button.PressedBackgroundColor.

Example:
    Expected Inputs:
        A borderless button "Test" at (2,2) width 10 with a distinct green Button.PressedBackgroundColor; its
        entry IsPressed is set to true before drawing.

    Expected Outputs:
        A fill cell inside the button (3,3) carries the green pressed background colour.
*/
func TestButtonBorderlessPressedUsesPressedColors(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	pressedBackgroundColor := GetRGBColor(0, 150, 0)
	styleEntry.Button.StyleMode = constants.ButtonStyleBorderless
	styleEntry.Button.PressedBackgroundColor = pressedBackgroundColor
	buttonInstance := layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	buttonEntry := Buttons.Get(buttonInstance.layerAlias, buttonInstance.controlAlias)
	buttonEntry.IsPressed = true
	UpdateDisplay(false)
	characterMemory := commonResource.screenLayer.CharacterMemory

	assert.Equal(test, pressedBackgroundColor, characterMemory[3][3].AttributeEntry.BackgroundColor)
}

/*
TestButtonBorderlessAllowsSingleRow is a test which verifies that a borderless button is not forced to the
three-row minimum that the framed styles need.

Example:
    Expected Inputs:
        A borderless button "Test" at (2,2) with a requested height of 1.

    Expected Outputs:
        Row 2 in the button's column range carries the button cell type; row 3 does not.
*/
func TestButtonBorderlessAllowsSingleRow(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	styleEntry.Button.StyleMode = constants.ButtonStyleBorderless
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 1, true)
	UpdateDisplay(false)
	characterMemory := commonResource.screenLayer.CharacterMemory

	assert.Equal(test, constants.CellTypeButton, characterMemory[2][3].AttributeEntry.CellType)
	assert.NotEqual(test, constants.CellTypeButton, characterMemory[3][3].AttributeEntry.CellType)
}
