package virtualterm

import (
	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
)

// basic TTY operations

func (term *Term) carriageReturn() {
	term.curPos.X = 0
}

func (term *Term) lineFeed() {
	debug.Log(term.curPos.Y)

	if term.csiMoveCursorDownwards(1) > 0 {
		term.csiScrollUp(1)
		term.csiMoveCursorDownwards(1)
	}

	if term._activeElement != nil {
		term._activeElement.ReadCell(nil)
	}

	go term.lfRedraw()
}

func (term *Term) ReverseLineFeed() {
	debug.Log(term.curPos.Y)

	if term.csiMoveCursorUpwards(1) > 0 {
		term.csiScrollDown(1)
		term.csiMoveCursorUpwards(1)
	}
}

/*
	csiMoveCursor[...] functions DON'T affect other contents in the grid
*/

// csiMoveCursorBackwards: -1 should default to 1.
// Returns how many additional cells were requested before hitting the edge of
// the screen.
func (term *Term) csiMoveCursorBackwards(i int32) (overflow int32) {
	debug.Log(i)

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

// csiMoveCursorForwards: -1 should default to 1.
// Returns how many additional cells were requested before hitting the edge of
// the screen.
func (term *Term) csiMoveCursorForwards(i int32) (overflow int32) {
	debug.Log(i)

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

// csiMoveCursorUpwards: -1 should default to 1.
// Returns how many additional cells were requested before hitting the edge of
// the screen.
func (term *Term) csiMoveCursorUpwards(i int32) (overflow int32) {
	debug.Log(i)

	if i < 0 {
		i = 1
	}

	top, _ := term.getScrollingRegion()

	term.curPos.Y -= i
	if term.curPos.Y <= top {
		overflow = term.curPos.Y * -1
		term.curPos.Y = top
	}

	return
}

// csiMoveCursorDownwards: -1 should default to 1.
// Returns how many additional cells were requested before hitting the edge of
// the screen.
func (term *Term) csiMoveCursorDownwards(i int32) (overflow int32) {
	debug.Log(i)

	if i < 0 {
		i = 1
	}

	term.curPos.Y += i

	_, bottom := term.getScrollingRegion()

	if term.curPos.Y > bottom {
		overflow = term.curPos.Y - (bottom)
		term.curPos.Y = bottom
	}

	return
}

// csiMoveCursorToPos: -1 values should default to current cursor position.
func (term *Term) csiMoveCursorToPos(x, y int32) {
	debug.Log(types.XY{X: x, Y: y})

	if x == -1 {
		x = term.curPos.X
	}
	if y == -1 {
		y = term.curPos.Y
	}

	top, bottom := term.getScrollingRegion()

	if x < 0 {
		x = 0
	} else if x >= term.size.X {
		x = term.size.X - 1
	}

	if y < top {
		y = top
	} else if y > bottom {
		y = bottom
	}

	term.curPos.X, term.curPos.Y = x, y
}

/*
	SCROLLING
*/

// csiSetScrollingRegion values should be offset by 1 (as seen in the ANSI
// escape codes). eg the top left corder would be `[]int32{1, 1}`.
func (term *Term) setScrollingRegion(region []int32) {
	debug.Log(region[0])
	debug.Log(region[1])

	term._scrollRegion = &scrollRegionT{
		Top:    region[0] - 1,
		Bottom: region[1] - 1,
	}
}

func (term *Term) getScrollingRegion() (top int32, bottom int32) {
	debug.Log(term._scrollRegion)

	if term._scrollRegion == nil || term._originMode {
		top = 0
		bottom = term.size.Y - 1
	} else {
		top = term._scrollRegion.Top
		bottom = term._scrollRegion.Bottom
	}

	return
}

// csiScrollUp: -1 should default to 1.
func (term *Term) csiScrollUp(n int32) {
	debug.Log(n)

	if n < 1 {
		n = 1
	}

	top, bottom := term.getScrollingRegion()

	if n > bottom-top {
		n = bottom - top
	}

	for i := top; i <= bottom; i++ {
		if i+n <= bottom {
			(*term.cells)[i] = (*term.cells)[i+n]
		} else {
			(*term.cells)[i] = term.makeRow()
		}
	}
}

// csiScrollUp: -1 should default to 1.
func (term *Term) csiScrollDown(n int32) {
	debug.Log(n)

	if n < 0 {
		n = 1
	}

	top, bottom := term.getScrollingRegion()

	/*if n > bottom-top {
		n = bottom - top
	}

	for i := bottom; i > top; i-- {
		if i+n <= bottom {
			(*term.cells)[i] = (*term.cells)[i-n]
		} else {
			(*term.cells)[i] = term.makeRow()
		}
	}*/
	term._scrollDown(top, bottom, 1)
}

// _scrollDown does not take into account term size nor scrolling region. Any
// error handling should be done by the calling function.
func (term *Term) _scrollDown(start, end, shift int32) {
	screen := term.makeScreen()

	if start+shift > end {
		shift = end - start
	}

	copy(screen[start+shift:end], (*term.cells)[start:end])
	copy((*term.cells)[start:end], screen[start:end])
}

// _scrollUp does not take into account term size nor scrolling region. Any
// error handling should be done by the calling function.
/*func (term *Term) _scrollUp(start, end, shift int32) {
	screen := term.makeScreen()

	if start-shift > start {
		shift = end - start
	}

	copy(screen[start-shift:end], (*term.cells)[start:end])
	copy((*term.cells)[start:end], screen[start:end])
}*/

/*
	INSERT
*/

// csiInsertLines: -1 should default to 1.
func (term *Term) csiInsertLines(n int32) {
	debug.Log(n)

	if n < 1 {
		n = 1
	}

	_, bottom := term.getScrollingRegion()

	/*start := term.curPos.Y + n
	if start > bottom {
		start = bottom - term.curPos.Y
	}*/

	/*screen := term.makeScreen()

	copy(screen[start:bottom], (*term.cells)[term.curPos.Y:bottom])
	copy((*term.cells)[term.curPos.Y:bottom], screen[term.curPos.Y:bottom])*/
	term._scrollDown(term.curPos.Y, bottom, n)
}

// csiInsertCharacters: -1 should default to 1.
func (term *Term) csiInsertCharacters(n int32) {
	debug.Log(n)

	if n < 1 {
		n = 1
	}

	insert := make([]types.Cell, n)
	for i := range insert {
		insert[i].Char = (*term.cells)[term.curPos.Y][term.curPos.X].Char
	}

	row := append((*term.cells)[term.curPos.Y][:term.curPos.X], insert...)
	row = append(row, (*term.cells)[term.curPos.Y][term.curPos.X:]...)
	(*term.cells)[term.curPos.Y] = row[:term.size.X]
}
