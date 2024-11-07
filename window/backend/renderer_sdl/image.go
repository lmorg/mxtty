package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

type image struct {
	surface   *sdl.Surface
	sr        *sdlRender
	sizeCells *types.XY
	rwops     *sdl.RWops
	texture   *sdl.Texture
}

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

	img.sizeCells = &types.XY{
		X: sr.glyphSize.X * size.X,
		Y: sr.glyphSize.Y * size.Y,
	}

	if size.X == 0 {
		img.sizeCells.X = int32((float64(img.surface.W) / float64(img.surface.H)) * float64(img.sizeCells.Y))
		size.X = int32((float64(img.sizeCells.X) / float64(sr.glyphSize.X)) + 1)
	}

	wPx, _ := sr.window.GetSize()
	if img.sizeCells.X > wPx {
		img.sizeCells.X = wPx
		img.sizeCells.Y = int32((float64(img.surface.H) / float64(img.surface.W)) * float64(img.sizeCells.X))
		size.X = int32((float64(img.sizeCells.X) / float64(sr.glyphSize.X)))
		size.Y = int32((float64(img.sizeCells.Y) / float64(sr.glyphSize.Y)) + 1)
	}

	img.texture, err = img.sr.renderer.CreateTextureFromSurface(img.surface)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func (img *image) Size() *types.XY {
	return img.sizeCells
}

func (img *image) Draw(size *types.XY, pos *types.XY) {
	srcRect := &sdl.Rect{
		W: img.surface.W,
		H: img.surface.H,
	}

	dstRect := &sdl.Rect{
		X: img.sr.border + (pos.X * img.sr.glyphSize.X),
		Y: img.sr.border + (pos.Y * img.sr.glyphSize.Y),
		W: size.X * img.sr.glyphSize.X,
		H: size.Y * img.sr.glyphSize.Y,
	}

	// for rendering
	err := img.sr.renderer.Copy(img.texture, srcRect, dstRect)
	if err != nil {
		img.sr.DisplayNotification(types.NOTIFY_ERROR, "Cannot render image: "+err.Error())
	}

	if img.sr.highlighter != nil {
		// for clipboard
		err := img.surface.BlitScaled(srcRect, img.sr.surface, dstRect)
		if err != nil {
			img.sr.DisplayNotification(types.NOTIFY_ERROR, "Cannot render image: "+err.Error())
		}
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

func (sr *sdlRender) AddRenderFnToStack(fn func()) {
	sr.fnStack = append(sr.fnStack, fn)
}
