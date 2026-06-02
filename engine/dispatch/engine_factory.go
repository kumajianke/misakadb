package engine_dispatch

import (
	engine_base "misakadb/engine/base"
	engine "misakadb/engine/tinydb"
)

func NewEngine(engineName string) engine_base.BaseEngineCore {
	if engineName == "tinydb" {
		return engine.NewTinyEngine()
	}

	return nil
}
