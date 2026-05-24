package nsystem

import (
	"autoworld/modules/enginetype"
)

type AlarmSystem struct {
}

func NewAlarmSystem() *AlarmSystem {
	return &AlarmSystem{}
}

func (s *AlarmSystem) Update(objectList []IObject) {
	for _, obj := range objectList {
		data := enginetype.GetComponent(obj, enginetype.Alarm)
		if data == nil {
			continue
		}

		for i := range data.Alarms {
			if data.Alarms[i].IsActive {
				data.Alarms[i].FramesLeft--
				if data.Alarms[i].FramesLeft <= 0 {
					data.Alarms[i].IsActive = false
					if data.Alarms[i].Callback != nil {
						data.Alarms[i].Callback()
					}
				}
			}
		}
	}
}
