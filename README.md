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

- 支持颜色 [X]
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
				w.SetTrace(true)
				w.SetBlockInput(false)
			case "untrace":
				w.SetTrace(false)
				w.SetBlockInput(false)
			case "ping":
				w.SendLine("pong")
				w.SetBlockInput(false)
			case "top":
				w.GotoTop()
				w.SetBlockInput(false)
			case "btm":
				w.GotoButtom()
				w.SetBlockInput(false)
			case "left":
				w.GotoLeft()
				w.SetBlockInput(false)
			case "exit":
				w.Stop()
				goto out
			case "clear":
				w.Clear()
				w.SendLine("你已经清空了屏幕")
				w.SetBlockInput(false)
			default:
				w.SendLine(fmt.Sprintf("unknown command %s", cmd))
				w.SetBlockInput(false)
			}
		case ev := <-w.GetEventChan():
			switch ev.(type) {
			case *interactive.EventTryToGetLower:
				w.SendLine("没有更多向下的消息了")
			case *interactive.EventTryToGetUpper:
				w.SendLine("没有更多向上的消息了")
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
```

---

UPDATE LOG

```
date: 2023.1.28
version: unstable
log:
	实现接口 Win.GotoLeft
	实现接口 Win.GetEventChan
	新增特性 可以通过Win.GetEventChan收到注册过的事件
	新增事件 上移事件EventMoveUp 下移事件EventMoveDown
	新增事件 已经处在最顶端但尝试上移的事件EventTryToGetUpper 已经处在最底端但尝试下移的事件EventTryToGetLower
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
