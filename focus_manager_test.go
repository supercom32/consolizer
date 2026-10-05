package consolizer

import (
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
	"testing"
)

/*
assertFocus is a method which allows you to fail the test unless GetFocusedControl reports the given layer alias,
control alias, and cell type. Passing empty aliases with constants.NullControlType asserts that nothing has focus.

Example:

	assertFocus(test, "after Tab", layerAlias, "b", constants.CellTypeButton)
*/
func assertFocus(test *testing.T, context string, expectedLayerAlias string, expectedControlAlias string, expectedControlType int) {
	test.Helper()
	layerAlias, controlAlias, controlType := GetFocusedControl()
	if layerAlias != expectedLayerAlias || controlAlias != expectedControlAlias || controlType != expectedControlType {
		test.Fatalf("%s: expected focus (%q, %q, %d), got (%q, %q, %d)", context, expectedLayerAlias, expectedControlAlias,
			expectedControlType, layerAlias, controlAlias, controlType)
	}
}

/*
addDialogLayer is a method which allows you to add a dialog layer at (20, 1), sized 18 by 10 and drawn above every layer
created by CommonTestSetup, holding an "ok" button at (1, 1) and a "cancel" button at (1, 4), both added to the
dialog's tab order in that order.

Example:

	dialog, okButton, cancelButton := addDialogLayer(test, styleEntry)
*/
func addDialogLayer(test *testing.T, styleEntry types.TuiStyleEntryType) (*LayerInstanceType, ButtonInstanceType, ButtonInstanceType) {
	test.Helper()
	dialog := AddLayer(20, 1, 18, 10, 4, nil)
	okButton := Button.Add(dialog.GetAlias(), "ok", "ok", styleEntry, 1, 1, 8, 3, true)
	cancelButton := Button.Add(dialog.GetAlias(), "cancel", "cancel", styleEntry, 1, 4, 8, 3, true)
	if err := okButton.AddToTabIndex(); err != nil {
		test.Fatalf("expected the ok button to join the tab order, got %v", err)
	}
	if err := cancelButton.AddToTabIndex(); err != nil {
		test.Fatalf("expected the cancel button to join the tab order, got %v", err)
	}
	return dialog, okButton, cancelButton
}

/*
TestTabAndShiftTabWrapWithinTheLayerScope is a test which allows you to verify that Tab and Shift+Tab wrap around the
focused control's own layer, and never move focus to a stop registered on a different layer.

Example:

	Expected Inputs:
		Tab order [a, b, c] on layer 1 and [x] on a second layer. a focused, then Tab three times, then Shift+Tab.
		Then x focused with SetFocus and Tab pressed.

	Expected Outputs:
		Focus is b, c, a after the three Tabs and c after Shift+Tab. From x, Tab keeps focus on x, the only stop on its
		layer.
*/
func TestTabAndShiftTabWrapWithinTheLayerScope(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c")
	otherLayer := AddLayer(20, 1, 15, 8, 4, nil)
	otherButton := Button.Add(otherLayer.GetAlias(), "x", "x", types.NewTuiStyleEntry(), 1, 1, 8, 3, true)
	for index := range buttons {
		buttons[index].AddToTabIndex()
	}
	otherButton.AddToTabIndex()
	SetFocus(&buttons[0])

	for _, expectedAlias := range []string{"b", "c", "a"} {
		pressTab(simScreen)
		assertFocus(test, "after Tab", layerAlias, expectedAlias, constants.CellTypeButton)
	}
	pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModNone)
	assertFocus(test, "after Shift+Tab", layerAlias, "c", constants.CellTypeButton)

	if err := SetFocus(&otherButton); err != nil {
		test.Fatalf("expected x to take focus, got %v", err)
	}
	pressTab(simScreen)
	assertFocus(test, "Tab from x", otherLayer.GetAlias(), "x", constants.CellTypeButton)
}

/*
TestFocusedControlLosingFocusabilityMovesToNextStop is a test which allows you to verify that disabling, hiding, or
deleting the focused control moves focus to the next stop in its layer, wrapping around, and that focus is cleared once
no stop can take it.

Example:

	Expected Inputs:
		Tab order [a, b, c, d] with b focused. b disabled, then c hidden, then d deleted, then a disabled.

	Expected Outputs:
		Focus is c, then d, then a (wrapping), then nothing.
*/
func TestFocusedControlLosingFocusabilityMovesToNextStop(test *testing.T) {
	layerAlias, buttons, _ := setupTabOrderTest(test, "a", "b", "c", "d")
	for index := range buttons {
		buttons[index].AddToTabIndex()
	}
	SetFocus(&buttons[1])

	buttons[1].SetEnabled(false)
	assertFocus(test, "after disabling b", layerAlias, "c", constants.CellTypeButton)
	buttons[2].SetVisible(false)
	assertFocus(test, "after hiding c", layerAlias, "d", constants.CellTypeButton)
	buttons[3].Delete()
	assertFocus(test, "after deleting d", layerAlias, "a", constants.CellTypeButton)
	buttons[0].SetEnabled(false)
	assertFocus(test, "after disabling a", "", "", constants.NullControlType)
}

/*
TestSetFocusThenTabContinuesFromIt is a test which allows you to verify that SetFocus and Tab share one focus state, so
Tab and Shift+Tab continue from the control given to SetFocus, and that SetFocus refuses a disabled control without
changing focus.

Example:

	Expected Inputs:
		Tab order [a, b, c]. SetFocus(c) then Tab, SetFocus(a) then Shift+Tab, then SetFocus on a disabled b.

	Expected Outputs:
		Focus is a after the Tab and c after the Shift+Tab. SetFocus on b returns an error and focus stays on c.
*/
func TestSetFocusThenTabContinuesFromIt(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c")
	for index := range buttons {
		buttons[index].AddToTabIndex()
	}

	SetFocus(&buttons[2])
	pressTab(simScreen)
	assertFocus(test, "Tab after SetFocus(c)", layerAlias, "a", constants.CellTypeButton)

	SetFocus(&buttons[0])
	pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModNone)
	assertFocus(test, "Shift+Tab after SetFocus(a)", layerAlias, "c", constants.CellTypeButton)

	buttons[1].SetEnabled(false)
	if err := SetFocus(&buttons[1]); err == nil {
		test.Fatalf("expected SetFocus on a disabled button to fail")
	}
	assertFocus(test, "after refused SetFocus", layerAlias, "c", constants.CellTypeButton)
}

