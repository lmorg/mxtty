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
		r = term.Pty.ReadRune()
		text = append(text, r)
		if r == codes.AsciiEscape {
			r = term.Pty.ReadRune()
			if r == '\\' { // ST (DCS terminator)
				break
			}
			text = append(text, r)
			continue
		}
	}

	//stack := strings.Split(string(text[:len(text)-1]), ";")

	log.Printf("Unhandled DCS code %s", string(text[:len(text)-2]))

}
