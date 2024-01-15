package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/codes"
)

func (term *Term) parseDcsCodes() {
	var (
		r    rune
		text []rune
	)

	for {
		r = term.Pty.Read()
		text = append(text, r)
		switch r {

		case 't': // gobble tmux passthrough codes: {ESC}Ptmux;{ESC}
			if len(text) == 1 {
				r = term.Pty.Read()
				text = append(text, r)
				if r != 'm' {
					continue
				}
				r = term.Pty.Read()
				text = append(text, r)
				if r != 'u' {
					continue
				}
				r = term.Pty.Read()
				text = append(text, r)
				if r != 'x' {
					continue
				}
				r = term.Pty.Read()
				text = append(text, r)
				if r != ';' {
					continue
				}
				r = term.Pty.Read()
				text = append(text, r)
				if r != codes.AsciiEscape {
					continue
				}
				term.Pty.TmuxPassthrough(true)
				return
			}

		case codes.AsciiEscape:
			r = term.Pty.Read()
			if r == '\\' { // ST (DCS terminator)
				goto parsed
			}
			text = append(text, r)
			continue

		case codes.AsciiCtrlG: // bell (xterm OSC terminator)
			goto parsed

		}

	}
parsed:
	text = text[:len(text)-1]

	//stack := strings.Split(string(text), ";")

	log.Printf("Unhandled DCS code %s", string(text))

}
