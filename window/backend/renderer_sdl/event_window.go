package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

func eventWindow(r types.Renderer, evt *sdl.WindowEvent, term types.Term) {
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
