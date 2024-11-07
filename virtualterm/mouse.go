package virtualterm

import (
	"fmt"
	"unsafe"

	"github.com/lmorg/mxtty/types"
)

func (term *Term) MouseClick(pos *types.XY, button uint8, clicks uint8, callback types.EventIgnoredCallback) {
	//log.Printf("DEBUG: MouseClick(%d: %s)", button, json.LazyLogging(pos))

	cells := term.visibleScreen()

	if cells[pos.Y][pos.X].Element == nil {
		callback()
		return
	}

	cells[pos.Y][pos.X].Element.MouseClick(cells[pos.Y][pos.X].ElementXY(), button, clicks, callback)
}

func (term *Term) MouseWheel(pos *types.XY, movement *types.XY) {
	//log.Printf("DEBUG: MouseScroll(%d: %s)", Y, json.LazyLogging(pos))

	cells := term.visibleScreen()

	if cells[pos.Y][pos.X].Element == nil {
		term._mouseWheelCallback(movement)
		return
	}

	cells[pos.Y][pos.X].Element.MouseWheel(
		cells[pos.Y][pos.X].ElementXY(),
		movement,
		func() { term._mouseWheelCallback(movement) },
	)
}

func (term *Term) _mouseWheelCallback(movement *types.XY) {
	if movement.Y == 0 {
		return
	}

	if unsafe.Pointer(term.cells) != unsafe.Pointer(&term._normBuf) {
		return
	}

	if len(term._scrollBuf) == 0 {
		return
	}

	term._scrollOffset += int(movement.Y * 2)
	term.updateScrollback()
}

func (term *Term) updateScrollback() {
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

func (term *Term) MouseMotion(pos *types.XY, movement *types.XY, callback types.EventIgnoredCallback) {
	cells := term.visibleScreen()

	if cells[pos.Y][pos.X].Element == nil {
		callback()
		return
	}

	cells[pos.Y][pos.X].Element.MouseMotion(cells[pos.Y][pos.X].ElementXY(), movement, callback)
}
