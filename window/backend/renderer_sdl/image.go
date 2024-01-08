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

func (img *image) Draw(rect *types.Rect) {
	srcRect := sdl.Rect{
		H: img.surface.H,
		W: img.surface.W,
	}

	dstRect := sdl.Rect{
		X: img.renderer.border,
		Y: img.renderer.border,
		W: img.renderer.surface.W,
		H: img.renderer.surface.H,
	}

	err := img.surface.BlitScaled(&srcRect, img.renderer.surface, &dstRect)
	if err != nil {
		log.Printf("ERROR: cannot blit image: %s", err.Error())
	}
}
