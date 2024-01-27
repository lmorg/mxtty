package rendersdl

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) eventWindow(evt *sdl.WindowEvent) {
	switch evt.Event {
	case sdl.WINDOWEVENT_RESIZED:
		//resizeTerm(r, term)
	}
}

/*func resizeTerm(r types.Renderer, term types.Term) {
	//r.
	term.Resize(rect)
}
*/
