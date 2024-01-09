package virtualterm

// basic TTY operations

func (term *Term) printTab() {
	indent := int(4 - (term.curPos.X % term._tabWidth))
	for i := 0; i < indent; i++ {
		term.writeCell(' ')
	}
}

func (term *Term) carriageReturn() {
	term.curPos.X = 0
}

func (term *Term) lineFeed() {
	if term.csiMoveCursorDownwards(1) > 0 {
		term.csiScrollUp(1)
		term.csiMoveCursorDownwards(1)
	}

	go term.lfRedraw()
}

func (term *Term) ReverseLineFeed() {
	if term.csiMoveCursorUpwards(1) > 0 {
		term.csiScrollDown(1)
		term.csiMoveCursorUpwards(1)
	}
}

// moveCursor functions DON'T affect other contents in the grid

func (term *Term) csiMoveCursorBackwards(i int32) (overflow int32) {
	if i < 0 {
		i = 1
	}

	term.curPos.X -= i
	if term.curPos.X < 0 {
		overflow = term.curPos.X * -1
		term.curPos.X = 0
	}

	//log.Printf("DEBUG: csiMoveCursorBackwards(%d) == %d [pos: %d]", i, overflow, term.curPos.X)

	return
}

func (term *Term) csiMoveCursorForwards(i int32) (overflow int32) {
	if i < 0 {
		i = 1
	}

	term.curPos.X += i
	if term.curPos.X >= term.size.X {
		overflow = term.curPos.X - (term.size.X - 1)
		term.curPos.X = term.size.X - 1
	}

	//log.Printf("DEBUG: csiMoveCursorForwards(%d) == %d [pos: %d]", i, overflow, term.curPos.X)

	return
}

func (term *Term) csiMoveCursorUpwards(i int32) (overflow int32) {
	if i < 0 {
		i = 1
	}

	top, _ := term.getScrollRegion()

	term.curPos.Y -= i
	if term.curPos.Y <= top {
		overflow = term.curPos.Y * -1
		term.curPos.Y = top
	}

	//log.Printf("DEBUG: csiMoveCursorUpwards(%d) == %d [pos: %d]", i, overflow, term.curPos.Y)

	return
}

func (term *Term) csiMoveCursorDownwards(i int32) (overflow int32) {
	if i < 0 {
		i = 1
	}

	term.curPos.Y += i

	_, bottom := term.getScrollRegion()

	if term.curPos.Y > bottom {
		overflow = term.curPos.Y - (bottom)
		term.curPos.Y = bottom
	}

	//log.Printf("DEBUG: csiMoveCursorDownwards(%d) == %d [pos: %d]", i, overflow, term.curPos.Y)

	return
}

func (term *Term) csiMoveCursorToPos(x, y int32) {
	if x < 0 {
		x = term.curPos.X
	} else if x >= term.size.X {
		x = term.size.X - 1
	}

	if y < 0 {
		y = term.curPos.Y
	} else if y >= term.size.Y {
		y = term.size.Y - 1
	}

	term.curPos.X, term.curPos.Y = x, y
}

/*
	SCROLLING
*/

func (term *Term) csiScrollUp(n int32) {
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

func (term *Term) csiScrollDown(n int32) {
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

func (term *Term) csiSetScrollingRegion(region []int32) {
	term._scrollRegion = &scrollRegionT{
		Top:    region[0],
		Bottom: region[1],
	}
}

/*
	INSERTING
*/

func (term *Term) csiInsertLines(n int32) {
	if n < 1 {
		n = 1
	}

	_, bottom := term.getScrollRegion()

	if term.curPos.Y+n > bottom {
		n = bottom - term.curPos.Y
	}

	for y, i := bottom, int32(0); y > term.curPos.Y; y-- {
		if i < n {
			(*term.cells)[y] = (*term.cells)[y-1]
			i++
		} else {
			(*term.cells)[y] = make([]cell, term.size.X)
		}
	}
}
