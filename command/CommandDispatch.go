package command

import (
	"misakadb/network/context"
	"reflect"
)

type CommandDispatch struct {
	GetServiceInfo func(serviceContext *context.ServiceConnContext) error `mapper:"get-service-info"`
	Exit           func(serviceContext *context.ServiceConnContext) error `mapper:"exit"`
}

func NewCommandDispatch() *CommandDispatch {
	dispatch := &CommandDispatch{}
	dispatch.GetServiceInfo = dispatch.ImpGetServiceInfo
	dispatch.Exit = dispatch.ImpExit
	return dispatch
}

func (dispatch *CommandDispatch) Dispatch(
	serviceContext *context.ServiceConnContext,
	command string,
) error {
	dispatchValue := reflect.ValueOf(dispatch).Elem()
	dispatchType := dispatchValue.Type()
	for i := 0; i < dispatchType.NumField(); i++ {
		field := dispatchType.Field(i)
		if field.Tag.Get("mapper") == command {
			context := dispatchValue.Field(i).Interface().(func(*context.ServiceConnContext) error)
			return context(serviceContext)
		}
	}

	return nil
}
