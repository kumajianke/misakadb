package share

import (
	"errors"
	"misakadb/network/context"
	"os"
	"strings"
)

func MiqlCreateDB(msonParse *MsonParse, serviceContext *context.ServiceConnContext) error {
	if msonParse.Active != "cre-dat" {
		return errors.New("Error Dispatch!")
	}

	engineName := msonParse.Engine    // 获取到对应的引擎名字
	dbEngine := NewEngine(engineName) // 数据库引擎
	path := dbEngine.Path()
	newPath := strings.Replace(path, "$name", msonParse.Name, 1)
	if _, erros_file := os.Stat(newPath); os.IsNotExist(erros_file) {
		os.Mkdir(newPath, 0700)
	} else {
		serviceContext.Send("[err]database is exist!")
		return nil
	}

	serviceContext.Send("[ok]create db is ok!")
	return nil
}
