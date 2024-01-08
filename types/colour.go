package types

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
