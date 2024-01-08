package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/window/backend/typeface"
	"github.com/veandco/go-sdl2/sdl"
)

type sdlRender struct{}

func (sr *sdlRender) Size() *types.Rect {
	return termSize
}

func (sr *sdlRender) Update() error {
	return window.UpdateSurface()
}

func (sr *sdlRender) Close() {
	typeface.Close()
	window.Destroy()
	sdl.Quit()
}

func (sr *sdlRender) SetWindowTitle(title string) {
	/*
		unsupported in SDL due to:
		NSWindow geometry should only be modified on the main thread!
	*/

	//window.SetTitle(fmt.Sprintf("%s - %s", title, app.Name))
}
