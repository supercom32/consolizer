package types

import (
	"sync"
)

type BaseControlType struct {
	Mutex         *sync.Mutex
	StyleEntry    TuiStyleEntryType
	Alias         string
	XLocation     int
	YLocation     int
	Width         int
	Height        int
	IsEnabled     bool
	IsVisible     bool
	TabIndex      int
	Label         string
	IsSelected    bool
	IsBorderDrawn bool
}

func NewBaseControl() BaseControlType {
	var baseControl BaseControlType
	baseControl.Mutex = &sync.Mutex{}
	baseControl.StyleEntry = NewTuiStyleEntry()
	baseControl.Alias = ""
	baseControl.XLocation = 0
	baseControl.YLocation = 0
	baseControl.Width = 0
	baseControl.Height = 0
	baseControl.IsEnabled = true
	baseControl.IsVisible = true
	baseControl.TabIndex = 0
	baseControl.Label = ""
	baseControl.IsSelected = false
	baseControl.IsBorderDrawn = false
	return baseControl
}

func (shared *BaseControlType) GetBounds() (int, int, int, int) {
	return shared.XLocation, shared.YLocation, shared.Width, shared.Height
}

func (shared *BaseControlType) SetPosition(x, y int) {
	shared.XLocation = x
	shared.YLocation = y
}

func (shared *BaseControlType) SetSize(width, height int) {
	shared.Width = width
	shared.Height = height
}

func (shared *BaseControlType) SetEnabled(enabled bool) {
	shared.IsEnabled = enabled
}

func (shared *BaseControlType) SetVisible(visible bool) {
	shared.IsVisible = visible
}

func (shared *BaseControlType) SetStyle(style TuiStyleEntryType) {
	shared.StyleEntry = style
}

func (shared *BaseControlType) SetTabIndex(index int) {
	shared.TabIndex = index
}

func (shared *BaseControlType) Lock() {
	if shared.Mutex != nil {
		shared.Mutex.Lock()
	}
}

func (shared *BaseControlType) Unlock() {
	if shared.Mutex != nil {
		shared.Mutex.Unlock()
	}
}
