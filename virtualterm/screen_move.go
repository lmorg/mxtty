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

	//log.Printf("DEBUG: top, bottom, n := %d, %d, %d", top, bottom, n)

	var i int32
	for i = top; i <= bottom-n; i++ {
		(*term.cells)[i] = (*term.cells)[i+1]
	}

	//log.Printf("DEBUG: top, bottom, n, i := %d, %d, %d, %d", top, bottom, n, i)
	for ; i <= bottom; i++ {
		(*term.cells)[i] = make([]cell, term.size.X)
		//(*term.cells)[i] = term._debug_FillRowWithDots()
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

	var i int32
	for i = bottom; i >= n; i-- {
		(*term.cells)[i] = (*term.cells)[i-1]
	}
	for ; i >= top; i-- {
		(*term.cells)[i] = make([]cell, term.size.X)
	}
}

func (term *Term) _debug_FillRowWithDots() []cell {
	row := make([]cell, term.size.X)
	for i := range row {
		row[i].char = '·'
		row[i].sgr = &sgr{
			fg: SGR_COLOUR_BLACK_BRIGHT,
			bg: SGR_DEFAULT.bg,
		}
	}

	return row
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
