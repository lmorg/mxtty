package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/codes"
)

/*
	Reference documentation used:
	- https://invisible-island.net/xterm/ctlseqs/ctlseqs.html
	- ChatGPT (when the documentation above was unclear)
*/

func (term *Term) parsePmCodes() {
	var (
		r    rune
		text []rune
	)

	for {
		r = term.Pty.Read()
		text = append(text, r)
		if r == codes.AsciiEscape {
			r = term.Pty.Read()
			if r == '\\' { // ST (PM terminator)
				break
			}
			text = append(text, r)
			continue
		}
	}

	log.Printf("DEBUG: Ignored PM code %s", string(text[:len(text)-2]))
}
