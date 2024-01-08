package virtualterm

import (
	"time"
)

func (term *Term) slowBlink() {
	for {
		time.Sleep(500 * time.Millisecond)
		term._slowBlinkState = !term._slowBlinkState
	}
}
