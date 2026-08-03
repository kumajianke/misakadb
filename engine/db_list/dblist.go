package dblist

import (
	"encoding/json"
	"errors"
	atomic_file_handler "misakadb/atomic/atomicFileHandler"
	"misakadb/clilog"
	"misakadb/lock/global_lock"
	"os"
	"path/filepath"
)

type DBStruct struct {
	DBName   string `json:"db-name"`
	DBEngine string `json:"db-engine"`
}

func StoreDB(dbname string, dbengine string) error {
	_, unlock, err := global_lock.GetOrStoreGlobalLock("dbList", "l")
	if err != nil {
		return err
	}

	defer unlock()

	dbDataDir := filepath.Join(".", "db-datas")
	filePath := filepath.Join(dbDataDir, "alldb.list")

	// 确保 db-datas 目录存在
	err = os.Mkdir(dbDataDir, 0700)
	if err != nil && !os.IsExist(err) {
		clilog.Error("[err]Failed to create db-datas directory: " + err.Error())
		return errors.New("failed to create db-datas directory: " + err.Error())
	}

	var dbList []DBStruct

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if len(fileData) > 0 {
		err = json.Unmarshal(fileData, &dbList)
		if err != nil {
			return err
		}
	}

	newDB := DBStruct{
		DBName:   dbname,
		DBEngine: dbengine,
	}

	for _, value := range dbList {
		if value.DBName == newDB.DBName {
			return errors.New("database is exist!")
		}
	}

	dbList = append(dbList, newDB)

	newData, err := json.MarshalIndent(dbList, "", "  ")
	if err != nil {
		return err
	}

	err = atomic_file_handler.AtomicSyncWriteFile(filePath, newData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func DropDB(name string) error {
	_, unlock, err := global_lock.GetOrStoreGlobalLock("dbList", "l")
	if err != nil {
		return err
	}
	defer unlock()

	filePath := filepath.Join(".", "db-datas", "alldb.list")

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("database list file not found")
		}
		return err
	}

	var dbList []DBStruct
	if len(fileData) > 0 {
		err = json.Unmarshal(fileData, &dbList)
		if err != nil {
			return err
		}
	}

	found := false
	newDBList := make([]DBStruct, 0, len(dbList))
	for _, db := range dbList {
		if db.DBName != name {
			newDBList = append(newDBList, db)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("database not found in list")
	}

	newData, err := json.MarshalIndent(newDBList, "", "  ")
	if err != nil {
		return err
	}

	err = atomic_file_handler.AtomicSyncWriteFile(filePath, newData, 0644)
	if err != nil {
		return err
	}

	return nil
}
