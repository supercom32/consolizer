package consolizer

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
	"time"
)

const SELECTOR_TEST_SUITE_NAME = "selector"

/*
TestSelectorRandomSelection is a test which verifies that a selector control correctly renders a selection made at
an arbitrary index.

Example:
    Expected Inputs:
        A selector with 4 items, where the second item (index 1) is programmatically selected.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the selected item highlighted.
*/
func TestSelectorRandomSelection(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	selectionEntry := NewSelectionEntry()
	selectionEntry.Add("Selection Alias 1", "Selection Text 1")
	selectionEntry.Add("Selection Alias 2", "Selection Text 2")
	selectionEntry.Add("Selection Alias 3", "Selection Text 3")
	selectionEntry.Add("Selection Alias 4", "Selection Text 4")
	selectorFieldInstance := layer1.AddSelector(styleEntry, selectionEntry, 2, 2, 4, 25, 1, 0, 1, true, true)
	setFocusedControl(layer1.layerAlias, selectorFieldInstance.controlAlias, constants.CellTypeTextField)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SELECTOR_TEST_SUITE_NAME, "TestSelectorRandomSelection", obtainedValue)
	expectedValue := LoadMasterImage(SELECTOR_TEST_SUITE_NAME, "TestSelectorRandomSelection")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestSelectorLongList is a test which verifies that a selector control correctly handles a long list of items,
including the rendering of a scrollbar when necessary.

Example:
    Expected Inputs:
        A selector with 8 items and a viewport height of 4.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing a subset of items and a scrollbar.
*/
func TestSelectorLongList(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	selectionEntry := NewSelectionEntry()
	selectionEntry.Add("Selection Alias 1", "Selection Text 1")
	selectionEntry.Add("Selection Alias 2", "Selection Text 2")
	selectionEntry.Add("Selection Alias 3", "Selection Text 3")
	selectionEntry.Add("Selection Alias 4", "Selection Text 4")
	selectionEntry.Add("Selection Alias 5", "Selection Text 5")
	selectionEntry.Add("Selection Alias 6", "Selection Text 6")
	selectionEntry.Add("Selection Alias 7", "Selection Text 7")
	selectionEntry.Add("Selection Alias 8", "Selection Text 8")
	selectorFieldInstance := layer1.AddSelector(styleEntry, selectionEntry, 2, 2, 4, 25, 1, 0, 0, true, true)
	setFocusedControl(layer1.layerAlias, selectorFieldInstance.controlAlias, constants.CellTypeTextField)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SELECTOR_TEST_SUITE_NAME, "TestSelectorLongList", obtainedValue)
	expectedValue := LoadMasterImage(SELECTOR_TEST_SUITE_NAME, "TestSelectorLongList")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestGetAllItems is a test which verifies that the GetAllItems method correctly retrieves all aliases and values
from a selector control.

Example:
    Expected Inputs:
        A selector containing 4 items with known aliases and values.
    Expected Outputs:
        The returned alias and value slices exactly match the expected input collections.
*/
func TestGetAllItems(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	selectionEntry := NewSelectionEntry()

	// Add some items to the selection entry
	expectedAliases := []string{"Alias 1", "Alias 2", "Alias 3", "Alias 4"}
	expectedValues := []string{"Value 1", "Value 2", "Value 3", "Value 4"}

	for i := 0; i < len(expectedAliases); i++ {
		selectionEntry.Add(expectedAliases[i], expectedValues[i])
	}

	// Create a selector with the selection entry
	selectorFieldInstance := layer1.AddSelector(styleEntry, selectionEntry, 2, 2, 4, 25, 1, 0, 0, true, true)

	// Call GetAllItems and verify the results
	aliases, values := selectorFieldInstance.GetAllItems()

	// Check that the returned arrays match the expected values
	assert.Equal(test, expectedAliases, aliases, "The returned aliases do not match the expected values")
	assert.Equal(test, expectedValues, values, "The returned values do not match the expected values")

	// Test with an empty selector
	emptySelectionEntry := NewSelectionEntry()
	emptySelectorFieldInstance := layer1.AddSelector(styleEntry, emptySelectionEntry, 10, 10, 4, 25, 1, 0, 0, true, true)
	emptyAliases, emptyValues := emptySelectorFieldInstance.GetAllItems()

	// Check that empty arrays are returned for an empty selector
	assert.Empty(test, emptyAliases, "The returned aliases should be empty for an empty selector")
	assert.Empty(test, emptyValues, "The returned values should be empty for an empty selector")
}

