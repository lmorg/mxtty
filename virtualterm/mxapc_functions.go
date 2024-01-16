package virtualterm

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/lmorg/mxtty/types"
)

func (term *Term) mxapcBegin(element types.ElementID, parameters *types.ApcSlice) {
	term._activeElement = term.renderer.NewElement(element, nil, nil)
	term._activeElement.Begin(parameters)
}

func (term *Term) mxapcEnd() {
	if term._activeElement == nil {
		return
	}
	term._activeElement.End()
	term._activeElement = nil
}
