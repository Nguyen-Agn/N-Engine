package components

import (
	"autoworld/domain"
	"autoworld/modules/enginetype"

	"github.com/yohamta/donburi"
)

var Alarm = enginetype.Alarm

func init() {
	enginetype.RegisterComponentInitializer("alr", func(entry *donburi.Entry) {
		donburi.SetValue(entry, enginetype.Alarm, domain.AlarmData{
			Alarms: make([]domain.Alarm, 0),
		})
	})
}

// AlarmComponent là Mixin để nhúng vào Custom Object.
type AlarmComponent struct {
	IObject
	data *domain.AlarmData
}

func (a *AlarmComponent) BindComponent(base IObject) {
	a.IObject = base
	a.data = enginetype.GetComponent(base, Alarm)
}

func (a AlarmComponent) SetAlarm(id string, frames int, callback func()) {
	if a.data == nil {
		return
	}
	// Update if exists
	for i := range a.data.Alarms {
		if a.data.Alarms[i].Id == id {
			a.data.Alarms[i].FramesLeft = frames
			a.data.Alarms[i].Callback = callback
			a.data.Alarms[i].IsActive = true
			return
		}
	}
	// Add new
	a.data.Alarms = append(a.data.Alarms, domain.Alarm{
		Id:         id,
		FramesLeft: frames,
		IsActive:   true,
		Callback:   callback,
	})
}

func (a AlarmComponent) GetAlarm(id string) int {
	if a.data == nil {
		return 0
	}
	for _, alarm := range a.data.Alarms {
		if alarm.Id == id && alarm.IsActive {
			return alarm.FramesLeft
		}
	}
	return 0
}

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
