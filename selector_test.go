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
	selectorEntry := getSelector(layer1.layerAlias, selectorFieldInstance.controlAlias)

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
	selectorEntry := getSelector(layer1.layerAlias, selectorFieldInstance.controlAlias)

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
	selectorEntry := getSelector(layer1.layerAlias, selectorFieldInstance.controlAlias)

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
	selectorEntry := getSelector(layer1.layerAlias, selectorFieldInstance.controlAlias)

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
	selectorEntry := getSelector(layer1.layerAlias, selectorFieldInstance.controlAlias)

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
moveMouseTo is a method which allows you to inject a single mouse event at a location, with the given buttons held,
into the simulation screen and process it through the event queue.

Example:

	moveMouseTo(simScreen, 4, 3, tcell.Button1)
*/
func moveMouseTo(simScreen tcell.SimulationScreen, xLocation int, yLocation int, buttons tcell.ButtonMask) {
	simScreen.InjectMouse(xLocation, yLocation, buttons, tcell.ModNone)
	UpdateEventQueues()
}

/*
assertNoNewSelection is a method which allows you to fail the test if the selector reports a new selection.

Example:

	assertNoNewSelection(test, selector, "after a held move")
*/
func assertNoNewSelection(test *testing.T, selector SelectorInstanceType, context string) {
	test.Helper()
	if selector.IsNewItemSelected() {
		_, index := selector.GetSelected()
		test.Fatalf("%s: expected no new selection, got item %d with source %d", context, index, selector.GetSelectionSource())
	}
}

/*
TestSelectorHeldDoubleClickIsReportedOnce is a test which allows you to verify that holding the button after a double
click and moving within the same item does not report the double click again, and that the double click still uses up
the pair so a fresh click afterward is a single click.

Example:

	Expected Inputs:
		Selector with items [a, b, c, d]. Click item 1, press item 1 again and read the selection, move with the button
		held to another column of item 1, release there, then click item 1 again.

	Expected Outputs:
		Item 1 with SingleClick, then item 1 with DoubleClick on the second press, no new selection after the held move
		or the release, then item 1 with SingleClick for the fresh click.
*/
func TestSelectorHeldDoubleClickIsReportedOnce(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := addTestSelector(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)

	clickAt(simScreen, 4, 3)
	assertSelection(test, selector, "first click", 1, constants.SelectionSourceSingleClick)

	moveMouseTo(simScreen, 4, 3, tcell.Button1)
	assertSelection(test, selector, "second press", 1, constants.SelectionSourceDoubleClick)

	moveMouseTo(simScreen, 6, 3, tcell.Button1)
	assertNoNewSelection(test, selector, "held move within the item")

	moveMouseTo(simScreen, 6, 3, tcell.ButtonNone)
	assertNoNewSelection(test, selector, "release")

	clickAt(simScreen, 4, 3)
	assertSelection(test, selector, "fresh click", 1, constants.SelectionSourceSingleClick)
}

/*
TestSelectorHeldSingleClickIsReportedOnce is a test which allows you to verify that a single click is reported once
on the press, and that further moves within the same item while the button is held, and the release, report nothing.

Example:

	Expected Inputs:
		Selector with items [a, b, c, d]. Press item 2, move with the button held to two other columns of item 2,
		release, then click item 3.

	Expected Outputs:
		Item 2 with SingleClick on the press, no new selection after either held move or the release, then item 3 with
		SingleClick for the fresh click.
*/
func TestSelectorHeldSingleClickIsReportedOnce(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := addTestSelector(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)

	moveMouseTo(simScreen, 4, 4, tcell.Button1)
	assertSelection(test, selector, "press", 2, constants.SelectionSourceSingleClick)

	moveMouseTo(simScreen, 5, 4, tcell.Button1)
	assertNoNewSelection(test, selector, "first held move within the item")

	moveMouseTo(simScreen, 7, 4, tcell.Button1)
	assertNoNewSelection(test, selector, "second held move within the item")

	moveMouseTo(simScreen, 7, 4, tcell.ButtonNone)
	assertNoNewSelection(test, selector, "release")

	clickAt(simScreen, 4, 5)
	assertSelection(test, selector, "fresh click", 3, constants.SelectionSourceSingleClick)
}

/*
TestSelectorHeldDragReportsEachNewItemOnce is a test which allows you to verify that a held drag inside the selector
it began on reports each item it moves onto exactly once, as a single click, and that a drag back onto the item first
pressed is reported again since it is a different item from the one last reported.

Example:

	Expected Inputs:
		Selector with items [a, b, c, d]. Press item 1, drag with the button held to item 2, move within item 2, drag
		back to item 1, release, then click item 1.

	Expected Outputs:
		Item 1 with SingleClick on the press, item 2 with SingleClick after the drag, no new selection after the move
		within item 2, item 1 with SingleClick after the drag back, no new selection after the release, and item 1 with
		SingleClick for the fresh click, since a drag never sets up a double click.
*/
func TestSelectorHeldDragReportsEachNewItemOnce(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := addTestSelector(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)

	moveMouseTo(simScreen, 4, 3, tcell.Button1)
	assertSelection(test, selector, "press", 1, constants.SelectionSourceSingleClick)

	moveMouseTo(simScreen, 4, 4, tcell.Button1)
	assertSelection(test, selector, "drag to item 2", 2, constants.SelectionSourceSingleClick)

	moveMouseTo(simScreen, 6, 4, tcell.Button1)
	assertNoNewSelection(test, selector, "held move within item 2")

	moveMouseTo(simScreen, 4, 3, tcell.Button1)
	assertSelection(test, selector, "drag back to item 1", 1, constants.SelectionSourceSingleClick)

	moveMouseTo(simScreen, 4, 3, tcell.ButtonNone)
	assertNoNewSelection(test, selector, "release")

	clickAt(simScreen, 4, 3)
	assertSelection(test, selector, "fresh click", 1, constants.SelectionSourceSingleClick)
}

