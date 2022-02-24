package consolizer

import (
	"fmt"
	"testing"
)

func TestGetTextboxClickCoordinates(test *testing.T) {
	xLocation, yLocation := getTextboxClickCoordinates(99, 10, 10)
	fmt.Printf("%d, %d\n", xLocation, yLocation)
}