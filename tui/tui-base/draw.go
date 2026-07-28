package tuibase

import (
	"misakadb/shares"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// DrawText 逐字符绘制（全角字符占2列），返回结束列号
func DrawText(x, y int, fg, bg termbox.Attribute, text string) int {
	col := x
	for _, ch := range text {
		width := runewidth.RuneWidth(ch)
		if width <= 0 {
			width = 1
		}

		if shares.IsWindows() {
			termbox.SetCell(col, y, ch, fg, bg)
			col += width
			continue
		}

		termbox.SetCell(col, y, ch, fg, bg)
		// 宽字符需要占满后续单元格，避免 Windows 下布局被残留内容冲坏。
		for i := 1; i < width; i++ {
			termbox.SetCell(col+i, y, ' ', fg, bg)
		}
		col += width
	}
	return col
}

// FillSpaces 用空格填充 [x, endX) 区间
func FillSpaces(x, endX, y int, fg, bg termbox.Attribute) {
	for c := x; c < endX; c++ {
		termbox.SetCell(c, y, ' ', fg, bg)
	}
}

func HorizontalLineRune() rune {
	if shares.IsWindows() {
		return '-'
	}
	return '─'
}

func VerticalLineRune() rune {
	if shares.IsWindows() {
		return '|'
	}
	return '│'
}

func BoxTopLeftRune() rune {
	if shares.IsWindows() {
		return '+'
	}
	return '╔'
}

func BoxTopRightRune() rune {
	if shares.IsWindows() {
		return '+'
	}
	return '╗'
}

func BoxBottomLeftRune() rune {
	if shares.IsWindows() {
		return '+'
	}
	return '╚'
}

func BoxBottomRightRune() rune {
	if shares.IsWindows() {
		return '+'
	}
	return '╝'
}

func BoxHorizontalRune() rune {
	if shares.IsWindows() {
		return '-'
	}
	return '═'
}

func BoxLeftDividerRune() rune {
	if shares.IsWindows() {
		return '+'
	}
	return '╠'
}

func BoxRightDividerRune() rune {
	if shares.IsWindows() {
		return '+'
	}
	return '╣'
}

func BoxVerticalRune() rune {
	if shares.IsWindows() {
		return '|'
	}
	return '║'
}
