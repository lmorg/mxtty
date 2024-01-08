package psuedotty

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/creack/pty"
	"github.com/lmorg/mxtty/virtualterm/types"
)

type PTY struct {
	Primary   *os.File
	Secondary *os.File
	buf       *bytes.Buffer
	b         []byte
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
	}

	go p.write()

	return p, err
}

func (p *PTY) write() {
	p.buf = bytes.NewBuffer(p.b)
	b := make([]byte, 10*1024)

	for {
		read, err := p.Secondary.Read(b)
		if err != nil {
			log.Printf("error reading from PTY (%d bytes dropped): %s", read, err.Error())
			continue
		}

		written, err := p.buf.Write(b[:read])
		if err != nil {
			log.Printf("error writing to buffer (%d bytes dropped): %s", written, err.Error())
			continue
		}

		if read != written {
			log.Printf("read and write buffer mismatch: read %d, written %d", read, written)
			continue
		}
	}
}

func (p *PTY) ReadRune() rune {
	for {
		r, _, err := p.buf.ReadRune()
		if err != nil {
			log.Printf("error reading from buffer: %s", err.Error())
		}
		return r
	}
}
