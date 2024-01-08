package virtualterm

import "log"

func isCsiTerminator(r rune) bool {
	return r >= 0x40 && r <= 0x7E
}

func parseCsiCodes(term *Term, text []rune) int {
	i := 2

	var (
		n     int32 = -1 // default value
		stack []int32
	)

	for ; i < len(text); i++ {
		if text[i] >= '0' && '9' >= text[i] {
			//n = (n * 10) + (text[i] - 48)
			multiplyN(&n, text[i])
			continue
		}

		switch text[i] {
		case 'A', 'E': // moveCursorUp
			term.moveCursorUpwards(n)

		case 'B', 'F': // moveCursorDown
			term.moveCursorDownwards(n)

		case 'C': // moveCursorForwards
			term.moveCursorForwards(n)

		case 'D': // moveCursorBackwards
			term.moveCursorBackwards(n)

		case 'm': // SGR
			lookupSgr(term, n)

		case 'J': // eraseDisplay...
			switch n {
			case -1, 0:
				term.eraseDisplayAfter()
			case 1:
				term.eraseDisplayBefore()
			case 2, 3:
				term.eraseDisplay() // TODO: 3 should erase scrollback buffer
			}

		case 'k', 'K': // clearLine...
			switch n {
			case -1, 0:
				term.eraseLineAfter()
			case 1:
				term.eraseLineBefore()
			case 2:
				term.eraseLine()
			}

		case '?': // private codes
			adjust, n, r := parseNumericAlphaCodes(i, text)
			log.Printf("CSI private code gobbled: '[?%d%s'", n, string(r))
			return i + adjust - 3

		case ';':
			stack = append(stack, n)
			//log.Printf("Unhandled CSI parameter: '%d;'", n)

		default:
			log.Printf("Unknown CSI code: '%d%s'", n, string(text[i]))
		}

		if isCsiTerminator(text[i]) {
			return i - 1
		}
	}
	return i
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

func lookupSgr(term *Term, n rune) {
	switch n {
	case 0: // reset / normal
		term.sgrReset()

	case 1: // bold
		term.sgrEffect(sgrBold)

	case 4: // underscore
		term.sgrEffect(sgrUnderscore)

	case 5: // blink
		term.sgrEffect(sgrBlink)

	case 7: // invert
		term.sgrEffect(sgrInvert)

		//
		// 4bit foreground colour:
		//

	case 30: // fg black
		term.sgrEffect(sgrFgColour4)
		term.sgr.fg.Red = sgrColour4Black

	case 31: // fg red
		term.sgrEffect(sgrFgColour4)
		term.sgr.fg.Red = sgrColour4Red

	case 32: // fg green
		term.sgrEffect(sgrFgColour4)
		term.sgr.fg.Red = sgrColour4Green

	case 33: // fg yellow
		term.sgrEffect(sgrFgColour4)
		term.sgr.fg.Red = sgrColour4Yellow

	case 34: // fg blue
		term.sgrEffect(sgrFgColour4)
		term.sgr.fg.Red = sgrColour4Blue

	case 35: // fg magenta
		term.sgrEffect(sgrFgColour4)
		term.sgr.fg.Red = sgrColour4Magenta

	case 36: // fg cyan
		term.sgrEffect(sgrFgColour4)
		term.sgr.fg.Red = sgrColour4Cyan

	case 37: // fg white
		term.sgrEffect(sgrFgColour4)
		term.sgr.fg.Red = sgrColour4White
	}
}
