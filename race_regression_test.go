package consolizer

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

/*
TestButtonDeleteDuringMouseEventDoesNotPanic is a test which allows you to verify that deleting a button on one
goroutine while consolizer's mouse event goroutine is reading it on another goroutine never panics. In addition, the
following should be noted:

  - This is a regression test for a check-then-get race in buttonType.updateStateMouse: IsExists and Get were two
    separate lookups, so a button removed between them made Get panic with "could not be obtained since it does not
    exist". The fix collapses both lookups into a single ControlMemoryManager.Lookup call.

  - The test is bounded to a few seconds rather than left to run indefinitely, since the unpatched code was observed
    to panic in under 1.5 seconds.

Example:

	go test -run TestButtonDeleteDuringMouseEventDoesNotPanic -count=1 .
*/
func TestButtonDeleteDuringMouseEventDoesNotPanic(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layerAlias := layer1.GetAlias()

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			default:
				Button.updateStates(true)
			}
		}
	}()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		Button.Add(layerAlias, "back", "BACK", styleEntry, 2, 2, 10, 3, true)
		UpdateDisplay(false)
		SetMouseStatus(5, 3, 0, "")
		layer1.DeleteAllButtons()
	}
	close(stop)
	<-done
}

/*
TestTooltipDeleteDuringLookupDoesNotPanic is a test which allows you to verify that deleting a button's tooltip on one
goroutine while another goroutine resolves that same tooltip through the button never panics. In addition, the
following should be noted:

  - This is a regression test for the same class of check-then-get race as
    TestButtonDeleteDuringMouseEventDoesNotPanic, but in the chained button-to-tooltip lookup inside
    tooltipType.getFromCharacterEntry, which runs on the periodic event goroutine roughly every 500ms in normal
    operation. The fix collapses each IsExists-then-Get pair into a single Lookup call.

  - The button itself is created once and left alone; only its tooltip is repeatedly removed and re-added under the
    same alias, so the race is isolated to the tooltip registry lookup rather than the already-covered button
    registry lookup.

  - The test drives getFromCharacterEntry directly with a synthetic character entry rather than going through mouse
    events, since that is the exact call the periodic goroutine makes.

Example:

	go test -run TestTooltipDeleteDuringLookupDoesNotPanic -count=1 .
*/
func TestTooltipDeleteDuringLookupDoesNotPanic(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layerAlias := layer1.GetAlias()

	buttonInstance := Button.Add(layerAlias, "back", "BACK", styleEntry, 2, 2, 10, 3, true)
	buttonEntry := buttonRegistry.Get(layerAlias, buttonInstance.GetAlias())
	tooltipAlias := buttonEntry.TooltipAlias

	characterEntry := types.NewCharacterEntry()
	characterEntry.LayerAlias = layerAlias
	characterEntry.AttributeEntry.CellType = constants.CellTypeButton
	characterEntry.AttributeEntry.CellControlAlias = "back"

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			default:
				Tooltip.getFromCharacterEntry(characterEntry)
			}
		}
	}()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		Tooltip.Add(layerAlias, tooltipAlias, "", styleEntry, 2, 2, 10, 3, 2, 6, 10, 3, false, true, constants.DefaultTooltipHoverTime)
		Tooltip.Delete(layerAlias, tooltipAlias)
	}
	close(stop)
	<-done
}

