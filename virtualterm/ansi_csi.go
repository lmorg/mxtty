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
	)

	for {
		r = term.Pty.ReadRune()
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

		case 's': // save cursor pos
			term.savedCurPos = term.curPos

		case 'u': // restore cursor pos
			term.curPos = term.savedCurPos

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
			log.Printf("Unknown CSI code %d: %v", *n, stack)
		}

		if isCsiTerminator(r) {
			return
		}
	}
}

func (term *Term) parsePrivateCodes() []rune {
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

func lookupSgr(sgr *sgr, n int32, stack []int32) {
	for _, i := range stack {
		switch i {
		case 0: // reset / normal
			sgr.Reset()

		case 1: // bold
			sgr.bitwise.Set(types.SGR_BOLD)

		case 3: // italic
			sgr.bitwise.Set(types.SGR_ITALIC)

		case 4: // underline
			sgr.bitwise.Set(types.SGR_UNDERLINE)

		case 5, 6: // blink
			sgr.bitwise.Set(types.SGR_SLOW_BLINK)

		case 7: // invert
			sgr.bitwise.Set(types.SGR_INVERT)

		case 22: // no bold
			sgr.bitwise.Unset(types.SGR_BOLD)

		case 23: // no italic
			sgr.bitwise.Unset(types.SGR_ITALIC)

		case 24: // no underline
			sgr.bitwise.Unset(types.SGR_UNDERLINE)

		case 25: // no blink
			sgr.bitwise.Unset(types.SGR_SLOW_BLINK)

		case 27: // no invert
			sgr.bitwise.Unset(types.SGR_INVERT)

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
			colour := _sgrEnhancedColour(n, stack)
			if colour != nil {
				sgr.fg = colour
			}
			return

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

		case 48:
			colour := _sgrEnhancedColour(n, stack)
			if colour != nil {
				sgr.bg = colour
			}
			return

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
}

func _sgrEnhancedColour(n int32, stack []int32) *types.Colour {
	if len(stack) < 2 {
		log.Printf("SGR error: too few parameters in %d: %v", n, stack)
		return nil
	}
	switch stack[1] {
	case 5:
		colour, ok := SGR_COLOUR_256[stack[2]]
		if !ok {
			log.Printf("SGR error: 256 value does not exist in %d: %v", n, stack)
			return nil
		}
		return colour

	case 2:
		if len(stack) != 5 {
			log.Printf("SGR error: too few parameters in %d (24bit): %v", n, stack)
			return nil
		}
		return &types.Colour{
			Red:   byte(stack[2]),
			Green: byte(stack[3]),
			Blue:  byte(stack[4]),
		}

	default:
		log.Printf("SGR error: unexpected value in %d: %v", n, stack)
		return nil
	}

}

func lookupPrivateCsi(term *Term, code []rune) {
	param := string(code[:len(code)-1])
	r := code[len(code)-1]
	switch r {
	case 'h':
		switch param {
		case "47": // alt screen buffer
			term.cells = &term.altBuf
		default:
			log.Printf("Private CSI parameter not implemented in %s: %v", string(r), param)
		}

	case 'l':
		switch param {
		case "47": // normal screen buffer
			term.cells = &term.normBuf
		default:
			log.Printf("Private CSI parameter not implemented in %s: %v", string(r), param)
		}

	default:
		log.Printf("Private CSI code not implemented: %s (%s)", string(r), string(code))
	}
}
