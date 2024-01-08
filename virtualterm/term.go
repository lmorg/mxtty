package virtualterm

import (
	"log"
	"sync"

	"github.com/lmorg/mxtty/types"
)

// Term is the display state of the virtual term
type Term struct {
	size     *types.XY
	curPos   types.XY
	sgr      *sgr
	renderer types.Renderer
	Pty      types.Pty
	_mutex   sync.Mutex

	_slowBlinkState bool

	cells    *[][]cell
	_normBuf [][]cell
	_altBuf  [][]cell

	// CSI states
	_tabWidth         int32
	_hideCursor       bool
	_savedCurPos      types.XY
	_scrollRegion     *scrollRegionT
	_windowTitleStack []string
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
		renderer:  renderer,
		_normBuf:  normBuf,
		_altBuf:   altBuf,
		size:      size,
		sgr:       SGR_DEFAULT.Copy(),
		_tabWidth: 8,
	}

	term.cells = &term._normBuf

	return term
}

// GetSize outputs mirror those from terminal and readline packages
func (term *Term) GetSize() *types.XY {
	return term.size
}

func (term *Term) cell() *cell {
	if term.curPos.X < 0 {
		//panic("This shouldn't happen")
		log.Printf("ERROR: term.curPos.X < 0(returning first cell) TODO fixme")
		term.curPos.X = 0
		//return &(*term.cells)[term.size.Y-1][term.curPos.X]
		//term.wrapCursorForwards()
	}

	if term.curPos.Y < 0 {
		//panic("This shouldn't happen")
		log.Printf("ERROR: term.curPos.Y < 0 (returning first cell) TODO fixme")
		term.curPos.Y = 0
		//return &(*term.cells)[term.size.Y-1][term.curPos.X]
		//term.wrapCursorForwards()
	}

	if term.curPos.X >= term.size.X {
		//panic("This shouldn't happen")
		log.Printf("ERROR: term.curPos.X >= term.size.X (returning last cell) TODO fixme")
		term.curPos.X = term.size.X - 1
		//return &(*term.cells)[term.size.Y-1][term.curPos.X]
		//term.wrapCursorForwards()
	}

	if term.curPos.Y >= term.size.Y {
		//panic("This shouldn't happen")
		log.Printf("ERROR: term.curPos.Y >= term.size.Y (returning last cell) TODO fixme")
		term.curPos.Y = term.size.Y - 1
		//return &(*term.cells)[term.size.Y-1][term.curPos.X]
		//term.wrapCursorForwards()
	}

	return &(*term.cells)[term.curPos.Y][term.curPos.X]
}

func (term *Term) previousCell() (*cell, *types.XY) {
	pos := term.curPos
	pos.X--

	if pos.X < 0 {
		pos.X = 0
		pos.Y--
	}

	if pos.Y < 0 {
		pos.Y = 0
	}

	return &(*term.cells)[pos.Y][pos.X], &pos
}

/*func (term *Term) CopyCells(src [][]cell) [][]cell {
	dst := make([][]cell, len(src))
	for y := range src {
		dst[y] = make([]cell, len(src[y]))
		for x := range src[y] {
			dst[y][x].char = src[y][x].char
			dst[y][x].sgr = src[y][x].sgr.Copy()
		}
	}

	return dst
}*/

type scrollRegionT struct {
	Top    int32
	Bottom int32
}

func (term *Term) getScrollRegion() (top int32, bottom int32) {
	if term._scrollRegion == nil {
		top = 0
		bottom = term.size.Y - 1
	} else {
		top = term._scrollRegion.Top - 1
		bottom = term._scrollRegion.Bottom - 1
	}

	return
}

func (term *Term) Return(b []byte) error {
	return term.Pty.Write(b)
}
