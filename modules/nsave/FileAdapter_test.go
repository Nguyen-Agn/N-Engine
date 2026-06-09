package nsave

import (
	"autoworld/domain"
	"os"
	"path/filepath"
	"testing"
)

func TestJsonFileAdapter_WriteReadDelete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "nsave_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	adapter := newJsonFileAdapter(tmpDir)
	savePath := filepath.Join(tmpDir, "test_save.json")

	// 1. Test Write
	snap := domain.SaveSnapshot{
		Version:        1,
		CurrentSceneID: "Level1",
	}

	err = adapter.Write(savePath, snap)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// 2. Test Exists
	if !adapter.Exists(savePath) {
		t.Errorf("Exists should be true after Write")
	}

	// 3. Test Read
	readSnap, err := adapter.Read(savePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if readSnap.CurrentSceneID != "Level1" || readSnap.Version != 1 {
		t.Errorf("Read data mismatch. Got %v", readSnap)
	}

	// 4. Test Delete
	err = adapter.Delete(savePath)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if adapter.Exists(savePath) {
		t.Errorf("Exists should be false after Delete")
	}
}

func TestJsonFileAdapter_ListAll(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "nsave_test_list")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	adapter := newJsonFileAdapter(tmpDir)
	
	// Create some fake files
	adapter.Write(filepath.Join(tmpDir, "save1.json"), domain.SaveSnapshot{})
	adapter.Write(filepath.Join(tmpDir, "save2.json"), domain.SaveSnapshot{})
	os.WriteFile(filepath.Join(tmpDir, "not_a_save.txt"), []byte("dummy"), 0644)

	slots := adapter.ListAll()
	if len(slots) != 2 {
		t.Errorf("ListAll() = %v slots, want 2", len(slots))
	}
	
	hasSave1 := false
	hasSave2 := false
	for _, s := range slots {
		if s == "save1" { hasSave1 = true }
		if s == "save2" { hasSave2 = true }
	}
	
	if !hasSave1 || !hasSave2 {
		t.Errorf("ListAll() missing saves. Got %v", slots)
	}
}
