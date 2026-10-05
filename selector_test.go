package consolizer

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
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

/*
getTestSelectionEntry is a method which allows you to build a selection entry holding a given number of items, whose
aliases are "item0", "item1" and so on and whose display values are "Item 0", "Item 1" and so on.

Example:

	selectionEntry := getTestSelectionEntry(4)
*/
func getTestSelectionEntry(itemCount int) types.SelectionEntryType {
	selectionEntry := NewSelectionEntry()
	for itemIndex := 0; itemIndex < itemCount; itemIndex++ {
		selectionEntry.Add(fmt.Sprintf("item%d", itemIndex), fmt.Sprintf("Item %d", itemIndex))
	}
	return selectionEntry
}

/*
assertSelectorIndexes is a method which allows you to fail the test unless a selector's highlighted item, selected item
and viewport position all hold the values expected.

Example:

	assertSelectorIndexes(test, "after Up", selectorEntry, 0, constants.SELECTED_NONE, 0)
*/
func assertSelectorIndexes(test *testing.T, context string, selectorEntry *types.SelectorEntryType, expectedHighlighted int, expectedSelected int, expectedViewportPosition int) {
	test.Helper()
	if selectorEntry.ItemHighlighted != expectedHighlighted || selectorEntry.ItemSelected != expectedSelected ||
		selectorEntry.ViewportPosition != expectedViewportPosition {
		test.Fatalf("%s: expected highlighted %d, selected %d, viewport %d, got highlighted %d, selected %d, viewport %d",
			context, expectedHighlighted, expectedSelected, expectedViewportPosition,
			selectorEntry.ItemHighlighted, selectorEntry.ItemSelected, selectorEntry.ViewportPosition)
	}
}

/*
TestSelectorEmptyListReceivesFocusAndArrowKeys is a test which allows you to verify the reported crash is fixed: moving
focus with Tab onto an empty selector and pressing Up used to store a viewport position of minus one, which made the
next screen update index the item list at minus one and panic on the event goroutine.

Example:

	Expected Inputs:
		Selector "available" with three items, a disabled button "skip", and selector "signed" with no items and no
		highlighted item, added to the tab order in that order. Focus starts on "available", then Tab, then Up, Down, Left, Right and Enter
		with the screen updated after each key.

	Expected Outputs:
		Tab lands on "signed". No key panics. "signed" keeps highlighted minus one, selected minus one and viewport
		zero, and reports no new selection. The four arrow keys are consumed, so only "enter" reaches the keyboard
		buffer, since Enter selects nothing on an empty list and is left for the application.
*/
func TestSelectorEmptyListReceivesFocusAndArrowKeys(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	availableSelector := Selector.Add(layerAlias, "available", styleEntry, getTestSelectionEntry(3), 2, 2, 4, 10, 1, 0, 0, false, true)
	skipButton := Button.Add(layerAlias, "skip", "skip", styleEntry, 30, 2, 8, 3, true)
	skipButton.SetEnabled(false)
	signedSelector := Selector.Add(layerAlias, "signed", styleEntry, NewSelectionEntry(), 16, 2, 4, 10, 1, 0, constants.NullItemSelection, false, true)
	for _, control := range []*BaseControlInstanceType{&availableSelector.BaseControlInstanceType, &skipButton.BaseControlInstanceType, &signedSelector.BaseControlInstanceType} {
		if err := control.AddToTabIndex(); err != nil {
			test.Fatalf("expected AddToTabIndex to succeed, got %v", err)
		}
	}
	simScreen := startInputSimulation(test)
	if err := SetFocus(&availableSelector); err != nil {
		test.Fatalf("expected SetFocus to succeed, got %v", err)
	}
	pressTab(simScreen)
	assertFocus(test, "after Tab", layerAlias, "signed", constants.CellTypeSelectorItem)

	signedEntry := Selectors.Get(layerAlias, "signed")
	for _, key := range []tcell.Key{tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight, tcell.KeyEnter} {
		pressKey(simScreen, key, 0, tcell.ModNone)
		UpdateDisplay(false)
		assertSelectorIndexes(test, fmt.Sprintf("after key %d", key), signedEntry, constants.NullItemSelection, constants.SELECTED_NONE, 0)
	}
	if signedSelector.IsNewItemSelected() {
		test.Fatalf("expected no new selection on an empty selector")
	}
	assertKeyboardBuffer(test, "after keys on the empty selector", "enter")
}

