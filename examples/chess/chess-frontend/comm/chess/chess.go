package chess

type ChessPieceType int

const (
	// 车
	ChessPieceTypeRook ChessPieceType = iota
	// 马
	ChessPieceTypeKnight
	// 象
	ChessPieceTypeBishop
	// 后
	ChessPieceTypeQueen
	// 王
	ChessPieceTypeKing
	// 兵
	ChessPieceTypePawn
)

type Side int

const (
	SideWhite Side = iota
	SideBlack
	// 用来指示平局
	SideBoth
)

// 棋子
type ChessPiece struct {
	// 棋子的类型
	PieceType ChessPieceType
	// 比如abcdefgh
	X rune
	// 比如12345678
	Y int
	// 游戏方
	GameSide Side
	// 是否移动过, 这个可以用来判定王车易位
	Moved bool

	// 给兵留的变量, 这是用来判定是否吃过路兵的
	PawnMovedTwoLastTime bool
}

func (p *ChessPiece) MustSetIndex(x int, y int) {
	X, Y := MustIndexToPosition(x, y)
	p.X = X
	p.Y = Y
}

// 棋盘类型
type ChessTable [64]*ChessPiece

// 确保传入的是有效的
func (ct *ChessTable) SetPosition(newPiece *ChessPiece) *ChessPiece {
	x, y := MustPositionToIndex(newPiece.X, newPiece.Y)
	oldPiece := ct[y*8+x]
	ct[y*8+x] = newPiece
	return oldPiece
}

func (ct *ChessTable) ClearPosition(X rune, Y int) *ChessPiece {
	x, y := MustPositionToIndex(X, Y)
	oldPiece := ct[y*8+x]
	ct[y*8+x] = nil
	return oldPiece
}

func (ct *ChessTable) ClearIndex(x int, y int) *ChessPiece {
	oldPiece := ct[y*8+x]
	ct[y*8+x] = nil
	return oldPiece
}

func (ct *ChessTable) GetPosition(X rune, Y int) *ChessPiece {
	x, y := MustPositionToIndex(X, Y)
	return ct[y*8+x]
}

func (ct *ChessTable) GetIndex(x int, y int) *ChessPiece {
	return ct[y*8+x]
}

func (ct *ChessTable) Copy() *ChessTable {
	var table ChessTable
	for i := 0; i < 64; i++ {
		if ct[i] != nil {
			newPiece := &ChessPiece{
				PieceType:            ct[i].PieceType,
				X:                    ct[i].X,
				Y:                    ct[i].Y,
				GameSide:             ct[i].GameSide,
				Moved:                ct[i].Moved,
				PawnMovedTwoLastTime: ct[i].PawnMovedTwoLastTime,
			}

			table[i] = newPiece
		}
	}

	return &table
}

// 测试upgrade
func NewTestTable1() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'a', Y: 7, PieceType: ChessPieceTypeQueen, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})

	return &table
}

// 之前模拟的一个局面, 用它测试出了bug, 现在已经修复
func NewTestTable2() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'a', Y: 8, PieceType: ChessPieceTypeQueen, GameSide: SideWhite, Moved: true})

	table.SetPosition(&ChessPiece{X: 'h', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true})

	table.SetPosition(&ChessPiece{X: 'a', Y: 3, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true})
	table.SetPosition(&ChessPiece{X: 'c', Y: 3, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true})
	table.SetPosition(&ChessPiece{X: 'e', Y: 3, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true})
	table.SetPosition(&ChessPiece{X: 'f', Y: 3, PieceType: ChessPieceTypeKnight, GameSide: SideWhite, Moved: true})
	table.SetPosition(&ChessPiece{X: 'g', Y: 3, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true})

	table.SetPosition(&ChessPiece{X: 'd', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'c', Y: 1, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'd', Y: 1, PieceType: ChessPieceTypeQueen, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 1, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'b', Y: 8, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'd', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})

	table.SetPosition(&ChessPiece{X: 'a', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 7, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: true})

	table.SetPosition(&ChessPiece{X: 'c', Y: 6, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'e', Y: 6, PieceType: ChessPieceTypeBishop, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'g', Y: 6, PieceType: ChessPieceTypeQueen, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'h', Y: 6, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: true})

	table.SetPosition(&ChessPiece{X: 'd', Y: 5, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'f', Y: 5, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'h', Y: 5, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true})

	table.SetPosition(&ChessPiece{X: 'g', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true})

	return &table
}

