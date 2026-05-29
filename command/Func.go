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
) error {
	sysConfigs := config.GetGlobalMisakaConfigure()
	jsonStr := config.ConvertConfigureToJSON(sysConfigs)
	err := serviceContext.Send(jsonStr)

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
) error {
	// 关闭连接
	onces.NewSafeConn((serviceContext.Conn)).Close()
	serviceContext = nil // 回收内存
	return nil
}
