package virtualterm

import (
	"log"
	"sync"
	"time"

	"github.com/lmorg/mxtty/window"
)

// Term is the display state of the virtual term
type Term struct {
	cells    [][]cell
	size     xy
	curPos   xy
	sgr      *sgr
	mutex    sync.Mutex
	tabWidth int32
}

type cell struct {
	char rune
	sgr  *sgr
}

type xy struct {
	X int32
	Y int32
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
		sgr: &sgr{
			fg: SGR_DEFAULT.fg,
			bg: SGR_DEFAULT.bg,
		},
		tabWidth: 8,
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

func (term *Term) cell() *cell { return &term.cells[term.curPos.Y][term.curPos.X] }
