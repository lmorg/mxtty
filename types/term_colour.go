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

// RGB32 combines RGB values into a 32-bit integer
func (c *Colour) RGB24() uint32 {
	return (uint32(c.Red) << 16) | (uint32(c.Green) << 8) | uint32(c.Blue)
}

// RGBA32 combines RGBA values into a 32-bit integer
func (c *Colour) RGBA32(alpha byte) uint32 {
	return (uint32(c.Red) << 24) | (uint32(c.Green) << 16) | (uint32(c.Blue) << 8) | uint32(alpha)
}

const _ALPHA_UINT32 = 255 << 8

// RGBA compatibility with color.Color
func (c *Colour) RGBA() (uint32, uint32, uint32, uint32) {
	return uint32(c.Red) << 8, uint32(c.Green) << 8, uint32(c.Blue) << 8, _ALPHA_UINT32
}
