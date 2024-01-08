package virtualterm

import (
	"log"
	"regexp"

	"github.com/lmorg/mxtty/codes"
)

func (term *Term) writeCell(r rune) {
	term.cell().char = r
	term.cell().sgr = term.sgr
	term.wrapCursorForwards()
}

var rxLazyCsiCheck = regexp.MustCompile(`\[[;\?a-zA-Z0-9]+`)

// Write multiple characters to the virtual terminal
func (term *Term) Write(text []rune) {
	var (
		escape bool
	)

	term.mutex.Lock()

	for i := 0; i < len(text); i++ {
		switch text[i] {
		case codes.AsciiEscape:
			escape = true
			continue

		case codes.AsciiBackspace, codes.IsoBackspace:
			_ = term.moveCursorBackwards(1)

		case codes.AsciiCtrlG: // bell
			// TODO: beep

		case '[':
			if !escape {
				term.writeCell(text[i])
				continue
			}
			escape = false
			start := i
			i += parseCsiCodes(term, text[i-1:])
			if !rxLazyCsiCheck.MatchString(string(text[start : i+1])) {
				log.Printf("Invalid CSI code parsed: %v", []byte(string(text[start:i+1])))
			}

		case ']': // TODO: OSC
			start := i
			for ; i < len(text); i++ {
				if text[i] == 'S' && i < len(text) && text[i+1] == 'T' {
					i += 2
					break
				}
			}
			log.Printf("TODO: OSC sequences: '%s'", string(text[start:i]))

		case '\t':
			indent := int(4 - (term.curPos.X % 4))
			for i := 0; i < indent; i++ {
				term.writeCell(' ')
			}

		case '\r':
			term.curPos.X = 0

		case '\n':
			if term.moveCursorDownwards(1) > 0 {
				term.moveContentsUp()
				term.moveCursorDownwards(1)
			}

		default:
			term.writeCell(text[i])
		}
		escape = false
	}

	term.mutex.Unlock()
}

func multiplyN(n *int32, r rune) {
	if *n < 0 {
		*n = 0
	}

	*n = (*n * 10) + (r - 48)
}
