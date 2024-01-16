package rendersdl

import (
	"log"

	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) loadImage(bmp []byte, size *types.XY) (types.Image, error) {
	rwops, err := sdl.RWFromMem(bmp)
	if err != nil {
		return nil, err
	}

	img := image{renderer: sr, rwops: rwops}

	img.surface, err = sdl.LoadBMPRW(rwops, true)
	if err != nil {
		return nil, err
	}

	if size.X == 0 {
		f := img.surface.W / img.surface.H
		size.X = size.Y * f * 2
	}

	img.size = &types.XY{
		X: sr.glyphSize.X * size.X,
		Y: sr.glyphSize.Y * size.Y,
	}

	return &img, nil
}

type image struct {
	surface  *sdl.Surface
	renderer *sdlRender
	size     *types.XY
	rwops    *sdl.RWops
}

func (img *image) Size() *types.XY {
	return img.size
}

func (img *image) Draw(size *types.XY, rect *types.Rect) {
	srcRect := sdl.Rect{
		W: img.surface.W,
		H: img.surface.H,
	}

	offset := (size.Y - (rect.End.Y - rect.Start.Y) - 1) * img.renderer.glyphSize.Y

	dstRect := sdl.Rect{
		X: img.renderer.border + (rect.Start.X * img.renderer.glyphSize.X),
		Y: img.renderer.border + (rect.Start.Y * img.renderer.glyphSize.Y) - offset,
		W: img.size.X,
		H: img.size.Y,
	}

	err := img.surface.BlitScaled(&srcRect, img.renderer.surface, &dstRect)
	if err != nil {
		log.Printf("ERROR: cannot blit image: %s", err.Error())
	}
}

func (img *image) Close() {
	img.surface.Free()
	img.rwops.Free()
}
