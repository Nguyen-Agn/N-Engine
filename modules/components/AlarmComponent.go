package components

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"

	"github.com/yohamta/donburi"
)

var Alarm = enginetype.Alarm

// init initializes the default data for the alarm component.
// It registers the "alr" component token with empty alarm list.
func init() {
	enginetype.RegisterComponentInitializer("alr", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Alarm, domain.AlarmData{
			Alarms: make([]domain.Alarm, 0),
		})
	})
}

// AlarmComponent is a Mixin component to embed into a Custom Object.
// It provides functionality for scheduling actions using frame-based timers.
type AlarmComponent struct {
	IObject
	data *domain.AlarmData
}

// BindComponent binds the base object and its ECS data to this component.
// Inputs:
//   - base: The base IObject to bind to.
func (a *AlarmComponent) BindComponent(base IObject) {
	a.IObject = base
	a.data = enginetype.GetComponent(base, Alarm)
}

// SetAlarm sets a timer to run a callback function after a certain number of frames.
// Purpose: Schedules a function to execute in the future.
// Inputs:
//   - name: The unique identifier for the timer.
//   - frames: The duration to wait, measured in frames.
//   - callback: The parameterless function to execute when the timer finishes.
func (a AlarmComponent) SetAlarm(name string, frames int, callback func()) {
	if a.data == nil {
		return
	}
	// Update if exists
	for i := range a.data.Alarms {
		if a.data.Alarms[i].Id == name {
			a.data.Alarms[i].FramesLeft = frames
			a.data.Alarms[i].Callback = callback
			a.data.Alarms[i].IsActive = true
			return
		}
	}
	// Add new
	a.data.Alarms = append(a.data.Alarms, domain.Alarm{
		Id:         name,
		FramesLeft: frames,
		IsActive:   true,
		Callback:   callback,
	})
}

// GetAlarm retrieves the remaining time of an active alarm.
// Purpose: Checks how many frames are left before the alarm triggers.
// Inputs:
//   - name: The unique identifier of the alarm.
//
// Outputs: Returns the remaining frames as an integer, or 0 if not found or inactive.
func (a AlarmComponent) GetAlarm(name string) int {
	if a.data == nil {
		return 0
	}
	for _, alarm := range a.data.Alarms {
		if alarm.Id == name && alarm.IsActive {
			return alarm.FramesLeft
		}
	}
	return 0
}

// StopAlarm stops an active alarm by its identifier.
// Purpose: Cancels a scheduled alarm so its callback will not run.
// Inputs:
//   - id: The unique identifier of the alarm to stop.
func (a AlarmComponent) StopAlarm(id string) {
	if a.data == nil {
		return
	}
	for i := range a.data.Alarms {
		if a.data.Alarms[i].Id == id {
			a.data.Alarms[i].IsActive = false
			return
		}
	}
}
