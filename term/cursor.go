package virtualterm

import (
	"time"

	"github.com/lmorg/mxtty/types"
)

func (term *Term) slowBlink() {
	for {
		select {
		case <-time.After(500 * time.Millisecond):
			if !term._hasFocus {
				continue
			}
			term._slowBlinkState = !term._slowBlinkState
			term.renderer.TriggerRedraw()

		case <-term._eventClose:
			term.Pty.Close()
			return

		case <-term._hasKeyPress:
			term._slowBlinkState = true
		}
	}
}

func (term *Term) ShowCursor(v bool) {
	term._hideCursor = !v
}

func (term *Term) _renderCursor() {
	if term._hideCursor || term._scrollOffset != 0 {
		return
	}

	if term._slowBlinkState {
		term.renderer.DrawHighlightRect(term.curPos(), &types.XY{1, 1})
		term.renderer.DrawHighlightRect(term.curPos(), &types.XY{1, 1})
	}
}
