package virtualterm

type sgrFlag uint32

// Flags
const (
	sgrReset sgrFlag = 0

	sgrBold sgrFlag = 1 << iota
	sgrItalic
	sgrUnderscore
	sgrBlink
	sgrInvert

	// colour bit pallets
	sgrFgColour4
	sgrFgColour8
	sgrFgColour24

	sgrBgColour4
	sgrBgColour8
	sgrBgColour24
)

const (
	sgrColour4Black = 0
	sgrColour4Red   = iota
	sgrColour4Green
	sgrColour4Yellow
	sgrColour4Blue
	sgrColour4Magenta
	sgrColour4Cyan
	sgrColour4White

	sgrColour4BlackBright
	sgrColour4RedBright
	sgrColour4GreenBright
	sgrColour4YellowBright
	sgrColour4BlueBright
	sgrColour4MagentaBright
	sgrColour4CyanBright
	sgrColour4WhiteBright
)
