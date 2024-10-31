package rendersdl

import (
	"github.com/lmorg/mxtty/config"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) eventWindow(evt *sdl.WindowEvent) {
	switch evt.Event {
	case sdl.WINDOWEVENT_RESIZED:
		//resizeTerm(r, term)

	case sdl.WINDOWEVENT_FOCUS_GAINED:
		sr.term.HasFocus(true)
		sr.window.SetWindowOpacity(float32(config.Config.Window.Opacity) / 100)
		sr.hkToggle = true

	case sdl.WINDOWEVENT_FOCUS_LOST:
		sr.term.HasFocus(false)
		sr.window.SetWindowOpacity(0.7)
		sr.hkToggle = false
	}
}

/*func resizeTerm(r types.Renderer, term types.Term) {
	//r.
	term.Resize(rect)
}
*/
