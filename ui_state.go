package consolizer

import "sync"

// uiStateMutex guards all user interface state: every layer, every control of every type, the scroll bars and
// tooltips attached to controls, and the control state the focus manager reads. The application goroutine changes
// this state through the exported functions and methods, the event goroutine changes it while handling input, the
// periodic goroutine changes it while timing tooltips, and drawing reads it. Each of these takes the lock once, at
// the point where it enters the package, and holds it for the whole operation, so no goroutine ever sees state that
// another is half way through changing. In addition, the following should be noted:
//
//   - Lock order: uiStateMutex, then commonResource.displayUpdate, then focusManager.mutex, then leaf locks only
//     (the event state lock, mouse and keyboard memory, the button history, and the control memory managers).
//     uiStateMutex is always the outermost lock, so it must never be taken while any of the others is held.
//
//   - The lock is not reentrant. Only the outermost entry point takes it. Everything that runs inside an entry
//     point, including the unexported helpers and the unexported versions of exported functions that the package
//     calls internally, expects it to already be held and never takes it again.
//
//   - Application code is never called while the lock is held. A hook such as a text field's OnValueChanged is
//     called with the lock released, so it may call any exported function.
//
//   - It must not be held while waiting for the event goroutines to stop, since they may be waiting for it.
var uiStateMutex sync.Mutex

/*
lockUIState is a method which allows you to take uiStateMutex at an entry point of the package, and returns the
function that releases it, so that a caller writes defer lockUIState()() as its first statement and holds the lock
for the rest of its body. In addition, the following should be noted:

  - If a panic passes through the caller, the release function unlocks the mutex, puts the terminal back into a
    usable state as described for restoreTerminalAfterFatalError, and then panics again with the same value, so
    the program still fails with the original error rather than leaving the lock held or the terminal unusable.

  - The release function must be called by defer, directly, for it to observe the panic.

Example:

	defer lockUIState()()
*/
func lockUIState() func() {
	uiStateMutex.Lock()
	return func() {
		if recovered := recover(); recovered != nil {
			uiStateMutex.Unlock()
			restoreTerminalAfterFatalError()
			panic(recovered)
		}
		uiStateMutex.Unlock()
	}
}

/*
restoreTerminalAfterFatalError is a method which allows you to put the terminal back into a usable state just before
the program fails with a fatal error, waiting for the display lock if another goroutine is drawing. The event goroutines
are told to stop but are not waited for, since the caller may be one of them, and the screen is then released, unless
it already has been. In addition, the following should be noted:

  - The caller must hold no lock at all, as is the case in the panic guard of lockUIState, since every other lock is
    taken inside uiStateMutex and has already been released when the guard runs.

  - Layers and controls are left in place, unlike RestoreTerminalSettings, since the program is about to fail.

Example:

	restoreTerminalAfterFatalError()
*/
func restoreTerminalAfterFatalError() {
	signalEventGoroutinesToStop()
	commonResource.displayUpdate.Lock()
	defer commonResource.displayUpdate.Unlock()
	releaseScreenAfterFatalError()
}

/*
tryRestoreTerminalAfterFatalError is a method which allows you to put the terminal back into a usable state as
restoreTerminalAfterFatalError does, but only if that can be done without waiting, so it is safe to call while any
lock is held. If the display lock is busy, for example because the error was found while drawing, the screen is left
for the panic guard of lockUIState to release once the panic has unwound the locks.

Example:

	tryRestoreTerminalAfterFatalError()
*/
func tryRestoreTerminalAfterFatalError() {
	signalEventGoroutinesToStop()
	if !commonResource.displayUpdate.TryLock() {
		return
	}
	defer commonResource.displayUpdate.Unlock()
	releaseScreenAfterFatalError()
}

/*
releaseScreenAfterFatalError is a method which allows you to release the screen, so the terminal returns to its normal
settings, and mark the session as terminated so nothing draws to the screen again. Nothing happens if the session is
already terminated or has no screen. The caller must hold the display lock.

Example:

	releaseScreenAfterFatalError()
*/
func releaseScreenAfterFatalError() {
	if commonResource.isTerminated || commonResource.screen == nil {
		return
	}
	commonResource.isTerminated = true
	commonResource.screen.DisableMouse()
	commonResource.screen.Fini()
}

/*
callWithUIStateUnlocked is a method which allows you to call application code, such as a hook the application has
registered, from inside an entry point that holds uiStateMutex. The lock is released for the duration of the call and
taken again before this returns, even if the call panics, so the application code may itself call any exported
function without deadlocking, and the entry point's own release still pairs with a held lock. In addition, the
following should be noted:

  - Any state read before the call may have been changed by the time it returns, by the application code or by
    another goroutine, so the caller must look its entries up again afterwards rather than reuse them.

Example:

	callWithUIStateUnlocked(func() { correctedValue = hook(currentValue) })
*/
func callWithUIStateUnlocked(callback func()) {
	uiStateMutex.Unlock()
	defer uiStateMutex.Lock()
	callback()
}
