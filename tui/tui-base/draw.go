package tuibase

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// DrawText 逐字符绘制（全角字符占2列），返回结束列号
func DrawText(x, y int, fg, bg termbox.Attribute, text string) int {
	col := x
	for _, ch := range text {
		termbox.SetCell(col, y, ch, fg, bg)
		col += runewidth.RuneWidth(ch)
	}
	return col
}

// FillSpaces 用空格填充 [x, endX) 区间
func FillSpaces(x, endX, y int, fg, bg termbox.Attribute) {
	for c := x; c < endX; c++ {
		termbox.SetCell(c, y, ' ', fg, bg)
	}
}
