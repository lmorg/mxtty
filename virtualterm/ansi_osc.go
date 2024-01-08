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
		r = term.Pty.Read()
		text = append(text, r)
		if r == codes.AsciiEscape {
			r = term.Pty.Read()
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

	case "1337":
		//$(osc)1337;File=inline=1:${base64 -i $file -o -}

	default:
		log.Printf("Unknown OSC code %s: %s", stack[0], string(text[:len(text)-1]))
	}
}
