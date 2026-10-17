package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const TOOLTIP_TEST_SUITE_NAME = "tooltip"

/*
TestTooltipThinLine is a test which verifies that a tooltip with a thin line border is rendered correctly. In
addition, the following should be noted:

  - Runs Tooltip.updateMouseEvent directly, through lockedTooltipMouseEvent, instead of going through the throttled
    UpdatePeriodicEvents, so this test drives the hover state machine's two required calls (arm, then reveal)
    deterministically. CommonTestSetup starts a background goroutine that also calls UpdatePeriodicEvents on its own
    10ms schedule, so going through that same throttled path here would make whether each call runs depend on the 500ms
    throttle it shares with that goroutine. Each call here takes the UI lock through lockedTooltipMouseEvent, as that
    goroutine does, so they never overlap.

Example:

	Expected Inputs:
		A tooltip with a thin line border added to a layer.
	Expected Outputs:
		Screen content matches expected ANSI string (Base64 encoded) showing the thin border and tooltip text.
*/
func TestTooltipThinLine(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddTooltip("This is a tooltip!", styleEntry, 0, 0, 20, 5, 3, 3, 25, 1, false, false, 0)
	UpdateDisplay(false)
	SetMouseStatus(0, 0, 0, "")
	lockedTooltipMouseEvent()
	time.Sleep(1 * time.Second)
	lockedTooltipMouseEvent()
	UpdateDisplay(false)
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
TestTooltipNoBorder is a test which verifies that a tooltip with no border is rendered correctly. In addition, the
following should be noted:

  - Runs Tooltip.updateMouseEvent directly, through lockedTooltipMouseEvent, instead of going through the throttled
    UpdatePeriodicEvents, so this test drives the hover state machine's two required calls (arm, then reveal)
    deterministically. CommonTestSetup starts a background goroutine that also calls UpdatePeriodicEvents on its own
    10ms schedule, so going through that same throttled path here would make whether each call runs depend on the 500ms
    throttle it shares with that goroutine. Each call here takes the UI lock through lockedTooltipMouseEvent, as that
    goroutine does, so they never overlap.

Example:

	Expected Inputs:
		A tooltip with no border added to a layer.
	Expected Outputs:
		Screen content matches expected ANSI string (Base64 encoded) showing only the tooltip text background.
*/
func TestTooltipNoBorder(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddTooltip("This is a tooltip and this works great!", styleEntry, 0, 0, 20, 5, 3, 3, 15, 10, false, false, 0)
	UpdateDisplay(false)
	SetMouseStatus(0, 0, 0, "")
	lockedTooltipMouseEvent()
	time.Sleep(1 * time.Second)
	lockedTooltipMouseEvent()
	UpdateDisplay(false)
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
TestTooltipNoDelay is a test which verifies that a tooltip with no delay is rendered correctly. In addition, the
following should be noted:

  - Runs Tooltip.updateMouseEvent directly, through lockedTooltipMouseEvent, instead of going through the throttled
    UpdatePeriodicEvents, so this test drives the hover state machine's two required calls (arm, then reveal)
    deterministically. CommonTestSetup starts a background goroutine that also calls UpdatePeriodicEvents on its own
    10ms schedule, so going through that same throttled path here would make whether each call runs depend on the 500ms
    throttle it shares with that goroutine. Each call here takes the UI lock through lockedTooltipMouseEvent, as that
    goroutine does, so they never overlap.

Example:

	Expected Inputs:
		A tooltip with zero hover display delay.
	Expected Outputs:
		Screen content matches expected ANSI string (Base64 encoded) with the tooltip showing immediately.
*/
func TestTooltipNoDelay(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddTooltip("This is a tooltip and this works great!", styleEntry, 0, 0, 20, 5, 3, 3, 15, 10, false, true, 0)
	UpdateDisplay(false)
	SetMouseStatus(0, 0, 0, "")
	lockedTooltipMouseEvent()
	time.Sleep(1 * time.Second)
	lockedTooltipMouseEvent()
	UpdateDisplay(false)
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
TestTooltipWithDelay is a test which verifies that a tooltip with a delay is rendered correctly. In addition, the
following should be noted:

  - Runs Tooltip.updateMouseEvent directly, through lockedTooltipMouseEvent, instead of going through the throttled
    UpdatePeriodicEvents, so this test drives the hover state machine's calls (arm, then check before and after the
    delay) deterministically. CommonTestSetup starts a background goroutine that also calls UpdatePeriodicEvents on its
    own 10ms schedule, so going through that same throttled path here would make whether each call runs depend on the
    500ms throttle it shares with that goroutine. Each call here takes the UI lock through lockedTooltipMouseEvent, as
    that goroutine does, so they never overlap.

Example:

	Expected Inputs:
		A tooltip with a 2000ms display delay, checked before and after the delay elapsed.
	Expected Outputs:
		Tooltip is hidden initially and then becomes visible after the 2-second delay.
*/
func TestTooltipWithDelay(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddTooltip("This is a tooltip and this works great!", styleEntry, 0, 0, 20, 5, 3, 3, 15, 10, false, true, 2000)
	UpdateDisplay(false)
	SetMouseStatus(0, 0, 0, "")
	lockedTooltipMouseEvent()
	time.Sleep(1 * time.Second)
	lockedTooltipMouseEvent()
	UpdateDisplay(false)
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
	lockedTooltipMouseEvent()
	UpdateDisplay(false)
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

/*
TestTooltipTextDelayAreaAndShowing is a test which allows you to verify that a tooltip's text, hover delay, hotspot
and display area can be read or set through its methods, with invalid values ignored or clamped, and that IsShowing
reports when the tooltip is shown. In addition, the following should be noted:

  - This also covers a bug where every hotspot was marked in the layer's top left corner rather than at its own
    location, so a tooltip whose hotspot did not start at (0, 0) could never be shown.

Example:

	Expected Inputs:
		A tooltip with the text "Tip", hotspot (0, 0) sized 5 by 1, and no delay. SetText("New"), SetHoverDelay with
		250 and -5, SetHotspot with (2, 3, 4, 1) and (1, 1, 0, 1), SetDisplayArea with (1, 5, 12, 2) and
		(0, 0, 3, 0). Then the mouse rested on (3, 3) and the hover state machine run twice.

	Expected Outputs:
		GetText reads "Tip", then "New". GetHoverDelay reads 0, 250, then 0. The hotspot is (2, 3, 4, 1) and the
		display area (1, 5, 12, 2) after the invalid sizes. IsShowing reads false before the mouse rests on the
		hotspot and true after.
*/
func TestTooltipTextDelayAreaAndShowing(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	layerInstance := &LayerInstanceType{layerAlias: layerAlias}
	tooltipInstance := layerInstance.AddTooltip("Tip", styleEntry, 0, 0, 5, 1, 0, 1, 10, 1, true, false, 0)
	if text := tooltipInstance.GetText(); text != "Tip" {
		test.Fatalf("expected the text %q, got %q", "Tip", text)
	}
	tooltipInstance.SetText("New")
	if text := tooltipInstance.GetText(); text != "New" {
		test.Fatalf("expected the text %q, got %q", "New", text)
	}
	for _, step := range []struct{ delay, expected int }{{250, 250}, {-5, 0}} {
		tooltipInstance.SetHoverDelay(step.delay)
		if delay := tooltipInstance.GetHoverDelay(); delay != step.expected {
			test.Fatalf("SetHoverDelay(%d): expected %d, got %d", step.delay, step.expected, delay)
		}
	}
	tooltipInstance.SetHotspot(2, 3, 4, 1)
	tooltipInstance.SetHotspot(1, 1, 0, 1)
	tooltipInstance.SetDisplayArea(1, 5, 12, 2)
	tooltipInstance.SetDisplayArea(0, 0, 3, 0)
	func() {
		defer lockUIState()()
		tooltipEntry := tooltipRegistry.Get(layerAlias, tooltipInstance.GetAlias())
		if tooltipEntry.HotspotXLocation != 2 || tooltipEntry.HotspotYLocation != 3 || tooltipEntry.HotspotWidth != 4 || tooltipEntry.HotspotHeight != 1 {
			test.Fatalf("expected the hotspot (2, 3, 4, 1), got (%d, %d, %d, %d)", tooltipEntry.HotspotXLocation,
				tooltipEntry.HotspotYLocation, tooltipEntry.HotspotWidth, tooltipEntry.HotspotHeight)
		}
		if tooltipEntry.TooltipXLocation != 1 || tooltipEntry.TooltipYLocation != 5 || tooltipEntry.TooltipWidth != 12 || tooltipEntry.TooltipHeight != 2 {
			test.Fatalf("expected the display area (1, 5, 12, 2), got (%d, %d, %d, %d)", tooltipEntry.TooltipXLocation,
				tooltipEntry.TooltipYLocation, tooltipEntry.TooltipWidth, tooltipEntry.TooltipHeight)
		}
	}()
	if tooltipInstance.IsShowing() {
		test.Fatalf("expected the tooltip to be hidden before the mouse rests on it")
	}
	UpdateDisplay(false)
	SetMouseStatus(3, 3, 0, "")
	lockedTooltipMouseEvent()
	time.Sleep(100 * time.Millisecond)
	lockedTooltipMouseEvent()
	if !tooltipInstance.IsShowing() {
		test.Fatalf("expected the tooltip to be showing once the mouse rested on its hotspot")
	}
}
