package virtualterm

var (
	sgrColour4Black   = &rgb{}
	sgrColour4Red     = &rgb{Red: 200}
	sgrColour4Green   = &rgb{Green: 200}
	sgrColour4Yellow  = &rgb{Red: 200, Green: 200}
	sgrColour4Blue    = &rgb{Blue: 200}
	sgrColour4Magenta = &rgb{Red: 200, Blue: 200}
	sgrColour4Cyan    = &rgb{Green: 200, Blue: 200}
	sgrColour4White   = &rgb{Red: 200, Green: 200, Blue: 200}

	sgrColour4BlackBright   = &rgb{Red: 100, Green: 100, Blue: 100}
	sgrColour4RedBright     = &rgb{Red: 255}
	sgrColour4GreenBright   = &rgb{Green: 255}
	sgrColour4YellowBright  = &rgb{Red: 255, Green: 255}
	sgrColour4BlueBright    = &rgb{Blue: 255}
	sgrColour4MagentaBright = &rgb{Red: 255, Blue: 255}
	sgrColour4CyanBright    = &rgb{Green: 255, Blue: 255}
	sgrColour4WhiteBright   = &rgb{Red: 255, Green: 255, Blue: 255}
)

var SGR_DEFAULT = &sgr{
	fg: sgrColour4White,
	bg: sgrColour4Black,
}