/*
TestModalLayerTrapsFocusAndRestoresIt is a test which allows you to verify that a modal layer takes focus when it opens,
keeps Tab, SetFocus, and mouse clicks inside it, and returns focus to the control that had it before when it is hidden
or deleted, and takes it again when shown again.

Example:

	Expected Inputs:
		Tab order [a, b] on layer 1 with b focused. A dialog with tab order [ok, cancel] is made modal. Tab twice,
		SetFocus(a), and a click on a follow. The dialog is hidden, shown again, and finally deleted.

	Expected Outputs:
		Focus moves to ok when the modal opens. Tab gives cancel then ok. SetFocus(a) fails, and the click on a neither
		focuses nor presses it. Hiding the dialog restores b, showing it moves focus back to ok, and deleting it
		restores b again.
*/
func TestModalLayerTrapsFocusAndRestoresIt(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b")
	buttons[0].AddToTabIndex()
	buttons[1].AddToTabIndex()
	SetFocus(&buttons[1])
	dialog, okButton, _ := addDialogLayer(test, types.NewTuiStyleEntry())
	UpdateDisplay(false)

	dialog.SetModal(true)
	assertFocus(test, "modal opened", dialog.GetAlias(), "ok", constants.CellTypeButton)
	pressTab(simScreen)
	assertFocus(test, "first Tab in modal", dialog.GetAlias(), "cancel", constants.CellTypeButton)
	pressTab(simScreen)
	assertFocus(test, "second Tab in modal", dialog.GetAlias(), "ok", constants.CellTypeButton)

	if err := SetFocus(&buttons[0]); err == nil {
		test.Fatalf("expected SetFocus outside the modal to fail")
	}
	clickAt(simScreen, 4, 2)
	assertFocus(test, "click outside modal", dialog.GetAlias(), "ok", constants.CellTypeButton)
	assertPressed(test, "click outside modal", okButton, "", "")

	dialog.SetIsVisible(false)
	assertFocus(test, "modal hidden", layerAlias, "b", constants.CellTypeButton)
	dialog.SetIsVisible(true)
	assertFocus(test, "modal shown again", dialog.GetAlias(), "ok", constants.CellTypeButton)
	dialog.Delete()
	assertFocus(test, "modal deleted", layerAlias, "b", constants.CellTypeButton)
}

/*
TestRadioGroupIsOneTabStopAndArrowsSelect is a test which allows you to verify that a radio group is a single tab stop,
entered at its selected button, and that the arrow keys move focus within the group, wrapping, and select the button
they land on.

Example:

	Expected Inputs:
		Tab order [before, r1, r2, r3, after] where r1 to r3 are one group stacked top to bottom and r2 is selected.
		before focused, then Tab, Tab, Shift+Tab, Down, Down, Up, and Right.

	Expected Outputs:
		Tab gives r2 then after, Shift+Tab gives r2. Down moves to r3 then wraps to r1, Up returns to r3, and Right
		wraps to r1, each selecting the button it lands on.
*/
func TestRadioGroupIsOneTabStopAndArrowsSelect(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "before", "after")
	styleEntry := types.NewTuiStyleEntry()
	radios := []RadioButtonInstanceType{
		radioButton.Add(layerAlias, "r1", "One", styleEntry, 20, 2, 1, false),
		radioButton.Add(layerAlias, "r2", "Two", styleEntry, 20, 3, 1, true),
		radioButton.Add(layerAlias, "r3", "Three", styleEntry, 20, 4, 1, false),
	}
	buttons[0].AddToTabIndex()
	for index := range radios {
		radios[index].AddToTabIndex()
	}
	buttons[1].AddToTabIndex()
	SetFocus(&buttons[0])

	pressTab(simScreen)
	assertFocus(test, "Tab into group", layerAlias, "r2", constants.CellTypeRadioButton)
	pressTab(simScreen)
	assertFocus(test, "Tab out of group", layerAlias, "after", constants.CellTypeButton)
	pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModNone)
	assertFocus(test, "Shift+Tab into group", layerAlias, "r2", constants.CellTypeRadioButton)

	for _, step := range []struct {
		key           tcell.Key
		expectedAlias string
	}{
		{tcell.KeyDown, "r3"}, {tcell.KeyDown, "r1"}, {tcell.KeyUp, "r3"}, {tcell.KeyRight, "r1"},
	} {
		pressKey(simScreen, step.key, 0, tcell.ModNone)
		assertFocus(test, "arrow in group", layerAlias, step.expectedAlias, constants.CellTypeRadioButton)
		if selectedAlias := radios[0].GetSelected(); selectedAlias != step.expectedAlias {
			test.Fatalf("expected %q to be selected, got %q", step.expectedAlias, selectedAlias)
		}
	}
	assertKeyboardBuffer(test, "after arrows")
}

