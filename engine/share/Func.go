package share

import (
	"errors"
	"misakadb/clilog"
	mson "misakadb/engine/Mson"
	engine_dispatch "misakadb/engine/dispatch"
	"misakadb/network/context"
	"os"
	"strings"
)

func RunMson(msonParse *mson.MsonParse, serviceContext *context.ServiceConnContext) error {
	if msonParse == nil {
		return errors.New("mson is nil")
	}

	switch msonParse.Active {
	case "cre-dat":
		return MiqlCreateDB(msonParse, serviceContext)
	default:
		return serviceContext.Send("[err]unknown miql!")
	}
}

func MiqlCreateDB(msonParse *mson.MsonParse, serviceContext *context.ServiceConnContext) error {
	if msonParse.Active != "cre-dat" {
		return errors.New("Error Dispatch!")
	}

	engineName := msonParse.Engine                                    // 获取到对应的引擎名字
	dbEngine := engine_dispatch.NewEngine(engineName, msonParse.Name) // 数据库引擎
	if dbEngine == nil {
		clilog.Error("未知的引擎诉求")
		if err := serviceContext.Send("[err]未知的引擎诉求"); err != nil {
			return err
		}
		return errors.New("unknown engine request")
	}
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
