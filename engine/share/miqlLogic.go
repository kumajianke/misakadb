package share

import (
	"errors"
	"misakadb/clilog"
	mson "misakadb/engine/Mson"
	engine_dispatch "misakadb/engine/dispatch"
	"misakadb/miusers"
	"misakadb/network/context"
)

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

	err := dbEngine.DBLoader().InitLoader(*msonParse) // 选择对应的数据库引擎进行初始化

	if err != nil {
		err_string := err.Error()
		serviceContext.Send("[err]" + err_string) // 错误信息的返回
		return err
	}

	serviceContext.Send("[ok]create db is ok!")
	return nil
}

func MiqlDropDB(msonPaese *mson.MsonParse, serviceContext *context.ServiceConnContext) error {
	if msonPaese.Active != "drp-dat" {
		return errors.New("Error Dispatch!")
	}
	username := serviceContext.LoginUser
	err := miusers.NewUserManager().VerifyRole(username, "root")
	if err != nil {
		serviceContext.Send("[err]" + err.Error())
	}

	return nil
}
