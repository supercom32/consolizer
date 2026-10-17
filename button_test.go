package consolizer

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

const BUTTON_TEST_SUITE_NAME = "button"

/*
TestButtonDefaultState is a test which verifies that a button is rendered correctly in its default state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button is rendered at (2,2) with label "Test" and width 10.
*/
func TestButtonDefaultState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonDefaultState", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonDefaultState") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonDisabledState is a test which verifies that a button is rendered correctly when it is in a disabled state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button is rendered with a disabled label color.
*/
func TestButtonDisabledState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, false)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonDisabledState", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonDisabledState") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonPressedState is a test which verifies that a button is rendered correctly when it is in a pressed state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button is rendered with a sunken frame style.
*/
func TestButtonPressedState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	buttonObject := layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)

	// Simulate a button press by directly setting the button state
	buttonEntry := buttonRegistry.Get(layer1.layerAlias, buttonObject.controlAlias)
	buttonEntry.IsPressed = true

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonPressedState", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonPressedState") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonSelectedState is a test which verifies that a button is rendered correctly when it is in a selected state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button label is underlined.
*/
func TestButtonSelectedState(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	buttonObject := layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)

	// Simulate a button selection by directly setting the button state
	buttonEntry := buttonRegistry.Get(layer1.layerAlias, buttonObject.controlAlias)
	buttonEntry.IsSelected = true

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonSelectedState", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonSelectedState") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonCustomDimensions is a test which verifies that a button is rendered correctly with custom width and height.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button is rendered with width 15 and height 5.
*/
func TestButtonCustomDimensions(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddButton("Test", styleEntry, 2, 2, 15, 5, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonCustomDimensions", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonCustomDimensions")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonLongLabel is a test which verifies that a button is rendered correctly when it has a label longer than its width.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button width expands to fit the long label.
*/
func TestButtonLongLabel(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddButton("This is a long button label", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonLongLabel", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonLongLabel") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonDelete is a test which verifies that a button is successfully removed when its Delete method is called.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button is absent from the rendered output after deletion.
*/
func TestButtonDelete(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	buttonInstance := layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)

	// Delete the button
	buttonInstance.Delete()

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonDelete", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonDelete") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonDeleteAll is a test which verifies that all buttons are successfully removed from a layer.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        All buttons are absent from the rendered output after calling DeleteAllButtons.
*/
func TestButtonDeleteAll(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddButton("Test 1", styleEntry, 2, 2, 10, 3, true)
	layer1.AddButton("Test 2", styleEntry, 2, 6, 10, 3, true)
	UpdateDisplay(false)

	// Delete all buttons
	Button.DeleteAll(layer1.layerAlias)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonDeleteAll", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonDeleteAll") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonFocus is a test which verifies that a button correctly handles focus state.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        Button focus state is reflected in the system.
*/
func TestButtonFocus(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)

	// Set focus to the button
	setFocusedControl(layer1.layerAlias, "testButton", constants.CellTypeButton)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonFocus", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonFocus") // This will be filled in after the first test run
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestButtonFlatStyle is a test which verifies that a flat button renders its frame without the two-tone bevel.

Example:
    Expected Inputs:
        styleEntry.Button.StyleMode set to constants.ButtonStyleFlat; a button "Test" at (2,2) width 10.

    Expected Outputs:
        Screen content matches the committed master image; the frame is a single uniform colour.
*/
func TestButtonFlatStyle(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	styleEntry.Button.StyleMode = constants.ButtonStyleFlat
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonFlatStyle", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonFlatStyle")
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", layerEntry.GetAnsiStringFromBase64(expectedValue))
		fmt.Println("Obtained:\n", layerEntry.GetAnsiStringFromBase64(obtainedValue))
	}
}

/*
TestButtonFlatFrameHasNoBevel is a test which verifies that a flat button's opposite frame corners share one
foreground colour, so there is no raised or sunken bevel, and that the colour is the button foreground.

Example:
    Expected Inputs:
        A flat button at (2,2) width 10 height 3 with a distinct grey Button.ForegroundColor.

    Expected Outputs:
        The top-left corner cell (2,2) and the bottom-right corner cell (11,4) both carry that grey foreground
        colour.
*/
func TestButtonFlatFrameHasNoBevel(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	frameColor := GetRGBColor(210, 210, 210)
	styleEntry.Button.StyleMode = constants.ButtonStyleFlat
	styleEntry.Button.ForegroundColor = frameColor
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)
	characterMemory := commonResource.screenLayer.CharacterMemory

	topLeftForeground := characterMemory[2][2].AttributeEntry.ForegroundColor
	bottomRightForeground := characterMemory[4][11].AttributeEntry.ForegroundColor
	assert.Equal(test, topLeftForeground, bottomRightForeground, "a flat frame has no bevel, so opposite corners match")
	assert.Equal(test, frameColor, topLeftForeground)
}

