package types

type Renderer interface {
	Start(Term)
	Size() *XY
	Resize() *XY
	PrintRuneColour(r rune, posX, posY int32, fg *Colour, bg *Colour, style SgrFlag) error
	CacheImage(bmp []byte) (Image, error)
	SetWindowTitle(string)
	GetWindowTitle() string
	Close()
}

type Image interface {
	Draw(*XY, *XY)
	Close()
}
