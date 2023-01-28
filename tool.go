package interactive

import (
	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

func reDraw(w *Win, resize bool) {
	s := w.handler
	s.Clear()
	maxLoff, outputLinesN := getMaxLoffAndOutputN(w.curmaxY, len(w.lines))
	if w.trace || w.loff > maxLoff {
		w.loff = maxLoff
	}

	// 开始输出界面
	for i, j := 0, w.loff; i < outputLinesN; i, j = i+1, j+1 {
		curwidth := 0
		curLine := w.lines[j]
		curLineSize := len(w.lines[j])
		for beg := w.coff; beg < curLineSize; beg++ {
			thisChar := curLine[beg]
			thisCharWidth := runewidth.RuneWidth(rune(thisChar))
			if w.curmaxX+1-curwidth >= thisCharWidth {
				s.SetContent(curwidth, i, thisChar, nil, tcell.StyleDefault)
				curwidth += thisCharWidth
			} else {
				break
			}
		}
	}

	s.SetContent(0, w.curmaxY, w.prompt, nil, tcell.StyleDefault)
	s.SetContent(w.promptWidth, w.curmaxY, ' ', nil, tcell.StyleDefault)
	offset := w.inputStart
	if resize {
		w.input = nil
		w.curwidth = 0
	}
	for i := 0; i < len(w.input); i++ {
		s.SetContent(offset, w.curmaxY, w.input[i], nil, tcell.StyleDefault)
		offset += runewidth.RuneWidth(w.input[i])
	}
	s.ShowCursor(offset, w.curmaxY)
	s.Show()
}

func maxwidthfrom(w *Win) int {
	n := w.coff + 1
	maxwidth := 0
	nLines := len(w.lines)
	for i := 0; i < nLines; i++ {
		thisLine := w.lines[i]
		thisWidth := 0
		for beg := n; beg < len(thisLine); beg++ {
			thisWidth += runewidth.RuneWidth(rune(thisLine[beg]))
		}
		if maxwidth < thisWidth {
			maxwidth = thisWidth
		}
	}
	return maxwidth
}

func getMaxLoffAndOutputN(curY, cntLines int) (x, y int) {
	if cntLines < curY {
		x = 0
		y = cntLines
	} else {
		x = cntLines - curY
		y = curY
	}
	return x, y
}
