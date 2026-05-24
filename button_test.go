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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
