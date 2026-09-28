package consolizer

import (
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"sync"
	"testing"
	"time"
)

/*
resetMouseEventState is a method which resets the package-level mouse throttle timer and interaction state back to
their zero values. In addition, the following should be noted:

  - eventStateMemory and lastMouseMoveTime are not reset by CommonTestSetup, so without this, state left behind by
    an earlier test in the same process could make a later test's first event throttled or its scrollbar-drag guard
    tripped unexpectedly.

Example:
    resetMouseEventState()
*/
func resetMouseEventState() {
	eventStateMemory.stateId = constants.EventStateNone
	lastMouseMoveTime = time.Time{}
}

/*
setupSimulationMouseScreen is a method which stops consolizer's background event goroutines and installs a tcell
simulation screen in their place, so a test can inject mouse events and drive them through UpdateEventQueues itself,
synchronously and one at a time. In addition, the following should be noted:

  - The background event goroutines are stopped first (rather than left running against the new screen) so that the
    test's own, synchronous UpdateEventQueues calls are never racing another goroutine's calls to the same function
    over the same screen and package-level mouse state.

  - This registers its own cleanup that finalizes the simulation screen and clears commonResource.screen back to
    nil, and does so before returning control to the test. Cleanups run last-registered-first, so this one runs
    before CommonTestSetup's own RestoreTerminalSettings cleanup, which would otherwise find commonResource.screen
    still pointing at a screen debug mode never created and, in a later test reusing this same process, at a screen
    an earlier test had already finalized, doubly closing it and panicking.

Example:
    simScreen := setupSimulationMouseScreen(test)
    simScreen.InjectMouse(3, 3, tcell.Button1, tcell.ModNone)
    UpdateEventQueues()
*/
func setupSimulationMouseScreen(test *testing.T) tcell.SimulationScreen {
	test.Helper()
	stopEventGoroutines()
	simScreen := tcell.NewSimulationScreen("")
	if err := simScreen.Init(); err != nil {
		test.Fatalf("failed to initialize simulation screen: %v", err)
	}
	simScreen.SetSize(commonResource.terminalWidth, commonResource.terminalHeight)
	simScreen.EnableMouse()
	commonResource.screen = simScreen
	test.Cleanup(func() {
		simScreen.Fini()
		commonResource.screen = nil
	})
	return simScreen
}

/*
findScrollbarHandle is a method which scans a vertical scrollbar's track for the cell currently occupied by its
handle, so tests do not need to duplicate the handle position formula from scrollbar.go.

Example:
    handleX, handleY := findScrollbarHandle(test, 5, 2, 10)
*/
func findScrollbarHandle(test *testing.T, xLocation int, yLocation int, length int) (int, int) {
	test.Helper()
	for currentYLocation := yLocation; currentYLocation < yLocation+length; currentYLocation++ {
		characterEntry := getCellInformationUnderMouseCursor(xLocation, currentYLocation)
		if characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar &&
			characterEntry.AttributeEntry.CellControlId == constants.CellControlIdScrollbarHandle {
			return xLocation, currentYLocation
		}
	}
	test.Fatalf("could not locate scrollbar handle within (%d, %d) to (%d, %d)", xLocation, yLocation, xLocation, yLocation+length-1)
	return 0, 0
}

