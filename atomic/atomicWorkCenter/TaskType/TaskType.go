package tasktype

import "slices"

type TaskType string

// 作业类型及参数
type TaskBooks struct {
	TaskType TaskType
	Params   []string
}

// 构造器
type taskBooksShipBuilder struct {
	tasks []*TaskBooks
	after []*TaskBooks
}

func NewTaskBooks(taskType TaskType, params ...string) *TaskBooks {
	return &TaskBooks{
		TaskType: taskType,
		Params:   params,
	}
}

// 创建新的构造器
func NewTaskBooksShipBuilder() *taskBooksShipBuilder {
	return &taskBooksShipBuilder{
		tasks: make([]*TaskBooks, 0),
	}
}

/*
*
DESCRIPTION

	通过已经预设了的taskBooks列表创建taskTypeBuild
*/
func NewTaskBooksShipBuilderWithTaskBooks(taskBooks []*TaskBooks) *taskBooksShipBuilder {
	return &taskBooksShipBuilder{
		tasks: taskBooks,
	}
}

// 添加任务类型
func (b *taskBooksShipBuilder) Add(taskType TaskType, params ...string) *taskBooksShipBuilder {
	b.tasks = append(b.tasks, &TaskBooks{
		TaskType: taskType,
		Params:   params,
	})
	return b
}

func (b *taskBooksShipBuilder) AddAfter(taskType TaskType, params ...string) *taskBooksShipBuilder {
	b.after = append(b.after, &TaskBooks{
		TaskType: taskType,
		Params:   params,
	})
	return b
}
func (b *taskBooksShipBuilder) AfterAndBuild(taskType TaskType, params ...string) []*TaskBooks {
	b.tasks = append(b.tasks, &TaskBooks{
		TaskType: taskType,
		Params:   params,
	})
	return b.Build()
}

// 构建任务列表
func (b *taskBooksShipBuilder) Build() []*TaskBooks {
	tasks := slices.Concat(b.tasks, b.after)
	taskBooksShipBuilder := NewTaskBooksShipBuilderWithTaskBooks(tasks)

	return taskBooksShipBuilder.tasks
}
