package virtualterm

import (
	"fmt"

	"github.com/lmorg/mxtty/types"
)

func (term *Term) MouseClick(pos *types.XY, button uint8, clicks uint8, pressed bool, callback types.EventIgnoredCallback) {
	//log.Printf("DEBUG: MouseClick(%d: %s)", button, json.LazyLogging(pos))

	term._mouseButtonDown = pressed

	if !pressed {
		return
	}

	cells := term.visibleScreen()

	if cells[pos.Y][pos.X].Element == nil {
		callback()
		return
	}

	cells[pos.Y][pos.X].Element.MouseClick(cells[pos.Y][pos.X].ElementXY(), button, clicks, callback)
	term._mouseButtonDown = false
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

	if term.IsAltBuf() {
		return
	}

	if len(term._scrollBuf) == 0 {
		return
	}

	term._scrollOffset += int(movement.Y * 2)
	term.updateScrollback()
}

func (term *Term) updateScrollback() {
	if term._scrollOffset > len(term._scrollBuf) {
		term._scrollOffset = len(term._scrollBuf)
	}

	if term._scrollOffset < 0 {
		term._scrollOffset = 0
	}

	if term._scrollOffset == 0 {
		term.ShowCursor(true)
		if term._scrollMsg != nil {
			term._scrollMsg.Close()
			term._scrollMsg = nil
		}

	} else {
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
		if term._mouseIn != nil {
			term._mouseIn.MouseOut()
		}
		callback()
		return
	}

	if cells[pos.Y][pos.X].Element != term._mouseIn {
		if term._mouseIn != nil {
			term._mouseIn.MouseOut()
		}
		term._mouseIn = cells[pos.Y][pos.X].Element
	}

	cells[pos.Y][pos.X].Element.MouseMotion(cells[pos.Y][pos.X].ElementXY(), movement, callback)
}
