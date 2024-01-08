package backend

import (
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/virtualterm"
	rendersdl "github.com/lmorg/mxtty/window/backend/renderer_sdl"
)

func Initialise(fontName string, fontSize int) types.Renderer {
	return rendersdl.Initialise(fontName, fontSize)
	//return rendererimgui.Initialise(fontName, fontSize)
}

func Start(r types.Renderer, term *virtualterm.Term) {
	rendersdl.Start(r, term)
	//rendererimgui.Start(term)
}
