package backend

import (
	"github.com/lmorg/mxtty/virtualterm/types"
	rendersdl "github.com/lmorg/mxtty/window/backend/renderer_sdl"
)

func Start() *types.Renderer {
	return rendersdl.Start()
}
