package global_lock

type LocksPrefix struct {
	DBFile string // 对文件的alldb.list占用锁
	DBArea string // 对数据库引擎的整个文件进行上锁
}

func GetLocksPrefix() *LocksPrefix {
	return &LocksPrefix{
		DBFile: "db-files:",
		DBArea: "db-area:",
	}
}
