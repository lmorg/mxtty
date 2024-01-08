package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/codes"
)

func (term *Term) parsePmCodes() {
	var (
		r    rune
		text []rune
	)

	for {
		r = term.Pty.ReadRune()
		text = append(text, r)
		if r == codes.AsciiEscape {
			r = term.Pty.ReadRune()
			if r == '\\' { // ST (PM terminator)
				break
			}
			text = append(text, r)
			continue
		}
	}

	log.Printf("DEBUG: Ignored PM code %s", string(text[:len(text)-2]))
}
