package virtualterm

import (
	"time"
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

		case <-term._hasKeypress:
			term._slowBlinkState = true
		}
	}
}
