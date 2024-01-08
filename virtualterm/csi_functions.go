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

func (term *Term) csiSetScrollingRegion(region []int32) {
	term._scrollRegion = &scrollRegionT{
		Top:    region[0],
		Bottom: region[1],
	}
}

func (term *Term) csiWindowTitleStackSaveTo() {
	term._windowTitleStack = append(term._windowTitleStack, term.renderer.GetWindowTitle())
}

func (term *Term) csiWindowTitleStackRestoreFrom() {
	title := term._windowTitleStack[len(term._windowTitleStack)-1]
	term.renderer.SetWindowTitle(title)
	term._windowTitleStack = term._windowTitleStack[:len(term._windowTitleStack)-1]
}
