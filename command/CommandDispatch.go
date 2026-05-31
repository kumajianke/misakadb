package command

import (
	"fmt"
	"misakadb/clilog"
	"misakadb/network/context"
	"reflect"
	"strings"
)

type MiqlCommDispatch struct {
	GetServiceInfo func(serviceContext *context.ServiceConnContext, args []string) error `mapper:"get-service-info"`
	Login          func(serviceContext *context.ServiceConnContext, args []string) error `mapper:"login"`
	Exit           func(serviceContext *context.ServiceConnContext, args []string) error `mapper:"exit"`
}

func NewMiqlCommDispatch() *MiqlCommDispatch {
	dispatch := &MiqlCommDispatch{}
	dispatch.GetServiceInfo = dispatch.ImpGetServiceInfo
	dispatch.Exit = dispatch.ImpExit
	dispatch.Login = dispatch.ImpLogin

	return dispatch
}

func (dispatch *MiqlCommDispatch) Dispatch(
	serviceContext *context.ServiceConnContext,
	command string,
) error {

	if strings.HasPrefix(command, "mq.") {
		// 这是一个 miql 语句
		if serviceContext.LoginUser == "" {
			serviceContext.Send("[error]you must login first")
			return nil
		}
		command = command[3:]
		clilog.Info(
			fmt.Sprintf("[%s] miql: %s",
				serviceContext.Conn.RemoteAddr(),
				command,
			))

		return nil
	}

	dispatchValue := reflect.ValueOf(dispatch).Elem()
	dispatchType := dispatchValue.Type()

	command_lst := strings.Split(command, " ")
	main_command := command_lst[0]
	var arg_command []string
	if len(command_lst) > 1 {
		arg_command = command_lst[1:]
	} else {
		arg_command = []string{}
	}

	for i := 0; i < dispatchType.NumField(); i++ {
		field := dispatchType.Field(i)
		if field.Tag.Get("mapper") == main_command {

			context := dispatchValue.Field(i).Interface().(func(*context.ServiceConnContext, []string) error)
			return context(serviceContext, arg_command)
		}
	}

	serviceContext.Send("[error]the miql or command is invalid input")
	return nil
}