/*
TestSelectorPressOutsideDoesNotSelect is a test which allows you to verify that a press which began outside the
selector, on an empty cell or on a button, never selects or highlights a selector item when it is dragged onto the
selector with the button held, and that a fresh click after the release still selects normally.

Example:

	Expected Inputs:
		Selector with items [a, b, c, d] highlighting item 0, and a button at (20, 2). For each case, a press on the
		empty cell (30, 5) or on the button at (22, 3), held moves onto item 1 and then item 2, a release over item 2,
		then a click on item 2.

	Expected Outputs:
		After the held moves, no new selection, item 0 still highlighted and item 0 still selected. After the release,
		no new selection and item 0 still selected, since the button is up and only hover highlighting applies. Item 2
		with SingleClick for the fresh click.
*/
func TestSelectorPressOutsideDoesNotSelect(test *testing.T) {
	testCases := []struct {
		name      string
		xLocation int
		yLocation int
	}{
		{name: "empty cell", xLocation: 30, yLocation: 5},
		{name: "button", xLocation: 22, yLocation: 3},
	}
	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			layerAlias, styleEntry := setupInputTest(test)
			selector := addTestSelector(layerAlias, styleEntry)
			Button.Add(layerAlias, "button", "OK", styleEntry, 20, 2, 8, 3, true)
			simScreen := startInputSimulation(test)
			selectorEntry := selectorRegistry.Get(layerAlias, "selector")
			highlightedBefore := selectorEntry.ItemHighlighted
			selectedBefore := selectorEntry.ItemSelected

			moveMouseTo(simScreen, testCase.xLocation, testCase.yLocation, tcell.Button1)
			moveMouseTo(simScreen, 4, 3, tcell.Button1)
			assertNoNewSelection(test, selector, "held move onto item 1")
			moveMouseTo(simScreen, 4, 4, tcell.Button1)
			assertNoNewSelection(test, selector, "held move onto item 2")
			if selectorEntry.ItemHighlighted != highlightedBefore || selectorEntry.ItemSelected != selectedBefore {
				test.Fatalf("expected highlighted %d and selected %d to be unchanged, got highlighted %d and selected %d",
					highlightedBefore, selectedBefore, selectorEntry.ItemHighlighted, selectorEntry.ItemSelected)
			}
			moveMouseTo(simScreen, 4, 4, tcell.ButtonNone)
			assertNoNewSelection(test, selector, "release over item 2")
			if selectorEntry.ItemSelected != selectedBefore {
				test.Fatalf("expected selected %d to be unchanged by the release, got %d", selectedBefore, selectorEntry.ItemSelected)
			}

			clickAt(simScreen, 4, 4)
			assertSelection(test, selector, "fresh click", 2, constants.SelectionSourceSingleClick)
		})
	}
}

/*
TestSelectorAddedUnderHeldPressDoesNotSelect is a test which allows you to verify that a press held while a selector is
added under the cursor does not select an item when the mouse moves over it, and that a fresh click afterward does.

Example:

	Expected Inputs:
		A press held over the empty cell (4, 3). Selector with items [a, b, c, d] then added and drawn so item 1 lies
		under the cursor, followed by held moves to (5, 3) and onto item 2, a release, then a click on item 1.

	Expected Outputs:
		No new selection after the held moves or the release. Item 1 with SingleClick for the fresh click.
*/
func TestSelectorAddedUnderHeldPressDoesNotSelect(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	simScreen := startInputSimulation(test)

	moveMouseTo(simScreen, 4, 3, tcell.Button1)
	selector := addTestSelector(layerAlias, styleEntry)
	UpdateDisplay(false)

	moveMouseTo(simScreen, 5, 3, tcell.Button1)
	assertNoNewSelection(test, selector, "held move over item 1")
	moveMouseTo(simScreen, 4, 4, tcell.Button1)
	assertNoNewSelection(test, selector, "held move onto item 2")
	moveMouseTo(simScreen, 4, 4, tcell.ButtonNone)
	assertNoNewSelection(test, selector, "release")

	clickAt(simScreen, 4, 3)
	assertSelection(test, selector, "fresh click", 1, constants.SelectionSourceSingleClick)
}

