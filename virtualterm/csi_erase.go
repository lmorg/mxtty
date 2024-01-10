package virtualterm

/*
	ERASE DISPLAY
*/

func (term *Term) csiEraseDisplayAfter() {
	for y := term.curPos.Y; y < term.size.Y; y++ {
		for x := term.curPos.X; x < term.size.X; x++ {
			(*term.cells)[y][x].Clear()
		}
	}
}

func (term *Term) csiEraseDisplayBefore() {
	for y := term.curPos.Y; y >= 0; y-- {
		for x := term.curPos.X; x >= 0; x-- {
			(*term.cells)[y][x].Clear()
		}
	}
}

func (term *Term) csiEraseDisplay() {
	var x, y int32
	for ; y < term.size.Y; y++ {
		for ; x < term.size.X; x++ {
			(*term.cells)[y][x].Clear()
		}
	}
}

/*
	ERASE LINE
*/

func (term *Term) csiEraseLineAfter() {
	for x := term.curPos.X; x < term.size.X; x++ {
		(*term.cells)[term.curPos.Y][x].Clear()
	}
}

func (term *Term) csiEraseLineBefore() {
	for x := term.curPos.X; x >= 0; x-- {
		(*term.cells)[term.curPos.Y][x].Clear()
	}
}

func (term *Term) csiEraseLine() {
	var x int32
	for ; x < term.size.X; x++ {
		(*term.cells)[term.curPos.Y][x].Clear()
	}
}

func (term *Term) csiEraseCharacters(n int32) {
	if n < 1 {
		n = 1
	}
	if term.curPos.X+n >= term.size.X {
		n = term.size.X - term.curPos.X
	}
	for x := int32(0); x < n; x++ {
		(*term.cells)[term.curPos.Y][x].Clear()
	}
}

func (term *Term) csiClearTab() {
	// TODO: this wouldn't actually work!
	term.csiEraseCharacters(term._tabWidth)
}

/*
	DELETE
*/

func (term *Term) csiDeleteCharacters(n int32) {
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
			(*term.cells)[term.curPos.Y][term.curPos.X+i].Clear()
		}
	}
}

func (term *Term) csiDeleteLines(n int32) {
	if n < 1 {
		n = 1
	}

	if term.curPos.Y+n >= term.size.Y {
		n = term.size.Y - term.curPos.Y
	}

	for i := int32(0); i < term.size.Y-term.curPos.Y; i++ {
		if term.curPos.Y+i+n < term.size.Y {
			(*term.cells)[term.curPos.Y+i] = (*term.cells)[term.curPos.Y]
		} else {
			(*term.cells)[term.curPos.Y] = term.newRow()
		}
	}
}