/*
TestButtonFlatPressedUsesPressedColors is a test which verifies that a pressed flat button, which has no bevel to
flip, shows the press through Button.PressedBackgroundColor on its cells.

Example:
    Expected Inputs:
        A flat button at (2,2) width 10 with Button.PressedBackgroundColor a distinct green; its entry
        IsPressed is set to true before drawing.

    Expected Outputs:
        A fill cell inside the button (3,3) carries the green pressed background colour.
*/
func TestButtonFlatPressedUsesPressedColors(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	pressedBackgroundColor := GetRGBColor(0, 150, 0)
	styleEntry.Button.StyleMode = constants.ButtonStyleFlat
	styleEntry.Button.PressedBackgroundColor = pressedBackgroundColor
	buttonInstance := layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	buttonEntry := buttonRegistry.Get(buttonInstance.layerAlias, buttonInstance.controlAlias)
	buttonEntry.IsPressed = true
	UpdateDisplay(false)
	characterMemory := commonResource.screenLayer.CharacterMemory

	assert.Equal(test, pressedBackgroundColor, characterMemory[3][3].AttributeEntry.BackgroundColor)
}

/*
TestButtonBorderlessStyle is a test which verifies that a borderless button renders its label with no frame at
all.

Example:
    Expected Inputs:
        styleEntry.Button.StyleMode set to constants.ButtonStyleBorderless; a button "Test" at (2,2) width 10.

    Expected Outputs:
        Screen content matches the committed master image.
*/
func TestButtonBorderlessStyle(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	styleEntry.Button.StyleMode = constants.ButtonStyleBorderless
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, BUTTON_TEST_SUITE_NAME, "TestButtonBorderlessStyle", obtainedValue)
	expectedValue := LoadMasterImage(BUTTON_TEST_SUITE_NAME, "TestButtonBorderlessStyle")
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", layerEntry.GetAnsiStringFromBase64(expectedValue))
		fmt.Println("Obtained:\n", layerEntry.GetAnsiStringFromBase64(obtainedValue))
	}
}

/*
TestButtonBorderlessHasNoFrame is a test which verifies that a borderless button draws no box-drawing characters
anywhere in its rectangle and still centres its label.

Example:
    Expected Inputs:
        A borderless button "Test" at (2,2) width 10 height 3.

    Expected Outputs:
        No cell in rows 2..4, columns 2..11 holds a box-drawing rune; the centred label "Test" starts at
        column 5 of row 3.
*/
func TestButtonBorderlessHasNoFrame(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	styleEntry.Button.StyleMode = constants.ButtonStyleBorderless
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)
	characterMemory := commonResource.screenLayer.CharacterMemory

	boxDrawingRunes := "┌┐└┘─│├┤┬┴┼╔╗╚╝═║"
	for rowIndex := 2; rowIndex <= 4; rowIndex++ {
		for columnIndex := 2; columnIndex <= 11; columnIndex++ {
			currentRune := characterMemory[rowIndex][columnIndex].Character
			assert.NotContainsf(test, boxDrawingRunes, string(currentRune), "borderless button drew a frame rune at (%d,%d)", columnIndex, rowIndex)
		}
	}
	assert.Equal(test, []rune("Test"), []rune{characterMemory[3][5].Character, characterMemory[3][6].Character, characterMemory[3][7].Character, characterMemory[3][8].Character})
}