/*
TestButtonFastClickAfterMovementReportsBothPresses is a test which verifies that a button release arriving within
the 50ms movement throttle window is still processed, and that the button can be clicked a second time right after.
In addition, the following should be noted:

  - Before the fix, the movement throttle also caught releases, so the release below would be dropped, leaving the
    button's IsPressed state stuck true and its next mouse down ignored since the press branch only fires when
    IsPressed is false.

Example:
    Expected Inputs:
        A mouse down on a button immediately followed by a release within 50ms of the last recorded movement, then
        a second down/release pair on the same button.

    Expected Outputs:
        Button.GetPressed() reports the button after each of the two release events.
*/
func TestButtonFastClickAfterMovementReportsBothPresses(test *testing.T) {
	resetMouseEventState()
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layerAlias := layer1.GetAlias()
	buttonInstance := Button.Add(layerAlias, "ok", "OK", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)

	simScreen := setupSimulationMouseScreen(test)

	buttonX, buttonY := 3, 3

	// A very recent "movement" establishes the throttle window that the release below must not be caught by.
	lastMouseMoveTime = time.Now()

	simScreen.InjectMouse(buttonX, buttonY, tcell.Button1, tcell.ModNone)
	UpdateEventQueues()
	if !buttonInstance.IsStatePressed() {
		test.Fatalf("expected button to be pressed after mouse down")
	}

	simScreen.InjectMouse(buttonX, buttonY, tcell.ButtonNone, tcell.ModNone)
	UpdateEventQueues()
	if buttonInstance.IsStatePressed() {
		test.Fatalf("expected the release to be processed even though it arrived within the movement throttle window")
	}

	firstLayerAlias, firstButtonAlias := buttonInstance.GetPressed()
	if firstLayerAlias != layerAlias || firstButtonAlias != "ok" {
		test.Fatalf("expected first press to be reported, got layer %q button %q", firstLayerAlias, firstButtonAlias)
	}

	simScreen.InjectMouse(buttonX, buttonY, tcell.Button1, tcell.ModNone)
	UpdateEventQueues()
	if !buttonInstance.IsStatePressed() {
		test.Fatalf("expected second mouse down to register a new press")
	}

	simScreen.InjectMouse(buttonX, buttonY, tcell.ButtonNone, tcell.ModNone)
	UpdateEventQueues()

	secondLayerAlias, secondButtonAlias := buttonInstance.GetPressed()
	if secondLayerAlias != layerAlias || secondButtonAlias != "ok" {
		test.Fatalf("expected second press to be reported, got layer %q button %q", secondLayerAlias, secondButtonAlias)
	}
}

/*
TestButtonHeldMoveWithinSameButtonStillClicksOnRelease is a test which verifies that moving the mouse to a different
cell of the same button while it remains held does not stop the eventual release over that button from being
reported as a click. In addition, the following should be noted:

  - The click is only ever recorded on release, so a held move by itself must not report anything yet; it is the
    release afterward that must still count, confirming the held move along the way did not clear the button's
    armed (IsPressed) state.

Example:
    Expected Inputs:
        A mouse down on a button, a held move to another cell within the same button's bounds, then a release over
        that same button.

    Expected Outputs:
        Button.GetPressed() reports nothing after the held move, but reports the button after the release.
*/
func TestButtonHeldMoveWithinSameButtonStillClicksOnRelease(test *testing.T) {
	resetMouseEventState()
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layerAlias := layer1.GetAlias()
	buttonInstance := Button.Add(layerAlias, "ok", "OK", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)

	SetMouseStatus(3, 3, 1, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)
	if !buttonInstance.IsStatePressed() {
		test.Fatalf("expected button to be pressed after mouse down")
	}

	// Held move to a different cell of the same button, before release.
	SetMouseStatus(7, 3, 1, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)
	if !buttonInstance.IsStatePressed() {
		test.Fatalf("expected button to remain pressed after a held move within its own bounds")
	}
	if layerAlias, buttonAlias := buttonInstance.GetPressed(); layerAlias != "" || buttonAlias != "" {
		test.Fatalf("expected no click to be reported before release, got layer %q button %q", layerAlias, buttonAlias)
	}

	SetMouseStatus(7, 3, 0, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)

	pressedLayerAlias, pressedButtonAlias := buttonInstance.GetPressed()
	if pressedLayerAlias != layerAlias || pressedButtonAlias != "ok" {
		test.Fatalf("expected the release over the same button to report a click, got layer %q button %q", pressedLayerAlias, pressedButtonAlias)
	}
}

