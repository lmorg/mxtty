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

/*
	ROWS
*/
/*
type Row struct {
	Cells []*Cell
	Meta  RowMetaFlag
}

type RowMetaFlag uint16

// Flags
const (
	ROW_META_NONE      RowMetaFlag = 0
	ROW_META_COLLAPSED RowMetaFlag = 1 << iota
)

func (f RowMetaFlag) Is(flag RowMetaFlag) bool {
	return f&flag != 0
}

func (f *RowMetaFlag) Set(flag RowMetaFlag) {
	*f |= flag
}

func (f *RowMetaFlag) Unset(flag RowMetaFlag) {
	*f &^= flag
}
*/
