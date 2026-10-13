package consolizer

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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

/*
TestCheckboxSpaceTogglesFocusedCheckbox is a test which allows you to verify that Space toggles the focused checkbox and
is consumed, and that a checkbox can be reached with Tab.

Example:

	Expected Inputs:
		Tab order [a, check] with a focused, then Tab, Space, and Space.

	Expected Outputs:
		The checkbox is focused, checked after the first Space and unchecked after the second, and the keyboard buffer
		is empty.
*/
func TestCheckboxSpaceTogglesFocusedCheckbox(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a")
	checkbox := Checkbox.Add(layerAlias, "check", "Check", types.NewTuiStyleEntry(), 20, 2, false, true)
	buttons[0].AddToTabIndex()
	checkbox.AddToTabIndex()
	SetFocus(&buttons[0])
	readKeyboardBuffer()

	pressTab(simScreen)
	assertFocus(test, "after Tab", layerAlias, "check", constants.CellTypeCheckbox)
	pressKey(simScreen, tcell.KeyRune, ' ', tcell.ModNone)
	if !checkbox.IsSelected() {
		test.Fatalf("expected Space to check the checkbox")
	}
	pressKey(simScreen, tcell.KeyRune, ' ', tcell.ModNone)
	if checkbox.IsSelected() {
		test.Fatalf("expected a second Space to uncheck the checkbox")
	}
	assertKeyboardBuffer(test, "after Space")
}

/*
assertCheckboxMatchesMasterImage is a method which allows you to render the current screen and fail the test unless it
matches the checkbox suite's master image for the given test case. When master image updates are requested through the
environment, the rendered screen is saved as the new master before it is compared.

Example:

	assertCheckboxMatchesMasterImage(test, "TestCheckboxLabelFirstSelected")
*/
func assertCheckboxMatchesMasterImage(test *testing.T, testCaseName string) {
	test.Helper()
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, CHECKBOX_TEST_SUITE_NAME, testCaseName, obtainedValue)
	expectedValue := LoadMasterImage(CHECKBOX_TEST_SUITE_NAME, testCaseName)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", layerEntry.GetAnsiStringFromBase64(expectedValue))
		fmt.Println("Obtained:\n", layerEntry.GetAnsiStringFromBase64(obtainedValue))
	}
}

/*
TestCheckboxLabelFirstSelected is a test which verifies that a selected checkbox in label-first mode draws its label at
its X location and its selected mark one space after the label.

Example:

	Expected Inputs:
		A selected checkbox labelled "Test Checkbox" at (2,2) with SetLabelFirst(true).

	Expected Outputs:
		Row 2 shows "Test Checkbox" in columns 2 to 14, the layer beneath in column 15, and the selected mark in
		column 16, matching the TestCheckboxLabelFirstSelected master image.
*/
func TestCheckboxLabelFirstSelected(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	checkboxInstance := Checkbox.Add(layer1.layerAlias, "testCheckbox", "Test Checkbox", styleEntry, 2, 2, true, true)
	checkboxInstance.SetLabelFirst(true)
	assertCheckboxMatchesMasterImage(test, "TestCheckboxLabelFirstSelected")
}

/*
TestCheckboxLabelFirstUnselected is a test which verifies that an unselected checkbox in label-first mode draws its
label at its X location and its unselected mark one space after the label.

Example:

	Expected Inputs:
		An unselected checkbox labelled "Test Checkbox" at (2,2) with SetLabelFirst(true).

	Expected Outputs:
		Row 2 shows "Test Checkbox" in columns 2 to 14, the layer beneath in column 15, and the unselected mark in
		column 16, matching the TestCheckboxLabelFirstUnselected master image.
*/
func TestCheckboxLabelFirstUnselected(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	checkboxInstance := Checkbox.Add(layer1.layerAlias, "testCheckbox", "Test Checkbox", styleEntry, 2, 2, false, true)
	checkboxInstance.SetLabelFirst(true)
	assertCheckboxMatchesMasterImage(test, "TestCheckboxLabelFirstUnselected")
}