/*
TestDefaultAndCancelButtons is a test which allows you to verify that Enter presses the layer's default button and Esc
its cancel button when the focused control does not use the key, and that controls which do use them keep them.

Example:

	Expected Inputs:
		A layer with text field "name", selector "list", dropdown "dropdown", and buttons "ok" (default) and "cancel"
		(cancel). Enter and Esc with the text field focused; Enter on the selector; Esc with the dropdown's tray open
		and then closed; Enter with the cancel button focused; Enter with the default button disabled.

	Expected Outputs:
		From the text field, Enter presses ok and Esc presses cancel. The selector picks its row instead. The open
		tray closes without a press, and Esc on the closed dropdown presses cancel. The focused cancel button wins
		over the default. With ok disabled, Enter reaches the keyboard buffer.
*/
func TestDefaultAndCancelButtons(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	layer1 := &LayerInstanceType{layerAlias: layerAlias}
	nameField := TextField.Add(layerAlias, "name", styleEntry, 2, 8, 10, 20, false, "", true)
	selector := Selector.Add(layerAlias, "list", styleEntry, func() types.SelectionEntryType {
		selectionEntry := types.NewSelectionEntry()
		selectionEntry.Add("a", "Item A")
		selectionEntry.Add("b", "Item B")
		return selectionEntry
	}(), 20, 2, 2, 10, 1, 0, 0, false, true)
	dropdownEntry := addTestDropdown(layerAlias, styleEntry)
	okButton := Button.Add(layerAlias, "ok", "ok", styleEntry, 20, 6, 8, 3, true)
	cancelButton := Button.Add(layerAlias, "cancel", "cancel", styleEntry, 30, 6, 8, 3, true)
	if err := layer1.SetDefaultButton(&okButton); err != nil {
		test.Fatalf("expected SetDefaultButton to succeed, got %v", err)
	}
	if err := layer1.SetCancelButton(&cancelButton); err != nil {
		test.Fatalf("expected SetCancelButton to succeed, got %v", err)
	}
	simScreen := startInputSimulation(test)

	SetFocus(&nameField)
	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	assertPressed(test, "Enter in text field", okButton, layerAlias, "ok")
	pressKey(simScreen, tcell.KeyEscape, 0, tcell.ModNone)
	assertPressed(test, "Esc in text field", okButton, layerAlias, "cancel")

	SetFocus(&selector)
	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	assertPressed(test, "Enter on selector", okButton, "", "")
	assertSelection(test, selector, "Enter on selector", 0, constants.SelectionSourceKeyboard)

	dropdownInstance := DropdownInstanceType{BaseControlInstanceType{layerAlias: layerAlias, controlAlias: "dropdown", controlType: constants.TYPE_DROPDOWN}}
	SetFocus(&dropdownInstance)
	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	if !dropdownEntry.IsTrayOpen {
		test.Fatalf("expected Enter to open the dropdown")
	}
	pressKey(simScreen, tcell.KeyEscape, 0, tcell.ModNone)
	if dropdownEntry.IsTrayOpen {
		test.Fatalf("expected Esc to close the tray")
	}
	assertPressed(test, "Esc on open tray", okButton, "", "")
	pressKey(simScreen, tcell.KeyEscape, 0, tcell.ModNone)
	assertPressed(test, "Esc on closed dropdown", okButton, layerAlias, "cancel")

	SetFocus(&cancelButton)
	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	assertPressed(test, "Enter on focused cancel", okButton, layerAlias, "cancel")

	okButton.SetEnabled(false)
	SetFocus(&nameField)
	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	assertPressed(test, "Enter with default disabled", okButton, "", "")
	assertKeyboardBuffer(test, "Enter with default disabled", "enter")

	otherLayer := AddLayer(0, 0, 5, 5, 5, nil)
	otherButton := Button.Add(otherLayer.GetAlias(), "other", "other", styleEntry, 0, 0, 5, 3, true)
	if err := layer1.SetDefaultButton(&otherButton); err == nil {
		test.Fatalf("expected a button from another layer to be refused as the default button")
	}
}

/*
TestFocusChangesAreQueued is a test which allows you to verify that every focus change is queued for GetFocusChange, in
order, with the control before and after the change, and that the queue is bounded.

Example:

	Expected Inputs:
		Tab order [a, b]. SetFocus(a), then Tab, then 300 alternating SetFocus calls.

	Expected Outputs:
		The first two changes are (none to a) and (a to b), followed by nothing. After the 300 calls, at most 256
		changes are queued and the last one ends on the control that has focus.
*/
func TestFocusChangesAreQueued(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b")
	buttons[0].AddToTabIndex()
	buttons[1].AddToTabIndex()

	SetFocus(&buttons[0])
	pressTab(simScreen)
	expectedChanges := []FocusChangeType{
		{PreviousControlType: constants.NullControlType, LayerAlias: layerAlias, ControlAlias: "a", ControlType: constants.CellTypeButton},
		{PreviousLayerAlias: layerAlias, PreviousControlAlias: "a", PreviousControlType: constants.CellTypeButton,
			LayerAlias: layerAlias, ControlAlias: "b", ControlType: constants.CellTypeButton},
	}
	for index, expectedChange := range expectedChanges {
		change, isFound := GetFocusChange()
		if !isFound || change != expectedChange {
			test.Fatalf("change %d: expected %+v, got %+v (found %v)", index, expectedChange, change, isFound)
		}
	}
	if change, isFound := GetFocusChange(); isFound {
		test.Fatalf("expected no further changes, got %+v", change)
	}

	for iteration := 0; iteration < 300; iteration++ {
		SetFocus(&buttons[iteration%2])
	}
	var lastChange FocusChangeType
	changeCount := 0
	for change, isFound := GetFocusChange(); isFound; change, isFound = GetFocusChange() {
		lastChange = change
		changeCount++
	}
	if changeCount > focusChangeQueueCapacity || changeCount == 0 {
		test.Fatalf("expected between 1 and %d queued changes, got %d", focusChangeQueueCapacity, changeCount)
	}
	if lastChange.ControlAlias != "b" {
		test.Fatalf("expected the last queued change to end on b, got %+v", lastChange)
	}
}

/*
TestClearTabIndexLeavesNoStaleState is a test which allows you to verify that ClearTabIndex removes every tab stop and
default button, clears focus, and discards queued focus changes, so nothing carries over to the next screen.

Example:

	Expected Inputs:
		Tab order [a, b] with b focused and a as the default button, then ClearTabIndex, Tab, and Enter.

	Expected Outputs:
		Nothing has focus and no change is queued after ClearTabIndex. Tab focuses nothing, and Enter presses nothing
		and reaches the keyboard buffer.
*/
func TestClearTabIndexLeavesNoStaleState(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b")
	layer1 := &LayerInstanceType{layerAlias: layerAlias}
	buttons[0].AddToTabIndex()
	buttons[1].AddToTabIndex()
	SetFocus(&buttons[1])
	layer1.SetDefaultButton(&buttons[0])
	readKeyboardBuffer()

	ClearTabIndex()
	assertFocus(test, "after ClearTabIndex", "", "", constants.NullControlType)
	if change, isFound := GetFocusChange(); isFound {
		test.Fatalf("expected no queued changes after ClearTabIndex, got %+v", change)
	}
	pressTab(simScreen)
	assertFocus(test, "Tab after ClearTabIndex", "", "", constants.NullControlType)
	pressKey(simScreen, tcell.KeyEnter, 0, tcell.ModNone)
	assertPressed(test, "Enter after ClearTabIndex", buttons[0], "", "")
	assertKeyboardBuffer(test, "Enter after ClearTabIndex", "enter")
}

