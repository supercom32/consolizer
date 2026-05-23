package consolizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
TestTimerFunctionality is a test which verifies the basic functionality of a timer, including expiration after a set duration.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        None
*/
func TestTimerFunctionality(test *testing.T) {
	InitializeTimerMemory()
	timerEntry := AddTimer(3000, true)
	SleepInMilliseconds(1000)
	assert.Equalf(test, timerEntry.IsExpired(), false, "The timer was flagged as expired when not enough time has elapsed.")
	SleepInMilliseconds(2500)
	assert.Equalf(test, timerEntry.IsExpired(), true, "The timer was not flagged as expired when more time has elapsed.")
}

/*
TestSetTimerFunctionality is a test which verifies that a timer can be correctly set with new values.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        None
*/
func TestSetTimerFunctionality(test *testing.T) {
	InitializeTimerMemory()
	timerEntry := AddTimer(9000, false)
	timerEntry.Set(3000, true)
	SleepInSeconds(1)
	assert.Equalf(test, timerEntry.IsExpired(), false, "The timer was flagged as expired when not enough time has elapsed.")

	Sleep(2500)
	assert.Equalf(test, timerEntry.IsExpired(), true, "The timer was not flagged as expired when more time has elapsed.")
}

/*
TestResetTimerFunctionality is a test which verifies that a timer can be correctly reset and restarted.

Example:
    Expected Inputs:
        None

    Expected Outputs:
        None
*/
func TestResetTimerFunctionality(test *testing.T) {
	InitializeTimerMemory()
	timerEntry := AddTimer(1000, true)
	SleepInMilliseconds(1500)
	assert.Equalf(test, true, timerEntry.IsExpired(), "The initial timer was not flagged as expired when more time has elapsed.")
	timerEntry.Start()
	SleepInMilliseconds(500)
	assert.Equalf(test, false, timerEntry.IsExpired(), "The reset timer was flagged as expired when not enough time has elapsed.")
	SleepInMilliseconds(1000)
	assert.Equalf(test, true, timerEntry.IsExpired(), "The reset timer was not flagged as expired when more time has elapsed.")
}

