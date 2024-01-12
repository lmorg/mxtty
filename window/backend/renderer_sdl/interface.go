package rendersdl

import (
	"sync/atomic"

	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/window/backend/typeface"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type sdlRender struct {
	window    *sdl.Window
	surface   *sdl.Surface
	glyphSize *types.XY
	termSize  *types.XY

	// preferences
	font       *ttf.Font
	_fontStyle types.SgrFlag
	border     int32
	dropShadow bool

	// title
	title       string
	updateTitle int32

	// audio
	bell *mix.Music

	// events
	_quit   chan bool
	_redraw chan bool
}

func (sr *sdlRender) triggerQuit()   { sr._quit <- true }
func (sr *sdlRender) TriggerRedraw() { sr._redraw <- true }

func (sr *sdlRender) TermSize() *types.XY {
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

	if sr.bell != nil {
		sr.bell.Free()
		mix.CloseAudio()
		mix.Quit()
	}

	sdl.Quit()
}

func (sr *sdlRender) SetWindowTitle(title string) {
	sr.title = title
	atomic.CompareAndSwapInt32(&sr.updateTitle, 0, 1)
}

func (sr *sdlRender) GetWindowTitle() string {
	return sr.window.GetTitle()
}
