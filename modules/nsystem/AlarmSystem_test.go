package nsystem

import (
	"testing"

	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"
	"github.com/yohamta/donburi"
)

func TestAlarmSystem(t *testing.T) {
	world := donburi.NewWorld()
	entity := world.Create(enginetype.Alarm)
	entry := world.Entry(entity)

	// Ensure data is initialized
	enginetype.InitializeComponent("ala", entry)

	data := enginetype.GetComponent(nil, enginetype.Alarm) // Need a proper object
	// Wait, GetComponent expects IObject (entryProvider)
	obj := &mockObject{entry: entry}
	data = enginetype.GetComponent(obj, enginetype.Alarm)
	if data == nil {
		t.Fatalf("Alarm data is nil")
	}

	callbackCalled := false
	// Add an alarm manually
	data.Alarms = append(data.Alarms, domain.Alarm{
		Id:         "test",
		FramesLeft: 2,
		IsActive:   true,
		Callback:   func() { callbackCalled = true },
	})

	sys := NewAlarmSystem()

	// Step 1: 2 -> 1
	sys.Update([]IObject{obj})
	if callbackCalled {
		t.Errorf("Callback should not be called yet")
	}

	// Step 2: 1 -> 0 (trigger)
	sys.Update([]IObject{obj})
	if !callbackCalled {
		t.Errorf("Callback should have been called")
	}
	if data.Alarms[0].IsActive {
		t.Errorf("Alarm should be inactive after triggering")
	}
}