// 用来测试平局
func NewTestTable3() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'e', Y: 7, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: true})
	table.SetPosition(&ChessPiece{X: 'g', Y: 6, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: true})

	table.SetPosition(&ChessPiece{X: 'h', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})

	return &table
}

// 王车易位
func NewTestTable4() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'a', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'g', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})

	return &table
}

// 王车易位, 目的地对方有车
func NewTestTable5() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'a', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'g', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'c', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: true})

	return &table
}

// 王车易位, 中途截胡
func NewTestTable6() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'a', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'g', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'f', Y: 3, PieceType: ChessPieceTypeQueen, GameSide: SideBlack, Moved: true})

	return &table
}

func NewTestTable7() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'a', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'g', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'f', Y: 3, PieceType: ChessPieceTypeQueen, GameSide: SideBlack, Moved: true})

	return &table
}

// 王车易位, 但是将军
func NewTestTable8() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'a', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'a', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'g', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'f', Y: 3, PieceType: ChessPieceTypeQueen, GameSide: SideBlack, Moved: true})

	return &table
}

// 胜利
func NewTestTable9() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'f', Y: 7, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'h', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})

	return &table
}

// 平局
func NewTestTable10() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'e', Y: 2, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: true})
	table.SetPosition(&ChessPiece{X: 'd', Y: 5, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: true})
	table.SetPosition(&ChessPiece{X: 'h', Y: 5, PieceType: ChessPieceTypeQueen, GameSide: SideWhite, Moved: true})
	table.SetPosition(&ChessPiece{X: 'a', Y: 8, PieceType: ChessPieceTypeQueen, GameSide: SideWhite, Moved: true})

	table.SetPosition(&ChessPiece{X: 'e', Y: 4, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})

	return &table
}

// 测bug的棋盘
func NewTestTable11() *ChessTable {
	var table ChessTable
	table.SetPosition(&ChessPiece{X: 'a', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 1, PieceType: ChessPieceTypeKnight, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 1, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'd', Y: 1, PieceType: ChessPieceTypeQueen, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 1, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 1, PieceType: ChessPieceTypeKnight, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'd', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'b', Y: 3, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 8, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 8, PieceType: ChessPieceTypeBishop, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 8, PieceType: ChessPieceTypeBishop, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 8, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: false})

	table.SetPosition(&ChessPiece{X: 'b', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'd', Y: 6, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 5, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 5, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 5, PieceType: ChessPieceTypeQueen, GameSide: SideBlack, Moved: true})

	return &table
}

// 测bug的棋盘
func NewTestTable12() *ChessTable {
	var table ChessTable
	table.SetPosition(&ChessPiece{X: 'b', Y: 1, PieceType: ChessPieceTypeKnight, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 1, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 1, PieceType: ChessPieceTypeKnight, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'b', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 3, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: true})
	table.SetPosition(&ChessPiece{X: 'f', Y: 3, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'd', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'd', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'e', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 8, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'd', Y: 7, PieceType: ChessPieceTypeBishop, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'f', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'g', Y: 6, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'b', Y: 5, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'b', Y: 4, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'f', Y: 4, PieceType: ChessPieceTypeQueen, GameSide: SideBlack, Moved: true})

	return &table
}

// 测bug的棋盘
func NewTestTable13() *ChessTable {
	var table ChessTable
	table.SetPosition(&ChessPiece{X: 'a', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 1, PieceType: ChessPieceTypeKnight, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 1, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 1, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 1, PieceType: ChessPieceTypeKnight, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'e', Y: 3, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'c', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'd', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'b', Y: 8, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'e', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'd', Y: 7, PieceType: ChessPieceTypeBishop, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'e', Y: 7, PieceType: ChessPieceTypeBishop, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'f', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'f', Y: 6, PieceType: ChessPieceTypeQueen, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'h', Y: 6, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: true})

	table.SetPosition(&ChessPiece{X: 'e', Y: 5, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 5, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true, PawnMovedTwoLastTime: false})
	return &table
}