/*
TestSelectorReplacedDuringHeldPressDoesNotSelect is a test which allows you to verify that a press which began on a
selector does not carry over to a new selector added under the same alias after the first was deleted, so a held drag
over the new selector selects nothing, and that a fresh click afterward does.

Example:

	Expected Inputs:
		Selector "selector" with items [a, b, c, d]. Press item 1, delete the selector, add a new one with the same
		alias and items and draw it, move with the button held onto item 2, release, then click item 2.

	Expected Outputs:
		Item 1 with SingleClick on the press. No new selection on the new selector after the held move or the release.
		Item 2 with SingleClick for the fresh click.
*/
func TestSelectorReplacedDuringHeldPressDoesNotSelect(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	originalSelector := addTestSelector(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)

	moveMouseTo(simScreen, 4, 3, tcell.Button1)
	assertSelection(test, originalSelector, "press", 1, constants.SelectionSourceSingleClick)

	originalSelector.Delete()
	replacementSelector := addTestSelector(layerAlias, styleEntry)
	UpdateDisplay(false)

	moveMouseTo(simScreen, 4, 4, tcell.Button1)
	assertNoNewSelection(test, replacementSelector, "held move onto item 2")
	moveMouseTo(simScreen, 4, 4, tcell.ButtonNone)
	assertNoNewSelection(test, replacementSelector, "release")

	clickAt(simScreen, 4, 4)
	assertSelection(test, replacementSelector, "fresh click", 2, constants.SelectionSourceSingleClick)
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

	signedEntry := selectorRegistry.Get(layerAlias, "signed")
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
	selectorEntry := selectorRegistry.Get(layerAlias, "list")
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
			selectorEntry := selectorRegistry.Get(layerAlias, "list")
			selectorEntry.ItemHighlighted = testCase.highlighted
			selectorEntry.ItemSelected = testCase.selected
			selectorEntry.ViewportPosition = testCase.viewportPosition
			_, isConsumed := Selector.updateKeyboardEventForSelector(layerAlias, "list", []rune(testCase.key))
			UpdateDisplay(false)
			assertSelectorIndexes(test, testCase.name, selectorEntry, testCase.expectedHighlighted, testCase.expectedSelected, testCase.expectedViewportPosition)
			if isConsumed != testCase.expectedConsumed {
				test.Fatalf("%s: expected consumed %t, got %t", testCase.name, testCase.expectedConsumed, isConsumed)
			}
			if scrollBarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, selectorEntry.ScrollbarAlias); isFound && scrollBarEntry.IsEnabled &&
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
	selectorEntry := getSelector(layerAlias, "list")
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
			assertSelectorIndexes(test, testCase.name, selectorRegistry.Get(layerAlias, "list"), testCase.expectedHighlighted, testCase.expectedSelected, 0)
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
	selector.SetViewport(7)
	selectorEntry := selectorRegistry.Get(layerAlias, "list")
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
		scrollBarEntry := scrollbarRegistry.Get(layerAlias, selectorEntry.ScrollbarAlias)
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
	scrollBarEntry := scrollbarRegistry.Get(layerAlias, selectorRegistry.Get(layerAlias, "list").ScrollbarAlias)
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
	selectorEntry := selectorRegistry.Get(layerAlias, "selector")
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
	selectorEntry := getSelector(layerAlias, "list")
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

/*
TestSelectorWheelScrollsViewport is a test which allows you to verify that the mouse wheel scrolls a selector one row
per notch within the range it can scroll to, that the scroll bar follows, that a selector which only highlights on a
click keeps its highlight, and that the wheel does nothing over a selector whose items all fit.

Example:

	Expected Inputs:
		Selector "list" at (2, 2) with ten items, four visible rows, a border, and highlighting on click only, so its
		scroll bar is at column 13. Wheel down seven times over row 3, wheel up once over the scroll bar, then wheel
		down over selector "short" at (20, 2) with two items and four visible rows.

	Expected Outputs:
		Viewport 6, the furthest it can scroll, after the seven notches, then 5 after wheel up, with the scroll bar
		value equal to the viewport. The highlight and the selection both stay 0. "short" keeps viewport 0.
*/
func TestSelectorWheelScrollsViewport(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	Selector.Add(layerAlias, "list", styleEntry, getTestSelectionEntry(10), 2, 2, 4, 10, 1, 0, 0, true, true)
	Selector.Add(layerAlias, "short", styleEntry, getTestSelectionEntry(2), 20, 2, 4, 10, 1, 0, 0, false, true)
	simScreen := startInputSimulation(test)
	selectorEntry := selectorRegistry.Get(layerAlias, "list")
	scrollBarEntry := scrollbarRegistry.Get(layerAlias, selectorEntry.ScrollbarAlias)

	for notch := 0; notch < 7; notch++ {
		injectWheel(simScreen, 4, 3, tcell.WheelDown)
	}
	assertSelectorIndexes(test, "after seven notches down", selectorEntry, 0, 0, 6)
	if scrollBarEntry.ScrollValue != 6 {
		test.Fatalf("after seven notches down: expected scroll bar value 6, got %d", scrollBarEntry.ScrollValue)
	}

	injectWheel(simScreen, 13, 3, tcell.WheelUp)
	assertSelectorIndexes(test, "after wheel up over the scroll bar", selectorEntry, 0, 0, 5)
	if scrollBarEntry.ScrollValue != 5 {
		test.Fatalf("after wheel up over the scroll bar: expected scroll bar value 5, got %d", scrollBarEntry.ScrollValue)
	}

	injectWheel(simScreen, 22, 3, tcell.WheelDown)
	if shortEntry := selectorRegistry.Get(layerAlias, "short"); shortEntry.ViewportPosition != 0 {
		test.Fatalf("expected the short selector to keep viewport 0, got %d", shortEntry.ViewportPosition)
	}
}

/*
getSelectorBorderSpan is a method which allows you to find the first and last screen columns on a row whose cells are
tagged as the border of the given selector. If no cell on the row belongs to that border, minus one is returned for
both columns.

Example:

	firstColumn, lastColumn := getSelectorBorderSpan(1, "selector")
*/
func getSelectorBorderSpan(yLocation int, selectorAlias string) (int, int) {
	firstColumn, lastColumn := -1, -1
	for xLocation, characterEntry := range commonResource.screenLayer.CharacterMemory[yLocation] {
		attributeEntry := characterEntry.AttributeEntry
		if attributeEntry.CellType != constants.CellTypeSelectorBorder || attributeEntry.CellControlAlias != selectorAlias {
			continue
		}
		if firstColumn == -1 {
			firstColumn = xLocation
		}
		lastColumn = xLocation
	}
	return firstColumn, lastColumn
}

