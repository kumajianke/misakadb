package engine_dispatch

import (
	engine_base "misakadb/engine/base"
	engine "misakadb/engine/tinydb"
)

func NewEngine(engineName string, db_name string) engine_base.BaseEngineCore {
	if engineName == "tinydb" {
		return engine.NewTinyEngine(db_name)
	}

	return nil
}

func GetDBEngine(dbName string) engine_base.BaseEngineCore {
	// TODO 通过数据库名字获取对应的引擎
	// 可以在RC优先注册一个数据库对应的引擎的map 只有找不到再去本地读取对应的文件信息
	return nil
}
