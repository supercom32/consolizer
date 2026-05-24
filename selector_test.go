package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

const SELECTOR_TEST_SUITE_NAME = "selector"

func TestSelectorRandomSelection(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
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

func TestSelectorLongList(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
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

func TestGetAllItems(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
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

func TestSelectorLongListWithColors(test *testing.T) {
	textStyleAlias := "red"
	attributeEntry := NewTextStyle()
	attributeEntry.ForegroundColor = GetRGBColor(255, 0, 0)
	AddTextStyle(textStyleAlias, attributeEntry)
	layer1, _, _, styleEntry := CommonTestSetup()
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

func TestFocusSelectionInitialPosition(test *testing.T) {
	// Setup
	layer1, _, _, styleEntry := CommonTestSetup()
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

func TestFocusSelectionMiddleItem(test *testing.T) {
	// Setup
	layer1, _, _, styleEntry := CommonTestSetup()
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

func TestFocusSelectionEndItem(test *testing.T) {
	// Setup
	layer1, _, _, styleEntry := CommonTestSetup()
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

func TestFocusSelectionBeginningItem(test *testing.T) {
	// Setup
	layer1, _, _, styleEntry := CommonTestSetup()
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

func TestFocusSelectionNonExistentItem(test *testing.T) {
	// Setup
	layer1, _, _, styleEntry := CommonTestSetup()
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