/*
assertSelectorBorderSpan is a method which allows you to fail the test unless every row of a selector's border, from
the row above its first item row to the row below its last, has its border cells running from exactly the expected
first column to the expected last column.

Example:

	assertSelectorBorderSpan(test, "selector", 1, 4, 1, 32)
*/
func assertSelectorBorderSpan(test *testing.T, selectorAlias string, topYLocation int, bottomYLocation int, expectedFirstColumn int, expectedLastColumn int) {
	test.Helper()
	for yLocation := topYLocation; yLocation <= bottomYLocation; yLocation++ {
		firstColumn, lastColumn := getSelectorBorderSpan(yLocation, selectorAlias)
		if firstColumn != expectedFirstColumn || lastColumn != expectedLastColumn {
			test.Fatalf("row %d: expected the border to span columns %d to %d, got %d to %d", yLocation,
				expectedFirstColumn, expectedLastColumn, firstColumn, lastColumn)
		}
	}
}

/*
getTestMultiColumnSelector is a method which allows you to add a bordered selector at (2, 2) with the six items "item0"
to "item5", an item width of 10, three columns and two rows, so every item is visible, the scroll bar stays disabled,
and item N is drawn on row 2+N/3 starting at column 2+10*(N%3).

Example:

	selector := getTestMultiColumnSelector(layerAlias, styleEntry)
*/
func getTestMultiColumnSelector(layerAlias string, styleEntry types.TuiStyleEntryType) SelectorInstanceType {
	return Selector.Add(layerAlias, "selector", styleEntry, getTestSelectionEntry(6), 2, 2, 2, 10, 3, 0, constants.SELECTED_NONE, false, true)
}

/*
TestSelectorMultiColumnBorderSpansAllColumns is a test which allows you to verify that a bordered selector with
several columns draws its border around every column, not just the first.

Example:

	Expected Inputs:
		Bordered selector at (2, 2) with six items, an item width of 10, three columns and two rows.

	Expected Outputs:
		On every row from 1 to 4, the selector's border cells run from column 1 to column 32, which is 32 cells, the
		item width times three columns plus one cell on each side.
*/
func TestSelectorMultiColumnBorderSpansAllColumns(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	getTestMultiColumnSelector(layerAlias, styleEntry)
	UpdateDisplay(false)
	assertSelectorBorderSpan(test, "selector", 1, 4, 1, 32)
	if firstColumn, lastColumn := getSelectorBorderSpan(1, "selector"); lastColumn-firstColumn+1 != 32 {
		test.Fatalf("expected the border to be 32 cells wide, got %d", lastColumn-firstColumn+1)
	}
}

/*
TestSelectorClickInThirdColumnSelectsItem is a test which allows you to verify that clicking an item in the third
column of a multi-column selector selects that item.

Example:

	Expected Inputs:
		Bordered selector at (2, 2) with items [item0 .. item5], an item width of 10, three columns and two rows.
		A click at (25, 3), inside the third column of the second row.

	Expected Outputs:
		Item 5 is selected by a single click, and GetSelected returns ("item5", 5).
*/
func TestSelectorClickInThirdColumnSelectsItem(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := getTestMultiColumnSelector(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)
	clickAt(simScreen, 25, 3)
	if !selector.IsNewItemSelected() {
		test.Fatalf("expected the click to select an item")
	}
	if source := selector.GetSelectionSource(); source != constants.SelectionSourceSingleClick {
		test.Fatalf("expected selection source %d, got %d", constants.SelectionSourceSingleClick, source)
	}
	if alias, index := selector.GetSelected(); alias != "item5" || index != 5 {
		test.Fatalf("expected (%q, %d), got (%q, %d)", "item5", 5, alias, index)
	}
}

/*
TestSelectorSingleColumnBorderIsUnchanged is a test which allows you to verify that a bordered single column selector
still draws its border exactly one cell outside its one column.

Example:

	Expected Inputs:
		Bordered selector at (2, 2) with items [Item A, Item B, Item C, Item D], an item width of 10, one column and
		four rows.

	Expected Outputs:
		On every row from 1 to 6, the selector's border cells run from column 1 to column 12, which is 12 cells.
*/
func TestSelectorSingleColumnBorderIsUnchanged(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	addTestSelector(layerAlias, styleEntry)
	UpdateDisplay(false)
	assertSelectorBorderSpan(test, "selector", 1, 6, 1, 12)
}

/*
assertSelectorAttachedLayout is a method which allows you to fail the test unless a selector's scroll bar sits at the
expected column and row with the expected length, and its tooltip's hotspot and tooltip area start at the selector's
location and are as wide as its item area plus two columns.

Example:

	assertSelectorAttachedLayout(test, "after SetPosition", selectorEntry, 33, 1, 4)
*/
func assertSelectorAttachedLayout(test *testing.T, context string, layerAlias string, selectorEntry *types.SelectorEntryType, expectedXLocation int, expectedYLocation int, expectedLength int) {
	test.Helper()
	scrollBarEntry := scrollbarRegistry.Get(layerAlias, selectorEntry.ScrollbarAlias)
	if scrollBarEntry.XLocation != expectedXLocation || scrollBarEntry.YLocation != expectedYLocation || scrollBarEntry.Length != expectedLength {
		test.Fatalf("%s: expected the scroll bar at (%d, %d) with length %d, got (%d, %d) with length %d", context,
			expectedXLocation, expectedYLocation, expectedLength, scrollBarEntry.XLocation, scrollBarEntry.YLocation, scrollBarEntry.Length)
	}
	tooltipEntry := tooltipRegistry.Get(layerAlias, selectorEntry.TooltipAlias)
	expectedWidth := selectorEntry.Width + 2
	if tooltipEntry.HotspotXLocation != selectorEntry.XLocation || tooltipEntry.HotspotYLocation != selectorEntry.YLocation ||
		tooltipEntry.HotspotWidth != expectedWidth || tooltipEntry.TooltipXLocation != selectorEntry.XLocation ||
		tooltipEntry.TooltipYLocation != selectorEntry.YLocation+1 || tooltipEntry.TooltipWidth != expectedWidth {
		test.Fatalf("%s: expected the tooltip hotspot at (%d, %d) and tooltip at (%d, %d), both %d wide, got (%d, %d) %d wide and (%d, %d) %d wide",
			context, selectorEntry.XLocation, selectorEntry.YLocation, selectorEntry.XLocation, selectorEntry.YLocation+1, expectedWidth,
			tooltipEntry.HotspotXLocation, tooltipEntry.HotspotYLocation, tooltipEntry.HotspotWidth,
			tooltipEntry.TooltipXLocation, tooltipEntry.TooltipYLocation, tooltipEntry.TooltipWidth)
	}
}

