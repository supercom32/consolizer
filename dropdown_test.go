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
        A dropdown control added to a layer with three items.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) for the closed dropdown.
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
        A dropdown control initialized with the second item (index 1) pre-selected.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) and selected value is "Item 2".
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
        A focused dropdown control with its tray visibility manually set to true.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the expanded selection tray.
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
        A manually opened dropdown where the selector index is changed from 0 to 2.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) and the dropdown value is updated to "Item 3".
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
        A dropdown control added with six items and a tray height restricted to three.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing an open tray with a visible scrollbar.
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
        An open dropdown tray with six items where the viewport position is programmatically set to 3.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the scrolled list of items in the tray.
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
        A layer containing one dropdown control which is سپس deleted.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) for an empty layer.
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
        A layer containing two dropdown controls followed by a DeleteAll call.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) for an empty layer after all controls are removed.
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
        A layer containing a dropdown where focus is programmatically set to that dropdown control.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the dropdown in its focused visual state.
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
