package types

type Renderer struct {
	Size           *Rect
	Close          func()
	Update         func() error
	PrintRuneColor func(r rune, posX, posY int32, fg *Colour, bg *Colour, style SgrFlag) error
}

type Colour struct {
	Red   byte
	Green byte
	Blue  byte
}

func (c *Colour) Copy() *Colour {
	return &Colour{
		Red:   c.Red,
		Green: c.Green,
		Blue:  c.Blue,
	}
}
