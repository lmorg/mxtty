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
	glyphSize *types.Rect
	termSize  *types.Rect

	title       string
	updateTitle int32
}

func (sr *sdlRender) Size() *types.Rect {
	return sr.termSize
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
