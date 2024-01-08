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

	//log.Printf("DEBUG: moveCursorBackwards(%d) == %d [pos: %d]", i, overflow, term.curPos.X)

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

	//log.Printf("DEBUG: moveCursorForwards(%d) == %d [pos: %d]", i, overflow, term.curPos.X)

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

	//log.Printf("DEBUG: moveCursorUpwards(%d) == %d [pos: %d]", i, overflow, term.curPos.Y)

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

	//log.Printf("DEBUG: moveCursorDownwards(%d) == %d [pos: %d]", i, overflow, term.curPos.Y)

	return
}

func (term *Term) scrollUp(n int32) {
	if n < 1 {
		n = 1
	}

	top, bottom := term.getScrollRegion()

	if n > bottom-top {
		n = bottom - top
	}

	for i := top; i <= bottom; i++ {
		if i+n <= bottom {
			(*term.cells)[i] = (*term.cells)[i+n]
		} else {
			(*term.cells)[i] = make([]cell, term.size.X)
		}
	}
}

func (term *Term) scrollDown(n int32) {
	if n < 0 {
		n = 1
	}

	top, bottom := term.getScrollRegion()

	if n > bottom-top {
		n = bottom - top
	}

	for i := bottom; i > top; i-- {
		if i+n <= bottom {
			(*term.cells)[i] = (*term.cells)[i-n]
		} else {
			(*term.cells)[i] = make([]cell, term.size.X)
		}
	}
}

func (term *Term) moveCursorToPos(x, y int32) {
	if x < 0 {
		x = 0 //term.curPos.X
	} else if x >= term.size.X {
		x = term.size.X - 1
	}

	if y < 0 {
		y = 0 //term.curPos.Y
	} else if y >= term.size.Y {
		y = term.size.Y - 1
	}

	term.curPos.X, term.curPos.Y = x, y
}
