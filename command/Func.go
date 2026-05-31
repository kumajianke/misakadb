package command

import (
	"errors"
	"fmt"
	"misakadb/clilog"
	"misakadb/config"
	"misakadb/miusers"
	onces "misakadb/network/Onces"
	"misakadb/network/context"
)

func (dispatch *MiqlCommDispatch) ImpGetServiceInfo(
	serviceContext *context.ServiceConnContext,
	arg_command []string,
) error {
	sysConfigs := config.GetGlobalMisakaConfigure()
	jsonStr := config.ConvertConfigureToJSON(sysConfigs)
	hideInfo := sysConfigs.Service.HideInfo

	clilog.Info(fmt.Sprintf("hideInfo %v", hideInfo))

	var err error
	if !hideInfo {
		err = serviceContext.Send(
			"[ok]" + jsonStr,
		)
	} else {
		clilog.Info("get-service-info命令已被禁用")
		err = serviceContext.Send("[error]service disable the function")
	}

	if err != nil {
		clilog.Error(
			fmt.Sprintf(
				"[%s] command `get-service-info` failed, err: %v",
				(serviceContext.Conn).RemoteAddr().String(),
				err,
			),
		)
		return err
	}

	clilog.Info(
		fmt.Sprintf(
			"[%s] command `get-service-info` success",
			(serviceContext.Conn).RemoteAddr().String(),
		),
	)

	return nil
}

func (dispatch *MiqlCommDispatch) ImpExit(
	serviceContext *context.ServiceConnContext,
	arg_command []string,
) error {
	// 关闭连接
	onces.NewSafeConn((serviceContext.Conn)).Close()
	serviceContext = nil // 回收内存
	return nil
}

func (dispatch *MiqlCommDispatch) ImpLogin(
	serviceContext *context.ServiceConnContext,
	arg_command []string,
) error {
	var (
		username string
		password string
	)
	if len(arg_command) == 2 {
		username = arg_command[0]
		password = arg_command[1]
	} else {
		serviceContext.Send("[err]the miql or command is invalid input")
		return errors.New("error arguments")
	}
	err := miusers.NewUserManager().VerifyPassword(username, password)
	if err != nil {
		serviceContext.Send("[err] username and password can not match.")
		return err
	}
	serviceContext.LoginUser = username // 记录上下文
	serviceContext.Send("[ok] login success.")
	return nil
}
