package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/types"
)

func (term *Term) Render() {
	term.mutex.Lock()

	var x, y int32
	var err error
	for y = 0; int(y) < len(*term.cells); y++ {
		for x = 0; int(x) < len((*term.cells)[y]); x++ {
			if (*term.cells)[y][x].char != 0 {
				fg, bg := term.sgrOpts((*term.cells)[y][x].sgr)
				err = term.renderer.PrintRuneColour((*term.cells)[y][x].char, x, y, fg, bg, (*term.cells)[y][x].sgr.bitwise)
			} else {
				fg, bg := term.sgrOpts(SGR_DEFAULT)
				err = term.renderer.PrintRuneColour(' ', x, y, fg, bg, 0)
			}
			if err != nil {
				log.Printf("error in %s [x: %d, y: %d, value: '%s']: %s", "(t *Term) Render()", x, y, string((*term.cells)[y][x].char), err.Error())
			}
		}
	}

	term._blinkCursor()

	term.mutex.Unlock()
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
	var (
		fg, bg *types.Colour
		style  types.SgrFlag
	)

	cell := term.cell()
	if cell == nil {
		return
	}

	r := cell.char
	if r == 0 {
		r = ' '
		fg, bg = blinkColour[true], blinkColour[false]
		style = 0
	} else {
		fg, bg = term.cell().sgr.fg, term.sgr.bg
		style = term.cell().sgr.bitwise
	}

	if term.slowBlinkState {
		fg, bg = bg, fg
	}

	err := term.renderer.PrintRuneColour(r, term.curPos.X, term.curPos.Y, fg, bg, style)
	if err != nil {
		log.Printf("error in %s [cursorBlink]: %s", "(t *Term) _blink()", err.Error())
	}
}
