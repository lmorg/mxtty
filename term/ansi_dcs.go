package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/codes"
)

/*
	Reference documentation used:
	- xterm: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Device-Control-functions
*/

func (term *Term) parseDcsCodes() {
	var (
		r    rune
		err  error
		text []rune
	)

	for {
		r, err = term.Pty.Read()
		if err != nil {
			return
		}
		text = append(text, r)
		switch r {

		case codes.AsciiEscape:
			r, err = term.Pty.Read()
			if err != nil {
				return
			}
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
