package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	elementImage "github.com/lmorg/mxtty/window/backend/renderer_sdl/element_image"
)

func (sr *sdlRender) NewElement(id types.ElementID) types.Element {
	switch id {
	case types.ELEMENT_ID_IMAGE:
		return elementImage.New(sr, sr.loadImage)

	//case types.ELEMENT_ID_TABLE:
	//	return elementTable.New(sr)

	default:
		return nil
	}
}
