package nsave

import (
	"autoworld/domain"
	"log"
)

type collector struct{}

// collectObjects collects save data from all eligible objects in a scene and groups them by their unique tags.
// Inputs: scene - the active game scene holding the maps.
// Outputs: nested map linking unique tags to their respective save state maps.
func (c *collector) collectObjects(scene domain.IScene) map[string]map[string]any {
	result := make(map[string]map[string]any)
	if scene == nil {
		return result
	}

	physicalMap := scene.GetMap()
	if physicalMap != nil {
		c.collectFromObjectList(physicalMap.GetObjects(), result, "")
	}

	guiMap := scene.GetGUIMap()
	if guiMap != nil {
		// prefix "gui:" to avoid collision with physical map objects
		c.collectFromObjectList(guiMap.GetObjects(), result, "gui:")
	}

	return result
}

// collectFromObjectList extracts data from a specific slice of objects, handling dead elements and adding prefixes.
// Inputs: objects - array of objects to inspect, result - map accumulating all states, prefix - tag prefix.
func (c *collector) collectFromObjectList(objects []domain.IObject, result map[string]map[string]any, prefix string) {
	for _, obj := range objects {
		// Skip objects that have been marked for removal
		if dead, ok := obj.(interface{ IsDead() bool }); ok && dead.IsDead() {
			continue
		}

		infor, ok := obj.(interface{ SaveTag() string })
		if !ok {
			continue
		}

		tag := infor.SaveTag()
		if tag == "" {
			continue // skip objects without a save tag
		}

		saveData := make(map[string]any)
		obj.OnSave(saveData)
		if len(saveData) == 0 {
			continue // object chose not to save anything
		}

		fullTag := prefix + tag
		if _, exists := result[fullTag]; exists {
			log.Printf("[nsave] WARNING: duplicate save tag found: %s. Overwriting data.\n", fullTag)
		}
		result[fullTag] = saveData
	}
}

// collectVariables extracts primitive variables from the global store for saving.
// Inputs: store - the global variable registry.
// Outputs: map containing primitive global variables.
func (c *collector) collectVariables(store domain.IGlobal) map[string]any {
	if store == nil {
		return nil
	}
	return store.DumpVariables()
}
