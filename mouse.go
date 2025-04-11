package consolizer

import (
	"sync"
)

/*
mouseMemoryType is a structure that holds information about the current mouse state
including position, button status, and wheel state.
*/
type mouseMemoryType struct {
	sync.Mutex
	xLocation     int
	yLocation     int
	buttonPressed uint
	wheelState    string
}

var MouseMemory mouseMemoryType
var PreviousMouseMemory mouseMemoryType

/*
ClearMouseMemory allows you to reset both the current and previous mouse memory
states to their default values. In addition, the following information should be noted:

- The default x and y locations are set to -1 (off-screen).
- The default button pressed state is 0 (no buttons pressed).
- The default wheel state is an empty string (no wheel movement).
*/
func ClearMouseMemory() {
	MouseMemory.Lock()
	defer func() {
		MouseMemory.Unlock()
	}()
	MouseMemory.xLocation = -1
	MouseMemory.yLocation = -1
	MouseMemory.buttonPressed = 0
	MouseMemory.wheelState = ""
	PreviousMouseMemory.Lock()
	defer func() {
		PreviousMouseMemory.Unlock()
	}()
	PreviousMouseMemory.xLocation = -1
	PreviousMouseMemory.yLocation = -1
	PreviousMouseMemory.buttonPressed = 0
	PreviousMouseMemory.wheelState = ""
}

/*
SetMouseStatus allows you to update the current mouse status while preserving the
previous state. In addition, the following information should be noted:

- The previous mouse state is updated with the current state before changing.
- The current mouse state is updated with the provided parameters.
- This method is thread-safe as it uses mutex locks to prevent race conditions.
*/
func SetMouseStatus(xLocation int, yLocation int, buttonPressed uint, wheelState string) {
	PreviousMouseMemory.Lock()
	defer func() {
		PreviousMouseMemory.Unlock()
	}()
	PreviousMouseMemory.xLocation = MouseMemory.xLocation
	PreviousMouseMemory.yLocation = MouseMemory.yLocation
	PreviousMouseMemory.buttonPressed = MouseMemory.buttonPressed
	PreviousMouseMemory.wheelState = MouseMemory.wheelState
	MouseMemory.Lock()
	defer func() {
		MouseMemory.Unlock()
	}()
	MouseMemory.xLocation = xLocation
	MouseMemory.yLocation = yLocation
	MouseMemory.buttonPressed = buttonPressed
	MouseMemory.wheelState = wheelState
}

/*
GetMouseStatus allows you to retrieve the current mouse status including position,
button state, and wheel state. In addition, the following information should be noted:

- Returns the x location, y location, button pressed state, and wheel state.
- This method is thread-safe as it uses mutex locks to prevent race conditions.
*/
func GetMouseStatus() (int, int, uint, string) {
	MouseMemory.Lock()
	defer func() {
		MouseMemory.Unlock()
	}()
	return MouseMemory.xLocation, MouseMemory.yLocation, MouseMemory.buttonPressed,
		MouseMemory.wheelState
}

/*
GetPreviousMouseStatus allows you to retrieve the previous mouse status before the
most recent update. In addition, the following information should be noted:

- Returns the previous x location, y location, button pressed state, and wheel state.
- This method is thread-safe as it uses mutex locks to prevent race conditions.
*/
func GetPreviousMouseStatus() (int, int, uint, string) {
	MouseMemory.Lock()
	defer func() {
		MouseMemory.Unlock()
	}()
	return PreviousMouseMemory.xLocation, PreviousMouseMemory.yLocation,
		PreviousMouseMemory.buttonPressed, PreviousMouseMemory.wheelState
}

/*
WaitForClickRelease allows you to pause execution until the user releases any
currently pressed mouse buttons. In addition, the following information should be noted:

- This method will block until the button pressed state becomes 0 (no buttons pressed).
- This is useful for implementing drag and drop operations or waiting for user input.
*/
func WaitForClickRelease() {
	for MouseMemory.buttonPressed != 0 {
	}
}

/*
IsMouseInBoundingBox allows you to check if the current mouse position is within a
specified rectangular area. In addition, the following information should be noted:

- Returns true if the mouse is within the bounding box, false otherwise.
- The bounding box is defined by its top-left corner (xLocation, yLocation) and its
  dimensions (width, height).
- This is useful for detecting mouse hover or click events on UI elements.
*/
func IsMouseInBoundingBox(xLocation int, yLocation int, width int, height int) bool {
	mouseXLocation, mouseYLocation, _, _ := GetMouseStatus()
	if mouseXLocation >= xLocation && mouseXLocation <= xLocation+width {
		if mouseYLocation >= yLocation && mouseYLocation <= yLocation+height {
			return true
		}
	}
	return false
}
