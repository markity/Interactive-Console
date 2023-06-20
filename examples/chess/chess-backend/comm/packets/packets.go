package packets

import (
	"chess-backend/comm/chess"
	"encoding/json"
)

type PacketTypeServerMoveRespType int

const (
	PacketTypeServerMoveRespTypeOK PacketTypeServerMoveRespType = iota
	PacketTypeServerMoveRespTypeFailed
	PacketTypeServerMoveRespTypePawnUpgrade
)

type PacketType int

const (
	// 心跳包
	PacketTypeHeartbeat PacketType = iota

	// 客户端要求开始匹配
	PacketTypeClientStartMatch

	// 服务端表示已经开始匹配, 这是用来回复StartMatch的
	PacketTypeServerMatching

	// 匹配完毕, 即将开始游戏
	PacketTypeServerMatchedOK

	// 客户端发送下棋的消息
	PacketTypeClientMove

	// 服务端告知用户下棋结果, 可能用户的输入不合法, 这里提示, 可能成功, 可能发生兵的升变, 要求用户继续输入
	PacketTypeServerMoveResp

	// 告知服务端兵升变成什么
	PacketTypeClientPawnUpgrade

	// 服务端告知升变完成
	PacketTypeServerUpgradeOK

	// 通知游戏结束
	PacketTypeServerGameOver

	// 通知对方掉线
	PacketTypeServerRemoteLoseConnection

	// 对方下棋下好了
	PacketTypeServerNotifyRemoteMove

	// 告知对方, 自己是否接受和棋
	PacketTypeClientWhetherAcceptDraw

	// 直接投降
	PacketTypeClientDoSurrender
)

type PacketHeader struct {
	Type *PacketType `json:"type"`
}

type PacketHeartbeat struct {
	PacketHeader
}

func (p *PacketHeartbeat) MustMarshalToBytes() []byte {
	i := PacketTypeHeartbeat
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketClientStartMatch struct {
	PacketHeader
}

func (p *PacketClientStartMatch) MustMarshalToBytes() []byte {
	i := PacketTypeClientStartMatch
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerMatching struct {
	PacketHeader
}

func (p *PacketServerMatching) MustMarshalToBytes() []byte {
	i := PacketTypeServerMatching
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerMatchedOK struct {
	PacketHeader
	Side  chess.Side        `json:"game_side"`
	Table *chess.ChessTable `json:"game_table"`
}

func (p *PacketServerMatchedOK) MustMarshalToBytes() []byte {
	i := PacketTypeServerMatchedOK
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketClientMove struct {
	PacketHeader
	FromX rune `json:"from_x"`
	FromY int  `json:"from_y"`
	ToX   rune `json:"to_x"`
	ToY   int  `json:"to_y"`

	// 是否发起和棋请求
	DoDraw bool `json:"do_draw"`
}

func (p *PacketClientMove) MustMarshalToBytes() []byte {
	i := PacketTypeClientMove
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketClientPawnUpgrade struct {
	PacketHeader
	ChessPieceType chess.ChessPieceType `json:"piece_type"`
}

func (p *PacketClientPawnUpgrade) MustMarshalToBytes() []byte {
	i := PacketTypeClientPawnUpgrade
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerMoveResp struct {
	PacketHeader
	// 有几种情况, 1. 移动OK, 2. 现在要求进行upgrade 3. 失败
	MoveRespType PacketTypeServerMoveRespType `json:"resp_type"`
	// 下面的字段只有在状态OK的时候出现
	TableOnOK *chess.ChessTable `json:"table,omitempty"`
	// 我认为兵移动到底线并升变是同一时间发生的, 没有中间过程。因此要升变之后再告知将军或者游戏结束的结果
	// 是否游戏结束, 游戏结束要求客户端等待gameover包
	Gameover bool `json:"gameover"`
	// 是否正在将对方的军, 如果失败或要求进行upgrade, 这个字段没意义
	KingThreat bool `json:"king_threat"`
}

func (p *PacketServerMoveResp) MustMarshalToBytes() []byte {
	i := PacketTypeServerMoveResp
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerGameOver struct {
	PacketHeader
	WinnerSide  chess.Side `json:"winner_side"`
	IsSurrender bool       `json:"is_surrender"`
	IsDraw      bool       `json:"is_draw"`
}

func (p *PacketServerGameOver) MustMarshalToBytes() []byte {
	i := PacketTypeServerGameOver
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerRemoteLoseConnection struct {
	PacketHeader
}

func (p *PacketServerRemoteLoseConnection) MustMarshalToBytes() []byte {
	i := PacketTypeServerRemoteLoseConnection
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerNotifyRemoteMove struct {
	PacketHeader
	Table             *chess.ChessTable `json:"table"`
	Gameover          bool              `json:"gameover"`
	KingThreat        bool              `json:"king_threat"`
	RemoteRequestDraw bool              `json:"remote_request_draw"`
}

func (p *PacketServerNotifyRemoteMove) MustMarshalToBytes() []byte {
	i := PacketTypeServerNotifyRemoteMove
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketClientWhetherAcceptDraw struct {
	PacketHeader
	AcceptDraw bool `json:"accept_draw"`
}

func (p *PacketClientWhetherAcceptDraw) MustMarshalToBytes() []byte {
	i := PacketTypeClientWhetherAcceptDraw
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketClientDoSurrender struct {
	PacketHeader
}

func (p *PacketClientDoSurrender) MustMarshalToBytes() []byte {
	i := PacketTypeClientDoSurrender
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerUpgradeOK struct {
	PacketHeader
	Table *chess.ChessTable `json:"table"`
	// 收到gameover后应该等待gameover包
	Gameover bool `json:"gameover"`
}

func (p *PacketServerUpgradeOK) MustMarshalToBytes() []byte {
	i := PacketTypeServerUpgradeOK
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}
