package types

type Sgr struct {
	Bitwise SgrFlag
	Fg      *Colour
	Bg      *Colour
}

func (s *Sgr) Reset() {
	s.Bitwise = 0
	s.Fg = SGR_DEFAULT.Fg
	s.Bg = SGR_DEFAULT.Bg
}

func (s *Sgr) Copy() *Sgr {
	return &Sgr{
		Fg:      s.Fg,
		Bg:      s.Bg,
		Bitwise: s.Bitwise,
	}
}