/*
TestButtonDeleteClearsStaleScreenReference is a test which allows you to verify that deleting a button immediately
resets its former cell's hit-testing metadata on the composited screen buffer, instead of leaving it to point at the
deleted alias until the next redraw. In addition, the following should be noted:

  - This is a regression test for the alias-reuse gap in the original fix: without clearing the stale cell at delete
    time, a new control created with the same alias and type before the next UpdateDisplay call would be reachable
    through the old, unredrawn cell location, even though that cell visually still shows the old control.

  - The check reads the composited buffer directly through getCellInformationUnderMouseCursor, the same call mouse
    dispatch makes, rather than through the memory manager, since the bug is specifically about what that buffer
    still reports after a deletion.

Example:

	go test -run TestButtonDeleteClearsStaleScreenReference -count=1 .
*/
func TestButtonDeleteClearsStaleScreenReference(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layerAlias := layer1.GetAlias()

	Button.Add(layerAlias, "back", "BACK", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)

	characterEntry := getCellInformationUnderMouseCursor(2, 2)
	if characterEntry.AttributeEntry.CellControlAlias != "back" || characterEntry.AttributeEntry.CellType != constants.CellTypeButton {
		test.Fatalf("expected the composited screen buffer to carry the button's alias before deletion, got alias %q and cell type %d",
			characterEntry.AttributeEntry.CellControlAlias, characterEntry.AttributeEntry.CellType)
	}

	layer1.DeleteAllButtons()

	characterEntry = getCellInformationUnderMouseCursor(2, 2)
	if characterEntry.AttributeEntry.CellControlAlias != "" {
		test.Fatalf("expected the deleted button's alias to be cleared from the screen buffer immediately, but it still reads %q",
			characterEntry.AttributeEntry.CellControlAlias)
	}
	if characterEntry.AttributeEntry.CellType != constants.NullCellType {
		test.Fatalf("expected the deleted button's cell type to be reset to NullCellType, got %d", characterEntry.AttributeEntry.CellType)
	}

	// Reusing the alias for a brand new button, without an intervening redraw, must not make the stale, unredrawn
	// cell resolve to it.
	Button.Add(layerAlias, "back", "BACK", styleEntry, 20, 10, 10, 3, true)
	characterEntry = getCellInformationUnderMouseCursor(2, 2)
	if characterEntry.AttributeEntry.CellControlAlias == "back" {
		test.Fatalf("a new button reusing the deleted alias must not be reachable through the old, unredrawn cell location")
	}
}

/*
TestClearStaleControlReferencesDoesNotModifyPublishedFrame is a test which allows you to verify that clearing a deleted
control's cells publishes a new frame instead of modifying the one readers may already hold, and that rows without a
matching cell are shared rather than copied. In addition, the following should be noted:

  - This is a regression test for a data race between clearStaleControlReferences, which rewrote the published
    frame's cells in place, and readers such as getCellInformationUnderMouseCursor, which copy the frame under a
    read lock and then read its cells after releasing it. The copy shares cell storage with the published frame, so
    the in-place write raced those reads. Before the fix, the snapshot taken below saw its cells change.

Example:

	Expected Inputs:
		Button "back" at (2, 2) on layer 1, displayed. A snapshot of the frame is taken under a read lock, then
		clearStaleControlReferences(layer 1, "back", CellTypeButton) runs.

	Expected Outputs:
		The snapshot's cell at (2, 2) still has alias "back" and CellTypeButton. The published frame's cell at (2, 2)
		has an empty alias and NullCellType. Row 0, which holds no button cell, is the same row in both frames.
*/
func TestClearStaleControlReferencesDoesNotModifyPublishedFrame(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layerAlias := layer1.GetAlias()
	Button.Add(layerAlias, "back", "BACK", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)

	commonResource.displayUpdate.RLock()
	snapshot := commonResource.screenLayer
	commonResource.displayUpdate.RUnlock()

	clearStaleControlReferences(layerAlias, "back", constants.CellTypeButton)

	snapshotCell := snapshot.CharacterMemory[2][2].AttributeEntry
	if snapshotCell.CellControlAlias != "back" || snapshotCell.CellType != constants.CellTypeButton {
		test.Fatalf("expected the earlier snapshot to be left unchanged, got alias %q and cell type %d",
			snapshotCell.CellControlAlias, snapshotCell.CellType)
	}

	commonResource.displayUpdate.RLock()
	published := commonResource.screenLayer
	commonResource.displayUpdate.RUnlock()
	publishedCell := published.CharacterMemory[2][2].AttributeEntry
	if publishedCell.CellControlAlias != "" || publishedCell.CellType != constants.NullCellType {
		test.Fatalf("expected the published frame to have the button cleared, got alias %q and cell type %d",
			publishedCell.CellControlAlias, publishedCell.CellType)
	}
	if &snapshot.CharacterMemory[0][0] != &published.CharacterMemory[0][0] {
		test.Fatalf("expected row 0, which has no button cell, to be shared between the old and new frames")
	}
}

