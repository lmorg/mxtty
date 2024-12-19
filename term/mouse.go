package virtualterm

import (
	"fmt"
	"strings"

	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
)

// MouseClick: pos X should be -1 when out of bounds
func (term *Term) MouseClick(pos *types.XY, button uint8, clicks uint8, pressed bool, callback types.EventIgnoredCallback) {
	screen := term.visibleScreen()

	// this is used to determine whether to override ligatures with default font rendering
	term._mouseButtonDown = pressed

	if pos == nil {
		// this just exists to reset ligatures
		return
	}

	if pressed {
		callback()
		return
	}

	if pos.X < 0 {
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
			if block[0] <= pos.Y && block[0]+block[1] > pos.Y {
				goto isOutputBlock
			}
		}

		// not an output block
		return

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

	if screen[pos.Y].Cells[pos.X].Element == nil {
		if h := term._mousePositionCodeFoldable(screen, pos); h != -1 {
			err := term.FoldAtIndent(pos)
			if err != nil {
				term.renderer.DisplayNotification(types.NOTIFY_WARN, err.Error())
			}
		}

		callback()
		return
	}

	screen[pos.Y].Cells[pos.X].Element.MouseClick(screen[pos.Y].Cells[pos.X].ElementXY(), button, clicks, pressed, callback)
}

func (term *Term) MouseWheel(pos *types.XY, movement *types.XY) {
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

func (term *Term) MousePosition(pos *types.XY) {
	if pos.X < 0 {
		return
	}

	screen := term.visibleScreen()

	if screen[pos.Y].Cells[pos.X].Element == nil {
		if height := term._mousePositionCodeFoldable(screen, pos); height >= 0 {
			term.renderer.DrawHighlightRect(
				&types.XY{pos.X, pos.Y},
				&types.XY{term.size.X - pos.X, height - pos.Y},
			)
			term.renderer.StatusBarText("[Click] Fold branch")
		}
	}
}

func (term *Term) _mousePositionCodeFoldable(screen types.Screen, pos *types.XY) int32 {
	if pos.Y >= term.curPos().Y {
		return -1
	}

	if screen[pos.Y].Cells[pos.X].Char == ' ' {
		return -1
	}

	if pos.X > 0 && screen[pos.Y].Cells[pos.X-1].Char != ' ' {
		pos.X--
	}

	for x := pos.X - 1; x >= 0; x-- {
		if screen[pos.Y].Cells[x].Char != ' ' {
			return -1
		}
	}

	absScreen := append(term._scrollBuf, term._normBuf...)
	absPos := term.convertRelPosToAbsPos(pos)

	height, err := outputBlockFoldIndent(term, absScreen, absPos, false)
	if err != nil {
		return -1
	}

	height = height - absPos.Y + pos.Y

	if height-pos.Y == 2 && strings.TrimSpace(string(*screen[height-1].Phrase)) == "" {
		return -1
	}

	return height
}