/*
TestSelectorRefilledWhileFocused is a test which allows you to verify that refilling a focused selector with
SetSelectionEntry leaves it with a valid highlight, where it used to reset the highlight to minus one so that the next
Up stored a viewport position of minus one and the following screen update panicked.

Example:

	Expected Inputs:
		A focused selector with four items and a four row tray, item 2 highlighted, refilled with six items, then Up,
		Down and Down with the screen updated after each key.

	Expected Outputs:
		After the refill the highlight is item 0, the selection is cleared and the viewport is 0. Up keeps item 0,
		then Down moves to items 1 and 2. Nothing panics.
*/
func TestSelectorRefilledWhileFocused(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := Selector.Add(layerAlias, "list", styleEntry, getTestSelectionEntry(4), 2, 2, 4, 10, 1, 0, 2, false, true)
	simScreen := startInputSimulation(test)
	if err := SetFocus(&selector); err != nil {
		test.Fatalf("expected SetFocus to succeed, got %v", err)
	}
	selector.SetSelectionEntry(getTestSelectionEntry(6))
	selectorEntry := Selectors.Get(layerAlias, "list")
	assertSelectorIndexes(test, "after refill", selectorEntry, 0, constants.SELECTED_NONE, 0)
	for _, step := range []struct {
		key                 tcell.Key
		expectedHighlighted int
	}{
		{tcell.KeyUp, 0},
		{tcell.KeyDown, 1},
		{tcell.KeyDown, 2},
	} {
		pressKey(simScreen, step.key, 0, tcell.ModNone)
		UpdateDisplay(false)
		assertSelectorIndexes(test, fmt.Sprintf("after key %d", step.key), selectorEntry, step.expectedHighlighted, constants.SELECTED_NONE, 0)
	}
}