func NewTestTable14() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 5, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: true})
	table.SetPosition(&ChessPiece{X: 'c', Y: 6, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: true})
	table.SetPosition(&ChessPiece{X: 'c', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'd', Y: 7, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'b', Y: 5, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: true})

	return &table
}

func NewTestTable15() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'a', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 1, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 1, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'd', Y: 2, PieceType: ChessPieceTypeKnight, GameSide: SideWhite, Moved: true})
	table.SetPosition(&ChessPiece{X: 'g', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'e', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'c', Y: 5, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 5, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 8, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'f', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'g', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})

	table.SetPosition(&ChessPiece{X: 'a', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'h', Y: 6, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: true})

	table.SetPosition(&ChessPiece{X: 'g', Y: 4, PieceType: ChessPieceTypeBishop, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'e', Y: 3, PieceType: ChessPieceTypeQueen, GameSide: SideBlack, Moved: true})

	table.SetPosition(&ChessPiece{X: 'd', Y: 7, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'b', Y: 5, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: true})

	return &table
}

func NewTestTable16() *ChessTable {
	var table ChessTable

	table.SetPosition(&ChessPiece{X: 'f', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 2, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: true})

	table.SetPosition(&ChessPiece{X: 'b', Y: 3, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'd', Y: 3, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'f', Y: 4, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: true, PawnMovedTwoLastTime: false})

	// ---

	table.SetPosition(&ChessPiece{X: 'c', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'f', Y: 8, PieceType: ChessPieceTypeBishop, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 7, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'f', Y: 7, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: true})
	table.SetPosition(&ChessPiece{X: 'g', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'c', Y: 5, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: true, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'g', Y: 3, PieceType: ChessPieceTypeQueen, GameSide: SideBlack, Moved: true})

	return &table
}

func NewChessTable() *ChessTable {
	var table ChessTable
	table.SetPosition(&ChessPiece{X: 'a', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 1, PieceType: ChessPieceTypeKnight, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 1, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'd', Y: 1, PieceType: ChessPieceTypeQueen, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 1, PieceType: ChessPieceTypeKing, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 1, PieceType: ChessPieceTypeBishop, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 1, PieceType: ChessPieceTypeKnight, GameSide: SideWhite, Moved: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 1, PieceType: ChessPieceTypeRook, GameSide: SideWhite, Moved: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'd', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 2, PieceType: ChessPieceTypePawn, GameSide: SideWhite, Moved: false, PawnMovedTwoLastTime: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 8, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 8, PieceType: ChessPieceTypeBishop, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'd', Y: 8, PieceType: ChessPieceTypeQueen, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 8, PieceType: ChessPieceTypeKing, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 8, PieceType: ChessPieceTypeBishop, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 8, PieceType: ChessPieceTypeKnight, GameSide: SideBlack, Moved: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 8, PieceType: ChessPieceTypeRook, GameSide: SideBlack, Moved: false})

	table.SetPosition(&ChessPiece{X: 'a', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'b', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'c', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'd', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'e', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'f', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'g', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})
	table.SetPosition(&ChessPiece{X: 'h', Y: 7, PieceType: ChessPieceTypePawn, GameSide: SideBlack, Moved: false, PawnMovedTwoLastTime: false})

	return &table
}

func MustPositionToIndex(X rune, Y int) (int, int) {
	var x int
	switch X {
	case 'a':
		x = 0
	case 'b':
		x = 1
	case 'c':
		x = 2
	case 'd':
		x = 3
	case 'e':
		x = 4
	case 'f':
		x = 5
	case 'g':
		x = 6
	case 'h':
		x = 7
	default:
		panic("unreachable")
	}

	var y = 0
	if Y < 1 || Y > 8 {
		panic("unreachable")
	}
	y = Y - 1

	return x, y
}

func MustIndexToPosition(x int, y int) (rune, int) {
	var X rune
	switch x {
	case 0:
		X = 'a'
	case 1:
		X = 'b'
	case 2:
		X = 'c'
	case 3:
		X = 'd'
	case 4:
		X = 'e'
	case 5:
		X = 'f'
	case 6:
		X = 'g'
	case 7:
		X = 'h'
	default:
		panic("unreachable")
	}

	if y < 0 || y > 7 {
		panic("unreachable")
	}
	return X, y + 1
}
