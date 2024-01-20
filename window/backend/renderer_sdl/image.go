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

	img.size = &types.XY{
		X: sr.glyphSize.X * size.X,
		Y: sr.glyphSize.Y * size.Y,
	}

	if size.X == 0 {
		img.size.X = int32((float64(img.surface.W) / float64(img.surface.H)) * float64(img.size.Y))
		size.X = int32((float64(img.size.X) / float64(sr.glyphSize.X)) + 1)
	}

	winW, _ := sr.window.GetSize()
	if img.size.X > winW {
		img.size.X = winW
		img.size.Y = int32((float64(img.surface.H) / float64(img.surface.W)) * float64(img.size.X))
		size.X = int32((float64(img.size.X) / float64(sr.glyphSize.X)))
		size.Y = int32((float64(img.size.Y) / float64(sr.glyphSize.Y)) + 1)
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
		img.sr.DisplayNotification(types.NOTIFY_ERROR, "Cannot render image: "+err.Error())
	}
}

func (img *image) Asset() any {
	return img.surface
}

func (img *image) Close() {
	img.texture.Destroy()
	img.surface.Free()
	img.rwops.Free()
}

func (sr *sdlRender) AddImageToStack(fn func()) {
	sr.imageStack = append(sr.imageStack, fn)
}
