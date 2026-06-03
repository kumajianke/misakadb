package engine

import (
	engine_base "misakadb/engine/base"
	"misakadb/engine/tinydb/components"
	"sync"
)

type TinyDBCore struct {
	TinyDBLoader     engine_base.BaseLoaderCore
	TinyDBBaker      engine_base.BaseBakerCore
	TinyMiQLExecutor engine_base.MiQLExecutorCore

	Name string
	Lock sync.Mutex
}

var _ engine_base.BaseEngineCore = (*TinyDBCore)(nil)

func (tinyDBCore *TinyDBCore) DBBaker() engine_base.BaseBakerCore {
	return tinyDBCore.TinyDBBaker
}

func (tinyDBCore *TinyDBCore) DBLoader() engine_base.BaseLoaderCore {
	return tinyDBCore.TinyDBLoader
}

func (tinyDBCore *TinyDBCore) MiQLExecutor() engine_base.MiQLExecutorCore {
	return tinyDBCore.TinyMiQLExecutor
}

func NewTinyEngine(db_name string) *TinyDBCore {

	return &TinyDBCore{
		TinyDBLoader: &components.TinyDBLoaderImp{
			DBName: db_name,
		},
		Name: db_name,
	}
}
