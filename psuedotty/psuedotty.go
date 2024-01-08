package psuedotty

import (
	"fmt"
	"log"
	"os"

	"github.com/creack/pty"
	"github.com/lmorg/mxtty/virtualterm/types"
)

type PTY struct {
	Primary   *os.File
	Secondary *os.File
	stream    chan rune
}

func NewPTY(size *types.Rect) (*PTY, error) {
	secondary, primary, err := pty.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open pty: %s", err.Error())
	}

	err = pty.Setsize(primary, &pty.Winsize{
		Cols: uint16(size.X),
		Rows: uint16(size.Y),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to set pty size: %s", err.Error())
	}

	p := &PTY{
		Primary:   primary,
		Secondary: secondary,
		stream:    make(chan rune),
	}

	go p.write()

	return p, err
}

func (p *PTY) write() {
	b := make([]byte, 10*1024)

	for {
		n, err := p.Secondary.Read(b)
		if err != nil {
			log.Printf("error reading from PTY (%d bytes dropped): %s", n, err.Error())
			continue
		}

		s := string(b[:n])
		for _, r := range s {
			p.stream <- r
		}
	}
}

func (p *PTY) ReadRune() rune {
	return <-p.stream
}
