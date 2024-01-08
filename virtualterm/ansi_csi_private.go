package virtualterm

import "log"

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

func lookupPrivateCsi(term *Term, code []rune) {
	param := string(code[:len(code)-1])
	r := code[len(code)-1]
	switch r {
	case 'h':
		switch param {
		case "12", "25": // Stop Blinking Cursor (att610) / Hide Cursor (DECTCEM)
			term.csiCursorShow()

		case "47", "1047": // alt screen buffer
			term.csiScreenBufferAlternative()

		case "1048":
			term.csiCursorPosSave()

		case "1049":
			term.csiCursorPosSave()
			term.csiScreenBufferAlternative()

		default:
			log.Printf("Private CSI parameter not implemented in %s: %v [param: %s]", string(r), string(code), param)
		}

	case 'l':
		switch param {
		case "12", "25": // Start Blinking Cursor (att610) / Show Cursor (DECTCEM)
			term.csiCursorShow()

		case "47", "1047": // normal screen buffer
			term.csiScreenBufferNormal()

		case "1048":
			term.csiCursorPosRestore()

		case "1049":
			term.csiScreenBufferNormal()
			term.csiCursorPosRestore()

		default:
			log.Printf("Private CSI parameter not implemented in %s: %v [param: %s]", string(r), string(code), param)
		}

	default:
		log.Printf("Private CSI code not implemented: %s (%s)", string(r), string(code))
	}
}
