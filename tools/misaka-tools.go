package main

import (
	"flag"
	"misakadb/clilog"
	toolscommands "misakadb/command/ToolsCommands"
	_ "misakadb/plugins/pluginbridge"
)

func main() {
	flag.Parse()
	command_all := flag.Args()

	if len(command_all) == 0 {
		clilog.Success("你好呀，有什么可以帮助你的?")
	}

	toolscommands.CommandExecute(command_all)
}