/*
TestSetTabIndexOrdersStops is a test which allows you to verify that stops given an explicit order with SetTabIndex
come first, in ascending order, before stops in registration order, and that invalid registrations are refused.

Example:

	Expected Inputs:
		a, b, c added to the tab order in that order, then c.SetTabIndex(0) and a.SetTabIndex(1). Tab pressed three
		times with nothing focused. SetTabIndex(-1) and AddToTabIndex on a label.

	Expected Outputs:
		Focus is c, a, then b. Both invalid calls return an error.
*/
func TestSetTabIndexOrdersStops(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c")
	for index := range buttons {
		buttons[index].AddToTabIndex()
	}
	buttons[2].SetTabIndex(0)
	buttons[0].SetTabIndex(1)

	for _, expectedAlias := range []string{"c", "a", "b"} {
		pressTab(simScreen)
		assertFocus(test, "after Tab", layerAlias, expectedAlias, constants.CellTypeButton)
	}
	if err := buttons[1].SetTabIndex(-1); err == nil {
		test.Fatalf("expected a negative tab index to be refused")
	}
	label := Label.Add(layerAlias, "caption", "Caption", types.NewTuiStyleEntry(), 30, 15, 10)
	if err := label.AddToTabIndex(); err == nil {
		test.Fatalf("expected a label to be refused as a tab stop")
	}
}

/*
TestFocusedButtonIsDrawnWithFocusedColors is a test which allows you to verify that a focused button is drawn with its
style's focused colours, and that a style that leaves them unset still shows focus by swapping the normal colours.

Example:

	Expected Inputs:
		Button a, rendered unfocused and then focused with the default style. Button b, whose style has no focused
		colours, rendered focused.

	Expected Outputs:
		Unfocused, a's label cell has the normal button colours. Focused, it has FocusedForegroundColor and
		FocusedBackgroundColor. Focused b has its normal foreground and background swapped.
*/
func TestFocusedButtonIsDrawnWithFocusedColors(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	focusButton := Button.Add(layerAlias, "a", "a", styleEntry, 1, 1, 8, 3, true)
	plainStyleEntry := styleEntry
	plainStyleEntry.Button.FocusedForegroundColor = 0
	plainStyleEntry.Button.FocusedBackgroundColor = 0
	plainButton := Button.Add(layerAlias, "b", "b", plainStyleEntry, 1, 5, 8, 3, true)
	getLabelColors := func(yLocation int) (constants.ColorType, constants.ColorType) {
		UpdateDisplay(false)
		attributeEntry := commonResource.screenLayer.CharacterMemory[yLocation][4].AttributeEntry
		return attributeEntry.ForegroundColor, attributeEntry.BackgroundColor
	}

	foregroundColor, backgroundColor := getLabelColors(2)
	if foregroundColor != styleEntry.Button.ForegroundColor || backgroundColor != styleEntry.Button.BackgroundColor {
		test.Fatalf("expected the unfocused button to use its normal colours")
	}
	SetFocus(&focusButton)
	foregroundColor, backgroundColor = getLabelColors(2)
	if foregroundColor != styleEntry.Button.FocusedForegroundColor || backgroundColor != styleEntry.Button.FocusedBackgroundColor {
		test.Fatalf("expected the focused button to use its focused colours")
	}
	SetFocus(&plainButton)
	foregroundColor, backgroundColor = getLabelColors(6)
	if foregroundColor != styleEntry.Button.BackgroundColor || backgroundColor != styleEntry.Button.ForegroundColor {
		test.Fatalf("expected a style without focused colours to show focus by swapping the normal colours")
	}
}

/*
TestGetFocusedColorsKeepsFocusVisible is a test which allows you to verify each fallback rule of getFocusedColors.

Example:

	Expected Inputs:
		Normal colours (white, black) with focused colours (blue, yellow), (unset, yellow), (unset, unset), and
		(white, black).

	Expected Outputs:
		(blue, yellow), (white, yellow), (black, white), and (black, white).
*/
func TestGetFocusedColorsKeepsFocusVisible(test *testing.T) {
	white := constants.AnsiColorByIndex[15]
	black := constants.AnsiColorByIndex[0]
	blue := constants.AnsiColorByIndex[4]
	yellow := constants.AnsiColorByIndex[11]
	testCases := []struct {
		name                       string
		focusedForeground          constants.ColorType
		focusedBackground          constants.ColorType
		expectedForeground         constants.ColorType
		expectedBackgroundExpected constants.ColorType
	}{
		{"both set", blue, yellow, blue, yellow},
		{"foreground unset", 0, yellow, white, yellow},
		{"both unset", 0, 0, black, white},
		{"same as normal", white, black, black, white},
	}
	for _, testCase := range testCases {
		foregroundColor, backgroundColor := getFocusedColors(white, black, testCase.focusedForeground, testCase.focusedBackground)
		if foregroundColor != testCase.expectedForeground || backgroundColor != testCase.expectedBackgroundExpected {
			test.Fatalf("%s: unexpected colours", testCase.name)
		}
	}
}

/*
TestClickFocusRulesAndHoverKeepFocus is a test which allows you to verify that a click focuses an enabled control, that
a click on a disabled control or empty space leaves focus alone, and that hovering over a selector never takes focus.

Example:

	Expected Inputs:
		Text field "name", disabled button "off", and a selector. Click the text field, click the disabled button,
		click empty space, then move the mouse over the selector without pressing.

	Expected Outputs:
		Focus is the text field after the first click and stays on it after every later step.
*/
func TestClickFocusRulesAndHoverKeepFocus(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	TextField.Add(layerAlias, "name", styleEntry, 2, 8, 10, 20, false, "", true)
	Button.Add(layerAlias, "off", "off", styleEntry, 20, 2, 8, 3, false)
	addTestSelector(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)

	clickAt(simScreen, 3, 8)
	assertFocus(test, "click text field", layerAlias, "name", constants.CellTypeTextField)
	clickAt(simScreen, 23, 3)
	assertFocus(test, "click disabled button", layerAlias, "name", constants.CellTypeTextField)
	clickAt(simScreen, 35, 18)
	assertFocus(test, "click empty space", layerAlias, "name", constants.CellTypeTextField)
	resetMouseEventState()
	simScreen.InjectMouse(4, 3, tcell.ButtonNone, tcell.ModNone)
	UpdateEventQueues()
	assertFocus(test, "hover selector", layerAlias, "name", constants.CellTypeTextField)
}

