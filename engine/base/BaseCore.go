package engine_base

import (
	"encoding/json"
	"errors"
	mson "misakadb/engine/Mson"
	"os"
	"path/filepath"
	"sync"
)

type BaseLockerCore interface {
	Lock() *sync.Mutex
	GetRowLock(name string) *sync.Mutex // 获取行级锁
}

type EngineLockerSupport struct {
	rowLocks sync.Map
	Locker   sync.Mutex
}

func (lockerCore *EngineLockerSupport) Lock() *sync.Mutex {
	return &lockerCore.Locker
}

func (lockerCore *EngineLockerSupport) GetRowLock(name string) *sync.Mutex {
	rowLock, _ := lockerCore.rowLocks.LoadOrStore(name, &sync.Mutex{})
	return rowLock.(*sync.Mutex)
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
 */
type BaseEngineCore interface {
	BaseLockerCore
	DBLoader() BaseLoaderCore
	DBBaker() BaseBakerCore
	MiQLExecutor() MiQLExecutorCore
	RemoveDB(dbname string) error
}

type BaseDBMeta struct {
	DBName     string   `json:"db_name"`
	AllTables  []string `json:"all_tables"`
	CreateTime string   `json:"create_time"`
	Engine     string   `json:"engine"`
}

func ShareLoaderDBMetaName(dbname string) (string, error) {

	path := filepath.Join(".", "db-datas", dbname, "meta.json")
	dbMeta := &BaseDBMeta{}
	content, err := os.ReadFile(path)
	if err != nil {
		return "", errors.New("can not read the db meta file")
	}

	err = json.Unmarshal([]byte(content), dbMeta)
	if err != nil {
		return "", errors.New("can not convert the db-meta to json")
	}

	return dbMeta.Engine, nil

}
