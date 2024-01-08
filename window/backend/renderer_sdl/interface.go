package rendersdl

import (
	"sync/atomic"

	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/window/backend/typeface"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type sdlRender struct {
	window    *sdl.Window
	surface   *sdl.Surface
	font      *ttf.Font
	glyphSize *types.XY
	termSize  *types.XY
	border    int32

	title       string
	updateTitle int32
}

func (sr *sdlRender) Size() *types.XY {
	return sr.termSize
}

func (sr *sdlRender) Resize() *types.XY {
	var err error
	sr.surface.Free()
	sr.surface, err = sr.window.GetSurface()
	if err != nil {
		panic(err) // TODO: this shouldn't panic!
	}

	return sr.getTermSize()
}

func (sr *sdlRender) Close() {
	typeface.Close()
	sr.window.Destroy()
	sdl.Quit()
}

func (sr *sdlRender) SetWindowTitle(title string) {
	sr.title = title
	atomic.CompareAndSwapInt32(&sr.updateTitle, 0, 1)
}

func (sr *sdlRender) GetWindowTitle() string {
	return sr.window.GetTitle()
}
