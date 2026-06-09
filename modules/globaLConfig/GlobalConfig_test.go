package globalconfig

import (
	"testing"

	"github.com/Nguyen-Agn/N-Engine/domain"
)

// Mock observer for testing
type mockObserver struct {
	notified bool
}

func (m *mockObserver) NotifyChange(source domain.IGlobalConfig) {
	m.notified = true
}

func TestGlobalConfig_BasicSetGet(t *testing.T) {
	cfg := NewGlobalConfig()

	// Test Set and Get
	cfg.SetValue("score", 100)
	val := cfg.GetValue("score")
	if val.(int) != 100 {
		t.Errorf("GetValue('score') = %v; want 100", val)
	}

	cfg.SetValue("name", "Player1")
	if cfg.GetValue("name").(string) != "Player1" {
		t.Errorf("GetValue('name') = %v; want 'Player1'", cfg.GetValue("name"))
	}
}

func TestGlobalConfig_Constants(t *testing.T) {
	cfg := NewGlobalConfig()

	ok := cfg.NewConst("MAX_HEALTH", 500)
	if !ok {
		t.Errorf("NewConst should return true for new constant")
	}

	val := cfg.GetConst("MAX_HEALTH")
	if val.(int) != 500 {
		t.Errorf("GetConst('MAX_HEALTH') = %v; want 500", val)
	}

	// Try to overwrite const
	ok2 := cfg.NewConst("MAX_HEALTH", 1000)
	if ok2 {
		t.Errorf("NewConst should return false when overwriting existing constant")
	}

	// Ensure value is unchanged
	val2 := cfg.GetConst("MAX_HEALTH")
	if val2.(int) != 500 {
		t.Errorf("GetConst('MAX_HEALTH') after failed overwrite = %v; want 500", val2)
	}
}

func TestGlobalConfig_Observer(t *testing.T) {
	cfg := NewGlobalConfig()
	obs := &mockObserver{}

	cfg.AddObserver(obs)
	cfg.NotifyChange()

	if !obs.notified {
		t.Errorf("Observer was not notified")
	}

	obs.notified = false
	cfg.RemoveObserver(obs)
	cfg.NotifyChange()

	if obs.notified {
		t.Errorf("Observer was notified after being removed")
	}
}

func TestGlobalConfig_DumpRestore(t *testing.T) {
	cfg := NewGlobalConfig()
	cfg.SetValue("level", 5)
	cfg.SetValue("gold", 1000)

	dump := cfg.DumpVariables()
	if dump["level"].(int) != 5 || dump["gold"].(int) != 1000 {
		t.Errorf("DumpVariables did not contain correct data")
	}

	// Change original values
	cfg.SetValue("level", 10)
	cfg.SetValue("gold", 0)

	// Restore
	cfg.RestoreVariables(dump)

	if cfg.GetValue("level").(int) != 5 || cfg.GetValue("gold").(int) != 1000 {
		t.Errorf("RestoreVariables did not restore correct data: %v", cfg.DumpVariables())
	}
}