/*
TestSelectorSizeIsItemArea is a test which allows you to verify that GetSize and GetBounds report a selector's item
area, its item width times its columns by its item rows, for both single and multi-column selectors.

Example:

	Expected Inputs:
		A bordered three column selector at (2, 2) with an item width of 10 and two rows, and a single column
		selector "single" at (2, 8) with an item width of 10 and four rows.

	Expected Outputs:
		GetSize returns (30, 2) and GetBounds returns (2, 2, 30, 2) for the first. GetSize returns (10, 4) for the
		second.
*/
func TestSelectorSizeIsItemArea(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := getTestMultiColumnSelector(layerAlias, styleEntry)
	if width, height := selector.GetSize(); width != 30 || height != 2 {
		test.Fatalf("expected size (30, 2), got (%d, %d)", width, height)
	}
	if xLocation, yLocation, width, height := selector.GetBounds(); xLocation != 2 || yLocation != 2 || width != 30 || height != 2 {
		test.Fatalf("expected bounds (2, 2, 30, 2), got (%d, %d, %d, %d)", xLocation, yLocation, width, height)
	}
	singleSelector := Selector.Add(layerAlias, "single", styleEntry, getTestSelectionEntry(4), 2, 8, 4, 10, 1, 0, 0, false, true)
	if width, height := singleSelector.GetSize(); width != 10 || height != 4 {
		test.Fatalf("expected size (10, 4), got (%d, %d)", width, height)
	}
}

/*
TestSelectorSetSizeResizesSelector is a test which allows you to verify that SetSize on a selector shares the width
between its columns, changes its rows, and brings its border, scroll bar, tooltip and viewport in line, and that sizes
which cannot be laid out are ignored.

Example:

	Expected Inputs:
		Bordered selector at (2, 2) with items [item0 .. item5], an item width of 10, three columns and two rows.
		SetSize(20, 1), then the viewport set to 3, then SetSize(2, 1), SetSize(18, 0), and finally SetSize(30, 2).

	Expected Outputs:
		After SetSize(20, 1): item width 6, size (18, 1), border columns 1 to 20 on rows 1 to 3, scroll bar enabled
		at (21, 1) with length 3 and maximum 3, tooltip 20 wide. The two invalid sizes leave the size (18, 1) and the
		viewport 3. After SetSize(30, 2): size (30, 2), viewport 0, scroll bar disabled at (33, 1) with length 4,
		and the border back at columns 1 to 32 on rows 1 to 4.
*/
func TestSelectorSetSizeResizesSelector(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := getTestMultiColumnSelector(layerAlias, styleEntry)
	selectorEntry := selectorRegistry.Get(layerAlias, "selector")
	scrollBarEntry := scrollbarRegistry.Get(layerAlias, selectorEntry.ScrollbarAlias)

	selector.SetSize(20, 1)
	if selectorEntry.ItemWidth != 6 {
		test.Fatalf("expected item width 6, got %d", selectorEntry.ItemWidth)
	}
	if width, height := selector.GetSize(); width != 18 || height != 1 {
		test.Fatalf("expected size (18, 1), got (%d, %d)", width, height)
	}
	assertSelectorAttachedLayout(test, "after SetSize(20, 1)", layerAlias, selectorEntry, 21, 1, 3)
	if !scrollBarEntry.IsEnabled || scrollBarEntry.MaxScrollValue != 3 {
		test.Fatalf("expected an enabled scroll bar with maximum 3, got %t and %d", scrollBarEntry.IsEnabled, scrollBarEntry.MaxScrollValue)
	}
	UpdateDisplay(false)
	assertSelectorBorderSpan(test, "selector", 1, 3, 1, 20)

	selector.SetViewport(3)
	selector.SetSize(2, 1)
	selector.SetSize(18, 0)
	if width, height := selector.GetSize(); width != 18 || height != 1 {
		test.Fatalf("expected invalid sizes to be ignored, got (%d, %d)", width, height)
	}
	if selectorEntry.ViewportPosition != 3 {
		test.Fatalf("expected invalid sizes to keep viewport 3, got %d", selectorEntry.ViewportPosition)
	}

	selector.SetSize(30, 2)
	if width, height := selector.GetSize(); width != 30 || height != 2 {
		test.Fatalf("expected size (30, 2), got (%d, %d)", width, height)
	}
	if selectorEntry.ViewportPosition != 0 || scrollBarEntry.IsEnabled {
		test.Fatalf("expected viewport 0 and a disabled scroll bar, got %d and %t", selectorEntry.ViewportPosition, scrollBarEntry.IsEnabled)
	}
	assertSelectorAttachedLayout(test, "after SetSize(30, 2)", layerAlias, selectorEntry, 33, 1, 4)
	UpdateDisplay(false)
	assertSelectorBorderSpan(test, "selector", 1, 4, 1, 32)
}

