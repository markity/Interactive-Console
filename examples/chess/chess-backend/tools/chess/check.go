package chess

func CheckChessPostsionVaild(x rune, y int) bool {
	if x != 'a' && x != 'b' && x != 'c' && x != 'd' && x != 'e' && x != 'f' && x != 'g' && x != 'h' {
		return false
	}

	if y < 1 || y > 8 {
		return false
	}

	return true
}

func CheckChessIndexValid(x int, y int) bool {
	return x >= 0 && x <= 7 && y >= 0 && y <= 7
}
