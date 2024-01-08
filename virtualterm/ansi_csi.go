package virtualterm

import (
	"log"
)

func isCsiTerminator(r rune) bool {
	return r >= 0x40 && r <= 0x7E
}

func multiplyN(n *int32, r rune) {
	if *n < 0 {
		*n = 0
	}

	*n = (*n * 10) + (r - 48)
}

func (term *Term) parseCsiCodes() {
	var (
		r       rune
		stack   = []int32{-1} // default value
		n       = &stack[0]
		cache   []rune
		unknown bool
	)

	for {
		r = term.Pty.ReadRune()
		cache = append(cache, r)
		if r >= '0' && '9' >= r {
			multiplyN(n, r)
			continue
		}

		switch r {
		case 'A', 'E': // moveCursorUp
			term.moveCursorUpwards(*n)

		case 'B', 'F': // moveCursorDown
			term.moveCursorDownwards(*n)

		case 'C': // moveCursorForwards
			term.moveCursorForwards(*n)

		case 'd':
			switch len(stack) {
			case 0:
				term.moveCursorToPos(-1, 0)
			case 1:
				term.moveCursorToPos(-1, *n-1)
			case 2:
				term.moveCursorToPos(stack[1]-1, stack[0]-1)
			default:
				term.moveCursorToPos(stack[1]-1, stack[0]-1)
				log.Printf("more parameters than expected for %s: %v (%s)", string(r), stack, string(cache))
			}

		case 'D': // moveCursorBackwards
			term.moveCursorBackwards(*n)

		case 'G':
			switch len(stack) {
			case 0:
				term.moveCursorToPos(0, -1)
			case 1:
				term.moveCursorToPos(*n-1, -1)
			case 2:
				term.moveCursorToPos(stack[0]-1, stack[1]-1)
			default:
				term.moveCursorToPos(stack[0]-1, stack[1]-1)
				log.Printf("more parameters than expected for %s: %v (%s)", string(r), stack, string(cache))
			}

		case 'H': // moveCursor
			if len(stack) != 2 {
				term.moveCursorToPos(0, 0)
			} else {
				term.moveCursorToPos(stack[1]-1, stack[0]-1)
			}

		case 'J': // eraseDisplay...
			switch *n {
			case -1, 0:
				term.eraseDisplayAfter()
			case 1:
				term.eraseDisplayBefore()
			case 2, 3:
				term.eraseDisplay() // TODO: 3 should also erase scrollback buffer
			}

		case 'K': // clearLine...
			switch *n {
			case -1, 0:
				term.eraseLineAfter()
			case 1:
				term.eraseLineBefore()
			case 2:
				term.eraseLine()
			}

		case 'm': // SGR
			lookupSgr(term.sgr, stack[0], stack)

		//case 'M': // Delete Ps Line(s) (default = 1) (DL).

		case 'n': // Device Status Report (DSR).
			switch *n {
			case 6:
				term.csiCallback("%d;%dR", term.curPos.Y+1, term.curPos.X+1)
			}

		case 'P': // Delete Ps Character(s) (default = 1) (DCH).
			term.deleteCharacters(*n)

		case 'r': // Set Scrolling Region [top;bottom] (default = full size of window) (DECSTBM)
			if len(stack) != 2 {
				log.Printf("Unexpected number of parameters in CSI r (%s): %v", string(cache), stack)
			} else {
				term.csiSetScrollingRegion(stack)
			}

		case 's': // save cursor pos
			term.csiCursorPosSave()

		case 'S': // scroll up
			term.scrollUp(*n)
			//term.moveCursorUpwards(*n)

		case 't': // Window manipulation (XTWINOPS)
			var p2 int32
			if len(stack) > 1 {
				p2 = stack[1]
			}
			switch stack[0] {
			case 22:
				switch p2 {
				case 0, 2:
					term.csiWindowTitleStackSaveTo()
				}
			case 23:
				switch p2 {
				case 0, 2:
					term.csiWindowTitleStackRestoreFrom()
				}
			default:
				log.Printf("Unknown CSI code %d: %v (%s)", *n, stack, string(cache))
			}

		case 'T': // scroll down
			term.scrollDown(*n)

		case 'u': // restore cursor pos
			term.csiCursorPosRestore()

		case '?': // private codes
			code := term.parseCsiExtendedCodes()
			lookupPrivateCsi(term, code)
			return

		case '>': // secondary codes
			code := term.parseCsiExtendedCodes()
			log.Printf("TODO: Secondary CSI code ignored: '%s%s'", string(cache), string(code))
			return

		case '=': // tertiary codes
			code := term.parseCsiExtendedCodes()
			log.Printf("TODO: Tertiary CSI code ignored: '%s%s'", string(cache), string(code))
			return

		//case '!': //

		case ':', ';':
			stack = append(stack, -1)
			n = &stack[len(stack)-1]

		default:
			unknown = true
		}

		if isCsiTerminator(r) {
			if unknown {
				log.Printf("Unknown CSI code %s: %v [string: %s]", string(r), cache, string(cache))
			}
			return
		}
	}
}

func (term *Term) parseCsiExtendedCodes() []rune {
	var (
		r    rune
		code []rune
	)

	for {
		r = term.Pty.ReadRune()
		code = append(code, r)
		if isCsiTerminator(r) {
			return code
		}
	}
}
