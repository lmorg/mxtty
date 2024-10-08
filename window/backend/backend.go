package backend

import (
	"github.com/lmorg/mxtty/types"
	rendersdl "github.com/lmorg/mxtty/window/backend/renderer_sdl"
)

func Initialise() types.Renderer {
	return rendersdl.Initialise()
}

func Start(r types.Renderer, term types.Term) {
	r.Start(term)
}
