package engine_base

/**
* 数据库日志核心，用于记录每个记录指定存储等信息
**/
type BaseLoaderCore interface {
	WriteLoader(log any) error // 写入日志
	ReadLoader(log any) error  // 读取日志
}

/**
 * 数据库备份核心，用于备份数据库
 */
type BaseBakerCore interface {
	TriggerBaker() error  // 触发备份的条件
	BackupDBLogic() error // 备份数据库的逻辑
}

/**
 * 数据库执行核心，用于执行数据库DML语句
 */
type MiQLExecutorCore interface {
	InsertDB(log any) error // 插入数据
	DeleteDB(log any) error // 删除数据
	UpdateDB(log any) error // 更新数据
	SearchDB(log any) error // 搜索数据库中的指定内容
}

/**
 * 数据库核心 不同的数据库指向了一个核心
 * TODO 所有核心的any只是暂时代替 后续会替换成对应的json结构
 */
type BaseEngineCore interface {
	Path() string
	DBLoader() BaseLoaderCore
	DBBaker() BaseBakerCore
	MiQLExecutor() MiQLExecutorCore
}
