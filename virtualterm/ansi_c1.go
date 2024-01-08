package virtualterm

import "log"

func (term *Term) parseC1Codes() {
	r := term.Pty.Read()
	switch r {
	case '[': // CSI code
		term.parseCsiCodes()

	case ']': // OSC code
		term.parseOscCodes()

	case 'P': // DCS code
		term.parseDcsCodes()

	case '^': // PM code
		term.parsePmCodes()

	case '_': // APC code
		term.parseApcCodes()

	case ' ', '#', '%': // 7/8bit controls / ANSI conformance level / character set
		param := term.Pty.Read() // ignore these sequences
		log.Printf("DEBUG: Ignored '{ESC}%s%s' sequence", string(r), string(param))

	case '(': // Designate G0 Character Set (ISO 2022)
		param := term.Pty.Read()
		log.Printf("TODO: Unhandled escape sequence: {ESC}%s%s", string(r), string(param))

	case ')': // Designate G1 Character Set (ISO 2022)
		param := term.Pty.Read()
		log.Printf("TODO: Unhandled escape sequence: {ESC}%s%s", string(r), string(param))

	case '*': // Designate G2 Character Set (ISO 2022)
		param := term.Pty.Read()
		log.Printf("TODO: Unhandled escape sequence: {ESC}%s%s", string(r), string(param))

	case '+': // Designate G3 Character Set (ISO 2022)
		param := term.Pty.Read()
		log.Printf("TODO: Unhandled escape sequence: {ESC}%s%s", string(r), string(param))

	case '.': // Designate G2 Character Set, VT300.
		param := term.Pty.Read()
		log.Printf("TODO: Unhandled escape sequence: {ESC}%s%s", string(r), string(param))

	case '/': // Designate G3 Character Set, VT300.
		param := term.Pty.Read()
		log.Printf("TODO: Unhandled escape sequence: {ESC}%s%s", string(r), string(param))

	case '=': // Application Keypad (DECPAM)
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case '>': // Normal Keypad (DECPNM)
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'F': // Cursor to lower left corner of screen (if enabled by the 'hpLowerleftBugCompat' resource).
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'c': // Full Reset (RIS)
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'l': // Memory Lock (per HP terminals). Locks memory above the cursor.
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'm': // Memory Unlock (per HP terminals)
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'n': // Invoke the G2 Character Set as GL (LS2).
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'o': // Invoke the G3 Character Set as GL (LS3).
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case '|': // Invoke the G3 Character Set as GR (LS3R).
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case '}': // Invoke the G2 Character Set as GR (LS2R).
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case '~': // Invoke the G1 Character Set as GR (LS1R).
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case '7':
		term.csiCursorPosSave()

	case '8':
		term.csiCursorPosRestore()

	default:
		log.Printf("WARNING: Unexpected rune after escape: %d", r)
		//term.writeCell(r)
	}
}
