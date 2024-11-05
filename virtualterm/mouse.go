package virtualterm

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/lmorg/murex/utils/json"
	"github.com/lmorg/mxtty/types"
)

func (term *Term) MouseClick(pos *types.XY, button uint8, callback types.EventIgnoredCallback) {
	log.Printf("DEBUG: MouseClick(%d: %s)", button, json.LazyLogging(pos))

	cells := term.visibleScreen()

	if cells[pos.Y][pos.X].Element == nil {
		callback()
		return
	}

	cells[pos.Y][pos.X].Element.MouseClick(cells[pos.Y][pos.X].ElementXY(), button, callback)
}

func (term *Term) MouseWheel(pos *types.XY, Y int) {
	cells := term.visibleScreen()

	if cells[pos.Y][pos.X].Element == nil {
		term._mouseWheelCallback(Y)
		return
	}

	cells[pos.Y][pos.X].Element.MouseWheel(
		cells[pos.Y][pos.X].ElementXY(),
		Y,
		func() { term._mouseWheelCallback(Y) },
	)
}

func (term *Term) _mouseWheelCallback(Y int) {
	if unsafe.Pointer(term.cells) != unsafe.Pointer(&term._normBuf) {
		return
	}

	if len(term._scrollBuf) == 0 {
		return
	}

	term._scrollOffset += Y * 2

	switch {
	case term._scrollOffset > len(term._scrollBuf):
		term._scrollOffset = len(term._scrollBuf)

	case term._scrollOffset < 0:
		term._scrollOffset = 0
		fallthrough

	case term._scrollOffset == 0:
		term.ShowCursor(true)
		if term._scrollMsg != nil {
			term._scrollMsg.Close()
			term._scrollMsg = nil
		}

	default:
		term.ShowCursor(false)
		msg := fmt.Sprintf("Viewing scrollback history. %d lines from end", term._scrollOffset)
		if term._scrollMsg == nil {
			term._scrollMsg = term.renderer.DisplaySticky(types.NOTIFY_SCROLL, msg)
		} else {
			term._scrollMsg.SetMessage(msg)
		}
	}
}
