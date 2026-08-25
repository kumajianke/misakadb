package atomicX

import (
	pluginsloader "misakadb/plugins/pluginsLoader"
)

// 原子： 删除数据库操作
func AtomicDropDB(dbname string) error {
	combo, err := pluginsloader.GetPluginsTaskCombo("drop_db")
	if err != nil {
		return err
	}
	if err, _ := combo([]string{dbname}); err != nil {
		return err
	}
	return nil
}

// TODO 原子:  创建数据库操作
func AtomicCreateDB(dbname string) error {
	return nil
}
