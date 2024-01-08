package virtualterm

// moveCursor functions DON'T affect other contents in the grid

func (term *Term) moveCursorBackwards(i int32) (overflow int32) {
	if i < 0 {
		i = 1
	}

	term.curPos.X -= i
	if term.curPos.X < 0 {
		overflow = term.curPos.X * -1
		term.curPos.X = 0
	}

	return
}

func (term *Term) moveCursorForwards(i int32) (overflow int32) {
	if i < 0 {
		i = 1
	}

	term.curPos.X += i
	if term.curPos.X >= term.size.X {
		overflow = term.curPos.X - (term.size.X - 1)
		term.curPos.X = term.size.X - 1
	}

	return
}

func (term *Term) moveCursorUpwards(i int32) (overflow int32) {
	if i < 0 {
		i = 1
	}

	term.curPos.Y -= i
	if term.curPos.Y < 0 {
		overflow = term.curPos.Y * -1
		term.curPos.Y = 0
	}

	return
}

func (term *Term) moveCursorDownwards(i int32) (overflow int32) {
	if i < 0 {
		i = 1
	}

	term.curPos.Y += i
	if term.curPos.Y >= term.size.Y {
		overflow = term.curPos.Y - (term.size.Y - 1)
		term.curPos.Y = term.size.Y - 1
	}

	return
}

// moveGridPos functions DO affect other contents in the grid

func (term *Term) moveContentsUp() {
	var i int32
	for ; i < term.size.Y-1; i++ {
		term.cells[i] = term.cells[i+1]
	}
	term.cells[i] = make([]cell, term.size.X, term.size.X)
}

func (term *Term) wrapCursorForwards() {
	term.curPos.X += 1

	if term.curPos.X > term.size.X {
		overflow := term.curPos.X - (term.size.X - 1)
		term.curPos.X = 0

		if overflow > 0 && term.moveCursorDownwards(1) > 0 {
			term.moveContentsUp()
			term.moveCursorDownwards(1)
		}
	}
}