/*
TestSelectorKeyboardEdgeCases is a test which allows you to verify every navigation key against the selector states
that used to produce an invalid index: an empty list, no highlight on a list with items, a highlight left out of
range, a multiple column list, where minus one modulo the column count is minus one, a list shorter than its height,
and a column count of zero, which used to divide by zero.

Example:

	Expected Inputs:
		Each case builds a selector with the given item count, height and column count, sets its highlighted item,
		selected item and viewport position, sends one key directly to the selector's keyboard handler, then
		updates the screen. For example: 7 items, height 2, 3 columns, highlight minus one, key "right".

	Expected Outputs:
		The highlighted item, selected item, viewport position and consumed flag match the case, and nothing panics.
		For the example: highlight 0 is chosen first and Right moves it to 1, selection minus one, viewport 0,
		consumed.
*/
func TestSelectorKeyboardEdgeCases(test *testing.T) {
	testCases := []struct {
		name                     string
		itemCount                int
		height                   int
		columnCount              int
		highlighted              int
		selected                 int
		viewportPosition         int
		key                      string
		expectedHighlighted      int
		expectedSelected         int
		expectedViewportPosition int
		expectedConsumed         bool
	}{
		{"empty up", 0, 4, 1, -1, -1, 0, "up", -1, -1, 0, true},
		{"empty down", 0, 4, 1, -1, -1, 0, "down", -1, -1, 0, true},
		{"empty left", 0, 4, 1, -1, -1, 0, "left", -1, -1, 0, true},
		{"empty right", 0, 4, 1, -1, -1, 0, "right", -1, -1, 0, true},
		{"empty enter", 0, 4, 1, -1, -1, 0, "enter", -1, -1, 0, false},
		{"no highlight up", 4, 4, 1, -1, -1, 0, "up", 0, -1, 0, true},
		{"no highlight down", 4, 4, 1, -1, -1, 0, "down", 1, -1, 0, true},
		{"no highlight left", 4, 4, 1, -1, -1, 0, "left", 0, -1, 0, false},
		{"no highlight right", 4, 4, 1, -1, -1, 0, "right", 0, -1, 0, false},
		{"no highlight enter", 4, 4, 1, -1, -1, 0, "enter", 0, 0, 0, true},
		{"no highlight uses selection", 4, 4, 1, -1, 2, 0, "down", 3, 2, 0, true},
		{"highlight past end", 4, 4, 1, 9, -1, 0, "up", 0, -1, 0, true},
		{"highlight past end uses selection", 4, 4, 1, 9, 3, 0, "enter", 3, 3, 0, true},
		{"multiple columns no highlight up", 7, 2, 3, -1, -1, 0, "up", 0, -1, 0, true},
		{"multiple columns no highlight down", 7, 2, 3, -1, -1, 0, "down", 3, -1, 0, true},
		{"multiple columns no highlight left", 7, 2, 3, -1, -1, 0, "left", 0, -1, 0, false},
		{"multiple columns no highlight right", 7, 2, 3, -1, -1, 0, "right", 1, -1, 0, true},
		{"multiple columns down onto missing item", 7, 2, 3, 4, -1, 0, "down", 4, -1, 0, true},
		{"multiple columns down to last row", 7, 2, 3, 3, -1, 0, "down", 6, -1, 3, true},
		{"multiple columns right past last item", 7, 2, 3, 6, -1, 3, "right", 6, -1, 3, true},
		{"multiple columns up scrolls back", 7, 2, 3, 3, -1, 3, "up", 0, -1, 0, true},
		{"multiple columns misaligned viewport", 7, 2, 3, 6, -1, 2, "left", 6, -1, 3, false},
		{"shorter than height down", 2, 5, 1, 1, -1, 0, "down", 1, -1, 0, true},
		{"shorter than height viewport past end", 2, 5, 1, 1, -1, 4, "up", 0, -1, 0, true},
		{"negative viewport", 6, 2, 1, 0, -1, -3, "down", 1, -1, 0, true},
		{"zero columns left", 4, 4, 0, 1, -1, 0, "left", 1, -1, 0, false},
		{"zero columns right", 4, 4, 0, 1, -1, 0, "right", 1, -1, 0, false},
	}
	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			layerAlias, styleEntry := setupInputTest(test)
			Selector.Add(layerAlias, "list", styleEntry, getTestSelectionEntry(testCase.itemCount), 2, 2, testCase.height, 6, testCase.columnCount, 0, 0, false, true)
			startInputSimulation(test)
			selectorEntry := Selectors.Get(layerAlias, "list")
			selectorEntry.ItemHighlighted = testCase.highlighted
			selectorEntry.ItemSelected = testCase.selected
			selectorEntry.ViewportPosition = testCase.viewportPosition
			_, isConsumed := Selector.updateKeyboardEventForSelector(layerAlias, "list", []rune(testCase.key))
			UpdateDisplay(false)
			assertSelectorIndexes(test, testCase.name, selectorEntry, testCase.expectedHighlighted, testCase.expectedSelected, testCase.expectedViewportPosition)
			if isConsumed != testCase.expectedConsumed {
				test.Fatalf("%s: expected consumed %t, got %t", testCase.name, testCase.expectedConsumed, isConsumed)
			}
			if scrollBarEntry, isFound := ScrollBars.Lookup(layerAlias, selectorEntry.ScrollbarAlias); isFound && scrollBarEntry.IsEnabled &&
				scrollBarEntry.ScrollValue != selectorEntry.ViewportPosition {
				test.Fatalf("%s: expected scroll value %d to follow the viewport, got %d", testCase.name, selectorEntry.ViewportPosition, scrollBarEntry.ScrollValue)
			}
		})
	}
}

/*
TestSelectorDrawWithOutOfRangeState is a test which allows you to verify that drawing a selector cannot panic whatever
indexes it holds, since rendering runs on the event goroutine where a panic cannot be recovered by the application.

Example:

	Expected Inputs:
		A selector with four items whose viewport position is set to minus five and whose highlighted item is set to
		nine through the exported entry, followed by a screen update.

	Expected Outputs:
		The screen update completes without panicking and the selector draws its first item at its top row.
*/
func TestSelectorDrawWithOutOfRangeState(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	Selector.Add(layerAlias, "list", styleEntry, getTestSelectionEntry(4), 2, 2, 4, 10, 1, 0, 0, false, true)
	selectorEntry := GetSelector(layerAlias, "list")
	selectorEntry.ViewportPosition = -5
	selectorEntry.ItemHighlighted = 9
	UpdateDisplay(false)
	characterEntry := getCellInformationUnderMouseCursor(2, 2)
	if characterEntry.AttributeEntry.CellType != constants.CellTypeSelectorItem || characterEntry.AttributeEntry.CellControlId != 0 {
		test.Fatalf("expected item 0 at the top row, got cell type %d and item %d", characterEntry.AttributeEntry.CellType, characterEntry.AttributeEntry.CellControlId)
	}
}

