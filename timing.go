package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/memory"
	"time"

	"github.com/supercom32/consolizer/types"
)

var Timers *memory.MemoryManager[types.TimerEntryType]

/*
InitializeTimerMemory is a method which initializes the memory manager for timers.

Example:
    InitializeTimerMemory()
*/
func InitializeTimerMemory() {
	Timers = memory.NewMemoryManager[types.TimerEntryType]()
}

/*
TimerType is a structure which represents a timer instance.
*/
type TimerType struct {
	timerAlias string
}

/*
AddTimer is a method which creates and returns a new TimerType instance with a generated UUID.

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
IsExpired is a method which checks if a created timer has expired or not. In addition, the following should be noted:

- If the specified timer has expired, then it will automatically be disabled.

- In order to activate the timer again, simply call 'Start'.

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
Set is a method which configures a timer with a new duration and enabled state. In addition, the following should be noted:

- If the timer is not enabled by default, you must call 'Start' when you wish for it to begin.

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
Start is a method which starts a timer that has already been previously created. In addition, the following should be noted:

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
Sleep is a method which pauses execution for a given amount of milliseconds. In addition, the following should be noted:

- This method is simply a convenient wrapper for the method 'SleepInMilliseconds'.

Example:
    Sleep(1000)
*/
func Sleep(timeInMilliseconds uint) {
	SleepInMilliseconds(timeInMilliseconds)
}

/*
SleepInSeconds is a method which pauses execution for a given amount of seconds.

Example:
    SleepInSeconds(2)
*/
func SleepInSeconds(timeInSeconds uint) {
	SleepInMilliseconds(timeInSeconds * 1000)
}

/*
SleepInMilliseconds is a method which pauses execution for a given amount of milliseconds.

Example:
    SleepInMilliseconds(500)
*/
func SleepInMilliseconds(timeInMilliseconds uint) {
	timeDuration := time.Duration(timeInMilliseconds)
	time.Sleep(timeDuration * time.Millisecond)
}

/*
GetCurrentTimeInMilliseconds is a method which gets the current epoch time in milliseconds.

Example:
    now := GetCurrentTimeInMilliseconds()
*/
func GetCurrentTimeInMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

/*
IsTimerExists is a method which checks if a timer with a given alias exists.

Example:
    exists := IsTimerExists("myTimer")
*/
func IsTimerExists(timerAlias string) bool {
	return Timers.IsExists(timerAlias)
}

/*
GetAllTimers is a method which retrieves all timers currently in memory.

Example:
    allTimers := GetAllTimers()
*/
func GetAllTimers() map[string]*types.TimerEntryType {
	return Timers.GetAllEntriesWithKeys()
}

/*
RemoveAllTimers is a method which removes all timers from memory.

Example:
    RemoveAllTimers()
*/
func RemoveAllTimers() {
	Timers.RemoveAll()
}
