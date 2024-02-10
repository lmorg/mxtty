package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/debug"
)

func (term *Term) readLoop() {
	var r rune

	for {
		r = term.Pty.Read()
		term._slowBlinkState = false

		term._mutex.Lock()
		term.readChar(r)
		term._mutex.Unlock()
	}
}

/*
	Reference documentation used:
	- ASCII table: https://upload.wikimedia.org/wikipedia/commons/thumb/1/1b/ASCII-Table-wide.svg/1280px-ASCII-Table-wide.svg.png
	- xterm: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Single-character-functions
*/

func (term *Term) readChar(r rune) {
	//writeDebuggingRune(r)

	switch r {

	case codes.AsciiCtrlG:
		// 7: {BELL}
		go term.renderer.Bell()

	case codes.AsciiBackspace, codes.IsoBackspace:
		// 8 / 127
		_ = term.csiMoveCursorBackwards(1)

	case codes.AsciiTab:
		// 9: horizontal tab, \t
		term.printTab()

	case codes.AsciiCtrlJ:
		// 10: line feed, \n
		term.lineFeed()

	case codes.AsciiCtrlK:
		// 11: vertical tab
		term.lineFeed()

	case codes.AsciiCtrlL:
		// 12: form feed
		term.lineFeed()

	case codes.AsciiCtrlM:
		// 13: carriage return, \r
		term.carriageReturn()

	case codes.AsciiEscape:
		// 27: escape, {ESC}
		term.parseC1Codes()

	default:
		if r < ' ' {

			log.Printf("WARNING: Unexpected ASCII control character (ignored): %d", r)
			return
		}
		term.writeCell(r)
	}
}

func writeDebuggingRune(r rune) {
	if !debug.Enabled {
		return
	}

	if r <= ' ' {
		debug.Log(r)
	} else {
		debug.Log(string(r))
	}
}