/*
TestTabAfterRebuiltOrderContinuesFromFocusedControl is a test which allows you to verify that after the tab order is
cleared and rebuilt, Tab continues from the control that has focus in the new order rather than from the stale position
left by the previous order. In addition, the following should be noted:

  - Before the fix, ClearTabIndex left the stored position at 1 from the old order, so the first Tab in the new
    two-entry order wrapped back to entry 0 (C) instead of moving on to D.

Example:

	Expected Inputs:
		Tab order [A, B], one Tab press, ClearTabIndex, tab order [C, D] with C focused via GetFocus, one Tab press.

	Expected Outputs:
		Focus is B after the first Tab press and D after the second.
*/
func TestTabAfterRebuiltOrderContinuesFromFocusedControl(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c", "d")
	buttons[0].AddToTabIndex()
	buttons[1].AddToTabIndex()
	buttons[0].GetFocus()

	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "b", "first order, after Tab")

	ClearTabIndex()
	buttons[2].AddToTabIndex()
	buttons[3].AddToTabIndex()
	buttons[2].GetFocus()

	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "d", "rebuilt order, after Tab")
}

/*
TestTabAfterRebuiltOrderWithoutFocusStartsAtFirstEntry is a test which allows you to verify that the first Tab after the
tab order is rebuilt focuses the first entry of the new order when the focused control is not part of it.

Example:

	Expected Inputs:
		Tab order [A, B] with nothing focused and three Tab presses (A, B, A), then ClearTabIndex and tab order [C, D].

	Expected Outputs:
		Focus is A after the three Tab presses and C after the next Tab press.
*/
func TestTabAfterRebuiltOrderWithoutFocusStartsAtFirstEntry(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c", "d")
	buttons[0].AddToTabIndex()
	buttons[1].AddToTabIndex()

	pressTab(simScreen)
	pressTab(simScreen)
	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "a", "first order, after three Tabs")

	ClearTabIndex()
	buttons[2].AddToTabIndex()
	buttons[3].AddToTabIndex()

	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "c", "rebuilt order, after Tab")
}

/*
TestTabContinuesFromFocusSetDirectly is a test which allows you to verify that when focus is moved with GetFocus instead
of Tab, the next Tab press continues from that control and wraps around at the end of the order.

Example:

	Expected Inputs:
		Tab order [A, B, C], B focused via GetFocus, then two Tab presses.

	Expected Outputs:
		Focus is C after the first Tab press and A after the second.
*/
func TestTabContinuesFromFocusSetDirectly(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c")
	for index := range buttons {
		buttons[index].AddToTabIndex()
	}
	buttons[1].GetFocus()

	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "c", "after first Tab")

	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "a", "after second Tab")
}

/*
TestTabWithFocusOutsideOrderFocusesFirstEntry is a test which allows you to verify that when the focused control is not
part of the tab order, Tab focuses the first entry of the order.

Example:

	Expected Inputs:
		Tab order [A, B], with C (not in the order) focused via GetFocus, then one Tab press.

	Expected Outputs:
		Focus is A.
*/
func TestTabWithFocusOutsideOrderFocusesFirstEntry(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c")
	buttons[0].AddToTabIndex()
	buttons[1].AddToTabIndex()
	buttons[2].GetFocus()

	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "a", "after Tab")
}

/*
TestTabWithEmptyOrderLeavesFocusUnchanged is a test which allows you to verify that pressing Tab with no registered tab
order neither panics nor changes which control has focus.

Example:

	Expected Inputs:
		An empty tab order, with A focused via GetFocus, then two Tab presses.

	Expected Outputs:
		Focus remains on A.
*/
func TestTabWithEmptyOrderLeavesFocusUnchanged(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a")
	buttons[0].GetFocus()

	pressTab(simScreen)
	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "a", "after Tab with empty order")
}

/*
TestTabWithSingleEntryKeepsFocusOnIt is a test which allows you to verify that with a single entry in the tab order, Tab
focuses that entry and keeps focus on it across repeated presses.

Example:

	Expected Inputs:
		Tab order [A] with nothing focused, then three Tab presses.

	Expected Outputs:
		Focus is A after every press.
*/
func TestTabWithSingleEntryKeepsFocusOnIt(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a")
	buttons[0].AddToTabIndex()

	for pressCount := 1; pressCount <= 3; pressCount++ {
		pressTab(simScreen)
		assertButtonFocused(test, layerAlias, "a", "after Tab")
	}
}

/*
TestTabContinuesFromFocusSetDirectlyAfterTabbing is a test which allows you to verify that when focus is moved with
GetFocus after the stored tab position has already been advanced by Tab, the next Tab press continues from the directly
focused control rather than from the stored position. In addition, the following should be noted:

  - Before the fix, the stored position stayed on C after the Tab presses, so the next Tab wrapped to A even though
    B had focus.

Example:

	Expected Inputs:
		Tab order [A, B, C], three Tab presses ending on C, B focused via GetFocus, then one Tab press.

	Expected Outputs:
		Focus is C after the three Tab presses and C again after the final Tab press.
*/
func TestTabContinuesFromFocusSetDirectlyAfterTabbing(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c")
	for index := range buttons {
		buttons[index].AddToTabIndex()
	}

	pressTab(simScreen)
	pressTab(simScreen)
	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "c", "after three Tabs")

	buttons[1].GetFocus()
	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "c", "after Tab from directly focused B")
}

/*
TestShiftTabMovesFocusToPreviousControl is a test which allows you to verify that Shift+Tab, reported by terminals as
Backtab, moves focus to the previous control in the tab order, wraps from the first entry to the last, and is not added
to the keyboard buffer.

Example:

	Expected Inputs:
		Tab order [A, B, C], B focused via GetFocus, then two Backtab presses.

	Expected Outputs:
		Focus is A after the first press and C after the second, and the keyboard buffer is empty.
*/
func TestShiftTabMovesFocusToPreviousControl(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c")
	readKeyboardBuffer()
	for index := range buttons {
		buttons[index].AddToTabIndex()
	}
	buttons[1].GetFocus()

	pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModNone)
	assertButtonFocused(test, layerAlias, "a", "after first Shift+Tab")

	pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModNone)
	assertButtonFocused(test, layerAlias, "c", "after second Shift+Tab")
	assertKeyboardBuffer(test, "after Shift+Tab")
}

