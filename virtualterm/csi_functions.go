package virtualterm

func (term *Term) _csiCursorPosSave() {
	term.savedCurPos = term.curPos
}

func (term *Term) _csiCursorPosRestore() {
	term.curPos = term.savedCurPos
}

func (term *Term) _csiScreenBufferAlternative() {
	term.cells = &term.altBuf
}

func (term *Term) _csiScreenBufferNormal() {
	term.cells = &term.normBuf
}
