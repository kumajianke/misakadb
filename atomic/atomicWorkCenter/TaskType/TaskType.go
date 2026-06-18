package tasktype

type TaskType string

const (
	TaskRemoveFile TaskType = "remove_file"
	TaskModFile    TaskType = "mod_file"
)

// 作业类型及参数
type TaskTypeShip struct {
	TaskType TaskType
	Params   []string
}

func NewTaskTypeShip(taskType TaskType, params ...string) *TaskTypeShip {
	return &TaskTypeShip{
		TaskType: taskType,
		Params:   params,
	}
}
