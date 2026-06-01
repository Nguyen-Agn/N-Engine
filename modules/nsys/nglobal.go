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

// get value type Int
// return defuatl = 0
func (g *globalStore) GetInt(key string) int {
	g.mu.Lock()
	defer g.mu.Unlock()
	if val, ok := g.variable[key].(int); ok {
		return val
	}
	return 0

}

// get value type Int64
// return defuatl = 0
func (g *globalStore) GetInt64(key string) int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	if val, ok := g.variable[key].(int64); ok {
		return val
	}
	return 0
}

// get value type String
// return defuatl = ""
func (g *globalStore) GetString(key string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	{
		if val, ok := g.variable[key].(string); ok {
			return val
		}
		return ""
	}
}

// get value type float32
// return defuatl = 0.0
func (g *globalStore) GetFloat32(key string) float32 {
	g.mu.Lock()
	defer g.mu.Unlock()
	if val, ok := g.variable[key].(float32); ok {
		return val
	}
	return 0.0
}

// get value type float64
// return defuatl = 0.0
func (g *globalStore) GetFloat64(key string) float64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	if val, ok := g.variable[key].(float64); ok {
		return val
	}
	return 0.0
}

// get value type bool
// return defuatl = false
func (g *globalStore) GetBool(key string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	{
		if val, ok := g.variable[key].(bool); ok {
			return val
		}
		return false
	}
}
