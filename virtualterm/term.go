package virtualterm

import (
	"log"
	"sync"
	"time"

	"github.com/lmorg/mxtty/window"
)

// Term is the display state of the virtual term
type Term struct {
	cells  [][]cell
	size   xy
	curPos xy
	sgr    *sgr
	mutex  sync.Mutex
}

type cell struct {
	char rune
	sgr  *sgr
}

type sgr struct {
	bitwise sgrFlag
	fg      rgb
	bg      rgb
}

func (c *cell) differs(oldChar rune, oldSgr *sgr) bool {
	if c.sgr.bitwise != oldSgr.bitwise {
		return true
	}

	if c.char == 0 && oldChar != 0 {
		return true
	}

	if c.sgr.fg.Red != oldSgr.fg.Red ||
		c.sgr.fg.Green != oldSgr.fg.Green ||
		c.sgr.fg.Blue != oldSgr.fg.Blue {
		return true
	}

	if c.sgr.bg.Red != oldSgr.bg.Red ||
		c.sgr.bg.Green != oldSgr.bg.Green ||
		c.sgr.bg.Blue != oldSgr.bg.Blue {
		return true
	}

	return false
}

func (sgr *sgr) checkFlag(flag sgrFlag) bool {
	return sgr.bitwise&flag != 0
}

type xy struct {
	X int32
	Y int32
}

type rgb struct {
	Red, Green, Blue byte
}

// NewTerminal creates a new virtual term
func NewTerminal(x, y int32) *Term {
	cells := make([][]cell, y, y)
	for i := range cells {
		cells[i] = make([]cell, x, x)
	}

	term := &Term{
		cells: cells,
		size:  xy{x, y},
		sgr:   new(sgr),
	}

	go term.blink()
	return term
}

func (term *Term) blink() {
	var (
		state bool
		err   error
	)

	for {
		time.Sleep(500 * time.Millisecond)

		err = window.PrintBlink(state, int32(term.curPos.X), int32(term.curPos.Y))
		if err != nil {
			log.Printf("error in %s: %s", "window.PrintBlink()", err.Error())
		}

		err = window.Update()
		if err != nil {
			log.Printf("error in %s: %s", "window.Update()", err.Error())
		}

		state = !state
	}
}

// GetSize outputs mirror those from terminal and readline packages
func (term *Term) GetSize() (int32, int32, error) {
	return term.size.X, term.size.Y, nil
}

// format

func (term *Term) sgrReset() {
	term.sgr.bitwise = 0
	term.sgr.fg = rgb{}
	term.sgr.bg = rgb{}
}

func (term *Term) sgrEffect(flag sgrFlag) {
	term.sgr.bitwise |= flag
}

func (c *cell) clear() {
	c.char = 0
	c.sgr = &sgr{}
}

// moveCursor functions DON'T effect other contents in the grid

func (term *Term) moveCursorBackwards(i int32) (overflow int32) {
	if i < 0 {
		i = 1
	}

	term.curPos.X -= i
	if term.curPos.X < 0 {
		overflow = term.curPos.X * -1
		term.curPos.X = 0
	}

	return
}

func (term *Term) moveCursorForwards(i int32) (overflow int32) {
	if i < 0 {
		i = 1
	}

	term.curPos.X += i
	if term.curPos.X >= term.size.X {
		overflow = term.curPos.X - (term.size.X - 1)
		term.curPos.X = term.size.X - 1
	}

	return
}

func (term *Term) moveCursorUpwards(i int32) (overflow int32) {
	if i < 0 {
		i = 1
	}

	term.curPos.Y -= i
	if term.curPos.Y < 0 {
		overflow = term.curPos.Y * -1
		term.curPos.Y = 0
	}

	return
}

func (term *Term) moveCursorDownwards(i int32) (overflow int32) {
	if i < 0 {
		i = 1
	}

	term.curPos.Y += i
	if term.curPos.Y >= term.size.Y {
		overflow = term.curPos.Y - (term.size.Y - 1)
		term.curPos.Y = term.size.Y - 1
	}

	return
}

func (term *Term) cell() *cell { return &term.cells[term.curPos.Y][term.curPos.X] }

// moveGridPos functions DO effect other contents in the grid

func (term *Term) moveContentsUp() {
	var i int32
	for ; i < term.size.Y-1; i++ {
		term.cells[i] = term.cells[i+1]
	}
	term.cells[i] = make([]cell, term.size.X, term.size.X)
}

func (term *Term) wrapCursorForwards() {
	term.curPos.X += 1

	if term.curPos.X >= term.size.X {
		overflow := term.curPos.X - (term.size.X - 1)
		term.curPos.X = 0

		if overflow > 0 && term.moveCursorDownwards(1) > 0 {
			term.moveContentsUp()
			term.moveCursorDownwards(1)
		}
	}
}
