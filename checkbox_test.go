package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

const CHECKBOX_TEST_SUITE_NAME = "checkbox"

/*
TestCheckboxDefaultState is a test which verifies that a checkbox is rendered correctly in its default state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Checkbox is rendered at (2,2) with label "Test Checkbox" and unselected.
*/
func TestCheckboxDefaultState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	Checkbox.Add(layer1.layerAlias, "testCheckbox", "Test Checkbox", styleEntry, 2, 2, false, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, CHECKBOX_TEST_SUITE_NAME, "TestCheckboxDefaultState", obtainedValue)
	expectedValue := LoadMasterImage(CHECKBOX_TEST_SUITE_NAME, "TestCheckboxDefaultState")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestCheckboxSelectedState is a test which verifies that a checkbox is rendered correctly when it is in a
selected state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Expects checkbox to be rendered with its selected character.
*/
func TestCheckboxSelectedState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	Checkbox.Add(layer1.layerAlias, "testCheckbox", "Test Checkbox", styleEntry, 2, 2, true, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, CHECKBOX_TEST_SUITE_NAME, "TestCheckboxSelectedState", obtainedValue)
	expectedValue := LoadMasterImage(CHECKBOX_TEST_SUITE_NAME, "TestCheckboxSelectedState")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestCheckboxDisabledState is a test which verifies that a checkbox is rendered correctly when it is in a
disabled state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Expects checkbox to be rendered with its disabled state appearance.
*/
func TestCheckboxDisabledState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	Checkbox.Add(layer1.layerAlias, "testCheckbox", "Test Checkbox", styleEntry, 2, 2, false, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, CHECKBOX_TEST_SUITE_NAME, "TestCheckboxDisabledState", obtainedValue)
	expectedValue := LoadMasterImage(CHECKBOX_TEST_SUITE_NAME, "TestCheckboxDisabledState")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestCheckboxToggleState is a test which verifies that a checkbox correctly toggles its selection state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Expects checkbox selection state to change after a simulated interaction.
*/
func TestCheckboxToggleState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	checkboxInstance := Checkbox.Add(layer1.layerAlias, "testCheckbox", "Test Checkbox", styleEntry, 2, 2, false, true)

	// Simulate a checkbox click by directly setting the checkbox state
	checkboxEntry := Checkboxes.Get(layer1.layerAlias, "testCheckbox")
	checkboxEntry.IsSelected = true

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, CHECKBOX_TEST_SUITE_NAME, "TestCheckboxToggleState", obtainedValue)
	expectedValue := LoadMasterImage(CHECKBOX_TEST_SUITE_NAME, "TestCheckboxToggleState")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the checkbox is selected
	if !assert.True(test, checkboxInstance.IsSelected(), "Checkbox should be selected") {
		fmt.Println("Checkbox selection state is incorrect")
	}
}

/*
TestCheckboxMultiple is a test which verifies that multiple checkboxes are rendered correctly on the same
layer.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Expects multiple checkboxes with different states to be rendered correctly.
*/
func TestCheckboxMultiple(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	Checkbox.Add(layer1.layerAlias, "testCheckbox1", "Checkbox 1", styleEntry, 2, 2, false, true)
	Checkbox.Add(layer1.layerAlias, "testCheckbox2", "Checkbox 2", styleEntry, 2, 4, true, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, CHECKBOX_TEST_SUITE_NAME, "TestCheckboxMultiple", obtainedValue)
	expectedValue := LoadMasterImage(CHECKBOX_TEST_SUITE_NAME, "TestCheckboxMultiple")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestCheckboxDelete is a test which verifies that a checkbox is successfully removed when its Delete method
is called.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Expects checkbox to be absent from the rendered output after deletion.
*/
func TestCheckboxDelete(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	checkboxInstance := Checkbox.Add(layer1.layerAlias, "testCheckbox", "Test Checkbox", styleEntry, 2, 2, false, true)
	UpdateDisplay(false)

	// Delete the checkbox
	checkboxInstance.Delete()

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, CHECKBOX_TEST_SUITE_NAME, "TestCheckboxDelete", obtainedValue)
	expectedValue := LoadMasterImage(CHECKBOX_TEST_SUITE_NAME, "TestCheckboxDelete")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestCheckboxDeleteAll is a test which verifies that all checkboxes are successfully removed from a layer.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Expects all checkboxes to be absent from the rendered output after calling DeleteAllCheckboxes.
*/
func TestCheckboxDeleteAll(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	Checkbox.Add(layer1.layerAlias, "testCheckbox1", "Checkbox 1", styleEntry, 2, 2, false, true)
	Checkbox.Add(layer1.layerAlias, "testCheckbox2", "Checkbox 2", styleEntry, 2, 4, true, true)
	UpdateDisplay(false)

	// Delete all checkboxes
	Checkbox.DeleteAll(layer1.layerAlias)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, CHECKBOX_TEST_SUITE_NAME, "TestCheckboxDeleteAll", obtainedValue)
	expectedValue := LoadMasterImage(CHECKBOX_TEST_SUITE_NAME, "TestCheckboxDeleteAll")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestCheckboxFocus is a test which verifies that a checkbox correctly handles focus state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Expects checkbox focus state to be reflected in the system.
*/
func TestCheckboxFocus(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	Checkbox.Add(layer1.layerAlias, "testCheckbox", "Test Checkbox", styleEntry, 2, 2, false, true)

	// Set focus to the checkbox
	setFocusedControl(layer1.layerAlias, "testCheckbox", constants.CellTypeCheckbox)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, CHECKBOX_TEST_SUITE_NAME, "TestCheckboxFocus", obtainedValue)
	expectedValue := LoadMasterImage(CHECKBOX_TEST_SUITE_NAME, "TestCheckboxFocus")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestCheckboxLongLabel is a test which verifies that a checkbox is rendered correctly when it has a very long
label.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Expects checkbox label to be rendered without breaking the layout.
*/
func TestCheckboxLongLabel(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	Checkbox.Add(layer1.layerAlias, "testCheckbox", "This is a very long checkbox label to test text wrapping", styleEntry, 2, 2, false, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, CHECKBOX_TEST_SUITE_NAME, "TestCheckboxLongLabel", obtainedValue)
	expectedValue := LoadMasterImage(CHECKBOX_TEST_SUITE_NAME, "TestCheckboxLongLabel")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}
