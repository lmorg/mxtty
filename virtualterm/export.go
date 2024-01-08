package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/virtualterm/types"
)

// ExportString returns a character map of the virtual terminal
func (term *Term) ExportString() string {
	term.mutex.Lock()

	gridLen := (term.size.X + 1) * term.size.Y
	r := make([]rune, gridLen, gridLen)
	var i int
	for y := range term.cells {
		for x := range term.cells[y] {
			if term.cells[y][x].char != 0 { // if cell contains no data then lets assume it's a space character
				r[i] = term.cells[y][x].char
			} else {
				r[i] = ' '
			}
			i++
		}
		r[i] = '\n'
		i++
	}

	term.mutex.Unlock()

	return string(r)
}

// ExportString returns a character map of the virtual terminal
func (term *Term) ExportMxTTY() {
	term.mutex.Lock()

	var x, y int32
	var err error
	for y = 0; int(y) < len(term.cells); y++ {
		for x = 0; int(x) < len(term.cells[y]); x++ {
			if term.cells[y][x].char != 0 {
				fg, bg := sgrOpts(term.cells[y][x].sgr)
				err = term.renderer.PrintRuneColor(term.cells[y][x].char, x, y, fg, bg)
			} else {
				fg, bg := sgrOpts(SGR_DEFAULT)
				err = term.renderer.PrintRuneColor(' ', x, y, fg, bg)
			}
			if err != nil {
				log.Printf("error in %s [x: %d, y: %d, value: '%s']: %s", "(t *Term) ExportMxTTY()", x, y, string(term.cells[y][x].char), err.Error())
			}
		}
	}

	term.mutex.Unlock()

	err = term.renderer.Update()
	if err != nil {
		log.Printf("error in %s [x: %d, y: %d]: %s", "(t *Term) ExportMxTTY()", x, y, err.Error())
	}
}

func sgrOpts(sgr *sgr) (fg *types.Colour, bg *types.Colour) {
	if sgr.Is(sgrInvert) {
		bg, fg = sgr.fg, sgr.bg
	} else {
		fg, bg = sgr.fg, sgr.bg
	}

	return fg, bg
}
