package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
	"testing"
)

const DROPDOWN_TEST_SUITE_NAME = "dropdown"

/*
TestDropdownDefaultState is a test which allows you to verify that a dropdown control is rendered correctly with its
default state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Output: Screen content matches expected ANSI string (Base64 encoded).
*/
func TestDropdownDefaultState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.SelectionAlias = []string{"item1", "item2", "item3"}
	selectionEntry.SelectionValue = []string{"Item 1", "Item 2", "Item 3"}
	layer1.AddDropdown(styleEntry, selectionEntry, 2, 2, 3, 10, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, DROPDOWN_TEST_SUITE_NAME, "TestDropdownDefaultState", obtainedValue)
	expectedValue := LoadMasterImage(DROPDOWN_TEST_SUITE_NAME, "TestDropdownDefaultState")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDropdownWithDefaultSelection is a test which allows you to verify that a dropdown control correctly displays a pre-
selected item upon initialization.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Output: Screen content matches expected ANSI string (Base64 encoded) and selected value/alias are correct.
*/
func TestDropdownWithDefaultSelection(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.SelectionAlias = []string{"item1", "item2", "item3"}
	selectionEntry.SelectionValue = []string{"Item 1", "Item 2", "Item 3"}
	dropdownInstance := layer1.AddDropdown(styleEntry, selectionEntry, 2, 2, 3, 10, 1)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, DROPDOWN_TEST_SUITE_NAME, "TestDropdownWithDefaultSelection", obtainedValue)
	expectedValue := LoadMasterImage(DROPDOWN_TEST_SUITE_NAME, "TestDropdownWithDefaultSelection")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the selected value
	if !assert.Equal(test, "Item 2", dropdownInstance.GetValue(), "Dropdown should have 'Item 2' selected") {
		fmt.Println("Dropdown selection value is incorrect")
	}

	// Verify the selected alias
	if !assert.Equal(test, "item2", dropdownInstance.GetAlias(), "Dropdown should have 'item2' alias selected") {
		fmt.Println("Dropdown selection alias is incorrect")
	}
}

/*
TestDropdownOpenState is a test which allows you to verify that a dropdown tray and its selector are correctly displayed
when the dropdown is opened.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Output: Screen content matches expected ANSI string (Base64 encoded) with dropdown tray open.
*/
func TestDropdownOpenState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.SelectionAlias = []string{"item1", "item2", "item3"}
	selectionEntry.SelectionValue = []string{"Item 1", "Item 2", "Item 3"}
	dropdownInstance := layer1.AddDropdown(styleEntry, selectionEntry, 2, 2, 3, 10, 0)

	// Set focus to the dropdown
	setFocusedControl(layer1.layerAlias, dropdownInstance.GetAlias(), constants.CellTypeDropdown)

	// Simulate opening the dropdown
	dropdownEntry := Dropdown.Get(layer1.layerAlias, dropdownInstance.controlAlias)
	dropdownEntry.IsTrayOpen = true

	// Make selector visible
	selectorEntry := Selectors.Get(layer1.layerAlias, dropdownEntry.SelectorAlias)
	selectorEntry.IsVisible = true

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, DROPDOWN_TEST_SUITE_NAME, "TestDropdownOpenState", obtainedValue)
	expectedValue := LoadMasterImage(DROPDOWN_TEST_SUITE_NAME, "TestDropdownOpenState")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDropdownChangeSelection is a test which allows you to verify that selecting a new item from an open dropdown tray
correctly updates the dropdown's value.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Output: Screen content matches expected ANSI string (Base64 encoded) after changing selection.
*/
func TestDropdownChangeSelection(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.SelectionAlias = []string{"item1", "item2", "item3"}
	selectionEntry.SelectionValue = []string{"Item 1", "Item 2", "Item 3"}
	dropdownInstance := layer1.AddDropdown(styleEntry, selectionEntry, 2, 2, 3, 10, 0)

	// Set focus to the dropdown
	setFocusedControl(layer1.layerAlias, dropdownInstance.GetAlias(), constants.CellTypeDropdown)

	// Simulate opening the dropdown
	dropdownEntry := Dropdown.Get(layer1.layerAlias, dropdownInstance.controlAlias)
	dropdownEntry.IsTrayOpen = true

	// Make selector visible
	selectorEntry := Selectors.Get(layer1.layerAlias, dropdownEntry.SelectorAlias)
	selectorEntry.IsVisible = true

	// Change selection
	selectorEntry.ItemSelected = 2

	// Simulate closing the dropdown and applying selection
	dropdownEntry.ItemSelected = selectorEntry.ItemSelected
	dropdownEntry.IsTrayOpen = false
	selectorEntry.IsVisible = false

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, DROPDOWN_TEST_SUITE_NAME, "TestDropdownChangeSelection", obtainedValue)
	expectedValue := LoadMasterImage(DROPDOWN_TEST_SUITE_NAME, "TestDropdownChangeSelection")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the selected value
	if !assert.Equal(test, "Item 3", dropdownInstance.GetValue(), "Dropdown should have 'Item 3' selected") {
		fmt.Println("Dropdown selection value is incorrect")
	}
}

