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

func GetDBEngine(dbname string) (engine_base.BaseEngineCore, error) {
	cache_data_string := GetDBEngineByRC(dbname)
	if cache_data_string != "" {
		return NewEngine(cache_data_string, dbname), nil // 缓存命中
	}

	engine_base, err := engine_base.ShareLoaderDBMetaName(dbname)
	if err != nil {
		return nil, err
	}

	RegisterCenter.RegisterCenterInstance.MapperDBEngine.Store("dbname", engine_base)
	// 回填数据

	return NewEngine(engine_base, dbname), nil
}

func GetDBEngineByRC(dbname string) string {
	rc := RegisterCenter.RegisterCenterInstance
	if RegisterCenter.RegisterCenterInstance == nil {
		return ""
	}
	cache_engine, flag := rc.MapperDBEngine.Load(dbname)
	if flag {
		engine, ok := cache_engine.(string)
		if ok {
			return engine
		}
	}
	return ""
}
