package engine_dispatch

import (
	engine_base "misakadb/engine/base"
	engine "misakadb/engine/tinydb"
	"misakadb/network/RegisterCenter"
)

func NewEngine(engineName string, db_name string) engine_base.BaseEngineCore {
	if engineName == "tinydb" {
		return engine.NewTinyEngine(db_name)
	}

	return nil
}

func GetDBEngineByRC(dbname string) engine_base.BaseEngineCore {
	rc := RegisterCenter.RegisterCenterInstance
	if RegisterCenter.RegisterCenterInstance == nil {
		return nil
	}
	cache_engine, flag := rc.MapperDBEngine.Load(dbname)
	if flag {
		engine, ok := cache_engine.(engine_base.BaseEngineCore)
		if ok {
			return engine
		}
	}
	return nil
}
