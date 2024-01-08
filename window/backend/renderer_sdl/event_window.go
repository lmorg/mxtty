package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/veandco/go-sdl2/sdl"
)

func eventWindow(r types.Renderer, evt *sdl.WindowEvent, term *virtualterm.Term) {
	switch evt.Event {
	case sdl.WINDOWEVENT_RESIZED:
		resizeTerm(r, term)
	}
}

func resizeTerm(r types.Renderer, term *virtualterm.Term) {
	rect := r.Resize()
	term.Resize(rect)
}
