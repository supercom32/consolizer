package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

const RADIO_BUTTON_TEST_SUITE_NAME = "radio_button"

/*
TestRadioButtonDefaultState is a test which verifies that a radio button control is rendered correctly with its
default (unselected) state.

Example:
    Expected Inputs:
        A radio button added to a layer with isSelected set to false.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the unselected character.
*/
func TestRadioButtonDefaultState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	radioButton.Add(layer1.layerAlias, "testRadioButton", "Test Radio Button", styleEntry, 2, 2, 1, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonDefaultState", obtainedValue)
	expectedValue := LoadMasterImage(RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonDefaultState")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestRadioButtonSelectedState is a test which verifies that a radio button control is rendered correctly with its
selected state.

Example:
    Expected Inputs:
        A radio button added to a layer with isSelected set to true.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the selected character.
*/
func TestRadioButtonSelectedState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	radioButton.Add(layer1.layerAlias, "testRadioButton", "Test Radio Button", styleEntry, 2, 2, 1, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonSelectedState", obtainedValue)
	expectedValue := LoadMasterImage(RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonSelectedState")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestRadioButtonGroup is a test which verifies that only one radio button within a single group can be selected at
a time.

Example:
    Expected Inputs:
        Three radio buttons added to the same group ID.
    Expected Outputs:
        Only the first radio button is reported as selected, and the others are false.
*/
func TestRadioButtonGroup(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	radioButtonInstance1 := radioButton.Add(layer1.layerAlias, "testRadioButton1", "Radio Button 1", styleEntry, 2, 2, 1, true)
	radioButtonInstance2 := radioButton.Add(layer1.layerAlias, "testRadioButton2", "Radio Button 2", styleEntry, 2, 4, 1, false)
	radioButtonInstance3 := radioButton.Add(layer1.layerAlias, "testRadioButton3", "Radio Button 3", styleEntry, 2, 6, 1, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonGroup", obtainedValue)
	expectedValue := LoadMasterImage(RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonGroup")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the first radio button is selected
	if !assert.True(test, radioButtonInstance1.IsSelected(), "Radio Button 1 should be selected") {
		fmt.Println("Radio Button 1 selection state is incorrect")
	}

	// Verify the other radio buttons are not selected
	if !assert.False(test, radioButtonInstance2.IsSelected(), "Radio Button 2 should not be selected") {
		fmt.Println("Radio Button 2 selection state is incorrect")
	}

	if !assert.False(test, radioButtonInstance3.IsSelected(), "Radio Button 3 should not be selected") {
		fmt.Println("Radio Button 3 selection state is incorrect")
	}
}

/*
TestRadioButtonMultipleGroups is a test which verifies that radio buttons in different groups operate independently.

Example:
    Expected Inputs:
        Two groups of radio buttons, each with two buttons.
    Expected Outputs:
        One button in each group can be selected simultaneously without affecting the other group.
*/
func TestRadioButtonMultipleGroups(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	radioButtonInstance1 := radioButton.Add(layer1.layerAlias, "testRadioButton1", "Group 1 - Button 1", styleEntry, 2, 2, 1, true)
	radioButtonInstance2 := radioButton.Add(layer1.layerAlias, "testRadioButton2", "Group 1 - Button 2", styleEntry, 2, 4, 1, false)
	radioButtonInstance3 := radioButton.Add(layer1.layerAlias, "testRadioButton3", "Group 2 - Button 1", styleEntry, 2, 6, 2, true)
	radioButtonInstance4 := radioButton.Add(layer1.layerAlias, "testRadioButton4", "Group 2 - Button 2", styleEntry, 2, 8, 2, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonMultipleGroups", obtainedValue)
	expectedValue := LoadMasterImage(RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonMultipleGroups")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the first radio button in each group is selected
	if !assert.True(test, radioButtonInstance1.IsSelected(), "Group 1 - Button 1 should be selected") {
		fmt.Println("Group 1 - Button 1 selection state is incorrect")
	}

	if !assert.True(test, radioButtonInstance3.IsSelected(), "Group 2 - Button 1 should be selected") {
		fmt.Println("Group 2 - Button 1 selection state is incorrect")
	}

	// Verify the second radio button in each group is not selected
	if !assert.False(test, radioButtonInstance2.IsSelected(), "Group 1 - Button 2 should not be selected") {
		fmt.Println("Group 1 - Button 2 selection state is incorrect")
	}

	if !assert.False(test, radioButtonInstance4.IsSelected(), "Group 2 - Button 2 should not be selected") {
		fmt.Println("Group 2 - Button 2 selection state is incorrect")
	}
}

/*
TestRadioButtonChangeSelection is a test which verifies that selecting a new radio button correctly updates the
selection state of all buttons in that group.

Example:
    Expected Inputs:
        A group of two buttons where the selection is programmatically changed from the first to the second.
    Expected Outputs:
        The first button becomes unselected and the second button becomes selected.
*/
func TestRadioButtonChangeSelection(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	radioButtonInstance1 := radioButton.Add(layer1.layerAlias, "testRadioButton1", "Radio Button 1", styleEntry, 2, 2, 1, true)
	radioButtonInstance2 := radioButton.Add(layer1.layerAlias, "testRadioButton2", "Radio Button 2", styleEntry, 2, 4, 1, false)

	// Verify initial state
	if !assert.True(test, radioButtonInstance1.IsSelected(), "Radio Button 1 should initially be selected") {
		fmt.Println("Radio Button 1 initial selection state is incorrect")
	}

	if !assert.False(test, radioButtonInstance2.IsSelected(), "Radio Button 2 should initially not be selected") {
		fmt.Println("Radio Button 2 initial selection state is incorrect")
	}

	// Change selection by selecting the second radio button
	selectRadioButton(layer1.layerAlias, "testRadioButton2")

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonChangeSelection", obtainedValue)
	expectedValue := LoadMasterImage(RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonChangeSelection")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the selection changed
	if !assert.False(test, radioButtonInstance1.IsSelected(), "Radio Button 1 should now be unselected") {
		fmt.Println("Radio Button 1 updated selection state is incorrect")
	}

	if !assert.True(test, radioButtonInstance2.IsSelected(), "Radio Button 2 should now be selected") {
		fmt.Println("Radio Button 2 updated selection state is incorrect")
	}
}

/*
TestRadioButtonGetSelected is a test which verifies that the GetSelected method correctly returns the alias of
the selected radio button in a group.

Example:
    Expected Inputs:
        A group of two radio buttons with the first one selected.
    Expected Outputs:
        GetSelected returns the alias of the first radio button.
*/
func TestRadioButtonGetSelected(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	radioButtonInstance1 := radioButton.Add(layer1.layerAlias, "testRadioButton1", "Radio Button 1", styleEntry, 2, 2, 1, true)
	radioButton.Add(layer1.layerAlias, "testRadioButton2", "Radio Button 2", styleEntry, 2, 4, 1, false)

	// Get the selected radio button
	selectedRadioButton := radioButtonInstance1.GetSelected()

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonGetSelected", obtainedValue)
	expectedValue := LoadMasterImage(RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonGetSelected")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}

	// Verify the correct radio button is reported as selected
	if !assert.Equal(test, "testRadioButton1", selectedRadioButton, "testRadioButton1 should be reported as selected") {
		fmt.Println("GetSelectedRadioButton returned incorrect value")
	}
}

/*
TestRadioButtonDelete is a test which verifies that a radio button control can be successfully removed from its
parent layer.

Example:
    Expected Inputs:
        A layer containing a single radio button that is subsequently deleted.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) for an empty layer.
*/
func TestRadioButtonDelete(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	radioButtonInstance := radioButton.Add(layer1.layerAlias, "testRadioButton", "Test Radio Button", styleEntry, 2, 2, 1, true)
	UpdateDisplay(false)

	// Delete the radio button
	radioButtonInstance.Delete()

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonDelete", obtainedValue)
	expectedValue := LoadMasterImage(RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonDelete")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestRadioButtonDeleteAll is a test which verifies that all radio button controls on a layer can be successfully
removed at once.

Example:
    Expected Inputs:
        A layer containing multiple radio buttons from different groups, followed by a DeleteAll call.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) for an empty layer after all controls are removed.
*/
func TestRadioButtonDeleteAll(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	radioButton.Add(layer1.layerAlias, "testRadioButton1", "Radio Button 1", styleEntry, 2, 2, 1, true)
	radioButton.Add(layer1.layerAlias, "testRadioButton2", "Radio Button 2", styleEntry, 2, 4, 1, false)
	radioButton.Add(layer1.layerAlias, "testRadioButton3", "Radio Button 3", styleEntry, 2, 6, 2, true)
	UpdateDisplay(false)

	// Delete all radio buttons
	radioButton.DeleteAll(layer1.layerAlias)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonDeleteAll", obtainedValue)
	expectedValue := LoadMasterImage(RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonDeleteAll")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestRadioButtonFocus is a test which verifies that a radio button control correctly handles gaining focus and
renders its focused state accordingly.

Example:
    Expected Inputs:
        A focused radio button control added to a layer.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the radio button in its focused visual state.
*/
func TestRadioButtonFocus(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	radioButton.Add(layer1.layerAlias, "testRadioButton", "Test Radio Button", styleEntry, 2, 2, 1, true)

	// Set focus to the radio button
	setFocusedControl(layer1.layerAlias, "testRadioButton", constants.CellTypeRadioButton)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonFocus", obtainedValue)
	expectedValue := LoadMasterImage(RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonFocus")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestRadioButtonLongLabel is a test which verifies that a radio button with a long label is rendered correctly,
even if it extends beyond standard screen boundaries.

Example:
    Expected Inputs:
        A radio button with a very long label string.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) and no panics occur during rendering.
*/
func TestRadioButtonLongLabel(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	// We color this at 1,1,1 because intellij has a color bug which makes black transparent.
	styleEntry.RadioButton.BackgroundColor = GetRGBColor(1, 1, 1)
	radioButton.Add(layer1.layerAlias, "testRadioButton", "This is a very long radio button label to test text wrapping", styleEntry, -10, 2, 1, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonLongLabel", obtainedValue)
	expectedValue := LoadMasterImage(RADIO_BUTTON_TEST_SUITE_NAME, "TestRadioButtonLongLabel")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}
