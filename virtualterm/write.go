package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/codes"
)

func (term *Term) writeCell(r rune) {
	if term.curPos.X >= term.size.X {
		overflow := term.curPos.X - (term.size.X - 1)
		term.curPos.X = 0

		if overflow > 0 && term.csiMoveCursorDownwards(1) > 0 {
			term.csiScrollUp(1)
			term.csiMoveCursorDownwards(1)
		}
	}

	term.cell().char = r
	term.cell().sgr = term.sgr.Copy()
	term.curPos.X++
}

// Write multiple characters to the virtual terminal
func (term *Term) printLoop() {
	var r rune

	for {
		r = term.Pty.Read()
		term._slowBlinkState = true
		//log.Printf("DEBUG: read rune %d [pos: %d:%d] [size: %d:%d]", r, term.curPos.X, term.curPos.Y, term.size.X, term.size.Y)

		term._mutex.Lock()
		term.readChar(r)
		term._mutex.Unlock()
	}
}

func (term *Term) readChar(r rune) {
	switch r {

	case codes.AsciiCtrlG: // bell (7)
		// TODO: beep

	case codes.AsciiBackspace, codes.IsoBackspace: // (10) / (127)
		_ = term.csiMoveCursorBackwards(1)

	case codes.AsciiTab: // \t (11)
		term.printTab()

	case codes.AsciiCtrlJ: // \n (12)
		term.lineFeed()

	case codes.AsciiCtrlM: // \r (13)
		term.carriageReturn()

	case codes.AsciiEscape: // (27)
		term.parseC1Codes()

	default:
		if r < 32 {
			log.Printf("Unexpected ASCII control character (ignored): %d", r)
			return
		}
		term.writeCell(r)
	}
}

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
}

func (term *Term) ReverseLineFeed() {
	if term.csiMoveCursorUpwards(1) > 0 {
		term.csiScrollDown(1)
		term.csiMoveCursorUpwards(1)
	}
}
