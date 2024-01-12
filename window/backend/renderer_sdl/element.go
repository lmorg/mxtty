package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	elementTable "github.com/lmorg/mxtty/window/backend/renderer_sdl/element_table"
)

func (sr *sdlRender) NewElement(id types.ElementID, size *types.XY, data []byte) types.Element {
	switch id {
	case types.ELEMENT_ID_IMAGE:
		return nil

	case types.ELEMENT_ID_TABLE:
		return elementTable.New(sr)

	default:
		return nil
	}
}
