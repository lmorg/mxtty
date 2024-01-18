package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

func eventMouseButton(evt *sdl.MouseButtonEvent, term types.Term, sr *sdlRender) {
	if evt.State == sdl.PRESSED {
		return
	}

	pos := types.XY{
		X: (evt.X - sr.border) / sr.glyphSize.X,
		Y: (evt.Y - sr.border) / sr.glyphSize.Y,
	}
	term.MouseClick(evt.Button, &pos)
}
