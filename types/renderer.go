package types

type Renderer interface {
	Start(Term)
	Size() *XY
	Resize() *XY
	PrintRuneColour(r rune, posX, posY int32, fg *Colour, bg *Colour, style SgrFlag) error
	CacheImage(bmp []byte) (Element, error)
	SetWindowTitle(string)
	GetWindowTitle() string
	Bell()
	Close()
}
