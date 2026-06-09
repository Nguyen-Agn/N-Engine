package nsave

import (
	"autoworld/domain"
	"log"
)

type applier struct{}

// applyObjects applies saved object data to the current scene by injecting state into matching objects.
// Inputs: scene - the active game scene, savedObjects - nested map containing the save payload for objects.
func (a *applier) applyObjects(scene domain.IScene, savedObjects map[string]map[string]any) {
	if scene == nil || savedObjects == nil {
		return
	}

	physicalMap := scene.GetMap()
	if physicalMap != nil {
		a.applyToObjectList(physicalMap.GetObjects(), savedObjects, "")
	}

	guiMap := scene.GetGUIMap()
	if guiMap != nil {
		a.applyToObjectList(guiMap.GetObjects(), savedObjects, "gui:")
	}

	// For keys in savedObjects that were not found in the scene, log a warning
	for key := range savedObjects {
		log.Printf("[nsave] INFO: Unmatched save tag '%s' in save data. Object might not be spawned yet.\n", key)
	}
}

// applyToObjectList iterates through a slice of objects and applies save data if their save tags match.
// Inputs: objects - a slice of game objects, savedObjects - the remaining save map, prefix - prefix to append to save tags (e.g., "gui:").
func (a *applier) applyToObjectList(objects []domain.IObject, savedObjects map[string]map[string]any, prefix string) {
	for _, obj := range objects {
		// Skip dead objects — they are queued for removal and must not be resurrected
		if dead, ok := obj.(interface{ IsDead() bool }); ok && dead.IsDead() {
			continue
		}

		infor, ok := obj.(interface{ SaveTag() string })
		if !ok {
			continue
		}

		tag := infor.SaveTag()
		if tag == "" {
			continue
		}

		fullTag := prefix + tag
		if saveData, exists := savedObjects[fullTag]; exists {
			obj.OnLoad(saveData)
			delete(savedObjects, fullTag) // mark as matched
		}
	}
}

// applyVariables restores global variable state into the main store from the saved data.
// Inputs: store - the global variable registry, vars - primitive value map containing saved variables.
func (a *applier) applyVariables(store domain.IGlobal, vars map[string]any) {
	if store == nil || vars == nil {
		return
	}
	store.RestoreVariables(vars)
}
