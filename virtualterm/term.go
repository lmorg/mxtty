package virtualterm

import (
	"log"
	"sync"
	"time"

	"github.com/lmorg/mxtty/virtualterm/types"
)

// Term is the display state of the virtual term
type Term struct {
	cells    [][]cell
	size     *types.Rect
	curPos   types.Rect
	sgr      *sgr
	mutex    sync.Mutex
	tabWidth int32
	renderer *types.Renderer
}

type cell struct {
	char rune
	sgr  *sgr
}

// NewTerminal creates a new virtual term
func NewTerminal(renderer *types.Renderer) *Term {
	cells := make([][]cell, renderer.Size.Y)
	for i := range cells {
		cells[i] = make([]cell, renderer.Size.X)
	}

	term := &Term{
		renderer: renderer,
		cells:    cells,
		size:     renderer.Size,
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

		err = term.renderer.PrintBlink(state, int32(term.curPos.X), int32(term.curPos.Y))
		if err != nil {
			log.Printf("error in %s: %s", "window.PrintBlink()", err.Error())
		}

		err = term.renderer.Update()
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

func (term *Term) cell() *cell {
	if term.curPos.X >= term.size.X {
		log.Printf("out of bounds caught: term.curPos.X >= term.size.X")
		term.curPos.X = term.size.X - 1
	}
	if term.curPos.Y >= term.size.Y {
		log.Printf("out of bounds caught: term.curPos.Y >= term.size.Y")
		term.curPos.Y = term.size.Y - 1
	}
	return &term.cells[term.curPos.Y][term.curPos.X]
}