/*
TestTooltipHoverStateConcurrentUpdatesAreSafe is a test which allows you to verify that a tooltip's hover state can be
updated by the event goroutine and the periodic-event goroutine at the same time, while UpdateDisplay renders it on a
third goroutine, and that the hover state machine still ends in the right state. In addition, the following should be
noted:

  - This is a regression test for a data race on the tooltip's HoverStartTime, HoverXLocation, HoverYLocation, and
    IsDrawn fields, which Tooltip.updateMouseEvent changed from two goroutines, and render read during
    UpdateDisplay, with no lock. They are now guarded by uiStateMutex. Run with -race to detect it.

  - The mouse is moved between two hotspot cells on every iteration, so each call keeps re-arming the hover timer
    rather than settling, which maximizes the overlap between the goroutines.

Example:

	Expected Inputs:
		Tooltip with hotspot (0, 0, 20, 5) and no display delay. Two goroutines each run the hover state machine
		500 times through lockedTooltipMouseEvent, holding the UI lock as the event and periodic goroutines do,
		while the test goroutine alternates the mouse between (1, 1) and (2, 1) and calls UpdateDisplay. The mouse
		is then left at (1, 1) and the state machine is run twice more.

	Expected Outputs:
		No data race is reported, and the tooltip is drawn after the final two calls.
*/
func TestTooltipHoverStateConcurrentUpdatesAreSafe(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	tooltipInstance := layer1.AddTooltip("Tooltip", styleEntry, 0, 0, 20, 5, 3, 3, 25, 1, false, false, 0)
	UpdateDisplay(false)

	done := make(chan struct{})
	for workerIndex := 0; workerIndex < 2; workerIndex++ {
		go func() {
			defer func() { done <- struct{}{} }()
			for iteration := 0; iteration < 500; iteration++ {
				lockedTooltipMouseEvent()
			}
		}()
	}
	for iteration := 0; iteration < 500; iteration++ {
		SetMouseStatus(1+iteration%2, 1, 0, "")
		UpdateDisplay(false)
	}
	<-done
	<-done

	SetMouseStatus(1, 1, 0, "")
	lockedTooltipMouseEvent()
	lockedTooltipMouseEvent()
	tooltipEntry := tooltipRegistry.Get(layer1.GetAlias(), tooltipInstance.GetAlias())
	isDrawn := func() bool {
		defer lockUIState()()
		return tooltipEntry.IsDrawn
	}()
	if !isDrawn {
		test.Fatalf("expected the tooltip to be drawn once the mouse rests on its hotspot")
	}
}

/*
TestInitializeTerminalStopsUnrestoredSession is a test which allows you to verify that calling InitializeTerminal while
a previous session is still running, because RestoreTerminalSettings was never called, stops that session's background
goroutines before the new session is set up. In addition, the following should be noted:

  - This is a regression test for a data race where the old session's setupEventUpdater and
    setupPeriodicEventUpdater kept running, reading updateDisplayChannel and eventIntervalTime while the new
    InitializeTerminal call replaced them. Before the fix, the first session's stop channel was never closed.

  - stopEventGoroutines closes the channel and then waits for both goroutines to exit, so a closed channel here
    means the old goroutines have already stopped.

Example:

	Expected Inputs:
		InitializeTerminal(40, 20), then InitializeTerminal(40, 20) again with no RestoreTerminalSettings between.

	Expected Outputs:
		The first session's updateDisplayChannel is closed, and the second session's is a different, open channel.
*/
func TestInitializeTerminalStopsUnrestoredSession(test *testing.T) {
	commonResource.isDebugEnabled = true
	InitializeTerminal(40, 20)
	test.Cleanup(RestoreTerminalSettings)
	firstSessionChannel := commonResource.updateDisplayChannel

	InitializeTerminal(40, 20)
	secondSessionChannel := commonResource.updateDisplayChannel

	select {
	case <-firstSessionChannel:
	default:
		test.Fatalf("expected the unrestored first session's goroutines to be stopped by the second InitializeTerminal")
	}
	if secondSessionChannel == firstSessionChannel {
		test.Fatalf("expected the second session to have its own stop channel")
	}
	select {
	case <-secondSessionChannel:
		test.Fatalf("expected the second session's goroutines to still be running")
	default:
	}
}

