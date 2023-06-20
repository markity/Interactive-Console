package main

import (
	commsettings "chinese-chess-backend/comm_settings"
	gamehandler "chinese-chess-backend/game_handler"
	"fmt"
	"runtime"
	"time"

	"github.com/Allenxuxu/gev"
)

func main() {
	server, err := gev.NewServer(&gamehandler.ConnHandler{},
		gev.Address(fmt.Sprintf("%s:%d", commsettings.ServerListenIP, commsettings.ServerListenPort)),
		gev.Network("tcp"),
		gev.LoadBalance(gev.RoundRobin()),
		gev.NumLoops(runtime.NumCPU()),
	)
	if err != nil {
		panic(err)
	}

	// 这个是全局的心跳检测器
	server.RunEvery(time.Millisecond*commsettings.HeartbeatInterval, gamehandler.OnTimeout)

	server.Start()
}
