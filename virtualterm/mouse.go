package virtualterm

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/lmorg/murex/utils/json"
	"github.com/lmorg/mxtty/types"
)

func (term *Term) MouseClick(button uint8, pos *types.XY) {
	log.Printf("DEBUG: MouseClick(%d: %s)", button, json.LazyLogging(pos))

	cells := term.visibleScreen()

	if cells[pos.Y][pos.X].Element != nil {
		//relPos := types.XY{X: pos.X, Y: pos.Y - term.findElementStart(pos)}
		cells[pos.Y][pos.X].Element.MouseClick(button, getElementXY(cells[pos.Y][pos.X].Char))
	}
}

func (term *Term) MouseWheel(Y int) {
	if unsafe.Pointer(term.cells) != unsafe.Pointer(&term._normBuf) {
		return
	}

	if len(term._scrollBuf) == 0 {
		return
	}

	term._scrollOffset += Y * 2

	switch {
	case term._scrollOffset < 0:
		term._scrollOffset = 0

	case term._scrollOffset > len(term._scrollBuf):
		term._scrollOffset = len(term._scrollBuf)

	case term._scrollOffset == 0:
		term.ShowCursor(true)
		if term._scrollMsg != nil {
			term._scrollMsg.Close()
			term._scrollMsg = nil
		}

	default:
		term.ShowCursor(false)
		if term._scrollMsg == nil {
			term._scrollMsg = term.renderer.DisplaySticky(types.NOTIFY_SCROLL, "")
		}
		term._scrollMsg.SetMessage(fmt.Sprintf("Viewing scrollback history. %d lines from end", term._scrollOffset))
	}
}