/*
TestSelectorLongListWithColors is a test which verifies that a selector control correctly renders items that
contain markup tags for colorization.

Example:
    Expected Inputs:
        A selector containing an item with "{{red}}Text{{/}}" markup.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the item with correct color formatting.
*/
func TestSelectorLongListWithColors(test *testing.T) {
	textStyleAlias := "red"
	attributeEntry := NewTextStyle()
	attributeEntry.ForegroundColor = GetRGBColor(255, 0, 0)
	AddTextStyle(textStyleAlias, attributeEntry)
	layer1, _, _, styleEntry := CommonTestSetup(test)
	selectionEntry := NewSelectionEntry()
	selectionEntry.Add("Selection Alias 1", "Selection Text 1")
	selectionEntry.Add("Selection Alias 2", "Selection Text 2")
	selectionEntry.Add("Selection Alias 3", "Selection {{red}}Text{{/}} 3")
	selectionEntry.Add("Selection Alias 4", "Selection Text 4")
	selectionEntry.Add("Selection Alias 5", "Selection Text 5")
	selectionEntry.Add("Selection Alias 6", "Selection Text 6")
	selectionEntry.Add("Selection Alias 7", "Selection Text 7")
	selectionEntry.Add("Selection Alias 8", "Selection Text 8")
	selectorFieldInstance := layer1.AddSelector(styleEntry, selectionEntry, 2, 2, 4, 25, 1, 0, 0, true, true)
	setFocusedControl(layer1.layerAlias, selectorFieldInstance.controlAlias, constants.CellTypeTextField)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, SELECTOR_TEST_SUITE_NAME, "TestSelectorLongListWithColors", obtainedValue)
	expectedValue := LoadMasterImage(SELECTOR_TEST_SUITE_NAME, "TestSelectorLongListWithColors")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestFocusSelectionInitialPosition is a test which verifies that a selector control initializes with the correct
viewport position of 0.

Example:
    Expected Inputs:
        A selector with 20 items and a viewport height of 4.
    Expected Outputs:
        The initial ViewportPosition is 0 and the screen content matches the default view.
