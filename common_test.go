package consolizer

import (
	"encoding/base64"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
	"os"
	"testing"
	"time"
)

/*
UpdateMasterImages is a method which allows you to update the master regression files for a test. In addition, the
following should be noted:

- If the environment variable UPDATE_MASTER_IMAGES is set to true, it will always perform the update.

Example:
    isUpdated := UpdateMasterImages(false, "TestName", "base64data")
*/
func UpdateMasterImages(isUpdateRequested bool, testSuiteName string, testCaseName string, ansiBase64 string) bool {
	if os.Getenv("UPDATE_MASTER_IMAGES") == "true" || isUpdateRequested {
		fullPath := constants.MasterImagesPath + testSuiteName + "/"
		os.MkdirAll(fullPath, 0755)
		os.WriteFile(fullPath+testCaseName+".base64", []byte(ansiBase64), 0644)
		ansiData, _ := base64.StdEncoding.DecodeString(ansiBase64)
		os.WriteFile(fullPath+testCaseName+".ansi", ansiData, 0644)
		fmt.Println("Updated master image for: " + testSuiteName + "/" + testCaseName)
		return true
	}
	return false
}

/*
LoadMasterImage is a method which allows you to load a master regression file for a test.

Example:
    expectedValue := LoadMasterImage("TestSuite", "TestCase")
*/
func LoadMasterImage(testSuiteName string, testCaseName string) string {
	expectedValueBytes, _ := os.ReadFile(constants.MasterImagesPath + testSuiteName + "/" + testCaseName + ".base64")
	return string(expectedValueBytes)
}

/*
CommonTestSetup is a method which initializes a standard testing environment with multiple layers and a
default TUI style. In addition, the following should be noted:

- Registers RestoreTerminalSettings as a cleanup on test, so the terminal session this call starts is always torn
  down when the calling test finishes, even on a failure or panic. Without this, the background goroutines
  InitializeTerminal starts keep running into later tests and racing their own terminal state.

Example:
    layer1, layer2, layer3, styleEntry := CommonTestSetup(test)
*/
func CommonTestSetup(test testing.TB) (*LayerInstanceType, *LayerInstanceType, *LayerInstanceType, types.TuiStyleEntryType) {
	test.Helper()
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	test.Cleanup(RestoreTerminalSettings)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(3, 10, layerWidth, layerHeight, 2, nil)
	layer3 := AddLayer(0, 0, layerWidth, layerHeight, 3, nil)
	layer1.Color(4, 6)
	layer1.FillLayer("a1a2a3a4a5")
	layer2.Color(11, 9)
	layer2.FillLayer("a1a2a3a4a5")
	return layer1, layer2, layer3, styleEntry
}

/*
CommonTestSetupImages is a method which initializes a standard testing environment for image-related tests. In
addition, the following should be noted:

- Registers RestoreTerminalSettings as a cleanup on test, so the terminal session this call starts is always torn
  down when the calling test finishes, even on a failure or panic. Without this, the background goroutines
  InitializeTerminal starts keep running into later tests and racing their own terminal state.

Example:
    layer1, layer2, layer3, styleEntry := CommonTestSetupImages(test)
*/
func CommonTestSetupImages(test testing.TB) (*LayerInstanceType, *LayerInstanceType, *LayerInstanceType, types.TuiStyleEntryType) {
	test.Helper()
	commonResource.isDebugEnabled = true
	layerWidth := 50
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	test.Cleanup(RestoreTerminalSettings)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(0, 0, layerWidth, layerHeight, 2, nil)
	layer3 := AddLayer(0, 0, layerWidth, layerHeight, 3, nil)
	return layer1, layer2, layer3, styleEntry
}

/*
CommonTestSetupHighResolutionImages is a method which initializes a standard testing environment for high resolution
image tests. In addition, the following should be noted:

- Registers RestoreTerminalSettings as a cleanup on test, so the terminal session this call starts is always torn
  down when the calling test finishes, even on a failure or panic. Without this, the background goroutines
  InitializeTerminal starts keep running into later tests and racing their own terminal state.

Example:
    layer1, layer2, layer3, styleEntry := CommonTestSetupHighResolutionImages(test)
*/
func CommonTestSetupHighResolutionImages(test testing.TB) (*LayerInstanceType, *LayerInstanceType, *LayerInstanceType, types.TuiStyleEntryType) {
	test.Helper()
	commonResource.isDebugEnabled = true
	layerWidth := 140
	layerHeight := 50
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	test.Cleanup(RestoreTerminalSettings)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(0, 0, layerWidth, layerHeight, 2, nil)
	layer3 := AddLayer(0, 0, layerWidth, layerHeight, 3, nil)
	return layer1, layer2, layer3, styleEntry
}

