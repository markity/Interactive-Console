package main

import (
	"fmt"
	"sync"

	interactive "github.com/markity/Interactive-Console"
)

func cmdHandler(w *interactive.Win, wait *sync.WaitGroup) {
	for {
		select {
		case cmd := <-w.GetCmdChan():
			switch cmd {
			case "trace":
				// SetTrace: 追踪最新的消息, 此时用户不允许上下移动浏览
				w.SetTrace(true)

				// 现在允许用户输入
				w.SetBlockInput(false)
			case "untrace":
				// 解除追踪
				w.SetTrace(false)
				w.SetBlockInput(false)
			case "ping":
				// 发送一行信息, 注意不能包含\n字符, 否则只保留\n前面的字符而丢弃后面的
				// 如果要发送多行, 多次调用SendLine, 这是异步安全的, 总是保证先发送的显示在前
				w.SendLine("pong")
				w.SetBlockInput(false)
			case "top":
				// GotoTop, GotoButtom, GotoLine, GotoNextLine, GotoPreviousLine前往指定行, 如果当前正在trace, 将解除trace
				// GotoTop 等价于 GotoLine(1)
				w.GotoLine(1)
				w.SetBlockInput(false)
			case "btm":
				w.GotoButtom()
				w.SetBlockInput(false)
			case "left":
				// GotoLeft前往最左端, 不会改变当前的trace状态
				w.GotoLeft()
				w.SetBlockInput(false)
			case "exit":
				// 停止显示并恢复原来的终端环境
				w.Stop()
				goto out
			case "clear":
				// Clear清除屏幕
				w.Clear()
				w.SendLine("你已经清空了屏幕")
				w.SetBlockInput(false)
			default:
				w.SendLine(fmt.Sprintf("unknown command %s", cmd))
				w.SetBlockInput(false)
			}
		case ev := <-w.GetEventChan():
			// 特别的事件将通过这个管道传输
			switch ev.(type) {
			// 如果当时已经在最后一行了, 用户仍然按下向下键, 那么产生这个事件
			case *interactive.EventTryToGetLower:
				w.SendLine("没有更多向下的消息了")
				w.SetTrace(true)
			// 如果当时已经在第一行了, 用户仍然按下向上键, 那么产生这个事件
			case *interactive.EventTryToGetUpper:
				w.SendLine("没有更多向上的消息了")
			// 如果当时处在trace状态, 用户却尝试按上下键
			case *interactive.EventTypeUpWhenTrace:
				w.SetTrace(false)
			case *interactive.EventTypeDownWhenTrace:
				w.SendLine("现在已经处于追踪模式了")
			}
		}
	}

out:
	wait.Done()
}

func main() {
	cfg := interactive.GetDefaultConfig()
	cfg.BlockInputAfterEnter = true
	cfg.TraceAfterRun = true
	cfg.SpecialEventHandleMask = interactive.EventMaskTryToGetUpper | interactive.EventMaskTryToGetLower | interactive.EventMaskTypeUpWhenTrace | interactive.EventMaskTypeDownWhenTrace
	win := interactive.Run(cfg)
	wait := sync.WaitGroup{}
	wait.Add(1)
	go cmdHandler(win, &wait)
	wait.Wait()
}
