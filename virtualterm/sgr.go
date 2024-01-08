package virtualterm

import "github.com/lmorg/mxtty/types"

type sgr struct {
	bitwise types.SgrFlag
	fg      *types.Colour
	bg      *types.Colour
}

func (s *sgr) Reset() {
	s.bitwise = 0
	s.fg = SGR_DEFAULT.fg
	s.bg = SGR_DEFAULT.bg
}

func (c *cell) clear() {
	c.char = 0
	c.sgr = &sgr{}
}

func (s *sgr) Copy() *sgr {
	return &sgr{
		fg:      s.fg.Copy(),
		bg:      s.bg.Copy(),
		bitwise: s.bitwise,
	}
}
