package types

type Cell struct {
	Char    rune
	Sgr     *Sgr
	Element Element
	Phrase  *[]rune
}

func (c *Cell) Clear() {
	c.Char = 0
	c.Sgr = &Sgr{}
	c.Element = nil
}

func (c *Cell) Rune() rune {
	switch {
	case c.Element != nil:
		return c.Element.Rune(c.ElementXY())

	case c.Char == 0:
		return ' '

	default:
		return c.Char
	}
}

const cellElementXyMask = (^int32(0)) << 16

func (c *Cell) ElementXY() *XY {
	return &XY{
		X: c.Char >> 16,
		Y: c.Char &^ cellElementXyMask,
	}
}
