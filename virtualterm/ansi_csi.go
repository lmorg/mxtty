package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/virtualterm/types"
)

func isCsiTerminator(r rune) bool {
	return r >= 0x40 && r <= 0x7E
}

func parseCsiCodes(term *Term, text []rune) int {
	i := 2

	var (
		stack = []int32{-1} // default value
		n     = &stack[0]
	)

	for {
		for ; i < len(text); i++ {
			if text[i] >= '0' && '9' >= text[i] {
				multiplyN(n, text[i])
				continue
			}

			switch text[i] {
			case 'A', 'E': // moveCursorUp
				term.moveCursorUpwards(*n)

			case 'B', 'F': // moveCursorDown
				term.moveCursorDownwards(*n)

			case 'C': // moveCursorForwards
				term.moveCursorForwards(*n)

			case 'D': // moveCursorBackwards
				term.moveCursorBackwards(*n)

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
					term.eraseDisplay() // TODO: 3 should erase scrollback buffer
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

			case '?': // private codes
				adjust, n, r := parseNumericAlphaCodes(i, text)
				log.Printf("CSI private code gobbled: '[?%d%s'", n, string(r))
				return i + adjust - 3

			case ';':
				stack = append(stack, -1)
				n = &stack[len(stack)-1]
				//log.Printf("Unhandled CSI parameter: '%d;'", n)

			default:
				log.Printf("Unknown CSI code: '%d%s'", n, string(text[i]))
			}

			if isCsiTerminator(text[i]) {
				return i - 1
			}
		}

		p := make([]byte, 10*1024)
		n, err := term.Pty.Read(p)
		if err != nil {
			log.Printf("error reading from buffer (%d bytes dropped): %s", n, err.Error())
			continue
		}
		r := []rune(string(p[:n]))
		text = append(text, r...)
	}

	//return i - 1
}

func parseNumericAlphaCodes(i int, text []rune) (int, int32, rune) {
	i++
	var n int32 = -1 // default value

	for ; i < len(text); i++ {
		if text[i] >= '0' && '9' >= text[i] {
			//n = (n * 10) + (text[i] - 48)
			multiplyN(&n, text[i])
			continue
		}

		if isCsiTerminator(text[i]) {
			return i, n, text[i]
		}

		log.Printf("Unexpected character in private CSI sequence: %s", string(text[i]))
		return i, n, text[i]
	}
	return i, n, 0
}

func lookupSgr(sgr *sgr, n int32, stack []int32) {
	switch n {
	case 0: // reset / normal
		sgr.Reset()

	case 1: // bold
		sgr.Set(SGR_BOLD)

	case 3: // italic
		sgr.Set(SGR_ITALIC)

	case 4: // underline
		sgr.Set(SGR_UNDERLINE)

	case 5, 6: // blink
		sgr.Set(SGR_BLINK)

	case 7: // invert
		sgr.Set(SGR_INVERT)

	case 23: // no italic
		sgr.Unset(SGR_ITALIC)

	case 24: // no underline
		sgr.Unset(SGR_UNDERLINE)

	case 25: // no blink
		sgr.Unset(SGR_BLINK)
	//
	// 3bit foreground colour:
	//

	case 30: // fg black
		sgr.fg = SGR_COLOUR_BLACK

	case 31: // fg red
		sgr.fg = SGR_COLOUR_RED

	case 32: // fg green
		sgr.fg = SGR_COLOUR_GREEN

	case 33: // fg yellow
		sgr.fg = SGR_COLOUR_YELLOW

	case 34: // fg blue
		sgr.fg = SGR_COLOUR_BLUE

	case 35: // fg magenta
		sgr.fg = SGR_COLOUR_MAGENTA

	case 36: // fg cyan
		sgr.fg = SGR_COLOUR_CYAN

	case 37: // fg white
		sgr.fg = SGR_COLOUR_WHITE

	case 38:
		log.Printf("colour: %d, %v", n, stack)
		if len(stack) < 2 {
			return
		}
		switch stack[1] {
		case 5:
		case 2:
			if len(stack) != 5 {
				return
			}
			sgr.fg = types.Colour{
				Red:   byte(stack[2]),
				Green: byte(stack[3]),
				Blue:  byte(stack[4]),
			}
		}

	case 39: // fg default
		sgr.fg = SGR_DEFAULT.fg

		//
		// 3bit background colour:
		//

	case 40: // bg black
		sgr.bg = SGR_COLOUR_BLACK

	case 41: // bg rede
		sgr.bg = SGR_COLOUR_RED

	case 42: // bg green
		sgr.bg = SGR_COLOUR_GREEN

	case 43: // bg yellow
		sgr.bg = SGR_COLOUR_YELLOW

	case 44: // bg blue
		sgr.bg = SGR_COLOUR_BLUE

	case 45: // bg magenta
		sgr.bg = SGR_COLOUR_MAGENTA

	case 46: // bg cyan
		sgr.bg = SGR_COLOUR_CYAN

	case 47: // bg white
		sgr.bg = SGR_COLOUR_WHITE

	case 49: // bg default
		sgr.bg = SGR_DEFAULT.bg

		//
	// 4bit foreground colour:
	//

	case 90: // fg black
		sgr.fg = SGR_COLOUR_BLACK_BRIGHT

	case 91: // fg red
		sgr.fg = SGR_COLOUR_RED_BRIGHT

	case 92: // fg green
		sgr.fg = SGR_COLOUR_GREEN_BRIGHT

	case 93: // fg yellow
		sgr.fg = SGR_COLOUR_YELLOW_BRIGHT

	case 94: // fg blue
		sgr.fg = SGR_COLOUR_BLUE_BRIGHT

	case 95: // fg magenta
		sgr.fg = SGR_COLOUR_MAGENTA_BRIGHT

	case 96: // fg cyan
		sgr.fg = SGR_COLOUR_CYAN_BRIGHT

	case 97: // fg white
		sgr.fg = SGR_COLOUR_WHITE_BRIGHT

		//
		// 4bit background colour:
		//

	case 100: // bg black
		sgr.bg = SGR_COLOUR_BLACK_BRIGHT

	case 101: // bg red
		sgr.bg = SGR_COLOUR_RED_BRIGHT

	case 102: // bg green
		sgr.bg = SGR_COLOUR_GREEN_BRIGHT

	case 103: // bg yellow
		sgr.bg = SGR_COLOUR_YELLOW_BRIGHT

	case 104: // bg blue
		sgr.bg = SGR_COLOUR_BLUE_BRIGHT

	case 105: // bg magenta
		sgr.bg = SGR_COLOUR_MAGENTA_BRIGHT

	case 106: // bg cyan
		sgr.bg = SGR_COLOUR_CYAN_BRIGHT

	case 107: // bg white
		sgr.bg = SGR_COLOUR_WHITE_BRIGHT

	default:
		log.Printf("Unknown SGR code: %d", n)
	}
}
