package virtualterm

import "github.com/lmorg/mxtty/virtualterm/types"

type sgrFlag uint32

// Flags
const (
	SGR_RESET sgrFlag = 0

	SGR_BOLD sgrFlag = 1 << iota
	SGR_ITALIC
	SGR_UNDERLINE
	SGR_BLINK
	SGR_INVERT
)

type sgr struct {
	bitwise sgrFlag
	fg      *types.Colour
	bg      *types.Colour
}

func (s *sgr) Is(flag sgrFlag) bool {
	return s.bitwise&flag != 0
}

func (s *sgr) Reset() {
	s.bitwise = 0
	s.fg = SGR_DEFAULT.fg
	s.bg = SGR_DEFAULT.bg
}

func (s *sgr) Set(flag sgrFlag) {
	s.bitwise |= flag
}

func (s *sgr) Unset(flag sgrFlag) {
	s.bitwise &^= flag
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