/*
TestSelectorAndDropdownStateConcurrentAccessIsSafe is a test which allows you to verify that the application goroutine
can call selector, dropdown, file menu and shared control methods while the event goroutine handles keyboard, wheel
and mouse input for the same controls and the screen is redrawn, without a data race and without leaving their state
inconsistent. In addition, the following should be noted:

  - The race itself is only detected when the test is run with the race detector, as in go test -race. Without it
    the test still checks that the final state is consistent.

  - The test goroutine plays the event goroutine by injecting events and calling UpdateEventQueues, as the
    simulation screen helpers do, while a second goroutine plays the application.

Example:

	Expected Inputs:
		A bordered selector "list" at (2, 2) with twelve items and three rows, focused, a dropdown at (20, 2) with
		four items and a three row tray, and a file menu at (30, 0). For 300 iterations, one goroutine calls the
		selector's SetHighlighted, SetViewport, Select, GetSelected, AddItem or DeleteItem, FocusSelection,
		SetPosition, SetSize, GetBounds, GetViewport, IsNewItemSelected, GetSelectionSource, SetVisible, SetEnabled
		and, every 50 iterations, SetSelectionEntry; the dropdown's SetSelectedItemIndex, GetSelected, GetValue,
		SetSize and, every 50 iterations, SetSelectionEntry; the file menu's GetSelectedItem and Unselect; and
		UpdateDisplay. Meanwhile the test goroutine presses Down and Up, scrolls the wheel over the selector, clicks
		the selector, opens the dropdown and clicks a tray item, and calls UpdateDisplay.

	Expected Outputs:
		No data race is reported. Afterwards the selector's highlight is inside its items or none, its viewport is
		between zero and the furthest it can scroll, and its scroll bar value equals its viewport.
*/
func TestSelectorAndDropdownStateConcurrentAccessIsSafe(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	selector := Selector.Add(layerAlias, "list", styleEntry, getTestSelectionEntry(12), 2, 2, 3, 10, 1, 0, constants.SELECTED_NONE, false, true)
	dropdown := Dropdown.Add(layerAlias, "dropdown", styleEntry, getTestSelectionEntry(4), 20, 2, 3, 8, 0)
	fileMenu := FileMenu.Add(layerAlias, "menu", styleEntry, []string{"File"}, []types.SelectionEntryType{getTestSelectionEntry(3)}, 30, 0, true)
	simScreen := startInputSimulation(test)
	if err := SetFocus(&selector); err != nil {
		test.Fatalf("expected the selector to take focus, got %v", err)
	}

	const iterationCount = 300
	done := make(chan struct{})
	go func() {
		defer close(done)
		for iteration := 0; iteration < iterationCount; iteration++ {
			selector.SetHighlighted(iteration%14 - 1)
			selector.SetViewport(iteration % 9)
			selector.Select(fmt.Sprintf("item%d", iteration%10))
			selector.GetSelected()
			if iteration%2 == 0 {
				selector.AddItem(fmt.Sprintf("extra%d", iteration), "Extra")
			} else {
				selector.DeleteItem(12)
			}
			selector.FocusSelection("item5")
			selector.SetPosition(2, 2)
			selector.SetSize(10, 3)
			selector.GetBounds()
			selector.GetViewport()
			selector.IsNewItemSelected()
			selector.GetSelectionSource()
			selector.SetVisible(true)
			selector.SetEnabled(true)
			if iteration%50 == 0 {
				selector.SetSelectionEntry(getTestSelectionEntry(12))
				dropdown.SetSelectionEntry(getTestSelectionEntry(4 + iteration%3))
			}
			dropdown.SetSelectedItemIndex(iteration % 4)
			dropdown.GetSelected()
			dropdown.GetValue()
			dropdown.SetSize(11+iteration%2, 3)
			fileMenu.GetSelectedItem()
			fileMenu.Unselect()
			UpdateDisplay(false)
		}
	}()
	for iteration := 0; iteration < iterationCount; iteration++ {
		if iteration%2 == 0 {
			pressKey(simScreen, tcell.KeyDown, 0, tcell.ModNone)
		} else {
			pressKey(simScreen, tcell.KeyUp, 0, tcell.ModNone)
		}
		injectWheel(simScreen, 4, 3, tcell.WheelDown)
		clickAt(simScreen, 4, 2+iteration%3)
		clickAt(simScreen, 20, 2)
		clickAt(simScreen, 22, 4)
		UpdateDisplay(false)
	}
	<-done

	highlighted := selector.GetHighlighted()
	itemAliases, _ := selector.GetAllItems()
	if highlighted < constants.SELECTED_NONE || highlighted >= len(itemAliases) {
		test.Fatalf("expected the highlight to be inside the %d items or none, got %d", len(itemAliases), highlighted)
	}
	unlockUIState := lockUIState()
	selectorEntry := selectorRegistry.Get(layerAlias, "list")
	viewportPosition := selectorEntry.ViewportPosition
	maxViewportPosition := Selector.getMaxViewportPosition(selectorEntry)
	scrollValue := scrollbarRegistry.Get(layerAlias, selectorEntry.ScrollbarAlias).ScrollValue
	unlockUIState()
	if viewportPosition < 0 || viewportPosition > maxViewportPosition {
		test.Fatalf("expected the viewport between 0 and %d, got %d", maxViewportPosition, viewportPosition)
	}
	if scrollValue != viewportPosition {
		test.Fatalf("expected the scroll bar value %d to equal the viewport %d", scrollValue, viewportPosition)
	}
}

