package consolizer

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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

/*
TestDropdownEscClosesOpenTrayWithoutChangingSelection is a test which allows you to verify that Esc on an open dropdown
closes the tray, keeps the committed selection, and is consumed, and that the item highlighted before Esc is not
committed by any later close. In addition, the following should be noted:

  - This pins existing behavior rather than reproducing a failure. Before the fix, Esc shared Enter's branch, but
    that branch committed the tray selector's ItemSelected, which arrow keys never change, so an open tray already
    closed without a selection change.

Example:

	Expected Inputs:
		Dropdown with items [Zero, One, Two, Three] and item 0 selected, focused. Keys Enter, Down, Esc, then Enter,
		Enter to reopen and close it again.

	Expected Outputs:
		After Esc the tray is closed, ItemSelected is 0 on both the dropdown and its tray selector, and the keyboard
		buffer is empty. After reopening, item 0 is highlighted, and after closing ItemSelected is still 0.
*/
func TestDropdownEscClosesOpenTrayWithoutChangingSelection(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	dropdownEntry := addTestDropdown(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)

	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	if !dropdownEntry.IsTrayOpen {
		test.Fatalf("expected Enter to open the tray")
	}
	pressKey(simScreen, tcell.KeyDown, 0, tcell.ModNone)
	pressKey(simScreen, tcell.KeyEscape, 0, tcell.ModNone)

	selectorEntry := Selectors.Get(layerAlias, dropdownEntry.SelectorAlias)
	if dropdownEntry.IsTrayOpen {
		test.Fatalf("expected Esc to close the tray")
	}
	if dropdownEntry.ItemSelected != 0 || selectorEntry.ItemSelected != 0 {
		test.Fatalf("expected selection to stay 0 after Esc, got dropdown %d and tray %d", dropdownEntry.ItemSelected, selectorEntry.ItemSelected)
	}
	assertKeyboardBuffer(test, "after Esc on an open tray")

	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	if selectorEntry.ItemHighlighted != 0 {
		test.Fatalf("expected reopened tray to highlight item 0, got %d", selectorEntry.ItemHighlighted)
	}
	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	if dropdownEntry.IsTrayOpen || dropdownEntry.ItemSelected != 0 {
		test.Fatalf("expected closed tray with item 0 selected, got open %v with item %d", dropdownEntry.IsTrayOpen, dropdownEntry.ItemSelected)
	}
}

/*
TestDropdownEscOnClosedTrayIsNotConsumed is a test which allows you to verify that Esc on a focused dropdown whose tray
is closed neither opens the tray nor changes the selection, and reaches the keyboard buffer so the application can
treat it as Back. In addition, the following should be noted:

  - Before the fix, Esc on a closed dropdown opened its tray and was consumed.

Example:

	Expected Inputs:
		Dropdown with items [Zero, One, Two, Three] and item 0 selected, focused with its tray closed. Key Esc.

	Expected Outputs:
		The tray stays closed, ItemSelected is 0, and the keyboard buffer is ["esc"].
*/
func TestDropdownEscOnClosedTrayIsNotConsumed(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	dropdownEntry := addTestDropdown(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)

	pressKey(simScreen, tcell.KeyEscape, 0, tcell.ModNone)

	if dropdownEntry.IsTrayOpen {
		test.Fatalf("expected Esc not to open a closed tray")
	}
	if dropdownEntry.ItemSelected != 0 {
		test.Fatalf("expected selection to stay 0, got %d", dropdownEntry.ItemSelected)
	}
	assertKeyboardBuffer(test, "after Esc on a closed tray", "esc")
}

/*
TestDropdownEnterCommitsHighlightedItem is a test which allows you to verify that Enter still opens a dropdown and then
commits the highlighted item when pressed again, and that neither press reaches the keyboard buffer.

Example:

	Expected Inputs:
		Dropdown with items [Zero, One, Two, Three] and item 0 selected, focused. Keys Enter, Down, Enter.

	Expected Outputs:
		The tray is closed, ItemSelected is 1, and the keyboard buffer is empty.
*/
func TestDropdownEnterCommitsHighlightedItem(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	dropdownEntry := addTestDropdown(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)

	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	pressKey(simScreen, tcell.KeyDown, 0, tcell.ModNone)
	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)

	if dropdownEntry.IsTrayOpen {
		test.Fatalf("expected the second Enter to close the tray")
	}
	if dropdownEntry.ItemSelected != 1 {
		test.Fatalf("expected item 1 to be committed, got %d", dropdownEntry.ItemSelected)
	}
	assertKeyboardBuffer(test, "after Enter, Down, Enter")
}

