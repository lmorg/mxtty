package virtualterm

import "log"

func (term *Term) parseC1Codes() {
	r := term.Pty.ReadRune()
	switch r {
	case '[': // CSI code
		term.parseCsiCodes()

	case ']': // OSC code
		term.parseOscCodes()

	case '(': // Designate G0 Character Set (ISO 2022)
		term.Pty.ReadRune() // TODO

	case ')': // Designate G1 Character Set (ISO 2022)
		term.Pty.ReadRune() // TODO

	case '*': // Designate G2 Character Set (ISO 2022)
		term.Pty.ReadRune() // TODO

	case '+': // Designate G3 Character Set (ISO 2022)
		term.Pty.ReadRune() // TODO

	case '=': // Application Keypad (DECPAM)
		// TODO

	case '>': // Normal Keypad (DECPNM)
		// TODO

	case 'F': //  Cursor to lower left corner of screen (if enabled by the 'hpLowerleftBugCompat' resource).
		// TODO

	case 'c': // Full Reset (RIS)
		// TODO

	case 'l': // Memory Lock (per HP terminals). Locks memory above the cursor.
		// TODO

	case 'm': // Memory Unlock (per HP terminals)
		// TODO

	case 'n': // Invoke the G2 Character Set as GL (LS2).
		// TODO

	case 'o': // Invoke the G3 Character Set as GL (LS3).
		// TODO

	case '|': // Invoke the G3 Character Set as GR (LS3R).
		// TODO

	case '}': // Invoke the G2 Character Set as GR (LS2R).
		// TODO

	case '~': // Invoke the G1 Character Set as GR (LS1R).
	// TODO

	default:
		log.Printf("rune after escape: %d", r)
		term.writeCell(r)
	}
}