/*
runConcurrentlyWithInput is a method which allows you to run the given application work on its own goroutine while the
test goroutine plays the event goroutine, so that the two touch the same user interface state at the same time. For
each of iterationCount iterations, the work is called with the iteration number on the application goroutine, while
the test goroutine calls eventStep with the iteration number and then UpdateDisplay. It returns once both have
finished.

Example:

	runConcurrentlyWithInput(300, func(iteration int) { button.SetLabel("OK") }, func(iteration int) { clickAt(simScreen, 3, 3) })
*/
func runConcurrentlyWithInput(iterationCount int, applicationStep func(iteration int), eventStep func(iteration int)) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for iteration := 0; iteration < iterationCount; iteration++ {
			applicationStep(iteration)
		}
	}()
	for iteration := 0; iteration < iterationCount; iteration++ {
		eventStep(iteration)
		UpdateDisplay(false)
	}
	<-done
}

/*
TestControlStateConcurrentAccessIsSafe is a test which allows you to verify that the application goroutine can change
and read every kind of control while the event goroutine handles clicks and typing on the same controls and the
screen is redrawn, without a data race and without a deadlock. In addition, the following should be noted:

  - The race itself is only detected when the test is run with the race detector, as in go test -race.

Example:

	Expected Inputs:
		A button at (2, 2), a checkbox at (2, 6), a text field at (14, 2), a textbox at (14, 5), a radio button at
		(2, 8), a label at (30, 2), a progress bar at (30, 4) and a tooltip. For 300 iterations, one goroutine sets and
		reads their labels, values, text and states, while the test goroutine clicks the button, the checkbox and
		the radio button, clicks into the text field and the textbox and types into them, and redraws.

	Expected Outputs:
		No data race is reported, no call deadlocks, and afterwards the textbox and text field return the text the
		application last set or that text with typed characters in it.
*/
func TestControlStateConcurrentAccessIsSafe(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	layerInstance := &LayerInstanceType{layerAlias: layerAlias}
	button := Button.Add(layerAlias, "button", "OK", styleEntry, 2, 2, 8, 3, true)
	checkbox := Checkbox.Add(layerAlias, "checkbox", "Check", styleEntry, 2, 6, false, true)
	textField := TextField.Add(layerAlias, "field", styleEntry, 14, 2, 10, 20, false, "abc", true)
	textboxInstance := textbox.Add(layerAlias, "box", styleEntry, 14, 5, 10, 3, true)
	radioButtonInstance := layerInstance.AddRadioButton("Radio", styleEntry, 2, 8, 1, false)
	label := layerInstance.AddLabel("Label", styleEntry, 30, 2, 8)
	progressBar := layerInstance.AddProgressBar("Bar", styleEntry, 30, 4, 8, 1, false, 0, 100, false)
	tooltipInstance := layerInstance.AddTooltip("Tip", styleEntry, 0, 0, 5, 1, 0, 1, 10, 1, false, false, 0)
	simScreen := startInputSimulation(test)

	runConcurrentlyWithInput(300, func(iteration int) {
		button.SetLabel("OK")
		button.IsPressed()
		checkbox.SetState(iteration%2 == 0)
		checkbox.IsSelected()
		textField.SetValue("value")
		textField.GetValue()
		textboxInstance.SetText("hello\nworld")
		textboxInstance.GetText()
		radioButtonInstance.IsSelected()
		label.SetTooltipText("tip")
		progressBar.SetValue(iteration % 100)
		progressBar.GetValueAsPercent()
		tooltipInstance.SetText("Tip")
		button.SetPosition(2, 2)
		textboxInstance.SetSize(10, 3)
	}, func(iteration int) {
		switch iteration % 5 {
		case 0:
			clickAt(simScreen, 3, 3)
		case 1:
			clickAt(simScreen, 3, 6)
		case 2:
			clickAt(simScreen, 3, 8)
		case 3:
			clickAt(simScreen, 16, 6)
			pressKey(simScreen, tcell.KeyRune, 'a', tcell.ModNone)
		case 4:
			clickAt(simScreen, 16, 2)
			pressKey(simScreen, tcell.KeyRune, 'b', tcell.ModNone)
		}
	})

	if text := textboxInstance.GetText(); !strings.Contains(strings.ReplaceAll(text, "a", ""), "hello") {
		test.Fatalf("expected the textbox to hold the text last set, possibly with typed characters, got %q", text)
	}
	if value := textField.GetValue(); strings.ReplaceAll(value, "b", "") != "value" {
		test.Fatalf("expected the text field to hold the value last set, possibly with typed characters, got %q", value)
	}
}

