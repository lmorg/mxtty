package types

type Renderer interface {
	Size() *Rect
	Close()
	Update() error
	PrintRuneColour(r rune, posX, posY int32, fg *Colour, bg *Colour, style SgrFlag) error
	SetWindowTitle(string)
}
