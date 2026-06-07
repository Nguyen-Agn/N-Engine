package domain

// ISaveManager l interface c?a h? th?ng Save/Load ton c?c.
type ISaveManager interface {
	// SaveGame luu ton b? tr?ng thi vo file theo du?ng d?n.
	SaveGame(path string) error

	// LoadGame t?i tr?ng thi t? file theo du?ng d?n v p d?ng vo game.
	LoadGame(path string) error

	// HasSave ki?m tra xem file save c t?n t?i khng.
	HasSave(path string) bool

	// DeleteSave xa save file t?i du?ng d?n cung c?p.
	DeleteSave(path string) error

	// ListSlots tr? v? danh sch cc t?n file save trong thu m?c luu m?c d?nh.
	ListSlots() []string

	// ReadSnapshot d?c d? li?u t? file nhung khng p d?ng.
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
