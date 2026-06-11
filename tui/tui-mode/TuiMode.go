package tuimode

type TuiMode int

// 全局终端类型 默认是未知 初始化之后默认CLIMode
var GlobalTuiMode TuiMode = UnKnown

const (
	CLIMode TuiMode = iota
	LockMode
	UnKnown
)

func GetTuiMode() TuiMode {
	if GlobalTuiMode == UnKnown {
		GlobalTuiMode = CLIMode
	}
	return GlobalTuiMode
}

func SetTuiMode(mode TuiMode) {
	GlobalTuiMode = mode
}

func (w TuiMode) String() string {
	return [...]string{"CLIMode", "LockMode", "UnKnown"}[w]
}