/*
TestShiftTabVariantsMoveFocusBackward is a test which allows you to verify that every way tcell can report Shift+Tab,
namely Backtab, Shift+Backtab, and Tab with the Shift modifier, moves focus to the previous control.

Example:

	Expected Inputs:
		Tab order [A, B, C, D], D focused via GetFocus, then Backtab, Backtab with Shift, and Tab with Shift.

	Expected Outputs:
		Focus is C, then B, then A.
*/
func TestShiftTabVariantsMoveFocusBackward(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c", "d")
	for index := range buttons {
		buttons[index].AddToTabIndex()
	}
	buttons[3].GetFocus()

	pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModNone)
	assertButtonFocused(test, layerAlias, "c", "after Backtab")

	pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModShift)
	assertButtonFocused(test, layerAlias, "b", "after Shift+Backtab")

	pressKey(simScreen, tcell.KeyTab, 0, tcell.ModShift)
	assertButtonFocused(test, layerAlias, "a", "after Shift+Tab")
}

/*
TestShiftTabWithFocusOutsideOrderFocusesLastEntry is a test which allows you to verify that when the focused control is
not part of the tab order, Shift+Tab focuses the last entry of the order.

Example:

	Expected Inputs:
		Tab order [A, B], with C (not in the order) focused via GetFocus, then one Backtab press.

	Expected Outputs:
		Focus is B.
*/
func TestShiftTabWithFocusOutsideOrderFocusesLastEntry(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c")
	buttons[0].AddToTabIndex()
	buttons[1].AddToTabIndex()
	buttons[2].GetFocus()

	pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModNone)
	assertButtonFocused(test, layerAlias, "b", "after Shift+Tab")
}

/*
TestTabSkipsDisabledHiddenAndDeletedControls is a test which allows you to verify that Tab and Shift+Tab skip controls
that are disabled, hidden, or have been deleted since they were registered in the tab order.

Example:

	Expected Inputs:
		Tab order [A, B, C, D, E] with B disabled, C hidden, and D deleted. A focused via GetFocus, then one Tab press
		and one Backtab press.

	Expected Outputs:
		Focus is E after Tab and A after Backtab.
*/
func TestTabSkipsDisabledHiddenAndDeletedControls(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c", "d", "e")
	for index := range buttons {
		buttons[index].AddToTabIndex()
	}
	buttons[1].SetEnabled(false)
	buttons[2].SetVisible(false)
	buttons[3].Delete()
	buttons[0].GetFocus()

	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "e", "after Tab")

	pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModNone)
	assertButtonFocused(test, layerAlias, "a", "after Shift+Tab")
}

/*
TestTabWithNoFocusableEntriesLeavesFocusUnchanged is a test which allows you to verify that when every entry in the tab
order is disabled, Tab and Shift+Tab leave focus where it was.

Example:

	Expected Inputs:
		Tab order [A, B] with both disabled, C focused via GetFocus, then one Tab and one Backtab press.

	Expected Outputs:
		Focus remains on C.
*/
func TestTabWithNoFocusableEntriesLeavesFocusUnchanged(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c")
	buttons[0].AddToTabIndex()
	buttons[1].AddToTabIndex()
	buttons[0].SetEnabled(false)
	buttons[1].SetEnabled(false)
	buttons[2].GetFocus()

	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "c", "after Tab")

	pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModNone)
	assertButtonFocused(test, layerAlias, "c", "after Shift+Tab")
}

/*
TestTextFieldAddToTabIndex is a test which allows you to verify that a text field registered with AddToTabIndex takes
part in Tab and Shift+Tab navigation, and is skipped while disabled.

Example:

	Expected Inputs:
		Tab order [button A, text field F, button B], A focused via GetFocus. Tab, then Backtab, then F disabled and
		Tab again.

	Expected Outputs:
		Focus is F after Tab, A after Backtab, and B after the final Tab.
*/
func TestTextFieldAddToTabIndex(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b")
	textField := TextField.Add(layerAlias, "field", types.NewTuiStyleEntry(), 12, 1, 10, 20, false, "", true)
	buttons[0].AddToTabIndex()
	textField.AddToTabIndex()
	buttons[1].AddToTabIndex()
	buttons[0].GetFocus()

	pressTab(simScreen)
	if !isControlCurrentlyFocused(layerAlias, "field", constants.CellTypeTextField) {
		test.Fatalf("expected the text field to be focused after Tab")
	}

	pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModNone)
	assertButtonFocused(test, layerAlias, "a", "after Shift+Tab")

	textField.SetEnabled(false)
	pressTab(simScreen)
	assertButtonFocused(test, layerAlias, "b", "after Tab past the disabled text field")
}

/*
TestGetFocusedControl is a test which allows you to verify that GetFocusedControl reports the control that has focus
as it changes through GetFocus and Tab, and reports no control when nothing is focused.

Example:

	Expected Inputs:
		Nothing focused, then button A focused via GetFocus, then Tab to text field F, then focus on an empty layer
		cell (layer alias set, control alias empty).

	Expected Outputs:
		("", "", NullControlType), then (layer, "a", CellTypeButton), then (layer, "field", CellTypeTextField), then
		("", "", NullControlType).
*/
func TestGetFocusedControl(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a")
	textField := TextField.Add(layerAlias, "field", types.NewTuiStyleEntry(), 12, 1, 10, 20, false, "", true)
	buttons[0].AddToTabIndex()
	textField.AddToTabIndex()

	assertFocusedControl := func(context string, expectedLayerAlias string, expectedControlAlias string, expectedControlType int) {
		test.Helper()
		focusedLayerAlias, focusedControlAlias, focusedControlType := GetFocusedControl()
		if focusedLayerAlias != expectedLayerAlias || focusedControlAlias != expectedControlAlias || focusedControlType != expectedControlType {
			test.Fatalf("%s: expected (%q, %q, %d), got (%q, %q, %d)", context, expectedLayerAlias, expectedControlAlias,
				expectedControlType, focusedLayerAlias, focusedControlAlias, focusedControlType)
		}
	}

	assertFocusedControl("nothing focused", "", "", constants.NullControlType)

	buttons[0].GetFocus()
	assertFocusedControl("after GetFocus", layerAlias, "a", constants.CellTypeButton)

	pressTab(simScreen)
	assertFocusedControl("after Tab", layerAlias, "field", constants.CellTypeTextField)

	setFocusedControl(layerAlias, "", constants.NullCellType)
	assertFocusedControl("empty layer cell focused", "", "", constants.NullControlType)
}

