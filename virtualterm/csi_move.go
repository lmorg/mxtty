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
	//debug.Log(term.curPos.Y)

	if term.csiMoveCursorDownwards(1) != 0 {
		term.appendScrollBuf()
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

	if term.csiMoveCursorUpwards(1) != 0 {
		term.csiScrollDown(1)
		term.csiMoveCursorUpwards(1)
	}
}

/*
	csiMoveCursor[...] functions DOESN'T affect other contents in the grid
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
		overflow = -term.curPos.X
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
		overflow = -term.curPos.Y
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

	term._scrollUp(top, bottom, n)
}

// _scrollDown does not take into account term size nor scrolling region. Any
// error handling should be done by the calling function.
func (term *Term) _scrollUp(top, bottom, shift int32) {
	for i := top; i <= bottom; i++ {
		if i+shift <= bottom {
			(*term.cells)[i] = (*term.cells)[i+shift]
		} else {
			(*term.cells)[i] = term.makeRow()
		}
	}
}

// csiScrollDown: -1 should default to 1.
func (term *Term) csiScrollDown(n int32) {
	debug.Log(n)

	if n < 1 {
		n = 1
	}

	top, bottom := term.getScrollingRegion()

	term._scrollDown(top, bottom, n)
}

// _scrollDown does not take into account term size nor scrolling region. Any
// error handling should be done by the calling function.
func (term *Term) _scrollDown(top, bottom, shift int32) {
	screen := term.makeScreen()

	if top+shift > bottom {
		copy((*term.cells)[top:bottom+1], screen)
		return
	}

	copy(screen[top+shift:], (*term.cells)[top:bottom+1])
	copy((*term.cells)[top:], screen[top:bottom+1])
}

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

	term._scrollDown(term.curPos.Y, bottom, n)
}

// csiInsertCharacters: -1 should default to 1.
func (term *Term) csiInsertCharacters(n int32) {
	debug.Log(n)

	if n < 1 {
		n = 1
	}

	insert := make([]types.Cell, n)
	row := term.makeRow()

	copy(row, (*term.cells)[term.curPos.Y][:term.curPos.X])
	copy(row[term.curPos.X:], insert)
	copy(row[term.curPos.X+n:], (*term.cells)[term.curPos.Y][term.curPos.X:])

	(*term.cells)[term.curPos.Y] = row
}
