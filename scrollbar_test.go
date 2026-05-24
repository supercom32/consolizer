package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

const SCROLLBAR_TEST_SUITE_NAME = "scrollbar"

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
