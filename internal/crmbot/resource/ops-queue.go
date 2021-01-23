package resource

type OperationType string

const (
	OperationType_CategoryAdd OperationType = "CategoryAdd"
)

type Operation struct {
	Name OperationType
	Step int
	Data interface{}
}

// OpsQueue - map to store all queues for operations
// map[userID]Operation
var OpsQueue = make(map[int]*Operation)