/*
TestButtonBorderlessPressedUsesPressedColors is a test which verifies that a pressed borderless button shows the
press through Button.PressedBackgroundColor.

Example:
    Expected Inputs:
        A borderless button "Test" at (2,2) width 10 with a distinct green Button.PressedBackgroundColor; its
        entry IsPressed is set to true before drawing.

    Expected Outputs:
        A fill cell inside the button (3,3) carries the green pressed background colour.
*/
func TestButtonBorderlessPressedUsesPressedColors(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	pressedBackgroundColor := GetRGBColor(0, 150, 0)
	styleEntry.Button.StyleMode = constants.ButtonStyleBorderless
	styleEntry.Button.PressedBackgroundColor = pressedBackgroundColor
	buttonInstance := layer1.AddButton("Test", styleEntry, 2, 2, 10, 3, true)
	buttonEntry := buttonRegistry.Get(buttonInstance.layerAlias, buttonInstance.controlAlias)
	buttonEntry.IsPressed = true
	UpdateDisplay(false)
	characterMemory := commonResource.screenLayer.CharacterMemory

	assert.Equal(test, pressedBackgroundColor, characterMemory[3][3].AttributeEntry.BackgroundColor)
}

/*
TestButtonBorderlessAllowsSingleRow is a test which verifies that a borderless button is not forced to the
three-row minimum that the framed styles need.

Example:
    Expected Inputs:
        A borderless button "Test" at (2,2) with a requested height of 1.

    Expected Outputs:
        Row 2 in the button's column range carries the button cell type; row 3 does not.
*/
func TestButtonBorderlessAllowsSingleRow(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	styleEntry.Button.StyleMode = constants.ButtonStyleBorderless
	layer1.AddButton("Test", styleEntry, 2, 2, 10, 1, true)
	UpdateDisplay(false)
	characterMemory := commonResource.screenLayer.CharacterMemory

	assert.Equal(test, constants.CellTypeButton, characterMemory[2][3].AttributeEntry.CellType)
	assert.NotEqual(test, constants.CellTypeButton, characterMemory[3][3].AttributeEntry.CellType)
}

/*
TestFocusedButtonPressedByEnterAndSpace is a test which allows you to verify that Enter and Space on a focused button
record a press that GetPressed reports, exactly like a mouse click, and that neither key reaches the keyboard buffer.

Example:

	Expected Inputs:
		Button A focused via GetFocus, then Enter, then Space.

	Expected Outputs:
		GetPressed returns (layer, "a") after each key, and the keyboard buffer is empty.
*/
func TestFocusedButtonPressedByEnterAndSpace(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a")
	readKeyboardBuffer()
	buttons[0].GetFocus()

	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	pressedLayerAlias, pressedButtonAlias := buttons[0].GetPressed()
	if pressedLayerAlias != layerAlias || pressedButtonAlias != "a" {
		test.Fatalf("after Enter: expected (%q, %q), got (%q, %q)", layerAlias, "a", pressedLayerAlias, pressedButtonAlias)
	}

	pressKey(simScreen, tcell.KeyRune, ' ', tcell.ModNone)
	pressedLayerAlias, pressedButtonAlias = buttons[0].GetPressed()
	if pressedLayerAlias != layerAlias || pressedButtonAlias != "a" {
		test.Fatalf("after Space: expected (%q, %q), got (%q, %q)", layerAlias, "a", pressedLayerAlias, pressedButtonAlias)
	}
	assertKeyboardBuffer(test, "after Enter and Space")
}

/*
TestDisabledOrUnfocusedButtonIgnoresEnter is a test which allows you to verify that Enter does not press a focused
button that is disabled, nor any button when no button has focus, and that in both cases Enter reaches the keyboard
buffer as an ordinary keystroke.

Example:

	Expected Inputs:
		Button A disabled and focused, Enter. Then A enabled, nothing focused, Enter.

	Expected Outputs:
		GetPressed returns ("", "") both times, and the keyboard buffer is ["enter"] both times.
*/
func TestDisabledOrUnfocusedButtonIgnoresEnter(test *testing.T) {
	_, buttons, simScreen := setupTabOrderTest(test, "a")
	readKeyboardBuffer()
	buttons[0].SetEnabled(false)
	buttons[0].GetFocus()

	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	if pressedLayerAlias, pressedButtonAlias := buttons[0].GetPressed(); pressedLayerAlias != "" || pressedButtonAlias != "" {
		test.Fatalf("disabled button: expected no press, got (%q, %q)", pressedLayerAlias, pressedButtonAlias)
	}
	assertKeyboardBuffer(test, "disabled button", "enter")

	buttons[0].SetEnabled(true)
	setFocusedControl("", "", constants.NullControlType)
	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	if pressedLayerAlias, pressedButtonAlias := buttons[0].GetPressed(); pressedLayerAlias != "" || pressedButtonAlias != "" {
		test.Fatalf("nothing focused: expected no press, got (%q, %q)", pressedLayerAlias, pressedButtonAlias)
	}
	assertKeyboardBuffer(test, "nothing focused", "enter")
}