*/
func TestFocusSelectionInitialPosition(test *testing.T) {
	// Setup
	layer1, _, _, styleEntry := CommonTestSetup(test)
	selectionEntry := NewSelectionEntry()

	// Create a selector with many items to ensure scrolling is necessary
	for i := 1; i <= 20; i++ {
		alias := fmt.Sprintf("Alias %d", i)
		value := fmt.Sprintf("Value %d", i)
		selectionEntry.Add(alias, value)
	}

	// Create a selector with a viewport height of 4 (can show 4 items at once)
	selectorFieldInstance := layer1.AddSelector(styleEntry, selectionEntry, 2, 2, 4, 25, 1, 0, 0, true, true)

	// Get the selector entry to check viewport position
	selectorEntry := GetSelector(layer1.layerAlias, selectorFieldInstance.controlAlias)

	// Initial viewport position should be 0
	assert.Equal(test, 0, selectorEntry.ViewportPosition, "Initial viewport position should be 0")

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	// Using a valid base64 string as a placeholder
	UpdateMasterImages(false, SELECTOR_TEST_SUITE_NAME, "TestFocusSelectionInitialPosition", obtainedValue)
	expectedValue := LoadMasterImage(SELECTOR_TEST_SUITE_NAME, "TestFocusSelectionInitialPosition")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestFocusSelectionMiddleItem is a test which verifies that the FocusSelection method correctly adjusts the
viewport position to center a selected item.

Example:
    Expected Inputs:
        A selector with 20 items where FocusSelection is called for "Alias 10".
    Expected Outputs:
        The ViewportPosition is adjusted to 7 to center the 10th item (index 9) in a height-4 viewport.
*/
func TestFocusSelectionMiddleItem(test *testing.T) {
	// Setup
	layer1, _, _, styleEntry := CommonTestSetup(test)
	selectionEntry := NewSelectionEntry()

	// Create a selector with many items to ensure scrolling is necessary
	for i := 1; i <= 20; i++ {
		alias := fmt.Sprintf("Alias %d", i)
		value := fmt.Sprintf("Value %d", i)
		selectionEntry.Add(alias, value)
	}

	// Create a selector with a viewport height of 4 (can show 4 items at once)
	selectorFieldInstance := layer1.AddSelector(styleEntry, selectionEntry, 2, 2, 4, 25, 1, 0, 0, true, true)

	// Get the selector entry to check viewport position
	selectorEntry := GetSelector(layer1.layerAlias, selectorFieldInstance.controlAlias)

	// Focus on an item in the middle (Alias 10)
	selectorFieldInstance.FocusSelection("Alias 10")

	// The viewport position should be adjusted to center Alias 10
	// With 4 visible items, and centering Alias 10 (index 9),
	// the viewport position should be 7 (9 - 4/2 = 7)
	expectedPosition := 7 // Integer division: 9 - (4/2) = 9 - 2 = 7
	assert.Equal(test, expectedPosition, selectorEntry.ViewportPosition,
		"Viewport position should be adjusted to center the selected item")

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	// Using a valid base64 string as a placeholder
	UpdateMasterImages(false, SELECTOR_TEST_SUITE_NAME, "TestFocusSelectionMiddleItem", obtainedValue)
	expectedValue := LoadMasterImage(SELECTOR_TEST_SUITE_NAME, "TestFocusSelectionMiddleItem")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestFocusSelectionEndItem is a test which verifies that FocusSelection correctly clamps the viewport position
to the bottom of the list.

Example:
    Expected Inputs:
        A selector with 20 items where FocusSelection is called for "Alias 18".
    Expected Outputs:
        The ViewportPosition is adjusted but correctly limited by the maximum possible value of 16.
*/
func TestFocusSelectionEndItem(test *testing.T) {
	// Setup
	layer1, _, _, styleEntry := CommonTestSetup(test)
	selectionEntry := NewSelectionEntry()

	// Create a selector with many items to ensure scrolling is necessary
	for i := 1; i <= 20; i++ {
		alias := fmt.Sprintf("Alias %d", i)
		value := fmt.Sprintf("Value %d", i)
		selectionEntry.Add(alias, value)
	}

	// Create a selector with a viewport height of 4 (can show 4 items at once)
	selectorFieldInstance := layer1.AddSelector(styleEntry, selectionEntry, 2, 2, 4, 25, 1, 0, 0, true, true)

	// Get the selector entry to check viewport position
	selectorEntry := GetSelector(layer1.layerAlias, selectorFieldInstance.controlAlias)

	// Focus on an item near the end (Alias 18)
	selectorFieldInstance.FocusSelection("Alias 18")

	// The viewport position should be adjusted, but limited by the max position
	// Max position = 20 items - 4 visible items = 16
	// For Alias 18 (index 17): 17 - (4/2) = 17 - 2 = 15
	expectedPosition := 15
	assert.Equal(test, expectedPosition, selectorEntry.ViewportPosition,
		"Viewport position should be adjusted but not exceed the maximum")

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	// Using a valid base64 string as a placeholder
	UpdateMasterImages(false, SELECTOR_TEST_SUITE_NAME, "TestFocusSelectionEndItem", obtainedValue)
	expectedValue := LoadMasterImage(SELECTOR_TEST_SUITE_NAME, "TestFocusSelectionEndItem")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestFocusSelectionBeginningItem is a test which verifies that FocusSelection correctly clamps the viewport
position to 0 when an item near the start is focused.

Example:
    Expected Inputs:
        A selector with 20 items where FocusSelection is called for "Alias 2".
    Expected Outputs:
        The ViewportPosition is set to 0 and not allowed to become negative.
*/
func TestFocusSelectionBeginningItem(test *testing.T) {
	// Setup
	layer1, _, _, styleEntry := CommonTestSetup(test)
	selectionEntry := NewSelectionEntry()

	// Create a selector with many items to ensure scrolling is necessary
	for i := 1; i <= 20; i++ {
		alias := fmt.Sprintf("Alias %d", i)
		value := fmt.Sprintf("Value %d", i)
		selectionEntry.Add(alias, value)
	}

	// Create a selector with a viewport height of 4 (can show 4 items at once)
	selectorFieldInstance := layer1.AddSelector(styleEntry, selectionEntry, 2, 2, 4, 25, 1, 0, 0, true, true)

	// Get the selector entry to check viewport position
	selectorEntry := GetSelector(layer1.layerAlias, selectorFieldInstance.controlAlias)

	// Focus on an item near the beginning (Alias 2)
	selectorFieldInstance.FocusSelection("Alias 2")

	// The viewport position should be adjusted, but limited by the min position (0)
	expectedPosition := 0
	assert.Equal(test, expectedPosition, selectorEntry.ViewportPosition,
		"Viewport position should be adjusted but not go below 0")

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	// Using a valid base64 string as a placeholder
	UpdateMasterImages(false, SELECTOR_TEST_SUITE_NAME, "TestFocusSelectionBeginningItem", obtainedValue)
	expectedValue := LoadMasterImage(SELECTOR_TEST_SUITE_NAME, "TestFocusSelectionBeginningItem")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestFocusSelectionNonExistentItem is a test which verifies that FocusSelection does not change the viewport
position when called with a non-existent item alias.

Example:
    Expected Inputs:
        A selector with 20 items where FocusSelection is called for a missing alias.
    Expected Outputs:
        The ViewportPosition remains at its previous value.
*/
func TestFocusSelectionNonExistentItem(test *testing.T) {
	// Setup
	layer1, _, _, styleEntry := CommonTestSetup(test)
	selectionEntry := NewSelectionEntry()

	// Create a selector with many items to ensure scrolling is necessary
	for i := 1; i <= 20; i++ {
		alias := fmt.Sprintf("Alias %d", i)
		value := fmt.Sprintf("Value %d", i)
		selectionEntry.Add(alias, value)
	}

	// Create a selector with a viewport height of 4 (can show 4 items at once)
	selectorFieldInstance := layer1.AddSelector(styleEntry, selectionEntry, 2, 2, 4, 25, 1, 0, 0, true, true)

	// Get the selector entry to check viewport position
	selectorEntry := GetSelector(layer1.layerAlias, selectorFieldInstance.controlAlias)

	// Set a known initial position
	selectorFieldInstance.FocusSelection("Alias 5")
	initialPosition := selectorEntry.ViewportPosition

	// Focus on a non-existent item
	selectorFieldInstance.FocusSelection("Non-existent Alias")

	// The viewport position should remain unchanged
	assert.Equal(test, initialPosition, selectorEntry.ViewportPosition,
		"Viewport position should remain unchanged when focusing on a non-existent item")

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	// Using a valid base64 string as a placeholder
	UpdateMasterImages(false, SELECTOR_TEST_SUITE_NAME, "TestFocusSelectionNonExistentItem", obtainedValue)
	expectedValue := LoadMasterImage(SELECTOR_TEST_SUITE_NAME, "TestFocusSelectionNonExistentItem")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestSelectorSelectionSourceKeyboardAndProgrammatic is a test which allows you to verify that a selector reports no
source before any selection, the keyboard source for Enter, and the programmatic source for Select.

Example:

	Expected Inputs:
		Selector with items [a, b, c, d], focused with SetFocus. Down then Enter, then Select("c").

	Expected Outputs:
		SelectionSourceNone before any selection, item 1 with SelectionSourceKeyboard after Enter, and item 2 with
		SelectionSourceProgrammatic after Select.
*/
func TestSelectorSelectionSourceKeyboardAndProgrammatic(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := addTestSelector(layerAlias, styleEntry)
	if err := SetFocus(&selector); err != nil {
		test.Fatalf("expected the selector to take focus, got %v", err)
	}
	simScreen := startInputSimulation(test)

	if source := selector.GetSelectionSource(); source != constants.SelectionSourceNone {
		test.Fatalf("expected no selection source before any selection, got %d", source)
	}

	pressKey(simScreen, tcell.KeyDown, 0, tcell.ModNone)
	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	assertSelection(test, selector, "after Enter", 1, constants.SelectionSourceKeyboard)

	selector.Select("c")
	assertSelection(test, selector, "after Select", 2, constants.SelectionSourceProgrammatic)
}

/*
TestSelectorSelectionSourceSingleAndDoubleClick is a test which allows you to verify that clicks on a selector item are
reported as a single click, then a double click for a second click on the same item within the interval, then a single
click again for a third click, and that a click on a different item is always a single click.

Example:

	Expected Inputs:
		Selector with items [a, b, c, d]. Three quick clicks on item 1, then one click on item 2.

	Expected Outputs:
		Item 1 with SingleClick, item 1 with DoubleClick, item 1 with SingleClick, then item 2 with SingleClick.
*/
func TestSelectorSelectionSourceSingleAndDoubleClick(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := addTestSelector(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)

	clickAt(simScreen, 4, 3)
	assertSelection(test, selector, "first click", 1, constants.SelectionSourceSingleClick)

	clickAt(simScreen, 4, 3)
	assertSelection(test, selector, "second click", 1, constants.SelectionSourceDoubleClick)

	clickAt(simScreen, 4, 3)
	assertSelection(test, selector, "third click", 1, constants.SelectionSourceSingleClick)

	clickAt(simScreen, 4, 4)
	assertSelection(test, selector, "click on another item", 2, constants.SelectionSourceSingleClick)
}

/*
TestSelectorDragToAnotherItemIsSingleClick is a test which allows you to verify that pressing on one selector item and
dragging with the button held onto another item selects the second item as a single click, and that a quick click on
that item afterward is not mistaken for a double click.

Example:

	Expected Inputs:
		Selector with items [a, b, c, d]. Press on item 1, move with the button held to item 2, release, then click
		item 2.

	Expected Outputs:
		Item 2 with SingleClick after the release, and item 2 with SingleClick after the click.
*/
func TestSelectorDragToAnotherItemIsSingleClick(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := addTestSelector(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)

	simScreen.InjectMouse(4, 3, tcell.Button1, tcell.ModNone)
	UpdateEventQueues()
	simScreen.InjectMouse(4, 4, tcell.Button1, tcell.ModNone)
	UpdateEventQueues()
	simScreen.InjectMouse(4, 4, tcell.ButtonNone, tcell.ModNone)
	UpdateEventQueues()
	assertSelection(test, selector, "after drag", 2, constants.SelectionSourceSingleClick)

	clickAt(simScreen, 4, 4)
	assertSelection(test, selector, "click after drag", 2, constants.SelectionSourceSingleClick)
}

/*
TestSelectorDoubleClickInterval is a test which allows you to verify that two clicks on the same item only count as a
double click when the second lands within the configured interval, that the interval can be changed, and that an
interval which is not greater than zero is rejected without changing the current one.

Example:

	Expected Inputs:
		Clicks on the same item 600ms apart with the default 500ms interval, then with the interval set to 1s.
		SetDoubleClickInterval with 0 and with -1ms.

	Expected Outputs:
		SingleClick then SingleClick with the default interval, SingleClick then DoubleClick with the 1s interval.
		Both invalid intervals return an error and the interval stays 1s.
*/
func TestSelectorDoubleClickInterval(test *testing.T) {
	setupInputTest(test)
	if interval := Selector.GetDoubleClickInterval(); interval != 500*time.Millisecond {
		test.Fatalf("expected default interval of 500ms, got %v", interval)
	}
	startTime := time.Now()

	firstSource := Selector.getMouseSelectionSource("layer", "selector", 1, startTime)
	secondSource := Selector.getMouseSelectionSource("layer", "selector", 1, startTime.Add(600*time.Millisecond))
	if firstSource != constants.SelectionSourceSingleClick || secondSource != constants.SelectionSourceSingleClick {
		test.Fatalf("default interval: expected two single clicks, got %d and %d", firstSource, secondSource)
	}

	if err := Selector.SetDoubleClickInterval(time.Second); err != nil {
		test.Fatalf("expected a 1s interval to be accepted, got %v", err)
	}
	selectorLastClick = selectorClickType{}
	firstSource = Selector.getMouseSelectionSource("layer", "selector", 1, startTime)
	secondSource = Selector.getMouseSelectionSource("layer", "selector", 1, startTime.Add(600*time.Millisecond))
	if firstSource != constants.SelectionSourceSingleClick || secondSource != constants.SelectionSourceDoubleClick {
		test.Fatalf("1s interval: expected single then double click, got %d and %d", firstSource, secondSource)
	}

	for _, invalidInterval := range []time.Duration{0, -time.Millisecond} {
		if err := Selector.SetDoubleClickInterval(invalidInterval); err == nil {
			test.Fatalf("expected interval %v to be rejected", invalidInterval)
		}
	}
	if interval := Selector.GetDoubleClickInterval(); interval != time.Second {
		test.Fatalf("expected rejected intervals to leave 1s in place, got %v", interval)
	}
}
