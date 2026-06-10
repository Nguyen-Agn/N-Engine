package nobject

import (
	"strconv"
	"strings"

	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"

	"github.com/yohamta/donburi"
)

type Object struct {
	entry *donburi.Entry
	pool  domain.IPool
}

// NewObject creates and returns a new Object instance using the given Donburi entry.
// Inputs: entry - pointer to the donburi ECS entry for this object.
// Outputs: a pointer to the newly created Object.
func NewObject(entry *donburi.Entry) *Object {
	return &Object{entry: entry}
}

// Entry returns the donburi.Entry inside the object.
// Used by napi.SetComponent and napi.GetComponent to assign custom component data.
// Outputs: pointer to the internal donburi.Entry.
func (this *Object) Entry() *donburi.Entry {
	return this.entry
}

// GetPool retrieves the object's parent pool.
// Outputs: the domain.IPool instance managing this object.
func (this *Object) GetPool() domain.IPool {
	return this.pool
}

// SetPool assigns the object's parent pool.
// Inputs: pool - the domain.IPool instance that will manage this object.
func (this *Object) SetPool(pool domain.IPool) {
	this.pool = pool
}

// #region Event

// OnCreate is triggered when the object is initially created or spawned.
func (this *Object) OnCreate() {}

// OnStep is triggered on every frame update to process the object's logic.
func (this *Object) OnStep() {}

// OnDestroy is triggered right before the object is completely removed from the world.
func (this *Object) OnDestroy() {}

// OnSave captures the object's specific save data for persistence.
// Inputs: data - a map where the object should store its persistable state.
func (this *Object) OnSave(data map[string]any) {}

// OnLoad restores the object's state from a loaded save data payload.
// Inputs: data - the previously saved map holding this object's specific state.
func (this *Object) OnLoad(data map[string]any) {}

// SetTokens configures the ECS component tokens bound to this object.
// Inputs: tokenClasses - a space-separated string of component types.
func (this *Object) SetTokens(tokenClasses string) {
	if this.entry == nil || tokenClasses == "" {
		return
	}
	tokens := strings.SplitSeq(tokenClasses, " ")
	for token := range tokens {
		parts := strings.SplitN(token, "-", 3)
		if len(parts) != 3 {
			continue
		}
		comp := parts[0]
		varName := parts[1]
		valStr := parts[2]

		switch comp {
		case "pos":
			if this.entry.HasComponent(enginetype.Position) {
				data := donburi.Get[domain.PositionData](this.entry, enginetype.Position)
				switch varName {
				case "x":
					if v, err := strconv.ParseFloat(valStr, 32); err == nil {
						data.X = float32(v)
					}
				case "y":
					if v, err := strconv.ParseFloat(valStr, 32); err == nil {
						data.Y = float32(v)
					}
				}
			}
		case "spr":
			if this.entry.HasComponent(enginetype.Sprite) {
				data := donburi.Get[domain.SpriteData](this.entry, enginetype.Sprite)
				switch varName {
				case "c":
					data.CurrentSprite = valStr
				case "idx":
					if v, err := strconv.Atoi(valStr); err == nil {
						data.SpriteIdx = v
					}
				case "z":
					if v, err := strconv.Atoi(valStr); err == nil {
						data.ZOrder = v
					}
				}
			}
		case "box":
			if this.entry.HasComponent(enginetype.Box) {
				data := donburi.Get[domain.BoxData](this.entry, enginetype.Box)
				switch varName {
				case "w":
					if v, err := strconv.ParseFloat(valStr, 32); err == nil {
						data.BoxW = float32(v)
					}
				case "h":
					if v, err := strconv.ParseFloat(valStr, 32); err == nil {
						data.BoxH = float32(v)
					}
				case "x":
					if v, err := strconv.ParseFloat(valStr, 32); err == nil {
						data.BoxX = float32(v)
					}
				case "y":
					if v, err := strconv.ParseFloat(valStr, 32); err == nil {
						data.BoxY = float32(v)
					}
				}
			}
		case "dir":
			if this.entry.HasComponent(enginetype.Direction) {
				data := donburi.Get[domain.DirectionData](this.entry, enginetype.Direction)
				if varName == "d" || varName == "dir" {
					if v, err := strconv.ParseFloat(valStr, 32); err == nil {
						data.Direction = float32(v)
					}
				}
			}
		case "velo":
			if this.entry.HasComponent(enginetype.Velocity) {
				data := donburi.Get[domain.VelocityData](this.entry, enginetype.Velocity)
				switch varName {
				case "vx":
					if v, err := strconv.ParseFloat(valStr, 32); err == nil {
						data.Vx = float32(v)
					}
				case "vy":
					if v, err := strconv.ParseFloat(valStr, 32); err == nil {
						data.Vy = float32(v)
					}
				case "max":
					if v, err := strconv.ParseFloat(valStr, 32); err == nil {
						data.MaxSpeed = float32(v)
					}
				}
			}
		default:
			// Fallback: try dynamic reflection-based setter for custom components (or exact field names)
			enginetype.SetComponentFieldByToken(this.entry, comp, varName, valStr)
		}
	}
}

// Remove marks this object for removal from the scene or pool.
func (this *Object) Remove() {}

// #endregion