/*
assertButtonArmed is a method which allows you to fail the test unless a button's visual pressed state, the armed state
a release over it turns into a click, is the one expected.

Example:

	assertButtonArmed(test, "after the press", buttonInstance, true)
*/
func assertButtonArmed(test *testing.T, context string, buttonInstance ButtonInstanceType, isExpectedArmed bool) {
	test.Helper()
	if isArmed := buttonInstance.IsStatePressed(); isArmed != isExpectedArmed {
		test.Fatalf("%s: expected button %q armed to be %t, got %t", context, buttonInstance.GetAlias(), isExpectedArmed, isArmed)
	}
}

/*
TestButtonClickFiresOnReleaseOverPressedButton is a test which allows you to verify that a press and release over the
same button reports exactly one click, and that the button is armed only while the press is held.

Example:

	Expected Inputs:
		Button "a" at (12, 2) with size 8x3. A press at (14, 3), then a release at (14, 3).

	Expected Outputs:
		Armed after the press with no click reported yet. Disarmed after the release, GetPressed reports ("layer",
		"a") once and then ("", "").
*/
func TestButtonClickFiresOnReleaseOverPressedButton(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	buttonA := Button.Add(layerAlias, "a", "A", styleEntry, 12, 2, 8, 3, true)
	simScreen := startInputSimulation(test)

	moveMouseTo(simScreen, 14, 3, tcell.Button1)
	assertButtonArmed(test, "press", buttonA, true)
	assertPressed(test, "press", buttonA, "", "")
	moveMouseTo(simScreen, 14, 3, tcell.ButtonNone)
	assertButtonArmed(test, "release", buttonA, false)
	assertPressed(test, "release", buttonA, layerAlias, "a")
	assertPressed(test, "second read", buttonA, "", "")
}

/*
TestButtonPressDraggedOffAndReleasedDoesNotClick is a test which allows you to verify that dragging a held press off the
button it began on disarms that button, and that releasing while off it reports no click for any button.

Example:

	Expected Inputs:
		Button "a" at (12, 2) and button "b" at (20, 2), both 8x3. For each case, a press on "a" at (14, 3), a held
		move off it onto the empty cell (5, 7) or onto "b" at (22, 3), then a release there.

	Expected Outputs:
		"a" armed after the press, and neither button armed after the held move or the release. GetPressed reports
		("", "") for both buttons.
*/
func TestButtonPressDraggedOffAndReleasedDoesNotClick(test *testing.T) {
	testCases := []struct {
		name      string
		xLocation int
		yLocation int
	}{
		{name: "empty cell", xLocation: 5, yLocation: 7},
		{name: "other button", xLocation: 22, yLocation: 3},
	}
	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			layerAlias, styleEntry := setupInputTest(test)
			buttonA := Button.Add(layerAlias, "a", "A", styleEntry, 12, 2, 8, 3, true)
			buttonB := Button.Add(layerAlias, "b", "B", styleEntry, 20, 2, 8, 3, true)
			simScreen := startInputSimulation(test)

			moveMouseTo(simScreen, 14, 3, tcell.Button1)
			assertButtonArmed(test, "press", buttonA, true)
			moveMouseTo(simScreen, testCase.xLocation, testCase.yLocation, tcell.Button1)
			assertButtonArmed(test, "held move off", buttonA, false)
			assertButtonArmed(test, "held move off", buttonB, false)
			moveMouseTo(simScreen, testCase.xLocation, testCase.yLocation, tcell.ButtonNone)
			assertButtonArmed(test, "release", buttonA, false)
			assertButtonArmed(test, "release", buttonB, false)
			assertPressed(test, "release", buttonA, "", "")
			assertPressed(test, "release", buttonB, "", "")
		})
	}
}

