package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
	"time"
)

type timerMemoryType struct {
	sync.Mutex
	Entries map[string]*types.TimerEntryType
}

var Timer timerMemoryType

func InitializeTimerMemory() {
	Timer.Entries = make(map[string]*types.TimerEntryType)
}

func AddTimer(timerAlias string, lengthOfTimerInMilliseconds int64, isTimerEnabled bool) {
	Timer.Lock()
	defer func() {
		Timer.Unlock()
	}()
	timerEntry := types.NewTimerEntry()
	timerEntry.IsTimerEnabled = isTimerEnabled
	timerEntry.StartTime = getCurrentTimeInMilliseconds()
	timerEntry.TimerLength = lengthOfTimerInMilliseconds
	Timer.Entries[timerAlias] = &timerEntry
}

func GetTimer(timerAlias string) *types.TimerEntryType {
	Timer.Lock()
	defer func() {
		Timer.Unlock()
	}()
	if Timer.Entries[timerAlias] == nil {
		panic(fmt.Sprintf("The requested timer with alias '%s' could not be returned since it does not exist.", timerAlias))
	}
	return Timer.Entries[timerAlias]
}

func DeleteTimer(timerAlias string) {
	Timer.Lock()
	defer func() {
		Timer.Unlock()
	}()
	delete(Timer.Entries, timerAlias)
}

func getCurrentTimeInMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
