package main

import (
	"chess-frontend/comm/chess"
	"chess-frontend/comm/packets"
	"chess-frontend/comm/settings"
	"chess-frontend/tools"
	"fmt"
	"net"
	"time"

	interactive "github.com/markity/Interactive-Console"
)

type GameState int

const (
	GameStateConnected GameState = iota
	GameStateMatching
	GameStateWaitSelfPut
	GameStateWaitPutResp
	GameStateWaitRemotePut
	GameStateWaitUpgrade
	GameStateWaitUpgradeOK
	GameStateWaitGameover
	GameStateWaitAcceptOrRefuse
)

var State GameState

// 下面的全局变量根据状态的不同而酌情使用
var LoseHeartbeat int
var Win *interactive.Win
var SelfSide chess.Side
var CmdChan chan string
var Table *chess.ChessTable
var KingThreat bool

func main() {
	// 连上服务端
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", settings.ServerListenIP, settings.ServerListenPort))
	if err != nil {
		fmt.Printf("failed to dial to server: %v\n", err)
		return
	}

	State = GameStateConnected
	LoseHeartbeat = 0

	// 连接完成后发送start_match指令
	startMatchPacket := packets.PacketClientStartMatch{}
	startMatchPacketBytesWithHeader := tools.DoPackWith4BytesHeader(startMatchPacket.MustMarshalToBytes())
	_, err = conn.Write(startMatchPacketBytesWithHeader)
	if err != nil {
		fmt.Printf("network error: %v\n", err)
		return
	}

	errChan := make(chan error, 2)
	writeToQueen := tools.NewUnboundedQueen()
	readFromConnChan := make(chan interface{})

	// writer
	go func() {
		for {
			bs := writeToQueen.PopBlock().([]byte)
			_, err := conn.Write(tools.DoPackWith4BytesHeader(bs))
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	// reader
	go func() {
		for {
			packetBytes, err := tools.ReadPacketBytesWith4BytesHeader(conn)
			if err != nil {
				errChan <- err
				return
			}

			readFromConnChan <- packets.ClientParse(packetBytes)
		}
	}()

	heartTicker := time.NewTicker(time.Millisecond * settings.HeartbeatInterval)
	// 现在state为Connected, 要等待到进入游戏为止
	for {
		select {
		case <-heartTicker.C:
			LoseHeartbeat++
			if LoseHeartbeat == settings.MaxLoseHeartbeat {
				fmt.Println("network error: heartbeat error")
				return
			}
			heartbeatPacket := packets.PacketHeartbeat{}
			writeToQueen.Push(heartbeatPacket.MustMarshalToBytes())
		case inpacket := <-readFromConnChan:
			switch packet := inpacket.(type) {
			case *packets.PacketHeartbeat:
				LoseHeartbeat = 0
			case *packets.PacketServerMatchedOK:
				if State != GameStateMatching {
					fmt.Printf("protocol error: waiting for PacketServerMatchedOK, but received unexpected packet\n")
				}
				SelfSide = packet.Side
				Table = packet.Table
				goto out
			case *packets.PacketServerMatching:
				if State != GameStateConnected {
					fmt.Printf("protocol error: waiting for PacketServerMatching, but received unexpected packet\n")
				}
				State = GameStateMatching
				println("matching...")
			default:
				fmt.Printf("protocol error: unexpected income bytes\n")
				return
			}
		case err := <-errChan:
			fmt.Printf("network error: %v\n", err)
			return
		}
	}

out:

	var msg string
	if SelfSide == chess.SideWhite {
		State = GameStateWaitSelfPut
		msg = "you are white, your turn"
	} else {
		State = GameStateWaitRemotePut
		msg = "you are black, remote user's turn"
	}
	cfg := interactive.GetDefaultConfig()
	cfg.BlockInputAfterEnter = true
	Win = interactive.Run(cfg)
	CmdChan = Win.GetCmdChan()
	tools.Draw(Win, Table, &msg)

	for {
		select {
		case <-heartTicker.C:
			LoseHeartbeat++
			if LoseHeartbeat == settings.MaxLoseHeartbeat {
				fmt.Println("network error: heartbeat error")
				return
			}
			heartbeatPacket := packets.PacketHeartbeat{}
			writeToQueen.Push(heartbeatPacket.MustMarshalToBytes())
		case err := <-errChan:
			Win.Stop()
			fmt.Printf("network error: %v\n", err)
			return
		case cmd := <-CmdChan:
			pattern := tools.ParseCommand(cmd)
			switch pattern.Type {
			case tools.CommandTypeEmpty:
				Win.SetBlockInput(false)
			case tools.CommandTypeUnkonwn:
				msg := "unknown command"
				tools.Draw(Win, Table, &msg)
				Win.SetBlockInput(false)
			case tools.CommandTypeSurrender:
				bs := packets.PacketClientDoSurrender{}
				writeToQueen.Push(bs.MustMarshalToBytes())
				State = GameStateWaitGameover
			case tools.CommandTypeSwitch:
				if State != GameStateWaitUpgrade {
					msg := "you cannot upgrade now"
					tools.Draw(Win, Table, &msg)
					Win.SetBlockInput(false)
					continue
				}

				upgradePacket := packets.PacketClientPawnUpgrade{
					ChessPieceType: pattern.Swi,
				}
				writeToQueen.Push(upgradePacket.MustMarshalToBytes())
				State = GameStateWaitUpgradeOK
			case tools.CommandTypeAccept:
				if State != GameStateWaitAcceptOrRefuse {
					msg := "no draw request yet"
					tools.Draw(Win, Table, &msg)
					Win.SetBlockInput(false)
					continue
				}

				drawRespPacket := packets.PacketClientWhetherAcceptDraw{
					AcceptDraw: true,
				}
				writeToQueen.Push(drawRespPacket.MustMarshalToBytes())
				State = GameStateWaitGameover
			case tools.CommandTypeRefuse:
				if State != GameStateWaitAcceptOrRefuse {
					msg := "no draw request yet"
					tools.Draw(Win, Table, &msg)
					Win.SetBlockInput(false)
					continue
				}

				drawRespPacket := packets.PacketClientWhetherAcceptDraw{
					AcceptDraw: false,
				}
				writeToQueen.Push(drawRespPacket.MustMarshalToBytes())
				State = GameStateWaitSelfPut
				msg := "your turn"
				tools.Draw(Win, Table, &msg)
				Win.SetBlockInput(false)
			case tools.CommandTypeMove, tools.CommandTypeMoveAndDraw:
				if State != GameStateWaitSelfPut {
					msg := "you cannot do this now"
					tools.Draw(Win, Table, &msg)
					Win.SetBlockInput(false)
					continue
				}

				if pattern.MoveFromX == pattern.MoveToX && pattern.MoveFromY == pattern.MoveToY {
					msg := "two positions must be different"
					tools.Draw(Win, Table, &msg)
					Win.SetBlockInput(false)
					continue
				}

				movPacket := packets.PacketClientMove{
					FromX:  pattern.MoveFromX,
					FromY:  pattern.MoveFromY,
					ToX:    pattern.MoveToX,
					ToY:    pattern.MoveToY,
					DoDraw: pattern.Type == tools.CommandTypeMoveAndDraw,
				}
				writeToQueen.Push(movPacket.MustMarshalToBytes())
				State = GameStateWaitPutResp
			}
		case pkgIface := <-readFromConnChan:
			switch packet := pkgIface.(type) {
			case *packets.PacketServerGameOver:
				var msg string
				if packet.IsDraw {
					msg = "draw, game will quit in 3 seconds"
				} else if packet.IsSurrender {
					if packet.WinnerSide == chess.SideWhite {
						msg = "black surrendered, winner is white, game will quit in 3 seconds"
					} else {
						msg = "white surrendered, winner is black, game will quit in 3 seconds"
					}
				} else {
					if packet.WinnerSide == chess.SideWhite {
						msg = "winner is white, game will quit in 3 seconds"
					} else {
						msg = "winner is black, game will quit in 3 seconds"
					}
				}

				tools.Draw(Win, Table, &msg)
				time.Sleep(time.Second * 3)
				Win.Stop()

				return
			case *packets.PacketServerMoveResp:
				if State != GameStateWaitPutResp {
					Win.Stop()
					fmt.Println("protocol error: received PacketServerMoveResp, but State is not GameStateWaitPutResp")
					return
				}

				if packet.MoveRespType == packets.PacketTypeServerMoveRespTypeFailed {
					msg := "invalid move, check it again"
					tools.Draw(Win, Table, &msg)
					State = GameStateWaitSelfPut
					Win.SetBlockInput(false)
					continue
				}

				Table = packet.TableOnOK

				if packet.Gameover {
					State = GameStateWaitGameover
					continue
				}

				if packet.MoveRespType == packets.PacketTypeServerMoveRespTypePawnUpgrade {
					msg := "type swi queen/bishop/knight/rook to do pawn upgrade"
					tools.Draw(Win, Table, &msg)
					State = GameStateWaitUpgrade
					Win.SetBlockInput(false)
					continue
				}

				var msg string
				if packet.KingThreat {
					msg = "remote user's turn, now is threating his king"
				} else {
					msg = "remote user's turn"
				}
				tools.Draw(Win, Table, &msg)

				State = GameStateWaitRemotePut
				Win.SetBlockInput(false)
			case *packets.PacketServerNotifyRemoteMove:
				if State != GameStateWaitRemotePut {
					Win.Stop()
					fmt.Println("protocol error: received PacketServerNotifyRemoteMove, but State is not GameStateWaitRemotefPut")
					return
				}

				Table = packet.Table

				if packet.Gameover {
					State = GameStateWaitGameover
					Win.SetBlockInput(true)
					continue
				}

				if packet.RemoteRequestDraw {
					KingThreat = packet.KingThreat
					msg := "remote user is requesting to draw, input accept or refuse"
					tools.Draw(Win, packet.Table, &msg)
					State = GameStateWaitAcceptOrRefuse
					continue
				}

				var msg string
				if packet.KingThreat {
					msg = "your turn, your king is under threat"
				} else {
					msg = "your turn"
				}
				tools.Draw(Win, Table, &msg)
				State = GameStateWaitSelfPut
			case *packets.PacketServerRemoteLoseConnection:
				msg := "remote user lost connection, you will quit in 3 seconds"
				tools.Draw(Win, Table, &msg)
				time.Sleep(time.Second * 3)
				Win.Stop()
				return
			case *packets.PacketServerUpgradeOK:
				if State != GameStateWaitUpgradeOK {
					Win.Stop()
					fmt.Println("protocol error: received PacketServerUpgradeOK, but State is not GameStateWaitUpgradeOK")
					return
				}

				Table = packet.Table

				if packet.Gameover {
					State = GameStateWaitGameover
					continue
				}

				msg := "remote user's turn"

				tools.Draw(Win, Table, &msg)

			case *packets.PacketHeartbeat:
				LoseHeartbeat = 0
			default:
				fmt.Printf("protocol error: unexpected income bytes\n")
				return
			}
		}
	}
}
