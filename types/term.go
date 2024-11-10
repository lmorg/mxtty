package types

import "os"

type EventIgnoredCallback func()

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
	MouseClick(*XY, uint8, uint8, bool, EventIgnoredCallback)
	MouseWheel(*XY, *XY)
	MouseMotion(*XY, *XY, EventIgnoredCallback)
	ShowCursor(bool)
	HasFocus(bool)
}

type Pty interface {
	File() *os.File
	Read() rune
	Write([]byte) error
}