/*
TestFocusStateConcurrentAccessIsConsistent is a test which allows you to verify that focus, tab order, drag state, and
modifier keys can be read and written from the application goroutine while the event goroutine changes them, without
ever observing a torn focus value. In addition, the following should be noted:

  - Run with -race to also check for data races. Before the fix, every one of these fields was read and written
    with no lock.

Example:

	Expected Inputs:
		One goroutine repeatedly tabbing, setting focus, drag state, highlight, and modifier keys, while the test
		goroutine repeatedly calls GetFocus, GetFocusedControl, AddToTabIndex, ClearTabIndex, and IsShiftPressed.

	Expected Outputs:
		Every GetFocusedControl result is ("", "", NullControlType), (layer, "a", CellTypeButton), or
		(layer, "field", CellTypeTextField).
*/
func TestFocusStateConcurrentAccessIsConsistent(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	buttonInstance := Button.Add(layerAlias, "a", "a", styleEntry, 1, 1, 8, 3, true)
	textField := TextField.Add(layerAlias, "field", styleEntry, 12, 1, 10, 20, false, "", true)

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			default:
				FocusNext()
				FocusPrevious()
				setFocusedControl(layerAlias, "field", constants.CellTypeTextField)
				setEventStateId(constants.EventStateDragAndDrop)
				setEventStateId(constants.EventStateNone)
				setPreviouslyHighlightedControl(layerAlias, "a", constants.CellTypeButton)
				setModifierKeys(tcell.ModShift)
			}
		}
	}()

	for iteration := 0; iteration < 2000; iteration++ {
		buttonInstance.GetFocus()
		buttonInstance.AddToTabIndex()
		textField.AddToTabIndex()
		focusedLayerAlias, focusedControlAlias, focusedControlType := GetFocusedControl()
		isConsistent := (focusedControlAlias == "" && focusedLayerAlias == "" && focusedControlType == constants.NullControlType) ||
			(focusedLayerAlias == layerAlias && focusedControlAlias == "a" && focusedControlType == constants.CellTypeButton) ||
			(focusedLayerAlias == layerAlias && focusedControlAlias == "field" && focusedControlType == constants.CellTypeTextField)
		if !isConsistent {
			close(stop)
			<-done
			test.Fatalf("torn focus value (%q, %q, %d)", focusedLayerAlias, focusedControlAlias, focusedControlType)
		}
		IsShiftPressed()
		getPreviouslyHighlightedControl()
		getEventStateId()
		ClearTabIndex()
	}
	close(stop)
	<-done
}

/*
getRenderedCellColors is a method which allows you to redraw the screen and obtain the foreground and background colours
of the screen cell at a given location.

Example:

	foregroundColor, backgroundColor := getRenderedCellColors(4, 2)
*/
func getRenderedCellColors(xLocation int, yLocation int) (constants.ColorType, constants.ColorType) {
	UpdateDisplay(false)
	attributeEntry := commonResource.screenLayer.CharacterMemory[yLocation][xLocation].AttributeEntry
	return attributeEntry.ForegroundColor, attributeEntry.BackgroundColor
}

/*
TestFocusIndicatorFollowsInputMethod is a test which allows you to verify that clicking a control focuses it without
drawing the keyboard focus colours, and that a keystroke then shows them, for every control that draws focus colours.

Example:

	Expected Inputs:
		One case per control: a button, a checkbox, a radio button, a dropdown and a bordered selector, each clicked
		once, then the key "x" is pressed, then the control is clicked again. A dropdown's tray is closed after each
		click. The cell sampled is the clicked cell, except for the selector, whose top left border cell is sampled.

	Expected Outputs:
		After each click the control has focus and the sampled cell has the control's normal colours. After "x" the
		same cell has the colours getFocusedColors resolves from the style's focused colours.
*/
func TestFocusIndicatorFollowsInputMethod(test *testing.T) {
	type styleColorsFunc func(styleEntry types.TuiStyleEntryType) (constants.ColorType, constants.ColorType, constants.ColorType, constants.ColorType)
	testCases := []struct {
		name            string
		addControl      func(layerAlias string, styleEntry types.TuiStyleEntryType)
		controlType     int
		clickXLocation  int
		clickYLocation  int
		sampleXLocation int
		sampleYLocation int
		getStyleColors  styleColorsFunc
	}{
		{
			name: "button",
			addControl: func(layerAlias string, styleEntry types.TuiStyleEntryType) {
				Button.Add(layerAlias, "control", "control", styleEntry, 1, 1, 12, 3, true)
			},
			controlType:    constants.CellTypeButton,
			clickXLocation: 4, clickYLocation: 2, sampleXLocation: 4, sampleYLocation: 2,
			getStyleColors: func(styleEntry types.TuiStyleEntryType) (constants.ColorType, constants.ColorType, constants.ColorType, constants.ColorType) {
				return styleEntry.Button.ForegroundColor, styleEntry.Button.BackgroundColor, styleEntry.Button.FocusedForegroundColor, styleEntry.Button.FocusedBackgroundColor
			},
		},
		{
			name: "checkbox",
			addControl: func(layerAlias string, styleEntry types.TuiStyleEntryType) {
				Checkbox.Add(layerAlias, "control", "control", styleEntry, 2, 2, false, true)
			},
			controlType:    constants.CellTypeCheckbox,
			clickXLocation: 2, clickYLocation: 2, sampleXLocation: 2, sampleYLocation: 2,
			getStyleColors: func(styleEntry types.TuiStyleEntryType) (constants.ColorType, constants.ColorType, constants.ColorType, constants.ColorType) {
				return styleEntry.Checkbox.ForegroundColor, styleEntry.Checkbox.BackgroundColor, styleEntry.Checkbox.FocusedForegroundColor, styleEntry.Checkbox.FocusedBackgroundColor
			},
		},
		{
			name: "radio button",
			addControl: func(layerAlias string, styleEntry types.TuiStyleEntryType) {
				radioButton.Add(layerAlias, "control", "control", styleEntry, 2, 2, 1, false)
			},
			controlType:    constants.CellTypeRadioButton,
			clickXLocation: 2, clickYLocation: 2, sampleXLocation: 2, sampleYLocation: 2,
			getStyleColors: func(styleEntry types.TuiStyleEntryType) (constants.ColorType, constants.ColorType, constants.ColorType, constants.ColorType) {
				return styleEntry.RadioButton.ForegroundColor, styleEntry.RadioButton.BackgroundColor, styleEntry.RadioButton.FocusedForegroundColor, styleEntry.RadioButton.FocusedBackgroundColor
			},
		},
		{
			name: "dropdown",
			addControl: func(layerAlias string, styleEntry types.TuiStyleEntryType) {
				Dropdown.Add(layerAlias, "control", styleEntry, getTestSelectionEntry(4), 2, 2, 3, 10, 0)
			},
			controlType:    constants.CellTypeDropdown,
			clickXLocation: 3, clickYLocation: 2, sampleXLocation: 3, sampleYLocation: 2,
			getStyleColors: func(styleEntry types.TuiStyleEntryType) (constants.ColorType, constants.ColorType, constants.ColorType, constants.ColorType) {
				return styleEntry.Dropdown.ForegroundColor, styleEntry.Dropdown.BackgroundColor, styleEntry.Dropdown.FocusedForegroundColor, styleEntry.Dropdown.FocusedBackgroundColor
			},
		},
		{
			name: "selector border",
			addControl: func(layerAlias string, styleEntry types.TuiStyleEntryType) {
				Selector.Add(layerAlias, "control", styleEntry, getTestSelectionEntry(4), 3, 3, 4, 10, 1, 0, 0, false, true)
			},
			controlType:    constants.CellTypeSelectorItem,
			clickXLocation: 3, clickYLocation: 3, sampleXLocation: 2, sampleYLocation: 2,
			getStyleColors: func(styleEntry types.TuiStyleEntryType) (constants.ColorType, constants.ColorType, constants.ColorType, constants.ColorType) {
				return styleEntry.Window.LineDrawingTextForegroundColor, styleEntry.Window.LineDrawingTextBackgroundColor, styleEntry.Selector.FocusedForegroundColor, styleEntry.Selector.FocusedBackgroundColor
			},
		},
	}
	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			layerAlias, styleEntry := setupInputTest(test)
			testCase.addControl(layerAlias, styleEntry)
			simScreen := startInputSimulation(test)
			normalForeground, normalBackground, focusedForeground, focusedBackground := testCase.getStyleColors(styleEntry)
			expectedFocusedForeground, expectedFocusedBackground := getFocusedColors(normalForeground, normalBackground, focusedForeground, focusedBackground)
			assertColors := func(context string, expectedForeground constants.ColorType, expectedBackground constants.ColorType) {
				test.Helper()
				foregroundColor, backgroundColor := getRenderedCellColors(testCase.sampleXLocation, testCase.sampleYLocation)
				if foregroundColor != expectedForeground || backgroundColor != expectedBackground {
					test.Fatalf("%s: expected colours (%d, %d), got (%d, %d)", context, expectedForeground, expectedBackground, foregroundColor, backgroundColor)
				}
			}
			clickAndCloseTray := func(context string) {
				test.Helper()
				// A click opens a dropdown's tray, whose border covers the sampled cell, so any tray is closed first.
				clickAt(simScreen, testCase.clickXLocation, testCase.clickYLocation)
				Dropdown.closeAllOpen()
				assertFocus(test, context, layerAlias, "control", testCase.controlType)
			}

			clickAndCloseTray("after the first click")
			assertColors("after the first click", normalForeground, normalBackground)
			pressKey(simScreen, tcell.KeyRune, 'x', tcell.ModNone)
			assertColors("after a keystroke", expectedFocusedForeground, expectedFocusedBackground)
			clickAndCloseTray("after the second click")
			assertColors("after the second click", normalForeground, normalBackground)
		})
	}
}

