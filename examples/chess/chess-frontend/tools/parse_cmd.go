package tools

import (
	"chess-frontend/comm/chess"
	"strings"
)

type CommandType int

const (
	CommandTypeEmpty CommandType = iota
	// mov
	CommandTypeMove
	// dmov
	CommandTypeMoveAndDraw
	// sur
	CommandTypeSurrender
	// swi 兵生变的指令
	CommandTypeSwitch
	CommandTypeUnkonwn
	// 是否同一平局
	CommandTypeAccept
	CommandTypeRefuse
)

type CommandPattern struct {
	Type CommandType

	MoveFromX rune
	MoveFromY int

	MoveToX rune
	MoveToY int

	Swi chess.ChessPieceType
}

func runeAlphaToIndex(i rune) (int, bool) {
	switch i {
	case 'a':
		return 1, true
	case 'b':
		return 2, true
	case 'c':
		return 3, true
	case 'd':
		return 4, true
	case 'e':
		return 5, true
	case 'f':
		return 6, true
	case 'g':
		return 7, true
	case 'h':
		return 8, true
	}

	return 0, false
}

func runeNumToIndex(i rune) (int, bool) {
	switch i {
	case '1':
		return 1, true
	case '2':
		return 2, true
	case '3':
		return 3, true
	case '4':
		return 4, true
	case '5':
		return 5, true
	case '6':
		return 6, true
	case '7':
		return 7, true
	case '8':
		return 8, true
	}

	return 0, false
}

func ParseCommand(s string) *CommandPattern {
	s = strings.ToLower(s)
	fields := strings.Fields(s)

	if len(fields) == 0 {
		return &CommandPattern{Type: CommandTypeEmpty}
	}

	if len(fields) == 1 {
		if fields[0] == "sur" {
			return &CommandPattern{Type: CommandTypeSurrender}
		}

		if fields[0] == "accept" {
			return &CommandPattern{Type: CommandTypeAccept}
		}

		if fields[0] == "refuse" {
			return &CommandPattern{Type: CommandTypeRefuse}
		}
		return &CommandPattern{Type: CommandTypeUnkonwn}
	}

	if len(fields) == 2 {
		if fields[0] != "swi" {
			return &CommandPattern{Type: CommandTypeSwitch}
		}

		switch fields[1] {
		case "rook":
			return &CommandPattern{Type: CommandTypeSwitch, Swi: chess.ChessPieceTypeRook}
		case "bishop":
			return &CommandPattern{Type: CommandTypeSwitch, Swi: chess.ChessPieceTypeBishop}
		case "knight":
			return &CommandPattern{Type: CommandTypeSwitch, Swi: chess.ChessPieceTypeKnight}
		case "queen":
			return &CommandPattern{Type: CommandTypeSwitch, Swi: chess.ChessPieceTypeQueen}
		default:
			return &CommandPattern{Type: CommandTypeUnkonwn}
		}
	}

	if len(fields) == 3 {
		if fields[0] != "mov" && fields[0] != "dmov" {
			return &CommandPattern{Type: CommandTypeUnkonwn}
		}

		from := []rune(fields[1])
		to := []rune(fields[2])

		// 检查两个坐标是否合法
		if len(from) != 2 || len(to) != 2 {
			return &CommandPattern{Type: CommandTypeUnkonwn}
		}

		_, ok := runeAlphaToIndex(from[0])
		if !ok {
			return &CommandPattern{Type: CommandTypeUnkonwn}
		}

		fromy, ok := runeNumToIndex(from[1])
		if !ok {
			return &CommandPattern{Type: CommandTypeUnkonwn}
		}

		_, ok = runeAlphaToIndex(to[0])
		if !ok {
			return &CommandPattern{Type: CommandTypeUnkonwn}
		}

		toy, ok := runeNumToIndex(to[1])
		if !ok {
			return &CommandPattern{Type: CommandTypeUnkonwn}
		}

		var p CommandType
		if fields[0] == "mov" {
			p = CommandTypeMove
		} else {
			p = CommandTypeMoveAndDraw
		}
		return &CommandPattern{Type: p, MoveFromX: from[0], MoveFromY: fromy, MoveToX: to[0], MoveToY: toy}
	}

	return &CommandPattern{Type: CommandTypeUnkonwn}
}
