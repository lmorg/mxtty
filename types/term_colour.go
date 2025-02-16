package types

type Colour struct {
	Red   byte
	Green byte
	Blue  byte
	Alpha byte
}

// RGB32 combines RGB values into a 32-bit integer
func (c *Colour) RGB24() uint32 {
	return (uint32(c.Red) << 16) | (uint32(c.Green) << 8) | uint32(c.Blue)
}

// RGBA32 combines RGBA values into a 32-bit integer
func (c *Colour) RGBA32() uint32 {
	return (uint32(c.Red) << 24) | (uint32(c.Green) << 16) | (uint32(c.Blue) << 8) | uint32(c.Alpha)
}

// RGBA compatibility with color.Color
func (c *Colour) RGBA() (uint32, uint32, uint32, uint32) {
	if c.Alpha == 0 {
		c.Alpha = 255
	}
	return uint32(c.Red) << 8, uint32(c.Green) << 8, uint32(c.Blue) << 8, uint32(c.Alpha) << 8
}
