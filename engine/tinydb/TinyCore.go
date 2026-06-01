package engine

import engine "misakadb/engine/base"

type TinyDBCore struct {
	engine.BaseEngineCore
}

func (tinyDBCore *TinyDBCore) Path() string {
	return "./db-datas/$name"
}
