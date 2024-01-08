package backend

import (
	"github.com/lmorg/mxtty/types"
	rendersdl "github.com/lmorg/mxtty/window/backend/renderer_sdl"
)

func Initialise(fontName string, fontSize int) types.Renderer {
	return rendersdl.Initialise(fontName, fontSize)
}

func Start(r types.Renderer, term types.Term) {
	r.Start(term)
}
