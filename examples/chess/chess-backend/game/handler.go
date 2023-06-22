package game

import (
	"chess-backend/comm/chess"
	"chess-backend/comm/packets"
	"chess-backend/comm/settings"

	chesstool "chess-backend/tools/chess"
	othertool "chess-backend/tools/other"

	"github.com/Allenxuxu/gev"
)

type ConnState int

const (
	// 连接上服务端的初始状态状态
	ConnStateNone ConnState = iota
	ConnStateMatching
	ConnStateGaming
)

type GameState int

const (
	GameStateWaitingWhitePut GameState = iota
	GameStateWaitingBlackPut
	// 等待黑方的兵升迁
	GameStateWaitingBlackUpgrade
	// 等待白方的兵升迁
	GameStateWaitingWhiteUpgrade
	// 等待响应, 是否接受和棋
	GameStateWaitingBlackAcceptDraw
	GameStateWaitingWhiteAcceptDraw
)

type ConnHandler struct{}

func (ch *ConnHandler) OnConnect(c *gev.Connection) {
	connID := int(AtomicIDIncrease.Add(1))
	connCtx := &ConnContext{ID: int(connID), LoseHertbeatCount: 0, Conn: c, ConnState: ConnStateNone, Gcontext: nil}

	ConnMapLock.Lock()
	ConnMap[connID] = connCtx
	ConnMapLock.Unlock()

	c.SetContext(connID)
}

func (ch *ConnHandler) OnClose(c *gev.Connection) {
	connID := c.Context().(int)
	ConnMapLock.Lock()
	if ConnMap[connID].ConnState == ConnStateGaming {
		// 需要告知游戏对端, 对手连接丢失
		var remoteConnContext *ConnContext
		if ConnMap[connID].Gcontext.BlackConnContext.ID == connID {
			remoteConnContext = ConnMap[connID].Gcontext.WhiteConnContext
		} else {
			remoteConnContext = ConnMap[connID].Gcontext.BlackConnContext
		}
		packet := packets.PacketServerRemoteLoseConnection{}
		remoteConnContext.Conn.Send(packet.MustMarshalToBytes())
		remoteConnContext.ConnState = ConnStateNone
	}
	delete(ConnMap, connID)
	ConnMapLock.Unlock()
}

