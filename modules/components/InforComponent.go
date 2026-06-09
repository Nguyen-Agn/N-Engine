package components

import "autoworld/modules/enginetype"

type InforComponent struct {
	IObject
}

var Infor = enginetype.Infor

// GetId retrieves the unique identifier of the entity.
// Purpose: Gets the entity's ID.
// Inputs: None.
// Outputs: The integer ID of the entity. Returns 0 if data is not initialized.
func (p InforComponent) GetId() int {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data == nil {
		return 0
	}
	return data.Id
}

// GetName retrieves the name of the entity.
// Purpose: Gets the entity's name.
// Inputs: None.
// Outputs: The string name of the entity. Returns empty string if data is not initialized.
func (p InforComponent) GetName() string {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data == nil {
		return ""
	}
	return data.Name
}

// AddTag adds a tag to the entity.
// Purpose: Associates a new string tag with the entity by hashing it.
// Inputs:
//   - tag: The string tag to add.
// Outputs: None.
func (p InforComponent) AddTag(tag string) {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data != nil {
		data.Tags = append(data.Tags, enginetype.HashString(tag))
	}
}

// HasTag checks if the entity has a specific tag.
// Purpose: Determines if the entity is associated with the given string tag.
// Inputs:
//   - tag: The string tag to check.
// Outputs: True if the tag is present, false otherwise.
func (p InforComponent) HasTag(tag string) bool {
	return p.HasTagHash(enginetype.HashString(tag))
}

// HasTagHash checks if the entity has a specific hashed tag.
// Purpose: Determines if the entity is associated with a pre-hashed tag.
// Inputs:
//   - hash: The uint64 hash of the tag to check.
// Outputs: True if the hashed tag is present, false otherwise.
func (p InforComponent) HasTagHash(hash uint64) bool {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data == nil {
		return false
	}
	for _, t := range data.Tags {
		if t == hash {
			return true
		}
	}
	return false
}

// IsDead checks if the entity is marked as dead.
// Purpose: Determines the lifecycle state of the entity.
// Inputs: None.
// Outputs: True if the entity is marked as dead, false otherwise.
func (p InforComponent) IsDead() bool {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data == nil {
		return false
	}
	return data.IsDead
}

// MarkDead marks the entity as dead.
// Purpose: Flags the entity for destruction or removal.
// Inputs: None.
// Outputs: None.
func (p InforComponent) MarkDead() {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data != nil {
		data.IsDead = true
	}
}

// SetIsDead manually sets the dead status of the entity.
// Purpose: Allows explicitly setting whether the entity is dead or alive.
// Inputs:
//   - dead: The boolean status to set.
// Outputs: None.
func (p InforComponent) SetIsDead(dead bool) {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data != nil {
		data.IsDead = dead
	}
}

// SaveTag retrieves the tag used for saving the entity state.
// Purpose: Gets the string identifier used for persistence.
// Inputs: None.
// Outputs: The string save tag. Returns empty string if not set.
func (p InforComponent) SaveTag() string {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data == nil {
		return ""
	}
	return data.SaveTag
}

// SetSaveTag sets the tag used for saving the entity state.
// Purpose: Assigns a string identifier for persistence purposes.
// Inputs:
//   - tag: The string save tag to set.
// Outputs: None.
func (p InforComponent) SetSaveTag(tag string) {
	data := enginetype.GetComponent(p.IObject, Infor)
	if data != nil {
		data.SaveTag = tag
	}
}
