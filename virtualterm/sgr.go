package virtualterm

import "github.com/lmorg/mxtty/virtualterm/types"

type sgr struct {
	bitwise sgrFlag
	fg      types.Colour
	bg      types.Colour
}

func (s *sgr) Is(flag sgrFlag) bool {
	return s.bitwise&flag != 0
}

func (s *sgr) sgrReset() {
	s.bitwise = 0
	s.fg = SGR_DEFAULT.fg
	s.bg = SGR_DEFAULT.bg
}

func (s *sgr) Set(flag sgrFlag) {
	s.bitwise |= flag
}

func (c *cell) clear() {
	c.char = 0
	c.sgr = &sgr{}
}

func (s *sgr) Copy() *sgr {
	return &sgr{
		fg:      s.fg,
		bg:      s.bg,
		bitwise: s.bitwise,
	}
}
