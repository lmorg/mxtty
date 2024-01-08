package virtualterm

import (
	"log"
	"strings"

	"github.com/lmorg/mxtty/codes"
)

func (term *Term) parseOscCodes() {
	var (
		r    rune
		text []rune
	)

	for {
		r = term.Pty.ReadRune()
		text = append(text, r)
		if r == codes.AsciiEscape {
			r = term.Pty.ReadRune()
			if r == '\\' { // ST (OSC terminator)
				break
			}
			text = append(text, r)
			continue
		}
		if r == codes.AsciiCtrlG { // bell (xterm OSC terminator)
			break
		}
	}

	stack := strings.Split(string(text[:len(text)-1]), ";")

	switch stack[0] {
	case "0": // change icon and window title
		term.renderer.SetWindowTitle(stack[1])

	case "2": // change window title
		term.renderer.SetWindowTitle(stack[1])

	default:
		log.Printf("Unknown OSC code %s: %s", stack[0], string(text[:len(text)-1]))
	}
}
