package psuedotty

import (
	"fmt"
	"os"

	"github.com/creack/pty"
	"github.com/lmorg/mxtty/virtualterm"
)

type PTY struct {
	Primary   *os.File
	Secondary *os.File
	virtTerm  *virtualterm.Term
}

func NewPTY(term *virtualterm.Term) (*PTY, error) {
	secondary, primary, err := pty.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open pty: %s", err.Error())
	}

	width, height, err := term.GetSize()
	if err != nil {
		return nil, err
	}

	size := pty.Winsize{
		Cols: uint16(width),
		Rows: uint16(height),
	}

	err = pty.Setsize(primary, &size)
	if err != nil {
		return nil, fmt.Errorf("unable to set pty size: %s", err.Error())
	}

	p := &PTY{
		Primary:   primary,
		Secondary: secondary,
		virtTerm:  term,
	}

	term.Pty = p.Secondary

	return p, err
}
