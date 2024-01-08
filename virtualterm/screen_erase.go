package virtualterm

/*
	ERASE DISPLAY
*/

func (term *Term) eraseDisplayAfter() {
	for y := term.curPos.Y; y < term.size.Y; y++ {
		for x := term.curPos.X; x < term.size.X; x++ {
			(*term.cells)[y][x].clear()
		}
	}
}

func (term *Term) eraseDisplayBefore() {
	for y := term.curPos.Y; y >= 0; y-- {
		for x := term.curPos.X; x >= 0; x-- {
			(*term.cells)[y][x].clear()
		}
	}
}

func (term *Term) eraseDisplay() {
	var x, y int32
	for ; y < term.size.Y; y++ {
		for ; x < term.size.X; x++ {
			(*term.cells)[y][x].clear()
		}
	}
}

/*
	ERASE LINE
*/

func (term *Term) eraseLineAfter() {
	for x := term.curPos.X; x < term.size.X; x++ {
		(*term.cells)[term.curPos.Y][x].clear()
	}
}

func (term *Term) eraseLineBefore() {
	for x := term.curPos.X; x >= 0; x-- {
		(*term.cells)[term.curPos.Y][x].clear()
	}
}

func (term *Term) eraseLine() {
	var x int32
	for ; x < term.size.X; x++ {
		(*term.cells)[term.curPos.Y][x].clear()
	}
}

/*
	DELETE
*/

func (term *Term) deleteCharacters(n int32) {
	if n < 1 {
		n = 1
	}

	if term.curPos.X+n >= term.size.X {
		n = term.size.X - term.curPos.X
	}

	for i := int32(0); i < term.size.X-term.curPos.X; i++ {
		if term.curPos.X+i+n < term.size.X {
			(*term.cells)[term.curPos.Y][term.curPos.X+i] = (*term.cells)[term.curPos.Y][term.curPos.X+i+n]
		} else {
			(*term.cells)[term.curPos.Y][term.curPos.X+i].clear()
		}
	}
}
