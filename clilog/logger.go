package clilog

import (
	"fmt"
	"io"
	"os"
)

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
	line := append([]any{prefix}, args...)
	fmt.Fprintln(writer, line...)
}
