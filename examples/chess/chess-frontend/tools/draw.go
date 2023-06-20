package tools

import (
	"chess-frontend/comm/chess"
	"fmt"

	interactive "github.com/markity/Interactive-Console"
)

func mustGetName(piece *chess.ChessPiece) string {
	if piece == nil {
		return "  "
	}
	switch piece.PieceType {
	case chess.ChessPieceTypeRook:
		return "车"
	case chess.ChessPieceTypeBishop:
		return "象"
	case chess.ChessPieceTypeKnight:
		return "马"
	case chess.ChessPieceTypeQueen:
		return "后"
	case chess.ChessPieceTypeKing:
		return "王"
	case chess.ChessPieceTypePawn:
		return "兵"
	default:
		panic("unreachable")
	}
}

func Draw(win *interactive.Win, table *chess.ChessTable, message *string) {
	win.Clear()
	style1 := interactive.GetDefaultSytleAttr()
	style1.Foreground = interactive.ColorForestGreen

	win.SendLineBackWithColor(style1, "   a b c d e f g h    ")

	for i := 7; i >= 0; i-- {
		tobeSend := make([]interface{}, 0)
		tobeSend = append(tobeSend, style1, " "+fmt.Sprint(i+1)+" ")
		for j := 0; j < 8; j++ {
			style2 := style1
			if table.GetIndex(j, i) != nil {
				if table.GetIndex(j, i).GameSide == chess.SideWhite {
					style2.Foreground = interactive.ColorGhostWhite
				} else {
					style2.Foreground = interactive.ColorDarkGrey
				}
			}

			tobeSend = append(tobeSend, style2, mustGetName(table.GetIndex(j, i)), style2)
		}
		tobeSend = append(tobeSend, style1, " "+fmt.Sprint(i+1)+" ")
		win.SendLineBackWithColor(tobeSend...)
	}

	win.SendLineBackWithColor(style1, "   a b c d e f g h    ")

	style2 := interactive.GetDefaultSytleAttr()
	style2.Foreground = interactive.ColorLightPink
	if message != nil {
		win.SendLineBackWithColor(style2, *message)
	}
}
