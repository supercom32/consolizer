package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

const RADIO_BUTTON_TEST_SUITE_NAME = "radio_button"

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
