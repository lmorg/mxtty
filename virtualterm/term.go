package virtualterm

import (
	"sync"

	"github.com/lmorg/mxtty/psuedotty"
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
	Pty      *psuedotty.PTY

	slowBlinkState bool
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
		sgr:      SGR_DEFAULT.Copy(),
		tabWidth: 8,
	}

	return term
}

// GetSize outputs mirror those from terminal and readline packages
func (term *Term) GetSize() *types.Rect {
	return term.size
}

func (term *Term) cell() *cell {
	if term.curPos.X >= term.size.X {
		//log.Printf("out of bounds caught: term.curPos.X >= term.size.X")
		//term.curPos.X = term.size.X - 1
		term.wrapCursorForwards()
	}
	//if term.curPos.Y >= term.size.Y {
	//	log.Printf("out of bounds caught: term.curPos.Y >= term.size.Y")
	//	term.curPos.Y = term.size.Y - 1
	//}
	return &term.cells[term.curPos.Y][term.curPos.X]
}
