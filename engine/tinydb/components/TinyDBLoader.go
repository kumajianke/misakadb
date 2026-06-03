package components

import (
	"encoding/json"
	"errors"
	"misakadb/clilog"
	mson "misakadb/engine/Mson"
	engine_base "misakadb/engine/base"
	filejson "misakadb/engine/tinydb/FileJson"
	generashares "misakadb/genera_shares"
	"os"
	"os/exec"
	"sync"
	"time"
)

/**
* 记录器对象 用来JSON序列化
 */
type TinyDBRecorder struct {
}

type TinyDBLoaderImp struct {
	engine_base.BaseLoaderCore

	Locker sync.Mutex
	DBName string
}

func (this *TinyDBLoaderImp) WriteLoader(log mson.MsonParse) error {
	this.Locker.Lock()
	defer this.Locker.Unlock()

	return nil
}

func (this *TinyDBLoaderImp) ReadLoader(log mson.MsonParse) error {
	return nil
}

func (this *TinyDBLoaderImp) InitLoader(log mson.MsonParse) error {
	this.Locker.Lock()
	defer this.Locker.Unlock()
	data_path := "./db-datas/" + log.Name
	all_files, err := os.ReadDir(data_path)
	if err != nil {
		clilog.Error("[err]InitLoader error: get files error")
		return err
	}
	if len(all_files) > 0 {
		clilog.Error("[err]new db folder has other files!")
		return errors.New("[err]new db folder has other files!")
	}
	err = os.Mkdir(data_path+"/.db", 0600)
	if err != nil {
		clilog.Error("[err]init db folder create error!")
		return errors.New("[err]init db folder create error!")
	}
	fileName := data_path + "/.db/meta.json"
	metaJson := filejson.NewTinyDBMeta(
		this.DBName,
		make([]string, 0),
		time.Now().Format("2006-1-2"),
	)
	jsonData, err := json.Marshal(metaJson)
	if err != nil {
		clilog.Error("[err]InitLoader error: JsonData error")
		return err
	}
	err = os.WriteFile(fileName, []byte(jsonData), 0600)
	if err != nil {
		clilog.Error("[err]InitLoader error: JsonData error")
		return err
	}
	if generashares.IsWindows() {
		err = exec.Command("attrib", "+h", data_path+"/.db").Run()
		if err != nil {
			clilog.Error("Window platform can not hide the .db folder ")
			return err
		}
	}

	return nil
}
