package tui

func listenAndGoto() {
	for {
		// TODO 监听输入 切换菜单
	}
}

func thread_start() {
	go func() {
		listenAndGoto()
	}()
}
