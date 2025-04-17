package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

func TestScrollbarDefaultState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbar.Add(layer1.layerAlias, "testScrollbar", styleEntry, 2, 2, 10, 100, 0, 1, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := ""
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

func TestScrollbarHorizontal(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbar.Add(layer1.layerAlias, "testScrollbar", styleEntry, 2, 2, 10, 100, 0, 1, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := ""
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

func TestScrollbarWithValue(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbarInstance := scrollbar.Add(layer1.layerAlias, "testScrollbar", styleEntry, 2, 2, 10, 100, 50, 1, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := ""
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
	scrollbarInstance := scrollbar.Add(layer1.layerAlias, "testScrollbar", styleEntry, 2, 2, 10, 100, 0, 1, false)

	// Set the scrollbar value
	scrollbarInstance.setScrollValue(75)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := ""
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the scrollbar value
	if !assert.Equal(test, 75, scrollbarInstance.getScrollValue(), "Scrollbar value should be 75") {
		fmt.Println("Scrollbar value is incorrect")
	}
}

func TestScrollbarSetHandlePosition(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbarInstance := scrollbar.Add(layer1.layerAlias, "testScrollbar", styleEntry, 2, 2, 10, 100, 0, 1, false)

	// Set the scrollbar handle position
	scrollbarInstance.setHandlePosition(5)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := ""
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

func TestScrollbarMaxValue(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbarInstance := scrollbar.Add(layer1.layerAlias, "testScrollbar", styleEntry, 2, 2, 10, 100, 0, 1, false)

	// Set the scrollbar to max value
	scrollbarInstance.setScrollValue(99)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := ""
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
	scrollbarInstance := scrollbar.Add(layer1.layerAlias, "testScrollbar", styleEntry, 2, 2, 10, 100, 50, 1, false)

	// Set the scrollbar to min value
	scrollbarInstance.setScrollValue(0)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := ""
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
	scrollbarInstance := scrollbar.Add(layer1.layerAlias, "testScrollbar", styleEntry, 2, 2, 10, 100, 0, 1, false)
	UpdateDisplay(false)

	// Delete the scrollbar
	scrollbarInstance.Delete()

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := ""
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

func TestScrollbarDeleteAll(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbar.Add(layer1.layerAlias, "testScrollbar1", styleEntry, 2, 2, 10, 100, 0, 1, false)
	scrollbar.Add(layer1.layerAlias, "testScrollbar2", styleEntry, 15, 2, 10, 100, 50, 1, true)
	UpdateDisplay(false)

	// Delete all scrollbars
	scrollbar.DeleteAllScrollbars(layer1.layerAlias)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := ""
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

func TestScrollbarFocus(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbar.Add(layer1.layerAlias, "testScrollbar", styleEntry, 2, 2, 10, 100, 0, 1, false)

	// Set focus to the scrollbar
	setFocusedControl(layer1.layerAlias, "testScrollbar", constants.CellTypeScrollbar)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := ""
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

func TestScrollbarKeyboardEvent(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	scrollbarInstance := scrollbar.Add(layer1.layerAlias, "testScrollbar", styleEntry, 2, 2, 10, 100, 50, 5, false)

	// Set focus to the scrollbar
	setFocusedControl(layer1.layerAlias, "testScrollbar", constants.CellTypeScrollbar)

	// Simulate keyboard events
	scrollbar.updateKeyboardEventManually(layer1.layerAlias, "testScrollbar", []rune("up"))

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	expectedValue := ""
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the scrollbar value decreased
	if !assert.Equal(test, 45, scrollbarInstance.getScrollValue(), "Scrollbar value should be 45") {
		fmt.Println("Scrollbar value is incorrect")
	}
}
