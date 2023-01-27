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
