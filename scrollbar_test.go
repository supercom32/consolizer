package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

const SCROLLBAR_TEST_SUITE_NAME = "scrollbar"

/*
TestScrollbarDefaultState is a test which verifies that a scrollbar control is rendered correctly with its default
(vertical) orientation and initial scroll value.

Example:
    Expected Inputs:
        A vertical scrollbar added to a layer with length 10 and max scroll value 100.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the track, arrows, and handle.
*/
func TestScrollbarDefaultState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	layer1.AddScrollbar(styleEntry, 5, 2, 10, 100, 0, 1, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarDefaultState", obtainedValue)
	expectedValue := LoadMasterImage(SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarDefaultState")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestScrollbarHorizontal is a test which verifies that a scrollbar control is rendered correctly with a horizontal
orientation and a specific scroll value.

Example:
    Expected Inputs:
        A horizontal scrollbar added with length 10 and current scroll value 70.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the horizontal track and handle position.
*/
func TestScrollbarHorizontal(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	layer1.AddScrollbar(styleEntry, 2, 2, 10, 100, 70, 1, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarHorizontal", obtainedValue)
	expectedValue := LoadMasterImage(SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarHorizontal")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestScrollbarWithValue is a test which verifies that a scrollbar initialized with a non-zero value reports that value
correctly through its getter method.

Example:
    Expected Inputs:
        A scrollbar initialized with a value of 50.
    Expected Outputs:
        The getScrollValue method returns exactly 50 and the handle is correctly positioned on screen.
*/
func TestScrollbarWithValue(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbarInstance := layer1.AddScrollbar(styleEntry, 2, 2, 10, 100, 50, 1, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarWithValue", obtainedValue)
	expectedValue := LoadMasterImage(SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarWithValue")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the scrollbar value
	if !assert.Equal(test, 50, scrollbarInstance.getScrollValue(), "Scrollbar value should be 50") {
		fmt.Println("Scrollbar value is incorrect")
	}
}

/*
TestScrollbarSetValue is a test which verifies that the current value of a scrollbar can be updated programmatically
and that it renders at the new position.

Example:
    Expected Inputs:
        A scrollbar updated from value 75 to 30 using setScrollValue.
    Expected Outputs:
        The scrollbar value becomes 30 and the handle moves to the corresponding track position.
*/
func TestScrollbarSetValue(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbarInstance := layer1.AddScrollbar(styleEntry, 2, 2, 10, 100, 75, 1, false)
	scrollbarInstance.setScrollValue(30)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarSetValue", obtainedValue)
	expectedValue := LoadMasterImage(SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarSetValue")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the scrollbar value
	if !assert.Equal(test, 30, scrollbarInstance.getScrollValue(), "Scrollbar value not what was expected") {
		fmt.Println("Scrollbar value is incorrect")
	}
}

/*
TestScrollbarSetHandlePosition is a test which verifies that manually setting the handle position correctly updates
the underlying scroll value.

Example:
    Expected Inputs:
        A scrollbar where the handle position is explicitly set to index 5.
    Expected Outputs:
        The scrollbar value is automatically calculated and updated based on the new handle position.
*/
func TestScrollbarSetHandlePosition(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbarInstance := layer1.AddScrollbar(styleEntry, 2, 2, 10, 100, 0, 1, false)

	// Set the scrollbar handle position
	scrollbarInstance.setHandlePosition(5)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarSetHandlePosition", obtainedValue)
	expectedValue := LoadMasterImage(SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarSetHandlePosition")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestScrollbarMaxValue is a test which verifies that the scrollbar correctly handles being set to its maximum value,
clamping it to the valid range if necessary.

Example:
    Expected Inputs:
        A scrollbar with max value 100 (range 0-99) where the value is set to 100.
    Expected Outputs:
        The scrollbar value is clamped to 99 and the handle is at the bottom-most position.
*/
func TestScrollbarMaxValue(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbarInstance := layer1.AddScrollbar(styleEntry, 2, 2, 10, 100, 0, 1, false)

	// Set the scrollbar to max value
	scrollbarInstance.setScrollValue(100)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarMaxValue", obtainedValue)
	expectedValue := LoadMasterImage(SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarMaxValue")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the scrollbar value
	if !assert.Equal(test, 99, scrollbarInstance.getScrollValue(), "Scrollbar value should be 99") {
		fmt.Println("Scrollbar value is incorrect")
	}
}

/*
TestScrollbarMinValue is a test which verifies that the scrollbar correctly handles being set to its minimum value of 0.

Example:
    Expected Inputs:
        A scrollbar at value 50 updated to value 0.
    Expected Outputs:
        The scrollbar value becomes 0 and the handle is at the top-most position.
*/
func TestScrollbarMinValue(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbarInstance := layer1.AddScrollbar(styleEntry, 2, 2, 10, 100, 50, 1, false)

	// Set the scrollbar to min value
	scrollbarInstance.setScrollValue(0)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarMinValue", obtainedValue)
	expectedValue := LoadMasterImage(SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarMinValue")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the scrollbar value
	if !assert.Equal(test, 0, scrollbarInstance.getScrollValue(), "Scrollbar value should be 0") {
		fmt.Println("Scrollbar value is incorrect")
	}
}

/*
TestScrollbarDelete is a test which verifies that a scrollbar control can be successfully removed from its parent layer.

Example:
    Expected Inputs:
        A layer containing a scrollbar followed by a Delete call.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) for an empty layer.
*/
func TestScrollbarDelete(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbarInstance := layer1.AddScrollbar(styleEntry, 2, 2, 10, 100, 0, 1, false)
	UpdateDisplay(false)

	// Delete the scrollbar
	scrollbarInstance.Delete()

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarDelete", obtainedValue)
	expectedValue := LoadMasterImage(SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarDelete")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestScrollbarDeleteAll is a test which verifies that all scrollbar controls on a layer can be successfully removed at once.

Example:
    Expected Inputs:
        A layer containing two scrollbars followed by a DeleteAll call.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) for an empty layer after all controls are removed.
*/
func TestScrollbarDeleteAll(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	layer1.AddScrollbar(styleEntry, 2, 2, 10, 100, 0, 1, false)
	layer1.AddScrollbar(styleEntry, 15, 2, 10, 100, 50, 1, true)
	UpdateDisplay(false)

	// Delete all scrollbars
	scrollbar.DeleteAll(layer1.layerAlias)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarDeleteAll", obtainedValue)
	expectedValue := LoadMasterImage(SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarDeleteAll")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestScrollbarFocus is a test which verifies that a scrollbar correctly handles gaining focus and renders its focused state.

Example:
    Expected Inputs:
        A focused scrollbar control on a layer with another control.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the scrollbar with focused highlighting.
*/
func TestScrollbarFocus(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbarInstance := layer1.AddScrollbar(styleEntry, 2, 2, 10, 100, 0, 1, false)
	layer1.AddButton("Testing", styleEntry, 6, 1, 10, 5, true)
	// Set focus to the scrollbar
	setFocusedControl(layer1.layerAlias, scrollbarInstance.GetAlias(), constants.CellTypeScrollbar)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarFocus", obtainedValue)
	expectedValue := LoadMasterImage(SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarFocus")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestScrollbarKeyboardEvent is a test which verifies that the scrollbar correctly reacts to keyboard navigation events.

Example:
    Expected Inputs:
        A focused scrollbar at value 50 that receives two "up" keystrokes (increment 5).
    Expected Outputs:
        The scroll value decreases to 40 and the handle position is updated accordingly.
*/
func TestScrollbarKeyboardEvent(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbarInstance := layer1.AddScrollbar(styleEntry, 2, 2, 10, 100, 50, 5, false)

	// Set focus to the scrollbar
	setFocusedControl(layer1.layerAlias, scrollbarInstance.GetAlias(), constants.CellTypeScrollbar)

	// Simulate keyboard events - Should move to 3rd scrollbar notch on scrollbar.
	scrollbar.updateKeyboardEventManually(layer1.layerAlias, scrollbarInstance.GetAlias(), []rune("up"))
	scrollbar.updateKeyboardEventManually(layer1.layerAlias, scrollbarInstance.GetAlias(), []rune("up"))

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarKeyboardEvent", obtainedValue)
	expectedValue := LoadMasterImage(SCROLLBAR_TEST_SUITE_NAME, "TestScrollbarKeyboardEvent")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the scrollbar value decreased
	if !assert.Equal(test, 40, scrollbarInstance.getScrollValue(), "Scrollbar value should be 40") {
		fmt.Println("Scrollbar value is incorrect")
	}
}
