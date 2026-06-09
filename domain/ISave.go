package domain

// ISaveManager l interface c?a h? th?ng Save/Load ton c?c.
type ISaveManager interface {
	// SaveGame luu ton b? tr?ng thi vo file theo du?ng d?n.
	// Purpose: Serializes and saves the current global game state to a file.
	// Inputs: path string - The file path where the save data will be written.
	// Outputs: error - Returns an error if file creation or serialization fails.
	SaveGame(path string) error

	// LoadGame t?i tr?ng thi t? file theo du?ng d?n v p d?ng vo game.
	// Purpose: Loads a saved game state from a file and applies it to the current running game.
	// Inputs: path string - The file path to load from.
	// Outputs: error - Returns an error if the file cannot be read, parsed, or applied.
	LoadGame(path string) error

	// HasSave ki?m tra xem file save c t?n t?i khng.
	// Purpose: Checks if a save file exists at the specified path.
	// Inputs: path string - The file path to check.
	// Outputs: bool - True if the save file exists.
	HasSave(path string) bool

	// DeleteSave xa save file t?i du?ng d?n cung c?p.
	// Purpose: Deletes a specific save file.
	// Inputs: path string - The file path of the save to delete.
	// Outputs: error - Returns an error if deletion fails.
	DeleteSave(path string) error

	// ListSlots tr? v? danh sch cc t?n file save trong thu m?c luu m?c d?nh.
	// Purpose: Retrieves a list of available save files in the default save directory.
	// Inputs: None.
	// Outputs: []string - A slice containing the filenames of all saves found.
	ListSlots() []string

	// ReadSnapshot d?c d? li?u t? file nhung khng p d?ng.
	// Purpose: Reads a save file and deserializes its data without applying it to the game state.
	// Inputs: path string - The file path to read from.
	// Outputs:
	//   1. SaveSnapshot - The parsed snapshot data structure.
	//   2. error - Returns an error if the read or deserialization fails.
	ReadSnapshot(path string) (SaveSnapshot, error)
}

// SaveSnapshot l c?u trc d? li?u du?c serialize ra file JSON.
type SaveSnapshot struct {
	Version        int                           `json:"version"`
	Path           string                        `json:"path"`
	SavedAt        int64                         `json:"saved_at"`
	CurrentSceneID string                        `json:"current_scene_id"`
	Variables      map[string]any                `json:"variables"` // T? IGlobal (n?u AutoSaveVars = true)
	Objects        map[string]map[string]any     `json:"objects"`   // key = SaveTag, value = OnSave()
}
