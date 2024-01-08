package virtualterm

func (term *Term) csiCursorPosSave() {
	term._savedCurPos = term.curPos
}

func (term *Term) csiCursorPosRestore() {
	term.curPos = term._savedCurPos
}

func (term *Term) csiScreenBufferAlternative() {
	term.cells = &term._altBuf
}

func (term *Term) csiScreenBufferNormal() {
	term.cells = &term._normBuf
}

func (term *Term) csiCursorHide() {
	term._hideCursor = true
}

func (term *Term) csiCursorShow() {
	term._hideCursor = false
}
