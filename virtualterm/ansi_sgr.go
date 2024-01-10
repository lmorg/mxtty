package virtualterm

import (
	"log"

	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/virtualterm/cell"
)

func lookupSgr(sgr *cell.Sgr, n int32, stack []int32) {
	for _, i := range stack {
		switch i {
		case -1, 0: // reset / normal
			sgr.Reset()

		case 1: // bold
			sgr.Bitwise.Set(types.SGR_BOLD)

		case 3: // italic
			sgr.Bitwise.Set(types.SGR_ITALIC)

		case 4: // underline
			sgr.Bitwise.Set(types.SGR_UNDERLINE)

		case 5, 6: // blink
			sgr.Bitwise.Set(types.SGR_SLOW_BLINK)

		case 7: // invert
			sgr.Bitwise.Set(types.SGR_INVERT)

		case 22: // no bold
			sgr.Bitwise.Unset(types.SGR_BOLD)

		case 23: // no italic
			sgr.Bitwise.Unset(types.SGR_ITALIC)

		case 24: // no underline
			sgr.Bitwise.Unset(types.SGR_UNDERLINE)

		case 25: // no blink
			sgr.Bitwise.Unset(types.SGR_SLOW_BLINK)

		case 27: // no invert
			sgr.Bitwise.Unset(types.SGR_INVERT)

		//
		// 3bit foreground colour:
		//

		case 30: // fg black
			sgr.Fg = cell.SGR_COLOUR_BLACK

		case 31: // fg red
			sgr.Fg = cell.SGR_COLOUR_RED

		case 32: // fg green
			sgr.Fg = cell.SGR_COLOUR_GREEN

		case 33: // fg yellow
			sgr.Fg = cell.SGR_COLOUR_YELLOW

		case 34: // fg blue
			sgr.Fg = cell.SGR_COLOUR_BLUE

		case 35: // fg magenta
			sgr.Fg = cell.SGR_COLOUR_MAGENTA

		case 36: // fg cyan
			sgr.Fg = cell.SGR_COLOUR_CYAN

		case 37: // fg white
			sgr.Fg = cell.SGR_COLOUR_WHITE

		case 38:
			colour := _sgrEnhancedColour(n, stack)
			if colour != nil {
				sgr.Fg = colour
			}
			return

		case 39: // fg default
			sgr.Fg = cell.SGR_DEFAULT.Fg

		//
		// 3bit background colour:
		//

		case 40: // bg black
			sgr.Bg = cell.SGR_COLOUR_BLACK

		case 41: // bg rede
			sgr.Bg = cell.SGR_COLOUR_RED

		case 42: // bg green
			sgr.Bg = cell.SGR_COLOUR_GREEN

		case 43: // bg yellow
			sgr.Bg = cell.SGR_COLOUR_YELLOW

		case 44: // bg blue
			sgr.Bg = cell.SGR_COLOUR_BLUE

		case 45: // bg magenta
			sgr.Bg = cell.SGR_COLOUR_MAGENTA

		case 46: // bg cyan
			sgr.Bg = cell.SGR_COLOUR_CYAN

		case 47: // bg white
			sgr.Bg = cell.SGR_COLOUR_WHITE

		case 48:
			colour := _sgrEnhancedColour(n, stack)
			if colour != nil {
				sgr.Bg = colour
			}
			return

		case 49: // bg default
			sgr.Bg = cell.SGR_DEFAULT.Bg

		//
		// 4bit foreground colour:
		//

		case 90: // fg black
			sgr.Fg = cell.SGR_COLOUR_BLACK_BRIGHT

		case 91: // fg red
			sgr.Fg = cell.SGR_COLOUR_RED_BRIGHT

		case 92: // fg green
			sgr.Fg = cell.SGR_COLOUR_GREEN_BRIGHT

		case 93: // fg yellow
			sgr.Fg = cell.SGR_COLOUR_YELLOW_BRIGHT

		case 94: // fg blue
			sgr.Fg = cell.SGR_COLOUR_BLUE_BRIGHT

		case 95: // fg magenta
			sgr.Fg = cell.SGR_COLOUR_MAGENTA_BRIGHT

		case 96: // fg cyan
			sgr.Fg = cell.SGR_COLOUR_CYAN_BRIGHT

		case 97: // fg white
			sgr.Fg = cell.SGR_COLOUR_WHITE_BRIGHT

		//
		// 4bit background colour:
		//

		case 100: // bg black
			sgr.Bg = cell.SGR_COLOUR_BLACK_BRIGHT

		case 101: // bg red
			sgr.Bg = cell.SGR_COLOUR_RED_BRIGHT

		case 102: // bg green
			sgr.Bg = cell.SGR_COLOUR_GREEN_BRIGHT

		case 103: // bg yellow
			sgr.Bg = cell.SGR_COLOUR_YELLOW_BRIGHT

		case 104: // bg blue
			sgr.Bg = cell.SGR_COLOUR_BLUE_BRIGHT

		case 105: // bg magenta
			sgr.Bg = cell.SGR_COLOUR_MAGENTA_BRIGHT

		case 106: // bg cyan
			sgr.Bg = cell.SGR_COLOUR_CYAN_BRIGHT

		case 107: // bg white
			sgr.Bg = cell.SGR_COLOUR_WHITE_BRIGHT

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
		colour, ok := cell.SGR_COLOUR_256[stack[2]]
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
