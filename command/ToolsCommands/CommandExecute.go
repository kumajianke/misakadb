package toolscommands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"misakadb/clilog"
	"misakadb/miusers"
	"misakadb/safe"
	"os"
	"strings"
)

func CommandExecute(command_all []string) {
	main_command := command_all[0]
	userManager := miusers.NewUserManager()
	reader := bufio.NewReader(os.Stdin)

	switch main_command {
	case "sys-init":
		_, userErr := os.Stat(miusers.UserFile)
		_, keyErr := os.Stat("./profiles/master.mikey")

		if userErr == nil {
			clilog.Error("数据库已经初始化了用户数据，无法再进行序列化！")
			os.Exit(0)
		}
		if userErr != nil && !errors.Is(userErr, os.ErrNotExist) {
			clilog.Error("无法检查用户文件状态:", userErr)
			os.Exit(0)
		}

		if keyErr != nil && !errors.Is(keyErr, os.ErrNotExist) {
			clilog.Error("无法检查 AES 密钥文件状态:", keyErr)
			os.Exit(0)
		}

		if errors.Is(userErr, os.ErrNotExist) && errors.Is(keyErr, os.ErrNotExist) {
			if err := os.MkdirAll("./profiles", 0700); err != nil {
				clilog.Error("创建 profiles 目录失败:", err)
				os.Exit(0)
			}
			safe.InitPassword()
		} else if errors.Is(userErr, os.ErrNotExist) && keyErr == nil {
			clilog.Info("检测到已有 AES 密钥，继续复用该密钥初始化用户数据。")
		} else if userErr == nil && errors.Is(keyErr, os.ErrNotExist) {
			clilog.Error("检测到 user.dat 存在但 master.mikey 缺失，数据状态异常，请手动修复后再执行。")
			os.Exit(0)
		}

		rootSecret, err := userManager.InitUser()
		if err != nil {
			clilog.Error("初始化用户失败:", err)
			os.Exit(0)
		}

		userload, err := userManager.LoadUserFile()
		if err != nil {
			clilog.Error("读取初始化后的用户数据失败:", err)
			os.Exit(0)
		}

		clilog.Success(
			fmt.Sprintf(
				"初始化用户信息表完毕, 这是 root 的密钥 %s。不包含前后符号，请立即执行：misaka-tools chpwd root",
				rootSecret,
			),
		)
		clilog.Info(fmt.Sprintf("一共初始化 %d 个用户.", len(userload)))

	case "add-user":
		if len(command_all) != 3 {
			clilog.Error("错误用法，正确用法misaka-tools add-user 用户名 密码")
			os.Exit(0)
		}
		username := command_all[1]
		password := command_all[2]
		allUser, err := userManager.AddUser(username, password)
		if err != nil {
			clilog.Error(err)
			os.Exit(0)
		}
		clilog.Success(fmt.Sprintf("添加用户 %s 成功，当前还有 %d 个用户", username, len(allUser)))

	case "change-password", "chpwd":
		if len(command_all) != 2 {
			clilog.Error("错误用法，正确用法misaka-tools chpwd 用户名")
			os.Exit(0)
		}
		if err := requireRootAuth(userManager, reader); err != nil {
			clilog.Error(err)
			os.Exit(0)
		}
		newPassword, err := promptInput(reader, "请输入新的用户密码: ")
		if err != nil {
			clilog.Error("读取新密码失败:", err)
			os.Exit(0)
		}
		if err := userManager.ChangePassword(command_all[1], newPassword); err != nil {
			clilog.Error(err)
			os.Exit(0)
		}
		clilog.Success(fmt.Sprintf("用户 %s 的密码修改成功", command_all[1]))

	case "chmod":
		if len(command_all) != 3 {
			clilog.Error("错误用法，正确用法misaka-tools chmod 用户名 角色")
			os.Exit(0)
		}
		if err := requireRootAuth(userManager, reader); err != nil {
			clilog.Error(err)
			os.Exit(0)
		}
		if err := userManager.ChangeRole(command_all[1], command_all[2]); err != nil {
			clilog.Error(err)
			os.Exit(0)
		}
		clilog.Success(fmt.Sprintf("用户 %s 的角色已修改为 %s", command_all[1], command_all[2]))

	case "remove":
		if len(command_all) != 2 {
			clilog.Error("错误用法，正确用法misaka-tools remove 用户名")
			os.Exit(0)
		}
		if err := requireRootAuth(userManager, reader); err != nil {
			clilog.Error(err)
			os.Exit(0)
		}
		if err := userManager.RemoveUser(command_all[1]); err != nil {
			clilog.Error(err)
			os.Exit(0)
		}
		clilog.Success(fmt.Sprintf("用户 %s 已删除", command_all[1]))

	case "remote":
		if len(command_all) != 3 {
			clilog.Error("错误用法，正确用法misaka-tools remote 用户名 true|false")
			os.Exit(0)
		}
		if err := requireRootAuth(userManager, reader); err != nil {
			clilog.Error(err)
			os.Exit(0)
		}
		flag, err := parseRemoteFlag(command_all[2])
		if err != nil {
			clilog.Error(err)
			os.Exit(0)
		}
		if err := userManager.SetRemote(command_all[1], flag); err != nil {
			clilog.Error(err)
			os.Exit(0)
		}
		clilog.Success(fmt.Sprintf("用户 %s 的远程登录已设置为 %t", command_all[1], flag))

	case "admin-cli":
		if err := requireRootAuth(userManager, reader); err != nil {
			clilog.Error(err)
			os.Exit(0)
		}
		runAdminCLI(userManager, reader)
	default:
		clilog.Error("错误的命令, 你无法这么实现。")
	}
}

