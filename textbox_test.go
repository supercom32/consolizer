package consolizer

import (
	"fmt"
	"testing"
)

func TestGetTextboxClickCoordinates(test *testing.T) {
	xLocation, yLocation := getTextboxClickCoordinates(39, 20)
	fmt.Printf("%d, %d\n", xLocation, yLocation)
}