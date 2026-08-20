package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const TOOLTIP_TEST_SUITE_NAME = "tooltip"

/*
TestTooltipThinLine is a test which verifies that a tooltip with a thin line border is rendered correctly.

Example:
    Expected Inputs:
        A tooltip with a thin line border added to a layer.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the thin border and tooltip text.
*/
func TestTooltipThinLine(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	layer1.AddTooltip("This is a tooltip!", styleEntry, 0, 0, 20, 5, 3, 3, 25, 1, false, false, 0)
	SetMouseStatus(0, 0, 0, "")
	time.Sleep(1 * time.Second)
	UpdateDisplay(false)
	UpdatePeriodicEvents()
	time.Sleep(1 * time.Second)
	UpdateDisplay(false)
	UpdatePeriodicEvents()
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TOOLTIP_TEST_SUITE_NAME, "TestTooltipThinLine", obtainedValue)
	expectedValue := LoadMasterImage(TOOLTIP_TEST_SUITE_NAME, "TestTooltipThinLine")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTooltipNoBorder is a test which verifies that a tooltip with no border is rendered correctly.

Example:
    Expected Inputs:
        A tooltip with no border added to a layer.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing only the tooltip text background.
*/
func TestTooltipNoBorder(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	layer1.AddTooltip("This is a tooltip and this works great!", styleEntry, 0, 0, 20, 5, 3, 3, 15, 10, false, false, 0)
	SetMouseStatus(0, 0, 0, "")
	time.Sleep(1 * time.Second)
	UpdateDisplay(false)
	UpdatePeriodicEvents()
	time.Sleep(1 * time.Second)
	UpdateDisplay(false)
	UpdatePeriodicEvents()
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TOOLTIP_TEST_SUITE_NAME, "TestTooltipNoBorder", obtainedValue)
	expectedValue := LoadMasterImage(TOOLTIP_TEST_SUITE_NAME, "TestTooltipNoBorder")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTooltipNoDelay is a test which verifies that a tooltip with no delay is rendered correctly.

Example:
    Expected Inputs:
        A tooltip with zero hover display delay.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) with the tooltip showing immediately.
*/
func TestTooltipNoDelay(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	layer1.AddTooltip("This is a tooltip and this works great!", styleEntry, 0, 0, 20, 5, 3, 3, 15, 10, false, true, 0)
	SetMouseStatus(0, 0, 0, "")
	time.Sleep(1 * time.Second)
	UpdateDisplay(false)
	UpdatePeriodicEvents()
	time.Sleep(1 * time.Second)
	UpdateDisplay(false)
	UpdatePeriodicEvents()
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TOOLTIP_TEST_SUITE_NAME, "TestTooltipNoDelay", obtainedValue)
	expectedValue := LoadMasterImage(TOOLTIP_TEST_SUITE_NAME, "TestTooltipNoDelay")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTooltipWithDelay is a test which verifies that a tooltip with a delay is rendered correctly.

Example:
    Expected Inputs:
        A tooltip with a 2000ms display delay, checked before and after the delay elapsed.
    Expected Outputs:
        Tooltip is hidden initially and then becomes visible after the 2-second delay.
*/
func TestTooltipWithDelay(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	layer1.AddTooltip("This is a tooltip and this works great!", styleEntry, 0, 0, 20, 5, 3, 3, 15, 10, false, true, 2000)
	SetMouseStatus(0, 0, 0, "")
	UpdateDisplay(false)
	UpdatePeriodicEvents()
	time.Sleep(1 * time.Second)
	UpdateDisplay(false)
	UpdatePeriodicEvents()
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TOOLTIP_TEST_SUITE_NAME, "TestTooltipWithDelay", obtainedValue)
	expectedValue := LoadMasterImage(TOOLTIP_TEST_SUITE_NAME, "TestTooltipWithDelay")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated first screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
	time.Sleep(2 * time.Second)
	UpdateDisplay(false)
	UpdatePeriodicEvents()
	layerEntry = commonResource.screenLayer
	obtainedValue = layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TOOLTIP_TEST_SUITE_NAME, "TestTooltipWithDelay_Showing", obtainedValue)
	expectedValue = LoadMasterImage(TOOLTIP_TEST_SUITE_NAME, "TestTooltipWithDelay_Showing")
	obtainedValueBase64 = layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 = layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated second screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}
