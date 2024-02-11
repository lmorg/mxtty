package virtualterm

import (
	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
)

/*
	ERASE DISPLAY
*/

func (term *Term) csiEraseDisplayAfter() {
	debug.Log(term.curPos)

	for y := term.curPos.Y; y < term.size.Y; y++ {
		for x := int32(0); x < term.size.X; x++ {
			(*term.cells)[y][x].Clear()
		}
	}
}

func (term *Term) csiEraseDisplayBefore() {
	debug.Log(term.curPos)

	for y := term.curPos.Y; y >= 0; y-- {
		for x := int32(0); x >= 0; x-- {
			(*term.cells)[y][x].Clear()
		}
	}
}

func (term *Term) csiEraseDisplay() {
	debug.Log(term.curPos)

	var x, y int32
	for ; y < term.size.Y; y++ {
		for x = 0; x < term.size.X; x++ {
			(*term.cells)[y][x].Clear()
		}
	}
}

/*
	ERASE LINE
*/

func (term *Term) csiEraseLineAfter() {
	debug.Log(term.curPos)

	for x := term.curPos.X; x < term.size.X; x++ {
		(*term.cells)[term.curPos.Y][x].Clear()
	}
}

func (term *Term) csiEraseLineBefore() {
	debug.Log(term.curPos)

	for x := term.curPos.X; x >= 0; x-- {
		(*term.cells)[term.curPos.Y][x].Clear()
	}
}

func (term *Term) csiEraseLine() {
	debug.Log(term.curPos)

	var x int32
	for ; x < term.size.X; x++ {
		(*term.cells)[term.curPos.Y][x].Clear()
	}
}

/*
	ERASE CHARACTERS
*/

func (term *Term) csiEraseCharacters(n int32) {
	debug.Log(n)

	if n < 1 {
		n = 1
	}

	clear := make([]types.Cell, n)
	copy((*term.cells)[term.curPos.Y][term.curPos.X:], clear)
}

/*
	DELETE
*/

func (term *Term) csiDeleteCharacters(n int32) {
	debug.Log(n)

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
	debug.Log(n)

	if n < 1 {
		n = 1
	}

	_, bottom := term.getScrollingRegion()

	if term.curPos.Y+n >= bottom {
		n = bottom - term.curPos.Y
	}

	for i := int32(0); i < bottom-term.curPos.Y; i++ {
		if term.curPos.Y+i+n <= bottom {
			(*term.cells)[term.curPos.Y+i] = (*term.cells)[term.curPos.Y]
		} else {
			(*term.cells)[term.curPos.Y] = term.makeRow()
		}
	}
}
