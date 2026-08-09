package base_unloader

import "misakadb/clilog"

func Register() error {
	if err := AddTaskType(); err != nil {
		return err
	}

	if err := AddTaskTypeAction(); err != nil {
		return err
	}

	clilog.Success("基础插件加载完毕.")
	return nil
}
