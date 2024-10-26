package virtualterm

import (
	"fmt"
	"log"

	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/debug"
)

func (term *Term) csiRepeatPreceding(n int32) {
	debug.Log(n)

	if n < 1 {
		n = 1
	}
	cell, _ := term.previousCell()
	for i := int32(0); i < n; i++ {
		term.writeCell(cell.Char)
	}
}

func (term *Term) csiCallback(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	err := term.Pty.Write([]byte(codes.Csi + msg))
	if err != nil {
		log.Printf("ERROR: writing callback message '%s': %s", msg, err.Error())
	}
}

/*
	SCREEN BUFFER
*/

func (term *Term) csiScreenBufferAlternative() {
	term.cells = &term._altBuf
}

func (term *Term) csiScreenBufferNormal() {
	term.cells = &term._normBuf
	for i := range term._altBuf {
		term._altBuf[i] = term.makeRow()
	}
}

/*
	CURSOR
*/

func (term *Term) csiCursorPosSave() {
	debug.Log(term._savedCurPos)
	debug.Log(term.curPos)
	term._savedCurPos = term._curPos
}

func (term *Term) csiCursorPosRestore() {
	debug.Log(term.curPos)
	debug.Log(term._savedCurPos)
	term._curPos = term._savedCurPos
}

func (term *Term) csiCursorHide() {
	term.ShowCursor(false)
}

func (term *Term) csiCursorShow() {
	term.ShowCursor(true)
}

/*
	WINDOW TITLE
*/

func (term *Term) csiWindowTitleStackSaveTo() {
	term._windowTitleStack = append(term._windowTitleStack, term.renderer.GetWindowTitle())
}

func (term *Term) csiWindowTitleStackRestoreFrom() {
	title := term._windowTitleStack[len(term._windowTitleStack)-1]
	term.renderer.SetWindowTitle(title)
	term._windowTitleStack = term._windowTitleStack[:len(term._windowTitleStack)-1]
}

/*
	MISC
*/

func (term *Term) csiNoAutoLineWrap(state bool) {
	term._noAutoLineWrap = state
}

func (term *Term) csiIrmInsertOrReplace(state _stateIrmT) {
	term._insertOrReplace = state
}
