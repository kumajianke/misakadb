package tui

import (
	"fmt"
	"misakadb/clilog"
	"misakadb/config"
	tuimode "misakadb/tui/tui-mode"
	"misakadb/tui/tui_tools"
	"os"

	"github.com/nsf/termbox-go"
)

func Menu() {
	println("[enable] Misaka DB Menu")
	println("[n] 日志模式")
	println("[l] 锁查看模式")
}

func listenAndGoto() {
	for {
		ev := termbox.PollEvent()

		// 改成 == 而不是 !=
		if ev.Type == termbox.EventKey {
			switch ev.Key {
			case termbox.KeySpace:
				fmt.Println("input detected")
				tui_tools.Clear()
				Menu()

			}

			// 处理字符按键
			switch ev.Ch {
			case 'n', 'N':
				tuimode.SetTuiMode(tuimode.CLIMode)
				tui_tools.Clear()
				clilog.Success("切换到日志模式")
			case 'l', 'L':

			case 'q', 'Q':
				fmt.Println(config.GetGlobalServiceConfigure().Debug)
				if config.GetGlobalServiceConfigure().Debug {
					os.Exit(0)
				}
			}

		}
	}
}

func Thread_start() {
	// 初始化 termbox
	err := termbox.Init()
	if err != nil {
		clilog.Error("[misaka termtools] failed to initialize termbox: " + err.Error())
		return
	}

	clilog.Success("[misaka termtools] enable")

	// 启动监听 goroutine
	go func() {
		defer termbox.Close() // 确保退出时清理
		listenAndGoto()
	}()
}
