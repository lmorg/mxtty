package virtualterm

import "github.com/lmorg/mxtty/virtualterm/types"

var (
	SGR_COLOUR_BLACK   = types.Colour{}
	SGR_COLOUR_RED     = types.Colour{Red: 200}
	SGR_COLOUR_GREEN   = types.Colour{Green: 200}
	SGR_COLOUR_YELLOW  = types.Colour{Red: 200, Green: 200}
	SGR_COLOUR_BLUE    = types.Colour{Blue: 200}
	SGR_COLOUR_MAGENTA = types.Colour{Red: 200, Blue: 200}
	SGR_COLOUR_CYAN    = types.Colour{Green: 200, Blue: 200}
	SGR_COLOUR_WHITE   = types.Colour{Red: 200, Green: 200, Blue: 200}

	SGR_COLOUR_BLACK_BRIGHT   = types.Colour{Red: 100, Green: 100, Blue: 100}
	SGR_COLOUR_RED_BRIGHT     = types.Colour{Red: 255}
	SGR_COLOUR_GREEN_BRIGHT   = types.Colour{Green: 255}
	SGR_COLOUR_YELLOW_BRIGHT  = types.Colour{Red: 255, Green: 255}
	SGR_COLOUR_BLUE_BRIGHT    = types.Colour{Blue: 255}
	SGR_COLOUR_MAGENTA_BRIGHT = types.Colour{Red: 255, Blue: 255}
	SGR_COLOUR_CYAN_BRIGHT    = types.Colour{Green: 255, Blue: 255}
	SGR_COLOUR_WHITE_BRIGHT   = types.Colour{Red: 255, Green: 255, Blue: 255}
)

var SGR_DEFAULT = &sgr{
	fg: SGR_COLOUR_WHITE,
	bg: SGR_COLOUR_BLACK,
}
