package virtualterm

import (
	"fmt"

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
	case "begin":
		switch apc.Index(1) {
		case "table":
			term.mxapcBegin(types.ELEMENT_ID_TABLE, apc)

		case "image":
			term.mxapcBegin(types.ELEMENT_ID_IMAGE, apc)

		default:
			term.renderer.DisplayNotification(types.NOTIFY_DEBUG,
				fmt.Sprintf("Unknown mxAPC code %s: %s", apc.Index(1), string(text[:len(text)-1])))
		}

	case "end":
		switch apc.Index(1) {
		case "table":
			term.mxapcEnd()

		case "image":
			term.mxapcEnd()

		default:
			term.renderer.DisplayNotification(types.NOTIFY_DEBUG,
				fmt.Sprintf("Unknown mxAPC code %s: %s", apc.Index(1), string(text[:len(text)-1])))
		}


	/*case "insert":
	switch apc.Index(1) {
	case "image":
		term.mxapcInsertImage(apc)
	}*/

	default:
		term.renderer.DisplayNotification(types.NOTIFY_DEBUG,
			fmt.Sprintf("Unknown mxAPC code %s: %s", apc.Index(1), string(text[:len(text)-1])))
	}
}
