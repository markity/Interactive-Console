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
func cmdHandler(w *interactive.Win, wait *sync.WaitGroup) {
	for {
		select {
		case cmd := <-w.GetCmdChan():
			switch cmd {
			case "trace":
				w.SetTrace(true)
				w.BlockInput(false)
			case "untrace":
				w.SetTrace(false)
				w.BlockInput(false)
			case "ping":
				w.SendLine("pong")
				w.BlockInput(false)
			case "exit":
				w.Stop()
				goto out
			case "clear":
				w.Clear()
				w.SendLine("你已经清空了屏幕")
				w.BlockInput(false)
			default:
				w.SendLine(fmt.Sprintf("unknown command %s", cmd))
				w.BlockInput(false)
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
	win := interactive.Run(cfg)
	wait := sync.WaitGroup{}
	wait.Add(1)
	go cmdHandler(win, &wait)
	wait.Wait()
}
```

