package tui

import (
	"fmt"
	"misakadb/clilog"
	"misakadb/config"
	tuimode "misakadb/tui/tui-mode"
	"os"
	"sync"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// ──────────────────────────────────────────────
// 菜单数据
// ──────────────────────────────────────────────

type menuItem struct {
	key  string
	desc string
}

var menuItems = []menuItem{
	{"SPACE", "Show The Menu"},
	{"N", "Log Mode"},
	{"L", "Lock Mode"},
	{"Q", "Exit (Debug Mode)"},
}

// ──────────────────────────────────────────────
// 日志缓冲区
// ──────────────────────────────────────────────

type logEntry struct {
	fg   termbox.Attribute
	text string
}

const maxLogBuffer = 500 // 全局日志缓冲条数

var (
	logBuffer   []logEntry
	logBufferMu sync.Mutex
)

func ansiToTermbox(ansiColor string) termbox.Attribute {
	switch ansiColor {
	case "\033[1;32m":
		return termbox.ColorGreen | termbox.AttrBold
	case "\033[1;34m":
		return termbox.ColorBlue | termbox.AttrBold
	case "\033[1;33m":
		return termbox.ColorYellow | termbox.AttrBold
	case "\033[1;31m":
		return termbox.ColorRed | termbox.AttrBold
	default:
		return termbox.ColorWhite
	}
}

// tuiLogRenderer 由任意 goroutine 调用，追加日志并用 Interrupt 唤醒渲染 goroutine
func tuiLogRenderer(color, level string, args ...any) {
	msg := fmt.Sprintf("[%s] %s", level, fmt.Sprint(args...))
	fg := ansiToTermbox(color)

	logBufferMu.Lock()
	logBuffer = append(logBuffer, logEntry{fg: fg, text: msg})
	if len(logBuffer) > maxLogBuffer {
		logBuffer = logBuffer[len(logBuffer)-maxLogBuffer:]
	}
	logBufferMu.Unlock()

	termbox.Interrupt() // 唤醒 PollEvent，由渲染 goroutine 刷新画面
}

// ──────────────────────────────────────────────
// termbox 绘制工具
// ──────────────────────────────────────────────

// drawText 逐字符绘制（全角字符占2列），返回结束列号
func drawText(x, y int, fg, bg termbox.Attribute, text string) int {
	col := x
	for _, ch := range text {
		termbox.SetCell(col, y, ch, fg, bg)
		col += runewidth.RuneWidth(ch)
	}
	return col
}

// fillSpaces 用空格填充 [x, endX) 区间
func fillSpaces(x, endX, y int, fg, bg termbox.Attribute) {
	for c := x; c < endX; c++ {
		termbox.SetCell(c, y, ' ', fg, bg)
	}
}

// ──────────────────────────────────────────────
// 日志模式渲染
// ──────────────────────────────────────────────

// renderLogMode 渲染日志查看模式
func renderLogMode() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()

	// 用空格填充整个屏幕（确保清空所有内容）
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	// 顶部提示栏（左对齐）
	// 日志模式提示信息中中英文混用，建议统一为中文或英文
	hint := "Log Mode - Press SPACE to return to the menu | Q to exit"
	drawText(0, 0, termbox.ColorCyan|termbox.AttrBold, termbox.ColorDefault, hint)

	// 分隔线
	for x := 0; x < w; x++ {
		termbox.SetCell(x, 1, '─', termbox.ColorDarkGray, termbox.ColorDefault)
	}

	// 渲染日志
	logBufferMu.Lock()
	defer logBufferMu.Unlock()

	// 可用行数（去掉顶部两行）
	availableLines := h - 2
	if availableLines < 1 {
		availableLines = 1
	}

	// 计算要显示的日志范围（从最新的日志开始往前取）
	startIdx := 0
	if len(logBuffer) > availableLines {
		startIdx = len(logBuffer) - availableLines
	}

	// 逐行显示日志（从左上角开始，左对齐）
	for i, entry := range logBuffer[startIdx:] {
		row := i + 2
		if row >= h {
			break
		}
		drawText(0, row, entry.fg, termbox.ColorDefault, entry.text)
	}

	termbox.Flush()
}

// ──────────────────────────────────────────────
// 菜单模式渲染
// ──────────────────────────────────────────────