/*
setupInputTest is a method which allows you to prepare a clean test environment for keyboard and focus tests,
returning the alias of the layer controls should be added to and the default style. Tab order, focus, mouse state,
the keyboard buffer, and the selector double click state are all reset before the test and again after it.

Example:

	layerAlias, styleEntry := setupInputTest(test)
*/
func setupInputTest(test *testing.T) (string, types.TuiStyleEntryType) {
	test.Helper()
	resetInputState()
	layer1, _, _, styleEntry := CommonTestSetup(test)
	test.Cleanup(resetInputState)
	return layer1.GetAlias(), styleEntry
}

/*
resetInputState is a method which allows you to reset every piece of package-level state that the keyboard and
focus tests depend on, including the modifier keys and highlighted control, so that no state leaks between tests
running in the same process.

Example:

	resetInputState()
*/
func resetInputState() {
	resetMouseEventState()
	ClearMouseMemory()
	ClearTabIndex()
	focusManager.mutex.Lock()
	focusManager.modalLayers = map[string]int{}
	focusManager.modalStack = nil
	focusManager.changeQueue = nil
	focusManager.isFocusIndicatorVisible = true
	focusManager.mutex.Unlock()
	setFocusedControl("", "", constants.NullControlType)
	setPreviouslyHighlightedControl("", "", constants.NullControlType)
	setModifierKeys(tcell.ModNone)
	readKeyboardBuffer()
	buttonHistory.clear()
	selectorLastClick = selectorClickType{}
	selectorDoubleClickInterval.Store(int64(constants.DefaultDoubleClickInterval * time.Millisecond))
}

/*
startInputSimulation is a method which allows you to render the controls added so far, so that mouse hit
testing can find them, and then install a simulation screen whose events the test drives itself through
UpdateEventQueues.

Example:

	simScreen := startInputSimulation(test)
*/
func startInputSimulation(test *testing.T) tcell.SimulationScreen {
	test.Helper()
	UpdateDisplay(false)
	return setupSimulationMouseScreen(test)
}

/*
pressKey is a method which allows you to inject a single key press into the simulation screen and process it through
the event queue, exactly as a real keypress would be handled.

Example:

	pressKey(simScreen, tcell.KeyTab, 0, tcell.ModShift)
*/
func pressKey(simScreen tcell.SimulationScreen, key tcell.Key, character rune, modifiers tcell.ModMask) {
	simScreen.InjectKey(key, character, modifiers)
	UpdateEventQueues()
}

/*
clickAt is a method which allows you to inject a complete left mouse click, a press followed by a release at the same
location, into the simulation screen and process both events through the event queue.

Example:

	clickAt(simScreen, 4, 3)
*/
func clickAt(simScreen tcell.SimulationScreen, xLocation int, yLocation int) {
	simScreen.InjectMouse(xLocation, yLocation, tcell.Button1, tcell.ModNone)
	UpdateEventQueues()
	simScreen.InjectMouse(xLocation, yLocation, tcell.ButtonNone, tcell.ModNone)
	UpdateEventQueues()
}

/*
readKeyboardBuffer is a method which allows you to drain the keyboard buffer and obtain every keystroke it held, in the
order they were added.

Example:

	keystrokes := readKeyboardBuffer()
*/
func readKeyboardBuffer() []string {
	var keystrokes []string
	for keystroke := KeyboardMemory.GetFromBuffer(); keystroke != nil; keystroke = KeyboardMemory.GetFromBuffer() {
		keystrokes = append(keystrokes, string(keystroke))
	}
	return keystrokes
}

/*
assertKeyboardBuffer is a method which allows you to fail the test unless the keyboard buffer holds exactly the
keystrokes expected, in order. The buffer is drained as part of the check.

Example:

	assertKeyboardBuffer(test, "after Esc", "esc")
*/
func assertKeyboardBuffer(test *testing.T, context string, expectedKeystrokes ...string) {
	test.Helper()
	keystrokes := readKeyboardBuffer()
	if len(keystrokes) != len(expectedKeystrokes) {
		test.Fatalf("%s: expected keyboard buffer %q, got %q", context, expectedKeystrokes, keystrokes)
	}
	for index := range keystrokes {
		if keystrokes[index] != expectedKeystrokes[index] {
			test.Fatalf("%s: expected keyboard buffer %q, got %q", context, expectedKeystrokes, keystrokes)
		}
	}
}

/*
addTestDropdown is a method which allows you to add a dropdown with the four items "Zero" to "Three" at
location (2, 2), with a three row tray and item 0 selected, and give it focus.

Example:

	dropdownEntry := addTestDropdown(layerAlias, styleEntry)
*/
func addTestDropdown(layerAlias string, styleEntry types.TuiStyleEntryType) *types.DropdownEntryType {
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.Add("zero", "Zero")
	selectionEntry.Add("one", "One")
	selectionEntry.Add("two", "Two")
	selectionEntry.Add("three", "Three")
	dropdownInstance := Dropdown.Add(layerAlias, "dropdown", styleEntry, selectionEntry, 2, 2, 3, 10, 0)
	dropdownInstance.GetFocus()
	return Dropdowns.Get(layerAlias, "dropdown")
}

