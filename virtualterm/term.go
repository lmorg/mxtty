package virtualterm

import (
	"log"
	"sync"

	"github.com/lmorg/mxtty/types"
)

// Term is the display state of the virtual term
type Term struct {
	size     *types.XY
	sgr      *types.Sgr
	renderer types.Renderer
	Pty      types.Pty
	_mutex   sync.Mutex

	cells    *[][]types.Cell
	_normBuf [][]types.Cell
	_altBuf  [][]types.Cell

	// line feed redraw
	_lfEnabled   bool
	_lfNum       int32
	_lfFrequency int32

	// tab stops
	_tabStops []int32
	_tabWidth int32

	// cursor and scrolling
	curPos        types.XY
	_originMode   bool // Origin Mode (DECOM), VT100.
	_hideCursor   bool
	_savedCurPos  types.XY
	_scrollRegion *scrollRegionT

	// state
	_activeElement  types.Element
	_slowBlinkState bool

	// misc CSI configs
	_windowTitleStack []string
}

func (term *Term) lfRedraw() {
	if !term._lfEnabled {
		return
	}

	term._lfNum++
	if term._lfNum >= term._lfFrequency {
		term._lfNum = 0
		term.renderer.TriggerRedraw()
	}
}

// NewTerminal creates a new virtual term
func NewTerminal(renderer types.Renderer) *Term {
	size := renderer.TermSize()

	term := &Term{
		renderer: renderer,
	}

	term.reset(size)

	return term
}

func (term *Term) reset(size *types.XY) {
	term.renderer.ResizeWindow(size)
	term.size = size

	term._normBuf = make([][]types.Cell, size.Y)
	for i := range term._normBuf {
		term._normBuf[i] = make([]types.Cell, size.X)
	}
	term._altBuf = make([][]types.Cell, size.Y)
	for i := range term._altBuf {
		term._altBuf[i] = make([]types.Cell, size.X)
	}

	term._tabWidth = 8
	term.csiResetTabStops()

	term.cells = &term._normBuf

	term.sgr = types.SGR_DEFAULT.Copy()

	term._lfFrequency = 2
	term._lfEnabled = true
}

func (term *Term) newRow() []types.Cell {
	return make([]types.Cell, term.size.X)
}
func (term *Term) GetSize() *types.XY {
	return term.size
}

func (term *Term) cell() *types.Cell {
	if term.curPos.X < 0 {
		log.Printf("ERROR: term.curPos.X < 0(returning first cell) TODO fixme")
		term.curPos.X = 0
	}

	if term.curPos.Y < 0 {
		log.Printf("ERROR: term.curPos.Y < 0 (returning first cell) TODO fixme")
		term.curPos.Y = 0
	}

	if term.curPos.X >= term.size.X {
		log.Printf("ERROR: term.curPos.X >= term.size.X (returning last cell) TODO fixme")
		term.curPos.X = term.size.X - 1
	}

	if term.curPos.Y >= term.size.Y {
		log.Printf("ERROR: term.curPos.Y >= term.size.Y (returning last cell) TODO fixme")
		term.curPos.Y = term.size.Y - 1
	}

	return &(*term.cells)[term.curPos.Y][term.curPos.X]
}

func (term *Term) previousCell() (*types.Cell, *types.XY) {
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

type scrollRegionT struct {
	Top    int32
	Bottom int32
}

func (term *Term) getScrollRegion() (top int32, bottom int32) {
	if term._scrollRegion == nil || !term._originMode {
		top = 0
		bottom = term.size.Y - 1
	} else {
		top = term._scrollRegion.Top - 1
		bottom = term._scrollRegion.Bottom - 1
	}

	return
}

func (term *Term) Reply(b []byte) error {
	return term.Pty.Write(b)
}

func (term *Term) Bg() *types.Colour {
	return types.SGR_DEFAULT.Bg
}

func (term *Term) copyCell(cell *types.Cell) *types.Cell {
	copy := new(types.Cell)
	copy.Char = cell.Char
	if term.cell().Sgr == nil {
		copy.Sgr = term.sgr.Copy()
	} else {
		copy.Sgr = cell.Sgr.Copy()
	}

	return copy
}