/*
TestSelectorSetPositionAndBorderMoveAttachedControls is a test which allows you to verify that moving a selector, or
turning its border off, moves its scroll bar and tooltip with it.

Example:

	Expected Inputs:
		Bordered single column selector at (2, 2) with ten items, an item width of 10 and three rows, so its scroll
		bar is enabled. SetPosition(5, 4), then SetBorderDrawn(false), then SetBorderDrawn(true).

	Expected Outputs:
		Scroll bar at (13, 1) with length 5 before moving, (16, 3) with length 5 after SetPosition, (14, 4) with
		length 3 without the border, and (16, 3) with length 5 with it again. The tooltip starts at the selector
		each time, and after moving the border is drawn from column 4 to 15 on rows 3 to 7.
*/
func TestSelectorSetPositionAndBorderMoveAttachedControls(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := Selector.Add(layerAlias, "selector", styleEntry, getTestSelectionEntry(10), 2, 2, 3, 10, 1, 0, constants.SELECTED_NONE, false, true)
	selectorEntry := selectorRegistry.Get(layerAlias, "selector")
	assertSelectorAttachedLayout(test, "after Add", layerAlias, selectorEntry, 13, 1, 5)

	selector.SetPosition(5, 4)
	assertSelectorAttachedLayout(test, "after SetPosition", layerAlias, selectorEntry, 16, 3, 5)
	UpdateDisplay(false)
	assertSelectorBorderSpan(test, "selector", 3, 7, 4, 15)

	selector.SetBorderDrawn(false)
	assertSelectorAttachedLayout(test, "without border", layerAlias, selectorEntry, 14, 4, 3)
	selector.SetBorderDrawn(true)
	assertSelectorAttachedLayout(test, "with border again", layerAlias, selectorEntry, 16, 3, 5)
}

/*
TestSelectorFocusSelectionWithNoColumns is a test which allows you to verify that FocusSelection centres the item on a
selector added with a column count of zero, which is laid out as a single column.

Example:

	Expected Inputs:
		Selector with items [item0 .. item9], three rows and a column count of 0. FocusSelection("item5").

	Expected Outputs:
		The viewport is 4, so item 5 is in the middle row, and the scroll bar value is 4.
*/
func TestSelectorFocusSelectionWithNoColumns(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := Selector.Add(layerAlias, "selector", styleEntry, getTestSelectionEntry(10), 2, 2, 3, 10, 0, 0, constants.SELECTED_NONE, false, true)
	selectorEntry := selectorRegistry.Get(layerAlias, "selector")
	selector.FocusSelection("item5")
	if selectorEntry.ViewportPosition != 4 {
		test.Fatalf("expected viewport 4, got %d", selectorEntry.ViewportPosition)
	}
	if scrollValue := scrollbarRegistry.Get(layerAlias, selectorEntry.ScrollbarAlias).ScrollValue; scrollValue != 4 {
		test.Fatalf("expected scroll bar value 4, got %d", scrollValue)
	}
}

/*
addViewportTestSelector is a method which allows you to add a bordered selector "list" at (2, 2) with the items
"item0" to "itemN-1", an item width of 6, the given number of rows and columns, nothing highlighted or selected, and
highlighting on hover, so item N is drawn on row 2+(N-viewport)/columns starting at column 2+6*(N%columns).

Example:

	selector := addViewportTestSelector(layerAlias, styleEntry, 10, 3, 1)
*/
func addViewportTestSelector(layerAlias string, styleEntry types.TuiStyleEntryType, itemCount int, rowCount int, columnCount int) SelectorInstanceType {
	return Selector.Add(layerAlias, "list", styleEntry, getTestSelectionEntry(itemCount), 2, 2, rowCount, 6, columnCount, 0, constants.SELECTED_NONE, false, true)
}

/*
assertSelectorViewport is a method which allows you to fail the test unless GetViewport returns the expected viewport
position and the selector's scroll bar value matches it.

Example:

	assertSelectorViewport(test, "after Down", layerAlias, selector, 2)
*/
func assertSelectorViewport(test *testing.T, context string, layerAlias string, selector SelectorInstanceType, expectedViewportPosition int) {
	test.Helper()
	if viewportPosition := selector.GetViewport(); viewportPosition != expectedViewportPosition {
		test.Fatalf("%s: expected viewport %d, got %d", context, expectedViewportPosition, viewportPosition)
	}
	selectorEntry := selectorRegistry.Get(layerAlias, selector.GetAlias())
	if scrollValue := scrollbarRegistry.Get(layerAlias, selectorEntry.ScrollbarAlias).ScrollValue; scrollValue != expectedViewportPosition {
		test.Fatalf("%s: expected scroll bar value %d, got %d", context, expectedViewportPosition, scrollValue)
	}
}

