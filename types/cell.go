package types

type Cell struct {
	Char    rune
	Sgr     *Sgr
	Element Element
}

const (
	CELL_NULL          = 0
	CELL_ELEMENT_BEGIN = 1
	CELL_ELEMENT_FILL  = 2
	CELL_ELEMENT_END   = 2
)

func (c *Cell) Clear() {
	c.Char = 0
	//c.sgr = &sgr{}
}
