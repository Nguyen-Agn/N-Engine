package globalconfig

// GlobalConfig triển khai domain.IGlobalConfig với Observer Pattern.
type GlobalConfig struct {
	values    map[string]any
	consts    map[string]any
	observers []IObserver
}

// SingleTon
var instance *GlobalConfig

func NewGlobalConfig() *GlobalConfig {
	if instance == nil {
		instance = &GlobalConfig{
			values: make(map[string]any),
			consts: make(map[string]any),
		}
	}
	return instance
}

// get global config
func (c *GlobalConfig) GetConfig() IGlobalConfig {
	return instance
}

// set value for key config
func (c *GlobalConfig) SetValue(key string, value any) {
	c.values[key] = value
}

// add new constanst
func (c *GlobalConfig) NewConst(key string, value any) bool {
	if c.consts[key] == nil {
		c.consts[key] = value
		return true
	}
	return false
}

// add observer
func (c *GlobalConfig) AddObserver(o IObserver) {
	c.observers = append(c.observers, o)
}

// remove observer
func (c *GlobalConfig) RemoveObserver(o IObserver) {
	for i, observer := range c.observers {
		if observer == o {
			c.observers = append(c.observers[:i], c.observers[i+1:]...)
			break
		}
	}
}

// notify to all observer
func (c *GlobalConfig) NotifyChange() {
	for _, o := range c.observers {
		o.NotifyChange(c)
	}
}

// get value type Int
// return defuatl = 0
func (c *GlobalConfig) GetInt(key string) int {
	if val, ok := c.values[key].(int); ok {
		return val
	}
	return 0
}

// get value type Int64
// return defuatl = 0
func (c *GlobalConfig) GetInt64(key string) int64 {
	if val, ok := c.values[key].(int64); ok {
		return val
	}
	return 0
}

// get value type String
// return defuatl = ""
func (c *GlobalConfig) GetString(key string) string {
	if val, ok := c.values[key].(string); ok {
		return val
	}
	return ""
}

// get value type float32
// return defuatl = 0.0
func (c *GlobalConfig) GetFloat32(key string) float32 {
	if val, ok := c.values[key].(float32); ok {
		return val
	}
	return 0.0
}

// get value type float64
// return defuatl = 0.0
func (c *GlobalConfig) GetFloat64(key string) float64 {
	if val, ok := c.values[key].(float64); ok {
		return val
	}
	return 0.0
}

// get value type bool
// return defuatl = false
func (c *GlobalConfig) GetBool(key string) bool {
	if val, ok := c.values[key].(bool); ok {
		return val
	}
	return false

}

func (c *GlobalConfig) GetConst(key string) any {
	return c.consts[key]
}
