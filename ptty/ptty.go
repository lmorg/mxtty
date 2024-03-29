package ptty

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/creack/pty"
	"github.com/lmorg/mxtty/types"
)

type PTY struct {
	primary   *os.File
	secondary *os.File
	stream    chan rune
}

func NewPTY(size *types.XY) (types.Pty, error) {
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
		primary:   primary,
		secondary: secondary,
		stream:    make(chan rune),
	}

	go p.read(secondary)

	return p, err
}

func (p *PTY) File() *os.File {
	return p.primary
}

func (p *PTY) Write(b []byte) error {
	_, err := p.secondary.Write(b)
	return err
}

func (p *PTY) read(f *os.File) {
	var (
		b    = make([]byte, 10*1024)
		utf8 []byte
		l    int
	)

	for {
		n, err := f.Read(b)
		if err != nil && err.Error() != io.EOF.Error() {
			log.Printf("ERROR: problem reading from PTY (%d bytes dropped): %s", n, err.Error())
			continue
		}

		for i := 0; i < n; i++ {
			if l == 0 {
				l = runeLength(b[i])
				utf8 = make([]byte, l)
			}

			utf8[len(utf8)-l] = b[i]

			if l == 1 {
				r := bytes.Runes(utf8)
				p.stream <- r[0]
			}
			l--
		}
	}
}

func runeLength(b byte) int {
	switch {
	case b&128 == 0:
		return 1
	case b&32 == 0:
		return 2
	case b&16 == 0:
		return 3
	case b&8 == 0:
		return 4
	default:
		return 0
	}
}

func (p *PTY) Read() rune {
	return <-p.stream
}