/*
TestDropdownWithManyItems is a test which allows you to verify that a dropdown correctly handles a large number of
items, including the rendering of a scrollbar when the tray is opened.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Output: Screen content matches expected ANSI string (Base64 encoded) with scrollbar visible.
*/
func TestDropdownWithManyItems(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.SelectionAlias = []string{"item1", "item2", "item3", "item4", "item5", "item6"}
	selectionEntry.SelectionValue = []string{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5", "Item 6"}
	dropdownInstance := layer1.AddDropdown(styleEntry, selectionEntry, 2, 2, 3, 10, 0)

	// Set focus to the dropdown
	setFocusedControl(layer1.layerAlias, dropdownInstance.GetAlias(), constants.CellTypeDropdown)

	// Simulate opening the dropdown
	dropdownEntry := Dropdown.Get(layer1.layerAlias, dropdownInstance.controlAlias)
	dropdownEntry.IsTrayOpen = true

	// Make selector visible
	selectorEntry := Selectors.Get(layer1.layerAlias, dropdownEntry.SelectorAlias)
	selectorEntry.IsVisible = true

	// Make scrollbar visible
	scrollBarEntry := ScrollBars.Get(layer1.layerAlias, dropdownEntry.ScrollbarAlias)
	scrollBarEntry.IsVisible = true

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, DROPDOWN_TEST_SUITE_NAME, "TestDropdownWithManyItems", obtainedValue)
	expectedValue := LoadMasterImage(DROPDOWN_TEST_SUITE_NAME, "TestDropdownWithManyItems")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDropdownScrolling is a test which allows you to verify that a dropdown tray correctly scrolls its items when a
scrollbar interaction occurs.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Output: Screen content matches expected ANSI string (Base64 encoded) after scrolling tray items.
*/
func TestDropdownScrolling(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.SelectionAlias = []string{"item1", "item2", "item3", "item4", "item5", "item6"}
	selectionEntry.SelectionValue = []string{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5", "Item 6"}
	dropdownInstance := layer1.AddDropdown(styleEntry, selectionEntry, 2, 2, 3, 10, 0)

	// Set focus to the dropdown
	setFocusedControl(layer1.layerAlias, dropdownInstance.GetAlias(), constants.CellTypeDropdown)

	// Simulate opening the dropdown
	dropdownEntry := Dropdown.Get(layer1.layerAlias, dropdownInstance.controlAlias)
	dropdownEntry.IsTrayOpen = true

	// Make selector visible
	selectorEntry := Selectors.Get(layer1.layerAlias, dropdownEntry.SelectorAlias)
	selectorEntry.IsVisible = true

	// Make scrollbar visible
	scrollBarEntry := ScrollBars.Get(layer1.layerAlias, dropdownEntry.ScrollbarAlias)
	scrollBarEntry.IsVisible = true

	// Scroll down
	scrollBarEntry.ScrollValue = 3
	selectorEntry.ViewportPosition = 3

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, DROPDOWN_TEST_SUITE_NAME, "TestDropdownScrolling", obtainedValue)
	expectedValue := LoadMasterImage(DROPDOWN_TEST_SUITE_NAME, "TestDropdownScrolling")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDropdownDelete is a test which allows you to verify that a dropdown control can be successfully deleted from its
parent layer.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Output: Screen content matches expected ANSI string (Base64 encoded) after deleting the dropdown.
*/
func TestDropdownDelete(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.SelectionAlias = []string{"item1", "item2", "item3"}
	selectionEntry.SelectionValue = []string{"Item 1", "Item 2", "Item 3"}
	dropdownInstance := layer1.AddDropdown(styleEntry, selectionEntry, 2, 2, 3, 10, 0)
	UpdateDisplay(false)

	// Delete the dropdown
	dropdownInstance.Delete()

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, DROPDOWN_TEST_SUITE_NAME, "TestDropdownDelete", obtainedValue)
	expectedValue := LoadMasterImage(DROPDOWN_TEST_SUITE_NAME, "TestDropdownDelete")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDropdownDeleteAll is a test which allows you to verify that all dropdown controls on a layer can be successfully
deleted at once.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Output: Screen content matches expected ANSI string (Base64 encoded) after deleting all dropdowns.
*/
func TestDropdownDeleteAll(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.SelectionAlias = []string{"item1", "item2", "item3"}
	selectionEntry.SelectionValue = []string{"Item 1", "Item 2", "Item 3"}
	layer1.AddDropdown(styleEntry, selectionEntry, 2, 2, 3, 10, 0)
	layer1.AddDropdown(styleEntry, selectionEntry, 2, 6, 3, 10, 1)
	UpdateDisplay(false)

	// Delete all dropdowns
	Dropdown.DeleteAll(layer1.layerAlias)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, DROPDOWN_TEST_SUITE_NAME, "TestDropdownDeleteAll", obtainedValue)
	expectedValue := LoadMasterImage(DROPDOWN_TEST_SUITE_NAME, "TestDropdownDeleteAll")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestDropdownFocus is a test which allows you to verify that a dropdown control correctly handles gaining focus and
renders accordingly.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Output: Screen content matches expected ANSI string (Base64 encoded) with dropdown in focus.
*/
func TestDropdownFocus(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.SelectionAlias = []string{"item1", "item2", "item3"}
	selectionEntry.SelectionValue = []string{"Item 1", "Item 2", "Item 3"}
	dropdownInstance := layer1.AddDropdown(styleEntry, selectionEntry, 2, 2, 3, 10, 0)

	// Set focus to the dropdown
	setFocusedControl(layer1.layerAlias, dropdownInstance.GetAlias(), constants.CellTypeDropdown)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, DROPDOWN_TEST_SUITE_NAME, "TestDropdownFocus", obtainedValue)
	expectedValue := LoadMasterImage(DROPDOWN_TEST_SUITE_NAME, "TestDropdownFocus")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}
