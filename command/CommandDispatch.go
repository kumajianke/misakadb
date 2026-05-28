package command

import (
	"misakadb/config"
	onces "misakadb/misaka_network/Onces"
	"misakadb/misaka_network/active"
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
	}
}