/*
TestDropdownSetSelectedItemIndexRejectsNegatives is a test which allows you to verify that SetSelectedItemIndex no
longer stores a negative index other than minus one. A stored minus five used to become the tray's highlight on
opening, and Up then set the tray's viewport to minus five, which panicked on the next screen update.

Example:

	Expected Inputs:
		The test dropdown with four items and item 0 selected, focused. SetSelectedItemIndex(-5), then (4), then
		(2), then (-1), then Enter to open the tray, Up, and Enter, updating the screen after each key.

	Expected Outputs:
		The index stays 0 after minus five and four, becomes 2, then minus one. Opening the tray and pressing Up
		does not panic, and Enter commits item 0.
*/
func TestDropdownSetSelectedItemIndexRejectsNegatives(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	addTestDropdown(layerAlias, styleEntry)
	dropdown := DropdownInstanceType{BaseControlInstanceType{layerAlias: layerAlias, controlAlias: "dropdown", controlType: constants.TYPE_DROPDOWN}}
	simScreen := startInputSimulation(test)
	for _, step := range []struct {
		itemIndex     int
		expectedIndex int
	}{
		{-5, 0},
		{4, 0},
		{2, 2},
		{constants.SELECTED_NONE, constants.SELECTED_NONE},
	} {
		dropdown.SetSelectedItemIndex(step.itemIndex)
		if selectedIndex := dropdown.GetSelectedItemIndex(); selectedIndex != step.expectedIndex {
			test.Fatalf("after SetSelectedItemIndex(%d): expected %d, got %d", step.itemIndex, step.expectedIndex, selectedIndex)
		}
	}
	for _, key := range []tcell.Key{tcell.KeyEnter, tcell.KeyUp, tcell.KeyEnter} {
		pressKey(simScreen, key, 0, tcell.ModNone)
		UpdateDisplay(false)
	}
	if selectedIndex := dropdown.GetSelectedItemIndex(); selectedIndex != 0 {
		test.Fatalf("expected item 0 to be committed, got %d", selectedIndex)
	}
}

/*
TestDropdownSetSelectionEntryResetsTray is a test which allows you to verify that replacing a dropdown's items resets
its tray and scroll bar to fit them, whether the new list is shorter, empty, or set while the tray is open.

Example:

	Expected Inputs:
		Case one: a dropdown with ten items and a three row tray, so its scroll bar is enabled, given two items. Case
		two: the test dropdown with its tray open, given an empty list, then Up, Down and Enter with the screen
		updated after each key.

	Expected Outputs:
		Case one: the scroll bar maximum is 0, it is disabled, and dragging its handle to the end leaves the tray's
		viewport at 0. Case two: no key panics, the selected index stays minus one, GetValue is empty, and the tray's
		scroll value is not negative.
*/
func TestDropdownSetSelectionEntryResetsTray(test *testing.T) {
	test.Run("shorter list", func(test *testing.T) {
		layerAlias, styleEntry := setupInputTest(test)
		dropdown := Dropdown.Add(layerAlias, "long", styleEntry, getTestSelectionEntry(10), 2, 2, 3, 10, 0)
		dropdownEntry := Dropdowns.Get(layerAlias, "long")
		dropdown.SetSelectionEntry(getTestSelectionEntry(2))
		scrollBarEntry := ScrollBars.Get(layerAlias, dropdownEntry.ScrollbarAlias)
		if scrollBarEntry.MaxScrollValue != 0 || scrollBarEntry.IsEnabled {
			test.Fatalf("expected maximum 0 and disabled, got %d and %t", scrollBarEntry.MaxScrollValue, scrollBarEntry.IsEnabled)
		}
		scrollBarEntry.IsEnabled = true
		scrollBarEntry.HandlePosition = scrollBarEntry.Length - 3
		scrollbar.computeValueByHandlePosition(layerAlias, dropdownEntry.ScrollbarAlias)
		if scrollBarEntry.ScrollValue != 0 {
			test.Fatalf("expected the scroll value to stay 0, got %d", scrollBarEntry.ScrollValue)
		}
	})
	test.Run("empty list while open", func(test *testing.T) {
		layerAlias, styleEntry := setupInputTest(test)
		dropdownEntry := addTestDropdown(layerAlias, styleEntry)
		dropdown := DropdownInstanceType{BaseControlInstanceType{layerAlias: layerAlias, controlAlias: "dropdown", controlType: constants.TYPE_DROPDOWN}}
		simScreen := startInputSimulation(test)
		pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
		if !dropdownEntry.IsTrayOpen {
			test.Fatalf("expected Enter to open the tray")
		}
		dropdown.SetSelectionEntry(NewSelectionEntry())
		for _, key := range []tcell.Key{tcell.KeyUp, tcell.KeyDown, tcell.KeyEnter} {
			pressKey(simScreen, key, 0, tcell.ModNone)
			UpdateDisplay(false)
		}
		if selectedIndex := dropdown.GetSelectedItemIndex(); selectedIndex != constants.SELECTED_NONE {
			test.Fatalf("expected no selection on an empty dropdown, got %d", selectedIndex)
		}
		if value := dropdown.GetValue(); value != "" {
			test.Fatalf("expected an empty value, got %q", value)
		}
		if scrollValue := ScrollBars.Get(layerAlias, dropdownEntry.ScrollbarAlias).ScrollValue; scrollValue < 0 {
			test.Fatalf("expected a scroll value of at least 0, got %d", scrollValue)
		}
	})
}

