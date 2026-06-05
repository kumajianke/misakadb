package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"misakadb/clilog"
	mson "misakadb/engine/Mson"
	engine_base "misakadb/engine/base"
	filejson "misakadb/engine/tinydb/FileJson"
	generashares "misakadb/genera_shares"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

/**
* 记录器对象 用来JSON序列化
 */
type TinyDBRecorder struct {
}

type TinyDBLoaderImp struct {
	Locker      engine_base.BaseLockerCore
	localLocker engine_base.EngineLockerSupport
	DBName      string
}

var _ engine_base.BaseLoaderCore = (*TinyDBLoaderImp)(nil)

func (this *TinyDBLoaderImp) lockerCore() engine_base.BaseLockerCore {
	if this.Locker != nil {
		return this.Locker
	}

	if this.localLocker.LockNamespace == "" {
		this.localLocker.LockNamespace = "tinydb:" + this.DBName
	}

	return &this.localLocker
}

func (this *TinyDBLoaderImp) WriteLoader(log mson.MsonParse) error {
	unlock, err := this.lockerCore().Lock()
	if err != nil {
		return err
	}
	defer unlock()

	return nil
}

func (this *TinyDBLoaderImp) ReadLoader(log mson.MsonParse) error {
	return nil
}

func (this *TinyDBLoaderImp) InitLoader(log mson.MsonParse) error {
	this.DBName = log.Name

	unlock, err := this.lockerCore().GetRowLock(this.DBName)
	if err != nil {
		return err
	}
	defer unlock()

	// 创建 数据库根目录
	newPath := filepath.Join(".", "db-datas", log.Name)
	_, erros_file := os.Stat(newPath)

	if erros_file == nil {
		return errors.New("database is exist!")
	} else if os.IsNotExist(erros_file) {
		dbRootPath := filepath.Join(".", "db-datas")
		dbRootInfo, rootErr := os.Stat(dbRootPath)
		if rootErr == nil && !dbRootInfo.IsDir() {
			dirErr := fmt.Errorf("%s is not a directory", dbRootPath)
			clilog.Error("[err] stat db root error: " + dirErr.Error())
			return dirErr
		}

		err = os.Mkdir(newPath, 0700)
		if err != nil {
			clilog.Error("[err] create dir error: " + err.Error())
			return err
		}
	} else {
		clilog.Error("[err] stat db dir error: " + erros_file.Error())
		return erros_file
	}

	// 创建内部 .db文件夹
	dbMetaDir := filepath.Join(newPath, ".db")
	err = os.Mkdir(dbMetaDir, 0700)
	if err != nil {
		clilog.Error("[err]init db folder create error!")
		return errors.New("init db folder create error!")
	}
	fileName := filepath.Join(dbMetaDir, "meta.json")
	metaJson := filejson.NewTinyDBMeta(
		this.DBName,
		make([]string, 0),
		time.Now().Format("2006-1-2"),
	)
	jsonData, err := json.Marshal(metaJson)
	if err != nil {
		clilog.Error("[err]InitLoader error: JsonData error")
		return errors.New("InitLoader error: JsonData error: " + err.Error())
	}
	err = os.WriteFile(fileName, []byte(jsonData), 0600)
	if err != nil {
		clilog.Error("[err]InitLoader error: JsonData error")
		return errors.New("InitLoader error: JsonData error: " + err.Error())
	}
	if generashares.IsWindows() {
		err = exec.Command("attrib", "+h", dbMetaDir).Run()
		if err != nil {
			clilog.Error("Window platform can not hide the .db folder ")
			return errors.New("Window platform can not hide the .db folder ")
		}
	}

	return nil
}