/*
TestButtonPressDraggedOffAndBackClicksOnce is a test which allows you to verify that a held press dragged off the button
it began on and back onto it re-arms that button, and that the release then reports exactly one click.

Example:

	Expected Inputs:
		Button "a" at (12, 2) and button "b" at (20, 2), both 8x3. For each case, a press on "a" at (14, 3), a held
		move off it onto the empty cell (5, 7) or onto "b" at (22, 3), a held move back onto "a" at (16, 3), then a
		release at (16, 3).

	Expected Outputs:
		"a" armed after the press, disarmed while off it, and armed again once back on it, with "b" never armed.
		GetPressed reports ("layer", "a") once after the release and then ("", ""), and ("", "") for "b".
*/
func TestButtonPressDraggedOffAndBackClicksOnce(test *testing.T) {
	testCases := []struct {
		name      string
		xLocation int
		yLocation int
	}{
		{name: "empty cell", xLocation: 5, yLocation: 7},
		{name: "other button", xLocation: 22, yLocation: 3},
	}
	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			layerAlias, styleEntry := setupInputTest(test)
			buttonA := Button.Add(layerAlias, "a", "A", styleEntry, 12, 2, 8, 3, true)
			buttonB := Button.Add(layerAlias, "b", "B", styleEntry, 20, 2, 8, 3, true)
			simScreen := startInputSimulation(test)

			moveMouseTo(simScreen, 14, 3, tcell.Button1)
			assertButtonArmed(test, "press", buttonA, true)
			moveMouseTo(simScreen, testCase.xLocation, testCase.yLocation, tcell.Button1)
			assertButtonArmed(test, "held move off", buttonA, false)
			assertButtonArmed(test, "held move off", buttonB, false)
			moveMouseTo(simScreen, 16, 3, tcell.Button1)
			assertButtonArmed(test, "held move back", buttonA, true)
			assertButtonArmed(test, "held move back", buttonB, false)
			assertPressed(test, "held move back", buttonA, "", "")
			moveMouseTo(simScreen, 16, 3, tcell.ButtonNone)
			assertButtonArmed(test, "release", buttonA, false)
			assertPressed(test, "release", buttonA, layerAlias, "a")
			assertPressed(test, "second read", buttonA, "", "")
			assertPressed(test, "release", buttonB, "", "")
		})
	}
}

/*
TestButtonPressFromElsewhereNeverArms is a test which allows you to verify that a press which began anywhere other than
a button, or on a different button, never arms a button it is dragged onto, that releasing over that button reports no
click, and that a fresh click on it afterward still does.

Example:

	Expected Inputs:
		Button "a" at (12, 2) and button "b" at (20, 2), both 8x3. For each case, a press on the empty cell (5, 7) or
		on "b" at (22, 3), held moves onto "a" at (14, 3) and (16, 3), a release at (16, 3), then a click at (14, 3).

	Expected Outputs:
		"a" never armed during the held moves or the release, and GetPressed reports ("", "") for both buttons after
		the release. GetPressed reports ("layer", "a") after the fresh click.
*/
func TestButtonPressFromElsewhereNeverArms(test *testing.T) {
	testCases := []struct {
		name      string
		xLocation int
		yLocation int
	}{
		{name: "empty cell", xLocation: 5, yLocation: 7},
		{name: "other button", xLocation: 22, yLocation: 3},
	}
	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			layerAlias, styleEntry := setupInputTest(test)
			buttonA := Button.Add(layerAlias, "a", "A", styleEntry, 12, 2, 8, 3, true)
			buttonB := Button.Add(layerAlias, "b", "B", styleEntry, 20, 2, 8, 3, true)
			simScreen := startInputSimulation(test)

			moveMouseTo(simScreen, testCase.xLocation, testCase.yLocation, tcell.Button1)
			moveMouseTo(simScreen, 14, 3, tcell.Button1)
			assertButtonArmed(test, "held move onto the button", buttonA, false)
			moveMouseTo(simScreen, 16, 3, tcell.Button1)
			assertButtonArmed(test, "held move within the button", buttonA, false)
			moveMouseTo(simScreen, 16, 3, tcell.ButtonNone)
			assertButtonArmed(test, "release", buttonA, false)
			assertPressed(test, "release", buttonA, "", "")
			assertPressed(test, "release", buttonB, "", "")

			clickAt(simScreen, 14, 3)
			assertPressed(test, "fresh click", buttonA, layerAlias, "a")
		})
	}
}