/*
TestCheckboxLabelFirstFocus is a test which verifies that a focused checkbox in label-first mode draws both its label
and its mark with the focused colours.

Example:

	Expected Inputs:
		An unselected checkbox labelled "Test Checkbox" at (2,2) with SetLabelFirst(true), then focused.

	Expected Outputs:
		The label in columns 2 to 14 and the mark in column 16 both carry the style's focused foreground and
		background colours, matching the TestCheckboxLabelFirstFocus master image.
*/
func TestCheckboxLabelFirstFocus(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	checkboxInstance := Checkbox.Add(layer1.layerAlias, "testCheckbox", "Test Checkbox", styleEntry, 2, 2, false, true)
	checkboxInstance.SetLabelFirst(true)
	setFocusedControl(layer1.layerAlias, "testCheckbox", constants.CellTypeCheckbox)
	assertCheckboxMatchesMasterImage(test, "TestCheckboxLabelFirstFocus")

	focusedForegroundColor, focusedBackgroundColor := getFocusedColors(styleEntry.Checkbox.ForegroundColor,
		styleEntry.Checkbox.BackgroundColor, styleEntry.Checkbox.FocusedForegroundColor, styleEntry.Checkbox.FocusedBackgroundColor)
	for _, xLocation := range []int{2, 14, 16} {
		attributeEntry := commonResource.screenLayer.CharacterMemory[2][xLocation].AttributeEntry
		assert.Equalf(test, focusedForegroundColor, attributeEntry.ForegroundColor, "column %d should use the focused foreground colour", xLocation)
		assert.Equalf(test, focusedBackgroundColor, attributeEntry.BackgroundColor, "column %d should use the focused background colour", xLocation)
	}
}

/*
TestCheckboxLabelFirstClickOnLabelToggles is a test which verifies that, in label-first mode, a mouse press on the
label toggles the checkbox, just as a press on the mark does.

Example:

	Expected Inputs:
		An unselected checkbox labelled "Test Checkbox" at (2,2) with SetLabelFirst(true), then a fresh left mouse
		press at (4,2), which is on the label.

	Expected Outputs:
		IsSelected reports true, and the redrawn screen shows the selected mark in column 16, matching the
		TestCheckboxLabelFirstClickOnLabelToggles master image.
*/
func TestCheckboxLabelFirstClickOnLabelToggles(test *testing.T) {
	resetMouseEventState()
	layer1, _, _, styleEntry := CommonTestSetup(test)
	checkboxInstance := Checkbox.Add(layer1.layerAlias, "testCheckbox", "Test Checkbox", styleEntry, 2, 2, false, true)
	checkboxInstance.SetLabelFirst(true)
	UpdateDisplay(false)

	SetMouseStatus(4, 2, 0, "")
	SetMouseStatus(4, 2, 1, "")
	Checkbox.updateMouseEvent()
	SetMouseStatus(4, 2, 0, "")

	assert.True(test, checkboxInstance.IsSelected(), "A press on the label should select the checkbox")
	assertCheckboxMatchesMasterImage(test, "TestCheckboxLabelFirstClickOnLabelToggles")
}

/*
TestCheckboxSwitchLabelFirstOnExistingCheckbox is a test which verifies that switching label-first mode on and off
for a checkbox that has already been drawn redraws it in the new layout, and that switching it off again restores
the default layout exactly.

Example:

	Expected Inputs:
		An unselected checkbox labelled "Test Checkbox" at (2,2), drawn once, then SetLabelFirst(true), then
		SetLabelFirst(false).

	Expected Outputs:
		After SetLabelFirst(true) the screen matches the TestCheckboxSwitchLabelFirstOnExistingCheckbox master image.
		After SetLabelFirst(false) the screen matches the TestCheckboxDefaultState master image again.
*/
func TestCheckboxSwitchLabelFirstOnExistingCheckbox(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	checkboxInstance := Checkbox.Add(layer1.layerAlias, "testCheckbox", "Test Checkbox", styleEntry, 2, 2, false, true)
	assertCheckboxMatchesMasterImage(test, "TestCheckboxDefaultState")

	checkboxInstance.SetLabelFirst(true)
	assertCheckboxMatchesMasterImage(test, "TestCheckboxSwitchLabelFirstOnExistingCheckbox")

	checkboxInstance.SetLabelFirst(false)
	assertCheckboxMatchesMasterImage(test, "TestCheckboxDefaultState")
}

