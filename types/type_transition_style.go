package types

import (
	"github.com/supercom32/consolizer/constants"
)

/*
TransitionStyleEntryType is a structure which contains the configuration for a spatial transition effect.

Example:
    transitionStyle := NewTransitionStyleEntry()
*/
type TransitionStyleEntryType struct {
	TransitionType constants.TransitionType
	Direction      constants.TransitionDirection
	SoftEdgeWidth  float32
	BlindCount     int
}

/*
NewTransitionStyleEntry is a constructor which creates a new transition style entry with default values.

Example:
    transitionStyle := NewTransitionStyleEntry()
*/
func NewTransitionStyleEntry(existingTransitionStyleEntry ...*TransitionStyleEntryType) TransitionStyleEntryType {
	var transitionStyle TransitionStyleEntryType
	if existingTransitionStyleEntry != nil {
		transitionStyle.TransitionType = existingTransitionStyleEntry[0].TransitionType
		transitionStyle.Direction = existingTransitionStyleEntry[0].Direction
		transitionStyle.SoftEdgeWidth = existingTransitionStyleEntry[0].SoftEdgeWidth
		transitionStyle.BlindCount = existingTransitionStyleEntry[0].BlindCount
	} else {
		transitionStyle.TransitionType = constants.TransitionTypeNone
		transitionStyle.Direction = constants.TransitionDirectionLeftToRight
		transitionStyle.SoftEdgeWidth = 0.1
		transitionStyle.BlindCount = 6
	}
	return transitionStyle
}
