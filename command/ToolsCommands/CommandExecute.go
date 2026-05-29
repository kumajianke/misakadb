package toolscommands

import (
	"fmt"
	"misakadb/clilog"
	"misakadb/miusers"
	"misakadb/safe"
	"os"
)

func CommandExecute(command_all []string) {
	main_command := command_all[0]
	userManager := miusers.NewUserManager()

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
		userManager.InitUser()
		userload := userManager.LoadUserFile()
		clilog.Info(fmt.Sprintf("一共初始化 %d 个用户.", len(userload)))

	case "add-user":
		if len(command_all) != 3 {
			clilog.Error("错误用法，正确用法misaka-tools add-user 用户名 密码")
			os.Exit(0)
		}
		username := command_all[1]
		password := command_all[2]
		all_user := userManager.AddUser(username, password)
		clilog.Success(fmt.Sprintf("添加用户 %s 成功，当前还有 %d 个用户", username, len(all_user)))
	default:
		clilog.Error("错误的命令, 你无法这么实现。")
	}
}
