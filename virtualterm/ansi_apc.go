package virtualterm

import (
	"log"
	"strings"

	"github.com/lmorg/mxtty/codes"
)

type apcSlice []string

func (as *apcSlice) Value(i int) string {
	if len(*as) <= i {
		return ""
	}
	return (*as)[i]
}

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

	stack := apcSlice(strings.Split(string(text[:len(text)-1]), ";"))

	switch stack.Value(0) {
	case "BEGIN":
		switch stack.Value(1) {
		case "table":
			term.mxapcTableBegin(stack)
		default:
			log.Printf("Unknown mxASC code %s: %s", stack[1], string(text[:len(text)-1]))
		}
	case "END":
		switch stack.Value(1) {
		case "table":
			term.mxapcTableEnd(stack)
		default:
			log.Printf("Unknown mxASC code %s: %s", stack[1], string(text[:len(text)-1]))
		}
	default:
		log.Printf("Unknown ASC code %s: %s", stack[0], string(text[:len(text)-1]))
	}
}
