package nsystem

import (
	"autoworld/modules/enginetype"
)

type AlarmSystem struct {
}

// NewAlarmSystem creates and returns a new instance of AlarmSystem.
// Outputs: Returns a pointer to a newly initialized AlarmSystem.
func NewAlarmSystem() *AlarmSystem {
	return &AlarmSystem{}
}

// Update iterates through the list of objects and processes their alarms.
// Inputs: objectList ([]IObject) - The list of objects to be updated.
// Purpose: For each object, it retrieves the Alarm component data. If the component exists, it decrements the frames left for active alarms. Once an alarm's frames reach 0 or below, it deactivates the alarm and triggers its callback function if one is provided.
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
