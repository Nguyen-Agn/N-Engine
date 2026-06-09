package nsystem

type LogicSystem struct {
	CreateQuery  []IObject
	DestroyQuery []IObject
}

// NewLogicSystem initializes and returns a new instance of LogicSystem.
// Outputs: Returns a pointer to a newly initialized LogicSystem.
func NewLogicSystem() *LogicSystem {
	return &LogicSystem{
		CreateQuery:  []IObject{},
		DestroyQuery: []IObject{},
	}
}

// Update processes lifecycle events for all objects.
// Inputs: objectList ([]IObject) - The active list of objects in the scene.
// Purpose: It calls OnCreate for newly added objects, OnStep for all active objects, and OnDestroy for removed objects. Finally, it clears the create and destroy queues.
func (this *LogicSystem) Update(objectList []IObject) {
	// Create
	for _, obj := range this.CreateQuery {
		obj.OnCreate()
	}

	// Step
	for _, obj := range objectList {
		obj.OnStep()
	}

	// Destroy
	for _, obj := range this.DestroyQuery {
		obj.OnDestroy()
	}

	//Clear
	this.CreateQuery = this.CreateQuery[:0]
	this.DestroyQuery = this.DestroyQuery[:0]
}

// AddObjectCreated queues an object to have its OnCreate method called in the next update.
// Inputs: object (IObject) - The object that was just created.
func (this *LogicSystem) AddObjectCreated(object IObject) {
	this.CreateQuery = append(this.CreateQuery, object)
}

// AddObjectDestroy queues an object to have its OnDestroy method called in the next update.
// Inputs: object (IObject) - The object that is pending destruction.
func (this *LogicSystem) AddObjectDestroy(object IObject) {
	this.DestroyQuery = append(this.DestroyQuery, object)
}
