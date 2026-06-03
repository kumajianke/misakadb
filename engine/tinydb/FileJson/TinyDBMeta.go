package filejson

type TinyDBMeta struct {
	DBName     string   `json:"db_name"`
	AllTables  []string `json:"all_tables"`
	CreateTime string   `json:"create_time"`
}

func NewTinyDBMeta(dbName string, allTables []string, createTime string) *TinyDBMeta {
	return &TinyDBMeta{
		DBName:     dbName,
		AllTables:  allTables,
		CreateTime: createTime,
	}
}
