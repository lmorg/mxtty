package virtualterm

import "github.com/lmorg/mxtty/virtualterm/types"

var (
	sgrColour4Black   = &types.Colour{}
	sgrColour4Red     = &types.Colour{Red: 200}
	sgrColour4Green   = &types.Colour{Green: 200}
	sgrColour4Yellow  = &types.Colour{Red: 200, Green: 200}
	sgrColour4Blue    = &types.Colour{Blue: 200}
	sgrColour4Magenta = &types.Colour{Red: 200, Blue: 200}
	sgrColour4Cyan    = &types.Colour{Green: 200, Blue: 200}
	sgrColour4White   = &types.Colour{Red: 200, Green: 200, Blue: 200}

	sgrColour4BlackBright   = &types.Colour{Red: 100, Green: 100, Blue: 100}
	sgrColour4RedBright     = &types.Colour{Red: 255}
	sgrColour4GreenBright   = &types.Colour{Green: 255}
	sgrColour4YellowBright  = &types.Colour{Red: 255, Green: 255}
	sgrColour4BlueBright    = &types.Colour{Blue: 255}
	sgrColour4MagentaBright = &types.Colour{Red: 255, Blue: 255}
	sgrColour4CyanBright    = &types.Colour{Green: 255, Blue: 255}
	sgrColour4WhiteBright   = &types.Colour{Red: 255, Green: 255, Blue: 255}
)

var SGR_DEFAULT = &sgr{
	fg: sgrColour4White,
	bg: sgrColour4Black,
}