/*
TestSelectorGetViewportFollowsNavigation is a test which allows you to verify that GetViewport reports the viewport
after it is moved by the arrow keys, the mouse wheel and FocusSelection, on single and multi-column selectors.

Example:

	Expected Inputs:
		Single column: ten items and three rows, focused, then Down four times, wheel down once over (3, 3), and
		FocusSelection("item9"). Multi-column: twelve items, three columns and two rows, focused, then Down twice,
		wheel down once over (3, 3), and FocusSelection("item1").

	Expected Outputs:
		Single column: viewport 0 on focus, 2 after the keys, 3 after the wheel, and 7 after FocusSelection, which
		is the furthest it can scroll. Multi-column: viewport 0 on focus, 3 after the keys, 6 after the wheel, and 0
		after FocusSelection. The scroll bar value matches the viewport each time.
*/
func TestSelectorGetViewportFollowsNavigation(test *testing.T) {
	test.Run("single column", func(test *testing.T) {
		layerAlias, styleEntry := setupInputTest(test)
		selector := addViewportTestSelector(layerAlias, styleEntry, 10, 3, 1)
		simScreen := startInputSimulation(test)
		if err := SetFocus(&selector); err != nil {
			test.Fatalf("expected the selector to take focus, got %v", err)
		}
		assertSelectorViewport(test, "on focus", layerAlias, selector, 0)
		for range 4 {
			pressKey(simScreen, tcell.KeyDown, 0, tcell.ModNone)
		}
		assertSelectorViewport(test, "after Down four times", layerAlias, selector, 2)
		injectWheel(simScreen, 3, 3, tcell.WheelDown)
		assertSelectorViewport(test, "after wheel down", layerAlias, selector, 3)
		selector.FocusSelection("item9")
		assertSelectorViewport(test, "after FocusSelection", layerAlias, selector, 7)
	})
	test.Run("multi-column", func(test *testing.T) {
		layerAlias, styleEntry := setupInputTest(test)
		selector := addViewportTestSelector(layerAlias, styleEntry, 12, 2, 3)
		simScreen := startInputSimulation(test)
		if err := SetFocus(&selector); err != nil {
			test.Fatalf("expected the selector to take focus, got %v", err)
		}
		assertSelectorViewport(test, "on focus", layerAlias, selector, 0)
		for range 2 {
			pressKey(simScreen, tcell.KeyDown, 0, tcell.ModNone)
		}
		assertSelectorViewport(test, "after Down twice", layerAlias, selector, 3)
		injectWheel(simScreen, 3, 3, tcell.WheelDown)
		assertSelectorViewport(test, "after wheel down", layerAlias, selector, 6)
		selector.FocusSelection("item1")
		assertSelectorViewport(test, "after FocusSelection", layerAlias, selector, 0)
	})
}

/*
TestSelectorSetViewportClampsAndAlignsToRows is a test which allows you to verify that SetViewport clamps the position
to the range the selector can scroll to, moves it back to the start of its row on a multi-column selector, moves the
scroll bar to match, and leaves the highlight alone.

Example:

	Expected Inputs:
		Single column: ten items and three rows with item 1 highlighted, then SetViewport with -4, 50 and 5.
		Multi-column: twelve items, three columns and two rows, then SetViewport with 4, 100 and -1.

	Expected Outputs:
		Single column: viewport 0, 7 and 5, with item 1 still highlighted. Multi-column: viewport 3, 6 and 0. The
		scroll bar value matches the viewport each time.
*/
func TestSelectorSetViewportClampsAndAlignsToRows(test *testing.T) {
	test.Run("single column", func(test *testing.T) {
		layerAlias, styleEntry := setupInputTest(test)
		selector := addViewportTestSelector(layerAlias, styleEntry, 10, 3, 1)
		selector.SetHighlighted(1)
		for _, step := range []struct{ viewportPosition, expected int }{{-4, 0}, {50, 7}, {5, 5}} {
			selector.SetViewport(step.viewportPosition)
			assertSelectorViewport(test, fmt.Sprintf("SetViewport(%d)", step.viewportPosition), layerAlias, selector, step.expected)
		}
		if highlighted := selector.GetHighlighted(); highlighted != 1 {
			test.Fatalf("expected item 1 to stay highlighted, got %d", highlighted)
		}
	})
	test.Run("multi-column", func(test *testing.T) {
		layerAlias, styleEntry := setupInputTest(test)
		selector := addViewportTestSelector(layerAlias, styleEntry, 12, 2, 3)
		for _, step := range []struct{ viewportPosition, expected int }{{4, 3}, {100, 6}, {-1, 0}} {
			selector.SetViewport(step.viewportPosition)
			assertSelectorViewport(test, fmt.Sprintf("SetViewport(%d)", step.viewportPosition), layerAlias, selector, step.expected)
		}
	})
}

/*
TestSelectorSetHighlightedScrollsIntoView is a test which allows you to verify that SetHighlighted highlights an item,
scrolls the fewest rows needed to show it, moves the scroll bar to match, clamps an index past the end to the last
item, clears the highlight for a negative index, and never changes the selection or the new selection flag.

Example:

	Expected Inputs:
		Single column: ten items and three rows, item 2 selected with the new selection flag set, then
		SetHighlighted with 8, 1, 40 and -3. Multi-column: twelve items, three columns and two rows, then
		SetHighlighted with 10 and 4.

	Expected Outputs:
		Single column, as (highlight, viewport): (8, 6), (1, 1), (9, 7), then (-1, 7), with the scroll bar value
		matching the viewport, and item 2 still selected with the flag still set after every step. Multi-column:
		(10, 6), then (4, 3).
*/
func TestSelectorSetHighlightedScrollsIntoView(test *testing.T) {
	test.Run("single column", func(test *testing.T) {
		layerAlias, styleEntry := setupInputTest(test)
		selector := addViewportTestSelector(layerAlias, styleEntry, 10, 3, 1)
		selectorEntry := selectorRegistry.Get(layerAlias, "list")
		selectorEntry.ItemSelected = 2
		selectorEntry.IsNewItemSelected = true
		for _, step := range []struct{ itemIndex, expectedHighlighted, expectedViewport int }{
			{8, 8, 6}, {1, 1, 1}, {40, 9, 7}, {-3, constants.SELECTED_NONE, 7},
		} {
			context := fmt.Sprintf("SetHighlighted(%d)", step.itemIndex)
			selector.SetHighlighted(step.itemIndex)
			if highlighted := selector.GetHighlighted(); highlighted != step.expectedHighlighted {
				test.Fatalf("%s: expected highlight %d, got %d", context, step.expectedHighlighted, highlighted)
			}
			assertSelectorViewport(test, context, layerAlias, selector, step.expectedViewport)
			if selectorEntry.ItemSelected != 2 || !selectorEntry.IsNewItemSelected {
				test.Fatalf("%s: expected item 2 to stay selected with the flag set, got %d and %t", context,
					selectorEntry.ItemSelected, selectorEntry.IsNewItemSelected)
			}
		}
	})
	test.Run("multi-column", func(test *testing.T) {
		layerAlias, styleEntry := setupInputTest(test)
		selector := addViewportTestSelector(layerAlias, styleEntry, 12, 2, 3)
		for _, step := range []struct{ itemIndex, expectedViewport int }{{10, 6}, {4, 3}} {
			context := fmt.Sprintf("SetHighlighted(%d)", step.itemIndex)
			selector.SetHighlighted(step.itemIndex)
			if highlighted := selector.GetHighlighted(); highlighted != step.itemIndex {
				test.Fatalf("%s: expected highlight %d, got %d", context, step.itemIndex, highlighted)
			}
			assertSelectorViewport(test, context, layerAlias, selector, step.expectedViewport)
		}
	})
}

