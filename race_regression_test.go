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