/*
TestButtonAddedUnderHeldPressNeverArms is a test which allows you to verify that a press already held when a button is
added under the cursor never arms that button, and that a fresh click on it afterward reports a click.

Example:

	Expected Inputs:
		A press held over the empty cell (14, 3). Button "a" then added at (12, 2) with size 8x3 and drawn, followed by
		a held move to (16, 3), a release there, then a click at (14, 3).

	Expected Outputs:
		"a" never armed during the held move or the release, and GetPressed reports ("", "") after the release.
		GetPressed reports ("layer", "a") after the fresh click.
*/
func TestButtonAddedUnderHeldPressNeverArms(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	simScreen := startInputSimulation(test)

	moveMouseTo(simScreen, 14, 3, tcell.Button1)
	buttonA := Button.Add(layerAlias, "a", "A", styleEntry, 12, 2, 8, 3, true)
	UpdateDisplay(false)

	moveMouseTo(simScreen, 16, 3, tcell.Button1)
	assertButtonArmed(test, "held move", buttonA, false)
	moveMouseTo(simScreen, 16, 3, tcell.ButtonNone)
	assertButtonArmed(test, "release", buttonA, false)
	assertPressed(test, "release", buttonA, "", "")

	clickAt(simScreen, 14, 3)
	assertPressed(test, "fresh click", buttonA, layerAlias, "a")
}

/*
TestButtonReplacedDuringHeldPressNeverArms is a test which allows you to verify that a press which began on a button
does not carry over to a new button added under the same alias after the first was deleted, so the new button is never
armed by it and releasing over it reports no click.

Example:

	Expected Inputs:
		Button "a" at (12, 2) with size 8x3. A press at (14, 3), then "a" deleted, added again at the same place and
		drawn, followed by a held move to (16, 3), a release there, then a click at (14, 3).

	Expected Outputs:
		The original "a" armed after the press. The replacement never armed during the held move or the release, and
		GetPressed reports ("", "") after the release. GetPressed reports ("layer", "a") after the fresh click.
*/
func TestButtonReplacedDuringHeldPressNeverArms(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	originalButton := Button.Add(layerAlias, "a", "A", styleEntry, 12, 2, 8, 3, true)
	simScreen := startInputSimulation(test)

	moveMouseTo(simScreen, 14, 3, tcell.Button1)
	assertButtonArmed(test, "press", originalButton, true)
	Button.Delete(layerAlias, "a")
	replacementButton := Button.Add(layerAlias, "a", "A", styleEntry, 12, 2, 8, 3, true)
	UpdateDisplay(false)

	moveMouseTo(simScreen, 16, 3, tcell.Button1)
	assertButtonArmed(test, "held move", replacementButton, false)
	moveMouseTo(simScreen, 16, 3, tcell.ButtonNone)
	assertButtonArmed(test, "release", replacementButton, false)
	assertPressed(test, "release", replacementButton, "", "")

	clickAt(simScreen, 14, 3)
	assertPressed(test, "fresh click", replacementButton, layerAlias, "a")
}

/*
TestButtonPressRejectedByModalLayerNeverArms is a test which allows you to verify that a press on a button outside an
active modal layer is not recorded as beginning on that button, so the button is never armed by it, even when the
modal layer stops being modal while the press is still held.

Example:

	Expected Inputs:
		Button "a" at (12, 2) with size 8x3 on the main layer, and a modal layer at (30, 12). A press on "a" at
		(14, 3), the modal layer then made non-modal, a held move to (16, 3), a release there, then a click at (14, 3).

	Expected Outputs:
		"a" never armed after the press, the held move or the release, and GetPressed reports ("", "") after the
		release. GetPressed reports ("layer", "a") after the fresh click.
*/
func TestButtonPressRejectedByModalLayerNeverArms(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	buttonA := Button.Add(layerAlias, "a", "A", styleEntry, 12, 2, 8, 3, true)
	modalLayer := AddLayer(30, 12, 6, 3, 10, nil)
	modalLayer.SetModal(true)
	simScreen := startInputSimulation(test)

	moveMouseTo(simScreen, 14, 3, tcell.Button1)
	assertButtonArmed(test, "press", buttonA, false)
	modalLayer.SetModal(false)
	moveMouseTo(simScreen, 16, 3, tcell.Button1)
	assertButtonArmed(test, "held move", buttonA, false)
	moveMouseTo(simScreen, 16, 3, tcell.ButtonNone)
	assertButtonArmed(test, "release", buttonA, false)
	assertPressed(test, "release", buttonA, "", "")

	clickAt(simScreen, 14, 3)
	assertPressed(test, "fresh click", buttonA, layerAlias, "a")
}
