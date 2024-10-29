package types

import "os"

type Term interface {
	Start(Pty)
	GetSize() *XY
	Resize(*XY)
	Render()
	CopyLines(int32, int32) []byte
	CopySquare(*XY, *XY) []byte
	Bg() *Colour
	Reply([]byte)
	MouseClick(uint8, *XY)
	MouseWheel(int)
	ShowCursor(bool)
}

type Pty interface {
	File() *os.File
	Read() rune
	Write([]byte) error
}
