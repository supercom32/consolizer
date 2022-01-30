package memory

import "testing"

func TestMouseMemory(test *testing.T) {
	SetMouseStatus(1,2,3,"down")
	xLocation, yLocation, buttonPressed, wheelState := GetMouseStatus()
	if (xLocation != 1 || yLocation != 2 || buttonPressed != 3 || wheelState != "down") {
		test.Errorf("Mouse state was saved, but not read back with the same values.")
	}
	SetMouseStatus(4,5,6,"up")
	xLocation, yLocation, buttonPressed, wheelState = GetPreviousMouseStatus()
	if (xLocation != 1 || yLocation != 2 || buttonPressed != 3 || wheelState != "down") {
		test.Errorf("New mouse state was saved, but the previous mouse state did not read back with the right values.")
	}
}

func TestMouseMemoryClear(test *testing.T) {
	SetMouseStatus(1,2,3,"down")
	SetMouseStatus(4,5,6,"up")
	ClearMouseMemory()
	xLocation, yLocation, buttonPressed, wheelState := GetMouseStatus()
	if (xLocation != -1 || yLocation != -1 || buttonPressed != 0 || wheelState != "") {
		test.Errorf("Current mouse status was expected to be clear when it wasn't.")
	}
	xLocation, yLocation, buttonPressed, wheelState = GetPreviousMouseStatus()
	if (xLocation != -1 || yLocation != -1 || buttonPressed != 0 || wheelState != "") {
		test.Errorf("Previous mouse status was expected to be clear when it wasn't.")
	}
}

func TestIsMouseInBoundingBox(test *testing.T) {
	SetMouseStatus(2,2,3,"down")
	if IsMouseInBoundingBox(1, 1, 10, 10) != true {
		test.Errorf("Mouse was in the bounding box, but was not detected as such.")
	}
	SetMouseStatus(0,0,3,"down")
	if IsMouseInBoundingBox(1, 1, 10, 10) != false {
		test.Errorf("Mouse was out of the bounding box, but was not detected as such.")
	}
}