/*
TestCheckboxLabelFirstWideRuneLabel is a test which verifies that label-first mode measures the label by its printed
width, so a label made of double-width runes places the mark after the full printed label.

Example:

	Expected Inputs:
		An unselected checkbox labelled "世界" (4 columns wide) at (2,2) with SetLabelFirst(true).

	Expected Outputs:
		Column 2 holds '世', column 4 holds '界', column 6 is not part of the checkbox, and column 7 holds the
		unselected mark with the checkbox cell type and control alias.
*/
func TestCheckboxLabelFirstWideRuneLabel(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	checkboxInstance := Checkbox.Add(layer1.layerAlias, "testCheckbox", "世界", styleEntry, 2, 2, false, true)
	checkboxInstance.SetLabelFirst(true)
	UpdateDisplay(false)

	row := commonResource.screenLayer.CharacterMemory[2]
	assert.Equal(test, '世', row[2].Character)
	assert.Equal(test, '界', row[4].Character)
	assert.NotEqual(test, constants.CellTypeCheckbox, row[6].AttributeEntry.CellType, "the separating column should not be drawn")
	assert.Equal(test, styleEntry.Checkbox.UnselectedCharacter, row[7].Character)
	assert.Equal(test, constants.CellTypeCheckbox, row[7].AttributeEntry.CellType)
	assert.Equal(test, "testCheckbox", row[7].AttributeEntry.CellControlAlias)
}

/*
TestCheckboxLabelFirstEmptyLabel is a test which verifies that a label-first checkbox with an empty label draws its
mark at its X location, exactly as the default mode does.

Example:

	Expected Inputs:
		An unselected checkbox with an empty label at (2,2) with SetLabelFirst(true).

	Expected Outputs:
		Column 2 holds the unselected mark with the checkbox cell type, and column 3 is not part of the checkbox.
*/
func TestCheckboxLabelFirstEmptyLabel(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	checkboxInstance := Checkbox.Add(layer1.layerAlias, "testCheckbox", "", styleEntry, 2, 2, false, true)
	checkboxInstance.SetLabelFirst(true)
	UpdateDisplay(false)

	row := commonResource.screenLayer.CharacterMemory[2]
	assert.Equal(test, styleEntry.Checkbox.UnselectedCharacter, row[2].Character)
	assert.Equal(test, constants.CellTypeCheckbox, row[2].AttributeEntry.CellType)
	assert.NotEqual(test, constants.CellTypeCheckbox, row[3].AttributeEntry.CellType)
}

/*
TestCheckboxSetLabelFirstOnDeletedCheckboxIsIgnored is a test which verifies that calling SetLabelFirst on a checkbox
instance whose control has been deleted is silently ignored and does not recreate the control.

Example:

	Expected Inputs:
		A checkbox that is added and then deleted through Checkbox.Delete, followed by SetLabelFirst(true) on the
		original instance.

	Expected Outputs:
		No panic, and the checkbox still does not exist in the memory manager.
*/
func TestCheckboxSetLabelFirstOnDeletedCheckboxIsIgnored(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	checkboxInstance := Checkbox.Add(layer1.layerAlias, "testCheckbox", "Test Checkbox", styleEntry, 2, 2, false, true)
	Checkbox.Delete(layer1.layerAlias, "testCheckbox")

	assert.NotPanics(test, func() { checkboxInstance.SetLabelFirst(true) })
	assert.False(test, Checkboxes.IsExists(layer1.layerAlias, "testCheckbox"))
}
