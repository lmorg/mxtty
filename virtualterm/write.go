package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/codes"
)

func (term *Term) writeCell(r rune) {
	if term.curPos.X >= term.size.X {
		overflow := term.curPos.X - (term.size.X - 1)
		term.curPos.X = 0

		if overflow > 0 && term.moveCursorDownwards(1) > 0 {
			term.moveContentsUp()
			term.moveCursorDownwards(1)
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
		r = term.Pty.ReadRune()
		term._slowBlinkState = true
		//log.Printf("DEBUG: read rune %d [pos: %d:%d] [size: %d:%d]", r, term.curPos.X, term.curPos.Y, term.size.X, term.size.Y)

		term._mutex.Lock()

		switch r {

		case codes.AsciiEscape:
			term.parseC1Codes()

		case codes.AsciiBackspace, codes.IsoBackspace:
			_ = term.moveCursorBackwards(1)

		case codes.AsciiCtrlG: // bell
			// TODO: beep

		case '\t':
			indent := int(4 - (term.curPos.X % term.tabWidth))
			for i := 0; i < indent; i++ {
				term.writeCell(' ')
			}

		case '\r':
			term.curPos.X = 0

		case '\n':
			//log.Printf("DEBUG: new line char")
			if term.moveCursorDownwards(1) > 0 {
				term.moveContentsUp()
				term.moveCursorDownwards(1)
			}
			//term.wrapCursorForwards()
			//.term.curPos.X = 0

		//case ' ':
		//	term.writeCell('·')

		default:
			if r < 32 {
				log.Printf("Unexpected ASCII control character: %d", r)
			}
			term.writeCell(r)
		}

		term._mutex.Unlock()
	}

}

func multiplyN(n *int32, r rune) {
	if *n < 0 {
		*n = 0
	}

	*n = (*n * 10) + (r - 48)
}
