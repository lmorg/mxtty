package virtualterm

import (
	"log"
	"sync"

	"github.com/lmorg/mxtty/psuedotty"
	"github.com/lmorg/mxtty/types"
)

// Term is the display state of the virtual term
type Term struct {
	cells       *[][]cell
	normBuf     [][]cell
	altBuf      [][]cell
	size        *types.Rect
	curPos      types.Rect
	savedCurPos types.Rect
	sgr         *sgr
	tabWidth    int32
	renderer    types.Renderer
	Pty         *psuedotty.PTY
	mutex       sync.Mutex

	slowBlinkState bool
}

type cell struct {
	char rune
	sgr  *sgr
}

// NewTerminal creates a new virtual term
func NewTerminal(renderer types.Renderer) *Term {
	size := renderer.Size()

	normBuf := make([][]cell, size.Y)
	for i := range normBuf {
		normBuf[i] = make([]cell, size.X)
	}
	altBuf := make([][]cell, size.Y)
	for i := range altBuf {
		altBuf[i] = make([]cell, size.X)
	}

	term := &Term{
		renderer: renderer,
		normBuf:  normBuf,
		altBuf:   altBuf,
		size:     size,
		sgr:      SGR_DEFAULT.Copy(),
		tabWidth: 8,
	}

	term.cells = &term.normBuf

	return term
}

// GetSize outputs mirror those from terminal and readline packages
func (term *Term) GetSize() *types.Rect {
	return term.size
}

func (term *Term) cell() *cell {
	if term.curPos.Y >= term.size.Y {
		//panic("This shouldn't happen")
		log.Printf("ERROR: term.curPos.Y >= term.size.Y (returning last cell)")
		term.curPos.Y = term.size.Y - 1
		//return &(*term.cells)[term.size.Y-1][term.curPos.X]
		//term.wrapCursorForwards()
	}

	return &(*term.cells)[term.curPos.Y][term.curPos.X]
}

func (term *Term) CopyCells(src [][]cell) [][]cell {
	dst := make([][]cell, len(src))
	for y := range src {
		dst[y] = make([]cell, len(src[y]))
		for x := range src[y] {
			dst[y][x].char = src[y][x].char
			dst[y][x].sgr = src[y][x].sgr.Copy()
		}
	}

	return dst
}
