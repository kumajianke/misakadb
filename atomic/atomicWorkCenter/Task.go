package atomic_work_center

type EnumTaskStatus int

const (
	Pending  EnumTaskStatus = iota // 待处理
	Running                        // 运行中
	Success                        // 完成
	BackRoll                       // 回滚
)

type Task struct {
	TaskStatus      EnumTaskStatus
	TaskReleaseTime int64
	TaskBody        []*TaskTypeShip
}
