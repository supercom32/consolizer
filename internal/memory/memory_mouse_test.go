package memory

import (
	"supercom32.net/consolizer"
	"testing"
)

func TestMouseMemory(test *testing.T) {
	consolizer.SetMouseStatus(1, 2, 3, "down")
	xLocation, yLocation, buttonPressed, wheelState := consolizer.GetMouseStatus()
	if xLocation != 1 || yLocation != 2 || buttonPressed != 3 || wheelState != "down" {
		test.Errorf("Mouse state was saved, but not read back with the same values.")
	}
	consolizer.SetMouseStatus(4, 5, 6, "up")
	xLocation, yLocation, buttonPressed, wheelState = consolizer.GetPreviousMouseStatus()
	if xLocation != 1 || yLocation != 2 || buttonPressed != 3 || wheelState != "down" {
		test.Errorf("New mouse state was saved, but the previous mouse state did not read back with the right values.")
	}
}

func TestMouseMemoryClear(test *testing.T) {
	consolizer.SetMouseStatus(1, 2, 3, "down")
	consolizer.SetMouseStatus(4, 5, 6, "up")
	consolizer.ClearMouseMemory()
	xLocation, yLocation, buttonPressed, wheelState := consolizer.GetMouseStatus()
	if xLocation != -1 || yLocation != -1 || buttonPressed != 0 || wheelState != "" {
		test.Errorf("Current mouse status was expected to be clear when it wasn't.")
	}
	xLocation, yLocation, buttonPressed, wheelState = consolizer.GetPreviousMouseStatus()
	if xLocation != -1 || yLocation != -1 || buttonPressed != 0 || wheelState != "" {
		test.Errorf("Previous mouse status was expected to be clear when it wasn't.")
	}
}

func TestIsMouseInBoundingBox(test *testing.T) {
	consolizer.SetMouseStatus(2, 2, 3, "down")
	if consolizer.IsMouseInBoundingBox(1, 1, 10, 10) != true {
		test.Errorf("Mouse was in the bounding box, but was not detected as such.")
	}
	consolizer.SetMouseStatus(0, 0, 3, "down")
	if consolizer.IsMouseInBoundingBox(1, 1, 10, 10) != false {
		test.Errorf("Mouse was out of the bounding box, but was not detected as such.")
	}
}
