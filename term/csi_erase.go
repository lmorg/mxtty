package virtualterm

import (
	"fmt"

	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
)

/*
	ERASE DISPLAY
*/

func (term *Term) csiEraseDisplayAfter() {
	debug.Log(term.curPos())

	_, bottom := term.getScrollingRegionExcOrigin()

	curPosY := term.curPos().Y
	if curPosY > bottom {
		debug.Log(fmt.Sprintf("curPos().Y[%d]>bottom[%d]", curPosY, bottom))
		return
	}

	for y := curPosY + 1; y <= bottom; y++ {
		(*term.cells)[y] = term.makeRow()
	}
	term.csiEraseLineAfter()
}

func (term *Term) csiEraseDisplayBefore() {
	debug.Log(term.curPos())

	top, _ := term.getScrollingRegionExcOrigin()
	for y := term.curPos().Y - 1; y >= top; y-- {
		(*term.cells)[y] = term.makeRow()
	}
	term.csiEraseLineBefore()
}

func (term *Term) csiEraseDisplay() {
	debug.Log(term.curPos())

	y, bottom := term.getScrollingRegionExcOrigin()
	for ; y <= bottom; y++ {
		(*term.cells)[y] = term.makeRow()
	}
}

func (term *Term) eraseScrollBack() {
	debug.Log(term.curPos())

	term._scrollOffset = 0
	term._scrollMsg = nil
	term._scrollBuf = [][]types.Cell{}
}

/*
	ERASE LINE
*/

func (term *Term) csiEraseLineAfter() {
	debug.Log(term.curPos())

	pos := term.curPos()
	n := term.size.X - pos.X
	clear := make([]types.Cell, n)
	copy((*term.cells)[pos.Y][pos.X:], clear)
}

func (term *Term) csiEraseLineBefore() {
	debug.Log(term.curPos())

	pos := term.curPos()
	n := pos.X + 1
	clear := make([]types.Cell, n)
	copy((*term.cells)[pos.Y], clear)
}

func (term *Term) csiEraseLine() {
	debug.Log(term.curPos())

	(*term.cells)[term.curPos().Y] = term.makeRow()
}

/*
	ERASE CHARACTERS
*/

func (term *Term) csiEraseCharacters(n int32) {
	debug.Log(n)

	if n < 1 {
		n = 1
	}

	pos := term.curPos()
	clear := make([]types.Cell, n)
	copy((*term.cells)[pos.Y][pos.X:], clear)
}

/*
	DELETE
*/

func (term *Term) csiDeleteCharacters(n int32) {
	debug.Log(n)

	if n < 1 {
		n = 1
	}

	pos := term.curPos()
	copy((*term.cells)[pos.Y][pos.X:], (*term.cells)[pos.Y][pos.X+n:])
	blank := make([]types.Cell, n)
	copy((*term.cells)[pos.Y][term.size.X-n:], blank)

}

func (term *Term) csiDeleteLines(n int32) {
	debug.Log(n)

	if n < 1 {
		n = 1
	}

	_, bottom := term.getScrollingRegionExcOrigin()

	term._scrollUp(term.curPos().Y, bottom, n)
}
