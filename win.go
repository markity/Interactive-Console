package interactive

import (
	"errors"
	"time"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

// 窗口对象, 一个窗口对象可以复用
// 即Run后再Stop, 之后可以再Run
type Win struct {
	// tcell的Screen句柄
	handler tcell.Screen

	// 输出行的数据
	lines []([]interface{})

	// 输入行的数据
	input []rune

	// 是否追踪最新输出
	trace bool

	// 命令提示符
	prompt rune

	// 命令提示符的宽度
	promptWidth int

	// line offset, 在追踪最新输出时, 这个没意义, 不保护
	loff int

	// column offset
	coff int

	// 当前输入的宽度
	curwidth int

	// 当前最大行的偏移, 当前最大列的偏移
	// 前者=行数-1
	// 后者=列数-1
	curmaxY int
	curmaxX int

	// 传递命令的管道
	cmdC chan string

	// 传递事件的管道
	specialEventC chan interface{}

	// 用户需要使用的事件掩码
	eventMask int64

	// 用于判断该窗体是否以及结束执行, 即为调用了Stop
	isStopped bool

	// 是否在键入回车之后禁止输入
	blockInputAfterEnter bool

	// 目前是否处于输入阻塞状态
	blockedNow bool

	// 用来通知关闭Win以及完成
	waitStopChan chan struct{}
}

// 运行窗体
func Run(cfg Config) *Win {
	s, _ := tcell.NewScreen()
	s.Init()
	x, y := s.Size()
	w := &Win{
		handler:              s,
		lines:                nil,
		input:                nil,
		trace:                cfg.TraceAfterRun,
		prompt:               cfg.Prompt,
		promptWidth:          runewidth.RuneWidth(cfg.Prompt),
		loff:                 0,
		coff:                 0,
		curwidth:             0,
		curmaxY:              y - 1,
		curmaxX:              x - 1,
		cmdC:                 make(chan string),
		isStopped:            false,
		blockInputAfterEnter: cfg.BlockInputAfterEnter,
		blockedNow:           cfg.BlockInputAfterRun,
		waitStopChan:         make(chan struct{}),
		specialEventC:        make(chan interface{}),
		eventMask:            cfg.SpecialEventHandleMask,
	}

	s.SetStyle(tcell.StyleDefault)
	s.Clear()
	s.SetContent(0, w.curmaxY, cfg.Prompt, nil, tcell.StyleDefault)
	s.ShowCursor(w.promptWidth+1, w.curmaxY)
	s.Show()
	go doListen(w)
	return w
}

// 会将用户输入的信息发送到这个channel, 永远不应该关闭这个channel
func (w *Win) GetCmdChan() chan string {
	return w.cmdC
}

// 会将特殊事件的信息发送到这个channel, 永远不应该关闭这个channel
func (w *Win) GetEventChan() chan interface{} {
	return w.specialEventC
}

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
			case tcell.KeyBackspace:
				fallthrough
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
			if w.curmaxX+1-w.promptWidth+1-w.curwidth > cWidth {
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
		case *stopEvent:
			w.isStopped = true
			w.handler.Fini()
			w.waitStopChan <- struct{}{}
			return
		case *setBlockInputChangeEvent:
			w.blockedNow = event.data
			w.input = nil
			w.curwidth = 0
			reDraw(w, false)
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
		}
	}
}

// 当向已经关闭的Win发送信息时返回error, 一个良好设计的程序不用检查这个error
func (w *Win) SendLineBack(s string) error {
	return w.SendLineBackWithColor(GetDefaultSytleAttr(), s)
}

func (w *Win) SendLineFront(s string) error {
	return w.SendLineFrontWithColor(GetDefaultSytleAttr(), s)
}

