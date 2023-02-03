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
		style := tcell.StyleDefault

		offset := 0
		for _, v := range curLine {
			str, ok := v.(string)
			if ok {
				for _, char := range str {
					if offset >= w.coff {
						charWidth := runewidth.RuneWidth(rune(char))
						if w.curmaxX+1-curwidth >= charWidth {
							s.SetContent(curwidth, i, char, nil, style)
							curwidth += charWidth
						} else {
							goto out
						}
					}
					offset++
				}
			} else {
				style = v.(tcell.Style)
			}
		}
	out:
	}

	s.SetContent(0, w.curmaxY, w.prompt, nil, styleAttr2TcellStyle(&w.promptStyle))
	s.SetContent(w.promptWidth, w.curmaxY, ' ', nil, tcell.StyleDefault)
	ioffset := w.promptWidth + 1
	if resize {
		w.input = nil
		w.curwidth = 0
	}
	for i := 0; i < len(w.input); i++ {
		s.SetContent(ioffset, w.curmaxY, w.input[i], nil, tcell.StyleDefault)
		ioffset += runewidth.RuneWidth(w.input[i])
	}
	s.ShowCursor(ioffset, w.curmaxY)

	s.Show()

}

func maxwidthfrom(w []([]interface{}), n int) int {
	maxwidth := 0

	for _, thisLine := range w {
		thisWidth := 0
		offset := 0
		for _, v := range thisLine {
			str, ok := v.(string)
			if ok {
				for _, char := range str {
					if offset >= n {
						thisWidth += runewidth.RuneWidth(char)
					}
					offset++
				}
			}
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
