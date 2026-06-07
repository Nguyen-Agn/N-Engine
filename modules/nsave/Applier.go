package nsave

import (
	"autoworld/domain"
	"log"
)

type applier struct{}

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

func (a *applier) applyVariables(store domain.IGlobal, vars map[string]any) {
	if store == nil || vars == nil {
		return
	}
	store.RestoreVariables(vars)
}
