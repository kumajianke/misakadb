package engine_base

import (
	mson "misakadb/engine/Mson"
	"sync"
)

type BaseLockerCore interface {
	Lock() *sync.Mutex
	GetRowLock(name string) *sync.Mutex // 获取行级锁
}

type EngineLockerSupport struct {
	RowLocks map[string]*sync.Mutex
	Locker   sync.Mutex
}

func (lockerCore *EngineLockerSupport) Lock() *sync.Mutex {
	return &lockerCore.Locker
}

func (lockerCore *EngineLockerSupport) GetRowLock(name string) *sync.Mutex {
	lockerCore.Locker.Lock()
	defer lockerCore.Locker.Unlock()

	if lockerCore.RowLocks == nil {
		lockerCore.RowLocks = make(map[string]*sync.Mutex)
	}

	rowLock := lockerCore.RowLocks[name]
	if rowLock == nil {
		lockerCore.RowLocks[name] = &sync.Mutex{}
	}

	return lockerCore.RowLocks[name]
}

/**
*用于数据库文件的操作IO等组件的使用
**/
type BaseLoaderCore interface {
	WriteLoader(log mson.MsonParse) error // 写入日志
	ReadLoader(log mson.MsonParse) error  // 读取日志
	InitLoader(log mson.MsonParse) error  // 初始化日志
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
	InsertDB(log mson.MsonParse) error // 插入数据
	DeleteDB(log mson.MsonParse) error // 删除数据
	UpdateDB(log mson.MsonParse) error // 更新数据
	SearchDB(log mson.MsonParse) error // 搜索数据库中的指定内容
}

/**
 * 数据库核心 不同的数据库指向了一个核心
 * TODO 所有核心的any只是暂时代替 后续会替换成对应的json结构
 */
type BaseEngineCore interface {
	BaseLockerCore
	DBLoader() BaseLoaderCore
	DBBaker() BaseBakerCore
	MiQLExecutor() MiQLExecutorCore
	RemoveDB() error
}
