package interactive

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

func doListen(w *Win) {
	s := w.handler
	for {
		ev := s.PollEvent()

		switch event := ev.(type) {
		case *tcell.EventKey:
			// 特殊键特殊处理
			// 回车
			switch event.Key() {
			case tcell.KeyEnter:
				if w.blockedNow {
					continue
				}

				// 清空最后一行
				for i, j := 0, w.promptWidth+1; i < w.curwidth; i, j = i+1, j+1 {
					s.SetContent(j, w.curmaxY, ' ', nil, tcell.StyleDefault)
				}
				s.ShowCursor(w.promptWidth+1, w.curmaxY)
				s.Show()
				stringCmd := string(w.input)
				go func() {
					w.cmdC <- stringCmd
				}()
				w.input = nil
				w.curwidth = 0
				if w.blockInputAfterEnter {
					w.blockedNow = true
				}
			case tcell.KeyBackspace2:
				if w.blockedNow {
					continue
				}

				if w.curwidth == 0 {
					break
				}

				// 清空最后一行
				for i, j := 0, w.promptWidth+1; i < w.curwidth; i, j = i+1, j+1 {
					s.SetContent(j, w.curmaxY, ' ', nil, tcell.StyleDefault)
				}

				// 更新数据结构
				w.curwidth -= runewidth.RuneWidth(w.input[len(w.input)-1])
				w.input = w.input[0 : len(w.input)-1]

				// 现在写已有的消息
				offset := w.promptWidth + 1
				for i := 0; i < len(w.input); i++ {
					s.SetContent(offset, w.curmaxY, w.input[i], nil, tcell.StyleDefault)
					offset += runewidth.RuneWidth(w.input[i])
				}
				s.ShowCursor(offset, w.curmaxY)
				s.Show()
			case tcell.KeyUp:
				if w.trace {
					if w.eventMask&EventMaskKeyUpWhenTrace == EventMaskKeyUpWhenTrace {
						go func() {
							w.specialEventC <- &EventTypeUpWhenTrace{When: time.Now()}
						}()
					}
					continue
				}
				if w.loff == 0 {
					if w.eventMask&EventMaskTryToMoveUpper == EventMaskTryToMoveUpper {
						go func() {
							w.specialEventC <- &EventTryToGetUpper{When: time.Now()}
						}()
					}
					continue
				}

				if w.eventMask&EventMaskKeyUp == EventMaskKeyUp {
					go func() {
						w.specialEventC <- &EventMoveUp{When: time.Now()}
					}()
				}
				w.loff -= 1
				reDraw(w, false)
			case tcell.KeyDown:
				maxloff, _ := getMaxLoffAndOutputN(w.curmaxY, len(w.lines))
				if w.trace {
					if w.eventMask&EventMaskKeyDownWhenTrace == EventMaskKeyDownWhenTrace {
						go func() {
							w.specialEventC <- &EventTypeDownWhenTrace{When: time.Now()}
						}()
					}
					continue
				}
				if w.loff == maxloff {
					if w.eventMask&EventMaskTryToMoveLower == EventMaskTryToMoveLower {
						go func() {
							w.specialEventC <- &EventTryToGetLower{When: time.Now()}
						}()
					}
					continue
				}

				if w.eventMask&EventMaskKeyDown == EventMaskKeyDown {
					go func() {
						w.specialEventC <- &EventMoveDown{When: time.Now()}
					}()
				}
				w.loff += 1
				reDraw(w, false)
			case tcell.KeyRight:
				if w.curmaxX+1 > maxwidthfrom(w.lines, w.coff+1) {
					continue
				}
				w.coff++
				reDraw(w, false)
			case tcell.KeyLeft:
				if w.coff == 0 {
					continue
				}
				w.coff--
				reDraw(w, false)
			case tcell.KeyCtrlSpace:
				if w.eventMask&EventMaskKeyCtrlSpace == EventMaskKeyCtrlSpace {
					w.specialEventC <- &EventKeyCtrlSpace{When: time.Now()}
				}
			case tcell.KeyCtrlA:
				if w.eventMask&EventMaskKeyCtrlA == EventMaskKeyCtrlA {
					w.specialEventC <- &EventKeyCtrlA{When: time.Now()}
				}
			case tcell.KeyCtrlB:
				if w.eventMask&EventMaskKeyCtrlB == EventMaskKeyCtrlB {
					w.specialEventC <- &EventKeyCtrlB{When: time.Now()}
				}
			case tcell.KeyCtrlC:
				if w.eventMask&EventMaskKeyCtrlC == EventMaskKeyCtrlC {
					w.specialEventC <- &EventKeyCtrlC{When: time.Now()}
				}
			case tcell.KeyCtrlD:
				if w.eventMask&EventMaskKeyCtrlD == EventMaskKeyCtrlD {
					w.specialEventC <- &EventKeyCtrlD{When: time.Now()}
				}
			case tcell.KeyCtrlE:
				if w.eventMask&EventMaskKeyCtrlE == EventMaskKeyCtrlE {
					w.specialEventC <- &EventKeyCtrlE{When: time.Now()}
				}
			case tcell.KeyCtrlF:
				if w.eventMask&EventMaskKeyCtrlF == EventMaskKeyCtrlF {
					w.specialEventC <- &EventKeyCtrlF{When: time.Now()}
				}
			case tcell.KeyCtrlG:
				if w.eventMask&EventMaskKeyCtrlG == EventMaskKeyCtrlG {
					w.specialEventC <- &EventKeyCtrlG{When: time.Now()}
				}
			case tcell.KeyCtrlH:
				if w.eventMask&EvnetMaskKeyCtrlH == EvnetMaskKeyCtrlH {
					w.specialEventC <- &EventKeyCtrlH{When: time.Now()}
				}
			case tcell.KeyCtrlI:
				if w.eventMask&EventMaskKeyCtrlI == EventMaskKeyCtrlI {
					w.specialEventC <- &EventKeyCtrlI{When: time.Now()}
				}
			case tcell.KeyCtrlJ:
				if w.eventMask&EventMaskKeyCtrlJ == EventMaskKeyCtrlJ {
					w.specialEventC <- &EventKeyCtrlJ{When: time.Now()}
				}
			case tcell.KeyCtrlK:
				if w.eventMask&EventMaskKeyCtrlK == EventMaskKeyCtrlK {
					w.specialEventC <- &EventKeyCtrlK{When: time.Now()}
				}
			case tcell.KeyCtrlL:
				if w.eventMask&EventMaskKeyCtrlL == EventMaskKeyCtrlL {
					w.specialEventC <- &EventKeyCtrlL{When: time.Now()}
				}
			case tcell.KeyCtrlN:
				if w.eventMask&EventMaskKeyCtrlN == EventMaskKeyCtrlN {
					w.specialEventC <- &EventKeyCtrlN{When: time.Now()}
				}
			case tcell.KeyCtrlO:
				if w.eventMask&EventMaskKeyCtrlO == EventMaskKeyCtrlO {
					w.specialEventC <- &EventKeyCtrlO{When: time.Now()}
				}
			case tcell.KeyCtrlP:
				if w.eventMask&EventMaskKeyCtrlP == EventMaskKeyCtrlP {
					w.specialEventC <- &EventKeyCtrlP{When: time.Now()}
				}
			case tcell.KeyCtrlQ:
				if w.eventMask&EventMaskKeyCtrlQ == EventMaskKeyCtrlQ {
					w.specialEventC <- &EventKeyCtrlQ{When: time.Now()}
				}
			case tcell.KeyCtrlR:
				if w.eventMask&EventMaskKeyCtrlR == EventMaskKeyCtrlR {
					w.specialEventC <- &EventKeyCtrlR{When: time.Now()}
				}
			case tcell.KeyCtrlS:
				if w.eventMask&EventMaskKeyCtrlS == EventMaskKeyCtrlS {
					w.specialEventC <- &EventKeyCtrlS{When: time.Now()}
				}
			case tcell.KeyCtrlT:
				if w.eventMask&EventMaskKeyCtrlT == EventMaskKeyCtrlT {
					w.specialEventC <- &EventKeyCtrlT{When: time.Now()}
				}
			case tcell.KeyCtrlU:
				if w.eventMask&EventMaskKeyCtrlU == EventMaskKeyCtrlU {
					w.specialEventC <- &EventKeyCtrlU{When: time.Now()}
				}
			case tcell.KeyCtrlV:
				if w.eventMask&EventMaskKeyCtrlV == EventMaskKeyCtrlV {
					w.specialEventC <- &EventKeyCtrlV{When: time.Now()}
				}
			case tcell.KeyCtrlW:
				if w.eventMask&EventMaskKeyCtrlW == EventMaskKeyCtrlW {
					w.specialEventC <- &EventKeyCtrlW{When: time.Now()}
				}
			case tcell.KeyCtrlX:
				if w.eventMask&EventMaskKeyCtrlX == EventMaskKeyCtrlX {
					w.specialEventC <- &EventKeyCtrlX{When: time.Now()}
				}
			case tcell.KeyCtrlY:
				if w.eventMask&EventMaskKeyCtrlY == EventMaskKeyCtrlY {
					w.specialEventC <- &EventKeyCtrlY{When: time.Now()}
				}
			case tcell.KeyCtrlZ:
				if w.eventMask&EventMaskKeyCtrlZ == EventMaskKeyCtrlZ {
					w.specialEventC <- &EventKeyCtrlZ{When: time.Now()}
				}
			}

			// 忽略非普通rune字符
			if event.Key() != tcell.KeyRune {
				continue
			}

			if w.blockedNow {
				continue
			}
			c := event.Rune()
			cWidth := runewidth.RuneWidth(c)
			if w.curmaxX+1-(w.promptWidth+1)-w.curwidth > cWidth {
				s.SetContent(w.curwidth+w.promptWidth+1, w.curmaxY, c, nil, tcell.StyleDefault)
				w.curwidth += cWidth
				s.ShowCursor(w.curwidth+w.promptWidth+1, w.curmaxY)
				s.Show()
				w.input = append(w.input, c)
			} else {
				s.Beep()
			}
		case *tcell.EventResize:
			x, y := s.Size()
			w.curmaxX, w.curmaxY = x-1, y-1
			reDraw(w, true)
			if w.eventMask&EventMaskWindowResize == EventMaskWindowResize {
				w.specialEventC <- &EventWindowResize{
					Height: y,
					Width:  x,
					When:   time.Now(),
				}
			}
		case *stopEvent:
			w.isStopped = true
			w.handler.Fini()
			w.waitStopChan <- struct{}{}
			return
		case *setBlockInputChangeEvent:
			w.blockedNow = event.data
			w.input = nil
			w.curwidth = 0
			// 开始先画一个命令提示符出来
			for i := 0; i < w.curmaxX; i++ {
				w.handler.SetContent(i, w.curmaxY, ' ', nil, tcell.StyleDefault)
			}
			s.SetContent(0, w.curmaxY, w.prompt, nil, styleAttr2TcellStyle(&w.promptStyle))
			s.ShowCursor(w.promptWidth+1, w.curmaxY)
			s.Show()
		case *clearEvent:
			w.lines = nil
			w.coff = 0
			w.loff = 0
			reDraw(w, false)
		case *gotoBottomEvent:
			w.trace = false
			maxloff, _ := getMaxLoffAndOutputN(w.curmaxY, len(w.lines))
			w.loff = maxloff
			reDraw(w, false)
		case *gotoTopEvent:
			w.trace = false
			w.loff = 0
			reDraw(w, false)
		case *gotoLeftEvent:
			if w.coff != 0 {
				w.coff = 0
				reDraw(w, false)
			}
		case *setTraceEvent:
			w.trace = event.data
			reDraw(w, false)
		case *setBlockInputAfterEnterEvent:
			w.blockInputAfterEnter = event.data
		case *gotoLineEvent:
			w.trace = false
			if event.data-1 == w.loff {
				continue
			}
			maxloff, _ := getMaxLoffAndOutputN(w.curmaxY, len(w.lines))
			if event.data <= 0 {
				w.loff = 0
			} else if event.data >= maxloff+1 {
				w.loff = maxloff
			} else {
				w.loff = event.data - 1
			}
			reDraw(w, false)
		case *gotoNextLineEvent:
			w.trace = false
			maxloff, _ := getMaxLoffAndOutputN(w.curmaxY, len(w.lines))
			if w.loff == maxloff {
				continue
			}
			w.loff++
			reDraw(w, false)
		case *gotoPreviousLineEvent:
			w.trace = false
			if w.loff == 0 {
				continue
			}
			w.loff--
			reDraw(w, false)
		case *sendLineFrontWithColorEvent:
			newLines := make([]([]interface{}), len(w.lines)+1, (len(w.lines)+1)*2)
			newLines[0] = event.data
			for i := 1; i <= len(w.lines); i++ {
				newLines[i] = w.lines[i-1]
			}
			w.lines = newLines

			if w.trace {
				continue
			}

			maxloff, _ := getMaxLoffAndOutputN(w.curmaxY, len(w.lines))
			if w.loff == maxloff {
				reDraw(w, false)
				continue
			}

			// TODO 是否合适?
			w.loff++
			reDraw(w, false)
		case *sendLineBackWithColorEvent:
			w.lines = append(w.lines, event.data)
			reDraw(w, false)
		case *popBackLineEvent:
			if len(w.lines) == 0 {
				continue
			}
			w.lines = w.lines[:len(w.lines)-1]
			maxloff, _ := getMaxLoffAndOutputN(w.curmaxY, len(w.lines))
			if w.loff > maxloff {
				w.loff = maxloff
			}
			reDraw(w, false)
		case *popFrontLineEvent:
			if len(w.lines) == 0 {
				continue
			}
			w.lines = w.lines[1:]
			if w.trace {
				continue
			}
			if w.loff >= 1 {
				w.loff--
			}
			reDraw(w, false)
		case *setPromptEvent:
			if event.dataRune != nil {
				w.prompt = *event.dataRune
			}
			if event.dataStyle != nil {
				w.promptStyle = *event.dataStyle
			}

			// 清空最后一行, 因为可能会发生命令提示符的宽度改变的情况
			// 这里需要一次性更新最后一行
			for i, j := 0, w.promptWidth+1; i < w.curwidth; i, j = i+1, j+1 {
				s.SetContent(j, w.curmaxY, ' ', nil, tcell.StyleDefault)
			}

			// 画新的命令提示符
			s.SetContent(0, w.curmaxY, w.prompt, nil, styleAttr2TcellStyle(&w.promptStyle))
			w.promptWidth = runewidth.RuneWidth(w.prompt)

			// 现在写已有的消息
			offset := w.promptWidth + 1
			for i := 0; i < len(w.input); i++ {
				s.SetContent(offset, w.curmaxY, w.input[i], nil, tcell.StyleDefault)
				offset += runewidth.RuneWidth(w.input[i])
			}
			s.ShowCursor(offset, w.curmaxY)
			s.Show()
		case *getWindowSizeEvent:
			resp := new(getWindowSizeEventResp)
			resp.height = w.curmaxY + 1
			resp.width = w.curmaxX + 1
			event.resp <- resp
		}
	}
}
