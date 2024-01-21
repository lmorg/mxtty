package virtualterm

import "log"

/*
	ERASE DISPLAY
*/

func (term *Term) csiEraseDisplayAfter() {
	log.Println("DEBUG: csiEraseDisplayAfter()")

	for y := term.curPos.Y; y < term.size.Y; y++ {
		for x := int32(0); x < term.size.X; x++ {
			(*term.cells)[y][x].Clear()
		}
	}
}

func (term *Term) csiEraseDisplayBefore() {
	log.Println("DEBUG: csiEraseDisplayBefore()")

	for y := term.curPos.Y; y >= 0; y-- {
		for x := int32(0); x >= 0; x-- {
			(*term.cells)[y][x].Clear()
		}
	}
}

func (term *Term) csiEraseDisplay() {
	log.Println("DEBUG: csiEraseDisplay()")

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
	log.Println("DEBUG: csiEraseLineAfter()")

	for x := term.curPos.X; x < term.size.X; x++ {
		(*term.cells)[term.curPos.Y][x].Clear()
	}
}

func (term *Term) csiEraseLineBefore() {
	log.Println("DEBUG: csiEraseLineBefore()")

	for x := term.curPos.X; x >= 0; x-- {
		(*term.cells)[term.curPos.Y][x].Clear()
	}
}

func (term *Term) csiEraseLine() {
	log.Println("DEBUG: csiEraseLine()")

	var x int32
	for ; x < term.size.X; x++ {
		(*term.cells)[term.curPos.Y][x].Clear()
	}
}

func (term *Term) csiEraseCharacters(n int32) {
	log.Println("DEBUG: csiEraseCharacters()")

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

/*
	DELETE
*/

func (term *Term) csiDeleteCharacters(n int32) {
	log.Println("DEBUG: csiDeleteCharacters()")

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
	log.Println("DEBUG: csiDeleteLines()")

	if n < 1 {
		n = 1
	}

	_, bottom := term.getScrollRegion()

	if term.curPos.Y+n >= bottom {
		n = bottom - term.curPos.Y
	}

	for i := int32(0); i < bottom-term.curPos.Y; i++ {
		if term.curPos.Y+i+n <= bottom {
			(*term.cells)[term.curPos.Y+i] = (*term.cells)[term.curPos.Y]
		} else {
			(*term.cells)[term.curPos.Y] = term.newRow()
		}
	}
}