func promptInput(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func requireRootAuth(userManager *miusers.UserManager, reader *bufio.Reader) error {
	username, err := promptInput(reader, "请输入执行用户名: ")
	if err != nil {
		return err
	}
	password, err := promptInput(reader, "请输入执行用户密码: ")
	if err != nil {
		return err
	}
	if err := userManager.VerifyRole(username, "root"); err != nil {
		return fmt.Errorf("root 权限鉴权失败: %w", err)
	}
	if err := userManager.VerifyPassword(username, password); err != nil {
		return fmt.Errorf("用户 %s 鉴权失败: %w", username, err)
	}
	return nil
}

func parseRemoteFlag(flag string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(flag)) {
	case "1", "true", "on", "enable", "enabled", "yes":
		return true, nil
	case "0", "false", "off", "disable", "disabled", "no":
		return false, nil
	default:
		return false, errors.New("remote 标记仅支持 true/false、on/off、1/0")
	}
}

func runAdminCLI(userManager *miusers.UserManager, reader *bufio.Reader) {
	clilog.Success("已进入 admin-cli，root 角色鉴权成功。输入 help 查看命令，输入 exit 退出。")

	for {
		commandLine, err := promptInput(reader, "misaka-admin> ")
		if err != nil {
			clilog.Error("读取命令失败:", err)
			return
		}
		if commandLine == "" {
			continue
		}

		args := strings.Fields(commandLine)
		switch args[0] {
		case "exit", "quit":
			clilog.Info("admin-cli 已退出。")
			return
		case "help":
			printAdminHelp()
		case "change-password", "chpwd":
			if len(args) != 2 {
				clilog.Error("正确用法: chpwd <username>")
				continue
			}
			newPassword, err := promptInput(reader, "请输入新的用户密码: ")
			if err != nil {
				clilog.Error("读取新密码失败:", err)
				continue
			}
			if err := userManager.ChangePassword(args[1], newPassword); err != nil {
				clilog.Error(err)
				continue
			}
			clilog.Success(fmt.Sprintf("用户 %s 的密码修改成功", args[1]))
		case "chmod":
			if len(args) != 3 {
				clilog.Error("正确用法: chmod <username> <role>")
				continue
			}
			if err := userManager.ChangeRole(args[1], args[2]); err != nil {
				clilog.Error(err)
				continue
			}
			clilog.Success(fmt.Sprintf("用户 %s 的角色已修改为 %s", args[1], args[2]))
		case "remove":
			if len(args) != 2 {
				clilog.Error("正确用法: remove <username>")
				continue
			}
			if err := userManager.RemoveUser(args[1]); err != nil {
				clilog.Error(err)
				continue
			}
			clilog.Success(fmt.Sprintf("用户 %s 已删除", args[1]))
		case "remote":
			if len(args) != 3 {
				clilog.Error("正确用法: remote <username> <flag>")
				continue
			}
			flag, err := parseRemoteFlag(args[2])
			if err != nil {
				clilog.Error(err)
				continue
			}
			if err := userManager.SetRemote(args[1], flag); err != nil {
				clilog.Error(err)
				continue
			}
			clilog.Success(fmt.Sprintf("用户 %s 的远程登录已设置为 %t", args[1], flag))
		case "add-user":
			if len(args) != 3 {
				clilog.Error("正确用法: add-user <username> <password>")
				continue
			}
			if _, err := userManager.AddUser(args[1], args[2]); err != nil {
				clilog.Error(err)
				continue
			}
			clilog.Success(fmt.Sprintf("用户 %s 添加成功", args[1]))
		default:
			clilog.Error("未知命令，输入 help 查看 admin-cli 支持的命令。")
		}
	}
}

func printAdminHelp() {
	clilog.Info("可用命令:")
	fmt.Println("  chpwd <username>")
	fmt.Println("  change-password <username>")
	fmt.Println("  chmod <username> <role>")
	fmt.Println("  remove <username>")
	fmt.Println("  remote <username> <flag>")
	fmt.Println("  add-user <username> <password>")
	fmt.Println("  help")
	fmt.Println("  exit")
}
