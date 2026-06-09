package globalconfig

import "maps"

// GlobalConfig implements domain.IGlobalConfig with the Observer Pattern.
type GlobalConfig struct {
	values    map[string]any
	consts    map[string]any
	observers []IObserver
}

// SingleTon
var instance *GlobalConfig

// NewGlobalConfig initializes and returns the singleton instance of GlobalConfig.
// Outputs: pointer to the global GlobalConfig instance.
func NewGlobalConfig() *GlobalConfig {
	if instance == nil {
		instance = &GlobalConfig{
			values: make(map[string]any),
			consts: make(map[string]any),
		}
	}
	return instance
}

// GetConfig returns the current singleton instance of GlobalConfig.
func (c *GlobalConfig) GetConfig() IGlobalConfig {
	return instance
}

// SetValue assigns a value to a specified configuration key.
// Inputs: key - string identifier, value - the data to store.
func (c *GlobalConfig) SetValue(key string, value any) {
	c.values[key] = value
}

// NewConst registers a new constant value if the key doesn't already exist.
// Returns true if successfully added, false otherwise.
func (c *GlobalConfig) NewConst(key string, value any) bool {
	if c.consts[key] == nil {
		c.consts[key] = value
		return true
	}
	return false
}

// AddObserver registers an observer to receive notifications on configuration changes.
func (c *GlobalConfig) AddObserver(o IObserver) {
	c.observers = append(c.observers, o)
}

// RemoveObserver unregisters an observer so it no longer receives notifications.
func (c *GlobalConfig) RemoveObserver(o IObserver) {
	for i, observer := range c.observers {
		if observer == o {
			c.observers = append(c.observers[:i], c.observers[i+1:]...)
			break
		}
	}
}

// NotifyChange broadcasts a change notification to all registered observers.
func (c *GlobalConfig) NotifyChange() {
	for _, o := range c.observers {
		o.NotifyChange(c)
	}
}

// GetValue retrieves an untyped configuration value by its key.
func (c *GlobalConfig) GetValue(key string) any {
	return c.values[key]
}

// GetConst retrieves an untyped constant configuration value by its key.
func (c *GlobalConfig) GetConst(key string) any {
	return c.consts[key]
}

// DumpVariables creates and returns a complete copy of all current configuration variables.
func (c *GlobalConfig) DumpVariables() map[string]any {
	dump := make(map[string]any)
	maps.Copy(dump, c.values)
	return dump
}

// RestoreVariables restores configuration state from external data and triggers a notification.
func (c *GlobalConfig) RestoreVariables(data map[string]any) {
	if data == nil {
		return
	}
	maps.Copy(c.values, data)
	c.NotifyChange()
}
