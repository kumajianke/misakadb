package tasktype

type TaskType string

// 作业类型及参数
type TaskBooks struct {
	TaskType TaskType
	Params   []string
}

// 构造器
type taskTypeBuilder struct {
	tasks []*TaskBooks
}

func NewTaskBooks(taskType TaskType, params ...string) *TaskBooks {
	return &TaskBooks{
		TaskType: taskType,
		Params:   params,
	}
}

// 创建新的构造器
func NewShipBuilder() *taskTypeBuilder {
	return &taskTypeBuilder{
		tasks: make([]*TaskBooks, 0),
	}
}

// 添加任务类型
func (b *taskTypeBuilder) Add(taskType TaskType, params ...string) *taskTypeBuilder {
	b.tasks = append(b.tasks, &TaskBooks{
		TaskType: taskType,
		Params:   params,
	})
	return b
}

// 构建任务列表
func (b *taskTypeBuilder) Build() []*TaskBooks {
	return b.tasks
}