/*
TestSelectorAddCorrectsOutOfRangeArguments is a test which allows you to verify that adding a selector with a viewport
position or highlighted item that does not fit its items stores corrected values, where the negative viewport used to
make the first screen update panic.

Example:

	Expected Inputs:
		Add calls with four items and a four row tray given viewport minus three and highlight nine, viewport
		seven and highlight two, and no items with highlight zero.

	Expected Outputs:
		Viewport 0 and highlight minus one; viewport 0 and highlight 2; viewport 0, highlight minus one and selection
		minus one. Each screen update completes without panicking.
*/
func TestSelectorAddCorrectsOutOfRangeArguments(test *testing.T) {
	testCases := []struct {
		name                string
		itemCount           int
		viewportPosition    int
		highlighted         int
		expectedHighlighted int
		expectedSelected    int
	}{
		{"negative viewport and highlight past end", 4, -3, 9, constants.NullItemSelection, 0},
		{"viewport past end", 4, 7, 2, 2, 0},
		{"empty list", 0, 0, 0, constants.NullItemSelection, constants.SELECTED_NONE},
	}
	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			layerAlias, styleEntry := setupInputTest(test)
			Selector.Add(layerAlias, "list", styleEntry, getTestSelectionEntry(testCase.itemCount), 2, 2, 4, 10, 1, testCase.viewportPosition, testCase.highlighted, false, true)
			UpdateDisplay(false)
			assertSelectorIndexes(test, testCase.name, Selectors.Get(layerAlias, "list"), testCase.expectedHighlighted, testCase.expectedSelected, 0)
		})
	}
}

/*
TestSelectorDeleteItemKeepsIndexesInRange is a test which allows you to verify that DeleteItem leaves no index outside
the shortened list and never writes into the slices the application passed in.

Example:

	Expected Inputs:
		A three row selector filled through SetSelectionEntry with ten items, scrolled to viewport 7 with item 9
		highlighted and selected, then items deleted from the front until one remains, then that item deleted.

	Expected Outputs:
		The application's own alias slice still reads "item0" to "item9". After each deletion the viewport is at
		most the item count minus three, and never below zero, and the highlight and selection are inside the list.
		After the last deletion highlight and selection are minus one and the viewport is 0. The scroll bar's
		maximum always equals the largest viewport position.
*/
func TestSelectorDeleteItemKeepsIndexesInRange(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := Selector.Add(layerAlias, "list", styleEntry, NewSelectionEntry(), 2, 2, 3, 10, 1, 0, 0, false, true)
	applicationEntry := getTestSelectionEntry(10)
	selector.SetSelectionEntry(applicationEntry)
	selector.setViewport(7)
	selectorEntry := Selectors.Get(layerAlias, "list")
	selectorEntry.ItemHighlighted = 9
	selectorEntry.ItemSelected = 9
	for remainingItems := 9; remainingItems >= 0; remainingItems-- {
		selector.DeleteItem(0)
		UpdateDisplay(false)
		maxViewportPosition := remainingItems - 3
		if maxViewportPosition < 0 {
			maxViewportPosition = 0
		}
		if selectorEntry.ViewportPosition < 0 || selectorEntry.ViewportPosition > maxViewportPosition {
			test.Fatalf("with %d items: expected viewport between 0 and %d, got %d", remainingItems, maxViewportPosition, selectorEntry.ViewportPosition)
		}
		if selectorEntry.ItemHighlighted >= remainingItems || selectorEntry.ItemSelected >= remainingItems {
			test.Fatalf("with %d items: highlight %d or selection %d is out of range", remainingItems, selectorEntry.ItemHighlighted, selectorEntry.ItemSelected)
		}
		scrollBarEntry := ScrollBars.Get(layerAlias, selectorEntry.ScrollbarAlias)
		if scrollBarEntry.MaxScrollValue != maxViewportPosition {
			test.Fatalf("with %d items: expected scroll bar maximum %d, got %d", remainingItems, maxViewportPosition, scrollBarEntry.MaxScrollValue)
		}
	}
	assertSelectorIndexes(test, "after deleting every item", selectorEntry, constants.NullItemSelection, constants.SELECTED_NONE, 0)
	for itemIndex, alias := range applicationEntry.SelectionAlias {
		if alias != fmt.Sprintf("item%d", itemIndex) {
			test.Fatalf("expected the application's slice to be unchanged, but index %d now holds %q", itemIndex, alias)
		}
	}
}

