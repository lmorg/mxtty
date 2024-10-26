package virtualterm

import (
	"github.com/creack/pty"
	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
	"golang.org/x/sys/unix"
)

func (term *Term) Resize(size *types.XY) {
	term.size = size
	term.resizePty()
}

func (term *Term) resizePty() {
	if term.Pty == nil { //|| term.process == nil {
		debug.Log("cannot resize pty: term.Pty is nil")
		return
	}

	err := pty.Setsize(term.Pty.File(), &pty.Winsize{
		Cols: uint16(term.size.X),
		Rows: uint16(term.size.Y),
	})
	if err != nil {
		debug.Log(err)
	}
	err = term.process.Signal(unix.SIGWINCH)
	if err != nil {
		debug.Log(err)
	}
}

func (term *Term) resize80() {
	term.reset(&types.XY{X: 80, Y: 24})
}

func (term *Term) resize132() {
	term.reset(&types.XY{X: 132, Y: 24})
}
