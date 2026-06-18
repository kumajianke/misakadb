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

// 构造器
type taskTypeBuilder struct {
	tasks []*TaskTypeShip
}

func NewTaskTypeShip(taskType TaskType, params ...string) *TaskTypeShip {
	return &TaskTypeShip{
		TaskType: taskType,
		Params:   params,
	}
}

// 创建新的构造器
func NewShipBuilder() *taskTypeBuilder {
	return &taskTypeBuilder{
		tasks: make([]*TaskTypeShip, 0),
	}
}

// 添加任务类型
func (b *taskTypeBuilder) Add(taskType TaskType, params ...string) *taskTypeBuilder {
	b.tasks = append(b.tasks, &TaskTypeShip{
		TaskType: taskType,
		Params:   params,
	})
	return b
}

// 构建任务列表
func (b *taskTypeBuilder) Build() []*TaskTypeShip {
	return b.tasks
}
