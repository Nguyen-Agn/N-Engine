package nsystem

type LogicSystem struct {
	CreateQuery  []IObject
	DestroyQuery []IObject
}

// NewLogicSystem khởi tạo LogicSystem.
func NewLogicSystem() *LogicSystem {
	return &LogicSystem{
		CreateQuery:  []IObject{},
		DestroyQuery: []IObject{},
	}
}

// Update duyệt qua tất cả thực thể và gọi StepUpdate.
func (this *LogicSystem) Update(objectList []IObject) {
	// Create
	for _, obj := range this.CreateQuery {
		obj.Create()
	}

	// Step
	for _, obj := range objectList {
		obj.StepUpdate()
	}

	// Destroy
	for _, obj := range this.DestroyQuery {
		obj.Destroy()
	}

	//Clear
	this.CreateQuery = this.CreateQuery[:0]
	this.DestroyQuery = this.DestroyQuery[:0]
}

func (this *LogicSystem) AddObjectCreated(object IObject) {
	this.CreateQuery = append(this.CreateQuery, object)
}

func (this *LogicSystem) AddObjectDestroy(object IObject) {
	this.DestroyQuery = append(this.DestroyQuery, object)
}
