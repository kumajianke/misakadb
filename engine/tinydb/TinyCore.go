package engine

import (
	engine_base "misakadb/engine/base"
	"misakadb/engine/tinydb/components"
)

type TinyDBCore struct {
	engine_base.EngineLockerSupport

	TinyDBLoader     engine_base.BaseLoaderCore
	TinyDBBaker      engine_base.BaseBakerCore
	TinyMiQLExecutor engine_base.MiQLExecutorCore

	Name string
}

var _ engine_base.BaseEngineCore = (*TinyDBCore)(nil)

func (this *TinyDBCore) RemoveDB() error {
	rowLock := this.GetRowLock(this.Name)
	rowLock.Lock()
	defer rowLock.Unlock()
	return nil
}

// 数据库引擎的备份组件
func (tinyDBCore *TinyDBCore) DBBaker() engine_base.BaseBakerCore {
	return tinyDBCore.TinyDBBaker
}

// 数据库的文件加载组件
func (tinyDBCore *TinyDBCore) DBLoader() engine_base.BaseLoaderCore {
	return tinyDBCore.TinyDBLoader
}

// 数据库引擎的Miql执行组件
func (tinyDBCore *TinyDBCore) MiQLExecutor() engine_base.MiQLExecutorCore {
	return tinyDBCore.TinyMiQLExecutor
}

// 获取一个TinyDB的实例
func NewTinyEngine(db_name string) *TinyDBCore {
	tinyDBCore := &TinyDBCore{
		Name: db_name,
	}

	tinyDBCore.TinyDBLoader = &components.TinyDBLoaderImp{
		DBName: db_name,
		Locker: tinyDBCore,
	}

	return tinyDBCore
}