/*
TestLayerStateConcurrentAccessIsSafe is a test which allows you to verify that the application goroutine can print
to, fill, move and hide a layer while the event goroutine handles clicks and the screen is redrawn, without a data
race and without a deadlock.

Example:

	Expected Inputs:
		The test layer. For 300 iterations, one goroutine locates the cursor, prints "hello world", fills the layer
		every 50 iterations, moves it between columns 0 and 2, and sets it visible, while the test goroutine clicks
		at (3, 3) and redraws.

	Expected Outputs:
		No data race is reported, no call deadlocks, and the layer ends visible.
*/
func TestLayerStateConcurrentAccessIsSafe(test *testing.T) {
	layerAlias, _ := setupInputTest(test)
	layerInstance := &LayerInstanceType{layerAlias: layerAlias}
	simScreen := startInputSimulation(test)

	runConcurrentlyWithInput(300, func(iteration int) {
		layerInstance.Locate(0, 0)
		layerInstance.Print("hello world")
		if iteration%50 == 0 {
			layerInstance.FillLayer("ab")
		}
		layerInstance.MoveLayerByAbsoluteValue(iteration%3, 0)
		layerInstance.SetIsVisible(true)
	}, func(iteration int) {
		clickAt(simScreen, 3, 3)
	})

	if !func() bool {
		defer lockUIState()()
		return layerRegistry.Get(layerAlias).IsVisible
	}() {
		test.Fatalf("expected the layer to end visible")
	}
}

