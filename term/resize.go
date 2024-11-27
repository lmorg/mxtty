package virtualterm

import (
	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
	"golang.org/x/sys/unix"
)

func (term *Term) Resize(size *types.XY) {
	xDiff := int(size.X - term.size.X)
	yDiff := int(size.Y - term.size.Y)

	term._mutex.Lock()

	term.size = size

	switch {
	case xDiff == 0:
		// nothing to do

	case xDiff > 0:
		// grow
		for y := range term._scrollBuf {
			term._scrollBuf[y] = append(term._scrollBuf[y], make([]types.Cell, xDiff)...)
		}
		for y := range term._normBuf {
			term._normBuf[y] = append(term._normBuf[y], make([]types.Cell, xDiff)...)
		}
		for y := range term._altBuf {
			term._altBuf[y] = append(term._altBuf[y], make([]types.Cell, xDiff)...)
		}

	case xDiff < 0:
		// crop (this is lazy, really we should reflow)
		xDiff = -xDiff
		for y := range term._scrollBuf {
			term._scrollBuf[y] = term._scrollBuf[y][:len(term._scrollBuf[y])-xDiff]
		}
		for y := range term._normBuf {
			term._normBuf[y] = term._normBuf[y][:len(term._normBuf[y])-xDiff]
		}
		for y := range term._altBuf {
			term._altBuf[y] = term._altBuf[y][:len(term._altBuf[y])-xDiff]
		}
	}

	switch {
	case yDiff == 0:
		// nothing to do

	case yDiff > 0:
		// grow
		for i := 0; i < yDiff; i++ {
			term._normBuf = append(term._normBuf, term.makeRow())
		}
		for i := 0; i < yDiff; i++ {
			term._altBuf = append(term._altBuf, term.makeRow())
		}

	case yDiff < 0:
		// shrink
		for i := 0; i < -yDiff; i++ {
			term.appendScrollBuf()
		}
		term._normBuf = term._normBuf[-yDiff:]
		term._altBuf = term._altBuf[-yDiff:]
	}

	term.resizePty()
	defer term._mutex.Unlock()
}

func (term *Term) resizePty() {
	if term.Pty == nil {
		debug.Log("cannot resize pt: term.Pty == nil")
		return
	}

	err := term.Pty.Resize(term.size)
	if err != nil {
		term.renderer.DisplayNotification(types.NOTIFY_ERROR, err.Error())
	}

	if term.process != nil {
		err = term.process.Signal(unix.SIGWINCH)
		if err != nil {
			debug.Log(err)
		}
	}
}

func (term *Term) resize80() {
	term.setSize(&types.XY{X: 80, Y: 24})
}

func (term *Term) resize132() {
	term.setSize(&types.XY{X: 132, Y: 24})
}

func (term *Term) setSize(size *types.XY) {
	term.reset(size)
	term.renderer.ResizeWindow(size)
}
