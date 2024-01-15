package virtualterm

import "github.com/lmorg/mxtty/types"

func (term *Term) mxapcTableBegin(parameters *types.ApcSlice) {
	term._activeElement = term.renderer.NewElement(types.ELEMENT_ID_TABLE, nil, nil)
	term._activeElement.Begin(parameters)
}

func (term *Term) mxapcTableEnd(parameters *types.ApcSlice) {
	if term._activeElement == nil {
		return
	}
	term._activeElement.End()
	term._activeElement = nil
}
