package types

import (
	"encoding/json"
)

/*
TimerEntryType is a structure which represents a timer control.

Example:

	var timer types.TimerEntryType
*/
type TimerEntryType struct {
	IsTimerEnabled bool
	StartTime      int64
	TimerLength    int64
}

/*
MarshalJSON is a method which allows you to convert a timer entry to JSON format. In addition, the following should be noted:

- Implements the json.Marshaler interface for TimerEntryType.

- Converts the timer entry to a JSON string representation.

- Used for serializing timer entries when saving state or transmitting data.

Example:

	instance.MarshalJSON()
*/
func (shared TimerEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		IsTimerEnabled bool
		StartTime      int64
		TimerLength    int64
	}{
		IsTimerEnabled: shared.IsTimerEnabled,
		StartTime:      shared.StartTime,
		TimerLength:    shared.TimerLength,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump is a method which allows you to get a JSON string representation of the timer entry. In addition, the following should be noted:

- Returns a formatted JSON string of the timer entry.

- Used for debugging and logging purposes.

- The output is human-readable and includes all timer properties.

Example:

	instance.GetEntryAsJsonDump()
*/
func (shared TimerEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewTimerEntry is a constructor which allows you to create a new timer entry. In addition, the following should be noted:

- Initializes a timer entry with default values.

- Sets up a timer with a specified duration and callback function.

- The timer can be started, stopped, and reset using the appropriate methods.

Example:

	NewTimerEntry(existingTimerEntry)
*/
func NewTimerEntry(existingTimerEntry ...*TimerEntryType) TimerEntryType {
	var timerEntry TimerEntryType
	if existingTimerEntry != nil {
		timerEntry.IsTimerEnabled = existingTimerEntry[0].IsTimerEnabled
		timerEntry.StartTime = existingTimerEntry[0].StartTime
		timerEntry.TimerLength = existingTimerEntry[0].TimerLength
	}
	return timerEntry
}
