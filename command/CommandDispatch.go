package command

import (
	"fmt"
	"misakadb/clilog"
	"misakadb/config"
	onces "misakadb/network/Onces"
	"misakadb/network/active"
)

type CommandDispatch struct {
}

func NewCommandDispatch() *CommandDispatch {
	return &CommandDispatch{}
}

func (dispatch *CommandDispatch) Dispatch(
	serviceHandler *active.ServiceConnHandler,
	command string,
) {

	switch command {
	case "exit":
		onces.NewSafeConn(*serviceHandler.Conn).ConnClose()
	case "get-service-info":
		sysConfigs := config.GetGlobalMisakaConfigure()
		jsonStr := config.ConvertConfigureToJSON(sysConfigs)
		serviceHandler.Send(jsonStr)
		clilog.Info(
			fmt.Sprintf(
				"[%s] command `get-service-info` success",
				(*serviceHandler.Conn).RemoteAddr().String(),
			),
		)
	}

}