/*
TestDropdownAddValidatesDefaultItem is a test which allows you to verify that a dropdown's default item is validated
when it is added, and that its tray starts out agreeing with it, so that closing the tray without picking anything
keeps the default.

Example:

	Expected Inputs:
		A dropdown with four items added with default item seven, and another added with default item two whose tray
		is opened with Enter and then closed by closing every open dropdown.

	Expected Outputs:
		The first reports index minus one. The second still reports index 2 after its tray closes.
*/
func TestDropdownAddValidatesDefaultItem(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	outOfRangeDropdown := Dropdown.Add(layerAlias, "outOfRange", styleEntry, getTestSelectionEntry(4), 2, 2, 3, 10, 7)
	if selectedIndex := outOfRangeDropdown.GetSelectedItemIndex(); selectedIndex != constants.SELECTED_NONE {
		test.Fatalf("expected an out of range default to be cleared, got %d", selectedIndex)
	}
	defaultDropdown := Dropdown.Add(layerAlias, "default", styleEntry, getTestSelectionEntry(4), 20, 2, 3, 10, 2)
	simScreen := startInputSimulation(test)
	if err := SetFocus(&defaultDropdown); err != nil {
		test.Fatalf("expected SetFocus to succeed, got %v", err)
	}
	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	Dropdown.closeAllOpen()
	if selectedIndex := defaultDropdown.GetSelectedItemIndex(); selectedIndex != 2 {
		test.Fatalf("expected closing the tray without a pick to keep item 2, got %d", selectedIndex)
	}
}

/*
injectWheel is a method which allows you to inject a single mouse wheel notch at a location into the simulation screen
and process it through the event queue, exactly as a real wheel movement would be handled.

Example:

	injectWheel(simScreen, 5, 4, tcell.WheelDown)
*/
func injectWheel(simScreen tcell.SimulationScreen, xLocation int, yLocation int, wheelButton tcell.ButtonMask) {
	simScreen.InjectMouse(xLocation, yLocation, wheelButton, tcell.ModNone)
	UpdateEventQueues()
}

/*
TestDropdownWheelScrollsOpenTray is a test which allows you to verify that the mouse wheel scrolls an open dropdown
tray one row per notch over its items, its scroll bar and its border, that the scroll bar follows the tray, that the
hover highlight follows the item scrolled under the cursor, and that scrolling neither closes the tray nor changes the
selection.

Example:

	Expected Inputs:
		Dropdown at (2, 2) with items [Zero, One, Two, Three], a three row tray drawn from row 3, and item 0 selected.
		A click on the dropdown opens the tray. Then wheel down twice and wheel up once over row 4, wheel down once
		over the tray's scroll bar at column 14, wheel up once over the left border at (2, 4), and wheel down once
		over the bottom border at (6, 6).

	Expected Outputs:
		Viewport 1 and highlight 2 after the first wheel down, still viewport 1 after the second, which is the
		furthest the tray can scroll, then viewport 0 and highlight 1 after wheel up, viewport 1 after wheel down
		over the scroll bar, viewport 0 over the left border, and viewport 1 over the bottom border, with the
		highlight staying 1 away from the items. The scroll bar value equals the viewport each time, the tray stays
		open, and item 0 stays selected.
*/
func TestDropdownWheelScrollsOpenTray(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	dropdownEntry := addTestDropdown(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)
	clickAt(simScreen, 2, 2)
	if !dropdownEntry.IsTrayOpen {
		test.Fatalf("expected the click to open the tray")
	}
	selectorEntry := Selectors.Get(layerAlias, dropdownEntry.SelectorAlias)
	scrollBarEntry := ScrollBars.Get(layerAlias, dropdownEntry.ScrollbarAlias)
	assertTray := func(context string, expectedViewportPosition int, expectedHighlighted int) {
		test.Helper()
		if !dropdownEntry.IsTrayOpen {
			test.Fatalf("%s: expected the tray to stay open", context)
		}
		if scrollBarEntry.ScrollValue != expectedViewportPosition {
			test.Fatalf("%s: expected scroll bar value %d, got %d", context, expectedViewportPosition, scrollBarEntry.ScrollValue)
		}
		if dropdownEntry.ItemSelected != 0 {
			test.Fatalf("%s: expected the dropdown selection to stay 0, got %d", context, dropdownEntry.ItemSelected)
		}
		assertSelectorIndexes(test, context, selectorEntry, expectedHighlighted, 0, expectedViewportPosition)
	}

	injectWheel(simScreen, 5, 4, tcell.WheelDown)
	assertTray("after wheel down", 1, 2)
	injectWheel(simScreen, 5, 4, tcell.WheelDown)
	assertTray("after wheel down at the end of the list", 1, 2)
	injectWheel(simScreen, 5, 4, tcell.WheelUp)
	assertTray("after wheel up", 0, 1)
	injectWheel(simScreen, 14, 4, tcell.WheelDown)
	assertTray("after wheel down over the scroll bar", 1, 1)
	injectWheel(simScreen, 2, 4, tcell.WheelUp)
	assertTray("after wheel up over the left border", 0, 1)
	injectWheel(simScreen, 6, 6, tcell.WheelDown)
	assertTray("after wheel down over the bottom border", 1, 1)
}
