package consolizer

import (
	"github.com/gdamore/tcell/v2"
	"sync"
	"testing"
	"time"
)

/*
setupOrderedKeyboardTest is a method which allows you to prepare a keyboard test around a focused selector with the
four items "a" to "d" and item 0 highlighted, installing a simulation screen the test drives itself, and returning the
selector and that screen.

Example:

	selector, simScreen := setupOrderedKeyboardTest(test)
*/
func setupOrderedKeyboardTest(test *testing.T) (SelectorInstanceType, tcell.SimulationScreen) {
	test.Helper()
	layerAlias, styleEntry := setupInputTest(test)
	selector := addTestSelector(layerAlias, styleEntry)
	if err := SetFocus(&selector); err != nil {
		test.Fatalf("expected the selector to take focus, got %v", err)
	}
	selector.SetHighlighted(0)
	return selector, startInputSimulation(test)
}

/*
pressDownSpaceDown is a method which allows you to type Down, Space, Down into the simulation screen, handling each
event through UpdateEventQueues as the event goroutine would.

Example:

	pressDownSpaceDown(simScreen)
*/
func pressDownSpaceDown(simScreen tcell.SimulationScreen) {
	pressKey(simScreen, tcell.KeyDown, 0, tcell.ModNone)
	pressKey(simScreen, tcell.KeyRune, ' ', tcell.ModNone)
	pressKey(simScreen, tcell.KeyDown, 0, tcell.ModNone)
}

/*
assertHighlighted is a method which allows you to fail the test unless the selector highlights the expected item.

Example:

	assertHighlighted(test, selector, "after Down", 1)
*/
func assertHighlighted(test *testing.T, selector SelectorInstanceType, context string, expectedIndex int) {
	test.Helper()
	if highlightedIndex := selector.GetHighlighted(); highlightedIndex != expectedIndex {
		test.Fatalf("%s: expected item %d highlighted, got %d", context, expectedIndex, highlightedIndex)
	}
}

/*
assertInkey is a method which allows you to fail the test unless the next Inkey call returns the expected keystroke,
where nil means no keystroke.

Example:

	assertInkey(test, "second call", []rune(" "))
*/
func assertInkey(test *testing.T, context string, expectedKeystroke []rune) {
	test.Helper()
	keystroke := Inkey()
	if (keystroke == nil) != (expectedKeystroke == nil) || string(keystroke) != string(expectedKeystroke) {
		test.Fatalf("%s: expected Inkey to return %q (nil %v), got %q (nil %v)", context, string(expectedKeystroke),
			expectedKeystroke == nil, string(keystroke), keystroke == nil)
	}
}

/*
makeInkeyStale is a method which allows you to pretend the application last called Inkey longer ago than ordered
keyboard mode waits for, without sleeping.

Example:

	makeInkeyStale()
*/
func makeInkeyStale() {
	orderedKeyboard.mutex.Lock()
	orderedKeyboard.lastInkeyTime = time.Now().Add(-2 * orderedKeyboardInkeyTimeout)
	orderedKeyboard.mutex.Unlock()
}

/*
TestOrderedKeyboardInkeyHandlesHeldKeysInOrder is a test which allows you to verify that, with ordered keyboard mode
on and the application calling Inkey, keys are held back and each Inkey call handles exactly one of them in the order
typed, returning nil for a key the selector used and the key itself otherwise.

Example:

	Expected Inputs:
		Focused selector [a, b, c, d] with item 0 highlighted, SetOrderedKeyboard(true), one Inkey call, then
		Down, Space, Down typed before the next Inkey call.

	Expected Outputs:
		Item 0 still highlighted after typing. Inkey then returns nil (item 1), " " (item 1), nil (item 2), and
		nil again with nothing left.
*/
func TestOrderedKeyboardInkeyHandlesHeldKeysInOrder(test *testing.T) {
	selector, simScreen := setupOrderedKeyboardTest(test)
	SetOrderedKeyboard(true)
	if !IsOrderedKeyboard() {
		test.Fatalf("expected ordered keyboard mode to be on")
	}
	assertInkey(test, "before typing", nil)

	pressDownSpaceDown(simScreen)
	assertHighlighted(test, selector, "after typing", 0)

	assertInkey(test, "first call", nil)
	assertHighlighted(test, selector, "after first call", 1)
	assertInkey(test, "second call", []rune(" "))
	assertHighlighted(test, selector, "after second call", 1)
	assertInkey(test, "third call", nil)
	assertHighlighted(test, selector, "after third call", 2)
	assertInkey(test, "fourth call", nil)
	assertHighlighted(test, selector, "after fourth call", 2)
}

