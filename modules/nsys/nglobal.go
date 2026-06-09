package nsys

import (
	"autoworld/domain"
)

// GetInstance trả về Singleton instance của IGlobal
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

func (g *globalStore) GetSprite(key string) domain.ISpriteLW {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.sprites[key]
}

func (g *globalStore) AddSprite(key string, sprite domain.ISpriteLW) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.sprites[key] = sprite
}

func (g *globalStore) UpdateSprite(key string, sprite domain.ISpriteLW) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.sprites[key] = sprite
}

func (g *globalStore) GetAudio(key string) domain.IAudioLW {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.sounds[key]
}

func (g *globalStore) AddAudio(key string, audio domain.IAudioLW) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.sounds[key] = audio
}

func (g *globalStore) UpdateAudio(key string, audio domain.IAudioLW) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.sounds[key] = audio
}

func (g *globalStore) GetObject(key string) domain.IObject {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.globalObjects[key]
}

func (g *globalStore) AddObject(key string, object domain.IObject) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.globalObjects[key] = object
}

func (g *globalStore) UpdateObject(key string, object domain.IObject) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.globalObjects[key] = object
}

func (g *globalStore) GetConst(key string) any {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.constant[key]
}

func (g *globalStore) NewConst(key string, value any) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.constant[key] == nil {
		g.constant[key] = value
		return true
	}
	return false
}

func (g *globalStore) UpdateConst(key string, value any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.constant[key] = value
}

func (g *globalStore) SetValue(key string, value any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.variable[key] = value
}

func (g *globalStore) GetValue(key string) any {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.variable[key]
}

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

func (g *globalStore) RestoreVariables(data map[string]any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for k, v := range data {
		g.variable[k] = v
	}
}
