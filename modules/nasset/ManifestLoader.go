package nasset

import (
	"fmt"
	"os"

	"github.com/Nguyen-Agn/N-Engine/domain"

	"github.com/BurntSushi/toml"
)

// ─── Manifest structs ────────────────────────────────────────────────────────

// spriteEntry mô tả một tài nguyên sprite trong file manifest TOML.
type spriteEntry struct {
	Key  string `toml:"key"`
	Path string `toml:"path"`
	// Mode: "single" hoặc "sheet"
	Mode string `toml:"mode"`
	// Các trường bên dưới chỉ dùng khi Mode = "sheet"
	FrameW int `toml:"frameW"`
	FrameH int `toml:"frameH"`
	Cols   int `toml:"cols"`
	Rows   int `toml:"rows"`
	GapX   int `toml:"gapX"`
	GapY   int `toml:"gapY"`
}

// audioEntry mô tả một tài nguyên âm thanh trong file manifest TOML.
type audioEntry struct {
	Key  string `toml:"key"`
	Path string `toml:"path"`
}

// varEntry mô tả một biến (variable) hoặc hằng số (constant) lưu trong nGlobal.
type varEntry struct {
	Key   string `toml:"key"`
	Value any    `toml:"value"`
}

// assetManifest là cấu trúc gốc tương ứng toàn bộ file manifest TOML.
type assetManifest struct {
	Sprites   []spriteEntry `toml:"sprites"`
	Audios    []audioEntry  `toml:"audios"`
	Vars      []varEntry    `toml:"vars"`
	Constants []varEntry    `toml:"constants"`
}

// ─── ManifestLoader ──────────────────────────────────────────────────────────

// ManifestLoader triển khai domain.IManifestLoader.
// Đọc file manifest TOML, load từng tài nguyên và lưu vào IGlobal store.
type ManifestLoader struct {
	// spriteLoader phụ trách việc decode ảnh từ đĩa.
	spriteLoader domain.ISpriteLoader
	// audioLoader phụ trách việc decode âm thanh từ đĩa (có thể nil nếu chưa cần).
	audioLoader domain.IAudioLoader
}

// Purpose: Creates a new ManifestLoader.
// Inputs: spriteLoader (domain.ISpriteLoader) - Loader for sprites, audioLoader (domain.IAudioLoader) - Loader for audios (can be nil).
// Outputs: (*ManifestLoader) - A pointer to the new ManifestLoader.
func NewManifestLoader(spriteLoader domain.ISpriteLoader, audioLoader domain.IAudioLoader) *ManifestLoader {
	return &ManifestLoader{
		spriteLoader: spriteLoader,
		audioLoader:  audioLoader,
	}
}

// Purpose: Reads a TOML manifest file, parses it, and loads all resources into the global store.
// Inputs: filePath (string) - The path to the TOML manifest file, store (domain.IGlobal) - The global store to hold the assets.
// Outputs: (error) - Error if reading or parsing fails, or if a critical asset fails to load.
func (m *ManifestLoader) LoadFromFile(filePath string, store domain.IGlobal) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("manifest: không thể đọc file '%s': %w", filePath, err)
	}

	var manifest assetManifest
	if err := toml.Unmarshal(data, &manifest); err != nil {
		return fmt.Errorf("manifest: lỗi parse TOML '%s': %w", filePath, err)
	}

	// Load sprites — lỗi ở đây là fatal (thiếu ảnh sẽ crash game)
	if err := m.loadSprites(manifest.Sprites, store); err != nil {
		return err
	}

	// Load audios — lỗi ở đây được trả về nhưng không chặn phần còn lại
	audioErr := m.loadAudios(manifest.Audios, store)

	// Vars và constants luôn load được (chỉ là giá trị JSON thuần)
	m.loadVars(manifest.Vars, store)
	m.loadConstants(manifest.Constants, store)

	// Trả về lỗi audio nếu có (caller quyết định có fatal hay chỉ warn)
	return audioErr
}

// Purpose: Iterates through sprite entries from the manifest and loads them into the store.
// Inputs: entries ([]spriteEntry) - List of sprite specifications, store (domain.IGlobal) - Global store.
// Outputs: (error) - Error if any sprite fails to load.
func (m *ManifestLoader) loadSprites(entries []spriteEntry, store domain.IGlobal) error {
	if m.spriteLoader == nil {
		return nil
	}
	for _, e := range entries {
		var (
			sprite domain.ISpriteLW
			err    error
		)
		switch e.Mode {
		case "sheet":
			sprite, err = m.spriteLoader.LoadSheet(e.Path, e.FrameW, e.FrameH, e.Cols, e.Rows, e.GapX, e.GapY)
		default:
			// Mặc định: "single"
			sprite, err = m.spriteLoader.LoadSingle(e.Path)
		}

		if err != nil {
			return fmt.Errorf("manifest: không thể load sprite key='%s' path='%s': %w", e.Key, e.Path, err)
		}
		store.AddSprite(e.Key, sprite)
		fmt.Println("Loaded Sprite ", e.Key)
	}
	return nil
}

// Purpose: Iterates through audio entries from the manifest and loads them into the store.
// Inputs: entries ([]audioEntry) - List of audio specifications, store (domain.IGlobal) - Global store.
// Outputs: (error) - Error if any audio fails to load.
func (m *ManifestLoader) loadAudios(entries []audioEntry, store domain.IGlobal) error {
	if m.audioLoader == nil {
		// audioLoader chưa được cung cấp — bỏ qua âm thanh, không lỗi.
		return nil
	}
	for _, e := range entries {
		audio, err := m.audioLoader.Load(e.Path)
		if err != nil {
			return fmt.Errorf("manifest: không thể load audio key='%s' path='%s': %w", e.Key, e.Path, err)
		}
		store.AddAudio(e.Key, audio)
		fmt.Println("Loaded Audio ", e.Key)
	}

	return nil
}

// Purpose: Loads global variables into the global store.
// Inputs: entries ([]varEntry) - List of variable specifications, store (domain.IGlobal) - Global store.
func (m *ManifestLoader) loadVars(entries []varEntry, store domain.IGlobal) {
	for _, e := range entries {
		store.SetValue(e.Key, e.Value)
	}
}

// Purpose: Loads constants into the global store.
// Inputs: entries ([]varEntry) - List of constant specifications, store (domain.IGlobal) - Global store.
func (m *ManifestLoader) loadConstants(entries []varEntry, store domain.IGlobal) {
	for _, e := range entries {
		store.NewConst(e.Key, e.Value)
	}
}
