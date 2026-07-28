package lockshow

import (
	"fmt"
	"misakadb/lock/global_lock"
	tuibase "misakadb/tui/tui-base"
	tuimode "misakadb/tui/tui-mode"
	"sort"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

var (
	lockModeMu     sync.Mutex
	lockModeStopCh chan struct{}
)

func RenderLockMode() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	title := "Lock Mode - Press SPACE to show menu | Q exit"
	tuibase.DrawText(0, 0, termbox.ColorCyan|termbox.AttrBold, termbox.ColorDefault, title)

	for x := 0; x < w; x++ {
		termbox.SetCell(x, 1, tuibase.HorizontalLineRune(), termbox.ColorDarkGray, termbox.ColorDefault)
	}

	line := fmt.Sprintf(
		"GlobalLockPool Viewer %s Refresh [%s]",
		time.Now().Format("2006-01-02 15:04:05"),
		tuimode.GetTuiMode(),
	)
	tuibase.DrawText(0, 2, termbox.ColorWhite, termbox.ColorDefault, line)

	for x := 0; x < w; x++ {
		termbox.SetCell(x, 3, tuibase.HorizontalLineRune(), termbox.ColorDarkGray, termbox.ColorDefault)
	}

	halfW := w / 2

	// 标题栏
	tuibase.DrawText(0, 4, termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault, "Young Pool")
	tuibase.DrawText(halfW, 4, termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault, "Old Pool")

	// 中间竖线
	for y := 4; y < h; y++ {
		termbox.SetCell(halfW-2, y, tuibase.VerticalLineRune(), termbox.ColorDarkGray, termbox.ColorDefault)
	}

	pools := global_lock.GetGlobalLockPool()
	youngMap := pools.GetYoungPoolSnapshot()
	oldMap := pools.GetOldPoolSnapshot()

	type lockItem struct {
		name string
		ref  int32
	}

	youngList := make([]lockItem, 0, len(youngMap))
	for k, v := range youngMap {
		youngList = append(youngList, lockItem{name: k, ref: v})
	}
	sort.Slice(youngList, func(i, j int) bool {
		return youngList[i].name < youngList[j].name
	})

	oldList := make([]lockItem, 0, len(oldMap))
	for k, v := range oldMap {
		oldList = append(oldList, lockItem{name: k, ref: v})
	}
	sort.Slice(oldList, func(i, j int) bool {
		return oldList[i].name < oldList[j].name
	})

	// 绘制 Young
	y := 5
	for _, item := range youngList {
		if y >= h {
			break
		}
		text := fmt.Sprintf("Lock: %s (ref: %d)", item.name, item.ref)
		if item.ref == 0 {
			text += " [即将降级]"
			tuibase.DrawText(0, y, termbox.ColorRed, termbox.ColorDefault, text)
		} else {
			tuibase.DrawText(0, y, termbox.ColorWhite, termbox.ColorDefault, text)
		}
		y++
	}

	// 绘制 Old
	y = 5
	for _, item := range oldList {
		if y >= h {
			break
		}
		text := fmt.Sprintf("Lock: %s (ref: %d)", item.name, item.ref)
		if item.ref == 0 {
			text += " [即将淘汰]"
			tuibase.DrawText(halfW, y, termbox.ColorRed, termbox.ColorDefault, text)
		} else {
			tuibase.DrawText(halfW, y, termbox.ColorWhite, termbox.ColorDefault, text)
		}
		y++
	}

	termbox.Flush()
}

func ModLockMode() {
	tuimode.SetTuiMode(tuimode.LockMode) // 切换模式

	lockModeMu.Lock()
	if lockModeStopCh != nil {
		close(lockModeStopCh)
	}
	stopCh := make(chan struct{})
	lockModeStopCh = stopCh
	lockModeMu.Unlock()

	RenderLockMode()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				if tuimode.GetTuiMode() != tuimode.LockMode {
					return
				}
				RenderLockMode()
			}
		}
	}()
}

func StopLockMode() {
	lockModeMu.Lock()
	defer lockModeMu.Unlock()

	if lockModeStopCh != nil {
		close(lockModeStopCh)
		lockModeStopCh = nil
	}
}