/*
TestSelectorScrollbarRangeMatchesItems is a test which allows you to verify that the scroll bar range set by
SetSelectionEntry and AddItem is the largest viewport position, where it used to be one more, so that scrolling to
the end showed a blank last row.

Example:

	Expected Inputs:
		A three row selector given ten items with SetSelectionEntry, then one more with AddItem, then an empty list
		with SetSelectionEntry.

	Expected Outputs:
		Scroll bar maximum 7 and enabled, then 8 and enabled, then 0 and disabled.
*/
func TestSelectorScrollbarRangeMatchesItems(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := Selector.Add(layerAlias, "list", styleEntry, NewSelectionEntry(), 2, 2, 3, 10, 1, 0, 0, false, true)
	scrollBarEntry := ScrollBars.Get(layerAlias, Selectors.Get(layerAlias, "list").ScrollbarAlias)
	selector.SetSelectionEntry(getTestSelectionEntry(10))
	if scrollBarEntry.MaxScrollValue != 7 || !scrollBarEntry.IsEnabled {
		test.Fatalf("after ten items: expected maximum 7 and enabled, got %d and %t", scrollBarEntry.MaxScrollValue, scrollBarEntry.IsEnabled)
	}
	selector.AddItem("extra", "Extra")
	if scrollBarEntry.MaxScrollValue != 8 || !scrollBarEntry.IsEnabled {
		test.Fatalf("after AddItem: expected maximum 8 and enabled, got %d and %t", scrollBarEntry.MaxScrollValue, scrollBarEntry.IsEnabled)
	}
	selector.SetSelectionEntry(NewSelectionEntry())
	if scrollBarEntry.MaxScrollValue != 0 || scrollBarEntry.IsEnabled {
		test.Fatalf("after an empty list: expected maximum 0 and disabled, got %d and %t", scrollBarEntry.MaxScrollValue, scrollBarEntry.IsEnabled)
	}
}

/*
TestSelectorClickOnStaleItemIsIgnored is a test which allows you to verify that a click on a selector cell drawn
before the list was shortened, and not yet redrawn, does not store the vanished item's index as the selection.

Example:

	Expected Inputs:
		A selector with four items is drawn, refilled with one item without a redraw, and then clicked on the row
		where item 3 was drawn.

	Expected Outputs:
		The selection stays minus one, the highlight stays inside the one item list, and no new selection is
		reported.
*/
func TestSelectorClickOnStaleItemIsIgnored(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := addTestSelector(layerAlias, styleEntry)
	UpdateDisplay(false)
	selector.SetSelectionEntry(getTestSelectionEntry(1))
	SetMouseStatus(3, 5, 0, "")
	SetMouseStatus(3, 5, 1, "")
	Selector.updateMouseEvent()
	selectorEntry := Selectors.Get(layerAlias, "selector")
	if selectorEntry.ItemSelected != constants.SELECTED_NONE || selectorEntry.ItemHighlighted > 0 || selector.IsNewItemSelected() {
		test.Fatalf("expected the stale click to be ignored, got highlighted %d, selected %d, new selection %t",
			selectorEntry.ItemHighlighted, selectorEntry.ItemSelected, selector.IsNewItemSelected())
	}
}

/*
TestSelectorFocusFixesStaleHighlight is a test which allows you to verify that giving focus to a selector whose
highlight lies past its items replaces the highlight with a valid one, and that GetAllItems returns copies.

Example:

	Expected Inputs:
		A selector with three items, two selected, whose highlight is set to ten through the exported entry, then
		focused with SetFocus. The aliases returned by GetAllItems are then overwritten.

	Expected Outputs:
		The highlight becomes 2, the selected item. The selector's own first alias is still "item0".
*/
func TestSelectorFocusFixesStaleHighlight(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := Selector.Add(layerAlias, "list", styleEntry, getTestSelectionEntry(3), 2, 2, 4, 10, 1, 0, 0, false, true)
	selector.Select("item2")
	selectorEntry := GetSelector(layerAlias, "list")
	selectorEntry.ItemHighlighted = 10
	if err := SetFocus(&selector); err != nil {
		test.Fatalf("expected SetFocus to succeed, got %v", err)
	}
	if selectorEntry.ItemHighlighted != 2 {
		test.Fatalf("expected highlight 2 after focus, got %d", selectorEntry.ItemHighlighted)
	}
	aliases, _ := selector.GetAllItems()
	aliases[0] = "changed"
	if selectorEntry.SelectionEntry.SelectionAlias[0] != "item0" {
		test.Fatalf("expected GetAllItems to return a copy, but the selector now holds %q", selectorEntry.SelectionEntry.SelectionAlias[0])
	}
}
