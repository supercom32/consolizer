package memory

import "sync"

type mouseMemoryType struct {
	mutex         sync.Mutex
	xLocation     int
	yLocation     int
	buttonPressed uint
	wheelState    string
}

var MouseMemory mouseMemoryType
var PreviousMouseMemory mouseMemoryType

func ClearMouseMemory() {
	MouseMemory.mutex.Lock()
	MouseMemory.xLocation = -1
	MouseMemory.yLocation = -1
	MouseMemory.buttonPressed = 0
	MouseMemory.wheelState = ""
	MouseMemory.mutex.Unlock()
	PreviousMouseMemory.mutex.Lock()
	PreviousMouseMemory.xLocation = -1
	PreviousMouseMemory.yLocation = -1
	PreviousMouseMemory.buttonPressed = 0
	PreviousMouseMemory.wheelState = ""
	PreviousMouseMemory.mutex.Unlock()
}

func SetMouseStatus(xLocation int, yLocation int, buttonPressed uint, wheelState string) {
	PreviousMouseMemory.mutex.Lock()
	PreviousMouseMemory.xLocation = MouseMemory.xLocation
	PreviousMouseMemory.yLocation = MouseMemory.yLocation
	PreviousMouseMemory.buttonPressed = MouseMemory.buttonPressed
	PreviousMouseMemory.wheelState = MouseMemory.wheelState
	PreviousMouseMemory.mutex.Unlock()
	MouseMemory.mutex.Lock()
	MouseMemory.xLocation = xLocation
	MouseMemory.yLocation = yLocation
	MouseMemory.buttonPressed = buttonPressed
	MouseMemory.wheelState = wheelState
	MouseMemory.mutex.Unlock()
}

func GetMouseStatus() (int, int, uint, string) {
	return MouseMemory.xLocation, MouseMemory.yLocation, MouseMemory.buttonPressed, MouseMemory.wheelState
}

func GetPreviousMouseStatus() (int, int, uint, string) {
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
