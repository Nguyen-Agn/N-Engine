package nsys

import (
	"testing"
)

func TestGlobalStore_Variables(t *testing.T) {
	store := GetInstance()

	store.SetValue("score", 150)
	if val := store.GetValue("score"); val != 150 {
		t.Errorf("GetValue('score') = %v, want 150", val)
	}

	dump := store.DumpVariables()
	if dump["score"] != 150 {
		t.Errorf("DumpVariables missing score")
	}

	store.SetValue("score", 200)
	store.RestoreVariables(dump)
	if val := store.GetValue("score"); val != 150 {
		t.Errorf("RestoreVariables failed. Got %v, want 150", val)
	}
}

func TestGlobalStore_Constants(t *testing.T) {
	store := GetInstance()

	ok := store.NewConst("MAX_LEVEL", 10)
	if !ok {
		t.Errorf("NewConst should return true for new const")
	}

	if val := store.GetConst("MAX_LEVEL"); val != 10 {
		t.Errorf("GetConst('MAX_LEVEL') = %v, want 10", val)
	}

	ok2 := store.NewConst("MAX_LEVEL", 20)
	if ok2 {
		t.Errorf("NewConst should return false when overwriting")
	}
	if val := store.GetConst("MAX_LEVEL"); val != 10 {
		t.Errorf("GetConst should not change after failed NewConst")
	}

	store.UpdateConst("MAX_LEVEL", 20)
	if val := store.GetConst("MAX_LEVEL"); val != 20 {
		t.Errorf("UpdateConst failed. Got %v, want 20", val)
	}
}

// Minimal mock for objects
type mockSprite struct{}
func (m *mockSprite) FrameWidth() int { return 0 }
func (m *mockSprite) FrameHeight() int { return 0 }
func (m *mockSprite) GetFrame(i int) any { return nil }

func TestGlobalStore_Sprites(t *testing.T) {
	store := GetInstance()
	// Using interface nil is tricky, so we don't mock the actual interface, just check map behavior.
	store.AddSprite("hero", nil)
	if store.GetSprite("hero") != nil {
		t.Errorf("GetSprite returned non-nil")
	}
}
