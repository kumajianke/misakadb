package tui

import (
	"fmt"
	"misakadb/clilog"
	"misakadb/config"
	pluginsloader "misakadb/plugins/pluginsLoader"
	pluginsX "misakadb/plugins/pluginsx"
	lockshow "misakadb/tui/lock-show"
	tuibase "misakadb/tui/tui-base"
	tuimode "misakadb/tui/tui-mode"
	"os"
	"strings"
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
	{"P", "Plugins Mode"},
	{"Q", "Exit (Debug Mode)"},
}

// ──────────────────────────────────────────────
// 日志缓冲区
// ──────────────────────────────────────────────

type logEntry struct {
	fg   termbox.Attribute
	text string
}

type pluginsInputState int

const (
	pluginsBrowseState pluginsInputState = iota
	pluginsCommandInputState
)

const maxLogBuffer = 500 // 全局日志缓冲条数

var (
	logBuffer         []logEntry
	logBufferMu       sync.Mutex
	pluginModeCommand string
	pluginModeOutput  []string
	pluginModeState   = pluginsBrowseState
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
	msg := fmt.Sprintf("[%s][%s] %s", color, level, fmt.Sprint(args...))
	fg := ansiToTermbox(color)

	logBufferMu.Lock()
	logBuffer = append(logBuffer, logEntry{fg: fg, text: msg})
	if len(logBuffer) > maxLogBuffer {
		logBuffer = logBuffer[len(logBuffer)-maxLogBuffer:]
	}
	logBufferMu.Unlock()

	termbox.Interrupt() // 唤醒 PollEvent，由渲染 goroutine 刷新画面
}

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
	tuibase.DrawText(0, 0, termbox.ColorCyan|termbox.AttrBold, termbox.ColorDefault, hint)

	// 分隔线
	for x := 0; x < w; x++ {
		termbox.SetCell(x, 1, tuibase.HorizontalLineRune(), termbox.ColorDarkGray, termbox.ColorDefault)
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
		tuibase.DrawText(0, row, entry.fg, termbox.ColorDefault, entry.text)
	}

	termbox.Flush()
}

func renderPluginsMode() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	title := "Plugins Mode - Press SPACE to show menu | Q to exit"
	tuibase.DrawText(0, 0, termbox.ColorCyan|termbox.AttrBold, termbox.ColorDefault, title)

	for x := 0; x < w; x++ {
		termbox.SetCell(x, 1, tuibase.HorizontalLineRune(), termbox.ColorDarkGray, termbox.ColorDefault)
	}

	configuredPlugins := make([]string, 0)
	if cfg := config.GetGlobalMisakaConfigure(); cfg != nil {
		configuredPlugins = append(configuredPlugins, cfg.Private.Storage.Plugins...)
	}
	loadedPlugins := pluginsloader.GetLoadedPluginsSnapshot()

	halfW := w / 2
	tuibase.DrawText(0, 2, termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault, fmt.Sprintf("Configured Plugins (%d)", len(configuredPlugins)))
	tuibase.DrawText(halfW, 2, termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault, fmt.Sprintf("Loaded Plugins (%d)", len(loadedPlugins)))

	for x := 0; x < w; x++ {
		termbox.SetCell(x, 3, tuibase.HorizontalLineRune(), termbox.ColorDarkGray, termbox.ColorDefault)
	}
	for y := 3; y < h; y++ {
		termbox.SetCell(halfW-2, y, tuibase.VerticalLineRune(), termbox.ColorDarkGray, termbox.ColorDefault)
	}

	row := 4
	if len(configuredPlugins) == 0 {
		tuibase.DrawText(0, row, termbox.ColorDarkGray, termbox.ColorDefault, "(empty)")
	} else {
		for _, pluginName := range configuredPlugins {
			if row >= h-4 {
				break
			}
			tuibase.DrawText(0, row, termbox.ColorWhite, termbox.ColorDefault, pluginName)
			row++
		}
	}

	row = 4
	if len(loadedPlugins) == 0 {
		tuibase.DrawText(halfW, row, termbox.ColorDarkGray, termbox.ColorDefault, "(no plugin registered)")
	} else {
		for _, pluginName := range loadedPlugins {
			if row >= h-4 {
				break
			}
			tuibase.DrawText(halfW, row, termbox.ColorWhite, termbox.ColorDefault, pluginName)
			row++
		}
	}

	commandRow := h - 3
	if commandRow >= 4 {
		for x := 0; x < w; x++ {
			termbox.SetCell(x, commandRow-1, tuibase.HorizontalLineRune(), termbox.ColorDarkGray, termbox.ColorDefault)
		}
		commandHint := "Command"
		if pluginModeState == pluginsBrowseState {
			commandHint = "Command (: enter)"
		} else {
			commandHint = "Command (input)"
		}
		tuibase.DrawText(0, commandRow, termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault, commandHint)
		tuibase.DrawText(10, commandRow, termbox.ColorWhite, termbox.ColorDefault, pluginModeCommand)
	}

	outputStartRow := h - 2
	if len(pluginModeOutput) > 0 && outputStartRow >= 4 {
		outputLines := pluginModeOutput
		if len(outputLines) > 2 {
			outputLines = outputLines[len(outputLines)-2:]
		}
		for idx, line := range outputLines {
			row := outputStartRow + idx
			if row >= h {
				break
			}
			tuibase.DrawText(0, row, termbox.ColorGreen, termbox.ColorDefault, line)
		}
	}

	termbox.Flush()
}

