package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) loadImage(bmp []byte, size *types.XY) (types.Image, error) {
	rwops, err := sdl.RWFromMem(bmp)
	if err != nil {
		return nil, err
	}

	img := image{sr: sr, rwops: rwops}

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

	img.texture, err = img.sr.renderer.CreateTextureFromSurface(img.surface)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

type image struct {
	surface *sdl.Surface
	sr      *sdlRender
	size    *types.XY
	rwops   *sdl.RWops
	texture *sdl.Texture
}

func (img *image) Size() *types.XY {
	return img.size
}

func (img *image) Draw(size *types.XY, rect *types.Rect) {
	srcRect := sdl.Rect{
		W: img.surface.W,
		H: img.surface.H,
	}

	offset := (size.Y - (rect.End.Y - rect.Start.Y) - 1) * img.sr.glyphSize.Y

	dstRect := sdl.Rect{
		X: img.sr.border + (rect.Start.X * img.sr.glyphSize.X),
		Y: img.sr.border + (rect.Start.Y * img.sr.glyphSize.Y) - offset,
		W: img.size.X,
		H: img.size.Y,
	}

	err := img.sr.renderer.Copy(img.texture, &srcRect, &dstRect)
	if err != nil {
		panic(err) //TODO: don't panic!
	}
}

func (img *image) Close() {
	img.texture.Destroy()
	img.surface.Free()
	img.rwops.Free()
}

func (sr *sdlRender) AddImageToStack(fn func()) {
	sr.imageStack = append(sr.imageStack, fn)
}
