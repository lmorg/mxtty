package types

type Cell struct {
	Char    rune
	Sgr     *Sgr
	Element Element
}

func (c *Cell) Clear() {
	c.Char = 0
	c.Sgr = &Sgr{}
	c.Element = nil
}

func (c *Cell) Rune() rune {
	switch {
	case c.Element != nil:
		return ' '

	case c.Char == 0:
		return ' '

	default:
		return c.Char
	}
}
