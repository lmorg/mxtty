package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/types"
)

func isCsiTerminator(r rune) bool {
	return r >= 0x40 && r <= 0x7E
}

func (term *Term) parseCsiCodes() {
	var (
		r     rune
		stack = []int32{-1} // default value
		n     = &stack[0]
		cache []rune
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
				term.moveCursorToPos(-1, *n)
			case 2:
				term.moveCursorToPos(stack[1], stack[0])
			default:
				term.moveCursorToPos(stack[1], stack[0])
				log.Printf("more parameters than expected for %s: %v (%s)", string(r), stack, string(cache))
			}

		case 'D': // moveCursorBackwards
			term.moveCursorBackwards(*n)

		case 'G':
			switch len(stack) {
			case 0:
				term.moveCursorToPos(0, -1)
			case 1:
				term.moveCursorToPos(*n, -1)
			case 2:
				term.moveCursorToPos(stack[0], stack[1])
			default:
				term.moveCursorToPos(stack[0], stack[1])
				log.Printf("more parameters than expected for %s: %v (%s)", string(r), stack, string(cache))
			}

		case 'H': // moveCursor
			if len(stack) != 2 {
				term.curPos = types.Rect{}
			} else {
				term.curPos = types.Rect{
					X: stack[0] + 1,
					Y: stack[1] + 1,
				}
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
			code := term.parsePrivateCodes()
			lookupPrivateCsi(term, code)
			//log.Printf("CSI private code gobbled: '[?%s'", string(code))
			return

		case ':', ';':
			stack = append(stack, -1)
			n = &stack[len(stack)-1]
			//log.Printf("Unhandled CSI parameter: '%d;'", n)

		default:
			log.Printf("Unknown CSI code %d: %v (%s)", *n, stack, string(cache))
		}

		if isCsiTerminator(r) {
			return
		}
	}
}