/*
TestButtonDragOffAndBackOnStillClicksOnRelease is a test which verifies that dragging off a pressed button cancels
it (no click on a release elsewhere), but dragging back onto that same button before releasing re-arms it so the
release is still reported as a click. In addition, the following should be noted:

  - A button only fires a click when the mouse-up lands back on the button that was pressed. Releasing off any
    button, or over a different one, must cancel silently instead.

Example:
    Expected Inputs:
        A mouse down on a button, a held move off it, a release while off it, then a fresh down/move-back-on/release
        sequence on the same button.

    Expected Outputs:
        The first release (off the button) reports no click. The second, after dragging back onto the button before
        releasing, reports the button.
*/
func TestButtonDragOffAndBackOnStillClicksOnRelease(test *testing.T) {
	resetMouseEventState()
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layerAlias := layer1.GetAlias()
	buttonInstance := Button.Add(layerAlias, "ok", "OK", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)

	SetMouseStatus(3, 3, 1, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)
	if !buttonInstance.IsStatePressed() {
		test.Fatalf("expected button to be pressed after mouse down")
	}

	// Held move onto a plain, non-button cell cancels the press.
	SetMouseStatus(25, 15, 1, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)
	if buttonInstance.IsStatePressed() {
		test.Fatalf("expected button to no longer be pressed after being dragged off")
	}

	// Releasing off the button must not report a click.
	SetMouseStatus(25, 15, 0, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)
	if layerAlias, buttonAlias := buttonInstance.GetPressed(); layerAlias != "" || buttonAlias != "" {
		test.Fatalf("expected no click to be reported for a release off the button, got layer %q button %q", layerAlias, buttonAlias)
	}

	// A fresh press, dragged back onto the same button before release, must still click.
	SetMouseStatus(3, 3, 1, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)
	if !buttonInstance.IsStatePressed() {
		test.Fatalf("expected button to be pressed again after a fresh mouse down")
	}

	SetMouseStatus(7, 3, 1, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)
	if !buttonInstance.IsStatePressed() {
		test.Fatalf("expected button to remain pressed after moving back onto it")
	}

	SetMouseStatus(7, 3, 0, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)

	pressedLayerAlias, pressedButtonAlias := buttonInstance.GetPressed()
	if pressedLayerAlias != layerAlias || pressedButtonAlias != "ok" {
		test.Fatalf("expected the release back over the same button to report a click, got layer %q button %q", pressedLayerAlias, pressedButtonAlias)
	}
}

/*
TestButtonReleaseOverDifferentButtonClearsOriginalPress is a test which verifies that releasing the mouse over a
different button than the one the press started on does not report a click for either button, clears the original
button's pressed state, and that the original button responds normally to its next (matched) click. In addition,
the following should be noted:

  - A click only counts when the release lands back on the button that was pressed. A press on A released over B
    must cancel silently: neither A nor B should be reported by GetPressed, and neither should remain pressed.

Example:
    Expected Inputs:
        A mouse down on button A followed by a release over button B.

    Expected Outputs:
        GetPressed reports nothing for that sequence. Button A is no longer pressed, and a following down/up on A
        alone is reported by GetPressed.
*/
func TestButtonReleaseOverDifferentButtonClearsOriginalPress(test *testing.T) {
	resetMouseEventState()
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layerAlias := layer1.GetAlias()
	buttonA := Button.Add(layerAlias, "a", "A", styleEntry, 2, 2, 8, 3, true)
	buttonB := Button.Add(layerAlias, "b", "B", styleEntry, 15, 2, 8, 3, true)
	UpdateDisplay(false)

	SetMouseStatus(3, 3, 1, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)
	if !buttonA.IsStatePressed() {
		test.Fatalf("expected button A to be pressed after mouse down")
	}

	// Release over button B instead of A, as would happen if a dialog closed and reopened between down and up.
	SetMouseStatus(16, 3, 0, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)

	if buttonA.IsStatePressed() {
		test.Fatalf("expected button A to no longer be pressed after releasing over a different button")
	}
	if buttonB.IsStatePressed() {
		test.Fatalf("expected button B, which was never pressed, to remain unpressed")
	}
	if layerAlias, buttonAlias := buttonA.GetPressed(); layerAlias != "" || buttonAlias != "" {
		test.Fatalf("expected no click to be reported when the release lands on a different button, got layer %q button %q", layerAlias, buttonAlias)
	}

	// Button A must respond to its next, matched click.
	SetMouseStatus(3, 3, 1, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)
	SetMouseStatus(3, 3, 0, "")
	bringLayerToFrontIfRequired()
	Button.updateStates(true)

	pressedLayerAlias, pressedButtonAlias := buttonA.GetPressed()
	if pressedLayerAlias != layerAlias || pressedButtonAlias != "a" {
		test.Fatalf("expected button A's next click to be reported, got layer %q button %q", pressedLayerAlias, pressedButtonAlias)
	}
}