/*
TestFocusStateConcurrentAccessIsSafe is a test which allows you to verify that the application goroutine can move and
read focus while the event goroutine handles Tab and Shift+Tab, without a data race and without a deadlock, and that
focus always ends on a registered stop.

Example:

	Expected Inputs:
		Buttons "a", "b" and "c" in the tab order. For 300 iterations, one goroutine calls SetFocus on one of them,
		FocusNext, FocusPrevious, GetFocusedControl and GetFocusChange, while the test goroutine presses Tab or
		Shift+Tab and redraws.

	Expected Outputs:
		No data race is reported, no call deadlocks, and the focused control is "a", "b" or "c".
*/
func TestFocusStateConcurrentAccessIsSafe(test *testing.T) {
	_, buttons, simScreen := setupTabOrderTest(test, "a", "b", "c")
	for index := range buttons {
		if err := buttons[index].AddToTabIndex(); err != nil {
			test.Fatalf("expected button %d to join the tab order, got %v", index, err)
		}
	}

	runConcurrentlyWithInput(300, func(iteration int) {
		SetFocus(&buttons[iteration%3])
		FocusNext()
		FocusPrevious()
		GetFocusedControl()
		GetFocusChange()
	}, func(iteration int) {
		if iteration%2 == 0 {
			pressKey(simScreen, tcell.KeyTab, 0, tcell.ModNone)
		} else {
			pressKey(simScreen, tcell.KeyBacktab, 0, tcell.ModShift)
		}
	})

	if _, controlAlias, _ := GetFocusedControl(); controlAlias != "a" && controlAlias != "b" && controlAlias != "c" {
		test.Fatalf("expected focus on a registered stop, got %q", controlAlias)
	}
}

/*
TestTextFieldHookMayCallExportedMethods is a test which allows you to verify that a text field's OnValueChanged hook,
which runs on the event goroutine, may call the package's exported methods, including ones on the same text field,
without deadlocking, and that its correction is still applied.

Example:

	Expected Inputs:
		A focused text field holding "ab" whose hook reads the field's value with GetValue, changes a second text
		field with SetValue, and returns its input in upper case. The key "c" is then pressed.

	Expected Outputs:
		The key is handled within five seconds, the hook saw "abc" through GetValue, the second field holds "abc",
		and the first field holds "ABC".
*/
func TestTextFieldHookMayCallExportedMethods(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	textField := TextField.Add(layerAlias, "field", styleEntry, 2, 2, 10, 20, false, "ab", true)
	otherTextField := TextField.Add(layerAlias, "other", styleEntry, 2, 4, 10, 20, false, "", true)
	simScreen := startInputSimulation(test)
	if err := SetFocus(&textField); err != nil {
		test.Fatalf("expected the text field to take focus, got %v", err)
	}
	textField.SetCursorPosition(2)
	var valueSeenByHook string
	textField.SetOnValueChanged(func(current string) string {
		valueSeenByHook = textField.GetValue()
		otherTextField.SetValue(current)
		return strings.ToUpper(current)
	})

	handled := make(chan struct{})
	go func() {
		defer close(handled)
		pressKey(simScreen, tcell.KeyRune, 'c', tcell.ModNone)
	}()
	select {
	case <-handled:
	case <-time.After(5 * time.Second):
		test.Fatalf("the keystroke was not handled within five seconds, so the hook deadlocked")
	}
	if valueSeenByHook != "abc" || otherTextField.GetValue() != "abc" || textField.GetValue() != "ABC" {
		test.Fatalf("expected the hook to see %q, the other field to hold %q and the field to hold %q, got %q, %q and %q",
			"abc", "abc", "ABC", valueSeenByHook, otherTextField.GetValue(), textField.GetValue())
	}
}