// Menu 全屏绘制居中快捷键 Panel
func Menu() {
	// 完全清空屏幕背景
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	w, h := termbox.Size()

	// 用空格填充整个屏幕（确保清空所有内容）
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	const panelW = 50
	panelH := 5 + len(menuItems)

	// 在完整屏幕中绝对居中
	startX := (w - panelW) / 2
	startY := (h - panelH) / 2
	if startY < 0 {
		startY = 0
	}

	colorBorder := termbox.ColorCyan
	colorTitle := termbox.ColorWhite | termbox.AttrBold
	colorKey := termbox.ColorYellow | termbox.AttrBold
	colorDesc := termbox.ColorWhite
	colorSep := termbox.ColorDarkGray
	bg := termbox.ColorDefault

	// ── 顶部边框 ──
	drawText(startX, startY, colorBorder, bg, "╔")
	for i := 1; i < panelW-1; i++ {
		drawText(startX+i, startY, colorBorder, bg, "═")
	}
	drawText(startX+panelW-1, startY, colorBorder, bg, "╗")

	// ── 标题 ──
	title := "Misaka DB —— 主菜单"
	titleRow := startY + 1
	titleW := runewidth.StringWidth(title)
	leftPad := (panelW - 2 - titleW) / 2

	drawText(startX, titleRow, colorBorder, bg, "║")
	fillSpaces(startX+1, startX+1+leftPad, titleRow, colorTitle, bg)
	afterTitle := drawText(startX+1+leftPad, titleRow, colorTitle, bg, title)
	fillSpaces(afterTitle, startX+panelW-1, titleRow, colorTitle, bg)
	drawText(startX+panelW-1, titleRow, colorBorder, bg, "║")

	// ── 分隔线 ──
	sepRow := startY + 2
	drawText(startX, sepRow, colorBorder, bg, "╠")
	for i := 1; i < panelW-1; i++ {
		drawText(startX+i, sepRow, colorSep, bg, "─")
	}
	drawText(startX+panelW-1, sepRow, colorBorder, bg, "╣")

	// ── 快捷键条目 ──
	for idx, item := range menuItems {
		row := startY + 3 + idx
		drawText(startX, row, colorBorder, bg, "║")

		col := startX + 2
		col = drawText(col, row, colorSep, bg, "[")
		col = drawText(col, row, colorKey, bg, item.key)
		col = drawText(col, row, colorSep, bg, "]")
		col = drawText(col, row, colorDesc, bg, "  ")
		afterDesc := drawText(col, row, colorDesc, bg, item.desc)
		fillSpaces(afterDesc, startX+panelW-1, row, colorDesc, bg)
		drawText(startX+panelW-1, row, colorBorder, bg, "║")
	}

	// ── 底部边框 ──
	bottomRow := startY + 3 + len(menuItems)
	drawText(startX, bottomRow, colorBorder, bg, "╚")
	for i := 1; i < panelW-1; i++ {
		drawText(startX+i, bottomRow, colorBorder, bg, "═")
	}
	drawText(startX+panelW-1, bottomRow, colorBorder, bg, "╝")

	termbox.Flush()
}

func modCliMode() {
	termbox.Close()
	if err := termbox.Init(); err != nil {
		return
	}

	tuimode.SetTuiMode(tuimode.CLIMode)
	renderLogMode() // 先渲染界面
	clilog.RegisterTUIRenderer(tuiLogRenderer)
}

// ──────────────────────────────────────────────
// 事件循环
// ──────────────────────────────────────────────

func listenAndGoto() {
	for {
		ev := termbox.PollEvent()

		switch ev.Type {

		case termbox.EventInterrupt:
			switch tuimode.GetTuiMode() {
			case tuimode.MenuMode:
				Menu()
			case tuimode.CLIMode:
				renderLogMode()
			case tuimode.LockMode:
			}

		case termbox.EventResize:
			switch tuimode.GetTuiMode() {
			case tuimode.MenuMode:
				Menu()
			case tuimode.CLIMode:
				renderLogMode()
			case tuimode.LockMode:
			}

		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeySpace:
				// 按 SPACE：重新初始化 termbox 以彻底清空屏幕
				termbox.Close()
				if err := termbox.Init(); err != nil {
					return
				}
				tuimode.SetTuiMode(tuimode.MenuMode)
				Menu()
			}

			switch ev.Ch {
			case 'n', 'N':
				modCliMode()
			case 'l', 'L':
				// TODO: 锁查看模式
				tuimode.SetTuiMode(tuimode.LockMode)

			case 'q', 'Q':
				// 退出：关闭 termbox 后恢复正常终端
				clilog.RegisterTUIRenderer(nil)
				termbox.Close()
				if config.GetGlobalServiceConfigure().Debug {
					os.Exit(0)
				}
			}
		}
	}
}

// ──────────────────────────────────────────────
// 启动入口
// ──────────────────────────────────────────────

func Thread_start() {
	if err := termbox.Init(); err != nil {
		clilog.Error("[misaka termtools] failed to initialize termbox: " + err.Error())
		return
	}

	// 注册 TUI 日志渲染器
	clilog.RegisterTUIRenderer(tuiLogRenderer)

	// 默认启动时不显示任何内容，保持 CLIMode（正常终端输出）
	tuimode.SetTuiMode(tuimode.CLIMode)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
	modCliMode()

	clilog.Success("[misaka termtools] enable - Press SPACE to show menu")

	go func() {
		defer func() {
			clilog.RegisterTUIRenderer(nil)
			termbox.Close()
		}()
		listenAndGoto()
	}()
}
