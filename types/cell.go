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
