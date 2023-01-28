# Interactive-Console

### 简介

这是一个命令行交互框架, 用于创建一个交互式终端框架, 功能大纲如下:

- 友好的用户输入, 输入输出分离, 给予用户一个分离的输入行, 避免输入输出混乱的问题
- 友好的命令接收方式, 每个输入回车后会异步发送channel信息
- 完备的控制, 可以随时禁止输入, 并且可以配置输入一行命令后自动关闭输入
- 支持浏览模式和追踪最新消息模式, 并且可以在两种模式之间随时切换
- 输出支持颜色, 支持可选的状态栏, 终端窗口大小改变时自动调节适配
- 线程安全的接口

### TODO

- 支持颜色 [已完成]
- 支持可选的状态栏 [X]

### 获得此库

```bash
go get github.com/markity/Interactive-Console
```

### 快速入门 - 一个简单的交互式程序

```go
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
				attr1 := interactive.GetDefaultSytleAttr()
				attr1.Foreground = interactive.ColorPurple
				attr1.Bold = true
				attr2 := attr1
				attr2.Italic = true
				w.SendLineBackWithColor(attr1, "pong", attr2, "!")
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
				attr := interactive.GetDefaultSytleAttr()
				attr.Foreground = interactive.ColorPink
				w.SendLineBackWithColor(attr, "你已经清空的屏幕")
				w.SetBlockInput(false)
			default:
				attr := interactive.GetDefaultSytleAttr()
				attr.Foreground = interactive.ColorRed
				w.SendLineBackWithColor(attr, fmt.Sprintf("unknown command %s", cmd))
				w.SetBlockInput(false)
			}
		case ev := <-w.GetEventChan():
			// 特别的事件将通过这个管道传输
			switch ev.(type) {
			// 如果当时已经在最后一行了, 用户仍然按下向下键, 那么产生这个事件
			case *interactive.EventTryToGetLower:
				w.SendLineBack("没有更多向下的消息了")
				w.SetTrace(true)
			// 如果当时已经在第一行了, 用户仍然按下向上键, 那么产生这个事件
			case *interactive.EventTryToGetUpper:
				w.SendLineBack("没有更多向上的消息了")
			// 如果当时处在trace状态, 用户却尝试按上下键
			case *interactive.EventTypeUpWhenTrace:
				w.SetTrace(false)
			case *interactive.EventTypeDownWhenTrace:
				w.SendLineBack("现在已经处于追踪模式了")
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
```

---

UPDATE LOG

```
date: 2023.1.28
version: unstable
log:
	实现接口 Win.GotoLeft
	实现接口 GotoLine, GotoNextLine, GotoPreviousLine
	实现接口 Win.GetEventChan
	实现接口 Win.SendLineFront, Win.PopFrontLine, Win.PopBackLine
	更名接口 Win.SendLine -> Win.SendLineBack
	新增特性 可以通过Win.GetEventChan收到注册过的事件
	新增事件 上移事件EventMoveUp 下移事件EventMoveDown
	新增事件 已经处在最顶端但尝试上移的事件EventTryToGetUpper 已经处在最底端但尝试下移的事件EventTryToGetLower
	新增特性 现在可以使用颜色了
	新增接口 SendLineBackWithColor, SendLineFrontWithColor
```

```
date: 2023.1.27
version: unstable
log:
	实现接口 Win.Run, Win.Stop
	实现接口 Win.SendLine, Win.Clear
	实现接口 Win.SetTrace, Win.SetBlockInput, Win.SetBlockInputAfterEnter
	实现接口 Win.GotoTop, Win.GotButtom
```
