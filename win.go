package interactive

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

// 窗口对象,
type Win struct {
	handler tcell.Screen

	// 每行的数据, 被锁保护
	lines []([]rune)

	// 输入行的数据, 不被锁保护
	input []rune

	// 是否追踪最新输出, 被锁保护
	trace bool

	// 命令提示符, 不修改, 不保护
	prompt rune
	// 命令提示符的宽度, 不修改, 不保护
	promptWidth int

	// line offset, 在追踪最新输出时, 这个没意义, 不保护
	loff int
	// column offset, 不保护
	coff int

	// 当前输入的宽度, 不保护
	curwidth int

	// 输入位置的起始位置, 不修改, 不保护
	intputStart int

	// 不保护
	curmaxY int
	curmaxX int

	// 传递命令
	cmdC chan string

	// TODO 是否需要保护?
	isStopped bool

	// 是否在键入回车之后禁止输入?
	blockInputAfterEnter bool

	// 目前是否处于输入阻塞状态
	blockedNow bool
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
		intputStart:          runewidth.RuneWidth(cfg.Prompt) + 1,
		curmaxY:              y - 1,
		curmaxX:              x - 1,
		cmdC:                 make(chan string),
		isStopped:            false,
		blockInputAfterEnter: cfg.BlockInputAfterEnter,
		blockedNow:           cfg.BlockInputAfterRun,
	}

	s.SetStyle(tcell.StyleDefault)
	s.Clear()
	s.SetContent(0, w.curmaxY, cfg.Prompt, nil, tcell.StyleDefault)
	s.SetContent(w.intputStart-1, w.curmaxY, ' ', nil, tcell.StyleDefault)
	s.ShowCursor(w.intputStart, w.curmaxY)
	s.Show()
	go doListen(w)
	return w
}

// 会将用户输入的信息发送到这个channel, 永远不应该关闭这个channel
func (w *Win) GetCmdChan() chan string {
	return w.cmdC
}

func doListen(w *Win) {
	s := w.handler
	defer func() {
		s.Fini()
	}()
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
				for i, j := 0, w.intputStart; i < w.curwidth; i, j = i+1, j+1 {
					s.SetContent(j, w.curmaxY, ' ', nil, tcell.StyleDefault)
				}
				s.ShowCursor(w.intputStart, w.curmaxY)
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
				for i, j := 0, w.intputStart; i < w.curwidth; i, j = i+1, j+1 {
					s.SetContent(j, w.curmaxY, ' ', nil, tcell.StyleDefault)
				}

				// 更新数据结构
				w.curwidth -= runewidth.RuneWidth(w.input[len(w.input)-1])
				w.input = w.input[0 : len(w.input)-1]

				// 现在写已有的消息
				offset := w.intputStart
				for i := 0; i < len(w.input); i++ {
					s.SetContent(offset, w.curmaxY, w.input[i], nil, tcell.StyleDefault)
					offset += runewidth.RuneWidth(w.input[i])
				}
				s.ShowCursor(offset, w.curmaxY)
				s.Show()
			case tcell.KeyUp:
				if w.loff == 0 {
					continue
				}
				w.loff -= 1
				reDraw(w, false)
			case tcell.KeyDown:
				maxloff, _ := getMaxLoffAndOutputN(w.curmaxY, len(w.lines))
				if w.loff == maxloff {
					continue
				}
				w.loff += 1
				reDraw(w, false)
			case tcell.KeyRight:
				if w.curmaxX+1 > maxwidthfrom(w) {
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
			if w.curmaxX+1-w.intputStart-w.curwidth > cWidth {
				s.SetContent(w.curwidth+w.intputStart, w.curmaxY, c, nil, tcell.StyleDefault)
				w.curwidth += cWidth
				s.ShowCursor(w.curwidth+w.intputStart, w.curmaxY)
				s.Show()
				w.input = append(w.input, c)
			} else {
				s.Beep()
			}
		case *sendEvent:
			w.lines = append(w.lines, event.data)
			reDraw(w, false)
		case *tcell.EventResize:
			x, y := s.Size()
			w.curmaxX, w.curmaxY = x-1, y-1
			reDraw(w, true)
		case *endEvent:
			w.isStopped = true
			return
		case *blockInputChangeEvent:
			w.blockedNow = event.data
			w.input = nil
			w.curwidth = 0
			reDraw(w, false)
		case *clearEvent:
			w.lines = nil
			w.coff = 0
			w.loff = 0
			reDraw(w, false)
		}
	}
}

// 当向已经关闭的Win发送信息时返回error, 一个良好设计的程序不用检查这个error
func (w *Win) SendLine(s string) error {
	if w.isStopped {
		return errors.New("send to a closed window")
	}
	data := make([]rune, 0, utf8.RuneCountInString(s))
	for _, v := range s {
		data = append(data, v)
		if v == '\n' {
			break
		}
	}
	w.handler.PostEventWait(&sendEvent{when: time.Now(), data: data})

	return nil
}

// TODO 支持带有颜色的输出
// func (w *Win) SendLineWithColor(...interface{}) error {
// 	return nil
// }

// 关闭窗口
func (w *Win) Stop() {
	w.handler.PostEventWait(&endEvent{when: time.Now()})
}

// 追踪最新输出, 此时不允许上下移动, 但允许左右移动
func (w *Win) SetTrace(enable bool) {
	w.trace = enable
}

// 是否禁止输入, 当禁止输入时, 用户输入将被清空
func (w *Win) BlockInput(ifBlock bool) {
	w.handler.PostEventWait(&blockInputChangeEvent{when: time.Now(), data: ifBlock})
}

func (w *Win) Clear() {
	w.handler.PostEventWait(&clearEvent{when: time.Now()})
}
