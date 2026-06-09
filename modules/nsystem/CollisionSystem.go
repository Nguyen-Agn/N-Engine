package nsystem

import (
	"github.com/Nguyen-Agn/N-Engine/domain"
	"github.com/Nguyen-Agn/N-Engine/modules/enginetype"

	"math"
)

var DefaultCollisionGridSize float32 = 64.0

type CollisionSystem struct {
	GridSize    float32
	collidables []domain.IObject
}

// NewCollisionSystem creates and returns a new instance of CollisionSystem with a default grid size.
// Outputs: Returns a pointer to a newly initialized CollisionSystem.
func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{
		GridSize:    DefaultCollisionGridSize,
		collidables: make([]domain.IObject, 0),
	}
}

type collisionEvent struct {
	objA    domain.IObject
	objB    domain.IObject
	handler func(other domain.IObject)
}

// AddObject registers an object to the collision system if it contains both Collision and Box components.
// Inputs: obj (domain.IObject) - The object to be added to the collidables list.
func (s *CollisionSystem) AddObject(obj domain.IObject) {
	if enginetype.GetComponent(obj, enginetype.Collision) != nil && enginetype.GetComponent(obj, enginetype.Box) != nil {
		s.collidables = append(s.collidables, obj)
	}
}

// Update performs spatial hashing based collision detection on all registered collidable objects.
// Inputs: objectList ([]domain.IObject) - The list of all objects (currently unused, it uses internal collidables list).
// Purpose: It first cleans up dead objects. Then it maps active collidables into a grid based on their bounding boxes. Next, it checks for AABB intersections between objects in the same grid cells. If a collision is detected and a handler exists for the target's tag, it queues and later executes the handler.
func (s *CollisionSystem) Update(objectList []domain.IObject) {

	alive := s.collidables[:0]
	for _, obj := range s.collidables {
		if info, ok := obj.(domain.IInfor); ok && info.IsDead() {
			continue
		}
		alive = append(alive, obj)
	}
	s.collidables = alive

	if len(s.collidables) < 2 {
		return
	}

	grid := make(map[int]map[int][]domain.IObject)

	getCell := func(x, y float32) (int, int) {
		return int(math.Floor(float64(x / s.GridSize))), int(math.Floor(float64(y / s.GridSize)))
	}

	for _, obj := range s.collidables {
		box := enginetype.GetComponent(obj, enginetype.Box)
		if !box.IsCollidable {
			continue
		}
		pos := enginetype.GetComponent(obj, enginetype.Position)

		startX, startY := getCell(pos.X+box.BoxX, pos.Y+box.BoxY)
		endX, endY := getCell(pos.X+box.BoxX+box.BoxW, pos.Y+box.BoxY+box.BoxH)

		for cx := startX; cx <= endX; cx++ {
			for cy := startY; cy <= endY; cy++ {
				if grid[cx] == nil {
					grid[cx] = make(map[int][]domain.IObject)
				}
				grid[cx][cy] = append(grid[cx][cy], obj)
			}
		}
	}

	var events []collisionEvent
	checked := make(map[int]map[int]bool)

	for _, col := range grid {
		for _, cellObjs := range col {
			n := len(cellObjs)
			for i := range n {
				objA := cellObjs[i]
				colDataA := enginetype.GetComponent(objA, enginetype.Collision)
				if colDataA == nil || len(colDataA.Handlers) == 0 {
					continue
				}

				infoA, okA := objA.(domain.IInfor)
				if !okA {
					continue
				}
				idA := infoA.GetId()

				for j := range n {
					if i == j {
						continue
					}
					objB := cellObjs[j]

					infoB, okB := objB.(domain.IInfor)
					if !okB {
						continue
					}
					idB := infoB.GetId()

					if checked[idA] != nil && checked[idA][idB] {
						continue
					}
					if checked[idA] == nil {
						checked[idA] = make(map[int]bool)
					}
					checked[idA][idB] = true

					var handler func(other domain.IObject)
					infoDataB := enginetype.GetComponent(objB, enginetype.Infor)
					if infoDataB != nil {
						for _, tag := range infoDataB.Tags {
							if h, ok := colDataA.Handlers[tag]; ok {
								handler = h
								break
							}
						}
					}

					if handler != nil {
						// AABB check
						boxA := enginetype.GetComponent(objA, enginetype.Box)
						posA := enginetype.GetComponent(objA, enginetype.Position)
						boxB := enginetype.GetComponent(objB, enginetype.Box)
						posB := enginetype.GetComponent(objB, enginetype.Position)

						if posA != nil && posB != nil {
							ax := posA.X + boxA.BoxX
							ay := posA.Y + boxA.BoxY
							bx := posB.X + boxB.BoxX
							by := posB.Y + boxB.BoxY

							if ax < bx+boxB.BoxW && ax+boxA.BoxW > bx &&
								ay < by+boxB.BoxH && ay+boxA.BoxH > by {

								events = append(events, collisionEvent{
									objA:    objA,
									objB:    objB,
									handler: handler,
								})
							}
						}
					}
				}
			}
		}
	}

	for _, ev := range events {
		infoA, okA := ev.objA.(domain.IInfor)
		infoB, okB := ev.objB.(domain.IInfor)

		if (okA && infoA.IsDead()) || (okB && infoB.IsDead()) {
			continue
		}

		ev.handler(ev.objB)
	}
}
