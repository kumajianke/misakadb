package share

import (
	engine_base "misakadb/engine/base"
	tinydb "misakadb/engine/tinydb"
)

func NewEngine(engineName string) engine_base.BaseEngineCore {
	if engineName == "tinydb" {
		return &tinydb.TinyDBCore{}
	}

	return nil
}
