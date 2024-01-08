package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/types"
)

func lookupSgr(sgr *sgr, n int32, stack []int32) {
	for _, i := range stack {
		switch i {
		case -1, 0: // reset / normal
			sgr.Reset()

		case 1: // bold
			sgr.bitwise.Set(types.SGR_BOLD)

		case 3: // italic
			sgr.bitwise.Set(types.SGR_ITALIC)

		case 4: // underline
			sgr.bitwise.Set(types.SGR_UNDERLINE)

		case 5, 6: // blink
			sgr.bitwise.Set(types.SGR_SLOW_BLINK)

		case 7: // invert
			sgr.bitwise.Set(types.SGR_INVERT)

		case 22: // no bold
			sgr.bitwise.Unset(types.SGR_BOLD)

		case 23: // no italic
			sgr.bitwise.Unset(types.SGR_ITALIC)

		case 24: // no underline
			sgr.bitwise.Unset(types.SGR_UNDERLINE)

		case 25: // no blink
			sgr.bitwise.Unset(types.SGR_SLOW_BLINK)

		case 27: // no invert
			sgr.bitwise.Unset(types.SGR_INVERT)

		//
		// 3bit foreground colour:
		//

		case 30: // fg black
			sgr.fg = SGR_COLOUR_BLACK

		case 31: // fg red
			sgr.fg = SGR_COLOUR_RED

		case 32: // fg green
			sgr.fg = SGR_COLOUR_GREEN

		case 33: // fg yellow
			sgr.fg = SGR_COLOUR_YELLOW

		case 34: // fg blue
			sgr.fg = SGR_COLOUR_BLUE

		case 35: // fg magenta
			sgr.fg = SGR_COLOUR_MAGENTA

		case 36: // fg cyan
			sgr.fg = SGR_COLOUR_CYAN

		case 37: // fg white
			sgr.fg = SGR_COLOUR_WHITE

		case 38:
			colour := _sgrEnhancedColour(n, stack)
			if colour != nil {
				sgr.fg = colour
			}
			return

		case 39: // fg default
			sgr.fg = SGR_DEFAULT.fg

		//
		// 3bit background colour:
		//

		case 40: // bg black
			sgr.bg = SGR_COLOUR_BLACK

		case 41: // bg rede
			sgr.bg = SGR_COLOUR_RED

		case 42: // bg green
			sgr.bg = SGR_COLOUR_GREEN

		case 43: // bg yellow
			sgr.bg = SGR_COLOUR_YELLOW

		case 44: // bg blue
			sgr.bg = SGR_COLOUR_BLUE

		case 45: // bg magenta
			sgr.bg = SGR_COLOUR_MAGENTA

		case 46: // bg cyan
			sgr.bg = SGR_COLOUR_CYAN

		case 47: // bg white
			sgr.bg = SGR_COLOUR_WHITE

		case 48:
			colour := _sgrEnhancedColour(n, stack)
			if colour != nil {
				sgr.bg = colour
			}
			return

		case 49: // bg default
			sgr.bg = SGR_DEFAULT.bg

		//
		// 4bit foreground colour:
		//

		case 90: // fg black
			sgr.fg = SGR_COLOUR_BLACK_BRIGHT

		case 91: // fg red
			sgr.fg = SGR_COLOUR_RED_BRIGHT

		case 92: // fg green
			sgr.fg = SGR_COLOUR_GREEN_BRIGHT

		case 93: // fg yellow
			sgr.fg = SGR_COLOUR_YELLOW_BRIGHT

		case 94: // fg blue
			sgr.fg = SGR_COLOUR_BLUE_BRIGHT

		case 95: // fg magenta
			sgr.fg = SGR_COLOUR_MAGENTA_BRIGHT

		case 96: // fg cyan
			sgr.fg = SGR_COLOUR_CYAN_BRIGHT

		case 97: // fg white
			sgr.fg = SGR_COLOUR_WHITE_BRIGHT

		//
		// 4bit background colour:
		//

		case 100: // bg black
			sgr.bg = SGR_COLOUR_BLACK_BRIGHT

		case 101: // bg red
			sgr.bg = SGR_COLOUR_RED_BRIGHT

		case 102: // bg green
			sgr.bg = SGR_COLOUR_GREEN_BRIGHT

		case 103: // bg yellow
			sgr.bg = SGR_COLOUR_YELLOW_BRIGHT

		case 104: // bg blue
			sgr.bg = SGR_COLOUR_BLUE_BRIGHT

		case 105: // bg magenta
			sgr.bg = SGR_COLOUR_MAGENTA_BRIGHT

		case 106: // bg cyan
			sgr.bg = SGR_COLOUR_CYAN_BRIGHT

		case 107: // bg white
			sgr.bg = SGR_COLOUR_WHITE_BRIGHT

		default:
			log.Printf("Unknown SGR code: %d", n)
		}
	}
}

func _sgrEnhancedColour(n int32, stack []int32) *types.Colour {
	if len(stack) < 2 {
		log.Printf("SGR error: too few parameters in %d: %v", n, stack)
		return nil
	}
	switch stack[1] {
	case 5:
		colour, ok := SGR_COLOUR_256[stack[2]]
		if !ok {
			log.Printf("SGR error: 256 value does not exist in %d: %v", n, stack)
			return nil
		}
		return colour

	case 2:
		if len(stack) != 5 {
			log.Printf("SGR error: too few parameters in %d (24bit): %v", n, stack)
			return nil
		}
		return &types.Colour{
			Red:   byte(stack[2]),
			Green: byte(stack[3]),
			Blue:  byte(stack[4]),
		}

	default:
		log.Printf("SGR error: unexpected value in %d: %v", n, stack)
		return nil
	}

}
