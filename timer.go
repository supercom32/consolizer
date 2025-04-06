package consolizer

import (
	"fmt"
	"supercom32.net/consolizer/types"
	"sync"
	"time"
)

type timerMemoryType struct {
	sync.Mutex
	Entries map[string]*types.TimerEntryType
}

var Timers timerMemoryType

func InitializeTimerMemory() {
	Timers.Entries = make(map[string]*types.TimerEntryType)
}

func (shared *timerMemoryType) Add(timerAlias string, lengthOfTimerInMilliseconds int64, isTimerEnabled bool) {
	Timers.Lock()
	defer func() {
		Timers.Unlock()
	}()
	timerEntry := types.NewTimerEntry()
	timerEntry.IsTimerEnabled = isTimerEnabled
	timerEntry.StartTime = Timers.GetCurrentTimeInMilliseconds()
	timerEntry.TimerLength = lengthOfTimerInMilliseconds
	Timers.Entries[timerAlias] = &timerEntry
}

func (shared *timerMemoryType) Get(timerAlias string) *types.TimerEntryType {
	Timers.Lock()
	defer func() {
		Timers.Unlock()
	}()
	if Timers.Entries[timerAlias] == nil {
		panic(fmt.Sprintf("The requested timer with alias '%s' could not be returned since it does not exist.", timerAlias))
	}
	return Timers.Entries[timerAlias]
}

func (shared *timerMemoryType) Remove(timerAlias string) {
	Timers.Lock()
	defer func() {
		Timers.Unlock()
	}()
	delete(Timers.Entries, timerAlias)
}

func (shared *timerMemoryType) GetCurrentTimeInMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
