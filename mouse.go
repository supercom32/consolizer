package consolizer

import (
	"sync"
)

/*
mouseMemoryType is a class which allows you to mouseMemoryType is a structure that holds information about the current
mouse state including position, button status, and wheel state.
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
ClearMouseMemory is a method which allows you to reset both the current and previous mouse memory states to their
default values. In addition, the following information should be noted:

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
SetMouseStatus is a method which allows you to update the current mouse status while preserving the previous state. In
addition, the following information should be noted:

- The previous mouse state is updated with the current state before changing.

- The current mouse state is updated with the provided parameters.

- This method is thread-safe as it uses mutex locks to prevent race conditions.

:param xLocation: The x coordinate of the mouse position.
:param yLocation: The y coordinate of the mouse position.
:param buttonPressed: The bitmask of currently pressed mouse buttons.
:param wheelState: The current state of the mouse wheel.
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
GetMouseStatus is a method which allows you to retrieve the current mouse status including position, button state, and
wheel state. In addition, the following information should be noted:

- This method is thread-safe as it uses mutex locks to prevent race conditions.

:return: The x location, y location, button pressed state, and wheel state of the mouse.
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
GetPreviousMouseStatus is a method which allows you to retrieve the previous mouse status before the most recent update.
In addition, the following information should be noted:

- This method is thread-safe as it uses mutex locks to prevent race conditions.

:return: The previous x location, y location, button pressed state, and wheel state of the mouse.
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
WaitForClickRelease is a method which allows you to pause execution until the user releases any currently pressed mouse
buttons. In addition, the following information should be noted:

- This method will block until the button pressed state becomes 0 (no buttons pressed).

- This is useful for implementing drag and drop operations or waiting for user input.
*/
func WaitForClickRelease() {
	for MouseMemory.buttonPressed != 0 {
	}
}

/*
IsMouseInBoundingBox is a method which allows you to check if the current mouse position is within a specified
rectangular area. In addition, the following information should be noted:

- The bounding box is defined by its top-left corner (xLocation, yLocation) and its dimensions (width, height).

- This is useful for detecting mouse hover or click events on UI elements.

:param xLocation: The x coordinate of the top-left corner of the bounding box.
:param yLocation: The y coordinate of the top-left corner of the bounding box.
:param width: The width of the bounding box.
:param height: The height of the bounding box.

:return: A boolean indicating whether the mouse is within the specified bounding box.
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

/*
GetLayerUnderMouseCursor is a method which allows you to retrieve the instance of the layer under the current mouse
cursor position. In addition, the following information should be noted:

- We create a new instance of the layer so the user can interact with it.

- This is a new instance, so any changes to the instance itself will not be reflected in the original layer system.

- Methods called on this instance will affect the original layer.

:return: An instance of the layer under the mouse cursor, or nil if no layer is found.
*/
func GetLayerUnderMouseCursor() *LayerInstanceType {
	mouseX, mouseY, _, _ := GetMouseStatus()
	if mouseX < 0 || mouseX >= commonResource.terminalWidth || mouseY < 0 || mouseY >= commonResource.terminalHeight {
		return nil
	}
	layerAlias := commonResource.screenLayer.CharacterMemory[mouseY][mouseX].LayerAlias
	if layerAlias == "" {
		return nil
	}
	layerEntry := Layers.Get(layerAlias)
	if layerEntry == nil {
		return nil
	}
	// We create a new instance of the layer so the user can interact with it.
	// Note that this is a new instance, so any changes to the instance itself
	// will not be reflected in the original layer system. However, methods
	// called on this instance will affect the original layer.
	layerInstance := &LayerInstanceType{
		layerAlias:  layerEntry.LayerAlias,
		parentAlias: layerEntry.ParentAlias,
	}
	return layerInstance
}
