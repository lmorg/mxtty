package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/window"
	"github.com/veandco/go-sdl2/sdl"
)

// ExportString returns a character map of the virtual terminal
func (t *Term) ExportString() string {
	t.mutex.Lock()

	gridLen := (t.size.X + 1) * t.size.Y
	r := make([]rune, gridLen, gridLen)
	var i int
	for y := range t.cells {
		for x := range t.cells[y] {
			if t.cells[y][x].char != 0 { // if cell contains no data then lets assume it's a space character
				r[i] = t.cells[y][x].char
			} else {
				r[i] = ' '
			}
			i++
		}
		r[i] = '\n'
		i++
	}

	t.mutex.Unlock()

	return string(r)
}

// ExportString returns a character map of the virtual terminal
func (t *Term) ExportMxTTY() {
	t.mutex.Lock()

	var x, y int32
	var err error
	for y = 0; int(y) < len(t.cells); y++ {
		for x = 0; int(x) < len(t.cells[y]); x++ {
			if t.cells[y][x].char != 0 {
				fg, bg := sdlOpts(t.cells[y][x].sgr)
				err = window.PrintRuneColour(t.cells[y][x].char, x, y, fg, bg)
			} else {
				fg, bg := sdlOpts(SGR_DEFAULT)
				err = window.PrintRuneColour(' ', x, y, fg, bg)
			}
			if err != nil {
				log.Printf("error in %s [x: %d, y: %d, value: '%s']: %s", "(t *Term) ExportMxTTY()", x, y, string(t.cells[y][x].char), err.Error())
			}
		}
	}

	t.mutex.Unlock()

	err = window.Update()
	if err != nil {
		log.Printf("error in %s [x: %d, y: %d]: %s", "(t *Term) ExportMxTTY()", x, y, err.Error())
	}
}

func sdlOpts(sgr *sgr) (*sdl.Color, *sdl.Color) {
	fg := &sdl.Color{
		R: sgr.fg.Red,
		G: sgr.fg.Green,
		B: sgr.fg.Blue,
		A: 255,
	}
	bg := &sdl.Color{
		R: sgr.bg.Red,
		G: sgr.bg.Green,
		B: sgr.bg.Blue,
		A: 255,
	}
	if sgr.Is(sgrInvert) {
		fg, bg = bg, fg
	}

	return fg, bg
}
