package chess

import "chess-backend/comm/chess"

type UpgradeResult struct {
	GameOver   bool
	KingThreat bool
	WinnerSide chess.Side
}

func DoUpgrade(table *chess.ChessTable, side chess.Side, remoteSide chess.Side, targetPieceType chess.ChessPieceType) (result UpgradeResult) {
	for _, v := range table {
		if v != nil && v.GameSide == chess.SideWhite && v.Y == 8 && v.PieceType == chess.ChessPieceTypePawn {
			v.PieceType = targetPieceType
		}

		if v != nil && v.GameSide == chess.SideBlack && v.Y == 1 && v.PieceType == chess.ChessPieceTypePawn {
			v.PieceType = targetPieceType
		}
	}

	remoteKing := findKing(table, remoteSide)

	// 是否将军
	kingThreat := checkPositionThreat(table, remoteSide, remoteKing.X, remoteKing.Y)

	// 王的8个单元格是否都受威胁
	movable := isMovable(table, remoteSide)

	// 赢
	if kingThreat && !movable {
		result.GameOver = true
		result.WinnerSide = side
		return
	}

	// 判断和棋
	if !kingThreat && !movable {
		result.GameOver = true
		result.WinnerSide = chess.SideBoth
		return
	}

	result.GameOver = false
	result.KingThreat = kingThreat
	return
}
