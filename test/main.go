package main

import (
	"fmt"
	"interactive"
	"sync"
)

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