/*
TestPanicInsideEntryPointReleasesUIStateLock is a test which allows you to verify that a panic raised while an exported
method holds the UI lock still reaches the caller with its original value, releases the lock, and leaves the
terminal restored, so the program can fail cleanly rather than hang the next caller.

Example:

	Expected Inputs:
		A session with a simulated screen. PrintDialog called at (999, 999), outside the test layer, which panics
		while holding the UI lock.

	Expected Outputs:
		The panic is recovered with a message naming the out of bounds location, the UI lock can then be taken at
		once, and the session is marked terminated, since its screen was released.
*/
func TestPanicInsideEntryPointReleasesUIStateLock(test *testing.T) {
	layerAlias, _ := setupInputTest(test)
	layerInstance := &LayerInstanceType{layerAlias: layerAlias}
	// The screen is set up here rather than with startInputSimulation, since the panic guard releases it and
	// tcell's simulation screen cannot be released a second time by that helper's cleanup.
	stopEventGoroutines()
	simScreen := tcell.NewSimulationScreen("")
	if err := simScreen.Init(); err != nil {
		test.Fatalf("failed to initialize simulation screen: %v", err)
	}
	commonResource.screen = simScreen
	test.Cleanup(func() { commonResource.screen = nil })
	recovered := func() (recovered interface{}) {
		defer func() { recovered = recover() }()
		layerInstance.PrintDialog(999, 999, 0, 0, false, "text")
		return nil
	}()
	if message, isString := recovered.(string); !isString || !strings.Contains(message, "(999, 999)") {
		test.Fatalf("expected the out of bounds panic to reach the caller, got %v", recovered)
	}
	if !uiStateMutex.TryLock() {
		test.Fatalf("expected the UI lock to be released after the panic")
	}
	uiStateMutex.Unlock()
	commonResource.displayUpdate.RLock()
	isTerminated := commonResource.isTerminated
	commonResource.displayUpdate.RUnlock()
	if !isTerminated {
		test.Fatalf("expected the terminal to be restored after the panic")
	}
}

/*
TestSafeSttyPanicDoesNotWaitWhileDrawing is a test which allows you to verify that a fatal error raised while another
goroutine holds the display lock, as happens when an error is found while drawing, panics at once instead of waiting
for the display lock, which used to hang the program instead of failing it.

Example:

	Expected Inputs:
		The display lock held by the test goroutine while another goroutine calls safeSttyPanic("fatal").

	Expected Outputs:
		Within five seconds the other goroutine recovers the panic with the value "fatal".
*/
func TestSafeSttyPanicDoesNotWaitWhileDrawing(test *testing.T) {
	setupInputTest(test)
	commonResource.displayUpdate.Lock()
	recoveredValue := make(chan interface{}, 1)
	go func() {
		defer func() { recoveredValue <- recover() }()
		safeSttyPanic("fatal")
	}()
	select {
	case value := <-recoveredValue:
		commonResource.displayUpdate.Unlock()
		if value != "fatal" {
			test.Fatalf("expected the panic value %q, got %v", "fatal", value)
		}
	case <-time.After(5 * time.Second):
		commonResource.displayUpdate.Unlock()
		test.Fatalf("safeSttyPanic waited for the display lock instead of panicking")
	}
}

/*
TestTypewriterDialogDoesNotHoldUIStateLock is a test which allows you to verify that a dialog printed with a typewriter
delay releases the UI lock between characters, so other goroutines, including the event goroutine that delivers the
key that skips the dialog, keep running while it prints.

Example:

	Expected Inputs:
		PrintDialog of 40 characters with a 50 millisecond delay, about two seconds in all, started on its own
		goroutine. 200 milliseconds later the test goroutine calls UpdateDisplay and GetLocation on the layer.

	Expected Outputs:
		Both calls return within 500 milliseconds, long before the dialog finishes, and the dialog then finishes.
*/
func TestTypewriterDialogDoesNotHoldUIStateLock(test *testing.T) {
	layerAlias, _ := setupInputTest(test)
	layerInstance := &LayerInstanceType{layerAlias: layerAlias}
	dialogDone := make(chan struct{})
	go func() {
		defer close(dialogDone)
		layerInstance.PrintDialog(0, 0, 0, 50, false, strings.Repeat("x", 40))
	}()
	time.Sleep(200 * time.Millisecond)
	startTime := time.Now()
	UpdateDisplay(false)
	layerInstance.GetLocation()
	if elapsedTime := time.Since(startTime); elapsedTime > 500*time.Millisecond {
		test.Fatalf("expected the calls to return while the dialog prints, but they took %v", elapsedTime)
	}
	select {
	case <-dialogDone:
	case <-time.After(10 * time.Second):
		test.Fatalf("the dialog did not finish")
	}
}