/*
addTestSelector is a method which allows you to add a bordered, single column selector with the four items
"a" to "d" at location (2, 2), so that item N is drawn on row 2+N starting at column 2. Adding a selector does not
give it focus.

Example:

	selector := addTestSelector(layerAlias, styleEntry)
*/
func addTestSelector(layerAlias string, styleEntry types.TuiStyleEntryType) SelectorInstanceType {
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.Add("a", "Item A")
	selectionEntry.Add("b", "Item B")
	selectionEntry.Add("c", "Item C")
	selectionEntry.Add("d", "Item D")
	return Selector.Add(layerAlias, "selector", styleEntry, selectionEntry, 2, 2, 4, 10, 1, 0, 0, false, true)
}

/*
assertSelection is a method which allows you to fail the test unless the selector reports a new selection of the
expected item made by the expected source. The new selection flag is cleared by the check, through GetSelected.

Example:

	assertSelection(test, selector, "first click", 1, constants.SelectionSourceSingleClick)
*/
func assertSelection(test *testing.T, selector SelectorInstanceType, context string, expectedIndex int, expectedSource int) {
	test.Helper()
	if !selector.IsNewItemSelected() {
		test.Fatalf("%s: expected a new item to be selected", context)
	}
	if source := selector.GetSelectionSource(); source != expectedSource {
		test.Fatalf("%s: expected selection source %d, got %d", context, expectedSource, source)
	}
	if _, index := selector.GetSelected(); index != expectedIndex {
		test.Fatalf("%s: expected item %d to be selected, got %d", context, expectedIndex, index)
	}
}

/*
setupTabOrderTest is a method which allows you to prepare a clean tab navigation environment by initializing a test
terminal, clearing any tab order and focus left behind by earlier tests, adding one button per alias given, and
installing a simulation screen so Tab presses can be injected and processed synchronously. The created buttons are
returned in the same order as their aliases, but none of them are added to the tab order. In addition, the following
should be noted:

  - Focus, tab order, modal state, mouse state, and the keyboard buffer are package-level state that
    InitializeTerminal does not fully reset, so resetInputState clears them both before the test runs and
    again in a cleanup, keeping them from leaking between tests.

Example:

	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c")
*/
func setupTabOrderTest(test *testing.T, buttonAliases ...string) (
	string, []ButtonInstanceType, tcell.SimulationScreen,
) {
	test.Helper()
	resetInputState()
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layerAlias := layer1.GetAlias()
	test.Cleanup(resetInputState)
	buttons := make([]ButtonInstanceType, 0, len(buttonAliases))
	for index, buttonAlias := range buttonAliases {
		buttons = append(buttons, Button.Add(layerAlias, buttonAlias, buttonAlias, styleEntry, 1, 1+index*3, 8, 3, true))
	}
	UpdateDisplay(false)
	simScreen := setupSimulationMouseScreen(test)
	return layerAlias, buttons, simScreen
}

/*
pressTab is a method which allows you to inject a single Tab key press into the simulation screen and process it through
the event queue, exactly as a real keypress would be handled.

Example:

	pressTab(simScreen)
*/
func pressTab(simScreen tcell.SimulationScreen) {
	simScreen.InjectKey(tcell.KeyTab, 0, tcell.ModNone)
	UpdateEventQueues()
}

/*
assertButtonFocused is a method which allows you to fail the test if the button with the given alias on the given layer
is not the control that currently has focus.

Example:

	assertButtonFocused(test, layerAlias, "b", "after first Tab")
*/
func assertButtonFocused(test *testing.T, layerAlias string, buttonAlias string, context string) {
	test.Helper()
	if !isControlCurrentlyFocused(layerAlias, buttonAlias, constants.CellTypeButton) {
		focusedControl := getFocusedControl()
		test.Fatalf("%s: expected button %q on layer %q to be focused, got control %q on layer %q with type %d",
			context, buttonAlias, layerAlias, focusedControl.controlAlias, focusedControl.layerAlias, focusedControl.controlType)
	}
}

/*
assertPressed is a method which allows you to fail the test unless GetPressed reports exactly the given button, then
clears it. Passing an empty button alias asserts that no button was pressed.

Example:

	assertPressed(test, "after Enter", buttons[0], layerAlias, "ok")
*/
func assertPressed(test *testing.T, context string, anyButton ButtonInstanceType, expectedLayerAlias string, expectedButtonAlias string) {
	test.Helper()
	pressedLayerAlias, pressedButtonAlias := anyButton.GetPressed()
	if expectedButtonAlias == "" {
		expectedLayerAlias = ""
	}
	if pressedLayerAlias != expectedLayerAlias || pressedButtonAlias != expectedButtonAlias {
		test.Fatalf("%s: expected pressed button (%q, %q), got (%q, %q)", context, expectedLayerAlias, expectedButtonAlias,
			pressedLayerAlias, pressedButtonAlias)
	}
}
