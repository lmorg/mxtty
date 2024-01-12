package virtualterm

import (
	"log"
	"strings"

	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/types"
)

func (term *Term) parseApcCodes() {
	var (
		r    rune
		text []rune
	)

	for {
		r = term.Pty.Read()
		text = append(text, r)
		if r == codes.AsciiEscape {
			r = term.Pty.Read()
			if r == '\\' { // ST (APC terminator)
				break
			}
			text = append(text, r)
			continue
		}
	}

	parameters := types.ApcSlice(strings.Split(string(text[:len(text)-1]), ";"))

	switch parameters.Value(0) {
	case "BEGIN":
		switch parameters.Value(1) {
		case "TABLE":
			term.mxapcTableBegin(parameters)
		default:
			log.Printf("Unknown mxASC code %s: %s", parameters[1], string(text[:len(text)-1]))
		}
	case "END":
		switch parameters.Value(1) {
		case "TABLE":
			term.mxapcTableEnd(parameters)
		default:
			log.Printf("Unknown mxASC code %s: %s", parameters[1], string(text[:len(text)-1]))
		}
	default:
		log.Printf("Unknown ASC code %s: %s", parameters[0], string(text[:len(text)-1]))
	}
}
