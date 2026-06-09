package nsys

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
)

// GetInstance returns the Singleton instance of the global store, initializing it if necessary.
// Outputs: domain.IGlobal instance.
func GetInstance() domain.IGlobal {
	once.Do(func() {
		store = &globalStore{
			sprites:       make(map[string]domain.ISpriteLW),
			sounds:        make(map[string]domain.IAudioLW),
			globalObjects: make(map[string]domain.IObject),
			variable:      make(map[string]any),
			constant:      make(map[string]any),
		}
	})
	return store
}

// GetSprite thread-safely retrieves a sprite by its key.
func (g *globalStore) GetSprite(key string) domain.ISpriteLW {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.sprites[key]
}

// AddSprite thread-safely adds a new sprite associated with a key.
func (g *globalStore) AddSprite(key string, sprite domain.ISpriteLW) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.sprites[key] = sprite
}

// UpdateSprite thread-safely updates an existing sprite under a key.
func (g *globalStore) UpdateSprite(key string, sprite domain.ISpriteLW) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.sprites[key] = sprite
}

// GetAudio thread-safely retrieves an audio instance by its key.
func (g *globalStore) GetAudio(key string) domain.IAudioLW {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.sounds[key]
}

// AddAudio thread-safely adds a new audio instance associated with a key.
func (g *globalStore) AddAudio(key string, audio domain.IAudioLW) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.sounds[key] = audio
}

// UpdateAudio thread-safely updates an existing audio instance under a key.
func (g *globalStore) UpdateAudio(key string, audio domain.IAudioLW) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.sounds[key] = audio
}

// GetObject thread-safely retrieves a globally stored object by its key.
func (g *globalStore) GetObject(key string) domain.IObject {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.globalObjects[key]
}

// AddObject thread-safely adds a new global object under a key.
func (g *globalStore) AddObject(key string, object domain.IObject) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.globalObjects[key] = object
}

// UpdateObject thread-safely updates an existing global object under a key.
func (g *globalStore) UpdateObject(key string, object domain.IObject) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.globalObjects[key] = object
}

// GetConst thread-safely retrieves a constant value by its key.
func (g *globalStore) GetConst(key string) any {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.constant[key]
}

// NewConst thread-safely creates a new constant if the key does not already exist.
// Returns true if successfully added, false if the key is already occupied.
func (g *globalStore) NewConst(key string, value any) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.constant[key] == nil {
		g.constant[key] = value
		return true
	}
	return false
}

// UpdateConst thread-safely overwrites the value of an existing constant.
func (g *globalStore) UpdateConst(key string, value any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.constant[key] = value
}

// SetValue thread-safely sets or overwrites a generic variable value under a key.
func (g *globalStore) SetValue(key string, value any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.variable[key] = value
}

// GetValue thread-safely retrieves a generic variable value by its key.
func (g *globalStore) GetValue(key string) any {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.variable[key]
}

// DumpVariables thread-safely extracts and returns all primitive variables for saving.
// Outputs: a map containing only primitive types safe for serialization.
func (g *globalStore) DumpVariables() map[string]any {
	g.mu.RLock()
	defer g.mu.RUnlock()
	res := make(map[string]any)
	for k, v := range g.variable {
		// Only dump primitives to ensure save compatibility
		switch v.(type) {
		case int, int64, float32, float64, string, bool:
			res[k] = v
		}
	}
	return res
}

// RestoreVariables thread-safely loads external variable data back into the store.
func (g *globalStore) RestoreVariables(data map[string]any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for k, v := range data {
		g.variable[k] = v
	}
}
