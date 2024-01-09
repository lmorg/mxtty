package types

import "os"

type Term interface {
	Start(Pty, string)
	GetSize() *XY
	Resize(*XY)
	Render()
	Bg() *Colour
	Reply([]byte) error
}

type Pty interface {
	File() *os.File
	Read() rune
	Write([]byte) error
}
