package rendersdl

import (
	"log"

	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) CacheImage(bmp []byte) (types.Image, error) {
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

	return &img, nil
}

type image struct {
	surface  *sdl.Surface
	renderer *sdlRender
}

func (img *image) Close() {
	img.surface.Free()
}

func (img *image) Draw(topLeft, bottomRight *types.XY) {
	srcRect := sdl.Rect{
		H: img.surface.H,
		W: img.surface.W,
	}

	dstRect := sdl.Rect{
		X: img.renderer.border + (topLeft.X * img.renderer.glyphSize.X),
		Y: img.renderer.border + (topLeft.Y * img.renderer.glyphSize.Y),
		W: img.renderer.border + (bottomRight.X * img.renderer.glyphSize.X),
		H: img.renderer.border + (bottomRight.Y * img.renderer.glyphSize.Y),
	}

	err := img.surface.BlitScaled(&srcRect, img.renderer.surface, &dstRect)
	if err != nil {
		log.Printf("ERROR: cannot blit image: %s", err.Error())
	}
}