// TODO 支持带有颜色的输出
func (w *Win) SendLineBackWithColor(s ...interface{}) error {
	if w.isStopped {
		return errors.New("send to a closed window")
	}

	for k, v := range s {
		attr, ok1 := v.(StyleAttr)
		_, ok2 := v.(string)
		// 简单的检查, 参数是否规范
		if !ok1 && !ok2 {
			return errors.New("invalid arguments")
		}
		if ok1 {
			style := tcell.Style(0)
			style = style.Background(tcell.Color(attr.Background))
			style = style.Foreground(tcell.Color(attr.Foreground))
			style = style.Blink(attr.Blink)
			style = style.Bold(attr.Bold)
			style = style.Dim(attr.Dim)
			style = style.Italic(attr.Italic)
			style = style.Reverse(attr.Reverse)
			style = style.Underline(attr.Underline)
			s[k] = style
		}
	}

	w.handler.PostEventWait(&sendLineBackWithColorEvent{when: time.Now(), data: s})
	return nil
}

func (w *Win) SendLineFrontWithColor(s ...interface{}) error {
	if w.isStopped {
		return errors.New("send to a closed window")
	}

	for k, v := range s {
		attr, ok1 := v.(StyleAttr)
		_, ok2 := v.(string)
		// 简单的检查, 参数是否规范
		if !ok1 && !ok2 {
			return errors.New("invalid arguments")
		}
		if ok1 {
			style := tcell.Style(0)
			style = style.Background(tcell.Color(attr.Background))
			style = style.Foreground(tcell.Color(attr.Foreground))
			style = style.Blink(attr.Blink)
			style = style.Bold(attr.Bold)
			style = style.Dim(attr.Dim)
			style = style.Italic(attr.Italic)
			style = style.Reverse(attr.Reverse)
			style = style.Underline(attr.Underline)
			s[k] = style
		}
	}

	w.handler.PostEventWait(&sendLineFrontWithColorEvent{when: time.Now(), data: s})
	return nil
}

// 关闭窗口
func (w *Win) Stop() {
	w.handler.PostEventWait(&stopEvent{when: time.Now()})
	<-w.waitStopChan
}

// 追踪最新输出, 此时不允许上下移动, 但允许左右移动
func (w *Win) SetTrace(enable bool) {
	w.handler.PostEventWait(&setTraceEvent{when: time.Now(), data: enable})
}

// 是否禁止输入, 当禁止输入时, 用户输入将被清空
func (w *Win) SetBlockInput(ifBlock bool) {
	w.handler.PostEventWait(&setBlockInputChangeEvent{when: time.Now(), data: ifBlock})
}

// 是否在发送一条命令后禁用输入, 直到手动调用BlockInput(false)才恢复下一条输入
func (w *Win) SetBlockInputAfterEnter(ifBlock bool) {
	w.handler.PostEventWait(&setBlockInputAfterEnterEvent{when: time.Now(), data: ifBlock})
}

// 清空窗体
func (w *Win) Clear() {
	w.handler.PostEventWait(&clearEvent{when: time.Now()})
}

// 移动到最后一行, 将取消trace状态
func (w *Win) GotoButtom() {
	w.handler.PostEventWait(&gotoBottomEvent{when: time.Now()})
}

// 移动到第一行, 将取消trace状态, 等价于GotoLine(1)
func (w *Win) GotoTop() {
	w.handler.PostEventWait(&gotoTopEvent{when: time.Now()})
}

// 移动到最左
func (w *Win) GotoLeft() {
	w.handler.PostEventWait(&gotoLeftEvent{when: time.Now()})
}

// 前往第n行, 将取消trace状态
func (w *Win) GotoLine(n int) {
	w.handler.PostEventWait(&gotoLineEvent{when: time.Now(), data: n})
}

// 前往下一行, 将取消trace状态
func (w *Win) GotoNextLine() {
	w.handler.PostEventWait(&gotoNextLineEvent{when: time.Now()})
}

// 前往上一行, 将取消trace状态
func (w *Win) GotoPreviousLine() {
	w.handler.PostEventWait(&gotoPreviousLineEvent{when: time.Now()})
}

// 删除第一行, 如果没有这一行则什么也不做
func (w *Win) PopFrontLine() {
	w.handler.PostEventWait(&popFrontLineEvent{when: time.Now()})
}

// 删除最后一行, 如果没有这一行则什么也不做
func (w *Win) PopBackLine() {
	w.handler.PostEventWait(&popBackLineEvent{when: time.Now()})
}
