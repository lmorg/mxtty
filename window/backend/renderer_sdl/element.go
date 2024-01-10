package rendersdl

import "github.com/lmorg/mxtty/types"

func (sr *sdlRender) NewElement(id types.ElementID, size *types.XY, data []byte) types.Element {
	switch id {
	case types.ELEMENT_ID_IMAGE:
		return nil

	case types.ELEMENT_ID_TABLE:
		return &ElementTable{
			size: size,
		}

	default:
		return nil
	}
}

type ElementTable struct {
	size *types.XY
}

func (el *ElementTable) Close() {
	// do nothing
}

func (el *ElementTable) fgh() {

}

func (el *ElementTable) Start() {
	
}