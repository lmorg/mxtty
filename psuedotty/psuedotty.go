package psuedotty

import (
	"fmt"
	"os"

	"github.com/creack/pty"
)

type PTY struct {
	Primary   *os.File
	Secondary *os.File
}

func NewPTY(width, height int32) (*PTY, error) {
	secondary, primary, err := pty.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open pty: %s", err.Error())
	}

	size := pty.Winsize{
		Cols: uint16(width),
		Rows: uint16(height),
	}

	err = pty.Setsize(primary, &size)
	if err != nil {
		return nil, fmt.Errorf("unable to set pty size: %s", err.Error())
	}

	p := new(PTY)
	p.Primary = primary
	p.Secondary = secondary
	return p, nil
}
