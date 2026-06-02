package engine

import (
	engine_base "misakadb/engine/base"
	"misakadb/engine/tinydb/components"
)

type TinyDBCore struct {
	TinyDBLoader     engine_base.BaseLoaderCore
	TinyDBBaker      engine_base.BaseBakerCore
	TinyMiQLExecutor engine_base.MiQLExecutorCore
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

func (tinyDBCore *TinyDBCore) Path() string {
	return "./db-datas/$name"
}

func NewTinyEngine() *TinyDBCore {
	return &TinyDBCore{TinyDBLoader: &components.TinyDBLoaderImp{}}
}
