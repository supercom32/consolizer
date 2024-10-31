package consolizer

import (
	"github.com/supercom32/consolizer/internal/memory"
	"time"
)

type TimerType struct {
	timerAlias string
}

// AddTimer creates and returns a new TimerType instance with a generated UUID.
func AddTimer(lengthOfTimerInMilliseconds int64, isTimerEnabled bool) *TimerType {
	timer := &TimerType{timerAlias: getUUID()}
	memory.AddTimer(timer.timerAlias, lengthOfTimerInMilliseconds, isTimerEnabled)
	return timer
}

/*
IsExpired allows you to check if a created timer has expired or not.
If the specified timer has expired, then it will automatically be disabled.
In order to activate the timer again, simply call 'StartTimer'.
*/
func (shared *TimerType) IsExpired() bool {
	timerEntry := memory.GetTimer(shared.timerAlias)
	if timerEntry.IsTimerEnabled {
		timeElapsed := GetCurrentTimeInMilliseconds() - timerEntry.StartTime
		if timeElapsed > timerEntry.TimerLength {
			timerEntry.IsTimerEnabled = false
			return true
		}
	}
	return false
}

/*
SetTimer allows you to create a new timer to measure time with. If the timer
is not enabled by default, you must call 'StartTimer' when you wish for it
to begin.
*/
func (shared *TimerType) SetTimer(durationInMilliseconds int64, isEnabled bool) {
	timerEntry := memory.GetTimer(shared.timerAlias)
	timerEntry.StartTime = GetCurrentTimeInMilliseconds()
	timerEntry.TimerLength = durationInMilliseconds
	timerEntry.IsTimerEnabled = isEnabled
}

/*
StartTimer allows you to start a timer that has already been previously
created. In addition, the following information should be noted:

- If you specify a timer that does not exist, then a panic will be
generated to fail as fast as possible.
*/
func (shared *TimerType) Start() {
	timerEntry := memory.GetTimer(shared.timerAlias)
	timerEntry.StartTime = GetCurrentTimeInMilliseconds()
	timerEntry.IsTimerEnabled = true
}

/*
Sleep allows you to pause execution for a given amount of milliseconds.
This method is simply a convenient wrapper for the method
'SleepInMilliseconds'.
*/
func Sleep(timeInMilliseconds uint) {
	SleepInMilliseconds(timeInMilliseconds)
}

/*
SleepInSeconds allows you to pause execution for a given amount of seconds.
*/
func SleepInSeconds(timeInSeconds uint) {
	SleepInMilliseconds(timeInSeconds * 1000)
}

/*
SleepInMilliseconds allows you to pause execution for a given amount of
milliseconds.
*/
func SleepInMilliseconds(timeInMilliseconds uint) {
	timeDuration := time.Duration(timeInMilliseconds)
	time.Sleep(timeDuration * time.Millisecond)
}

/*
GetCurrentTimeInMilliseconds allows you to get the current epoch
time in milliseconds.
*/
func GetCurrentTimeInMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
