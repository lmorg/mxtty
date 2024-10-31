package virtualterm

import (
	"time"
)

func (term *Term) slowBlink() {
	for {
		time.Sleep(500 * time.Millisecond)
		if !term._hasFocus {
			continue
		}

		term._slowBlinkState = !term._slowBlinkState
		term.renderer.TriggerRedraw()
	}
}