func executePluginsModeCommand() {
	command := strings.TrimSpace(pluginModeCommand)
	pluginModeOutput = pluginModeOutput[:0]
	pluginModeState = pluginsBrowseState

	switch command {
	case "", ":":
		return
	case ":out-alias":
		aliasDocs := pluginsX.GetPluginsBus().GetTaskAliasDocsSnapshot()
		if len(aliasDocs) == 0 {
			pluginModeOutput = append(pluginModeOutput, "当前没有已注册的别名规则")
			return
		}
		for _, aliasDoc := range aliasDocs {
			pluginModeOutput = append(pluginModeOutput, fmt.Sprintf(
				"%s => %s (%s)",
				aliasDoc.Alias,
				aliasDoc.TaskType,
				aliasDoc.Desc,
			))
		}
	default:
		pluginModeOutput = append(pluginModeOutput, "未知命令: "+command)
	}
}

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
	tuibase.DrawText(startX, startY, colorBorder, bg, string(tuibase.BoxTopLeftRune()))
	for i := 1; i < panelW-1; i++ {
		tuibase.DrawText(startX+i, startY, colorBorder, bg, string(tuibase.BoxHorizontalRune()))
	}
	tuibase.DrawText(startX+panelW-1, startY, colorBorder, bg, string(tuibase.BoxTopRightRune()))

	// ── 标题 ──
	title := "Misaka DB —— Main Menu"
	titleRow := startY + 1
	titleW := runewidth.StringWidth(title)
	leftPad := (panelW - 2 - titleW) / 2

	tuibase.DrawText(startX, titleRow, colorBorder, bg, string(tuibase.BoxVerticalRune()))
	tuibase.FillSpaces(startX+1, startX+1+leftPad, titleRow, colorTitle, bg)
	afterTitle := tuibase.DrawText(startX+1+leftPad, titleRow, colorTitle, bg, title)
	tuibase.FillSpaces(afterTitle, startX+panelW-1, titleRow, colorTitle, bg)
	tuibase.DrawText(startX+panelW-1, titleRow, colorBorder, bg, string(tuibase.BoxVerticalRune()))

	// ── 分隔线 ──
	sepRow := startY + 2
	tuibase.DrawText(startX, sepRow, colorBorder, bg, string(tuibase.BoxLeftDividerRune()))
	for i := 1; i < panelW-1; i++ {
		tuibase.DrawText(startX+i, sepRow, colorSep, bg, string(tuibase.HorizontalLineRune()))
	}
	tuibase.DrawText(startX+panelW-1, sepRow, colorBorder, bg, string(tuibase.BoxRightDividerRune()))

	// ── 快捷键条目 ──
	for idx, item := range menuItems {
		row := startY + 3 + idx
		tuibase.DrawText(startX, row, colorBorder, bg, string(tuibase.BoxVerticalRune()))

		col := startX + 2
		col = tuibase.DrawText(col, row, colorSep, bg, "[")
		col = tuibase.DrawText(col, row, colorKey, bg, item.key)
		col = tuibase.DrawText(col, row, colorSep, bg, "]")
		col = tuibase.DrawText(col, row, colorDesc, bg, "  ")
		afterDesc := tuibase.DrawText(col, row, colorDesc, bg, item.desc)
		tuibase.FillSpaces(afterDesc, startX+panelW-1, row, colorDesc, bg)
		tuibase.DrawText(startX+panelW-1, row, colorBorder, bg, string(tuibase.BoxVerticalRune()))
	}

	// ── 底部边框 ──
	bottomRow := startY + 3 + len(menuItems)
	tuibase.DrawText(startX, bottomRow, colorBorder, bg, string(tuibase.BoxBottomLeftRune()))
	for i := 1; i < panelW-1; i++ {
		tuibase.DrawText(startX+i, bottomRow, colorBorder, bg, string(tuibase.BoxHorizontalRune()))
	}
	tuibase.DrawText(startX+panelW-1, bottomRow, colorBorder, bg, string(tuibase.BoxBottomRightRune()))

	termbox.Flush()
}

