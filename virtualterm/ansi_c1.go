package virtualterm

import "log"

/*
	Documentation:
	* https://en.wikipedia.org/wiki/C0_and_C1_control_codes
*/

func (term *Term) parseC1Codes() {
	r := term.Pty.Read()
	switch r {
	case '[':
		// CSI code
		term.parseCsiCodes()

	case ']':
		// OSC code
		term.parseOscCodes()

	case 'P':
		// DCS code
		term.parseDcsCodes()

	case '^': // PM code
		term.parsePmCodes()

	case '_':
		// APC code
		term.parseApcCodes()

	case ' ', '#', '%':
		// 7/8bit controls
		// ANSI conformance level
		// character set
		param := term.Pty.Read() // ignore these sequences
		log.Printf("DEBUG: Ignored '{ESC}%s%s' sequence", string(r), string(param))

	case '(': // Designate G0 Character Set (ISO 2022)
		param := term.Pty.Read()
		log.Printf("TODO: Unhandled escape sequence: {ESC}%s%s", string(r), string(param))

	case ')':
		// Designate G1 Character Set (ISO 2022)
		param := term.Pty.Read()
		log.Printf("TODO: Unhandled escape sequence: {ESC}%s%s", string(r), string(param))

	case '*':
		// Designate G2 Character Set (ISO 2022)
		param := term.Pty.Read()
		log.Printf("TODO: Unhandled escape sequence: {ESC}%s%s", string(r), string(param))

	case '+':
		// Designate G3 Character Set (ISO 2022)
		param := term.Pty.Read()
		log.Printf("TODO: Unhandled escape sequence: {ESC}%s%s", string(r), string(param))

	case '.':
		// Designate G2 Character Set, VT300.
		param := term.Pty.Read()
		log.Printf("TODO: Unhandled escape sequence: {ESC}%s%s", string(r), string(param))

	case '/':
		// Designate G3 Character Set, VT300.
		param := term.Pty.Read()
		log.Printf("TODO: Unhandled escape sequence: {ESC}%s%s", string(r), string(param))

	case '=':
		// Application Keypad (DECPAM)
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case '>':
		// Normal Keypad (DECPNM)
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'c':
		// Full Reset (RIS)
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'l':
		// Memory Lock (per HP terminals). Locks memory above the cursor.
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'm':
		// Memory Unlock (per HP terminals)
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'n':
		// Invoke the G2 Character Set as GL (LS2).
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'o':
		// Invoke the G3 Character Set as GL (LS3).
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case '|':
		// Invoke the G3 Character Set as GR (LS3R).
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case '}':
		// Invoke the G2 Character Set as GR (LS2R).
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case '~':
		// Invoke the G1 Character Set as GR (LS1R).
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case '7':
		term.csiCursorPosSave()

	case '8':
		term.csiCursorPosRestore()

	case '@':
		// Padding Character
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'A':
		// High Octet Preset
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'B':
		// Break Permitted Here
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'C':
		// No Break Here
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'D':
		// Index
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'E':
		// Next Line
		term.carriageReturn()
		term.lineFeed()

	case 'F', 'G':
		// Start of Selected Area
		// Start of Selected Area
		// Cursor to lower left corner of screen (if enabled by the 'hpLowerleftBugCompat' resource).
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'H':
		// Character Tabulation Set
		// Horizontal Tabulation Set
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'I':
		// Character Tabulation With Justification
		// Horizontal Tabulation With Justification
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'J':
		// Line Tabulation Set
		// Vertical Tabulation Set
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'K':
		// Partial Line Forward
		// Partial Line Down
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'L':
		// Partial Line Backward
		// Partial Line Up
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'M':
		// Reverse Line Feed
		// Reverse Index
		term.ReverseLineFeed()

	case 'N':
		// Single-Shift 2
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'O':
		// Single-Shift 3
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'Q':
		// Private Use 1
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'R':
		// Private Use 2
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'S':
		// Set Transmit State
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'T':
		// Cancel character
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'U':
		// Message Waiting
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'V':
		// Start of Protected Area
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'W':
		// End of Protected Area
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'X':
		// Start of String
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'Y':
		// Single Graphic Character Introducer
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case 'Z':
		// Single Character Introducer
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	case '\\':
		// String Terminator
		log.Printf("TODO: Unhandled C1 code: %s", string(r))

	/////

	default:
		log.Printf("WARNING: Unexpected rune after escape: %d", r)
		//term.writeCell(r)
	}
}
