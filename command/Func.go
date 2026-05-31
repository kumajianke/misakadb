package command

import (
	"fmt"
	"misakadb/clilog"
	"misakadb/config"
	onces "misakadb/network/Onces"
	"misakadb/network/context"
)

func (dispatch *CommandDispatch) ImpGetServiceInfo(
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

func (dispatch *CommandDispatch) ImpExit(
	serviceContext *context.ServiceConnContext,
	arg_command []string,
) error {
	// 关闭连接
	onces.NewSafeConn((serviceContext.Conn)).Close()
	serviceContext = nil // 回收内存
	return nil
}
