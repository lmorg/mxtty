package types

import "os"

type MouseClickCallback func()

type Term interface {
	Start(Pty)
	GetSize() *XY
	Resize(*XY)
	Render()
	CopyRange(*XY, *XY) []byte
	CopyLines(int32, int32) []byte
	CopySquare(*XY, *XY) []byte
	Bg() *Colour
	Reply([]byte)
	MouseClick(uint8, *XY, MouseClickCallback)
	MouseWheel(int)
	ShowCursor(bool)
	HasFocus(bool)
}

type Pty interface {
	File() *os.File
	Read() rune
	Write([]byte) error
}
