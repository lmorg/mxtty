package types

type Renderer interface {
	Size() *Rect
	Resize() *Rect
	PrintRuneColour(r rune, posX, posY int32, fg *Colour, bg *Colour, style SgrFlag) error
	SetWindowTitle(string)
	GetWindowTitle() string
	Close()
}
