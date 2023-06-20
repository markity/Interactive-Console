package game

import (
	"chess-backend/comm/chess"
	"sync"
	"sync/atomic"

	"github.com/Allenxuxu/gev"
)

type ConnContext struct {
	ID                int
	LoseHertbeatCount int
	Conn              *gev.Connection
	ConnState         ConnState

	// 下面的字段只有在ConnState为Gaming时有意义
	Gcontext *GameContext
}

type GameContext struct {
	BlackConnContext *ConnContext
	WhiteConnContext *ConnContext
	Gstate           GameState
	Table            *chess.ChessTable
	DrawAfterUpgrade bool
}

// 包含所有连接的上下文, 用锁保护
var ConnMap map[int]*ConnContext
var ConnMapLock sync.Mutex

func init() {
	ConnMap = make(map[int]*ConnContext)
}

// 用来做自增连接id的计数器
var AtomicIDIncrease atomic.Int32