/*
TestScrollbarDragReleaseWithinThrottleWindowUnsticksState is a test which verifies that releasing the mouse during a
scrollbar drag, within 50ms of the last recorded drag movement, still ends the drag and lets subsequent button
clicks register. In addition, the following should be noted:

  - Before the fix, any event while eventStateMemory.stateId was EventStateDragAndDropScrollbar was subject to the
    movement throttle, including the release meant to end the drag. Dropping that release left the drag state
    stuck, and buttonType.updateStateMouse bails out unconditionally while that state is set.

Example:
    Expected Inputs:
        A mouse down on a scrollbar handle, a held drag move, then a release within 50ms of that move.

    Expected Outputs:
        eventStateMemory.stateId returns to EventStateNone, and a following button click is reported.
*/
func TestScrollbarDragReleaseWithinThrottleWindowUnsticksState(test *testing.T) {
	resetMouseEventState()
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layerAlias := layer1.GetAlias()
	layer1.AddScrollbar(styleEntry, 5, 2, 10, 100, 0, 1, false)
	// Placed clear of layer2's (3,10)-(43,30) opaque area from CommonTestSetup, which would otherwise render over
	// this button in the composited screen and make it unclickable.
	buttonInstance := Button.Add(layerAlias, "ok", "OK", styleEntry, 20, 2, 10, 3, true)
	UpdateDisplay(false)

	handleX, handleY := findScrollbarHandle(test, 5, 2, 10)

	simScreen := setupSimulationMouseScreen(test)

	// Press the scrollbar handle to begin a drag.
	simScreen.InjectMouse(handleX, handleY, tcell.Button1, tcell.ModNone)
	UpdateEventQueues()
	if eventStateMemory.stateId != constants.EventStateDragAndDropScrollbar {
		test.Fatalf("expected scrollbar drag state to begin, got state %d", eventStateMemory.stateId)
	}

	// Drag the handle down by one row while still held; this also refreshes the movement throttle timer.
	simScreen.InjectMouse(handleX, handleY+1, tcell.Button1, tcell.ModNone)
	UpdateEventQueues()

	// Release immediately, well within the 50ms throttle window the drag movement above just set.
	simScreen.InjectMouse(handleX, handleY+1, tcell.ButtonNone, tcell.ModNone)
	UpdateEventQueues()
	if eventStateMemory.stateId != constants.EventStateNone {
		test.Fatalf("expected scrollbar drag state to end after release, got state %d", eventStateMemory.stateId)
	}

	// A previously stuck drag state makes Button.updateStateMouse bail out unconditionally, so verify a click now
	// registers normally.
	buttonX, buttonY := 21, 3
	simScreen.InjectMouse(buttonX, buttonY, tcell.Button1, tcell.ModNone)
	UpdateEventQueues()
	if !buttonInstance.IsStatePressed() {
		test.Fatalf("expected button press to register after the scrollbar drag ended")
	}

	simScreen.InjectMouse(buttonX, buttonY, tcell.ButtonNone, tcell.ModNone)
	UpdateEventQueues()

	pressedLayerAlias, pressedButtonAlias := buttonInstance.GetPressed()
	if pressedLayerAlias != layerAlias || pressedButtonAlias != "ok" {
		test.Fatalf("expected button click to be reported, got layer %q button %q", pressedLayerAlias, pressedButtonAlias)
	}
}

/*
TestButtonHistoryConcurrentAccess is a test which allows you to verify that buttonHistory can be read via
Button.GetPressed on one goroutine while simulated clicks are written to it on another without triggering a data
race. In addition, the following should be noted:

  - This is a regression test for buttonHistory being a bare, unsynchronized struct that the event goroutine, the
    periodic-event goroutine, and the application's own goroutine could all touch at once.

  - Run with go test -race for this test to be meaningful; it does not assert on the reported values themselves,
    only that concurrent access is safe.

Example:
    go test -run TestButtonHistoryConcurrentAccess -race -count=1 .
*/
func TestButtonHistoryConcurrentAccess(test *testing.T) {
	resetMouseEventState()
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layerAlias := layer1.GetAlias()
	buttonInstance := Button.Add(layerAlias, "ok", "OK", styleEntry, 2, 2, 10, 3, true)
	UpdateDisplay(false)

	stop := make(chan struct{})
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()
		for {
			select {
			case <-stop:
				return
			default:
				buttonInstance.GetPressed()
			}
		}
	}()

	go func() {
		defer waitGroup.Done()
		isPressed := false
		deadline := time.Now().Add(300 * time.Millisecond)
		for time.Now().Before(deadline) {
			isPressed = !isPressed
			if isPressed {
				SetMouseStatus(3, 3, 1, "")
			} else {
				SetMouseStatus(3, 3, 0, "")
			}
			bringLayerToFrontIfRequired()
			Button.updateStates(true)
		}
		close(stop)
	}()

	waitGroup.Wait()
}