/*
TestOrderedKeyboardOffHandlesKeysImmediately is a test which allows you to verify that with ordered keyboard mode off,
the default, keys a control uses are handled on arrival and unused keys wait in the keyboard buffer, as before the mode
existed, even while the application is calling Inkey.

Example:

	Expected Inputs:
		Focused selector [a, b, c, d] with item 0 highlighted, ordered mode off, one Inkey call, then Down, Space,
		Down.

	Expected Outputs:
		Item 2 highlighted straight after typing, then Inkey returns " " followed by nil.
*/
func TestOrderedKeyboardOffHandlesKeysImmediately(test *testing.T) {
	selector, simScreen := setupOrderedKeyboardTest(test)
	if IsOrderedKeyboard() {
		test.Fatalf("expected ordered keyboard mode to be off by default")
	}
	assertInkey(test, "before typing", nil)

	pressDownSpaceDown(simScreen)
	assertHighlighted(test, selector, "after typing", 2)
	assertInkey(test, "first call", []rune(" "))
	assertInkey(test, "second call", nil)
}

/*
TestOrderedKeyboardWithoutInkeyHandlesKeysImmediately is a test which allows you to verify that ordered keyboard mode
does not hold keys back for an application that has never called Inkey, such as one that only polls its buttons.

Example:

	Expected Inputs:
		Focused selector [a, b, c, d] with item 0 highlighted, SetOrderedKeyboard(true), no Inkey call, then Down,
		Space, Down.

	Expected Outputs:
		Item 2 highlighted straight after typing and the keyboard buffer holding [" "].
*/
func TestOrderedKeyboardWithoutInkeyHandlesKeysImmediately(test *testing.T) {
	selector, simScreen := setupOrderedKeyboardTest(test)
	SetOrderedKeyboard(true)

	pressDownSpaceDown(simScreen)
	assertHighlighted(test, selector, "after typing", 2)
	assertKeyboardBuffer(test, "after typing", " ")
}

/*
TestOrderedKeyboardFlushesWhenInkeyStops is a test which allows you to verify that key events held back are handed to
the controls in order once the application has not called Inkey for 250 milliseconds, both by the periodic update and
by the next key to arrive, with unused keys placed in the keyboard buffer ahead of anything typed later.

Example:

	Expected Inputs:
		Focused selector [a, b, c, d] with item 0 highlighted and ordered mode on. Case one: Inkey, Down, Space, Down,
		Inkey goes stale, UpdatePeriodicEvents. Case two: Inkey, Down, "x", Inkey goes stale, then "y" typed.

	Expected Outputs:
		Case one: item 0 until the periodic update, then item 2 and buffer [" "]. Case two: item 1 and buffer
		["x", "y"] after "y" arrives.
*/
func TestOrderedKeyboardFlushesWhenInkeyStops(test *testing.T) {
	selector, simScreen := setupOrderedKeyboardTest(test)
	SetOrderedKeyboard(true)

	assertInkey(test, "periodic: before typing", nil)
	pressDownSpaceDown(simScreen)
	UpdatePeriodicEvents()
	assertHighlighted(test, selector, "periodic: while Inkey is recent", 0)
	makeInkeyStale()
	UpdatePeriodicEvents()
	assertHighlighted(test, selector, "periodic: after Inkey went stale", 2)
	assertKeyboardBuffer(test, "periodic: after Inkey went stale", " ")

	selector.SetHighlighted(0)
	assertInkey(test, "next key: before typing", nil)
	pressKey(simScreen, tcell.KeyDown, 0, tcell.ModNone)
	pressKey(simScreen, tcell.KeyRune, 'x', tcell.ModNone)
	assertHighlighted(test, selector, "next key: while Inkey is recent", 0)
	makeInkeyStale()
	pressKey(simScreen, tcell.KeyRune, 'y', tcell.ModNone)
	assertHighlighted(test, selector, "next key: after Inkey went stale", 1)
	assertKeyboardBuffer(test, "next key: after Inkey went stale", "x", "y")
}

