package filejson

import (
	engine_base "misakadb/engine/base"
)

type TinyDBMeta struct {
	engine_base.BaseDBMeta
}

func NewTinyDBMeta(dbName string, allTables []string, createTime string) *TinyDBMeta {
	return &TinyDBMeta{
		BaseDBMeta: engine_base.BaseDBMeta{
			DBName:     dbName,
			AllTables:  allTables,
			CreateTime: createTime,
			Engine:     "tinydb",
		},
	}
}
