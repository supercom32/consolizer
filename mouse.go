package consolizer

import (
	"sync"
)

type mouseMemoryType struct {
	sync.Mutex
	xLocation     int
	yLocation     int
	buttonPressed uint
	wheelState    string
}

var MouseMemory mouseMemoryType
var PreviousMouseMemory mouseMemoryType

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

func GetMouseStatus() (int, int, uint, string) {
	MouseMemory.Lock()
	defer func() {
		MouseMemory.Unlock()
	}()
	return MouseMemory.xLocation, MouseMemory.yLocation, MouseMemory.buttonPressed, MouseMemory.wheelState
}

func GetPreviousMouseStatus() (int, int, uint, string) {
	MouseMemory.Lock()
	defer func() {
		MouseMemory.Unlock()
	}()
	return PreviousMouseMemory.xLocation, PreviousMouseMemory.yLocation, PreviousMouseMemory.buttonPressed, PreviousMouseMemory.wheelState
}

func WaitForClickRelease() {
	for MouseMemory.buttonPressed != 0 {
	}
}

func IsMouseInBoundingBox(xLocation int, yLocation int, width int, height int) bool {
	mouseXLocation, mouseYLocation, _, _ := GetMouseStatus()
	if mouseXLocation >= xLocation && mouseXLocation <= xLocation+width {
		if mouseYLocation >= yLocation && mouseYLocation <= yLocation+height {
			return true
		}
	}
	return false
}
