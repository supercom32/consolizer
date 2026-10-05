package consolizer

import (
	"testing"
	"time"

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
    same alias, so the race is isolated to the Tooltips lookup rather than the already-covered Buttons lookup.

  - The test drives getFromCharacterEntry directly with a synthetic character entry rather than going through mouse
    events, since that is the exact call the periodic goroutine makes.

Example:

	go test -run TestTooltipDeleteDuringLookupDoesNotPanic -count=1 .
*/
func TestTooltipDeleteDuringLookupDoesNotPanic(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	layerAlias := layer1.GetAlias()

	buttonInstance := Button.Add(layerAlias, "back", "BACK", styleEntry, 2, 2, 10, 3, true)
	buttonEntry := Buttons.Get(layerAlias, buttonInstance.GetAlias())
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
    UpdateDisplay, with no lock. Run with -race to detect it.

  - The mouse is moved between two hotspot cells on every iteration, so each call keeps re-arming the hover timer
    rather than settling, which maximizes the overlap between the goroutines.

Example:

	Expected Inputs:
		Tooltip with hotspot (0, 0, 20, 5) and no display delay. Two goroutines each call Tooltip.updateMouseEvent
		500 times while the test goroutine alternates the mouse between (1, 1) and (2, 1) and calls UpdateDisplay.
		The mouse is then left at (1, 1) and Tooltip.updateMouseEvent is called twice.

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
				Tooltip.updateMouseEvent()
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
	Tooltip.updateMouseEvent()
	Tooltip.updateMouseEvent()
	tooltipEntry := Tooltips.Get(layer1.GetAlias(), tooltipInstance.GetAlias())
	tooltipStateMutex.Lock()
	isDrawn := tooltipEntry.IsDrawn
	tooltipStateMutex.Unlock()
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