/*
TestFocusIndicatorIgnoresMovementAndProgrammaticFocus is a test which allows you to verify which input changes whether
the keyboard focus indicator is shown: Tab shows it, mouse movement and the wheel leave it alone, a mouse press hides
it, and focus moved by the application keeps whatever the user's last input decided.

Example:

	Expected Inputs:
		Buttons a and b in the tab order. Tab onto a, move the mouse over b with no button held, scroll the wheel,
		press and release the mouse on empty space, call SetFocus on b, then press Shift+Tab.

	Expected Outputs:
		After Tab, a's label cell has the focused colours, and still does after the movement and the wheel. After the
		press on empty space, a is drawn in its normal colours. After SetFocus, b has focus but is drawn in its
		normal colours. After Shift+Tab, a has focus and is drawn in the focused colours again.
*/
func TestFocusIndicatorIgnoresMovementAndProgrammaticFocus(test *testing.T) {
	layerAlias, buttons, simScreen := setupTabOrderTest(test, "a", "b")
	for index := range buttons {
		if err := buttons[index].AddToTabIndex(); err != nil {
			test.Fatalf("expected AddToTabIndex to succeed, got %v", err)
		}
	}
	styleEntry := Buttons.Get(layerAlias, "a").StyleEntry
	focusedForeground, focusedBackground := getFocusedColors(styleEntry.Button.ForegroundColor, styleEntry.Button.BackgroundColor,
		styleEntry.Button.FocusedForegroundColor, styleEntry.Button.FocusedBackgroundColor)
	assertLabelColors := func(context string, yLocation int, isFocusShown bool) {
		test.Helper()
		expectedForeground, expectedBackground := styleEntry.Button.ForegroundColor, styleEntry.Button.BackgroundColor
		if isFocusShown {
			expectedForeground, expectedBackground = focusedForeground, focusedBackground
		}
		foregroundColor, backgroundColor := getRenderedCellColors(4, yLocation)
		if foregroundColor != expectedForeground || backgroundColor != expectedBackground {
			test.Fatalf("%s: expected colours (%d, %d), got (%d, %d)", context, expectedForeground, expectedBackground, foregroundColor, backgroundColor)
		}
	}

	pressTab(simScreen)
	assertFocus(test, "after Tab", layerAlias, "a", constants.CellTypeButton)
	assertLabelColors("after Tab", 2, true)
	simScreen.InjectMouse(4, 5, tcell.ButtonNone, tcell.ModNone)
	UpdateEventQueues()
	simScreen.InjectMouse(4, 5, tcell.WheelDown, tcell.ModNone)
	UpdateEventQueues()
	assertLabelColors("after mouse movement and the wheel", 2, true)
	clickAt(simScreen, 30, 2)
	assertLabelColors("after a press on empty space", 2, false)
	if err := SetFocus(&buttons[1]); err != nil {
		test.Fatalf("expected SetFocus to succeed, got %v", err)
	}
	assertFocus(test, "after SetFocus", layerAlias, "b", constants.CellTypeButton)
	assertLabelColors("after SetFocus following a click", 5, false)
	pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModShift)
	assertFocus(test, "after Shift+Tab", layerAlias, "a", constants.CellTypeButton)
	assertLabelColors("after Shift+Tab", 2, true)
}
