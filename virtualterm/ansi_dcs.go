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

	log.Printf("WARNING: Unhandled DCS code %s", string(text))

}
