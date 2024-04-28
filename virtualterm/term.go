package virtualterm

import (
	"sync"

	"github.com/lmorg/mxtty/charset"
	"github.com/lmorg/mxtty/types"
)

/*
	The virtual terminal emulator

	There is a distinct lack of unit tests in this package. That will change
	over time, however it is worth noting that the hardest part of the problem
	is understanding _what_ correct behaviour should look like, as opposed to
	any logic itself being complex. Therefore this package has has extensive
	manual testing against it via the following CLI applications:
	- vttest: https://invisible-island.net/vttest/vttest.html
	- vim
	- tmux
	- murex: https://murex.rocks
	- bash

	...as well as heavy reliance on documentation, as described in each source
	file.
*/

// Term is the display state of the virtual term
type Term struct {
	size     *types.XY
	sgr      *types.Sgr
	renderer types.Renderer
	Pty      types.Pty
	_mutex   sync.Mutex

	cells         *[][]types.Cell
	_normBuf      [][]types.Cell
	_altBuf       [][]types.Cell
	_scrollBuf    [][]types.Cell
	_scrollOffset int
	_scrollMsg    types.Notification

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
	_vtMode          _stateVtMode
	_activeElement   types.Element
	_slowBlinkState  bool
	_insertOrReplace _stateIrmT

	// character sets
	_activeCharSet int
	_charSetG      [4]map[rune]rune

	// misc CSI configs
	_windowTitleStack []string
	_noAutoLineWrap   bool // No Auto-Wrap Mode (DECAWM), VT100.
}

type _stateVtMode int

const (
	_VT100   = 0
	_VT52    = 1
	_TEK4014 = 2
)

type _stateIrmT int

const (
	_STATE_IRM_REPLACE = 0
	_STATE_IRM_INSERT  = 1
)

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
	var size *types.XY

	if renderer != nil {
		size = renderer.TermSize()
	}

	term := &Term{
		renderer: renderer,
	}

	term.reset(size)

	return term
}

func (term *Term) Start(pty types.Pty, shell string) {
	term.Pty = pty

	go term.exec(shell)
	go term.readLoop()
	go term.slowBlink()
}

func (term *Term) reset(size *types.XY) {
	if term.renderer != nil {
		term.renderer.AddRenderFnToStack(func() {
			term.renderer.ResizeWindow(size)
		})
	}

	term.size = size
	term.curPos = types.XY{}

	term._normBuf = term.makeScreen()
	term._altBuf = term.makeScreen()
	term.eraseScrollBack()

	term._tabWidth = 8
	term.csiResetTabStops()

	term.cells = &term._normBuf

	term.sgr = types.SGR_DEFAULT.Copy()

	term._charSetG[1] = charset.DecSpecialChar

	term._lfFrequency = 2
	term._lfEnabled = true
}

func (term *Term) makeScreen() [][]types.Cell {
	screen := make([][]types.Cell, term.size.Y)
	for i := range screen {
		screen[i] = term.makeRow()
	}
	return screen
}

func (term *Term) makeRow() []types.Cell {
	return make([]types.Cell, term.size.X)
}

func (term *Term) GetSize() *types.XY {
	return term.size
}

func (term *Term) cell() *types.Cell {
	if term.curPos.X < 0 {
		//term.renderer.DisplayNotification(types.NOTIFY_DEBUG, "term.curPos.X < 0 (returning first cell)")
		term.curPos.X = 0
		//term.lineFeed()
	}

	if term.curPos.Y < 0 {
		//term.renderer.DisplayNotification(types.NOTIFY_DEBUG, "term.curPos.Y < 0 (returning first cell)")
		term.curPos.Y = 0
	}

	if term.curPos.X >= term.size.X {
		//term.renderer.DisplayNotification(types.NOTIFY_DEBUG, "term.curPos.X >= term.size.X (returning last cell)")
		term.curPos.X = 0
		term.lineFeed()
	}

	if term.curPos.Y >= term.size.Y {
		//term.renderer.DisplayNotification(types.NOTIFY_DEBUG, "term.curPos.Y >= term.size.Y (returning last cell)")
		term.curPos.Y = term.size.Y - 1
	}

	return &(*term.cells)[term.curPos.Y][term.curPos.X]
}

func (term *Term) previousCell() (*types.Cell, *types.XY) {
	pos := term.curPos
	pos.X--

	if pos.X < 0 {
		pos.X = term.size.X - 1
		pos.Y--
	} else if pos.X >= term.size.X {
		pos.X = term.size.X - 1
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

func (term *Term) ShowCursor(v bool) {
	term._hideCursor = !v
}
