package clilog

import (
	"fmt"
	"io"
	tuimode "misakadb/tui/tui-mode"
	"os"
	"sync"
)

// 使用带缓冲的 channel，避免非 CLI 模式下阻塞
var chan_history_cli = make(chan string, 100)

const reset = "\033[0m"

// tuiRenderer 是由 TUI 模块注册的日志渲染回调。
// 签名：func(ansiColor, level string, args ...any)
// 非 CLIMode 时若此函数不为 nil，则调用它替代存 channel。
var (
	tuiRenderer   func(color, level string, args ...any)
	tuiRendererMu sync.RWMutex
)

// RegisterTUIRenderer 注册（或注销）termbox 日志渲染器。
// 传入 nil 表示注销，恢复 channel 缓存模式。
func RegisterTUIRenderer(fn func(color, level string, args ...any)) {
	tuiRendererMu.Lock()
	defer tuiRendererMu.Unlock()
	tuiRenderer = fn
}

func Info(args ...any) {
	printlnWithStyle(os.Stdout, "\033[1;34m", "INFO", args...)
}

func Success(args ...any) {
	printlnWithStyle(os.Stdout, "\033[1;32m", "SUCCESS", args...)
}

func Warning(args ...any) {
	printlnWithStyle(os.Stdout, "\033[1;33m", "WARNING", args...)
}

func Error(args ...any) {
	printlnWithStyle(os.Stderr, "\033[1;31m", "ERROR", args...)
}

func printlnWithStyle(writer io.Writer, color string, level string, args ...any) {
	if tuimode.GetTuiMode() != tuimode.CLIMode {
		tuiRendererMu.RLock()
		renderer := tuiRenderer
		tuiRendererMu.RUnlock()

		if renderer != nil {
			// 交由 TUI 模块在 termbox 画布上渲染
			renderer(color, level, args...)
		} else {
			// 无渲染器：暂存到 channel（不阻塞）
			prefix := fmt.Sprintf("%s[%s]%s", color, level, reset)
			select {
			case chan_history_cli <- prefix:
			default: // channel 已满，丢弃
			}
		}
		return
	}

	// CLIMode：直接输出到终端
	prefix := fmt.Sprintf("%s[%s]%s", color, level, reset)
	line := append([]any{prefix}, args...)
	fmt.Fprintln(writer, line...)
}

// GetHistoryChannel 返回历史记录 channel，供 TUI 模式消费
func GetHistoryChannel() <-chan string {
	return chan_history_cli
}

// DrainHistory 清空历史记录 channel 中的所有待处理消息
func DrainHistory() []string {
	var history []string
	for {
		select {
		case msg := <-chan_history_cli:
			history = append(history, msg)
		default:
			return history
		}
	}
}
