package virtualterm

import (
	"log"
	"unsafe"

	"github.com/lmorg/mxtty/types"
)

func (term *Term) Render() {
	term._mutex.Lock()

	var x, y int32
	var err error
	for y = 0; int(y) < len(*term.cells); y++ {
		for x = 0; int(x) < len((*term.cells)[y]); x++ {
			switch (*term.cells)[y][x].Char {
			default:
				fg, bg := term.sgrOpts((*term.cells)[y][x].Sgr)
				err = term.renderer.PrintRuneColour((*term.cells)[y][x].Char, x, y, fg, bg, (*term.cells)[y][x].Sgr.Bitwise)
			case types.CELL_NULL:
				//fg, bg := term.sgrOpts(SGR_DEFAULT)
				//err = term.renderer.PrintRuneColour(' ', x, y, fg, bg, 0)
			case types.CELL_ELEMENT_BEGIN:
				err = term.drawElement(&(*term.cells)[y][x])
			case types.CELL_ELEMENT_FILL:
				continue
			}
			if err != nil {
				log.Printf("error in %s [x: %d, y: %d, value: '%s']: %s", "(t *Term) Render()", x, y, string((*term.cells)[y][x].Char), err.Error())
			}
		}
	}

	term._blinkCursor()

	term._mutex.Unlock()
}

func (term *Term) sgrOpts(sgr *types.Sgr) (fg *types.Colour, bg *types.Colour) {
	if sgr.Bitwise.Is(types.SGR_INVERT) {
		bg, fg = sgr.Fg, sgr.Bg
	} else {
		fg, bg = sgr.Fg, sgr.Bg
	}

	if unsafe.Pointer(bg) == unsafe.Pointer(types.SGR_DEFAULT.Bg) {
		bg = nil
	}

	return fg, bg
}

func (term *Term) _blinkCursor() {
	if term._hideCursor {
		return
	}

	var (
		fg, bg *types.Colour
		style  types.SgrFlag
	)

	r := term.cell().Char
	if r == 0 {
		r = ' '
		fg, bg = types.BlinkColour[true], types.BlinkColour[false]
		style = 0
	} else {
		fg, bg = term.cell().Sgr.Fg, term.sgr.Bg
		style = term.cell().Sgr.Bitwise
	}

	if term._slowBlinkState {
		fg, bg = bg, fg
	}

	err := term.renderer.PrintRuneColour(r, term.curPos.X, term.curPos.Y, fg, bg, style)
	if err != nil {
		log.Printf("error in %s [cursorBlink]: %s", "(t *Term) _blink()", err.Error())
	}
}
