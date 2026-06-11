package clilog

import (
	"fmt"
	"io"
	tuimode "misakadb/tui/tui-mode"
	"os"
)

// 使用带缓冲的 channel，避免非 CLI 模式下阻塞
// 缓冲区大小为 100，可根据实际需求调整
var chan_history_cli = make(chan string, 100)

const reset = "\033[0m"

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
	prefix := fmt.Sprintf("%s[%s]%s", color, level, reset)

	if tuimode.GetTuiMode() != tuimode.CLIMode {
		// 非 CLI 模式：尝试将历史记录发送到 channel（不阻塞）
		select {
		case chan_history_cli <- prefix:
			// 成功发送
		default:
			// channel 已满，丢弃（避免阻塞）
		}
		return
	}

	// CLI 模式：直接输出
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
