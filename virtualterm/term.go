package virtualterm

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lmorg/mxtty/charset"
	"github.com/lmorg/mxtty/config"
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
	process  *os.Process
	_mutex   sync.Mutex

	cells         *[][]types.Cell
	_normBuf      [][]types.Cell
	_altBuf       [][]types.Cell
	_scrollBuf    [][]types.Cell
	_scrollOffset int
	_scrollMsg    types.Notification

	// smooth scroll
	_ssCounter   int32
	_ssFrequency int32

	// tab stops
	_tabStops []int32
	_tabWidth int32

	// cursor and scrolling
	_curPos       types.XY
	_originMode   bool // Origin Mode (DECOM), VT100.
	_hideCursor   bool
	_savedCurPos  types.XY
	_scrollRegion *scrollRegionT

	// state
	_vtMode          _stateVtMode
	_slowBlinkState  bool
	_insertOrReplace _stateIrmT
	_hasFocus        bool
	_activeElement   types.Element

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
	if term.renderer == nil {
		return
	}

	term._ssCounter++
	if term._ssCounter >= term._ssFrequency {
		term._ssCounter = 0
		term.renderer.TriggerRedraw()
	}
}

// NewTerminal creates a new virtual term
func NewTerminal(renderer types.Renderer) *Term {
	var size *types.XY

	if renderer != nil {
		size = renderer.GetTermSize()
	}

	term := &Term{
		renderer: renderer,
	}

	term.reset(size)

	return term
}

func (term *Term) Start(pty types.Pty) {
	term.Pty = pty

	go term.exec()
	go term.readLoop()
	go term.slowBlink()
	go term.refreshInterval()
}

func (term *Term) refreshInterval() {
	if config.Config.Terminal.RefreshInterval == 0 {
		return
	}

	d := time.Duration(config.Config.Terminal.RefreshInterval) * time.Millisecond
	//time.Sleep(3 * time.Second) // lets let everything start first
	for {
		time.Sleep(d)
		term.renderer.TriggerRedraw()
	}
}

func (term *Term) reset(size *types.XY) {
	if term.renderer != nil {
		term.renderer.AddRenderFnToStack(func() {
			term.renderer.ResizeWindow(size)
		})
	}

	term.size = size
	term.resizePty()
	term._curPos = types.XY{}

	term._normBuf = term.makeScreen()
	term._altBuf = term.makeScreen()
	term.eraseScrollBack()

	term._tabWidth = 8
	term.csiResetTabStops()

	term.cells = &term._normBuf

	term.sgr = types.SGR_DEFAULT.Copy()

	term._charSetG[1] = charset.DecSpecialChar

	term.setJumpScroll()
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

func (term *Term) currentCell() *types.Cell {
	pos := term.curPos()

	return &(*term.cells)[pos.Y][pos.X]
}

func (term *Term) previousCell() (*types.Cell, *types.XY) {
	pos := term.curPos()
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

	return &(*term.cells)[pos.Y][pos.X], pos
}

func (term *Term) curPos() *types.XY {
	var y int32
	switch {
	case term._curPos.Y < 0:
		y = 0
	case term._curPos.Y > term.size.Y:
		y = term.size.Y - 1
		//term.lineFeed()
	default:
		y = term._curPos.Y
	}

	var x int32
	switch {
	case term._curPos.X < 0:
		x = 0
	case term._curPos.X >= term.size.X:
		x = term.size.X - 1
	default:
		x = term._curPos.X
	}

	return &types.XY{X: x, Y: y}
}

type scrollRegionT struct {
	Top    int32
	Bottom int32
}

func (term *Term) Reply(b []byte) {
	err := term.Pty.Write(b)
	if err != nil {
		term.renderer.DisplayNotification(types.NOTIFY_ERROR, fmt.Sprintf("Cannot write to PTY: %s", err.Error()))
	}
}

func (term *Term) Bg() *types.Colour {
	if term._hasFocus {
		return types.SGR_DEFAULT.Bg
	}

	return types.BgUnfocused
}

func (term *Term) copyCurrentCell(cell *types.Cell) *types.Cell {
	copy := new(types.Cell)
	copy.Char = cell.Char
	if term.currentCell().Sgr == nil {
		copy.Sgr = term.sgr.Copy()
	} else {
		copy.Sgr = cell.Sgr.Copy()
	}

	return copy
}

func (term *Term) ShowCursor(v bool) {
	term._hideCursor = !v
}

func (term *Term) visibleScreen() [][]types.Cell {
	if term._scrollOffset == 0 {
		return *term.cells
	}

	// render scrollback buffer
	start := len(term._scrollBuf) - term._scrollOffset
	cells := term._scrollBuf[start:]
	if len(cells) < int(term.size.Y) {
		cells = append(cells, term._normBuf...)
	}

	return cells
}

func (term *Term) CopyRange(topLeft, bottomRight *types.XY) []byte {
	// This is some ugly ass code. Sorry!
	// It is also called infrequently and not worth my time optimizing right now
	var (
		ix, iy int
		x, y   int32
		cells  = term.visibleScreen()
		b      []byte
		line   string
	)

	for iy = range cells {
		for ix = range cells[y] {
			x, y = int32(ix), int32(iy)
			switch {
			case bottomRight.Y < topLeft.Y: // select up
				// start multiline
				if (x <= topLeft.X && y == topLeft.Y) ||
					// middle multiline
					(y < topLeft.Y && y > bottomRight.Y) ||
					// end multiline
					(x >= bottomRight.X && y == bottomRight.Y) {
					line += string(cells[y][x].Rune())
				}

			case topLeft.Y == bottomRight.Y:
				// midline
				if bottomRight.X < topLeft.X { //backwards
					if x <= topLeft.X && x >= bottomRight.X && y == topLeft.Y {
						line += string(cells[y][x].Rune())
					}
				} else { // forwards
					if x >= topLeft.X && x <= bottomRight.X && y == topLeft.Y {
						line += string(cells[y][x].Rune())
					}
				}

			default: // select down
				// start multiline
				if (x >= topLeft.X && y == topLeft.Y) ||
					// middle multiline
					(y > topLeft.Y && y < bottomRight.Y) ||
					// end multiline
					(x <= bottomRight.X && y == bottomRight.Y) {
					line += string(cells[y][x].Rune())
				}
			}
		}
		if len(line) > 0 {
			line = strings.TrimRight(line, " ")
			b = append(b, []byte(line+"\n")...)
			line = ""
		}
	}

	if len(b) > 0 {
		return b[:len(b)-1]
	}
	return b
}

func (term *Term) CopyLines(top, bottom int32) []byte {
	cells := term.visibleScreen()
	var b []byte

	for y := top; y <= bottom; y++ {
		var line string
		for x := range cells[y] {
			line += string(cells[y][x].Rune())
		}
		line = strings.TrimRight(line, " ") + "\n"
		b = append(b, []byte(line)...)
	}

	if len(b) > 0 {
		return b[:len(b)-1]
	}
	return b
}

func (term *Term) CopySquare(begin *types.XY, end *types.XY) []byte {
	cells := term.visibleScreen()
	var b []byte

	for y := begin.Y; y <= end.Y; y++ {
		var line string
		for x := begin.X; x <= end.X; x++ {
			line += string(cells[y][x].Rune())
		}
		line = strings.TrimRight(line, " ") + "\n"
		b = append(b, []byte(line)...)
	}

	if len(b) > 0 {
		return b[:len(b)-1]
	}
	return b
}

func (term *Term) HasFocus(state bool) {
	term._hasFocus = state
	term._slowBlinkState = true
}