/*
TestOrderedKeyboardTurnedOffFlushesHeldKeys is a test which allows you to verify that turning ordered keyboard mode
off hands every key event held back to the controls in order, and that keys typed afterwards are handled immediately.

Example:

	Expected Inputs:
		Focused selector [a, b, c, d] with item 0 highlighted, ordered mode on, Inkey, Down, Space, Down, then
		SetOrderedKeyboard(false), then Down.

	Expected Outputs:
		Item 2 and buffer [" "] once the mode is turned off, then item 3 straight after the last Down.
*/
func TestOrderedKeyboardTurnedOffFlushesHeldKeys(test *testing.T) {
	selector, simScreen := setupOrderedKeyboardTest(test)
	SetOrderedKeyboard(true)
	assertInkey(test, "before typing", nil)
	pressDownSpaceDown(simScreen)
	assertHighlighted(test, selector, "while held", 0)

	SetOrderedKeyboard(false)
	if IsOrderedKeyboard() {
		test.Fatalf("expected ordered keyboard mode to be off")
	}
	assertHighlighted(test, selector, "after turning the mode off", 2)
	assertKeyboardBuffer(test, "after turning the mode off", " ")

	pressKey(simScreen, tcell.KeyDown, 0, tcell.ModNone)
	assertHighlighted(test, selector, "after Down with the mode off", 3)
}

/*
TestOrderedKeyboardConcurrentInkey is a test which allows you to verify, under the race detector, that keys typed
while another goroutine keeps calling Inkey all reach the application and the controls exactly once, with the unused
keys in the order typed.

Example:

	Expected Inputs:
		Focused selector [a, b, c, d] with item 0 highlighted and ordered mode on. One goroutine calls Inkey in a loop
		while Down, "a", Down, "b" are typed. The mode is then turned off and the keyboard buffer drained.

	Expected Outputs:
		The keys returned by Inkey followed by the buffer are exactly ["a", "b"], and item 2 is highlighted.
*/
func TestOrderedKeyboardConcurrentInkey(test *testing.T) {
	selector, simScreen := setupOrderedKeyboardTest(test)
	SetOrderedKeyboard(true)
	assertInkey(test, "before typing", nil)

	var receivedKeys []string
	stopReader := make(chan struct{})
	var readerGroup sync.WaitGroup
	readerGroup.Add(1)
	go func() {
		defer readerGroup.Done()
		for {
			select {
			case <-stopReader:
				return
			default:
			}
			if keystroke := Inkey(); keystroke != nil {
				receivedKeys = append(receivedKeys, string(keystroke))
			}
			time.Sleep(time.Millisecond)
		}
	}()
	pressKey(simScreen, tcell.KeyDown, 0, tcell.ModNone)
	pressKey(simScreen, tcell.KeyRune, 'a', tcell.ModNone)
	pressKey(simScreen, tcell.KeyDown, 0, tcell.ModNone)
	pressKey(simScreen, tcell.KeyRune, 'b', tcell.ModNone)
	time.Sleep(20 * time.Millisecond)
	close(stopReader)
	readerGroup.Wait()

	SetOrderedKeyboard(false)
	receivedKeys = append(receivedKeys, readKeyboardBuffer()...)
	if len(receivedKeys) != 2 || receivedKeys[0] != "a" || receivedKeys[1] != "b" {
		test.Fatalf("expected the application to receive [a b], got %q", receivedKeys)
	}
	assertHighlighted(test, selector, "after all keys", 2)
}
