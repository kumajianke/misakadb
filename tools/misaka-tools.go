package main

import (
	"flag"
	"fmt"
	"misakadb/clilog"
	"misakadb/miusers"
	"misakadb/safe"
	"os"
)

func main() {
	flag.Parse()
	command_all := flag.Args()
	if len(command_all) == 0 {
		clilog.Success("你好呀，有什么可以帮助你的?")
	}
	main_command := command_all[0]
	switch main_command {
	case "sys-init":
		_, user_err := os.Stat(miusers.UserFile)
		if user_err == nil {
			clilog.Error("数据库已经初始化了用户数据，无法再进行序列化！")
			os.Exit(0)
		}

		_, password_err := os.Stat("./profiles/master.mikey")
		if password_err == nil {
			clilog.Error("数据库已经初始化了密钥数据，无法再进行序列化！")
			os.Exit(0)
		}
		// 密钥初始化
		safe.InitPassword()

		// 用户初始化
		userManager := miusers.NewUserManager()
		userManager.InitUser()
		userload := userManager.LoadUserFile()
		clilog.Info(fmt.Sprintf("一共初始化 %d 个用户.", len(userload)))

	default:
		clilog.Error("错误的命令, 你无法这么实现。")
	}
}
