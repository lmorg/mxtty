package virtualterm

import (
	"fmt"
	"log"

	"github.com/lmorg/mxtty/types"
)

func (term *Term) Render() {
	term._mutex.Lock()

	var x, y int32
	var err error
	for y = 0; int(y) < len(*term.cells); y++ {
		for x = 0; int(x) < len((*term.cells)[y]); x++ {
			switch (*term.cells)[y][x].char {
			default:
				fg, bg := term.sgrOpts((*term.cells)[y][x].sgr)
				err = term.renderer.PrintRuneColour((*term.cells)[y][x].char, x, y, fg, bg, (*term.cells)[y][x].sgr.bitwise)
			case CELL_NULL:
				fg, bg := term.sgrOpts(SGR_DEFAULT)
				err = term.renderer.PrintRuneColour(' ', x, y, fg, bg, 0)
			case CELL_ELEMENT_START:
				e := (*term.cells)[y][x].element
				if e == nil {
					err = fmt.Errorf("nil pointer to element")
				}
				e.Draw(nil) // TODO: this shouldn't be nil
			case CELL_ELEMENT_FILL:
				continue
			}
			if err != nil {
				log.Printf("error in %s [x: %d, y: %d, value: '%s']: %s", "(t *Term) Render()", x, y, string((*term.cells)[y][x].char), err.Error())
			}
		}
	}

	term._blinkCursor()

	term._mutex.Unlock()
}

func (term *Term) sgrOpts(sgr *sgr) (fg *types.Colour, bg *types.Colour) {
	if sgr.bitwise.Is(types.SGR_INVERT) {
		bg, fg = sgr.fg, sgr.bg
	} else {
		fg, bg = sgr.fg, sgr.bg
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

	r := term.cell().char
	if r == 0 {
		r = ' '
		fg, bg = blinkColour[true], blinkColour[false]
		style = 0
	} else {
		fg, bg = term.cell().sgr.fg, term.sgr.bg
		style = term.cell().sgr.bitwise
	}

	if term._slowBlinkState {
		fg, bg = bg, fg
	}

	err := term.renderer.PrintRuneColour(r, term.curPos.X, term.curPos.Y, fg, bg, style)
	if err != nil {
		log.Printf("error in %s [cursorBlink]: %s", "(t *Term) _blink()", err.Error())
	}
}