/*
TestSelectorHighlightEdgeCases is a test which allows you to verify that SetHighlighted on a selector with no items
leaves nothing highlighted, and that GetHighlighted reports an out of range highlight as none without changing it.

Example:

	Expected Inputs:
		A selector "empty" with no items given SetHighlighted(0). A selector with ten items whose stored highlight is
		set directly to 20.

	Expected Outputs:
		The empty selector reports constants.SELECTED_NONE and viewport 0. The other reports constants.SELECTED_NONE
		while its stored highlight stays 20.
*/
func TestSelectorHighlightEdgeCases(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	emptySelector := Selector.Add(layerAlias, "empty", styleEntry, NewSelectionEntry(), 2, 8, 1, 6, 1, 0, constants.SELECTED_NONE, false, true)
	emptySelector.SetHighlighted(0)
	if highlighted := emptySelector.GetHighlighted(); highlighted != constants.SELECTED_NONE {
		test.Fatalf("expected no highlight on an empty selector, got %d", highlighted)
	}
	if viewportPosition := emptySelector.GetViewport(); viewportPosition != 0 {
		test.Fatalf("expected viewport 0 on an empty selector, got %d", viewportPosition)
	}

	selector := addViewportTestSelector(layerAlias, styleEntry, 10, 3, 1)
	selectorEntry := selectorRegistry.Get(layerAlias, "list")
	selectorEntry.ItemHighlighted = 20
	if highlighted := selector.GetHighlighted(); highlighted != constants.SELECTED_NONE {
		test.Fatalf("expected an out of range highlight to be reported as none, got %d", highlighted)
	}
	if selectorEntry.ItemHighlighted != 20 {
		test.Fatalf("expected GetHighlighted to leave the stored highlight at 20, got %d", selectorEntry.ItemHighlighted)
	}
}

/*
TestSelectorViewportAndHighlightOnMissingSelector is a test which allows you to verify that GetViewport,
SetViewport, GetHighlighted and SetHighlighted do nothing, and do not panic, on a selector that was never added or
has been deleted.

Example:

	Expected Inputs:
		An instance for the alias "missing", which was never added, and an instance for "list", added and then
		deleted. Each calls SetViewport(3) and SetHighlighted(2), then GetViewport and GetHighlighted.

	Expected Outputs:
		No panic. GetViewport returns 0 and GetHighlighted returns constants.SELECTED_NONE for both.
*/
func TestSelectorViewportAndHighlightOnMissingSelector(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	deletedSelector := addViewportTestSelector(layerAlias, styleEntry, 10, 3, 1)
	Selector.Delete(layerAlias, "list")
	missingSelector := SelectorInstanceType{BaseControlInstanceType{layerAlias: layerAlias, controlAlias: "missing", controlType: constants.TYPE_SELECTOR}}
	for _, selector := range []SelectorInstanceType{missingSelector, deletedSelector} {
		selector.SetViewport(3)
		selector.SetHighlighted(2)
		if viewportPosition := selector.GetViewport(); viewportPosition != 0 {
			test.Fatalf("%q: expected viewport 0, got %d", selector.GetAlias(), viewportPosition)
		}
		if highlighted := selector.GetHighlighted(); highlighted != constants.SELECTED_NONE {
			test.Fatalf("%q: expected no highlight, got %d", selector.GetAlias(), highlighted)
		}
	}
}

/*
TestSelectorClickOnlyAndColumnMethods is a test which allows you to verify that a selector's highlight on click only
setting and its column count can be read and set through its methods, that changing the column count resizes the
selector, and that a count below one is ignored.

Example:

	Expected Inputs:
		The bordered three column selector with an item width of 10. SetHighlightOnClickOnly(true), then
		SetColumnCount(2), then SetColumnCount(0).

	Expected Outputs:
		IsHighlightOnClickOnly reads false, then true. GetColumnCount reads 3, then 2 with a width of 20, and still 2
		after the invalid count.
*/
func TestSelectorClickOnlyAndColumnMethods(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := getTestMultiColumnSelector(layerAlias, styleEntry)
	if selector.IsHighlightOnClickOnly() {
		test.Fatalf("expected highlight on click only to start off")
	}
	selector.SetHighlightOnClickOnly(true)
	if !selector.IsHighlightOnClickOnly() {
		test.Fatalf("expected highlight on click only to read true once set")
	}
	if columnCount := selector.GetColumnCount(); columnCount != 3 {
		test.Fatalf("expected 3 columns, got %d", columnCount)
	}
	selector.SetColumnCount(2)
	if width, _ := selector.GetSize(); selector.GetColumnCount() != 2 || width != 20 {
		test.Fatalf("expected 2 columns and a width of 20, got %d and %d", selector.GetColumnCount(), width)
	}
	selector.SetColumnCount(0)
	if columnCount := selector.GetColumnCount(); columnCount != 2 {
		test.Fatalf("expected an invalid column count to be ignored, got %d", columnCount)
	}
}