func modCliMode() {
	lockshow.StopLockMode()
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
				lockshow.RenderLockMode()
			case tuimode.PluginsMode:
				renderPluginsMode()
			}

		case termbox.EventResize:
			switch tuimode.GetTuiMode() {
			case tuimode.MenuMode:
				Menu()
			case tuimode.CLIMode:
				renderLogMode()
			case tuimode.LockMode:
				lockshow.RenderLockMode()
			case tuimode.PluginsMode:
				renderPluginsMode()
			}

		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeySpace:
				if tuimode.GetTuiMode() == tuimode.PluginsMode && pluginModeState == pluginsCommandInputState {
					pluginModeCommand += " "
					renderPluginsMode()
					continue
				}
				lockshow.StopLockMode()
				// 按 SPACE：重新初始化 termbox 以彻底清空屏幕
				termbox.Close()
				if err := termbox.Init(); err != nil {
					return
				}
				tuimode.SetTuiMode(tuimode.MenuMode)
				Menu()
			case termbox.KeyEnter:
				if tuimode.GetTuiMode() == tuimode.PluginsMode {
					executePluginsModeCommand()
					pluginModeCommand = ""
					renderPluginsMode()
				}
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if tuimode.GetTuiMode() == tuimode.PluginsMode && len(pluginModeCommand) > 0 {
					pluginModeCommand = pluginModeCommand[:len(pluginModeCommand)-1]
					renderPluginsMode()
				}
			case termbox.KeyEsc:
				if tuimode.GetTuiMode() == tuimode.PluginsMode && pluginModeState == pluginsCommandInputState {
					pluginModeState = pluginsBrowseState
					pluginModeCommand = ""
					renderPluginsMode()
				}
			}

			if tuimode.GetTuiMode() == tuimode.PluginsMode && pluginModeState == pluginsCommandInputState {
				if ev.Ch != 0 {
					pluginModeCommand += string(ev.Ch)
					renderPluginsMode()
				}
				continue
			}

			switch ev.Ch {
			case 'n', 'N':
				modCliMode()
			case 'l', 'L':
				lockshow.StopLockMode()
				termbox.Close()
				if err := termbox.Init(); err != nil {
					return
				}
				lockshow.ModLockMode()
			case 'p', 'P':
				lockshow.StopLockMode()
				termbox.Close()
				if err := termbox.Init(); err != nil {
					return
				}
				tuimode.SetTuiMode(tuimode.PluginsMode)
				pluginModeCommand = ""
				pluginModeOutput = pluginModeOutput[:0]
				pluginModeState = pluginsBrowseState
				renderPluginsMode()

			case 'q', 'Q':
				lockshow.StopLockMode()
				// 退出：关闭 termbox 后恢复正常终端
				clilog.RegisterTUIRenderer(nil)
				termbox.Close()
				if config.GetGlobalServiceConfigure().Debug {
					os.Exit(0)
				}
			case ':':
				if tuimode.GetTuiMode() == tuimode.PluginsMode {
					pluginModeState = pluginsCommandInputState
					pluginModeCommand = ":"
					renderPluginsMode()
				}
			default:
				continue
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
