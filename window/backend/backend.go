package backend

import (
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/virtualterm"
	rendersdl "github.com/lmorg/mxtty/window/backend/renderer_sdl"
)

func Initialise() types.Renderer {
	return rendersdl.Initialise()
}

func Start(term *virtualterm.Term) {
	rendersdl.Start(term)
}
