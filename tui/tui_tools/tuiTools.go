package tui_tools

import (
	"fmt"
	"time"
)

// 清屏
func Clear() {
	fmt.Print("\033[2J\033[H")
}

// 清除当前行
func ClearLine() {
	fmt.Print("\033[2K\r")
}

// 移动光标到指定位置 (x, y)
func MoveCursor(x, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

// 隐藏光标
func HideCursor() {
	fmt.Print("\033[?25l")
}

// 显示光标
func ShowCursor() {
	fmt.Print("\033[?25h")
}

func main() {
	Clear()
	fmt.Println("第一行")
	fmt.Println("第二行")

	// 2秒后清屏
	time.Sleep(2 * time.Second)
	Clear()

	fmt.Println("已清屏")
}
