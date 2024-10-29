package rendersdl

import (
	"sync"
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
	renderer  *sdl.Renderer
	glyphSize *types.XY
	termSize  *types.XY
	term      types.Term
	limiter   sync.Mutex

	// preferences
	font       *ttf.Font
	_fontStyle types.SgrFlag
	border     int32

	// title
	title       string
	updateTitle int32

	// audio
	bell *mix.Music

	// events
	_quit   chan bool
	_redraw chan bool

	// notifications
	notifications  notifyT
	notifyIcon     map[int]types.Image
	notifyIconSize *types.XY

	// widgets
	termWidget  *termWidgetT
	highlighter *highlighterT
	inputBox    *inputBoxT
	menu        *menuT

	// render function stack (AddRenderFnToStack)
	fnStack []func()

	// state
	blinkState  bool
	keyModifier uint16
}

func (sr *sdlRender) TriggerQuit()  { go sr._triggerQuit() }
func (sr *sdlRender) _triggerQuit() { sr._quit <- true }

func (sr *sdlRender) TriggerRedraw() { go sr._triggerRedraw() }
func (sr *sdlRender) _triggerRedraw() {
	if sr.limiter.TryLock() {
		sr._redraw <- true
	}
}

func (sr *sdlRender) TermSize() *types.XY {
	return sr.termSize
}

func (sr *sdlRender) resize() *types.XY {
	/*var err error
	//sr.surface.Free()
	//sr.surface, err = sr.window.GetSurface()
	sr.renderer, err = sr.window.GetRenderer()
	if err != nil {
		panic(err) // TODO: this shouldn't panic!
	}*/

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

func (sr *sdlRender) FocusWindow() {
	sr.window.Raise()
}

func (sr *sdlRender) GetWindowMeta() any {
	return sr.window
}

func (sr *sdlRender) ResizeWindow(size *types.XY) {
	w := (size.X * sr.glyphSize.X) + (sr.border * 2)
	h := (size.Y * sr.glyphSize.Y) + (sr.border * 2)
	sr.window.SetSize(w, h)
}
