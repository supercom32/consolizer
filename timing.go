package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/memory"
	"time"

	"github.com/supercom32/consolizer/types"
)

var Timers *memory.MemoryManager[types.TimerEntryType]

/*
InitializeTimerMemory is a method which allows you to initialize the memory manager for timers.

Example:

	InitializeTimerMemory()
*/
func InitializeTimerMemory() {
	Timers = memory.NewMemoryManager[types.TimerEntryType]()
}

type TimerType struct {
	timerAlias string
}

/*
AddTimer is a method which allows you to create and return a new TimerType instance with a generated UUID.

:param lengthOfTimerInMilliseconds: The length of the timer in milliseconds.
:param isTimerEnabled: Whether the timer should be enabled by default.

:return: A pointer to the created TimerType instance.

Example:

	timer := AddTimer(1000, true)
*/
func AddTimer(lengthOfTimerInMilliseconds int64, isTimerEnabled bool) *TimerType {
	timer := &TimerType{timerAlias: getUUID()}
	timerEntry := types.NewTimerEntry()
	timerEntry.IsTimerEnabled = isTimerEnabled
	timerEntry.StartTime = GetCurrentTimeInMilliseconds()
	timerEntry.TimerLength = lengthOfTimerInMilliseconds
	Timers.Add(timer.timerAlias, &timerEntry)
	return timer
}

/*
IsExpired is a method which allows you to check if a created timer has expired or not. In addition, the following should
be noted:

- If the specified timer has expired, then it will automatically be disabled.

- In order to activate the timer again, simply call 'Start'.

:return: True if the timer has expired, false otherwise.

Example:

	if timer.IsExpired() {
	    fmt.Println("Timer expired!")
	}
*/
func (shared *TimerType) IsExpired() bool {
	timerEntry := Timers.Get(shared.timerAlias)
	if timerEntry == nil {
		safeSttyPanic(fmt.Sprintf("The requested timer with alias '%s' could not be returned since it does not exist.", shared.timerAlias))
	}
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
Set is a method which allows you to configure a timer with a new duration and enabled state. In addition, the
following should be noted:

- If the timer is not enabled by default, you must call 'Start' when you wish for it to begin.

:param durationInMilliseconds: The new duration for the timer in milliseconds.
:param isEnabled: Whether the timer should be enabled.

Example:

	timer.Set(2000, true)
*/
func (shared *TimerType) Set(durationInMilliseconds int64, isEnabled bool) {
	timerEntry := Timers.Get(shared.timerAlias)
	if timerEntry == nil {
		safeSttyPanic(fmt.Sprintf("The requested timer with alias '%s' could not be returned since it does not exist.", shared.timerAlias))
	}
	timerEntry.StartTime = GetCurrentTimeInMilliseconds()
	timerEntry.TimerLength = durationInMilliseconds
	timerEntry.IsTimerEnabled = isEnabled
}

/*
Start is a method which allows you to start a timer that has already been previously created. In addition, the following
should be noted:

- If you attempt to start a timer that does not exist, then a panic will be generated to fail as fast as possible.

Example:

	timer.Start()
*/
func (shared *TimerType) Start() {
	timerEntry := Timers.Get(shared.timerAlias)
	if timerEntry == nil {
		safeSttyPanic(fmt.Sprintf("The requested timer with alias '%s' could not be returned since it does not exist.", shared.timerAlias))
	}
	timerEntry.StartTime = GetCurrentTimeInMilliseconds()
	timerEntry.IsTimerEnabled = true
}

/*
Sleep is a method which allows you to pause execution for a given amount of milliseconds. In addition, the following
should be noted:

- This method is simply a convenient wrapper for the method 'SleepInMilliseconds'.

:param timeInMilliseconds: The amount of time to sleep in milliseconds.

Example:

	Sleep(1000)
*/
func Sleep(timeInMilliseconds uint) {
	SleepInMilliseconds(timeInMilliseconds)
}

/*
SleepInSeconds is a method which allows you to pause execution for a given amount of seconds.

:param timeInSeconds: The amount of time to sleep in seconds.

Example:

	SleepInSeconds(2)
*/
func SleepInSeconds(timeInSeconds uint) {
	SleepInMilliseconds(timeInSeconds * 1000)
}

/*
SleepInMilliseconds is a method which allows you to pause execution for a given amount of milliseconds.

:param timeInMilliseconds: The amount of time to sleep in milliseconds.

Example:

	SleepInMilliseconds(500)
*/
func SleepInMilliseconds(timeInMilliseconds uint) {
	timeDuration := time.Duration(timeInMilliseconds)
	time.Sleep(timeDuration * time.Millisecond)
}

/*
GetCurrentTimeInMilliseconds is a method which allows you to get the current epoch time in milliseconds.

:return: The current time in milliseconds.

Example:

	now := GetCurrentTimeInMilliseconds()
*/
func GetCurrentTimeInMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

/*
IsTimerExists is a method which allows you to check if a timer with a given alias exists.

:param timerAlias: The alias of the timer to check.

:return: True if the timer exists, false otherwise.

Example:

	exists := IsTimerExists("myTimer")
*/
func IsTimerExists(timerAlias string) bool {
	return Timers.IsExists(timerAlias)
}

/*
GetAllTimers is a method which allows you to retrieve all timers currently in memory.

:return: A map of timer aliases to their corresponding timer entries.

Example:

	allTimers := GetAllTimers()
*/
func GetAllTimers() map[string]*types.TimerEntryType {
	return Timers.GetAllEntriesWithKeys()
}

/*
RemoveAllTimers is a method which allows you to remove all timers from memory.

Example:

	RemoveAllTimers()
*/
func RemoveAllTimers() {
	Timers.RemoveAll()
}
