package rendersdl

import (
	"log"

	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) CacheImage(bmp []byte) (types.Element, error) {
	rwops, err := sdl.RWFromMem(bmp)
	if err != nil {
		return nil, err
	}

	defer rwops.Free()

	img := image{renderer: sr}

	img.surface, err = sdl.LoadBMPRW(rwops, true)
	if err != nil {
		return nil, err
	}

	img.size = &types.XY{
		X: sr.termSize.X / 4,
		Y: sr.termSize.Y / 4,
	}

	return &img, nil
}

type image struct {
	surface  *sdl.Surface
	renderer *sdlRender
	size     *types.XY
}

func (img *image) Size() *types.XY {
	return img.size
}

func (img *image) Draw(topLeft *types.XY) {
	srcRect := sdl.Rect{
		H: img.surface.H,
		W: img.surface.W,
	}

	dstRect := sdl.Rect{
		X: img.renderer.border + (topLeft.X * img.renderer.glyphSize.X),
		Y: img.renderer.border + (topLeft.Y * img.renderer.glyphSize.Y),
		W: img.renderer.border + img.size.X,
		H: img.renderer.border + img.size.Y,
	}

	err := img.surface.BlitScaled(&srcRect, img.renderer.surface, &dstRect)
	if err != nil {
		log.Printf("ERROR: cannot blit image: %s", err.Error())
	}
}

func (img *image) Close() {
	img.surface.Free()
}
