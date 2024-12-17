package virtualterm

import (
	"fmt"

	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
)

// MouseClick: pos X should be -1 when out of bounds
func (term *Term) MouseClick(pos *types.XY, button uint8, clicks uint8, pressed bool, callback types.EventIgnoredCallback) {
	//log.Printf("DEBUG: MouseClick(%d: %s)", button, json.LazyLogging(pos))

	term._mouseButtonDown = pressed

	screen := term.visibleScreen()

	if pos != nil && pos.X < 0 {
		absPos := int32(len(term._scrollBuf)) - int32(term._scrollOffset) + pos.Y

		if len(screen[pos.Y].Hidden) > 0 {
			err := term.UnhideRows(absPos)
			if err != nil {
				term.renderer.DisplayNotification(types.NOTIFY_WARN, err.Error())
			}
			return
		}

		var block []int32
		for _, block = range term._cacheBlock {
			if block[0] <= pos.Y && block[0]+block[1] >= pos.Y {
				goto isOutputBlock
			}
		}

	isOutputBlock:

		blockPos, _, err := term.outputBlocksFindStartEnd(absPos)
		debug.Log(blockPos)
		if err != nil {
			term.renderer.DisplayNotification(types.NOTIFY_WARN, err.Error())
			return
		}

		if err = term.HideRows(blockPos[0], blockPos[1]+1); err != nil {
			term.renderer.DisplayNotification(types.NOTIFY_WARN, err.Error())
		}
		return
	}

	if !pressed {
		return
	}

	if screen[pos.Y].Cells[pos.X].Element == nil {
		callback()
		return
	}

	screen[pos.Y].Cells[pos.X].Element.MouseClick(screen[pos.Y].Cells[pos.X].ElementXY(), button, clicks, callback)
	term._mouseButtonDown = false
}

func (term *Term) MouseWheel(pos *types.XY, movement *types.XY) {
	//log.Printf("DEBUG: MouseScroll(%d: %s)", Y, json.LazyLogging(pos))

	screen := term.visibleScreen()

	if screen[pos.Y].Cells[pos.X].Element == nil {
		term._mouseWheelCallback(movement)
		return
	}

	screen[pos.Y].Cells[pos.X].Element.MouseWheel(
		screen[pos.Y].Cells[pos.X].ElementXY(),
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
	screen := term.visibleScreen()

	if screen[pos.Y].Cells[pos.X].Element == nil {
		if term._mouseIn != nil {
			term._mouseIn.MouseOut()
		}
		callback()
		return
	}

	if screen[pos.Y].Cells[pos.X].Element != term._mouseIn {
		if term._mouseIn != nil {
			term._mouseIn.MouseOut()
		}
		term._mouseIn = screen[pos.Y].Cells[pos.X].Element
	}

	screen[pos.Y].Cells[pos.X].Element.MouseMotion(screen[pos.Y].Cells[pos.X].ElementXY(), movement, callback)
}
