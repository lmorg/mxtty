package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) eventMouseButton(evt *sdl.MouseButtonEvent) {
	if sr.inputBoxActive {
		return
	}

	if evt.State == sdl.PRESSED {
		return
	}

	pos := types.XY{
		X: (evt.X - sr.border) / sr.glyphSize.X,
		Y: (evt.Y - sr.border) / sr.glyphSize.Y,
	}
	sr.term.MouseClick(evt.Button, &pos)
}