func (ch *ConnHandler) OnMessage(c *gev.Connection, ctx interface{}, data []byte) interface{} {
	// 没有收到消息, 继续等待消息传完
	if data == nil {
		return nil
	}

	connID := c.Context().(int)

	packIface := packets.ServerParse(data)

	ConnMapLock.Lock()
	defer ConnMapLock.Unlock()

	switch packet := packIface.(type) {
	case *packets.PacketHeartbeat:
		// 清0丢失心跳计数
		ConnMap[connID].LoseHertbeatCount = 0
	case *packets.PacketClientStartMatch:
		// 协议错误
		if ConnMap[connID].ConnState != ConnStateNone {
			c.Close()
		}

		// 找一个正在match的连接
		for _, v := range ConnMap {
			if v.ID != connID && v.ConnState == ConnStateMatching {
				// 随机摇game side
				var whiteConnContext *ConnContext
				var blackConnContext *ConnContext
				if othertool.RandGetBool() {
					whiteConnContext = ConnMap[connID]
					blackConnContext = ConnMap[v.ID]
				} else {
					blackConnContext = ConnMap[connID]
					whiteConnContext = ConnMap[v.ID]
				}

				// 创建一个默认棋盘
				table := chess.NewChessTable()

				// 建立游戏上下文
				gameContext := GameContext{
					WhiteConnContext: whiteConnContext,
					BlackConnContext: blackConnContext,
					Gstate:           GameStateWaitingWhitePut,
					Table:            table,
					DrawAfterUpgrade: false,
				}

				// matching已经发送给v.ID的conn了, 重复发送可能造成协议错误
				matchingPacket := packets.PacketServerMatching{}
				matchingPacketBytes := matchingPacket.MustMarshalToBytes()
				c.Send(matchingPacketBytes)

				packetForBlack := packets.PacketServerMatchedOK{Side: chess.SideBlack, Table: table}
				blackConnContext.ConnState = ConnStateGaming
				blackConnContext.Gcontext = &gameContext
				blackConnContext.Conn.Send(packetForBlack.MustMarshalToBytes())

				packetForWhite := packets.PacketServerMatchedOK{Side: chess.SideWhite, Table: table}
				whiteConnContext.ConnState = ConnStateGaming
				whiteConnContext.Gcontext = &gameContext
				whiteConnContext.Conn.Send(packetForWhite.MustMarshalToBytes())
				return nil
			}
		}

		// 找不到一个匹配的, 那么标记为正在匹配
		ConnMap[connID].ConnState = ConnStateMatching
		retPacket := packets.PacketServerMatching{}
		ConnMap[connID].Conn.Send(retPacket.MustMarshalToBytes())
		return nil
	case *packets.PacketClientMove:
		// 协议判断
		if ConnMap[connID].ConnState != ConnStateGaming {
			ConnMap[connID].Conn.Close()
			return nil
		}

		// 建立一些信息, 方便写代码
		var gameContext = ConnMap[connID].Gcontext
		var selfContext *ConnContext = ConnMap[connID]
		var selfSide chess.Side
		var remoteContext *ConnContext
		var remoteSide chess.Side
		othertool.Ignore(remoteContext)
		othertool.Ignore(remoteSide)
		if gameContext.BlackConnContext == selfContext {
			remoteContext = gameContext.WhiteConnContext
			selfSide = chess.SideBlack
			remoteSide = chess.SideWhite
		} else {
			remoteContext = gameContext.BlackConnContext
			selfSide = chess.SideWhite
			remoteSide = chess.SideBlack
		}

		// 协议判断, 要求发送方确实是下棋的一方
		if (selfSide == chess.SideBlack && gameContext.Gstate != GameStateWaitingBlackPut) ||
			(selfSide == chess.SideWhite && gameContext.Gstate != GameStateWaitingWhitePut) {
			ConnMap[connID].Conn.Close()
			return nil
		}

		// 协议判断, 输入格式判断, 要求输入格式确实正确
		// 注意x,y两两相等的情况也是不合法的, 这点应该在客户端得到保障
		if !chesstool.CheckChessPostsionVaild(packet.FromX, packet.FromY) ||
			!chesstool.CheckChessPostsionVaild(packet.ToX, packet.ToY) ||
			(packet.FromX == packet.ToX && packet.FromY == packet.ToY) {
			ConnMap[connID].Conn.Close()
			return nil
		}

		result := chesstool.DoMove(gameContext.Table, selfSide, packet.FromX, packet.FromY, packet.ToX, packet.ToY)
		// result.OK 移动是否有效
		if !result.OK {
			moveFailedPacket := packets.PacketServerMoveResp{
				MoveRespType: packets.PacketTypeServerMoveRespTypeFailed,
				TableOnOK:    nil,
			}
			selfContext.Conn.Send(moveFailedPacket.MustMarshalToBytes())
			return nil
		}

		if result.PawnUpgrade {
			if packet.DoDraw {
				gameContext.DrawAfterUpgrade = true
			}

			moveOKPacket := packets.PacketServerMoveResp{
				MoveRespType: packets.PacketTypeServerMoveRespTypePawnUpgrade,
				TableOnOK:    gameContext.Table,
				// 下面的字段此时没有意义
				KingThreat: result.KingThreat,
				Gameover:   false,
			}

			selfContext.Conn.Send(moveOKPacket.MustMarshalToBytes())

			if selfSide == chess.SideWhite {
				gameContext.Gstate = GameStateWaitingWhiteUpgrade
			} else {
				gameContext.Gstate = GameStateWaitingBlackUpgrade
			}

			return nil
		}

		if result.GameOver {
			moveOKPacket := packets.PacketServerMoveResp{
				MoveRespType: packets.PacketTypeServerMoveRespTypeOK,
				TableOnOK:    gameContext.Table,
				Gameover:     true,
				// 这个字段此时没有意义
				KingThreat: false,
			}

			remoteMovePacket := packets.PacketServerNotifyRemoteMove{
				Table:    gameContext.Table,
				Gameover: true,
				// 这两个字段此时没有意义
				KingThreat:        false,
				RemoteRequestDraw: false,
			}

			gameoverPacket := packets.PacketServerGameOver{
				WinnerSide:  result.GameWinner,
				IsSurrender: false,
				// 这个字段的意思是主动和棋
				IsDraw: false,
			}

			selfContext.Conn.Send(moveOKPacket.MustMarshalToBytes())
			selfContext.Conn.Send(gameoverPacket.MustMarshalToBytes())

			remoteContext.Conn.Send(remoteMovePacket.MustMarshalToBytes())
			remoteContext.Conn.Send(gameoverPacket.MustMarshalToBytes())

			selfContext.ConnState = ConnStateNone
			selfContext.Gcontext = nil
			remoteContext.ConnState = ConnStateNone
			remoteContext.Gcontext = nil
			return nil
		}

		// 游戏没有结束

		moveOKPacket := packets.PacketServerMoveResp{
			MoveRespType: packets.PacketTypeServerMoveRespTypeOK,
			TableOnOK:    gameContext.Table,
			KingThreat:   result.KingThreat,
			Gameover:     false,
		}

		selfContext.Conn.Send(moveOKPacket.MustMarshalToBytes())

		notifyRemoteMovePacket := packets.PacketServerNotifyRemoteMove{
			Table:             gameContext.Table,
			Gameover:          false,
			KingThreat:        result.KingThreat,
			RemoteRequestDraw: packet.DoDraw,
		}

		remoteContext.Conn.Send(notifyRemoteMovePacket.MustMarshalToBytes())

		if packet.DoDraw {
			if selfSide == chess.SideWhite {
				gameContext.Gstate = GameStateWaitingBlackAcceptDraw
			} else {
				gameContext.Gstate = GameStateWaitingWhiteAcceptDraw
			}
		} else {
			if selfSide == chess.SideWhite {
				gameContext.Gstate = GameStateWaitingBlackPut
			} else {
				gameContext.Gstate = GameStateWaitingWhitePut
			}
		}

		return nil
	case *packets.PacketClientPawnUpgrade:
		// 协议判断
		if ConnMap[connID].ConnState != ConnStateGaming {
			ConnMap[connID].Conn.Close()
			return nil
		}

		// 拿到一些信息
		var gameContext = ConnMap[connID].Gcontext
		var selfContext *ConnContext = ConnMap[connID]
		var selfSide chess.Side
		var remoteContext *ConnContext
		var remoteSide chess.Side
		othertool.Ignore(remoteContext)
		othertool.Ignore(remoteSide)
		if gameContext.BlackConnContext == selfContext {
			remoteContext = gameContext.WhiteConnContext
			selfSide = chess.SideBlack
			remoteSide = chess.SideWhite
		} else {
			remoteContext = gameContext.BlackConnContext
			selfSide = chess.SideWhite
			remoteSide = chess.SideBlack
		}

		// 协议判断
		if selfSide == chess.SideWhite && gameContext.Gstate != GameStateWaitingWhiteUpgrade {
			selfContext.Conn.Close()
			return nil
		}
		if selfSide == chess.SideBlack && gameContext.Gstate != GameStateWaitingBlackUpgrade {
			selfContext.Conn.Close()
			return nil
		}

		// 协议判断, 检查升变的棋子是否合法, 只允许以下4种棋子
		if typ := packet.ChessPieceType; typ != chess.ChessPieceTypeRook && typ != chess.ChessPieceTypeBishop &&
			typ != chess.ChessPieceTypeKnight && typ != chess.ChessPieceTypeQueen {
			selfContext.Conn.Close()
			return nil
		}

		result := chesstool.DoUpgrade(gameContext.Table, selfSide, remoteSide, packet.ChessPieceType)
		if result.GameOver {
			notifyRemoteMovePacket := packets.PacketServerNotifyRemoteMove{
				Table:    gameContext.Table,
				Gameover: true,
				// 下面两个字段此时无意义
				KingThreat:        false,
				RemoteRequestDraw: false,
			}
			notifyUpgradeOKPacket := packets.PacketServerUpgradeOK{
				Table:    gameContext.Table,
				Gameover: true,
			}
			selfContext.Conn.Send(notifyUpgradeOKPacket.MustMarshalToBytes())
			remoteContext.Conn.Send(notifyRemoteMovePacket.MustMarshalToBytes())

			gameoverPacket := packets.PacketServerGameOver{
				WinnerSide:  result.WinnerSide,
				IsSurrender: false,
				IsDraw:      false,
			}

			selfContext.Conn.Send(gameoverPacket.MustMarshalToBytes())
			remoteContext.Conn.Send(gameoverPacket.MustMarshalToBytes())

			selfContext.ConnState = ConnStateNone
			selfContext.Gcontext = nil
			remoteContext.ConnState = ConnStateNone
			remoteContext.Gcontext = nil

			return nil
		}

		if gameContext.DrawAfterUpgrade {
			gameContext.DrawAfterUpgrade = false
			if selfSide == chess.SideWhite {
				gameContext.Gstate = GameStateWaitingBlackAcceptDraw
			} else {
				gameContext.Gstate = GameStateWaitingWhiteAcceptDraw
			}

			notifyUpgradeOKPacket := packets.PacketServerUpgradeOK{
				Table:    gameContext.Table,
				Gameover: false,
			}
			selfContext.Conn.Send(notifyUpgradeOKPacket.MustMarshalToBytes())

			notifyRemoteMovePacket := packets.PacketServerNotifyRemoteMove{
				Table:             gameContext.Table,
				Gameover:          false,
				RemoteRequestDraw: true,
				KingThreat:        result.KingThreat,
			}
			remoteContext.Conn.Send(notifyRemoteMovePacket.MustMarshalToBytes())

			return nil
		}

		notifyRemoteMovePacket := packets.PacketServerNotifyRemoteMove{
			Table:             gameContext.Table,
			Gameover:          false,
			KingThreat:        result.KingThreat,
			RemoteRequestDraw: false,
		}
		notifyUpgradeOKPacket := packets.PacketServerUpgradeOK{
			Table:    gameContext.Table,
			Gameover: false,
		}
		selfContext.Conn.Send(notifyUpgradeOKPacket.MustMarshalToBytes())
		remoteContext.Conn.Send(notifyRemoteMovePacket.MustMarshalToBytes())

		if selfSide == chess.SideWhite {
			gameContext.Gstate = GameStateWaitingBlackPut
		} else {
			gameContext.Gstate = GameStateWaitingWhitePut
		}

		return nil
	case *packets.PacketClientDoSurrender:
		// 协议判断
		if ConnMap[connID].ConnState != ConnStateGaming {
			ConnMap[connID].Conn.Close()
			return nil
		}

		var gameContext = ConnMap[connID].Gcontext
		var selfContext *ConnContext = ConnMap[connID]
		var selfSide chess.Side
		othertool.Ignore(selfSide)
		var remoteContext *ConnContext
		var remoteSide chess.Side
		othertool.Ignore(remoteContext)
		othertool.Ignore(remoteSide)
		if gameContext.BlackConnContext == selfContext {
			remoteContext = gameContext.WhiteConnContext
			selfSide = chess.SideBlack
			remoteSide = chess.SideWhite
		} else {
			remoteContext = gameContext.BlackConnContext
			selfSide = chess.SideWhite
			remoteSide = chess.SideBlack
		}

		gameOverPacket := packets.PacketServerGameOver{
			WinnerSide:  remoteSide,
			IsSurrender: true,
		}
		selfContext.Conn.Send(gameOverPacket.MustMarshalToBytes())
		remoteContext.Conn.Send(gameOverPacket.MustMarshalToBytes())
		selfContext.Gcontext = nil
		selfContext.ConnState = ConnStateNone
		remoteContext.Gcontext = nil
		remoteContext.ConnState = ConnStateNone
		return nil
	case *packets.PacketClientWhetherAcceptDraw:
		// 协议判断
		if ConnMap[connID].ConnState != ConnStateGaming {
			ConnMap[connID].Conn.Close()
			return nil
		}

		var gameContext = ConnMap[connID].Gcontext
		var selfContext *ConnContext = ConnMap[connID]
		var selfSide chess.Side
		var remoteContext *ConnContext
		var remoteSide chess.Side
		othertool.Ignore(remoteContext)
		othertool.Ignore(remoteSide)
		if gameContext.BlackConnContext == selfContext {
			remoteContext = gameContext.WhiteConnContext
			selfSide = chess.SideBlack
			remoteSide = chess.SideWhite
		} else {
			remoteContext = gameContext.BlackConnContext
			selfSide = chess.SideWhite
			remoteSide = chess.SideBlack
		}

		// 判断更多协议错误
		if selfSide == chess.SideWhite && gameContext.Gstate != GameStateWaitingWhiteAcceptDraw {
			ConnMap[connID].Conn.Close()
			return nil
		}
		if selfSide == chess.SideBlack && gameContext.Gstate != GameStateWaitingBlackAcceptDraw {
			ConnMap[connID].Conn.Close()
			return nil
		}

		if packet.AcceptDraw {
			gameOver := packets.PacketServerGameOver{
				WinnerSide:  chess.SideBoth,
				IsSurrender: false,
				IsDraw:      true,
			}
			selfContext.Conn.Send(gameOver.MustMarshalToBytes())
			remoteContext.Conn.Send(gameOver.MustMarshalToBytes())
			selfContext.ConnState = ConnStateNone
			selfContext.Gcontext = nil
			remoteContext.ConnState = ConnStateNone
			remoteContext.Gcontext = nil
			return nil
		}

		// 不接受和棋
		if selfSide == chess.SideWhite {
			gameContext.Gstate = GameStateWaitingWhitePut
		} else {
			gameContext.Gstate = GameStateWaitingBlackPut
		}
	case nil:
		// 协议错误, 直接关闭
		c.Close()
	}
	return nil
}

func OnTimeout() {
	var packet = packets.PacketHeartbeat{}

	ConnMapLock.Lock()
	for k := range ConnMap {
		ConnMap[k].Conn.Send(packet.MustMarshalToBytes())
		ConnMap[k].LoseHertbeatCount++
		if ConnMap[k].LoseHertbeatCount >= settings.MaxLoseHeartbeat {
			ConnMap[k].Conn.Close()
		}
	}
	ConnMapLock.Unlock()
}
