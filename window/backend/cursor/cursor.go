package cursor

import "github.com/veandco/go-sdl2/sdl"

var cursor int

const (
	arrow int = iota + 1
	ibeam
	hand
)

func Arrow() {
	if cursor == arrow {
		return
	}

	sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW))
	cursor = arrow
}

func Ibeam() {
	if cursor == ibeam {
		return
	}

	sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_IBEAM))
	cursor = ibeam
}

func Hand() {
	if cursor == hand {
		return
	}

	sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_HAND))
	cursor = hand
}
