package types

import "os"

type Term interface {
	Start(Pty, string)
	GetSize() *XY
	Resize(*XY)
	Render()
	Bg() *Colour
	Reply([]byte) error
	MouseClick(uint8, *XY)
	MouseWheel(int)
	ShowCursor(bool)
}

type Pty interface {
	File() *os.File
	Read() rune
	Write([]byte) error
}
