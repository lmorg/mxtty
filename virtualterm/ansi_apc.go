package virtualterm

import (
	"log"

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
				text = text[:len(text)-1]
				break
			}
			text = append(text, r)
			continue
		}
		if r == codes.AsciiCtrlG { // bell (xterm OSC terminator)
			text = text[:len(text)-1]
			break
		}
	}

	apc := types.NewApcSlice(text)

	switch apc.Index(0) {
	case "BEGIN":
		log.Println("BEGIN", apc.Index(1))
		switch apc.Index(1) {
		case "TABLE":
			term.mxapcTableBegin(apc)
		default:
			log.Printf("Unknown mxAPC code %s: %s", apc.Index(1), string(text[:len(text)-1]))
		}
	case "END":
		switch apc.Index(1) {
		case "TABLE":
			term.mxapcTableEnd(apc)
		default:
			log.Printf("Unknown mxAPC code %s: %s", apc.Index(1), string(text[:len(text)-1]))
		}
	default:
		log.Printf("Unknown APC code %s: %s", apc.Index(0), string(text[:len(text)-1]))
	}
}
